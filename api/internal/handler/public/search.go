package public

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
)

// SearchHandler 搜索处理器
type SearchHandler struct {
	playerRepo repository.PlayerRepository
	gameRepo   repository.GameRepository
	userRepo   repository.UserRepository
}

// NewSearchHandler 创建搜索处理器
func NewSearchHandler(
	playerRepo repository.PlayerRepository,
	gameRepo repository.GameRepository,
	userRepo repository.UserRepository,
) *SearchHandler {
	return &SearchHandler{
		playerRepo: playerRepo,
		gameRepo:   gameRepo,
		userRepo:   userRepo,
	}
}

// RegisterSearchRoutes 注册搜索路由
func RegisterSearchRoutes(router gin.IRouter, h *SearchHandler) {
	router.GET("/search", h.search)
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Query    string `form:"q" binding:"required,min=1,max=50"`
	Type     string `form:"type"` // player, game, all (default: all)
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

// SearchResultItem 搜索结果项
type SearchResultItem struct {
	ID          uint64  `json:"id"`
	Type        string  `json:"type"` // player, game
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	ImageURL    string  `json:"imageUrl,omitempty"`
	Extra       any     `json:"extra,omitempty"`
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Players []SearchResultItem `json:"players,omitempty"`
	Games   []SearchResultItem `json:"games,omitempty"`
	Total   int64              `json:"total"`
}

// search 搜索
// @Summary      搜索
// @Description  搜索陪玩师和游戏
// @Tags         Public - Search
// @Accept       json
// @Produce      json
// @Param        q         query     string  true   "搜索关键词"
// @Param        type      query     string  false  "搜索类型: player, game, all" default(all)
// @Param        page      query     int     false  "页码" default(1)
// @Param        pageSize  query     int     false  "每页数量" default(20)
// @Success      200  {object}  resp.SuccessResponse
// @Failure      400  {object}  apierr.APIError
// @Failure      500  {object}  apierr.APIError
// @Router       /public/search [get]
func (h *SearchHandler) search(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.Error(c, apierr.BadRequest("搜索关键词不能为空"))
		return
	}

	// 默认值
	if req.Type == "" {
		req.Type = "all"
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 50 {
		req.PageSize = 20
	}

	result := SearchResponse{
		Players: []SearchResultItem{},
		Games:   []SearchResultItem{},
	}

	// 搜索陪玩师
	if req.Type == "all" || req.Type == "player" {
		players, err := h.searchPlayers(c, req.Query, req.Page, req.PageSize)
		if err == nil {
			result.Players = players
			result.Total += int64(len(players))
		}
	}

	// 搜索游戏
	if req.Type == "all" || req.Type == "game" {
		games, err := h.searchGames(c, req.Query, req.Page, req.PageSize)
		if err == nil {
			result.Games = games
			result.Total += int64(len(games))
		}
	}

	resp.OK(c, result)
}

// searchPlayers 搜索陪玩师
func (h *SearchHandler) searchPlayers(c *gin.Context, query string, page, pageSize int) ([]SearchResultItem, error) {
	// 获取所有陪玩师并过滤
	players, _, err := h.playerRepo.ListPaged(c.Request.Context(), 1, 1000)
	if err != nil {
		return nil, err
	}

	var results []SearchResultItem
	for _, p := range players {
		// 只搜索已审核通过的陪玩师
		if p.VerificationStatus != model.VerificationVerified {
			continue
		}

		// 模糊匹配昵称
		if !containsIgnoreCase(p.Nickname, query) && !containsIgnoreCase(p.Bio, query) {
			continue
		}

		// 获取用户头像
		var avatarURL string
		if user, err := h.userRepo.Get(c.Request.Context(), p.UserID); err == nil {
			avatarURL = user.AvatarURL
		}

		results = append(results, SearchResultItem{
			ID:          p.ID,
			Type:        "player",
			Name:        p.Nickname,
			Description: p.Bio,
			ImageURL:    avatarURL,
			Extra: map[string]any{
				"rank":            p.Rank,
				"ratingAverage":   p.RatingAverage,
				"hourlyRateCents": p.HourlyRateCents,
			},
		})
	}

	// 分页
	start := (page - 1) * pageSize
	if start >= len(results) {
		return []SearchResultItem{}, nil
	}
	end := start + pageSize
	if end > len(results) {
		end = len(results)
	}

	return results[start:end], nil
}

// searchGames 搜索游戏
func (h *SearchHandler) searchGames(c *gin.Context, query string, page, pageSize int) ([]SearchResultItem, error) {
	// 获取所有游戏并过滤
	games, _, err := h.gameRepo.ListPaged(c.Request.Context(), 1, 1000)
	if err != nil {
		return nil, err
	}

	var results []SearchResultItem
	for _, g := range games {
		// 只搜索启用的游戏
		if !g.IsActive {
			continue
		}

		// 模糊匹配名称
		if !containsIgnoreCase(g.Name, query) && !containsIgnoreCase(g.Description, query) {
			continue
		}

		results = append(results, SearchResultItem{
			ID:          g.ID,
			Type:        "game",
			Name:        g.Name,
			Description: g.Description,
			ImageURL:    g.IconURL,
		})
	}

	// 分页
	start := (page - 1) * pageSize
	if start >= len(results) {
		return []SearchResultItem{}, nil
	}
	end := start + pageSize
	if end > len(results) {
		end = len(results)
	}

	return results[start:end], nil
}

// containsIgnoreCase 忽略大小写的包含检查
func containsIgnoreCase(s, substr string) bool {
	if s == "" || substr == "" {
		return false
	}
	return contains(toLower(s), toLower(substr))
}

// toLower 转小写（简单实现）
func toLower(s string) string {
	b := make([]byte, len(s))
	for i := range s {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		} else {
			b[i] = c
		}
	}
	return string(b)
}

// contains 检查字符串包含
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
