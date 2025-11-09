package role

import (
	"context"

	"gamelink/internal/model"
	"gamelink/internal/repository"

	"github.com/stretchr/testify/mock"
)

// MockRoleRepository is a mock for RoleRepository
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) List(ctx context.Context) ([]model.RoleModel, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.RoleModel), args.Error(1)
}

func (m *MockRoleRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.RoleModel, int64, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.RoleModel), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoleRepository) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, isSystem *bool) ([]model.RoleModel, int64, error) {
	args := m.Called(ctx, page, pageSize, keyword, isSystem)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.RoleModel), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoleRepository) ListWithPermissions(ctx context.Context) ([]model.RoleModel, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.RoleModel), args.Error(1)
}

func (m *MockRoleRepository) Get(ctx context.Context, id uint64) (*model.RoleModel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RoleModel), args.Error(1)
}

func (m *MockRoleRepository) GetWithPermissions(ctx context.Context, id uint64) (*model.RoleModel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RoleModel), args.Error(1)
}

func (m *MockRoleRepository) GetBySlug(ctx context.Context, slug string) (*model.RoleModel, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RoleModel), args.Error(1)
}

func (m *MockRoleRepository) Create(ctx context.Context, role *model.RoleModel) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) Update(ctx context.Context, role *model.RoleModel) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoleRepository) AssignPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	args := m.Called(ctx, roleID, permissionIDs)
	return args.Error(0)
}

func (m *MockRoleRepository) AddPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	args := m.Called(ctx, roleID, permissionIDs)
	return args.Error(0)
}

func (m *MockRoleRepository) RemovePermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	args := m.Called(ctx, roleID, permissionIDs)
	return args.Error(0)
}

func (m *MockRoleRepository) AssignToUser(ctx context.Context, userID uint64, roleIDs []uint64) error {
	args := m.Called(ctx, userID, roleIDs)
	return args.Error(0)
}

func (m *MockRoleRepository) RemoveFromUser(ctx context.Context, userID uint64, roleIDs []uint64) error {
	args := m.Called(ctx, userID, roleIDs)
	return args.Error(0)
}

func (m *MockRoleRepository) ListByUserID(ctx context.Context, userID uint64) ([]model.RoleModel, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.RoleModel), args.Error(1)
}

func (m *MockRoleRepository) CheckUserHasRole(ctx context.Context, userID uint64, roleSlug string) (bool, error) {
	args := m.Called(ctx, userID, roleSlug)
	return args.Bool(0), args.Error(1)
}

var _ repository.RoleRepository = (*MockRoleRepository)(nil)

