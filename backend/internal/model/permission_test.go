package model_test

import (
	"encoding/json"
	"testing"

	"gamelink/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestPermissionModel(t *testing.T) {
	permission := &model.Permission{
		Base: model.Base{
			ID: 1,
		},
		Method:      model.HTTPMethodPOST,
		Path:        "/api/v1/users",
		Code:        "user.create",
		Group:       "用户管理",
		Description: "允许创建新用户账户",
	}

	assert.Equal(t, uint64(1), permission.ID)
	assert.Equal(t, model.HTTPMethodPOST, permission.Method)
	assert.Equal(t, "/api/v1/users", permission.Path)
	assert.Equal(t, "user.create", permission.Code)
	assert.Equal(t, "用户管理", permission.Group)
	assert.Equal(t, "允许创建新用户账户", permission.Description)
}

func TestPermissionJSONSerialization(t *testing.T) {
	permission := &model.Permission{
		Base: model.Base{
			ID: 1,
		},
		Method:      model.HTTPMethodGET,
		Path:        "/api/v1/orders",
		Code:        "order.manage",
		Group:       "订单管理",
		Description: "允许管理所有订单",
	}

	// 序列化
	data, err := json.Marshal(permission)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "order.manage")
	assert.Contains(t, string(data), "订单管理")

	// 反序列化
	var decoded model.Permission
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, permission.ID, decoded.ID)
	assert.Equal(t, permission.Method, decoded.Method)
	assert.Equal(t, permission.Path, decoded.Path)
	assert.Equal(t, permission.Code, decoded.Code)
	assert.Equal(t, permission.Group, decoded.Group)
	assert.Equal(t, permission.Description, decoded.Description)
}

func TestPermissionTableName(t *testing.T) {
	permission := model.Permission{}
	tableName := permission.TableName()
	assert.Equal(t, "permissions", tableName)
}

func TestPermissionZeroValues(t *testing.T) {
	permission := &model.Permission{
		Method:      "",
		Path:        "",
		Code:        "",
		Group:       "",
		Description: "",
	}

	assert.Equal(t, model.HTTPMethod(""), permission.Method)
	assert.Equal(t, "", permission.Path)
	assert.Equal(t, "", permission.Code)
	assert.Equal(t, "", permission.Group)
	assert.Equal(t, "", permission.Description)
}

func TestPermissionHTTPMethods(t *testing.T) {
	// 测试所有HTTP方法
	methods := []model.HTTPMethod{
		model.HTTPMethodGET,
		model.HTTPMethodPOST,
		model.HTTPMethodPUT,
		model.HTTPMethodPATCH,
		model.HTTPMethodDELETE,
	}

	for _, method := range methods {
		permission := &model.Permission{
			Method: method,
		}
		assert.Equal(t, method, permission.Method)
	}
}

func TestPermissionEdgeCases(t *testing.T) {
	// 测试长文本
	longPath := "/api/v1/very/long/path/that/might/be/used/for/testing/purposes/and/should/handle/long/paths/without/issues"
	longCode := "very.long.permission.code.with.multiple.dots.and.sections.for.testing.purposes"
	longGroup := "这是一个非常长的API分组名称，用于测试字符串长度的边界情况"
	longDescription := "这是一个非常长的权限描述，可以包含很多详细信息，比如权限的具体用途、适用场景、注意事项等等。这种长文本测试可以确保我们的模型能够处理各种长度的输入。"

	permission1 := &model.Permission{
		Method:      model.HTTPMethodPOST,
		Path:        longPath,
		Code:        longCode,
		Group:       longGroup,
		Description: longDescription,
	}

	assert.Equal(t, longPath, permission1.Path)
	assert.Equal(t, longCode, permission1.Code)
	assert.Equal(t, longGroup, permission1.Group)
	assert.Equal(t, longDescription, permission1.Description)

	// 测试特殊字符
	permission2 := &model.Permission{
		Method:      model.HTTPMethodGET,
		Path:        "/api/test@special#path",
		Code:        "permission@123#test",
		Group:       "分组@#$%^&*()",
		Description: "描述包含特殊字符：<>{}[]|\\\"quotes\" and 'apostrophes' and @#$%^&*()_+-=[]{}|;':\",./<>?😊🚀",
	}

	assert.Equal(t, "/api/test@special#path", permission2.Path)
	assert.Equal(t, "permission@123#test", permission2.Code)
	assert.Equal(t, "分组@#$%^&*()", permission2.Group)
	assert.Contains(t, permission2.Description, "特殊字符")
}

