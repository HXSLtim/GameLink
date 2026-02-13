package admin

import (
	"context"
	"errors"
	"strings"
	"time"

	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/logging"
)

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

func mapUserError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, repository.ErrNotFound) {
		return ErrUserNotFound
	}
	return WrapError(err, "操作用户数据失败")
}
