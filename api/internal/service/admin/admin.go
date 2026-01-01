package admin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/common"
	repoiface "gamelink/internal/repository/interfaces"
	feedservice "gamelink/internal/service/feed"
	"gamelink/pkg/apierr"
	"gamelink/pkg/cache"
	"gamelink/pkg/logging"
)

var (
	// ErrValidation 表示输入校验失败。
	ErrValidation = apierr.BadRequest("validation failed")
	// ErrUserNotFound 用于统一标识用户不存在的场景。
	ErrUserNotFound = apierr.NotFound("user not found")
	// ErrOrderInvalidTransition 代表订单状态流转不合法。
	ErrOrderInvalidTransition = apierr.BadRequest("invalid order status transition")
	// ErrUnauthorized 表示无权操作。
	ErrUnauthorized = apierr.Forbidden("unauthorized")

	// ErrNotFound 暴露仓储的未找到错误，便于 handler 判定。
	ErrNotFound = repository.ErrNotFound
)

// AdminService 聚合后台管理所需的业务逻辑。
type AdminService struct {
	games          repository.GameRepository
	users          repository.UserRepository
	players        repository.PlayerRepository
	orders         repoiface.OrderRepository
	payments       repository.PaymentRepository
	roles          repository.RoleRepository
	serviceItems   repository.ServiceItemRepository // 服务项目仓库
	permissions    repository.PermissionRepository
	menus          repository.MenuRepository
	stats          repository.StatsRepository
	wallets        repository.WalletRepository       // 用户钱包仓库
	gameCategories repository.GameCategoryRepository // 游戏分类仓库
	cache          cache.Cache
	tx             TxManager
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
	orders repoiface.OrderRepository,
	payments repository.PaymentRepository,
	roles repository.RoleRepository,
	serviceItems repository.ServiceItemRepository,
	permissions repository.PermissionRepository,
	menus repository.MenuRepository,
	stats repository.StatsRepository,
	wallets repository.WalletRepository,
	gameCategories repository.GameCategoryRepository,
	cache cache.Cache,
) *AdminService {
	return &AdminService{
		games:          games,
		users:          users,
		players:        players,
		orders:         orders,
		payments:       payments,
		roles:          roles,
		serviceItems:   serviceItems,
		permissions:    permissions,
		menus:          menus,
		stats:          stats,
		wallets:        wallets,
		gameCategories: gameCategories,
		cache:          cache,
	}
}

// TxManager abstracts UnitOfWork for transactional operations.
type TxManager interface {
	WithTx(ctx context.Context, fn func(r *common.Repos) error) error
}

// PermissionService 获取权限服务
func (s *AdminService) PermissionService() *PermissionService {
	return NewPermissionService(s.permissions, s.cache)
}

// RoleService 获取角色服务
func (s *AdminService) RoleService() *RoleService {
	return NewRoleService(s.roles, s.cache)
}

// MenuService 获取菜单服务
func (s *AdminService) MenuService() *MenuService {
	return NewMenuService(s.menus)
}

// StatsService 获取统计服务
func (s *AdminService) StatsService() *StatsService {
	return NewStatsService(s.stats)
}

// SetTxManager injects a transaction manager.
func (s *AdminService) SetTxManager(tx TxManager) { s.tx = tx }

// UpdatePlayerSkillTags 替换玩家技能标签集合（需要 TxManager）。
func (s *AdminService) UpdatePlayerSkillTags(ctx context.Context, playerID uint64, tags []string) error {
	if s.tx == nil {
		return apierr.InternalError("事务管理器未配置")
	}
	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
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
		return nil, nil, apierr.InternalError("事务管理器未配置")
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

	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
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

// ListGamesPagedWithFilter 返回带筛选的分页游戏列表。
func (s *AdminService) ListGamesPagedWithFilter(ctx context.Context, page, pageSize int, keyword string) ([]model.Game, *model.Pagination, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	items, total, err := s.games.ListPagedWithFilter(ctx, page, pageSize, keyword)
	if err != nil {
		return nil, nil, err
	}
	p := buildPagination(page, pageSize, total)
	return items, &p, nil
}

// BatchDeleteGames 批量删除游戏。
func (s *AdminService) BatchDeleteGames(ctx context.Context, ids []uint64) (int64, error) {
	if len(ids) == 0 {
		return 0, apierr.BadRequest("no game ids provided")
	}
	deleted, err := s.games.BatchDelete(ctx, ids)
	if err != nil {
		return 0, WrapError(err, "batch delete games")
	}
	s.invalidateCache(ctx, cacheKeyGames)
	return deleted, nil
}

// GetGame 获取单个游戏详情。
func (s *AdminService) GetGame(ctx context.Context, id uint64) (*model.Game, error) {
	game, err := s.games.Get(ctx, id)
	if err != nil {
		return nil, WrapError(err, "get game")
	}
	return game, nil
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
		return nil, WrapError(err, "create game")
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
		return nil, WrapError(err, "get game")
	}

	game.Key = strings.TrimSpace(input.Key)
	game.Name = strings.TrimSpace(input.Name)
	game.Category = strings.TrimSpace(input.Category)
	game.IconURL = strings.TrimSpace(input.IconURL)
	game.Description = strings.TrimSpace(input.Description)

	if err := s.games.Update(ctx, game); err != nil {
		return nil, WrapError(err, "update game")
	}

	s.invalidateCache(ctx, cacheKeyGames)
	// audit
	s.appendLogAsync(ctx, string(model.OpEntityGame), game.ID, string(model.OpActionUpdate), map[string]any{"key": game.Key})

	return game, nil
}

