package server

import (
	"database/sql"

	"github.com/cgy/webhook-tester/internal/config"
	"github.com/cgy/webhook-tester/internal/database/sqlc"
	"github.com/cgy/webhook-tester/internal/ratelimit"
	"github.com/cgy/webhook-tester/internal/sse"
)

type Server struct {
	db        *sql.DB
	queries   *sqlc.Queries
	config    *config.Config
	hub       *sse.Hub
	whLimiter *ratelimit.Limiter
	ipLimiter *ratelimit.Limiter
}

func New(db *sql.DB, cfg *config.Config) *Server {
	return &Server{
		db:        db,
		queries:   sqlc.New(db),
		config:    cfg,
		hub:       sse.NewHub(),
		whLimiter: ratelimit.NewLimiter(cfg.RateLimitPerWH),
		ipLimiter: ratelimit.NewLimiter(cfg.RateLimitPerIP),
	}
}
