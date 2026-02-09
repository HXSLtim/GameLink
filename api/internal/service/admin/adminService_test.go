package admin

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	repoiface "gamelink/internal/repository/interfaces"
	"gamelink/pkg/apierr"
)

// MockGameRepository is a mock implementation of GameRepository
type MockGameRepository struct {
	mock.Mock
}

func (m *MockGameRepository) Create(ctx context.Context, game *model.Game) error {
	args := m.Called(ctx, game)
	return args.Error(0)
}

func (m *MockGameRepository) Get(ctx context.Context, id uint64) (*model.Game, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Game), args.Error(1)
}

func (m *MockGameRepository) List(ctx context.Context) ([]model.Game, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Game), args.Error(1)
}

func (m *MockGameRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Game, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]model.Game), args.Get(1).(int64), args.Error(2)
}

func (m *MockGameRepository) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string) ([]model.Game, int64, error) {
	args := m.Called(ctx, page, pageSize, keyword)
	return args.Get(0).([]model.Game), args.Get(1).(int64), args.Error(2)
}

func (m *MockGameRepository) Update(ctx context.Context, game *model.Game) error {
	args := m.Called(ctx, game)
	return args.Error(0)
}

func (m *MockGameRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockGameRepository) BatchDelete(ctx context.Context, ids []uint64) (int64, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockGameRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error) {
	args := m.Called(ctx, ids, isActive)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockGameRepository) BatchUpdateSortOrder(ctx context.Context, updates map[uint64]int) (int64, error) {
	args := m.Called(ctx, updates)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockGameRepository) BatchUpdateCategory(ctx context.Context, ids []uint64, category string) (int64, error) {
	args := m.Called(ctx, ids, category)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockGameRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.Game, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Game), args.Error(1)
}

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
	users map[uint64]*model.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[uint64]*model.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	if user.ID == 0 {
		user.ID = 1 // Assign a default ID for testing
	}
	// Only store in map if it's initialized (when using NewMockUserRepository)
	if m.users != nil {
		m.users[user.ID] = user
	}
	return nil
}

func (m *MockUserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context) ([]model.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) Count(ctx context.Context, opts repository.UserListOptions) (int, error) {
	args := m.Called(ctx, opts)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.User, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	args := m.Called(ctx, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	return m.FindByPhone(ctx, phone)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return m.FindByEmail(ctx, email)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, userID uint64, newPassword string) error {
	args := m.Called(ctx, userID, newPassword)
	return args.Error(0)
}

func (m *MockUserRepository) GetByWeChatOpenID(ctx context.Context, openID string) (*model.User, error) {
	args := m.Called(ctx, openID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByWeChatUnionID(ctx context.Context, unionID string) (*model.User, error) {
	args := m.Called(ctx, unionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

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

func (m *MockRoleRepository) SetParent(ctx context.Context, roleID uint64, parentID *uint64) error {
	return nil
}

func (m *MockRoleRepository) GetInheritanceChain(ctx context.Context, roleID uint64) ([]model.RoleModel, error) {
	return []model.RoleModel{}, nil
}

func (m *MockRoleRepository) GetChildRoles(ctx context.Context, roleID uint64) ([]model.RoleModel, error) {
	return []model.RoleModel{}, nil
}

func (m *MockRoleRepository) UpdateLevel(ctx context.Context, roleID uint64, level int) error {
	return nil
}

func (m *MockRoleRepository) GetUserIDsByRoleID(ctx context.Context, roleID uint64) ([]uint64, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uint64), args.Error(1)
}

// MockPlayerRepository is a mock implementation of PlayerRepository
type MockPlayerRepository struct {
	mock.Mock
}

func (m *MockPlayerRepository) Create(ctx context.Context, player *model.Player) error {
	args := m.Called(ctx, player)
	return args.Error(0)
}

func (m *MockPlayerRepository) Get(ctx context.Context, id uint64) (*model.Player, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Player), args.Error(1)
}

func (m *MockPlayerRepository) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Player), args.Error(1)
}

func (m *MockPlayerRepository) List(ctx context.Context) ([]model.Player, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Player), args.Error(1)
}

func (m *MockPlayerRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]model.Player), args.Get(1).(int64), args.Error(2)
}

func (m *MockPlayerRepository) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, status *model.VerificationStatus) ([]model.Player, int64, error) {
	args := m.Called(ctx, page, pageSize, keyword, status)
	return args.Get(0).([]model.Player), args.Get(1).(int64), args.Error(2)
}

func (m *MockPlayerRepository) Update(ctx context.Context, player *model.Player) error {
	args := m.Called(ctx, player)
	return args.Error(0)
}

func (m *MockPlayerRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPlayerRepository) BatchUpdateRank(ctx context.Context, ids []uint64, rank string) (int64, error) {
	args := m.Called(ctx, ids, rank)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPlayerRepository) BatchUpdateHourlyRate(ctx context.Context, ids []uint64, rateCents int64) (int64, error) {
	args := m.Called(ctx, ids, rateCents)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPlayerRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, status model.VerificationStatus) (int64, error) {
	args := m.Called(ctx, ids, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPlayerRepository) BatchDelete(ctx context.Context, ids []uint64) (int64, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPlayerRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.Player, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Player), args.Error(1)
}

func (m *MockPlayerRepository) ListFeatured(ctx context.Context, limit int, status *model.VerificationStatus) ([]model.Player, int64, error) {
	args := m.Called(ctx, limit, status)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.Player), args.Get(1).(int64), args.Error(2)
}

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) Get(ctx context.Context, id uint64) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepository) Update(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrderRepository) UpdateWithCondition(ctx context.Context, orderID uint64, expectedStatus model.OrderStatus, updates map[string]any) (bool, error) {
	args := m.Called(ctx, orderID, expectedStatus, updates)
	return args.Bool(0), args.Error(1)
}

func (m *MockOrderRepository) List(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.Order), args.Get(1).(int64), args.Error(2)
}

type MockServiceItemRepository struct {
	mock.Mock
}

func (m *MockServiceItemRepository) Create(ctx context.Context, item *model.ServiceItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockServiceItemRepository) Get(ctx context.Context, id uint64) (*model.ServiceItem, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ServiceItem), args.Error(1)
}

func (m *MockServiceItemRepository) GetByCode(ctx context.Context, itemCode string) (*model.ServiceItem, error) {
	args := m.Called(ctx, itemCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ServiceItem), args.Error(1)
}

func (m *MockServiceItemRepository) List(ctx context.Context, opts repository.ServiceItemListOptions) ([]model.ServiceItem, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.ServiceItem), args.Get(1).(int64), args.Error(2)
}

func (m *MockServiceItemRepository) Update(ctx context.Context, item *model.ServiceItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockServiceItemRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockServiceItemRepository) BatchDelete(ctx context.Context, ids []uint64) (int64, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockServiceItemRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, isActive bool) error {
	args := m.Called(ctx, ids, isActive)
	return args.Error(0)
}

func (m *MockServiceItemRepository) BatchUpdatePrice(ctx context.Context, ids []uint64, basePriceCents int64) error {
	args := m.Called(ctx, ids, basePriceCents)
	return args.Error(0)
}

func (m *MockServiceItemRepository) BatchUpdateCommission(ctx context.Context, ids []uint64, commissionRate float64) error {
	args := m.Called(ctx, ids, commissionRate)
	return args.Error(0)
}

func (m *MockServiceItemRepository) GetGifts(ctx context.Context, page, pageSize int) ([]model.ServiceItem, int64, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.ServiceItem), args.Get(1).(int64), args.Error(2)
}

func (m *MockServiceItemRepository) GetGameServices(ctx context.Context, gameID uint64, subCategory *model.ServiceItemSubCategory) ([]model.ServiceItem, error) {
	args := m.Called(ctx, gameID, subCategory)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.ServiceItem), args.Error(1)
}

type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.Payment), args.Get(1).(int64), args.Error(2)
}

func (m *MockPaymentRepository) Get(ctx context.Context, id uint64) (*model.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

func (m *MockPaymentRepository) GetWithRelations(ctx context.Context, id uint64) (*model.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

func (m *MockPaymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPaymentRepository) GetByOrderID(ctx context.Context, orderID uint64) ([]model.Payment, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Payment), args.Error(1)
}

func (m *MockPaymentRepository) GetByRequestID(ctx context.Context, requestID string) (*model.Payment, error) {
	args := m.Called(ctx, requestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) List(ctx context.Context) ([]model.Permission, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Permission, int64, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.Permission), args.Get(1).(int64), args.Error(2)
}

func (m *MockPermissionRepository) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword, method, group string, isSystem *bool) ([]model.Permission, int64, error) {
	args := m.Called(ctx, page, pageSize, keyword, method, group, isSystem)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.Permission), args.Get(1).(int64), args.Error(2)
}

func (m *MockPermissionRepository) ListByGroup(ctx context.Context) (map[string][]model.Permission, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string][]model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) ListGroups(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockPermissionRepository) Get(ctx context.Context, id uint64) (*model.Permission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) GetByResource(ctx context.Context, resource, action string) (*model.Permission, error) {
	args := m.Called(ctx, resource, action)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) GetByCode(ctx context.Context, code string) (*model.Permission, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) GetByMethodAndPath(ctx context.Context, method, path string) (*model.Permission, error) {
	args := m.Called(ctx, method, path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) Create(ctx context.Context, perm *model.Permission) error {
	args := m.Called(ctx, perm)
	return args.Error(0)
}

func (m *MockPermissionRepository) Update(ctx context.Context, perm *model.Permission) error {
	args := m.Called(ctx, perm)
	return args.Error(0)
}

func (m *MockPermissionRepository) UpsertByMethodPath(ctx context.Context, perm *model.Permission) error {
	args := m.Called(ctx, perm)
	return args.Error(0)
}

func (m *MockPermissionRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPermissionRepository) ListByRoleID(ctx context.Context, roleID uint64) ([]model.Permission, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) ListByUserID(ctx context.Context, userID uint64) ([]model.Permission, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) ListWithChildren(ctx context.Context) ([]model.Permission, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) GetWithChildren(ctx context.Context, id uint64) (*model.Permission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) CountRoleReferences(ctx context.Context, permissionID uint64) (int64, error) {
	args := m.Called(ctx, permissionID)
	return args.Get(0).(int64), args.Error(1)
}

type MockMenuRepository struct {
	mock.Mock
}

func (m *MockMenuRepository) Create(ctx context.Context, menu *model.Menu) error {
	args := m.Called(ctx, menu)
	return args.Error(0)
}

func (m *MockMenuRepository) Update(ctx context.Context, menu *model.Menu) error {
	args := m.Called(ctx, menu)
	return args.Error(0)
}

func (m *MockMenuRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMenuRepository) Get(ctx context.Context, id uint64) (*model.Menu, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Menu), args.Error(1)
}

func (m *MockMenuRepository) List(ctx context.Context, parentID *uint64) ([]model.Menu, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Menu), args.Error(1)
}

func (m *MockMenuRepository) ListPaged(ctx context.Context, page, pageSize int, parentID *uint64) ([]model.Menu, int64, error) {
	args := m.Called(ctx, page, pageSize, parentID)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.Menu), args.Get(1).(int64), args.Error(2)
}

func (m *MockMenuRepository) ListByPermission(ctx context.Context, codes []string) ([]model.Menu, error) {
	args := m.Called(ctx, codes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Menu), args.Error(1)
}

func (m *MockMenuRepository) HasChildren(ctx context.Context, parentID uint64) (bool, error) {
	args := m.Called(ctx, parentID)
	return args.Bool(0), args.Error(1)
}

type MockStatsRepository struct {
	mock.Mock
}

func (m *MockStatsRepository) Dashboard(ctx context.Context) (repository.Dashboard, error) {
	args := m.Called(ctx)
	return args.Get(0).(repository.Dashboard), args.Error(1)
}

func (m *MockStatsRepository) RevenueTrend(ctx context.Context, days int) ([]repository.DateValue, error) {
	args := m.Called(ctx, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.DateValue), args.Error(1)
}

func (m *MockStatsRepository) UserGrowth(ctx context.Context, days int) ([]repository.DateValue, error) {
	args := m.Called(ctx, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.DateValue), args.Error(1)
}

func (m *MockStatsRepository) OrdersByStatus(ctx context.Context) (map[string]int64, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *MockStatsRepository) TopPlayers(ctx context.Context, limit int) ([]repository.PlayerTop, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.PlayerTop), args.Error(1)
}

func (m *MockStatsRepository) AuditOverview(ctx context.Context, from, to *time.Time) (map[string]int64, map[string]int64, error) {
	args := m.Called(ctx, from, to)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(map[string]int64), args.Get(1).(map[string]int64), args.Error(2)
}

func (m *MockStatsRepository) AuditTrend(ctx context.Context, from, to *time.Time, entity, action string) ([]repository.DateValue, error) {
	args := m.Called(ctx, from, to, entity, action)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.DateValue), args.Error(1)
}

// MockWalletRepository is a mock implementation of WalletRepository
type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) Save(ctx context.Context, wallet *model.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) SaveWithOptimisticLock(ctx context.Context, wallet *model.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) UpdateBalanceWithLock(ctx context.Context, userID uint64, delta int64, maxRetries int) (*model.Wallet, error) {
	args := m.Called(ctx, userID, delta, maxRetries)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

func (m *MockWalletRepository) Get(ctx context.Context, userID uint64) (*model.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

func (m *MockWalletRepository) GetByUserID(ctx context.Context, userID uint64) (*model.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

// MockGameCategoryRepository is a mock implementation of GameCategoryRepository
type MockGameCategoryRepository struct {
	mock.Mock
}

func (m *MockGameCategoryRepository) Create(ctx context.Context, category *model.GameCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockGameCategoryRepository) Get(ctx context.Context, id uint64) (*model.GameCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.GameCategory), args.Error(1)
}

func (m *MockGameCategoryRepository) List(ctx context.Context, opts repository.GameCategoryListOptions) ([]*model.GameCategory, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]*model.GameCategory), args.Get(1).(int64), args.Error(2)
}

func (m *MockGameCategoryRepository) Update(ctx context.Context, category *model.GameCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockGameCategoryRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockGameCategoryRepository) GetByName(ctx context.Context, name string) (*model.GameCategory, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.GameCategory), args.Error(1)
}

func (m *MockGameCategoryRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, isActive bool) error {
	args := m.Called(ctx, ids, isActive)
	return args.Error(0)
}

func (m *MockGameCategoryRepository) BatchDelete(ctx context.Context, ids []uint64) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockGameCategoryRepository) Exists(ctx context.Context, id uint64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockGameCategoryRepository) CountGames(ctx context.Context, categoryID uint64) (int64, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockGameCategoryRepository) CountServiceItems(ctx context.Context, categoryID uint64) (int64, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).(int64), args.Error(1)
}

// MockCache is a mock implementation of Cache
type MockCache struct {
	mock.Mock
}

// NewMockCacheWithDefaults creates a MockCache with default expectations for all cache operations
func NewMockCacheWithDefaults() *MockCache {
	c := &MockCache{}
	// Set up default expectations for common cache operations
	c.On("Get", mock.Anything, mock.Anything).Return("", false, nil).Maybe()
	c.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	c.On("Delete", mock.Anything, mock.Anything).Return(nil).Maybe()
	return c
}

func (m *MockCache) Get(ctx context.Context, key string) (string, bool, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Bool(1), args.Error(2)
}

func (m *MockCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockCache) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCache) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCache) GetRedisClient() interface{} {
	args := m.Called()
	return args.Get(0)
}

// Helper function to create a hashed password
func hashPasswordForTest(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}

// ============================================================================
// Game Management Tests
// ============================================================================

func TestAdminService_CreateGame(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		input       CreateGameInput
		setupMocks  func(*MockGameRepository, *MockCache)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - create game",
			input: CreateGameInput{
				Key:         "test-game",
				Name:        "Test Game",
				Category:    "moba",
				IconURL:     "https://example.com/icon.png",
				Description: "A test game",
			},
			setupMocks: func(m *MockGameRepository, c *MockCache) {
				m.On("Create", ctx, mock.AnythingOfType("*model.Game")).Return(nil)
				c.On("Delete", ctx, "admin:games").Return(nil)
			},
			wantErr: false,
		},
		{
			name: "validation error - empty key",
			input: CreateGameInput{
				Key:  "",
				Name: "Test Game",
			},
			setupMocks:  func(m *MockGameRepository, c *MockCache) {},
			wantErr:     true,
			errContains: "validation",
		},
		{
			name: "validation error - empty name",
			input: CreateGameInput{
				Key:  "test-game",
				Name: "",
			},
			setupMocks:  func(m *MockGameRepository, c *MockCache) {},
			wantErr:     true,
			errContains: "validation",
		},
		{
			name: "repository error",
			input: CreateGameInput{
				Key:  "test-game",
				Name: "Test Game",
			},
			setupMocks: func(m *MockGameRepository, c *MockCache) {
				m.On("Create", ctx, mock.AnythingOfType("*model.Game")).Return(errors.New("db error"))
			},
			wantErr:     true,
			errContains: "create game",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			games := &MockGameRepository{}
			cache := NewMockCacheWithDefaults()
			tt.setupMocks(games, cache)

			svc := NewAdminService(
				games, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				cache,
			)

			game, err := svc.CreateGame(ctx, tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, game)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, game)
				assert.Equal(t, tt.input.Key, game.Key)
				assert.Equal(t, tt.input.Name, game.Name)
			}

			games.AssertExpectations(t)
		})
	}
}

