package reviewreply

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// NewReviewReplyRepository creates a repository for review replies.
func NewReviewReplyRepository(db *gorm.DB) repository.ReviewReplyRepository {
	return &gormReviewReplyRepository{db: db}
}

type gormReviewReplyRepository struct {
	db *gorm.DB
}

func (r *gormReviewReplyRepository) Create(ctx context.Context, reply *model.ReviewReply) error {
	return r.db.WithContext(ctx).Create(reply).Error
}

func (r *gormReviewReplyRepository) ListByReview(ctx context.Context, reviewID uint64) ([]model.ReviewReply, error) {
	var replies []model.ReviewReply
	err := r.db.WithContext(ctx).Where("review_id = ?", reviewID).Order("created_at ASC").Find(&replies).Error
	return replies, err
}

func (r *gormReviewReplyRepository) UpdateStatus(ctx context.Context, replyID uint64, status string, note string) error {
	tx := r.db.WithContext(ctx).Model(&model.ReviewReply{}).Where("id = ?", replyID).Updates(map[string]any{
		"status":          status,
		"moderation_note": note,
		"moderated_at":    gorm.Expr("CURRENT_TIMESTAMP"),
	})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}
