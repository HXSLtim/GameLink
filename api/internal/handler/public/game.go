// Package public provides public API handlers for game listing without authentication.
package public

import (
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
)

// GameHandler 公共游戏处理器
type GameHandler struct {
	games          repository.GameRepository
	gameCategories repository.GameCategoryRepository
}

// NewGameHandler 创建公共游戏处理器
func NewGameHandler(games repository.GameRepository, gameCategories repository.GameCategoryRepository) *GameHandler {
	return &GameHandler{
		games:          games,
		gameCategories: gameCategories,
	}
}

// PublicGameInfo 公开的游戏信息
type PublicGameInfo struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	Category string `json:"category"`
}

// GameListResponse 游戏列表响应
type GameListResponse struct {
	Games    []PublicGameInfo `json:"games"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
}

// PublicGameCategory 公开的游戏分类信息
type PublicGameCategory struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	IconURL   string `json:"iconUrl,omitempty"`
	SortOrder int    `json:"sortOrder,omitempty"`
}

// GameCategoryListResponse 游戏分类列表响应
type GameCategoryListResponse struct {
	Categories []PublicGameCategory `json:"categories"`
	Total      int64                `json:"total"`
}

// ListCategories 获取游戏分类列表（公开）
// @Summary 获取游戏分类列表
// @Description 获取游戏分类列表，无需登录
// @Tags 公共-游戏
// @Accept json
// @Produce json
// @Success 200 {object} GameCategoryListResponse
// @Router /public/games/categories [get]
func (h *GameHandler) ListCategories(c *gin.Context) {
	// 优先使用 game_categories 表
	if h.gameCategories != nil {
		active := true
		categories, total, err := h.gameCategories.List(c.Request.Context(), repository.GameCategoryListOptions{
			Page:     1,
			PageSize: 200,
			IsActive: &active,
		})
		if err != nil {
			resp.Error(c, apierr.InternalError("获取游戏分类失败"))
			return
		}

		result := make([]PublicGameCategory, 0, len(categories))
		for _, item := range categories {
			result = append(result, PublicGameCategory{
				ID:        item.ID,
				Name:      item.Name,
				IconURL:   item.IconURL,
				SortOrder: item.SortOrder,
			})
		}

		resp.OK(c, GameCategoryListResponse{
			Categories: result,
			Total:      total,
		})
		return
	}

	// 兜底：从 games.category 去重构建分类
	games, err := h.games.List(c.Request.Context())
	if err != nil {
		resp.Error(c, apierr.InternalError("获取游戏分类失败"))
		return
	}
	seen := map[string]PublicGameCategory{}
	for _, game := range games {
		name := game.Category
		if name == "" {
			name = "未分类"
		}
		if _, ok := seen[name]; ok {
			continue
		}
		id := uint64(len(seen) + 1)
		if game.CategoryID != nil && *game.CategoryID > 0 {
			id = *game.CategoryID
		}
		seen[name] = PublicGameCategory{
			ID:   id,
			Name: name,
		}
	}

	result := make([]PublicGameCategory, 0, len(seen))
	for _, item := range seen {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	resp.OK(c, GameCategoryListResponse{
		Categories: result,
		Total:      int64(len(result)),
	})
}

// ListGames 获取游戏列表（公开）
// @Summary 获取游戏列表
// @Description 获取平台支持的游戏列表，无需登录
// @Tags 公共-游戏
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(50)
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} GameListResponse
// @Router /public/games [get]
func (h *GameHandler) ListGames(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	games, total, err := h.games.ListPagedWithFilter(c.Request.Context(), page, pageSize, keyword)
	if err != nil {
		resp.Error(c, apierr.InternalError("获取游戏列表失败"))
		return
	}

	// 转换为公开信息
	result := make([]PublicGameInfo, 0, len(games))
	for _, g := range games {
		info := PublicGameInfo{
			ID:       g.ID,
			Name:     g.Name,
			Icon:     g.IconURL,
			Category: g.Category,
		}
		result = append(result, info)
	}

	resp.OK(c, GameListResponse{
		Games:    result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// GetGame 获取游戏详情（公开）
// @Summary 获取游戏详情
// @Description 获取指定游戏的详细信息，无需登录
// @Tags 公共-游戏
// @Accept json
// @Produce json
// @Param id path int true "游戏ID"
// @Success 200 {object} PublicGameInfo
// @Failure 404 {object}  apierr.APIError
// @Router /public/games/{id} [get]
func (h *GameHandler) GetGame(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		resp.Error(c, apierr.BadRequest("无效的游戏ID"))
		return
	}

	game, err := h.games.Get(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			resp.Error(c, apierr.NotFound("游戏不存在"))
			return
		}
		resp.Error(c, apierr.InternalError("获取游戏信息失败"))
		return
	}

	info := PublicGameInfo{
		ID:       game.ID,
		Name:     game.Name,
		Icon:     game.IconURL,
		Category: game.Category,
	}

	resp.OK(c, info)
}

// RegisterRoutes 注册公共游戏路由
func (h *GameHandler) RegisterRoutes(rg *gin.RouterGroup) {
	games := rg.Group("/games")
	{
		games.GET("", h.ListGames)
		games.GET("/categories", h.ListCategories) // must be before /:id
		games.GET("/:id", h.GetGame)
	}
}
