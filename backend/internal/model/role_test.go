package model_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gamelink/internal/model"
)

func TestRoleModel(t *testing.T) {
	permission1 := model.Permission{Base: model.Base{ID: 1}, Code: "user.create", Description: "创建用户"}
	permission2 := model.Permission{Base: model.Base{ID: 2}, Code: "user.read", Description: "查看用户"}
	user1 := model.User{Base: model.Base{ID: 100}, Name: "用户1"}
	user2 := model.User{Base: model.Base{ID: 101}, Name: "用户2"}

	role := &model.RoleModel{
		Base: model.Base{
			ID: 1,
		},
		Slug:        "admin",
		Name:        "管理员",
		Description: "系统管理员角色，拥有所有权限",
		IsSystem:    true,
		Permissions: []model.Permission{permission1, permission2},
		Users:       []model.User{user1, user2},
	}

	assert.Equal(t, uint64(1), role.ID)
	assert.Equal(t, "admin", role.Slug)
	assert.Equal(t, "管理员", role.Name)
	assert.Equal(t, "系统管理员角色，拥有所有权限", role.Description)
	assert.True(t, role.IsSystem)
	assert.Len(t, role.Permissions, 2)
	assert.Equal(t, "user.create", role.Permissions[0].Code)
	assert.Equal(t, "user.read", role.Permissions[1].Code)
	assert.Len(t, role.Users, 2)
	assert.Equal(t, "用户1", role.Users[0].Name)
	assert.Equal(t, "用户2", role.Users[1].Name)
}

func TestRoleJSONSerialization(t *testing.T) {
	permission := model.Permission{Base: model.Base{ID: 1}, Code: "test.permission", Description: "测试权限"}
	user := model.User{Base: model.Base{ID: 100}, Name: "测试用户"}

	role := &model.RoleModel{
		Base: model.Base{
			ID: 1,
		},
		Slug:        "test_role",
		Name:        "测试角色",
		Description: "这是一个测试角色",
		IsSystem:    false,
		Permissions: []model.Permission{permission},
		Users:       []model.User{user},
	}

	// 序列化
	data, err := json.Marshal(role)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "test_role")
	assert.Contains(t, string(data), "测试角色")

	// 反序列化
	var decoded model.RoleModel
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, role.ID, decoded.ID)
	assert.Equal(t, role.Slug, decoded.Slug)
	assert.Equal(t, role.Name, decoded.Name)
	assert.Equal(t, role.Description, decoded.Description)
	assert.Equal(t, role.IsSystem, decoded.IsSystem)
}

func TestRoleTableName(t *testing.T) {
	role := model.RoleModel{}
	tableName := role.TableName()
	assert.Equal(t, "roles", tableName)
}

func TestRoleConstants(t *testing.T) {
	// 测试角色标识符常量
	assert.Equal(t, model.RoleSlug("super_admin"), model.RoleSlugSuperAdmin)
	assert.Equal(t, model.RoleSlug("admin"), model.RoleSlugAdmin)
	assert.Equal(t, model.RoleSlug("player"), model.RoleSlugPlayer)
	assert.Equal(t, model.RoleSlug("user"), model.RoleSlugUser)
}

func TestRoleIsSystemRole(t *testing.T) {
	// 测试系统角色
	role1 := &model.RoleModel{
		IsSystem: true,
	}
	assert.True(t, role1.IsSystemRole())

	// 测试非系统角色
	role2 := &model.RoleModel{
		IsSystem: false,
	}
	assert.False(t, role2.IsSystemRole())
}

func TestRoleIsSuperAdmin(t *testing.T) {
	// 测试超级管理员角色
	role1 := &model.RoleModel{
		Slug: "super_admin",
	}
	assert.True(t, role1.IsSuperAdmin())

	// 测试普通管理员角色
	role2 := &model.RoleModel{
		Slug: "admin",
	}
	assert.False(t, role2.IsSuperAdmin())

	// 测试其他角色
	role3 := &model.RoleModel{
		Slug: "user",
	}
	assert.False(t, role3.IsSuperAdmin())

	role4 := &model.RoleModel{
		Slug: "player",
	}
	assert.False(t, role4.IsSuperAdmin())

	// 测试空角色标识符
	role5 := &model.RoleModel{
		Slug: "",
	}
	assert.False(t, role5.IsSuperAdmin())
}

func TestRoleZeroValues(t *testing.T) {
	role := &model.RoleModel{
		Slug:        "",
		Name:        "",
		Description: "",
		IsSystem:    false,
		Permissions: nil,
		Users:       nil,
	}

	assert.Equal(t, "", role.Slug)
	assert.Equal(t, "", role.Name)
	assert.Equal(t, "", role.Description)
	assert.False(t, role.IsSystem)
	assert.Nil(t, role.Permissions)
	assert.Nil(t, role.Users)
}