func TestAdminService_UpdateGame(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		gameID      uint64
		input       UpdateGameInput
		setupMocks  func(*MockGameRepository, *MockCache)
		wantErr     bool
		errContains string
	}{
		{
			name:   "success - update game",
			gameID: 1,
			input: UpdateGameInput{
				Key:         "updated-game",
				Name:        "Updated Game",
				Category:    "fps",
				IconURL:     "https://example.com/new-icon.png",
				Description: "Updated description",
			},
			setupMocks: func(m *MockGameRepository, c *MockCache) {
				m.On("Get", ctx, uint64(1)).Return(&model.Game{
					Base: model.Base{ID: 1, ExtJSON: "{}"},
					Key:  "old-game",
					Name: "Old Game",
				}, nil)
				m.On("Update", ctx, mock.AnythingOfType("*model.Game")).Return(nil)
				c.On("Delete", ctx, "admin:games").Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "game not found",
			gameID: 999,
			input: UpdateGameInput{
				Key:  "test-game",
				Name: "Test Game",
			},
			setupMocks: func(m *MockGameRepository, c *MockCache) {
				m.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name:   "validation error - empty key",
			gameID: 1,
			input: UpdateGameInput{
				Key:  "",
				Name: "Test Game",
			},
			setupMocks:  func(m *MockGameRepository, c *MockCache) {},
			wantErr:     true,
			errContains: "validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			games := &MockGameRepository{}
			cache := NewMockCacheWithDefaults()
			tt.setupMocks(games, cache)

			svc := NewAdminService(
				games, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				cache,
			)

			game, err := svc.UpdateGame(ctx, tt.gameID, tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, game)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, game)
			}

			games.AssertExpectations(t)
		})
	}
}

func TestAdminService_DeleteGame(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		gameID      uint64
		setupMocks  func(*MockGameRepository, *MockCache)
		wantErr     bool
		errContains string
	}{
		{
			name:   "success - delete game",
			gameID: 1,
			setupMocks: func(m *MockGameRepository, c *MockCache) {
				m.On("Delete", ctx, uint64(1)).Return(nil)
				c.On("Delete", ctx, "admin:games").Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "game not found",
			gameID: 999,
			setupMocks: func(m *MockGameRepository, c *MockCache) {
				m.On("Delete", ctx, uint64(999)).Return(repository.ErrNotFound)
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name:   "repository error",
			gameID: 1,
			setupMocks: func(m *MockGameRepository, c *MockCache) {
				m.On("Delete", ctx, uint64(1)).Return(errors.New("db error"))
			},
			wantErr:     true,
			errContains: "delete game",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			games := &MockGameRepository{}
			cache := NewMockCacheWithDefaults()
			tt.setupMocks(games, cache)

			svc := NewAdminService(
				games, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				cache,
			)

			err := svc.DeleteGame(ctx, tt.gameID)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}

			games.AssertExpectations(t)
		})
	}
}

func TestAdminService_ListGames(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		setupMocks func(*MockGameRepository, *MockCache)
		wantErr    bool
		wantCount  int
	}{
		{
			name: "success - list games",
			setupMocks: func(m *MockGameRepository, c *MockCache) {
				// Cache miss, then fetch from DB
				c.On("Get", ctx, "admin:games").Return("", false, nil)
				m.On("List", ctx).Return([]model.Game{
					{Base: model.Base{ID: 1, ExtJSON: "{}"}, Key: "game1", Name: "Game 1"},
					{Base: model.Base{ID: 2, ExtJSON: "{}"}, Key: "game2", Name: "Game 2"},
				}, nil)
				c.On("Set", ctx, "admin:games", mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil)
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name: "repository error",
			setupMocks: func(m *MockGameRepository, c *MockCache) {
				c.On("Get", ctx, "admin:games").Return("", false, nil)
				m.On("List", ctx).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			games := &MockGameRepository{}
			cache := NewMockCacheWithDefaults()
			tt.setupMocks(games, cache)

			svc := NewAdminService(
				games, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				cache,
			)

			result, err := svc.ListGames(ctx)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.wantCount)
			}

			games.AssertExpectations(t)
		})
	}
}

func TestAdminService_BatchDeleteGames(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		ids         []uint64
		setupMocks  func(*MockGameRepository, *MockCache)
		wantErr     bool
		errContains string
		wantDeleted int64
	}{
		{
			name: "success - batch delete",
			ids:  []uint64{1, 2, 3},
			setupMocks: func(m *MockGameRepository, c *MockCache) {
				m.On("BatchDelete", ctx, []uint64{1, 2, 3}).Return(int64(3), nil)
				c.On("Delete", ctx, "admin:games").Return(nil)
			},
			wantErr:     false,
			wantDeleted: 3,
		},
		{
			name:        "empty ids",
			ids:         []uint64{},
			setupMocks:  func(m *MockGameRepository, c *MockCache) {},
			wantErr:     true,
			errContains: "no game ids",
		},
		{
			name: "repository error",
			ids:  []uint64{1, 2},
			setupMocks: func(m *MockGameRepository, c *MockCache) {
				m.On("BatchDelete", ctx, []uint64{1, 2}).Return(int64(0), errors.New("db error"))
			},
			wantErr:     true,
			errContains: "batch delete games",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			games := &MockGameRepository{}
			cache := NewMockCacheWithDefaults()
			tt.setupMocks(games, cache)

			svc := NewAdminService(
				games, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				cache,
			)

			deleted, err := svc.BatchDeleteGames(ctx, tt.ids)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantDeleted, deleted)
			}

			games.AssertExpectations(t)
		})
	}
}

// ============================================================================
// User Management Tests
// ============================================================================

func TestAdminService_CreateUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		input       CreateUserInput
		setupMocks  func(*MockUserRepository, *MockWalletRepository, *MockRoleRepository, *MockCache)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - create user",
			input: CreateUserInput{
				Name:     "Test User",
				Phone:    "13800138000",
				Email:    "test@example.com",
				Password: "Test123!@#",
				Role:     model.RoleUser,
				Status:   model.UserStatusActive,
			},
			setupMocks: func(u *MockUserRepository, w *MockWalletRepository, r *MockRoleRepository, c *MockCache) {
				u.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(nil)
				w.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).Return(nil)
				r.On("GetBySlug", ctx, string(model.RoleSlugUser)).Return(&model.RoleModel{
					Base: model.Base{ID: 1, ExtJSON: "{}"},
					Slug: string(model.RoleSlugUser),
				}, nil)
				r.On("AssignToUser", ctx, mock.AnythingOfType("uint64"), mock.AnythingOfType("[]uint64")).Return(nil)
				c.On("Delete", ctx, "admin:users").Return(nil)
			},
			wantErr: false,
		},
		{
			name: "validation error - empty name",
			input: CreateUserInput{
				Name:     "",
				Phone:    "13800138000",
				Email:    "test@example.com",
				Password: "Test123!@#",
				Role:     model.RoleUser,
				Status:   model.UserStatusActive,
			},
			setupMocks:  func(u *MockUserRepository, w *MockWalletRepository, r *MockRoleRepository, c *MockCache) {},
			wantErr:     true,
			errContains: "validation",
		},
		{
			name: "validation error - weak password",
			input: CreateUserInput{
				Name:     "Test User",
				Phone:    "13800138000",
				Email:    "test@example.com",
				Password: "weak",
				Role:     model.RoleUser,
				Status:   model.UserStatusActive,
			},
			setupMocks:  func(u *MockUserRepository, w *MockWalletRepository, r *MockRoleRepository, c *MockCache) {},
			wantErr:     true,
			errContains: "validation",
		},
		{
			name: "validation error - empty role",
			input: CreateUserInput{
				Name:     "Test User",
				Phone:    "13800138000",
				Email:    "test@example.com",
				Password: "Test123!@#",
				Role:     "",
				Status:   model.UserStatusActive,
			},
			setupMocks:  func(u *MockUserRepository, w *MockWalletRepository, r *MockRoleRepository, c *MockCache) {},
			wantErr:     true,
			errContains: "validation",
		},
		{
			name: "repository error",
			input: CreateUserInput{
				Name:     "Test User",
				Phone:    "13800138000",
				Email:    "test@example.com",
				Password: "Test123!@#",
				Role:     model.RoleUser,
				Status:   model.UserStatusActive,
			},
			setupMocks: func(u *MockUserRepository, w *MockWalletRepository, r *MockRoleRepository, c *MockCache) {
				u.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &MockUserRepository{}
			wallets := &MockWalletRepository{}
			roles := &MockRoleRepository{}
			cache := NewMockCacheWithDefaults()
			tt.setupMocks(users, wallets, roles, cache)

			svc := NewAdminService(
				nil, users, nil, nil, nil, roles, nil, nil, nil, nil, wallets, nil,
				cache,
			)

			user, err := svc.CreateUser(ctx, tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.input.Name, user.Name)
				assert.Equal(t, tt.input.Phone, user.Phone)
			}

			users.AssertExpectations(t)
			wallets.AssertExpectations(t)
			roles.AssertExpectations(t)
		})
	}
}

