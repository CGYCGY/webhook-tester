package server

import (
	"database/sql"

	"github.com/cgy/webhook-tester/internal/config"
	"github.com/cgy/webhook-tester/internal/database/sqlc"
)

type Server struct {
	db      *sql.DB
	queries *sqlc.Queries
	config  *config.Config
}

func New(db *sql.DB, cfg *config.Config) *Server {
	return &Server{
		db:      db,
		queries: sqlc.New(db),
		config:  cfg,
	}
}
