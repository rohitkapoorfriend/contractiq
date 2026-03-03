package contract_test

import (
	"testing"

	"github.com/contractiq/contractiq/internal/domain/contract"
	"github.com/stretchr/testify/assert"
)

func TestStatusFSM(t *testing.T) {
	tests := []struct {
		name    string
		from    contract.Status
		to      contract.Status
		allowed bool
	}{
		{"draft to pending review", contract.StatusDraft, contract.StatusPendingReview, true},
		{"pending to approved", contract.StatusPendingReview, contract.StatusApproved, true},
		{"pending to draft (reject)", contract.StatusPendingReview, contract.StatusDraft, true},
		{"approved to active", contract.StatusApproved, contract.StatusActive, true},
		{"active to expired", contract.StatusActive, contract.StatusExpired, true},
		{"active to terminated", contract.StatusActive, contract.StatusTerminated, true},
		{"draft to active (skip)", contract.StatusDraft, contract.StatusActive, false},
		{"draft to terminated", contract.StatusDraft, contract.StatusTerminated, false},
		{"expired to active", contract.StatusExpired, contract.StatusActive, false},
		{"terminated to draft", contract.StatusTerminated, contract.StatusDraft, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.allowed, tt.from.CanTransitionTo(tt.to))
		})
	}
}

func TestStatusIsValid(t *testing.T) {
	assert.True(t, contract.StatusDraft.IsValid())
	assert.True(t, contract.StatusActive.IsValid())
	assert.False(t, contract.Status("INVALID").IsValid())
}

func TestStatusIsTerminal(t *testing.T) {
	assert.True(t, contract.StatusExpired.IsTerminal())
	assert.True(t, contract.StatusTerminated.IsTerminal())
	assert.False(t, contract.StatusDraft.IsTerminal())
	assert.False(t, contract.StatusActive.IsTerminal())
}
