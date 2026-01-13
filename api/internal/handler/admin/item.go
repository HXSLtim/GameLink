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

		// 新增批量操作路由
		group.POST("/batch/status", func(c *gin.Context) { batchUpdateItemStatusHandler(c, svc) })
		group.POST("/batch/delete", func(c *gin.Context) { batchDeleteItemsHandler(c, svc) })
		group.POST("/batch/commission", func(c *gin.Context) { batchUpdateItemCommissionHandler(c, svc) })
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
// @Success      200            {object}  ServiceItem
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
// @Success      200            {object}  item.ServiceItemListResponse
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
// @Success      200            {object}  item.ServiceItemDTO
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
// @Param        id             path      int  true  "服务项目ID"
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

// ============================================================================
// 新增批量操作 Handler
// ============================================================================

// BatchUpdateItemStatusRequest 批量更新服务项目状态请求
type BatchUpdateItemStatusRequest struct {
	ItemIDs  []uint64 `json:"itemIds" binding:"required,min=1,max=100"`
	IsActive bool     `json:"isActive"`
}

// BatchDeleteItemsRequest 批量删除服务项目请求
type BatchDeleteItemsRequest struct {
	ItemIDs []uint64 `json:"itemIds" binding:"required,min=1,max=100"`
}

// BatchUpdateItemCommissionRequest 批量更新佣金比例请求
type BatchUpdateItemCommissionRequest struct {
	ItemIDs        []uint64 `json:"itemIds" binding:"required,min=1,max=100"`
	CommissionRate float64  `json:"commissionRate" binding:"required,min=0,max=1"`
}

// batchUpdateItemStatusHandler 批量更新服务项目状态（启用/禁用）
// @Summary      批量更新服务项目状态
// @Description  批量启用/禁用多个服务项目，最多100条
// @Tags         Admin - ServiceItem
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchUpdateItemStatusRequest  true  "批量更新状态请求"
// @Success      200  {object}  item.BatchOperationResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      401  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /admin/service-items/batch/status [post]
func batchUpdateItemStatusHandler(c *gin.Context, svc *item.ServiceItemService) {
	var req BatchUpdateItemStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if len(req.ItemIDs) == 0 {
		respondAPIError(c, apierr.BadRequest("itemIds is required"))
		return
	}
	if len(req.ItemIDs) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 items per batch"))
		return
	}

	result, err := svc.BatchUpdateItemStatus(c.Request.Context(), item.BatchUpdateItemStatusRequest{
		ItemIDs:  req.ItemIDs,
		IsActive: req.IsActive,
	})
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch update status failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}

// batchDeleteItemsHandler 批量删除服务项目
// @Summary      批量删除服务项目
// @Description  批量删除多个服务项目，最多100条。检查是否有订单使用这些服务项目
// @Tags         Admin - ServiceItem
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchDeleteItemsRequest  true  "批量删除请求"
// @Success      200  {object}  item.BatchOperationResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      401  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /admin/service-items/batch/delete [post]
func batchDeleteItemsHandler(c *gin.Context, svc *item.ServiceItemService) {
	var req BatchDeleteItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if len(req.ItemIDs) == 0 {
		respondAPIError(c, apierr.BadRequest("itemIds is required"))
		return
	}
	if len(req.ItemIDs) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 items per batch"))
		return
	}

	result, err := svc.BatchDeleteItems(c.Request.Context(), item.BatchDeleteItemsRequest{
		ItemIDs: req.ItemIDs,
	})
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch delete failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}

// batchUpdateItemCommissionHandler 批量更新服务项目佣金比例
// @Summary      批量更新服务项目佣金比例
// @Description  批量更新多个服务项目的佣金比例，最多100条
// @Tags         Admin - ServiceItem
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchUpdateItemCommissionRequest  true  "批量更新佣金比例请求"
// @Success      200  {object}  item.BatchOperationResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      401  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /admin/service-items/batch/commission [post]
func batchUpdateItemCommissionHandler(c *gin.Context, svc *item.ServiceItemService) {
	var req BatchUpdateItemCommissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if len(req.ItemIDs) == 0 {
		respondAPIError(c, apierr.BadRequest("itemIds is required"))
		return
	}
	if len(req.ItemIDs) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 items per batch"))
		return
	}
	if req.CommissionRate < 0 || req.CommissionRate > 1 {
		respondAPIError(c, apierr.BadRequest("commissionRate must be between 0 and 1"))
		return
	}

	result, err := svc.BatchUpdateItemCommission(c.Request.Context(), item.BatchUpdateItemCommissionRequest{
		ItemIDs:        req.ItemIDs,
		CommissionRate: req.CommissionRate,
	})
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch update commission failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}
