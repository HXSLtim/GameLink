package player

import (
	"context"
	"encoding/json"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
)

// ScheduleService manages player schedules.
type ScheduleService struct {
	schedules repository.PlayerScheduleRepository
	players   repository.PlayerRepository
}

// ScheduleDay represents availability for a day.
type ScheduleDay struct {
	Enabled bool   `json:"enabled"`
	Start   string `json:"start"`
	End     string `json:"end"`
}

// SchedulePayload represents player schedule.
type SchedulePayload struct {
	WeeklySchedule  map[string]ScheduleDay `json:"weeklySchedule"`
	AutoOffline     bool                   `json:"autoOffline"`
	MaxOrdersPerDay int                    `json:"maxOrdersPerDay"`
}

// NewScheduleService creates schedule service.
func NewScheduleService(
	schedules repository.PlayerScheduleRepository,
	players repository.PlayerRepository,
) *ScheduleService {
	return &ScheduleService{
		schedules: schedules,
		players:   players,
	}
}

// GetSchedule returns player schedule or default.
func (s *ScheduleService) GetSchedule(ctx context.Context, userID uint64) (*SchedulePayload, error) {
	player, err := s.players.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	schedule, err := s.schedules.GetByPlayerID(ctx, player.ID)
	if err != nil && !repository.IsNotFound(err) {
		return nil, err
	}
	if schedule == nil {
		return defaultSchedule(), nil
	}
	return toSchedulePayload(schedule), nil
}

// UpdateSchedule saves player schedule.
func (s *ScheduleService) UpdateSchedule(ctx context.Context, userID uint64, payload SchedulePayload) (*SchedulePayload, error) {
	player, err := s.players.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if payload.WeeklySchedule == nil {
		payload.WeeklySchedule = map[string]ScheduleDay{}
	}
	if err := validateWeeklySchedule(payload.WeeklySchedule); err != nil {
		return nil, err
	}
	if payload.MaxOrdersPerDay < 0 {
		payload.MaxOrdersPerDay = 0
	}
	weeklyJSON, _ := json.Marshal(payload.WeeklySchedule)
	schedule := &model.PlayerSchedule{
		PlayerID:        player.ID,
		WeeklySchedule:  string(weeklyJSON),
		AutoOffline:     payload.AutoOffline,
		MaxOrdersPerDay: payload.MaxOrdersPerDay,
	}
	if err := s.schedules.Upsert(ctx, schedule); err != nil {
		return nil, err
	}
	return &payload, nil
}

func defaultSchedule() *SchedulePayload {
	return &SchedulePayload{
		WeeklySchedule:  map[string]ScheduleDay{},
		AutoOffline:     true,
		MaxOrdersPerDay: 0,
	}
}

func toSchedulePayload(schedule *model.PlayerSchedule) *SchedulePayload {
	payload := defaultSchedule()
	if schedule == nil {
		return payload
	}
	payload.AutoOffline = schedule.AutoOffline
	payload.MaxOrdersPerDay = schedule.MaxOrdersPerDay
	if schedule.WeeklySchedule != "" {
		var weekly map[string]ScheduleDay
		if err := json.Unmarshal([]byte(schedule.WeeklySchedule), &weekly); err == nil && weekly != nil {
			payload.WeeklySchedule = weekly
		}
	}
	return payload
}

func validateWeeklySchedule(schedule map[string]ScheduleDay) error {
	for _, day := range schedule {
		if !day.Enabled {
			continue
		}
		if day.Start == "" || day.End == "" {
			return apierr.BadRequest("invalid schedule time range")
		}
		start, err := time.Parse("15:04", day.Start)
		if err != nil {
			return apierr.BadRequest("invalid schedule start time")
		}
		end, err := time.Parse("15:04", day.End)
		if err != nil {
			return apierr.BadRequest("invalid schedule end time")
		}
		if !start.Before(end) {
			return apierr.BadRequest("schedule start must be before end")
		}
	}
	return nil
}
