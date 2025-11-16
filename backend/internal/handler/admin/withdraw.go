package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	withdrawrepo "gamelink/internal/repository/withdraw"
)

// RegisterWithdrawRoutes 注册管理端提现管理路由
func RegisterWithdrawRoutes(router gin.IRouter, withdrawRepo withdrawrepo.WithdrawRepository) {
	group := router.Group("/withdraws")
	{
		group.GET("", func(c *gin.Context) { listWithdrawsHandler(c, withdrawRepo) })
		group.GET("/:id", func(c *gin.Context) { getWithdrawHandler(c, withdrawRepo) })
		group.POST("/:id/approve", func(c *gin.Context) { approveWithdrawHandler(c, withdrawRepo) })
		group.POST("/:id/reject", func(c *gin.Context) { rejectWithdrawHandler(c, withdrawRepo) })
		group.POST("/:id/complete", func(c *gin.Context) { completeWithdrawHandler(c, withdrawRepo) })
	}
}

// listWithdrawsHandler 获取提现申请列表
// @Summary      获取提现申请列表
// @Description  管理员查看所有提现申请
// @Tags         Admin - Withdraw
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        status         query     string  false  "Status filter"
// @Param        playerId       query     int     false  "陪玩师ID"
// @Param        page           query     int     false  "页码"
// @Param        pageSize       query     int     false  "每页数量"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/withdraws [get]
func listWithdrawsHandler(c *gin.Context, repo withdrawrepo.WithdrawRepository) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if err != nil || pageSize <= 0 {
		pageSize = 20
	} else if pageSize > 100 {
		pageSize = 100
	}

	opts := withdrawrepo.WithdrawListOptions{
		Page:     page,
		PageSize: pageSize,
	}

	// 状态筛选
	if status := c.Query("status"); status != "" {
		s := model.WithdrawStatus(status)
		opts.Status = &s
	}

	// 陪玩师筛选
	if playerIDStr := c.Query("playerId"); playerIDStr != "" {
		if playerID, err := strconv.ParseUint(playerIDStr, 10, 64); err == nil {
			opts.PlayerID = &playerID
		}
	}

	withdraws, total, err := repo.List(c.Request.Context(), opts)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data: map[string]interface{}{
			"withdraws": withdraws,
			"total":     total,
		},
	})
}

// getWithdrawHandler 获取提现详情
// @Summary      获取提现详情
// @Description  API endpoint
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      int     true  "提现ID"
// @Success      200            {object}  model.APIResponse[model.Withdraw]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/withdraws/{id} [get]
func getWithdrawHandler(c *gin.Context, repo withdrawrepo.WithdrawRepository) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "Invalid withdraw ID")
		return
	}

	withdraw, err := repo.Get(c.Request.Context(), id)
	if err != nil {
		writeJSONError(c, http.StatusNotFound, "Withdraw not found")
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[model.Withdraw]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    *withdraw,
	})
}

// ApproveWithdrawRequest 批准提现请求
type ApproveWithdrawRequest struct {
	Remark string `json:"remark"`
}

// approveWithdrawHandler 批准提现
// @Summary      批准提现
// @Description  API endpoint
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                   true  "Bearer {token}"
// @Param        id             path      int                      true  "提现ID"
// @Param        request        body      ApproveWithdrawRequest  false  "审核备注"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/withdraws/{id}/approve [post]
func approveWithdrawHandler(c *gin.Context, repo withdrawrepo.WithdrawRepository) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "Invalid withdraw ID")
		return
	}

	adminVal, ok := c.Get("user_id")
	if !ok {
		writeJSONError(c, http.StatusUnauthorized, "missing admin user")
		return
	}
	adminUserID, ok := adminVal.(uint64)
	if !ok {
		writeJSONError(c, http.StatusInternalServerError, "invalid admin user id type")
		return
	}

	var req ApproveWithdrawRequest
	c.ShouldBindJSON(&req)

	withdraw, err := repo.Get(c.Request.Context(), id)
	if err != nil {
		writeJSONError(c, http.StatusNotFound, "Withdraw not found")
		return
	}

	// 只能审批待处理的提现
	if withdraw.Status != model.WithdrawStatusPending {
		writeJSONError(c, http.StatusBadRequest, "Can only approve pending withdraws")
		return
	}

	// 更新状态
	now := time.Now()
	withdraw.Status = model.WithdrawStatusApproved
	withdraw.ProcessedBy = &adminUserID
	withdraw.ProcessedAt = &now
	withdraw.AdminRemark = req.Remark

	if err := repo.Update(c.Request.Context(), withdraw); err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "Withdraw approved successfully",
	})
}

// RejectWithdrawRequest 拒绝提现请求
type RejectWithdrawRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// rejectWithdrawHandler 拒绝提现
// @Summary      拒绝提现
// @Description  API endpoint
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                  true  "Bearer {token}"
// @Param        id             path      int                     true  "提现ID"
// @Param        request        body      RejectWithdrawRequest  true  "拒绝原因"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/withdraws/{id}/reject [post]
func rejectWithdrawHandler(c *gin.Context, repo withdrawrepo.WithdrawRepository) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "Invalid withdraw ID")
		return
	}

	adminVal, ok := c.Get("user_id")
	if !ok {
		writeJSONError(c, http.StatusUnauthorized, "missing admin user")
		return
	}
	adminUserID, ok := adminVal.(uint64)
	if !ok {
		writeJSONError(c, http.StatusInternalServerError, "invalid admin user id type")
		return
	}

	var req RejectWithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	withdraw, err := repo.Get(c.Request.Context(), id)
	if err != nil {
		writeJSONError(c, http.StatusNotFound, "Withdraw not found")
		return
	}

	// 只能审批待处理的提现
	if withdraw.Status != model.WithdrawStatusPending {
		writeJSONError(c, http.StatusBadRequest, "Can only reject pending withdraws")
		return
	}

	// 更新状态
	now := time.Now()
	withdraw.Status = model.WithdrawStatusRejected
	withdraw.ProcessedBy = &adminUserID
	withdraw.ProcessedAt = &now
	withdraw.RejectReason = req.Reason

	if err := repo.Update(c.Request.Context(), withdraw); err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "Withdraw rejected",
	})
}

// completeWithdrawHandler 完成提现（已打款）
// @Summary      完成提现
// @Description  API endpoint
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      int     true  "提现ID"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/withdraws/{id}/complete [post]
func completeWithdrawHandler(c *gin.Context, repo withdrawrepo.WithdrawRepository) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "Invalid withdraw ID")
		return
	}

	adminVal, ok := c.Get("user_id")
	if !ok {
		writeJSONError(c, http.StatusUnauthorized, "missing admin user")
		return
	}
	adminUserID, ok := adminVal.(uint64)
	if !ok {
		writeJSONError(c, http.StatusInternalServerError, "invalid admin user id type")
		return
	}

	withdraw, err := repo.Get(c.Request.Context(), id)
	if err != nil {
		writeJSONError(c, http.StatusNotFound, "Withdraw not found")
		return
	}

	// 只能完成已批准的提现
	if withdraw.Status != model.WithdrawStatusApproved {
		writeJSONError(c, http.StatusBadRequest, "Can only complete approved withdraws")
		return
	}

	// 更新状态
	now := time.Now()
	withdraw.Status = model.WithdrawStatusCompleted
	withdraw.CompletedAt = &now
	if withdraw.ProcessedBy == nil {
		withdraw.ProcessedBy = &adminUserID
	}

	if err := repo.Update(c.Request.Context(), withdraw); err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "Withdraw completed successfully",
	})
}
