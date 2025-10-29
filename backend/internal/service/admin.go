package service

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sort"
	"strings"
	"time"

	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"gamelink/internal/cache"
	"gamelink/internal/logging"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/gormrepo"
)

var (
	// ErrValidation 表示输入校验失败。
	ErrValidation = errors.New("validation failed")
	// ErrUserNotFound 用于统一标识用户不存在的场景。
	ErrUserNotFound = errors.New("user not found")
	// ErrOrderInvalidTransition 代表订单状态流转不合法。
	ErrOrderInvalidTransition = errors.New("invalid order status transition")

	// ErrNotFound 暴露仓储的未找到错误，便于 handler 判定。
	ErrNotFound = repository.ErrNotFound
)

// AdminService 聚合后台管理所需的业务逻辑。
type AdminService struct {
	games    repository.GameRepository
	users    repository.UserRepository
	players  repository.PlayerRepository
	orders   repository.OrderRepository
	payments repository.PaymentRepository
	cache    cache.Cache
	tx       TxManager
}

const (
	cacheKeyGames    = "admin:games"
	cacheKeyUsers    = "admin:users"
	cacheKeyPlayers  = "admin:players"
	cacheKeyOrders   = "admin:orders"
	cacheKeyPayments = "admin:payments"
)

var listCacheTTL = readListCacheTTL()

func readListCacheTTL() time.Duration {
	if v := strings.TrimSpace(os.Getenv("ADMIN_LIST_TTL")); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return 30 * time.Second
}

// NewAdminService 创建服务实例。
func NewAdminService(
	games repository.GameRepository,
	users repository.UserRepository,
	players repository.PlayerRepository,
	orders repository.OrderRepository,
	payments repository.PaymentRepository,
	cache cache.Cache,
) *AdminService {
	return &AdminService{
		games:    games,
		users:    users,
		players:  players,
		orders:   orders,
		payments: payments,
		cache:    cache,
	}
}

// TxManager abstracts UnitOfWork for transactional operations.
type TxManager interface {
	WithTx(ctx context.Context, fn func(r *gormrepo.Repos) error) error
}

// SetTxManager injects a transaction manager.
func (s *AdminService) SetTxManager(tx TxManager) { s.tx = tx }

// UpdatePlayerSkillTags 替换玩家技能标签集合（需要 TxManager）。
func (s *AdminService) UpdatePlayerSkillTags(ctx context.Context, playerID uint64, tags []string) error {
	if s.tx == nil {
		return errors.New("transaction manager not configured")
	}
	err := s.tx.WithTx(ctx, func(r *gormrepo.Repos) error {
		// ensure player exists
		if _, err := r.Players.Get(ctx, playerID); err != nil {
			return err
		}
		return r.Tags.ReplaceTags(ctx, playerID, tags)
	})
	if err == nil {
		s.appendLogAsync(ctx, string(model.OpEntityPlayer), playerID, string(model.OpActionUpdate), map[string]any{"tags_count": len(tags)})
	}
	return err
}

