package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/contractiq/contractiq/internal/domain/approval"
	"github.com/contractiq/contractiq/pkg/apperror"
	"github.com/jmoiron/sqlx"
)

// ApprovalRepository implements approval.Repository using PostgreSQL.
type ApprovalRepository struct {
	db sqlx.ExtContext
}

// NewApprovalRepository creates a new PostgreSQL approval repository.
func NewApprovalRepository(db sqlx.ExtContext) *ApprovalRepository {
	return &ApprovalRepository{db: db}
}

type approvalRow struct {
	ID         string     `db:"id"`
	ContractID string     `db:"contract_id"`
	ReviewerID string     `db:"reviewer_id"`
	Decision   string     `db:"decision"`
	Comment    string     `db:"comment"`
	DecidedAt  *time.Time `db:"decided_at"`
	CreatedAt  time.Time  `db:"created_at"`
}

func (r *ApprovalRepository) FindByID(ctx context.Context, id string) (*approval.Approval, error) {
	var row approvalRow
	query := `SELECT id, contract_id, reviewer_id, decision, comment, decided_at, created_at
	          FROM approvals WHERE id = $1`

	err := sqlx.GetContext(ctx, r.db, &row, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("approval", id)
		}
		return nil, fmt.Errorf("failed to find approval: %w", err)
	}

	return rowToApproval(row), nil
}

func (r *ApprovalRepository) FindByContractID(ctx context.Context, contractID string) ([]*approval.Approval, error) {
	query := `SELECT id, contract_id, reviewer_id, decision, comment, decided_at, created_at
	          FROM approvals WHERE contract_id = $1 ORDER BY created_at DESC`

	var rows []approvalRow
	err := sqlx.SelectContext(ctx, r.db, &rows, query, contractID)
	if err != nil {
		return nil, fmt.Errorf("failed to list approvals: %w", err)
	}

	result := make([]*approval.Approval, 0, len(rows))
	for _, row := range rows {
		result = append(result, rowToApproval(row))
	}
	return result, nil
}

func (r *ApprovalRepository) Save(ctx context.Context, a *approval.Approval) error {
	query := `INSERT INTO approvals (id, contract_id, reviewer_id, decision, comment, decided_at, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		a.ID(), a.ContractID(), a.ReviewerID(), string(a.Decision()),
		a.Comment(), a.DecidedAt(), a.CreatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to save approval: %w", err)
	}
	return nil
}

func (r *ApprovalRepository) Update(ctx context.Context, a *approval.Approval) error {
	query := `UPDATE approvals SET decision = $1, comment = $2, decided_at = $3
	          WHERE id = $4`

	_, err := r.db.ExecContext(ctx, query,
		string(a.Decision()), a.Comment(), a.DecidedAt(), a.ID(),
	)
	if err != nil {
		return fmt.Errorf("failed to update approval: %w", err)
	}
	return nil
}

func rowToApproval(row approvalRow) *approval.Approval {
	return approval.Reconstitute(
		row.ID, row.ContractID, row.ReviewerID,
		approval.Decision(row.Decision),
		row.Comment, row.DecidedAt, row.CreatedAt,
	)
}
