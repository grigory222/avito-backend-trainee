package common

import (
	"encoding/json"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/dto"
	"net/http"
)

type contextKey string

const ClaimsKey = contextKey("claims")
const DTOKey = contextKey("dto")

func WriteError(w http.ResponseWriter, status int, errorDto dto.ErrorDto) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorDto)
}
