package clock

import "time"

// Clock provides an abstraction over time for testability.
type Clock interface {
	Now() time.Time
}

// RealClock returns the actual current time.
type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now().UTC() }

// MockClock returns a fixed time, useful for testing.
type MockClock struct {
	FixedTime time.Time
}

func (m MockClock) Now() time.Time { return m.FixedTime }

// New returns a real clock for production use.
func New() Clock { return RealClock{} }
