package service

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gamelink/internal/cache"
	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// ErrValidation 表示输入校验失败。
var ErrValidation = errors.New("validation failed")

// ErrNotFound 暴露仓储的未找到错误，便于 handler 判定。
var ErrNotFound = repository.ErrNotFound

// AdminService 聚合后台管理所需的业务逻辑。
type AdminService struct {
	games    repository.GameRepository
	users    repository.UserRepository
	players  repository.PlayerRepository
	orders   repository.OrderRepository
	payments repository.PaymentRepository
	cache    cache.Cache
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

	return game, nil
}

// DeleteGame 删除游戏。
func (s *AdminService) DeleteGame(ctx context.Context, id uint64) error {
	if err := s.games.Delete(ctx, id); err != nil {
		return err
	}
	s.invalidateCache(ctx, cacheKeyGames)
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

// GetUser 返回指定用户。
func (s *AdminService) GetUser(ctx context.Context, id uint64) (*model.User, error) {
	return s.users.Get(ctx, id)
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
	return user, nil
}

// UpdateUser 更新用户基础信息。
func (s *AdminService) UpdateUser(ctx context.Context, id uint64, input UpdateUserInput) (*model.User, error) {
	user, err := s.users.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := validateUserInput(input.Name, input.Role, input.Status, optionalPassword(input.Password)); err != nil {
		return nil, err
	}

	user.Phone = strings.TrimSpace(input.Phone)
	user.Email = strings.TrimSpace(input.Email)
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
	return user, nil
}

// DeleteUser 删除用户。
func (s *AdminService) DeleteUser(ctx context.Context, id uint64) error {
	if err := s.users.Delete(ctx, id); err != nil {
		return err
	}
	s.invalidateCache(ctx, cacheKeyUsers)
	return nil
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
	HourlyRateCents    int64
	MainGameID         uint64
	VerificationStatus model.VerificationStatus
}

// UpdatePlayerInput 更新陪玩资料。
type UpdatePlayerInput struct {
	Nickname           string
	Bio                string
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
	return player, nil
}

// DeletePlayer 删除陪玩档案。
func (s *AdminService) DeletePlayer(ctx context.Context, id uint64) error {
	if err := s.players.Delete(ctx, id); err != nil {
		return err
	}
	s.invalidateCache(ctx, cacheKeyPlayers)
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

// UpdateOrderInput 用于更新订单状态。
type UpdateOrderInput struct {
	Status         model.OrderStatus
	PriceCents     int64
	Currency       model.Currency
	ScheduledStart *time.Time
	ScheduledEnd   *time.Time
	CancelReason   string
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

	order.Status = input.Status
	order.PriceCents = input.PriceCents
	order.Currency = input.Currency
	order.ScheduledStart = input.ScheduledStart
	order.ScheduledEnd = input.ScheduledEnd
	order.CancelReason = strings.TrimSpace(input.CancelReason)

	if err := s.orders.Update(ctx, order); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, cacheKeyOrders)
	return order, nil
}

// DeleteOrder 删除订单（软删）。
func (s *AdminService) DeleteOrder(ctx context.Context, id uint64) error {
	if err := s.orders.Delete(ctx, id); err != nil {
		return err
	}
	s.invalidateCache(ctx, cacheKeyOrders)
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

	payment.Status = input.Status
	payment.ProviderTradeNo = strings.TrimSpace(input.ProviderTradeNo)
	payment.ProviderRaw = input.ProviderRaw
	payment.PaidAt = input.PaidAt
	payment.RefundedAt = input.RefundedAt

	if err := s.payments.Update(ctx, payment); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, cacheKeyPayments)
	return payment, nil
}

// DeletePayment 删除支付记录。
func (s *AdminService) DeletePayment(ctx context.Context, id uint64) error {
	if err := s.payments.Delete(ctx, id); err != nil {
		return err
	}
	s.invalidateCache(ctx, cacheKeyPayments)
	return nil
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

func isValidPaymentStatus(status model.PaymentStatus) bool {
	switch status {
	case model.PaymentStatusPending, model.PaymentStatusPaid, model.PaymentStatusFailed, model.PaymentStatusRefunded:
		return true
	default:
		return false
	}
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
