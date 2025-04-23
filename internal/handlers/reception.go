package handlers

import (
	"encoding/json"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/common"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/dto"
	"github.com/grigory222/avito-backend-trainee/pkg/logger"
	"net/http"
)

type ReceptionHandler struct {
	service ReceptionService
	logger  *logger.Logger
}

func NewReceptionHandler(service ReceptionService, logger *logger.Logger) *ReceptionHandler {
	return &ReceptionHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ReceptionHandler) AddReception(w http.ResponseWriter, r *http.Request) {
	recDto, ok := r.Context().Value(common.DTOKey).(*dto.CreateReceptionRequestDto)
	if !ok || recDto == nil {
		h.logger.Error("CreateReceptionRequestDto called without dto")
		common.WriteError(w, http.StatusInternalServerError, dto.UnauthorizedError)
		return
	}

	rec, err := h.service.AddReception(recDto.PVZId)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, dto.ErrorDto{Message: err.Error()})
		return
	}

	response := dto.ReceptionDto{
		Id:       rec.Id,
		DateTime: rec.DateTime,
		PVZId:    rec.PVZId,
		Status:   rec.Status,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
