package model

// PlayerService represents player service offerings.
type PlayerService struct {
	Base
	PlayerID    uint64 `json:"playerId" gorm:"column:player_id;index;not null"`
	GameID      uint64 `json:"gameId" gorm:"column:game_id;index;not null"`
	RankID      uint64 `json:"rankId" gorm:"column:rank_id;index;not null"`
	Description string `json:"description" gorm:"type:text"`
	IsActive    bool   `json:"isActive" gorm:"column:is_active;default:true;index"`

	// Relations
	Game *Game     `json:"game,omitempty" gorm:"foreignKey:GameID"`
	Rank *GameRank `json:"rank,omitempty" gorm:"foreignKey:RankID"`
}

// TableName specifies table name.
func (PlayerService) TableName() string {
	return "player_services"
}
