package admin

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	mw "gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	svc "gamelink/internal/service/withdraw"
)

// parseUint64 parses a string to uint64
func parseUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

// WithdrawRoutingHandler 处理提现分流管理接口
// Requirements: 14.1-14.5
type WithdrawRoutingHandler struct {
	statsSvc *svc.WithdrawRoutingStatsService
}

// NewWithdrawRoutingHandler 创建Handler
func NewWithdrawRoutingHandler(statsSvc *svc.WithdrawRoutingStatsService) *WithdrawRoutingHandler {
	return &WithdrawRoutingHandler{statsSvc: statsSvc}
}

// ListWithdrawalsByCompany
// @Summary      按结算公司查询提现列表
// @Tags         Admin/WithdrawRouting
// @Security     BearerAuth
// @Param        settlementCompanyId  query  int     false  "结算公司ID"
// @Param        status               query  string  false  "状态过滤" Enums(pending, approved, rejected, completed, failed)
// @Param        dateFrom             query  string  false  "开始日期 (YYYY-MM-DD)"
// @Param        dateTo               query  string  false  "结束日期 (YYYY-MM-DD)"
// @Param        page                 query  int     false  "页码"
// @Param        pageSize             query  int     false  "每页数量"
// @Produce      json
// @Success      200  {object}  model.ListWithdrawsByCompanyResponse
// @Router       /admin/withdrawals/by-company [get]
func (h *WithdrawRoutingHandler) ListWithdrawalsByCompany(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	req := &model.ListWithdrawsByCompanyRequest{
		Page:     page,
		PageSize: pageSize,
	}

	// Parse settlement company ID
	if companyIDStr := strings.TrimSpace(c.Query("settlementCompanyId")); companyIDStr != "" {
		companyID, err := parseUint64(companyIDStr)
		if err != nil {
			respondBadRequest(c, "invalid settlementCompanyId")
			return
		}
		req.SettlementCompanyID = &companyID
	}

	// Parse status
	if statusStr := strings.TrimSpace(c.Query("status")); statusStr != "" {
		status := model.WithdrawStatus(statusStr)
		req.Status = &status
	}

	// Parse date range
	if dateFromStr := strings.TrimSpace(c.Query("dateFrom")); dateFromStr != "" {
		dateFrom, err := time.Parse("2006-01-02", dateFromStr)
		if err == nil {
			req.DateFrom = &dateFrom
		}
	}
	if dateToStr := strings.TrimSpace(c.Query("dateTo")); dateToStr != "" {
		dateTo, err := time.Parse("2006-01-02", dateToStr)
		if err == nil {
			// Add one day to include the end date
			dateTo = dateTo.AddDate(0, 0, 1)
			req.DateTo = &dateTo
		}
	}

	resp, err := h.statsSvc.ListWithdrawalsByCompany(c.Request.Context(), req)
	if err != nil {
		respondError(c, err)
		return
	}

	totalPages := int((resp.Total + int64(resp.PageSize) - 1) / int64(resp.PageSize))
	pagination := &model.Pagination{
		Page:       resp.Page,
		PageSize:   resp.PageSize,
		Total:      int(resp.Total),
		TotalPages: totalPages,
		HasNext:    resp.Page < totalPages,
		HasPrev:    resp.Page > 1,
	}
	respondList(c, resp.Withdraws, pagination)
}

// GetWithdrawRoutingStats
// @Summary      获取提现分流统计
// @Tags         Admin/WithdrawRouting
// @Security     BearerAuth
// @Param        dateFrom  query  string  false  "开始日期 (YYYY-MM-DD)"
// @Param        dateTo    query  string  false  "结束日期 (YYYY-MM-DD)"
// @Produce      json
// @Success      200  {object}  model.WithdrawRoutingStatsResponse
// @Router       /admin/withdrawals/routing-stats [get]
func (h *WithdrawRoutingHandler) GetWithdrawRoutingStats(c *gin.Context) {
	req := &model.WithdrawRoutingStatsRequest{}

	// Parse date range
	if dateFromStr := strings.TrimSpace(c.Query("dateFrom")); dateFromStr != "" {
		dateFrom, err := time.Parse("2006-01-02", dateFromStr)
		if err == nil {
			req.DateFrom = &dateFrom
		}
	}
	if dateToStr := strings.TrimSpace(c.Query("dateTo")); dateToStr != "" {
		dateTo, err := time.Parse("2006-01-02", dateToStr)
		if err == nil {
			// Add one day to include the end date
			dateTo = dateTo.AddDate(0, 0, 1)
			req.DateTo = &dateTo
		}
	}

	stats, err := h.statsSvc.GetRoutingStats(c.Request.Context(), req)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, stats)
}