// DeleteGame 删除游戏。
func (s *AdminService) DeleteGame(ctx context.Context, id uint64) error {
	if err := s.games.Delete(ctx, id); err != nil {
		return WrapError(err, "delete game")
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

// UserStatsResponse 用户统计信息响应
type UserStatsResponse struct {
	Total               int            `json:"total"`
	ByRole              map[string]int `json:"byRole"`
	ByStatus            map[string]int `json:"byStatus"`
	RecentRegistrations int            `json:"recentRegistrations"`
}

// GetUserStats 获取用户统计信息
func (s *AdminService) GetUserStats(ctx context.Context) (*UserStatsResponse, error) {
	// 获取用户总数
	total, err := s.users.Count(ctx, repository.UserListOptions{})
	if err != nil {
		return nil, WrapError(err, "count total users")
	}

	// 按角色统计
	byRole := make(map[string]int)
	for _, role := range []model.Role{model.RoleUser, model.RolePlayer, model.RoleAdmin} {
		count, err := s.users.Count(ctx, repository.UserListOptions{
			Roles: []model.Role{role},
		})
		if err != nil {
			return nil, WrapError(err, "count users by role")
		}
		byRole[string(role)] = count
	}

	// 按状态统计
	byStatus := make(map[string]int)
	for _, status := range []model.UserStatus{model.UserStatusActive, model.UserStatusBanned, model.UserStatusSuspended} {
		count, err := s.users.Count(ctx, repository.UserListOptions{
			Statuses: []model.UserStatus{status},
		})
		if err != nil {
			return nil, WrapError(err, "count users by status")
		}
		byStatus[string(status)] = count
	}

	// 最近7天注册数
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	recentCount, err := s.users.Count(ctx, repository.UserListOptions{
		DateFrom: &sevenDaysAgo,
	})
	if err != nil {
		return nil, WrapError(err, "count recent registrations")
	}

	return &UserStatsResponse{
		Total:               total,
		ByRole:              byRole,
		ByStatus:            byStatus,
		RecentRegistrations: recentCount,
	}, nil
}

// GetUser 返回指定用户。
func (s *AdminService) GetUser(ctx context.Context, id uint64) (*model.User, error) {
	user, err := s.users.Get(ctx, id)
	if err != nil {
		return nil, mapUserError(err)
	}
	return user, nil
}

// GetUsersByIDs 批量获取用户信息
func (s *AdminService) GetUsersByIDs(ctx context.Context, ids []uint64) ([]model.User, error) {
	if len(ids) == 0 {
		return []model.User{}, nil
	}
	return s.users.GetByIDs(ctx, ids)
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

	// 为新用户创建钱包记录（初始余额为0）
	wallet := &model.Wallet{
		UserID:       user.ID,
		BalanceCents: 0,
		FrozenCents:  0,
	}
	if err := s.wallets.Save(ctx, wallet); err != nil {
		slog.Warn("failed to create wallet for new user", slog.Uint64("user_id", user.ID), slog.String("error", err.Error()))
		// 不中断流程，继续执行
	}

	// 同步 user.Role 到 user_roles 表
	if err := s.syncUserRoleToTable(ctx, user.ID, user.Role); err != nil {
		slog.Warn("failed to sync user_role to table", slog.Uint64("user_id", user.ID), slog.String("error", err.Error()))
		// 不中断流程，继续执行
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
		return nil, WrapError(err, "update user")
	}

	// 同步 user.Role 到 user_roles 表
	if err := s.syncUserRoleToTable(ctx, user.ID, user.Role); err != nil {
		slog.Warn("failed to sync user_role to table", slog.Uint64("user_id", user.ID), slog.String("error", err.Error()))
		// 不中断流程，继续执行
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

	// 同步 user.Role 到 user_roles 表
	if err := s.syncUserRoleToTable(ctx, user.ID, user.Role); err != nil {
		slog.Warn("failed to sync user_role to table", slog.Uint64("user_id", user.ID), slog.String("error", err.Error()))
		// 不中断流程，继续执行
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

// syncUserRoleToTable 同步 user.Role 到 user_roles 多对多表。
// 根据 user.Role 字段的值，在 user_roles 表中创建对应的关联记录。
func (s *AdminService) syncUserRoleToTable(ctx context.Context, userID uint64, role model.Role) error {
	// 根据 role 字段查找对应的 RoleModel
	var roleSlug string
	switch role {
	case model.RoleAdmin:
		roleSlug = string(model.RoleSlugAdmin)
	case model.RolePlayer:
		roleSlug = string(model.RoleSlugPlayer)
	case model.RoleUser:
		roleSlug = string(model.RoleSlugUser)
	default:
		// 未知角色，记录日志但不报错
		slog.Warn("unknown user role, skipping user_roles sync", slog.String("role", string(role)), slog.Uint64("user_id", userID))
		return nil
	}

	roleModel, err := s.roles.GetBySlug(ctx, roleSlug)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			slog.Warn("role not found in database, skipping user_roles sync", slog.String("slug", roleSlug), slog.Uint64("user_id", userID))
			return nil
		}
		return err
	}

	// 为用户分配该角色（替换现有所有角色，保持与 user.Role 字段一致）
	if err := s.roles.AssignToUser(ctx, userID, []uint64{roleModel.ID}); err != nil {
		return err
	}

	slog.Info("user_role_synced_to_table", slog.Uint64("user_id", userID), slog.String("role", string(role)), slog.Uint64("role_id", roleModel.ID))
	return nil
}

// validPassword 验证密码强度
// ✅ 密码安全修复: 增强密码复杂度要求
// 要求:
// - 最小长度8位 (原来6位)
// - 必须包含大写字母
// - 必须包含小写字母
// - 必须包含数字
// - 必须包含特殊字符 (!@#$%^&*()_+-=[]{}|;:,.<>?)
func validPassword(pw string) bool {
	// 最小长度检查: 8位
	if len(pw) < 8 {
		return false
	}

	// 最大长度检查: 防止DoS攻击
	if len(pw) > 128 {
		return false
	}

	// 字符类型计数
	var (
		hasUppercase = false
		hasLowercase = false
		hasDigit     = false
		hasSpecial   = false
	)

	// 允许的特殊字符
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"

	for _, r := range pw {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUppercase = true
		case r >= 'a' && r <= 'z':
			hasLowercase = true
		case r >= '0' && r <= '9':
			hasDigit = true
		case containsRune(specialChars, r):
			hasSpecial = true
		}
	}

	// 必须同时满足所有条件
	return hasUppercase && hasLowercase && hasDigit && hasSpecial
}

// containsRune 检查字符串是否包含指定字符
func containsRune(s string, r rune) bool {
	for _, c := range s {
		if c == r {
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

// UpdateVerificationInput 审核陪玩师请求参数
type UpdateVerificationInput struct {
	Nickname           string
	Bio                string
	HourlyRateCents    int64
	MainGameID         uint64
	VerificationStatus model.VerificationStatus
	VerifiedBy         uint64 // 审核人ID
	Remark             string // 审核备注
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

// ListPlayersPagedWithFilter 返回带筛选的分页陪玩列表。
func (s *AdminService) ListPlayersPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, status *model.VerificationStatus) ([]model.Player, *model.Pagination, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	items, total, err := s.players.ListPagedWithFilter(ctx, page, pageSize, keyword, status)
	if err != nil {
		return nil, nil, err
	}
	p := buildPagination(page, pageSize, total)
	return items, &p, nil
}

// BatchUpdatePlayerStatus 批量更新陪玩师状态。
func (s *AdminService) BatchUpdatePlayerStatus(ctx context.Context, ids []uint64, status model.VerificationStatus) (int64, error) {
	if len(ids) == 0 {
		return 0, apierr.BadRequest("no player ids provided")
	}
	updated, err := s.players.BatchUpdateStatus(ctx, ids, status)
	if err != nil {
		return 0, WrapError(err, "batch update player status")
	}
	s.invalidateCache(ctx, cacheKeyPlayers)
	return updated, nil
}

// BatchDeletePlayers 批量删除陪玩师。
func (s *AdminService) BatchDeletePlayers(ctx context.Context, ids []uint64) (int64, error) {
	if len(ids) == 0 {
		return 0, apierr.BadRequest("no player ids provided")
	}
	deleted, err := s.players.BatchDelete(ctx, ids)
	if err != nil {
		return 0, WrapError(err, "batch delete players")
	}
	s.invalidateCache(ctx, cacheKeyPlayers)
	return deleted, nil
}

// GetPlayer 返回陪玩详情。
func (s *AdminService) GetPlayer(ctx context.Context, id uint64) (*model.Player, error) {
	player, err := s.players.Get(ctx, id)
	if err != nil {
		return nil, WrapError(err, "get player")
	}
	return player, nil
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

// UpdatePlayerVerification 审核陪玩师（保存审核记录）
func (s *AdminService) UpdatePlayerVerification(ctx context.Context, id uint64, input UpdateVerificationInput) (*model.Player, error) {
	player, err := s.players.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := validatePlayerInput(player.UserID, input.VerificationStatus); err != nil {
		return nil, err
	}

	// 保留原有信息
	player.Nickname = strings.TrimSpace(input.Nickname)
	player.Bio = strings.TrimSpace(input.Bio)
	player.HourlyRateCents = input.HourlyRateCents
	player.MainGameID = input.MainGameID

	// 更新审核状态和记录
	oldStatus := player.VerificationStatus
	player.VerificationStatus = input.VerificationStatus

	// 记录审核信息
	now := time.Now()
	player.VerifiedAt = &now
	player.VerifiedBy = &input.VerifiedBy
	player.VerifyRemark = strings.TrimSpace(input.Remark)

	// 如果是拒绝，保存拒绝原因
	if input.VerificationStatus == model.VerificationRejected {
		player.RejectReason = strings.TrimSpace(input.Remark)
	} else {
		player.RejectReason = ""
	}

	if err := s.players.Update(ctx, player); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, cacheKeyPlayers)

	// 审计日志
	s.appendLogAsync(ctx, string(model.OpEntityPlayer), player.ID, "verify", map[string]any{
		"old_status":  oldStatus,
		"new_status":  input.VerificationStatus,
		"verified_by": input.VerifiedBy,
		"remark":      input.Remark,
	})

	// TODO: 发送通知给陪玩师
	// s.sendVerificationNotification(ctx, player, input.VerificationStatus)

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
	UserID          uint64
	PlayerID        *uint64
	GameID          uint64
	ItemID          uint64 // 服务项目ID (必填)
	Title           string
	Description     string
	TotalPriceCents int64
	Currency        model.Currency
	ScheduledStart  *time.Time
	ScheduledEnd    *time.Time
}

// CreateOrder 新建订单，默认状态为 pending。
func (s *AdminService) CreateOrder(ctx context.Context, in CreateOrderInput) (*model.Order, error) {
	// ✅ 数据安全修复: 验证所有必填字段，包括ItemID
	if in.UserID == 0 || in.GameID == 0 || in.ItemID == 0 || in.TotalPriceCents < 0 || !model.IsValidCurrency(in.Currency) {
		return nil, ErrValidation
	}
	if in.ScheduledStart != nil && in.ScheduledEnd != nil && in.ScheduledEnd.Before(*in.ScheduledStart) {
		return nil, ErrValidation
	}

	// 验证服务项目是否存在
	serviceItem, err := s.serviceItems.Get(ctx, in.ItemID)
	if err != nil {
		return nil, apierr.BadRequest("服务项目不存在")
	}

	// 验证服务项目是否激活
	if !serviceItem.IsActive {
		return nil, apierr.BadRequest("服务项目已停用")
	}

	// 可选: 验证服务项目与游戏的关联性
	if serviceItem.GameID != nil && *serviceItem.GameID != in.GameID {
		return nil, apierr.BadRequest("服务项目与游戏不匹配")
	}

	// 验证陪玩师是否存在
	if in.PlayerID != nil && *in.PlayerID != 0 {
		if _, err := s.players.Get(ctx, *in.PlayerID); err != nil {
			return nil, err
		}
	}

	gameID := in.GameID
	order := &model.Order{
		OrderNo:         model.GenerateEscortOrderNo(),
		UserID:          in.UserID,
		ItemID:          in.ItemID, // ✅ 修复: 使用传入的ItemID而不是硬编码
		GameID:          &gameID,
		Quantity:        1,
		UnitPriceCents:  in.TotalPriceCents,
		TotalPriceCents: in.TotalPriceCents,
		Currency:        in.Currency,
		Status:          model.OrderStatusPending,
		Title:           strings.TrimSpace(in.Title),
		Description:     strings.TrimSpace(in.Description),
		ScheduledStart:  in.ScheduledStart,
		ScheduledEnd:    in.ScheduledEnd,
	}
	if in.PlayerID != nil {
		order.PlayerID = in.PlayerID
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
	order.SetPlayerID(playerID)
	if err := s.orders.Update(ctx, order); err != nil {
		return nil, WrapError(err, "update order")
	}
	s.invalidateCache(ctx, cacheKeyOrders)
	s.appendLogAsync(ctx, string(model.OpEntityOrder), order.ID, string(model.OpActionAssignPlayer), map[string]any{"player_id": playerID})
	return order, nil
}

// UpdateOrderInput 用于更新订单状态。
type UpdateOrderInput struct {
	Status            model.OrderStatus
	TotalPriceCents   int64
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
func (s *AdminService) ListOrders(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, *model.Pagination, error) {
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
	order, err := s.orders.Get(ctx, id)
	if err != nil {
		return nil, WrapError(err, "get order")
	}
	return order, nil
}

// UpdateOrder 更新订单信息。
func (s *AdminService) UpdateOrder(ctx context.Context, id uint64, input UpdateOrderInput) (*model.Order, error) {
	order, err := s.orders.Get(ctx, id)
	if err != nil {
		return nil, WrapError(err, "get order")
	}

	if !isValidOrderStatus(input.Status) {
		return nil, ErrValidation
	}
	if !model.IsValidCurrency(input.Currency) {
		return nil, ErrValidation
	}
	if input.TotalPriceCents < 0 {
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
	order.TotalPriceCents = input.TotalPriceCents
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
	switch order.Status {
	case model.OrderStatusCanceled:
		action = model.OpActionCancel
	case model.OrderStatusRefunded:
		action = model.OpActionRefund
	default:
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
		Status:          model.OrderStatusConfirmed,
		TotalPriceCents: order.TotalPriceCents,
		Currency:        order.Currency,
		ScheduledStart:  order.ScheduledStart,
		ScheduledEnd:    order.ScheduledEnd,
		CancelReason:    order.CancelReason,
		StartedAt:       order.StartedAt,
		CompletedAt:     order.CompletedAt,
		RefundReason:    order.RefundReason,
		RefundedAt:      order.RefundedAt,
		Note:            note,
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
		Status:          model.OrderStatusInProgress,
		TotalPriceCents: order.TotalPriceCents,
		Currency:        order.Currency,
		ScheduledStart:  order.ScheduledStart,
		ScheduledEnd:    order.ScheduledEnd,
		CancelReason:    order.CancelReason,
		StartedAt:       &startedAt,
		CompletedAt:     order.CompletedAt,
		RefundReason:    order.RefundReason,
		RefundedAt:      order.RefundedAt,
		Note:            note,
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
		Status:          model.OrderStatusCompleted,
		TotalPriceCents: order.TotalPriceCents,
		Currency:        order.Currency,
		ScheduledStart:  order.ScheduledStart,
		ScheduledEnd:    order.ScheduledEnd,
		CancelReason:    order.CancelReason,
		StartedAt:       order.StartedAt,
		CompletedAt:     &completedAt,
		RefundReason:    order.RefundReason,
		RefundedAt:      order.RefundedAt,
		Note:            note,
	})
}

// RefundOrder 执行退款并记录退款信息。
func (s *AdminService) RefundOrder(ctx context.Context, id uint64, input RefundOrderInput) (*model.Order, error) {
	order, err := s.orders.Get(ctx, id)
	if err != nil {
		return nil, WrapError(err, "get order")
	}
	reason := strings.TrimSpace(input.Reason)
	if reason == "" {
		return nil, apierr.BadRequest("reason is required")
	}
	switch order.Status {
	case model.OrderStatusCompleted, model.OrderStatusInProgress, model.OrderStatusConfirmed:
		// allowed
	default:
		return nil, apierr.BadRequest("invalid order status for refund")
	}
	amount := order.TotalPriceCents
	if input.AmountCents != nil {
		if *input.AmountCents <= 0 || *input.AmountCents > order.TotalPriceCents {
			return nil, apierr.BadRequest("invalid refund amount")
		}
		amount = *input.AmountCents
	}
	refundedAt := time.Now().UTC()
	note := strings.TrimSpace(input.Note)
	updatedOrder, err := s.UpdateOrder(ctx, id, UpdateOrderInput{
		Status:            model.OrderStatusRefunded,
		TotalPriceCents:   order.TotalPriceCents,
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
		return nil, WrapError(err, "update order")
	}

	// 更新相关支付为已退款状态（若存在）
	payments, err := s.listPaymentsByOrder(ctx, id)
	if err != nil {
		return nil, WrapError(err, "list payments by order")
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
				return nil, WrapError(updErr, "update payment")
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
		return WrapError(err, "delete order")
	}
	s.invalidateCache(ctx, cacheKeyOrders)
	s.appendLogAsync(ctx, string(model.OpEntityOrder), id, string(model.OpActionDelete), nil)
	return nil
}

// --- Payment management ---

// UpdatePaymentInput 调整支付状态。
type UpdatePaymentInput struct {
	Status              model.PaymentStatus
	ProviderTradeNo     string
	ProviderRaw         json.RawMessage
	PaidAt              *time.Time
	RefundedAt          *time.Time
	RefundedAmountCents *int64 // 已退款金额（分）
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
		return nil, WrapError(err, "get order")
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
		return nil, WrapError(err, "create payment")
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
		return nil, WrapError(err, "get payment")
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
	payment, err := s.payments.Get(ctx, id)
	if err != nil {
		return nil, WrapError(err, "get payment")
	}
	return payment, nil
}

// GetPaymentWithRelations 获取支付详情及关联的订单和用户信息。
func (s *AdminService) GetPaymentWithRelations(ctx context.Context, id uint64) (*model.Payment, error) {
	payment, err := s.payments.GetWithRelations(ctx, id)
	if err != nil {
		return nil, WrapError(err, "get payment with relations")
	}
	return payment, nil
}

// GetPaymentsByOrderID 根据订单ID获取所有支付记录。
func (s *AdminService) GetPaymentsByOrderID(ctx context.Context, orderID uint64) ([]model.Payment, error) {
	payments, err := s.payments.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, WrapError(err, "get payments by order id")
	}
	return payments, nil
}

// UpdatePayment 更新支付状态。
func (s *AdminService) UpdatePayment(ctx context.Context, id uint64, input UpdatePaymentInput) (*model.Payment, error) {
	payment, err := s.payments.Get(ctx, id)
	if err != nil {
		return nil, WrapError(err, "get payment")
	}

	if !isValidPaymentStatus(input.Status) {
		return nil, apierr.BadRequest("invalid payment status")
	}

	if !isAllowedPaymentTransition(payment.Status, input.Status) {
		return nil, apierr.BadRequest("invalid payment status transition")
	}

	payment.Status = input.Status
	payment.ProviderTradeNo = strings.TrimSpace(input.ProviderTradeNo)
	payment.ProviderRaw = input.ProviderRaw
	payment.PaidAt = input.PaidAt
	payment.RefundedAt = input.RefundedAt

	if err := s.payments.Update(ctx, payment); err != nil {
		return nil, WrapError(err, "update payment")
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

// UpdatePaymentWithRefund processes a refund with amount validation and logging.
// Requirements: 2.1, 2.2, 2.3, 2.4, 9.1, 9.2, 9.3
func (s *AdminService) UpdatePaymentWithRefund(ctx context.Context, id uint64, input UpdatePaymentInput, refundAmount int64, reason string, operatorID *uint64) (*model.Payment, error) {
	payment, err := s.payments.Get(ctx, id)
	if err != nil {
		return nil, WrapError(err, "get payment")
	}

	// Validate refund amount
	if err := payment.ValidateRefundAmount(refundAmount); err != nil {
		return nil, apierr.BadRequest(err.Error())
	}

	// Validate status transition if status is changing
	if input.Status != "" && input.Status != payment.Status {
		if !isAllowedPaymentTransition(payment.Status, input.Status) {
			return nil, apierr.BadRequest("invalid payment status transition")
		}
		payment.Status = input.Status
	}

	// Update refunded amount
	payment.RefundedAmountCents += refundAmount
	payment.ProviderTradeNo = strings.TrimSpace(input.ProviderTradeNo)
	payment.ProviderRaw = input.ProviderRaw
	payment.RefundedAt = input.RefundedAt

	// Check if fully refunded and update status
	if payment.IsFullyRefunded() && payment.Status != model.PaymentStatusRefunded {
		payment.Status = model.PaymentStatusRefunded
	}

	if err := s.payments.Update(ctx, payment); err != nil {
		return nil, WrapError(err, "update payment")
	}

	s.invalidateCache(ctx, cacheKeyPayments)

	// Log the refund operation with detailed metadata
	s.appendLogAsync(ctx, string(model.OpEntityPayment), payment.ID, string(model.OpActionRefund), map[string]any{
		"refund_amount_cents":  refundAmount,
		"total_refunded_cents": payment.RefundedAmountCents,
		"remaining_cents":      payment.RemainingRefundableAmount(),
		"reason":               reason,
		"is_full_refund":       payment.IsFullyRefunded(),
		"operator_id":          operatorID,
		"status":               payment.Status,
	})

	if rid, ok := logging.RequestIDFromContext(ctx); ok {
		slog.Info("payment_refunded",
			slog.Uint64("payment_id", payment.ID),
			slog.Int64("refund_amount", refundAmount),
			slog.Int64("total_refunded", payment.RefundedAmountCents),
			slog.String("status", string(payment.Status)),
			slog.String("request_id", rid))
	} else {
		slog.Info("payment_refunded",
			slog.Uint64("payment_id", payment.ID),
			slog.Int64("refund_amount", refundAmount),
			slog.Int64("total_refunded", payment.RefundedAmountCents),
			slog.String("status", string(payment.Status)))
	}

	return payment, nil
}

// DeletePayment 删除支付记录。
func (s *AdminService) DeletePayment(ctx context.Context, id uint64) error {
	if err := s.payments.Delete(ctx, id); err != nil {
		return WrapError(err, "delete payment")
	}
	s.invalidateCache(ctx, cacheKeyPayments)
	s.appendLogAsync(ctx, string(model.OpEntityPayment), id, string(model.OpActionDelete), nil)
	return nil
}

// ============================================================================
// Batch Payment Operations
// ============================================================================

// BatchCaptureResult 批量收款操作结果
type BatchCaptureResult struct {
	SuccessCount int                   `json:"successCount"`
	FailedCount  int                   `json:"failedCount"`
	FailedIDs    []uint64              `json:"failedIds,omitempty"`
	Errors       []BatchOperationError `json:"errors,omitempty"`
}

// BatchOperationError 批量操作错误详情
type BatchOperationError struct {
	PaymentID uint64 `json:"paymentId"`
	Message   string `json:"message"`
}

// BatchCaptureRequest 批量收款请求
type BatchCaptureRequest struct {
	PaymentIDs      []uint64   `json:"paymentIds" binding:"required,min=1,max=500"`
	ProviderTradeNo string     `json:"providerTradeNo,omitempty"`
	PaidAt          *time.Time `json:"paidAt,omitempty"`
}

// BatchCapture 批量收款 - 将多个pending状态的支付标记为已支付
// 业务规则：
// 1. 只能处理pending状态的支付
// 2. 支付会设置为paid状态
// 3. 返回成功/失败计数及错误详情
func (s *AdminService) BatchCapture(ctx context.Context, req BatchCaptureRequest) (*BatchCaptureResult, error) {
	if len(req.PaymentIDs) == 0 {
		return nil, apierr.BadRequest("payment ids cannot be empty")
	}
	if len(req.PaymentIDs) > 500 {
		return nil, apierr.BadRequest("maximum 500 payments allowed per batch")
	}

	result := &BatchCaptureResult{
		FailedIDs: make([]uint64, 0),
		Errors:    make([]BatchOperationError, 0),
	}

	paidAt := req.PaidAt
	if paidAt == nil {
		now := time.Now().UTC()
		paidAt = &now
	}

	for _, paymentID := range req.PaymentIDs {
		payment, err := s.payments.Get(ctx, paymentID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				result.FailedCount++
				result.FailedIDs = append(result.FailedIDs, paymentID)
				result.Errors = append(result.Errors, BatchOperationError{
					PaymentID: paymentID,
					Message:   "payment not found",
				})
				continue
			}
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("failed to get payment: %v", err),
			})
			continue
		}

		// 验证状态：只能capture pending状态的支付
		if payment.Status != model.PaymentStatusPending {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("invalid status for capture: %s (expected: pending)", payment.Status),
			})
			continue
		}

		// 更新支付状态
		payment.Status = model.PaymentStatusPaid
		payment.PaidAt = paidAt
		if req.ProviderTradeNo != "" {
			payment.ProviderTradeNo = strings.TrimSpace(req.ProviderTradeNo)
		} else {
			payment.ProviderTradeNo = fmt.Sprintf("batch_capture_%d_%d", paymentID, time.Now().Unix())
		}

		if err := s.payments.Update(ctx, payment); err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("failed to update payment: %v", err),
			})
			continue
		}

		result.SuccessCount++

		// 异步记录日志
		s.appendLogAsync(ctx, string(model.OpEntityPayment), payment.ID, string(model.OpActionCapture), map[string]any{
			"batch_operation": true,
			"trade_no":        payment.ProviderTradeNo,
		})
	}

	s.invalidateCache(ctx, cacheKeyPayments)
	return result, nil
}

// BatchRefundRequest 批量退款请求
type BatchRefundRequest struct {
	PaymentIDs []uint64   `json:"paymentIds" binding:"required,min=1,max=500"`
	Reason     string     `json:"reason" binding:"required,max=500"`
	RefundedAt *time.Time `json:"refundedAt,omitempty"`
}

// BatchRefund 批量退款 - 退款多个已支付的支付
// 业务规则：
// 1. 只能退款paid状态的支付
// 2. 全额退款，状态更新为refunded
// 3. 订单状态也会更新为refunded
func (s *AdminService) BatchRefund(ctx context.Context, req BatchRefundRequest) (*BatchCaptureResult, error) {
	if len(req.PaymentIDs) == 0 {
		return nil, apierr.BadRequest("payment ids cannot be empty")
	}
	if len(req.PaymentIDs) > 500 {
		return nil, apierr.BadRequest("maximum 500 payments allowed per batch")
	}
	if strings.TrimSpace(req.Reason) == "" {
		return nil, apierr.BadRequest("refund reason is required")
	}

	result := &BatchCaptureResult{
		FailedIDs: make([]uint64, 0),
		Errors:    make([]BatchOperationError, 0),
	}

	refundedAt := req.RefundedAt
	if refundedAt == nil {
		now := time.Now().UTC()
		refundedAt = &now
	}

	for _, paymentID := range req.PaymentIDs {
		payment, err := s.payments.Get(ctx, paymentID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				result.FailedCount++
				result.FailedIDs = append(result.FailedIDs, paymentID)
				result.Errors = append(result.Errors, BatchOperationError{
					PaymentID: paymentID,
					Message:   "payment not found",
				})
				continue
			}
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("failed to get payment: %v", err),
			})
			continue
		}

		// 验证状态：只能退款paid状态的支付
		if payment.Status != model.PaymentStatusPaid {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("invalid status for refund: %s (expected: paid)", payment.Status),
			})
			continue
		}

		// 检查是否已经全额退款
		if payment.IsFullyRefunded() {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   "payment is already fully refunded",
			})
			continue
		}

		// 更新支付状态
		payment.Status = model.PaymentStatusRefunded
		payment.RefundedAt = refundedAt
		payment.RefundedAmountCents = payment.AmountCents
		payment.ProviderTradeNo = fmt.Sprintf("batch_refund_%d_%d", paymentID, time.Now().Unix())

		if err := s.payments.Update(ctx, payment); err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("failed to update payment: %v", err),
			})
			continue
		}

		// 更新关联订单状态
		order, err := s.orders.Get(ctx, payment.OrderID)
		if err == nil {
			order.Status = model.OrderStatusRefunded
			order.RefundAmountCents = payment.AmountCents
			order.RefundReason = req.Reason
			order.RefundedAt = refundedAt
			_ = s.orders.Update(ctx, order)
		}

		result.SuccessCount++

		// 异步记录日志
		s.appendLogAsync(ctx, string(model.OpEntityPayment), payment.ID, string(model.OpActionRefund), map[string]any{
			"batch_operation":     true,
			"refund_amount_cents": payment.AmountCents,
			"reason":              req.Reason,
		})
	}

	s.invalidateCache(ctx, cacheKeyPayments)
	return result, nil
}

