package repository

import (
    "context"
    "time"

    "gamelink/internal/model"
)

// UserRepository 定义后台对用户的管理操作。
type UserRepository interface {
    List(ctx context.Context) ([]model.User, error)
    // ListPaged 支持分页列表与总数统计
    ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error)
    // ListWithFilters 支持按角色/状态/时间范围/关键字的分页筛选
    ListWithFilters(ctx context.Context, opts UserListOptions) ([]model.User, int64, error)
    Get(ctx context.Context, id uint64) (*model.User, error)
	// FindByEmail 通过邮箱查找用户。
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	// FindByPhone 通过手机号查找用户。
	FindByPhone(ctx context.Context, phone string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
    Update(ctx context.Context, user *model.User) error
    Delete(ctx context.Context, id uint64) error
}

// UserListOptions 用户高级查询参数
type UserListOptions struct {
    Page       int
    PageSize   int
    Roles      []model.Role
    Statuses   []model.UserStatus
    DateFrom   *time.Time
    DateTo     *time.Time
    Keyword    string // 匹配 name/email/phone
}
