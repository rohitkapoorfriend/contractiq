package contract

import (
	"fmt"
	"time"
)

// Money represents a monetary value with amount in cents and currency code.
type Money struct {
	AmountCents int64  `json:"amount_cents"`
	Currency    string `json:"currency"`
}

// NewMoney creates a new Money value object with validation.
func NewMoney(amountCents int64, currency string) (Money, error) {
	if amountCents < 0 {
		return Money{}, fmt.Errorf("amount cannot be negative")
	}
	if len(currency) != 3 {
		return Money{}, fmt.Errorf("currency must be a 3-letter ISO code")
	}
	return Money{AmountCents: amountCents, Currency: currency}, nil
}

// Dollars returns the amount as a floating point dollar value.
func (m Money) Dollars() float64 {
	return float64(m.AmountCents) / 100
}

func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", m.Dollars(), m.Currency)
}

// Equals compares two Money values.
func (m Money) Equals(other Money) bool {
	return m.AmountCents == other.AmountCents && m.Currency == other.Currency
}

// Clause represents a single clause within a contract.
type Clause struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Order   int    `json:"order"`
}

// NewClause creates a validated clause.
func NewClause(title, content string, order int) (Clause, error) {
	if title == "" {
		return Clause{}, fmt.Errorf("clause title is required")
	}
	if content == "" {
		return Clause{}, fmt.Errorf("clause content is required")
	}
	if order < 0 {
		return Clause{}, fmt.Errorf("clause order must be non-negative")
	}
	return Clause{Title: title, Content: content, Order: order}, nil
}

// DateRange represents a time period with start and end dates.
type DateRange struct {
	Start time.Time `json:"start_date"`
	End   time.Time `json:"end_date"`
}

// NewDateRange creates a validated date range.
func NewDateRange(start, end time.Time) (DateRange, error) {
	if end.Before(start) {
		return DateRange{}, fmt.Errorf("end date must be after start date")
	}
	return DateRange{Start: start, End: end}, nil
}

// Contains checks if a given time falls within the date range.
func (d DateRange) Contains(t time.Time) bool {
	return (t.Equal(d.Start) || t.After(d.Start)) && (t.Equal(d.End) || t.Before(d.End))
}

// DurationDays returns the number of days in the range.
func (d DateRange) DurationDays() int {
	return int(d.End.Sub(d.Start).Hours() / 24)
}

// IsExpired checks if the date range has passed relative to the given time.
func (d DateRange) IsExpired(now time.Time) bool {
	return now.After(d.End)
}
