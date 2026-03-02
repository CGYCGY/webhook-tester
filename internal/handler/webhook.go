package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cgy/webhook-tester/internal/auth"
	"github.com/cgy/webhook-tester/internal/config"
	"github.com/cgy/webhook-tester/internal/database/sqlc"
	"github.com/cgy/webhook-tester/internal/templates"
	"github.com/cgy/webhook-tester/internal/templates/components"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type WebhookHandler struct {
	queries *sqlc.Queries
	config  *config.Config
}

func NewWebhook(queries *sqlc.Queries, cfg *config.Config) *WebhookHandler {
	return &WebhookHandler{
		queries: queries,
		config:  cfg,
	}
}

func baseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if fwd := r.Header.Get("X-Forwarded-Proto"); fwd != "" {
		scheme = fwd
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

func (h *WebhookHandler) ListWebhooks(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	webhooks, err := h.queries.ListWebhooksByUserID(r.Context(), claims.UserID)
	if err != nil {
		http.Error(w, "Failed to load webhooks", http.StatusInternalServerError)
		return
	}

	base := baseURL(r)
	views := make([]templates.WebhookView, 0, len(webhooks))
	for _, wh := range webhooks {
		count, err := h.queries.GetWebhookRequestCount(r.Context(), wh.ID)
		if err != nil {
			count = 0
		}
		views = append(views, templates.WebhookView{
			ID:           wh.ID,
			Name:         wh.Name,
			Description:  wh.Description,
			URL:          fmt.Sprintf("%s/hook/%s", base, wh.ID),
			RequestCount: count,
			CreatedAt:    wh.CreatedAt.Format("Jan 2, 2006"),
		})
	}

	templates.DashboardPage(views).Render(r.Context(), w)
}

func (h *WebhookHandler) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	description := strings.TrimSpace(r.FormValue("description"))

	if errMsg := validateWebhookFields(name, description); errMsg != "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		components.FormError(errMsg).Render(r.Context(), w)
		return
	}

	id := uuid.New().String()
	now := time.Now().UTC()

	if err := h.queries.CreateWebhook(r.Context(), sqlc.CreateWebhookParams{
		ID:          id,
		UserID:      claims.UserID,
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}); err != nil {
		http.Error(w, "Failed to create webhook", http.StatusInternalServerError)
		return
	}

	base := baseURL(r)
	view := templates.WebhookView{
		ID:           id,
		Name:         name,
		Description:  description,
		URL:          fmt.Sprintf("%s/hook/%s", base, id),
		RequestCount: 0,
		CreatedAt:    now.Format("Jan 2, 2006"),
	}

	if r.Header.Get("HX-Request") == "true" {
		components.WebhookCard(view).Render(r.Context(), w)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func (h *WebhookHandler) EditWebhook(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	id := chi.URLParam(r, "uuid")

	wh, err := h.queries.GetWebhookByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Webhook not found", http.StatusNotFound)
		return
	}

	if wh.UserID != claims.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	description := strings.TrimSpace(r.FormValue("description"))

	if errMsg := validateWebhookFields(name, description); errMsg != "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		components.FormError(errMsg).Render(r.Context(), w)
		return
	}

	now := time.Now().UTC()
	if err := h.queries.UpdateWebhook(r.Context(), sqlc.UpdateWebhookParams{
		Name:        name,
		Description: description,
		UpdatedAt:   now,
		ID:          id,
	}); err != nil {
		http.Error(w, "Failed to update webhook", http.StatusInternalServerError)
		return
	}

	count, _ := h.queries.GetWebhookRequestCount(r.Context(), id)
	base := baseURL(r)
	view := templates.WebhookView{
		ID:           id,
		Name:         name,
		Description:  description,
		URL:          fmt.Sprintf("%s/hook/%s", base, id),
		RequestCount: count,
		CreatedAt:    wh.CreatedAt.Format("Jan 2, 2006"),
	}

	components.WebhookCard(view).Render(r.Context(), w)
}

