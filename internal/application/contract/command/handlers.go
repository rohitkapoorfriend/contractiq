package command

import (
	"context"

	"github.com/contractiq/contractiq/internal/application/contract/dto"
	"github.com/contractiq/contractiq/internal/application/unitofwork"
	"github.com/contractiq/contractiq/internal/domain/contract"
	"github.com/contractiq/contractiq/internal/domain/event"
	"github.com/contractiq/contractiq/pkg/apperror"
	"github.com/contractiq/contractiq/pkg/clock"
)

// Handler processes contract write commands.
type Handler struct {
	uow       unitofwork.UnitOfWork
	clock     clock.Clock
	publisher event.Publisher
}

// NewHandler creates a new contract command handler.
func NewHandler(uow unitofwork.UnitOfWork, clock clock.Clock, publisher event.Publisher) *Handler {
	return &Handler{uow: uow, clock: clock, publisher: publisher}
}

// Create creates a new draft contract.
func (h *Handler) Create(ctx context.Context, ownerID string, req dto.CreateContractRequest) (*dto.ContractResponse, error) {
	now := h.clock.Now()

	money, err := contract.NewMoney(req.Value.AmountCents, req.Value.Currency)
	if err != nil {
		return nil, apperror.NewValidation(err.Error())
	}

	dateRange, err := contract.NewDateRange(req.StartDate, req.EndDate)
	if err != nil {
		return nil, apperror.NewValidation(err.Error())
	}

	var c *contract.Contract

	dbErr := h.uow.Do(ctx, func(repos unitofwork.Repositories) error {
		var clauses []contract.Clause
		if req.TemplateID != "" {
			tmpl, err := repos.Templates.FindByID(ctx, req.TemplateID)
			if err != nil {
				return err
			}
			clauses = tmpl.Clauses()

			c, err = contract.NewContractFromTemplate(
				req.Title, req.Description, ownerID, req.TemplateID,
				clauses, money, dateRange, now,
			)
			if err != nil {
				return apperror.NewValidation(err.Error())
			}
		} else {
			var buildErr error
			c, buildErr = contract.NewContract(req.Title, req.Description, ownerID, money, dateRange, now)
			if buildErr != nil {
				return apperror.NewValidation(buildErr.Error())
			}
		}

		// Build clauses from request if provided directly (and not from template)
		if req.TemplateID == "" && len(req.Clauses) > 0 {
			domainClauses := make([]contract.Clause, 0, len(req.Clauses))
			for _, cl := range req.Clauses {
				dc, err := contract.NewClause(cl.Title, cl.Content, cl.Order)
				if err != nil {
					return apperror.NewValidation(err.Error())
				}
				domainClauses = append(domainClauses, dc)
			}
			if err := c.SetClauses(domainClauses, now); err != nil {
				return apperror.NewValidation(err.Error())
			}
		}

		if req.PartyID != "" {
			if _, err := repos.Parties.FindByID(ctx, req.PartyID); err != nil {
				return err
			}
			if err := c.SetParty(req.PartyID); err != nil {
				return apperror.NewValidation(err.Error())
			}
		}

		return repos.Contracts.Save(ctx, c)
	})

	if dbErr != nil {
		return nil, dbErr
	}

	h.publisher.Publish(c.Events()...)
	return toResponse(c), nil
}

// Update modifies a draft contract.
func (h *Handler) Update(ctx context.Context, id string, req dto.UpdateContractRequest) (*dto.ContractResponse, error) {
	now := h.clock.Now()

	money, err := contract.NewMoney(req.Value.AmountCents, req.Value.Currency)
	if err != nil {
		return nil, apperror.NewValidation(err.Error())
	}

	dateRange, err := contract.NewDateRange(req.StartDate, req.EndDate)
	if err != nil {
		return nil, apperror.NewValidation(err.Error())
	}

	var c *contract.Contract

	dbErr := h.uow.Do(ctx, func(repos unitofwork.Repositories) error {
		var findErr error
		c, findErr = repos.Contracts.FindByID(ctx, id)
		if findErr != nil {
			return findErr
		}

		if c.Version() != req.Version {
			return apperror.NewConcurrencyConflict("contract")
		}

		if err := c.Update(req.Title, req.Description, money, dateRange, now); err != nil {
			return apperror.NewValidation(err.Error())
		}

		if len(req.Clauses) > 0 {
			domainClauses := make([]contract.Clause, 0, len(req.Clauses))
			for _, cl := range req.Clauses {
				dc, err := contract.NewClause(cl.Title, cl.Content, cl.Order)
				if err != nil {
					return apperror.NewValidation(err.Error())
				}
				domainClauses = append(domainClauses, dc)
			}
			if err := c.SetClauses(domainClauses, now); err != nil {
				return apperror.NewValidation(err.Error())
			}
		}

		return repos.Contracts.Update(ctx, c)
	})

	if dbErr != nil {
		return nil, dbErr
	}

	return toResponse(c), nil
}

