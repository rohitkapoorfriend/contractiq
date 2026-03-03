package event

import "time"

// Event represents a domain event that occurred in the system.
type Event interface {
	EventName() string
	OccurredAt() time.Time
	AggregateID() string
}

// Base provides common fields for all domain events.
type Base struct {
	Name      string    `json:"event_name"`
	Timestamp time.Time `json:"occurred_at"`
	AggrID    string    `json:"aggregate_id"`
}

func (b Base) EventName() string    { return b.Name }
func (b Base) OccurredAt() time.Time { return b.Timestamp }
func (b Base) AggregateID() string   { return b.AggrID }

// Publisher dispatches domain events to registered handlers.
type Publisher interface {
	Publish(events ...Event)
}

// Handler processes a specific domain event.
type Handler interface {
	Handle(event Event)
}