// RegisterUserAndPlayer creates a user and a player profile in a single transaction.
func (s *AdminService) RegisterUserAndPlayer(ctx context.Context, u CreateUserInput, p CreatePlayerInput) (*model.User, *model.Player, error) {
	if s.tx == nil {
		return nil, nil, errors.New("transaction manager not configured")
	}
	// basic validations reuse existing ones
	if err := validateUserInput(u.Name, u.Role, u.Status, u.Password); err != nil {
		return nil, nil, err
	}
	// For registration flow, user will be created first; only verify player fields except UserID.
	if p.VerificationStatus == "" {
		return nil, nil, ErrValidation
	}

	var createdUser *model.User
	var createdPlayer *model.Player

	err := s.tx.WithTx(ctx, func(r *gormrepo.Repos) error {
		// hash password
		hashed, err := hashPassword(u.Password)
		if err != nil {
			return err
		}
		user := &model.User{
			Phone:        strings.TrimSpace(u.Phone),
			Email:        strings.TrimSpace(u.Email),
			PasswordHash: hashed,
			Name:         strings.TrimSpace(u.Name),
			AvatarURL:    strings.TrimSpace(u.AvatarURL),
			Role:         u.Role,
			Status:       u.Status,
		}
		if err := r.Users.Create(ctx, user); err != nil {
			return err
		}
		createdUser = user

		player := &model.Player{
			UserID:             user.ID,
			Nickname:           strings.TrimSpace(p.Nickname),
			Bio:                strings.TrimSpace(p.Bio),
			HourlyRateCents:    p.HourlyRateCents,
			MainGameID:         p.MainGameID,
			VerificationStatus: p.VerificationStatus,
		}
		if err := r.Players.Create(ctx, player); err != nil {
			return err
		}
		createdPlayer = player
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	// invalidate relevant caches
	s.invalidateCache(ctx, cacheKeyUsers, cacheKeyPlayers)
	return createdUser, createdPlayer, nil
}

// --- Game management ---

// CreateGameInput 创建游戏时使用的参数。
type CreateGameInput struct {
	Key         string
	Name        string
	Category    string
	IconURL     string
	Description string
}

// UpdateGameInput 修改游戏资料。
type UpdateGameInput struct {
	Key         string
	Name        string
	Category    string
	IconURL     string
	Description string
}

// ListGames 返回全部游戏。
func (s *AdminService) ListGames(ctx context.Context) ([]model.Game, error) {
	return getCachedList(ctx, s.cache, cacheKeyGames, listCacheTTL, func() ([]model.Game, error) {
		return s.games.List(ctx)
	})
}

// ListGamesPaged 返回分页游戏列表。
func (s *AdminService) ListGamesPaged(ctx context.Context, page, pageSize int) ([]model.Game, *model.Pagination, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	items, total, err := s.games.ListPaged(ctx, page, pageSize)
	if err != nil {
		return nil, nil, err
	}
	p := buildPagination(page, pageSize, total)
	return items, &p, nil
}

// GetGame 获取单个游戏详情。
func (s *AdminService) GetGame(ctx context.Context, id uint64) (*model.Game, error) {
	return s.games.Get(ctx, id)
}

// CreateGame 创建游戏。
func (s *AdminService) CreateGame(ctx context.Context, input CreateGameInput) (*model.Game, error) {
	if err := validateGameInput(input.Key, input.Name); err != nil {
		return nil, err
	}

	game := &model.Game{
		Key:         strings.TrimSpace(input.Key),
		Name:        strings.TrimSpace(input.Name),
		Category:    strings.TrimSpace(input.Category),
		IconURL:     strings.TrimSpace(input.IconURL),
		Description: strings.TrimSpace(input.Description),
	}

	if err := s.games.Create(ctx, game); err != nil {
		return nil, err
	}

	s.invalidateCache(ctx, cacheKeyGames)
	// audit
	s.appendLogAsync(ctx, string(model.OpEntityGame), game.ID, string(model.OpActionCreate), map[string]any{"key": game.Key})

	return game, nil
}

// UpdateGame 更新游戏。
func (s *AdminService) UpdateGame(ctx context.Context, id uint64, input UpdateGameInput) (*model.Game, error) {
	if err := validateGameInput(input.Key, input.Name); err != nil {
		return nil, err
	}

	game, err := s.games.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	game.Key = strings.TrimSpace(input.Key)
	game.Name = strings.TrimSpace(input.Name)
	game.Category = strings.TrimSpace(input.Category)
	game.IconURL = strings.TrimSpace(input.IconURL)
	game.Description = strings.TrimSpace(input.Description)

	if err := s.games.Update(ctx, game); err != nil {
		return nil, err
	}

	s.invalidateCache(ctx, cacheKeyGames)
	// audit
	s.appendLogAsync(ctx, string(model.OpEntityGame), game.ID, string(model.OpActionUpdate), map[string]any{"key": game.Key})

	return game, nil
}

// DeleteGame 删除游戏。
func (s *AdminService) DeleteGame(ctx context.Context, id uint64) error {
	if err := s.games.Delete(ctx, id); err != nil {
		return err
	}
	s.invalidateCache(ctx, cacheKeyGames)
	// audit
	s.appendLogAsync(ctx, string(model.OpEntityGame), id, string(model.OpActionDelete), nil)
	return nil
}

func validateGameInput(key, name string) error {
	if strings.TrimSpace(key) == "" || strings.TrimSpace(name) == "" {
		return ErrValidation
	}
	return nil
}

// --- User management ---

// CreateUserInput 定义创建用户的请求。
type CreateUserInput struct {
	Phone     string
	Email     string
	Password  string
	Name      string
	AvatarURL string
	Role      model.Role
	Status    model.UserStatus
}

// UpdateUserInput 定义更新用户资料的请求。
type UpdateUserInput struct {
	Phone     string
	Email     string
	Name      string
	AvatarURL string
	Role      model.Role
	Status    model.UserStatus
	Password  *string
}

// ListUsers 返回全部用户。
func (s *AdminService) ListUsers(ctx context.Context) ([]model.User, error) {
	return getCachedList(ctx, s.cache, cacheKeyUsers, listCacheTTL, func() ([]model.User, error) {
		return s.users.List(ctx)
	})
}

// ListUsersPaged 返回分页用户列表。
func (s *AdminService) ListUsersPaged(ctx context.Context, page, pageSize int) ([]model.User, *model.Pagination, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	items, total, err := s.users.ListPaged(ctx, page, pageSize)
	if err != nil {
		return nil, nil, err
	}
	p := buildPagination(page, pageSize, total)
	return items, &p, nil
}

// ListUsersWithOptions 返回带筛选的分页用户列表。
func (s *AdminService) ListUsersWithOptions(ctx context.Context, opts repository.UserListOptions) ([]model.User, *model.Pagination, error) {
	normalized := opts
	normalized.Page = repository.NormalizePage(opts.Page)
	normalized.PageSize = repository.NormalizePageSize(opts.PageSize)
	items, total, err := s.users.ListWithFilters(ctx, normalized)
	if err != nil {
		return nil, nil, err
	}
	p := buildPagination(normalized.Page, normalized.PageSize, total)
	return items, &p, nil
}

// GetUser 返回指定用户。
func (s *AdminService) GetUser(ctx context.Context, id uint64) (*model.User, error) {
	user, err := s.users.Get(ctx, id)
	if err != nil {
		return nil, mapUserError(err)
	}
	return user, nil
}

// CreateUser 新建用户并对密码加密。
func (s *AdminService) CreateUser(ctx context.Context, input CreateUserInput) (*model.User, error) {
	if err := validateUserInput(input.Name, input.Role, input.Status, input.Password); err != nil {
		return nil, err
	}

	hashed, err := hashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Phone:        strings.TrimSpace(input.Phone),
		Email:        strings.TrimSpace(input.Email),
		PasswordHash: hashed,
		Name:         strings.TrimSpace(input.Name),
		AvatarURL:    strings.TrimSpace(input.AvatarURL),
		Role:         input.Role,
		Status:       input.Status,
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, cacheKeyUsers)
	// audit
	s.appendLogAsync(ctx, string(model.OpEntityUser), user.ID, string(model.OpActionCreate), map[string]any{"role": user.Role, "status": user.Status})
	if rid, ok := logging.RequestIDFromContext(ctx); ok {
		slog.Info("user_created", slog.Uint64("user_id", user.ID), slog.String("role", string(user.Role)), slog.String("status", string(user.Status)), slog.String("request_id", rid))
	} else {
		slog.Info("user_created", slog.Uint64("user_id", user.ID), slog.String("role", string(user.Role)), slog.String("status", string(user.Status)))
	}
	return user, nil
}

// UpdateUser 更新用户基础信息。
func (s *AdminService) UpdateUser(ctx context.Context, id uint64, input UpdateUserInput) (*model.User, error) {
	user, err := s.users.Get(ctx, id)
	if err != nil {
		return nil, mapUserError(err)
	}

	if err := validateUserInput(input.Name, input.Role, input.Status, optionalPassword(input.Password)); err != nil {
		return nil, err
	}

	// 避免将唯一字段更新为空字符串导致唯一索引冲突；空值保持原值
	if v := strings.TrimSpace(input.Phone); v != "" {
		user.Phone = v
	}
	if v := strings.TrimSpace(input.Email); v != "" {
		user.Email = v
	}
	user.Name = strings.TrimSpace(input.Name)
	user.AvatarURL = strings.TrimSpace(input.AvatarURL)
	user.Role = input.Role
	user.Status = input.Status

	if input.Password != nil && strings.TrimSpace(*input.Password) != "" {
		hash, err := hashPassword(*input.Password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = hash
	}

	if err := s.users.Update(ctx, user); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, cacheKeyUsers)
	// audit
	s.appendLogAsync(ctx, string(model.OpEntityUser), user.ID, string(model.OpActionUpdate), map[string]any{"role": user.Role, "status": user.Status})
	if rid, ok := logging.RequestIDFromContext(ctx); ok {
		slog.Info("user_updated", slog.Uint64("user_id", user.ID), slog.String("request_id", rid))
	} else {
		slog.Info("user_updated", slog.Uint64("user_id", user.ID))
	}
	return user, nil
}

// DeleteUser 删除用户。
func (s *AdminService) DeleteUser(ctx context.Context, id uint64) error {
	if err := s.users.Delete(ctx, id); err != nil {
		return mapUserError(err)
	}
	s.invalidateCache(ctx, cacheKeyUsers)
	// audit
	s.appendLogAsync(ctx, string(model.OpEntityUser), id, string(model.OpActionDelete), nil)
	if rid, ok := logging.RequestIDFromContext(ctx); ok {
		slog.Info("user_deleted", slog.Uint64("user_id", id), slog.String("request_id", rid))
	} else {
		slog.Info("user_deleted", slog.Uint64("user_id", id))
	}
	return nil
}

// UpdateUserStatus 单独更新用户状态并记录审计。
func (s *AdminService) UpdateUserStatus(ctx context.Context, id uint64, status model.UserStatus) (*model.User, error) {
	user, err := s.users.Get(ctx, id)
	if err != nil {
		return nil, mapUserError(err)
	}
	if err := validateUserInput(user.Name, user.Role, status, ""); err != nil {
		return nil, err
	}
	user.Status = status
	if err := s.users.Update(ctx, user); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, cacheKeyUsers)
	s.appendLogAsync(ctx, string(model.OpEntityUser), user.ID, string(model.OpActionUpdateStatus), map[string]any{"status": user.Status})
	return user, nil
}

