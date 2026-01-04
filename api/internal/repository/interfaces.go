package repository

import (
	"context"
	"time"

	"gamelink/internal/model"
	repoiface "gamelink/internal/repository/interfaces"
)

// GameRepository defines game data access operations.
type GameRepository interface {
	List(ctx context.Context) ([]model.Game, error)
	ListPaged(ctx context.Context, page, pageSize int) ([]model.Game, int64, error)
	ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string) ([]model.Game, int64, error)
	Get(ctx context.Context, id uint64) (*model.Game, error)
	GetByIDs(ctx context.Context, ids []uint64) ([]model.Game, error)
	Create(ctx context.Context, game *model.Game) error
	Update(ctx context.Context, game *model.Game) error
	Delete(ctx context.Context, id uint64) error
	BatchDelete(ctx context.Context, ids []uint64) (int64, error)
	BatchUpdateStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error)
	BatchUpdateSortOrder(ctx context.Context, updates map[uint64]int) (int64, error)
	BatchUpdateCategory(ctx context.Context, ids []uint64, category string) (int64, error)
}

// UserRepository defines user data access operations.
type UserRepository interface {
	List(ctx context.Context) ([]model.User, error)
	ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error)
	ListWithFilters(ctx context.Context, opts UserListOptions) ([]model.User, int64, error)
	Count(ctx context.Context, opts UserListOptions) (int, error)
	Get(ctx context.Context, id uint64) (*model.User, error)
	GetByIDs(ctx context.Context, ids []uint64) ([]model.User, error)
	GetByPhone(ctx context.Context, phone string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByPhone(ctx context.Context, phone string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	UpdatePassword(ctx context.Context, userID uint64, newPassword string) error
	Delete(ctx context.Context, id uint64) error
}

// PlayerRepository defines player data access operations.
type PlayerRepository interface {
	List(ctx context.Context) ([]model.Player, error)
	ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error)
	ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, status *model.VerificationStatus) ([]model.Player, int64, error)
	Get(ctx context.Context, id uint64) (*model.Player, error)
	GetByIDs(ctx context.Context, ids []uint64) ([]model.Player, error)
	GetByUserID(ctx context.Context, userID uint64) (*model.Player, error)
	Create(ctx context.Context, player *model.Player) error
	Update(ctx context.Context, player *model.Player) error
	Delete(ctx context.Context, id uint64) error
	BatchUpdateRank(ctx context.Context, ids []uint64, rank string) (int64, error)
	BatchUpdateHourlyRate(ctx context.Context, ids []uint64, rateCents int64) (int64, error)
	BatchUpdateStatus(ctx context.Context, ids []uint64, status model.VerificationStatus) (int64, error)
	BatchDelete(ctx context.Context, ids []uint64) (int64, error)
}

// Order repository interfaces now live in the interfaces subpackage. Keep
// aliases here to avoid forcing callers to update immediately.
type (
	OrderReader      = repoiface.OrderReader
	OrderWriter      = repoiface.OrderWriter
	OrderQuery       = repoiface.OrderQuery
	OrderReadWriter  = repoiface.OrderReadWriter
	OrderRepository  = repoiface.OrderRepository
	OrderListOptions = repoiface.OrderListOptions
)

// PaymentRepository defines payment data access operations.
type PaymentRepository interface {
	Create(ctx context.Context, payment *model.Payment) error
	List(ctx context.Context, opts PaymentListOptions) ([]model.Payment, int64, error)
	Get(ctx context.Context, id uint64) (*model.Payment, error)
	GetWithRelations(ctx context.Context, id uint64) (*model.Payment, error) // 获取支付记录及关联的订单和用户信息
	Update(ctx context.Context, payment *model.Payment) error
	Delete(ctx context.Context, id uint64) error
	GetByOrderID(ctx context.Context, orderID uint64) ([]model.Payment, error) // 根据订单ID获取支付记录
}

// RefundRecordRepository defines refund record data access operations.
type RefundRecordRepository interface {
	Create(ctx context.Context, record *model.RefundRecord) error
	Get(ctx context.Context, id uint64) (*model.RefundRecord, error)
	Update(ctx context.Context, record *model.RefundRecord) error
	ListByPaymentID(ctx context.Context, paymentID uint64) ([]model.RefundRecord, error)
	ListByOrderID(ctx context.Context, orderID uint64) ([]model.RefundRecord, error)
}

// WalletRepository defines wallet data access operations.
type WalletRepository interface {
	GetByUserID(ctx context.Context, userID uint64) (*model.Wallet, error)
	Save(ctx context.Context, wallet *model.Wallet) error
}

// PermissionRepository defines permission data access operations.
type PermissionRepository interface {
	List(ctx context.Context) ([]model.Permission, error)
	ListPaged(ctx context.Context, page, pageSize int) ([]model.Permission, int64, error)
	ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword, method, group string, isSystem *bool) ([]model.Permission, int64, error)
	ListByGroup(ctx context.Context) (map[string][]model.Permission, error)
	ListGroups(ctx context.Context) ([]string, error)
	Get(ctx context.Context, id uint64) (*model.Permission, error)
	GetByResource(ctx context.Context, resource, action string) (*model.Permission, error)
	GetByCode(ctx context.Context, code string) (*model.Permission, error)
	GetByMethodAndPath(ctx context.Context, method, path string) (*model.Permission, error)
	Create(ctx context.Context, perm *model.Permission) error
	Update(ctx context.Context, perm *model.Permission) error
	UpsertByMethodPath(ctx context.Context, perm *model.Permission) error
	Delete(ctx context.Context, id uint64) error
	ListByRoleID(ctx context.Context, roleID uint64) ([]model.Permission, error)
	ListByUserID(ctx context.Context, userID uint64) ([]model.Permission, error)
	// Tree structure methods
	ListWithChildren(ctx context.Context) ([]model.Permission, error)
	GetWithChildren(ctx context.Context, id uint64) (*model.Permission, error)
	// Reference check methods
	CountRoleReferences(ctx context.Context, permissionID uint64) (int64, error)
}

