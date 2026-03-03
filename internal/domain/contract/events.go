package contract

import (
	"time"

	"github.com/contractiq/contractiq/internal/domain/event"
)

// ContractCreated is emitted when a new contract is created.
type ContractCreated struct {
	event.Base
	Title string `json:"title"`
}

func NewContractCreatedEvent(id, title string, at time.Time) ContractCreated {
	return ContractCreated{
		Base:  event.Base{Name: "contract.created", Timestamp: at, AggrID: id},
		Title: title,
	}
}

// ContractSubmitted is emitted when a contract is submitted for review.
type ContractSubmitted struct {
	event.Base
}

func NewContractSubmittedEvent(id string, at time.Time) ContractSubmitted {
	return ContractSubmitted{
		Base: event.Base{Name: "contract.submitted", Timestamp: at, AggrID: id},
	}
}

// ContractApproved is emitted when a contract is approved.
type ContractApproved struct {
	event.Base
	ApprovedBy string `json:"approved_by"`
}

func NewContractApprovedEvent(id, approvedBy string, at time.Time) ContractApproved {
	return ContractApproved{
		Base:       event.Base{Name: "contract.approved", Timestamp: at, AggrID: id},
		ApprovedBy: approvedBy,
	}
}

// ContractSigned is emitted when a contract is signed and becomes active.
type ContractSigned struct {
	event.Base
	SignedBy string `json:"signed_by"`
}

func NewContractSignedEvent(id, signedBy string, at time.Time) ContractSigned {
	return ContractSigned{
		Base:     event.Base{Name: "contract.signed", Timestamp: at, AggrID: id},
		SignedBy: signedBy,
	}
}

// ContractTerminated is emitted when a contract is terminated.
type ContractTerminated struct {
	event.Base
	Reason string `json:"reason"`
}

func NewContractTerminatedEvent(id, reason string, at time.Time) ContractTerminated {
	return ContractTerminated{
		Base:   event.Base{Name: "contract.terminated", Timestamp: at, AggrID: id},
		Reason: reason,
	}
}
