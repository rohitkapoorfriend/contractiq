package handler

import (
	"net/http"

	"github.com/contractiq/contractiq/internal/infrastructure/auth"
	"github.com/contractiq/contractiq/internal/interfaces/http/response"
	"github.com/contractiq/contractiq/internal/interfaces/http/validation"
	"github.com/contractiq/contractiq/pkg/apperror"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	userService *auth.UserService
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(userService *auth.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// Register handles POST /api/v1/auth/register.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req auth.RegisterRequest
	if err := validation.DecodeAndValidate(r, &req); err != nil {
		response.Error(w, apperror.NewValidation(err.Error()))
		return
	}

	result, err := h.userService.Register(r.Context(), req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Created(w, result)
}

// Login handles POST /api/v1/auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest
	if err := validation.DecodeAndValidate(r, &req); err != nil {
		response.Error(w, apperror.NewValidation(err.Error()))
		return
	}

	result, err := h.userService.Login(r.Context(), req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, result)
}