// UpdateUserRole 单独更新用户角色并记录审计。
func (s *AdminService) UpdateUserRole(ctx context.Context, id uint64, role model.Role) (*model.User, error) {
	user, err := s.users.Get(ctx, id)
	if err != nil {
		return nil, mapUserError(err)
	}
	if err := validateUserInput(user.Name, role, user.Status, ""); err != nil {
		return nil, err
	}
	user.Role = role
	if err := s.users.Update(ctx, user); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, cacheKeyUsers)
	s.appendLogAsync(ctx, string(model.OpEntityUser), user.ID, string(model.OpActionUpdateRole), map[string]any{"role": user.Role})
	return user, nil
}

func validateUserInput(name string, role model.Role, status model.UserStatus, password string) error {
	if strings.TrimSpace(name) == "" {
		return ErrValidation
	}
	if role == "" || status == "" {
		return ErrValidation
	}
	if password != "" && !validPassword(password) {
		return ErrValidation
	}
	return nil
}

func validPassword(pw string) bool {
	if len(pw) < 6 {
		return false
	}
	hasLetter := false
	hasDigit := false
	for _, r := range pw {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			hasLetter = true
		}
		if r >= '0' && r <= '9' {
			hasDigit = true
		}
		if hasLetter && hasDigit {
			return true
		}
	}
	return false
}

func optionalPassword(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func hashPassword(raw string) (string, error) {
	if strings.TrimSpace(raw) == "" {
		return "", ErrValidation
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// --- Player management ---

// CreatePlayerInput 创建陪玩资料。
type CreatePlayerInput struct {
	UserID             uint64
	Nickname           string
	Bio                string
	Rank               string
	HourlyRateCents    int64
	MainGameID         uint64
	VerificationStatus model.VerificationStatus
}

// UpdatePlayerInput 更新陪玩资料。
type UpdatePlayerInput struct {
	Nickname           string
	Bio                string
	Rank               string
	HourlyRateCents    int64
	MainGameID         uint64
	VerificationStatus model.VerificationStatus
}

// ListPlayers 返回陪玩列表。
func (s *AdminService) ListPlayers(ctx context.Context) ([]model.Player, error) {
	return getCachedList(ctx, s.cache, cacheKeyPlayers, listCacheTTL, func() ([]model.Player, error) {
		return s.players.List(ctx)
	})
}

// ListPlayersPaged 返回分页陪玩列表。
func (s *AdminService) ListPlayersPaged(ctx context.Context, page, pageSize int) ([]model.Player, *model.Pagination, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	items, total, err := s.players.ListPaged(ctx, page, pageSize)
	if err != nil {
		return nil, nil, err
	}
	p := buildPagination(page, pageSize, total)
	return items, &p, nil
}

// GetPlayer 返回陪玩详情。
func (s *AdminService) GetPlayer(ctx context.Context, id uint64) (*model.Player, error) {
	return s.players.Get(ctx, id)
}

// CreatePlayer 新建陪玩档案。
func (s *AdminService) CreatePlayer(ctx context.Context, input CreatePlayerInput) (*model.Player, error) {
	if err := validatePlayerInput(input.UserID, input.VerificationStatus); err != nil {
		return nil, err
	}

	player := &model.Player{
		UserID:             input.UserID,
		Nickname:           strings.TrimSpace(input.Nickname),
		Bio:                strings.TrimSpace(input.Bio),
		HourlyRateCents:    input.HourlyRateCents,
		MainGameID:         input.MainGameID,
		VerificationStatus: input.VerificationStatus,
	}

	if err := s.players.Create(ctx, player); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, cacheKeyPlayers)
	// audit
	s.appendLogAsync(ctx, string(model.OpEntityPlayer), player.ID, string(model.OpActionCreate), map[string]any{"user_id": player.UserID})
	return player, nil
}

// UpdatePlayer 调整陪玩信息。
func (s *AdminService) UpdatePlayer(ctx context.Context, id uint64, input UpdatePlayerInput) (*model.Player, error) {
	player, err := s.players.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := validatePlayerInput(player.UserID, input.VerificationStatus); err != nil {
		return nil, err
	}

	player.Nickname = strings.TrimSpace(input.Nickname)
	player.Bio = strings.TrimSpace(input.Bio)
	player.HourlyRateCents = input.HourlyRateCents
	player.MainGameID = input.MainGameID
	player.VerificationStatus = input.VerificationStatus

	if err := s.players.Update(ctx, player); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, cacheKeyPlayers)
	// audit
	s.appendLogAsync(ctx, string(model.OpEntityPlayer), player.ID, string(model.OpActionUpdate), map[string]any{"main_game_id": player.MainGameID})
	return player, nil
}

// DeletePlayer 删除陪玩档案。
func (s *AdminService) DeletePlayer(ctx context.Context, id uint64) error {
	if err := s.players.Delete(ctx, id); err != nil {
		return err
	}
	s.invalidateCache(ctx, cacheKeyPlayers)
	// audit
	s.appendLogAsync(ctx, string(model.OpEntityPlayer), id, string(model.OpActionDelete), nil)
	return nil
}

func validatePlayerInput(userID uint64, verification model.VerificationStatus) error {
	if userID == 0 {
		return ErrValidation
	}
	if verification == "" {
		return ErrValidation
	}
	return nil
}

// --- Order management ---

// CreateOrderInput 创建订单请求。
type CreateOrderInput struct {
	UserID         uint64
	PlayerID       *uint64
	GameID         uint64
	Title          string
	Description    string
	PriceCents     int64
	Currency       model.Currency
	ScheduledStart *time.Time
	ScheduledEnd   *time.Time
}

// CreateOrder 新建订单，默认状态为 pending。
func (s *AdminService) CreateOrder(ctx context.Context, in CreateOrderInput) (*model.Order, error) {
	if in.UserID == 0 || in.GameID == 0 || in.PriceCents < 0 || !model.IsValidCurrency(in.Currency) {
		return nil, ErrValidation
	}
	if in.ScheduledStart != nil && in.ScheduledEnd != nil && in.ScheduledEnd.Before(*in.ScheduledStart) {
		return nil, ErrValidation
	}
	if in.PlayerID != nil && *in.PlayerID != 0 {
		if _, err := s.players.Get(ctx, *in.PlayerID); err != nil {
			return nil, err
		}
	}
	order := &model.Order{
		UserID:         in.UserID,
		PlayerID:       0,
		GameID:         in.GameID,
		Title:          strings.TrimSpace(in.Title),
		Description:    strings.TrimSpace(in.Description),
		Status:         model.OrderStatusPending,
		PriceCents:     in.PriceCents,
		Currency:       in.Currency,
		ScheduledStart: in.ScheduledStart,
		ScheduledEnd:   in.ScheduledEnd,
	}
	if in.PlayerID != nil {
		order.PlayerID = *in.PlayerID
	}
	if err := s.orders.Create(ctx, order); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, cacheKeyOrders)
	s.appendLogAsync(ctx, string(model.OpEntityOrder), order.ID, string(model.OpActionCreate), map[string]any{"status": order.Status})
	return order, nil
}

