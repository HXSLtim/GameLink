package admin

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/service/item"
	apierr "gamelink/pkg/apierr"
)

// ServiceItem 服务项目模型（类型别名）
type ServiceItem = model.ServiceItem

// RegisterServiceItemRoutes 注册管理端服务项目管理路由
func RegisterServiceItemRoutes(router gin.IRouter, svc *item.ServiceItemService) {
	group := router.Group("/service-items")
	{
		group.POST("", func(c *gin.Context) { createServiceItemHandler(c, svc) })
		group.GET("", func(c *gin.Context) { listServiceItemsHandler(c, svc) })
		group.GET("/:id", func(c *gin.Context) { getServiceItemHandler(c, svc) })
		group.PUT("/:id", func(c *gin.Context) { updateServiceItemHandler(c, svc) })
		group.DELETE("/:id", func(c *gin.Context) { deleteServiceItemHandler(c, svc) })
		group.POST("/batch-update-status", func(c *gin.Context) { batchUpdateStatusHandler(c, svc) })
		group.POST("/batch-update-price", func(c *gin.Context) { batchUpdatePriceHandler(c, svc) })
	}
}

// createServiceItemHandler 创建服务项目
// @Summary      创建服务项目
// @Description  管理员创建服务项目（护航服务或礼物）
// @Tags         Admin - ServiceItem
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                                true  "Bearer {token}"
// @Param        request        body      item.CreateServiceItemRequest  true  "服务项目信息"
// @Success      200            {object}  model.APIResponse[ServiceItem]
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/service-items [post]
func createServiceItemHandler(c *gin.Context, svc *item.ServiceItemService) {
	var req item.CreateServiceItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	item, err := svc.CreateServiceItem(c.Request.Context(), req)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, *item)
}

// listServiceItemsHandler 获取服务项目列表
// @Summary      获取服务项目列表
// @Description  API endpoint// @Tags         Admin - ServiceItem
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        category       query     string  false  "分类"
// @Param        subCategory    query     string  false  "Sub category (solo/team/gift)"
// @Param        gameId         query     int     false  "游戏ID"
// @Param        isActive       query    bool         false  "Is active"// @Param        page           query     int     false  "页码"
// @Param        pageSize       query     int     false  "每页数量"
// @Success      200            {object}  model.APIResponse[item.ServiceItemListResponse]
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/service-items [get]
func listServiceItemsHandler(c *gin.Context, svc *item.ServiceItemService) {
	var req item.ListServiceItemsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	resp, err := svc.ListServiceItems(c.Request.Context(), req)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, *resp)
}

// getServiceItemHandler 获取服务项目详情
// @Summary      获取服务项目详情
// @Description  API endpoint// @Tags         Admin - ServiceItem
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      int     true  "服务项目ID"
// @Success      200            {object}  model.APIResponse[item.ServiceItemDTO]
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/service-items/{id} [get]
func getServiceItemHandler(c *gin.Context, svc *item.ServiceItemService) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	resp, err := svc.GetServiceItem(c.Request.Context(), id)
	if err != nil {
		if err == item.ErrNotFound {
			respondError(c, apierr.NotFound(apierr.ErrServiceItemNotFound))
			return
		}
		respondError(c, err)
		return
	}

	respondSuccess(c, *resp)
}

// updateServiceItemHandler 更新服务项目
// @Summary      更新服务项目
// @Description  API endpoint// @Tags         Admin - ServiceItem
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                                 true  "Bearer {token}"
// @Param        id             path      int                                    true  "服务项目ID"
// @Param        request        body      item.UpdateServiceItemRequest  true  "更新信息"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/service-items/{id} [put]
func updateServiceItemHandler(c *gin.Context, svc *item.ServiceItemService) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req item.UpdateServiceItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	err := svc.UpdateServiceItem(c.Request.Context(), id, req)
	if err != nil {
		respondError(c, err)
		return
	}

	respondMsg(c, "Service item updated successfully")
}

// deleteServiceItemHandler 删除服务项目
// @Summary      删除服务项目
// @Description  API endpoint// @Tags         Admin - ServiceItem
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      int     true  "服务项目ID"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/service-items/{id} [delete]
func deleteServiceItemHandler(c *gin.Context, svc *item.ServiceItemService) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	err := svc.DeleteServiceItem(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondMsg(c, "Service item deleted successfully")
}

// batchUpdateStatusHandler 批量更新状
// @Summary      批量更新状
// @Description  API endpoint// @Tags         Admin - ServiceItem
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                                  true  "Bearer {token}"
// @Param        request        body      item.BatchUpdateStatusRequest  true  "批量更新请求"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/service-items/batch-update-status [post]
func batchUpdateStatusHandler(c *gin.Context, svc *item.ServiceItemService) {
	var req item.BatchUpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	err := svc.BatchUpdateStatus(c.Request.Context(), req)
	if err != nil {
		respondError(c, err)
		return
	}

	respondMsg(c, "Status updated successfully")
}

// batchUpdatePriceHandler 批量更新价格
// @Summary      批量更新价格
// @Description  API endpoint// @Tags         Admin - ServiceItem
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                                 true  "Bearer {token}"
// @Param        request        body      item.BatchUpdatePriceRequest  true  "批量更新请求"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/service-items/batch-update-price [post]
func batchUpdatePriceHandler(c *gin.Context, svc *item.ServiceItemService) {
	var req item.BatchUpdatePriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	err := svc.BatchUpdatePrice(c.Request.Context(), req)
	if err != nil {
		respondError(c, err)
		return
	}

	respondMsg(c, "Price updated successfully")
}
