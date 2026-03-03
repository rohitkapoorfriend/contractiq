package contract

import "context"

// Repository defines the persistence interface for the Contract aggregate.
// Implementations live in the infrastructure layer.
type Repository interface {
	// FindByID retrieves a contract by its unique identifier.
	FindByID(ctx context.Context, id string) (*Contract, error)

	// FindAll retrieves contracts matching the given filter criteria.
	FindAll(ctx context.Context, filter Filter) (*PageResult, error)

	// Save persists a new contract.
	Save(ctx context.Context, contract *Contract) error

	// Update persists changes to an existing contract with optimistic concurrency.
	Update(ctx context.Context, contract *Contract) error

	// Delete removes a contract. Only draft contracts may be deleted.
	Delete(ctx context.Context, id string) error
}
