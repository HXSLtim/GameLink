package model

// Review captures a user's rating and feedback for a completed order/player.
type Review struct {
	Base
	OrderID  uint64 `json:"order_id" gorm:"index"`
	UserID   uint64 `json:"user_id" gorm:"index"`
	PlayerID uint64 `json:"player_id" gorm:"index"`
	Score    Rating `json:"score" gorm:"type:tinyint"`          // 1-5
	Content  string `json:"content,omitempty" gorm:"type:text"` // optional comment
}
