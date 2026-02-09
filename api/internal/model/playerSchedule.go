package model

// PlayerSchedule stores player availability configuration.
type PlayerSchedule struct {
	Base
	PlayerID        uint64 `json:"playerId" gorm:"column:player_id;uniqueIndex;not null"`
	WeeklySchedule  string `json:"weeklySchedule" gorm:"column:weekly_schedule;type:json;default:'{}'"`
	AutoOffline     bool   `json:"autoOffline" gorm:"column:auto_offline;default:true"`
	MaxOrdersPerDay int    `json:"maxOrdersPerDay" gorm:"column:max_orders_per_day;default:0"`
}

// TableName specifies table name.
func (PlayerSchedule) TableName() string {
	return "player_schedules"
}
