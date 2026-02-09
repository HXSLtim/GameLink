package user

import (
	"context"
	"encoding/json"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
)

// SettingsService manages user settings.
type SettingsService struct {
	repo repository.UserSettingsRepository
}

// SettingsPayload represents user settings.
type SettingsPayload struct {
	Theme         string          `json:"theme"`
	Language      string          `json:"language"`
	Notifications map[string]bool `json:"notifications"`
	Privacy       map[string]bool `json:"privacy"`
}

// NewSettingsService creates settings service.
func NewSettingsService(repo repository.UserSettingsRepository) *SettingsService {
	return &SettingsService{repo: repo}
}

// GetSettings returns settings, falling back to defaults.
func (s *SettingsService) GetSettings(ctx context.Context, userID uint64) (*SettingsPayload, error) {
	settings, err := s.repo.GetByUserID(ctx, userID)
	if err != nil && !repository.IsNotFound(err) {
		return nil, err
	}
	if settings == nil {
		return defaultSettings(), nil
	}
	return toSettingsPayload(settings), nil
}

// UpdateSettings saves new settings.
func (s *SettingsService) UpdateSettings(ctx context.Context, userID uint64, payload SettingsPayload) (*SettingsPayload, error) {
	if payload.Theme != "" && !isValidTheme(payload.Theme) {
		return nil, apierr.BadRequest("invalid theme")
	}
	if payload.Language != "" && !isValidLanguage(payload.Language) {
		return nil, apierr.BadRequest("invalid language")
	}
	if payload.Theme == "" {
		payload.Theme = "auto"
	}
	if payload.Language == "" {
		payload.Language = "zh-CN"
	}
	if payload.Notifications == nil {
		payload.Notifications = map[string]bool{}
	}
	if payload.Privacy == nil {
		payload.Privacy = map[string]bool{}
	}

	notificationsJSON, _ := json.Marshal(payload.Notifications)
	privacyJSON, _ := json.Marshal(payload.Privacy)

	settings := &model.UserSettings{
		UserID:        userID,
		Theme:         payload.Theme,
		Language:      payload.Language,
		Notifications: string(notificationsJSON),
		Privacy:       string(privacyJSON),
	}
	if err := s.repo.Upsert(ctx, settings); err != nil {
		return nil, err
	}
	return &payload, nil
}

func defaultSettings() *SettingsPayload {
	return &SettingsPayload{
		Theme:    "auto",
		Language: "zh-CN",
		Notifications: map[string]bool{
			"orderReminder": true,
			"newMessage":    true,
			"systemNotice":  true,
		},
		Privacy: map[string]bool{
			"showOnlineStatus":     true,
			"allowStrangerMessage": true,
		},
	}
}

func toSettingsPayload(settings *model.UserSettings) *SettingsPayload {
	payload := defaultSettings()
	if settings == nil {
		return payload
	}
	if settings.Theme != "" {
		payload.Theme = settings.Theme
	}
	if settings.Language != "" {
		payload.Language = settings.Language
	}
	if settings.Notifications != "" {
		var notifications map[string]bool
		if err := json.Unmarshal([]byte(settings.Notifications), &notifications); err == nil && notifications != nil {
			payload.Notifications = notifications
		}
	}
	if settings.Privacy != "" {
		var privacy map[string]bool
		if err := json.Unmarshal([]byte(settings.Privacy), &privacy); err == nil && privacy != nil {
			payload.Privacy = privacy
		}
	}
	return payload
}

func isValidTheme(theme string) bool {
	switch theme {
	case "auto", "light", "dark":
		return true
	default:
		return false
	}
}

func isValidLanguage(language string) bool {
	switch language {
	case "zh-CN", "en-US":
		return true
	default:
		return false
	}
}
