package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/common"
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
	addProductRequestDto, ok := r.Context().Value(common.DTOKey).(*dto.AddProductRequestDto)
	if !ok || addProductRequestDto == nil {
		ph.logger.Error("AddProductToReception called without addProductRequestDto")
		common.WriteError(w, http.StatusInternalServerError, dto.UnauthorizedError)
		return
	}

	product, err := ph.service.AddProduct(addProductRequestDto.Type, addProductRequestDto.PVZId)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, dto.ErrorDto{Message: err.Error()})
		return
	}

	addProductResponseDto := &dto.AddProductResponseDto{
		Id:          product.Id,
		DateTime:    product.DateTime,
		Type:        product.Type,
		ReceptionId: product.ReceptionId,
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(addProductResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
