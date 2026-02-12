package model

import "errors"

// RoleSlug defines system predefined role identifiers.
type RoleSlug string

// RoleSlug values define system predefined roles.
const (
	RoleSlugSuperAdmin      RoleSlug = "superAdmin"      // 超级管理员
	RoleSlugAdmin           RoleSlug = "admin"           // 管理员/店长
	RoleSlugFinance         RoleSlug = "finance"         // 财务
	RoleSlugCustomerService RoleSlug = "customerService" // 客服
	RoleSlugCSLeader        RoleSlug = "csLeader"        // 客服主管
	RoleSlugCSAgent         RoleSlug = "csAgent"         // 普通客服
	RoleSlugPlayer          RoleSlug = "player"          // 陪玩师
	RoleSlugUser            RoleSlug = "user"            // 普通用户
)

// MaxRoleInheritanceDepth defines the maximum allowed inheritance depth for roles.
const MaxRoleInheritanceDepth = 5

// Role inheritance errors
var (
	ErrRoleMaxDepthExceeded    = errors.New("role inheritance depth exceeds maximum limit")
	ErrRoleCircularInheritance = errors.New("circular role inheritance detected")
)

// RoleModel represents a platform role with custom permissions.
// Supports both system-predefined roles and user-defined custom roles.
// Supports role inheritance via ParentID for hierarchical permission management.
type RoleModel struct {
	Base
	Slug        string  `json:"slug" gorm:"size:64;uniqueIndex;not null;comment:角色标识"`
	Name        string  `json:"name" gorm:"size:128;not null;comment:角色名称"`
	Description string  `json:"description" gorm:"size:255;comment:角色描述"`
	IsSystem    bool    `json:"isSystem" gorm:"column:is_system;default:false;comment:是否系统角色（不可删除）"`
	ParentID    *uint64 `json:"parentId" gorm:"index;comment:父角色ID，用于继承"`
	Priority    int     `json:"priority" gorm:"default:0;comment:优先级，数值越大优先级越高，用于权限冲突解决"`
	Level       int     `json:"level" gorm:"default:0;comment:继承层级，根角色为0，最大5层"`

	// Self-referential relationship for inheritance
	Parent   *RoleModel  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []RoleModel `json:"children,omitempty" gorm:"foreignKey:ParentID"`

	// 关联关系
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;foreignKey:ID;joinForeignKey:RoleID;References:ID;joinReferences:PermissionID"`
	Users       []User       `json:"users" gorm:"many2many:user_roles;foreignKey:ID;joinForeignKey:RoleID;References:ID;joinReferences:UserID"`
}

// TableName specifies the table name for RoleModel.
func (RoleModel) TableName() string {
	return "roles"
}

// IsSystemRole checks if this is a system-predefined role.
func (r *RoleModel) IsSystemRole() bool {
	return r.IsSystem
}

// IsSuperAdmin checks if this role is super_admin.
func (r *RoleModel) IsSuperAdmin() bool {
	return r.Slug == string(RoleSlugSuperAdmin)
}

// HasParent checks if this role has a parent role.
func (r *RoleModel) HasParent() bool {
	return r.ParentID != nil && *r.ParentID > 0
}

// ValidateInheritanceDepth checks if the role's inheritance level is within limits.
func (r *RoleModel) ValidateInheritanceDepth() error {
	if r.Level > MaxRoleInheritanceDepth {
		return ErrRoleMaxDepthExceeded
	}
	return nil
}

// CalculateLevel calculates the inheritance level based on parent.
// Root roles have level 0, children have parent.Level + 1.
func (r *RoleModel) CalculateLevel(parentLevel int) int {
	if r.ParentID == nil {
		return 0
	}
	return parentLevel + 1
}
