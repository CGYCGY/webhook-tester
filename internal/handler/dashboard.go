package handler

import (
	"net/http"

	"github.com/cgy/webhook-tester/internal/templates"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	templates.DashboardPage().Render(r.Context(), w)
}

func Settings(w http.ResponseWriter, r *http.Request) {
	templates.SettingsPage().Render(r.Context(), w)
}
