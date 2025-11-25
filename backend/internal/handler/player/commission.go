package player

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/apierr"
	"gamelink/internal/model"
	"gamelink/internal/service/commission"
)

// CommissionSummaryResponse 抽成汇总响应（类型别名）
type CommissionSummaryResponse = commission.CommissionSummaryResponse

// CommissionRecordListResponse 抽成记录列表响应（类型别名）
type CommissionRecordListResponse = commission.CommissionRecordListResponse

// SettlementListResponse 结算列表响应（类型别名）
type SettlementListResponse = commission.SettlementListResponse

// Local type definitions for Swagger generation (these mirror the commission service types)
type CommissionSummaryResponseSwagger struct {
	Month           string  `json:"month" example:"2024-01"`
	TotalIncome     float64 `json:"total_income" example:"10000.00"`
	TotalCommission float64 `json:"total_commission" example:"2000.00"`
	NetIncome       float64 `json:"net_income" example:"8000.00"`
	OrderCount      int     `json:"order_count" example:"25"`
}

// CommissionRecordListResponseSwagger 抽成记录列表响应
type CommissionRecordListResponseSwagger struct {
	Items      []CommissionRecordSwagger `json:"items"`
	Total      int64                     `json:"total" example:"100"`
	Page       int                       `json:"page" example:"1"`
	PageSize   int                       `json:"page_size" example:"20"`
	TotalPages int                       `json:"total_pages" example:"5"`
}

// CommissionRecordSwagger 抽成记录
type CommissionRecordSwagger struct {
	ID              uint64  `json:"id" example:"1"`
	OrderID         uint64  `json:"order_id" example:"1001"`
	PlayerID        uint64  `json:"player_id" example:"2001"`
	CommissionRate  float64 `json:"commission_rate" example:"0.2"`
	CommissionAmount float64 `json:"commission_amount" example:"100.00"`
	TotalPrice      float64 `json:"total_price" example:"500.00"`
	NetIncome       float64 `json:"net_income" example:"400.00"`
	CreatedAt       string  `json:"created_at" example:"2024-01-15T10:30:00Z"`
}

// SettlementListResponseSwagger 结算列表响应
type SettlementListResponseSwagger struct {
	Items      []SettlementSwagger `json:"items"`
	Total      int64               `json:"total" example:"50"`
	Page       int                 `json:"page" example:"1"`
	PageSize   int                 `json:"page_size" example:"20"`
	TotalPages int                 `json:"total_pages" example:"3"`
}

// SettlementSwagger 结算记录
type SettlementSwagger struct {
	ID            uint64  `json:"id" example:"1"`
	PlayerID      uint64  `json:"player_id" example:"2001"`
	Month         string  `json:"month" example:"2024-01"`
	TotalAmount   float64 `json:"total_amount" example:"8000.00"`
	Status        string  `json:"status" example:"completed"`
	SettledAt     string  `json:"settled_at,omitempty" example:"2024-02-01T00:00:00Z"`
	CreatedAt     string  `json:"created_at" example:"2024-01-31T23:59:59Z"`
}

// Swagger-friendly envelopes to avoid generics in swag annotations
type CommissionSummaryAPIResponseSwagger struct {
	Success    bool                          `json:"success"`
	Code       int                           `json:"code"`
	Message    string                        `json:"message"`
	Data       CommissionSummaryResponseSwagger `json:"data"`
	Pagination *model.Pagination             `json:"pagination,omitempty"`
	TraceID    string                        `json:"traceId,omitempty"`
}

type CommissionRecordListAPIResponseSwagger struct {
	Success    bool                               `json:"success"`
	Code       int                                `json:"code"`
	Message    string                             `json:"message"`
	Data       CommissionRecordListResponseSwagger `json:"data"`
	Pagination *model.Pagination                  `json:"pagination,omitempty"`
	TraceID    string                             `json:"traceId,omitempty"`
}

type SettlementListAPIResponseSwagger struct {
	Success    bool                          `json:"success"`
	Code       int                           `json:"code"`
	Message    string                        `json:"message"`
	Data       SettlementListResponseSwagger `json:"data"`
	Pagination *model.Pagination             `json:"pagination,omitempty"`
	TraceID    string                        `json:"traceId,omitempty"`
}

