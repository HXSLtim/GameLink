package model

import (
	"regexp"
	"strings"
)

// HTTPMethod defines standard HTTP methods for permission control.
type HTTPMethod string

// HTTPMethod values define standard HTTP methods.
const (
	HTTPMethodGET    HTTPMethod = "GET"
	HTTPMethodPOST   HTTPMethod = "POST"
	HTTPMethodPUT    HTTPMethod = "PUT"
	HTTPMethodPATCH  HTTPMethod = "PATCH"
	HTTPMethodDELETE HTTPMethod = "DELETE"
)

// PermissionCodePattern defines the regex pattern for permission code validation.
// Format: module.resource.action (three dot-separated segments)
var PermissionCodePattern = regexp.MustCompile(`^[a-z][a-z0-9]*\.[a-z][a-z0-9]*\.[a-z][a-z0-9]*$`)

// Permission represents a backend API resource that can be accessed.
// It records method + path as a unique identifier for fine-grained authorization.
// Supports tree structure via ParentID for hierarchical permission organization.
type Permission struct {
	Base
	Method      HTTPMethod `json:"method" gorm:"size:16;not null;uniqueIndex:idx_method_path"`
	Path        string     `json:"path" gorm:"size:255;not null;uniqueIndex:idx_method_path"`
	Code        string     `json:"code" gorm:"size:128;uniqueIndex;comment:语义化标识，如 admin.games.read"`
	Group       string     `json:"group" gorm:"size:64;index;comment:API 分组，如 /admin/games"`
	Description string     `json:"description" gorm:"size:255"`
	ParentID    *uint64    `json:"parentId" gorm:"index;comment:父权限ID，用于树形结构"`
	SortOrder   int        `json:"sortOrder" gorm:"default:0;comment:排序顺序"`
	IsSystem    bool       `json:"isSystem" gorm:"default:false;comment:是否系统权限（不可删除）"`

	// Self-referential relationship for tree structure
	Parent   *Permission  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Permission `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

// TableName specifies the table name for Permission.
func (Permission) TableName() string {
	return "permissions"
}

// ValidateCode validates the permission code format.
// Returns true if the code matches the pattern: module.resource.action
func (p *Permission) ValidateCode() bool {
	if p.Code == "" {
		return false
	}
	return PermissionCodePattern.MatchString(p.Code)
}

// GetModule returns the module part of the permission code.
func (p *Permission) GetModule() string {
	parts := strings.Split(p.Code, ".")
	if len(parts) >= 1 {
		return parts[0]
	}
	return ""
}

// GetResource returns the resource part of the permission code.
func (p *Permission) GetResource() string {
	parts := strings.Split(p.Code, ".")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}

// GetAction returns the action part of the permission code.
func (p *Permission) GetAction() string {
	parts := strings.Split(p.Code, ".")
	if len(parts) >= 3 {
		return parts[2]
	}
	return ""
}

// IsSystemPermission checks if this is a system permission that cannot be deleted.
func (p *Permission) IsSystemPermission() bool {
	return p.IsSystem
}

// PermissionTreeNode represents a node in the permission tree structure.
// Used for hierarchical display of permissions in the frontend.
type PermissionTreeNode struct {
	ID          uint64                `json:"id"`
	Code        string                `json:"code"`
	Description string                `json:"description"`
	Group       string                `json:"group"`
	Method      HTTPMethod            `json:"method"`
	Path        string                `json:"path"`
	ParentID    *uint64               `json:"parentId,omitempty"`
	SortOrder   int                   `json:"sortOrder"`
	IsSystem    bool                  `json:"isSystem"`
	Children    []*PermissionTreeNode `json:"children,omitempty"`
}

// PermissionGroup represents a group of permissions organized by module.
type PermissionGroup struct {
	Group       string                `json:"group"`
	Permissions []*PermissionTreeNode `json:"permissions"`
}

// ToTreeNode converts a Permission to a PermissionTreeNode.
func (p *Permission) ToTreeNode() *PermissionTreeNode {
	return &PermissionTreeNode{
		ID:          p.ID,
		Code:        p.Code,
		Description: p.Description,
		Group:       p.Group,
		Method:      p.Method,
		Path:        p.Path,
		ParentID:    p.ParentID,
		SortOrder:   p.SortOrder,
		IsSystem:    p.IsSystem,
		Children:    make([]*PermissionTreeNode, 0),
	}
}

// BuildPermissionTree builds a tree structure from a flat list of permissions.
// It organizes permissions by their parent-child relationships.
func BuildPermissionTree(permissions []Permission) []*PermissionTreeNode {
	if len(permissions) == 0 {
		return nil
	}

	// Create a map for quick lookup
	nodeMap := make(map[uint64]*PermissionTreeNode)
	var roots []*PermissionTreeNode

	// First pass: create all nodes
	for i := range permissions {
		node := permissions[i].ToTreeNode()
		nodeMap[permissions[i].ID] = node
	}

	// Second pass: build tree structure
	for i := range permissions {
		node := nodeMap[permissions[i].ID]
		if permissions[i].ParentID == nil || *permissions[i].ParentID == 0 {
			roots = append(roots, node)
		} else {
			if parent, ok := nodeMap[*permissions[i].ParentID]; ok {
				parent.Children = append(parent.Children, node)
			} else {
				// Parent not found, treat as root
				roots = append(roots, node)
			}
		}
	}

	return roots
}

// BuildPermissionTreeByGroup builds a tree structure grouped by permission group.
func BuildPermissionTreeByGroup(permissions []Permission) []PermissionGroup {
	if len(permissions) == 0 {
		return nil
	}

	// Group permissions by their group field
	groupMap := make(map[string][]Permission)
	groupOrder := make([]string, 0)

	for _, p := range permissions {
		if _, exists := groupMap[p.Group]; !exists {
			groupOrder = append(groupOrder, p.Group)
		}
		groupMap[p.Group] = append(groupMap[p.Group], p)
	}

	// Build tree for each group
	result := make([]PermissionGroup, 0, len(groupOrder))
	for _, group := range groupOrder {
		perms := groupMap[group]
		tree := BuildPermissionTree(perms)
		result = append(result, PermissionGroup{
			Group:       group,
			Permissions: tree,
		})
	}

	return result
}