// BatchCancelRequest 批量取消支付请求
type BatchCancelRequest struct {
	PaymentIDs []uint64 `json:"paymentIds" binding:"required,min=1,max=500"`
}

// BatchCancel 批量取消支付 - 取消多个pending状态的支付
// 业务规则：
// 1. 只能取消pending状态的支付
// 2. 支付状态更新为failed
func (s *AdminService) BatchCancel(ctx context.Context, req BatchCancelRequest) (*BatchCaptureResult, error) {
	if len(req.PaymentIDs) == 0 {
		return nil, apierr.BadRequest("payment ids cannot be empty")
	}
	if len(req.PaymentIDs) > 500 {
		return nil, apierr.BadRequest("maximum 500 payments allowed per batch")
	}

	result := &BatchCaptureResult{
		FailedIDs: make([]uint64, 0),
		Errors:    make([]BatchOperationError, 0),
	}

	for _, paymentID := range req.PaymentIDs {
		payment, err := s.payments.Get(ctx, paymentID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				result.FailedCount++
				result.FailedIDs = append(result.FailedIDs, paymentID)
				result.Errors = append(result.Errors, BatchOperationError{
					PaymentID: paymentID,
					Message:   "payment not found",
				})
				continue
			}
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("failed to get payment: %v", err),
			})
			continue
		}

		// 验证状态：只能取消pending状态的支付
		if payment.Status != model.PaymentStatusPending {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("invalid status for cancel: %s (expected: pending)", payment.Status),
			})
			continue
		}

		// 更新支付状态为failed（表示已取消）
		payment.Status = model.PaymentStatusFailed

		if err := s.payments.Update(ctx, payment); err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("failed to update payment: %v", err),
			})
			continue
		}

		result.SuccessCount++

		// 异步记录日志
		s.appendLogAsync(ctx, string(model.OpEntityPayment), payment.ID, string(model.OpActionCancel), map[string]any{
			"batch_operation": true,
		})
	}

	s.invalidateCache(ctx, cacheKeyPayments)
	return result, nil
}

