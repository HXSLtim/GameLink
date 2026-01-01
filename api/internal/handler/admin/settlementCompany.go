package admin

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	mw "gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	svc "gamelink/internal/service/settlementcompany"
	"gamelink/pkg/apierr"
)

// SettlementCompanyHandler 处理结算公司管理接口
// Requirements: 11.1-11.5, 12.1-12.5
type SettlementCompanyHandler struct {
	svc *svc.SettlementCompanyService
}

// NewSettlementCompanyHandler 创建Handler
func NewSettlementCompanyHandler(svc *svc.SettlementCompanyService) *SettlementCompanyHandler {
	return &SettlementCompanyHandler{svc: svc}
}

// ListSettlementCompanies
// @Summary      列出结算公司
// @Tags         Admin/SettlementCompanies
// @Security     BearerAuth
// @Param        page       query  int     false  "页码"
// @Param        pageSize   query  int     false  "每页数量"
// @Param        status     query  string  false  "状态过滤" Enums(active, inactive)
// @Param        keyword    query  string  false  "关键词搜索"
// @Param        sortBy     query  string  false  "排序字段" Enums(name, created_at, player_count)
// @Param        sortOrder  query  string  false  "排序方向" Enums(asc, desc)
// @Produce      json
// @Success      200  {object}  model.APIResponse[model.ListSettlementCompaniesResponse]
// @Router       /admin/settlement-companies [get]
func (h *SettlementCompanyHandler) ListSettlementCompanies(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	req := &model.ListSettlementCompaniesRequest{
		Page:      page,
		PageSize:  pageSize,
		Status:    model.CompanyStatus(strings.TrimSpace(c.Query("status"))),
		Keyword:   strings.TrimSpace(c.Query("keyword")),
		SortBy:    strings.TrimSpace(c.Query("sortBy")),
		SortOrder: strings.TrimSpace(c.Query("sortOrder")),
	}

	resp, err := h.svc.ListCompanies(c.Request.Context(), req)
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
	respondList(c, resp.Companies, pagination)
}

// GetSettlementCompany
// @Summary      获取结算公司详情
// @Tags         Admin/SettlementCompanies
// @Security     BearerAuth
// @Param        id   path  int  true  "结算公司ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[model.SettlementCompany]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/settlement-companies/{id} [get]
func (h *SettlementCompanyHandler) GetSettlementCompany(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	company, err := h.svc.GetCompany(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, company)
}

// CreateSettlementCompany
// @Summary      创建结算公司
// @Tags         Admin/SettlementCompanies
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  model.CreateSettlementCompanyRequest  true  "结算公司信息"
// @Success      201  {object}  model.APIResponse[model.SettlementCompany]
// @Failure      400  {object}  model.ErrorResponse
// @Failure      409  {object}  model.ErrorResponse
// @Router       /admin/settlement-companies [post]
func (h *SettlementCompanyHandler) CreateSettlementCompany(c *gin.Context) {
	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req model.CreateSettlementCompanyRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	company, err := h.svc.CreateCompany(c.Request.Context(), &req, adminID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, company)
}

