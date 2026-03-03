package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/contractiq/contractiq/internal/domain/contract"
	"github.com/contractiq/contractiq/pkg/apperror"
	"github.com/jmoiron/sqlx"
)

// ContractRepository implements contract.Repository using PostgreSQL.
type ContractRepository struct {
	db sqlx.ExtContext
}

// NewContractRepository creates a new PostgreSQL contract repository.
func NewContractRepository(db sqlx.ExtContext) *ContractRepository {
	return &ContractRepository{db: db}
}

type contractRow struct {
	ID          string         `db:"id"`
	Title       string         `db:"title"`
	Description string         `db:"description"`
	Status      string         `db:"status"`
	AmountCents int64          `db:"amount_cents"`
	Currency    string         `db:"currency"`
	Clauses     []byte         `db:"clauses"`
	StartDate   time.Time      `db:"start_date"`
	EndDate     time.Time      `db:"end_date"`
	OwnerID     string         `db:"owner_id"`
	PartyID     sql.NullString `db:"party_id"`
	TemplateID  sql.NullString `db:"template_id"`
	Version     int            `db:"version"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

func (r *ContractRepository) FindByID(ctx context.Context, id string) (*contract.Contract, error) {
	var row contractRow
	query := `SELECT id, title, description, status, amount_cents, currency, clauses,
	          start_date, end_date, owner_id, party_id, template_id, version, created_at, updated_at
	          FROM contracts WHERE id = $1`

	err := sqlx.GetContext(ctx, r.db, &row, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("contract", id)
		}
		return nil, fmt.Errorf("failed to find contract: %w", err)
	}

	return rowToContract(row)
}

func (r *ContractRepository) FindAll(ctx context.Context, filter contract.Filter) (*contract.PageResult, error) {
	where := "WHERE 1=1"
	args := make([]interface{}, 0)
	argIdx := 1

	if filter.OwnerID != nil {
		where += fmt.Sprintf(" AND owner_id = $%d", argIdx)
		args = append(args, *filter.OwnerID)
		argIdx++
	}
	if filter.Status != nil {
		where += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, string(*filter.Status))
		argIdx++
	}
	if filter.PartyID != nil {
		where += fmt.Sprintf(" AND party_id = $%d", argIdx)
		args = append(args, *filter.PartyID)
		argIdx++
	}
	if filter.Search != nil {
		where += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+*filter.Search+"%")
		argIdx++
	}

	// Count total
	var totalCount int
	countQuery := "SELECT COUNT(*) FROM contracts " + where
	err := sqlx.GetContext(ctx, r.db, &totalCount, countQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to count contracts: %w", err)
	}

	// Fetch page
	query := fmt.Sprintf(`SELECT id, title, description, status, amount_cents, currency, clauses,
	          start_date, end_date, owner_id, party_id, template_id, version, created_at, updated_at
	          FROM contracts %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		where, argIdx, argIdx+1)
	args = append(args, filter.PageSize, filter.Offset())

	var rows []contractRow
	err = sqlx.SelectContext(ctx, r.db, &rows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list contracts: %w", err)
	}

	items := make([]*contract.Contract, 0, len(rows))
	for _, row := range rows {
		c, err := rowToContract(row)
		if err != nil {
			return nil, err
		}
		items = append(items, c)
	}

	return &contract.PageResult{
		Items:      items,
		TotalCount: totalCount,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
	}, nil
}

func (r *ContractRepository) Save(ctx context.Context, c *contract.Contract) error {
	clausesJSON, err := json.Marshal(c.Clauses())
	if err != nil {
		return fmt.Errorf("failed to marshal clauses: %w", err)
	}

	query := `INSERT INTO contracts (id, title, description, status, amount_cents, currency, clauses,
	          start_date, end_date, owner_id, party_id, template_id, version, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	_, err = r.db.ExecContext(ctx, query,
		c.ID(), c.Title(), c.Description(), string(c.Status()),
		c.Value().AmountCents, c.Value().Currency, clausesJSON,
		c.DateRange().Start, c.DateRange().End,
		c.OwnerID(), nullString(c.PartyID()), nullString(c.TemplateID()),
		c.Version(), c.CreatedAt(), c.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to save contract: %w", err)
	}
	return nil
}

func (r *ContractRepository) Update(ctx context.Context, c *contract.Contract) error {
	clausesJSON, err := json.Marshal(c.Clauses())
	if err != nil {
		return fmt.Errorf("failed to marshal clauses: %w", err)
	}

	query := `UPDATE contracts SET title = $1, description = $2, status = $3,
	          amount_cents = $4, currency = $5, clauses = $6,
	          start_date = $7, end_date = $8, party_id = $9, template_id = $10,
	          version = version + 1, updated_at = $11
	          WHERE id = $12 AND version = $13`

	result, err := r.db.ExecContext(ctx, query,
		c.Title(), c.Description(), string(c.Status()),
		c.Value().AmountCents, c.Value().Currency, clausesJSON,
		c.DateRange().Start, c.DateRange().End,
		nullString(c.PartyID()), nullString(c.TemplateID()),
		c.UpdatedAt(), c.ID(), c.Version(),
	)
	if err != nil {
		return fmt.Errorf("failed to update contract: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return apperror.NewConcurrencyConflict("contract")
	}
	return nil
}

func (r *ContractRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM contracts WHERE id = $1 AND status = 'DRAFT'`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete contract: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return apperror.NewNotFound("contract", id)
	}
	return nil
}

func rowToContract(row contractRow) (*contract.Contract, error) {
	var clauses []contract.Clause
	if row.Clauses != nil {
		if err := json.Unmarshal(row.Clauses, &clauses); err != nil {
			return nil, fmt.Errorf("failed to unmarshal clauses: %w", err)
		}
	}

	var partyID, templateID *string
	if row.PartyID.Valid {
		partyID = &row.PartyID.String
	}
	if row.TemplateID.Valid {
		templateID = &row.TemplateID.String
	}

	return contract.Reconstitute(
		row.ID, row.Title, row.Description,
		contract.Status(row.Status),
		contract.Money{AmountCents: row.AmountCents, Currency: row.Currency},
		clauses,
		contract.DateRange{Start: row.StartDate, End: row.EndDate},
		row.OwnerID, partyID, templateID,
		row.Version, row.CreatedAt, row.UpdatedAt,
	), nil
}

func nullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}
