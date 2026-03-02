package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/cgy/webhook-tester/internal/config"
	"github.com/cgy/webhook-tester/internal/database/sqlc"
	"github.com/cgy/webhook-tester/internal/ratelimit"
	"github.com/cgy/webhook-tester/internal/sse"
	"github.com/cgy/webhook-tester/internal/templates/components"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type HookHandler struct {
	queries   *sqlc.Queries
	config    *config.Config
	hub       *sse.Hub
	whLimiter *ratelimit.Limiter
	ipLimiter *ratelimit.Limiter
}

func NewHook(queries *sqlc.Queries, cfg *config.Config, hub *sse.Hub, whLimiter, ipLimiter *ratelimit.Limiter) *HookHandler {
	return &HookHandler{
		queries:   queries,
		config:    cfg,
		hub:       hub,
		whLimiter: whLimiter,
		ipLimiter: ipLimiter,
	}
}

func (h *HookHandler) CaptureRequest(w http.ResponseWriter, r *http.Request) {
	webhookID := chi.URLParam(r, "uuid")

	// Check max body size from Content-Length header
	if r.ContentLength > h.config.MaxBodySize {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Request body too large",
		})
		return
	}

	// Per-IP rate limit
	if !h.ipLimiter.Allow(r.RemoteAddr) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Rate limit exceeded",
		})
		return
	}

	wh, err := h.queries.GetWebhookByID(r.Context(), webhookID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Webhook not found",
		})
		return
	}

	// Per-webhook rate limit (checked after webhook existence validation)
	if !h.whLimiter.Allow(webhookID) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Rate limit exceeded",
		})
		return
	}

	// Read body with MaxBytesReader as safety net
	r.Body = http.MaxBytesReader(w, r.Body, h.config.MaxBodySize)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Request body too large",
		})
		return
	}

	headersMap := make(map[string]string, len(r.Header))
	for k, v := range r.Header {
		if len(v) > 0 {
			headersMap[k] = v[0]
		}
	}
	headersJSON, _ := json.Marshal(headersMap)

	queryMap := make(map[string]string)
	for k, v := range r.URL.Query() {
		if len(v) > 0 {
			queryMap[k] = v[0]
		}
	}
	queryJSON, _ := json.Marshal(queryMap)

	reqID := uuid.New().String()
	now := time.Now().UTC()

	if err := h.queries.CreateRequest(r.Context(), sqlc.CreateRequestParams{
		ID:            reqID,
		WebhookID:     webhookID,
		Method:        r.Method,
		Path:          r.URL.Path,
		QueryParams:   string(queryJSON),
		Headers:       string(headersJSON),
		Body:          string(body),
		ContentType:   r.Header.Get("Content-Type"),
		SourceIp:      r.RemoteAddr,
		ContentLength: r.ContentLength,
		CreatedAt:     now,
	}); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Failed to store request",
		})
		return
	}

	_ = h.queries.DeleteOldRequests(r.Context(), webhookID)

	// Publish to SSE hub for real-time updates
	path := r.URL.Path
	if path == "" {
		path = "/"
	}
	reqView := components.RequestView{
		ID:            reqID,
		WebhookID:     webhookID,
		Method:        r.Method,
		Path:          path,
		SourceIP:      r.RemoteAddr,
		ContentType:   r.Header.Get("Content-Type"),
		ContentLength: r.ContentLength,
		CreatedAt:     now.Format("Jan 2, 2006 3:04 PM"),
	}
	var buf bytes.Buffer
	_ = components.RequestRow(reqView).Render(context.Background(), &buf)
	h.hub.Publish(webhookID, buf.Bytes())

	cfg, configured := parseResponseConfig(wh.ResponseConfig)
	if !configured {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "Request captured",
		})
		return
	}

	ct := cfg.ContentType
	if ct == "" {
		ct = "text/plain"
	}
	w.Header().Set("Content-Type", ct)
	w.WriteHeader(cfg.Status)
	if cfg.Body != "" {
		resolved := resolveTemplate(cfg.Body, r, body)
		w.Write([]byte(resolved))
	}
}