// AssignOrder 指派陪玩师。
func (s *AdminService) AssignOrder(ctx context.Context, id uint64, playerID uint64) (*model.Order, error) {
	if playerID == 0 {
		return nil, ErrValidation
	}
	if _, err := s.players.Get(ctx, playerID); err != nil {
		return nil, err
	}
	order, err := s.orders.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	// 不允许在完成/取消/退款后指派
	switch order.Status {
	case model.OrderStatusCompleted, model.OrderStatusCanceled, model.OrderStatusRefunded:
		return nil, ErrValidation
	}
	order.PlayerID = playerID
	if err := s.orders.Update(ctx, order); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, cacheKeyOrders)
	s.appendLogAsync(ctx, string(model.OpEntityOrder), order.ID, string(model.OpActionAssignPlayer), map[string]any{"player_id": playerID})
	return order, nil
}

// UpdateOrderInput 用于更新订单状态。
type UpdateOrderInput struct {
	Status            model.OrderStatus
	PriceCents        int64
	Currency          model.Currency
	ScheduledStart    *time.Time
	ScheduledEnd      *time.Time
	CancelReason      string
	StartedAt         *time.Time
	CompletedAt       *time.Time
	RefundAmountCents *int64
	RefundReason      string
	RefundedAt        *time.Time
	Note              string
}

// RefundOrderInput 描述退款请求。
type RefundOrderInput struct {
	Reason      string
	AmountCents *int64
	Note        string
}

// OrderTimelineItem 组合订单历史时间线。
type OrderTimelineItem struct {
	ID           uint64         `json:"id"`
	OrderID      uint64         `json:"order_id"`
	PaymentID    *uint64        `json:"payment_id,omitempty"`
	EventType    string         `json:"event_type"`
	Title        string         `json:"title"`
	Description  string         `json:"description,omitempty"`
	Operator     string         `json:"operator,omitempty"`
	OperatorRole string         `json:"operator_role,omitempty"`
	OperatorID   *uint64        `json:"operator_id,omitempty"`
	StatusBefore string         `json:"status_before,omitempty"`
	StatusAfter  string         `json:"status_after,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
}

// OrderRefundItem 描述订单退款记录。
type OrderRefundItem struct {
	ID          uint64     `json:"id"`
	OrderID     uint64     `json:"order_id"`
	PaymentID   uint64     `json:"payment_id"`
	AmountCents int64      `json:"amount_cents"`
	Reason      string     `json:"reason,omitempty"`
	Status      string     `json:"status"`
	Method      string     `json:"refund_method"`
	Note        string     `json:"note,omitempty"`
	RefundedAt  *time.Time `json:"refunded_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ListOrders 列出订单。
func (s *AdminService) ListOrders(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, *model.Pagination, error) {
	normalized := opts
	normalized.Page = repository.NormalizePage(opts.Page)
	normalized.PageSize = repository.NormalizePageSize(opts.PageSize)

	orders, total, err := s.orders.List(ctx, normalized)
	if err != nil {
		return nil, nil, err
	}

	pagination := buildPagination(normalized.Page, normalized.PageSize, total)
	return orders, &pagination, nil
}

// GetOrder 获取订单详情。
func (s *AdminService) GetOrder(ctx context.Context, id uint64) (*model.Order, error) {
	return s.orders.Get(ctx, id)
}

// UpdateOrder 更新订单信息。
func (s *AdminService) UpdateOrder(ctx context.Context, id uint64, input UpdateOrderInput) (*model.Order, error) {
	order, err := s.orders.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if !isValidOrderStatus(input.Status) {
		return nil, ErrValidation
	}
	if !model.IsValidCurrency(input.Currency) {
		return nil, ErrValidation
	}
	if input.PriceCents < 0 {
		return nil, ErrValidation
	}
	if input.ScheduledStart != nil && input.ScheduledEnd != nil && input.ScheduledEnd.Before(*input.ScheduledStart) {
		return nil, ErrValidation
	}

	// state machine guard
	if !isAllowedOrderTransition(order.Status, input.Status) {
		return nil, ErrOrderInvalidTransition
	}

	prevStatus := order.Status

	order.Status = input.Status
	order.PriceCents = input.PriceCents
	order.Currency = input.Currency
	order.ScheduledStart = input.ScheduledStart
	order.ScheduledEnd = input.ScheduledEnd
	order.CancelReason = strings.TrimSpace(input.CancelReason)
	if input.StartedAt != nil {
		order.StartedAt = input.StartedAt
	}
	if input.CompletedAt != nil {
		order.CompletedAt = input.CompletedAt
	}
	if input.RefundAmountCents != nil {
		order.RefundAmountCents = *input.RefundAmountCents
	}
	if input.RefundReason != "" || input.RefundAmountCents != nil {
		order.RefundReason = strings.TrimSpace(input.RefundReason)
	}
	if input.RefundedAt != nil {
		order.RefundedAt = input.RefundedAt
	}

	if err := s.orders.Update(ctx, order); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, cacheKeyOrders)
	action := model.OpActionUpdateStatus
	if order.Status == model.OrderStatusCanceled {
		action = model.OpActionCancel
	} else if order.Status == model.OrderStatusRefunded {
		action = model.OpActionRefund
	} else {
		switch {
		case prevStatus == model.OrderStatusPending && order.Status == model.OrderStatusConfirmed:
			action = model.OpActionConfirm
		case prevStatus == model.OrderStatusConfirmed && order.Status == model.OrderStatusInProgress:
			action = model.OpActionStart
		case prevStatus == model.OrderStatusInProgress && order.Status == model.OrderStatusCompleted:
			action = model.OpActionComplete
		}
	}
	meta := map[string]any{
		"status":      order.Status,
		"from_status": prevStatus,
	}
	if order.CancelReason != "" {
		meta["reason"] = order.CancelReason
	}
	if input.Note != "" {
		meta["note"] = strings.TrimSpace(input.Note)
	}
	if order.StartedAt != nil {
		meta["started_at"] = order.StartedAt.Format(time.RFC3339)
	}
	if order.CompletedAt != nil {
		meta["completed_at"] = order.CompletedAt.Format(time.RFC3339)
	}
	if input.RefundAmountCents != nil {
		meta["refund_amount_cents"] = order.RefundAmountCents
	}
	if order.RefundReason != "" {
		meta["refund_reason"] = order.RefundReason
	}
	if order.RefundedAt != nil {
		meta["refunded_at"] = order.RefundedAt.Format(time.RFC3339)
	}
	s.appendLogAsync(ctx, string(model.OpEntityOrder), order.ID, string(action), meta)
	if rid, ok := logging.RequestIDFromContext(ctx); ok {
		slog.Info("order_status_changed", slog.Uint64("order_id", order.ID), slog.String("status", string(order.Status)), slog.String("request_id", rid))
	} else {
		slog.Info("order_status_changed", slog.Uint64("order_id", order.ID), slog.String("status", string(order.Status)))
	}
	return order, nil
}

