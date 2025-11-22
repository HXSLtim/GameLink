package user

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/apierr"
	"gamelink/internal/service/player"
)

// RegisterPlayerRoutes 注册用户端陪玩师路由
func RegisterPlayerRoutes(router gin.IRouter, svc *player.PlayerService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/user/players")
	{
		// 公开接口（不需要认证）
		group.GET("", func(c *gin.Context) { listPlayersHandler(c, svc) })
		group.GET("/:id", func(c *gin.Context) { getPlayerDetailHandler(c, svc) })
	}
}

// listPlayersHandler 获取陪玩师列表
// @Summary      获取陪玩师列表
// @Description  获取陪玩师列表，支持过滤和排序
// @Tags         User - Players
// @Accept       json
// @Produce      json
// @Param        gameId      query     uint64    false  "Game ID"
// @Param        minPrice    query     int64     false  "Min price (cents)"
// @Param        maxPrice    query     int64     false  "Max price (cents)"
// @Param        minRating   query     float32   false  "Min rating"
// @Param        onlineOnly  query     bool      false  "Online only"
// @Param        sortBy      query     string    false  "Sort by" Enums(price,rating,orders)
// @Param        page        query     int       false  "Page number" default(1)
// @Param        pageSize    query     int       false  "Page size" default(20)
// @Success      200         {object}  model.APIResponse[player.PlayerListResponse]
// @Failure      400         {object}  apierr.APIError
// @Failure      500         {object}  apierr.APIError
// @Router       /user/players [get]
func listPlayersHandler(c *gin.Context, svc *player.PlayerService) {
	var req player.PlayerListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("无效的查询参数").WithDetails(err.Error()))
		return
	}

	resp, err := svc.ListPlayers(c.Request.Context(), req)
	if err != nil {
		respondAPIError(c, apierr.InternalError("获取陪玩师列表失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "OK", *resp)
}

// @Description  API endpoint// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "陪玩师ID"
// @Success      200  {object}  model.APIResponse[player.PlayerDetailResponse]
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /user/players/{id} [get]
func getPlayerDetailHandler(c *gin.Context, svc *player.PlayerService) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidID))
		return
	}

	resp, err := svc.GetPlayerDetail(c.Request.Context(), id)
	if err != nil {
		if err == player.ErrNotFound {
			respondAPIError(c, apierr.NotFound(err.Error()))
			return
		}
		respondAPIError(c, apierr.InternalError("获取陪玩师详情失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "OK", *resp)
}