// BatchUpdateStatusRequest 批量更新支付状态请求
type BatchUpdateStatusRequest struct {
	PaymentIDs []uint64            `json:"paymentIds" binding:"required,min=1,max=500"`
	Status     model.PaymentStatus `json:"-" binding:"-"` // Not used for binding, set from handler
}

// BatchUpdateStatus 批量更新支付状态
// 业务规则：
// 1. 验证状态转换是否有效
// 2. 只允许有效的状态转换
func (s *AdminService) BatchUpdateStatus(ctx context.Context, req BatchUpdateStatusRequest) (*BatchCaptureResult, error) {
	if len(req.PaymentIDs) == 0 {
		return nil, apierr.BadRequest("payment ids cannot be empty")
	}
	if len(req.PaymentIDs) > 500 {
		return nil, apierr.BadRequest("maximum 500 payments allowed per batch")
	}

	// 验证状态是否有效
	if !isValidPaymentStatus(req.Status) {
		return nil, apierr.BadRequest("invalid payment status")
	}

	result := &BatchCaptureResult{
		FailedIDs: make([]uint64, 0),
		Errors:    make([]BatchOperationError, 0),
	}

	for _, paymentID := range req.PaymentIDs {
		payment, err := s.payments.Get(ctx, paymentID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				result.FailedCount++
				result.FailedIDs = append(result.FailedIDs, paymentID)
				result.Errors = append(result.Errors, BatchOperationError{
					PaymentID: paymentID,
					Message:   "payment not found",
				})
				continue
			}
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("failed to get payment: %v", err),
			})
			continue
		}

		// 验证状态转换
		if !isAllowedPaymentTransition(payment.Status, req.Status) {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("invalid status transition from %s to %s", payment.Status, req.Status),
			})
			continue
		}

		// 如果是转换为paid，设置PaidAt
		if req.Status == model.PaymentStatusPaid && payment.PaidAt == nil {
			now := time.Now().UTC()
			payment.PaidAt = &now
		}

		// 如果是转换为refunded，设置RefundedAt
		if req.Status == model.PaymentStatusRefunded && payment.RefundedAt == nil {
			now := time.Now().UTC()
			payment.RefundedAt = &now
			payment.RefundedAmountCents = payment.AmountCents
		}

		// 更新支付状态
		payment.Status = req.Status

		if err := s.payments.Update(ctx, payment); err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, paymentID)
			result.Errors = append(result.Errors, BatchOperationError{
				PaymentID: paymentID,
				Message:   fmt.Sprintf("failed to update payment: %v", err),
			})
			continue
		}

		result.SuccessCount++

		// 异步记录日志
		action := model.OpActionUpdateStatus
		if req.Status == model.PaymentStatusPaid {
			action = model.OpActionCapture
		} else if req.Status == model.PaymentStatusRefunded {
			action = model.OpActionRefund
		} else if req.Status == model.PaymentStatusFailed {
			action = model.OpActionCancel
		}
		s.appendLogAsync(ctx, string(model.OpEntityPayment), payment.ID, string(action), map[string]any{
			"batch_operation": true,
			"old_status":      string(payment.Status),
			"new_status":      string(req.Status),
		})
	}

	s.invalidateCache(ctx, cacheKeyPayments)
	return result, nil
}

