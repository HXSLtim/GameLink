package model

// Review captures a user's rating and feedback for a completed order/player.
type Review struct {
	Base
	OrderID  uint64 `json:"orderId" gorm:"column:order_id;index"`
	UserID   uint64 `json:"reviewerId" gorm:"column:user_id;index"` // camelCase: reviewerId
	PlayerID uint64 `json:"playerId" gorm:"column:player_id;index"`
	Score    Rating `json:"rating" gorm:"column:score;type:tinyint"`           // 1-5
	Content  string `json:"comment,omitempty" gorm:"column:content;type:text"` // optional comment
}
