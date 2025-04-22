package middlewares

import (
	"context"
	"encoding/json"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/common"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/dto"
	"net/http"
)

func DecoderMiddleware[DTO any]() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var inputDto DTO
			err := json.NewDecoder(r.Body).Decode(&inputDto)
			if err != nil {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				errorDto := &dto.ErrorDto{
					Message: "Некорректные данные",
				}
				err = json.NewEncoder(w).Encode(errorDto)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}
			ctx := context.WithValue(r.Context(), common.DTOKey, &inputDto)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
