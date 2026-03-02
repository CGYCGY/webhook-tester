package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/cgy/webhook-tester/internal/auth"
	"github.com/cgy/webhook-tester/internal/config"
	"github.com/cgy/webhook-tester/internal/database/sqlc"
	"github.com/cgy/webhook-tester/internal/templates"
	"github.com/cgy/webhook-tester/internal/templates/components"
	"golang.org/x/crypto/bcrypt"
)

type SettingsHandler struct {
	queries *sqlc.Queries
	config  *config.Config
}

func NewSettings(queries *sqlc.Queries, cfg *config.Config) *SettingsHandler {
	return &SettingsHandler{
		queries: queries,
		config:  cfg,
	}
}

func (h *SettingsHandler) SettingsPage(w http.ResponseWriter, r *http.Request) {
	templates.SettingsPage("", "").Render(r.Context(), w)
}

func (h *SettingsHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		components.FormError("Invalid form submission.").Render(r.Context(), w)
		return
	}

	current := strings.TrimSpace(r.FormValue("current_password"))
	newPass := strings.TrimSpace(r.FormValue("new_password"))
	confirm := strings.TrimSpace(r.FormValue("confirm_password"))

	if current == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		components.FormError("Current password is required.").Render(r.Context(), w)
		return
	}
	if len(newPass) < 8 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		components.FormError("New password must be at least 8 characters.").Render(r.Context(), w)
		return
	}
	if newPass != confirm {
		w.WriteHeader(http.StatusUnprocessableEntity)
		components.FormError("New password and confirmation do not match.").Render(r.Context(), w)
		return
	}

	user, err := h.queries.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		components.FormError("An error occurred. Please try again.").Render(r.Context(), w)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(current)); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		components.FormError("Current password is incorrect.").Render(r.Context(), w)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		components.FormError("An error occurred. Please try again.").Render(r.Context(), w)
		return
	}

	if err := h.queries.UpdateUserPassword(r.Context(), sqlc.UpdateUserPasswordParams{
		Password:  string(hashed),
		UpdatedAt: time.Now().UTC(),
		ID:        claims.UserID,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		components.FormError("Failed to update password. Please try again.").Render(r.Context(), w)
		return
	}

	templates.PasswordChangeSuccess().Render(r.Context(), w)
}
