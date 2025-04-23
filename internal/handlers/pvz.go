package handlers

import (
	"encoding/json"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/common"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/dto"
	"github.com/grigory222/avito-backend-trainee/pkg/logger"
	"net/http"
)

type PVZHandler struct {
	service PVZService
	logger  *logger.Logger
}

func NewPVZHandler(service PVZService, logger *logger.Logger) *PVZHandler {
	return &PVZHandler{
		service: service,
		logger:  logger,
	}
}

func (h *PVZHandler) AddPVZ(w http.ResponseWriter, r *http.Request) {
	pvzDto, ok := r.Context().Value(common.DTOKey).(*dto.PVZDto)
	if !ok || pvzDto == nil {
		h.logger.Error("AddPvz called without dto")
		common.WriteError(w, http.StatusInternalServerError, dto.UnauthorizedError)
		return
	}
	pvz, err := h.service.AddPVZ(pvzDto.City)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, dto.ErrorDto{Message: err.Error()})
		return
	}

	pvzResponse := dto.PVZDto{
		Id:               pvz.Id,
		RegistrationDate: pvz.RegistrationDate,
		City:             pvz.City,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(pvzResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *PVZHandler) GetPVZ(w http.ResponseWriter, r *http.Request) {

}
