package middleware

import (
	"context"
	"net/http"
	"strings"
	"todo-app/internal/pkg/auth"
)

type contextKey string
const UserContextKey contextKey = "user"

type AuthMiddleware struct{}

func NewAuthMiddleware() *AuthMiddleware {
    return &AuthMiddleware{}
}

func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
            return
        }

        token := parts[1]

        claims, err := auth.ValidateToken(token)
        if err != nil {
            http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), UserContextKey, claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}

func GetUserFromContext(ctx context.Context) (*auth.Claims, bool) {
    claims, ok := ctx.Value(UserContextKey).(*auth.Claims)
    return claims, ok
}