package party

import (
	"context"
	"time"

	"github.com/contractiq/contractiq/internal/domain/party"
	"github.com/contractiq/contractiq/pkg/apperror"
	"github.com/contractiq/contractiq/pkg/clock"
)

// CreateRequest is the input for creating a party.
type CreateRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=200"`
	Email   string `json:"email" validate:"required,email"`
	Type    string `json:"type" validate:"required,oneof=ORGANIZATION INDIVIDUAL"`
	Company string `json:"company" validate:"max=200"`
	Phone   string `json:"phone" validate:"max=20"`
	Address string `json:"address" validate:"max=500"`
}

// UpdateRequest is the input for updating a party.
type UpdateRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=200"`
	Email   string `json:"email" validate:"required,email"`
	Company string `json:"company" validate:"max=200"`
	Phone   string `json:"phone" validate:"max=20"`
	Address string `json:"address" validate:"max=500"`
	Version int    `json:"version" validate:"required,min=1"`
}

// Response is the output representation of a party.
type Response struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Type      string    `json:"type"`
	Company   string    `json:"company"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedBy string    `json:"created_by"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Service handles party use cases.
type Service struct {
	repo  party.Repository
	clock clock.Clock
}

// NewService creates a new party service.
func NewService(repo party.Repository, clock clock.Clock) *Service {
	return &Service{repo: repo, clock: clock}
}

// Create creates a new party.
func (s *Service) Create(ctx context.Context, userID string, req CreateRequest) (*Response, error) {
	now := s.clock.Now()

	p, err := party.NewParty(req.Name, req.Email, party.PartyType(req.Type), req.Company, req.Phone, req.Address, userID, now)
	if err != nil {
		return nil, apperror.NewValidation(err.Error())
	}

	if err := s.repo.Save(ctx, p); err != nil {
		return nil, err
	}

	return toResponse(p), nil
}

// GetByID retrieves a party.
func (s *Service) GetByID(ctx context.Context, id string) (*Response, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toResponse(p), nil
}

// List retrieves all parties for a user.
func (s *Service) List(ctx context.Context, userID string) ([]Response, error) {
	parties, err := s.repo.FindAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]Response, 0, len(parties))
	for _, p := range parties {
		result = append(result, *toResponse(p))
	}
	return result, nil
}

// Update modifies a party.
func (s *Service) Update(ctx context.Context, id string, req UpdateRequest) (*Response, error) {
	now := s.clock.Now()

	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if p.Version() != req.Version {
		return nil, apperror.NewConcurrencyConflict("party")
	}

	if err := p.Update(req.Name, req.Email, req.Company, req.Phone, req.Address, now); err != nil {
		return nil, apperror.NewValidation(err.Error())
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}

	return toResponse(p), nil
}

// Delete removes a party.
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func toResponse(p *party.Party) *Response {
	return &Response{
		ID:        p.ID(),
		Name:      p.Name(),
		Email:     p.Email(),
		Type:      string(p.Type()),
		Company:   p.Company(),
		Phone:     p.Phone(),
		Address:   p.Address(),
		CreatedBy: p.CreatedBy(),
		Version:   p.Version(),
		CreatedAt: p.CreatedAt(),
		UpdatedAt: p.UpdatedAt(),
	}
}