func TestPermissionCommonPermissions(t *testing.T) {
	// 测试常见权限
	permissions := []struct {
		method      model.HTTPMethod
		path        string
		code        string
		group       string
		description string
	}{
		{model.HTTPMethodGET, "/api/v1/users", "user.read", "用户管理", "查看用户信息"},
		{model.HTTPMethodPOST, "/api/v1/users", "user.create", "用户管理", "创建新用户"},
		{model.HTTPMethodPUT, "/api/v1/users/:id", "user.update", "用户管理", "更新用户信息"},
		{model.HTTPMethodDELETE, "/api/v1/users/:id", "user.delete", "用户管理", "删除用户"},
		{model.HTTPMethodGET, "/api/v1/orders", "order.read", "订单管理", "查看订单信息"},
		{model.HTTPMethodPOST, "/api/v1/orders", "order.create", "订单管理", "创建新订单"},
		{model.HTTPMethodGET, "/api/v1/players", "player.read", "陪玩师管理", "查看陪玩师信息"},
		{model.HTTPMethodPOST, "/api/v1/players", "player.create", "陪玩师管理", "创建陪玩师"},
	}

	for _, perm := range permissions {
		permission := &model.Permission{
			Method:      perm.method,
			Path:        perm.path,
			Code:        perm.code,
			Group:       perm.group,
			Description: perm.description,
		}
		assert.Equal(t, perm.method, permission.Method)
		assert.Equal(t, perm.path, permission.Path)
		assert.Equal(t, perm.code, permission.Code)
		assert.Equal(t, perm.group, permission.Group)
		assert.Equal(t, perm.description, permission.Description)
	}
}

func TestPermissionWithBaseFields(t *testing.T) {
	permission := &model.Permission{
		Base: model.Base{
			ID: 123,
		},
		Method: model.HTTPMethodGET,
		Path:   "/api/test",
		Code:   "test.permission",
	}

	assert.Equal(t, uint64(123), permission.ID)
	assert.Equal(t, model.HTTPMethodGET, permission.Method)
	assert.Equal(t, "/api/test", permission.Path)
	assert.Equal(t, "test.permission", permission.Code)
}

func TestPermissionJSONFields(t *testing.T) {
	permission := &model.Permission{
		Base: model.Base{
			ID: 1,
		},
		Method:      model.HTTPMethodPOST,
		Path:        "/api/v1/users",
		Code:        "user.create",
		Group:       "用户管理",
		Description: "创建用户",
	}

	data, err := json.Marshal(permission)
	assert.NoError(t, err)

	// 验证JSON结构
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// 检查必需的字段
	assert.Contains(t, result, "id")
	assert.Contains(t, result, "method")
	assert.Contains(t, result, "path")
	assert.Contains(t, result, "code")
	assert.Contains(t, result, "group")
	assert.Contains(t, result, "description")

	// 验证值
	assert.Equal(t, float64(1), result["id"])
	assert.Equal(t, "POST", result["method"])
	assert.Equal(t, "/api/v1/users", result["path"])
	assert.Equal(t, "user.create", result["code"])
	assert.Equal(t, "用户管理", result["group"])
	assert.Equal(t, "创建用户", result["description"])
}

func TestPermissionEmptyFields(t *testing.T) {
	permission := &model.Permission{}

	assert.Equal(t, uint64(0), permission.ID)
	assert.Equal(t, model.HTTPMethod(""), permission.Method)
	assert.Equal(t, "", permission.Path)
	assert.Equal(t, "", permission.Code)
	assert.Equal(t, "", permission.Group)
	assert.Equal(t, "", permission.Description)
}

func TestPermissionMultilingualContent(t *testing.T) {
	// 测试多语言内容
	permissions := []struct {
		method      model.HTTPMethod
		path        string
		code        string
		group       string
		description string
		lang        string
	}{
		{model.HTTPMethodGET, "/api/v1/users", "user.read", "User Management", "Allow reading user information", "English"},
		{model.HTTPMethodPOST, "/api/v1/users", "用户.创建", "用户管理", "允许创建新用户", "Chinese"},
		{model.HTTPMethodGET, "/api/v1/utilisateurs", "utilisateur.lire", "Gestion des utilisateurs", "Permet de lire les informations utilisateur", "French"},
		{model.HTTPMethodPOST, "/api/v1/usuarios", "usuario.crear", "Gestión de usuarios", "Permite crear nuevos usuarios", "Spanish"},
		{model.HTTPMethodGET, "/api/v1/benutzer", "benutzer.lesen", "Benutzerverwaltung", "Ermöglicht das Lesen von Benutzerinformationen", "German"},
	}

	for _, perm := range permissions {
		permission := &model.Permission{
			Method:      perm.method,
			Path:        perm.path,
			Code:        perm.code,
			Group:       perm.group,
			Description: perm.description,
		}
		assert.Equal(t, perm.method, permission.Method, "Failed for language: %s", perm.lang)
		assert.Equal(t, perm.path, permission.Path, "Failed for language: %s", perm.lang)
		assert.Equal(t, perm.code, permission.Code, "Failed for language: %s", perm.lang)
		assert.Equal(t, perm.group, permission.Group, "Failed for language: %s", perm.lang)
		assert.Equal(t, perm.description, permission.Description, "Failed for language: %s", perm.lang)
	}
}

