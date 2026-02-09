package player

import (
	"context"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
)

// ServiceManagement handles player service CRUD.
type ServiceManagement struct {
	services repository.PlayerServiceRepository
	players  repository.PlayerRepository
	games    repository.GameRepository
	ranks    repository.GameRankRepository
}

// NewServiceManagement creates service management.
func NewServiceManagement(
	services repository.PlayerServiceRepository,
	players repository.PlayerRepository,
	games repository.GameRepository,
	ranks repository.GameRankRepository,
) *ServiceManagement {
	return &ServiceManagement{
		services: services,
		players:  players,
		games:    games,
		ranks:    ranks,
	}
}

// CreateServiceRequest represents create request.
type CreateServiceRequest struct {
	GameID      uint64 `json:"gameId" binding:"required"`
	RankID      uint64 `json:"rankId" binding:"required"`
	Description string `json:"description"`
}

// UpdateServiceRequest represents update request.
type UpdateServiceRequest struct {
	GameID      *uint64 `json:"gameId"`
	RankID      *uint64 `json:"rankId"`
	Description *string `json:"description"`
}

// UpdateServiceStatusRequest represents status update.
type UpdateServiceStatusRequest struct {
	IsActive bool `json:"isActive"`
}

// ListServices lists services for current player.
func (s *ServiceManagement) ListServices(ctx context.Context, userID uint64) ([]model.PlayerService, error) {
	player, err := s.players.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.services.ListByPlayer(ctx, player.ID)
}

// CreateService creates a new service for current player.
func (s *ServiceManagement) CreateService(ctx context.Context, userID uint64, req CreateServiceRequest) (*model.PlayerService, error) {
	player, err := s.players.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if _, err := s.games.Get(ctx, req.GameID); err != nil {
		return nil, apierr.BadRequest("无效的游戏ID").WithDetails(err.Error())
	}
	rank, err := s.ranks.Get(ctx, req.RankID)
	if err != nil {
		return nil, apierr.BadRequest("无效的段位ID").WithDetails(err.Error())
	}
	if rank.GameID != req.GameID {
		return nil, apierr.BadRequest("段位与游戏不匹配")
	}

	service := &model.PlayerService{
		PlayerID:    player.ID,
		GameID:      req.GameID,
		RankID:      req.RankID,
		Description: req.Description,
		IsActive:    true,
	}
	if err := s.services.Create(ctx, service); err != nil {
		return nil, err
	}
	return s.services.Get(ctx, service.ID)
}

// UpdateService updates a service.
func (s *ServiceManagement) UpdateService(ctx context.Context, userID uint64, serviceID uint64, req UpdateServiceRequest) (*model.PlayerService, error) {
	player, err := s.players.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	service, err := s.services.Get(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	if service.PlayerID != player.ID {
		return nil, apierr.Unauthorized("无权操作该服务")
	}
	if req.GameID != nil {
		if _, err := s.games.Get(ctx, *req.GameID); err != nil {
			return nil, apierr.BadRequest("无效的游戏ID").WithDetails(err.Error())
		}
		service.GameID = *req.GameID
	}
	if req.RankID != nil {
		rank, err := s.ranks.Get(ctx, *req.RankID)
		if err != nil {
			return nil, apierr.BadRequest("无效的段位ID").WithDetails(err.Error())
		}
		expectedGameID := service.GameID
		if req.GameID != nil {
			expectedGameID = *req.GameID
		}
		if rank.GameID != expectedGameID {
			return nil, apierr.BadRequest("段位与游戏不匹配")
		}
		service.RankID = *req.RankID
	}
	if req.Description != nil {
		service.Description = *req.Description
	}
	if req.GameID != nil && req.RankID == nil {
		rank, err := s.ranks.Get(ctx, service.RankID)
		if err == nil && rank.GameID != service.GameID {
			return nil, apierr.BadRequest("段位与游戏不匹配")
		}
	}
	if err := s.services.Update(ctx, service); err != nil {
		return nil, err
	}
	return s.services.Get(ctx, service.ID)
}

// DeleteService deletes a service.
func (s *ServiceManagement) DeleteService(ctx context.Context, userID uint64, serviceID uint64) error {
	player, err := s.players.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	service, err := s.services.Get(ctx, serviceID)
	if err != nil {
		return err
	}
	if service.PlayerID != player.ID {
		return apierr.Unauthorized("无权操作该服务")
	}
	return s.services.Delete(ctx, serviceID)
}

// UpdateServiceStatus updates service active status.
func (s *ServiceManagement) UpdateServiceStatus(ctx context.Context, userID uint64, serviceID uint64, isActive bool) (*model.PlayerService, error) {
	player, err := s.players.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	service, err := s.services.Get(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	if service.PlayerID != player.ID {
		return nil, apierr.Unauthorized("无权操作该服务")
	}
	service.IsActive = isActive
	if err := s.services.Update(ctx, service); err != nil {
		return nil, err
	}
	return s.services.Get(ctx, service.ID)
}
