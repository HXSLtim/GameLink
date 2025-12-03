package model

// Review captures a user's rating and feedback for a completed order/player.
// @Description 订单评价模型，记录用户对陪玩师的服务评价
// @Example [{"id": 1, "orderId": 1001, "reviewerId": 2001, "playerId": 3001, "rating": 5, "comment": "服务非常棒，技术高超！", "createdAt": "2024-01-15T10:30:00Z", "updatedAt": "2024-01-15T10:30:00Z"}]
type Review struct {
	Base
	// 订单ID
	// @Example 1001
	OrderID uint64 `json:"orderId" gorm:"column:order_id;index"`
	// 评价者ID（用户ID）
	// @Example 2001
	UserID uint64 `json:"reviewerId" gorm:"column:user_id;index"`
	// 被评价的陪玩师ID
	// @Example 3001
	PlayerID uint64 `json:"playerId" gorm:"column:player_id;index"`
	// 评分，1-5分
	// @Enum 1, 2, 3, 4, 5
	// @Example 5
	Score Rating `json:"rating" gorm:"column:score;type:tinyint"`
	// 评价内容，可选
	// @Example 陪玩师技术很好，服务态度也很棒，下次还会选择！
	Content string `json:"comment,omitempty" gorm:"column:content;type:text"`
}
