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

type UserRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo repository.UserRepository
	ctx  context.Context
}

func (s *UserRepositoryTestSuite) SetupSuite() {
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

func (s *UserRepositoryTestSuite) TearDownSuite() {
	sqlDB, err := s.db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func (s *UserRepositoryTestSuite) SetupTest() {
	// 清空用户表
	s.db.Exec("DELETE FROM users")
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (s *UserRepositoryTestSuite) TestCreate() {
	// 测试创建用户成功
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
	s.NotZero(user.CreatedAt)
	s.NotZero(user.UpdatedAt)

	// 测试重复邮箱应该失败
	user2 := &model.User{
		Email:        "test@example.com", // 重复的邮箱
		Phone:        "13800138001",
		Name:         "Test User 2",
		PasswordHash: "hashed-password-2",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
	}

	err = s.repo.Create(s.ctx, user2)
	s.Error(err)

	// 测试重复手机号应该失败
	user3 := &model.User{
		Email:        "unique@example.com",
		Phone:        "13800138000", // 重复的手机号
		Name:         "Test User 3",
		PasswordHash: "hashed-password-3",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
	}

	err = s.repo.Create(s.ctx, user3)
	s.Error(err)
}

func (s *UserRepositoryTestSuite) TestGet() {
	// 先创建一个用户
	user := &model.User{
		Email:        "get@example.com",
		Phone:        "13800138004",
		Name:         "Get User",
		PasswordHash: "hashed-password",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
	}

	err := s.repo.Create(s.ctx, user)
	s.NoError(err)

	// 测试获取存在的用户
	found, err := s.repo.Get(s.ctx, user.ID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(user.ID, found.ID)
	s.Equal(user.Email, found.Email)
	s.Equal(user.Phone, found.Phone)
	s.Equal(user.Name, found.Name)
	s.Equal(user.Role, found.Role)
	s.Equal(user.Status, found.Status)

	// 测试获取不存在的用户
	found, err = s.repo.Get(s.ctx, 99999)
	s.Error(err)
	s.Nil(found)
	s.Equal(repository.ErrNotFound, err)
}

func (s *UserRepositoryTestSuite) TestFindByEmail() {
	// 先创建一个用户
	user := &model.User{
		Email:        "findbyemail@example.com",
		Phone:        "13800138005",
		Name:         "FindByEmail User",
		PasswordHash: "hashed-password",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
	}

	err := s.repo.Create(s.ctx, user)
	s.NoError(err)

	// 测试通过邮箱查找存在的用户
	found, err := s.repo.FindByEmail(s.ctx, user.Email)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(user.ID, found.ID)
	s.Equal(user.Email, found.Email)

	// 测试通过邮箱查找不存在的用户
	found, err = s.repo.FindByEmail(s.ctx, "nonexistent@example.com")
	s.Error(err)
	s.Nil(found)
	s.Equal(repository.ErrNotFound, err)
}

func (s *UserRepositoryTestSuite) TestGetByPhone() {
	// 先创建一个用户
	user := &model.User{
		Email:        "getbyphone@example.com",
		Phone:        "13800138006",
		Name:         "GetByPhone User",
		PasswordHash: "hashed-password",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
	}

	err := s.repo.Create(s.ctx, user)
	s.NoError(err)

	// 测试通过手机查找存在的用户
	found, err := s.repo.GetByPhone(s.ctx, user.Phone)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(user.ID, found.ID)
	s.Equal(user.Phone, found.Phone)

	// 测试通过手机查找不存在的用户
	found, err = s.repo.GetByPhone(s.ctx, "13800000000")
	s.Error(err)
	s.Nil(found)
	s.Equal(repository.ErrNotFound, err)
}

func (s *UserRepositoryTestSuite) TestUpdate() {
	// 先创建一个用户
	user := &model.User{
		Email:        "update@example.com",
		Phone:        "13800138008",
		Name:         "Original Name",
		PasswordHash: "hashed-password",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		AvatarURL:    "https://example.com/original.jpg",
	}

	err := s.repo.Create(s.ctx, user)
	s.NoError(err)

	// 更新用户
	user.Name = "Updated Name"
	user.Phone = "13800138009"
	user.Email = "updated@example.com"
	user.AvatarURL = "https://example.com/updated.jpg"
	user.Role = model.RolePlayer
	user.Status = model.UserStatusSuspended
	user.PasswordHash = "updated-hashed-password"

	err = s.repo.Update(s.ctx, user)
	s.NoError(err)

	// 验证更新
	updated, err := s.repo.Get(s.ctx, user.ID)
	s.NoError(err)
	s.Equal("Updated Name", updated.Name)
	s.Equal("13800138009", updated.Phone)
	s.Equal("updated@example.com", updated.Email)
	s.Equal("https://example.com/updated.jpg", updated.AvatarURL)
	s.Equal(model.RolePlayer, updated.Role)
	s.Equal(model.UserStatusSuspended, updated.Status)
	s.Equal("updated-hashed-password", updated.PasswordHash)
}

func (s *UserRepositoryTestSuite) TestDelete() {
	// 先创建一个用户
	user := &model.User{
		Email:        "delete@example.com",
		Phone:        "13800138010",
		Name:         "Delete User",
		PasswordHash: "hashed-password",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
	}

	err := s.repo.Create(s.ctx, user)
	s.NoError(err)

	// 删除用户
	err = s.repo.Delete(s.ctx, user.ID)
	s.NoError(err)

	// 验证用户已被删除
	found, err := s.repo.Get(s.ctx, user.ID)
	s.Error(err)
	s.Nil(found)
	s.Equal(repository.ErrNotFound, err)

	// 测试删除不存在的用户
	err = s.repo.Delete(s.ctx, 99999)
	s.Error(err)
	s.Equal(repository.ErrNotFound, err)
}

func (s *UserRepositoryTestSuite) TestList() {
	// 清空表
	s.db.Exec("DELETE FROM users")

	// 创建多个用户
	users := []*model.User{
		{
			Email:        "list1@example.com",
			Phone:        "13800138011",
			Name:         "List User 1",
			PasswordHash: "hashed-password-1",
			Role:         model.RoleUser,
			Status:       model.UserStatusActive,
		},
		{
			Email:        "list2@example.com",
			Phone:        "13800138012",
			Name:         "List User 2",
			PasswordHash: "hashed-password-2",
			Role:         model.RolePlayer,
			Status:       model.UserStatusActive,
		},
		{
			Email:        "list3@example.com",
			Phone:        "13800138013",
			Name:         "List User 3",
			PasswordHash: "hashed-password-3",
			Role:         model.RoleAdmin,
			Status:       model.UserStatusSuspended,
		},
	}

	for _, user := range users {
		err := s.repo.Create(s.ctx, user)
		s.NoError(err)
	}

	// 获取用户列表
	list, err := s.repo.List(s.ctx)
	s.NoError(err)
	s.Len(list, 3)

	// 测试空列表
	s.db.Exec("DELETE FROM users")
	list, err = s.repo.List(s.ctx)
	s.NoError(err)
	s.Empty(list)
}

func (s *UserRepositoryTestSuite) TestListPaged() {
	// 清空表
	s.db.Exec("DELETE FROM users")

	// 创建5个用户
	for i := 0; i < 5; i++ {
		user := &model.User{
			Email:        "test@example.com",
			Phone:        "13800138000",
			Name:         "Test User",
			PasswordHash: "hashed-password",
			Role:         model.RoleUser,
			Status:       model.UserStatusActive,
		}
		err := s.repo.Create(s.ctx, user)
		s.NoError(err)
	}

	// 测试分页
	users, total, err := s.repo.ListPaged(s.ctx, 1, 2)
	s.NoError(err)
	s.Len(users, 2)
	s.Equal(int64(5), total)

	// 测试第二页
	users, total, err = s.repo.ListPaged(s.ctx, 2, 2)
	s.NoError(err)
	s.Len(users, 2)
	s.Equal(int64(5), total)

	// 测试最后一页
	users, total, err = s.repo.ListPaged(s.ctx, 3, 2)
	s.NoError(err)
	s.Len(users, 1)
	s.Equal(int64(5), total)
}

func (s *UserRepositoryTestSuite) TestListWithFilters() {
	// 清空表
	s.db.Exec("DELETE FROM users")

	// 创建测试用户
	users := []*model.User{
		{
			Email:        "filter1@example.com",
			Phone:        "13800138040",
			Name:         "Filter User 1",
			PasswordHash: "hashed-password-1",
			Role:         model.RoleUser,
			Status:       model.UserStatusActive,
		},
		{
			Email:        "filter2@example.com",
			Phone:        "13800138041",
			Name:         "Filter User 2",
			PasswordHash: "hashed-password-2",
			Role:         model.RolePlayer,
			Status:       model.UserStatusSuspended,
		},
		{
			Email:        "filter3@example.com",
			Phone:        "13800138042",
			Name:         "Filter User 3",
			PasswordHash: "hashed-password-3",
			Role:         model.RoleAdmin,
			Status:       model.UserStatusActive,
		},
	}

	for _, user := range users {
		err := s.repo.Create(s.ctx, user)
		s.NoError(err)
	}

	// 按角色过滤
	role := model.RoleUser
	opts := repository.UserListOptions{
		Page:     1,
		PageSize: 10,
		Role:     role,
	}
	filtered, total, err := s.repo.ListWithFilters(s.ctx, opts)
	s.NoError(err)
	s.Len(filtered, 1)
	s.Equal(int64(1), total)
	s.Equal(model.RoleUser, filtered[0].Role)

	// 按状态过滤
	status := model.UserStatusSuspended
	opts = repository.UserListOptions{
		Page:     1,
		PageSize: 10,
		Status:   status,
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
		Keyword:  "Filter User",
	}
	filtered, total, err = s.repo.ListWithFilters(s.ctx, opts)
	s.NoError(err)
	s.Len(filtered, 3)
	s.Equal(int64(3), total)

	// 组合过滤条件
	opts = repository.UserListOptions{
		Page:     1,
		PageSize: 10,
		Roles:    []model.Role{model.RoleUser, model.RolePlayer, model.RoleAdmin},
		Keyword:  "Filter",
	}
	filtered, total, err = s.repo.ListWithFilters(s.ctx, opts)
	s.NoError(err)
	s.Len(filtered, 3)
	s.Equal(int64(3), total)
}
