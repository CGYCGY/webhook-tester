package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const claimsKey contextKey = "auth_claims"

func Middleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenValue string

			cookie, err := r.Cookie("token")
			if err == nil {
				tokenValue = cookie.Value
			} else {
				authHeader := r.Header.Get("Authorization")
				if strings.HasPrefix(authHeader, "Bearer ") {
					tokenValue = strings.TrimPrefix(authHeader, "Bearer ")
				}
			}

			if tokenValue == "" {
				respondUnauthorized(w, r)
				return
			}

			claims, err := ValidateToken(tokenValue, jwtSecret)
			if err != nil {
				respondUnauthorized(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(ctx context.Context) *Claims {
	claims, _ := ctx.Value(claimsKey).(*Claims)
	return claims
}

func respondUnauthorized(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "text/html") {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":"unauthorized"}`))
}
