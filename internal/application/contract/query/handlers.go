package query

import (
	"context"

	"github.com/contractiq/contractiq/internal/application/contract/dto"
	"github.com/contractiq/contractiq/internal/domain/contract"
	"github.com/contractiq/contractiq/pkg/apperror"
)

// PagedResult is the query-layer paginated result for contracts.
type PagedResult struct {
	Items      []*dto.ContractResponse
	TotalCount int
	Page       int
	PageSize   int
	TotalPages int
}

// Handler processes contract read queries.
type Handler struct {
	repo contract.Repository
}

// NewHandler creates a new contract query handler.
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

// List retrieves a paginated, filtered list of contracts for a user.
func (h *Handler) List(ctx context.Context, ownerID string, req dto.ListContractsRequest) (*PagedResult, error) {
	filter := contract.DefaultFilter()
	filter.OwnerID = &ownerID

	if req.Page > 0 {
		filter.Page = req.Page
	}
	if req.PageSize > 0 {
		filter.PageSize = req.PageSize
	}

	if req.Status != "" {
		s := contract.Status(req.Status)
		if !s.IsValid() {
			return nil, apperror.NewValidation("invalid status filter: " + req.Status)
		}
		filter.Status = &s
	}

	if req.PartyID != "" {
		filter.PartyID = &req.PartyID
	}

	if req.Search != "" {
		filter.Search = &req.Search
	}

	result, err := h.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	items := make([]*dto.ContractResponse, 0, len(result.Items))
	for _, c := range result.Items {
		items = append(items, toResponse(c))
	}

	return &PagedResult{
		Items:      items,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages(),
	}, nil
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
		Status:      string(c.Status()),
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