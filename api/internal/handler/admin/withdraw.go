package admin

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	withdrawrepo "gamelink/internal/repository/withdraw"
	withdrawservice "gamelink/internal/service/withdraw"
	apierr "gamelink/pkg/apierr"
)

// Withdraw 提现模型（类型别名）
type Withdraw = model.Withdraw

// RegisterWithdrawRoutes 注册管理端提现管理路
func RegisterWithdrawRoutes(router gin.IRouter, withdrawRepo withdrawrepo.WithdrawRepository, withdrawService *withdrawservice.WithdrawRoutingService) {
	group := router.Group("/withdraws")
	{
		group.GET("", func(c *gin.Context) { listWithdrawsHandler(c, withdrawRepo) })
		group.GET("/:id", func(c *gin.Context) { getWithdrawHandler(c, withdrawRepo) })
		group.POST("/:id/approve", func(c *gin.Context) { approveWithdrawHandler(c, withdrawRepo) })
		group.POST("/:id/reject", func(c *gin.Context) { rejectWithdrawHandler(c, withdrawRepo) })
		group.POST("/:id/complete", func(c *gin.Context) { completeWithdrawHandler(c, withdrawRepo) })

		// 批量操作路由
		group.POST("/batch/approve", func(c *gin.Context) { batchApproveWithdrawalsHandler(c, withdrawService) })
		group.POST("/batch/reject", func(c *gin.Context) { batchRejectWithdrawalsHandler(c, withdrawService) })
		group.POST("/batch/complete", func(c *gin.Context) { batchCompleteWithdrawalsHandler(c, withdrawService) })
	}

	// 提现分流路由 (前端使用 /withdrawals 路径)
	withdrawalsGroup := router.Group("/withdrawals")
	{
		withdrawalsGroup.GET("/by-company", func(c *gin.Context) { listWithdrawsByCompanyHandler(c, withdrawRepo) })
		withdrawalsGroup.GET("/routing-stats", func(c *gin.Context) { getWithdrawRoutingStatsHandler(c, withdrawRepo) })
	}
}

// listWithdrawsHandler 获取提现申请列表
// @Summary      获取提现申请列表
// @Description  管理员查看所有提现申请，支持状态筛选、陪玩师筛选、分
// @Tags         Admin - Withdraw
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        status         query     string  false  "Status filter"
// @Param        playerId       query     int     false  "陪玩师ID"
// @Param        page           query     int     false  "页码"
// @Param        pageSize       query     int     false  "每页数量"
// @Success      200            {array}   model.Withdraw
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/withdraws [get]
func listWithdrawsHandler(c *gin.Context, repo withdrawrepo.WithdrawRepository) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}
	// 对提现列表额外加一层保护，避免 pageSize 过大
	if pageSize <= 0 {
		pageSize = 20
	} else if pageSize > 100 {
		pageSize = 100
	}

	opts := withdrawrepo.WithdrawListOptions{
		Page:     page,
		PageSize: pageSize,
	}

	// 状态筛
	if status := c.Query("status"); status != "" {
		s := model.WithdrawStatus(status)
		opts.Status = &s
	}

	// 陪玩师筛
	if playerIDStr := c.Query("playerId"); playerIDStr != "" {
		if playerID, err := strconv.ParseUint(playerIDStr, 10, 64); err == nil {
			opts.PlayerID = &playerID
		}
	}

	withdraws, total, err := repo.List(c.Request.Context(), opts)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, map[string]interface{}{
		"withdraws": withdraws,
		"total":     total,
	})
}

// getWithdrawHandler 获取提现详情
// @Summary      获取提现详情
// @Description  根据 ID 获取单个提现申请的详细信
// @Tags         Admin - Withdraw
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id             path      int     true  "提现ID"
// @Success      200            {object}  Withdraw
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/withdraws/{id} [get]
func getWithdrawHandler(c *gin.Context, repo withdrawrepo.WithdrawRepository) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	withdraw, err := repo.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, apierr.NotFound(apierr.ErrWithdrawNotFound))
		return
	}

	respondSuccess(c, *withdraw)
}

// ApproveWithdrawRequest 批准提现请求
type ApproveWithdrawRequest struct {
	Remark string `json:"remark"`
}

