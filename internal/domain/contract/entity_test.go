package contract_test

import (
	"testing"
	"time"

	"github.com/contractiq/contractiq/internal/domain/contract"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestContract(t *testing.T) *contract.Contract {
	t.Helper()
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	value, err := contract.NewMoney(100000, "USD")
	require.NoError(t, err)
	dateRange, err := contract.NewDateRange(now, now.AddDate(1, 0, 0))
	require.NoError(t, err)

	c, err := contract.NewContract("Test Contract", "A test contract", "owner-123", value, dateRange, now)
	require.NoError(t, err)
	return c
}

func TestNewContract(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	value, _ := contract.NewMoney(100000, "USD")
	dateRange, _ := contract.NewDateRange(now, now.AddDate(1, 0, 0))

	t.Run("creates contract in draft status", func(t *testing.T) {
		c, err := contract.NewContract("NDA Agreement", "Non-disclosure", "owner-1", value, dateRange, now)
		require.NoError(t, err)
		assert.NotEmpty(t, c.ID())
		assert.Equal(t, "NDA Agreement", c.Title())
		assert.Equal(t, contract.StatusDraft, c.Status())
		assert.Equal(t, 1, c.Version())
	})

	t.Run("emits ContractCreated event", func(t *testing.T) {
		c, err := contract.NewContract("NDA", "desc", "owner-1", value, dateRange, now)
		require.NoError(t, err)
		events := c.Events()
		require.Len(t, events, 1)
		assert.Equal(t, "contract.created", events[0].EventName())
	})

	t.Run("rejects empty title", func(t *testing.T) {
		_, err := contract.NewContract("", "desc", "owner-1", value, dateRange, now)
		assert.Error(t, err)
	})

	t.Run("rejects empty owner", func(t *testing.T) {
		_, err := contract.NewContract("Title", "desc", "", value, dateRange, now)
		assert.Error(t, err)
	})
}

func TestContractSubmit(t *testing.T) {
	now := time.Now().UTC()

	t.Run("submits draft contract", func(t *testing.T) {
		c := newTestContract(t)
		_ = c.Events() // clear creation event

		err := c.Submit(now)
		require.NoError(t, err)
		assert.Equal(t, contract.StatusPendingReview, c.Status())

		events := c.Events()
		require.Len(t, events, 1)
		assert.Equal(t, "contract.submitted", events[0].EventName())
	})

	t.Run("rejects submit on non-draft", func(t *testing.T) {
		c := newTestContract(t)
		require.NoError(t, c.Submit(now))

		err := c.Submit(now) // already pending review
		assert.Error(t, err)
	})
}

func TestContractApprove(t *testing.T) {
	now := time.Now().UTC()

	t.Run("approves pending contract", func(t *testing.T) {
		c := newTestContract(t)
		require.NoError(t, c.Submit(now))
		_ = c.Events()

		err := c.Approve("reviewer-1", now)
		require.NoError(t, err)
		assert.Equal(t, contract.StatusApproved, c.Status())
	})

	t.Run("rejects approve on draft", func(t *testing.T) {
		c := newTestContract(t)
		err := c.Approve("reviewer-1", now)
		assert.Error(t, err)
	})
}

func TestContractSign(t *testing.T) {
	now := time.Now().UTC()

	t.Run("signs approved contract", func(t *testing.T) {
		c := newTestContract(t)
		require.NoError(t, c.Submit(now))
		require.NoError(t, c.Approve("reviewer-1", now))
		_ = c.Events()

		err := c.Sign("signer-1", now)
		require.NoError(t, err)
		assert.Equal(t, contract.StatusActive, c.Status())

		events := c.Events()
		require.Len(t, events, 1)
		assert.Equal(t, "contract.signed", events[0].EventName())
	})
}

func TestContractTerminate(t *testing.T) {
	now := time.Now().UTC()

	t.Run("terminates active contract", func(t *testing.T) {
		c := newTestContract(t)
		require.NoError(t, c.Submit(now))
		require.NoError(t, c.Approve("reviewer-1", now))
		require.NoError(t, c.Sign("signer-1", now))
		_ = c.Events()

		err := c.Terminate("breach of terms", now)
		require.NoError(t, err)
		assert.Equal(t, contract.StatusTerminated, c.Status())

		events := c.Events()
		require.Len(t, events, 1)
		assert.Equal(t, "contract.terminated", events[0].EventName())
	})

	t.Run("rejects terminate on draft", func(t *testing.T) {
		c := newTestContract(t)
		err := c.Terminate("reason", now)
		assert.Error(t, err)
	})
}

func TestContractUpdate(t *testing.T) {
	now := time.Now().UTC()
	value, _ := contract.NewMoney(200000, "USD")
	dateRange, _ := contract.NewDateRange(now, now.AddDate(2, 0, 0))

	t.Run("updates draft contract", func(t *testing.T) {
		c := newTestContract(t)
		err := c.Update("Updated Title", "Updated desc", value, dateRange, now)
		require.NoError(t, err)
		assert.Equal(t, "Updated Title", c.Title())
		assert.Equal(t, value, c.Value())
	})

	t.Run("rejects update on non-draft", func(t *testing.T) {
		c := newTestContract(t)
		require.NoError(t, c.Submit(now))

		err := c.Update("Updated", "desc", value, dateRange, now)
		assert.Error(t, err)
	})
}

func TestContractFullLifecycle(t *testing.T) {
	now := time.Now().UTC()

	c := newTestContract(t)
	assert.Equal(t, contract.StatusDraft, c.Status())

	require.NoError(t, c.Submit(now))
	assert.Equal(t, contract.StatusPendingReview, c.Status())

	require.NoError(t, c.Approve("reviewer-1", now))
	assert.Equal(t, contract.StatusApproved, c.Status())

	require.NoError(t, c.Sign("signer-1", now))
	assert.Equal(t, contract.StatusActive, c.Status())

	require.NoError(t, c.Terminate("completed", now))
	assert.Equal(t, contract.StatusTerminated, c.Status())

	// Terminal state - no more transitions
	assert.Error(t, c.Expire(now))
}
