package contract

import "fmt"

// Status represents the lifecycle state of a contract.
type Status string

const (
	StatusDraft         Status = "DRAFT"
	StatusPendingReview Status = "PENDING_REVIEW"
	StatusApproved      Status = "APPROVED"
	StatusActive        Status = "ACTIVE"
	StatusExpired       Status = "EXPIRED"
	StatusTerminated    Status = "TERMINATED"
)

// AllStatuses returns all valid contract statuses.
func AllStatuses() []Status {
	return []Status{
		StatusDraft,
		StatusPendingReview,
		StatusApproved,
		StatusActive,
		StatusExpired,
		StatusTerminated,
	}
}

func (s Status) String() string { return string(s) }

// IsValid checks if the status is a known value.
func (s Status) IsValid() bool {
	for _, v := range AllStatuses() {
		if v == s {
			return true
		}
	}
	return false
}

// IsTerminal returns true if the contract cannot transition further.
func (s Status) IsTerminal() bool {
	return s == StatusExpired || s == StatusTerminated
}

// transitions defines the valid finite state machine for contract status.
var transitions = map[Status][]Status{
	StatusDraft:         {StatusPendingReview},
	StatusPendingReview: {StatusApproved, StatusDraft},
	StatusApproved:      {StatusActive},
	StatusActive:        {StatusExpired, StatusTerminated},
}

// CanTransitionTo checks if the given transition is allowed.
func (s Status) CanTransitionTo(target Status) bool {
	allowed, ok := transitions[s]
	if !ok {
		return false
	}
	for _, a := range allowed {
		if a == target {
			return true
		}
	}
	return false
}

// TransitionTo attempts a status transition and returns an error if invalid.
func (s Status) TransitionTo(target Status) (Status, error) {
	if !s.CanTransitionTo(target) {
		return s, fmt.Errorf("invalid status transition from %s to %s", s, target)
	}
	return target, nil
}