// MenuRepository defines admin menu / front-end route persistence.
type MenuRepository interface {
	Create(ctx context.Context, menu *model.Menu) error
	Update(ctx context.Context, menu *model.Menu) error
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, id uint64) (*model.Menu, error)
	List(ctx context.Context, parentID *uint64) ([]model.Menu, error)
	ListPaged(ctx context.Context, page, pageSize int, parentID *uint64) ([]model.Menu, int64, error)
	ListByPermission(ctx context.Context, codes []string) ([]model.Menu, error)
	HasChildren(ctx context.Context, parentID uint64) (bool, error)
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
	// Role inheritance methods
	SetParent(ctx context.Context, roleID uint64, parentID *uint64) error
	GetInheritanceChain(ctx context.Context, roleID uint64) ([]model.RoleModel, error)
	GetChildRoles(ctx context.Context, roleID uint64) ([]model.RoleModel, error)
	UpdateLevel(ctx context.Context, roleID uint64, level int) error
	// Cache invalidation support
	GetUserIDsByRoleID(ctx context.Context, roleID uint64) ([]uint64, error)
}

// PlayerTagRepository defines player tag data access operations.
type PlayerTagRepository interface {
	GetTags(ctx context.Context, playerID uint64) ([]string, error)
	ReplaceTags(ctx context.Context, playerID uint64, tags []string) error
}

// ReviewRepository defines review data access operations.
type ReviewRepository interface {
	List(ctx context.Context, opts ReviewListOptions) ([]model.Review, int64, error)
	ListPending(ctx context.Context, page, pageSize int) ([]model.Review, int64, error)
	Get(ctx context.Context, id uint64) (*model.Review, error)
	Create(ctx context.Context, review *model.Review) error
	Update(ctx context.Context, review *model.Review) error
	UpdateStatus(ctx context.Context, id uint64, status model.ReviewStatus, rejectionReason string) error
	BatchUpdateStatus(ctx context.Context, ids []uint64, status model.ReviewStatus, rejectionReason string) error
	Delete(ctx context.Context, id uint64) error
	// Statistics methods
	GetStats(ctx context.Context) (ReviewStats, error)
	GetTrend(ctx context.Context, days int) ([]DateValue, error)
	GetTopPlayersByReviewCount(ctx context.Context, limit int) ([]PlayerReviewStats, error)
	GetTopPlayersByRating(ctx context.Context, limit int) ([]PlayerReviewStats, error)
	GetGameStats(ctx context.Context) ([]GameReviewStats, error)
}

// OperationLogRepository defines operation log data access operations.
type OperationLogRepository interface {
	Append(ctx context.Context, log *model.OperationLog) error
	ListByEntity(ctx context.Context, entityType string, entityID uint64, opts OperationLogListOptions) ([]model.OperationLog, int64, error)
	List(ctx context.Context, opts OperationLogSearchOptions) ([]model.OperationLog, int64, error)
}

// ChatGroupRepository defines chat group data access operations.
type ChatGroupRepository interface {
	Create(ctx context.Context, group *model.ChatGroup) error
	Get(ctx context.Context, id uint64) (*model.ChatGroup, error)
	GetByRelatedOrderID(ctx context.Context, orderID uint64) (*model.ChatGroup, error)
	ListByUser(ctx context.Context, userID uint64, opts ChatGroupListOptions) ([]model.ChatGroup, int64, error)
	ListMembers(ctx context.Context, groupID uint64, opts ChatGroupMemberListOptions) ([]model.ChatGroupMember, int64, error)
	Update(ctx context.Context, group *model.ChatGroup) error
	Deactivate(ctx context.Context, id uint64) error
	ListDeactivatedBefore(ctx context.Context, cutoff time.Time, limit int) ([]model.ChatGroup, error)
	DeleteByIDs(ctx context.Context, ids []uint64) error
}

// ChatMemberRepository defines membership access operations.
type ChatMemberRepository interface {
	Add(ctx context.Context, member *model.ChatGroupMember) error
	AddBatch(ctx context.Context, members []*model.ChatGroupMember) error
	Update(ctx context.Context, member *model.ChatGroupMember) error
	Remove(ctx context.Context, groupID, userID uint64) error
	Get(ctx context.Context, groupID, userID uint64) (*model.ChatGroupMember, error)
}

// ChatMessageRepository defines chat message storage operations.
type ChatMessageRepository interface {
	Create(ctx context.Context, message *model.ChatMessage) error
	CreateBatch(ctx context.Context, messages []*model.ChatMessage) error
	ListByGroup(ctx context.Context, opts ChatMessageListOptions) ([]model.ChatMessage, int64, error)
	Get(ctx context.Context, id uint64) (*model.ChatMessage, error)
	MarkDeleted(ctx context.Context, id uint64, deletedBy uint64) error
	ListForModeration(ctx context.Context, opts ChatMessageModerationListOptions) ([]model.ChatMessage, int64, error)
	UpdateAuditStatus(ctx context.Context, id uint64, status model.ChatMessageAuditStatus, moderatorID *uint64, reason string) error
	DeleteByGroupIDs(ctx context.Context, groupIDs []uint64) error
}

