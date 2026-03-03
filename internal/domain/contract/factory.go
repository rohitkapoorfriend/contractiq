package contract

import (
	"fmt"
	"time"

	"github.com/contractiq/contractiq/pkg/identifier"
)

// NewContract creates a new Contract aggregate in Draft status with validation.
func NewContract(title, description, ownerID string, value Money, dateRange DateRange, now time.Time) (*Contract, error) {
	if title == "" {
		return nil, fmt.Errorf("contract title is required")
	}
	if ownerID == "" {
		return nil, fmt.Errorf("contract owner is required")
	}

	id := identifier.New()
	c := &Contract{
		id:          id,
		title:       title,
		description: description,
		status:      StatusDraft,
		value:       value,
		dateRange:   dateRange,
		ownerID:     ownerID,
		version:     1,
		createdAt:   now,
		updatedAt:   now,
	}
	c.recordEvent(NewContractCreatedEvent(id, title, now))
	return c, nil
}

// NewContractFromTemplate creates a contract pre-populated from a template.
func NewContractFromTemplate(
	title, description, ownerID, templateID string,
	clauses []Clause,
	value Money,
	dateRange DateRange,
	now time.Time,
) (*Contract, error) {
	c, err := NewContract(title, description, ownerID, value, dateRange, now)
	if err != nil {
		return nil, err
	}
	c.templateID = &templateID
	c.clauses = clauses
	return c, nil
}

// SetParty assigns a party to the contract.
func (c *Contract) SetParty(partyID string) error {
	if partyID == "" {
		return fmt.Errorf("party ID is required")
	}
	c.partyID = &partyID
	return nil
}
