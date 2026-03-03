package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	domaincontract "github.com/contractiq/contractiq/internal/domain/contract"
	"github.com/contractiq/contractiq/internal/domain/template"
	"github.com/contractiq/contractiq/pkg/apperror"
	"github.com/jmoiron/sqlx"
)

// TemplateRepository implements template.Repository using PostgreSQL.
type TemplateRepository struct {
	db sqlx.ExtContext
}

// NewTemplateRepository creates a new PostgreSQL template repository.
func NewTemplateRepository(db sqlx.ExtContext) *TemplateRepository {
	return &TemplateRepository{db: db}
}

type templateRow struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Clauses     []byte    `db:"clauses"`
	IsActive    bool      `db:"is_active"`
	CreatedBy   string    `db:"created_by"`
	Version     int       `db:"version"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (r *TemplateRepository) FindByID(ctx context.Context, id string) (*template.Template, error) {
	var row templateRow
	query := `SELECT id, name, description, clauses, is_active, created_by, version, created_at, updated_at
	          FROM templates WHERE id = $1`

	err := sqlx.GetContext(ctx, r.db, &row, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("template", id)
		}
		return nil, fmt.Errorf("failed to find template: %w", err)
	}

	return rowToTemplate(row)
}

func (r *TemplateRepository) FindAll(ctx context.Context, activeOnly bool) ([]*template.Template, error) {
	query := `SELECT id, name, description, clauses, is_active, created_by, version, created_at, updated_at
	          FROM templates`
	if activeOnly {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY created_at DESC"

	var rows []templateRow
	err := sqlx.SelectContext(ctx, r.db, &rows, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}

	result := make([]*template.Template, 0, len(rows))
	for _, row := range rows {
		t, err := rowToTemplate(row)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (r *TemplateRepository) Save(ctx context.Context, t *template.Template) error {
	clausesJSON, err := json.Marshal(t.Clauses())
	if err != nil {
		return fmt.Errorf("failed to marshal clauses: %w", err)
	}

	query := `INSERT INTO templates (id, name, description, clauses, is_active, created_by, version, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err = r.db.ExecContext(ctx, query,
		t.ID(), t.Name(), t.Description(), clausesJSON,
		t.IsActive(), t.CreatedBy(), t.Version(), t.CreatedAt(), t.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to save template: %w", err)
	}
	return nil
}

func (r *TemplateRepository) Update(ctx context.Context, t *template.Template) error {
	clausesJSON, err := json.Marshal(t.Clauses())
	if err != nil {
		return fmt.Errorf("failed to marshal clauses: %w", err)
	}

	query := `UPDATE templates SET name = $1, description = $2, clauses = $3, is_active = $4,
	          version = version + 1, updated_at = $5
	          WHERE id = $6 AND version = $7`

	result, err := r.db.ExecContext(ctx, query,
		t.Name(), t.Description(), clausesJSON, t.IsActive(),
		t.UpdatedAt(), t.ID(), t.Version(),
	)
	if err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return apperror.NewConcurrencyConflict("template")
	}
	return nil
}

func (r *TemplateRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM templates WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return apperror.NewNotFound("template", id)
	}
	return nil
}

func rowToTemplate(row templateRow) (*template.Template, error) {
	var clauses []domaincontract.Clause
	if row.Clauses != nil {
		if err := json.Unmarshal(row.Clauses, &clauses); err != nil {
			return nil, fmt.Errorf("failed to unmarshal clauses: %w", err)
		}
	}

	return template.Reconstitute(
		row.ID, row.Name, row.Description,
		clauses, row.IsActive, row.CreatedBy,
		row.Version, row.CreatedAt, row.UpdatedAt,
	), nil
}