func TestAdminService_UpdateUser(t *testing.T) {
	ctx := context.Background()

	// Helper to create a fresh existingUser for each test
	newExistingUser := func() *model.User {
		return &model.User{
			Base:         model.Base{ID: 1, ExtJSON: "{}"},
			Name:         "Old Name",
			Phone:        "13800138000",
			Email:        "old@example.com",
			Role:         model.RoleUser,
			Status:       model.UserStatusActive,
			PasswordHash: hashPasswordForTest("Old123!@#"),
		}
	}

	tests := []struct {
		name        string
		userID      uint64
		input       UpdateUserInput
		setupMocks  func(*MockUserRepository, *MockRoleRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:   "success - update user",
			userID: 1,
			input: UpdateUserInput{
				Name:   "Updated Name",
				Phone:  "13900139000",
				Email:  "updated@example.com",
				Role:   model.RolePlayer,
				Status: model.UserStatusActive,
			},
			setupMocks: func(u *MockUserRepository, r *MockRoleRepository) {
				u.On("Get", ctx, uint64(1)).Return(newExistingUser(), nil)
				u.On("Update", ctx, mock.AnythingOfType("*model.User")).Return(nil)
				r.On("GetBySlug", ctx, string(model.RoleSlugPlayer)).Return(&model.RoleModel{
					Base: model.Base{ID: 2, ExtJSON: "{}"},
					Slug: string(model.RoleSlugPlayer),
				}, nil)
				r.On("AssignToUser", ctx, mock.AnythingOfType("uint64"), mock.AnythingOfType("[]uint64")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "user not found",
			userID: 999,
			input: UpdateUserInput{
				Name:   "Test Name",
				Role:   model.RoleUser,
				Status: model.UserStatusActive,
			},
			setupMocks: func(u *MockUserRepository, r *MockRoleRepository) {
				u.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name:   "update password",
			userID: 1,
			input: UpdateUserInput{
				Name:     "Updated Name",
				Role:     model.RoleUser,
				Status:   model.UserStatusActive,
				Password: stringPtr("New123!@#"),
			},
			setupMocks: func(u *MockUserRepository, r *MockRoleRepository) {
				u.On("Get", ctx, uint64(1)).Return(newExistingUser(), nil)
				u.On("Update", ctx, mock.AnythingOfType("*model.User")).Return(nil)
				r.On("GetBySlug", ctx, string(model.RoleSlugUser)).Return(&model.RoleModel{
					Base: model.Base{ID: 1, ExtJSON: "{}"},
					Slug: string(model.RoleSlugUser),
				}, nil)
				r.On("AssignToUser", ctx, mock.AnythingOfType("uint64"), mock.AnythingOfType("[]uint64")).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &MockUserRepository{}
			roles := &MockRoleRepository{}
			tt.setupMocks(users, roles)

			svc := NewAdminService(
				nil, users, nil, nil, nil, roles, nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			user, err := svc.UpdateUser(ctx, tt.userID, tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, user)
			}

			users.AssertExpectations(t)
			roles.AssertExpectations(t)
		})
	}
}

func TestAdminService_DeleteUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		userID      uint64
		setupMocks  func(*MockUserRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:   "success - delete user",
			userID: 1,
			setupMocks: func(u *MockUserRepository) {
				u.On("Delete", ctx, uint64(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "user not found",
			userID: 999,
			setupMocks: func(u *MockUserRepository) {
				u.On("Delete", ctx, uint64(999)).Return(repository.ErrNotFound)
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &MockUserRepository{}
			tt.setupMocks(users)

			svc := NewAdminService(
				nil, users, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			err := svc.DeleteUser(ctx, tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}

			users.AssertExpectations(t)
		})
	}
}

func TestAdminService_ListUsers(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		setupMocks func(*MockUserRepository)
		wantErr    bool
		wantCount  int
	}{
		{
			name: "success - list users",
			setupMocks: func(u *MockUserRepository) {
				u.On("List", ctx).Return([]model.User{
					{Base: model.Base{ID: 1, ExtJSON: "{}"}, Name: "User 1"},
					{Base: model.Base{ID: 2, ExtJSON: "{}"}, Name: "User 2"},
				}, nil)
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name: "repository error",
			setupMocks: func(u *MockUserRepository) {
				u.On("List", ctx).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &MockUserRepository{}
			tt.setupMocks(users)

			svc := NewAdminService(
				nil, users, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			result, err := svc.ListUsers(ctx)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.wantCount)
			}

			users.AssertExpectations(t)
		})
	}
}

func TestAdminService_GetUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		userID      uint64
		setupMocks  func(*MockUserRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:   "success - get user",
			userID: 1,
			setupMocks: func(u *MockUserRepository) {
				u.On("Get", ctx, uint64(1)).Return(&model.User{
					Base: model.Base{ID: 1, ExtJSON: "{}"},
					Name: "Test User",
				}, nil)
			},
			wantErr: false,
		},
		{
			name:   "user not found",
			userID: 999,
			setupMocks: func(u *MockUserRepository) {
				u.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &MockUserRepository{}
			tt.setupMocks(users)

			svc := NewAdminService(
				nil, users, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			user, err := svc.GetUser(ctx, tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.NotNil(t, user)
			}

			users.AssertExpectations(t)
		})
	}
}

func TestAdminService_UpdateUserStatus(t *testing.T) {
	ctx := context.Background()

	existingUser := &model.User{
		Base:   model.Base{ID: 1, ExtJSON: "{}"},
		Name:   "Test User",
		Role:   model.RoleUser,
		Status: model.UserStatusActive,
	}

	tests := []struct {
		name        string
		userID      uint64
		status      model.UserStatus
		setupMocks  func(*MockUserRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:   "success - update status",
			userID: 1,
			status: model.UserStatusBanned,
			setupMocks: func(u *MockUserRepository) {
				u.On("Get", ctx, uint64(1)).Return(existingUser, nil)
				u.On("Update", ctx, mock.MatchedBy(func(u *model.User) bool {
					return u.Status == model.UserStatusBanned
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "user not found",
			userID: 999,
			status: model.UserStatusBanned,
			setupMocks: func(u *MockUserRepository) {
				u.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &MockUserRepository{}
			tt.setupMocks(users)

			svc := NewAdminService(
				nil, users, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			user, err := svc.UpdateUserStatus(ctx, tt.userID, tt.status)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.status, user.Status)
			}

			users.AssertExpectations(t)
		})
	}
}

func TestAdminService_UpdateUserRole(t *testing.T) {
	ctx := context.Background()

	existingUser := &model.User{
		Base:   model.Base{ID: 1, ExtJSON: "{}"},
		Name:   "Test User",
		Role:   model.RoleUser,
		Status: model.UserStatusActive,
	}

	tests := []struct {
		name        string
		userID      uint64
		role        model.Role
		setupMocks  func(*MockUserRepository, *MockRoleRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:   "success - update role",
			userID: 1,
			role:   model.RolePlayer,
			setupMocks: func(u *MockUserRepository, r *MockRoleRepository) {
				u.On("Get", ctx, uint64(1)).Return(existingUser, nil)
				u.On("Update", ctx, mock.MatchedBy(func(u *model.User) bool {
					return u.Role == model.RolePlayer
				})).Return(nil)
				r.On("GetBySlug", ctx, string(model.RoleSlugPlayer)).Return(&model.RoleModel{
					Base: model.Base{ID: 2, ExtJSON: "{}"},
					Slug: string(model.RoleSlugPlayer),
				}, nil)
				r.On("AssignToUser", ctx, mock.AnythingOfType("uint64"), mock.AnythingOfType("[]uint64")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "user not found",
			userID: 999,
			role:   model.RolePlayer,
			setupMocks: func(u *MockUserRepository, r *MockRoleRepository) {
				u.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &MockUserRepository{}
			roles := &MockRoleRepository{}
			tt.setupMocks(users, roles)

			svc := NewAdminService(
				nil, users, nil, nil, nil, roles, nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			user, err := svc.UpdateUserRole(ctx, tt.userID, tt.role)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.role, user.Role)
			}

			users.AssertExpectations(t)
			roles.AssertExpectations(t)
		})
	}
}

func TestAdminService_GetUserStats(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		setupMocks func(*MockUserRepository)
		wantErr    bool
		wantStats  *UserStatsResponse
	}{
		{
			name: "success - get stats",
			setupMocks: func(u *MockUserRepository) {
				u.On("Count", ctx, repository.UserListOptions{}).Return(100, nil)
				u.On("Count", ctx, mock.MatchedBy(func(opts repository.UserListOptions) bool {
					return len(opts.Roles) == 1 && opts.Roles[0] == model.RoleUser
				})).Return(60, nil)
				u.On("Count", ctx, mock.MatchedBy(func(opts repository.UserListOptions) bool {
					return len(opts.Roles) == 1 && opts.Roles[0] == model.RolePlayer
				})).Return(30, nil)
				u.On("Count", ctx, mock.MatchedBy(func(opts repository.UserListOptions) bool {
					return len(opts.Roles) == 1 && opts.Roles[0] == model.RoleAdmin
				})).Return(10, nil)
				u.On("Count", ctx, mock.MatchedBy(func(opts repository.UserListOptions) bool {
					return len(opts.Statuses) == 1 && opts.Statuses[0] == model.UserStatusActive
				})).Return(80, nil)
				u.On("Count", ctx, mock.MatchedBy(func(opts repository.UserListOptions) bool {
					return len(opts.Statuses) == 1 && opts.Statuses[0] == model.UserStatusBanned
				})).Return(15, nil)
				u.On("Count", ctx, mock.MatchedBy(func(opts repository.UserListOptions) bool {
					return len(opts.Statuses) == 1 && opts.Statuses[0] == model.UserStatusSuspended
				})).Return(5, nil)
				u.On("Count", ctx, mock.MatchedBy(func(opts repository.UserListOptions) bool {
					return opts.DateFrom != nil
				})).Return(20, nil)
			},
			wantErr: false,
			wantStats: &UserStatsResponse{
				Total:               100,
				ByRole:              map[string]int{"user": 60, "player": 30, "admin": 10},
				ByStatus:            map[string]int{"active": 80, "banned": 15, "suspended": 5},
				RecentRegistrations: 20,
			},
		},
		{
			name: "repository error",
			setupMocks: func(u *MockUserRepository) {
				u.On("Count", ctx, mock.AnythingOfType("repository.UserListOptions")).Return(0, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &MockUserRepository{}
			tt.setupMocks(users)

			svc := NewAdminService(
				nil, users, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			stats, err := svc.GetUserStats(ctx)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, stats)
				assert.Equal(t, tt.wantStats.Total, stats.Total)
				assert.Equal(t, tt.wantStats.RecentRegistrations, stats.RecentRegistrations)
			}

			users.AssertExpectations(t)
		})
	}
}

// ============================================================================
// Player Management Tests
// ============================================================================

func TestAdminService_CreatePlayer(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		input       CreatePlayerInput
		setupMocks  func(*MockPlayerRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - create player",
			input: CreatePlayerInput{
				UserID:             1,
				Nickname:           "TestPlayer",
				Bio:                "Test bio",
				HourlyRateCents:    5000,
				MainGameID:         1,
				VerificationStatus: model.VerificationPending,
			},
			setupMocks: func(p *MockPlayerRepository) {
				p.On("Create", ctx, mock.AnythingOfType("*model.Player")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "validation error - zero user id",
			input: CreatePlayerInput{
				UserID:             0,
				Nickname:           "TestPlayer",
				VerificationStatus: model.VerificationPending,
			},
			setupMocks:  func(p *MockPlayerRepository) {},
			wantErr:     true,
			errContains: "validation",
		},
		{
			name: "validation error - empty verification status",
			input: CreatePlayerInput{
				UserID:             1,
				Nickname:           "TestPlayer",
				VerificationStatus: "",
			},
			setupMocks:  func(p *MockPlayerRepository) {},
			wantErr:     true,
			errContains: "validation",
		},
		{
			name: "repository error",
			input: CreatePlayerInput{
				UserID:             1,
				Nickname:           "TestPlayer",
				VerificationStatus: model.VerificationPending,
			},
			setupMocks: func(p *MockPlayerRepository) {
				p.On("Create", ctx, mock.AnythingOfType("*model.Player")).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			tt.setupMocks(players)

			svc := NewAdminService(
				nil, nil, players, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			player, err := svc.CreatePlayer(ctx, tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, player)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, player)
			}

			players.AssertExpectations(t)
		})
	}
}

func TestAdminService_UpdatePlayer(t *testing.T) {
	ctx := context.Background()

	existingPlayer := &model.Player{
		Base:               model.Base{ID: 1, ExtJSON: "{}"},
		UserID:             1,
		Nickname:           "OldNickname",
		HourlyRateCents:    3000,
		VerificationStatus: model.VerificationPending,
	}

	tests := []struct {
		name        string
		playerID    uint64
		input       UpdatePlayerInput
		setupMocks  func(*MockPlayerRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:     "success - update player",
			playerID: 1,
			input: UpdatePlayerInput{
				Nickname:           "UpdatedNickname",
				Bio:                "Updated bio",
				HourlyRateCents:    5000,
				MainGameID:         2,
				VerificationStatus: model.VerificationVerified,
			},
			setupMocks: func(p *MockPlayerRepository) {
				p.On("Get", ctx, uint64(1)).Return(existingPlayer, nil)
				p.On("Update", ctx, mock.MatchedBy(func(p *model.Player) bool {
					return p.Nickname == "UpdatedNickname"
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "player not found",
			playerID: 999,
			input: UpdatePlayerInput{
				Nickname:           "Test",
				VerificationStatus: model.VerificationPending,
			},
			setupMocks: func(p *MockPlayerRepository) {
				p.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			tt.setupMocks(players)

			svc := NewAdminService(
				nil, nil, players, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			player, err := svc.UpdatePlayer(ctx, tt.playerID, tt.input)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, player)
			}

			players.AssertExpectations(t)
		})
	}
}

func TestAdminService_DeletePlayer(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		playerID    uint64
		setupMocks  func(*MockPlayerRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:     "success - delete player",
			playerID: 1,
			setupMocks: func(p *MockPlayerRepository) {
				p.On("Delete", ctx, uint64(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "player not found",
			playerID: 999,
			setupMocks: func(p *MockPlayerRepository) {
				p.On("Delete", ctx, uint64(999)).Return(repository.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			tt.setupMocks(players)

			svc := NewAdminService(
				nil, nil, players, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			err := svc.DeletePlayer(ctx, tt.playerID)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			players.AssertExpectations(t)
		})
	}
}

func TestAdminService_ListPlayers(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		setupMocks func(*MockPlayerRepository)
		wantErr    bool
		wantCount  int
	}{
		{
			name: "success - list players",
			setupMocks: func(p *MockPlayerRepository) {
				p.On("List", ctx).Return([]model.Player{
					{Base: model.Base{ID: 1, ExtJSON: "{}"}, Nickname: "Player 1"},
					{Base: model.Base{ID: 2, ExtJSON: "{}"}, Nickname: "Player 2"},
				}, nil)
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name: "repository error",
			setupMocks: func(p *MockPlayerRepository) {
				p.On("List", ctx).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			tt.setupMocks(players)

			svc := NewAdminService(
				nil, nil, players, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			result, err := svc.ListPlayers(ctx)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.wantCount)
			}

			players.AssertExpectations(t)
		})
	}
}

func TestAdminService_BatchUpdatePlayerStatus(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		ids         []uint64
		status      model.VerificationStatus
		setupMocks  func(*MockPlayerRepository)
		wantErr     bool
		errContains string
		wantUpdated int64
	}{
		{
			name:   "success - batch update",
			ids:    []uint64{1, 2, 3},
			status: model.VerificationVerified,
			setupMocks: func(p *MockPlayerRepository) {
				p.On("BatchUpdateStatus", ctx, []uint64{1, 2, 3}, model.VerificationVerified).Return(int64(3), nil)
			},
			wantErr:     false,
			wantUpdated: 3,
		},
		{
			name:        "empty ids",
			ids:         []uint64{},
			status:      model.VerificationVerified,
			setupMocks:  func(p *MockPlayerRepository) {},
			wantErr:     true,
			errContains: "no player ids",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			tt.setupMocks(players)

			svc := NewAdminService(
				nil, nil, players, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			updated, err := svc.BatchUpdatePlayerStatus(ctx, tt.ids, tt.status)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantUpdated, updated)
			}

			players.AssertExpectations(t)
		})
	}
}

func TestAdminService_BatchDeletePlayers(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		ids         []uint64
		setupMocks  func(*MockPlayerRepository)
		wantErr     bool
		errContains string
		wantDeleted int64
	}{
		{
			name: "success - batch delete",
			ids:  []uint64{1, 2, 3},
			setupMocks: func(p *MockPlayerRepository) {
				p.On("BatchDelete", ctx, []uint64{1, 2, 3}).Return(int64(3), nil)
			},
			wantErr:     false,
			wantDeleted: 3,
		},
		{
			name:        "empty ids",
			ids:         []uint64{},
			setupMocks:  func(p *MockPlayerRepository) {},
			wantErr:     true,
			errContains: "no player ids",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			tt.setupMocks(players)

			svc := NewAdminService(
				nil, nil, players, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			deleted, err := svc.BatchDeletePlayers(ctx, tt.ids)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantDeleted, deleted)
			}

			players.AssertExpectations(t)
		})
	}
}

// ============================================================================
// Order Management Tests
// ============================================================================

func TestAdminService_CreateOrder(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name        string
		input       CreateOrderInput
		setupMocks  func(*MockOrderRepository, *MockServiceItemRepository, *MockPlayerRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - create order",
			input: CreateOrderInput{
				UserID:          1,
				GameID:          1,
				ItemID:          1,
				TotalPriceCents: 5000,
				Currency:        model.CurrencyCNY,
				Title:           "Test Order",
			},
			setupMocks: func(o *MockOrderRepository, s *MockServiceItemRepository, p *MockPlayerRepository) {
				s.On("Get", ctx, uint64(1)).Return(&model.ServiceItem{
					ID:       1,
					IsActive: true,
					GameID:   nil,
				}, nil)
				o.On("Create", ctx, mock.AnythingOfType("*model.Order")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "validation error - zero user id",
			input: CreateOrderInput{
				UserID:          0,
				GameID:          1,
				ItemID:          1,
				TotalPriceCents: 5000,
				Currency:        model.CurrencyCNY,
			},
			setupMocks:  func(o *MockOrderRepository, s *MockServiceItemRepository, p *MockPlayerRepository) {},
			wantErr:     true,
			errContains: "validation",
		},
		{
			name: "validation error - invalid currency",
			input: CreateOrderInput{
				UserID:          1,
				GameID:          1,
				ItemID:          1,
				TotalPriceCents: 5000,
				Currency:        "INVALID",
			},
			setupMocks:  func(o *MockOrderRepository, s *MockServiceItemRepository, p *MockPlayerRepository) {},
			wantErr:     true,
			errContains: "validation",
		},
		{
			name: "service item not found",
			input: CreateOrderInput{
				UserID:          1,
				GameID:          1,
				ItemID:          999,
				TotalPriceCents: 5000,
				Currency:        model.CurrencyCNY,
			},
			setupMocks: func(o *MockOrderRepository, s *MockServiceItemRepository, p *MockPlayerRepository) {
				s.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			wantErr:     true,
			errContains: "服务项目不存在",
		},
		{
			name: "service item not active",
			input: CreateOrderInput{
				UserID:          1,
				GameID:          1,
				ItemID:          1,
				TotalPriceCents: 5000,
				Currency:        model.CurrencyCNY,
			},
			setupMocks: func(o *MockOrderRepository, s *MockServiceItemRepository, p *MockPlayerRepository) {
				s.On("Get", ctx, uint64(1)).Return(&model.ServiceItem{
					ID:       1,
					IsActive: false,
				}, nil)
			},
			wantErr:     true,
			errContains: "服务项目已停用",
		},
		{
			name: "validation error - scheduled end before start",
			input: CreateOrderInput{
				UserID:          1,
				GameID:          1,
				ItemID:          1,
				TotalPriceCents: 5000,
				Currency:        model.CurrencyCNY,
				ScheduledStart:  &now,
				ScheduledEnd:    func() *time.Time { t := now.Add(-1 * time.Hour); return &t }(),
			},
			setupMocks:  func(o *MockOrderRepository, s *MockServiceItemRepository, p *MockPlayerRepository) {},
			wantErr:     true,
			errContains: "validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders := &MockOrderRepository{}
			serviceItems := &MockServiceItemRepository{}
			players := &MockPlayerRepository{}
			tt.setupMocks(orders, serviceItems, players)

			svc := NewAdminService(
				nil, nil, players, orders, nil, nil, serviceItems, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			order, err := svc.CreateOrder(ctx, tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, order)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, order)
			}

			orders.AssertExpectations(t)
			serviceItems.AssertExpectations(t)
		})
	}
}

func TestAdminService_UpdateOrder(t *testing.T) {
	ctx := context.Background()

	existingOrder := &model.Order{
		Base:            model.Base{ID: 1, ExtJSON: "{}"},
		UserID:          1,
		ItemID:          1,
		TotalPriceCents: 5000,
		Currency:        model.CurrencyCNY,
		Status:          model.OrderStatusPending,
		OrderConfig:     "{}",
	}

	tests := []struct {
		name        string
		orderID     uint64
		input       UpdateOrderInput
		setupMocks  func(*MockOrderRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:    "success - update order status",
			orderID: 1,
			input: UpdateOrderInput{
				Status:          model.OrderStatusConfirmed,
				TotalPriceCents: 5000,
				Currency:        model.CurrencyCNY,
			},
			setupMocks: func(o *MockOrderRepository) {
				o.On("Get", ctx, uint64(1)).Return(existingOrder, nil)
				o.On("Update", ctx, mock.MatchedBy(func(o *model.Order) bool {
					return o.Status == model.OrderStatusConfirmed
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "order not found",
			orderID: 999,
			input: UpdateOrderInput{
				Status:          model.OrderStatusConfirmed,
				TotalPriceCents: 5000,
				Currency:        model.CurrencyCNY,
			},
			setupMocks: func(o *MockOrderRepository) {
				o.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name:    "invalid status transition",
			orderID: 1,
			input: UpdateOrderInput{
				Status:          model.OrderStatusPending,
				TotalPriceCents: 5000,
				Currency:        model.CurrencyCNY,
			},
			setupMocks: func(o *MockOrderRepository) {
				o.On("Get", ctx, uint64(1)).Return(&model.Order{
					Base:            model.Base{ID: 1, ExtJSON: "{}"},
					Status:          model.OrderStatusCompleted,
					TotalPriceCents: 5000,
					Currency:        model.CurrencyCNY,
					OrderConfig:     "{}",
				}, nil)
			},
			wantErr:     true,
			errContains: "invalid order status transition",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders := &MockOrderRepository{}
			tt.setupMocks(orders)

			svc := NewAdminService(
				nil, nil, nil, orders, nil, nil, nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			order, err := svc.UpdateOrder(ctx, tt.orderID, tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.NotNil(t, order)
			}

			orders.AssertExpectations(t)
		})
	}
}

// ============================================================================
// WrapError Tests
// ============================================================================

func TestWrapError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		operation string
		wantErr   bool
		errType   string
	}{
		{
			name:      "nil error",
			err:       nil,
			operation: "test operation",
			wantErr:   false,
		},
		{
			name:      "ErrNotFound",
			err:       repository.ErrNotFound,
			operation: "get record",
			wantErr:   true,
			errType:   "not found",
		},
		{
			name:      "ErrValidation",
			err:       ErrValidation,
			operation: "validate input",
			wantErr:   true,
			errType:   "bad request",
		},
		{
			name:      "generic error",
			err:       errors.New("some error"),
			operation: "do something",
			wantErr:   true,
			errType:   "internal error",
		},
		{
			name:      "already apierr",
			err:       apierr.BadRequest("bad input"),
			operation: "validate",
			wantErr:   true,
			errType:   "bad request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapError(tt.err, tt.operation)

			if tt.wantErr {
				require.Error(t, result)
				switch tt.errType {
				case "not found":
					assert.True(t, apierr.IsNotFound(result))
				case "bad request":
					assert.True(t, apierr.IsBadRequest(result))
				case "internal error":
					assert.True(t, apierr.IsInternalError(result))
				}
			} else {
				assert.NoError(t, result)
			}
		})
	}
}

// ============================================================================
// Password Validation Tests
// ============================================================================

func TestValidPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantBool bool
	}{
		{
			name:     "valid password with all requirements",
			password: "Test123!@#",
			wantBool: true,
		},
		{
			name:     "valid password with different special chars",
			password: "Passw0rd+",
			wantBool: true,
		},
		{
			name:     "too short - less than 8 chars",
			password: "Test1!",
			wantBool: false,
		},
		{
			name:     "missing uppercase",
			password: "test123!@#",
			wantBool: false,
		},
		{
			name:     "missing lowercase",
			password: "TEST123!@#",
			wantBool: false,
		},
		{
			name:     "missing digit",
			password: "TestPass!@#",
			wantBool: false,
		},
		{
			name:     "missing special char",
			password: "Test12345",
			wantBool: false,
		},
		{
			name:     "too long - more than 128 chars",
			password: string(make([]byte, 129)),
			wantBool: false,
		},
		{
			name:     "exactly 8 chars with all requirements",
			password: "T3st!@#a",
			wantBool: true,
		},
		{
			name:     "all lowercase - fails uppercase requirement",
			password: "test123!@#abc",
			wantBool: false,
		},
		{
			name:     "all uppercase - fails lowercase requirement",
			password: "TEST123!@#ABC",
			wantBool: false,
		},
		{
			name:     "only letters and special - fails digit requirement",
			password: "TestPass!@#",
			wantBool: false,
		},
		{
			name:     "only letters and digits - fails special char requirement",
			password: "Test12345",
			wantBool: false,
		},
		{
			name:     "valid with multiple special chars",
			password: "Secure123!@#$%",
			wantBool: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validPassword(tt.password)
			assert.Equal(t, tt.wantBool, result, "validPassword(%q) = %v, want %v", tt.password, result, tt.wantBool)
		})
	}
}

// ============================================================================
// Service Initialization Tests
// ============================================================================

func TestNewAdminService(t *testing.T) {
	games := &MockGameRepository{}
	users := &MockUserRepository{}
	players := &MockPlayerRepository{}
	orders := &MockOrderRepository{}
	payments := &MockPaymentRepository{}
	roles := &MockRoleRepository{}
	serviceItems := &MockServiceItemRepository{}
	permissions := &MockPermissionRepository{}
	menus := &MockMenuRepository{}
	stats := &MockStatsRepository{}
	wallets := &MockWalletRepository{}
	gameCategories := &MockGameCategoryRepository{}
	cache := NewMockCacheWithDefaults()

	svc := NewAdminService(
		games,
		users,
		players,
		orders,
		payments,
		roles,
		serviceItems,
		permissions,
		menus,
		stats,
		wallets,
		gameCategories,
		cache,
	)

	assert.NotNil(t, svc)
	assert.NotNil(t, svc.PermissionService())
	assert.NotNil(t, svc.RoleService())
	assert.NotNil(t, svc.MenuService())
	assert.NotNil(t, svc.StatsService())
}

// ============================================================================
// Helper Functions
// ============================================================================

func stringPtr(s string) *string {
	return &s
}

// ============================================================================
// StatsService Tests
// ============================================================================

func TestStatsService_Dashboard(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		setupMock func(*MockStatsRepository)
		wantErr   bool
	}{
		{
			name: "success - get dashboard",
			setupMock: func(s *MockStatsRepository) {
				s.On("Dashboard", ctx).Return(repository.Dashboard{
					TotalUsers:   100,
					TotalPlayers: 50,
					TotalOrders:  200,
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			setupMock: func(s *MockStatsRepository) {
				s.On("Dashboard", ctx).Return(repository.Dashboard{}, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := &MockStatsRepository{}
			tt.setupMock(stats)

			svc := NewStatsService(stats)
			result, err := svc.Dashboard(ctx)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, int64(100), result.TotalUsers)
			}

			stats.AssertExpectations(t)
		})
	}
}

func TestStatsService_RevenueTrend(t *testing.T) {
	ctx := context.Background()

	stats := &MockStatsRepository{}
	stats.On("RevenueTrend", ctx, 7).Return([]repository.DateValue{
		{Date: "2024-01-01", Value: 1000},
		{Date: "2024-01-02", Value: 1500},
	}, nil)

	svc := NewStatsService(stats)
	result, err := svc.RevenueTrend(ctx, 7)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	stats.AssertExpectations(t)
}

func TestStatsService_UserGrowth(t *testing.T) {
	ctx := context.Background()

	stats := &MockStatsRepository{}
	stats.On("UserGrowth", ctx, 30).Return([]repository.DateValue{
		{Date: "2024-01-01", Value: 10},
	}, nil)

	svc := NewStatsService(stats)
	result, err := svc.UserGrowth(ctx, 30)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	stats.AssertExpectations(t)
}

func TestStatsService_OrdersByStatus(t *testing.T) {
	ctx := context.Background()

	stats := &MockStatsRepository{}
	stats.On("OrdersByStatus", ctx).Return(map[string]int64{
		"pending":   10,
		"completed": 50,
	}, nil)

	svc := NewStatsService(stats)
	result, err := svc.OrdersByStatus(ctx)

	require.NoError(t, err)
	assert.Equal(t, int64(10), result["pending"])
	assert.Equal(t, int64(50), result["completed"])
	stats.AssertExpectations(t)
}

func TestStatsService_TopPlayers(t *testing.T) {
	ctx := context.Background()

	stats := &MockStatsRepository{}
	stats.On("TopPlayers", ctx, 10).Return([]repository.PlayerTop{
		{PlayerID: 1, Nickname: "Player1", RatingAverage: 4.8, RatingCount: 100},
		{PlayerID: 2, Nickname: "Player2", RatingAverage: 4.5, RatingCount: 80},
	}, nil)

	svc := NewStatsService(stats)
	result, err := svc.TopPlayers(ctx, 10)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	stats.AssertExpectations(t)
}

func TestStatsService_AuditOverview(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	stats := &MockStatsRepository{}
	stats.On("AuditOverview", ctx, &now, &now).Return(
		map[string]int64{"user": 10},
		map[string]int64{"create": 5},
		nil,
	)

	svc := NewStatsService(stats)
	byEntity, byAction, err := svc.AuditOverview(ctx, &now, &now)

	require.NoError(t, err)
	assert.Equal(t, int64(10), byEntity["user"])
	assert.Equal(t, int64(5), byAction["create"])
	stats.AssertExpectations(t)
}

func TestStatsService_AuditTrend(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	stats := &MockStatsRepository{}
	stats.On("AuditTrend", ctx, &now, &now, "user", "create").Return([]repository.DateValue{
		{Date: "2024-01-01", Value: 5},
	}, nil)

	svc := NewStatsService(stats)
	result, err := svc.AuditTrend(ctx, &now, &now, "user", "create")

	require.NoError(t, err)
	assert.Len(t, result, 1)
	stats.AssertExpectations(t)
}

func TestStatsService_UserBehaviorStats(t *testing.T) {
	ctx := context.Background()

	stats := &MockStatsRepository{}
	svc := NewStatsService(stats)

	result, err := svc.UserBehaviorStats(ctx)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1128), result.DAU)
}

// ============================================================================
// RoleService Tests
// ============================================================================

func TestRoleService_ListRoles(t *testing.T) {
	ctx := context.Background()

	roles := &MockRoleRepository{}
	roles.On("List", ctx).Return([]model.RoleModel{
		{Base: model.Base{ID: 1}, Slug: "admin", Name: "Admin"},
		{Base: model.Base{ID: 2}, Slug: "user", Name: "User"},
	}, nil)

	cache := NewMockCacheWithDefaults()
	svc := NewRoleService(roles, cache)
	result, err := svc.ListRoles(ctx)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	roles.AssertExpectations(t)
}

func TestRoleService_ListRolesPaged(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		page     int
		pageSize int
		expPage  int
		expSize  int
	}{
		{"normal pagination", 1, 10, 1, 10},
		{"default page", 0, 10, 1, 10},
		{"default page size", 1, 0, 1, 20},
		{"page size too large", 1, 200, 1, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roles := &MockRoleRepository{}
			roles.On("ListPaged", ctx, tt.expPage, tt.expSize).Return([]model.RoleModel{}, int64(0), nil)

			cache := NewMockCacheWithDefaults()
			svc := NewRoleService(roles, cache)
			_, _, err := svc.ListRolesPaged(ctx, tt.page, tt.pageSize)

			require.NoError(t, err)
			roles.AssertExpectations(t)
		})
	}
}

func TestRoleService_ListRolesPagedWithFilter(t *testing.T) {
	ctx := context.Background()
	isSystem := true

	roles := &MockRoleRepository{}
	roles.On("ListPagedWithFilter", ctx, 1, 20, "admin", &isSystem).Return([]model.RoleModel{
		{Base: model.Base{ID: 1}, Slug: "admin", Name: "Admin", IsSystem: true},
	}, int64(1), nil)

	cache := NewMockCacheWithDefaults()
	svc := NewRoleService(roles, cache)
	result, total, err := svc.ListRolesPagedWithFilter(ctx, 1, 20, "admin", &isSystem)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	roles.AssertExpectations(t)
}

func TestRoleService_ListRolesWithPermissions(t *testing.T) {
	ctx := context.Background()

	roles := &MockRoleRepository{}
	roles.On("ListWithPermissions", ctx).Return([]model.RoleModel{
		{Base: model.Base{ID: 1}, Slug: "admin", Name: "Admin"},
	}, nil)

	cache := NewMockCacheWithDefaults()
	svc := NewRoleService(roles, cache)
	result, err := svc.ListRolesWithPermissions(ctx)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	roles.AssertExpectations(t)
}

func TestRoleService_GetRole(t *testing.T) {
	ctx := context.Background()

	roles := &MockRoleRepository{}
	roles.On("Get", ctx, uint64(1)).Return(&model.RoleModel{
		Base: model.Base{ID: 1},
		Slug: "admin",
		Name: "Admin",
	}, nil)

	cache := NewMockCacheWithDefaults()
	svc := NewRoleService(roles, cache)
	result, err := svc.GetRole(ctx, 1)

	require.NoError(t, err)
	assert.Equal(t, "admin", result.Slug)
	roles.AssertExpectations(t)
}

func TestRoleService_GetRoleWithPermissions(t *testing.T) {
	ctx := context.Background()

	roles := &MockRoleRepository{}
	roles.On("GetWithPermissions", ctx, uint64(1)).Return(&model.RoleModel{
		Base: model.Base{ID: 1},
		Slug: "admin",
		Name: "Admin",
	}, nil)

	cache := NewMockCacheWithDefaults()
	svc := NewRoleService(roles, cache)
	result, err := svc.GetRoleWithPermissions(ctx, 1)

	require.NoError(t, err)
	assert.Equal(t, "admin", result.Slug)
	roles.AssertExpectations(t)
}

func TestRoleService_GetRoleBySlug(t *testing.T) {
	ctx := context.Background()

	roles := &MockRoleRepository{}
	roles.On("GetBySlug", ctx, "admin").Return(&model.RoleModel{
		Base: model.Base{ID: 1},
		Slug: "admin",
		Name: "Admin",
	}, nil)

	cache := NewMockCacheWithDefaults()
	svc := NewRoleService(roles, cache)
	result, err := svc.GetRoleBySlug(ctx, "admin")

	require.NoError(t, err)
	assert.Equal(t, uint64(1), result.ID)
	roles.AssertExpectations(t)
}

func TestRoleService_CreateRole(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		role        *model.RoleModel
		setupMock   func(*MockRoleRepository, *MockCache)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - create role",
			role: &model.RoleModel{Slug: "newrole", Name: "New Role"},
			setupMock: func(r *MockRoleRepository, c *MockCache) {
				r.On("GetBySlug", ctx, "newrole").Return(nil, repository.ErrNotFound)
				r.On("Create", ctx, mock.AnythingOfType("*model.RoleModel")).Return(nil)
				c.On("Delete", ctx, "admin:roles").Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "validation error - empty slug",
			role:        &model.RoleModel{Slug: "", Name: "New Role"},
			setupMock:   func(r *MockRoleRepository, c *MockCache) {},
			wantErr:     true,
			errContains: "slug and name are required",
		},
		{
			name:        "validation error - empty name",
			role:        &model.RoleModel{Slug: "newrole", Name: ""},
			setupMock:   func(r *MockRoleRepository, c *MockCache) {},
			wantErr:     true,
			errContains: "slug and name are required",
		},
		{
			name: "slug already exists",
			role: &model.RoleModel{Slug: "admin", Name: "Admin"},
			setupMock: func(r *MockRoleRepository, c *MockCache) {
				r.On("GetBySlug", ctx, "admin").Return(&model.RoleModel{
					Base: model.Base{ID: 1},
					Slug: "admin",
				}, nil)
			},
			wantErr:     true,
			errContains: "角色标识已存在",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roles := &MockRoleRepository{}
			cache := NewMockCacheWithDefaults()
			tt.setupMock(roles, cache)

			svc := NewRoleService(roles, cache)
			err := svc.CreateRole(ctx, tt.role)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}

			roles.AssertExpectations(t)
		})
	}
}

