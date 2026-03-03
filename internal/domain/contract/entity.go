package contract

import (
	"time"

	"github.com/contractiq/contractiq/internal/domain/event"
)

// Contract is the core aggregate root of the system.
type Contract struct {
	id          string
	title       string
	description string
	status      Status
	value       Money
	clauses     []Clause
	dateRange   DateRange
	ownerID     string
	partyID     *string
	templateID  *string
	version     int
	createdAt   time.Time
	updatedAt   time.Time

	events []event.Event
}

// --- Getters (exported, read-only access to unexported fields) ---

func (c *Contract) ID() string          { return c.id }
func (c *Contract) Title() string        { return c.title }
func (c *Contract) Description() string  { return c.description }
func (c *Contract) Status() Status       { return c.status }
func (c *Contract) Value() Money         { return c.value }
func (c *Contract) Clauses() []Clause    { return c.clauses }
func (c *Contract) DateRange() DateRange { return c.dateRange }
func (c *Contract) OwnerID() string      { return c.ownerID }
func (c *Contract) PartyID() *string     { return c.partyID }
func (c *Contract) TemplateID() *string  { return c.templateID }
func (c *Contract) Version() int         { return c.version }
func (c *Contract) CreatedAt() time.Time { return c.createdAt }
func (c *Contract) UpdatedAt() time.Time { return c.updatedAt }

// Events returns and clears the uncommitted domain events.
func (c *Contract) Events() []event.Event {
	events := c.events
	c.events = nil
	return events
}

func (c *Contract) recordEvent(e event.Event) {
	c.events = append(c.events, e)
}

// --- Behavior (state transitions via methods) ---

// Update modifies mutable fields of a draft contract.
func (c *Contract) Update(title, description string, value Money, dateRange DateRange, now time.Time) error {
	if c.status != StatusDraft {
		return &InvalidOperationError{
			Operation: "update",
			Status:    c.status,
		}
	}
	c.title = title
	c.description = description
	c.value = value
	c.dateRange = dateRange
	c.updatedAt = now
	return nil
}

// SetClauses replaces all clauses on a draft contract.
func (c *Contract) SetClauses(clauses []Clause, now time.Time) error {
	if c.status != StatusDraft {
		return &InvalidOperationError{Operation: "set clauses", Status: c.status}
	}
	c.clauses = clauses
	c.updatedAt = now
	return nil
}

// Submit moves the contract from Draft to PendingReview.
func (c *Contract) Submit(now time.Time) error {
	newStatus, err := c.status.TransitionTo(StatusPendingReview)
	if err != nil {
		return &InvalidOperationError{Operation: "submit", Status: c.status}
	}
	c.status = newStatus
	c.updatedAt = now
	c.recordEvent(NewContractSubmittedEvent(c.id, now))
	return nil
}

// Approve moves the contract from PendingReview to Approved.
func (c *Contract) Approve(approvedBy string, now time.Time) error {
	newStatus, err := c.status.TransitionTo(StatusApproved)
	if err != nil {
		return &InvalidOperationError{Operation: "approve", Status: c.status}
	}
	c.status = newStatus
	c.updatedAt = now
	c.recordEvent(NewContractApprovedEvent(c.id, approvedBy, now))
	return nil
}

// Sign activates the contract after approval.
func (c *Contract) Sign(signedBy string, now time.Time) error {
	newStatus, err := c.status.TransitionTo(StatusActive)
	if err != nil {
		return &InvalidOperationError{Operation: "sign", Status: c.status}
	}
	c.status = newStatus
	c.updatedAt = now
	c.recordEvent(NewContractSignedEvent(c.id, signedBy, now))
	return nil
}

// Terminate ends an active contract with a reason.
func (c *Contract) Terminate(reason string, now time.Time) error {
	newStatus, err := c.status.TransitionTo(StatusTerminated)
	if err != nil {
		return &InvalidOperationError{Operation: "terminate", Status: c.status}
	}
	if reason == "" {
		reason = "no reason provided"
	}
	c.status = newStatus
	c.updatedAt = now
	c.recordEvent(NewContractTerminatedEvent(c.id, reason, now))
	return nil
}

// Expire marks an active contract as expired.
func (c *Contract) Expire(now time.Time) error {
	newStatus, err := c.status.TransitionTo(StatusExpired)
	if err != nil {
		return &InvalidOperationError{Operation: "expire", Status: c.status}
	}
	c.status = newStatus
	c.updatedAt = now
	return nil
}

// InvalidOperationError indicates an attempt to perform an invalid operation for the current status.
type InvalidOperationError struct {
	Operation string
	Status    Status
}

func (e *InvalidOperationError) Error() string {
	return "cannot " + e.Operation + " contract in status " + string(e.Status)
}

// Reconstitute creates a Contract from persisted data (bypasses validation).
// Used by repository implementations to hydrate from database.
func Reconstitute(
	id, title, description string,
	status Status,
	value Money,
	clauses []Clause,
	dateRange DateRange,
	ownerID string,
	partyID, templateID *string,
	version int,
	createdAt, updatedAt time.Time,
) *Contract {
	return &Contract{
		id:          id,
		title:       title,
		description: description,
		status:      status,
		value:       value,
		clauses:     clauses,
		dateRange:   dateRange,
		ownerID:     ownerID,
		partyID:     partyID,
		templateID:  templateID,
		version:     version,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}
