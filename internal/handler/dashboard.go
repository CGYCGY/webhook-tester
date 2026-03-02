package handler

import (
	"net/http"

	"github.com/cgy/webhook-tester/internal/templates"
)

func Settings(w http.ResponseWriter, r *http.Request) {
	templates.SettingsPage().Render(r.Context(), w)
}
