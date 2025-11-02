package item

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	serviceitemrepo "gamelink/internal/repository/serviceitem"
)

var (
	// ErrNotFound 服务项目不存�?
	ErrNotFound = repository.ErrNotFound
	// ErrValidation 表示输入校验失败
	ErrValidation = errors.New("validation failed")
)

// ServiceItemService 服务项目服务（统一管理护航服务和礼物）
type ServiceItemService struct {
	items  serviceitemrepo.ServiceItemRepository
	games  repository.GameRepository
	players repository.PlayerRepository
}

// NewServiceItemService 创建服务项目服务
func NewServiceItemService(
	items serviceitemrepo.ServiceItemRepository,
	games repository.GameRepository,
	players repository.PlayerRepository,
) *ServiceItemService {
	return &ServiceItemService{
		items:  items,
		games:  games,
		players: players,
	}
}

// CreateServiceItemRequest 创建服务项目请求
type CreateServiceItemRequest struct {
	ItemCode        string                         `json:"itemCode" binding:"required,max=32"`
	Name            string                         `json:"name" binding:"required,max=128"`
	Description     string                         `json:"description"`
	SubCategory     model.ServiceItemSubCategory   `json:"subCategory" binding:"required,oneof=solo team gift"`
	GameID          *uint64                        `json:"gameId"`
	PlayerID        *uint64                        `json:"playerId"`
	RankLevel       string                         `json:"rankLevel"`
	BasePriceCents  int64                          `json:"basePriceCents" binding:"required,min=0"`
	ServiceHours    int                            `json:"serviceHours" binding:"min=0"`
	CommissionRate  float64                        `json:"commissionRate" binding:"required,min=0,max=1"`
	MinUsers        int                            `json:"minUsers" binding:"min=1"`
	MaxPlayers      int                            `json:"maxPlayers" binding:"min=1"`
	Tags            string                         `json:"tags"`
	IconURL         string                         `json:"iconUrl"`
	SortOrder       int                            `json:"sortOrder"`
}

// CreateServiceItem 创建服务项目
func (s *ServiceItemService) CreateServiceItem(ctx context.Context, req CreateServiceItemRequest) (*model.ServiceItem, error) {
	// 验证游戏ID（如果提供）
	if req.GameID != nil {
		_, err := s.games.Get(ctx, *req.GameID)
		if err != nil {
			return nil, fmt.Errorf("invalid game_id: %w", err)
		}
	}

	// 验证陪玩师ID（如果提供）
	if req.PlayerID != nil {
		_, err := s.players.Get(ctx, *req.PlayerID)
		if err != nil {
			return nil, fmt.Errorf("invalid player_id: %w", err)
		}
	}

	// 验证礼物的service_hours必须�?
	if req.SubCategory == model.SubCategoryGift && req.ServiceHours != 0 {
		return nil, errors.New("gift items must have service_hours = 0")
	}

	item := &model.ServiceItem{
		ItemCode:       req.ItemCode,
		Name:           req.Name,
		Description:    req.Description,
		Category:       "escort", // 统一�?escort
		SubCategory:    req.SubCategory,
		GameID:         req.GameID,
		PlayerID:       req.PlayerID,
		RankLevel:      req.RankLevel,
		BasePriceCents: req.BasePriceCents,
		ServiceHours:   req.ServiceHours,
		CommissionRate: req.CommissionRate,
		MinUsers:       req.MinUsers,
		MaxPlayers:     req.MaxPlayers,
		Tags:           req.Tags,
		IconURL:        req.IconURL,
		IsActive:       true,
		SortOrder:      req.SortOrder,
	}

	if err := s.items.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

// UpdateServiceItemRequest 更新服务项目请求
type UpdateServiceItemRequest struct {
	Name           *string  `json:"name"`
	Description    *string  `json:"description"`
	BasePriceCents *int64   `json:"basePriceCents"`
	ServiceHours   *int     `json:"serviceHours"`
	CommissionRate *float64 `json:"commissionRate"`
	RankLevel      *string  `json:"rankLevel"`
	Tags           *string  `json:"tags"`
	IconURL        *string  `json:"iconUrl"`
	IsActive       *bool    `json:"isActive"`
	SortOrder      *int     `json:"sortOrder"`
}

// UpdateServiceItem 更新服务项目
func (s *ServiceItemService) UpdateServiceItem(ctx context.Context, id uint64, req UpdateServiceItemRequest) error {
	item, err := s.items.Get(ctx, id)
	if err != nil {
		return err
	}

	// 更新字段
	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	if req.BasePriceCents != nil {
		if *req.BasePriceCents < 0 {
			return errors.New("base price must be >= 0")
		}
		item.BasePriceCents = *req.BasePriceCents
	}
	if req.ServiceHours != nil {
		// 礼物的service_hours必须�?
		if item.IsGift() && *req.ServiceHours != 0 {
			return errors.New("gift items must have service_hours = 0")
		}
		item.ServiceHours = *req.ServiceHours
	}
	if req.CommissionRate != nil {
		if *req.CommissionRate < 0 || *req.CommissionRate > 1 {
			return errors.New("commission rate must be between 0 and 1")
		}
		item.CommissionRate = *req.CommissionRate
	}
	if req.RankLevel != nil {
		item.RankLevel = *req.RankLevel
	}
	if req.Tags != nil {
		item.Tags = *req.Tags
	}
	if req.IconURL != nil {
		item.IconURL = *req.IconURL
	}
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		item.SortOrder = *req.SortOrder
	}

	return s.items.Update(ctx, item)
}

