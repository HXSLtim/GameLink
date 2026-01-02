package user

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/service/gift"
	"gamelink/internal/service/item"
	"gamelink/pkg/apierr"
)

// ServiceItemListResponse 服务列表响应（类型别名）
type ServiceItemListResponse = item.ServiceItemListResponse

// GiftOrderResponse 礼物订单响应（类型别名）
type GiftOrderResponse = gift.GiftOrderResponse

// Swagger-friendly envelopes to avoid generics in swag annotations
type ServiceItemListAPIResponseSwagger struct {
	Success    bool                    `json:"success"`
	Code       int                     `json:"code"`
	Message    string                  `json:"message"`
	Data       ServiceItemListResponse `json:"data"`
	Pagination *model.Pagination       `json:"pagination,omitempty"`
	TraceID    string                  `json:"traceId,omitempty"`
}

type GiftOrderAPIResponseSwagger struct {
	Success    bool              `json:"success"`
	Code       int               `json:"code"`
	Message    string            `json:"message"`
	Data       GiftOrderResponse `json:"data"`
	Pagination *model.Pagination `json:"pagination,omitempty"`
	TraceID    string            `json:"traceId,omitempty"`
}

type SentGiftsAPIResponseSwagger struct {
	Success    bool              `json:"success"`
	Code       int               `json:"code"`
	Message    string            `json:"message"`
	Data       map[string]any    `json:"data"`
	Pagination *model.Pagination `json:"pagination,omitempty"`
	TraceID    string            `json:"traceId,omitempty"`
}

// RegisterGiftRoutes Register user gift routes
func RegisterGiftRoutes(router gin.IRouter, giftSvc *gift.GiftService, itemSvc *item.ServiceItemService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/gifts")
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
// @Success      200       {object}  ServiceItemListAPIResponseSwagger
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

// sendGiftHandler 赠送礼物
// @Summary      赠送礼物
// @Description  赠送礼物给陪玩师
// @Tags         User - Gift
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      gift.SendGiftRequest  true  "赠送礼物请求"
// @Success      200      {object}  GiftOrderAPIResponseSwagger
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
		if err == gift.ErrValidation || err == gift.ErrInvalidGiftItem {
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
// @Success      200       {object}  SentGiftsAPIResponseSwagger
// @Failure      400       {object}  apierr.APIError
// @Failure      401       {object}  apierr.APIError
// @Failure      500       {object}  apierr.APIError
// @Router       /user/gifts/sent [get]
func getSentGiftsHandler(c *gin.Context, svc *gift.GiftService) {
	userID := getUserIDFromContext(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// Sent gifts query will be implemented when gift service has GetUserSentGifts method
	// For now, return empty data
	_ = userID
	_ = page
	_ = pageSize

	respondSuccess(c, "OK", map[string]interface{}{"gifts": []interface{}{}, "total": 0})
}
