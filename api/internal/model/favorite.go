package model

// Favorite 用户收藏陪玩师
type Favorite struct {
	Base
	UserID   uint64 `json:"userId" gorm:"column:user_id;uniqueIndex:idx_user_player"`
	PlayerID uint64 `json:"playerId" gorm:"column:player_id;uniqueIndex:idx_user_player;index"`

	// 关联
	User   *User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Player *Player `json:"player,omitempty" gorm:"foreignKey:PlayerID"`
}

// TableName 表名
func (Favorite) TableName() string {
	return "favorites"
}
