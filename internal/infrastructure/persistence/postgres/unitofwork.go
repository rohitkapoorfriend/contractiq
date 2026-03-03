package postgres

import (
	"context"
	"fmt"

	"github.com/contractiq/contractiq/internal/application/unitofwork"
	"github.com/jmoiron/sqlx"
)

// UnitOfWork implements unitofwork.UnitOfWork using PostgreSQL transactions.
type UnitOfWork struct {
	db *sqlx.DB
}

// NewUnitOfWork creates a new PostgreSQL unit of work.
func NewUnitOfWork(db *sqlx.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

// Do executes the function within a database transaction.
func (u *UnitOfWork) Do(ctx context.Context, fn func(repos unitofwork.Repositories) error) error {
	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	repos := unitofwork.Repositories{
		Contracts: NewContractRepository(tx),
		Templates: NewTemplateRepository(tx),
		Parties:   NewPartyRepository(tx),
		Approvals: NewApprovalRepository(tx),
	}

	if err := fn(repos); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback failed: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
