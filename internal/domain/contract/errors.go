package contract

import "github.com/contractiq/contractiq/pkg/apperror"

var (
	ErrNotFound             = apperror.NewNotFound("contract", "")
	ErrInvalidTransition    = apperror.NewConflict("invalid contract status transition")
	ErrAlreadySigned        = apperror.NewConflict("contract has already been signed")
	ErrNotActive            = apperror.NewConflict("contract is not active")
	ErrConcurrencyConflict  = apperror.NewConcurrencyConflict("contract")
)

func NewErrNotFound(id string) *apperror.Error {
	return apperror.NewNotFound("contract", id)
}
