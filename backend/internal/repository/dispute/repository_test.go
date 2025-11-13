package dispute

import (
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

func setupDisputeTest(t *testing.T) repository.DisputeRepository {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&model.OrderDispute{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return NewDisputeRepository(db)
}

func TestDisputeRepository_Create(t *testing.T) {
	repo := setupDisputeTest(t)
	ctx := context.Background()

	dispute := &model.OrderDispute{
		OrderID:     1,
		UserID:      100,
		Reason:      "Product defective",
		Description: "Item arrived broken",
		Status:      model.DisputeStatusPending,
	}

	err := repo.Create(ctx, dispute)
	assert.NoError(t, err)
	assert.NotZero(t, dispute.ID)
}

func TestDisputeRepository_Get(t *testing.T) {
	repo := setupDisputeTest(t)
	ctx := context.Background()

	// Create a dispute
	dispute := &model.OrderDispute{
		OrderID:     1,
		UserID:      100,
		Reason:      "Product defective",
		Status:      model.DisputeStatusPending,
	}
	err := repo.Create(ctx, dispute)
	assert.NoError(t, err)

	// Get the dispute
	retrieved, err := repo.Get(ctx, dispute.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, dispute.OrderID, retrieved.OrderID)

	// Get non-existent dispute
	_, err = repo.Get(ctx, 99999)
	assert.Equal(t, repository.ErrNotFound, err)
}

func TestDisputeRepository_GetByOrderID(t *testing.T) {
	repo := setupDisputeTest(t)
	ctx := context.Background()

	// Create a dispute
	dispute := &model.OrderDispute{
		OrderID:     42,
		UserID:      100,
		Reason:      "Product defective",
		Status:      model.DisputeStatusPending,
	}
	err := repo.Create(ctx, dispute)
	assert.NoError(t, err)

	// Get by order ID
	retrieved, err := repo.GetByOrderID(ctx, 42)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, uint64(42), retrieved.OrderID)

	// Get non-existent order ID
	_, err = repo.GetByOrderID(ctx, 99999)
	assert.Equal(t, repository.ErrNotFound, err)
}

func TestDisputeRepository_Update(t *testing.T) {
	repo := setupDisputeTest(t)
	ctx := context.Background()

	// Create a dispute
	dispute := &model.OrderDispute{
		OrderID:     1,
		UserID:      100,
		Reason:      "Product defective",
		Status:      model.DisputeStatusPending,
	}
	err := repo.Create(ctx, dispute)
	assert.NoError(t, err)

	// Update the dispute
	dispute.Status = model.DisputeStatusResolved
	dispute.Resolution = model.ResolutionRefund
	err = repo.Update(ctx, dispute)
	assert.NoError(t, err)

	// Verify update
	retrieved, err := repo.Get(ctx, dispute.ID)
	assert.NoError(t, err)
	assert.Equal(t, model.DisputeStatusResolved, retrieved.Status)
	assert.Equal(t, model.ResolutionRefund, retrieved.Resolution)
}

func TestDisputeRepository_List(t *testing.T) {
	repo := setupDisputeTest(t)
	ctx := context.Background()

	// Create multiple disputes
	for i := 0; i < 5; i++ {
		dispute := &model.OrderDispute{
			OrderID:     uint64(i + 1),
			UserID:      100,
			Reason:      "Reason " + string(rune(i)),
			Status:      model.DisputeStatusPending,
		}
		err := repo.Create(ctx, dispute)
		assert.NoError(t, err)
	}

	// List all disputes
	disputes, total, err := repo.List(ctx, repository.DisputeListOptions{
		Page:     1,
		PageSize: 10,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, disputes, 5)

	// List with status filter
	disputes, total, err = repo.List(ctx, repository.DisputeListOptions{
		Statuses: []model.DisputeStatus{model.DisputeStatusPending},
		Page:     1,
		PageSize: 10,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)

	// List with pagination
	disputes, total, err = repo.List(ctx, repository.DisputeListOptions{
		Page:     1,
		PageSize: 2,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, disputes, 2)
}

func TestDisputeRepository_ListSLABreached(t *testing.T) {
	repo := setupDisputeTest(t)
	ctx := context.Background()

	// Create disputes with SLA breached (sla_breached=false but deadline passed)
	pastTime := time.Now().Add(-1 * time.Hour)
	for i := 0; i < 2; i++ {
		dispute := &model.OrderDispute{
			OrderID:     uint64(i + 1),
			UserID:      100,
			Reason:      "Reason",
			Status:      model.DisputeStatusPending,
			SLABreached: false,
			SLADeadline: &pastTime,
		}
		err := repo.Create(ctx, dispute)
		assert.NoError(t, err)
	}

	// List SLA breached
	disputes, err := repo.ListSLABreached(ctx)
	assert.NoError(t, err)
	assert.Len(t, disputes, 2)
}

func TestDisputeRepository_MarkSLABreached(t *testing.T) {
	repo := setupDisputeTest(t)
	ctx := context.Background()

	// Create a dispute
	dispute := &model.OrderDispute{
		OrderID:     1,
		UserID:      100,
		Reason:      "Reason",
		Status:      model.DisputeStatusPending,
		SLABreached: false,
	}
	err := repo.Create(ctx, dispute)
	assert.NoError(t, err)

	// Mark as SLA breached
	err = repo.MarkSLABreached(ctx, dispute.ID)
	assert.NoError(t, err)

	// Verify
	retrieved, err := repo.Get(ctx, dispute.ID)
	assert.NoError(t, err)
	assert.True(t, retrieved.SLABreached)
}

func TestDisputeRepository_Delete(t *testing.T) {
	repo := setupDisputeTest(t)
	ctx := context.Background()

	// Create a dispute
	dispute := &model.OrderDispute{
		OrderID:     1,
		UserID:      100,
		Reason:      "Reason",
		Status:      model.DisputeStatusPending,
	}
	err := repo.Create(ctx, dispute)
	assert.NoError(t, err)

	// Delete
	err = repo.Delete(ctx, dispute.ID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = repo.Get(ctx, dispute.ID)
	assert.Equal(t, repository.ErrNotFound, err)
}

func TestDisputeRepository_CountByStatus(t *testing.T) {
	repo := setupDisputeTest(t)
	ctx := context.Background()

	// Create disputes with different statuses
	for i := 0; i < 3; i++ {
		dispute := &model.OrderDispute{
			OrderID:     uint64(i + 1),
			UserID:      100,
			Reason:      "Reason",
			Status:      model.DisputeStatusPending,
		}
		err := repo.Create(ctx, dispute)
		assert.NoError(t, err)
	}

	// Count by status
	count, err := repo.CountByStatus(ctx, model.DisputeStatusPending)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

func TestDisputeRepository_GetPendingCount(t *testing.T) {
	repo := setupDisputeTest(t)
	ctx := context.Background()

	// Create pending disputes
	for i := 0; i < 2; i++ {
		dispute := &model.OrderDispute{
			OrderID:     uint64(i + 1),
			UserID:      100,
			Reason:      "Reason",
			Status:      model.DisputeStatusPending,
		}
		err := repo.Create(ctx, dispute)
		assert.NoError(t, err)
	}

	// Get pending count
	count, err := repo.GetPendingCount(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}
