package user

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// NewUserSettingsRepository returns settings repository.
func NewUserSettingsRepository(db *gorm.DB) repository.UserSettingsRepository {
	return &userSettingsRepository{db: db}
}

type userSettingsRepository struct {
	db *gorm.DB
}

func (r *userSettingsRepository) GetByUserID(ctx context.Context, userID uint64) (*model.UserSettings, error) {
	var settings model.UserSettings
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&settings).Error; err != nil {
		return nil, repository.WrapNotFound(err)
	}
	return &settings, nil
}

func (r *userSettingsRepository) Upsert(ctx context.Context, settings *model.UserSettings) error {
	if settings == nil {
		return nil
	}
	var existing model.UserSettings
	tx := r.db.WithContext(ctx).Where("user_id = ?", settings.UserID).First(&existing)
	if tx.Error == nil {
		settings.ID = existing.ID
	}
	return r.db.WithContext(ctx).Save(settings).Error
}
