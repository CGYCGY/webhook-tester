package server

import (
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
