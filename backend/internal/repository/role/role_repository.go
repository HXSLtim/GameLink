package role

import (
	"context"

	"gamelink/internal/model"
)

// RoleRepository å®ä¹è§è²çæ°æ®è®¿é®æä½ã?
type RoleRepository interface {
	// List è·åææè§è²åè¡?
	List(ctx context.Context) ([]model.RoleModel, error)

	// ListPaged åé¡µè·åè§è²åè¡¨
	ListPaged(ctx context.Context, page, pageSize int) ([]model.RoleModel, int64, error)

	// ListWithPermissions è·åè§è²åè¡¨ï¼é¢å è½½æé
	ListWithPermissions(ctx context.Context) ([]model.RoleModel, error)

	// Get æ ¹æ®IDè·åè§è²
	Get(ctx context.Context, id uint64) (*model.RoleModel, error)

	// GetWithPermissions æ ¹æ®IDè·åè§è²ï¼é¢å è½½æé
	GetWithPermissions(ctx context.Context, id uint64) (*model.RoleModel, error)

	// GetBySlug æ ¹æ®Slugè·åè§è²
	GetBySlug(ctx context.Context, slug string) (*model.RoleModel, error)

	// Create åå»ºè§è²
	Create(ctx context.Context, role *model.RoleModel) error

	// Update æ´æ°è§è²
	Update(ctx context.Context, role *model.RoleModel) error

	// Delete å é¤è§è²ï¼ç³»ç»è§è²ä¸å¯å é¤ï¼
	Delete(ctx context.Context, id uint64) error

	// AssignPermissions ä¸ºè§è²åéæéï¼æ¿æ¢ç°ææéï¼?
	AssignPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error

	// AddPermissions ä¸ºè§è²æ·»å æéï¼è¿½å ï¼?
	AddPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error

	// RemovePermissions ç§»é¤è§è²çæé?
	RemovePermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error

	// ListByUserID è·åç¨æ·æ¥æçææè§è?
	ListByUserID(ctx context.Context, userID uint64) ([]model.RoleModel, error)

	// AssignToUser ä¸ºç¨æ·åéè§è?
	AssignToUser(ctx context.Context, userID uint64, roleIDs []uint64) error

	// RemoveFromUser ç§»é¤ç¨æ·çè§è?
	RemoveFromUser(ctx context.Context, userID uint64, roleIDs []uint64) error

	// CheckUserHasRole æ£æ¥ç¨æ·æ¯å¦æ¥ææå®è§è?
	CheckUserHasRole(ctx context.Context, userID uint64, roleSlug string) (bool, error)
}