func TestRoleService_UpdateRole(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		role        *model.RoleModel
		setupMock   func(*MockRoleRepository, *MockCache)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - update role",
			role: &model.RoleModel{Base: model.Base{ID: 1}, Slug: "admin", Name: "Updated Admin"},
			setupMock: func(r *MockRoleRepository, c *MockCache) {
				r.On("Get", ctx, uint64(1)).Return(&model.RoleModel{
					Base:     model.Base{ID: 1},
					Slug:     "admin",
					Name:     "Admin",
					IsSystem: false,
				}, nil)
				r.On("Update", ctx, mock.AnythingOfType("*model.RoleModel")).Return(nil)
				c.On("Delete", ctx, "admin:roles").Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "validation error - zero ID",
			role:        &model.RoleModel{Base: model.Base{ID: 0}, Slug: "admin", Name: "Admin"},
			setupMock:   func(r *MockRoleRepository, c *MockCache) {},
			wantErr:     true,
			errContains: "role ID is required",
		},
		{
			name: "system role - only description updated",
			role: &model.RoleModel{Base: model.Base{ID: 1}, Slug: "newslug", Name: "New Name", Description: "New Desc"},
			setupMock: func(r *MockRoleRepository, c *MockCache) {
				r.On("Get", ctx, uint64(1)).Return(&model.RoleModel{
					Base:     model.Base{ID: 1},
					Slug:     "admin",
					Name:     "Admin",
					IsSystem: true,
				}, nil)
				r.On("Update", ctx, mock.MatchedBy(func(role *model.RoleModel) bool {
					return role.Slug == "admin" && role.Name == "Admin" && role.Description == "New Desc"
				})).Return(nil)
				c.On("Delete", ctx, "admin:roles").Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roles := &MockRoleRepository{}
			cache := NewMockCacheWithDefaults()
			tt.setupMock(roles, cache)

			svc := NewRoleService(roles, cache)
			err := svc.UpdateRole(ctx, tt.role)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}

			roles.AssertExpectations(t)
		})
	}
}

func TestRoleService_DeleteRole(t *testing.T) {
	ctx := context.Background()

	roles := &MockRoleRepository{}
	roles.On("Delete", ctx, uint64(1)).Return(nil)

	cache := NewMockCacheWithDefaults()
	cache.On("Delete", ctx, "admin:roles").Return(nil)

	svc := NewRoleService(roles, cache)
	err := svc.DeleteRole(ctx, 1)

	require.NoError(t, err)
	roles.AssertExpectations(t)
}

func TestRoleService_AssignPermissionsToRole(t *testing.T) {
	ctx := context.Background()

	roles := &MockRoleRepository{}
	roles.On("AssignPermissions", ctx, uint64(1), []uint64{1, 2, 3}).Return(nil)
	roles.On("GetUserIDsByRoleID", ctx, uint64(1)).Return([]uint64{10, 20}, nil)

	cache := NewMockCacheWithDefaults()
	cache.On("Delete", ctx, "admin:roles").Return(nil)
	cache.On("Delete", ctx, mock.AnythingOfType("string")).Return(nil)

	svc := NewRoleService(roles, cache)
	err := svc.AssignPermissionsToRole(ctx, 1, []uint64{1, 2, 3})

	require.NoError(t, err)
	roles.AssertExpectations(t)
}

func TestRoleService_AddPermissionsToRole(t *testing.T) {
	ctx := context.Background()

	roles := &MockRoleRepository{}
	roles.On("AddPermissions", ctx, uint64(1), []uint64{4, 5}).Return(nil)
	roles.On("GetUserIDsByRoleID", ctx, uint64(1)).Return([]uint64{10}, nil)

	cache := NewMockCacheWithDefaults()
	cache.On("Delete", ctx, "admin:roles").Return(nil)
	cache.On("Delete", ctx, mock.AnythingOfType("string")).Return(nil)

	svc := NewRoleService(roles, cache)
	err := svc.AddPermissionsToRole(ctx, 1, []uint64{4, 5})

	require.NoError(t, err)
	roles.AssertExpectations(t)
}

func TestRoleService_RemovePermissionsFromRole(t *testing.T) {
	ctx := context.Background()

	roles := &MockRoleRepository{}
	roles.On("RemovePermissions", ctx, uint64(1), []uint64{1, 2}).Return(nil)
	roles.On("GetUserIDsByRoleID", ctx, uint64(1)).Return([]uint64{}, nil)

	cache := NewMockCacheWithDefaults()
	cache.On("Delete", ctx, "admin:roles").Return(nil)

	svc := NewRoleService(roles, cache)
	err := svc.RemovePermissionsFromRole(ctx, 1, []uint64{1, 2})

	require.NoError(t, err)
	roles.AssertExpectations(t)
}

func TestRoleService_ListRolesByUserID(t *testing.T) {
	ctx := context.Background()

	roles := &MockRoleRepository{}
	roles.On("ListByUserID", ctx, uint64(1)).Return([]model.RoleModel{
		{Base: model.Base{ID: 1}, Slug: "admin", Name: "Admin"},
	}, nil)

	cache := NewMockCacheWithDefaults()
	cache.On("Get", ctx, "admin:roles:user:1").Return("", false, nil)
	cache.On("Set", ctx, "admin:roles:user:1", mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil)

	svc := NewRoleService(roles, cache)
	result, err := svc.ListRolesByUserID(ctx, 1)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	roles.AssertExpectations(t)
}

func TestRoleService_AssignRolesToUser(t *testing.T) {
	ctx := context.Background()

	roles := &MockRoleRepository{}
	roles.On("AssignToUser", ctx, uint64(1), []uint64{1, 2}).Return(nil)

	cache := NewMockCacheWithDefaults()
	cache.On("Delete", ctx, "admin:roles:user:1").Return(nil)
	cache.On("Delete", ctx, mock.AnythingOfType("string")).Return(nil)

	svc := NewRoleService(roles, cache)
	err := svc.AssignRolesToUser(ctx, 1, []uint64{1, 2})

	require.NoError(t, err)
	roles.AssertExpectations(t)
}

func TestRoleService_RemoveRolesFromUser(t *testing.T) {
	ctx := context.Background()

	roles := &MockRoleRepository{}
	roles.On("RemoveFromUser", ctx, uint64(1), []uint64{1}).Return(nil)

	cache := NewMockCacheWithDefaults()
	cache.On("Delete", ctx, "admin:roles:user:1").Return(nil)
	cache.On("Delete", ctx, mock.AnythingOfType("string")).Return(nil)

	svc := NewRoleService(roles, cache)
	err := svc.RemoveRolesFromUser(ctx, 1, []uint64{1})

	require.NoError(t, err)
	roles.AssertExpectations(t)
}

func TestRoleService_CheckUserHasRole(t *testing.T) {
	ctx := context.Background()

	roles := &MockRoleRepository{}
	roles.On("CheckUserHasRole", ctx, uint64(1), "admin").Return(true, nil)

	cache := NewMockCacheWithDefaults()
	svc := NewRoleService(roles, cache)
	result, err := svc.CheckUserHasRole(ctx, 1, "admin")

	require.NoError(t, err)
	assert.True(t, result)
	roles.AssertExpectations(t)
}

func TestRoleService_CheckUserIsSuperAdmin(t *testing.T) {
	ctx := context.Background()

	roles := &MockRoleRepository{}
	roles.On("CheckUserHasRole", ctx, uint64(1), "superAdmin").Return(true, nil)

	cache := NewMockCacheWithDefaults()
	svc := NewRoleService(roles, cache)
	result, err := svc.CheckUserIsSuperAdmin(ctx, 1)

	require.NoError(t, err)
	assert.True(t, result)
	roles.AssertExpectations(t)
}

func TestRoleService_GetRolePermissionIDs(t *testing.T) {
	ctx := context.Background()

	roles := &MockRoleRepository{}
	roles.On("GetWithPermissions", ctx, uint64(1)).Return(&model.RoleModel{
		Base: model.Base{ID: 1},
		Slug: "admin",
		Permissions: []model.Permission{
			{Base: model.Base{ID: 1}},
			{Base: model.Base{ID: 2}},
		},
	}, nil)

	cache := NewMockCacheWithDefaults()
	svc := NewRoleService(roles, cache)
	result, err := svc.GetRolePermissionIDs(ctx, 1)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Contains(t, result, uint64(1))
	assert.Contains(t, result, uint64(2))
	roles.AssertExpectations(t)
}

func TestRoleService_GetUserIDsByRoleID(t *testing.T) {
	ctx := context.Background()

	roles := &MockRoleRepository{}
	roles.On("GetUserIDsByRoleID", ctx, uint64(1)).Return([]uint64{10, 20, 30}, nil)

	cache := NewMockCacheWithDefaults()
	svc := NewRoleService(roles, cache)
	result, err := svc.GetUserIDsByRoleID(ctx, 1)

	require.NoError(t, err)
	assert.Len(t, result, 3)
	roles.AssertExpectations(t)
}

