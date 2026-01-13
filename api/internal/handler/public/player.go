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
	RatingAverage      float32 `json:"ratingAverage"`
	RatingCount        uint32  `json:"ratingCount"`
	OnlineStatus       string  `json:"onlineStatus"`
	VerificationStatus string  `json:"verificationStatus"`
}

// PlayerListResponse 陪玩师列表响应
type PlayerListResponse struct {
	Players []PublicPlayerInfo `json:"players"`
	Total   int64              `json:"total"`
	Page    int                `json:"page"`
	PageSize int               `json:"pageSize"`
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
		info := PublicPlayerInfo{
			ID:                 p.ID,
			UserID:             p.UserID,
			Nickname:           p.Nickname,
			Avatar:             avatar,
			Bio:                p.Bio,
			Rank:               p.Rank,
			RatingAverage:      p.RatingAverage,
			RatingCount:        p.RatingCount,
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

	info := PublicPlayerInfo{
		ID:                 player.ID,
		UserID:             player.UserID,
		Nickname:           player.Nickname,
		Avatar:             avatar,
		Bio:                player.Bio,
		Rank:               player.Rank,
		RatingAverage:      player.RatingAverage,
		RatingCount:        player.RatingCount,
		OnlineStatus:       string(player.OnlineStatus),
		VerificationStatus: string(player.VerificationStatus),
	}

	resp.OK(c, info)
}

// RegisterRoutes 注册公共陪玩师路由
func (h *PlayerHandler) RegisterRoutes(rg *gin.RouterGroup) {
	players := rg.Group("/players")
	{
		players.GET("", h.ListPlayers)
		players.GET("/:id", h.GetPlayer)
	}
}
