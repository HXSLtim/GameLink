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
		user.ID = uint64(len(m.users) + 1)
	}
	m.users[user.ID] = user
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

func (m *MockPlayerRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, status model.VerificationStatus) (int64, error) {
	args := m.Called(ctx, ids, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPlayerRepository) BatchDelete(ctx context.Context, ids []uint64) (int64, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPlayerRepository) BatchUpdateHourlyRate(ctx context.Context, ids []uint64, rateCents int64) (int64, error) {
	args := m.Called(ctx, ids, rateCents)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPlayerRepository) BatchUpdateRank(ctx context.Context, ids []uint64, rank string) (int64, error) {
	args := m.Called(ctx, ids, rank)
	return args.Get(0).(int64), args.Error(1)
}

// MockOrderRepository is a mock implementation of OrderRepository
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

func (m *MockOrderRepository) List(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.Order), args.Get(1).(int64), args.Error(2)
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

// MockPaymentRepository is a mock implementation of PaymentRepository
type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
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

func (m *MockPaymentRepository) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.Payment), args.Get(1).(int64), args.Error(2)
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
		return []model.Payment{}, args.Error(1)
	}
	return args.Get(0).([]model.Payment), args.Error(1)
}

// MockRoleRepository is a mock implementation of RoleRepository
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) Create(ctx context.Context, role *model.RoleModel) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) Get(ctx context.Context, id uint64) (*model.RoleModel, error) {
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

func (m *MockRoleRepository) List(ctx context.Context) ([]model.RoleModel, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.RoleModel), args.Error(1)
}

func (m *MockRoleRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.RoleModel, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]model.RoleModel), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoleRepository) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, isSystem *bool) ([]model.RoleModel, int64, error) {
	args := m.Called(ctx, page, pageSize, keyword, isSystem)
	return args.Get(0).([]model.RoleModel), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoleRepository) Update(ctx context.Context, role *model.RoleModel) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
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

func (m *MockRoleRepository) GetChildRoles(ctx context.Context, parentID uint64) ([]model.RoleModel, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).([]model.RoleModel), args.Error(1)
}

func (m *MockRoleRepository) ListByUserID(ctx context.Context, userID uint64) ([]model.RoleModel, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.RoleModel), args.Error(1)
}

func (m *MockRoleRepository) CheckUserHasRole(ctx context.Context, userID uint64, roleSlug string) (bool, error) {
	args := m.Called(ctx, userID, roleSlug)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoleRepository) GetUserIDsByRoleID(ctx context.Context, roleID uint64) ([]uint64, error) {
	args := m.Called(ctx, roleID)
	return args.Get(0).([]uint64), args.Error(1)
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

func (m *MockRoleRepository) GetWithPermissions(ctx context.Context, id uint64) (*model.RoleModel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RoleModel), args.Error(1)
}

func (m *MockRoleRepository) ListWithPermissions(ctx context.Context) ([]model.RoleModel, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.RoleModel), args.Error(1)
}

func (m *MockRoleRepository) SetParent(ctx context.Context, roleID uint64, parentID *uint64) error {
	args := m.Called(ctx, roleID, parentID)
	return args.Error(0)
}

func (m *MockRoleRepository) GetInheritanceChain(ctx context.Context, roleID uint64) ([]model.RoleModel, error) {
	args := m.Called(ctx, roleID)
	return args.Get(0).([]model.RoleModel), args.Error(1)
}

func (m *MockRoleRepository) UpdateLevel(ctx context.Context, roleID uint64, level int) error {
	args := m.Called(ctx, roleID, level)
	return args.Error(0)
}

// MockServiceItemRepository is a mock implementation of ServiceItemRepository
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
	return args.Get(0).([]model.ServiceItem), args.Get(1).(int64), args.Error(2)
}

