package dto

import "time"

// CreateContractRequest is the input for creating a new contract.
type CreateContractRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=200"`
	Description string `json:"description" validate:"max=2000"`
	PartyID     string `json:"party_id" validate:"omitempty,uuid"`
	TemplateID  string `json:"template_id" validate:"omitempty,uuid"`
	Value       MoneyDTO   `json:"value" validate:"required"`
	StartDate   time.Time  `json:"start_date" validate:"required"`
	EndDate     time.Time  `json:"end_date" validate:"required,gtfield=StartDate"`
	Clauses     []ClauseDTO `json:"clauses" validate:"dive"`
}

// UpdateContractRequest is the input for updating an existing contract.
type UpdateContractRequest struct {
	Title       string     `json:"title" validate:"required,min=3,max=200"`
	Description string     `json:"description" validate:"max=2000"`
	Value       MoneyDTO   `json:"value" validate:"required"`
	StartDate   time.Time  `json:"start_date" validate:"required"`
	EndDate     time.Time  `json:"end_date" validate:"required,gtfield=StartDate"`
	Clauses     []ClauseDTO `json:"clauses" validate:"dive"`
	Version     int        `json:"version" validate:"required,min=1"`
}

// TerminateRequest is the input for terminating a contract.
type TerminateRequest struct {
	Reason string `json:"reason" validate:"required,min=5,max=500"`
}

// MoneyDTO represents the monetary value in API requests/responses.
type MoneyDTO struct {
	AmountCents int64  `json:"amount_cents" validate:"min=0"`
	Currency    string `json:"currency" validate:"required,len=3"`
}

// ClauseDTO represents a clause in API requests/responses.
type ClauseDTO struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	Order   int    `json:"order" validate:"min=0"`
}

// ContractResponse is the output representation of a contract.
type ContractResponse struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      string       `json:"status"`
	Value       MoneyDTO     `json:"value"`
	Clauses     []ClauseDTO  `json:"clauses"`
	StartDate   time.Time    `json:"start_date"`
	EndDate     time.Time    `json:"end_date"`
	OwnerID     string       `json:"owner_id"`
	PartyID     *string      `json:"party_id,omitempty"`
	TemplateID  *string      `json:"template_id,omitempty"`
	Version     int          `json:"version"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// ListContractsRequest represents query parameters for listing contracts.
type ListContractsRequest struct {
	Status   string `json:"status" validate:"omitempty,oneof=DRAFT PENDING_REVIEW APPROVED ACTIVE EXPIRED TERMINATED"`
	PartyID  string `json:"party_id" validate:"omitempty,uuid"`
	Search   string `json:"search" validate:"omitempty,max=100"`
	Page     int    `json:"page" validate:"omitempty,min=1"`
	PageSize int    `json:"page_size" validate:"omitempty,min=1,max=100"`
}