// approveWithdrawHandler 批准提现
// @Summary      批准提现
// @Description  批准一个待处理的提现申
// @Tags         Admin - Withdraw
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id             path      int                      true  "提现ID"
// @Param        request        body      ApproveWithdrawRequest  false  "审核备注"
// @Success      200            {object}  string
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/withdraws/{id}/approve [post]
func approveWithdrawHandler(c *gin.Context, repo withdrawrepo.WithdrawRepository) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	adminUserID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req ApproveWithdrawRequest
	_ = c.ShouldBindJSON(&req) // 可选字段，忽略绑定错误

	withdraw, err := repo.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, apierr.NotFound(apierr.ErrWithdrawNotFound))
		return
	}

	// 只能审批待处理的提现
	if withdraw.Status != model.WithdrawStatusPending {
		respondBadRequest(c, apierr.ErrWithdrawApproveInvalidStatus)
		return
	}

	// 更新状
	now := time.Now()
	withdraw.Status = model.WithdrawStatusApproved
	withdraw.ProcessedBy = &adminUserID
	withdraw.ProcessedAt = &now
	withdraw.AdminRemark = req.Remark

	if err := repo.Update(c.Request.Context(), withdraw); err != nil {
		respondError(c, err)
		return
	}

	respondMsg(c, "Withdraw approved successfully")
}

// RejectWithdrawRequest 拒绝提现请求
type RejectWithdrawRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// rejectWithdrawHandler 拒绝提现
// @Summary      拒绝提现
// @Description  拒绝一个待处理的提现申
// @Tags         Admin - Withdraw
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id             path      int                     true  "提现ID"
// @Param        request        body      RejectWithdrawRequest  true  "拒绝原因"
// @Success      200            {object}  string
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/withdraws/{id}/reject [post]
func rejectWithdrawHandler(c *gin.Context, repo withdrawrepo.WithdrawRepository) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	adminUserID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req RejectWithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	withdraw, err := repo.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, apierr.NotFound(apierr.ErrWithdrawNotFound))
		return
	}

	// 只能审批待处理的提现
	if withdraw.Status != model.WithdrawStatusPending {
		respondBadRequest(c, apierr.ErrWithdrawRejectInvalidStatus)
		return
	}

	// 更新状
	now := time.Now()
	withdraw.Status = model.WithdrawStatusRejected
	withdraw.ProcessedBy = &adminUserID
	withdraw.ProcessedAt = &now
	withdraw.RejectReason = req.Reason

	if err := repo.Update(c.Request.Context(), withdraw); err != nil {
		respondError(c, err)
		return
	}

	respondMsg(c, "Withdraw rejected")
}

// completeWithdrawHandler 完成提现（已打款
// @Summary      完成提现
// @Description  将一个已批准的提现申请标记为已完成（已打款）
// @Tags         Admin - Withdraw
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id             path      int     true  "提现ID"
// @Success      200            {object}  string
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/withdraws/{id}/complete [post]
func completeWithdrawHandler(c *gin.Context, repo withdrawrepo.WithdrawRepository) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	adminUserID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	withdraw, err := repo.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, apierr.NotFound(apierr.ErrWithdrawNotFound))
		return
	}

	// 只能完成已批准的提现
	if withdraw.Status != model.WithdrawStatusApproved {
		respondBadRequest(c, apierr.ErrWithdrawCompleteInvalidStatus)
		return
	}

	// 更新状
	now := time.Now()
	withdraw.Status = model.WithdrawStatusCompleted
	withdraw.CompletedAt = &now
	if withdraw.ProcessedBy == nil {
		withdraw.ProcessedBy = &adminUserID
	}

	if err := repo.Update(c.Request.Context(), withdraw); err != nil {
		respondError(c, err)
		return
	}

	respondMsg(c, "Withdraw completed successfully")
}

