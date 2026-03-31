package server

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/cgy/webhook-tester/internal/auth"
	"github.com/cgy/webhook-tester/internal/handler"
	"github.com/cgy/webhook-tester/internal/static"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	// Static files
	staticFS, _ := fs.Sub(static.Files, ".")
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	authHandler := handler.NewAuth(s.queries, s.config)

	// Root redirect
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	})

	// Public routes
	r.Get("/login", authHandler.LoginPage)
	r.Post("/login", authHandler.Login)
	r.Post("/api/login", authHandler.APILogin)
	r.Get("/llms.txt", func(w http.ResponseWriter, r *http.Request) {
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		if fwd := r.Header.Get("X-Forwarded-Proto"); fwd != "" {
			scheme = fwd
		}
		base := fmt.Sprintf("%s://%s", scheme, r.Host)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, `# Webhook Tester API

> Base URL: %s

A tool for creating webhook endpoints that capture and inspect incoming HTTP requests.

## Authentication

All API endpoints (except login and hook capture) require a Bearer token.

POST %s/api/login
Content-Type: application/json
{"email": "<email>", "password": "<password>"}
-> 200 {"token": "<jwt>"}
-> 401 {"error": "invalid email or password"}

Use the token in subsequent requests:
Authorization: Bearer <jwt>

Token expires after 24 hours.

## Endpoints

### List Webhooks
GET %s/api/webhooks
Authorization: Bearer <jwt>
-> 200 [{"id": "<uuid>", "name": "...", "description": "...", "url": "%s/hook/<uuid>", "request_count": 5, "created_at": "2006-01-02T15:04:05Z", "response_config": {"status": 200, "content_type": "application/json", "body": "..."}}]

### Send a Request to a Webhook (public, no auth)
Any HTTP method: GET, POST, PUT, PATCH, DELETE, etc.
%s/hook/<uuid>
%s/hook/<uuid>/any/sub/path
The request is captured and stored. Headers, query params, body are all recorded.
`, base, base, base, base, base, base)
	})

	webhookHandler := handler.NewWebhook(s.queries, s.config)
	hookHandler := handler.NewHook(s.queries, s.config, s.hub, s.whLimiter, s.ipLimiter)
	sseHandler := handler.NewSSE(s.queries, s.config, s.hub)
	settingsHandler := handler.NewSettings(s.queries, s.config)

	// Public hook capture endpoint
	r.HandleFunc("/hook/{uuid}", hookHandler.CaptureRequest)
	r.HandleFunc("/hook/{uuid}/*", hookHandler.CaptureRequest)

	// Authenticated routes
	r.Group(func(r chi.Router) {
		r.Use(auth.Middleware(s.config.JWTSecret))

		r.Post("/logout", authHandler.Logout)
		r.Get("/dashboard", webhookHandler.ListWebhooks)
		r.Get("/settings", settingsHandler.SettingsPage)
		r.Post("/settings/password", settingsHandler.ChangePassword)
		r.Get("/api/webhooks", webhookHandler.APIListWebhooks)
		r.Post("/webhooks", webhookHandler.CreateWebhook)
		r.Get("/webhooks/{uuid}", webhookHandler.ViewWebhook)
		r.Get("/webhooks/{uuid}/sse", sseHandler.Stream)
		r.Get("/webhooks/{uuid}/requests/{requestID}", webhookHandler.ViewRequest)
		r.Put("/webhooks/{uuid}", webhookHandler.EditWebhook)
		r.Patch("/webhooks/{uuid}/response", webhookHandler.UpdateResponseConfig)
		r.Delete("/webhooks/{uuid}", webhookHandler.DeleteWebhook)
	})

	return r
}