func TestPermissionConstants(t *testing.T) {
	// 测试HTTP方法常量
	assert.Equal(t, model.HTTPMethod("GET"), model.HTTPMethodGET)
	assert.Equal(t, model.HTTPMethod("POST"), model.HTTPMethodPOST)
	assert.Equal(t, model.HTTPMethod("PUT"), model.HTTPMethodPUT)
	assert.Equal(t, model.HTTPMethod("PATCH"), model.HTTPMethodPATCH)
	assert.Equal(t, model.HTTPMethod("DELETE"), model.HTTPMethodDELETE)
}

// TestPermissionValidateCode tests the permission code validation
// Requirements: 1.1, 1.3
func TestPermissionValidateCode(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		// Valid codes
		{"valid three-part code", "admin.users.create", true},
		{"valid with numbers", "admin2.users3.read4", true},
		{"valid lowercase", "module.resource.action", true},

		// Invalid codes
		{"empty code", "", false},
		{"single part", "admin", false},
		{"two parts", "admin.users", false},
		{"four parts", "admin.users.create.extra", false},
		{"uppercase", "Admin.Users.Create", false},
		{"starts with number", "1admin.users.create", false},
		{"contains special chars", "admin.users.create!", false},
		{"contains spaces", "admin .users.create", false},
		{"contains underscore", "admin_module.users.create", false},
		{"contains hyphen", "admin-module.users.create", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			permission := &model.Permission{Code: tt.code}
			result := permission.ValidateCode()
			assert.Equal(t, tt.expected, result, "ValidateCode() for code %q", tt.code)
		})
	}
}

// TestPermissionCodeParts tests the GetModule, GetResource, GetAction methods
// Requirements: 1.1
func TestPermissionCodeParts(t *testing.T) {
	tests := []struct {
		name             string
		code             string
		expectedModule   string
		expectedResource string
		expectedAction   string
	}{
		{
			name:             "valid three-part code",
			code:             "admin.users.create",
			expectedModule:   "admin",
			expectedResource: "users",
			expectedAction:   "create",
		},
		{
			name:             "empty code",
			code:             "",
			expectedModule:   "",
			expectedResource: "",
			expectedAction:   "",
		},
		{
			name:             "single part",
			code:             "admin",
			expectedModule:   "admin",
			expectedResource: "",
			expectedAction:   "",
		},
		{
			name:             "two parts",
			code:             "admin.users",
			expectedModule:   "admin",
			expectedResource: "users",
			expectedAction:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			permission := &model.Permission{Code: tt.code}
			assert.Equal(t, tt.expectedModule, permission.GetModule())
			assert.Equal(t, tt.expectedResource, permission.GetResource())
			assert.Equal(t, tt.expectedAction, permission.GetAction())
		})
	}
}

// TestPermissionIsSystemPermission tests the IsSystemPermission method
// Requirements: 1.5
func TestPermissionIsSystemPermission(t *testing.T) {
	tests := []struct {
		name     string
		isSystem bool
		expected bool
	}{
		{"system permission", true, true},
		{"non-system permission", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			permission := &model.Permission{IsSystem: tt.isSystem}
			assert.Equal(t, tt.expected, permission.IsSystemPermission())
		})
	}
}

// TestPermissionTreeStructure tests the tree structure fields
// Requirements: 1.1
func TestPermissionTreeStructure(t *testing.T) {
	parentID := uint64(1)

	// Test permission with parent
	childPermission := &model.Permission{
		Base:      model.Base{ID: 2},
		Code:      "admin.users.create",
		ParentID:  &parentID,
		SortOrder: 1,
	}

	assert.NotNil(t, childPermission.ParentID)
	assert.Equal(t, uint64(1), *childPermission.ParentID)
	assert.Equal(t, 1, childPermission.SortOrder)

	// Test root permission (no parent)
	rootPermission := &model.Permission{
		Base:      model.Base{ID: 1},
		Code:      "admin.users.read",
		ParentID:  nil,
		SortOrder: 0,
	}

	assert.Nil(t, rootPermission.ParentID)
	assert.Equal(t, 0, rootPermission.SortOrder)
}

