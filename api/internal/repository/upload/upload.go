package upload

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// NewUploadRepository returns a GORM-based upload repository.
func NewUploadRepository(db *gorm.DB) repository.UploadRepository {
	return &gormUploadRepository{db: db}
}

type gormUploadRepository struct {
	db *gorm.DB
}

func (r *gormUploadRepository) Create(ctx context.Context, upload *model.Upload) error {
	return r.db.WithContext(ctx).Create(upload).Error
}
