package model_test

import (
	"encoding/json"
	"testing"

	"gamelink/internal/model"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, model.RoleSlug("superAdmin"), model.RoleSlugSuperAdmin)
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
		Slug: "superAdmin",
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
		{"superAdmin", "超级管理员", "拥有系统所有权限的超级管理员角色", true},
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

// TestRoleHasParent tests the HasParent method
// Requirements: 10.2
func TestRoleHasParent(t *testing.T) {
	tests := []struct {
		name     string
		parentID *uint64
		expected bool
	}{
		{"nil parent", nil, false},
		{"zero parent", func() *uint64 { v := uint64(0); return &v }(), false},
		{"valid parent", func() *uint64 { v := uint64(1); return &v }(), true},
		{"large parent ID", func() *uint64 { v := uint64(999999); return &v }(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role := &model.RoleModel{ParentID: tt.parentID}
			assert.Equal(t, tt.expected, role.HasParent())
		})
	}
}

// TestRoleCalculateLevel tests the CalculateLevel method
// Requirements: 10.2
func TestRoleCalculateLevel(t *testing.T) {
	tests := []struct {
		name          string
		parentID      *uint64
		parentLevel   int
		expectedLevel int
	}{
		{"root role (nil parent)", nil, 0, 0},
		{"root role (nil parent, any parent level)", nil, 5, 0},
		{"child of root", func() *uint64 { v := uint64(1); return &v }(), 0, 1},
		{"grandchild", func() *uint64 { v := uint64(2); return &v }(), 1, 2},
		{"level 3", func() *uint64 { v := uint64(3); return &v }(), 2, 3},
		{"level 4", func() *uint64 { v := uint64(4); return &v }(), 3, 4},
		{"level 5 (max)", func() *uint64 { v := uint64(5); return &v }(), 4, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role := &model.RoleModel{ParentID: tt.parentID}
			result := role.CalculateLevel(tt.parentLevel)
			assert.Equal(t, tt.expectedLevel, result)
		})
	}
}

// TestRoleValidateInheritanceDepth tests the ValidateInheritanceDepth method
// Requirements: 10.2
func TestRoleValidateInheritanceDepth(t *testing.T) {
	tests := []struct {
		name        string
		level       int
		expectError bool
	}{
		{"level 0 (root)", 0, false},
		{"level 1", 1, false},
		{"level 2", 2, false},
		{"level 3", 3, false},
		{"level 4", 4, false},
		{"level 5 (max allowed)", 5, false},
		{"level 6 (exceeds max)", 6, true},
		{"level 10 (far exceeds max)", 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role := &model.RoleModel{Level: tt.level}
			err := role.ValidateInheritanceDepth()
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, model.ErrRoleMaxDepthExceeded, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestRoleInheritanceFields tests the inheritance-related fields
// Requirements: 10.2
func TestRoleInheritanceFields(t *testing.T) {
	parentID := uint64(1)
	role := &model.RoleModel{
		Base:     model.Base{ID: 2},
		Slug:     "child_role",
		Name:     "子角色",
		ParentID: &parentID,
		Priority: 10,
		Level:    1,
	}

	assert.NotNil(t, role.ParentID)
	assert.Equal(t, uint64(1), *role.ParentID)
	assert.Equal(t, 10, role.Priority)
	assert.Equal(t, 1, role.Level)
}

// TestRoleMaxInheritanceDepthConstant tests the MaxRoleInheritanceDepth constant
// Requirements: 10.2
func TestRoleMaxInheritanceDepthConstant(t *testing.T) {
	assert.Equal(t, 5, model.MaxRoleInheritanceDepth)
}

// TestRoleInheritanceErrors tests the inheritance error constants
// Requirements: 10.2, 10.5
func TestRoleInheritanceErrors(t *testing.T) {
	assert.NotNil(t, model.ErrRoleMaxDepthExceeded)
	assert.NotNil(t, model.ErrRoleCircularInheritance)
	assert.Equal(t, "role inheritance depth exceeds maximum limit", model.ErrRoleMaxDepthExceeded.Error())
	assert.Equal(t, "circular role inheritance detected", model.ErrRoleCircularInheritance.Error())
}

// TestRolePriorityForConflictResolution tests the Priority field for conflict resolution
// Requirements: 10.4
func TestRolePriorityForConflictResolution(t *testing.T) {
	// Higher priority should override lower priority
	lowPriorityRole := &model.RoleModel{
		Slug:     "low_priority",
		Priority: 1,
	}
	highPriorityRole := &model.RoleModel{
		Slug:     "high_priority",
		Priority: 10,
	}

	assert.Less(t, lowPriorityRole.Priority, highPriorityRole.Priority)
}

// TestRoleInheritanceChainScenario tests a complete inheritance chain scenario
// Requirements: 10.2
func TestRoleInheritanceChainScenario(t *testing.T) {
	// Create a 3-level inheritance chain: superAdmin -> admin -> moderator
	superAdminID := uint64(1)
	adminID := uint64(2)

	superAdmin := &model.RoleModel{
		Base:     model.Base{ID: superAdminID},
		Slug:     "superAdmin",
		Name:     "超级管理员",
		ParentID: nil,
		Level:    0,
	}

	admin := &model.RoleModel{
		Base:     model.Base{ID: adminID},
		Slug:     "admin",
		Name:     "管理员",
		ParentID: &superAdminID,
		Level:    1,
	}

	moderator := &model.RoleModel{
		Base:     model.Base{ID: 3},
		Slug:     "moderator",
		Name:     "版主",
		ParentID: &adminID,
		Level:    2,
	}

	// Verify inheritance relationships
	assert.False(t, superAdmin.HasParent())
	assert.True(t, admin.HasParent())
	assert.True(t, moderator.HasParent())

	// Verify levels
	assert.Equal(t, 0, superAdmin.Level)
	assert.Equal(t, 1, admin.Level)
	assert.Equal(t, 2, moderator.Level)

	// Verify level calculation
	assert.Equal(t, 0, superAdmin.CalculateLevel(0))
	assert.Equal(t, 1, admin.CalculateLevel(superAdmin.Level))
	assert.Equal(t, 2, moderator.CalculateLevel(admin.Level))

	// All should pass depth validation
	assert.NoError(t, superAdmin.ValidateInheritanceDepth())
	assert.NoError(t, admin.ValidateInheritanceDepth())
	assert.NoError(t, moderator.ValidateInheritanceDepth())
}
