package reconciliation

import (
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&model.VipLevel{},
		&model.User{},
		&model.Reconciliation{},
		&model.ReconciliationDetail{},
	))
	return db
}

func TestRepository_CreateAndGetWithDetails(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	now := time.Now().UTC()
	rec := &model.Reconciliation{
		ReconciliationNo:   "RCN-TEST-001",
		ReconciliationDate: now,
		Type:               model.ReconciliationTypeManual,
		Status:             model.ReconciliationStatusPending,
		PeriodStart:        now.Add(-24 * time.Hour),
		PeriodEnd:          now,
		Abstract:           "daily reconciliation",
		Details: []model.ReconciliationDetail{
			{
				LineNo:           1,
				ExternalType:     "payment",
				ExternalNo:       "OUT-1",
				ExternalAmount:   1000,
				ExternalDate:     now,
				InternalType:     "payment",
				InternalNo:       "IN-1",
				InternalAmount:   1000,
				InternalDate:     now,
				Status:           "matched",
				DifferenceAmount: 0,
			},
		},
	}

	require.NoError(t, repo.Create(ctx, rec))
	require.NotZero(t, rec.ID)

	got, err := repo.Get(ctx, rec.ID, true)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, rec.ReconciliationNo, got.ReconciliationNo)
	require.Len(t, got.Details, 1)
	assert.Equal(t, "OUT-1", got.Details[0].ExternalNo)
}

func TestRepository_ListWithFilters(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)
	ctx := context.Background()
	now := time.Now().UTC()

	pending := &model.Reconciliation{
		ReconciliationNo:   "RCN-LIST-001",
		ReconciliationDate: now,
		Type:               model.ReconciliationTypePayment,
		Status:             model.ReconciliationStatusPending,
		PeriodStart:        now.Add(-48 * time.Hour),
		PeriodEnd:          now,
	}
	success := &model.Reconciliation{
		ReconciliationNo:   "RCN-LIST-002",
		ReconciliationDate: now,
		Type:               model.ReconciliationTypeBank,
		Status:             model.ReconciliationStatusSuccess,
		PeriodStart:        now.Add(-48 * time.Hour),
		PeriodEnd:          now,
	}
	require.NoError(t, repo.Create(ctx, pending))
	require.NoError(t, repo.Create(ctx, success))

	items, total, err := repo.List(ctx, repository.ReconciliationListOptions{
		Page:     1,
		PageSize: 10,
		Status:   ptrStatus(model.ReconciliationStatusPending),
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, items, 1)
	assert.Equal(t, model.ReconciliationStatusPending, items[0].Status)
}

func TestRepository_ExecuteUpdatesStatusAndSummary(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)
	ctx := context.Background()
	now := time.Now().UTC()

	rec := &model.Reconciliation{
		ReconciliationNo:   "RCN-EXEC-001",
		ReconciliationDate: now,
		Type:               model.ReconciliationTypeInternal,
		Status:             model.ReconciliationStatusPending,
		PeriodStart:        now.Add(-24 * time.Hour),
		PeriodEnd:          now,
		Details: []model.ReconciliationDetail{
			{
				LineNo:           1,
				ExternalType:     "payment",
				ExternalNo:       "OUT-1",
				ExternalAmount:   1000,
				ExternalDate:     now,
				InternalType:     "payment",
				InternalNo:       "IN-1",
				InternalAmount:   1000,
				InternalDate:     now,
				Status:           "matched",
				DifferenceAmount: 0,
			},
			{
				LineNo:           2,
				ExternalType:     "payment",
				ExternalNo:       "OUT-2",
				ExternalAmount:   2000,
				ExternalDate:     now,
				InternalType:     "payment",
				InternalNo:       "IN-2",
				InternalAmount:   1500,
				InternalDate:     now,
				Status:           "mismatch",
				DifferenceAmount: 500,
			},
		},
	}
	require.NoError(t, repo.Create(ctx, rec))

	executed, err := repo.Execute(ctx, rec.ID, repository.ReconciliationExecuteOptions{
		ProcessedBy: 42,
	})
	require.NoError(t, err)
	require.NotNil(t, executed)
	assert.Equal(t, model.ReconciliationStatusException, executed.Status)
	assert.Equal(t, 2, executed.TotalRecords)
	assert.Equal(t, 1, executed.MatchedRecords)
	assert.Equal(t, int64(500), executed.DifferenceAmount)
	require.NotNil(t, executed.ProcessedBy)
	assert.Equal(t, uint64(42), *executed.ProcessedBy)
	require.NotNil(t, executed.ProcessedAt)

	_, err = repo.Execute(ctx, rec.ID, repository.ReconciliationExecuteOptions{
		ProcessedBy: 42,
	})
	require.Error(t, err)
}

func ptrStatus(v model.ReconciliationStatus) *model.ReconciliationStatus {
	return &v
}
