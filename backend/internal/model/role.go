package model

// RoleSlug defines system predefined role identifiers.
type RoleSlug string

// RoleSlug values define system predefined roles.
const (
	RoleSlugSuperAdmin RoleSlug = "super_admin"
	RoleSlugAdmin      RoleSlug = "admin"
	RoleSlugPlayer     RoleSlug = "player"
	RoleSlugUser       RoleSlug = "user"
)

// RoleModel represents a platform role with custom permissions.
// Supports both system-predefined roles and user-defined custom roles.
type RoleModel struct {
	Base
	Slug        string `json:"slug" gorm:"size:64;uniqueIndex;not null;comment:角色标识"`
	Name        string `json:"name" gorm:"size:128;not null;comment:角色名称"`
	Description string `json:"description" gorm:"size:255;comment:角色描述"`
	IsSystem    bool   `json:"isSystem" gorm:"column:is_system;default:false;comment:是否系统角色（不可删除）"`

	// 关联关系
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
	Users       []User       `json:"users,omitempty" gorm:"many2many:user_roles;"`
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

