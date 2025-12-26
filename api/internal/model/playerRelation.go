package model

// PlayerGame 表示陪玩与可接游戏的关联。
type PlayerGame struct {
	Base
	PlayerID uint64 `json:"player_id" gorm:"index:idx_player_game,unique;not null"`
	GameID   uint64 `json:"game_id" gorm:"index:idx_player_game,unique;not null"`
	IsMain   bool   `json:"is_main"`

	// Relations (optional, enable FK with cascade)
	Player Player `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:PlayerID;references:ID"`
	Game   Game   `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:GameID;references:ID"`
}

// PlayerSkillTag 表示陪玩的技能标签。
type PlayerSkillTag struct {
	Base
	PlayerID uint64 `json:"player_id" gorm:"index:idx_player_tag,unique;not null"`
	Tag      string `json:"tag" gorm:"size:32;index:idx_player_tag,unique;not null"`

	Player Player `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:PlayerID;references:ID"`
}
