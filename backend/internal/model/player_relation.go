package model

// PlayerGame 表示陪玩与可接游戏的关联。
type PlayerGame struct {
	Base
	PlayerID uint64 `json:"player_id" gorm:"index:idx_player_game,unique"`
	GameID   uint64 `json:"game_id" gorm:"index:idx_player_game,unique"`
	IsMain   bool   `json:"is_main"`
}

// PlayerSkillTag 表示陪玩的技能标签。
type PlayerSkillTag struct {
	Base
	PlayerID uint64 `json:"player_id" gorm:"index:idx_player_tag,unique"`
	Tag      string `json:"tag" gorm:"size:32;index:idx_player_tag,unique"`
}
