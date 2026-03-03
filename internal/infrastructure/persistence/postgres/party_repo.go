package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/contractiq/contractiq/internal/domain/party"
	"github.com/contractiq/contractiq/pkg/apperror"
	"github.com/jmoiron/sqlx"
)

// PartyRepository implements party.Repository using PostgreSQL.
type PartyRepository struct {
	db sqlx.ExtContext
}

// NewPartyRepository creates a new PostgreSQL party repository.
func NewPartyRepository(db sqlx.ExtContext) *PartyRepository {
	return &PartyRepository{db: db}
}

type partyRow struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	PartyType string    `db:"party_type"`
	Company   string    `db:"company"`
	Phone     string    `db:"phone"`
	Address   string    `db:"address"`
	CreatedBy string    `db:"created_by"`
	Version   int       `db:"version"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r *PartyRepository) FindByID(ctx context.Context, id string) (*party.Party, error) {
	var row partyRow
	query := `SELECT id, name, email, party_type, company, phone, address, created_by, version, created_at, updated_at
	          FROM parties WHERE id = $1`

	err := sqlx.GetContext(ctx, r.db, &row, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("party", id)
		}
		return nil, fmt.Errorf("failed to find party: %w", err)
	}

	return rowToParty(row), nil
}

func (r *PartyRepository) FindAll(ctx context.Context, createdBy string) ([]*party.Party, error) {
	query := `SELECT id, name, email, party_type, company, phone, address, created_by, version, created_at, updated_at
	          FROM parties WHERE created_by = $1 ORDER BY created_at DESC`

	var rows []partyRow
	err := sqlx.SelectContext(ctx, r.db, &rows, query, createdBy)
	if err != nil {
		return nil, fmt.Errorf("failed to list parties: %w", err)
	}

	result := make([]*party.Party, 0, len(rows))
	for _, row := range rows {
		result = append(result, rowToParty(row))
	}
	return result, nil
}

func (r *PartyRepository) Save(ctx context.Context, p *party.Party) error {
	query := `INSERT INTO parties (id, name, email, party_type, company, phone, address, created_by, version, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.ExecContext(ctx, query,
		p.ID(), p.Name(), p.Email(), string(p.Type()),
		p.Company(), p.Phone(), p.Address(), p.CreatedBy(),
		p.Version(), p.CreatedAt(), p.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to save party: %w", err)
	}
	return nil
}

func (r *PartyRepository) Update(ctx context.Context, p *party.Party) error {
	query := `UPDATE parties SET name = $1, email = $2, company = $3, phone = $4, address = $5,
	          version = version + 1, updated_at = $6
	          WHERE id = $7 AND version = $8`

	result, err := r.db.ExecContext(ctx, query,
		p.Name(), p.Email(), p.Company(), p.Phone(), p.Address(),
		p.UpdatedAt(), p.ID(), p.Version(),
	)
	if err != nil {
		return fmt.Errorf("failed to update party: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return apperror.NewConcurrencyConflict("party")
	}
	return nil
}

func (r *PartyRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM parties WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete party: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return apperror.NewNotFound("party", id)
	}
	return nil
}

func rowToParty(row partyRow) *party.Party {
	return party.Reconstitute(
		row.ID, row.Name, row.Email,
		party.PartyType(row.PartyType),
		row.Company, row.Phone, row.Address, row.CreatedBy,
		row.Version, row.CreatedAt, row.UpdatedAt,
	)
}
