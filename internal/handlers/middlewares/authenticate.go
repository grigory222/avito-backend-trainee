package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/dto"
	"net/http"
	"strings"
)

type Claims struct {
	// Встроенные стандартные claims
	jwt.RegisteredClaims

	// Кастомные claims
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type contextKey string

const UserIDKey = contextKey("user_id")

func writeError(w http.ResponseWriter, status int, errorDto dto.ErrorDto) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorDto)
}

func AuthenticateMiddleware(secret string) func(http.Handler) http.Handler {
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
				writeError(w, http.StatusUnauthorized, dto.MissingOrInvalidToken)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			// Разбираем и проверяем токен с кастомными claims
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				// Убедимся, что используется HMAC
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				writeError(w, http.StatusUnauthorized, dto.UnauthorizedError)
				return
			}

			// Кладём user_id в контекст
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
