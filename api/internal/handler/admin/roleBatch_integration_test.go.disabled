package admin_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/handler/admin"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	roleservice "gamelink/internal/service/admin"
	"gamelink/pkg/apierr"
)

// TestRoleBatchHandler_BatchDeleteRoles 测试批量删除角色
func TestRoleBatchHandler_BatchDeleteRoles(t *testing.T) {
	// 跳过如果没有测试数据库
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// 创建测试角色
	role1 := CreateTestRole(t, db, "test-role-1", "Test Role 1")
	role2 := CreateTestRole(t, db, "test-role-2", "Test Role 2")
	role3 := CreateTestRole(t, db, "test-role-3", "Test Role 3")

	// 创建系统角色（不应该被删除）
	systemRole := &model.RoleModel{
		Slug:     "system-role-batch-test",
		Name:     "System Role Batch Test",
		IsSystem: true,
	}
	require.NoError(t, db.Create(systemRole).Error)

	// 创建 handler
	roleRepo := repository.NewRoleRepository(db)
	roleSvc := roleservice.NewRoleService(roleRepo, nil)
	handler := admin.NewRoleBatchHandler(roleSvc)

	// 设置 Gin 测试模式
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/roles/batch/delete", func(c *gin.Context) {
		// 模拟认证中间件
		c.Set("user_id", uint64(1))
		handler.BatchDeleteRoles(c)
	})

	// 测试用例
	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		validate       func(t *testing.T, body []byte, db *gorm.DB)
	}{
		{
			name: "成功批量删除角色",
			requestBody: map[string][]uint64{
				"roleIds": {role1.ID, role2.ID, role3.ID},
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, body []byte, db *gorm.DB) {
				var response admin.BatchOperationResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)

				assert.Equal(t, 3, response.SuccessCount)
				assert.Equal(t, 0, response.FailedCount)
				assert.Equal(t, 3, response.TotalCount)

				// 验证角色已被删除
				var count int64
				db.Model(&model.RoleModel{}).Where("id IN ?", []uint64{role1.ID, role2.ID, role3.ID}).Count(&count)
				assert.Equal(t, int64(0), count)
			},
		},
		{
			name: "系统角色不可删除",
			requestBody: map[string][]uint64{
				"roleIds": {systemRole.ID},
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, body []byte, db *gorm.DB) {
				var response admin.BatchOperationResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)

				assert.Equal(t, 0, response.SuccessCount)
				assert.Equal(t, 1, response.FailedCount)
				assert.Equal(t, 1, response.TotalCount)
				assert.Len(t, response.FailedItems, 1)
				assert.Contains(t, response.FailedItems[0].Message, "系统角色")
			},
		},
		{
			name: "混合删除（部分成功）",
			requestBody: map[string][]uint64{
				"roleIds": {systemRole.ID, CreateTestRole(t, db, "temp-role", "Temp Role").ID},
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, body []byte, db *gorm.DB) {
				var response admin.BatchOperationResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)

				assert.Equal(t, 1, response.SuccessCount)
				assert.Equal(t, 1, response.FailedCount)
				assert.Equal(t, 2, response.TotalCount)
			},
		},
		{
			name:           "空角色ID列表",
			requestBody:    map[string][]uint64{"roleIds": {}},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "超过100个角色限制",
			requestBody:    map[string][]uint64{"roleIds": make([]uint64, 101)},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/roles/batch/delete", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.validate != nil {
				tt.validate(t, w.Body.Bytes(), db)
			}
		})
	}
}

// TestRoleBatchHandler_BatchAssignPermissions 测试批量分配权限
func TestRoleBatchHandler_BatchAssignPermissions(t *testing.T) {
	// 跳过如果没有测试数据库
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// 创建测试角色和权限
	role1 := CreateTestRole(t, db, "test-role-perm-1", "Test Role Perm 1")
	role2 := CreateTestRole(t, db, "test-role-perm-2", "Test Role Perm 2")

	perm1 := CreateTestPermission(t, db, "test.perm.1", "GET", "/api/test/1", "test.perm.1")
	perm2 := CreateTestPermission(t, db, "test.perm.2", "POST", "/api/test/2", "test.perm.2")
	perm3 := CreateTestPermission(t, db, "test.perm.3", "PUT", "/api/test/3", "test.perm.3")

	// 创建 handler
	roleRepo := repository.NewRoleRepository(db)
	roleSvc := roleservice.NewRoleService(roleRepo, nil)
	handler := admin.NewRoleBatchHandler(roleSvc)

	// 设置 Gin 测试模式
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/roles/batch/assign-permissions", func(c *gin.Context) {
		// 模拟认证中间件
		c.Set("user_id", uint64(1))
		handler.BatchAssignPermissions(c)
	})

	// 测试用例
	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		validate       func(t *testing.T, body []byte, db *gorm.DB, role1, role2 uint64)
	}{
		{
			name: "成功批量分配权限",
			requestBody: map[string]interface{}{
				"roleIds":       []uint64{role1.ID, role2.ID},
				"permissionIds": []uint64{perm1.ID, perm2.ID, perm3.ID},
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, body []byte, db *gorm.DB, r1, r2 uint64) {
				var response admin.BatchOperationResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)

				assert.Equal(t, 2, response.SuccessCount)
				assert.Equal(t, 0, response.FailedCount)
				assert.Equal(t, 2, response.TotalCount)

				// 验证权限已分配
				var count1, count2 int64
				db.Table("role_permissions").Where("role_id = ?", r1).Count(&count1)
				db.Table("role_permissions").Where("role_id = ?", r2).Count(&count2)
				assert.Equal(t, int64(3), count1)
				assert.Equal(t, int64(3), count2)
			},
		},
		{
			name: "空角色ID列表",
			requestBody: map[string]interface{}{
				"roleIds":       []uint64{},
				"permissionIds": []uint64{perm1.ID},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "空权限ID列表",
			requestBody: map[string]interface{}{
				"roleIds":       []uint64{role1.ID},
				"permissionIds": []uint64{},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "超过100个角色限制",
			requestBody: map[string]interface{}{
				"roleIds":       make([]uint64, 101),
				"permissionIds": []uint64{perm1.ID},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/roles/batch/assign-permissions", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.validate != nil {
				tt.validate(t, w.Body.Bytes(), db, role1.ID, role2.ID)
			}
		})
	}
}

