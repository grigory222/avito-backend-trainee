package handlers

import (
	"encoding/json"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/common"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/dto"
	"github.com/grigory222/avito-backend-trainee/pkg/logger"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	startDateParam, startDateParsingError := getQueryParam(r, "startDate", "")
	endDateParam, endDateParsingError := getQueryParam(r, "endDate", "")
	if startDateParsingError != nil || endDateParsingError != nil {
		common.WriteError(w, http.StatusBadRequest, dto.ErrorDto{Message: "Ошибка чтения параметра даты"})
		return
	}

	startDate, startDateConvError := convertToDate(startDateParam)
	endDate, endDateConvError := convertToDate(endDateParam)
	if startDateConvError != nil || endDateConvError != nil {
		common.WriteError(w, http.StatusBadRequest, dto.ErrorDto{Message: "Время должно быть в формате RFC3339: 2006-01-02T15:04:05Z07:00"})
		return
	}

	page, pageParsingError := getQueryParam(r, "page", 1)
	limit, limitParsingError := getQueryParam(r, "limit", 10)

	if pageParsingError != nil {
		common.WriteError(w, http.StatusBadRequest, dto.ErrorDto{Message: "Ошибка в параметре page"})
		return
	}
	if limitParsingError != nil {
		common.WriteError(w, http.StatusBadRequest, dto.ErrorDto{Message: "Ошибка в параметре limit"})
		return
	}

	pagination, err := h.service.GetPVZWithPagination(&startDate, &endDate, page, limit)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, dto.ErrorDto{Message: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(pagination)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (h *PVZHandler) CloseLastReception(w http.ResponseWriter, r *http.Request) {
	prefix := "/pvz/"
	suffix := "/close_last_reception"

	path := r.URL.Path
	if !strings.HasPrefix(path, prefix) || !strings.HasSuffix(path, suffix) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	pvzIdStr := strings.TrimSuffix(strings.TrimPrefix(path, prefix), suffix)
	pvzIdStr = strings.Trim(pvzIdStr, "/")

	if pvzIdStr == "" {
		http.Error(w, "pvzId is missing", http.StatusBadRequest)
		return
	}

	reception, err := h.service.CloseLastReception(pvzIdStr)
	if err != nil {
		h.logger.Error("Failed to close last reception: ", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(dto.ErrorDto{Message: "Неверный запрос или приемка уже закрыта"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(reception)
}

func (h *PVZHandler) DeleteLastProduct(w http.ResponseWriter, r *http.Request) {
	prefix := "/pvz/"
	suffix := "/delete_last_product"

	path := r.URL.Path
	if !strings.HasPrefix(path, prefix) || !strings.HasSuffix(path, suffix) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	pvzIdStr := strings.TrimSuffix(strings.TrimPrefix(path, prefix), suffix)
	pvzIdStr = strings.Trim(pvzIdStr, "/")

	if pvzIdStr == "" {
		http.Error(w, "pvzId is missing", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteLastProduct(pvzIdStr)
	if err != nil {
		h.logger.Error("Failed to delete last product: ", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(dto.ErrorDto{Message: "Неверный запрос, нет активной приемки или нет товаров для удаления"})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func convertToDate(s string) (time.Time, error) {
	date, err := time.Parse(time.RFC3339, s)
	return date, err
}

func getQueryParam[T any](r *http.Request, key string, defaultValue T) (T, error) {
	value := r.URL.Query().Get(key)

	if value == "" {
		return defaultValue, nil
	}

	switch any(defaultValue).(type) {
	case string:
		return any(value).(T), nil
	case int:
		num, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue, err
		}
		return any(num).(T), nil
	default:
		return defaultValue, nil
	}
}
