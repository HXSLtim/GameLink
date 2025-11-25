package permission_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	permissionRepo "gamelink/internal/repository/permission"
)

type PermissionRepositorySimpleTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo repository.PermissionRepository
	ctx  context.Context
}

func (s *PermissionRepositorySimpleTestSuite) SetupSuite() {
	// 使用SQLite内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Skipf("sqlite driver unavailable: %v", err)
		return
	}

	// 迁移表结构 - 需要权限相关的表
	err = db.AutoMigrate(&model.Permission{}, &model.RolePermission{}, &model.UserRole{})
	if err != nil {
		s.T().Skipf("sqlite migration unavailable: %v", err)
		return
	}

	s.db = db
	s.repo = permissionRepo.NewPermissionRepository(db)
	s.ctx = context.Background()
}

func (s *PermissionRepositorySimpleTestSuite) TearDownSuite() {
	sqlDB, err := s.db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func (s *PermissionRepositorySimpleTestSuite) SetupTest() {
	if s.db == nil {
		s.T().Skip("sqlite unavailable")
		return
	}
	// 清空权限表
	s.db.Exec("DELETE FROM permissions")
	s.db.Exec("DELETE FROM role_permissions")
	s.db.Exec("DELETE FROM user_roles")
}

func TestPermissionRepositorySimpleTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionRepositorySimpleTestSuite))
}

