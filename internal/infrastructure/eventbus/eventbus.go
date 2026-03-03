package eventbus

import (
	"sync"

	"github.com/contractiq/contractiq/internal/domain/event"
	"go.uber.org/zap"
)

// InMemoryPublisher is a simple in-memory event publisher for domain events.
type InMemoryPublisher struct {
	mu       sync.RWMutex
	handlers map[string][]event.Handler
	logger   *zap.Logger
}

// NewInMemoryPublisher creates a new in-memory event publisher.
func NewInMemoryPublisher(logger *zap.Logger) *InMemoryPublisher {
	return &InMemoryPublisher{
		handlers: make(map[string][]event.Handler),
		logger:   logger,
	}
}

// Subscribe registers a handler for a specific event name.
func (p *InMemoryPublisher) Subscribe(eventName string, handler event.Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.handlers[eventName] = append(p.handlers[eventName], handler)
}

// Publish dispatches events to all registered handlers.
func (p *InMemoryPublisher) Publish(events ...event.Event) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, evt := range events {
		p.logger.Debug("publishing domain event",
			zap.String("event", evt.EventName()),
			zap.String("aggregate_id", evt.AggregateID()),
		)

		handlers, ok := p.handlers[evt.EventName()]
		if !ok {
			continue
		}
		for _, h := range handlers {
			h.Handle(evt)
		}
	}
}
