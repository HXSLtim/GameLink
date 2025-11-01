package repository

import (
	"context"
	"time"

	"gamelink/internal/model"
)

// GameRepository defines game data access operations.
type GameRepository interface {
	List(ctx context.Context) ([]model.Game, error)
	ListPaged(ctx context.Context, page, pageSize int) ([]model.Game, int64, error)
	Get(ctx context.Context, id uint64) (*model.Game, error)
	Create(ctx context.Context, game *model.Game) error
	Update(ctx context.Context, game *model.Game) error
	Delete(ctx context.Context, id uint64) error
}

// UserRepository defines user data access operations.
type UserRepository interface {
	List(ctx context.Context) ([]model.User, error)
	ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error)
	ListWithFilters(ctx context.Context, opts UserListOptions) ([]model.User, int64, error)
	Get(ctx context.Context, id uint64) (*model.User, error)
	GetByPhone(ctx context.Context, phone string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByPhone(ctx context.Context, phone string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint64) error
}

// PlayerRepository defines player data access operations.
type PlayerRepository interface {
	List(ctx context.Context) ([]model.Player, error)
	ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error)
	Get(ctx context.Context, id uint64) (*model.Player, error)
	Create(ctx context.Context, player *model.Player) error
	Update(ctx context.Context, player *model.Player) error
	Delete(ctx context.Context, id uint64) error
}

// OrderRepository defines order data access operations.
type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	List(ctx context.Context, opts OrderListOptions) ([]model.Order, int64, error)
	Get(ctx context.Context, id uint64) (*model.Order, error)
	Update(ctx context.Context, order *model.Order) error
	Delete(ctx context.Context, id uint64) error
}

// PaymentRepository defines payment data access operations.
type PaymentRepository interface {
	Create(ctx context.Context, payment *model.Payment) error
	List(ctx context.Context, opts PaymentListOptions) ([]model.Payment, int64, error)
	Get(ctx context.Context, id uint64) (*model.Payment, error)
	Update(ctx context.Context, payment *model.Payment) error
	Delete(ctx context.Context, id uint64) error
}

// PermissionRepository defines permission data access operations.
type PermissionRepository interface {
	List(ctx context.Context) ([]model.Permission, error)
	ListPaged(ctx context.Context, page, pageSize int) ([]model.Permission, int64, error)
	ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword, method, group string) ([]model.Permission, int64, error)
	ListByGroup(ctx context.Context) (map[string][]model.Permission, error)
	ListGroups(ctx context.Context) ([]string, error)
	Get(ctx context.Context, id uint64) (*model.Permission, error)
	GetByResource(ctx context.Context, resource, action string) (*model.Permission, error)
	GetByMethodAndPath(ctx context.Context, method, path string) (*model.Permission, error)
	Create(ctx context.Context, perm *model.Permission) error
	Update(ctx context.Context, perm *model.Permission) error
	UpsertByMethodPath(ctx context.Context, perm *model.Permission) error
	Delete(ctx context.Context, id uint64) error
	ListByRoleID(ctx context.Context, roleID uint64) ([]model.Permission, error)
	ListByUserID(ctx context.Context, userID uint64) ([]model.Permission, error)
}

// RoleRepository defines role data access operations.
type RoleRepository interface {
	List(ctx context.Context) ([]model.RoleModel, error)
	ListPaged(ctx context.Context, page, pageSize int) ([]model.RoleModel, int64, error)
	ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, isSystem *bool) ([]model.RoleModel, int64, error)
	ListWithPermissions(ctx context.Context) ([]model.RoleModel, error)
	Get(ctx context.Context, id uint64) (*model.RoleModel, error)
	GetWithPermissions(ctx context.Context, id uint64) (*model.RoleModel, error)
	GetBySlug(ctx context.Context, slug string) (*model.RoleModel, error)
	Create(ctx context.Context, role *model.RoleModel) error
	Update(ctx context.Context, role *model.RoleModel) error
	Delete(ctx context.Context, id uint64) error
	AssignPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error
	AddPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error
	RemovePermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error
	AssignToUser(ctx context.Context, userID uint64, roleIDs []uint64) error
	RemoveFromUser(ctx context.Context, userID uint64, roleIDs []uint64) error
	ListByUserID(ctx context.Context, userID uint64) ([]model.RoleModel, error)
	CheckUserHasRole(ctx context.Context, userID uint64, roleSlug string) (bool, error)
}

