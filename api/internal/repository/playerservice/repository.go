package playerservice

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// NewPlayerServiceRepository creates player service repository.
func NewPlayerServiceRepository(db *gorm.DB) repository.PlayerServiceRepository {
	return &gormPlayerServiceRepository{db: db}
}

type gormPlayerServiceRepository struct {
	db *gorm.DB
}

func (r *gormPlayerServiceRepository) Get(ctx context.Context, id uint64) (*model.PlayerService, error) {
	var service model.PlayerService
	if err := r.db.WithContext(ctx).
		Preload("Game").
		Preload("Rank").
		First(&service, id).Error; err != nil {
		return nil, repository.WrapNotFound(err)
	}
	return &service, nil
}

func (r *gormPlayerServiceRepository) ListByPlayer(ctx context.Context, playerID uint64) ([]model.PlayerService, error) {
	var services []model.PlayerService
	if err := r.db.WithContext(ctx).
		Where("player_id = ?", playerID).
		Preload("Game").
		Preload("Rank").
		Order("updated_at DESC").
		Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}

func (r *gormPlayerServiceRepository) Create(ctx context.Context, service *model.PlayerService) error {
	return r.db.WithContext(ctx).Create(service).Error
}

func (r *gormPlayerServiceRepository) Update(ctx context.Context, service *model.PlayerService) error {
	return r.db.WithContext(ctx).Save(service).Error
}

func (r *gormPlayerServiceRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.PlayerService{}, id).Error
}
