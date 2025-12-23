package model

import "time"

// Role defines platform roles for access control.
// @Enum user, player, admin
type Role string

// Role values define platform roles for access control.
const (
	RoleUser   Role = "user"
	RolePlayer Role = "player"
	RoleAdmin  Role = "admin"
)

// UserStatus indicates account state.
// @Enum active, suspended, banned
type UserStatus string

// UserStatus values indicate account state.
const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusBanned    UserStatus = "banned"
)

// User represents a platform account.
type User struct {
	Base
	Phone        string     `json:"phone,omitempty" gorm:"size:32;uniqueIndex"`
	Email        string     `json:"email,omitempty" gorm:"size:128;uniqueIndex"`
	PasswordHash string     `json:"-" gorm:"column:password_hash;size:255"`
	Name         string     `json:"name" gorm:"size:64;index"` // 添加索引，用于搜索
	AvatarURL    string     `json:"avatarUrl,omitempty" gorm:"column:avatar_url;size:255"`
	Role         Role       `json:"role" gorm:"size:32;comment:主要角色（向后兼容）"`
	Status       UserStatus `json:"status" gorm:"size:32;index;index:idx_status_last_login,priority:1"`                       // 复合索引第一部分
	LastLoginAt  *time.Time `json:"lastLoginAt,omitempty" gorm:"column:last_login_at;index:idx_status_last_login,priority:2"` // 复合索引第二部分

	// 多角色支持（新增）
	Roles []RoleModel `json:"roles,omitempty" gorm:"many2many:user_roles;"`

	// 用户钱包（积分就是余额，从钱包读取）
	Wallet *Wallet `json:"wallet,omitempty" gorm:"foreignKey:UserID"`

	// VIP 相关字段
	VipLevelID          *uint64    `json:"vipLevelId,omitempty" gorm:"column:vip_level_id;index"`              // 当前VIP等级ID
	VipUnlocked         bool       `json:"vipUnlocked" gorm:"column:vip_unlocked;default:false"`               // VIP是否已解锁
	VipExp              int64      `json:"vipExp" gorm:"column:vip_exp;default:0"`                             // VIP经验（累计消费，分）
	TotalRechargeCents  int64      `json:"totalRechargeCents" gorm:"column:total_recharge_cents;default:0"`    // 累计充值（分）
	VipUnlockedAt       *time.Time `json:"vipUnlockedAt,omitempty" gorm:"column:vip_unlocked_at"`              // VIP解锁时间
	VipExpireAt         *time.Time `json:"vipExpireAt,omitempty" gorm:"column:vip_expire_at"`                  // VIP过期时间（nil=永久）
	LastMonthlyCouponAt *time.Time `json:"lastMonthlyCouponAt,omitempty" gorm:"column:last_monthly_coupon_at"` // 上次发放月度券时间

	// VIP 等级关联
	VipLevel *VipLevel `json:"vipLevel,omitempty" gorm:"foreignKey:VipLevelID"`
}