// PlayerTagRepository defines player tag data access operations.
type PlayerTagRepository interface {
	GetTags(ctx context.Context, playerID uint64) ([]string, error)
	ReplaceTags(ctx context.Context, playerID uint64, tags []string) error
}

// ReviewRepository defines review data access operations.
type ReviewRepository interface {
	List(ctx context.Context, opts ReviewListOptions) ([]model.Review, int64, error)
	Get(ctx context.Context, id uint64) (*model.Review, error)
	Create(ctx context.Context, review *model.Review) error
	Update(ctx context.Context, review *model.Review) error
	Delete(ctx context.Context, id uint64) error
}

// OperationLogRepository defines operation log data access operations.
type OperationLogRepository interface {
	Append(ctx context.Context, log *model.OperationLog) error
	ListByEntity(ctx context.Context, entityType string, entityID uint64, opts OperationLogListOptions) ([]model.OperationLog, int64, error)
}

// StatsRepository provides statistical query capabilities.
type StatsRepository interface {
	Dashboard(ctx context.Context) (Dashboard, error)
	RevenueTrend(ctx context.Context, days int) ([]DateValue, error)
	UserGrowth(ctx context.Context, days int) ([]DateValue, error)
	OrdersByStatus(ctx context.Context) (map[string]int64, error)
	TopPlayers(ctx context.Context, limit int) ([]PlayerTop, error)
	AuditOverview(ctx context.Context, from, to *time.Time) (map[string]int64, map[string]int64, error)
	AuditTrend(ctx context.Context, from, to *time.Time, entity, action string) ([]DateValue, error)
}

// UserListOptions contains filtering options for user queries.
type UserListOptions struct {
	Page     int
	PageSize int
	Role     model.Role
	Roles    []model.Role
	Status   model.UserStatus
	Statuses []model.UserStatus
	Keyword  string
	DateFrom *time.Time
	DateTo   *time.Time
}

// OrderListOptions contains filtering options for order queries.
type OrderListOptions struct {
	Page     int
	PageSize int
	UserID   *uint64
	PlayerID *uint64
	GameID   *uint64
	Statuses []model.OrderStatus
	Keyword  string
	DateFrom *time.Time
	DateTo   *time.Time
}

// PaymentListOptions contains filtering options for payment queries.
type PaymentListOptions struct {
	Page     int
	PageSize int
	OrderID  *uint64
	UserID   *uint64
	Method   *model.PaymentMethod
	Methods  []model.PaymentMethod
	Status   *model.PaymentStatus
	Statuses []model.PaymentStatus
	DateFrom *time.Time
	DateTo   *time.Time
}

// ReviewListOptions contains filtering options for review queries.
type ReviewListOptions struct {
	Page     int
	PageSize int
	OrderID  *uint64
	UserID   *uint64
	PlayerID *uint64
	DateFrom *time.Time
	DateTo   *time.Time
}

// OperationLogListOptions contains filtering options for operation log queries.
type OperationLogListOptions struct {
	Page        int
	PageSize    int
	Action      string
	ActorUserID *uint64
	DateFrom    *time.Time
	DateTo      *time.Time
}

// Dashboard aggregates summary data for the homepage.
type Dashboard struct {
	TotalUsers           int64            `json:"totalUsers"`
	TotalPlayers         int64            `json:"totalPlayers"`
	TotalGames           int64            `json:"totalGames"`
	TotalOrders          int64            `json:"totalOrders"`
	OrdersByStatus       map[string]int64 `json:"ordersByStatus"`
	PaymentsByStatus     map[string]int64 `json:"paymentsByStatus"`
	TotalPaidAmountCents int64            `json:"totalPaidAmountCents"`
}

// DateValue represents a value aggregated by date.
type DateValue struct {
	Date  string `json:"date"`
	Value int64  `json:"value"`
}

// PlayerTop represents a leaderboard entry.
type PlayerTop struct {
	PlayerID      uint64  `json:"playerId"`
	Nickname      string  `json:"nickname"`
	RatingAverage float32 `json:"ratingAverage"`
	RatingCount   uint32  `json:"ratingCount"`
}
