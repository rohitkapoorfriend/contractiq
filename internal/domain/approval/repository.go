package approval

import "context"

// Repository defines the persistence interface for the Approval aggregate.
type Repository interface {
	FindByID(ctx context.Context, id string) (*Approval, error)
	FindByContractID(ctx context.Context, contractID string) ([]*Approval, error)
	Save(ctx context.Context, approval *Approval) error
	Update(ctx context.Context, approval *Approval) error
}
