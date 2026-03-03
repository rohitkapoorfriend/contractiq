package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	apptemplate "github.com/contractiq/contractiq/internal/application/template"
	"github.com/contractiq/contractiq/internal/interfaces/http/middleware"
	"github.com/contractiq/contractiq/internal/interfaces/http/response"
	"github.com/contractiq/contractiq/internal/interfaces/http/validation"
	"github.com/contractiq/contractiq/pkg/apperror"
	"github.com/contractiq/contractiq/pkg/identifier"
)

// TemplateHandler handles template REST endpoints.
type TemplateHandler struct {
	service *apptemplate.Service
}

// NewTemplateHandler creates a new template handler.
func NewTemplateHandler(service *apptemplate.Service) *TemplateHandler {
	return &TemplateHandler{service: service}
}

// Create handles POST /api/v1/templates.
func (h *TemplateHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req apptemplate.CreateRequest
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

// Get handles GET /api/v1/templates/{id}.
func (h *TemplateHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !identifier.IsValid(id) {
		response.Error(w, apperror.NewValidation("invalid template ID"))
		return
	}

	result, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, result)
}

// List handles GET /api/v1/templates.
func (h *TemplateHandler) List(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active_only") == "true"

	result, err := h.service.List(r.Context(), activeOnly)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, result)
}

// Update handles PUT /api/v1/templates/{id}.
func (h *TemplateHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !identifier.IsValid(id) {
		response.Error(w, apperror.NewValidation("invalid template ID"))
		return
	}

	var req apptemplate.UpdateRequest
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

// Delete handles DELETE /api/v1/templates/{id}.
func (h *TemplateHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !identifier.IsValid(id) {
		response.Error(w, apperror.NewValidation("invalid template ID"))
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}

	response.NoContent(w)
}
