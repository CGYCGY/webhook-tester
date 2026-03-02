package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/cgy/webhook-tester/internal/auth"
	"github.com/cgy/webhook-tester/internal/config"
	"github.com/cgy/webhook-tester/internal/database/sqlc"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	queries *sqlc.Queries
	config  *config.Config
}

func NewAuth(queries *sqlc.Queries, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		queries: queries,
		config:  cfg,
	}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	// If user already has a valid token, redirect to dashboard.
	if cookie, err := r.Cookie("token"); err == nil {
		if _, err := auth.ValidateToken(cookie.Value, h.config.JWTSecret); err == nil {
			http.Redirect(w, r, "/dashboard", http.StatusFound)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := loginTmpl.Execute(w, loginData{}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderLoginError(w, "Invalid form submission.")
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		h.renderLoginError(w, "Email and password are required.")
		return
	}

	user, err := h.queries.GetUserByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.renderLoginError(w, "Invalid email or password.")
			return
		}
		h.renderLoginError(w, "An error occurred. Please try again.")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		h.renderLoginError(w, "Invalid email or password.")
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, h.config.JWTSecret)
	if err != nil {
		h.renderLoginError(w, "An error occurred. Please try again.")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		MaxAge:   86400,
	})

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (h *AuthHandler) renderLoginError(w http.ResponseWriter, errMsg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	if err := loginTmpl.Execute(w, loginData{Error: errMsg}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
