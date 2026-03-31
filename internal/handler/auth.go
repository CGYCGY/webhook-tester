package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cgy/webhook-tester/internal/auth"
	"github.com/cgy/webhook-tester/internal/config"
	"github.com/cgy/webhook-tester/internal/database/sqlc"
	"github.com/cgy/webhook-tester/internal/templates"
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
	if cookie, err := r.Cookie("token"); err == nil {
		if _, err := auth.ValidateToken(cookie.Value, h.config.JWTSecret); err == nil {
			http.Redirect(w, r, "/dashboard", http.StatusFound)
			return
		}
	}

	templates.LoginPage("").Render(r.Context(), w)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderLoginError(w, r, "Invalid form submission.")
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		h.renderLoginError(w, r, "Email and password are required.")
		return
	}

	user, err := h.queries.GetUserByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.renderLoginError(w, r, "Invalid email or password.")
			return
		}
		h.renderLoginError(w, r, "An error occurred. Please try again.")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		h.renderLoginError(w, r, "Invalid email or password.")
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, h.config.JWTSecret)
	if err != nil {
		h.renderLoginError(w, r, "An error occurred. Please try again.")
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

func (h *AuthHandler) APILogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	if body.Email == "" || body.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "email and password are required"})
		return
	}

	user, err := h.queries.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid email or password"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid email or password"})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, h.config.JWTSecret)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) renderLoginError(w http.ResponseWriter, r *http.Request, errMsg string) {
	w.WriteHeader(http.StatusUnauthorized)
	templates.LoginPage(errMsg).Render(r.Context(), w)
}