// ConfirmOrder 将订单从 pending 确认到 confirmed。
func (s *AdminService) ConfirmOrder(ctx context.Context, id uint64, note string) (*model.Order, error) {
	order, err := s.orders.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	note = strings.TrimSpace(note)
	return s.UpdateOrder(ctx, id, UpdateOrderInput{
		Status:         model.OrderStatusConfirmed,
		PriceCents:     order.PriceCents,
		Currency:       order.Currency,
		ScheduledStart: order.ScheduledStart,
		ScheduledEnd:   order.ScheduledEnd,
		CancelReason:   order.CancelReason,
		StartedAt:      order.StartedAt,
		CompletedAt:    order.CompletedAt,
		RefundReason:   order.RefundReason,
		RefundedAt:     order.RefundedAt,
		Note:           note,
	})
}

// StartOrder 将订单置为进行中，并记录实际开始时间。
func (s *AdminService) StartOrder(ctx context.Context, id uint64, note string) (*model.Order, error) {
	order, err := s.orders.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	note = strings.TrimSpace(note)
	startedAt := time.Now().UTC()
	return s.UpdateOrder(ctx, id, UpdateOrderInput{
		Status:         model.OrderStatusInProgress,
		PriceCents:     order.PriceCents,
		Currency:       order.Currency,
		ScheduledStart: order.ScheduledStart,
		ScheduledEnd:   order.ScheduledEnd,
		CancelReason:   order.CancelReason,
		StartedAt:      &startedAt,
		CompletedAt:    order.CompletedAt,
		RefundReason:   order.RefundReason,
		RefundedAt:     order.RefundedAt,
		Note:           note,
	})
}

// CompleteOrder 完成订单服务，并记录完成时间。
func (s *AdminService) CompleteOrder(ctx context.Context, id uint64, note string) (*model.Order, error) {
	order, err := s.orders.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	note = strings.TrimSpace(note)
	completedAt := time.Now().UTC()
	return s.UpdateOrder(ctx, id, UpdateOrderInput{
		Status:         model.OrderStatusCompleted,
		PriceCents:     order.PriceCents,
		Currency:       order.Currency,
		ScheduledStart: order.ScheduledStart,
		ScheduledEnd:   order.ScheduledEnd,
		CancelReason:   order.CancelReason,
		StartedAt:      order.StartedAt,
		CompletedAt:    &completedAt,
		RefundReason:   order.RefundReason,
		RefundedAt:     order.RefundedAt,
		Note:           note,
	})
}

