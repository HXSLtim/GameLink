package item

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
	"gamelink/pkg/cache"
)

var (
	// ErrNotFound 服务项目不存
	ErrNotFound = repository.ErrNotFound
	// ErrValidation 表示输入校验失败
	ErrValidation = apierr.BadRequest("输入参数验证失败")
)

// Cache key patterns and TTL
const (
	cacheKeyActiveItems = "service_items:active"
	cacheKeyGifts       = "service_items:gifts"
	cacheKeyItem        = "service_item:%d"
	cacheKeyGame        = "game:%d"
	cacheTTLItems       = 10 * time.Minute
	cacheTTLItem        = 5 * time.Minute
	cacheTTLGame        = 30 * time.Minute
)

// ServiceItemService 服务项目服务(统一管理护航服务和礼物)
type ServiceItemService struct {
	items   repository.ServiceItemRepository
	games   repository.GameRepository
	players repository.PlayerRepository
	cache   cache.Cache
}

// NewServiceItemService 创建服务项目服务
func NewServiceItemService(
	items repository.ServiceItemRepository,
	games repository.GameRepository,
	players repository.PlayerRepository,
) *ServiceItemService {
	return &ServiceItemService{
		items:   items,
		games:   games,
		players: players,
	}
}

// NewServiceItemServiceWithCache 创建带缓存的服务项目服务
func NewServiceItemServiceWithCache(
	items repository.ServiceItemRepository,
	games repository.GameRepository,
	players repository.PlayerRepository,
	c cache.Cache,
) *ServiceItemService {
	return &ServiceItemService{
		items:   items,
		games:   games,
		players: players,
		cache:   c,
	}
}

// SetCache 设置缓存实例
func (s *ServiceItemService) SetCache(c cache.Cache) {
	s.cache = c
}

// invalidateItemCaches 清除服务项目相关缓存
func (s *ServiceItemService) invalidateItemCaches(ctx context.Context, itemID uint64) {
	if s.cache == nil {
		return
	}
	// 清除列表缓存
	_ = s.cache.Delete(ctx, cacheKeyActiveItems)
	_ = s.cache.Delete(ctx, cacheKeyGifts)
	// 清除单个项目缓存
	if itemID > 0 {
		_ = s.cache.Delete(ctx, fmt.Sprintf(cacheKeyItem, itemID))
	}
}