// UpdateSettlementCompany
// @Summary      更新结算公司
// @Tags         Admin/SettlementCompanies
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                                   true  "结算公司ID"
// @Param        request  body  model.UpdateSettlementCompanyRequest  true  "结算公司信息"
// @Success      200  {object}  model.APIResponse[model.SettlementCompany]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/settlement-companies/{id} [put]
func (h *SettlementCompanyHandler) UpdateSettlementCompany(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req model.UpdateSettlementCompanyRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	company, err := h.svc.UpdateCompany(c.Request.Context(), id, &req, adminID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondUpdated(c, company)
}

// ToggleSettlementCompanyStatus
// @Summary      启用/禁用结算公司
// @Tags         Admin/SettlementCompanies
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                true  "结算公司ID"
// @Param        request  body  ToggleStatusPayload  true  "状态信息"
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/settlement-companies/{id}/toggle [post]
func (h *SettlementCompanyHandler) ToggleSettlementCompanyStatus(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var payload ToggleStatusPayload
	if !ValidateAndRespond(c, &payload) {
		return
	}

	err := h.svc.ToggleCompanyStatus(c.Request.Context(), id, payload.Enabled, adminID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccessWithMsg[any](c, "status updated", nil)
}

// GetSettlementCompanyHistory
// @Summary      获取结算公司修改历史
// @Tags         Admin/SettlementCompanies
// @Security     BearerAuth
// @Param        id   path  int  true  "结算公司ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[[]model.SettlementCompanyHistory]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/settlement-companies/{id}/history [get]
func (h *SettlementCompanyHandler) GetSettlementCompanyHistory(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	histories, err := h.svc.GetCompanyHistory(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, histories)
}

// AssignPlayerToCompany
// @Summary      分配陪玩师到结算公司
// @Tags         Admin/SettlementCompanies
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                                true  "陪玩师ID"
// @Param        request  body  AssignPlayerToCompanyPayload       true  "分配信息"
// @Success      200  {object}  model.APIResponse[model.PlayerCompanyAssignment]
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/players/{id}/assign-company [post]
func (h *SettlementCompanyHandler) AssignPlayerToCompany(c *gin.Context) {
	playerID, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var payload AssignPlayerToCompanyPayload
	if !ValidateAndRespond(c, &payload) {
		return
	}

	req := &model.AssignPlayerToCompanyRequest{
		PlayerID:            playerID,
		SettlementCompanyID: payload.SettlementCompanyID,
		EffectiveDate:       payload.EffectiveDate,
		Reason:              payload.Reason,
	}

	assignment, err := h.svc.AssignPlayerToCompany(c.Request.Context(), req, adminID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, assignment)
}

// BatchAssignPlayersToCompany
// @Summary      批量分配陪玩师到结算公司
// @Tags         Admin/SettlementCompanies
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  model.BatchAssignPlayersRequest  true  "批量分配信息"
// @Success      200  {object}  model.APIResponse[BatchAssignResult]
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/players/batch-assign-company [post]
func (h *SettlementCompanyHandler) BatchAssignPlayersToCompany(c *gin.Context) {
	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req model.BatchAssignPlayersRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	count, err := h.svc.BatchAssignPlayers(c.Request.Context(), &req, adminID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, BatchAssignResult{
		AssignedCount: count,
		Message:       "players assigned successfully",
	})
}

// GetPlayerCurrentAssignment
// @Summary      获取陪玩师当前结算公司分配
// @Tags         Admin/SettlementCompanies
// @Security     BearerAuth
// @Param        id   path  int  true  "陪玩师ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[model.PlayerCompanyAssignment]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/players/{id}/current-company [get]
func (h *SettlementCompanyHandler) GetPlayerCurrentAssignment(c *gin.Context) {
	playerID, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	assignment, err := h.svc.GetCurrentAssignment(c.Request.Context(), playerID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, assignment)
}

// GetPlayerAssignmentHistory
// @Summary      获取陪玩师结算公司分配历史
// @Tags         Admin/SettlementCompanies
// @Security     BearerAuth
// @Param        id   path  int  true  "陪玩师ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[model.PlayerAssignmentHistoryResponse]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/players/{id}/company-history [get]
func (h *SettlementCompanyHandler) GetPlayerAssignmentHistory(c *gin.Context) {
	playerID, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	resp, err := h.svc.GetAssignmentHistory(c.Request.Context(), playerID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, resp)
}

// ToggleStatusPayload 启用/禁用状态请求体
type ToggleStatusPayload struct {
	Enabled bool `json:"enabled"`
}

// AssignPlayerToCompanyPayload 分配陪玩师到结算公司请求体
type AssignPlayerToCompanyPayload struct {
	SettlementCompanyID uint64    `json:"settlementCompanyId" binding:"required"`
	EffectiveDate       time.Time `json:"effectiveDate" binding:"required"`
	Reason              string    `json:"reason" binding:"required,max=500"`
}

// BatchAssignResult 批量分配结果
type BatchAssignResult struct {
	AssignedCount int    `json:"assignedCount"`
	Message       string `json:"message"`
}

// BatchOperationItem 批量操作项（用于Swagger文档）
type BatchOperationItem struct {
	ID      uint64 `json:"id"`
	Message string `json:"message,omitempty"`
}

// BatchOperationResult 批量操作结果（用于Swagger文档）
type BatchOperationResult struct {
	SuccessCount int                  `json:"success_count"`
	FailedCount  int                  `json:"failed_count"`
	TotalCount   int                  `json:"total_count"`
	FailedItems  []BatchOperationItem `json:"failed_items,omitempty"`
	SuccessItems []uint64             `json:"success_items,omitempty"`
}

// BatchUpdateCompanyStatus 批量更新结算公司状态
// @Summary      批量更新结算公司状态
// @Description  批量启用/禁用结算公司
// @Tags         Admin/SettlementCompanies
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchUpdateCompanyStatusRequest  true  "批量更新状态请求"
// @Success      200  {object}  model.APIResponse[BatchOperationResult]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/settlement-companies/batch/status [post]
func (h *SettlementCompanyHandler) BatchUpdateCompanyStatus(c *gin.Context) {
	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req BatchUpdateCompanyStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if len(req.CompanyIDs) == 0 {
		respondAPIError(c, apierr.BadRequest("company_ids is required"))
		return
	}
	if len(req.CompanyIDs) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 companies per batch"))
		return
	}

	result, err := h.svc.BatchUpdateCompanyStatus(c.Request.Context(), req.CompanyIDs, req.IsActive, adminID)
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch update company status failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}

// BatchDeleteCompanies 批量删除结算公司
// @Summary      批量删除结算公司
// @Description  批量删除结算公司（检查是否有陪玩师关联）
// @Tags         Admin/SettlementCompanies
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchDeleteCompaniesRequest  true  "批量删除请求"
// @Success      200  {object}  model.APIResponse[BatchOperationResult]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/settlement-companies/batch/delete [post]
func (h *SettlementCompanyHandler) BatchDeleteCompanies(c *gin.Context) {
	var req BatchDeleteCompaniesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if len(req.CompanyIDs) == 0 {
		respondAPIError(c, apierr.BadRequest("company_ids is required"))
		return
	}
	if len(req.CompanyIDs) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 companies per batch"))
		return
	}

	result, err := h.svc.BatchDeleteCompanies(c.Request.Context(), req.CompanyIDs)
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch delete companies failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}