// RefundOrder 执行退款并记录退款信息。
func (s *AdminService) RefundOrder(ctx context.Context, id uint64, input RefundOrderInput) (*model.Order, error) {
	order, err := s.orders.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	reason := strings.TrimSpace(input.Reason)
	if reason == "" {
		return nil, ErrValidation
	}
	switch order.Status {
	case model.OrderStatusCompleted, model.OrderStatusInProgress, model.OrderStatusConfirmed:
		// allowed
	default:
		return nil, ErrValidation
	}
	amount := order.PriceCents
	if input.AmountCents != nil {
		if *input.AmountCents <= 0 || *input.AmountCents > order.PriceCents {
			return nil, ErrValidation
		}
		amount = *input.AmountCents
	}
	refundedAt := time.Now().UTC()
	note := strings.TrimSpace(input.Note)
	updatedOrder, err := s.UpdateOrder(ctx, id, UpdateOrderInput{
		Status:            model.OrderStatusRefunded,
		PriceCents:        order.PriceCents,
		Currency:          order.Currency,
		ScheduledStart:    order.ScheduledStart,
		ScheduledEnd:      order.ScheduledEnd,
		CancelReason:      order.CancelReason,
		StartedAt:         order.StartedAt,
		CompletedAt:       order.CompletedAt,
		RefundAmountCents: &amount,
		RefundReason:      reason,
		RefundedAt:        &refundedAt,
		Note:              note,
	})
	if err != nil {
		return nil, err
	}

	// 更新相关支付为已退款状态（若存在）
	payments, err := s.listPaymentsByOrder(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, pay := range payments {
		if pay.Status == model.PaymentStatusRefunded {
			continue
		}
		if pay.Status == model.PaymentStatusPaid || pay.Status == model.PaymentStatusPending {
			_, updErr := s.UpdatePayment(ctx, pay.ID, UpdatePaymentInput{
				Status:          model.PaymentStatusRefunded,
				ProviderTradeNo: pay.ProviderTradeNo,
				ProviderRaw:     pay.ProviderRaw,
				PaidAt:          pay.PaidAt,
				RefundedAt:      &refundedAt,
			})
			if updErr != nil && !errors.Is(updErr, ErrValidation) {
				return nil, updErr
			}
		}
	}
	return updatedOrder, nil
}

// GetOrderPayments 返回订单下的所有支付记录。
func (s *AdminService) GetOrderPayments(ctx context.Context, orderID uint64) ([]model.Payment, error) {
	return s.listPaymentsByOrder(ctx, orderID)
}

// GetOrderRefunds 汇总订单退款记录（基于支付信息与订单字段）。
func (s *AdminService) GetOrderRefunds(ctx context.Context, orderID uint64) ([]OrderRefundItem, error) {
	order, err := s.orders.Get(ctx, orderID)
	if err != nil {
		return nil, err
	}
	payments, err := s.listPaymentsByOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	result := make([]OrderRefundItem, 0)
	for _, pay := range payments {
		if pay.Status != model.PaymentStatusRefunded {
			continue
		}
		item := OrderRefundItem{
			ID:          pay.ID,
			OrderID:     orderID,
			PaymentID:   pay.ID,
			AmountCents: pay.AmountCents,
			Method:      string(pay.Method),
			Status:      mapRefundStatus(pay.Status),
			RefundedAt:  pay.RefundedAt,
			CreatedAt:   pay.CreatedAt,
			Reason:      order.RefundReason,
			Note:        order.RefundReason,
		}
		result = append(result, item)
	}

	// 如果订单存在退款金额但支付记录未覆盖，则补充一条摘要信息
	if order.RefundAmountCents > 0 {
		hasSummary := false
		for _, item := range result {
			if item.AmountCents == order.RefundAmountCents {
				hasSummary = true
				break
			}
		}
		if !hasSummary {
			createdAt := order.UpdatedAt
			if order.RefundedAt != nil {
				createdAt = *order.RefundedAt
			}
			item := OrderRefundItem{
				ID:          orderID*10 + 1,
				OrderID:     orderID,
				PaymentID:   0,
				AmountCents: order.RefundAmountCents,
				Method:      "unknown",
				Status:      "success",
				Reason:      order.RefundReason,
				RefundedAt:  order.RefundedAt,
				CreatedAt:   createdAt,
				Note:        order.RefundReason,
			}
			result = append(result, item)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	return result, nil
}

// GetOrderReviews 返回订单相关的全部评价。
func (s *AdminService) GetOrderReviews(ctx context.Context, orderID uint64) ([]model.Review, error) {
	reviews := make([]model.Review, 0)
	page := 1
	orderIDCopy := orderID
	for {
		opts := repository.ReviewListOptions{
			Page:     page,
			PageSize: 200,
			OrderID:  &orderIDCopy,
		}
		items, pagination, err := s.ListReviews(ctx, opts)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, items...)
		if pagination == nil || !pagination.HasNext {
			break
		}
		page++
	}
	return reviews, nil
}

// GetOrderTimeline 汇总订单的状态流转与关键事件。
func (s *AdminService) GetOrderTimeline(ctx context.Context, orderID uint64) ([]OrderTimelineItem, error) {
	order, err := s.orders.Get(ctx, orderID)
	if err != nil {
		return nil, err
	}
	logs, err := s.collectOperationLogs(ctx, string(model.OpEntityOrder), orderID)
	if err != nil {
		return nil, err
	}

	userCache := make(map[uint64]*model.User)
	items := make([]OrderTimelineItem, 0, len(logs))
	for _, logEntry := range logs {
		meta := map[string]any{}
		if len(logEntry.MetadataJSON) > 0 {
			_ = json.Unmarshal(logEntry.MetadataJSON, &meta)
		}
		item := OrderTimelineItem{
			ID:        logEntry.ID,
			OrderID:   orderID,
			EventType: mapTimelineEventType(logEntry.Action),
			Title:     mapTimelineTitle(logEntry.Action),
			Metadata:  meta,
			CreatedAt: logEntry.CreatedAt,
		}
		if note, ok := meta["note"].(string); ok && strings.TrimSpace(note) != "" {
			item.Description = strings.TrimSpace(note)
		} else if reason, ok := meta["reason"].(string); ok && strings.TrimSpace(reason) != "" {
			item.Description = strings.TrimSpace(reason)
		}
		if before, ok := meta["from_status"].(string); ok {
			item.StatusBefore = before
		}
		if after, ok := meta["status"].(string); ok {
			item.StatusAfter = after
		}
		if logEntry.ActorUserID != nil {
			if user := s.resolveUser(ctx, userCache, *logEntry.ActorUserID); user != nil {
				item.Operator = user.Name
				item.OperatorRole = string(user.Role)
				id := user.ID
				item.OperatorID = &id
			}
		}
		items = append(items, item)
	}

	// 追加支付关键事件
	payments, err := s.listPaymentsByOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	for _, pay := range payments {
		if pay.PaidAt != nil {
			item := OrderTimelineItem{
				ID:        pay.ID*10 + 1,
				OrderID:   orderID,
				PaymentID: ptrUint64(pay.ID),
				EventType: "action",
				Title:     "支付确认",
				Metadata: map[string]any{
					"payment_status": pay.Status,
					"payment_method": pay.Method,
					"amount_cents":   pay.AmountCents,
				},
				CreatedAt: *pay.PaidAt,
			}
			items = append(items, item)
		}
		if pay.RefundedAt != nil {
			item := OrderTimelineItem{
				ID:          pay.ID*10 + 2,
				OrderID:     orderID,
				PaymentID:   ptrUint64(pay.ID),
				EventType:   "status_change",
				Title:       "支付退款",
				Description: strings.TrimSpace(order.RefundReason),
				Metadata: map[string]any{
					"payment_status": pay.Status,
					"payment_method": pay.Method,
					"amount_cents":   pay.AmountCents,
				},
				CreatedAt:   *pay.RefundedAt,
				StatusAfter: string(model.OrderStatusRefunded),
			}
			items = append(items, item)
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})

	return items, nil
}

// DeleteOrder 删除订单（软删）。
func (s *AdminService) DeleteOrder(ctx context.Context, id uint64) error {
	if err := s.orders.Delete(ctx, id); err != nil {
		return err
	}
	s.invalidateCache(ctx, cacheKeyOrders)
	s.appendLogAsync(ctx, string(model.OpEntityOrder), id, string(model.OpActionDelete), nil)
	return nil
}

// --- Payment management ---

// UpdatePaymentInput 调整支付状态。
type UpdatePaymentInput struct {
	Status          model.PaymentStatus
	ProviderTradeNo string
	ProviderRaw     json.RawMessage
	PaidAt          *time.Time
	RefundedAt      *time.Time
}

// CreatePaymentInput 创建支付记录。
type CreatePaymentInput struct {
	OrderID     uint64
	UserID      uint64
	Method      model.PaymentMethod
	AmountCents int64
	Currency    model.Currency
	ProviderRaw json.RawMessage
}

// CreatePayment 新建支付记录，默认状态 pending。
func (s *AdminService) CreatePayment(ctx context.Context, in CreatePaymentInput) (*model.Payment, error) {
	if in.OrderID == 0 || in.UserID == 0 || in.AmountCents <= 0 || !model.IsValidCurrency(in.Currency) {
		return nil, ErrValidation
	}
	if in.Method == "" {
		return nil, ErrValidation
	}
	if _, err := s.orders.Get(ctx, in.OrderID); err != nil {
		return nil, err
	}
	if _, err := s.users.Get(ctx, in.UserID); err != nil {
		return nil, mapUserError(err)
	}
	pay := &model.Payment{
		OrderID:     in.OrderID,
		UserID:      in.UserID,
		Method:      in.Method,
		AmountCents: in.AmountCents,
		Currency:    in.Currency,
		Status:      model.PaymentStatusPending,
		ProviderRaw: in.ProviderRaw,
	}
	if err := s.payments.Create(ctx, pay); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, cacheKeyPayments)
	s.appendLogAsync(ctx, string(model.OpEntityPayment), pay.ID, string(model.OpActionCreate), map[string]any{"status": pay.Status, "method": pay.Method})
	return pay, nil
}

// CapturePaymentInput 确认入账。
type CapturePaymentInput struct {
	ProviderTradeNo string
	ProviderRaw     json.RawMessage
	PaidAt          *time.Time
}

// CapturePayment 将支付置为 paid。
func (s *AdminService) CapturePayment(ctx context.Context, id uint64, in CapturePaymentInput) (*model.Payment, error) {
	pay, err := s.payments.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !isAllowedPaymentTransition(pay.Status, model.PaymentStatusPaid) {
		return nil, ErrValidation
	}
	pay.Status = model.PaymentStatusPaid
	pay.ProviderTradeNo = strings.TrimSpace(in.ProviderTradeNo)
	pay.ProviderRaw = in.ProviderRaw
	if in.PaidAt != nil {
		pay.PaidAt = in.PaidAt
	} else {
		now := time.Now().UTC()
		pay.PaidAt = &now
	}
	if err := s.payments.Update(ctx, pay); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, cacheKeyPayments)
	s.appendLogAsync(ctx, string(model.OpEntityPayment), pay.ID, string(model.OpActionCapture), map[string]any{"trade_no": pay.ProviderTradeNo})
	return pay, nil
}