// GetPaymentLogs returns operation logs for a payment.
// Requirements: 2.5
func (s *AdminService) GetPaymentLogs(ctx context.Context, paymentID uint64, opts repository.OperationLogListOptions) ([]model.OperationLog, int64, error) {
	if s.tx == nil {
		return nil, 0, apierr.InternalError("transaction manager not configured")
	}

	var logs []model.OperationLog
	var total int64

	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		var err error
		logs, total, err = r.OpLogs.ListByEntity(ctx, string(model.OpEntityPayment), paymentID, opts)
		return err
	})

	if err != nil {
		return nil, 0, WrapError(err, "get payment logs")
	}

	return logs, total, nil
}

// appendLogAsync 追加操作日志（尽力而为，不影响主流程）。
func (s *AdminService) appendLogAsync(ctx context.Context, entity string, id uint64, action string, meta map[string]any) {
	if s.tx == nil {
		return
	}
	_ = s.tx.WithTx(ctx, func(r *common.Repos) error {
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
	return WrapError(err, "操作用户数据失败")
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
		return nil, nil, apierr.InternalError("事务管理器未配置")
	}
	var logs []model.OperationLog
	var total int64
	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
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

// SearchOperationLogs 搜索操作日志（支持跨实体搜索）
func (s *AdminService) SearchOperationLogs(ctx context.Context, opts repository.OperationLogSearchOptions) ([]model.OperationLog, *model.Pagination, error) {
	if s.tx == nil {
		return nil, nil, apierr.InternalError("事务管理器未配置")
	}
	var logs []model.OperationLog
	var total int64
	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		norm := repository.OperationLogSearchOptions{
			Page:        repository.NormalizePage(opts.Page),
			PageSize:    repository.NormalizePageSize(opts.PageSize),
			EntityType:  opts.EntityType,
			EntityID:    opts.EntityID,
			Action:      opts.Action,
			ActorUserID: opts.ActorUserID,
			DateFrom:    opts.DateFrom,
			DateTo:      opts.DateTo,
		}
		items, cnt, err := r.OpLogs.List(ctx, norm)
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
		return nil, nil, apierr.InternalError("事务管理器未配置")
	}
	var items []model.Review
	var total int64
	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
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
		return nil, apierr.InternalError("事务管理器未配置")
	}
	var item *model.Review
	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		var err error
		item, err = r.Reviews.Get(ctx, id)
		return err
	})
	if err != nil {
		return nil, WrapError(err, "get review")
	}
	return item, nil
}

