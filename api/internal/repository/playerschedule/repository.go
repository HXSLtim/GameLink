package playerschedule

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// NewPlayerScheduleRepository creates player schedule repo.
func NewPlayerScheduleRepository(db *gorm.DB) repository.PlayerScheduleRepository {
	return &gormPlayerScheduleRepository{db: db}
}

type gormPlayerScheduleRepository struct {
	db *gorm.DB
}

func (r *gormPlayerScheduleRepository) GetByPlayerID(ctx context.Context, playerID uint64) (*model.PlayerSchedule, error) {
	var schedule model.PlayerSchedule
	if err := r.db.WithContext(ctx).Where("player_id = ?", playerID).First(&schedule).Error; err != nil {
		return nil, repository.WrapNotFound(err)
	}
	return &schedule, nil
}

func (r *gormPlayerScheduleRepository) Upsert(ctx context.Context, schedule *model.PlayerSchedule) error {
	if schedule == nil {
		return nil
	}
	var existing model.PlayerSchedule
	tx := r.db.WithContext(ctx).Where("player_id = ?", schedule.PlayerID).First(&existing)
	if tx.Error == nil {
		schedule.ID = existing.ID
	}
	return r.db.WithContext(ctx).Save(schedule).Error
}
