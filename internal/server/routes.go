package server

import (
	"github.com/cgy/webhook-tester/internal/auth"
	"github.com/cgy/webhook-tester/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	authHandler := handler.NewAuth(s.queries, s.config)

	// Public routes
	r.Get("/login", authHandler.LoginPage)
	r.Post("/login", authHandler.Login)

	// Authenticated routes
	r.Group(func(r chi.Router) {
		r.Use(auth.Middleware(s.config.JWTSecret))

		r.Post("/logout", authHandler.Logout)
		r.Get("/dashboard", handler.NotImplemented)
		r.Get("/settings", handler.NotImplemented)
	})

	return r
}