// CreateReview 新建评价。
func (s *AdminService) CreateReview(ctx context.Context, r model.Review) (*model.Review, error) {
	if !r.Score.Valid() || r.OrderID == 0 || r.UserID == 0 || r.PlayerID == 0 {
		return nil, ErrValidation
	}
	if s.tx == nil {
		return nil, apierr.InternalError("事务管理器未配置")
	}
	err := s.tx.WithTx(ctx, func(txr *common.Repos) error { return txr.Reviews.Create(ctx, &r) })
	if err != nil {
		return nil, WrapError(err, "create review")
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
		return nil, apierr.InternalError("事务管理器未配置")
	}
	var item *model.Review
	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		obj, err := r.Reviews.Get(ctx, id)
		if err != nil {
			return WrapError(err, "get review")
		}
		obj.Score = score
		obj.Content = strings.TrimSpace(content)
		if err := r.Reviews.Update(ctx, obj); err != nil {
			return WrapError(err, "update review")
		}
		item = obj
		return nil
	})
	if err != nil {
		return nil, WrapError(err, "update review transaction")
	}
	s.appendLogAsync(ctx, string(model.OpEntityReview), id, string(model.OpActionUpdate), nil)
	return item, nil
}

// DeleteReview 删除评价。
func (s *AdminService) DeleteReview(ctx context.Context, id uint64) error {
	if s.tx == nil {
		return errors.New("transaction manager not configured")
	}
	err := s.tx.WithTx(ctx, func(r *common.Repos) error { return r.Reviews.Delete(ctx, id) })
	return WrapError(err, "delete review")
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

// ReviewReportDTO 举报信息DTO
type ReviewReportDTO struct {
	ID           uint64                   `json:"id"`
	ReviewID     uint64                   `json:"reviewId"`
	ReporterID   uint64                   `json:"reporterId"`
	ReporterName string                   `json:"reporterName"`
	Reason       string                   `json:"reason"`
	Evidence     string                   `json:"evidence,omitempty"`
	Status       model.ReviewReportStatus `json:"status"`
	HandledBy    *uint64                  `json:"handledBy,omitempty"`
	HandlerName  string                   `json:"handlerName,omitempty"`
	HandledAt    *time.Time               `json:"handledAt,omitempty"`
	HandlingNote string                   `json:"handlingNote,omitempty"`
	CreatedAt    time.Time                `json:"createdAt"`
}

// ReportReview 举报评价
func (s *AdminService) ReportReview(ctx context.Context, reviewID, reporterID uint64, reason, evidence string) (uint64, error) {
	if s.tx == nil {
		return 0, apierr.InternalError("事务管理器未配置")
	}

	var reportID uint64
	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		// 验证评价是否存在
		_, err := r.Reviews.Get(ctx, reviewID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}
			return err
		}

		// 创建举报记录
		report := &model.ReviewReport{
			ReviewID:   reviewID,
			ReporterID: reporterID,
			Reason:     reason,
			Evidence:   evidence,
			Status:     model.ReviewReportStatusPending,
		}

		if err := r.ReviewReports.Create(ctx, report); err != nil {
			return err
		}

		reportID = report.ID

		// 标记评价为已举报
		review, err := r.Reviews.Get(ctx, reviewID)
		if err != nil {
			return err
		}
		review.IsReported = true
		return r.Reviews.Update(ctx, review)
	})

	if err != nil {
		return 0, WrapError(err, "report review")
	}

	return reportID, nil
}

// ListReviewReports 列出举报
func (s *AdminService) ListReviewReports(ctx context.Context, page, pageSize int, reviewID, reporterID *uint64, status *model.ReviewReportStatus, dateFrom, dateTo *time.Time) ([]ReviewReportDTO, *model.Pagination, error) {
	if s.tx == nil {
		return nil, nil, apierr.InternalError("事务管理器未配置")
	}

	var reports []model.ReviewReport
	var total int64

	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		var err error
		reports, total, err = r.ReviewReports.List(ctx, repository.ReviewReportListOptions{
			Page:       page,
			PageSize:   pageSize,
			ReviewID:   reviewID,
			ReporterID: reporterID,
			Status:     status,
			DateFrom:   dateFrom,
			DateTo:     dateTo,
		})
		return err
	})

	if err != nil {
		return nil, nil, WrapError(err, "list review reports")
	}

	// 转换为DTO
	reportDTOs := make([]ReviewReportDTO, 0, len(reports))
	for _, report := range reports {
		dto := ReviewReportDTO{
			ID:           report.ID,
			ReviewID:     report.ReviewID,
			ReporterID:   report.ReporterID,
			Reason:       report.Reason,
			Evidence:     report.Evidence,
			Status:       report.Status,
			HandledBy:    report.HandledBy,
			HandledAt:    report.HandledAt,
			HandlingNote: report.HandlingNote,
			CreatedAt:    report.CreatedAt,
		}

		// 获取举报人信息
		if reporter, err := s.users.Get(ctx, report.ReporterID); err == nil {
			dto.ReporterName = reporter.Name
		}

		// 获取处理人信息
		if report.HandledBy != nil {
			if handler, err := s.users.Get(ctx, *report.HandledBy); err == nil {
				dto.HandlerName = handler.Name
			}
		}

		reportDTOs = append(reportDTOs, dto)
	}

	p := &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	}

	return reportDTOs, p, nil
}

// GetReviewReport 获取举报详情
func (s *AdminService) GetReviewReport(ctx context.Context, id uint64) (*ReviewReportDTO, error) {
	if s.tx == nil {
		return nil, apierr.InternalError("事务管理器未配置")
	}

	var report *model.ReviewReport
	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		var err error
		report, err = r.ReviewReports.Get(ctx, id)
		return err
	})

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, WrapError(err, "get review report")
	}

	dto := &ReviewReportDTO{
		ID:           report.ID,
		ReviewID:     report.ReviewID,
		ReporterID:   report.ReporterID,
		Reason:       report.Reason,
		Evidence:     report.Evidence,
		Status:       report.Status,
		HandledBy:    report.HandledBy,
		HandledAt:    report.HandledAt,
		HandlingNote: report.HandlingNote,
		CreatedAt:    report.CreatedAt,
	}

	// 获取举报人信息
	if reporter, err := s.users.Get(ctx, report.ReporterID); err == nil {
		dto.ReporterName = reporter.Name
	}

	// 获取处理人信息
	if report.HandledBy != nil {
		if handler, err := s.users.Get(ctx, *report.HandledBy); err == nil {
			dto.HandlerName = handler.Name
		}
	}

	return dto, nil
}

