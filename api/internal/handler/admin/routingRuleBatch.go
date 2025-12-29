package admin

import (
	"github.com/gin-gonic/gin"

	"gamelink/pkg/apierr"
)

// BatchUpdateRoutingRuleStatus 批量更新分流规则状态
// @Summary      批量更新分流规则状态
// @Description  批量启用/禁用分流规则
// @Tags         Admin/RoutingRules
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchUpdateRoutingRuleStatusRequest  true  "规则ID列表和状态"
// @Success      200  {object}  model.APIResponse[BatchOperationResponse]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/routing-rules/batch/status [post]
func (h *RoutingRuleHandler) BatchUpdateRoutingRuleStatus(c *gin.Context) {
	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req BatchUpdateRoutingRuleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if len(req.RuleIDs) == 0 {
		respondAPIError(c, apierr.BadRequest("rule_ids is required"))
		return
	}
	if len(req.RuleIDs) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 rules per batch"))
		return
	}

	result, err := h.svc.BatchUpdateRuleStatus(c.Request.Context(), req.RuleIDs, req.IsActive, adminID)
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch update rule status failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}

// BatchDeleteRoutingRules 批量删除分流规则
// @Summary      批量删除分流规则
// @Description  批量删除分流规则
// @Tags         Admin/RoutingRules
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchDeleteRoutingRulesRequest  true  "规则ID列表"
// @Success      200  {object}  model.APIResponse[BatchOperationResponse]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/routing-rules/batch/delete [post]
func (h *RoutingRuleHandler) BatchDeleteRoutingRules(c *gin.Context) {
	var req BatchDeleteRoutingRulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if len(req.RuleIDs) == 0 {
		respondAPIError(c, apierr.BadRequest("rule_ids is required"))
		return
	}
	if len(req.RuleIDs) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 rules per batch"))
		return
	}

	result, err := h.svc.BatchDeleteRules(c.Request.Context(), req.RuleIDs)
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch delete rules failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}

// Batch Operation Request DTOs

// BatchUpdateRoutingRuleStatusRequest 批量更新分流规则状态请求
type BatchUpdateRoutingRuleStatusRequest struct {
	RuleIDs  []uint64 `json:"rule_ids" binding:"required,min=1,max=100"`
	IsActive bool     `json:"is_active" binding:"required"`
}

// BatchDeleteRoutingRulesRequest 批量删除分流规则请求
type BatchDeleteRoutingRulesRequest struct {
	RuleIDs []uint64 `json:"rule_ids" binding:"required,min=1,max=100"`
}
