package gormrepo

import (
    "context"
    "strings"

    "gorm.io/gorm"

    "gamelink/internal/model"
    "gamelink/internal/repository"
)

// UserRepository 实现用户管理仓储。
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储。
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// List returns all users ordered by creation time.
func (r *UserRepository) List(ctx context.Context) ([]model.User, error) {
	var users []model.User
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// ListPaged 返回分页用户列表与总数。
// ListPaged returns a page of users and the total count.
func (r *UserRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&model.User{})

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []model.User
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// ListWithFilters returns a page of users with filters and the total count.
func (r *UserRepository) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
    page := repository.NormalizePage(opts.Page)
    pageSize := repository.NormalizePageSize(opts.PageSize)
    offset := (page - 1) * pageSize

    q := r.db.WithContext(ctx).Model(&model.User{})
    if len(opts.Roles) > 0 {
        q = q.Where("role IN ?", opts.Roles)
    }
    if len(opts.Statuses) > 0 {
        q = q.Where("status IN ?", opts.Statuses)
    }
    if opts.DateFrom != nil {
        q = q.Where("created_at >= ?", *opts.DateFrom)
    }
    if opts.DateTo != nil {
        q = q.Where("created_at <= ?", *opts.DateTo)
    }
    if kw := strings.TrimSpace(opts.Keyword); kw != "" {
        like := "%" + kw + "%"
        q = q.Where("name LIKE ? OR email LIKE ? OR phone LIKE ?", like, like, like)
    }

    var total int64
    if err := q.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    var users []model.User
    if err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
        return nil, 0, err
    }
    return users, total, nil
}

// Get returns a user by id.
func (r *UserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail returns a user by unique email.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByPhone returns a user by unique phone.
func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Create inserts a new user.
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Update updates editable fields of a user.
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	tx := r.db.WithContext(ctx).Model(user).Updates(map[string]any{
		"phone":         user.Phone,
		"email":         user.Email,
		"name":          user.Name,
		"avatar_url":    user.AvatarURL,
		"role":          user.Role,
		"status":        user.Status,
		"password_hash": user.PasswordHash,
	})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// Delete soft-deletes a user by id.
func (r *UserRepository) Delete(ctx context.Context, id uint64) error {
	tx := r.db.WithContext(ctx).Delete(&model.User{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

var _ repository.UserRepository = (*UserRepository)(nil)
