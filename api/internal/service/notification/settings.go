package notification

import (
	"context"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// SettingsService manages notification settings.
type SettingsService struct {
	repo repository.UserNotificationSettingRepository
}

// SettingsPayload represents notification settings used by the app.
type SettingsPayload struct {
	OrderReminder bool `json:"orderReminder"`
	NewMessage    bool `json:"newMessage"`
	SystemNotice  bool `json:"systemNotice"`
}

// NewSettingsService creates notification settings service.
func NewSettingsService(repo repository.UserNotificationSettingRepository) *SettingsService {
	return &SettingsService{repo: repo}
}

// GetSettings returns notification settings, applying defaults if missing.
func (s *SettingsService) GetSettings(ctx context.Context, userID uint64) (*SettingsPayload, error) {
	settings, err := s.repo.GetByUserID(ctx, userID)
	if err != nil && !repository.IsNotFound(err) {
		return nil, err
	}
	if settings == nil {
		return &SettingsPayload{
			OrderReminder: true,
			NewMessage:    true,
			SystemNotice:  true,
		}, nil
	}
	return &SettingsPayload{
		OrderReminder: settings.OrderStatusEnabled,
		NewMessage:    settings.ChatEnabled,
		SystemNotice:  settings.SystemEnabled,
	}, nil
}

// UpdateSettings saves notification settings.
func (s *SettingsService) UpdateSettings(ctx context.Context, userID uint64, payload SettingsPayload) (*SettingsPayload, error) {
	settings, err := s.repo.GetByUserID(ctx, userID)
	if err != nil && !repository.IsNotFound(err) {
		return nil, err
	}
	if settings == nil {
		settings = &model.UserNotificationSetting{
			UserID:             userID,
			InAppEnabled:       true,
			PushEnabled:        true,
			OrderStatusEnabled: true,
			ChatEnabled:        true,
			SystemEnabled:      true,
		}
	}
	settings.OrderStatusEnabled = payload.OrderReminder
	settings.ChatEnabled = payload.NewMessage
	settings.SystemEnabled = payload.SystemNotice

	if err := s.repo.Upsert(ctx, settings); err != nil {
		return nil, err
	}
	return &payload, nil
}