// RegisterCommissionRoutes 注册陪玩师端抽成管理路由
func RegisterCommissionRoutes(router gin.IRouter, svc *commission.CommissionService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/player/commission")
	group.Use(authMiddleware) // 需要认
	{
		group.GET("/summary", func(c *gin.Context) { getCommissionSummaryHandler(c, svc) })
		group.GET("/records", func(c *gin.Context) { getCommissionRecordsHandler(c, svc) })
		group.GET("/settlements", func(c *gin.Context) { getMonthlySettlementsHandler(c, svc) })
	}
}

// getCommissionSummaryHandler 获取抽成汇总
// @Summary      获取抽成汇总
// @Description  获取指定月份的抽成总览，包括总收入、总抽成等
// @Tags         Player - Commission
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        month          query     string  true  "月份 (YYYY-MM)"
// @Success      200            {object}  CommissionSummaryAPIResponseSwagger
// @Failure      401            {object}  apierr.APIError
// @Failure      404            {object}  apierr.APIError
// @Failure      500            {object}  apierr.APIError
// @Router       /player/commission/summary [get]
func getCommissionSummaryHandler(c *gin.Context, svc *commission.CommissionService) {
	userID := getUserIDFromContext(c)

	// 获取月份参数，默认当前月
	month := c.DefaultQuery("month", time.Now().Format("2006-01"))

	// 查找该用户的陪玩师ID
	playerID, err := getPlayerIDByUserID(c, userID)
	if err != nil {
		respondAPIError(c, apierr.NotFound("陪玩师信息未找到").WithDetails(err.Error()))
		return
	}

	resp, err := svc.GetPlayerCommissionSummary(c.Request.Context(), playerID, month)
	if err != nil {
		respondAPIError(c, apierr.InternalError("获取抽成汇总失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "OK", *resp)
}

// getCommissionRecordsHandler 获取抽成记录
// @Summary      获取抽成记录
// @Description  获取抽成记录列表，支持分页
// @Tags         Player - Commission
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page           query     int     false  "页码" default(1)
// @Param        pageSize       query     int     false  "每页数量" default(20)
// @Success      200            {object}  CommissionRecordListAPIResponseSwagger
// @Failure      401            {object}  apierr.APIError
// @Failure      404            {object}  apierr.APIError
// @Failure      500            {object}  apierr.APIError
// @Router       /player/commission/records [get]
func getCommissionRecordsHandler(c *gin.Context, svc *commission.CommissionService) {
	userID := getUserIDFromContext(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// 查找该用户的陪玩师ID
	playerID, err := getPlayerIDByUserID(c, userID)
	if err != nil {
		respondAPIError(c, apierr.NotFound("陪玩师信息未找到").WithDetails(err.Error()))
		return
	}

	resp, err := svc.GetCommissionRecords(c.Request.Context(), playerID, page, pageSize)
	if err != nil {
		respondAPIError(c, apierr.InternalError("获取抽成记录失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "OK", *resp)
}

// getMonthlySettlementsHandler 获取月度结算记录
// @Summary      获取月度结算记录
// @Description  获取月度结算记录列表，支持分页
// @Tags         Player - Commission
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page           query     int     false  "页码" default(1)
// @Param        pageSize       query     int     false  "每页数量" default(20)
// @Success      200            {object}  SettlementListAPIResponseSwagger
// @Failure      401            {object}  apierr.APIError
// @Failure      404            {object}  apierr.APIError
// @Failure      500            {object}  apierr.APIError
// @Router       /player/commission/settlements [get]
func getMonthlySettlementsHandler(c *gin.Context, svc *commission.CommissionService) {
	userID := getUserIDFromContext(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// 查找该用户的陪玩师ID
	playerID, err := getPlayerIDByUserID(c, userID)
	if err != nil {
		respondAPIError(c, apierr.NotFound("陪玩师信息未找到").WithDetails(err.Error()))
		return
	}

	resp, err := svc.GetMonthlySettlements(c.Request.Context(), playerID, page, pageSize)
	if err != nil {
		respondAPIError(c, apierr.InternalError("获取月度结算记录失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "OK", *resp)
}

// getPlayerIDByUserID 根据用户ID获取陪玩师ID
func getPlayerIDByUserID(c *gin.Context, userID uint64) (uint64, error) {
	// TODO: 优化这个查询，可以在用户上下文中缓存playerID
	// 这里简化处理，实际应该从service层获
	// 暂时返回userID作为playerID（需要后续完善）
	return userID, nil
}