// ============================================================================
// Additional AdminService Tests for Coverage Improvement
// ============================================================================

func TestAdminService_GetGame(t *testing.T) {
	tests := []struct {
		name        string
		gameID      uint64
		setupMock   func(*MockGameRepository)
		expectError bool
	}{
		{
			name:   "success - get game",
			gameID: 1,
			setupMock: func(g *MockGameRepository) {
				g.On("Get", mock.Anything, uint64(1)).Return(&model.Game{
					Base: model.Base{ID: 1},
					Key:  "lol",
					Name: "League of Legends",
				}, nil)
			},
			expectError: false,
		},
		{
			name:   "game not found",
			gameID: 999,
			setupMock: func(g *MockGameRepository) {
				g.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			games := &MockGameRepository{}
			tt.setupMock(games)

			svc := &AdminService{games: games}
			result, err := svc.GetGame(context.Background(), tt.gameID)

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
			}
			games.AssertExpectations(t)
		})
	}
}

func TestAdminService_GetPlayer(t *testing.T) {
	tests := []struct {
		name        string
		playerID    uint64
		setupMock   func(*MockPlayerRepository)
		expectError bool
	}{
		{
			name:     "success - get player",
			playerID: 1,
			setupMock: func(p *MockPlayerRepository) {
				p.On("Get", mock.Anything, uint64(1)).Return(&model.Player{
					Base:   model.Base{ID: 1},
					UserID: 100,
				}, nil)
			},
			expectError: false,
		},
		{
			name:     "player not found",
			playerID: 999,
			setupMock: func(p *MockPlayerRepository) {
				p.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			tt.setupMock(players)

			svc := &AdminService{players: players}
			result, err := svc.GetPlayer(context.Background(), tt.playerID)

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
			}
			players.AssertExpectations(t)
		})
	}
}

func TestAdminService_ListGamesPaged(t *testing.T) {
	tests := []struct {
		name        string
		page        int
		pageSize    int
		setupMock   func(*MockGameRepository)
		expectError bool
		expectLen   int
		expectTotal int
	}{
		{
			name:     "success - list games paged",
			page:     1,
			pageSize: 10,
			setupMock: func(g *MockGameRepository) {
				g.On("ListPaged", mock.Anything, 1, 10).Return([]model.Game{
					{Base: model.Base{ID: 1}, Key: "lol", Name: "League of Legends"},
					{Base: model.Base{ID: 2}, Key: "dota2", Name: "Dota 2"},
				}, int64(2), nil)
			},
			expectError: false,
			expectLen:   2,
			expectTotal: 2,
		},
		{
			name:     "default page and page size",
			page:     0,
			pageSize: 0,
			setupMock: func(g *MockGameRepository) {
				g.On("ListPaged", mock.Anything, 1, 20).Return([]model.Game{}, int64(0), nil)
			},
			expectError: false,
			expectLen:   0,
			expectTotal: 0,
		},
		{
			name:     "page size too large",
			page:     1,
			pageSize: 200,
			setupMock: func(g *MockGameRepository) {
				g.On("ListPaged", mock.Anything, 1, 100).Return([]model.Game{}, int64(0), nil)
			},
			expectError: false,
			expectLen:   0,
			expectTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			games := &MockGameRepository{}
			tt.setupMock(games)

			svc := &AdminService{games: games}
			result, pagination, err := svc.ListGamesPaged(context.Background(), tt.page, tt.pageSize)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.expectLen)
				assert.NotNil(t, pagination)
				assert.Equal(t, tt.expectTotal, pagination.Total)
			}
			games.AssertExpectations(t)
		})
	}
}

func TestAdminService_ListGamesPagedWithFilter(t *testing.T) {
	tests := []struct {
		name        string
		page        int
		pageSize    int
		keyword     string
		setupMock   func(*MockGameRepository)
		expectError bool
		expectLen   int
	}{
		{
			name:     "success - list with filter",
			page:     1,
			pageSize: 10,
			keyword:  "league",
			setupMock: func(g *MockGameRepository) {
				g.On("ListPagedWithFilter", mock.Anything, 1, 10, "league").Return([]model.Game{
					{Base: model.Base{ID: 1}, Key: "lol", Name: "League of Legends"},
				}, int64(1), nil)
			},
			expectError: false,
			expectLen:   1,
		},
		{
			name:     "empty keyword",
			page:     1,
			pageSize: 10,
			keyword:  "",
			setupMock: func(g *MockGameRepository) {
				g.On("ListPagedWithFilter", mock.Anything, 1, 10, "").Return([]model.Game{}, int64(0), nil)
			},
			expectError: false,
			expectLen:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			games := &MockGameRepository{}
			tt.setupMock(games)

			svc := &AdminService{games: games}
			result, _, err := svc.ListGamesPagedWithFilter(context.Background(), tt.page, tt.pageSize, tt.keyword)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.expectLen)
			}
			games.AssertExpectations(t)
		})
	}
}

func TestAdminService_ListPlayersPaged(t *testing.T) {
	tests := []struct {
		name        string
		page        int
		pageSize    int
		setupMock   func(*MockPlayerRepository)
		expectError bool
		expectLen   int
	}{
		{
			name:     "success - list players paged",
			page:     1,
			pageSize: 10,
			setupMock: func(p *MockPlayerRepository) {
				p.On("ListPaged", mock.Anything, 1, 10).Return([]model.Player{
					{Base: model.Base{ID: 1}, UserID: 100},
					{Base: model.Base{ID: 2}, UserID: 101},
				}, int64(2), nil)
			},
			expectError: false,
			expectLen:   2,
		},
		{
			name:     "default pagination",
			page:     0,
			pageSize: 0,
			setupMock: func(p *MockPlayerRepository) {
				p.On("ListPaged", mock.Anything, 1, 20).Return([]model.Player{}, int64(0), nil)
			},
			expectError: false,
			expectLen:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			tt.setupMock(players)

			svc := &AdminService{players: players}
			result, _, err := svc.ListPlayersPaged(context.Background(), tt.page, tt.pageSize)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.expectLen)
			}
			players.AssertExpectations(t)
		})
	}
}

func TestAdminService_ListPlayersPagedWithFilter(t *testing.T) {
	verified := model.VerificationVerified
	tests := []struct {
		name        string
		page        int
		pageSize    int
		keyword     string
		status      *model.VerificationStatus
		setupMock   func(*MockPlayerRepository)
		expectError bool
		expectLen   int
	}{
		{
			name:     "success - list with keyword filter",
			page:     1,
			pageSize: 10,
			keyword:  "test",
			status:   nil,
			setupMock: func(p *MockPlayerRepository) {
				p.On("ListPagedWithFilter", mock.Anything, 1, 10, "test", (*model.VerificationStatus)(nil)).Return([]model.Player{
					{Base: model.Base{ID: 1}, UserID: 100},
				}, int64(1), nil)
			},
			expectError: false,
			expectLen:   1,
		},
		{
			name:     "success - list with status filter",
			page:     1,
			pageSize: 10,
			keyword:  "",
			status:   &verified,
			setupMock: func(p *MockPlayerRepository) {
				p.On("ListPagedWithFilter", mock.Anything, 1, 10, "", &verified).Return([]model.Player{
					{Base: model.Base{ID: 1}, UserID: 100, VerificationStatus: model.VerificationVerified},
				}, int64(1), nil)
			},
			expectError: false,
			expectLen:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			tt.setupMock(players)

			svc := &AdminService{players: players}
			result, _, err := svc.ListPlayersPagedWithFilter(context.Background(), tt.page, tt.pageSize, tt.keyword, tt.status)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.expectLen)
			}
			players.AssertExpectations(t)
		})
	}
}

func TestAdminService_ListUsersPaged(t *testing.T) {
	tests := []struct {
		name        string
		page        int
		pageSize    int
		setupMock   func(*MockUserRepository)
		expectError bool
		expectLen   int
	}{
		{
			name:     "success - list users paged",
			page:     1,
			pageSize: 10,
			setupMock: func(u *MockUserRepository) {
				u.On("ListPaged", mock.Anything, 1, 10).Return([]model.User{
					{Base: model.Base{ID: 1}, Name: "User1"},
					{Base: model.Base{ID: 2}, Name: "User2"},
				}, int64(2), nil)
			},
			expectError: false,
			expectLen:   2,
		},
		{
			name:     "default pagination",
			page:     0,
			pageSize: 0,
			setupMock: func(u *MockUserRepository) {
				u.On("ListPaged", mock.Anything, 1, 20).Return([]model.User{}, int64(0), nil)
			},
			expectError: false,
			expectLen:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &MockUserRepository{}
			tt.setupMock(users)

			svc := &AdminService{users: users}
			result, _, err := svc.ListUsersPaged(context.Background(), tt.page, tt.pageSize)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.expectLen)
			}
			users.AssertExpectations(t)
		})
	}
}

func TestAdminService_GetUsersByIDs(t *testing.T) {
	tests := []struct {
		name        string
		ids         []uint64
		setupMock   func(*MockUserRepository)
		expectError bool
		expectLen   int
	}{
		{
			name: "success - get users by IDs",
			ids:  []uint64{1, 2, 3},
			setupMock: func(u *MockUserRepository) {
				u.On("GetByIDs", mock.Anything, []uint64{1, 2, 3}).Return([]model.User{
					{Base: model.Base{ID: 1}, Name: "User1"},
					{Base: model.Base{ID: 2}, Name: "User2"},
					{Base: model.Base{ID: 3}, Name: "User3"},
				}, nil)
			},
			expectError: false,
			expectLen:   3,
		},
		{
			name: "empty IDs - returns early without calling repo",
			ids:  []uint64{},
			setupMock: func(u *MockUserRepository) {
				// No mock setup - implementation returns early for empty IDs
			},
			expectError: false,
			expectLen:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &MockUserRepository{}
			tt.setupMock(users)

			svc := &AdminService{users: users}
			result, err := svc.GetUsersByIDs(context.Background(), tt.ids)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.expectLen)
			}
			users.AssertExpectations(t)
		})
	}
}

// ============================================================================
// Order-related Tests
// ============================================================================

func TestAdminService_ListOrders(t *testing.T) {
	tests := []struct {
		name        string
		opts        repoiface.OrderListOptions
		setupMock   func(*MockOrderRepository)
		expectError bool
		expectLen   int
	}{
		{
			name: "success - list orders",
			opts: repoiface.OrderListOptions{
				Page:     1,
				PageSize: 10,
			},
			setupMock: func(o *MockOrderRepository) {
				o.On("List", mock.Anything, mock.MatchedBy(func(opts repoiface.OrderListOptions) bool {
					return opts.Page == 1 && opts.PageSize == 10
				})).Return([]model.Order{
					{Base: model.Base{ID: 1}, UserID: 100},
					{Base: model.Base{ID: 2}, UserID: 101},
				}, int64(2), nil)
			},
			expectError: false,
			expectLen:   2,
		},
		{
			name: "default pagination",
			opts: repoiface.OrderListOptions{
				Page:     0,
				PageSize: 0,
			},
			setupMock: func(o *MockOrderRepository) {
				o.On("List", mock.Anything, mock.MatchedBy(func(opts repoiface.OrderListOptions) bool {
					return opts.Page == 1 && opts.PageSize == 20
				})).Return([]model.Order{}, int64(0), nil)
			},
			expectError: false,
			expectLen:   0,
		},
		{
			name: "repository error",
			opts: repoiface.OrderListOptions{
				Page:     1,
				PageSize: 10,
			},
			setupMock: func(o *MockOrderRepository) {
				o.On("List", mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders := &MockOrderRepository{}
			tt.setupMock(orders)

			svc := &AdminService{orders: orders}
			result, pagination, err := svc.ListOrders(context.Background(), tt.opts)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.expectLen)
				assert.NotNil(t, pagination)
			}
			orders.AssertExpectations(t)
		})
	}
}

func TestAdminService_GetOrder(t *testing.T) {
	tests := []struct {
		name        string
		orderID     uint64
		setupMock   func(*MockOrderRepository)
		expectError bool
	}{
		{
			name:    "success - get order",
			orderID: 1,
			setupMock: func(o *MockOrderRepository) {
				o.On("Get", mock.Anything, uint64(1)).Return(&model.Order{
					Base:   model.Base{ID: 1},
					UserID: 100,
				}, nil)
			},
			expectError: false,
		},
		{
			name:    "order not found",
			orderID: 999,
			setupMock: func(o *MockOrderRepository) {
				o.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders := &MockOrderRepository{}
			tt.setupMock(orders)

			svc := &AdminService{orders: orders}
			result, err := svc.GetOrder(context.Background(), tt.orderID)

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
			}
			orders.AssertExpectations(t)
		})
	}
}

func TestAdminService_DeleteOrder(t *testing.T) {
	tests := []struct {
		name        string
		orderID     uint64
		setupMock   func(*MockOrderRepository)
		expectError bool
	}{
		{
			name:    "success - delete order",
			orderID: 1,
			setupMock: func(o *MockOrderRepository) {
				o.On("Delete", mock.Anything, uint64(1)).Return(nil)
			},
			expectError: false,
		},
		{
			name:    "order not found",
			orderID: 999,
			setupMock: func(o *MockOrderRepository) {
				o.On("Delete", mock.Anything, uint64(999)).Return(repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders := &MockOrderRepository{}
			tt.setupMock(orders)

			svc := &AdminService{orders: orders}
			err := svc.DeleteOrder(context.Background(), tt.orderID)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			orders.AssertExpectations(t)
		})
	}
}

func TestAdminService_ListUsersWithOptions(t *testing.T) {
	tests := []struct {
		name        string
		opts        repository.UserListOptions
		setupMock   func(*MockUserRepository)
		expectError bool
		expectLen   int
	}{
		{
			name: "success - list with options",
			opts: repository.UserListOptions{
				Page:     1,
				PageSize: 10,
			},
			setupMock: func(u *MockUserRepository) {
				u.On("ListWithFilters", mock.Anything, mock.MatchedBy(func(opts repository.UserListOptions) bool {
					return opts.Page == 1 && opts.PageSize == 10
				})).Return([]model.User{
					{Base: model.Base{ID: 1}, Name: "User1"},
				}, int64(1), nil)
			},
			expectError: false,
			expectLen:   1,
		},
		{
			name: "default pagination",
			opts: repository.UserListOptions{
				Page:     0,
				PageSize: 0,
			},
			setupMock: func(u *MockUserRepository) {
				u.On("ListWithFilters", mock.Anything, mock.MatchedBy(func(opts repository.UserListOptions) bool {
					return opts.Page == 1 && opts.PageSize == 20
				})).Return([]model.User{}, int64(0), nil)
			},
			expectError: false,
			expectLen:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &MockUserRepository{}
			tt.setupMock(users)

			svc := &AdminService{users: users}
			result, _, err := svc.ListUsersWithOptions(context.Background(), tt.opts)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.expectLen)
			}
			users.AssertExpectations(t)
		})
	}
}

func TestAdminService_SetTxManager(t *testing.T) {
	svc := &AdminService{}
	assert.Nil(t, svc.tx)

	// Create a mock TxManager (we just need to verify it's set)
	// Since TxManager is an interface, we can't easily mock it without more setup
	// But we can test that the method doesn't panic
	svc.SetTxManager(nil)
	assert.Nil(t, svc.tx)
}

func TestAdminService_ConfirmOrder(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		orderID     uint64
		note        string
		setupMock   func(*MockOrderRepository)
		expectError bool
	}{
		{
			name:    "success - confirm order",
			orderID: 1,
			note:    "Confirmed by admin",
			setupMock: func(o *MockOrderRepository) {
				order := &model.Order{
					Base:            model.Base{ID: 1},
					UserID:          100,
					Status:          model.OrderStatusPending,
					TotalPriceCents: 10000,
					Currency:        "CNY",
					ScheduledStart:  &now,
					ScheduledEnd:    &now,
				}
				o.On("Get", mock.Anything, uint64(1)).Return(order, nil).Once()
				o.On("Get", mock.Anything, uint64(1)).Return(order, nil).Once()
				o.On("Update", mock.Anything, mock.MatchedBy(func(o *model.Order) bool {
					return o.Status == model.OrderStatusConfirmed
				})).Return(nil)
			},
			expectError: false,
		},
		{
			name:    "order not found",
			orderID: 999,
			note:    "",
			setupMock: func(o *MockOrderRepository) {
				o.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders := &MockOrderRepository{}
			tt.setupMock(orders)

			svc := &AdminService{orders: orders}
			result, err := svc.ConfirmOrder(context.Background(), tt.orderID, tt.note)

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
			}
			orders.AssertExpectations(t)
		})
	}
}

func TestAdminService_StartOrder(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		orderID     uint64
		note        string
		setupMock   func(*MockOrderRepository)
		expectError bool
	}{
		{
			name:    "success - start order",
			orderID: 1,
			note:    "Started by admin",
			setupMock: func(o *MockOrderRepository) {
				order := &model.Order{
					Base:            model.Base{ID: 1},
					UserID:          100,
					Status:          model.OrderStatusConfirmed,
					TotalPriceCents: 10000,
					Currency:        "CNY",
					ScheduledStart:  &now,
					ScheduledEnd:    &now,
				}
				o.On("Get", mock.Anything, uint64(1)).Return(order, nil).Once()
				o.On("Get", mock.Anything, uint64(1)).Return(order, nil).Once()
				o.On("Update", mock.Anything, mock.MatchedBy(func(o *model.Order) bool {
					return o.Status == model.OrderStatusInProgress && o.StartedAt != nil
				})).Return(nil)
			},
			expectError: false,
		},
		{
			name:    "order not found",
			orderID: 999,
			note:    "",
			setupMock: func(o *MockOrderRepository) {
				o.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders := &MockOrderRepository{}
			tt.setupMock(orders)

			svc := &AdminService{orders: orders}
			result, err := svc.StartOrder(context.Background(), tt.orderID, tt.note)

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
			}
			orders.AssertExpectations(t)
		})
	}
}