// TestRoleBatchService_BatchDeleteRoles 测试服务层批量删除角色
func TestRoleBatchService_BatchDeleteRoles(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	roleRepo := repository.NewRoleRepository(db)
	roleSvc := roleservice.NewRoleService(roleRepo, nil)

	// 创建测试角色
	role1 := CreateTestRole(t, db, "svc-test-role-1", "SVC Test Role 1")
	role2 := CreateTestRole(t, db, "svc-test-role-2", "SVC Test Role 2")
	systemRole := &model.RoleModel{
		Slug:     "system-role-svc-test",
		Name:     "System Role SVC Test",
		IsSystem: true,
	}
	require.NoError(t, db.Create(systemRole).Error)

	t.Run("成功批量删除角色", func(t *testing.T) {
		result, err := roleSvc.BatchDeleteRoles(db.Statement.Context, []uint64{role1.ID, role2.ID})
		require.NoError(t, err)
		assert.Equal(t, 2, result.SuccessCount)
		assert.Equal(t, 0, result.FailedCount)

		// 验证角色已被删除
		_, err = roleRepo.Get(db.Statement.Context, role1.ID)
		assert.Error(t, err)
		assert.True(t, apierr.IsNotFound(err))
	})

	t.Run("系统角色不可删除", func(t *testing.T) {
		result, err := roleSvc.BatchDeleteRoles(db.Statement.Context, []uint64{systemRole.ID})
		require.NoError(t, err)
		assert.Equal(t, 0, result.SuccessCount)
		assert.Equal(t, 1, result.FailedCount)
		assert.Len(t, result.FailedRoles, 1)
		assert.Contains(t, result.FailedRoles[0].Reason, "系统角色")
	})

	t.Run("空角色ID列表", func(t *testing.T) {
		_, err := roleSvc.BatchDeleteRoles(db.Statement.Context, []uint64{})
		assert.Error(t, err)
	})
}

// TestRoleBatchService_BatchAssignPermissionsToMultipleRoles 测试服务层批量分配权限
func TestRoleBatchService_BatchAssignPermissionsToMultipleRoles(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	roleRepo := repository.NewRoleRepository(db)
	roleSvc := roleservice.NewRoleService(roleRepo, nil)

	// 创建测试角色和权限
	role1 := CreateTestRole(t, db, "perm-test-role-1", "Perm Test Role 1")
	role2 := CreateTestRole(t, db, "perm-test-role-2", "Perm Test Role 2")
	perm1 := CreateTestPermission(t, db, "svc.test.perm.1", "GET", "/api/svc/test/1", "svc.test.perm.1")
	perm2 := CreateTestPermission(t, db, "svc.test.perm.2", "POST", "/api/svc/test/2", "svc.test.perm.2")

	t.Run("成功批量分配权限", func(t *testing.T) {
		result, err := roleSvc.BatchAssignPermissionsToMultipleRoles(
			db.Statement.Context,
			[]uint64{role1.ID, role2.ID},
			[]uint64{perm1.ID, perm2.ID},
		)
		require.NoError(t, err)
		assert.Equal(t, 2, result.SuccessCount)
		assert.Equal(t, 0, result.FailedCount)

		// 验证权限已分配
		permIDs, err := roleSvc.GetRolePermissionIDs(db.Statement.Context, role1.ID)
		require.NoError(t, err)
		assert.Len(t, permIDs, 2)
	})

	t.Run("空角色ID列表", func(t *testing.T) {
		_, err := roleSvc.BatchAssignPermissionsToMultipleRoles(
			db.Statement.Context,
			[]uint64{},
			[]uint64{perm1.ID},
		)
		assert.Error(t, err)
	})

	t.Run("权限ID列表为空（应该成功，清空权限）", func(t *testing.T) {
		role3 := CreateTestRole(t, db, "perm-test-role-3", "Perm Test Role 3")
		// 先分配一些权限
		_ = roleSvc.AssignPermissionsToRole(db.Statement.Context, role3.ID, []uint64{perm1.ID})

		// 然后批量清空
		result, err := roleSvc.BatchAssignPermissionsToMultipleRoles(
			db.Statement.Context,
			[]uint64{role3.ID},
			[]uint64{},
		)
		require.NoError(t, err)
		assert.Equal(t, 1, result.SuccessCount)

		// 验证权限已被清空
		permIDs, err := roleSvc.GetRolePermissionIDs(db.Statement.Context, role3.ID)
		require.NoError(t, err)
		assert.Len(t, permIDs, 0)
	})
}
