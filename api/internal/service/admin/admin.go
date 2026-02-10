package admin

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"
	
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/common"
	repoiface "gamelink/internal/repository/interfaces"
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
