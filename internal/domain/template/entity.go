package template

import (
	"fmt"
	"time"

	"github.com/contractiq/contractiq/internal/domain/contract"
	"github.com/contractiq/contractiq/pkg/identifier"
)

// Template represents a reusable contract template.
type Template struct {
	id          string
	name        string
	description string
	clauses     []contract.Clause
	isActive    bool
	createdBy   string
	version     int
	createdAt   time.Time
	updatedAt   time.Time
}

func (t *Template) ID() string              { return t.id }
func (t *Template) Name() string             { return t.name }
func (t *Template) Description() string      { return t.description }
func (t *Template) Clauses() []contract.Clause { return t.clauses }
func (t *Template) IsActive() bool           { return t.isActive }
func (t *Template) CreatedBy() string        { return t.createdBy }
func (t *Template) Version() int             { return t.version }
func (t *Template) CreatedAt() time.Time     { return t.createdAt }
func (t *Template) UpdatedAt() time.Time     { return t.updatedAt }

// NewTemplate creates a new template with validation.
func NewTemplate(name, description, createdBy string, clauses []contract.Clause, now time.Time) (*Template, error) {
	if name == "" {
		return nil, fmt.Errorf("template name is required")
	}
	if createdBy == "" {
		return nil, fmt.Errorf("template creator is required")
	}
	return &Template{
		id:          identifier.New(),
		name:        name,
		description: description,
		clauses:     clauses,
		isActive:    true,
		createdBy:   createdBy,
		version:     1,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// Update modifies the template fields.
func (t *Template) Update(name, description string, clauses []contract.Clause, now time.Time) error {
	if name == "" {
		return fmt.Errorf("template name is required")
	}
	t.name = name
	t.description = description
	t.clauses = clauses
	t.updatedAt = now
	return nil
}

// Deactivate marks the template as inactive.
func (t *Template) Deactivate(now time.Time) {
	t.isActive = false
	t.updatedAt = now
}

// Activate marks the template as active.
func (t *Template) Activate(now time.Time) {
	t.isActive = true
	t.updatedAt = now
}

// Reconstitute creates a Template from persisted data.
func Reconstitute(
	id, name, description string,
	clauses []contract.Clause,
	isActive bool,
	createdBy string,
	version int,
	createdAt, updatedAt time.Time,
) *Template {
	return &Template{
		id:          id,
		name:        name,
		description: description,
		clauses:     clauses,
		isActive:    isActive,
		createdBy:   createdBy,
		version:     version,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}
