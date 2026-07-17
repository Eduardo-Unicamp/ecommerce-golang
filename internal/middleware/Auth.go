package middleware

import (
	"context"
	"first-api/internal/auth"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "userID"

func RequireAuth(jwtConfig *auth.JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{error:"acesso não autorizado"}`, http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := auth.ValidateToken(tokenStr, jwtConfig)
			if err != nil {
				http.Error(w, `{"error":"Token inválido ou expirado"}`, http.StatusUnauthorized)
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.Subject)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromToken(ctx context.Context) string {

	if id, ok := ctx.Value(UserIDKey).(string); ok {
		return id
	}
	return ""
}
