package party

import "context"

// Repository defines the persistence interface for the Party aggregate.
type Repository interface {
	FindByID(ctx context.Context, id string) (*Party, error)
	FindAll(ctx context.Context, createdBy string) ([]*Party, error)
	Save(ctx context.Context, party *Party) error
	Update(ctx context.Context, party *Party) error
	Delete(ctx context.Context, id string) error
}