// DeleteServiceItem 删除服务项目
func (s *ServiceItemService) DeleteServiceItem(ctx context.Context, id uint64) error {
	_, err := s.items.Get(ctx, id)
	if err != nil {
		return err
	}

	return s.items.Delete(ctx, id)
}

// ServiceItemDTO 服务项目DTO
type ServiceItemDTO struct {
	ID             uint64    `json:"id"`
	ItemCode       string    `json:"itemCode"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Category       string    `json:"category"`
	SubCategory    string    `json:"subCategory"`
	GameID         *uint64   `json:"gameId"`
	GameName       string    `json:"gameName,omitempty"`
	PlayerID       *uint64   `json:"playerId"`
	PlayerNickname string    `json:"playerNickname,omitempty"`
	RankLevel      string    `json:"rankLevel"`
	BasePriceCents int64     `json:"basePriceCents"`
	ServiceHours   int       `json:"serviceHours"`
	CommissionRate float64   `json:"commissionRate"`
	MinUsers       int       `json:"minUsers"`
	MaxPlayers     int       `json:"maxPlayers"`
	Tags           string    `json:"tags"`
	IconURL        string    `json:"iconUrl"`
	IsActive       bool      `json:"isActive"`
	SortOrder      int       `json:"sortOrder"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// GetServiceItem 获取服务项目详情
func (s *ServiceItemService) GetServiceItem(ctx context.Context, id uint64) (*ServiceItemDTO, error) {
	item, err := s.items.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toDTO(ctx, item), nil
}

// ListServiceItems 获取服务项目列表
func (s *ServiceItemService) ListServiceItems(ctx context.Context, req ListServiceItemsRequest) (*ServiceItemListResponse, error) {
	items, total, err := s.items.List(ctx, serviceitemrepo.ServiceItemListOptions{
		Category:    req.Category,
		SubCategory: req.SubCategory,
		GameID:      req.GameID,
		PlayerID:    req.PlayerID,
		IsActive:    req.IsActive,
		Page:        req.Page,
		PageSize:    req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	itemDTOs := make([]ServiceItemDTO, 0, len(items))
	for _, item := range items {
		itemDTOs = append(itemDTOs, *s.toDTO(ctx, &item))
	}

	return &ServiceItemListResponse{
		Items: itemDTOs,
		Total: total,
	}, nil
}

// ListServiceItemsRequest 服务项目列表请求
type ListServiceItemsRequest struct {
	Category    *string                        `form:"category"`
	SubCategory *model.ServiceItemSubCategory  `form:"subCategory"`
	GameID      *uint64                        `form:"gameId"`
	PlayerID    *uint64                        `form:"playerId"`
	IsActive    *bool                          `form:"isActive"`
	Page        int                            `form:"page"`
	PageSize    int                            `form:"pageSize"`
}

// ServiceItemListResponse 服务项目列表响应
type ServiceItemListResponse struct {
	Items []ServiceItemDTO `json:"items"`
	Total int64            `json:"total"`
}

// GetGiftList 获取礼物列表（用户端�?
func (s *ServiceItemService) GetGiftList(ctx context.Context, page, pageSize int) (*ServiceItemListResponse, error) {
	items, total, err := s.items.GetGifts(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	itemDTOs := make([]ServiceItemDTO, 0, len(items))
	for _, item := range items {
		itemDTOs = append(itemDTOs, *s.toDTO(ctx, &item))
	}

	return &ServiceItemListResponse{
		Items: itemDTOs,
		Total: total,
	}, nil
}

// toDTO 转换为DTO
func (s *ServiceItemService) toDTO(ctx context.Context, item *model.ServiceItem) *ServiceItemDTO {
	dto := &ServiceItemDTO{
		ID:             item.ID,
		ItemCode:       item.ItemCode,
		Name:           item.Name,
		Description:    item.Description,
		Category:       item.Category,
		SubCategory:    string(item.SubCategory),
		GameID:         item.GameID,
		PlayerID:       item.PlayerID,
		RankLevel:      item.RankLevel,
		BasePriceCents: item.BasePriceCents,
		ServiceHours:   item.ServiceHours,
		CommissionRate: item.CommissionRate,
		MinUsers:       item.MinUsers,
		MaxPlayers:     item.MaxPlayers,
		Tags:           item.Tags,
		IconURL:        item.IconURL,
		IsActive:       item.IsActive,
		SortOrder:      item.SortOrder,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}

	// 获取游戏名称
	if item.GameID != nil {
		game, err := s.games.Get(ctx, *item.GameID)
		if err == nil {
			dto.GameName = game.Name
		}
	}

	// 获取陪玩师昵�?
	if item.PlayerID != nil {
		player, err := s.players.Get(ctx, *item.PlayerID)
		if err == nil {
			dto.PlayerNickname = player.Nickname
		}
	}

	return dto
}

// BatchUpdateStatusRequest 批量更新状态请�?
type BatchUpdateStatusRequest struct {
	IDs      []uint64 `json:"ids" binding:"required"`
	IsActive bool     `json:"isActive"`
}

// BatchUpdateStatus 批量更新状�?
func (s *ServiceItemService) BatchUpdateStatus(ctx context.Context, req BatchUpdateStatusRequest) error {
	if len(req.IDs) == 0 {
		return errors.New("no item ids provided")
	}
	return s.items.BatchUpdateStatus(ctx, req.IDs, req.IsActive)
}

// BatchUpdatePriceRequest 批量更新价格请求
type BatchUpdatePriceRequest struct {
	IDs            []uint64 `json:"ids" binding:"required"`
	BasePriceCents int64    `json:"basePriceCents" binding:"required,min=0"`
}

// BatchUpdatePrice 批量更新价格
func (s *ServiceItemService) BatchUpdatePrice(ctx context.Context, req BatchUpdatePriceRequest) error {
	if len(req.IDs) == 0 {
		return errors.New("no item ids provided")
	}
	return s.items.BatchUpdatePrice(ctx, req.IDs, req.BasePriceCents)
}

