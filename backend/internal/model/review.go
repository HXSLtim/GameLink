package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// ReviewStatus 表示评价的审核状态
type ReviewStatus string

const (
	ReviewStatusPending  ReviewStatus = "pending"  // 待审核
	ReviewStatusApproved ReviewStatus = "approved" // 已通过
	ReviewStatusRejected ReviewStatus = "rejected" // 已拒绝
	ReviewStatusDeleted  ReviewStatus = "deleted"  // 已删除
)

// Valid 检查评价状态是否合法
func (rs ReviewStatus) Valid() bool {
	switch rs {
	case ReviewStatusPending, ReviewStatusApproved, ReviewStatusRejected, ReviewStatusDeleted:
		return true
	default:
		return false
	}
}

// StringArray 用于存储字符串数组到数据库
type StringArray []string

// Scan 实现 sql.Scanner 接口
func (sa *StringArray) Scan(value interface{}) error {
	if value == nil {
		*sa = []string{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("failed to scan StringArray: value is not []byte or string")
	}

	if len(bytes) == 0 {
		*sa = []string{}
		return nil
	}

	return json.Unmarshal(bytes, sa)
}

// Value 实现 driver.Valuer 接口
func (sa StringArray) Value() (driver.Value, error) {
	if len(sa) == 0 {
		return "[]", nil
	}
	return json.Marshal(sa)
}

// Review captures a user's rating and feedback for a completed order/player.
// @Description 订单评价模型，记录用户对陪玩师的服务评价
// @Example [{"id": 1, "orderId": 1001, "reviewerId": 2001, "playerId": 3001, "rating": 5, "comment": "服务非常棒，技术高超！", "status": "approved", "isReported": false, "images": [], "createdAt": "2024-01-15T10:30:00Z", "updatedAt": "2024-01-15T10:30:00Z"}]
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
	Score Rating `json:"rating" gorm:"column:score;type:smallint"` // PostgreSQL uses smallint instead of tinyint
	// 评价内容，可选
	// @Example 陪玩师技术很好，服务态度也很棒，下次还会选择！
	Content string `json:"comment,omitempty" gorm:"column:content;type:text"`
	// 审核状态
	// @Enum pending, approved, rejected, deleted
	// @Example approved
	Status ReviewStatus `json:"status" gorm:"column:status;type:varchar(20);default:'pending';index"`
	// 是否被举报
	// @Example false
	IsReported bool `json:"isReported" gorm:"column:is_reported;default:false;index"`
	// 评价图片URL数组
	// @Example ["https://example.com/image1.jpg", "https://example.com/image2.jpg"]
	Images StringArray `json:"images,omitempty" gorm:"column:images;type:json"`
	// 拒绝原因（当状态为rejected时）
	// @Example 评价内容包含敏感词
	RejectionReason string `json:"rejectionReason,omitempty" gorm:"column:rejection_reason;type:text"`
}