// ListPayments 列出支付记录。
func (s *AdminService) ListPayments(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, *model.Pagination, error) {
	normalized := opts
	normalized.Page = repository.NormalizePage(opts.Page)
	normalized.PageSize = repository.NormalizePageSize(opts.PageSize)

	payments, total, err := s.payments.List(ctx, normalized)
	if err != nil {
		return nil, nil, err
	}

	pagination := buildPagination(normalized.Page, normalized.PageSize, total)
	return payments, &pagination, nil
}

// GetPayment 获取支付详情。
func (s *AdminService) GetPayment(ctx context.Context, id uint64) (*model.Payment, error) {
	return s.payments.Get(ctx, id)
}

// UpdatePayment 更新支付状态。
func (s *AdminService) UpdatePayment(ctx context.Context, id uint64, input UpdatePaymentInput) (*model.Payment, error) {
	payment, err := s.payments.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if !isValidPaymentStatus(input.Status) {
		return nil, ErrValidation
	}

	if !isAllowedPaymentTransition(payment.Status, input.Status) {
		return nil, ErrValidation
	}

	payment.Status = input.Status
	payment.ProviderTradeNo = strings.TrimSpace(input.ProviderTradeNo)
	payment.ProviderRaw = input.ProviderRaw
	payment.PaidAt = input.PaidAt
	payment.RefundedAt = input.RefundedAt

	if err := s.payments.Update(ctx, payment); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, cacheKeyPayments)
	payAction := model.OpActionUpdateStatus
	if input.Status == model.PaymentStatusRefunded {
		payAction = model.OpActionRefund
	}
	s.appendLogAsync(ctx, string(model.OpEntityPayment), payment.ID, string(payAction), map[string]any{"status": payment.Status})
	if rid, ok := logging.RequestIDFromContext(ctx); ok {
		slog.Info("payment_status_changed", slog.Uint64("payment_id", payment.ID), slog.String("status", string(payment.Status)), slog.String("request_id", rid))
	} else {
		slog.Info("payment_status_changed", slog.Uint64("payment_id", payment.ID), slog.String("status", string(payment.Status)))
	}
	return payment, nil
}

// DeletePayment 删除支付记录。
func (s *AdminService) DeletePayment(ctx context.Context, id uint64) error {
	if err := s.payments.Delete(ctx, id); err != nil {
		return err
	}
	s.invalidateCache(ctx, cacheKeyPayments)
	s.appendLogAsync(ctx, string(model.OpEntityPayment), id, string(model.OpActionDelete), nil)
	return nil
}

// appendLogAsync 追加操作日志（尽力而为，不影响主流程）。
func (s *AdminService) appendLogAsync(ctx context.Context, entity string, id uint64, action string, meta map[string]any) {
	if s.tx == nil {
		return
	}
	_ = s.tx.WithTx(ctx, func(r *gormrepo.Repos) error {
		var raw []byte
		if meta != nil {
			if b, err := json.Marshal(meta); err == nil {
				raw = b
			}
		}
		var actorPtr *uint64
		if uid, ok := logging.ActorUserIDFromContext(ctx); ok {
			actorID := uid
			actorPtr = &actorID
		}
		log := &model.OperationLog{EntityType: entity, EntityID: id, Action: action, ActorUserID: actorPtr, MetadataJSON: raw}
		return r.OpLogs.Append(ctx, log)
	})
}

func isValidOrderStatus(status model.OrderStatus) bool {
	switch status {
	case model.OrderStatusPending, model.OrderStatusConfirmed, model.OrderStatusInProgress,
		model.OrderStatusCompleted, model.OrderStatusCanceled, model.OrderStatusRefunded:
		return true
	default:
		return false
	}
}

func isAllowedOrderTransition(prev, next model.OrderStatus) bool {
	if prev == next {
		return true
	}
	switch prev {
	case model.OrderStatusPending:
		return next == model.OrderStatusConfirmed || next == model.OrderStatusCanceled || next == model.OrderStatusRefunded
	case model.OrderStatusConfirmed:
		return next == model.OrderStatusInProgress || next == model.OrderStatusCanceled || next == model.OrderStatusRefunded
	case model.OrderStatusInProgress:
		return next == model.OrderStatusCompleted || next == model.OrderStatusCanceled || next == model.OrderStatusRefunded
	case model.OrderStatusCompleted:
		return next == model.OrderStatusRefunded
	case model.OrderStatusCanceled, model.OrderStatusRefunded:
		return false
	default:
		return false
	}
}

func isValidPaymentStatus(status model.PaymentStatus) bool {
	switch status {
	case model.PaymentStatusPending, model.PaymentStatusPaid, model.PaymentStatusFailed, model.PaymentStatusRefunded:
		return true
	default:
		return false
	}
}

func isAllowedPaymentTransition(prev, next model.PaymentStatus) bool {
	if prev == next {
		return true
	}
	switch prev {
	case model.PaymentStatusPending:
		return next == model.PaymentStatusPaid || next == model.PaymentStatusFailed || next == model.PaymentStatusRefunded
	case model.PaymentStatusPaid:
		return next == model.PaymentStatusRefunded
	case model.PaymentStatusFailed, model.PaymentStatusRefunded:
		return false
	default:
		return false
	}
}

func (s *AdminService) listPaymentsByOrder(ctx context.Context, orderID uint64) ([]model.Payment, error) {
	result := make([]model.Payment, 0)
	page := 1
	for {
		opts := repository.PaymentListOptions{
			Page:     page,
			PageSize: 200,
		}
		orderIDCopy := orderID
		opts.OrderID = &orderIDCopy
		items, pagination, err := s.ListPayments(ctx, opts)
		if err != nil {
			return nil, err
		}
		result = append(result, items...)
		if pagination == nil || !pagination.HasNext {
			break
		}
		page++
	}
	return result, nil
}

func (s *AdminService) collectOperationLogs(ctx context.Context, entityType string, entityID uint64) ([]model.OperationLog, error) {
	all := make([]model.OperationLog, 0)
	page := 1
	for {
		opts := repository.OperationLogListOptions{
			Page:     page,
			PageSize: 200,
		}
		items, pagination, err := s.ListOperationLogs(ctx, entityType, entityID, opts)
		if err != nil {
			return nil, err
		}
		all = append(all, items...)
		if pagination == nil || !pagination.HasNext {
			break
		}
		page++
	}
	return all, nil
}

func (s *AdminService) resolveUser(ctx context.Context, cache map[uint64]*model.User, id uint64) *model.User {
	if user, ok := cache[id]; ok {
		return user
	}
	user, err := s.users.Get(ctx, id)
	if err != nil {
		cache[id] = nil
		return nil
	}
	cache[id] = user
	return user
}

