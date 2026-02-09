package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost+2) // cost=12 for better security
}

// CheckPassword compares a password with its hash
func CheckPassword(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}

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

// LoginType 登录方式
type LoginType string

const (
	LoginTypePassword LoginType = "password" // 密码登录
	LoginTypeSMS      LoginType = "sms"      // 短信验证码登录
	LoginTypeEmail    LoginType = "email"    // 邮箱验证码登录
	LoginTypeOAuth    LoginType = "oauth"    // 第三方OAuth登录（预留）
)

// User represents a platform account.
type User struct {
	Base
	Phone        string     `json:"phone,omitempty" gorm:"size:32;uniqueIndex"`
	Email        string     `json:"email,omitempty" gorm:"size:128;uniqueIndex"`
	PasswordHash string     `json:"-" gorm:"column:password_hash;size:255"`
	Name         string     `json:"name" gorm:"size:64;index"`     // 真实姓名（身份识别）
	Nickname     string     `json:"nickname" gorm:"size:64;index"` // 昵称（社交展示）
	AvatarURL    string     `json:"avatarUrl,omitempty" gorm:"column:avatar_url;size:255"`
	Role         Role       `json:"role" gorm:"size:32;comment:主要角色（向后兼容）"`
	Status       UserStatus `json:"status" gorm:"size:32;index;index:idx_status_last_login,priority:1"`                       // 复合索引第一部分
	LastLoginAt  *time.Time `json:"lastLoginAt,omitempty" gorm:"column:last_login_at;index:idx_status_last_login,priority:2"` // 复合索引第二部分

	// 封禁相关字段
	BanReason string     `json:"banReason,omitempty" gorm:"column:ban_reason;size:500"` // 封禁原因
	BannedAt  *time.Time `json:"bannedAt,omitempty" gorm:"column:banned_at"`            // 封禁时间
	BannedBy  *uint64    `json:"bannedBy,omitempty" gorm:"column:banned_by;index"`      // 封禁操作人ID

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

	// 微信小程序字段
	// 注意：使用 *string 允许 NULL 值，条件唯一索引只对非 NULL 值生效
	WeChatOpenID  *string `json:"-" gorm:"column:wechat_open_id;size:64;uniqueIndex:,where:wechat_open_id IS NOT NULL"` // 微信 OpenID
	WeChatUnionID string  `json:"-" gorm:"column:wechat_union_id;size:64;index"`                                        // 微信 UnionID（跨应用）
}
