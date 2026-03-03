package query

import (
	"context"

	"github.com/contractiq/contractiq/internal/application/contract/dto"
	"github.com/contractiq/contractiq/internal/domain/contract"
)

// Handler processes contract queries (read operations).
type Handler struct {
	repo contract.Repository
}

// NewHandler creates a new query handler.
func NewHandler(repo contract.Repository) *Handler {
	return &Handler{repo: repo}
}

// GetByID retrieves a single contract by ID.
func (h *Handler) GetByID(ctx context.Context, id string) (*dto.ContractResponse, error) {
	c, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toResponse(c), nil
}

// List retrieves contracts matching filter criteria.
func (h *Handler) List(ctx context.Context, ownerID string, req dto.ListContractsRequest) (*PagedResponse, error) {
	filter := contract.DefaultFilter()
	filter.OwnerID = &ownerID

	if req.Status != "" {
		status := contract.Status(req.Status)
		filter.Status = &status
	}
	if req.PartyID != "" {
		filter.PartyID = &req.PartyID
	}
	if req.Search != "" {
		filter.Search = &req.Search
	}
	if req.Page > 0 {
		filter.Page = req.Page
	}
	if req.PageSize > 0 {
		filter.PageSize = req.PageSize
	}

	result, err := h.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ContractResponse, 0, len(result.Items))
	for _, c := range result.Items {
		items = append(items, *toResponse(c))
	}

	return &PagedResponse{
		Items:      items,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages(),
	}, nil
}

// PagedResponse wraps a page of contract responses with metadata.
type PagedResponse struct {
	Items      []dto.ContractResponse `json:"items"`
	TotalCount int                    `json:"total_count"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	TotalPages int                    `json:"total_pages"`
}

func toResponse(c *contract.Contract) *dto.ContractResponse {
	clauses := make([]dto.ClauseDTO, 0, len(c.Clauses()))
	for _, cl := range c.Clauses() {
		clauses = append(clauses, dto.ClauseDTO{
			Title:   cl.Title,
			Content: cl.Content,
			Order:   cl.Order,
		})
	}

	return &dto.ContractResponse{
		ID:          c.ID(),
		Title:       c.Title(),
		Description: c.Description(),
		Status:      c.Status().String(),
		Value: dto.MoneyDTO{
			AmountCents: c.Value().AmountCents,
			Currency:    c.Value().Currency,
		},
		Clauses:    clauses,
		StartDate:  c.DateRange().Start,
		EndDate:    c.DateRange().End,
		OwnerID:    c.OwnerID(),
		PartyID:    c.PartyID(),
		TemplateID: c.TemplateID(),
		Version:    c.Version(),
		CreatedAt:  c.CreatedAt(),
		UpdatedAt:  c.UpdatedAt(),
	}
}