func (h *WebhookHandler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	id := chi.URLParam(r, "uuid")

	wh, err := h.queries.GetWebhookByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Webhook not found", http.StatusNotFound)
		return
	}

	if wh.UserID != claims.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.queries.DeleteWebhook(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete webhook", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/dashboard")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func (h *WebhookHandler) ViewWebhook(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	id := chi.URLParam(r, "uuid")

	wh, err := h.queries.GetWebhookByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Webhook not found", http.StatusNotFound)
		return
	}

	if wh.UserID != claims.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	reqs, err := h.queries.ListRequestsByWebhookID(r.Context(), sqlc.ListRequestsByWebhookIDParams{
		WebhookID: id,
		Limit:     100,
	})
	if err != nil {
		http.Error(w, "Failed to load requests", http.StatusInternalServerError)
		return
	}

	base := baseURL(r)
	webhookView := templates.WebhookView{
		ID:           wh.ID,
		Name:         wh.Name,
		Description:  wh.Description,
		URL:          fmt.Sprintf("%s/hook/%s", base, wh.ID),
		RequestCount: int64(len(reqs)),
		CreatedAt:    wh.CreatedAt.Format("Jan 2, 2006"),
	}

	requestViews := make([]templates.RequestView, 0, len(reqs))
	for _, req := range reqs {
		path := req.Path
		if path == "" {
			path = "/"
		}
		requestViews = append(requestViews, templates.RequestView{
			ID:            req.ID,
			WebhookID:     req.WebhookID,
			Method:        req.Method,
			Path:          path,
			SourceIP:      req.SourceIP,
			ContentType:   req.ContentType,
			ContentLength: req.ContentLength,
			CreatedAt:     req.CreatedAt.Format("Jan 2, 2006 3:04 PM"),
		})
	}

	templates.WebhookDetailPage(webhookView, requestViews).Render(r.Context(), w)
}

func (h *WebhookHandler) ViewRequest(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	webhookID := chi.URLParam(r, "uuid")
	requestID := chi.URLParam(r, "requestID")

	wh, err := h.queries.GetWebhookByID(r.Context(), webhookID)
	if err != nil {
		http.Error(w, "Webhook not found", http.StatusNotFound)
		return
	}

	if wh.UserID != claims.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	req, err := h.queries.GetRequestByID(r.Context(), requestID)
	if err != nil {
		http.Error(w, "Request not found", http.StatusNotFound)
		return
	}

	if req.WebhookID != webhookID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var headers map[string]string
	if err := json.Unmarshal([]byte(req.Headers), &headers); err != nil {
		headers = make(map[string]string)
	}

	var queryParams map[string]string
	if err := json.Unmarshal([]byte(req.QueryParams), &queryParams); err != nil {
		queryParams = make(map[string]string)
	}

	base := baseURL(r)
	webhookView := templates.WebhookView{
		ID:          wh.ID,
		Name:        wh.Name,
		Description: wh.Description,
		URL:         fmt.Sprintf("%s/hook/%s", base, wh.ID),
		CreatedAt:   wh.CreatedAt.Format("Jan 2, 2006"),
	}

	path := req.Path
	if path == "" {
		path = "/"
	}

	detailView := templates.DetailRequestView{
		ID:            req.ID,
		WebhookID:     req.WebhookID,
		Method:        req.Method,
		Path:          path,
		Headers:       headers,
		QueryParams:   queryParams,
		Body:          req.Body,
		ContentType:   req.ContentType,
		SourceIP:      req.SourceIP,
		ContentLength: req.ContentLength,
		CreatedAt:     req.CreatedAt.Format("Jan 2, 2006 3:04 PM"),
		HeadersJSON:   req.Headers,
		QueryJSON:     req.QueryParams,
	}

	templates.RequestDetailPage(webhookView, detailView).Render(r.Context(), w)
}

func validateWebhookFields(name, description string) string {
	if name == "" {
		return "Name is required."
	}
	if len(name) > 100 {
		return "Name must be 100 characters or fewer."
	}
	if len(description) > 500 {
		return "Description must be 500 characters or fewer."
	}
	return ""
}
