package handlers

import (
	"context"
	"net/http"
	"strings"

	"x86trade_backend/utils"
)

type ctxKey string

const ctxUserIDKey ctxKey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "authorization required", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}
		token := parts[1]
		claims, err := utils.ParseAccessToken(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		// put user id into context
		ctx := context.WithValue(r.Context(), ctxUserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) (int, bool) {
	v := ctx.Value(ctxUserIDKey)
	if v == nil {
		return 0, false
	}
	id, ok := v.(int)
	return id, ok
}
