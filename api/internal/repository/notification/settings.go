package notification

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// NewNotificationSettingRepository returns notification settings repository.
func NewNotificationSettingRepository(db *gorm.DB) repository.UserNotificationSettingRepository {
	return &notificationSettingRepository{db: db}
}

type notificationSettingRepository struct {
	db *gorm.DB
}

func (r *notificationSettingRepository) GetByUserID(ctx context.Context, userID uint64) (*model.UserNotificationSetting, error) {
	var settings model.UserNotificationSetting
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&settings).Error; err != nil {
		return nil, repository.WrapNotFound(err)
	}
	return &settings, nil
}

func (r *notificationSettingRepository) Upsert(ctx context.Context, settings *model.UserNotificationSetting) error {
	if settings == nil {
		return nil
	}
	var existing model.UserNotificationSetting
	tx := r.db.WithContext(ctx).Where("user_id = ?", settings.UserID).First(&existing)
	if tx.Error == nil {
		settings.ID = existing.ID
	}
	return r.db.WithContext(ctx).Save(settings).Error
}
