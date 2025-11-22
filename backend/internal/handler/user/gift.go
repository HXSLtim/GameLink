package user

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/apierr"
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
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        page      query     int  false  "页码" default(1)
// @Param        pageSize  query     int  false  "每页数量" default(20)
// @Success      200       {object}  model.APIResponse[item.ServiceItemListResponse]
// @Failure      400       {object}  apierr.APIError
// @Failure      401       {object}  apierr.APIError
// @Failure      500       {object}  apierr.APIError
// @Router       /user/gifts [get]
func listGiftsHandler(c *gin.Context, svc *item.ServiceItemService) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	resp, err := svc.GetGiftList(c.Request.Context(), page, pageSize)
	if err != nil {
		respondAPIError(c, apierr.InternalError("获取礼物列表失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "OK", *resp)
}

// sendGiftHandler 赠送礼�?
// @Summary      赠送礼�?
// @Description  赠送礼物给陪玩师
// @Tags         User - Gift
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      gift.SendGiftRequest  true  "赠送礼物请求"
// @Success      200      {object}  model.APIResponse[gift.GiftOrderResponse]
// @Failure      400      {object}  apierr.APIError
// @Failure      401      {object}  apierr.APIError
// @Failure      404      {object}  apierr.APIError
// @Failure      409      {object}  apierr.APIError
// @Failure      500      {object}  apierr.APIError
// @Router       /user/gifts/send [post]
func sendGiftHandler(c *gin.Context, svc *gift.GiftService) {
	userID := getUserIDFromContext(c)

	var req gift.SendGiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidJSONPayload).WithDetails(err.Error()))
		return
	}

	resp, err := svc.SendGift(c.Request.Context(), userID, req)
	if err != nil {
		if err == gift.ErrNotFound {
			respondAPIError(c, apierr.NotFound(err.Error()))
			return
		}
		if err == gift.ErrValidation {
			respondAPIError(c, apierr.BadRequest(err.Error()))
			return
		}
		if err == gift.ErrInvalidGiftItem {
			respondAPIError(c, apierr.BadRequest(err.Error()))
			return
		}
		respondAPIError(c, apierr.InternalError("赠送礼物失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "礼物赠送成功", *resp)
}

// getSentGiftsHandler 获取已赠送的礼物记录
// @Summary      获取已赠送的礼物记录
// @Description  用户查看自己赠送的礼物记录
// @Tags         User - Gift
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        page      query     int  false  "页码" default(1)
// @Param        pageSize  query     int  false  "每页数量" default(20)
// @Success      200       {object}  model.APIResponse[any]
// @Failure      400       {object}  apierr.APIError
// @Failure      401       {object}  apierr.APIError
// @Failure      500       {object}  apierr.APIError
// @Router       /user/gifts/sent [get]
func getSentGiftsHandler(c *gin.Context, svc *gift.GiftService) {
	userID := getUserIDFromContext(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// TODO: 实现获取用户赠送的礼物记录
	_ = userID
	_ = page
	_ = pageSize

	respondSuccess(c, "OK", map[string]interface{}{"gifts": []interface{}{}, "total": 0})
}
