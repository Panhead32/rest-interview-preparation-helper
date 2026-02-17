package auth

import (
	"context"
	model "interview-project/internal/models/response"
	"interview-project/pkg/utils"
	"net/http"
	"strings"
)

// AuthMiddleware validates JWT token and extracts userID into request context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteJSON(w, http.StatusUnauthorized, model.NewErrorResponse("authorization required", "missing authorization header"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.WriteJSON(w, http.StatusUnauthorized, model.NewErrorResponse("authorization required", "invalid authorization header format"))
			return
		}

		tokenString := parts[1]

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, model.NewErrorResponse("authorization failed", err.Error()))
			return
		}

		ctx := context.WithValue(r.Context(), utils.UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