// HandleReportResponse 处理举报响应
type HandleReportResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// HandleReviewReport 处理举报
func (s *AdminService) HandleReviewReport(ctx context.Context, reportID, handlerID uint64, action, note string) (*HandleReportResponse, error) {
	if s.tx == nil {
		return nil, apierr.InternalError("事务管理器未配置")
	}

	var message string
	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		// 获取举报记录
		report, err := r.ReviewReports.Get(ctx, reportID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}
			return err
		}

		// 检查举报是否已处理
		if report.Status != model.ReviewReportStatusPending {
			return apierr.BadRequest("report already handled")
		}

		// 获取被举报的评价
		review, err := r.Reviews.Get(ctx, report.ReviewID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}
			return err
		}

		now := time.Now()

		switch action {
		case "delete":
			// 删除评价
			oldStatus := review.Status
			review.Status = model.ReviewStatusDeleted
			if err := r.Reviews.Update(ctx, review); err != nil {
				return err
			}

			// 记录删除评价的操作日志
			if r.OpLogs != nil {
				metadata := fmt.Sprintf(`{"old_status":"%s","new_status":"%s","report_id":%d,"action":"delete"}`, oldStatus, model.ReviewStatusDeleted, reportID)
				log := &model.OperationLog{
					EntityType:   string(model.OpEntityReview),
					EntityID:     review.ID,
					ActorUserID:  &handlerID,
					Action:       string(model.OpActionDelete),
					Reason:       fmt.Sprintf("处理举报：%s", note),
					MetadataJSON: []byte(metadata),
				}
				_ = r.OpLogs.Append(ctx, log)
			}

			// 更新举报状态为已通过
			report.Status = model.ReviewReportStatusApproved
			report.HandledBy = &handlerID
			report.HandledAt = &now
			report.HandlingNote = note
			if report.HandlingNote == "" {
				report.HandlingNote = "评价已删除"
			}
			message = "评价已删除"

		case "warn":
			// 警告评价者（保留评价，但标记为已处理）
			report.Status = model.ReviewReportStatusApproved
			report.HandledBy = &handlerID
			report.HandledAt = &now
			report.HandlingNote = note
			if report.HandlingNote == "" {
				report.HandlingNote = "已警告评价者"
			}
			message = "已警告评价者"

		case "reject":
			// 驳回举报
			report.Status = model.ReviewReportStatusRejected
			report.HandledBy = &handlerID
			report.HandledAt = &now
			report.HandlingNote = note
			if report.HandlingNote == "" {
				report.HandlingNote = "举报不成立"
			}
			message = "举报已驳回"

		default:
			return apierr.BadRequest("invalid action")
		}

		// 更新举报记录
		if err := r.ReviewReports.Update(ctx, report); err != nil {
			return err
		}

		// 检查是否还有其他待处理的举报，如果没有则取消评价的举报标记
		pendingStatus := model.ReviewReportStatusPending
		reviewIDPtr := &report.ReviewID
		pendingReports, _, err := r.ReviewReports.List(ctx, repository.ReviewReportListOptions{
			ReviewID: reviewIDPtr,
			Status:   &pendingStatus,
			Page:     1,
			PageSize: 1,
		})
		if err == nil && len(pendingReports) == 0 && review.IsReported {
			review.IsReported = false
			if err := r.Reviews.Update(ctx, review); err != nil {
				// 记录错误但不影响举报处理
				slog.Warn("failed to update review reported status", slog.Any("error", err))
			}
		}

		// 记录处理举报的操作日志
		if r.OpLogs != nil {
			metadata := fmt.Sprintf(`{"report_id":%d,"action":"%s","note":"%s"}`, reportID, action, note)
			log := &model.OperationLog{
				EntityType:   string(model.OpEntityReview),
				EntityID:     review.ID,
				ActorUserID:  &handlerID,
				Action:       string(model.OpActionHandleReport),
				Reason:       fmt.Sprintf("处理举报：%s", message),
				MetadataJSON: []byte(metadata),
			}
			_ = r.OpLogs.Append(ctx, log)
		}

		return nil
	})

	if err != nil {
		return nil, WrapError(err, "handle review report")
	}

	return &HandleReportResponse{
		Status:  "success",
		Message: message,
	}, nil
}

// ListPendingReviews 获取待审核评价列表
func (s *AdminService) ListPendingReviews(ctx context.Context, page, pageSize int) ([]model.Review, int64, error) {
	if s.tx == nil {
		return nil, 0, apierr.InternalError("事务管理器未配置")
	}

	var reviews []model.Review
	var total int64
	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		var err error
		reviews, total, err = r.Reviews.ListPending(ctx, page, pageSize)
		return err
	})

	if err != nil {
		return nil, 0, WrapError(err, "list pending reviews")
	}

	return reviews, total, nil
}

// ApproveReview 批准评价
func (s *AdminService) ApproveReview(ctx context.Context, reviewID uint64, reason string, actorUserID *uint64) error {
	if s.tx == nil {
		return apierr.InternalError("事务管理器未配置")
	}

	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		// 获取评价
		review, err := r.Reviews.Get(ctx, reviewID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}
			return err
		}

		// 检查状态：只有待审核的评价可以批准
		if review.Status != model.ReviewStatusPending {
			return apierr.BadRequest("只能批准待审核的评价")
		}

		oldStatus := review.Status

		// 更新状态为已通过
		if err := r.Reviews.UpdateStatus(ctx, reviewID, model.ReviewStatusApproved, ""); err != nil {
			return err
		}

		// 记录操作日志
		if r.OpLogs != nil {
			metadata := fmt.Sprintf(`{"old_status":"%s","new_status":"%s"}`, oldStatus, model.ReviewStatusApproved)
			log := &model.OperationLog{
				EntityType:   string(model.OpEntityReview),
				EntityID:     reviewID,
				ActorUserID:  actorUserID,
				Action:       string(model.OpActionApprove),
				Reason:       reason,
				MetadataJSON: []byte(metadata),
			}
			_ = r.OpLogs.Append(ctx, log)
		}

		return nil
	})

	if err != nil {
		return WrapError(err, "approve review")
	}

	return nil
}

// RejectReview 拒绝评价
func (s *AdminService) RejectReview(ctx context.Context, reviewID uint64, reason string, actorUserID *uint64) error {
	if s.tx == nil {
		return apierr.InternalError("事务管理器未配置")
	}

	// 验证拒绝原因
	if reason == "" {
		return apierr.BadRequest("拒绝原因不能为空")
	}

	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		// 获取评价
		review, err := r.Reviews.Get(ctx, reviewID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}
			return err
		}

		// 检查状态：只有待审核的评价可以拒绝
		if review.Status != model.ReviewStatusPending {
			return apierr.BadRequest("只能拒绝待审核的评价")
		}

		oldStatus := review.Status

		// 更新状态为已拒绝
		if err := r.Reviews.UpdateStatus(ctx, reviewID, model.ReviewStatusRejected, reason); err != nil {
			return err
		}

		// 记录操作日志
		if r.OpLogs != nil {
			metadata := fmt.Sprintf(`{"old_status":"%s","new_status":"%s","rejection_reason":"%s"}`, oldStatus, model.ReviewStatusRejected, reason)
			log := &model.OperationLog{
				EntityType:   string(model.OpEntityReview),
				EntityID:     reviewID,
				ActorUserID:  actorUserID,
				Action:       string(model.OpActionReject),
				Reason:       reason,
				MetadataJSON: []byte(metadata),
			}
			_ = r.OpLogs.Append(ctx, log)
		}

		return nil
	})

	if err != nil {
		return WrapError(err, "reject review")
	}

	return nil
}

// BatchApproveReviews 批量批准评价
func (s *AdminService) BatchApproveReviews(ctx context.Context, reviewIDs []uint64, actorUserID *uint64) error {
	if s.tx == nil {
		return apierr.InternalError("事务管理器未配置")
	}

	if len(reviewIDs) == 0 {
		return apierr.BadRequest("评价ID列表不能为空")
	}

	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		// 验证所有评价都是待审核状态
		for _, id := range reviewIDs {
			review, err := r.Reviews.Get(ctx, id)
			if err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					return apierr.NotFound("评价不存在: " + string(rune(id)))
				}
				return err
			}
			if review.Status != model.ReviewStatusPending {
				return apierr.BadRequest("评价不是待审核状态: " + string(rune(id)))
			}
		}

		// 批量更新状态
		if err := r.Reviews.BatchUpdateStatus(ctx, reviewIDs, model.ReviewStatusApproved, ""); err != nil {
			return err
		}

		// 记录操作日志
		if r.OpLogs != nil {
			for _, id := range reviewIDs {
				metadata := fmt.Sprintf(`{"old_status":"%s","new_status":"%s","batch":true}`, model.ReviewStatusPending, model.ReviewStatusApproved)
				log := &model.OperationLog{
					EntityType:   string(model.OpEntityReview),
					EntityID:     id,
					ActorUserID:  actorUserID,
					Action:       string(model.OpActionApprove),
					Reason:       "批量批准评价",
					MetadataJSON: []byte(metadata),
				}
				_ = r.OpLogs.Append(ctx, log)
			}
		}

		return nil
	})

	if err != nil {
		return WrapError(err, "batch approve reviews")
	}

	return nil
}

