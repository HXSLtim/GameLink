package model

import "time"

// VerificationStatus indicates identity/skill verification flow.
// @Enum pending, verified, rejected
type VerificationStatus string

// VerificationStatus values indicate the verification flow state.
const (
	VerificationPending  VerificationStatus = "pending"
	VerificationVerified VerificationStatus = "verified"
	VerificationRejected VerificationStatus = "rejected"
)

// PlayerOnlineStatus 陪玩师在线状态
type PlayerOnlineStatus string

const (
	PlayerOnlineStatusOnline  PlayerOnlineStatus = "online"  // 在线
	PlayerOnlineStatusOffline PlayerOnlineStatus = "offline" // 离线
	PlayerOnlineStatusBusy    PlayerOnlineStatus = "busy"    // 忙碌
)

// Player is the pro/陪玩 profile bound to a User.
type Player struct {
	Base
	UserID             uint64             `json:"userId" gorm:"column:user_id;uniqueIndex"`
	Nickname           string             `json:"nickname,omitempty" gorm:"size:64"`
	Bio                string             `json:"bio,omitempty" gorm:"type:text"`
	Rank               string             `json:"rank,omitempty" gorm:"size:32"`
	RatingAverage      float32            `json:"ratingAverage" gorm:"column:rating_average;default:0;check:rating_average >= 0 AND rating_average <= 5" validate:"min=0,max=5"`
	RatingCount        uint32             `json:"ratingCount" gorm:"column:rating_count;default:0"`
	HourlyRateCents    int64              `json:"hourlyRateCents" gorm:"column:hourly_rate_cents"`
	MainGameID         uint64             `json:"mainGameId,omitempty" gorm:"column:main_game_id;index"`
	VerificationStatus VerificationStatus `json:"verificationStatus" gorm:"column:verification_status;size:32;index"`

	// 在线状态相关字段
	OnlineStatus    PlayerOnlineStatus `json:"onlineStatus" gorm:"column:online_status;size:32;index;default:'offline'"` // 在线状态
	AcceptingOrders bool               `json:"acceptingOrders" gorm:"column:accepting_orders;index;default:false"`       // 接单开关
	LastOnlineAt    *time.Time         `json:"lastOnlineAt,omitempty" gorm:"column:last_online_at"`                      // 最后在线时间

	// 审核相关字段
	VerifiedAt   *time.Time `json:"verifiedAt,omitempty" gorm:"column:verified_at"`              // 审核时间
	VerifiedBy   *uint64    `json:"verifiedBy,omitempty" gorm:"column:verified_by"`              // 审核人ID
	VerifyRemark string     `json:"verifyRemark,omitempty" gorm:"column:verify_remark;size:500"` // 审核备注
	RejectReason string     `json:"rejectReason,omitempty" gorm:"column:reject_reason;size:500"` // 拒绝原因
}