// TestPermissionToTreeNode tests the ToTreeNode conversion
// Requirements: 1.1
func TestPermissionToTreeNode(t *testing.T) {
	parentID := uint64(1)
	permission := &model.Permission{
		Base:        model.Base{ID: 2},
		Method:      model.HTTPMethodPOST,
		Path:        "/api/admin/users",
		Code:        "admin.users.create",
		Group:       "用户管理",
		Description: "创建用户",
		ParentID:    &parentID,
		SortOrder:   1,
		IsSystem:    true,
	}

	node := permission.ToTreeNode()

	assert.Equal(t, uint64(2), node.ID)
	assert.Equal(t, "admin.users.create", node.Code)
	assert.Equal(t, "创建用户", node.Description)
	assert.Equal(t, "用户管理", node.Group)
	assert.Equal(t, model.HTTPMethodPOST, node.Method)
	assert.Equal(t, "/api/admin/users", node.Path)
	assert.NotNil(t, node.ParentID)
	assert.Equal(t, uint64(1), *node.ParentID)
	assert.Equal(t, 1, node.SortOrder)
	assert.True(t, node.IsSystem)
	assert.NotNil(t, node.Children)
	assert.Len(t, node.Children, 0)
}

// TestBuildPermissionTree tests the BuildPermissionTree function
// Requirements: 1.1
func TestBuildPermissionTree(t *testing.T) {
	t.Run("empty permissions", func(t *testing.T) {
		result := model.BuildPermissionTree([]model.Permission{})
		assert.Nil(t, result)
	})

	t.Run("single root permission", func(t *testing.T) {
		permissions := []model.Permission{
			{Base: model.Base{ID: 1}, Code: "admin.users.read", ParentID: nil},
		}
		result := model.BuildPermissionTree(permissions)
		assert.Len(t, result, 1)
		assert.Equal(t, uint64(1), result[0].ID)
		assert.Len(t, result[0].Children, 0)
	})

	t.Run("parent-child relationship", func(t *testing.T) {
		parentID := uint64(1)
		permissions := []model.Permission{
			{Base: model.Base{ID: 1}, Code: "admin.users.read", ParentID: nil},
			{Base: model.Base{ID: 2}, Code: "admin.users.create", ParentID: &parentID},
		}
		result := model.BuildPermissionTree(permissions)
		assert.Len(t, result, 1)
		assert.Equal(t, uint64(1), result[0].ID)
		assert.Len(t, result[0].Children, 1)
		assert.Equal(t, uint64(2), result[0].Children[0].ID)
	})

	t.Run("orphan child becomes root", func(t *testing.T) {
		nonExistentParentID := uint64(999)
		permissions := []model.Permission{
			{Base: model.Base{ID: 1}, Code: "admin.users.read", ParentID: nil},
			{Base: model.Base{ID: 2}, Code: "admin.users.create", ParentID: &nonExistentParentID},
		}
		result := model.BuildPermissionTree(permissions)
		// Both should be roots since parent 999 doesn't exist
		assert.Len(t, result, 2)
	})
}

// TestBuildPermissionTreeByGroup tests the BuildPermissionTreeByGroup function
// Requirements: 1.1
func TestBuildPermissionTreeByGroup(t *testing.T) {
	t.Run("empty permissions", func(t *testing.T) {
		result := model.BuildPermissionTreeByGroup([]model.Permission{})
		assert.Nil(t, result)
	})

	t.Run("single group", func(t *testing.T) {
		permissions := []model.Permission{
			{Base: model.Base{ID: 1}, Code: "admin.users.read", Group: "用户管理"},
			{Base: model.Base{ID: 2}, Code: "admin.users.create", Group: "用户管理"},
		}
		result := model.BuildPermissionTreeByGroup(permissions)
		assert.Len(t, result, 1)
		assert.Equal(t, "用户管理", result[0].Group)
		assert.Len(t, result[0].Permissions, 2)
	})

	t.Run("multiple groups", func(t *testing.T) {
		permissions := []model.Permission{
			{Base: model.Base{ID: 1}, Code: "admin.users.read", Group: "用户管理"},
			{Base: model.Base{ID: 2}, Code: "admin.orders.read", Group: "订单管理"},
		}
		result := model.BuildPermissionTreeByGroup(permissions)
		assert.Len(t, result, 2)
		// Groups should be in order of first appearance
		assert.Equal(t, "用户管理", result[0].Group)
		assert.Equal(t, "订单管理", result[1].Group)
	})
}