// BatchRejectReviews 批量拒绝评价
func (s *AdminService) BatchRejectReviews(ctx context.Context, reviewIDs []uint64, reason string, actorUserID *uint64) error {
	if s.tx == nil {
		return apierr.InternalError("事务管理器未配置")
	}

	if len(reviewIDs) == 0 {
		return apierr.BadRequest("评价ID列表不能为空")
	}

	if reason == "" {
		return apierr.BadRequest("拒绝原因不能为空")
	}

	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		// 验证所有评价都是待审核状态
		for _, id := range reviewIDs {
			review, err := r.Reviews.Get(ctx, id)
			if err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					return apierr.NotFound("评价不存在: " + string(rune(id)))
				}
				return err
			}
			if review.Status != model.ReviewStatusPending {
				return apierr.BadRequest("评价不是待审核状态: " + string(rune(id)))
			}
		}

		// 批量更新状态
		if err := r.Reviews.BatchUpdateStatus(ctx, reviewIDs, model.ReviewStatusRejected, reason); err != nil {
			return err
		}

		// 记录操作日志
		if r.OpLogs != nil {
			for _, id := range reviewIDs {
				metadata := fmt.Sprintf(`{"old_status":"%s","new_status":"%s","rejection_reason":"%s","batch":true}`, model.ReviewStatusPending, model.ReviewStatusRejected, reason)
				log := &model.OperationLog{
					EntityType:   string(model.OpEntityReview),
					EntityID:     id,
					ActorUserID:  actorUserID,
					Action:       string(model.OpActionReject),
					Reason:       reason,
					MetadataJSON: []byte(metadata),
				}
				_ = r.OpLogs.Append(ctx, log)
			}
		}

		return nil
	})

	if err != nil {
		return WrapError(err, "batch reject reviews")
	}

	return nil
}

// UpdateReviewReply 更新评价回复
func (s *AdminService) UpdateReviewReply(ctx context.Context, userID, replyID uint64, content string) (map[string]interface{}, error) {
	if s.tx == nil {
		return nil, apierr.InternalError("事务管理器未配置")
	}

	var result map[string]interface{}
	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		// 获取回复
		reply, err := r.ReviewReplies.Get(ctx, replyID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}
			return err
		}

		// 权限检查：只能更新自己的回复
		if reply.AuthorID != userID {
			return ErrUnauthorized
		}

		oldContent := reply.Content

		// 更新回复内容
		reply.Content = strings.TrimSpace(content)

		// 重新进行内容审核
		engine := feedservice.NewDefaultModerationEngine()
		moderationResult, err := engine.Evaluate(ctx, feedservice.ModerationInput{Content: reply.Content})
		if err != nil {
			return err
		}

		var status model.ReviewReplyStatus
		note := moderationResult.Reason
		switch moderationResult.Decision {
		case feedservice.ModerationDecisionApprove:
			status = model.ReviewReplyStatusApproved
		case feedservice.ModerationDecisionReject:
			status = model.ReviewReplyStatusRejected
		case feedservice.ModerationDecisionManual:
			status = model.ReviewReplyStatusPending
		default:
			status = model.ReviewReplyStatusPending
		}

		reply.Status = status
		reply.ModerationNote = note

		// 更新回复
		if err := r.ReviewReplies.Update(ctx, reply); err != nil {
			return err
		}

		// 记录操作日志
		if r.OpLogs != nil {
			metadata := fmt.Sprintf(`{"reply_id":%d,"old_content":"%s","new_content":"%s","status":"%s"}`, replyID, oldContent, reply.Content, status)
			log := &model.OperationLog{
				EntityType:   string(model.OpEntityReview),
				EntityID:     reply.ReviewID,
				ActorUserID:  &userID,
				Action:       string(model.OpActionUpdateReply),
				Reason:       "更新回复",
				MetadataJSON: []byte(metadata),
			}
			_ = r.OpLogs.Append(ctx, log)
		}

		// 发送通知给评价者
		review, err := r.Reviews.Get(ctx, reply.ReviewID)
		if err == nil && review.UserID != userID {
			notification := &model.NotificationEvent{
				UserID:        review.UserID,
				Title:         "评价回复已更新",
				Message:       "陪玩师更新了对您评价的回复",
				Channel:       "web",
				Priority:      model.NotificationPriorityNormal,
				ReferenceType: "review_reply",
				ReferenceID:   &replyID,
			}
			_ = r.Notifications.Create(ctx, notification)
		}

		result = map[string]interface{}{
			"replyId": reply.ID,
			"status":  string(reply.Status),
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteReviewReply 删除评价回复
func (s *AdminService) DeleteReviewReply(ctx context.Context, userID, replyID uint64) error {
	if s.tx == nil {
		return apierr.InternalError("事务管理器未配置")
	}

	return s.tx.WithTx(ctx, func(r *common.Repos) error {
		// 获取回复
		reply, err := r.ReviewReplies.Get(ctx, replyID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}
			return err
		}

		// 权限检查：只能删除自己的回复
		if reply.AuthorID != userID {
			return ErrUnauthorized
		}

		reviewID := reply.ReviewID

		// 删除回复
		if err := r.ReviewReplies.Delete(ctx, replyID); err != nil {
			return err
		}

		// 记录操作日志
		if r.OpLogs != nil {
			metadata := fmt.Sprintf(`{"reply_id":%d,"content":"%s"}`, replyID, reply.Content)
			log := &model.OperationLog{
				EntityType:   string(model.OpEntityReview),
				EntityID:     reviewID,
				ActorUserID:  &userID,
				Action:       string(model.OpActionDeleteReply),
				Reason:       "删除回复",
				MetadataJSON: []byte(metadata),
			}
			_ = r.OpLogs.Append(ctx, log)
		}

		// 发送通知给评价者
		review, err := r.Reviews.Get(ctx, reviewID)
		if err == nil && review.UserID != userID {
			notification := &model.NotificationEvent{
				UserID:        review.UserID,
				Title:         "评价回复已删除",
				Message:       "陪玩师删除了对您评价的回复",
				Channel:       "web",
				Priority:      model.NotificationPriorityNormal,
				ReferenceType: "review_reply",
				ReferenceID:   &replyID,
			}
			_ = r.Notifications.Create(ctx, notification)
		}

		return nil
	})
}

// BatchUpdateGamesStatus 批量更新游戏状态（启用/禁用）。
func (s *AdminService) BatchUpdateGamesStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error) {
	if len(ids) == 0 {
		return 0, apierr.BadRequest("no game ids provided")
	}
	updated, err := s.games.BatchUpdateStatus(ctx, ids, isActive)
	if err != nil {
		return 0, WrapError(err, "batch update games status")
	}
	s.invalidateCache(ctx, cacheKeyGames)
	for _, id := range ids {
		s.appendLogAsync(ctx, string(model.OpEntityGame), id, string(model.OpActionUpdate), map[string]any{"is_active": isActive})
	}
	return updated, nil
}

// BatchUpdateGamesSortOrder 批量更新游戏排序。
func (s *AdminService) BatchUpdateGamesSortOrder(ctx context.Context, updates map[uint64]int) (int64, error) {
	if len(updates) == 0 {
		return 0, apierr.BadRequest("no updates provided")
	}
	updated, err := s.games.BatchUpdateSortOrder(ctx, updates)
	if err != nil {
		return 0, WrapError(err, "batch update games sort order")
	}
	s.invalidateCache(ctx, cacheKeyGames)
	for id := range updates {
		s.appendLogAsync(ctx, string(model.OpEntityGame), id, string(model.OpActionUpdate), map[string]any{"sort_order": updates[id]})
	}
	return updated, nil
}
