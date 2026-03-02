package handler

import (
	"fmt"
	"net/http"

	"github.com/cgy/webhook-tester/internal/auth"
	"github.com/cgy/webhook-tester/internal/config"
	"github.com/cgy/webhook-tester/internal/database/sqlc"
	"github.com/cgy/webhook-tester/internal/sse"
	"github.com/go-chi/chi/v5"
)

type SSEHandler struct {
	queries *sqlc.Queries
	config  *config.Config
	hub     *sse.Hub
}

func NewSSE(queries *sqlc.Queries, cfg *config.Config, hub *sse.Hub) *SSEHandler {
	return &SSEHandler{
		queries: queries,
		config:  cfg,
		hub:     hub,
	}
}

func (h *SSEHandler) Stream(w http.ResponseWriter, r *http.Request) {
	webhookID := chi.URLParam(r, "uuid")

	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	webhook, err := h.queries.GetWebhookByID(r.Context(), webhookID)
	if err != nil {
		http.Error(w, "webhook not found", http.StatusNotFound)
		return
	}

	if webhook.UserID != claims.UserID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch, unsubscribe := h.hub.Subscribe(webhookID)
	defer unsubscribe()

	for {
		select {
		case data := <-ch:
			fmt.Fprintf(w, "event: new-request\ndata: %s\n\n", data)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