func TestAdminService_CompleteOrder(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		orderID     uint64
		note        string
		setupMock   func(*MockOrderRepository)
		expectError bool
	}{
		{
			name:    "success - complete order",
			orderID: 1,
			note:    "Completed by admin",
			setupMock: func(o *MockOrderRepository) {
				order := &model.Order{
					Base:            model.Base{ID: 1},
					UserID:          100,
					Status:          model.OrderStatusInProgress,
					TotalPriceCents: 10000,
					Currency:        "CNY",
					ScheduledStart:  &now,
					ScheduledEnd:    &now,
					StartedAt:       &now,
				}
				o.On("Get", mock.Anything, uint64(1)).Return(order, nil).Once()
				o.On("Get", mock.Anything, uint64(1)).Return(order, nil).Once()
				o.On("Update", mock.Anything, mock.MatchedBy(func(o *model.Order) bool {
					return o.Status == model.OrderStatusCompleted && o.CompletedAt != nil
				})).Return(nil)
			},
			expectError: false,
		},
		{
			name:    "order not found",
			orderID: 999,
			note:    "",
			setupMock: func(o *MockOrderRepository) {
				o.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders := &MockOrderRepository{}
			tt.setupMock(orders)

			svc := &AdminService{orders: orders}
			result, err := svc.CompleteOrder(context.Background(), tt.orderID, tt.note)

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
			}
			orders.AssertExpectations(t)
		})
	}
}

func TestAdminService_UpdatePlayerVerification(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		playerID    uint64
		input       UpdateVerificationInput
		setupMock   func(*MockPlayerRepository, *MockCache)
		expectError bool
	}{
		{
			name:     "success - verify player",
			playerID: 1,
			input: UpdateVerificationInput{
				Nickname:           "TestPlayer",
				Bio:                "Test bio",
				HourlyRateCents:    5000,
				MainGameID:         1,
				VerificationStatus: model.VerificationVerified,
				VerifiedBy:         100,
				Remark:             "Approved",
			},
			setupMock: func(p *MockPlayerRepository, c *MockCache) {
				player := &model.Player{
					Base:               model.Base{ID: 1},
					UserID:             10,
					Nickname:           "OldNickname",
					VerificationStatus: model.VerificationPending,
				}
				p.On("Get", ctx, uint64(1)).Return(player, nil)
				p.On("Update", ctx, mock.AnythingOfType("*model.Player")).Return(nil)
				c.On("Delete", ctx, "admin:players").Return(nil)
			},
			expectError: false,
		},
		{
			name:     "success - reject player",
			playerID: 1,
			input: UpdateVerificationInput{
				Nickname:           "TestPlayer",
				Bio:                "Test bio",
				HourlyRateCents:    5000,
				MainGameID:         1,
				VerificationStatus: model.VerificationRejected,
				VerifiedBy:         100,
				Remark:             "Rejected due to invalid proof",
			},
			setupMock: func(p *MockPlayerRepository, c *MockCache) {
				player := &model.Player{
					Base:               model.Base{ID: 1},
					UserID:             10,
					Nickname:           "OldNickname",
					VerificationStatus: model.VerificationPending,
				}
				p.On("Get", ctx, uint64(1)).Return(player, nil)
				p.On("Update", ctx, mock.MatchedBy(func(pl *model.Player) bool {
					return pl.VerificationStatus == model.VerificationRejected && pl.RejectReason != ""
				})).Return(nil)
				c.On("Delete", ctx, "admin:players").Return(nil)
			},
			expectError: false,
		},
		{
			name:     "player not found",
			playerID: 999,
			input: UpdateVerificationInput{
				Nickname:           "TestPlayer",
				VerificationStatus: model.VerificationVerified,
				VerifiedBy:         100,
			},
			setupMock: func(p *MockPlayerRepository, c *MockCache) {
				p.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
		{
			name:     "validation error - empty status",
			playerID: 1,
			input: UpdateVerificationInput{
				Nickname:           "TestPlayer",
				VerificationStatus: "",
				VerifiedBy:         100,
			},
			setupMock: func(p *MockPlayerRepository, c *MockCache) {
				player := &model.Player{
					Base:   model.Base{ID: 1},
					UserID: 10,
				}
				p.On("Get", ctx, uint64(1)).Return(player, nil)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockPlayerRepository{}
			cache := NewMockCacheWithDefaults()
			tt.setupMock(players, cache)

			svc := NewAdminService(
				nil, nil, players, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				cache,
			)

			result, err := svc.UpdatePlayerVerification(ctx, tt.playerID, tt.input)

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
			}
			players.AssertExpectations(t)
		})
	}
}

func TestIsValidOrderStatus(t *testing.T) {
	tests := []struct {
		status model.OrderStatus
		valid  bool
	}{
		{model.OrderStatusPending, true},
		{model.OrderStatusConfirmed, true},
		{model.OrderStatusInProgress, true},
		{model.OrderStatusCompleted, true},
		{model.OrderStatusCanceled, true},
		{model.OrderStatusRefunded, true},
		{model.OrderStatus("invalid"), false},
		{model.OrderStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			result := isValidOrderStatus(tt.status)
			assert.Equal(t, tt.valid, result)
		})
	}
}

