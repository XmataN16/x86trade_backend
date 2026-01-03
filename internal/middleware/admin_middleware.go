package middleware

import (
	"net/http"
	"x86trade_backend/internal/repository"
)

// AdminOnly проверяет, что текущий пользователь — админ.
// Требует, чтобы AuthMiddleware уже положил user id в контекст (UserIDFromContext).
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "authorization required", http.StatusUnauthorized)
			return
		}
		u, err := repository.GetUserByID(r.Context(), userID)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if u == nil || !u.IsAdmin {
			http.Error(w, "admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
