// Package admin provides admin handlers for the API.
package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	monitorservice "gamelink/internal/service/monitor"
	"gamelink/internal/ws"
)

// MonitorHandler 监控相关处理器
type MonitorHandler struct {
	realtimeSvc *monitorservice.RealtimeService
	alertRepo   model.AlertRepository
}

// NewMonitorHandler 创建监控处理器
func NewMonitorHandler(realtimeSvc *monitorservice.RealtimeService, alertRepo model.AlertRepository) *MonitorHandler {
	return &MonitorHandler{
		realtimeSvc: realtimeSvc,
		alertRepo:   alertRepo,
	}
}

// GetSystemStatus 获取系统状态快照
// @Summary      获取系统状态快照
// @Description  获取当前系统CPU、内存、协程、数据库连接等状态
// @Tags         Admin/Monitor
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  model.APIResponse[ws.SystemStatus]
// @Router       /admin/monitor/system-status [get]
func (h *MonitorHandler) GetSystemStatus(c *gin.Context) {
	status := h.realtimeSvc.GetSystemStatus()
	c.JSON(http.StatusOK, model.APIResponse[*ws.SystemStatus]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    status,
	})
}

// GetOnlineUsers 获取在线用户统计
// @Summary      获取在线用户统计
// @Description  获取当前在线用户数、峰值、分角色统计
// @Tags         Admin/Monitor
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  model.APIResponse[ws.OnlineUsers]
// @Router       /admin/monitor/online-users [get]
func (h *MonitorHandler) GetOnlineUsers(c *gin.Context) {
	users := h.realtimeSvc.GetOnlineUsers()
	c.JSON(http.StatusOK, model.APIResponse[*ws.OnlineUsers]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    users,
	})
}

// GetOrderQueue 获取订单队列状态
// @Summary      获取订单队列状态
// @Description  获取待处理、处理中订单数量及处理速度
// @Tags         Admin/Monitor
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  model.APIResponse[ws.OrderQueue]
// @Router       /admin/monitor/order-queue [get]
func (h *MonitorHandler) GetOrderQueue(c *gin.Context) {
	queue := h.realtimeSvc.GetOrderQueue(c.Request.Context())
	c.JSON(http.StatusOK, model.APIResponse[*ws.OrderQueue]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    queue,
	})
}

// AlertQueryParams 告警查询参数
type AlertQueryParams struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Level    string `form:"level"`
	Type     string `form:"type"`
	IsRead   *bool  `form:"is_read"`
	DateFrom string `form:"date_from"`
	DateTo   string `form:"date_to"`
}

// GetAlerts 获取告警列表
// @Summary      获取告警列表
// @Description  获取系统告警列表，支持分页和筛选
// @Tags         Admin/Monitor
// @Security     BearerAuth
// @Param        page       query  int     false  "页码"
// @Param        page_size  query  int     false  "每页数量"
// @Param        level      query  string  false  "告警级别" Enums(high,medium,low)
// @Param        type       query  string  false  "告警类型" Enums(system,business,security)
// @Param        is_read    query  bool    false  "是否已读"
// @Param        date_from  query  string  false  "开始日期"
// @Param        date_to    query  string  false  "结束日期"
// @Produce      json
// @Success      200  {object}  model.APIResponse[[]model.Alert]
// @Router       /admin/monitor/alerts [get]
func (h *MonitorHandler) GetAlerts(c *gin.Context) {
	var params AlertQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 设置默认值
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}

	// 构建查询选项
	opts := model.AlertQueryOptions{
		Page:     params.Page,
		PageSize: params.PageSize,
		Level:    params.Level,
		Type:     params.Type,
		IsRead:   params.IsRead,
	}

	if params.DateFrom != "" {
		if t, err := time.Parse("2006-01-02", params.DateFrom); err == nil {
			opts.DateFrom = &t
		}
	}
	if params.DateTo != "" {
		if t, err := time.Parse("2006-01-02", params.DateTo); err == nil {
			opts.DateTo = &t
		}
	}

	alerts, total, err := h.alertRepo.List(c.Request.Context(), opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "获取告警列表失败: " + err.Error(),
		})
		return
	}

	totalPages := (int(total) + params.PageSize - 1) / params.PageSize
	c.JSON(http.StatusOK, model.APIResponse[[]model.Alert]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    alerts,
		Pagination: &model.Pagination{
			Page:       params.Page,
			PageSize:   params.PageSize,
			Total:      int(total),
			TotalPages: totalPages,
			HasNext:    params.Page < totalPages,
			HasPrev:    params.Page > 1,
		},
	})
}

// MarkAlertRead 标记告警已读
// @Summary      标记告警已读
// @Description  将指定告警标记为已读
// @Tags         Admin/Monitor
// @Security     BearerAuth
// @Param        id   path  string  true  "告警ID"
// @Produce      json
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/monitor/alerts/{id}/read [put]
func (h *MonitorHandler) MarkAlertRead(c *gin.Context) {
	alertID := c.Param("id")
	if alertID == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "告警ID不能为空",
		})
		return
	}

	id, err := strconv.ParseUint(alertID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "无效的告警ID",
		})
		return
	}

	if err := h.alertRepo.MarkAsRead(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "标记已读失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
	})
}

// BatchMarkReadPayload 批量标记已读请求体
type BatchMarkReadPayload struct {
	IDs []string `json:"ids" binding:"required"`
}

// BatchMarkAlertsRead 批量标记告警已读
// @Summary      批量标记告警已读
// @Description  批量将多个告警标记为已读
// @Tags         Admin/Monitor
// @Security     BearerAuth
// @Accept       json
// @Param        request  body  BatchMarkReadPayload  true  "告警ID列表"
// @Produce      json
// @Success      200  {object}  model.SuccessResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/monitor/alerts/batch-read [put]
func (h *MonitorHandler) BatchMarkAlertsRead(c *gin.Context) {
	var payload BatchMarkReadPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 转换ID
	ids := make([]uint, 0, len(payload.IDs))
	for _, idStr := range payload.IDs {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, uint(id))
	}

	if len(ids) == 0 {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "没有有效的告警ID",
		})
		return
	}

	if err := h.alertRepo.BatchMarkAsRead(c.Request.Context(), ids); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "批量标记已读失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
	})
}