// listWithdrawsByCompanyHandler 按结算公司查询提现列表
// @Summary      按结算公司查询提现列表
// @Description  按结算公司筛选提现记录
// @Tags         Admin - Withdraw
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        settlementCompanyId  query     int     false  "结算公司ID"
// @Param        status               query     string  false  "状态"
// @Param        dateFrom             query     string  false  "开始日期"
// @Param        dateTo               query     string  false  "结束日期"
// @Param        page                 query     int     false  "页码"
// @Param        pageSize             query     int     false  "每页数量"
// @Success      200                  {array}   model.Withdraw
// @Router       /admin/withdrawals/by-company [get]
func listWithdrawsByCompanyHandler(c *gin.Context, repo withdrawrepo.WithdrawRepository) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	opts := withdrawrepo.WithdrawByCompanyOptions{
		Page:     page,
		PageSize: pageSize,
	}

	// 结算公司ID筛选
	if companyIDStr := c.Query("settlementCompanyId"); companyIDStr != "" {
		if companyID, err := strconv.ParseUint(companyIDStr, 10, 64); err == nil {
			opts.SettlementCompanyID = &companyID
		}
	}

	// 状态筛选
	if status := c.Query("status"); status != "" {
		s := model.WithdrawStatus(status)
		opts.Status = &s
	}

	// 日期筛选
	if dateFrom := c.Query("dateFrom"); dateFrom != "" {
		if t, err := time.Parse("2006-01-02", dateFrom); err == nil {
			opts.DateFrom = &t
		}
	}
	if dateTo := c.Query("dateTo"); dateTo != "" {
		if t, err := time.Parse("2006-01-02", dateTo); err == nil {
			nextDay := t.AddDate(0, 0, 1)
			opts.DateTo = &nextDay
		}
	}

	withdraws, total, err := repo.ListByCompany(c.Request.Context(), opts)
	if err != nil {
		respondError(c, err)
		return
	}

	// 转换为前端期望的格式
	type withdrawResponse struct {
		ID                    uint64  `json:"id"`
		PlayerID              uint64  `json:"playerId"`
		PlayerName            string  `json:"playerName"`
		Amount                float64 `json:"amount"`
		Status                string  `json:"status"`
		SettlementCompanyID   *uint64 `json:"settlementCompanyId"`
		SettlementCompanyName string  `json:"settlementCompanyName"`
		CreatedAt             string  `json:"createdAt"`
	}

	data := make([]withdrawResponse, len(withdraws))
	for i, w := range withdraws {
		playerName := "陪玩师 #" + strconv.FormatUint(w.PlayerID, 10)
		data[i] = withdrawResponse{
			ID:                    w.ID,
			PlayerID:              w.PlayerID,
			PlayerName:            playerName,
			Amount:                float64(w.AmountCents) / 100,
			Status:                string(w.Status),
			SettlementCompanyID:   w.SettlementCompanyID,
			SettlementCompanyName: w.SettlementCompanyName,
			CreatedAt:             w.CreatedAt.Format(time.RFC3339),
		}
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	respondList(c, data, &model.Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      int(total),
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	})
}

// getWithdrawRoutingStatsHandler 获取提现分流统计
// @Summary      获取提现分流统计
// @Description  获取提现分流统计数据
// @Tags         Admin - Withdraw
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        dateFrom  query     string  false  "开始日期"
// @Param        dateTo    query     string  false  "结束日期"
// @Success      200       {object}  model.WithdrawRoutingStatsResponse
// @Router       /admin/withdrawals/routing-stats [get]
func getWithdrawRoutingStatsHandler(c *gin.Context, repo withdrawrepo.WithdrawRepository) {
	var dateFrom, dateTo *time.Time

	if df := c.Query("dateFrom"); df != "" {
		if t, err := time.Parse("2006-01-02", df); err == nil {
			dateFrom = &t
		}
	}
	if dt := c.Query("dateTo"); dt != "" {
		if t, err := time.Parse("2006-01-02", dt); err == nil {
			nextDay := t.AddDate(0, 0, 1)
			dateTo = &nextDay
		}
	}

	stats, err := repo.GetRoutingStats(c.Request.Context(), dateFrom, dateTo)
	if err != nil {
		respondError(c, err)
		return
	}

	// 转换为前端期望的格式
	type companyStats struct {
		CompanyID   uint64  `json:"companyId"`
		CompanyName string  `json:"companyName"`
		Amount      float64 `json:"amount"`
		Count       int64   `json:"count"`
	}

	byCompany := make([]companyStats, len(stats.ByCompany))
	for i, cs := range stats.ByCompany {
		byCompany[i] = companyStats{
			CompanyID:   cs.SettlementCompanyID,
			CompanyName: cs.SettlementCompanyName,
			Amount:      float64(cs.TotalAmountCents) / 100,
			Count:       cs.TotalWithdrawals,
		}
	}

	// 计算待处理金额（需要额外查询）
	pendingOpts := withdrawrepo.WithdrawByCompanyOptions{
		Page:     1,
		PageSize: 1,
	}
	pendingStatus := model.WithdrawStatusPending
	pendingOpts.Status = &pendingStatus
	_, pendingTotal, _ := repo.ListByCompany(c.Request.Context(), pendingOpts)

	response := map[string]interface{}{
		"totalAmount":     float64(stats.TotalAmountCents) / 100,
		"totalCount":      stats.TotalWithdrawals,
		"completedAmount": float64(stats.TotalActualAmountCents) / 100,
		"completedCount":  stats.TotalWithdrawals,
		"pendingAmount":   0.0,
		"pendingCount":    pendingTotal,
		"byCompany":       byCompany,
	}

	respondSuccess(c, response)
}

// ============================================================================
// 批量操作处理器
// ============================================================================

