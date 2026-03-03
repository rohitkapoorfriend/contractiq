package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	appparty "github.com/contractiq/contractiq/internal/application/party"
	"github.com/contractiq/contractiq/internal/interfaces/http/middleware"
	"github.com/contractiq/contractiq/internal/interfaces/http/response"
	"github.com/contractiq/contractiq/internal/interfaces/http/validation"
	"github.com/contractiq/contractiq/pkg/apperror"
	"github.com/contractiq/contractiq/pkg/identifier"
)

// PartyHandler handles party REST endpoints.
type PartyHandler struct {
	service *appparty.Service
}

// NewPartyHandler creates a new party handler.
func NewPartyHandler(service *appparty.Service) *PartyHandler {
	return &PartyHandler{service: service}
}

// Create handles POST /api/v1/parties.
func (h *PartyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req appparty.CreateRequest
	if err := validation.DecodeAndValidate(r, &req); err != nil {
		response.Error(w, apperror.NewValidation(err.Error()))
		return
	}

	userID := middleware.GetUserID(r.Context())
	result, err := h.service.Create(r.Context(), userID, req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Created(w, result)
}

// Get handles GET /api/v1/parties/{id}.
func (h *PartyHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !identifier.IsValid(id) {
		response.Error(w, apperror.NewValidation("invalid party ID"))
		return
	}

	result, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, result)
}

// List handles GET /api/v1/parties.
func (h *PartyHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	result, err := h.service.List(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, result)
}

// Update handles PUT /api/v1/parties/{id}.
func (h *PartyHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !identifier.IsValid(id) {
		response.Error(w, apperror.NewValidation("invalid party ID"))
		return
	}

	var req appparty.UpdateRequest
	if err := validation.DecodeAndValidate(r, &req); err != nil {
		response.Error(w, apperror.NewValidation(err.Error()))
		return
	}

	result, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, result)
}

// Delete handles DELETE /api/v1/parties/{id}.
func (h *PartyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !identifier.IsValid(id) {
		response.Error(w, apperror.NewValidation("invalid party ID"))
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}

	response.NoContent(w)
}
