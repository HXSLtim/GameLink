package model

import "time"

// RolePermission represents the many-to-many relationship between roles and permissions.
type RolePermission struct {
	RoleID       uint64    `json:"roleId" gorm:"column:role_id;primaryKey;not null"`
	PermissionID uint64    `json:"permissionId" gorm:"column:permission_id;primaryKey;not null"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`

	// 外键关联
	Role       RoleModel  `json:"role,omitempty" gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
	Permission Permission `json:"permission,omitempty" gorm:"foreignKey:PermissionID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for RolePermission.
func (RolePermission) TableName() string {
	return "role_permissions"
}



