package contract_test

import (
	"testing"
	"time"

	"github.com/contractiq/contractiq/internal/domain/contract"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMoney(t *testing.T) {
	t.Run("creates valid money", func(t *testing.T) {
		m, err := contract.NewMoney(10050, "USD")
		require.NoError(t, err)
		assert.Equal(t, int64(10050), m.AmountCents)
		assert.Equal(t, "USD", m.Currency)
		assert.Equal(t, 100.50, m.Dollars())
		assert.Equal(t, "100.50 USD", m.String())
	})

	t.Run("rejects negative amount", func(t *testing.T) {
		_, err := contract.NewMoney(-100, "USD")
		assert.Error(t, err)
	})

	t.Run("rejects invalid currency", func(t *testing.T) {
		_, err := contract.NewMoney(100, "US")
		assert.Error(t, err)
	})

	t.Run("equality", func(t *testing.T) {
		m1, _ := contract.NewMoney(100, "USD")
		m2, _ := contract.NewMoney(100, "USD")
		m3, _ := contract.NewMoney(200, "USD")
		assert.True(t, m1.Equals(m2))
		assert.False(t, m1.Equals(m3))
	})
}

func TestDateRange(t *testing.T) {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)

	t.Run("creates valid date range", func(t *testing.T) {
		dr, err := contract.NewDateRange(start, end)
		require.NoError(t, err)
		assert.Equal(t, 364, dr.DurationDays())
	})

	t.Run("rejects end before start", func(t *testing.T) {
		_, err := contract.NewDateRange(end, start)
		assert.Error(t, err)
	})

	t.Run("contains", func(t *testing.T) {
		dr, _ := contract.NewDateRange(start, end)
		assert.True(t, dr.Contains(time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)))
		assert.False(t, dr.Contains(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))
	})

	t.Run("expired", func(t *testing.T) {
		dr, _ := contract.NewDateRange(start, end)
		assert.True(t, dr.IsExpired(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))
		assert.False(t, dr.IsExpired(time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)))
	})
}

func TestClause(t *testing.T) {
	t.Run("creates valid clause", func(t *testing.T) {
		c, err := contract.NewClause("Confidentiality", "All information shall...", 1)
		require.NoError(t, err)
		assert.Equal(t, "Confidentiality", c.Title)
	})

	t.Run("rejects empty title", func(t *testing.T) {
		_, err := contract.NewClause("", "content", 1)
		assert.Error(t, err)
	})

	t.Run("rejects empty content", func(t *testing.T) {
		_, err := contract.NewClause("Title", "", 1)
		assert.Error(t, err)
	})
}
