package user

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/favorite"
	"gamelink/pkg/apierr"
)

// FavoriteHandler 收藏处理器
type FavoriteHandler struct {
	favoriteRepo *favorite.Repository
	playerRepo   repository.PlayerRepository
}

// NewFavoriteHandler 创建收藏处理器
func NewFavoriteHandler(favoriteRepo *favorite.Repository, playerRepo repository.PlayerRepository) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteRepo: favoriteRepo,
		playerRepo:   playerRepo,
	}
}

// RegisterFavoriteRoutes 注册收藏路由
func RegisterFavoriteRoutes(router gin.IRouter, h *FavoriteHandler, authMiddleware gin.HandlerFunc) {
	group := router.Group("/favorites/players")
	group.Use(authMiddleware)
	group.GET("", h.listFavorites)
	group.POST("/:id", h.addFavorite)
	group.DELETE("/:id", h.removeFavorite)
	group.GET("/:id/check", h.checkFavorite)
}

// FavoritePlayerDTO 收藏的陪玩师信息
type FavoritePlayerDTO struct {
	ID              uint64  `json:"id"`
	PlayerID        uint64  `json:"playerId"`
	Nickname        string  `json:"nickname"`
	AvatarURL       string  `json:"avatarUrl"`
	Bio             string  `json:"bio"`
	Rank            string  `json:"rank"`
	RatingAverage   float32 `json:"ratingAverage"`
	HourlyRateCents int64   `json:"hourlyRateCents"`
	CreatedAt       string  `json:"createdAt"`
}

// listFavorites 获取收藏列表
// @Summary      获取收藏列表
// @Description  获取用户收藏的陪玩师列表
// @Tags         User - Favorites
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        page      query     int  false  "页码" default(1)
// @Param        pageSize  query     int  false  "每页数量" default(20)
// @Success      200  {object}  resp.SuccessResponse
// @Failure      401  {object}  apierr.APIError
// @Failure      500  {object}  apierr.APIError
// @Router       /user/favorites/players [get]
func (h *FavoriteHandler) listFavorites(c *gin.Context) {
	userID := resp.GetUserID(c)
	page, pageSize := parsePagination(c)

	favorites, total, err := h.favoriteRepo.ListByUserID(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		resp.Error(c, apierr.InternalError("获取收藏列表失败").WithDetails(err.Error()))
		return
	}

	// 转换为 DTO
	items := make([]FavoritePlayerDTO, 0, len(favorites))
	for _, fav := range favorites {
		if fav.Player == nil {
			continue
		}
		items = append(items, FavoritePlayerDTO{
			ID:              fav.ID,
			PlayerID:        fav.PlayerID,
			Nickname:        fav.Player.Nickname,
			Bio:             fav.Player.Bio,
			Rank:            fav.Player.Rank,
			RatingAverage:   fav.Player.RatingAverage,
			HourlyRateCents: fav.Player.HourlyRateCents,
			CreatedAt:       fav.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	resp.List(c, items, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}

// addFavorite 添加收藏
// @Summary      添加收藏
// @Description  收藏陪玩师
// @Tags         User - Favorites
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id  path  uint64  true  "陪玩师ID"
// @Success      200  {object}  resp.SuccessResponse
// @Failure      400  {object}  apierr.APIError
// @Failure      401  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Failure      409  {object}  apierr.APIError
// @Failure      500  {object}  apierr.APIError
// @Router       /user/favorites/players/{id} [post]
func (h *FavoriteHandler) addFavorite(c *gin.Context) {
	userID := resp.GetUserID(c)

	playerID, err := resp.ParseUintParam(c, "id")
	if err != nil {
		resp.Error(c, apierr.BadRequest("无效的陪玩师ID"))
		return
	}

	// 检查陪玩师是否存在
	_, err = h.playerRepo.Get(c.Request.Context(), playerID)
	if err != nil {
		if err == repository.ErrNotFound {
			resp.Error(c, apierr.NotFound("陪玩师不存在"))
			return
		}
		resp.Error(c, apierr.InternalError("查询陪玩师失败"))
		return
	}

	// 检查是否已收藏
	exists, err := h.favoriteRepo.Exists(c.Request.Context(), userID, playerID)
	if err != nil {
		resp.Error(c, apierr.InternalError("检查收藏状态失败"))
		return
	}
	if exists {
		resp.Error(c, apierr.Conflict("已收藏该陪玩师"))
		return
	}

	// 创建收藏
	fav := &model.Favorite{
		UserID:   userID,
		PlayerID: playerID,
	}
	if err := h.favoriteRepo.Create(c.Request.Context(), fav); err != nil {
		resp.Error(c, apierr.InternalError("添加收藏失败").WithDetails(err.Error()))
		return
	}

	resp.Created(c, map[string]interface{}{
		"id":       fav.ID,
		"playerId": playerID,
	})
}

// removeFavorite 取消收藏
// @Summary      取消收藏
// @Description  取消收藏陪玩师
// @Tags         User - Favorites
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id  path  uint64  true  "陪玩师ID"
// @Success      200  {object}  resp.SuccessResponse
// @Failure      400  {object}  apierr.APIError
// @Failure      401  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Failure      500  {object}  apierr.APIError
// @Router       /user/favorites/players/{id} [delete]
func (h *FavoriteHandler) removeFavorite(c *gin.Context) {
	userID := resp.GetUserID(c)

	playerID, err := resp.ParseUintParam(c, "id")
	if err != nil {
		resp.Error(c, apierr.BadRequest("无效的陪玩师ID"))
		return
	}

	if err := h.favoriteRepo.Delete(c.Request.Context(), userID, playerID); err != nil {
		if err == repository.ErrNotFound {
			resp.Error(c, apierr.NotFound("未收藏该陪玩师"))
			return
		}
		resp.Error(c, apierr.InternalError("取消收藏失败").WithDetails(err.Error()))
		return
	}

	resp.Deleted(c)
}

// checkFavorite 检查是否已收藏
// @Summary      检查收藏状态
// @Description  检查是否已收藏某陪玩师
// @Tags         User - Favorites
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id  path  uint64  true  "陪玩师ID"
// @Success      200  {object}  resp.SuccessResponse
// @Failure      400  {object}  apierr.APIError
// @Failure      401  {object}  apierr.APIError
// @Failure      500  {object}  apierr.APIError
// @Router       /user/favorites/players/{id}/check [get]
func (h *FavoriteHandler) checkFavorite(c *gin.Context) {
	userID := resp.GetUserID(c)

	playerID, err := resp.ParseUintParam(c, "id")
	if err != nil {
		resp.Error(c, apierr.BadRequest("无效的陪玩师ID"))
		return
	}

	exists, err := h.favoriteRepo.Exists(c.Request.Context(), userID, playerID)
	if err != nil {
		resp.Error(c, apierr.InternalError("检查收藏状态失败"))
		return
	}

	resp.OK(c, map[string]interface{}{
		"isFavorite": exists,
		"playerId":   playerID,
	})
}
