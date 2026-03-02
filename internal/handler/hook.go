package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/cgy/webhook-tester/internal/config"
	"github.com/cgy/webhook-tester/internal/database/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type HookHandler struct {
	queries *sqlc.Queries
	config  *config.Config
}

func NewHook(queries *sqlc.Queries, cfg *config.Config) *HookHandler {
	return &HookHandler{
		queries: queries,
		config:  cfg,
	}
}

func (h *HookHandler) CaptureRequest(w http.ResponseWriter, r *http.Request) {
	webhookID := chi.URLParam(r, "uuid")

	_, err := h.queries.GetWebhookByID(r.Context(), webhookID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Webhook not found",
		})
		return
	}

	body, _ := io.ReadAll(r.Body)

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
		SourceIP:      r.RemoteAddr,
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

	_ = h.queries.DeleteOldRequests(r.Context(), sqlc.DeleteOldRequestsParams{
		WebhookID:      webhookID,
		WebhookIDInner: webhookID,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"message": "Request captured",
	})
}