// getCachedGame 从缓存获取游戏信息
func (s *ServiceItemService) getCachedGame(ctx context.Context, gameID uint64) (*model.Game, error) {
	if s.cache == nil {
		return s.games.Get(ctx, gameID)
	}

	key := fmt.Sprintf(cacheKeyGame, gameID)
	if val, ok, _ := s.cache.Get(ctx, key); ok {
		var game model.Game
		if err := json.Unmarshal([]byte(val), &game); err == nil {
			return &game, nil
		}
	}

	game, err := s.games.Get(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(game); err == nil {
		_ = s.cache.Set(ctx, key, string(data), cacheTTLGame)
	}
	return game, nil
}

// CreateServiceItemRequest 创建服务项目请求
type CreateServiceItemRequest struct {
	ItemCode       string                       `json:"itemCode" binding:"required,max=32"`
	Name           string                       `json:"name" binding:"required,max=128"`
	Description    string                       `json:"description"`
	SubCategory    model.ServiceItemSubCategory `json:"subCategory" binding:"required,oneof=solo team gift"`
	GameID         *uint64                      `json:"gameId"`
	PlayerID       *uint64                      `json:"playerId"`
	RankLevel      string                       `json:"rankLevel"`
	BasePriceCents int64                        `json:"basePriceCents" binding:"required,min=0"`
	ServiceHours   int                          `json:"serviceHours" binding:"min=0"`
	CommissionRate float64                      `json:"commissionRate" binding:"required,min=0,max=1"`
	MinUsers       int                          `json:"minUsers" binding:"min=1"`
	MaxPlayers     int                          `json:"maxPlayers" binding:"min=1"`
	Tags           string                       `json:"tags"`
	IconURL        string                       `json:"iconUrl"`
	SortOrder      int                          `json:"sortOrder"`
}

// CreateServiceItem 创建服务项目
func (s *ServiceItemService) CreateServiceItem(ctx context.Context, req CreateServiceItemRequest) (*model.ServiceItem, error) {
	// 验证游戏ID（如果提供）
	if req.GameID != nil {
		_, err := s.games.Get(ctx, *req.GameID)
		if err != nil {
			return nil, apierr.BadRequest("游戏ID无效").WithDetails(err.Error())
		}
	}

	// 验证陪玩师ID（如果提供）
	if req.PlayerID != nil {
		_, err := s.players.Get(ctx, *req.PlayerID)
		if err != nil {
			return nil, apierr.BadRequest("陪玩师ID无效").WithDetails(err.Error())
		}
	}

	// 验证礼物的service_hours必须
	if req.SubCategory == model.SubCategoryGift && req.ServiceHours != 0 {
		return nil, apierr.BadRequest("礼物类项目的服务时长必须为0")
	}

	item := &model.ServiceItem{
		ItemCode:       req.ItemCode,
		Name:           req.Name,
		Description:    req.Description,
		Category:       "escort", // 统一escort
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

	// 清除缓存
	s.invalidateItemCaches(ctx, item.ID)

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

	if err := applyServiceItemUpdates(item, req); err != nil {
		return err
	}

	if err := s.items.Update(ctx, item); err != nil {
		return err
	}

	// 清除缓存
	s.invalidateItemCaches(ctx, id)

	return nil
}

// applyServiceItemUpdates 应用服务项目更新
func applyServiceItemUpdates(item *model.ServiceItem, req UpdateServiceItemRequest) error {
	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	if err := applyPriceUpdate(item, req); err != nil {
		return err
	}
	if err := applyServiceHoursUpdate(item, req); err != nil {
		return err
	}
	if err := applyCommissionRateUpdate(item, req); err != nil {
		return err
	}
	applyOptionalFields(item, req)
	return nil
}

func applyPriceUpdate(item *model.ServiceItem, req UpdateServiceItemRequest) error {
	if req.BasePriceCents == nil {
		return nil
	}
	if *req.BasePriceCents < 0 {
		return apierr.BadRequest("基础价格必须大于等于0")
	}
	item.BasePriceCents = *req.BasePriceCents
	return nil
}

func applyServiceHoursUpdate(item *model.ServiceItem, req UpdateServiceItemRequest) error {
	if req.ServiceHours == nil {
		return nil
	}
	if item.IsGift() && *req.ServiceHours != 0 {
		return apierr.BadRequest("礼物类项目的服务时长必须为0")
	}
	item.ServiceHours = *req.ServiceHours
	return nil
}

func applyCommissionRateUpdate(item *model.ServiceItem, req UpdateServiceItemRequest) error {
	if req.CommissionRate == nil {
		return nil
	}
	if *req.CommissionRate < 0 || *req.CommissionRate > 1 {
		return apierr.BadRequest("抽成比例必须在0-1之间")
	}
	item.CommissionRate = *req.CommissionRate
	return nil
}

func applyOptionalFields(item *model.ServiceItem, req UpdateServiceItemRequest) {
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
}

// DeleteServiceItem 删除服务项目
func (s *ServiceItemService) DeleteServiceItem(ctx context.Context, id uint64) error {
	_, err := s.items.Get(ctx, id)
	if err != nil {
		return err
	}

	if err := s.items.Delete(ctx, id); err != nil {
		return err
	}

	// 清除缓存
	s.invalidateItemCaches(ctx, id)

	return nil
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
	items, total, err := s.items.List(ctx, repository.ServiceItemListOptions{
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
	Category    *string                       `form:"category"`
	SubCategory *model.ServiceItemSubCategory `form:"subCategory"`
	GameID      *uint64                       `form:"gameId"`
	PlayerID    *uint64                       `form:"playerId"`
	IsActive    *bool                         `form:"isActive"`
	Page        int                           `form:"page"`
	PageSize    int                           `form:"pageSize"`
}

// ServiceItemListResponse 服务项目列表响应
type ServiceItemListResponse struct {
	Items []ServiceItemDTO `json:"items"`
	Total int64            `json:"total"`
}

// GetGiftList 获取礼物列表（用户端
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

	// 获取游戏名称（使用缓存）
	if item.GameID != nil {
		game, err := s.getCachedGame(ctx, *item.GameID)
		if err == nil {
			dto.GameName = game.Name
		}
	}

	// 获取陪玩师昵
	if item.PlayerID != nil {
		player, err := s.players.Get(ctx, *item.PlayerID)
		if err == nil {
			dto.PlayerNickname = player.Nickname
		}
	}

	return dto
}

// BatchUpdateStatusRequest 批量更新状态请
type BatchUpdateStatusRequest struct {
	IDs      []uint64 `json:"ids" binding:"required"`
	IsActive bool     `json:"isActive"`
}

// BatchUpdateStatus 批量更新状
func (s *ServiceItemService) BatchUpdateStatus(ctx context.Context, req BatchUpdateStatusRequest) error {
	if len(req.IDs) == 0 {
		return apierr.BadRequest("未提供项目ID")
	}
	if err := s.items.BatchUpdateStatus(ctx, req.IDs, req.IsActive); err != nil {
		return err
	}
	// 清除缓存
	s.invalidateItemCaches(ctx, 0)
	return nil
}

// BatchUpdatePriceRequest 批量更新价格请求
type BatchUpdatePriceRequest struct {
	IDs            []uint64 `json:"ids" binding:"required"`
	BasePriceCents int64    `json:"basePriceCents" binding:"required,min=0"`
}

// BatchUpdatePrice 批量更新价格
func (s *ServiceItemService) BatchUpdatePrice(ctx context.Context, req BatchUpdatePriceRequest) error {
	if len(req.IDs) == 0 {
		return apierr.BadRequest("未提供项目ID")
	}
	if err := s.items.BatchUpdatePrice(ctx, req.IDs, req.BasePriceCents); err != nil {
		return err
	}
	// 清除缓存
	s.invalidateItemCaches(ctx, 0)
	return nil
}

// ============================================================================
// 新增批量操作：删除、更新佣金比例
// ============================================================================

// BatchOperationResponse 批量操作响应
type BatchOperationResponse struct {
	SuccessCount int                       `json:"success_count"`
	FailedCount  int                       `json:"failed_count"`
	TotalCount   int                       `json:"total_count"`
	FailedItems  []BatchOperationErrorItem `json:"failed_items,omitempty"`
	SuccessItems []uint64                  `json:"success_items,omitempty"`
}

// BatchOperationErrorItem 单个操作错误详情
type BatchOperationErrorItem struct {
	ID      uint64 `json:"id"`
	Message string `json:"message"`
}

// BatchDeleteItemsRequest 批量删除服务项目请求
type BatchDeleteItemsRequest struct {
	ItemIDs []uint64 `json:"itemIds" binding:"required,min=1,max=100"`
}

// BatchDeleteItems 批量删除服务项目
// 检查是否有订单使用这些服务项目，如果有则拒绝删除
func (s *ServiceItemService) BatchDeleteItems(ctx context.Context, req BatchDeleteItemsRequest) (*BatchOperationResponse, error) {
	if len(req.ItemIDs) == 0 {
		return nil, apierr.BadRequest("itemIds is required")
	}
	if len(req.ItemIDs) > 100 {
		return nil, apierr.BadRequest("maximum 100 items per batch")
	}

	response := &BatchOperationResponse{
		TotalCount:   len(req.ItemIDs),
		SuccessItems: make([]uint64, 0),
		FailedItems:  make([]BatchOperationErrorItem, 0),
	}

	// Order reference check to be implemented: check if any orders use these service items
	// Requires injecting OrderRepository and querying orders by service_item_id
	// For now, implements basic deletion logic without order reference checking

	for _, itemID := range req.ItemIDs {
		// 检查项目是否存在
		_, err := s.items.Get(ctx, itemID)
		if err != nil {
			if err == ErrNotFound {
				response.FailedItems = append(response.FailedItems, BatchOperationErrorItem{
					ID:      itemID,
					Message: "service item not found",
				})
				response.FailedCount++
				continue
			}
			response.FailedItems = append(response.FailedItems, BatchOperationErrorItem{
				ID:      itemID,
				Message: "get service item failed",
			})
			response.FailedCount++
			continue
		}

		// 删除前检查是否存在订单引用
		orderRefCount, err := s.items.CountOrderReferences(ctx, itemID)
		if err != nil {
			response.FailedItems = append(response.FailedItems, BatchOperationErrorItem{
				ID:      itemID,
				Message: "check order references failed",
			})
			response.FailedCount++
			continue
		}
		if orderRefCount > 0 {
			response.FailedItems = append(response.FailedItems, BatchOperationErrorItem{
				ID:      itemID,
				Message: fmt.Sprintf("service item is referenced by %d order(s)", orderRefCount),
			})
			response.FailedCount++
			continue
		}

		// 删除项目
		err = s.items.Delete(ctx, itemID)
		if err != nil {
			response.FailedItems = append(response.FailedItems, BatchOperationErrorItem{
				ID:      itemID,
				Message: "delete failed",
			})
			response.FailedCount++
			continue
		}

		response.SuccessItems = append(response.SuccessItems, itemID)
		response.SuccessCount++
	}

	// 清除缓存
	if response.SuccessCount > 0 {
		s.invalidateItemCaches(ctx, 0)
	}

	return response, nil
}

// BatchUpdateItemCommissionRequest 批量更新佣金比例请求
type BatchUpdateItemCommissionRequest struct {
	ItemIDs        []uint64 `json:"itemIds" binding:"required,min=1,max=100"`
	CommissionRate float64  `json:"commissionRate" binding:"required,min=0,max=1"`
}

// BatchUpdateItemCommission 批量更新服务项目佣金比例
func (s *ServiceItemService) BatchUpdateItemCommission(ctx context.Context, req BatchUpdateItemCommissionRequest) (*BatchOperationResponse, error) {
	if len(req.ItemIDs) == 0 {
		return nil, apierr.BadRequest("itemIds is required")
	}
	if len(req.ItemIDs) > 100 {
		return nil, apierr.BadRequest("maximum 100 items per batch")
	}
	if req.CommissionRate < 0 || req.CommissionRate > 1 {
		return nil, apierr.BadRequest("commission rate must be between 0 and 1")
	}

	response := &BatchOperationResponse{
		TotalCount:   len(req.ItemIDs),
		SuccessItems: make([]uint64, 0),
		FailedItems:  make([]BatchOperationErrorItem, 0),
	}

	// 使用批量更新提高性能
	err := s.items.BatchUpdateCommission(ctx, req.ItemIDs, req.CommissionRate)
	if err != nil {
		return nil, err
	}

	// 批量更新成功，所有项目都成功
	response.SuccessItems = req.ItemIDs
	response.SuccessCount = len(req.ItemIDs)

	// 清除缓存
	s.invalidateItemCaches(ctx, 0)

	return response, nil
}

// BatchUpdateItemStatusRequest 批量更新服务项目状态请求
type BatchUpdateItemStatusRequest struct {
	ItemIDs  []uint64 `json:"itemIds" binding:"required,min=1,max=100"`
	IsActive bool     `json:"isActive"`
}

// BatchUpdateItemStatus 批量更新服务项目状态（启用/禁用）
func (s *ServiceItemService) BatchUpdateItemStatus(ctx context.Context, req BatchUpdateItemStatusRequest) (*BatchOperationResponse, error) {
	if len(req.ItemIDs) == 0 {
		return nil, apierr.BadRequest("itemIds is required")
	}
	if len(req.ItemIDs) > 100 {
		return nil, apierr.BadRequest("maximum 100 items per batch")
	}

	response := &BatchOperationResponse{
		TotalCount:   len(req.ItemIDs),
		SuccessItems: make([]uint64, 0),
		FailedItems:  make([]BatchOperationErrorItem, 0),
	}

	// 使用批量更新提高性能
	err := s.items.BatchUpdateStatus(ctx, req.ItemIDs, req.IsActive)
	if err != nil {
		return nil, err
	}

	// 批量更新成功，所有项目都成功
	response.SuccessItems = req.ItemIDs
	response.SuccessCount = len(req.ItemIDs)

	// 清除缓存
	s.invalidateItemCaches(ctx, 0)

	return response, nil
}