// Batch Operation Request DTOs

// BatchUpdateCompanyStatusRequest 批量更新结算公司状态请求
type BatchUpdateCompanyStatusRequest struct {
	CompanyIDs []uint64 `json:"companyIds" binding:"required,min=1,max=100"`
	IsActive   bool     `json:"isActive" binding:"required"`
}

// BatchDeleteCompaniesRequest 批量删除结算公司请求
type BatchDeleteCompaniesRequest struct {
	CompanyIDs []uint64 `json:"companyIds" binding:"required,min=1,max=100"`
}

// RegisterSettlementCompanyRoutes 注册结算公司管理路由
// Requirements: 11.1-11.5, 12.1-12.5
func RegisterSettlementCompanyRoutes(router gin.IRouter, svc *svc.SettlementCompanyService, pm *mw.PermissionMiddleware) {
	handler := NewSettlementCompanyHandler(svc)

	// 结算公司管理
	// @Summary      列出结算公司
	// @Tags         Admin/SettlementCompanies
	// @Security     BearerAuth
	// @Router       /admin/settlement-companies [get]
	router.GET("/settlement-companies", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/settlement-companies"), handler.ListSettlementCompanies)

	// @Summary      获取结算公司详情
	// @Tags         Admin/SettlementCompanies
	// @Security     BearerAuth
	// @Router       /admin/settlement-companies/{id} [get]
	router.GET("/settlement-companies/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/settlement-companies/:id"), handler.GetSettlementCompany)

	// @Summary      创建结算公司
	// @Tags         Admin/SettlementCompanies
	// @Security     BearerAuth
	// @Router       /admin/settlement-companies [post]
	router.POST("/settlement-companies", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/settlement-companies"), handler.CreateSettlementCompany)

	// @Summary      更新结算公司
	// @Tags         Admin/SettlementCompanies
	// @Security     BearerAuth
	// @Router       /admin/settlement-companies/{id} [put]
	router.PUT("/settlement-companies/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/settlement-companies/:id"), handler.UpdateSettlementCompany)

	// @Summary      启用/禁用结算公司
	// @Tags         Admin/SettlementCompanies
	// @Security     BearerAuth
	// @Router       /admin/settlement-companies/{id}/toggle [post]
	router.POST("/settlement-companies/:id/toggle", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/settlement-companies/:id/toggle"), handler.ToggleSettlementCompanyStatus)

	// @Summary      获取结算公司修改历史
	// @Tags         Admin/SettlementCompanies
	// @Security     BearerAuth
	// @Router       /admin/settlement-companies/{id}/history [get]
	router.GET("/settlement-companies/:id/history", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/settlement-companies/:id/history"), handler.GetSettlementCompanyHistory)

	// 陪玩师结算公司分配
	// @Summary      分配陪玩师到结算公司
	// @Tags         Admin/SettlementCompanies
	// @Security     BearerAuth
	// @Router       /admin/players/{id}/assign-company [post]
	router.POST("/players/:id/assign-company", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/players/:id/assign-company"), handler.AssignPlayerToCompany)

	// @Summary      批量分配陪玩师到结算公司
	// @Tags         Admin/SettlementCompanies
	// @Security     BearerAuth
	// @Router       /admin/players/batch-assign-company [post]
	router.POST("/players/batch-assign-company", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/players/batch-assign-company"), handler.BatchAssignPlayersToCompany)

	// @Summary      获取陪玩师当前结算公司分配
	// @Tags         Admin/SettlementCompanies
	// @Security     BearerAuth
	// @Router       /admin/players/{id}/current-company [get]
	router.GET("/players/:id/current-company", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/players/:id/current-company"), handler.GetPlayerCurrentAssignment)

	// @Summary      获取陪玩师结算公司分配历史
	// @Tags         Admin/SettlementCompanies
	// @Security     BearerAuth
	// @Router       /admin/players/{id}/company-history [get]
	router.GET("/players/:id/company-history", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/players/:id/company-history"), handler.GetPlayerAssignmentHistory)

	// 批量操作
	// @Summary      批量更新结算公司状态
	// @Tags         Admin/SettlementCompanies
	// @Security     BearerAuth
	// @Router       /admin/settlement-companies/batch/status [post]
	router.POST("/settlement-companies/batch/status", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/settlement-companies/batch/status"), handler.BatchUpdateCompanyStatus)

	// @Summary      批量删除结算公司
	// @Tags         Admin/SettlementCompanies
	// @Security     BearerAuth
	// @Router       /admin/settlement-companies/batch/delete [post]
	router.POST("/settlement-companies/batch/delete", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/settlement-companies/batch/delete"), handler.BatchDeleteCompanies)
}