func mapUserError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, repository.ErrNotFound) {
		return ErrUserNotFound
	}
	return err
}

func mapRefundStatus(status model.PaymentStatus) string {
	switch status {
	case model.PaymentStatusRefunded:
		return "success"
	case model.PaymentStatusPending:
		return "pending"
	case model.PaymentStatusFailed:
		return "failed"
	default:
		return strings.ToLower(string(status))
	}
}

func mapTimelineEventType(action string) string {
	switch action {
	case string(model.OpActionCreate):
		return "system"
	case string(model.OpActionAssignPlayer):
		return "action"
	case string(model.OpActionConfirm), string(model.OpActionStart), string(model.OpActionComplete),
		string(model.OpActionUpdateStatus), string(model.OpActionCancel), string(model.OpActionRefund):
		return "status_change"
	default:
		return "action"
	}
}

func mapTimelineTitle(action string) string {
	switch action {
	case string(model.OpActionCreate):
		return "订单创建"
	case string(model.OpActionAssignPlayer):
		return "指派陪玩师"
	case string(model.OpActionConfirm):
		return "订单确认"
	case string(model.OpActionStart):
		return "开始服务"
	case string(model.OpActionComplete):
		return "完成订单"
	case string(model.OpActionCancel):
		return "订单取消"
	case string(model.OpActionRefund):
		return "订单退款"
	case string(model.OpActionUpdateStatus):
		return "状态更新"
	default:
		return strings.ReplaceAll(action, "_", " ")
	}
}

func ptrUint64(id uint64) *uint64 {
	return &id
}

func (s *AdminService) invalidateCache(ctx context.Context, keys ...string) {
	if s.cache == nil {
		return
	}
	for _, key := range keys {
		_ = s.cache.Delete(ctx, key)
	}
}

func buildPagination(page, pageSize int, total int64) model.Pagination {
	totalPages := 0
	if pageSize > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	return model.Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// ListOperationLogs 返回实体的操作日志。
func (s *AdminService) ListOperationLogs(ctx context.Context, entityType string, entityID uint64, opts repository.OperationLogListOptions) ([]model.OperationLog, *model.Pagination, error) {
	if s.tx == nil {
		return nil, nil, errors.New("transaction manager not configured")
	}
	var logs []model.OperationLog
	var total int64
	err := s.tx.WithTx(ctx, func(r *gormrepo.Repos) error {
		norm := repository.OperationLogListOptions{
			Page:        repository.NormalizePage(opts.Page),
			PageSize:    repository.NormalizePageSize(opts.PageSize),
			Action:      opts.Action,
			ActorUserID: opts.ActorUserID,
			DateFrom:    opts.DateFrom,
			DateTo:      opts.DateTo,
		}
		items, cnt, err := r.OpLogs.ListByEntity(ctx, entityType, entityID, norm)
		if err != nil {
			return err
		}
		logs, total = items, cnt
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	p := buildPagination(repository.NormalizePage(opts.Page), repository.NormalizePageSize(opts.PageSize), total)
	return logs, &p, nil
}

// --- Review management ---

// ListReviews 列出评价。
func (s *AdminService) ListReviews(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, *model.Pagination, error) {
	if s.tx == nil {
		return nil, nil, errors.New("transaction manager not configured")
	}
	var items []model.Review
	var total int64
	err := s.tx.WithTx(ctx, func(r *gormrepo.Repos) error {
		page := repository.NormalizePage(opts.Page)
		size := repository.NormalizePageSize(opts.PageSize)
		out, cnt, err := r.Reviews.List(ctx, repository.ReviewListOptions{
			Page: page, PageSize: size, OrderID: opts.OrderID, UserID: opts.UserID, PlayerID: opts.PlayerID, DateFrom: opts.DateFrom, DateTo: opts.DateTo,
		})
		if err != nil {
			return err
		}
		items, total = out, cnt
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	p := buildPagination(repository.NormalizePage(opts.Page), repository.NormalizePageSize(opts.PageSize), total)
	return items, &p, nil
}

// GetReview 返回评价详情。
func (s *AdminService) GetReview(ctx context.Context, id uint64) (*model.Review, error) {
	if s.tx == nil {
		return nil, errors.New("transaction manager not configured")
	}
	var item *model.Review
	err := s.tx.WithTx(ctx, func(r *gormrepo.Repos) error {
		var err error
		item, err = r.Reviews.Get(ctx, id)
		return err
	})
	return item, err
}

// CreateReview 新建评价。
func (s *AdminService) CreateReview(ctx context.Context, r model.Review) (*model.Review, error) {
	if !r.Score.Valid() || r.OrderID == 0 || r.UserID == 0 || r.PlayerID == 0 {
		return nil, ErrValidation
	}
	if s.tx == nil {
		return nil, errors.New("transaction manager not configured")
	}
	err := s.tx.WithTx(ctx, func(txr *gormrepo.Repos) error { return txr.Reviews.Create(ctx, &r) })
	if err != nil {
		return nil, err
	}
	s.appendLogAsync(ctx, string(model.OpEntityReview), r.ID, string(model.OpActionCreate), map[string]any{"order_id": r.OrderID, "player_id": r.PlayerID})
	return &r, nil
}

// UpdateReview 修改评价分数/内容。
func (s *AdminService) UpdateReview(ctx context.Context, id uint64, score model.Rating, content string) (*model.Review, error) {
	if !score.Valid() {
		return nil, ErrValidation
	}
	if s.tx == nil {
		return nil, errors.New("transaction manager not configured")
	}
	var item *model.Review
	err := s.tx.WithTx(ctx, func(r *gormrepo.Repos) error {
		obj, err := r.Reviews.Get(ctx, id)
		if err != nil {
			return err
		}
		obj.Score = score
		obj.Content = strings.TrimSpace(content)
		if err := r.Reviews.Update(ctx, obj); err != nil {
			return err
		}
		item = obj
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.appendLogAsync(ctx, string(model.OpEntityReview), id, string(model.OpActionUpdate), nil)
	return item, nil
}

// DeleteReview 删除评价。
func (s *AdminService) DeleteReview(ctx context.Context, id uint64) error {
	if s.tx == nil {
		return errors.New("transaction manager not configured")
	}
	return s.tx.WithTx(ctx, func(r *gormrepo.Repos) error { return r.Reviews.Delete(ctx, id) })
}

func getCachedList[T any](ctx context.Context, c cache.Cache, key string, ttl time.Duration, fetch func() ([]T, error)) ([]T, error) {
	if c != nil {
		if raw, ok, err := c.Get(ctx, key); err == nil && ok {
			var cached []T
			if err := json.Unmarshal([]byte(raw), &cached); err == nil {
				return cached, nil
			}
			_ = c.Delete(ctx, key)
		}
	}

	result, err := fetch()
	if err != nil {
		return nil, err
	}

	if c != nil {
		if data, err := json.Marshal(result); err == nil {
			_ = c.Set(ctx, key, string(data), ttl)
		}
	}

	return result, nil
}
