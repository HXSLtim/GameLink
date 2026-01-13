// Package admin provides admin handlers for the API.
package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	kpiservice "gamelink/internal/service/kpi"
)

// KPIHandler handles KPI-related requests.
type KPIHandler struct {
	kpiSvc *kpiservice.KPIService
}

// NewKPIHandler creates a new KPI handler.
func NewKPIHandler(kpiSvc *kpiservice.KPIService) *KPIHandler {
	return &KPIHandler{kpiSvc: kpiSvc}
}

// KPIQueryParams represents KPI query parameters.
type KPIQueryParams struct {
	Period    string `form:"period"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Compare   string `form:"compare"`
}

// parseKPIQueryParams parses KPI query parameters.
func parseKPIQueryParams(c *gin.Context) kpiservice.QueryParams {
	var qp KPIQueryParams
	_ = c.ShouldBindQuery(&qp)

	params := kpiservice.QueryParams{}

	switch qp.Period {
	case "today":
		params.Period = kpiservice.PeriodToday
	case "week":
		params.Period = kpiservice.PeriodWeek
	case "month":
		params.Period = kpiservice.PeriodMonth
	case "custom":
		params.Period = kpiservice.PeriodCustom
	default:
		params.Period = kpiservice.PeriodMonth
	}

	if qp.StartDate != "" {
		if t, err := time.Parse("2006-01-02", qp.StartDate); err == nil {
			params.StartDate = &t
		}
	}
	if qp.EndDate != "" {
		if t, err := time.Parse("2006-01-02", qp.EndDate); err == nil {
			params.EndDate = &t
		}
	}

	switch qp.Compare {
	case "yoy":
		params.Compare = kpiservice.CompareYoY
	default:
		params.Compare = kpiservice.CompareMoM
	}

	return params
}

// GetOverview returns KPI overview data.
// @Summary      获取 KPI 概览
// @Description  获取 GMV、订单数、新用户等 KPI 指标概览
// @Tags         Admin/KPI
// @Security     BearerAuth
// @Param        period     query  string  false  "时间周期" Enums(today,week,month,custom)
// @Param        start_date query  string  false  "开始日期 (YYYY-MM-DD)"
// @Param        end_date   query  string  false  "结束日期 (YYYY-MM-DD)"
// @Param        compare    query  string  false  "对比类型" Enums(mom,yoy)
// @Produce      json
// @Success      200  {object}  kpiservice.KPIOverview
// @Router       /admin/kpi/overview [get]
func (h *KPIHandler) GetOverview(c *gin.Context) {
	params := parseKPIQueryParams(c)

	overview, err := h.kpiSvc.GetOverview(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "获取 KPI 概览失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse[*kpiservice.KPIOverview]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    overview,
	})
}

// GetTrend returns trend data for a specific KPI metric.
// @Summary      获取 KPI 趋势
// @Description  获取指定 KPI 指标的趋势数据
// @Tags         Admin/KPI
// @Security     BearerAuth
// @Param        metric     query  string  true   "指标名称" Enums(gmv,orders,new_users,dau)
// @Param        period     query  string  false  "时间周期" Enums(today,week,month,custom)
// @Param        start_date query  string  false  "开始日期 (YYYY-MM-DD)"
// @Param        end_date   query  string  false  "结束日期 (YYYY-MM-DD)"
// @Produce      json
// @Success      200  {array}   kpiservice.TrendPoint
// @Router       /admin/kpi/trend [get]
func (h *KPIHandler) GetTrend(c *gin.Context) {
	metric := c.Query("metric")
	if metric == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "指标名称不能为空",
		})
		return
	}

	params := parseKPIQueryParams(c)

	trend, err := h.kpiSvc.GetTrend(c.Request.Context(), metric, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "获取趋势数据失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse[[]kpiservice.TrendPoint]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    trend,
	})
}

// GetTargets returns KPI target configurations.
// @Summary      获取 KPI 目标列表
// @Description  获取 KPI 目标配置列表
// @Tags         Admin/KPI
// @Security     BearerAuth
// @Param        period_type  query  string  false  "周期类型" Enums(daily,weekly,monthly)
// @Param        metric_name  query  string  false  "指标名称"
// @Produce      json
// @Success      200  {array}   model.KPITarget
// @Router       /admin/kpi/targets [get]
func (h *KPIHandler) GetTargets(c *gin.Context) {
	periodType := c.Query("period_type")
	metricName := c.Query("metric_name")

	targets, err := h.kpiSvc.GetTargets(c.Request.Context(), periodType, metricName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "获取目标列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse[[]model.KPITarget]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    targets,
	})
}

// CreateTargetPayload represents the request body for creating a KPI target.
type CreateTargetPayload struct {
	PeriodType  string  `json:"periodType" binding:"required"`
	MetricName  string  `json:"metricName" binding:"required"`
	TargetValue float64 `json:"targetValue" binding:"required"`
	StartDate   string  `json:"startDate" binding:"required"`
	EndDate     string  `json:"endDate" binding:"required"`
}

// CreateTarget creates a new KPI target.
// @Summary      创建 KPI 目标
// @Description  创建新的 KPI 目标配置
// @Tags         Admin/KPI
// @Security     BearerAuth
// @Accept       json
// @Param        request  body  CreateTargetPayload  true  "目标配置"
// @Produce      json
// @Success      201  {object}  model.KPITarget
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/kpi/targets [post]
func (h *KPIHandler) CreateTarget(c *gin.Context) {
	var payload CreateTargetPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", payload.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "开始日期格式错误",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", payload.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "结束日期格式错误",
		})
		return
	}

	// Get user ID from context
	userIDValue, _ := c.Get("user_id")
	var userID uint64
	switch v := userIDValue.(type) {
	case uint64:
		userID = v
	case float64:
		userID = uint64(v)
	case int:
		userID = uint64(v)
	}

	target := &model.KPITarget{
		PeriodType:  payload.PeriodType,
		MetricName:  payload.MetricName,
		TargetValue: payload.TargetValue,
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedBy:   userID,
	}

	if err := h.kpiSvc.CreateTarget(c.Request.Context(), target); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "创建目标失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, model.APIResponse[*model.KPITarget]{
		Success: true,
		Code:    http.StatusCreated,
		Message: "OK",
		Data:    target,
	})
}

// UpdateTargetPayload represents the request body for updating a KPI target.
type UpdateTargetPayload struct {
	PeriodType  *string  `json:"periodType"`
	MetricName  *string  `json:"metricName"`
	TargetValue *float64 `json:"targetValue"`
	StartDate   *string  `json:"startDate"`
	EndDate     *string  `json:"endDate"`
}

// UpdateTarget updates an existing KPI target.
// @Summary      更新 KPI 目标
// @Description  更新 KPI 目标配置
// @Tags         Admin/KPI
// @Security     BearerAuth
// @Accept       json
// @Param        id       path  int                  true  "目标ID"
// @Param        request  body  UpdateTargetPayload  true  "更新内容"
// @Produce      json
// @Success      200  {object}  model.SuccessResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/kpi/targets/{id} [put]
func (h *KPIHandler) UpdateTarget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "无效的目标ID",
		})
		return
	}

	var payload UpdateTargetPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	updates := make(map[string]interface{})
	if payload.PeriodType != nil {
		updates["period_type"] = *payload.PeriodType
	}
	if payload.MetricName != nil {
		updates["metric_name"] = *payload.MetricName
	}
	if payload.TargetValue != nil {
		updates["target_value"] = *payload.TargetValue
	}
	if payload.StartDate != nil {
		if t, err := time.Parse("2006-01-02", *payload.StartDate); err == nil {
			updates["start_date"] = t
		}
	}
	if payload.EndDate != nil {
		if t, err := time.Parse("2006-01-02", *payload.EndDate); err == nil {
			updates["end_date"] = t
		}
	}

	if err := h.kpiSvc.UpdateTarget(c.Request.Context(), uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "更新目标失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
	})
}

// DeleteTarget deletes a KPI target.
// @Summary      删除 KPI 目标
// @Description  删除 KPI 目标配置
// @Tags         Admin/KPI
// @Security     BearerAuth
// @Param        id  path  int  true  "目标ID"
// @Produce      json
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/kpi/targets/{id} [delete]
func (h *KPIHandler) DeleteTarget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "无效的目标ID",
		})
		return
	}

	if err := h.kpiSvc.DeleteTarget(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "删除目标失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
	})
}
