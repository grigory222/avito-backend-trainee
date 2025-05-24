package handlers

import (
	"encoding/json"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/common"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/dto"
	"github.com/grigory222/avito-backend-trainee/internal/myjwt"
	"github.com/grigory222/avito-backend-trainee/pkg/logger"
	"net/http"
)

type UserHandler struct {
	service  UserService
	provider *myjwt.Provider
	logger   *logger.Logger
}

func NewUserHandler(service UserService, provider *myjwt.Provider, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		service:  service,
		provider: provider,
		logger:   logger,
	}
}

func (h *UserHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	roleStruct, ok := r.Context().Value(common.DTOKey).(*dto.RoleStruct)
	if !ok || roleStruct == nil {
		h.logger.Error("DummyLogin called without dto")
		common.WriteError(w, http.StatusBadRequest, dto.UnauthorizedError)
		return
	}

	if roleStruct.Role != "employee" && roleStruct.Role != "moderator" {
		common.WriteError(w, http.StatusBadRequest, dto.ErrorDto{Message: "Wrong role provided"})
		return
	}

	token, err := h.provider.GenerateTokenWithRole(roleStruct.Role)
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, dto.UnauthorizedError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	userRequest, ok := r.Context().Value(common.DTOKey).(*dto.UserCreateRequest)
	if !ok || userRequest == nil {
		h.logger.Error("Register called without dto")
		common.WriteError(w, http.StatusBadRequest, dto.ErrorDto{Message: "Wrong data"})
		return
	}
	user, err := h.service.UserRegister(userRequest.Email, userRequest.Password, userRequest.Role)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, dto.ErrorDto{Message: err.Error()})
		return
	}

	userResponse := &dto.UserCreateResponse{
		Id:    user.Id,
		Email: user.Email,
		Role:  user.Role,
	}

	err = json.NewEncoder(w).Encode(userResponse)
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, dto.ErrorDto{Message: err.Error()})
		return
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	userRequest, ok := r.Context().Value(common.DTOKey).(*dto.UserLoginRequest)
	if !ok || userRequest == nil {
		h.logger.Error("Login called without dto")
		common.WriteError(w, http.StatusBadRequest, dto.ErrorDto{Message: "Wrong data"})
		return
	}

	user, err := h.service.UserLogin(userRequest.Email, userRequest.Password)

	if err != nil {
		common.WriteError(w, http.StatusBadRequest, dto.ErrorDto{Message: err.Error()})
		return
	}

	token, err := h.provider.GenerateTokenWithRole(user.Role)
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, dto.ErrorDto{Message: err.Error()})
		return
	}

	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, dto.ErrorDto{Message: err.Error()})
		return
	}
}
