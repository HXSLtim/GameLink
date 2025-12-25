package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
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

// ReviewReplyStatus 评价回复状态
type ReviewReplyStatus string

const (
	ReviewReplyStatusPending  ReviewReplyStatus = "pending"  // 待审核
	ReviewReplyStatusApproved ReviewReplyStatus = "approved" // 已通过
	ReviewReplyStatusRejected ReviewReplyStatus = "rejected" // 已拒绝
)

// AppealStatus 申诉状态
type AppealStatus string

const (
	AppealStatusPending  AppealStatus = "pending"  // 待处理
	AppealStatusApproved AppealStatus = "approved" // 已通过
	AppealStatusRejected AppealStatus = "rejected" // 已拒绝
)

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
	// 订单明细ID（多人订单时关联具体座位）
	OrderItemID *uint64 `json:"orderItemId,omitempty" gorm:"column:order_item_id;index"`
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
	// 评价图片URL数组（最多3张）
	// @Example ["https://example.com/image1.jpg", "https://example.com/image2.jpg"]
	Images StringArray `json:"images,omitempty" gorm:"column:images;type:json"`
	// 拒绝原因（当状态为rejected时）
	// @Example 评价内容包含敏感词
	RejectionReason string `json:"rejectionReason,omitempty" gorm:"column:rejection_reason;type:text"`

	// 公开与匿名设置
	IsPublic    bool `json:"isPublic" gorm:"column:is_public;default:false"`       // 是否公开（默认不公开）
	IsAnonymous bool `json:"isAnonymous" gorm:"column:is_anonymous;default:false"` // 是否匿名

	// 修改记录
	EditCount  int        `json:"editCount" gorm:"column:edit_count;default:0"`    // 修改次数（最多3次）
	LastEditAt *time.Time `json:"lastEditAt,omitempty" gorm:"column:last_edit_at"` // 最后修改时间

	// 评价窗口
	ExpireAt time.Time `json:"expireAt" gorm:"column:expire_at;index"` // 评价截止时间（订单完成后7天）

	// Relations
	OrderItem *OrderItem `json:"orderItem,omitempty" gorm:"foreignKey:OrderItemID"`
}
