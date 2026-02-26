package reconciliation

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

type fakeRepo struct {
	listFn    func(ctx context.Context, opts repository.ReconciliationListOptions) ([]model.Reconciliation, int64, error)
	getFn     func(ctx context.Context, id uint64, withDetails bool) (*model.Reconciliation, error)
	createFn  func(ctx context.Context, rec *model.Reconciliation) error
	executeFn func(ctx context.Context, id uint64, opts repository.ReconciliationExecuteOptions) (*model.Reconciliation, error)
}

func (f *fakeRepo) List(ctx context.Context, opts repository.ReconciliationListOptions) ([]model.Reconciliation, int64, error) {
	if f.listFn != nil {
		return f.listFn(ctx, opts)
	}
	return nil, 0, nil
}

func (f *fakeRepo) Get(ctx context.Context, id uint64, withDetails bool) (*model.Reconciliation, error) {
	if f.getFn != nil {
		return f.getFn(ctx, id, withDetails)
	}
	return nil, repository.ErrNotFound
}

func (f *fakeRepo) Create(ctx context.Context, rec *model.Reconciliation) error {
	if f.createFn != nil {
		return f.createFn(ctx, rec)
	}
	return nil
}

func (f *fakeRepo) Execute(ctx context.Context, id uint64, opts repository.ReconciliationExecuteOptions) (*model.Reconciliation, error) {
	if f.executeFn != nil {
		return f.executeFn(ctx, id, opts)
	}
	return nil, nil
}

func TestService_CreateSetsDefaultsAndSummary(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	var captured *model.Reconciliation

	repo := &fakeRepo{
		createFn: func(_ context.Context, rec *model.Reconciliation) error {
			captured = rec
			rec.ID = 100
			return nil
		},
	}
	svc := NewService(repo)

	created, err := svc.Create(ctx, CreateInput{
		ReconciliationDate: now,
		Type:               model.ReconciliationTypeManual,
		PeriodStart:        now.Add(-24 * time.Hour),
		PeriodEnd:          now,
		Abstract:           "manual check",
		Details: []CreateDetailInput{
			{
				ExternalType:   "payment",
				ExternalNo:     "EXT-1",
				ExternalAmount: 1000,
				ExternalDate:   now,
				InternalType:   "payment",
				InternalNo:     "INT-1",
				InternalAmount: 1000,
				InternalDate:   now,
			},
			{
				ExternalType:   "payment",
				ExternalNo:     "EXT-2",
				ExternalAmount: 1500,
				ExternalDate:   now,
				InternalType:   "payment",
				InternalNo:     "INT-2",
				InternalAmount: 1000,
				InternalDate:   now,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotNil(t, captured)
	assert.NotEmpty(t, captured.ReconciliationNo)
	assert.Equal(t, model.ReconciliationStatusPending, captured.Status)
	assert.Equal(t, 2, captured.TotalRecords)
	assert.Equal(t, 1, captured.MatchedRecords)
	assert.Equal(t, int64(500), captured.DifferenceAmount)
	require.Len(t, captured.Details, 2)
	assert.Equal(t, 1, captured.Details[0].LineNo)
	assert.Equal(t, int64(500), captured.Details[1].DifferenceAmount)
}

func TestService_ListNormalizesPagination(t *testing.T) {
	ctx := context.Background()

	repo := &fakeRepo{
		listFn: func(_ context.Context, opts repository.ReconciliationListOptions) ([]model.Reconciliation, int64, error) {
			assert.Equal(t, 1, opts.Page)
			assert.Equal(t, 20, opts.PageSize)
			return []model.Reconciliation{{}}, 1, nil
		},
	}
	svc := NewService(repo)

	items, pagination, err := svc.List(ctx, ListInput{})
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.NotNil(t, pagination)
	assert.Equal(t, 1, pagination.Page)
	assert.Equal(t, 20, pagination.PageSize)
	assert.Equal(t, 1, pagination.Total)
}

func TestService_ExecuteRejectsInvalidTargetStatus(t *testing.T) {
	ctx := context.Background()
	svc := NewService(&fakeRepo{})

	_, err := svc.Execute(ctx, 1, 10, ExecuteInput{
		TargetStatus: ptrExecStatus(model.ReconciliationStatusPending),
	})
	require.Error(t, err)
}

func TestService_CreateRejectsInvalidPeriodRange(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	svc := NewService(&fakeRepo{})

	_, err := svc.Create(ctx, CreateInput{
		ReconciliationDate: now,
		Type:               model.ReconciliationTypeManual,
		PeriodStart:        now,
		PeriodEnd:          now.Add(-time.Minute),
	})
	require.Error(t, err)
}

func TestService_ExecuteSuccess(t *testing.T) {
	ctx := context.Background()
	repo := &fakeRepo{
		executeFn: func(_ context.Context, id uint64, opts repository.ReconciliationExecuteOptions) (*model.Reconciliation, error) {
			require.Equal(t, uint64(88), id)
			require.Equal(t, uint64(33), opts.ProcessedBy)
			return &model.Reconciliation{
				Base:   model.Base{ID: id},
				Status: model.ReconciliationStatusSuccess,
			}, nil
		},
	}

	svc := NewService(repo)
	status := model.ReconciliationStatusSuccess
	got, err := svc.Execute(ctx, 88, 33, ExecuteInput{
		TargetStatus: &status,
	})
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, uint64(88), got.ID)
	assert.Equal(t, model.ReconciliationStatusSuccess, got.Status)
}

func ptrExecStatus(v model.ReconciliationStatus) *model.ReconciliationStatus {
	return &v
}
