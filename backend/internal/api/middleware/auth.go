package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/services"
)

type contextKey string

const UserContextKey contextKey = "user"

func Auth(userService services.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_token")
			if err != nil {
				writeJSONError(w, "Unauthorized: Session cookie missing", http.StatusUnauthorized)
				return
			}

			user, err := userService.Authenticate(cookie.Value)
			if err != nil {
				writeJSONError(w, "Unauthorized: Invalid or expired session", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext helper to retrieve user model from context.
func GetUserFromContext(ctx context.Context) (*models.User, bool) {
	u, ok := ctx.Value(UserContextKey).(*models.User)
	return u, ok
}

func writeJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error":message})
}