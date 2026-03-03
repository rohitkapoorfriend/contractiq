package party

import (
	"fmt"
	"time"

	"github.com/contractiq/contractiq/pkg/identifier"
)

// PartyType distinguishes between organizations and individuals.
type PartyType string

const (
	PartyTypeOrganization PartyType = "ORGANIZATION"
	PartyTypeIndividual   PartyType = "INDIVIDUAL"
)

// Party represents an external party involved in contracts.
type Party struct {
	id        string
	name      string
	email     string
	partyType PartyType
	company   string
	phone     string
	address   string
	createdBy string
	version   int
	createdAt time.Time
	updatedAt time.Time
}

func (p *Party) ID() string          { return p.id }
func (p *Party) Name() string        { return p.name }
func (p *Party) Email() string       { return p.email }
func (p *Party) Type() PartyType     { return p.partyType }
func (p *Party) Company() string     { return p.company }
func (p *Party) Phone() string       { return p.phone }
func (p *Party) Address() string     { return p.address }
func (p *Party) CreatedBy() string   { return p.createdBy }
func (p *Party) Version() int        { return p.version }
func (p *Party) CreatedAt() time.Time { return p.createdAt }
func (p *Party) UpdatedAt() time.Time { return p.updatedAt }

// NewParty creates a new Party with validation.
func NewParty(name, email string, partyType PartyType, company, phone, address, createdBy string, now time.Time) (*Party, error) {
	if name == "" {
		return nil, fmt.Errorf("party name is required")
	}
	if email == "" {
		return nil, fmt.Errorf("party email is required")
	}
	if partyType != PartyTypeOrganization && partyType != PartyTypeIndividual {
		return nil, fmt.Errorf("invalid party type: %s", partyType)
	}
	return &Party{
		id:        identifier.New(),
		name:      name,
		email:     email,
		partyType: partyType,
		company:   company,
		phone:     phone,
		address:   address,
		createdBy: createdBy,
		version:   1,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// Update modifies party details.
func (p *Party) Update(name, email, company, phone, address string, now time.Time) error {
	if name == "" {
		return fmt.Errorf("party name is required")
	}
	if email == "" {
		return fmt.Errorf("party email is required")
	}
	p.name = name
	p.email = email
	p.company = company
	p.phone = phone
	p.address = address
	p.updatedAt = now
	return nil
}

// Reconstitute creates a Party from persisted data.
func Reconstitute(
	id, name, email string,
	partyType PartyType,
	company, phone, address, createdBy string,
	version int,
	createdAt, updatedAt time.Time,
) *Party {
	return &Party{
		id:        id,
		name:      name,
		email:     email,
		partyType: partyType,
		company:   company,
		phone:     phone,
		address:   address,
		createdBy: createdBy,
		version:   version,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}
