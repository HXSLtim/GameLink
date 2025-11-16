package order

import (
	"context"

	"gamelink/internal/model"
)

// validateCreateOrder 用于校验创建订单所需的基础数据，并返回陪玩师信息
func (s *OrderService) validateCreateOrder(ctx context.Context, req CreateOrderRequest) (*model.Player, error) {
	player, err := s.players.Get(ctx, req.PlayerID)
	if err != nil {
		return nil, err
	}

	if _, err := s.games.Get(ctx, req.GameID); err != nil {
		return nil, err
	}

	return player, nil
}
