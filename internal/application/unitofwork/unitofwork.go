package unitofwork

import (
	"context"

	"github.com/contractiq/contractiq/internal/domain/approval"
	"github.com/contractiq/contractiq/internal/domain/contract"
	"github.com/contractiq/contractiq/internal/domain/party"
	"github.com/contractiq/contractiq/internal/domain/template"
)

// UnitOfWork provides a transactional boundary for coordinating multiple repositories.
type UnitOfWork interface {
	// Do executes the given function within a database transaction.
	// If the function returns an error, the transaction is rolled back.
	Do(ctx context.Context, fn func(repos Repositories) error) error
}

// Repositories provides access to all domain repositories within a transaction.
type Repositories struct {
	Contracts contract.Repository
	Templates template.Repository
	Parties   party.Repository
	Approvals approval.Repository
}