func (m *MockServiceItemRepository) GetGameServices(ctx context.Context, gameID uint64, subCategory *model.ServiceItemSubCategory) ([]model.ServiceItem, error) {
	args := m.Called(ctx, gameID, subCategory)
	if args.Get(0) == nil {
		return []model.ServiceItem{}, args.Error(1)
	}
	return args.Get(0).([]model.ServiceItem), args.Error(1)
}

// MockPermissionRepository is a mock implementation of PermissionRepository
type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) Create(ctx context.Context, perm *model.Permission) error {
	args := m.Called(ctx, perm)
	return args.Error(0)
}

func (m *MockPermissionRepository) Get(ctx context.Context, id uint64) (*model.Permission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) List(ctx context.Context) ([]model.Permission, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Permission, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]model.Permission), args.Get(1).(int64), args.Error(2)
}

func (m *MockPermissionRepository) Update(ctx context.Context, perm *model.Permission) error {
	args := m.Called(ctx, perm)
	return args.Error(0)
}

func (m *MockPermissionRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPermissionRepository) ListByUserID(ctx context.Context, userID uint64) ([]model.Permission, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) ListByRoleID(ctx context.Context, roleID uint64) ([]model.Permission, error) {
	args := m.Called(ctx, roleID)
	return args.Get(0).([]model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) CountRoleReferences(ctx context.Context, id uint64) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPermissionRepository) GetByMethodAndPath(ctx context.Context, method, path string) (*model.Permission, error) {
	args := m.Called(ctx, method, path)
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

func (m *MockPermissionRepository) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword, method, group string, isSystem *bool) ([]model.Permission, int64, error) {
	args := m.Called(ctx, page, pageSize, keyword, method, group, isSystem)
	return args.Get(0).([]model.Permission), args.Get(1).(int64), args.Error(2)
}

func (m *MockPermissionRepository) ListByGroup(ctx context.Context) (map[string][]model.Permission, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string][]model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) UpsertByMethodPath(ctx context.Context, perm *model.Permission) error {
	args := m.Called(ctx, perm)
	return args.Error(0)
}

func (m *MockPermissionRepository) ListWithChildren(ctx context.Context) ([]model.Permission, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) ListGroups(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockPermissionRepository) GetByResource(ctx context.Context, resource, action string) (*model.Permission, error) {
	args := m.Called(ctx, resource, action)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) GetWithChildren(ctx context.Context, id uint64) (*model.Permission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Permission), args.Error(1)
}

// MockMenuRepository is a mock implementation of MenuRepository
type MockMenuRepository struct {
	mock.Mock
}

func (m *MockMenuRepository) Create(ctx context.Context, menu *model.Menu) error {
	args := m.Called(ctx, menu)
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
	return args.Get(0).([]model.Menu), args.Error(1)
}

func (m *MockMenuRepository) ListPaged(ctx context.Context, page, pageSize int, parentID *uint64) ([]model.Menu, int64, error) {
	args := m.Called(ctx, page, pageSize, parentID)
	return args.Get(0).([]model.Menu), args.Get(1).(int64), args.Error(2)
}

func (m *MockMenuRepository) Update(ctx context.Context, menu *model.Menu) error {
	args := m.Called(ctx, menu)
	return args.Error(0)
}

func (m *MockMenuRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMenuRepository) HasChildren(ctx context.Context, id uint64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockMenuRepository) ListByPermission(ctx context.Context, codes []string) ([]model.Menu, error) {
	args := m.Called(ctx, codes)
	return args.Get(0).([]model.Menu), args.Error(1)
}

// MockStatsRepository is a mock implementation of StatsRepository
type MockStatsRepository struct {
	mock.Mock
}

func (m *MockStatsRepository) Dashboard(ctx context.Context) (repository.Dashboard, error) {
	args := m.Called(ctx)
	return args.Get(0).(repository.Dashboard), args.Error(1)
}

func (m *MockStatsRepository) RevenueTrend(ctx context.Context, days int) ([]repository.DateValue, error) {
	args := m.Called(ctx, days)
	return args.Get(0).([]repository.DateValue), args.Error(1)
}

