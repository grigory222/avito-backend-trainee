package middlewares

import (
	"context"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/common"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/dto"
	myjwt "github.com/grigory222/avito-backend-trainee/internal/jwt"
	"net/http"
	"strings"
)

func AuthenticateMiddleware(provider *myjwt.Provider) func(http.Handler) http.Handler {
	skipPaths := map[string]bool{
		"/login":      true,
		"/register":   true,
		"/dummyLogin": true,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if skipPaths[r.URL.Path] {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				common.WriteError(w, http.StatusUnauthorized, dto.MissingOrInvalidToken)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := provider.VerifyToken(tokenStr)
			if err != nil {
				common.WriteError(w, http.StatusUnauthorized, dto.UnauthorizedError)
				return
			}

			ctx := context.WithValue(r.Context(), common.ClaimsKey, claims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
