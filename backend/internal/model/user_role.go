package model

import "time"

// UserRole represents the many-to-many relationship between users and roles.
// One user can have multiple roles.
type UserRole struct {
	UserID    uint64    `json:"userId" gorm:"column:user_id;primaryKey;not null"`
	RoleID    uint64    `json:"roleId" gorm:"column:role_id;primaryKey;not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`

	// 外键关联
	User User      `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Role RoleModel `json:"role,omitempty" gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for UserRole.
func (UserRole) TableName() string {
	return "user_roles"
}

