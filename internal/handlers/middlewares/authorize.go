package middlewares

import (
	"github.com/grigory222/avito-backend-trainee/internal/handlers/dto"
	"net/http"
)

func AuthorizeMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(ClaimsKey).(*Claims)
			if !ok || claims == nil {
				writeError(w, http.StatusUnauthorized, dto.UnauthorizedError)
				return
			}
			if claims.UserID == "" {
				writeError(w, http.StatusUnauthorized, dto.NoUserIDProvided)
				return
			}

			// Проверка роли
			if claims.Role != requiredRole {
				writeError(w, http.StatusForbidden, dto.ForbiddenError)
				return
			}

			// Передаем управление дальше
			next.ServeHTTP(w, r)
		})
	}
}
