// Package admin provides admin handlers for the API.
package admin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	analyticsservice "gamelink/internal/service/analytics"
)

// AnalyticsHandler handles analytics-related requests.
type AnalyticsHandler struct {
	analyticsSvc *analyticsservice.AnalyticsService
}

// NewAnalyticsHandler creates a new analytics handler.
func NewAnalyticsHandler(analyticsSvc *analyticsservice.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsSvc: analyticsSvc}
}

// AnalyticsQueryParams represents common query parameters for analytics.
type AnalyticsQueryParams struct {
	StartDate   string `form:"start_date" binding:"required"`
	EndDate     string `form:"end_date" binding:"required"`
	Granularity string `form:"granularity"`
}

// parseQueryParams parses and validates analytics query parameters.
func parseQueryParams(c *gin.Context) (*analyticsservice.DateRange, analyticsservice.Granularity, error) {
	var params AnalyticsQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		return nil, "", err
	}

	startDate, err := time.Parse("2006-01-02", params.StartDate)
	if err != nil {
		return nil, "", err
	}

	endDate, err := time.Parse("2006-01-02", params.EndDate)
	if err != nil {
		return nil, "", err
	}

	// Add end of day to end date
	endDate = endDate.Add(24*time.Hour - time.Second)

	granularity := analyticsservice.GranularityDay
	switch params.Granularity {
	case "week":
		granularity = analyticsservice.GranularityWeek
	case "month":
		granularity = analyticsservice.GranularityMonth
	}

	return &analyticsservice.DateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}, granularity, nil
}

// GetActiveUsers returns active users statistics.
// @Summary      获取活跃用户数据
// @Description  获取 DAU、WAU、MAU 及趋势数据
// @Tags         Admin/Analytics
// @Security     BearerAuth
// @Param        start_date   query  string  true   "开始日期 (YYYY-MM-DD)"
// @Param        end_date     query  string  true   "结束日期 (YYYY-MM-DD)"
// @Param        granularity  query  string  false  "时间粒度" Enums(day,week,month)
// @Produce      json
// @Success      200  {object}  model.APIResponse[analyticsservice.ActiveUsersData]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/analytics/active-users [get]
func (h *AnalyticsHandler) GetActiveUsers(c *gin.Context) {
	dateRange, granularity, err := parseQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	data, err := h.analyticsSvc.GetActiveUsers(c.Request.Context(), *dateRange, granularity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "获取活跃用户数据失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse[*analyticsservice.ActiveUsersData]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
	})
}

// GetRetention returns user retention statistics.
// @Summary      获取用户留存数据
// @Description  获取 D1、D7、D30 留存率及留存矩阵
// @Tags         Admin/Analytics
// @Security     BearerAuth
// @Param        start_date   query  string  true   "开始日期 (YYYY-MM-DD)"
// @Param        end_date     query  string  true   "结束日期 (YYYY-MM-DD)"
// @Param        granularity  query  string  false  "时间粒度" Enums(day,week,month)
// @Produce      json
// @Success      200  {object}  model.APIResponse[analyticsservice.RetentionData]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/analytics/retention [get]
func (h *AnalyticsHandler) GetRetention(c *gin.Context) {
	dateRange, granularity, err := parseQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	data, err := h.analyticsSvc.GetRetention(c.Request.Context(), *dateRange, granularity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "获取留存数据失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse[*analyticsservice.RetentionData]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
	})
}

// GetPaymentAnalytics returns payment analytics data.
// @Summary      获取付费分析数据
// @Description  获取付费率、ARPU、ARPPU 及收入趋势
// @Tags         Admin/Analytics
// @Security     BearerAuth
// @Param        start_date   query  string  true   "开始日期 (YYYY-MM-DD)"
// @Param        end_date     query  string  true   "结束日期 (YYYY-MM-DD)"
// @Param        granularity  query  string  false  "时间粒度" Enums(day,week,month)
// @Produce      json
// @Success      200  {object}  model.APIResponse[analyticsservice.PaymentData]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/analytics/payment [get]
func (h *AnalyticsHandler) GetPaymentAnalytics(c *gin.Context) {
	dateRange, granularity, err := parseQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	data, err := h.analyticsSvc.GetPaymentAnalytics(c.Request.Context(), *dateRange, granularity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "获取付费分析数据失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse[*analyticsservice.PaymentData]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
	})
}

// GetConversionFunnel returns conversion funnel data.
// @Summary      获取转化漏斗数据
// @Description  获取用户转化漏斗各阶段数据
// @Tags         Admin/Analytics
// @Security     BearerAuth
// @Param        start_date   query  string  true   "开始日期 (YYYY-MM-DD)"
// @Param        end_date     query  string  true   "结束日期 (YYYY-MM-DD)"
// @Param        granularity  query  string  false  "时间粒度" Enums(day,week,month)
// @Produce      json
// @Success      200  {object}  model.APIResponse[analyticsservice.ConversionFunnel]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/analytics/conversion [get]
func (h *AnalyticsHandler) GetConversionFunnel(c *gin.Context) {
	dateRange, granularity, err := parseQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	data, err := h.analyticsSvc.GetConversionFunnel(c.Request.Context(), *dateRange, granularity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "获取转化漏斗数据失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse[*analyticsservice.ConversionFunnel]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
	})
}
