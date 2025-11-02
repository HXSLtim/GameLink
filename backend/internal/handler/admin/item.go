package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

    "gamelink/internal/model"
    "gamelink/internal/service/item"
)

// RegisterServiceItemRoutes 注册管理端服务项目管理路�?
func RegisterServiceItemRoutes(router gin.IRouter, svc *item.ServiceItemService) {
	group := router.Group("/admin/service-items")
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
// @Param        request        body      serviceitem.CreateServiceItemRequest  true  "服务项目信息"
// @Success      200            {object}  model.APIResponse[model.ServiceItem]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/service-items [post]
func createServiceItemHandler(c *gin.Context, svc *item.ServiceItemService) {
    var req item.CreateServiceItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	item, err := svc.CreateServiceItem(c.Request.Context(), req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

    respondJSON(c, http.StatusOK, model.APIResponse[model.ServiceItem]{
		Success: true,
		Code:    http.StatusOK,
		Message: "Service item created successfully",
		Data:    *item,
	})
}

// listServiceItemsHandler 获取服务项目列表
// @Summary      获取服务项目列表
// @Description  管理员查看所有服务项�?
// @Tags         Admin - ServiceItem
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        category       query     string  false  "分类"
// @Param        subCategory    query     string  false  "子分�?solo/team/gift)"
// @Param        gameId         query     int     false  "游戏ID"
// @Param        isActive       query     bool    false  "是否激�?
// @Param        page           query     int     false  "页码"
// @Param        pageSize       query     int     false  "每页数量"
// @Success      200            {object}  model.APIResponse[serviceitem.ServiceItemListResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/service-items [get]
func listServiceItemsHandler(c *gin.Context, svc *item.ServiceItemService) {
    var req item.ListServiceItemsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := svc.ListServiceItems(c.Request.Context(), req)
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

// getServiceItemHandler 获取服务项目详情
// @Summary      获取服务项目详情
// @Description  管理员查看服务项目详�?
// @Tags         Admin - ServiceItem
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      int     true  "服务项目ID"
// @Success      200            {object}  model.APIResponse[serviceitem.ServiceItemDTO]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/service-items/{id} [get]
func getServiceItemHandler(c *gin.Context, svc *item.ServiceItemService) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid service item ID")
		return
	}

	resp, err := svc.GetServiceItem(c.Request.Context(), id)
    if err != nil {
        if err == item.ErrNotFound {
			respondError(c, http.StatusNotFound, "Service item not found")
			return
		}
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

    respondJSON(c, http.StatusOK, model.APIResponse[item.ServiceItemDTO]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    *resp,
	})
}

// updateServiceItemHandler 更新服务项目
// @Summary      更新服务项目
// @Description  管理员更新服务项�?
// @Tags         Admin - ServiceItem
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                                 true  "Bearer {token}"
// @Param        id             path      int                                    true  "服务项目ID"
// @Param        request        body      serviceitem.UpdateServiceItemRequest  true  "更新信息"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/service-items/{id} [put]
func updateServiceItemHandler(c *gin.Context, svc *item.ServiceItemService) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid service item ID")
		return
	}

    var req item.UpdateServiceItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	err = svc.UpdateServiceItem(c.Request.Context(), id, req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "Service item updated successfully",
	})
}

// deleteServiceItemHandler 删除服务项目
// @Summary      删除服务项目
// @Description  管理员删除服务项�?
// @Tags         Admin - ServiceItem
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      int     true  "服务项目ID"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/service-items/{id} [delete]
func deleteServiceItemHandler(c *gin.Context, svc *item.ServiceItemService) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid service item ID")
		return
	}

	err = svc.DeleteServiceItem(c.Request.Context(), id)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "Service item deleted successfully",
	})
}

// batchUpdateStatusHandler 批量更新状�?
// @Summary      批量更新状�?
// @Description  管理员批量启�?禁用服务项目
// @Tags         Admin - ServiceItem
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                                  true  "Bearer {token}"
// @Param        request        body      serviceitem.BatchUpdateStatusRequest  true  "批量更新请求"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/service-items/batch-update-status [post]
func batchUpdateStatusHandler(c *gin.Context, svc *item.ServiceItemService) {
    var req item.BatchUpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	err := svc.BatchUpdateStatus(c.Request.Context(), req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "Status updated successfully",
	})
}

// batchUpdatePriceHandler 批量更新价格
// @Summary      批量更新价格
// @Description  管理员批量调整服务项目价�?
// @Tags         Admin - ServiceItem
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                                 true  "Bearer {token}"
// @Param        request        body      serviceitem.BatchUpdatePriceRequest  true  "批量更新请求"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/service-items/batch-update-price [post]
func batchUpdatePriceHandler(c *gin.Context, svc *item.ServiceItemService) {
    var req item.BatchUpdatePriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	err := svc.BatchUpdatePrice(c.Request.Context(), req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "Price updated successfully",
	})
}