func TestIsAllowedOrderTransition(t *testing.T) {
	tests := []struct {
		name string
		from model.OrderStatus
		to   model.OrderStatus
		want bool
	}{
		{"pending to confirmed", model.OrderStatusPending, model.OrderStatusConfirmed, true},
		{"pending to canceled", model.OrderStatusPending, model.OrderStatusCanceled, true},
		{"confirmed to in_progress", model.OrderStatusConfirmed, model.OrderStatusInProgress, true},
		{"in_progress to completed", model.OrderStatusInProgress, model.OrderStatusCompleted, true},
		{"completed to refunded", model.OrderStatusCompleted, model.OrderStatusRefunded, true},
		{"pending to completed", model.OrderStatusPending, model.OrderStatusCompleted, false},
		{"canceled to confirmed", model.OrderStatusCanceled, model.OrderStatusConfirmed, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAllowedOrderTransition(tt.from, tt.to)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestMapUserError(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		expectNil   bool
		expectError bool
	}{
		{"nil error", nil, true, false},
		{"not found error", repository.ErrNotFound, false, true},
		{"duplicate email", errors.New("duplicate email"), false, true},
		{"duplicate phone", errors.New("duplicate phone"), false, true},
		{"other error", errors.New("some error"), false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapUserError(tt.err)
			if tt.expectNil {
				assert.Nil(t, result)
			} else if tt.expectError {
				assert.Error(t, result)
			}
		})
	}
}

func TestContainsRune(t *testing.T) {
	tests := []struct {
		s        string
		r        rune
		expected bool
	}{
		{"hello", 'e', true},
		{"xyz", 'a', false},
		{"Hello123", '1', true},
		{"abc", 'A', false},
		{"", 'a', false},
	}

	for _, tt := range tests {
		t.Run(tt.s+"_"+string(tt.r), func(t *testing.T) {
			result := containsRune(tt.s, tt.r)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOptionalPassword(t *testing.T) {
	tests := []struct {
		name     string
		password *string
		expected string
	}{
		{"nil password", nil, ""},
		{"empty password", strPtr(""), ""},
		{"valid password", strPtr("Password123!"), "Password123!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := optionalPassword(tt.password)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func strPtr(s string) *string {
	return &s
}

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"valid password", "Password123!", false},
		{"empty password", "", true},
		{"whitespace only", "   ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := hashPassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
				// Verify it's a valid bcrypt hash
				assert.True(t, len(result) > 50)
			}
		})
	}
}

func TestBuildPagination(t *testing.T) {
	tests := []struct {
		name      string
		page      int
		pageSize  int
		total     int64
		wantPage  int
		wantSize  int
		wantTotal int
	}{
		{"normal pagination", 1, 10, 100, 1, 10, 100},
		{"page 2", 2, 10, 100, 2, 10, 100},
		{"large total", 1, 20, 500, 1, 20, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildPagination(tt.page, tt.pageSize, tt.total)
			assert.Equal(t, tt.wantPage, result.Page)
			assert.Equal(t, tt.wantSize, result.PageSize)
			assert.Equal(t, tt.wantTotal, result.Total)
		})
	}
}

func TestPtrUint64(t *testing.T) {
	val := uint64(123)
	result := ptrUint64(val)
	assert.NotNil(t, result)
	assert.Equal(t, val, *result)
}

func TestMenuService_NewMenuService(t *testing.T) {
	menus := &MockMenuRepository{}

	svc := NewMenuService(menus)

	assert.NotNil(t, svc)
}

func TestPermissionService_NewPermissionService(t *testing.T) {
	perms := &MockPermissionRepository{}
	cache := NewMockCacheWithDefaults()

	svc := NewPermissionService(perms, cache)

	assert.NotNil(t, svc)
}

func TestRoleService_InvalidateRolePermissionsAndPropagateToUsers(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		roleID      uint64
		setupMock   func(*MockRoleRepository, *MockCache)
		expectError bool
	}{
		{
			name:   "success - invalidate and propagate",
			roleID: 1,
			setupMock: func(r *MockRoleRepository, c *MockCache) {
				r.On("GetUserIDsByRoleID", ctx, uint64(1)).Return([]uint64{10, 20, 30}, nil)
				c.On("Delete", ctx, "role:1:permissions").Return(nil)
				c.On("Delete", ctx, "user:10:permissions").Return(nil)
				c.On("Delete", ctx, "user:20:permissions").Return(nil)
				c.On("Delete", ctx, "user:30:permissions").Return(nil)
			},
			expectError: false,
		},
		{
			name:   "error getting user IDs",
			roleID: 1,
			setupMock: func(r *MockRoleRepository, c *MockCache) {
				r.On("GetUserIDsByRoleID", ctx, uint64(1)).Return([]uint64{}, errors.New("db error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roles := &MockRoleRepository{}
			cache := NewMockCacheWithDefaults()
			tt.setupMock(roles, cache)

			svc := NewRoleService(roles, cache)
			err := svc.InvalidateRolePermissionsAndPropagateToUsers(ctx, tt.roleID)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestWrapRepositoryError(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		err       error
		expectNil bool
	}{
		{"nil error", "test", nil, true},
		{"not found error", "get user", repository.ErrNotFound, false},
		{"validation error", "create user", ErrValidation, false},
		{"generic error", "update user", errors.New("db error"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wrapRepositoryError(tt.operation, tt.err)
			if tt.expectNil {
				assert.Nil(t, result)
			} else {
				assert.Error(t, result)
			}
		})
	}
}

func TestWrapErrorFunc(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		err       error
		expectNil bool
	}{
		{"nil error", "test", nil, true},
		{"not found error", "get user", repository.ErrNotFound, false},
		{"validation error", "create user", ErrValidation, false},
		{"order invalid transition", "update order", ErrOrderInvalidTransition, false},
		{"apierr error", "test", apierr.NotFound("test"), false},
		{"generic error", "update user", errors.New("db error"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wrapError(tt.operation, tt.err)
			if tt.expectNil {
				assert.Nil(t, result)
			} else {
				assert.Error(t, result)
			}
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"repository not found", repository.ErrNotFound, true},
		{"user not found", ErrUserNotFound, true},
		{"apierr not found", apierr.NotFound("test"), true},
		{"other error", errors.New("some error"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNotFound(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"validation error", ErrValidation, true},
		{"apierr validation", apierr.BadRequest("validation failed"), true},
		{"other error", errors.New("some error"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidationError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsBadRequest(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"validation error", ErrValidation, true},
		{"order invalid transition", ErrOrderInvalidTransition, true},
		{"apierr bad request", apierr.BadRequest("bad request"), true},
		{"other error", errors.New("some error"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isBadRequest(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWrapVoid(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		err       error
		expectNil bool
	}{
		{"nil error", "test", nil, true},
		{"not found error", "delete user", repository.ErrNotFound, false},
		{"generic error", "update user", errors.New("db error"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapVoid(tt.err, tt.operation)
			if tt.expectNil {
				assert.Nil(t, result)
			} else {
				assert.Error(t, result)
			}
		})
	}
}

func TestAdminService_WrapOrder(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		orderFunc   func() (*model.Order, error)
		operation   string
		expectError bool
		expectNil   bool
	}{
		{
			name: "success",
			orderFunc: func() (*model.Order, error) {
				return &model.Order{Base: model.Base{ID: 1}}, nil
			},
			operation:   "get order",
			expectError: false,
			expectNil:   false,
		},
		{
			name: "not found error",
			orderFunc: func() (*model.Order, error) {
				return nil, repository.ErrNotFound
			},
			operation:   "get order",
			expectError: true,
			expectNil:   true,
		},
		{
			name: "generic error",
			orderFunc: func() (*model.Order, error) {
				return nil, errors.New("db error")
			},
			operation:   "get order",
			expectError: true,
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &AdminService{}
			result, err := svc.WrapOrder(ctx, tt.orderFunc, tt.operation)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
			}
		})
	}
}

func TestSyncUserRoleToTable(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		userID      uint64
		role        model.Role
		setupMock   func(*MockRoleRepository)
		expectError bool
	}{
		{
			name:   "user role - syncs to user role",
			userID: 1,
			role:   model.RoleUser,
			setupMock: func(r *MockRoleRepository) {
				roleModel := &model.RoleModel{Name: "User", Slug: "user"}
				roleModel.ID = 1
				r.On("GetBySlug", mock.Anything, "user").Return(roleModel, nil)
				r.On("AssignToUser", mock.Anything, uint64(1), []uint64{uint64(1)}).Return(nil)
			},
			expectError: false,
		},
		{
			name:   "admin role - syncs to admin role",
			userID: 1,
			role:   model.RoleAdmin,
			setupMock: func(r *MockRoleRepository) {
				roleModel := &model.RoleModel{Name: "Admin", Slug: "admin"}
				roleModel.ID = 2
				r.On("GetBySlug", mock.Anything, "admin").Return(roleModel, nil)
				r.On("AssignToUser", mock.Anything, uint64(1), []uint64{uint64(2)}).Return(nil)
			},
			expectError: false,
		},
		{
			name:   "player role - syncs to player role",
			userID: 1,
			role:   model.RolePlayer,
			setupMock: func(r *MockRoleRepository) {
				roleModel := &model.RoleModel{Name: "Player", Slug: "player"}
				roleModel.ID = 3
				r.On("GetBySlug", mock.Anything, "player").Return(roleModel, nil)
				r.On("AssignToUser", mock.Anything, uint64(1), []uint64{uint64(3)}).Return(nil)
			},
			expectError: false,
		},
		{
			name:   "role not found in database - no error",
			userID: 1,
			role:   model.RoleUser,
			setupMock: func(r *MockRoleRepository) {
				r.On("GetBySlug", mock.Anything, "user").Return(nil, repository.ErrNotFound)
			},
			expectError: false,
		},
		{
			name:   "unknown role - no action",
			userID: 1,
			role:   model.Role("unknown"),
			setupMock: func(r *MockRoleRepository) {
				// No mock needed for unknown role
			},
			expectError: false,
		},
		{
			name:   "repository error",
			userID: 1,
			role:   model.RoleUser,
			setupMock: func(r *MockRoleRepository) {
				r.On("GetBySlug", mock.Anything, "user").Return(nil, errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roles := &MockRoleRepository{}
			tt.setupMock(roles)

			svc := NewAdminService(
				nil, nil, nil, nil, nil, roles, nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			err := svc.syncUserRoleToTable(ctx, tt.userID, tt.role)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ============================================================================
// Payment Validation Tests
// ============================================================================

func TestIsValidPaymentStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   model.PaymentStatus
		expected bool
	}{
		{
			name:     "pending is valid",
			status:   model.PaymentStatusPending,
			expected: true,
		},
		{
			name:     "paid is valid",
			status:   model.PaymentStatusPaid,
			expected: true,
		},
		{
			name:     "failed is valid",
			status:   model.PaymentStatusFailed,
			expected: true,
		},
		{
			name:     "refunded is valid",
			status:   model.PaymentStatusRefunded,
			expected: true,
		},
		{
			name:     "unknown status is invalid",
			status:   model.PaymentStatus("unknown"),
			expected: false,
		},
		{
			name:     "empty status is invalid",
			status:   model.PaymentStatus(""),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidPaymentStatus(tt.status)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsAllowedPaymentTransition(t *testing.T) {
	tests := []struct {
		name     string
		prev     model.PaymentStatus
		next     model.PaymentStatus
		expected bool
	}{
		// Same status transitions
		{
			name:     "pending to pending",
			prev:     model.PaymentStatusPending,
			next:     model.PaymentStatusPending,
			expected: true,
		},
		{
			name:     "paid to paid",
			prev:     model.PaymentStatusPaid,
			next:     model.PaymentStatusPaid,
			expected: true,
		},
		// Valid transitions from pending
		{
			name:     "pending to paid",
			prev:     model.PaymentStatusPending,
			next:     model.PaymentStatusPaid,
			expected: true,
		},
		{
			name:     "pending to failed",
			prev:     model.PaymentStatusPending,
			next:     model.PaymentStatusFailed,
			expected: true,
		},
		{
			name:     "pending to refunded",
			prev:     model.PaymentStatusPending,
			next:     model.PaymentStatusRefunded,
			expected: true,
		},
		// Valid transitions from paid
		{
			name:     "paid to refunded",
			prev:     model.PaymentStatusPaid,
			next:     model.PaymentStatusRefunded,
			expected: true,
		},
		// Invalid transitions from paid
		{
			name:     "paid to pending - invalid",
			prev:     model.PaymentStatusPaid,
			next:     model.PaymentStatusPending,
			expected: false,
		},
		{
			name:     "paid to failed - invalid",
			prev:     model.PaymentStatusPaid,
			next:     model.PaymentStatusFailed,
			expected: false,
		},
		// Invalid transitions from failed
		{
			name:     "failed to pending - invalid",
			prev:     model.PaymentStatusFailed,
			next:     model.PaymentStatusPending,
			expected: false,
		},
		{
			name:     "failed to paid - invalid",
			prev:     model.PaymentStatusFailed,
			next:     model.PaymentStatusPaid,
			expected: false,
		},
		// Invalid transitions from refunded
		{
			name:     "refunded to pending - invalid",
			prev:     model.PaymentStatusRefunded,
			next:     model.PaymentStatusPending,
			expected: false,
		},
		{
			name:     "refunded to paid - invalid",
			prev:     model.PaymentStatusRefunded,
			next:     model.PaymentStatusPaid,
			expected: false,
		},
		// Unknown status
		{
			name:     "unknown to pending - invalid",
			prev:     model.PaymentStatus("unknown"),
			next:     model.PaymentStatusPending,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAllowedPaymentTransition(tt.prev, tt.next)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ============================================================================
// AssignOrder Tests
// ============================================================================

func TestAdminService_AssignOrder(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		orderID     uint64
		playerID    uint64
		setupMock   func(*MockOrderRepository, *MockPlayerRepository)
		expectError bool
	}{
		{
			name:     "success - assign player to order",
			orderID:  1,
			playerID: 10,
			setupMock: func(orders *MockOrderRepository, players *MockPlayerRepository) {
				player := &model.Player{Nickname: "Test Player"}
				player.ID = 10
				players.On("Get", mock.Anything, uint64(10)).Return(player, nil)

				order := &model.Order{
					Status: model.OrderStatusPending,
				}
				order.ID = 1
				orders.On("Get", mock.Anything, uint64(1)).Return(order, nil)
				orders.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name:     "error - player ID is zero",
			orderID:  1,
			playerID: 0,
			setupMock: func(orders *MockOrderRepository, players *MockPlayerRepository) {
				// No mock needed
			},
			expectError: true,
		},
		{
			name:     "error - player not found",
			orderID:  1,
			playerID: 10,
			setupMock: func(orders *MockOrderRepository, players *MockPlayerRepository) {
				players.On("Get", mock.Anything, uint64(10)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
		{
			name:     "error - order not found",
			orderID:  1,
			playerID: 10,
			setupMock: func(orders *MockOrderRepository, players *MockPlayerRepository) {
				player := &model.Player{Nickname: "Test Player"}
				player.ID = 10
				players.On("Get", mock.Anything, uint64(10)).Return(player, nil)
				orders.On("Get", mock.Anything, uint64(1)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
		{
			name:     "error - order already completed",
			orderID:  1,
			playerID: 10,
			setupMock: func(orders *MockOrderRepository, players *MockPlayerRepository) {
				player := &model.Player{Nickname: "Test Player"}
				player.ID = 10
				players.On("Get", mock.Anything, uint64(10)).Return(player, nil)

				order := &model.Order{
					Status: model.OrderStatusCompleted,
				}
				order.ID = 1
				orders.On("Get", mock.Anything, uint64(1)).Return(order, nil)
			},
			expectError: true,
		},
		{
			name:     "error - order already canceled",
			orderID:  1,
			playerID: 10,
			setupMock: func(orders *MockOrderRepository, players *MockPlayerRepository) {
				player := &model.Player{Nickname: "Test Player"}
				player.ID = 10
				players.On("Get", mock.Anything, uint64(10)).Return(player, nil)

				order := &model.Order{
					Status: model.OrderStatusCanceled,
				}
				order.ID = 1
				orders.On("Get", mock.Anything, uint64(1)).Return(order, nil)
			},
			expectError: true,
		},
		{
			name:     "error - order already refunded",
			orderID:  1,
			playerID: 10,
			setupMock: func(orders *MockOrderRepository, players *MockPlayerRepository) {
				player := &model.Player{Nickname: "Test Player"}
				player.ID = 10
				players.On("Get", mock.Anything, uint64(10)).Return(player, nil)

				order := &model.Order{
					Status: model.OrderStatusRefunded,
				}
				order.ID = 1
				orders.On("Get", mock.Anything, uint64(1)).Return(order, nil)
			},
			expectError: true,
		},
		{
			name:     "error - update fails",
			orderID:  1,
			playerID: 10,
			setupMock: func(orders *MockOrderRepository, players *MockPlayerRepository) {
				player := &model.Player{Nickname: "Test Player"}
				player.ID = 10
				players.On("Get", mock.Anything, uint64(10)).Return(player, nil)

				order := &model.Order{
					Status: model.OrderStatusPending,
				}
				order.ID = 1
				orders.On("Get", mock.Anything, uint64(1)).Return(order, nil)
				orders.On("Update", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders := &MockOrderRepository{}
			players := &MockPlayerRepository{}
			tt.setupMock(orders, players)

			svc := NewAdminService(
				nil, nil, players, orders, nil, nil, nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			result, err := svc.AssignOrder(ctx, tt.orderID, tt.playerID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

// ============================================================================
// Payment Service Tests
// ============================================================================

func TestAdminService_ListPayments(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		opts        repository.PaymentListOptions
		setupMock   func(*MockPaymentRepository)
		expectError bool
		expectCount int
	}{
		{
			name: "success - list payments",
			opts: repository.PaymentListOptions{
				Page:     1,
				PageSize: 10,
			},
			setupMock: func(payments *MockPaymentRepository) {
				paymentList := []model.Payment{
					{Status: model.PaymentStatusPaid},
					{Status: model.PaymentStatusPending},
				}
				paymentList[0].ID = 1
				paymentList[1].ID = 2
				payments.On("List", mock.Anything, mock.Anything).Return(paymentList, int64(2), nil)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name: "success - empty list",
			opts: repository.PaymentListOptions{
				Page:     1,
				PageSize: 10,
			},
			setupMock: func(payments *MockPaymentRepository) {
				payments.On("List", mock.Anything, mock.Anything).Return([]model.Payment{}, int64(0), nil)
			},
			expectError: false,
			expectCount: 0,
		},
		{
			name: "error - repository error",
			opts: repository.PaymentListOptions{
				Page:     1,
				PageSize: 10,
			},
			setupMock: func(payments *MockPaymentRepository) {
				payments.On("List", mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payments := &MockPaymentRepository{}
			tt.setupMock(payments)

			svc := NewAdminService(
				nil, nil, nil, nil, payments, nil, nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			result, pagination, err := svc.ListPayments(ctx, tt.opts)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectCount)
				assert.NotNil(t, pagination)
			}
		})
	}
}

func TestAdminService_GetPayment(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		paymentID   uint64
		setupMock   func(*MockPaymentRepository)
		expectError bool
	}{
		{
			name:      "success - get payment",
			paymentID: 1,
			setupMock: func(payments *MockPaymentRepository) {
				payment := &model.Payment{
					Status:   model.PaymentStatusPaid,
					Currency: model.CurrencyCNY,
				}
				payment.ID = 1
				payments.On("Get", mock.Anything, uint64(1)).Return(payment, nil)
			},
			expectError: false,
		},
		{
			name:      "error - payment not found",
			paymentID: 999,
			setupMock: func(payments *MockPaymentRepository) {
				payments.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
		{
			name:      "error - repository error",
			paymentID: 1,
			setupMock: func(payments *MockPaymentRepository) {
				payments.On("Get", mock.Anything, uint64(1)).Return(nil, errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payments := &MockPaymentRepository{}
			tt.setupMock(payments)

			svc := NewAdminService(
				nil, nil, nil, nil, payments, nil, nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			result, err := svc.GetPayment(ctx, tt.paymentID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.paymentID, result.ID)
			}
		})
	}
}

func TestAdminService_GetPaymentWithRelations(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		paymentID   uint64
		setupMock   func(*MockPaymentRepository)
		expectError bool
	}{
		{
			name:      "success - get payment with relations",
			paymentID: 1,
			setupMock: func(payments *MockPaymentRepository) {
				payment := &model.Payment{
					Status:   model.PaymentStatusPaid,
					Currency: model.CurrencyCNY,
				}
				payment.ID = 1
				payments.On("GetWithRelations", mock.Anything, uint64(1)).Return(payment, nil)
			},
			expectError: false,
		},
		{
			name:      "error - payment not found",
			paymentID: 999,
			setupMock: func(payments *MockPaymentRepository) {
				payments.On("GetWithRelations", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payments := &MockPaymentRepository{}
			tt.setupMock(payments)

			svc := NewAdminService(
				nil, nil, nil, nil, payments, nil, nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			result, err := svc.GetPaymentWithRelations(ctx, tt.paymentID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestAdminService_GetPaymentsByOrderID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		orderID     uint64
		setupMock   func(*MockPaymentRepository)
		expectError bool
		expectCount int
	}{
		{
			name:    "success - get payments by order ID",
			orderID: 1,
			setupMock: func(payments *MockPaymentRepository) {
				paymentList := []model.Payment{
					{Status: model.PaymentStatusPaid},
					{Status: model.PaymentStatusRefunded},
				}
				paymentList[0].ID = 1
				paymentList[1].ID = 2
				payments.On("GetByOrderID", mock.Anything, uint64(1)).Return(paymentList, nil)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name:    "success - no payments for order",
			orderID: 2,
			setupMock: func(payments *MockPaymentRepository) {
				payments.On("GetByOrderID", mock.Anything, uint64(2)).Return([]model.Payment{}, nil)
			},
			expectError: false,
			expectCount: 0,
		},
		{
			name:    "error - repository error",
			orderID: 1,
			setupMock: func(payments *MockPaymentRepository) {
				payments.On("GetByOrderID", mock.Anything, uint64(1)).Return(nil, errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payments := &MockPaymentRepository{}
			tt.setupMock(payments)

			svc := NewAdminService(
				nil, nil, nil, nil, payments, nil, nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			result, err := svc.GetPaymentsByOrderID(ctx, tt.orderID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectCount)
			}
		})
	}
}

func TestAdminService_DeletePayment(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		paymentID   uint64
		setupMock   func(*MockPaymentRepository)
		expectError bool
	}{
		{
			name:      "success - delete payment",
			paymentID: 1,
			setupMock: func(payments *MockPaymentRepository) {
				payments.On("Delete", mock.Anything, uint64(1)).Return(nil)
			},
			expectError: false,
		},
		{
			name:      "error - payment not found",
			paymentID: 999,
			setupMock: func(payments *MockPaymentRepository) {
				payments.On("Delete", mock.Anything, uint64(999)).Return(repository.ErrNotFound)
			},
			expectError: true,
		},
		{
			name:      "error - repository error",
			paymentID: 1,
			setupMock: func(payments *MockPaymentRepository) {
				payments.On("Delete", mock.Anything, uint64(1)).Return(errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payments := &MockPaymentRepository{}
			tt.setupMock(payments)

			svc := NewAdminService(
				nil, nil, nil, nil, payments, nil, nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			err := svc.DeletePayment(ctx, tt.paymentID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ============================================================================
// Helper Function Tests
// ============================================================================

func TestMapRefundStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   model.PaymentStatus
		expected string
	}{
		{
			name:     "refunded status",
			status:   model.PaymentStatusRefunded,
			expected: "success",
		},
		{
			name:     "pending status",
			status:   model.PaymentStatusPending,
			expected: "pending",
		},
		{
			name:     "failed status",
			status:   model.PaymentStatusFailed,
			expected: "failed",
		},
		{
			name:     "paid status - default",
			status:   model.PaymentStatusPaid,
			expected: "paid",
		},
		{
			name:     "unknown status - default",
			status:   model.PaymentStatus("UNKNOWN"),
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapRefundStatus(tt.status)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapTimelineEventType(t *testing.T) {
	tests := []struct {
		name     string
		action   string
		expected string
	}{
		{
			name:     "create action",
			action:   string(model.OpActionCreate),
			expected: "system",
		},
		{
			name:     "assign player action",
			action:   string(model.OpActionAssignPlayer),
			expected: "action",
		},
		{
			name:     "confirm action",
			action:   string(model.OpActionConfirm),
			expected: "status_change",
		},
		{
			name:     "start action",
			action:   string(model.OpActionStart),
			expected: "status_change",
		},
		{
			name:     "complete action",
			action:   string(model.OpActionComplete),
			expected: "status_change",
		},
		{
			name:     "update status action",
			action:   string(model.OpActionUpdateStatus),
			expected: "status_change",
		},
		{
			name:     "cancel action",
			action:   string(model.OpActionCancel),
			expected: "status_change",
		},
		{
			name:     "refund action",
			action:   string(model.OpActionRefund),
			expected: "status_change",
		},
		{
			name:     "unknown action - default",
			action:   "unknown_action",
			expected: "action",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapTimelineEventType(tt.action)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapTimelineTitle(t *testing.T) {
	tests := []struct {
		name     string
		action   string
		expected string
	}{
		{
			name:     "create action",
			action:   string(model.OpActionCreate),
			expected: "订单创建",
		},
		{
			name:     "assign player action",
			action:   string(model.OpActionAssignPlayer),
			expected: "指派陪玩师",
		},
		{
			name:     "confirm action",
			action:   string(model.OpActionConfirm),
			expected: "订单确认",
		},
		{
			name:     "start action",
			action:   string(model.OpActionStart),
			expected: "开始服务",
		},
		{
			name:     "complete action",
			action:   string(model.OpActionComplete),
			expected: "完成订单",
		},
		{
			name:     "cancel action",
			action:   string(model.OpActionCancel),
			expected: "订单取消",
		},
		{
			name:     "refund action",
			action:   string(model.OpActionRefund),
			expected: "订单退款",
		},
		{
			name:     "update status action",
			action:   string(model.OpActionUpdateStatus),
			expected: "状态更新",
		},
		{
			name:     "unknown action - default",
			action:   "some_action",
			expected: "some action",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapTimelineTitle(tt.action)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ============================================================================
// Order Batch Tests - normalizeReason, normalizeNote, trim
// ============================================================================

func TestNormalizeReason(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal reason",
			input:    "Customer requested cancellation",
			expected: "Customer requested cancellation",
		},
		{
			name:     "reason with whitespace - not trimmed",
			input:    "  Customer requested cancellation  ",
			expected: "  Customer requested cancellation  ",
		},
		{
			name:     "empty reason",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeReason(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizeNote(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal note",
			input:    "Internal note for admin",
			expected: "Internal note for admin",
		},
		{
			name:     "note with whitespace - not trimmed",
			input:    "  Internal note  ",
			expected: "  Internal note  ",
		},
		{
			name:     "empty note",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeNote(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTrim(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal string",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "string with spaces - not trimmed",
			input:    "  hello world  ",
			expected: "  hello world  ",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trim(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ============================================================================
// Player Batch Tests
// ============================================================================

func TestIsValidVerificationStatusTransition(t *testing.T) {
	tests := []struct {
		name     string
		current  model.VerificationStatus
		new      model.VerificationStatus
		expected bool
	}{
		// From pending
		{
			name:     "pending to verified",
			current:  model.VerificationPending,
			new:      model.VerificationVerified,
			expected: true,
		},
		{
			name:     "pending to rejected",
			current:  model.VerificationPending,
			new:      model.VerificationRejected,
			expected: true,
		},
		{
			name:     "pending to pending - invalid",
			current:  model.VerificationPending,
			new:      model.VerificationPending,
			expected: false,
		},
		// From verified
		{
			name:     "verified to rejected",
			current:  model.VerificationVerified,
			new:      model.VerificationRejected,
			expected: true,
		},
		{
			name:     "verified to pending - revoke",
			current:  model.VerificationVerified,
			new:      model.VerificationPending,
			expected: true,
		},
		{
			name:     "verified to verified - invalid",
			current:  model.VerificationVerified,
			new:      model.VerificationVerified,
			expected: false,
		},
		// From rejected
		{
			name:     "rejected to pending",
			current:  model.VerificationRejected,
			new:      model.VerificationPending,
			expected: true,
		},
		{
			name:     "rejected to verified - re-verification",
			current:  model.VerificationRejected,
			new:      model.VerificationVerified,
			expected: true,
		},
		{
			name:     "rejected to rejected - invalid",
			current:  model.VerificationRejected,
			new:      model.VerificationRejected,
			expected: false,
		},
		// Unknown status
		{
			name:     "unknown status - invalid",
			current:  model.VerificationStatus("unknown"),
			new:      model.VerificationVerified,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidVerificationStatusTransition(tt.current, tt.new)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ============================================================================
// Stats Service Tests
// ============================================================================

func TestStatsService_SetUserBehaviorRepo(t *testing.T) {
	stats := &MockStatsRepository{}
	svc := NewStatsService(stats)

	// SetUserBehaviorRepo should not panic
	svc.SetUserBehaviorRepo(nil)
}

func TestStatsService_SetUserLoginHistoryRepo(t *testing.T) {
	stats := &MockStatsRepository{}
	svc := NewStatsService(stats)

	// SetUserLoginHistoryRepo should not panic
	svc.SetUserLoginHistoryRepo(nil)
}

// ============================================================================
// Menu Service Tests
// ============================================================================

func TestMenuService_Create(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		menu        *model.Menu
		setupMock   func(*MockMenuRepository)
		expectError bool
	}{
		{
			name: "success - create menu",
			menu: &model.Menu{
				Name: "Dashboard",
				Path: "/dashboard",
			},
			setupMock: func(menus *MockMenuRepository) {
				menus.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "error - repository error",
			menu: &model.Menu{
				Name: "Dashboard",
				Path: "/dashboard",
			},
			setupMock: func(menus *MockMenuRepository) {
				menus.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menus := &MockMenuRepository{}
			tt.setupMock(menus)

			svc := NewMenuService(menus)
			err := svc.Create(ctx, tt.menu)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMenuService_Update(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		menu        *model.Menu
		setupMock   func(*MockMenuRepository)
		expectError bool
	}{
		{
			name: "success - update menu",
			menu: &model.Menu{
				Name: "Dashboard Updated",
				Path: "/dashboard",
			},
			setupMock: func(menus *MockMenuRepository) {
				menus.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "error - menu ID is zero",
			menu: &model.Menu{
				Name: "Dashboard",
				Path: "/dashboard",
			},
			setupMock: func(menus *MockMenuRepository) {
				// No mock needed - should fail before calling repo
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menus := &MockMenuRepository{}
			tt.setupMock(menus)

			// Set ID for success case
			if tt.name == "success - update menu" {
				tt.menu.ID = 1
			}

			svc := NewMenuService(menus)
			err := svc.Update(ctx, tt.menu)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMenuService_Delete(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		menuID      uint64
		setupMock   func(*MockMenuRepository)
		expectError bool
	}{
		{
			name:   "success - delete menu",
			menuID: 1,
			setupMock: func(menus *MockMenuRepository) {
				menus.On("Delete", mock.Anything, uint64(1)).Return(nil)
			},
			expectError: false,
		},
		{
			name:   "error - menu not found",
			menuID: 999,
			setupMock: func(menus *MockMenuRepository) {
				menus.On("Delete", mock.Anything, uint64(999)).Return(repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menus := &MockMenuRepository{}
			tt.setupMock(menus)

			svc := NewMenuService(menus)
			err := svc.Delete(ctx, tt.menuID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMenuService_Get(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		menuID      uint64
		setupMock   func(*MockMenuRepository)
		expectError bool
	}{
		{
			name:   "success - get menu",
			menuID: 1,
			setupMock: func(menus *MockMenuRepository) {
				menu := &model.Menu{Name: "Dashboard", Path: "/dashboard"}
				menu.ID = 1
				menus.On("Get", mock.Anything, uint64(1)).Return(menu, nil)
			},
			expectError: false,
		},
		{
			name:   "error - menu not found",
			menuID: 999,
			setupMock: func(menus *MockMenuRepository) {
				menus.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menus := &MockMenuRepository{}
			tt.setupMock(menus)

			svc := NewMenuService(menus)
			result, err := svc.Get(ctx, tt.menuID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestMenuService_List(t *testing.T) {
	ctx := context.Background()
	parentID := uint64(1)

	tests := []struct {
		name        string
		parentID    *uint64
		setupMock   func(*MockMenuRepository)
		expectError bool
		expectCount int
	}{
		{
			name:     "success - list menus with parent",
			parentID: &parentID,
			setupMock: func(menus *MockMenuRepository) {
				menuList := []model.Menu{
					{Name: "Menu 1"},
					{Name: "Menu 2"},
				}
				menus.On("List", mock.Anything, &parentID).Return(menuList, nil)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name:     "success - list root menus",
			parentID: nil,
			setupMock: func(menus *MockMenuRepository) {
				menuList := []model.Menu{
					{Name: "Root Menu"},
				}
				menus.On("List", mock.Anything, (*uint64)(nil)).Return(menuList, nil)
			},
			expectError: false,
			expectCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menus := &MockMenuRepository{}
			tt.setupMock(menus)

			svc := NewMenuService(menus)
			result, err := svc.List(ctx, tt.parentID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectCount)
			}
		})
	}
}

func TestMenuService_ListPaged(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		page        int
		pageSize    int
		parentID    *uint64
		setupMock   func(*MockMenuRepository)
		expectError bool
		expectCount int
	}{
		{
			name:     "success - list paged menus",
			page:     1,
			pageSize: 10,
			parentID: nil,
			setupMock: func(menus *MockMenuRepository) {
				menuList := []model.Menu{
					{Name: "Menu 1"},
					{Name: "Menu 2"},
				}
				menus.On("ListPaged", mock.Anything, 1, 10, (*uint64)(nil)).Return(menuList, int64(2), nil)
			},
			expectError: false,
			expectCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menus := &MockMenuRepository{}
			tt.setupMock(menus)

			svc := NewMenuService(menus)
			result, total, err := svc.ListPaged(ctx, tt.page, tt.pageSize, tt.parentID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectCount)
				assert.Equal(t, int64(2), total)
			}
		})
	}
}

func TestMenuService_ListAccessible(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		codes       []string
		setupMock   func(*MockMenuRepository)
		expectError bool
		expectCount int
	}{
		{
			name:  "success - list accessible menus",
			codes: []string{"dashboard:view", "users:view"},
			setupMock: func(menus *MockMenuRepository) {
				menuList := []model.Menu{
					{Name: "Dashboard"},
					{Name: "Users"},
				}
				menus.On("ListByPermission", mock.Anything, []string{"dashboard:view", "users:view"}).Return(menuList, nil)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name:  "empty codes - return empty list",
			codes: []string{},
			setupMock: func(menus *MockMenuRepository) {
				// No mock needed - should return empty list without calling repo
			},
			expectError: false,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menus := &MockMenuRepository{}
			tt.setupMock(menus)

			svc := NewMenuService(menus)
			result, err := svc.ListAccessible(ctx, tt.codes)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectCount)
			}
		})
	}
}

// ============================================================================
// Permission Service Tests
// ============================================================================

func TestPermissionService_ListPermissions(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupMock   func(*MockPermissionRepository)
		expectError bool
		expectCount int
	}{
		{
			name: "success - list permissions",
			setupMock: func(perms *MockPermissionRepository) {
				permList := []model.Permission{
					{Code: "users:read", Method: "GET", Path: "/users"},
					{Code: "users:write", Method: "POST", Path: "/users"},
				}
				perms.On("List", mock.Anything).Return(permList, nil)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name: "error - repository error",
			setupMock: func(perms *MockPermissionRepository) {
				perms.On("List", mock.Anything).Return(nil, errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			perms := &MockPermissionRepository{}
			tt.setupMock(perms)

			svc := NewPermissionService(perms, NewMockCacheWithDefaults())
			result, err := svc.ListPermissions(ctx)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectCount)
			}
		})
	}
}

func TestPermissionService_ListPermissionsPaged(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		page        int
		pageSize    int
		setupMock   func(*MockPermissionRepository)
		expectError bool
		expectCount int
	}{
		{
			name:     "success - list paged permissions",
			page:     1,
			pageSize: 10,
			setupMock: func(perms *MockPermissionRepository) {
				permList := []model.Permission{
					{Code: "users:read"},
					{Code: "users:write"},
				}
				perms.On("ListPaged", mock.Anything, 1, 10).Return(permList, int64(2), nil)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name:     "normalize page - negative page",
			page:     -1,
			pageSize: 10,
			setupMock: func(perms *MockPermissionRepository) {
				permList := []model.Permission{{Code: "test"}}
				perms.On("ListPaged", mock.Anything, 1, 10).Return(permList, int64(1), nil)
			},
			expectError: false,
			expectCount: 1,
		},
		{
			name:     "normalize pageSize - too large",
			page:     1,
			pageSize: 200,
			setupMock: func(perms *MockPermissionRepository) {
				permList := []model.Permission{{Code: "test"}}
				perms.On("ListPaged", mock.Anything, 1, 20).Return(permList, int64(1), nil)
			},
			expectError: false,
			expectCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			perms := &MockPermissionRepository{}
			tt.setupMock(perms)

			svc := NewPermissionService(perms, NewMockCacheWithDefaults())
			result, _, err := svc.ListPermissionsPaged(ctx, tt.page, tt.pageSize)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectCount)
			}
		})
	}
}

func TestPermissionService_GetPermission(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		permID      uint64
		setupMock   func(*MockPermissionRepository)
		expectError bool
	}{
		{
			name:   "success - get permission",
			permID: 1,
			setupMock: func(perms *MockPermissionRepository) {
				perm := &model.Permission{Code: "users:read", Method: "GET", Path: "/users"}
				perm.ID = 1
				perms.On("Get", mock.Anything, uint64(1)).Return(perm, nil)
			},
			expectError: false,
		},
		{
			name:   "error - permission not found",
			permID: 999,
			setupMock: func(perms *MockPermissionRepository) {
				perms.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			perms := &MockPermissionRepository{}
			tt.setupMock(perms)

			svc := NewPermissionService(perms, NewMockCacheWithDefaults())
			result, err := svc.GetPermission(ctx, tt.permID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestPermissionService_ListPermissionsByGroup(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		group       string
		setupMock   func(*MockPermissionRepository)
		expectError bool
		expectCount int
	}{
		{
			name:  "success - list permissions by group",
			group: "users",
			setupMock: func(perms *MockPermissionRepository) {
				grouped := map[string][]model.Permission{
					"users": {
						{Code: "users:read"},
						{Code: "users:write"},
					},
					"orders": {
						{Code: "orders:read"},
					},
				}
				perms.On("ListByGroup", mock.Anything).Return(grouped, nil)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name:  "success - group not found returns empty",
			group: "nonexistent",
			setupMock: func(perms *MockPermissionRepository) {
				grouped := map[string][]model.Permission{
					"users": {{Code: "users:read"}},
				}
				perms.On("ListByGroup", mock.Anything).Return(grouped, nil)
			},
			expectError: false,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			perms := &MockPermissionRepository{}
			tt.setupMock(perms)

			svc := NewPermissionService(perms, NewMockCacheWithDefaults())
			result, err := svc.ListPermissionsByGroup(ctx, tt.group)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectCount)
			}
		})
	}
}

// ============================================================================
// Order Related Tests - GetOrderPayments, GetOrderRefunds
// ============================================================================

func TestAdminService_GetOrderPayments(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		orderID     uint64
		setupMock   func(*MockPaymentRepository)
		expectError bool
		expectCount int
	}{
		{
			name:    "success - get order payments",
			orderID: 1,
			setupMock: func(payments *MockPaymentRepository) {
				paymentList := []model.Payment{
					{Method: model.PaymentMethodWeChat, AmountCents: 1000, Status: model.PaymentStatusPaid},
					{Method: model.PaymentMethodAlipay, AmountCents: 500, Status: model.PaymentStatusPending},
				}
				paymentList[0].ID = 1
				paymentList[1].ID = 2
				payments.On("List", mock.Anything, mock.MatchedBy(func(opts repository.PaymentListOptions) bool {
					return opts.OrderID != nil && *opts.OrderID == uint64(1)
				})).Return(paymentList, int64(2), nil)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name:    "success - no payments found",
			orderID: 2,
			setupMock: func(payments *MockPaymentRepository) {
				payments.On("List", mock.Anything, mock.MatchedBy(func(opts repository.PaymentListOptions) bool {
					return opts.OrderID != nil && *opts.OrderID == uint64(2)
				})).Return([]model.Payment{}, int64(0), nil)
			},
			expectError: false,
			expectCount: 0,
		},
		{
			name:    "error - repository error",
			orderID: 3,
			setupMock: func(payments *MockPaymentRepository) {
				payments.On("List", mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payments := &MockPaymentRepository{}
			tt.setupMock(payments)

			svc := NewAdminService(
				&MockGameRepository{},
				&MockUserRepository{},
				&MockPlayerRepository{},
				&MockOrderRepository{},
				payments,
				&MockRoleRepository{},
				nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			result, err := svc.GetOrderPayments(ctx, tt.orderID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectCount)
			}
		})
	}
}

func TestAdminService_GetOrderRefunds(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name        string
		orderID     uint64
		setupMock   func(*MockOrderRepository, *MockPaymentRepository)
		expectError bool
		expectCount int
	}{
		{
			name:    "success - get refunds from payments",
			orderID: 1,
			setupMock: func(orders *MockOrderRepository, payments *MockPaymentRepository) {
				order := &model.Order{
					TotalPriceCents:   2000,
					RefundAmountCents: 0,
					RefundReason:      "",
					Status:            model.OrderStatusRefunded,
				}
				order.ID = 1
				orders.On("Get", mock.Anything, uint64(1)).Return(order, nil)

				refundedAt := now
				paymentList := []model.Payment{
					{Method: model.PaymentMethodWeChat, AmountCents: 1000, Status: model.PaymentStatusRefunded, RefundedAt: &refundedAt},
				}
				paymentList[0].ID = 1
				payments.On("List", mock.Anything, mock.MatchedBy(func(opts repository.PaymentListOptions) bool {
					return opts.OrderID != nil && *opts.OrderID == uint64(1)
				})).Return(paymentList, int64(1), nil)
			},
			expectError: false,
			expectCount: 1,
		},
		{
			name:    "success - get refunds from order field",
			orderID: 2,
			setupMock: func(orders *MockOrderRepository, payments *MockPaymentRepository) {
				refundedAt := now
				order := &model.Order{
					TotalPriceCents:   2000,
					RefundAmountCents: 1500,
					RefundReason:      "User requested",
					RefundedAt:        &refundedAt,
					Status:            model.OrderStatusRefunded,
				}
				order.ID = 2
				orders.On("Get", mock.Anything, uint64(2)).Return(order, nil)

				payments.On("List", mock.Anything, mock.MatchedBy(func(opts repository.PaymentListOptions) bool {
					return opts.OrderID != nil && *opts.OrderID == uint64(2)
				})).Return([]model.Payment{}, int64(0), nil)
			},
			expectError: false,
			expectCount: 1,
		},
		{
			name:    "error - order not found",
			orderID: 999,
			setupMock: func(orders *MockOrderRepository, payments *MockPaymentRepository) {
				orders.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders := &MockOrderRepository{}
			payments := &MockPaymentRepository{}
			tt.setupMock(orders, payments)

			svc := NewAdminService(
				&MockGameRepository{},
				&MockUserRepository{},
				&MockPlayerRepository{},
				orders,
				payments,
				&MockRoleRepository{},
				nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			result, err := svc.GetOrderRefunds(ctx, tt.orderID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectCount)
			}
		})
	}
}

// ============================================================================
// Payment Service Tests - CreatePayment, CapturePayment
// ============================================================================

func TestAdminService_CreatePayment(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		input       CreatePaymentInput
		setupMock   func(*MockOrderRepository, *MockUserRepository, *MockPaymentRepository)
		expectError bool
	}{
		{
			name: "success - create payment",
			input: CreatePaymentInput{
				OrderID:     1,
				UserID:      1,
				Method:      model.PaymentMethodWeChat,
				AmountCents: 1000,
				Currency:    model.CurrencyCNY,
			},
			setupMock: func(orders *MockOrderRepository, users *MockUserRepository, payments *MockPaymentRepository) {
				order := &model.Order{Status: model.OrderStatusPending}
				order.ID = 1
				orders.On("Get", mock.Anything, uint64(1)).Return(order, nil)

				user := &model.User{Name: "Test User"}
				user.ID = 1
				users.On("Get", mock.Anything, uint64(1)).Return(user, nil)

				payments.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "error - invalid order ID",
			input: CreatePaymentInput{
				OrderID:     0,
				UserID:      1,
				Method:      model.PaymentMethodWeChat,
				AmountCents: 1000,
				Currency:    model.CurrencyCNY,
			},
			setupMock: func(orders *MockOrderRepository, users *MockUserRepository, payments *MockPaymentRepository) {
			},
			expectError: true,
		},
		{
			name: "error - invalid user ID",
			input: CreatePaymentInput{
				OrderID:     1,
				UserID:      0,
				Method:      model.PaymentMethodWeChat,
				AmountCents: 1000,
				Currency:    model.CurrencyCNY,
			},
			setupMock: func(orders *MockOrderRepository, users *MockUserRepository, payments *MockPaymentRepository) {
			},
			expectError: true,
		},
		{
			name: "error - invalid amount",
			input: CreatePaymentInput{
				OrderID:     1,
				UserID:      1,
				Method:      model.PaymentMethodWeChat,
				AmountCents: 0,
				Currency:    model.CurrencyCNY,
			},
			setupMock: func(orders *MockOrderRepository, users *MockUserRepository, payments *MockPaymentRepository) {
			},
			expectError: true,
		},
		{
			name: "error - empty method",
			input: CreatePaymentInput{
				OrderID:     1,
				UserID:      1,
				Method:      "",
				AmountCents: 1000,
				Currency:    model.CurrencyCNY,
			},
			setupMock: func(orders *MockOrderRepository, users *MockUserRepository, payments *MockPaymentRepository) {
			},
			expectError: true,
		},
		{
			name: "error - order not found",
			input: CreatePaymentInput{
				OrderID:     999,
				UserID:      1,
				Method:      model.PaymentMethodWeChat,
				AmountCents: 1000,
				Currency:    model.CurrencyCNY,
			},
			setupMock: func(orders *MockOrderRepository, users *MockUserRepository, payments *MockPaymentRepository) {
				orders.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
		{
			name: "error - user not found",
			input: CreatePaymentInput{
				OrderID:     1,
				UserID:      999,
				Method:      model.PaymentMethodWeChat,
				AmountCents: 1000,
				Currency:    model.CurrencyCNY,
			},
			setupMock: func(orders *MockOrderRepository, users *MockUserRepository, payments *MockPaymentRepository) {
				order := &model.Order{Status: model.OrderStatusPending}
				order.ID = 1
				orders.On("Get", mock.Anything, uint64(1)).Return(order, nil)
				users.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders := &MockOrderRepository{}
			users := &MockUserRepository{}
			payments := &MockPaymentRepository{}
			tt.setupMock(orders, users, payments)

			svc := NewAdminService(
				&MockGameRepository{},
				users,
				&MockPlayerRepository{},
				orders,
				payments,
				&MockRoleRepository{},
				nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			result, err := svc.CreatePayment(ctx, tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestAdminService_CapturePayment(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		paymentID   uint64
		input       CapturePaymentInput
		setupMock   func(*MockPaymentRepository)
		expectError bool
	}{
		{
			name:      "success - capture payment",
			paymentID: 1,
			input: CapturePaymentInput{
				ProviderTradeNo: "TRADE123",
			},
			setupMock: func(payments *MockPaymentRepository) {
				payment := &model.Payment{
					Status:      model.PaymentStatusPending,
					AmountCents: 1000,
				}
				payment.ID = 1
				payments.On("Get", mock.Anything, uint64(1)).Return(payment, nil)
				payments.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name:      "error - payment not found",
			paymentID: 999,
			input: CapturePaymentInput{
				ProviderTradeNo: "TRADE123",
			},
			setupMock: func(payments *MockPaymentRepository) {
				payments.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
		{
			name:      "error - invalid status transition",
			paymentID: 2,
			input: CapturePaymentInput{
				ProviderTradeNo: "TRADE123",
			},
			setupMock: func(payments *MockPaymentRepository) {
				payment := &model.Payment{
					Status:      model.PaymentStatusRefunded,
					AmountCents: 1000,
				}
				payment.ID = 2
				payments.On("Get", mock.Anything, uint64(2)).Return(payment, nil)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payments := &MockPaymentRepository{}
			tt.setupMock(payments)

			svc := NewAdminService(
				&MockGameRepository{},
				&MockUserRepository{},
				&MockPlayerRepository{},
				&MockOrderRepository{},
				payments,
				&MockRoleRepository{},
				nil, nil, nil, nil, nil, nil,
				NewMockCacheWithDefaults(),
			)

			result, err := svc.CapturePayment(ctx, tt.paymentID, tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, model.PaymentStatusPaid, result.Status)
			}
		})
	}
}
