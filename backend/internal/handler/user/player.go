package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/apierr"
	"gamelink/internal/model"
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

// listPlayersHandler 获取陪玩师列�?// @Summary      获取陪玩师列�?// @Description  获取陪玩师列表，支持筛选和排序
// @Tags         User - Players
// @Accept       json
// @Produce      json
// @Param        gameId      query     int     false  "游戏ID"
// @Param        minPrice    query     int     false  "最低价格（分）"
// @Param        maxPrice    query     int     false  "最高价格（分）"
// @Param        minRating   query     number  false  "最低评�?
// @Param        onlineOnly  query     bool    false  "仅在�?
// @Param        sortBy      query     string  false  "排序方式" Enums(price, rating, orders)
// @Param        page        query     int     false  "页码"
// @Param        pageSize    query     int     false  "每页数量"
// @Success      200         {object}  model.APIResponse[player.PlayerListResponse]
// @Failure      400         {object}  model.APIResponse[any]
// @Router       /user/players [get]
func listPlayersHandler(c *gin.Context, svc *player.PlayerService) {
	var req player.PlayerListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := svc.ListPlayers(c.Request.Context(), req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[player.PlayerListResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    *resp,
	})
}

// getPlayerDetailHandler 获取陪玩师详�?// @Summary      获取陪玩师详�?// @Description  获取陪玩师详细信息，包括评价和统�?// @Tags         User - Players
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "陪玩师ID"
// @Success      200  {object}  model.APIResponse[player.PlayerDetailResponse]
// @Failure      400  {object}  model.APIResponse[any]
// @Failure      404  {object}  model.APIResponse[any]
// @Router       /user/players/{id} [get]
func getPlayerDetailHandler(c *gin.Context, svc *player.PlayerService) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}

	resp, err := svc.GetPlayerDetail(c.Request.Context(), id)
	if err != nil {
		if err == player.ErrNotFound {
			respondError(c, http.StatusNotFound, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[player.PlayerDetailResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    *resp,
	})
}
