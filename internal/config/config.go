package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port           int
	DataDir        string
	AdminEmail     string
	AdminPassword  string
	JWTSecret      string
	MaxBodySize    int64
	RateLimitPerWH int
	RateLimitPerIP int
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:           8090,
		DataDir:        "/data",
		MaxBodySize:    1 << 20, // 1MB
		RateLimitPerWH: 60,
		RateLimitPerIP: 120,
	}

	if v := os.Getenv("PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid PORT: %w", err)
		}
		cfg.Port = p
	}

	if v := os.Getenv("DATA_DIR"); v != "" {
		cfg.DataDir = v
	}

	cfg.AdminEmail = os.Getenv("ADMIN_EMAIL")
	cfg.AdminPassword = os.Getenv("ADMIN_PASSWORD")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	if v := os.Getenv("MAX_BODY_SIZE"); v != "" {
		s, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid MAX_BODY_SIZE: %w", err)
		}
		cfg.MaxBodySize = s
	}

	if v := os.Getenv("RATE_LIMIT_PER_WEBHOOK"); v != "" {
		r, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid RATE_LIMIT_PER_WEBHOOK: %w", err)
		}
		cfg.RateLimitPerWH = r
	}

	if v := os.Getenv("RATE_LIMIT_PER_IP"); v != "" {
		r, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid RATE_LIMIT_PER_IP: %w", err)
		}
		cfg.RateLimitPerIP = r
	}

	return cfg, nil
}
