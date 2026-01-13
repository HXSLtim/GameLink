// Package public provides public API handlers for service items without authentication.
package public

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
)

// ServiceItemHandler 公共服务项目处理器
type ServiceItemHandler struct {
	items repository.ServiceItemRepository
}

// NewServiceItemHandler 创建公共服务项目处理器
func NewServiceItemHandler(items repository.ServiceItemRepository) *ServiceItemHandler {
	return &ServiceItemHandler{
		items: items,
	}
}

// PublicServiceItemInfo 公开的服务项目信息
type PublicServiceItemInfo struct {
	ID             uint64  `json:"id"`
	Name           string  `json:"name"`
	GameID         *uint64 `json:"gameId"`
	Category       string  `json:"category"`
	SubCategory    string  `json:"subCategory"`
	BasePriceCents int64   `json:"basePriceCents"`
	ServiceHours   int     `json:"serviceHours"`
	Description    string  `json:"description"`
	Icon           string  `json:"icon"`
}

// ServiceItemListResponse 服务项目列表响应
type ServiceItemListResponse struct {
	Items    []PublicServiceItemInfo `json:"items"`
	Total    int64                   `json:"total"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"pageSize"`
}

// ListServiceItems 获取服务项目列表（公开）
// @Summary 获取服务项目列表
// @Description 获取平台服务项目列表，无需登录
// @Tags 公共-服务项目
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(50)
// @Param gameId query int false "游戏ID筛选"
// @Success 200 {object} ServiceItemListResponse
// @Router /public/service-items [get]
func (h *ServiceItemHandler) ListServiceItems(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
	gameIDStr := c.Query("gameId")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	opts := repository.ServiceItemListOptions{
		Page:     page,
		PageSize: pageSize,
	}

	// 只查询激活的服务项目
	isActive := true
	opts.IsActive = &isActive

	if gameIDStr != "" {
		gameID, err := strconv.ParseUint(gameIDStr, 10, 64)
		if err == nil {
			opts.GameID = &gameID
		}
	}

	items, total, err := h.items.List(c.Request.Context(), opts)
	if err != nil {
		resp.Error(c, apierr.InternalError("获取服务项目列表失败"))
		return
	}

	// 转换为公开信息
	result := make([]PublicServiceItemInfo, 0, len(items))
	for _, item := range items {
		info := PublicServiceItemInfo{
			ID:             item.ID,
			Name:           item.Name,
			GameID:         item.GameID,
			Category:       item.Category,
			SubCategory:    string(item.SubCategory),
			BasePriceCents: item.BasePriceCents,
			ServiceHours:   item.ServiceHours,
			Description:    item.Description,
			Icon:           item.IconURL,
		}
		result = append(result, info)
	}

	resp.OK(c, ServiceItemListResponse{
		Items:    result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// GetServiceItem 获取服务项目详情（公开）
// @Summary 获取服务项目详情
// @Description 获取指定服务项目的详细信息，无需登录
// @Tags 公共-服务项目
// @Accept json
// @Produce json
// @Param id path int true "服务项目ID"
// @Success 200 {object} PublicServiceItemInfo
// @Failure 404 {object}  apierr.APIError
// @Router /public/service-items/{id} [get]
func (h *ServiceItemHandler) GetServiceItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		resp.Error(c, apierr.BadRequest("无效的服务项目ID"))
		return
	}

	item, err := h.items.Get(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			resp.Error(c, apierr.NotFound("服务项目不存在"))
			return
		}
		resp.Error(c, apierr.InternalError("获取服务项目信息失败"))
		return
	}

	// 只返回激活的服务项目
	if !item.IsActive {
		resp.Error(c, apierr.NotFound("服务项目不存在"))
		return
	}

	info := PublicServiceItemInfo{
		ID:             item.ID,
		Name:           item.Name,
		GameID:         item.GameID,
		Category:       item.Category,
		SubCategory:    string(item.SubCategory),
		BasePriceCents: item.BasePriceCents,
		ServiceHours:   item.ServiceHours,
		Description:    item.Description,
		Icon:           item.IconURL,
	}

	resp.OK(c, info)
}

// RegisterRoutes 注册公共服务项目路由
func (h *ServiceItemHandler) RegisterRoutes(rg *gin.RouterGroup) {
	items := rg.Group("/service-items")
	{
		items.GET("", h.ListServiceItems)
		items.GET("/:id", h.GetServiceItem)
	}
}
