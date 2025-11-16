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

// @Description  API endpoint// @Tags         User - Players
// @Accept       json
// @Produce      json
// @Param        gameId      query     int     false  "Game ID"
// @Param        minPrice    query     int     false  "Min price (cents)"
// @Param        maxPrice    query     int     false  "Max price (cents)"
// @Param        minRating   query     number  false  "Min rating"
// @Param        onlineOnly  query     bool    false  "Online only"
// @Param        sortBy      query     string  false  "Sort by" Enums(price, rating, orders)
// @Param        page        query     int     false  "Page number"
// @Param        pageSize    query     int     false  "Page size"
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

// @Description  API endpoint// @Accept       json
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
