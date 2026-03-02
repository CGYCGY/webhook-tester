package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cgy/webhook-tester/internal/config"
	"github.com/cgy/webhook-tester/internal/database"
	"github.com/cgy/webhook-tester/internal/seed"
	"github.com/cgy/webhook-tester/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		log.Fatalf("failed to create data dir: %v", err)
	}

	dbPath := filepath.Join(cfg.DataDir, "webhook-tester.db")
	db, err := database.Open(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	if cfg.AdminEmail != "" && cfg.AdminPassword != "" {
		if err := seed.SeedAdmin(db, cfg.AdminEmail, cfg.AdminPassword); err != nil {
			log.Fatalf("failed to seed admin: %v", err)
		}
	}

	srv := server.New(db, cfg)
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, srv.Router()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
