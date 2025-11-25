package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	userRepo "gamelink/internal/repository/user"
)

type UserRepositorySimpleTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo repository.UserRepository
	ctx  context.Context
}

func (s *UserRepositorySimpleTestSuite) SetupSuite() {
	// 使用SQLite内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Skipf("sqlite driver unavailable: %v", err)
		return
	}

	// 迁移表结构 - 需要用户相关的表
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		s.T().Skipf("sqlite migration unavailable: %v", err)
		return
	}

	s.db = db
	s.repo = userRepo.NewUserRepository(db)
	s.ctx = context.Background()
}

func (s *UserRepositorySimpleTestSuite) TearDownSuite() {
	sqlDB, err := s.db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func (s *UserRepositorySimpleTestSuite) SetupTest() {
	// 清空用户表
	s.db.Exec("DELETE FROM users")
}

func TestUserRepositorySimpleTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositorySimpleTestSuite))
}

func (s *UserRepositorySimpleTestSuite) TestBasicCRUD() {
	// 测试创建用户
	user := &model.User{
		Email:        "test@example.com",
		Phone:        "13800138000",
		Name:         "Test User",
		PasswordHash: "hashed-password",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		AvatarURL:    "https://example.com/avatar.jpg",
	}

	err := s.repo.Create(s.ctx, user)
	s.NoError(err)
	s.NotZero(user.ID)

	// 测试获取用户
	found, err := s.repo.Get(s.ctx, user.ID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(user.ID, found.ID)
	s.Equal(user.Email, found.Email)

	// 测试通过邮箱查找
	foundByEmail, err := s.repo.FindByEmail(s.ctx, user.Email)
	s.NoError(err)
	s.NotNil(foundByEmail)
	s.Equal(user.ID, foundByEmail.ID)

	// 测试通过手机查找
	foundByPhone, err := s.repo.GetByPhone(s.ctx, user.Phone)
	s.NoError(err)
	s.NotNil(foundByPhone)
	s.Equal(user.ID, foundByPhone.ID)

	// 测试更新用户
	user.Name = "Updated Name"
	err = s.repo.Update(s.ctx, user)
	s.NoError(err)

	updated, err := s.repo.Get(s.ctx, user.ID)
	s.NoError(err)
	s.Equal("Updated Name", updated.Name)

	// 测试删除用户
	err = s.repo.Delete(s.ctx, user.ID)
	s.NoError(err)

	deleted, err := s.repo.Get(s.ctx, user.ID)
	s.Error(err)
	s.Nil(deleted)
	s.Equal(repository.ErrNotFound, err)
}

func (s *UserRepositorySimpleTestSuite) TestListOperations() {
	// 创建多个用户
	users := []*model.User{
		{
			Email:        "user1@example.com",
			Phone:        "13800138001",
			Name:         "User 1",
			PasswordHash: "hash1",
			Role:         model.RoleUser,
			Status:       model.UserStatusActive,
		},
		{
			Email:        "user2@example.com",
			Phone:        "13800138002",
			Name:         "User 2",
			PasswordHash: "hash2",
			Role:         model.RolePlayer,
			Status:       model.UserStatusActive,
		},
		{
			Email:        "user3@example.com",
			Phone:        "13800138003",
			Name:         "User 3",
			PasswordHash: "hash3",
			Role:         model.RoleAdmin,
			Status:       model.UserStatusSuspended,
		},
	}

	for _, user := range users {
		err := s.repo.Create(s.ctx, user)
		s.NoError(err)
	}

	// 测试列表所有用户
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

func (s *UserRepositorySimpleTestSuite) TestFilters() {
	// 创建测试用户
	users := []*model.User{
		{
			Email:        "filter1@example.com",
			Phone:        "13800138010",
			Name:         "Alice User",
			PasswordHash: "hash1",
			Role:         model.RoleUser,
			Status:       model.UserStatusActive,
		},
		{
			Email:        "filter2@example.com",
			Phone:        "13800138020",
			Name:         "Bob Player",
			PasswordHash: "hash2",
			Role:         model.RolePlayer,
			Status:       model.UserStatusSuspended,
		},
		{
			Email:        "filter3@example.com",
			Phone:        "13800138030",
			Name:         "Charlie Admin",
			PasswordHash: "hash3",
			Role:         model.RoleAdmin,
			Status:       model.UserStatusActive,
		},
	}

	for _, user := range users {
		err := s.repo.Create(s.ctx, user)
		s.NoError(err)
	}

	// 按角色过滤（使用复数字段）
	opts := repository.UserListOptions{
		Page:     1,
		PageSize: 10,
		Roles:    []model.Role{model.RoleUser},
	}
	filtered, total, err := s.repo.ListWithFilters(s.ctx, opts)
	s.NoError(err)
	s.Len(filtered, 1)
	s.Equal(int64(1), total)
	s.Equal(model.RoleUser, filtered[0].Role)

	// 按状态过滤（使用复数字段）
	opts = repository.UserListOptions{
		Page:     1,
		PageSize: 10,
		Statuses: []model.UserStatus{model.UserStatusSuspended},
	}
	filtered, total, err = s.repo.ListWithFilters(s.ctx, opts)
	s.NoError(err)
	s.Len(filtered, 1)
	s.Equal(int64(1), total)
	s.Equal(model.UserStatusSuspended, filtered[0].Status)

	// 按多个角色过滤
	opts = repository.UserListOptions{
		Page:     1,
		PageSize: 10,
		Roles:    []model.Role{model.RoleUser, model.RolePlayer},
	}
	filtered, total, err = s.repo.ListWithFilters(s.ctx, opts)
	s.NoError(err)
	s.Len(filtered, 2)
	s.Equal(int64(2), total)

	// 按关键词过滤
	opts = repository.UserListOptions{
		Page:     1,
		PageSize: 10,
		Keyword:  "Alice",
	}
	filtered, total, err = s.repo.ListWithFilters(s.ctx, opts)
	s.NoError(err)
	s.Len(filtered, 1)
	s.Equal(int64(1), total)
	s.Contains(filtered[0].Name, "Alice")
}