func (m *MockStatsRepository) UserGrowth(ctx context.Context, days int) ([]repository.DateValue, error) {
	args := m.Called(ctx, days)
	return args.Get(0).([]repository.DateValue), args.Error(1)
}

func (m *MockStatsRepository) OrdersByStatus(ctx context.Context) (map[string]int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *MockStatsRepository) TopPlayers(ctx context.Context, limit int) ([]repository.PlayerTop, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]repository.PlayerTop), args.Error(1)
}

func (m *MockStatsRepository) AuditOverview(ctx context.Context, from, to *time.Time) (map[string]int64, map[string]int64, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).(map[string]int64), args.Get(1).(map[string]int64), args.Error(2)
}

func (m *MockStatsRepository) AuditTrend(ctx context.Context, from, to *time.Time, entity, action string) ([]repository.DateValue, error) {
	args := m.Called(ctx, from, to, entity, action)
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
		setupMocks  func(*MockGameRepository)
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
			setupMocks: func(m *MockGameRepository) {
				m.On("Create", ctx, mock.AnythingOfType("*model.Game")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "validation error - empty key",
			input: CreateGameInput{
				Key:  "",
				Name: "Test Game",
			},
			setupMocks:  func(m *MockGameRepository) {},
			wantErr:     true,
			errContains: "validation",
		},
		{
			name: "validation error - empty name",
			input: CreateGameInput{
				Key:  "test-game",
				Name: "",
			},
			setupMocks:  func(m *MockGameRepository) {},
			wantErr:     true,
			errContains: "validation",
		},
		{
			name: "repository error",
			input: CreateGameInput{
				Key:  "test-game",
				Name: "Test Game",
			},
			setupMocks: func(m *MockGameRepository) {
				m.On("Create", ctx, mock.AnythingOfType("*model.Game")).Return(errors.New("db error"))
			},
			wantErr:     true,
			errContains: "create game",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			games := &MockGameRepository{}
			tt.setupMocks(games)

			svc := NewAdminService(
				games, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				&MockCache{},
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
		setupMocks  func(*MockGameRepository)
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
			setupMocks: func(m *MockGameRepository) {
				m.On("Get", ctx, uint64(1)).Return(&model.Game{
					Base: model.Base{ID: 1, ExtJSON: "{}"},
					Key:  "old-game",
					Name: "Old Game",
				}, nil)
				m.On("Update", ctx, mock.AnythingOfType("*model.Game")).Return(nil)
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
			setupMocks: func(m *MockGameRepository) {
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
			setupMocks:  func(m *MockGameRepository) {},
			wantErr:     true,
			errContains: "validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			games := &MockGameRepository{}
			tt.setupMocks(games)

			svc := NewAdminService(
				games, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				&MockCache{},
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
		setupMocks  func(*MockGameRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:   "success - delete game",
			gameID: 1,
			setupMocks: func(m *MockGameRepository) {
				m.On("Delete", ctx, uint64(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "game not found",
			gameID: 999,
			setupMocks: func(m *MockGameRepository) {
				m.On("Delete", ctx, uint64(999)).Return(repository.ErrNotFound)
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name:   "repository error",
			gameID: 1,
			setupMocks: func(m *MockGameRepository) {
				m.On("Delete", ctx, uint64(1)).Return(errors.New("db error"))
			},
			wantErr:     true,
			errContains: "delete game",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			games := &MockGameRepository{}
			tt.setupMocks(games)

			svc := NewAdminService(
				games, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				&MockCache{},
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
		setupMocks func(*MockGameRepository)
		wantErr    bool
		wantCount  int
	}{
		{
			name: "success - list games",
			setupMocks: func(m *MockGameRepository) {
				m.On("List", ctx).Return([]model.Game{
					{Base: model.Base{ID: 1, ExtJSON: "{}"}, Key: "game1", Name: "Game 1"},
					{Base: model.Base{ID: 2, ExtJSON: "{}"}, Key: "game2", Name: "Game 2"},
				}, nil)
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name: "repository error",
			setupMocks: func(m *MockGameRepository) {
				m.On("List", ctx).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			games := &MockGameRepository{}
			tt.setupMocks(games)

			svc := NewAdminService(
				games, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				&MockCache{},
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
		setupMocks  func(*MockGameRepository)
		wantErr     bool
		errContains string
		wantDeleted int64
	}{
		{
			name: "success - batch delete",
			ids:  []uint64{1, 2, 3},
			setupMocks: func(m *MockGameRepository) {
				m.On("BatchDelete", ctx, []uint64{1, 2, 3}).Return(int64(3), nil)
			},
			wantErr:     false,
			wantDeleted: 3,
		},
		{
			name:        "empty ids",
			ids:         []uint64{},
			setupMocks:  func(m *MockGameRepository) {},
			wantErr:     true,
			errContains: "no game ids",
		},
		{
			name: "repository error",
			ids:  []uint64{1, 2},
			setupMocks: func(m *MockGameRepository) {
				m.On("BatchDelete", ctx, []uint64{1, 2}).Return(int64(0), errors.New("db error"))
			},
			wantErr:     true,
			errContains: "batch delete games",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			games := &MockGameRepository{}
			tt.setupMocks(games)

			svc := NewAdminService(
				games, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				&MockCache{},
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
		setupMocks  func(*MockUserRepository, *MockWalletRepository, *MockRoleRepository)
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
			setupMocks: func(u *MockUserRepository, w *MockWalletRepository, r *MockRoleRepository) {
				u.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(nil)
				w.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).Return(nil)
				r.On("GetBySlug", ctx, string(model.RoleSlugUser)).Return(&model.RoleModel{
					Base: model.Base{ID: 1, ExtJSON: "{}"},
					Slug: string(model.RoleSlugUser),
				}, nil)
				r.On("AssignToUser", ctx, mock.AnythingOfType("uint64"), mock.AnythingOfType("[]uint64")).Return(nil)
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
			setupMocks:  func(u *MockUserRepository, w *MockWalletRepository, r *MockRoleRepository) {},
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
			setupMocks:  func(u *MockUserRepository, w *MockWalletRepository, r *MockRoleRepository) {},
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
			setupMocks:  func(u *MockUserRepository, w *MockWalletRepository, r *MockRoleRepository) {},
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
			setupMocks: func(u *MockUserRepository, w *MockWalletRepository, r *MockRoleRepository) {
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
			tt.setupMocks(users, wallets, roles)

			svc := NewAdminService(
				nil, users, nil, nil, nil, roles, nil, nil, nil, nil, wallets, nil,
				&MockCache{},
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

	existingUser := &model.User{
		Base:         model.Base{ID: 1, ExtJSON: "{}"},
		Name:         "Old Name",
		Phone:        "13800138000",
		Email:        "old@example.com",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: hashPasswordForTest("Old123!@#"),
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
				u.On("Get", ctx, uint64(1)).Return(existingUser, nil)
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
				u.On("Get", ctx, uint64(1)).Return(existingUser, nil)
				u.On("Update", ctx, mock.MatchedBy(func(u *model.User) bool {
					return u.PasswordHash != existingUser.PasswordHash
				})).Return(nil)
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
				&MockCache{},
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
				&MockCache{},
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
				&MockCache{},
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
				&MockCache{},
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
				&MockCache{},
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
				&MockCache{},
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
				&MockCache{},
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
				&MockCache{},
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
				&MockCache{},
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
				&MockCache{},
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
				&MockCache{},
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
				&MockCache{},
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
				&MockCache{},
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
				&MockCache{},
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
			errContains: "invalid transition",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders := &MockOrderRepository{}
			tt.setupMocks(orders)

			svc := NewAdminService(
				nil, nil, nil, orders, nil, nil, nil, nil, nil, nil, nil, nil,
				&MockCache{},
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
	cache := &MockCache{}

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
