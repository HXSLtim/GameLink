package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

    "gamelink/internal/model"
    "gamelink/internal/service/gift"
    "gamelink/internal/service/item"
)

// RegisterGiftRoutes Register user gift routes
func RegisterGiftRoutes(router gin.IRouter, giftSvc *gift.GiftService, itemSvc *item.ServiceItemService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/user/gifts")
	group.Use(authMiddleware)
	{
		group.GET("", func(c *gin.Context) { listGiftsHandler(c, itemSvc) })
		group.POST("/send", func(c *gin.Context) { sendGiftHandler(c, giftSvc) })
		group.GET("/sent", func(c *gin.Context) { getSentGiftsHandler(c, giftSvc) })
	}
}

// listGiftsHandler 获取礼物列表
// @Summary      获取礼物列表
// @Description  用户浏览可赠送的礼物
// @Tags         User - Gift
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        page           query     int     false  "页码"
// @Param        pageSize       query     int     false  "每页数量"
// @Success      200            {object}  model.APIResponse[serviceitem.ServiceItemListResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /user/gifts [get]
func listGiftsHandler(c *gin.Context, svc *item.ServiceItemService) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	resp, err := svc.GetGiftList(c.Request.Context(), page, pageSize)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

    respondJSON(c, http.StatusOK, model.APIResponse[item.ServiceItemListResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    *resp,
	})
}

// sendGiftHandler 赠送礼�?
// @Summary      赠送礼�?
// @Description  用户给陪玩师赠送礼�?
// @Tags         User - Gift
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                  true  "Bearer {token}"
// @Param        request        body      gift.SendGiftRequest  true  "赠送礼物请�?
// @Success      200            {object}  model.APIResponse[gift.GiftOrderResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /user/gifts/send [post]
func sendGiftHandler(c *gin.Context, svc *gift.GiftService) {
	userID := getUserIDFromContext(c)

	var req gift.SendGiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := svc.SendGift(c.Request.Context(), userID, req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[gift.GiftOrderResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "Gift sent successfully",
		Data:    *resp,
	})
}

// getSentGiftsHandler 获取已赠送的礼物记录
// @Summary      获取已赠送的礼物记录
// @Description  用户查看自己赠送的礼物记录
// @Tags         User - Gift
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        page           query     int     false  "页码"
// @Param        pageSize       query     int     false  "每页数量"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /user/gifts/sent [get]
func getSentGiftsHandler(c *gin.Context, svc *gift.GiftService) {
	userID := getUserIDFromContext(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// TODO: 实现获取用户赠送的礼物记录
	_ = userID
	_ = page
	_ = pageSize

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    map[string]interface{}{"gifts": []interface{}{}, "total": 0},
	})
}
