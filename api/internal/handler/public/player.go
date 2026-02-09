// Package public provides public API handlers for player listing without authentication.
package public

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
)

// PlayerHandler 公共陪玩师处理器
type PlayerHandler struct {
	players repository.PlayerRepository
	users   repository.UserRepository
}

// NewPlayerHandler 创建公共陪玩师处理器
func NewPlayerHandler(players repository.PlayerRepository, users repository.UserRepository) *PlayerHandler {
	return &PlayerHandler{
		players: players,
		users:   users,
	}
}

// PublicPlayerInfo 公开的陪玩师信息
type PublicPlayerInfo struct {
	ID                 uint64  `json:"id"`
	UserID             uint64  `json:"userId"`
	Nickname           string  `json:"nickname"`
	Avatar             string  `json:"avatar"`
	Bio                string  `json:"bio"`
	Rank               string  `json:"rank"`
	MainGame           string  `json:"mainGame"`
	HourlyRateCents    int64   `json:"hourlyRateCents"`
	OrderCount         uint32  `json:"orderCount"`
	RatingAverage      float32 `json:"ratingAverage"`
	RatingCount        uint32  `json:"ratingCount"`
	IsOnline           bool    `json:"isOnline"`
	IsVerified         bool    `json:"isVerified"`
	OnlineStatus       string  `json:"onlineStatus"`
	VerificationStatus string  `json:"verificationStatus"`
}

// PlayerListResponse 陪玩师列表响应
type PlayerListResponse struct {
	Players  []PublicPlayerInfo `json:"players"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"pageSize"`
}

// ListPlayers 获取陪玩师列表（公开）
// @Summary 获取陪玩师列表
// @Description 获取已认证的陪玩师列表，无需登录
// @Tags 公共-陪玩师
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(20)
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} PlayerListResponse
// @Router /public/players [get]
func (h *PlayerHandler) ListPlayers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 只查询已认证的陪玩师
	status := model.VerificationVerified
	players, total, err := h.players.ListPagedWithFilter(c.Request.Context(), page, pageSize, keyword, &status)
	if err != nil {
		resp.Error(c, apierr.InternalError("获取陪玩师列表失败"))
		return
	}

	// 转换为公开信息
	result := make([]PublicPlayerInfo, 0, len(players))
	for _, p := range players {
		avatar := ""
		if p.User != nil {
			avatar = p.User.AvatarURL
		}
		mainGameName := ""
		if p.MainGame != nil {
			mainGameName = p.MainGame.Name
		}
		info := PublicPlayerInfo{
			ID:                 p.ID,
			UserID:             p.UserID,
			Nickname:           p.Nickname,
			Avatar:             avatar,
			Bio:                p.Bio,
			Rank:               p.Rank,
			MainGame:           mainGameName,
			HourlyRateCents:    p.HourlyRateCents,
			OrderCount:         p.OrderCount,
			RatingAverage:      p.RatingAverage,
			RatingCount:        p.RatingCount,
			IsOnline:           p.OnlineStatus == model.PlayerOnlineStatusOnline,
			IsVerified:         p.VerificationStatus == model.VerificationVerified,
			OnlineStatus:       string(p.OnlineStatus),
			VerificationStatus: string(p.VerificationStatus),
		}
		result = append(result, info)
	}

	resp.OK(c, PlayerListResponse{
		Players:  result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// GetPlayer 获取陪玩师详情（公开）
// @Summary 获取陪玩师详情
// @Description 获取指定陪玩师的公开信息，无需登录
// @Tags 公共-陪玩师
// @Accept json
// @Produce json
// @Param id path int true "陪玩师ID"
// @Success 200 {object} PublicPlayerInfo
// @Failure 404 {object}  apierr.APIError
// @Router /public/players/{id} [get]
func (h *PlayerHandler) GetPlayer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		resp.Error(c, apierr.BadRequest("无效的陪玩师ID"))
		return
	}

	player, err := h.players.Get(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			resp.Error(c, apierr.NotFound("陪玩师不存在"))
			return
		}
		resp.Error(c, apierr.InternalError("获取陪玩师信息失败"))
		return
	}

	// 只返回已认证的陪玩师
	if player.VerificationStatus != model.VerificationVerified {
		resp.Error(c, apierr.NotFound("陪玩师不存在"))
		return
	}

	avatar := ""
	if player.User != nil {
		avatar = player.User.AvatarURL
	}
	mainGameName := ""
	if player.MainGame != nil {
		mainGameName = player.MainGame.Name
	}

	info := PublicPlayerInfo{
		ID:                 player.ID,
		UserID:             player.UserID,
		Nickname:           player.Nickname,
		Avatar:             avatar,
		Bio:                player.Bio,
		Rank:               player.Rank,
		MainGame:           mainGameName,
		HourlyRateCents:    player.HourlyRateCents,
		OrderCount:         player.OrderCount,
		RatingAverage:      player.RatingAverage,
		RatingCount:        player.RatingCount,
		IsOnline:           player.OnlineStatus == model.PlayerOnlineStatusOnline,
		IsVerified:         player.VerificationStatus == model.VerificationVerified,
		OnlineStatus:       string(player.OnlineStatus),
		VerificationStatus: string(player.VerificationStatus),
	}

	resp.OK(c, info)
}

// ListFeaturedPlayers 获取精选陪玩师列表（公开）
// @Summary 获取精选陪玩师列表
// @Description 获取评分高、订单多的精选陪玩师，无需登录
// @Tags 公共-陪玩师
// @Accept json
// @Produce json
// @Param limit query int false "返回数量" default(10)
// @Success 200 {array} PublicPlayerInfo
// @Router /public/players/featured [get]
func (h *PlayerHandler) ListFeaturedPlayers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	// 只查询已认证的陪玩师
	status := model.VerificationVerified
	players, _, err := h.players.ListFeatured(c.Request.Context(), limit, &status)
	if err != nil {
		resp.Error(c, apierr.InternalError("获取精选陪玩师列表失败"))
		return
	}

	// 转换为公开信息
	result := make([]PublicPlayerInfo, 0, len(players))
	for _, p := range players {
		avatar := ""
		if p.User != nil {
			avatar = p.User.AvatarURL
		}
		mainGameName := ""
		if p.MainGame != nil {
			mainGameName = p.MainGame.Name
		}
		info := PublicPlayerInfo{
			ID:                 p.ID,
			UserID:             p.UserID,
			Nickname:           p.Nickname,
			Avatar:             avatar,
			Bio:                p.Bio,
			Rank:               p.Rank,
			MainGame:           mainGameName,
			HourlyRateCents:    p.HourlyRateCents,
			OrderCount:         p.OrderCount,
			RatingAverage:      p.RatingAverage,
			RatingCount:        p.RatingCount,
			IsOnline:           p.OnlineStatus == model.PlayerOnlineStatusOnline,
			IsVerified:         p.VerificationStatus == model.VerificationVerified,
			OnlineStatus:       string(p.OnlineStatus),
			VerificationStatus: string(p.VerificationStatus),
		}
		result = append(result, info)
	}

	resp.OK(c, result)
}

// RegisterRoutes 注册公共陪玩师路由
func (h *PlayerHandler) RegisterRoutes(rg *gin.RouterGroup) {
	players := rg.Group("/players")
	{
		players.GET("", h.ListPlayers)
		players.GET("/featured", h.ListFeaturedPlayers)
		players.GET("/:id", h.GetPlayer)
	}
}
