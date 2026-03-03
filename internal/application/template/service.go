package template

import (
	"context"
	"time"

	"github.com/contractiq/contractiq/internal/domain/contract"
	"github.com/contractiq/contractiq/internal/domain/template"
	"github.com/contractiq/contractiq/pkg/apperror"
	"github.com/contractiq/contractiq/pkg/clock"
)

// CreateRequest is the input for creating a template.
type CreateRequest struct {
	Name        string      `json:"name" validate:"required,min=3,max=200"`
	Description string      `json:"description" validate:"max=2000"`
	Clauses     []ClauseDTO `json:"clauses" validate:"dive"`
}

// UpdateRequest is the input for updating a template.
type UpdateRequest struct {
	Name        string      `json:"name" validate:"required,min=3,max=200"`
	Description string      `json:"description" validate:"max=2000"`
	Clauses     []ClauseDTO `json:"clauses" validate:"dive"`
	Version     int         `json:"version" validate:"required,min=1"`
}

// ClauseDTO is a clause for templates.
type ClauseDTO struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	Order   int    `json:"order" validate:"min=0"`
}

// Response is the output representation of a template.
type Response struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Clauses     []ClauseDTO `json:"clauses"`
	IsActive    bool        `json:"is_active"`
	CreatedBy   string      `json:"created_by"`
	Version     int         `json:"version"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// Service handles template use cases.
type Service struct {
	repo  template.Repository
	clock clock.Clock
}

// NewService creates a new template service.
func NewService(repo template.Repository, clock clock.Clock) *Service {
	return &Service{repo: repo, clock: clock}
}

// Create creates a new template.
func (s *Service) Create(ctx context.Context, userID string, req CreateRequest) (*Response, error) {
	now := s.clock.Now()

	clauses := make([]contract.Clause, 0, len(req.Clauses))
	for _, c := range req.Clauses {
		clause, err := contract.NewClause(c.Title, c.Content, c.Order)
		if err != nil {
			return nil, apperror.NewValidation(err.Error())
		}
		clauses = append(clauses, clause)
	}

	t, err := template.NewTemplate(req.Name, req.Description, userID, clauses, now)
	if err != nil {
		return nil, apperror.NewValidation(err.Error())
	}

	if err := s.repo.Save(ctx, t); err != nil {
		return nil, err
	}

	return toResponse(t), nil
}

// GetByID retrieves a template.
func (s *Service) GetByID(ctx context.Context, id string) (*Response, error) {
	t, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toResponse(t), nil
}

// List retrieves all templates.
func (s *Service) List(ctx context.Context, activeOnly bool) ([]Response, error) {
	templates, err := s.repo.FindAll(ctx, activeOnly)
	if err != nil {
		return nil, err
	}
	result := make([]Response, 0, len(templates))
	for _, t := range templates {
		result = append(result, *toResponse(t))
	}
	return result, nil
}

// Update modifies a template.
func (s *Service) Update(ctx context.Context, id string, req UpdateRequest) (*Response, error) {
	now := s.clock.Now()

	t, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if t.Version() != req.Version {
		return nil, apperror.NewConcurrencyConflict("template")
	}

	clauses := make([]contract.Clause, 0, len(req.Clauses))
	for _, c := range req.Clauses {
		clause, err := contract.NewClause(c.Title, c.Content, c.Order)
		if err != nil {
			return nil, apperror.NewValidation(err.Error())
		}
		clauses = append(clauses, clause)
	}

	if err := t.Update(req.Name, req.Description, clauses, now); err != nil {
		return nil, apperror.NewValidation(err.Error())
	}

	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}

	return toResponse(t), nil
}

// Delete removes a template.
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func toResponse(t *template.Template) *Response {
	clauses := make([]ClauseDTO, 0, len(t.Clauses()))
	for _, c := range t.Clauses() {
		clauses = append(clauses, ClauseDTO{
			Title:   c.Title,
			Content: c.Content,
			Order:   c.Order,
		})
	}
	return &Response{
		ID:          t.ID(),
		Name:        t.Name(),
		Description: t.Description(),
		Clauses:     clauses,
		IsActive:    t.IsActive(),
		CreatedBy:   t.CreatedBy(),
		Version:     t.Version(),
		CreatedAt:   t.CreatedAt(),
		UpdatedAt:   t.UpdatedAt(),
	}
}
