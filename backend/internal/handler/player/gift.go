package player

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/service/gift"
)

// RegisterGiftRoutes 注册陪玩师端礼物路由
func RegisterGiftRoutes(router gin.IRouter, svc *gift.GiftService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/player/gifts")
	group.Use(authMiddleware)
	{
		group.GET("/received", func(c *gin.Context) { getReceivedGiftsHandler(c, svc) })
		group.GET("/stats", func(c *gin.Context) { getGiftStatsHandler(c, svc) })
	}
}

// getReceivedGiftsHandler 获取收到的礼�?
// @Summary      获取收到的礼�?
// @Description  陪玩师查看收到的礼物列表
// @Tags         Player - Gift
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        page           query     int     false  "页码"
// @Param        pageSize       query     int     false  "每页数量"
// @Success      200            {object}  model.APIResponse[gift.ReceivedGiftsResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /player/gifts/received [get]
func getReceivedGiftsHandler(c *gin.Context, svc *gift.GiftService) {
	userID := getUserIDFromContext(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// TODO: 需要从userID获取playerID
	// 暂时使用userID作为playerID
	playerID := userID

	resp, err := svc.GetPlayerReceivedGifts(c.Request.Context(), playerID, page, pageSize)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[gift.ReceivedGiftsResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    *resp,
	})
}

// getGiftStatsHandler 获取礼物统计
// @Summary      获取礼物统计
// @Description  陪玩师查看礼物收入统�?
// @Tags         Player - Gift
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  model.APIResponse[gift.GiftStatsResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /player/gifts/stats [get]
func getGiftStatsHandler(c *gin.Context, svc *gift.GiftService) {
	userID := getUserIDFromContext(c)

	// TODO: 需要从userID获取playerID
	playerID := userID

	resp, err := svc.GetGiftStats(c.Request.Context(), playerID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[gift.GiftStatsResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    *resp,
	})
}
