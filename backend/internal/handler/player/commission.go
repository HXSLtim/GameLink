package player

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/service/commission"
)

// RegisterCommissionRoutes 注册陪玩师端抽成管理路由
func RegisterCommissionRoutes(router gin.IRouter, svc *commission.CommissionService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/player/commission")
	group.Use(authMiddleware) // 需要认�?
	{
		group.GET("/summary", func(c *gin.Context) { getCommissionSummaryHandler(c, svc) })
		group.GET("/records", func(c *gin.Context) { getCommissionRecordsHandler(c, svc) })
		group.GET("/settlements", func(c *gin.Context) { getMonthlySettlementsHandler(c, svc) })
	}
}

// getCommissionSummaryHandler 获取抽成汇�?
// @Summary      获取抽成汇�?
// @Description  获取陪玩师抽成收入汇�?
// @Tags         Player - Commission
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        month          query     string  true  "月份 (YYYY-MM)"
// @Success      200            {object}  model.APIResponse[commission.CommissionSummaryResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /player/commission/summary [get]
func getCommissionSummaryHandler(c *gin.Context, svc *commission.CommissionService) {
	userID := getUserIDFromContext(c)

	// 获取月份参数，默认当前月
	month := c.DefaultQuery("month", time.Now().Format("2006-01"))

	// 查找该用户的陪玩师ID
	playerID, err := getPlayerIDByUserID(c, userID)
	if err != nil {
		respondError(c, http.StatusNotFound, "Player not found")
		return
	}

	resp, err := svc.GetPlayerCommissionSummary(c.Request.Context(), playerID, month)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[commission.CommissionSummaryResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    *resp,
	})
}

// getCommissionRecordsHandler 获取抽成记录
// @Summary      获取抽成记录
// @Description  获取陪玩师抽成记录列�?
// @Tags         Player - Commission
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        page           query     int     false  "页码"
// @Param        pageSize       query     int     false  "每页数量"
// @Success      200            {object}  model.APIResponse[commission.CommissionRecordListResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /player/commission/records [get]
func getCommissionRecordsHandler(c *gin.Context, svc *commission.CommissionService) {
	userID := getUserIDFromContext(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// 查找该用户的陪玩师ID
	playerID, err := getPlayerIDByUserID(c, userID)
	if err != nil {
		respondError(c, http.StatusNotFound, "Player not found")
		return
	}

	resp, err := svc.GetCommissionRecords(c.Request.Context(), playerID, page, pageSize)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[commission.CommissionRecordListResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    *resp,
	})
}

// getMonthlySettlementsHandler 获取月度结算记录
// @Summary      获取月度结算记录
// @Description  获取陪玩师月度结算记�?
// @Tags         Player - Commission
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        page           query     int     false  "页码"
// @Param        pageSize       query     int     false  "每页数量"
// @Success      200            {object}  model.APIResponse[commission.SettlementListResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /player/commission/settlements [get]
func getMonthlySettlementsHandler(c *gin.Context, svc *commission.CommissionService) {
	userID := getUserIDFromContext(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// 查找该用户的陪玩师ID
	playerID, err := getPlayerIDByUserID(c, userID)
	if err != nil {
		respondError(c, http.StatusNotFound, "Player not found")
		return
	}

	resp, err := svc.GetMonthlySettlements(c.Request.Context(), playerID, page, pageSize)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[commission.SettlementListResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    *resp,
	})
}

// getPlayerIDByUserID 根据用户ID获取陪玩师ID
func getPlayerIDByUserID(c *gin.Context, userID uint64) (uint64, error) {
	// TODO: 优化这个查询，可以在用户上下文中缓存playerID
	// 这里简化处理，实际应该从service层获�?
	// 暂时返回userID作为playerID（需要后续完善）
	return userID, nil
}

