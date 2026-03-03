package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/contractiq/contractiq/internal/application/contract/command"
	"github.com/contractiq/contractiq/internal/application/contract/dto"
	"github.com/contractiq/contractiq/internal/application/contract/query"
	"github.com/contractiq/contractiq/internal/interfaces/http/middleware"
	"github.com/contractiq/contractiq/internal/interfaces/http/response"
	"github.com/contractiq/contractiq/internal/interfaces/http/validation"
	"github.com/contractiq/contractiq/pkg/apperror"
	"github.com/contractiq/contractiq/pkg/identifier"
)

// ContractHandler handles contract REST endpoints.
type ContractHandler struct {
	cmdHandler   *command.Handler
	queryHandler *query.Handler
}

// NewContractHandler creates a new contract handler.
func NewContractHandler(cmd *command.Handler, q *query.Handler) *ContractHandler {
	return &ContractHandler{cmdHandler: cmd, queryHandler: q}
}

// Create handles POST /api/v1/contracts.
func (h *ContractHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateContractRequest
	if err := validation.DecodeAndValidate(r, &req); err != nil {
		response.Error(w, apperror.NewValidation(err.Error()))
		return
	}

	userID := middleware.GetUserID(r.Context())
	result, err := h.cmdHandler.Create(r.Context(), userID, req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Created(w, result)
}

// Get handles GET /api/v1/contracts/{id}.
func (h *ContractHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !identifier.IsValid(id) {
		response.Error(w, apperror.NewValidation("invalid contract ID"))
		return
	}

	result, err := h.queryHandler.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, result)
}

// List handles GET /api/v1/contracts.
func (h *ContractHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	req := dto.ListContractsRequest{
		Status:  r.URL.Query().Get("status"),
		PartyID: r.URL.Query().Get("party_id"),
		Search:  r.URL.Query().Get("search"),
	}

	if p := r.URL.Query().Get("page"); p != "" {
		page, _ := strconv.Atoi(p)
		req.Page = page
	}
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		pageSize, _ := strconv.Atoi(ps)
		req.PageSize = pageSize
	}

	result, err := h.queryHandler.List(r.Context(), userID, req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, response.Paginated{
		Items: result.Items,
		Pagination: response.Pagination{
			TotalCount: result.TotalCount,
			Page:       result.Page,
			PageSize:   result.PageSize,
			TotalPages: result.TotalPages,
		},
	})
}

// Update handles PUT /api/v1/contracts/{id}.
func (h *ContractHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !identifier.IsValid(id) {
		response.Error(w, apperror.NewValidation("invalid contract ID"))
		return
	}

	var req dto.UpdateContractRequest
	if err := validation.DecodeAndValidate(r, &req); err != nil {
		response.Error(w, apperror.NewValidation(err.Error()))
		return
	}

	result, err := h.cmdHandler.Update(r.Context(), id, req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, result)
}

// Submit handles POST /api/v1/contracts/{id}/submit.
func (h *ContractHandler) Submit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !identifier.IsValid(id) {
		response.Error(w, apperror.NewValidation("invalid contract ID"))
		return
	}

	result, err := h.cmdHandler.Submit(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, result)
}

// Approve handles POST /api/v1/contracts/{id}/approve.
func (h *ContractHandler) Approve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !identifier.IsValid(id) {
		response.Error(w, apperror.NewValidation("invalid contract ID"))
		return
	}

	userID := middleware.GetUserID(r.Context())
	result, err := h.cmdHandler.Approve(r.Context(), id, userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, result)
}

// Sign handles POST /api/v1/contracts/{id}/sign.
func (h *ContractHandler) Sign(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !identifier.IsValid(id) {
		response.Error(w, apperror.NewValidation("invalid contract ID"))
		return
	}

	userID := middleware.GetUserID(r.Context())
	result, err := h.cmdHandler.Sign(r.Context(), id, userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, result)
}

// Terminate handles POST /api/v1/contracts/{id}/terminate.
func (h *ContractHandler) Terminate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !identifier.IsValid(id) {
		response.Error(w, apperror.NewValidation("invalid contract ID"))
		return
	}

	var req dto.TerminateRequest
	if err := validation.DecodeAndValidate(r, &req); err != nil {
		response.Error(w, apperror.NewValidation(err.Error()))
		return
	}

	result, err := h.cmdHandler.Terminate(r.Context(), id, req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, result)
}