// ChatReportRepository defines access operations for chat reports.
type ChatReportRepository interface {
	Create(ctx context.Context, report *model.ChatReport) error
	Get(ctx context.Context, id uint64) (*model.ChatReport, error)
	Update(ctx context.Context, report *model.ChatReport) error
	List(ctx context.Context, opts ChatReportListOptions) ([]model.ChatReport, int64, error)
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

// FeedRepository defines persistence methods for community feeds.
type FeedRepository interface {
	Create(ctx context.Context, feed *model.Feed) error
	Get(ctx context.Context, id uint64) (*model.Feed, error)
	List(ctx context.Context, opts FeedListOptions) ([]model.Feed, error)
	ListPaged(ctx context.Context, opts FeedPagedListOptions) ([]model.Feed, int64, error)
	Update(ctx context.Context, feed *model.Feed) error
	Delete(ctx context.Context, id uint64) error
	// UpdateModeration updates feed moderation status. moderatorID is optional (nil for auto-moderation).
	UpdateModeration(ctx context.Context, feedID uint64, status model.FeedModerationStatus, note string, moderatorID *uint64) error
	BatchUpdateModeration(ctx context.Context, feedIDs []uint64, status model.FeedModerationStatus, note string, moderatorID *uint64) error
	CreateReport(ctx context.Context, report *model.FeedReport) error
	GetReport(ctx context.Context, id uint64) (*model.FeedReport, error)
	ListReports(ctx context.Context, opts FeedReportListOptions) ([]model.FeedReport, int64, error)
	UpdateReport(ctx context.Context, report *model.FeedReport) error
	// Statistics
	CountByStatus(ctx context.Context) (map[model.FeedModerationStatus]int64, error)
	GetTrend(ctx context.Context, days int) ([]DateValue, error)
}

// NotificationRepository defines persistence for notification events.
type NotificationRepository interface {
	ListByUser(ctx context.Context, opts NotificationListOptions) ([]model.NotificationEvent, int64, error)
	MarkRead(ctx context.Context, userID uint64, ids []uint64) error
	MarkAllRead(ctx context.Context, userID uint64) error
	CountUnread(ctx context.Context, userID uint64) (int64, error)
	Create(ctx context.Context, event *model.NotificationEvent) error
	Delete(ctx context.Context, userID uint64, id uint64) error
}

// ReviewReplyRepository defines data access for review replies.
type ReviewReplyRepository interface {
	Create(ctx context.Context, reply *model.ReviewReply) error
	Get(ctx context.Context, replyID uint64) (*model.ReviewReply, error)
	ListByReview(ctx context.Context, reviewID uint64) ([]model.ReviewReply, error)
	Update(ctx context.Context, reply *model.ReviewReply) error
	Delete(ctx context.Context, replyID uint64) error
	UpdateStatus(ctx context.Context, replyID uint64, status string, note string) error
}

// ReviewReportRepository defines data access for review reports.
type ReviewReportRepository interface {
	Create(ctx context.Context, report *model.ReviewReport) error
	Get(ctx context.Context, id uint64) (*model.ReviewReport, error)
	List(ctx context.Context, opts ReviewReportListOptions) ([]model.ReviewReport, int64, error)
	Update(ctx context.Context, report *model.ReviewReport) error
}

// DisputeRepository defines data access operations for order disputes.
type DisputeRepository interface {
	Create(ctx context.Context, dispute *model.OrderDispute) error
	Get(ctx context.Context, id uint64) (*model.OrderDispute, error)
	GetByOrderID(ctx context.Context, orderID uint64) (*model.OrderDispute, error)
	Update(ctx context.Context, dispute *model.OrderDispute) error
	List(ctx context.Context, opts DisputeListOptions) ([]model.OrderDispute, int64, error)
	ListPendingAssignment(ctx context.Context, page, pageSize int) ([]model.OrderDispute, int64, error)
	ListSLABreached(ctx context.Context) ([]model.OrderDispute, error)
	MarkSLABreached(ctx context.Context, disputeID uint64) error
	Delete(ctx context.Context, id uint64) error
	CountByStatus(ctx context.Context, status model.DisputeStatus) (int64, error)
	GetPendingCount(ctx context.Context) (int64, error)
	GetStats(ctx context.Context) (map[string]int64, error)
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

// FeedListOptions describes feed query filters (cursor-based).
type FeedListOptions struct {
	Limit        int
	CursorBefore *uint64
	AuthorID     *uint64
	Visibility   []model.FeedVisibility
	OnlyApproved bool
}

// FeedPagedListOptions describes feed query filters (page-based for admin).
type FeedPagedListOptions struct {
	Page             int
	PageSize         int
	AuthorID         *uint64
	CategoryID       *uint64
	Keyword          string
	ModerationStatus *model.FeedModerationStatus
	Visibility       *model.FeedVisibility
	DateFrom         *time.Time
	DateTo           *time.Time
}

// FeedReportListOptions describes feed report query filters.
type FeedReportListOptions struct {
	Page       int
	PageSize   int
	FeedID     *uint64
	ReporterID *uint64
	Status     *string
	DateFrom   *time.Time
	DateTo     *time.Time
}

// NotificationListOptions describes notification queries.
type NotificationListOptions struct {
	Page     int
	PageSize int
	UserID   uint64
	Unread   *bool
	Priority []model.NotificationPriority
}

// PaymentListOptions contains filtering options for payment queries.
type PaymentListOptions struct {
	Page               int
	PageSize           int
	OrderID            *uint64
	UserID             *uint64
	Method             *model.PaymentMethod
	Methods            []model.PaymentMethod
	Status             *model.PaymentStatus
	Statuses           []model.PaymentStatus
	DateFrom           *time.Time
	DateTo             *time.Time
	CollectionEntityID *uint64 // 收款主体ID筛选
	MerchantNo         string  // 商户号筛选
	ProviderTradeNo    string  // 第三方交易号筛选
	MinAmountCents     *int64  // 最小金额筛选
	MaxAmountCents     *int64  // 最大金额筛选
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

// ReviewReportListOptions contains filtering options for review report queries.
type ReviewReportListOptions struct {
	Page       int
	PageSize   int
	ReviewID   *uint64
	ReporterID *uint64
	Status     *model.ReviewReportStatus
	DateFrom   *time.Time
	DateTo     *time.Time
}

// DisputeListOptions contains filtering options for dispute queries.
type DisputeListOptions struct {
	Page             int
	PageSize         int
	UserID           *uint64
	OrderID          *uint64
	AssignedToUserID *uint64
	Statuses         []model.DisputeStatus
	SLABreached      *bool
	Keyword          string
	OrderNo          string // 订单号筛选
	DateFrom         *time.Time
	DateTo           *time.Time
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

// OperationLogSearchOptions contains filtering options for general operation log search.
type OperationLogSearchOptions struct {
	Page        int
	PageSize    int
	EntityType  string
	EntityID    *uint64
	Action      string
	ActorUserID *uint64
	DateFrom    *time.Time
	DateTo      *time.Time
}

// ChatGroupListOptions defines filters for querying chat groups by user.
type ChatGroupListOptions struct {
	Page            int
	PageSize        int
	GroupType       *model.ChatGroupType
	IncludeInactive bool
	Keyword         string
	RelatedOrderID  *uint64
}

// ChatGroupMemberListOptions defines pagination for group member listing.
type ChatGroupMemberListOptions struct {
	Page     int
	PageSize int
	Role     string
	Keyword  string
}

// ChatMessageListOptions defines filters for listing chat messages.
type ChatMessageListOptions struct {
	Page          int
	PageSize      int
	GroupID       uint64
	BeforeID      *uint64
	AfterID       *uint64
	DateFrom      *time.Time
	DateTo        *time.Time
	MessageType   *model.ChatMessageType
	AuditStatuses []model.ChatMessageAuditStatus
}

// ChatMessageModerationListOptions defines filters for moderation queue.
type ChatMessageModerationListOptions struct {
	Page        int
	PageSize    int
	GroupID     *uint64
	SenderID    *uint64
	AuditStatus *model.ChatMessageAuditStatus
	DateFrom    *time.Time
	DateTo      *time.Time
}

// ChatReportListOptions defines filters for querying chat reports.
type ChatReportListOptions struct {
	Page       int
	PageSize   int
	Status     string
	ReporterID *uint64
	MessageID  *uint64
	DateFrom   *time.Time
	DateTo     *time.Time
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

// WithdrawRepository 提现记录仓储接口
type WithdrawRepository interface {
	Create(ctx context.Context, withdraw *model.Withdraw) error
	Get(ctx context.Context, id uint64) (*model.Withdraw, error)
	Update(ctx context.Context, withdraw *model.Withdraw) error
	List(ctx context.Context, opts any) ([]model.Withdraw, int64, error)
	GetPlayerBalance(ctx context.Context, playerID uint64) (any, error)
}

// ServiceItemRepository 服务项目仓储接口
type ServiceItemRepository interface {
	Create(ctx context.Context, item *model.ServiceItem) error
	Get(ctx context.Context, id uint64) (*model.ServiceItem, error)
	GetByCode(ctx context.Context, itemCode string) (*model.ServiceItem, error)
	List(ctx context.Context, opts ServiceItemListOptions) ([]model.ServiceItem, int64, error)
	Update(ctx context.Context, item *model.ServiceItem) error
	Delete(ctx context.Context, id uint64) error
	BatchDelete(ctx context.Context, ids []uint64) (int64, error)
	BatchUpdateStatus(ctx context.Context, ids []uint64, isActive bool) error
	BatchUpdatePrice(ctx context.Context, ids []uint64, basePriceCents int64) error
	BatchUpdateCommission(ctx context.Context, ids []uint64, commissionRate float64) error
	GetGifts(ctx context.Context, page, pageSize int) ([]model.ServiceItem, int64, error)
	GetGameServices(ctx context.Context, gameID uint64, subCategory *model.ServiceItemSubCategory) ([]model.ServiceItem, error)
}

// CommissionRepository 抽成记录仓储接口
type CommissionRepository interface {
	// 抽成规则
	CreateRule(ctx context.Context, rule *model.CommissionRule) error
	GetRule(ctx context.Context, id uint64) (*model.CommissionRule, error)
	GetDefaultRule(ctx context.Context) (*model.CommissionRule, error)
	GetRuleForOrder(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error)
	ListRules(ctx context.Context, opts any) ([]model.CommissionRule, int64, error)
	UpdateRule(ctx context.Context, rule *model.CommissionRule) error
	DeleteRule(ctx context.Context, id uint64) error
	// 抽成记录
	CreateRecord(ctx context.Context, record *model.CommissionRecord) error
	GetRecord(ctx context.Context, id uint64) (*model.CommissionRecord, error)
	GetRecordByOrderID(ctx context.Context, orderID uint64) (*model.CommissionRecord, error)
	ListRecords(ctx context.Context, opts any) ([]model.CommissionRecord, int64, error)
	UpdateRecord(ctx context.Context, record *model.CommissionRecord) error
	// 月度结算
	CreateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error
	GetSettlement(ctx context.Context, id uint64) (*model.MonthlySettlement, error)
	GetSettlementByPlayerMonth(ctx context.Context, playerID uint64, month string) (*model.MonthlySettlement, error)
	ListSettlements(ctx context.Context, opts any) ([]model.MonthlySettlement, int64, error)
	UpdateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error
	// 统计查询
	GetMonthlyStats(ctx context.Context, month string) (any, error)
	GetPlayerMonthlyIncome(ctx context.Context, playerID uint64, month string) (int64, error)
}

// RankingCommissionRepository 排名抽成配置仓储
type RankingCommissionRepository interface {
	CreateConfig(ctx context.Context, config *model.RankingCommissionConfig) error
	GetConfig(ctx context.Context, id uint64) (*model.RankingCommissionConfig, error)
	GetActiveConfigForMonth(ctx context.Context, rankingType model.RankingType, month string) (*model.RankingCommissionConfig, error)
	ListConfigs(ctx context.Context, opts any) ([]model.RankingCommissionConfig, int64, error)
	UpdateConfig(ctx context.Context, config *model.RankingCommissionConfig) error
	DeleteConfig(ctx context.Context, id uint64) error
}

// ServiceItemListOptions 服务项列表查询选项
type ServiceItemListOptions struct {
	Page        int
	PageSize    int
	GameID      *uint64
	PlayerID    *uint64
	Category    *string                       // 使用指针表示可选
	SubCategory *model.ServiceItemSubCategory // 使用指针表示可选
	IsActive    *bool
}

// CommissionRuleListOptions 抽成规则列表选项
type CommissionRuleListOptions struct {
	GameID      *uint64
	PlayerID    *uint64
	ServiceType string
	IsActive    *bool
	Page        int
	PageSize    int
}

// CommissionRecordListOptions 抽成记录列表选项
type CommissionRecordListOptions struct {
	PlayerID        *uint64
	StartTime       *time.Time
	EndTime         *time.Time
	SettlementMonth *string
	Status          string
	Page            int
	PageSize        int
}

// SettlementListOptions 结算列表选项
type SettlementListOptions struct {
	PlayerID        *uint64
	SettlementMonth *string
	Status          *string
	Page            int
	PageSize        int
}

// MonthlyStats 月度统计
type MonthlyStats struct {
	Month                  string
	TotalOrders            int64
	TotalRevenueCents      int64
	TotalCommissionCents   int64
	TotalPlayerIncomeCents int64
}

// UserTagRepository 用户标签仓储接口
// 错误约定：当资源不存在时返回 repository.ErrNotFound
type UserTagRepository interface {
	// 标签管理
	CreateTag(ctx context.Context, tag *model.UserTag) error
	GetTag(ctx context.Context, id uint64) (*model.UserTag, error)
	ListTags(ctx context.Context) ([]model.UserTag, error)
	UpdateTag(ctx context.Context, tag *model.UserTag) error
	DeleteTag(ctx context.Context, id uint64) error

	// 用户标签操作
	AddTagToUser(ctx context.Context, userID uint64, tagID uint64) error
	RemoveTagFromUser(ctx context.Context, userID uint64, tagID uint64) error
	GetUserTags(ctx context.Context, userID uint64) ([]model.UserTag, error)
	BatchSetUserTags(ctx context.Context, userID uint64, tagIDs []uint64) error
	GetUsersByTag(ctx context.Context, tagID uint64, page, pageSize int) ([]model.User, int64, error)
}

// UserLoginHistoryRepository 登录历史仓储接口
// 错误约定：当资源不存在时返回 repository.ErrNotFound
type UserLoginHistoryRepository interface {
	Create(ctx context.Context, history *model.UserLoginHistory) error
	GetByUserID(ctx context.Context, userID uint64, page, pageSize int) ([]model.UserLoginHistory, int64, error)
	GetByUserIDAndDate(ctx context.Context, userID uint64, dateFrom, dateTo time.Time) ([]model.UserLoginHistory, error)
}

// UserBehaviorRepository 用户行为仓储接口
// 错误约定：当资源不存在时返回 repository.ErrNotFound
type UserBehaviorRepository interface {
	Create(ctx context.Context, behavior *model.UserBehavior) error
	GetUserBehaviors(ctx context.Context, userID uint64, page, pageSize int) ([]model.UserBehavior, int64, error)
	GetBehaviorStats(ctx context.Context, days int) (map[string]int64, error)
}

// SensitiveWordRepository 敏感词仓储接口
// 错误约定：当资源不存在时返回 repository.ErrNotFound
type SensitiveWordRepository interface {
	Create(ctx context.Context, word *model.SensitiveWord) error
	Get(ctx context.Context, id uint64) (*model.SensitiveWord, error)
	List(ctx context.Context, opts SensitiveWordListOptions) ([]model.SensitiveWord, int64, error)
	Update(ctx context.Context, word *model.SensitiveWord) error
	Delete(ctx context.Context, id uint64) error
	GetAll(ctx context.Context) ([]model.SensitiveWord, error)
}

// SensitiveWordListOptions 敏感词列表查询选项
type SensitiveWordListOptions struct {
	Page      int
	PageSize  int
	Keyword   string
	Category  *model.SensitiveWordCategory
	Severity  *model.SensitiveWordSeverity
	MatchType *model.SensitiveWordMatchType
	IsActive  *bool
}

// ReviewStats 评价统计数据
type ReviewStats struct {
	TotalReviews       int64         `json:"totalReviews"`
	AverageRating      float64       `json:"averageRating"`
	RatingDistribution map[int]int64 `json:"ratingDistribution"` // 1-5分的分布
}

// PlayerReviewStats 陪玩师评价统计
type PlayerReviewStats struct {
	PlayerID      uint64  `json:"playerId"`
	PlayerName    string  `json:"playerName"`
	ReviewCount   int64   `json:"reviewCount"`
	AverageRating float64 `json:"averageRating"`
}

// GameReviewStats 游戏评价统计
type GameReviewStats struct {
	GameID        uint64  `json:"gameId"`
	GameName      string  `json:"gameName"`
	ReviewCount   int64   `json:"reviewCount"`
	AverageRating float64 `json:"averageRating"`
}

// ReviewDisplaySettingsRepository 评价展示设置仓储接口
// 错误约定：当资源不存在时返回 repository.ErrNotFound
type ReviewDisplaySettingsRepository interface {
	// Get 获取当前设置（单例模式，只有一条记录）
	Get(ctx context.Context) (*model.ReviewDisplaySettings, error)
	// Save 保存设置（创建或更新）
	Save(ctx context.Context, settings *model.ReviewDisplaySettings) error
}

// ContentCategoryRepository 内容分类仓储接口
// 错误约定：当资源不存在时返回 repository.ErrNotFound
type ContentCategoryRepository interface {
	Create(ctx context.Context, category *model.ContentCategory) error
	Get(ctx context.Context, id uint64) (*model.ContentCategory, error)
	GetByName(ctx context.Context, name string) (*model.ContentCategory, error)
	List(ctx context.Context, opts ContentCategoryListOptions) ([]model.ContentCategory, int64, error)
	Update(ctx context.Context, category *model.ContentCategory) error
	Delete(ctx context.Context, id uint64) error
	// GetFeedCount 获取分类下的动态数量
	GetFeedCount(ctx context.Context, categoryID uint64) (int64, error)
	// MigrateFeeds 将分类下的动态迁移到另一个分类
	MigrateFeeds(ctx context.Context, fromCategoryID, toCategoryID uint64) error
}

// ContentCategoryListOptions 内容分类列表查询选项
type ContentCategoryListOptions struct {
	Page     int
	PageSize int
	Keyword  string
	Status   *model.ContentCategoryStatus
}

// ============================================================================
// 陪玩师等级/认证模块接口
// ============================================================================

// GameRankRepository 游戏段位配置仓储接口
// 错误约定：当资源不存在时返回 repository.ErrNotFound
type GameRankRepository interface {
	Create(ctx context.Context, rank *model.GameRank) error
	Get(ctx context.Context, id uint64) (*model.GameRank, error)
	GetWithGame(ctx context.Context, id uint64) (*model.GameRank, error)
	List(ctx context.Context) ([]model.GameRank, error)
	ListByGameID(ctx context.Context, gameID uint64) ([]model.GameRank, error)
	ListPaged(ctx context.Context, opts GameRankListOptions) ([]model.GameRank, int64, error)
	Update(ctx context.Context, rank *model.GameRank) error
	Delete(ctx context.Context, id uint64) error
	BatchDelete(ctx context.Context, ids []uint64) (int64, error)
	BatchUpdateStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error)
}

// GameRankListOptions 游戏段位列表查询选项
type GameRankListOptions struct {
	Page     int
	PageSize int
	GameID   *uint64
	Keyword  string
	IsActive *bool
}

// PlayerRankRepository 陪玩师段位认证仓储接口
// 错误约定：当资源不存在时返回 repository.ErrNotFound
type PlayerRankRepository interface {
	Create(ctx context.Context, record *model.PlayerRankRecord) error
	Get(ctx context.Context, id uint64) (*model.PlayerRankRecord, error)
	GetWithRelations(ctx context.Context, id uint64) (*model.PlayerRankRecord, error)
	GetByPlayerAndGame(ctx context.Context, playerID, gameID uint64) (*model.PlayerRankRecord, error)
	ListByPlayerID(ctx context.Context, playerID uint64) ([]model.PlayerRankRecord, error)
	ListPaged(ctx context.Context, opts PlayerRankListOptions) ([]model.PlayerRankRecord, int64, error)
	ListPending(ctx context.Context, page, pageSize int) ([]model.PlayerRankRecord, int64, error)
	Update(ctx context.Context, record *model.PlayerRankRecord) error
	UpdateStatus(ctx context.Context, id uint64, status model.PlayerRankStatus, verifiedBy *uint64, rejectReason string) error
	Delete(ctx context.Context, id uint64) error
	CountByStatus(ctx context.Context) (map[model.PlayerRankStatus]int64, error)
	GetPendingCount(ctx context.Context) (int64, error)
}

// PlayerRankListOptions 陪玩师段位认证列表查询选项
type PlayerRankListOptions struct {
	Page     int
	PageSize int
	PlayerID *uint64
	GameID   *uint64
	Status   *model.PlayerRankStatus
	Statuses []model.PlayerRankStatus
}

// PlayerCertificationRepository 陪玩师实名认证仓储接口
// 错误约定：当资源不存在时返回 repository.ErrNotFound
type PlayerCertificationRepository interface {
	Create(ctx context.Context, cert *model.PlayerCertification) error
	Get(ctx context.Context, id uint64) (*model.PlayerCertification, error)
	GetWithPlayer(ctx context.Context, id uint64) (*model.PlayerCertification, error)
	GetByPlayerID(ctx context.Context, playerID uint64) (*model.PlayerCertification, error)
	ListPaged(ctx context.Context, opts PlayerCertificationListOptions) ([]model.PlayerCertification, int64, error)
	ListPending(ctx context.Context, page, pageSize int) ([]model.PlayerCertification, int64, error)
	Update(ctx context.Context, cert *model.PlayerCertification) error
	UpdateStatus(ctx context.Context, id uint64, status model.CertificationStatus, verifiedBy *uint64, rejectReason string) error
	Delete(ctx context.Context, id uint64) error
	CountByStatus(ctx context.Context) (map[model.CertificationStatus]int64, error)
	GetPendingCount(ctx context.Context) (int64, error)
}

// PlayerCertificationListOptions 陪玩师实名认证列表查询选项
type PlayerCertificationListOptions struct {
	Page     int
	PageSize int
	PlayerID *uint64
	Status   *model.CertificationStatus
	Statuses []model.CertificationStatus
	Keyword  string
}

// ============================================================================
// 订单超时处理模块接口
// ============================================================================

// OrderTimeoutRepository 订单超时仓储接口
// 错误约定：当资源不存在时返回 repository.ErrNotFound
type OrderTimeoutRepository interface {
	// 配置管理
	GetConfig(ctx context.Context, key string) (*model.OrderTimeoutConfig, error)
	ListConfigs(ctx context.Context) ([]model.OrderTimeoutConfig, error)
	SaveConfig(ctx context.Context, config *model.OrderTimeoutConfig) error
	DeleteConfig(ctx context.Context, key string) error

	// 超时日志
	CreateLog(ctx context.Context, log *model.OrderTimeoutLog) error
	GetLog(ctx context.Context, id uint64) (*model.OrderTimeoutLog, error)
	GetLogWithOrder(ctx context.Context, id uint64) (*model.OrderTimeoutLog, error)
	ListLogsByOrderID(ctx context.Context, orderID uint64) ([]model.OrderTimeoutLog, error)
	ListLogsPaged(ctx context.Context, opts OrderTimeoutLogListOptions) ([]model.OrderTimeoutLog, int64, error)
	GetLogStats(ctx context.Context) (map[model.OrderTimeoutType]int64, error)

	// 客服分配
	CreateAssignment(ctx context.Context, assignment *model.OrderServiceAssignment) error
	GetAssignment(ctx context.Context, id uint64) (*model.OrderServiceAssignment, error)
	GetAssignmentWithRelations(ctx context.Context, id uint64) (*model.OrderServiceAssignment, error)
	GetAssignmentByOrderID(ctx context.Context, orderID uint64) (*model.OrderServiceAssignment, error)
	ListAssignmentsByServiceUser(ctx context.Context, serviceUserID uint64, status *model.ServiceAssignmentStatus) ([]model.OrderServiceAssignment, error)
	ListAssignmentsPaged(ctx context.Context, opts ServiceAssignmentListOptions) ([]model.OrderServiceAssignment, int64, error)
	UpdateAssignment(ctx context.Context, assignment *model.OrderServiceAssignment) error
	UpdateAssignmentStatus(ctx context.Context, id uint64, status model.ServiceAssignmentStatus) error
	DeleteAssignment(ctx context.Context, id uint64) error
	GetAssignmentStats(ctx context.Context) (map[model.ServiceAssignmentStatus]int64, error)
	GetActiveAssignmentCount(ctx context.Context, serviceUserID uint64) (int64, error)
}

// OrderTimeoutLogListOptions 订单超时日志列表查询选项
type OrderTimeoutLogListOptions struct {
	Page        int
	PageSize    int
	OrderID     *uint64
	TimeoutType *model.OrderTimeoutType
	Action      *model.OrderTimeoutAction
	DateFrom    *time.Time
	DateTo      *time.Time
}

// ServiceAssignmentListOptions 客服分配列表查询选项
type ServiceAssignmentListOptions struct {
	Page          int
	PageSize      int
	OrderID       *uint64
	ServiceUserID *uint64
	Status        *model.ServiceAssignmentStatus
	AssignType    string
	DateFrom      *time.Time
	DateTo        *time.Time
}

// ============================================================================
// 用户拉黑模块接口
// ============================================================================

// UserBlockRepository 用户拉黑仓储接口
// 错误约定：当资源不存在时返回 repository.ErrNotFound
type UserBlockRepository interface {
	Create(ctx context.Context, block *model.UserBlock) error
	Get(ctx context.Context, id uint64) (*model.UserBlock, error)
	GetWithRelations(ctx context.Context, id uint64) (*model.UserBlock, error)
	GetByBlockerAndBlocked(ctx context.Context, blockerID, blockedID uint64) (*model.UserBlock, error)
	GetActiveByBlockerAndBlocked(ctx context.Context, blockerID, blockedID uint64) (*model.UserBlock, error)
	IsBlocked(ctx context.Context, userID1, userID2 uint64) (bool, error)
	IsBlockedBy(ctx context.Context, blockerID, blockedID uint64) (bool, error)
	ListByBlockerID(ctx context.Context, blockerID uint64, status *model.BlockStatus) ([]model.UserBlock, error)
	ListByBlockedID(ctx context.Context, blockedID uint64, status *model.BlockStatus) ([]model.UserBlock, error)
	ListPaged(ctx context.Context, opts UserBlockListOptions) ([]model.UserBlock, int64, error)
	Update(ctx context.Context, block *model.UserBlock) error
	UpdateStatus(ctx context.Context, id uint64, status model.BlockStatus, canceledBy *uint64, adminRemark string) error
	Delete(ctx context.Context, id uint64) error
	GetBlockedUserIDs(ctx context.Context, blockerID uint64) ([]uint64, error)
	GetBlockerUserIDs(ctx context.Context, blockedID uint64) ([]uint64, error)
	GetAllBlockRelatedUserIDs(ctx context.Context, userID uint64) ([]uint64, error)
	CountByStatus(ctx context.Context) (map[model.BlockStatus]int64, error)
	GetActiveCount(ctx context.Context) (int64, error)
}

// UserBlockListOptions 用户拉黑列表查询选项
type UserBlockListOptions struct {
	Page        int
	PageSize    int
	BlockerID   *uint64
	BlockedID   *uint64
	BlockerType *model.BlockUserType
	BlockedType *model.BlockUserType
	Status      *model.BlockStatus
	DateFrom    *time.Time
	DateTo      *time.Time
}

// ============================================================================
// VIP会员模块接口
// ============================================================================

// VipRepository VIP仓储接口
// 错误约定：当资源不存在时返回 repository.ErrNotFound
type VipRepository interface {
	// VIP等级管理
	CreateLevel(ctx context.Context, level *model.VipLevel) error
	GetLevel(ctx context.Context, id uint64) (*model.VipLevel, error)
	GetLevelBySlug(ctx context.Context, slug string) (*model.VipLevel, error)
	GetDefaultLevel(ctx context.Context) (*model.VipLevel, error)
	ListLevels(ctx context.Context) ([]model.VipLevel, error)
	ListActiveLevels(ctx context.Context) ([]model.VipLevel, error)
	ListLevelsPaged(ctx context.Context, opts VipLevelListOptions) ([]model.VipLevel, int64, error)
	UpdateLevel(ctx context.Context, level *model.VipLevel) error
	DeleteLevel(ctx context.Context, id uint64) error
	SetDefaultLevel(ctx context.Context, id uint64) error
	GetLevelByExp(ctx context.Context, exp int64) (*model.VipLevel, error)
	BatchUpdateLevelStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error)
	BatchDeleteLevels(ctx context.Context, ids []uint64) (int64, error)

	// VIP配置管理
	GetConfig(ctx context.Context, key string) (*model.VipConfig, error)
	ListConfigs(ctx context.Context) ([]model.VipConfig, error)
	SaveConfig(ctx context.Context, config *model.VipConfig) error
	DeleteConfig(ctx context.Context, key string) error
}

// VipLevelListOptions VIP等级列表查询选项
type VipLevelListOptions struct {
	Page     int
	PageSize int
	Keyword  string
	IsActive *bool
}

// GameCategoryRepository 游戏分类仓储接口
// 错误约定：当资源不存在时返回 repository.ErrNotFound
type GameCategoryRepository interface {
	// Create 创建分类
	Create(ctx context.Context, category *model.GameCategory) error
	// Get 获取分类详情
	Get(ctx context.Context, id uint64) (*model.GameCategory, error)
	// List 获取分类列表（支持分页、筛选）
	List(ctx context.Context, opts GameCategoryListOptions) ([]*model.GameCategory, int64, error)
	// Update 更新分类
	Update(ctx context.Context, category *model.GameCategory) error
	// Delete 删除分类（软删除）
	Delete(ctx context.Context, id uint64) error
	// GetByName 根据名称获取分类
	GetByName(ctx context.Context, name string) (*model.GameCategory, error)
	// BatchUpdateStatus 批量更新状态
	BatchUpdateStatus(ctx context.Context, ids []uint64, isActive bool) error
	// BatchDelete 批量删除
	BatchDelete(ctx context.Context, ids []uint64) error
	// Exists 检查分类是否存在
	Exists(ctx context.Context, id uint64) (bool, error)
	// CountGames 统计分类下的游戏数量
	CountGames(ctx context.Context, categoryID uint64) (int64, error)
	// CountServiceItems 统计分类下的服务项目数量
	CountServiceItems(ctx context.Context, categoryID uint64) (int64, error)
}

// GameCategoryListOptions 游戏分类列表查询选项
type GameCategoryListOptions struct {
	Page     int
	PageSize int
	IsActive *bool
	Keyword  string
}
