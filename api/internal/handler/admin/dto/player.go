package dto

import (
	"time"

	"gamelink/internal/model"
)

// ==================== Response DTOs ====================

// PlayerResponse 陪玩师响应 DTO
type PlayerResponse struct {
	ID                 uint64                   `json:"id"`
	UserID             uint64                   `json:"userId"`
	Nickname           string                   `json:"nickname,omitempty"`
	Bio                string                   `json:"bio,omitempty"`
	Rank               string                   `json:"rank,omitempty"`
	RatingAverage      float32                  `json:"ratingAverage"`
	RatingCount        uint32                   `json:"ratingCount"`
	OrderCount         uint32                   `json:"orderCount"`
	HourlyRateCents    int64                    `json:"hourlyRateCents"`
	MainGameID         uint64                   `json:"mainGameId,omitempty"`
	VerificationStatus model.VerificationStatus `json:"verificationStatus"`
	OnlineStatus       model.PlayerOnlineStatus `json:"onlineStatus"`
	AcceptingOrders    bool                     `json:"acceptingOrders"`
	LastOnlineAt       *time.Time               `json:"lastOnlineAt,omitempty"`
	CreatedAt          time.Time                `json:"createdAt"`

	// 审核信息（管理后台可见）
	VerifiedAt   *time.Time `json:"verifiedAt,omitempty"`
	VerifiedBy   *uint64    `json:"verifiedBy,omitempty"`
	VerifyRemark string     `json:"verifyRemark,omitempty"`
	RejectReason string     `json:"rejectReason,omitempty"`

	// 关联简要信息
	User     *UserBrief `json:"user,omitempty"`
	MainGame *GameBrief `json:"mainGame,omitempty"`
}

// PlayerBrief 陪玩师简要信息（嵌入其他 DTO 用）
type PlayerBrief struct {
	ID       uint64 `json:"id"`
	UserID   uint64 `json:"userId"`
	Nickname string `json:"nickname"`
}

// GameBrief 游戏简要信息
type GameBrief struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

// PlayerListResponse 陪玩师列表响应
type PlayerListResponse struct {
	Items      []PlayerResponse `json:"items"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"pageSize"`
	TotalPages int              `json:"totalPages"`
}

// ==================== 转换函数 ====================

// ToPlayerResponse 将 model.Player 转换为 PlayerResponse
func ToPlayerResponse(player *model.Player) *PlayerResponse {
	if player == nil {
		return nil
	}

	resp := &PlayerResponse{
		ID:                 player.ID,
		UserID:             player.UserID,
		Nickname:           player.Nickname,
		Bio:                player.Bio,
		Rank:               player.Rank,
		RatingAverage:      player.RatingAverage,
		RatingCount:        player.RatingCount,
		OrderCount:         player.OrderCount,
		HourlyRateCents:    player.HourlyRateCents,
		MainGameID:         player.MainGameID,
		VerificationStatus: player.VerificationStatus,
		OnlineStatus:       player.OnlineStatus,
		AcceptingOrders:    player.AcceptingOrders,
		LastOnlineAt:       player.LastOnlineAt,
		CreatedAt:          player.CreatedAt,
		VerifiedAt:         player.VerifiedAt,
		VerifiedBy:         player.VerifiedBy,
		VerifyRemark:       player.VerifyRemark,
		RejectReason:       player.RejectReason,
	}

	// 关联用户信息
	if player.User != nil {
		resp.User = &UserBrief{
			ID:        player.User.ID,
			Name:      player.User.Name,
			AvatarURL: player.User.AvatarURL,
		}
	}

	// 关联游戏信息
	if player.MainGame != nil {
		resp.MainGame = &GameBrief{
			ID:   player.MainGame.ID,
			Name: player.MainGame.Name,
		}
	}

	return resp
}

// ToPlayerResponseList 批量转换
func ToPlayerResponseList(players []model.Player) []PlayerResponse {
	responses := make([]PlayerResponse, 0, len(players))
	for i := range players {
		if resp := ToPlayerResponse(&players[i]); resp != nil {
			responses = append(responses, *resp)
		}
	}
	return responses
}
