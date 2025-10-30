package model

import "time"

// Role defines platform roles for access control.
type Role string

// Role values define platform roles for access control.
const (
	RoleUser   Role = "user"
	RolePlayer Role = "player"
	RoleAdmin  Role = "admin"
)

// UserStatus indicates account state.
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
	Name         string     `json:"name" gorm:"size:64"`
	AvatarURL    string     `json:"avatarUrl,omitempty" gorm:"column:avatar_url;size:255"`
	Role         Role       `json:"role" gorm:"size:32;comment:主要角色（向后兼容）"`
	Status       UserStatus `json:"status" gorm:"size:32;index"`
	LastLoginAt  *time.Time `json:"lastLoginAt,omitempty" gorm:"column:last_login_at"`

	// 多角色支持（新增）
	Roles []RoleModel `json:"roles,omitempty" gorm:"many2many:user_roles;"`
}
