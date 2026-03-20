package template

import "context"

// Repository defines the persistence interface for the Template aggregate.
type Repository interface {
	FindByID(ctx context.Context, id string) (*Template, error)
	FindAll(ctx context.Context, activeOnly bool) ([]*Template, error)
	Save(ctx context.Context, template *Template) error
	Update(ctx context.Context, template *Template) error
	Delete(ctx context.Context, id string) error
}