func (s *PermissionRepositorySimpleTestSuite) TestBasicCRUD() {
	// 测试创建权限
	permission := &model.Permission{
		Code:        "user:create",
		Group:       "user",
		Method:      model.HTTPMethodPOST,
		Path:        "/api/v1/users",
		Description: "Create user permission",
	}

	err := s.repo.Create(s.ctx, permission)
	s.NoError(err)
	s.NotZero(permission.ID)

	// 测试获取权限
	found, err := s.repo.Get(s.ctx, permission.ID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(permission.ID, found.ID)
	s.Equal(permission.Code, found.Code)
	s.Equal(permission.Group, found.Group)
	s.Equal(permission.Method, found.Method)
	s.Equal(permission.Path, found.Path)

	// 测试通过code获取权限
	foundByCode, err := s.repo.GetByCode(s.ctx, permission.Code)
	s.NoError(err)
	s.NotNil(foundByCode)
	s.Equal(permission.ID, foundByCode.ID)

	// 测试通过method和path获取权限
	foundByMethodPath, err := s.repo.GetByMethodAndPath(s.ctx, string(permission.Method), permission.Path)
	s.NoError(err)
	s.NotNil(foundByMethodPath)
	s.Equal(permission.ID, foundByMethodPath.ID)

	// 测试更新权限
	permission.Description = "Updated description"
	err = s.repo.Update(s.ctx, permission)
	s.NoError(err)

	updated, err := s.repo.Get(s.ctx, permission.ID)
	s.NoError(err)
	s.Equal("Updated description", updated.Description)

	// 测试删除权限
	err = s.repo.Delete(s.ctx, permission.ID)
	s.NoError(err)

	deleted, err := s.repo.Get(s.ctx, permission.ID)
	s.Error(err)
	s.Nil(deleted)
	s.Equal(repository.ErrNotFound, err)
}

func (s *PermissionRepositorySimpleTestSuite) TestListOperations() {
	// 创建多个权限
	permissions := []*model.Permission{
		{
			Code:        "permission1",
			Group:       "user",
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/users",
			Description: "Get users",
		},
		{
			Code:        "permission2",
			Group:       "user",
			Method:      model.HTTPMethodPOST,
			Path:        "/api/v1/users",
			Description: "Create user",
		},
		{
			Code:        "permission3",
			Group:       "admin",
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin",
			Description: "Admin permission",
		},
	}

	for _, permission := range permissions {
		err := s.repo.Create(s.ctx, permission)
		s.NoError(err)
	}

	// 测试列表所有权限
	list, err := s.repo.List(s.ctx)
	s.NoError(err)
	s.Len(list, 3)

	// 测试分页
	pagedList, total, err := s.repo.ListPaged(s.ctx, 1, 2)
	s.NoError(err)
	s.Len(pagedList, 2)
	s.Equal(int64(3), total)

	// 第二页
	pagedList, total, err = s.repo.ListPaged(s.ctx, 2, 2)
	s.NoError(err)
	s.Len(pagedList, 1)
	s.Equal(int64(3), total)
}

func (s *PermissionRepositorySimpleTestSuite) TestFilters() {
	// 创建测试权限
	permissions := []*model.Permission{
		{
			Code:        "user-get",
			Group:       "user",
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/users",
			Description: "Get users",
		},
		{
			Code:        "user-post",
			Group:       "user",
			Method:      model.HTTPMethodPOST,
			Path:        "/api/v1/users",
			Description: "Create user",
		},
		{
			Code:        "admin-put",
			Group:       "admin",
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin",
			Description: "Admin permission",
		},
	}

	for _, permission := range permissions {
		err := s.repo.Create(s.ctx, permission)
		s.NoError(err)
	}

	// 按关键词过滤
	filtered, total, err := s.repo.ListPagedWithFilter(s.ctx, 1, 10, "user", "", "")
	s.NoError(err)
	s.Len(filtered, 2)
	s.Equal(int64(2), total)

	// 按方法过滤
	filtered, total, err = s.repo.ListPagedWithFilter(s.ctx, 1, 10, "", "GET", "")
	s.NoError(err)
	s.Len(filtered, 1)
	s.Equal(int64(1), total)
	s.Equal(model.HTTPMethodGET, filtered[0].Method)

	// 按分组过滤
	filtered, total, err = s.repo.ListPagedWithFilter(s.ctx, 1, 10, "", "", "admin")
	s.NoError(err)
	s.Len(filtered, 1)
	s.Equal(int64(1), total)
	s.Equal("admin", filtered[0].Group)
}

func (s *PermissionRepositorySimpleTestSuite) TestListByGroup() {
	// 创建权限
	permissions := []*model.Permission{
		{
			Code:        "user-perm-1",
			Group:       "user",
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/users",
			Description: "User permission 1",
		},
		{
			Code:        "user-perm-2",
			Group:       "user",
			Method:      model.HTTPMethodPOST,
			Path:        "/api/v1/users",
			Description: "User permission 2",
		},
		{
			Code:        "admin-perm-1",
			Group:       "admin",
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin",
			Description: "Admin permission 1",
		},
	}

	for _, permission := range permissions {
		err := s.repo.Create(s.ctx, permission)
		s.NoError(err)
	}

	// 按分组获取权限
	grouped, err := s.repo.ListByGroup(s.ctx)
	s.NoError(err)
	s.Len(grouped, 2) // user和admin两个分组

	// 验证user分组
	userPermissions, exists := grouped["user"]
	s.True(exists)
	s.Len(userPermissions, 2)

	// 验证admin分组
	adminPermissions, exists := grouped["admin"]
	s.True(exists)
	s.Len(adminPermissions, 1)
}

func (s *PermissionRepositorySimpleTestSuite) TestListGroups() {
	// 创建权限（包含重复分组）
	permissions := []*model.Permission{
		{
			Code:        "perm-1",
			Group:       "user",
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/users",
			Description: "Permission 1",
		},
		{
			Code:        "perm-2",
			Group:       "user",
			Method:      model.HTTPMethodPOST,
			Path:        "/api/v1/users",
			Description: "Permission 2",
		},
		{
			Code:        "perm-3",
			Group:       "admin",
			Method:      model.HTTPMethodPUT,
			Path:        "/api/v1/admin",
			Description: "Permission 3",
		},
		{
			Code:        "perm-4",
			Group:       "system",
			Method:      model.HTTPMethodDELETE,
			Path:        "/api/v1/system",
			Description: "Permission 4",
		},
	}

	for _, permission := range permissions {
		err := s.repo.Create(s.ctx, permission)
		s.NoError(err)
	}

	// 获取分组列表
	groups, err := s.repo.ListGroups(s.ctx)
	s.NoError(err)
	s.Len(groups, 3) // user, admin, system

	// 验证分组内容
	expectedGroups := []string{"admin", "system", "user"}
	s.Equal(expectedGroups, groups)
}

func (s *PermissionRepositorySimpleTestSuite) TestUpsertByMethodPath() {
	// 测试upsert新权限
	permission := &model.Permission{
		Code:        "upsert-new",
		Group:       "api",
		Method:      model.HTTPMethodPATCH,
		Path:        "/api/v1/resource/:id",
		Description: "New permission for upsert",
	}

	// Upsert新权限
	err := s.repo.UpsertByMethodPath(s.ctx, permission)
	s.NoError(err)
	s.NotZero(permission.ID)

	// 验证权限被创建
	found, err := s.repo.Get(s.ctx, permission.ID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(permission.Code, found.Code)
	s.Equal(permission.Method, found.Method)
	s.Equal(permission.Path, found.Path)

	// 测试upsert现有权限
	existingID := permission.ID
	permission.Code = "upsert-updated"
	permission.Group = "updated"
	permission.Description = "Updated permission"

	// Upsert现有权限
	err = s.repo.UpsertByMethodPath(s.ctx, permission)
	s.NoError(err)
	s.Equal(existingID, permission.ID) // ID应该保持不变

	// 验证权限被更新
	updated, err := s.repo.Get(s.ctx, existingID)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal("upsert-updated", updated.Code)
	s.Equal("updated", updated.Group)
	s.Equal("Updated permission", updated.Description)
}