// GenerateWithdrawRoutingReport
// @Summary      生成提现分流报表
// @Tags         Admin/WithdrawRouting
// @Security     BearerAuth
// @Param        reportType  query  string  true   "报表类型" Enums(monthly, quarterly, yearly)
// @Param        year        query  int     true   "年份"
// @Param        month       query  int     false  "月份 (月报必填)"
// @Param        quarter     query  int     false  "季度 (季报必填)"
// @Produce      json
// @Success      200  {object}  model.WithdrawRoutingReport
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/withdrawals/routing-report [get]
func (h *WithdrawRoutingHandler) GenerateWithdrawRoutingReport(c *gin.Context) {
	var req model.WithdrawRoutingReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	report, err := h.statsSvc.GenerateRoutingReport(c.Request.Context(), &req)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, report)
}

// GetCompanyWithdrawalStats
// @Summary      获取单个结算公司的提现统计
// @Tags         Admin/WithdrawRouting
// @Security     BearerAuth
// @Param        id        path   int     true   "结算公司ID"
// @Param        dateFrom  query  string  false  "开始日期 (YYYY-MM-DD)"
// @Param        dateTo    query  string  false  "结束日期 (YYYY-MM-DD)"
// @Produce      json
// @Success      200  {object}  model.WithdrawRoutingStats
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/settlement-companies/{id}/withdrawal-stats [get]
func (h *WithdrawRoutingHandler) GetCompanyWithdrawalStats(c *gin.Context) {
	companyID, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var dateFrom, dateTo *time.Time

	// Parse date range
	if dateFromStr := strings.TrimSpace(c.Query("dateFrom")); dateFromStr != "" {
		df, err := time.Parse("2006-01-02", dateFromStr)
		if err == nil {
			dateFrom = &df
		}
	}
	if dateToStr := strings.TrimSpace(c.Query("dateTo")); dateToStr != "" {
		dt, err := time.Parse("2006-01-02", dateToStr)
		if err == nil {
			// Add one day to include the end date
			dt = dt.AddDate(0, 0, 1)
			dateTo = &dt
		}
	}

	stats, err := h.statsSvc.GetCompanyWithdrawalStats(c.Request.Context(), companyID, dateFrom, dateTo)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, stats)
}

// RegisterWithdrawRoutingRoutes 注册提现分流管理路由
// Requirements: 14.1-14.5
func RegisterWithdrawRoutingRoutes(router gin.IRouter, statsSvc *svc.WithdrawRoutingStatsService, pm *mw.PermissionMiddleware) {
	handler := NewWithdrawRoutingHandler(statsSvc)

	// 提现分流统计
	// @Summary      按结算公司查询提现列表
	// @Tags         Admin/WithdrawRouting
	// @Security     BearerAuth
	// @Router       /admin/withdrawals/by-company [get]
	router.GET("/withdrawals/by-company", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/withdrawals/by-company"), handler.ListWithdrawalsByCompany)

	// @Summary      获取提现分流统计
	// @Tags         Admin/WithdrawRouting
	// @Security     BearerAuth
	// @Router       /admin/withdrawals/routing-stats [get]
	router.GET("/withdrawals/routing-stats", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/withdrawals/routing-stats"), handler.GetWithdrawRoutingStats)

	// @Summary      生成提现分流报表
	// @Tags         Admin/WithdrawRouting
	// @Security     BearerAuth
	// @Router       /admin/withdrawals/routing-report [get]
	router.GET("/withdrawals/routing-report", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/withdrawals/routing-report"), handler.GenerateWithdrawRoutingReport)

	// @Summary      获取单个结算公司的提现统计
	// @Tags         Admin/WithdrawRouting
	// @Security     BearerAuth
	// @Router       /admin/settlement-companies/{id}/withdrawal-stats [get]
	router.GET("/settlement-companies/:id/withdrawal-stats", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/settlement-companies/:id/withdrawal-stats"), handler.GetCompanyWithdrawalStats)
}
