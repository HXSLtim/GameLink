package wallet

import (
	"context"

	"gamelink/internal/model"
	"gamelink/internal/repository"

	"gorm.io/gorm"
)

type Repository interface {
	GetByUserID(ctx context.Context, userID uint64) (*model.Wallet, error)
	Save(ctx context.Context, wallet *model.Wallet) error
}

type gormWalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) Repository {
	return &gormWalletRepository{db: db}
}

func (r *gormWalletRepository) GetByUserID(ctx context.Context, userID uint64) (*model.Wallet, error) {
	var w model.Wallet
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&w).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &w, nil
}

func (r *gormWalletRepository) Save(ctx context.Context, wallet *model.Wallet) error {
	return r.db.WithContext(ctx).Save(wallet).Error
}
