package admin

import (
	"context"
	"fmt"

	"gamelink/internal/model"
	"gamelink/pkg/apierr"
)

// BatchUpdateGamesStatusWithResponse 批量更新游戏状态（启用/禁用），返回详细响应。
func (s *AdminService) BatchUpdateGamesStatusWithResponse(ctx context.Context, ids []uint64, isActive bool) (*BatchOperationResponse, error) {
	response := &BatchOperationResponse{
		TotalCount:   len(ids),
		SuccessItems: make([]uint64, 0),
		FailedItems:  make([]BatchErrorItem, 0),
	}

	if len(ids) == 0 {
		response.FailedCount = len(ids)
		return response, apierr.BadRequest("no game ids provided")
	}

	for _, gameID := range ids {
		// 验证游戏是否存在
		_, err := s.games.Get(ctx, gameID)
		if err != nil {
			if apierr.IsNotFound(err) {
				response.FailedItems = append(response.FailedItems, BatchErrorItem{
					ID:      gameID,
					Message: "game not found",
				})
				response.FailedCount++
				continue
			}
			response.FailedItems = append(response.FailedItems, BatchErrorItem{
				ID:      gameID,
				Message: fmt.Sprintf("get game failed: %v", err),
			})
			response.FailedCount++
			continue
		}

		// 使用 repository 的批量更新方法（单次调用更高效）
		// 但为了详细响应，这里逐个验证后批量更新
		response.SuccessItems = append(response.SuccessItems, gameID)
		response.SuccessCount++
	}

	if response.SuccessCount > 0 {
		_, err := s.games.BatchUpdateStatus(ctx, response.SuccessItems, isActive)
		if err != nil {
			// 批量更新失败，将所有成功的标记为失败
			for _, id := range response.SuccessItems {
				response.FailedItems = append(response.FailedItems, BatchErrorItem{
					ID:      id,
					Message: fmt.Sprintf("batch update failed: %v", err),
				})
			}
			response.FailedCount += response.SuccessCount
			response.SuccessCount = 0
			response.SuccessItems = make([]uint64, 0)
			return response, nil
		}

		s.invalidateCache(ctx, cacheKeyGames)
		for _, id := range response.SuccessItems {
			s.appendLogAsync(ctx, string(model.OpEntityGame), id, string(model.OpActionUpdate), map[string]any{"is_active": isActive})
		}
	}

	return response, nil
}

// BatchUpdateGamesCategory 批量更新游戏分类。
func (s *AdminService) BatchUpdateGamesCategory(ctx context.Context, ids []uint64, category string) (*BatchOperationResponse, error) {
	response := &BatchOperationResponse{
		TotalCount:   len(ids),
		SuccessItems: make([]uint64, 0),
		FailedItems:  make([]BatchErrorItem, 0),
	}

	if len(ids) == 0 {
		response.FailedCount = len(ids)
		return response, apierr.BadRequest("no game ids provided")
	}

	// 验证分类名称
	category = trim(category)
	if category == "" {
		response.FailedCount = len(ids)
		for _, id := range ids {
			response.FailedItems = append(response.FailedItems, BatchErrorItem{
				ID:      id,
				Message: "category cannot be empty",
			})
		}
		return response, nil
	}

	for _, gameID := range ids {
		// 验证游戏是否存在
		game, err := s.games.Get(ctx, gameID)
		if err != nil {
			if apierr.IsNotFound(err) {
				response.FailedItems = append(response.FailedItems, BatchErrorItem{
					ID:      gameID,
					Message: "game not found",
				})
				response.FailedCount++
				continue
			}
			response.FailedItems = append(response.FailedItems, BatchErrorItem{
				ID:      gameID,
				Message: fmt.Sprintf("get game failed: %v", err),
			})
			response.FailedCount++
			continue
		}

		// 记录旧分类用于日志
		oldCategory := game.Category

		// 更新分类
		game.Category = category
		err = s.games.Update(ctx, game)
		if err != nil {
			response.FailedItems = append(response.FailedItems, BatchErrorItem{
				ID:      gameID,
				Message: fmt.Sprintf("update category failed: %v", err),
			})
			response.FailedCount++
			continue
		}

		response.SuccessItems = append(response.SuccessItems, gameID)
		response.SuccessCount++

		// 记录日志
		s.appendLogAsync(ctx, string(model.OpEntityGame), gameID, string(model.OpActionUpdate), map[string]any{
			"category":     category,
			"old_category": oldCategory,
		})
	}

	if response.SuccessCount > 0 {
		s.invalidateCache(ctx, cacheKeyGames)
	}

	return response, nil
}