// Submit moves a contract from Draft to PendingReview.
func (h *Handler) Submit(ctx context.Context, id string) (*dto.ContractResponse, error) {
	now := h.clock.Now()
	var c *contract.Contract

	dbErr := h.uow.Do(ctx, func(repos unitofwork.Repositories) error {
		var err error
		c, err = repos.Contracts.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if err := c.Submit(now); err != nil {
			return apperror.NewConflict(err.Error())
		}
		return repos.Contracts.Update(ctx, c)
	})

	if dbErr != nil {
		return nil, dbErr
	}

	h.publisher.Publish(c.Events()...)
	return toResponse(c), nil
}

// Approve moves a contract from PendingReview to Approved.
func (h *Handler) Approve(ctx context.Context, id, approverID string) (*dto.ContractResponse, error) {
	now := h.clock.Now()
	var c *contract.Contract

	dbErr := h.uow.Do(ctx, func(repos unitofwork.Repositories) error {
		var err error
		c, err = repos.Contracts.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if err := c.Approve(approverID, now); err != nil {
			return apperror.NewConflict(err.Error())
		}
		return repos.Contracts.Update(ctx, c)
	})

	if dbErr != nil {
		return nil, dbErr
	}

	h.publisher.Publish(c.Events()...)
	return toResponse(c), nil
}

// Sign activates an approved contract.
func (h *Handler) Sign(ctx context.Context, id, signerID string) (*dto.ContractResponse, error) {
	now := h.clock.Now()
	var c *contract.Contract

	dbErr := h.uow.Do(ctx, func(repos unitofwork.Repositories) error {
		var err error
		c, err = repos.Contracts.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if err := c.Sign(signerID, now); err != nil {
			return apperror.NewConflict(err.Error())
		}
		return repos.Contracts.Update(ctx, c)
	})

	if dbErr != nil {
		return nil, dbErr
	}

	h.publisher.Publish(c.Events()...)
	return toResponse(c), nil
}

// Terminate ends an active contract.
func (h *Handler) Terminate(ctx context.Context, id string, req dto.TerminateRequest) (*dto.ContractResponse, error) {
	now := h.clock.Now()
	var c *contract.Contract

	dbErr := h.uow.Do(ctx, func(repos unitofwork.Repositories) error {
		var err error
		c, err = repos.Contracts.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if err := c.Terminate(req.Reason, now); err != nil {
			return apperror.NewConflict(err.Error())
		}
		return repos.Contracts.Update(ctx, c)
	})

	if dbErr != nil {
		return nil, dbErr
	}

	h.publisher.Publish(c.Events()...)
	return toResponse(c), nil
}

// toResponse maps a contract aggregate to its DTO representation.
func toResponse(c *contract.Contract) *dto.ContractResponse {
	clauses := make([]dto.ClauseDTO, 0, len(c.Clauses()))
	for _, cl := range c.Clauses() {
		clauses = append(clauses, dto.ClauseDTO{
			Title:   cl.Title,
			Content: cl.Content,
			Order:   cl.Order,
		})
	}

	resp := &dto.ContractResponse{
		ID:          c.ID(),
		Title:       c.Title(),
		Description: c.Description(),
		Status:      string(c.Status()),
		Value: dto.MoneyDTO{
			AmountCents: c.Value().AmountCents,
			Currency:    c.Value().Currency,
		},
		Clauses:   clauses,
		StartDate: c.DateRange().Start,
		EndDate:   c.DateRange().End,
		OwnerID:   c.OwnerID(),
		PartyID:   c.PartyID(),
		TemplateID: c.TemplateID(),
		Version:   c.Version(),
		CreatedAt: c.CreatedAt(),
		UpdatedAt: c.UpdatedAt(),
	}
	return resp
}