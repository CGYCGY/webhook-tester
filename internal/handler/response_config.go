package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ResponseConfig holds the custom response settings for a webhook.
type ResponseConfig struct {
	Status      int    `json:"status"`
	ContentType string `json:"content_type"`
	Body        string `json:"body"`
}

// parseResponseConfig parses the stored JSON into a ResponseConfig.
// Returns (cfg, false) when the raw value is empty/unconfigured ("{}").
func parseResponseConfig(raw string) (ResponseConfig, bool) {
	if raw == "" || raw == "{}" {
		return ResponseConfig{}, false
	}
	var cfg ResponseConfig
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		return ResponseConfig{}, false
	}
	if cfg.Status == 0 {
		cfg.Status = 200
	}
	return cfg, true
}

// resolveTemplate replaces template placeholders in the body string with
// values from the incoming request:
//   - {{query.X}}  → query parameter X
//   - {{header.X}} → request header X (canonical form, e.g. Content-Type)
//   - {{body}}     → raw request body
func resolveTemplate(tmpl string, r *http.Request, body []byte) string {
	result := tmpl
	for k, vs := range r.URL.Query() {
		result = strings.ReplaceAll(result, "{{query."+k+"}}", vs[0])
	}
	for k, vs := range r.Header {
		result = strings.ReplaceAll(result, "{{header."+k+"}}", vs[0])
	}
	result = strings.ReplaceAll(result, "{{body}}", string(body))
	return result
}
