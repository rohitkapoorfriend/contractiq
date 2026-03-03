package approval

import (
	"fmt"
	"time"

	"github.com/contractiq/contractiq/pkg/identifier"
)

// Decision represents the outcome of an approval review.
type Decision string

const (
	DecisionPending  Decision = "PENDING"
	DecisionApproved Decision = "APPROVED"
	DecisionRejected Decision = "REJECTED"
)

// Approval represents an approval request for a contract.
type Approval struct {
	id         string
	contractID string
	reviewerID string
	decision   Decision
	comment    string
	decidedAt  *time.Time
	createdAt  time.Time
}

func (a *Approval) ID() string          { return a.id }
func (a *Approval) ContractID() string   { return a.contractID }
func (a *Approval) ReviewerID() string   { return a.reviewerID }
func (a *Approval) Decision() Decision   { return a.decision }
func (a *Approval) Comment() string      { return a.comment }
func (a *Approval) DecidedAt() *time.Time { return a.decidedAt }
func (a *Approval) CreatedAt() time.Time  { return a.createdAt }

// NewApproval creates a new pending approval request.
func NewApproval(contractID, reviewerID string, now time.Time) (*Approval, error) {
	if contractID == "" {
		return nil, fmt.Errorf("contract ID is required")
	}
	if reviewerID == "" {
		return nil, fmt.Errorf("reviewer ID is required")
	}
	return &Approval{
		id:         identifier.New(),
		contractID: contractID,
		reviewerID: reviewerID,
		decision:   DecisionPending,
		createdAt:  now,
	}, nil
}

// Approve records an approval decision.
func (a *Approval) Approve(comment string, now time.Time) error {
	if a.decision != DecisionPending {
		return fmt.Errorf("approval has already been decided")
	}
	a.decision = DecisionApproved
	a.comment = comment
	a.decidedAt = &now
	return nil
}

// Reject records a rejection decision.
func (a *Approval) Reject(comment string, now time.Time) error {
	if a.decision != DecisionPending {
		return fmt.Errorf("approval has already been decided")
	}
	if comment == "" {
		return fmt.Errorf("rejection comment is required")
	}
	a.decision = DecisionRejected
	a.comment = comment
	a.decidedAt = &now
	return nil
}

// IsPending checks if the approval is still awaiting a decision.
func (a *Approval) IsPending() bool {
	return a.decision == DecisionPending
}

// Reconstitute creates an Approval from persisted data.
func Reconstitute(
	id, contractID, reviewerID string,
	decision Decision,
	comment string,
	decidedAt *time.Time,
	createdAt time.Time,
) *Approval {
	return &Approval{
		id:         id,
		contractID: contractID,
		reviewerID: reviewerID,
		decision:   decision,
		comment:    comment,
		decidedAt:  decidedAt,
		createdAt:  createdAt,
	}
}
