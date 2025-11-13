package orderdispute

import (
        "context"

        "gorm.io/gorm"

        "gamelink/internal/model"
        "gamelink/internal/repository"
)

// gormOrderDisputeRepository 实现 OrderDisputeRepository。
type gormOrderDisputeRepository struct {
        db *gorm.DB
}

// NewRepository 创建争议仓储。
func NewRepository(db *gorm.DB) repository.OrderDisputeRepository {
        return &gormOrderDisputeRepository{db: db}
}

func (r *gormOrderDisputeRepository) Create(ctx context.Context, dispute *model.OrderDispute) error {
        return r.db.WithContext(ctx).Create(dispute).Error
}

func (r *gormOrderDisputeRepository) Update(ctx context.Context, dispute *model.OrderDispute) error {
        tx := r.db.WithContext(ctx).Model(dispute).Where("id = ?", dispute.ID).Updates(map[string]any{
                "status":               dispute.Status,
                "resolution":           dispute.Resolution,
                "resolution_note":      dispute.ResolutionNote,
                "refund_amount_cents":  dispute.RefundAmountCents,
                "handled_by_id":        dispute.HandledByID,
                "handled_at":           dispute.HandledAt,
                "response_deadline":    dispute.ResponseDeadline,
                "responded_at":         dispute.RespondedAt,
        })
        if tx.Error != nil {
                return tx.Error
        }
        if tx.RowsAffected == 0 {
                return repository.ErrNotFound
        }
        return nil
}

func (r *gormOrderDisputeRepository) ListByOrder(ctx context.Context, orderID uint64) ([]model.OrderDispute, error) {
        var disputes []model.OrderDispute
        if err := r.db.WithContext(ctx).
                Where("order_id = ?", orderID).
                Order("created_at DESC").
                Find(&disputes).Error; err != nil {
                return nil, err
        }
        return disputes, nil
}

func (r *gormOrderDisputeRepository) GetLatestByOrder(ctx context.Context, orderID uint64) (*model.OrderDispute, error) {
        var dispute model.OrderDispute
        err := r.db.WithContext(ctx).
                Where("order_id = ?", orderID).
                Order("created_at DESC").
                First(&dispute).Error
        if err != nil {
                if err == gorm.ErrRecordNotFound {
                        return nil, repository.ErrNotFound
                }
                return nil, err
        }
        return &dispute, nil
}
