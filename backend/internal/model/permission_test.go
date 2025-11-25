package model_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gamelink/internal/model"
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