func TestRoleCommonRoles(t *testing.T) {
	// 测试常见角色
	roles := []struct {
		slug        string
		name        string
		description string
		isSystem    bool
	}{
		{"super_admin", "超级管理员", "拥有系统所有权限的超级管理员角色", true},
		{"admin", "管理员", "系统管理员，拥有大部分管理权限", true},
		{"player_manager", "陪玩师管理员", "专门管理陪玩师的管理员", true},
		{"order_manager", "订单管理员", "专门管理订单的管理员", true},
		{"finance_manager", "财务管理员", "专门管理财务的管理员", true},
		{"customer_service", "客服", "处理用户问题和争议的角色", true},
		{"player", "陪玩师", "提供游戏陪玩服务的用户角色", true},
		{"user", "普通用户", "使用平台服务的普通用户角色", true},
		{"vip_user", "VIP用户", "高级用户，享有特殊权益", false},
		{"trial_user", "试用用户", "正在试用平台的用户", false},
	}

	for _, roleData := range roles {
		role := &model.RoleModel{
			Slug:        roleData.slug,
			Name:        roleData.name,
			Description: roleData.description,
			IsSystem:    roleData.isSystem,
		}
		assert.Equal(t, roleData.slug, role.Slug)
		assert.Equal(t, roleData.name, role.Name)
		assert.Equal(t, roleData.description, role.Description)
		assert.Equal(t, roleData.isSystem, role.IsSystem)
	}
}

func TestRoleWithPermissionsAndUsers(t *testing.T) {
	// 创建权限
	perms := []model.Permission{
		{Base: model.Base{ID: 1}, Code: "user.create", Description: "创建用户"},
		{Base: model.Base{ID: 2}, Code: "user.read", Description: "查看用户"},
		{Base: model.Base{ID: 3}, Code: "user.update", Description: "更新用户"},
		{Base: model.Base{ID: 4}, Code: "user.delete", Description: "删除用户"},
	}

	// 创建用户
	users := []model.User{
		{Base: model.Base{ID: 100}, Name: "张三"},
		{Base: model.Base{ID: 101}, Name: "李四"},
		{Base: model.Base{ID: 102}, Name: "王五"},
	}

	role := &model.RoleModel{
		Slug:        "user_manager",
		Name:        "用户管理员",
		Description: "专门管理用户的角色",
		Permissions: perms,
		Users:       users,
	}

	assert.Len(t, role.Permissions, 4)
	assert.Len(t, role.Users, 3)
	assert.Equal(t, "user.create", role.Permissions[0].Code)
	assert.Equal(t, "张三", role.Users[0].Name)
}

func TestRoleJSONWithRelations(t *testing.T) {
	permission := model.Permission{Base: model.Base{ID: 1}, Code: "test.permission", Description: "测试权限"}
	user := model.User{Base: model.Base{ID: 100}, Name: "测试用户"}

	role := &model.RoleModel{
		Base: model.Base{
			ID: 1,
		},
		Slug:        "test_role",
		Name:        "测试角色",
		Description: "这是一个测试角色",
		Permissions: []model.Permission{permission},
		Users:       []model.User{user},
	}

	data, err := json.Marshal(role)
	assert.NoError(t, err)

	// 验证JSON结构包含关联数据
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// 检查关联字段
	assert.Contains(t, result, "permissions")
	assert.Contains(t, result, "users")

	// 验证权限数据
	permissions := result["permissions"].([]interface{})
	assert.Len(t, permissions, 1)
	perm := permissions[0].(map[string]interface{})
	assert.Equal(t, "test.permission", perm["code"])
	assert.Equal(t, "测试权限", perm["description"])

	// 验证用户数据
	users := result["users"].([]interface{})
	assert.Len(t, users, 1)
	usr := users[0].(map[string]interface{})
	assert.Equal(t, "测试用户", usr["name"])
}

func TestRoleEmptyRelations(t *testing.T) {
	role := &model.RoleModel{
		Slug:        "test_role",
		Name:        "测试角色",
		Permissions: []model.Permission{},
		Users:       []model.User{},
	}

	assert.NotNil(t, role.Permissions)
	assert.NotNil(t, role.Users)
	assert.Len(t, role.Permissions, 0)
	assert.Len(t, role.Users, 0)

	// 测试序列化空关联
	data, err := json.Marshal(role)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	assert.Contains(t, result, "permissions")
	assert.Contains(t, result, "users")
	
	permissions := result["permissions"].([]interface{})
	users := result["users"].([]interface{})
	assert.Len(t, permissions, 0)
	assert.Len(t, users, 0)
}