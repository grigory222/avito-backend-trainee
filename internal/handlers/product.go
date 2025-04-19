package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/dto"
	"github.com/grigory222/avito-backend-trainee/pkg/logger"
	"net/http"
)

type ProductHandler struct {
	service ProductService
	logger  *logger.Logger
}

func NewProductHandler(service ProductService, logger *logger.Logger) *ProductHandler {
	return &ProductHandler{
		service: service,
		logger:  logger,
	}
}

func (ph *ProductHandler) Hello(w http.ResponseWriter, r *http.Request) {
	// Логирование запроса
	ph.logger.Debug("Hello endpoint called")

	// Ответ на запрос
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintln(w, `{"message": "Hello, world!"}`)
	if err != nil {
		return
	}
}

func (ph *ProductHandler) AddProductToReception(w http.ResponseWriter, r *http.Request) {
	var addProductRequestDto dto.AddProductRequestDto

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&addProductRequestDto)
	if err != nil {
		ph.logger.Error("failed to decode request body ", err.Error())
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

}
