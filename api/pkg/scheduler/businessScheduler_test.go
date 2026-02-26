package scheduler

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"gamelink/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleDisputeSLA2hBreach_CreateAlertWhenBreached(t *testing.T) {
	repo := &fakeAlertRepo{}
	scheduler := &BusinessScheduler{
		logger:   slog.Default(),
		alertRepo: repo,
	}

	scheduler.handleDisputeSLA2hBreach(context.Background(), 3)

	require.Len(t, repo.created, 1)
	alert := repo.created[0]
	assert.Equal(t, model.AlertLevelHigh, alert.Level)
	assert.Equal(t, model.AlertTypeBusiness, alert.Type)
	assert.Equal(t, "争议 SLA 超时告警", alert.Title)
	assert.Contains(t, alert.Message, "3")
	assert.Equal(t, "scheduler:dispute:sla", alert.Source)
}

func TestHandleDisputeSLA2hBreach_SkipWhenNoBreach(t *testing.T) {
	repo := &fakeAlertRepo{}
	scheduler := &BusinessScheduler{
		logger:   slog.Default(),
		alertRepo: repo,
	}

	scheduler.handleDisputeSLA2hBreach(context.Background(), 0)

	require.Len(t, repo.created, 0)
}

func TestHandleDisputeSLA2hBreach_AlertCreateError(t *testing.T) {
	repo := &fakeAlertRepo{createErr: errors.New("db down")}
	scheduler := &BusinessScheduler{
		logger:   slog.Default(),
		alertRepo: repo,
	}

	// Should not panic even if alert creation fails.
	require.NotPanics(t, func() {
		scheduler.handleDisputeSLA2hBreach(context.Background(), 2)
	})
	require.Len(t, repo.created, 1)
}

type fakeAlertRepo struct {
	created   []model.Alert
	createErr error
}

func (f *fakeAlertRepo) Create(_ context.Context, alert *model.Alert) error {
	if alert != nil {
		f.created = append(f.created, *alert)
	}
	return f.createErr
}

func (f *fakeAlertRepo) GetByID(_ context.Context, _ uint) (*model.Alert, error) {
	return nil, nil
}

func (f *fakeAlertRepo) List(_ context.Context, _ model.AlertQueryOptions) ([]model.Alert, int64, error) {
	return nil, 0, nil
}

func (f *fakeAlertRepo) MarkAsRead(_ context.Context, _ uint) error {
	return nil
}

func (f *fakeAlertRepo) BatchMarkAsRead(_ context.Context, _ []uint) error {
	return nil
}

func (f *fakeAlertRepo) GetUnreadCount(_ context.Context) (int64, error) {
	return int64(0), nil
}
