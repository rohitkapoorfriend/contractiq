package command

import (
	"context"
	"time"

	"github.com/contractiq/contractiq/internal/application/contract/dto"
	"github.com/contractiq/contractiq/internal/application/unitofwork"
	"github.com/contractiq/contractiq/internal/domain/contract"
	"github.com/contractiq/contractiq/internal/domain/event"
	"github.com/contractiq/contractiq/pkg/apperror"
	"github.com/contractiq/contractiq/pkg/clock"
)

// Handler processes contract commands (write operations).
type Handler struct {
	uow       unitofwork.UnitOfWork
	clock     clock.Clock
	publisher event.Publisher
}

// NewHandler creates a new command handler.
func NewHandler(uow unitofwork.UnitOfWork, clock clock.Clock, publisher event.Publisher) *Handler {
	return &Handler{uow: uow, clock: clock, publisher: publisher}
}

// Create creates a new contract.
func (h *Handler) Create(ctx context.Context, ownerID string, req dto.CreateContractRequest) (*dto.ContractResponse, error) {
	now := h.clock.Now()

	value, err := contract.NewMoney(req.Value.AmountCents, req.Value.Currency)
	if err != nil {
		return nil, apperror.NewValidation(err.Error())
	}

	dateRange, err := contract.NewDateRange(req.StartDate, req.EndDate)
	if err != nil {
		return nil, apperror.NewValidation(err.Error())
	}

	clauses := make([]contract.Clause, 0, len(req.Clauses))
	for _, c := range req.Clauses {
		clause, err := contract.NewClause(c.Title, c.Content, c.Order)
		if err != nil {
			return nil, apperror.NewValidation(err.Error())
		}
		clauses = append(clauses, clause)
	}

	var result *dto.ContractResponse

	err = h.uow.Do(ctx, func(repos unitofwork.Repositories) error {
		var c *contract.Contract
		var createErr error

		if req.TemplateID != "" {
			c, createErr = contract.NewContractFromTemplate(
				req.Title, req.Description, ownerID, req.TemplateID,
				clauses, value, dateRange, now,
			)
		} else {
			c, createErr = contract.NewContract(req.Title, req.Description, ownerID, value, dateRange, now)
			if createErr == nil && len(clauses) > 0 {
				createErr = c.SetClauses(clauses, now)
			}
		}
		if createErr != nil {
			return apperror.NewValidation(createErr.Error())
		}

		if req.PartyID != "" {
			if err := c.SetParty(req.PartyID); err != nil {
				return apperror.NewValidation(err.Error())
			}
		}

		if err := repos.Contracts.Save(ctx, c); err != nil {
			return err
		}

		h.publisher.Publish(c.Events()...)
		result = toContractResponse(c)
		return nil
	})

	return result, err
}

// Update modifies an existing draft contract.
func (h *Handler) Update(ctx context.Context, id string, req dto.UpdateContractRequest) (*dto.ContractResponse, error) {
	now := h.clock.Now()

	value, err := contract.NewMoney(req.Value.AmountCents, req.Value.Currency)
	if err != nil {
		return nil, apperror.NewValidation(err.Error())
	}

	dateRange, err := contract.NewDateRange(req.StartDate, req.EndDate)
	if err != nil {
		return nil, apperror.NewValidation(err.Error())
	}

	clauses := make([]contract.Clause, 0, len(req.Clauses))
	for _, c := range req.Clauses {
		clause, err := contract.NewClause(c.Title, c.Content, c.Order)
		if err != nil {
			return nil, apperror.NewValidation(err.Error())
		}
		clauses = append(clauses, clause)
	}

	var result *dto.ContractResponse

	err = h.uow.Do(ctx, func(repos unitofwork.Repositories) error {
		c, err := repos.Contracts.FindByID(ctx, id)
		if err != nil {
			return err
		}

		if c.Version() != req.Version {
			return apperror.NewConcurrencyConflict("contract")
		}

		if err := c.Update(req.Title, req.Description, value, dateRange, now); err != nil {
			return apperror.NewConflict(err.Error())
		}

		if err := c.SetClauses(clauses, now); err != nil {
			return apperror.NewConflict(err.Error())
		}

		if err := repos.Contracts.Update(ctx, c); err != nil {
			return err
		}

		result = toContractResponse(c)
		return nil
	})

	return result, err
}

// Submit transitions a contract to PendingReview.
func (h *Handler) Submit(ctx context.Context, id string) (*dto.ContractResponse, error) {
	return h.doTransition(ctx, id, func(c *contract.Contract, now time.Time) error {
		return c.Submit(now)
	})
}

// Approve transitions a contract to Approved.
func (h *Handler) Approve(ctx context.Context, id, approverID string) (*dto.ContractResponse, error) {
	return h.doTransition(ctx, id, func(c *contract.Contract, now time.Time) error {
		return c.Approve(approverID, now)
	})
}

// Sign transitions a contract to Active.
func (h *Handler) Sign(ctx context.Context, id, signerID string) (*dto.ContractResponse, error) {
	return h.doTransition(ctx, id, func(c *contract.Contract, now time.Time) error {
		return c.Sign(signerID, now)
	})
}

// Terminate ends an active contract.
func (h *Handler) Terminate(ctx context.Context, id string, req dto.TerminateRequest) (*dto.ContractResponse, error) {
	return h.doTransition(ctx, id, func(c *contract.Contract, now time.Time) error {
		return c.Terminate(req.Reason, now)
	})
}

func (h *Handler) doTransition(ctx context.Context, id string, transition func(*contract.Contract, time.Time) error) (*dto.ContractResponse, error) {
	now := h.clock.Now()
	var result *dto.ContractResponse

	err := h.uow.Do(ctx, func(repos unitofwork.Repositories) error {
		c, err := repos.Contracts.FindByID(ctx, id)
		if err != nil {
			return err
		}

		if err := transition(c, now); err != nil {
			return apperror.NewConflict(err.Error())
		}

		if err := repos.Contracts.Update(ctx, c); err != nil {
			return err
		}

		h.publisher.Publish(c.Events()...)
		result = toContractResponse(c)
		return nil
	})

	return result, err
}

func toContractResponse(c *contract.Contract) *dto.ContractResponse {
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
