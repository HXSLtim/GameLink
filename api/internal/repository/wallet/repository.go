package wallet

import (
	"context"
	"errors"

	"gamelink/internal/model"
	"gamelink/internal/repository"

	"gorm.io/gorm"
)

// ErrOptimisticLock 乐观锁冲突错误，表示并发更新冲突
var ErrOptimisticLock = errors.New("optimistic lock conflict: wallet was modified by another transaction")

type Repository interface {
	GetByUserID(ctx context.Context, userID uint64) (*model.Wallet, error)
	Save(ctx context.Context, wallet *model.Wallet) error
	// SaveWithOptimisticLock 使用乐观锁保存钱包，防止并发更新冲突
	// 如果版本号不匹配，返回 ErrOptimisticLock
	SaveWithOptimisticLock(ctx context.Context, wallet *model.Wallet) error
	// UpdateBalanceWithLock 原子更新余额，使用乐观锁防止并发冲突
	// delta 可以是正数（增加）或负数（减少）
	// 返回更新后的钱包或错误
	UpdateBalanceWithLock(ctx context.Context, userID uint64, delta int64, maxRetries int) (*model.Wallet, error)
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

// SaveWithOptimisticLock 使用乐观锁保存钱包
// 只有当数据库中的 version 与传入的 wallet.Version 匹配时才更新
// 更新成功后自动递增 version
func (r *gormWalletRepository) SaveWithOptimisticLock(ctx context.Context, wallet *model.Wallet) error {
	if wallet.ID == 0 {
		// 新记录，直接保存
		wallet.Version = 1
		return r.db.WithContext(ctx).Create(wallet).Error
	}

	// 使用乐观锁更新：WHERE id = ? AND version = ?
	result := r.db.WithContext(ctx).Model(wallet).
		Where("id = ? AND version = ?", wallet.ID, wallet.Version).
		Updates(map[string]interface{}{
			"balance_cents": wallet.BalanceCents,
			"frozen_cents":  wallet.FrozenCents,
			"version":       wallet.Version + 1,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrOptimisticLock
	}

	// 更新成功，递增本地版本号
	wallet.Version++
	return nil
}

// UpdateBalanceWithLock 原子更新余额，带重试机制
// 使用乐观锁防止并发冲突，失败时自动重试
func (r *gormWalletRepository) UpdateBalanceWithLock(ctx context.Context, userID uint64, delta int64, maxRetries int) (*model.Wallet, error) {
	if maxRetries <= 0 {
		maxRetries = 3
	}

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		// 获取当前钱包
		wallet, err := r.GetByUserID(ctx, userID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) && delta > 0 {
				// 钱包不存在且是增加余额，创建新钱包
				wallet = &model.Wallet{
					UserID:       userID,
					BalanceCents: delta,
					Version:      1,
				}
				if err := r.db.WithContext(ctx).Create(wallet).Error; err != nil {
					lastErr = err
					continue
				}
				return wallet, nil
			}
			return nil, err
		}

		// 检查余额是否足够（如果是扣款）
		if delta < 0 && wallet.BalanceCents+delta < 0 {
			return nil, errors.New("insufficient balance")
		}

		// 更新余额
		wallet.BalanceCents += delta

		// 使用乐观锁保存
		if err := r.SaveWithOptimisticLock(ctx, wallet); err != nil {
			if errors.Is(err, ErrOptimisticLock) {
				lastErr = err
				continue // 重试
			}
			return nil, err
		}

		return wallet, nil
	}

	return nil, lastErr
}
