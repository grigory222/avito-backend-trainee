package middlewares

import (
	"github.com/grigory222/avito-backend-trainee/internal/handlers/common"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/dto"
	"github.com/grigory222/avito-backend-trainee/internal/myjwt"
	"net/http"
)

func AuthorizeMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(common.ClaimsKey).(*myjwt.Claims)
			if !ok || claims == nil {
				common.WriteError(w, http.StatusUnauthorized, dto.UnauthorizedError)
				return
			}

			if _, err := claims.GetSubject(); err != nil {
				common.WriteError(w, http.StatusUnauthorized, dto.NoUserIDProvided)
				return
			}

			// Проверка роли
			if claims.Role != requiredRole {
				common.WriteError(w, http.StatusForbidden, dto.ForbiddenError)
				return
			}

			// Передаем управление дальше
			next.ServeHTTP(w, r)
		})
	}
}