// BatchApproveWithdrawalsRequest 批量批准提现请求
type BatchApproveWithdrawalsRequest struct {
	WithdrawIDs []uint64 `json:"withdrawIds" binding:"required,min=1,max=100"`
	Remark      string   `json:"remark" binding:"max=500"`
}

// batchApproveWithdrawalsHandler 批量批准提现申请
// @Summary      批量批准提现申请
// @Description  批量批准多个待处理的提现申请（最多100个）
// @Tags         Admin - Withdraw
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  BatchApproveWithdrawalsRequest  true  "批量批准请求"
// @Success      200      {object}  withdrawservice.BatchOperationResult
// @Failure      400      {object}  model.ErrorResponse
// @Failure      401      {object}  model.ErrorResponse
// @Router       /admin/withdraws/batch/approve [post]
func batchApproveWithdrawalsHandler(c *gin.Context, svc *withdrawservice.WithdrawRoutingService) {
	adminUserID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req BatchApproveWithdrawalsRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	result, err := svc.BatchApprove(c.Request.Context(), &withdrawservice.BatchApproveRequest{
		WithdrawIDs: req.WithdrawIDs,
		Remark:      req.Remark,
	}, adminUserID)

	if err != nil {
		respondError(c, err)
		return
	}

	message := "批量批准提现成功"
	if result.FailedCount > 0 {
		message = fmt.Sprintf("批量批准提现完成，成功%d个，失败%d个", result.SuccessCount, result.FailedCount)
	}

	respondSuccessWithMsg(c, message, result)
}

// BatchRejectWithdrawalsRequest 批量拒绝提现请求
type BatchRejectWithdrawalsRequest struct {
	WithdrawIDs []uint64 `json:"withdrawIds" binding:"required,min=1,max=100"`
	Reason      string   `json:"reason" binding:"required,min=1,max=500"`
}

// batchRejectWithdrawalsHandler 批量拒绝提现申请
// @Summary      批量拒绝提现申请
// @Description  批量拒绝多个待处理的提现申请（最多100个）
// @Tags         Admin - Withdraw
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  BatchRejectWithdrawalsRequest  true  "批量拒绝请求"
// @Success      200      {object}  withdrawservice.BatchOperationResult
// @Failure      400      {object}  model.ErrorResponse
// @Failure      401      {object}  model.ErrorResponse
// @Router       /admin/withdraws/batch/reject [post]
func batchRejectWithdrawalsHandler(c *gin.Context, svc *withdrawservice.WithdrawRoutingService) {
	adminUserID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req BatchRejectWithdrawalsRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	result, err := svc.BatchReject(c.Request.Context(), &withdrawservice.BatchRejectRequest{
		WithdrawIDs: req.WithdrawIDs,
		Reason:      req.Reason,
	}, adminUserID)

	if err != nil {
		respondError(c, err)
		return
	}

	message := "批量拒绝提现成功"
	if result.FailedCount > 0 {
		message = fmt.Sprintf("批量拒绝提现完成，成功%d个，失败%d个", result.SuccessCount, result.FailedCount)
	}

	respondSuccessWithMsg(c, message, result)
}

// BatchCompleteWithdrawalsRequest 批量完成提现请求
type BatchCompleteWithdrawalsRequest struct {
	WithdrawIDs []uint64 `json:"withdrawIds" binding:"required,min=1,max=100"`
}

// batchCompleteWithdrawalsHandler 批量完成提现（标记为已打款）
// @Summary      批量完成提现
// @Description  批量将多个已批准的提现申请标记为已完成（已打款）（最多100个）
// @Tags         Admin - Withdraw
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  BatchCompleteWithdrawalsRequest  true  "批量完成请求"
// @Success      200      {object}  withdrawservice.BatchOperationResult
// @Failure      400      {object}  model.ErrorResponse
// @Failure      401      {object}  model.ErrorResponse
// @Router       /admin/withdraws/batch/complete [post]
func batchCompleteWithdrawalsHandler(c *gin.Context, svc *withdrawservice.WithdrawRoutingService) {
	adminUserID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req BatchCompleteWithdrawalsRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	result, err := svc.BatchComplete(c.Request.Context(), &withdrawservice.BatchCompleteRequest{
		WithdrawIDs: req.WithdrawIDs,
	}, adminUserID)

	if err != nil {
		respondError(c, err)
		return
	}

	message := "批量完成提现成功"
	if result.FailedCount > 0 {
		message = fmt.Sprintf("批量完成提现完成，成功%d个，失败%d个", result.SuccessCount, result.FailedCount)
	}

	respondSuccessWithMsg(c, message, result)
}
