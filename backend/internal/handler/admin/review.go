package admin

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	adminservice "gamelink/internal/service/admin"
	apierr "gamelink/pkg/apierr"
)

// Review 评价模型（类型别名）
type Review = model.Review

// OperationLogWithActor 带操作员名称的操作日志
type OperationLogWithActor struct {
	model.OperationLog
	ActorName string `json:"actorName,omitempty"`
}

// ReviewHandler 管理评价接口
type ReviewHandler struct{ svc *adminservice.AdminService }

func NewReviewHandler(s *adminservice.AdminService) *ReviewHandler { return &ReviewHandler{svc: s} }

// ListReviews
// @Summary      评价列表
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Produce      json
// @Param        page       query  int  false  "页码"
// @Param        pageSize   query     int       false  "每页数量"
// @Param        orderId     query     int       false  "订单ID"
// @Param        userId     query     int       false  "用户ID"
// @Param        playerId   query     int       false  "陪玩师ID"
// @Param        dateFrom       query    string       false  "Start date (YYYY-MM-DD)"// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
// @Success      200  {object}  model.APIResponse[[]Review]
// @Router       /admin/reviews [get]
func (h *ReviewHandler) ListReviews(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}
	var orderID, userID, playerID *uint64
	if v, err := queryUint64Ptr(c, "order_id"); err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidOrderID)
		return
	} else {
		orderID = v
	}
	if v, err := queryUint64Ptr(c, "user_id"); err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidUserID)
		return
	} else {
		userID = v
	}
	if v, err := queryUint64Ptr(c, "player_id"); err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidPlayerID)
		return
	} else {
		playerID = v
	}
	dateFrom, err := queryTimePtr(c, "date_from")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidDateFrom)
		return
	}
	dateTo, err := queryTimePtr(c, "date_to")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidDateTo)
		return
	}
	items, p, err := h.svc.ListReviews(c.Request.Context(), repository.ReviewListOptions{Page: page, PageSize: pageSize, OrderID: orderID, UserID: userID, PlayerID: playerID, DateFrom: dateFrom, DateTo: dateTo})
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}
	items = ensureSlice(items)
	writeJSON(c, 200, model.APIResponse[[]model.Review]{Success: true, Code: 200, Message: "OK", Data: items, Pagination: p})
}

// GetReview
// @Summary      获取评价
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "评价ID"
// @Success      200  {object}  model.APIResponse[Review]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/reviews/{id} [get]
func (h *ReviewHandler) GetReview(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidID)
		return
	}
	item, err := h.svc.GetReview(c.Request.Context(), id)
	if errors.Is(err, adminservice.ErrNotFound) {
		_ = c.Error(adminservice.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}
	writeJSON(c, 200, model.APIResponse[*model.Review]{Success: true, Code: 200, Message: "OK", Data: item})
}

// CreateReview
// @Summary      创建评价
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  CreateReviewPayload  true  "评价"
// @Success      201  {object}  model.APIResponse[Review]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var p CreateReviewPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidJSONPayload)
		return
	}
	r := model.Review{OrderID: p.OrderID, UserID: p.UserID, PlayerID: p.PlayerID, Score: model.Rating(p.Score), Content: strings.TrimSpace(p.Content)}
	out, err := h.svc.CreateReview(c.Request.Context(), r)
	if errors.Is(err, apierr.BadRequest("validation failed")) {
		_ = c.Error(apierr.BadRequest("validation failed"))
		return
	}
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}
	writeJSON(c, 201, model.APIResponse[*model.Review]{Success: true, Code: 201, Message: "created", Data: out})
}

// UpdateReview
// @Summary      更新评价
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                    true  "评价ID"
// @Param        request  body  UpdateReviewPayload    true  "评价"
// @Success      200  {object}  model.APIResponse[Review]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/reviews/{id} [put]
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidID)
		return
	}
	var p UpdateReviewPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidJSONPayload)
		return
	}
	out, err := h.svc.UpdateReview(c.Request.Context(), id, model.Rating(p.Score), p.Content)
	if errors.Is(err, apierr.BadRequest("validation failed")) {
		_ = c.Error(apierr.BadRequest("validation failed"))
		return
	}
	if errors.Is(err, adminservice.ErrNotFound) {
		_ = c.Error(adminservice.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}
	writeJSON(c, 200, model.APIResponse[*model.Review]{Success: true, Code: 200, Message: "updated", Data: out})
}

// DeleteReview
// @Summary      删除评价
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "评价ID"
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/reviews/{id} [delete]
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidID)
		return
	}
	err = h.svc.DeleteReview(c.Request.Context(), id)
	if errors.Is(err, adminservice.ErrNotFound) {
		_ = c.Error(adminservice.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}
	writeJSON(c, 200, model.APIResponse[any]{Success: true, Code: 200, Message: "deleted"})
}

// ListReviewLogs
// @Summary      获取评价操作日志
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Produce      json
// @Param        id           path   int  true  "评价ID"
// @Param        page         query  int  false "页码"
// @Param        pageSize   query     int       false  "每页数量"
// @Param        action       query  string false "动作过滤" Enums(create,update,delete)
// @Param        actor_user_id query int   false "操作者用户ID"
// @Param        dateFrom       query    string       false  "Start date (YYYY-MM-DD)"// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
// @Param        export       query  string false "导出格式" Enums(csv)
// @Param        fields         query    string       false  "Export fields (comma separated)"// @Param        header_lang  query  string false "列头语言" Enums(en,zh)
// @Success      200  {object}  model.APIResponse[[]model.OperationLog]
// @Router       /admin/reviews/{id}/logs [get]
func (h *ReviewHandler) ListReviewLogs(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidID)
		return
	}
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}
	var actorID *uint64
	if v, err := queryUint64Ptr(c, "actor_user_id"); err == nil {
		actorID = v
	}
	var dateFrom, dateTo *time.Time
	if v, err := queryTimePtr(c, "date_from"); err == nil {
		dateFrom = v
	} else {
		writeJSONError(c, 400, apierr.ErrInvalidDateFrom)
		return
	}
	if v, err := queryTimePtr(c, "date_to"); err == nil {
		dateTo = v
	} else {
		writeJSONError(c, 400, apierr.ErrInvalidDateTo)
		return
	}
	opts := repository.OperationLogListOptions{Page: page, PageSize: pageSize, Action: strings.TrimSpace(c.Query("action")), ActorUserID: actorID, DateFrom: dateFrom, DateTo: dateTo}
	items, p, err := h.svc.ListOperationLogs(c.Request.Context(), "review", id, opts)
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}
	if strings.EqualFold(strings.TrimSpace(c.Query("export")), "csv") {
		exportOperationLogsCSV(c, "review", id, items)
		return
	}

	// 填充操作员名称
	result := make([]OperationLogWithActor, 0, len(items))
	userIDs := make([]uint64, 0)
	for _, item := range items {
		if item.ActorUserID != nil {
			userIDs = append(userIDs, *item.ActorUserID)
		}
	}

	// 批量获取用户名称
	userNames := make(map[uint64]string)
	if len(userIDs) > 0 {
		users, _ := h.svc.GetUsersByIDs(c.Request.Context(), userIDs)
		for _, u := range users {
			userNames[u.ID] = u.Name
		}
	}

	for _, item := range items {
		logWithActor := OperationLogWithActor{OperationLog: item}
		if item.ActorUserID != nil {
			if name, ok := userNames[*item.ActorUserID]; ok {
				logWithActor.ActorName = name
			}
		}
		result = append(result, logWithActor)
	}

	writeJSON(c, 200, model.APIResponse[[]OperationLogWithActor]{Success: true, Code: 200, Message: "OK", Data: result, Pagination: p})
}

// ListPlayerReviews
// @Summary      获取陪玩师的评价
// @Tags         Admin/Players
// @Security     BearerAuth
// @Produce      json
// @Param        id         path   int  true  "陪玩师ID"
// @Param        page       query  int  false  "页码"
// @Param        pageSize   query     int       false  "每页数量"
// @Success      200  {object}  model.APIResponse[[]Review]
// @Router       /admin/players/{id}/reviews [get]
func (h *ReviewHandler) ListPlayerReviews(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidID)
		return
	}
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}
	pid := id
	items, p, err := h.svc.ListReviews(c.Request.Context(), repository.ReviewListOptions{Page: page, PageSize: pageSize, PlayerID: &pid})
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}
	items = ensureSlice(items)
	writeJSON(c, 200, model.APIResponse[[]model.Review]{Success: true, Code: 200, Message: "OK", Data: items, Pagination: p})
}

type CreateReviewPayload struct {
	OrderID  uint64 `json:"order_id" binding:"required"`
	UserID   uint64 `json:"user_id" binding:"required"`
	PlayerID uint64 `json:"player_id" binding:"required"`
	Score    uint8  `json:"score" binding:"required"`
	Content  string `json:"content"`
}

type UpdateReviewPayload struct {
	Score   uint8  `json:"score" binding:"required"`
	Content string `json:"content"`
}

// CreateReviewReport
// @Summary      创建评价举报
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                      true  "评价ID"
// @Param        request  body  CreateReviewReportPayload  true  "举报信息"
// @Success      201  {object}  model.APIResponse[CreateReviewReportResponse]
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/reviews/{id}/reports [post]
func (h *ReviewHandler) CreateReviewReport(c *gin.Context) {
	reviewID, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidID)
		return
	}

	var p CreateReviewReportPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidJSONPayload)
		return
	}

	// Get reporter ID from context (authenticated user)
	reporterID, exists := c.Get("user_id")
	if !exists {
		writeJSONError(c, 401, "unauthorized")
		return
	}

	out, err := h.svc.ReportReview(c.Request.Context(), reviewID, reporterID.(uint64), p.Reason, p.Evidence)
	if errors.Is(err, adminservice.ErrNotFound) {
		_ = c.Error(adminservice.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}

	writeJSON(c, 201, model.APIResponse[CreateReviewReportResponse]{
		Success: true,
		Code:    201,
		Message: "report created",
		Data:    CreateReviewReportResponse{ReportID: out},
	})
}

// ListReviewReports
// @Summary      列出评价举报
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Produce      json
// @Param        page         query  int     false  "页码"
// @Param        pageSize     query  int     false  "每页数量"
// @Param        review_id    query  int     false  "评价ID"
// @Param        reporter_id  query  int     false  "举报人ID"
// @Param        status       query  string  false  "状态" Enums(pending,approved,rejected)
// @Param        date_from    query  string  false  "开始日期 (YYYY-MM-DD)"
// @Param        date_to      query  string  false  "结束日期 (YYYY-MM-DD)"
// @Success      200  {object}  model.APIResponse[[]ReviewReportDTO]
// @Router       /admin/review-reports [get]
func (h *ReviewHandler) ListReviewReports(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	var reviewID, reporterID *uint64
	if v, err := queryUint64Ptr(c, "review_id"); err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidID)
		return
	} else {
		reviewID = v
	}

	if v, err := queryUint64Ptr(c, "reporter_id"); err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidUserID)
		return
	} else {
		reporterID = v
	}

	var status *model.ReviewReportStatus
	if statusStr := strings.TrimSpace(c.Query("status")); statusStr != "" {
		s := model.ReviewReportStatus(statusStr)
		if !s.Valid() {
			writeJSONError(c, 400, "invalid status")
			return
		}
		status = &s
	}

	dateFrom, err := queryTimePtr(c, "date_from")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidDateFrom)
		return
	}

	dateTo, err := queryTimePtr(c, "date_to")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidDateTo)
		return
	}

	items, p, err := h.svc.ListReviewReports(c.Request.Context(), page, pageSize, reviewID, reporterID, status, dateFrom, dateTo)
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}

	items = ensureSlice(items)
	writeJSON(c, 200, model.APIResponse[[]adminservice.ReviewReportDTO]{
		Success:    true,
		Code:       200,
		Message:    "OK",
		Data:       items,
		Pagination: p,
	})
}

// GetReviewReport
// @Summary      获取举报详情
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "举报ID"
// @Success      200  {object}  model.APIResponse[ReviewReportDTO]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/review-reports/{id} [get]
func (h *ReviewHandler) GetReviewReport(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidID)
		return
	}

	item, err := h.svc.GetReviewReport(c.Request.Context(), id)
	if errors.Is(err, adminservice.ErrNotFound) {
		_ = c.Error(adminservice.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}

	writeJSON(c, 200, model.APIResponse[*adminservice.ReviewReportDTO]{
		Success: true,
		Code:    200,
		Message: "OK",
		Data:    item,
	})
}

// HandleReviewReport
// @Summary      处理评价举报
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                       true  "举报ID"
// @Param        request  body  HandleReviewReportPayload  true  "处理信息"
// @Success      200  {object}  model.APIResponse[HandleReviewReportResponse]
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/review-reports/{id}/handle [put]
func (h *ReviewHandler) HandleReviewReport(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidID)
		return
	}

	var p HandleReviewReportPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidJSONPayload)
		return
	}

	// Validate action
	if p.Action != "delete" && p.Action != "warn" && p.Action != "reject" {
		writeJSONError(c, 400, "invalid action, must be one of: delete, warn, reject")
		return
	}

	// Get handler ID from context (authenticated admin)
	handlerID, exists := c.Get("user_id")
	if !exists {
		writeJSONError(c, 401, "unauthorized")
		return
	}

	out, err := h.svc.HandleReviewReport(c.Request.Context(), id, handlerID.(uint64), p.Action, p.Note)
	if errors.Is(err, adminservice.ErrNotFound) {
		_ = c.Error(adminservice.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}

	writeJSON(c, 200, model.APIResponse[HandleReviewReportResponse]{
		Success: true,
		Code:    200,
		Message: out.Message,
		Data:    HandleReviewReportResponse{Status: out.Status, Message: out.Message},
	})
}

// Payload types for review reports
type CreateReviewReportPayload struct {
	Reason   string `json:"reason" binding:"required,max=500"`
	Evidence string `json:"evidence" binding:"max=1000"`
}

type CreateReviewReportResponse struct {
	ReportID uint64 `json:"reportId"`
}

type ReviewReportDTO = adminservice.ReviewReportDTO

type HandleReviewReportPayload struct {
	Action string `json:"action" binding:"required,oneof=delete warn reject"`
	Note   string `json:"note" binding:"max=500"`
}

type HandleReviewReportResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// ListPendingReviews
// @Summary      获取待审核评价列表
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Produce      json
// @Param        page       query  int  false  "页码"
// @Param        pageSize   query  int  false  "每页数量"
// @Success      200  {object}  model.APIResponse[[]Review]
// @Router       /admin/reviews/pending [get]
func (h *ReviewHandler) ListPendingReviews(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	items, total, err := h.svc.ListPendingReviews(c.Request.Context(), page, pageSize)
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}

	items = ensureSlice(items)
	p := &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	}
	writeJSON(c, 200, model.APIResponse[[]model.Review]{
		Success:    true,
		Code:       200,
		Message:    "OK",
		Data:       items,
		Pagination: p,
	})
}

// ApproveReview
// @Summary      批准评价
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                    true  "评价ID"
// @Param        request  body  ApproveReviewPayload   false "批准信息"
// @Success      200  {object}  model.APIResponse[any]
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/reviews/{id}/approve [put]
func (h *ReviewHandler) ApproveReview(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidID)
		return
	}

	// Parse optional reason from body
	var p ApproveReviewPayload
	_ = c.ShouldBindJSON(&p) // Ignore error, reason is optional

	reason := p.Reason
	if reason == "" {
		reason = "批准评价"
	}

	// Get actor user ID from context
	var actorUserID *uint64
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint64)
		actorUserID = &uid
	}

	err = h.svc.ApproveReview(c.Request.Context(), id, reason, actorUserID)
	if errors.Is(err, adminservice.ErrNotFound) {
		_ = c.Error(adminservice.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}

	writeJSON(c, 200, model.APIResponse[any]{
		Success: true,
		Code:    200,
		Message: "review approved",
	})
}

// RejectReview
// @Summary      拒绝评价
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                   true  "评价ID"
// @Param        request  body  RejectReviewPayload   true  "拒绝信息"
// @Success      200  {object}  model.APIResponse[any]
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/reviews/{id}/reject [put]
func (h *ReviewHandler) RejectReview(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidID)
		return
	}

	var p RejectReviewPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidJSONPayload)
		return
	}

	// Get actor user ID from context
	var actorUserID *uint64
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint64)
		actorUserID = &uid
	}

	err = h.svc.RejectReview(c.Request.Context(), id, p.Reason, actorUserID)
	if errors.Is(err, adminservice.ErrNotFound) {
		_ = c.Error(adminservice.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}

	writeJSON(c, 200, model.APIResponse[any]{
		Success: true,
		Code:    200,
		Message: "review rejected",
	})
}

// BatchApproveReviews
// @Summary      批量批准评价
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchApprovePayload  true  "批量批准信息"
// @Success      200  {object}  model.APIResponse[any]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/reviews/batch-approve [put]
func (h *ReviewHandler) BatchApproveReviews(c *gin.Context) {
	var p BatchApprovePayload
	if err := c.ShouldBindJSON(&p); err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidJSONPayload)
		return
	}

	// Get actor user ID from context
	var actorUserID *uint64
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint64)
		actorUserID = &uid
	}

	err := h.svc.BatchApproveReviews(c.Request.Context(), p.ReviewIDs, actorUserID)
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}

	writeJSON(c, 200, model.APIResponse[any]{
		Success: true,
		Code:    200,
		Message: "reviews approved",
	})
}

// BatchRejectReviews
// @Summary      批量拒绝评价
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchRejectPayload  true  "批量拒绝信息"
// @Success      200  {object}  model.APIResponse[any]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/reviews/batch-reject [put]
func (h *ReviewHandler) BatchRejectReviews(c *gin.Context) {
	var p BatchRejectPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidJSONPayload)
		return
	}

	// Get actor user ID from context
	var actorUserID *uint64
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint64)
		actorUserID = &uid
	}

	err := h.svc.BatchRejectReviews(c.Request.Context(), p.ReviewIDs, p.Reason, actorUserID)
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}

	writeJSON(c, 200, model.APIResponse[any]{
		Success: true,
		Code:    200,
		Message: "reviews rejected",
	})
}

// Payload types for review moderation
type ApproveReviewPayload struct {
	Reason string `json:"reason" binding:"max=500"`
}

type RejectReviewPayload struct {
	Reason string `json:"reason" binding:"required,max=500"`
}

type BatchApprovePayload struct {
	ReviewIDs []uint64 `json:"reviewIds" binding:"required,min=1"`
}

type BatchRejectPayload struct {
	ReviewIDs []uint64 `json:"reviewIds" binding:"required,min=1"`
	Reason    string   `json:"reason" binding:"required,max=500"`
}

// UpdateReplyPayload 更新回复请求
type UpdateReplyPayload struct {
	Content string `json:"content" binding:"required,max=500"`
}

// UpdateReply
// @Summary      更新评价回复
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                  true  "回复ID"
// @Param        request  body  UpdateReplyPayload   true  "回复内容"
// @Success      200  {object}  model.APIResponse[map[string]interface{}]
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/review-replies/{id} [put]
func (h *ReviewHandler) UpdateReply(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidID)
		return
	}

	var p UpdateReplyPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		writeJSONError(c, 400, err.Error())
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		writeJSONError(c, 401, "未授权")
		return
	}

	result, err := h.svc.UpdateReviewReply(c.Request.Context(), userID.(uint64), id, p.Content)
	if errors.Is(err, adminservice.ErrNotFound) {
		_ = c.Error(adminservice.ErrNotFound)
		return
	}
	if errors.Is(err, adminservice.ErrUnauthorized) {
		writeJSONError(c, 403, "无权操作")
		return
	}
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}

	writeJSON(c, 200, model.APIResponse[map[string]interface{}]{
		Success: true,
		Code:    200,
		Message: "回复更新成功",
		Data:    result,
	})
}

// DeleteReply
// @Summary      删除评价回复
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "回复ID"
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/review-replies/{id} [delete]
func (h *ReviewHandler) DeleteReply(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidID)
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		writeJSONError(c, 401, "未授权")
		return
	}

	err = h.svc.DeleteReviewReply(c.Request.Context(), userID.(uint64), id)
	if errors.Is(err, adminservice.ErrNotFound) {
		_ = c.Error(adminservice.ErrNotFound)
		return
	}
	if errors.Is(err, adminservice.ErrUnauthorized) {
		writeJSONError(c, 403, "无权操作")
		return
	}
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}

	c.JSON(200, model.SuccessResponse{Success: true, Code: 200, Message: "回复删除成功"})
}

// SearchOperationLogs
// @Summary      搜索操作日志
// @Tags         Admin/OperationLogs
// @Security     BearerAuth
// @Produce      json
// @Param        page          query  int     false  "页码"
// @Param        pageSize      query  int     false  "每页数量"
// @Param        entity_type   query  string  false  "实体类型" Enums(review,order,payment,user,player)
// @Param        entity_id     query  int     false  "实体ID"
// @Param        action        query  string  false  "动作过滤"
// @Param        actor_user_id query  int     false  "操作者用户ID"
// @Param        date_from     query  string  false  "开始日期 (YYYY-MM-DD)"
// @Param        date_to       query  string  false  "结束日期 (YYYY-MM-DD)"
// @Param        export        query  string  false  "导出格式" Enums(csv)
// @Success      200  {object}  model.APIResponse[[]model.OperationLog]
// @Router       /admin/operation-logs [get]
func (h *ReviewHandler) SearchOperationLogs(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	var entityID, actorID *uint64
	if v, err := queryUint64Ptr(c, "entity_id"); err == nil {
		entityID = v
	}
	if v, err := queryUint64Ptr(c, "actor_user_id"); err == nil {
		actorID = v
	}

	var dateFrom, dateTo *time.Time
	if v, err := queryTimePtr(c, "date_from"); err == nil {
		dateFrom = v
	} else if strings.TrimSpace(c.Query("date_from")) != "" {
		writeJSONError(c, 400, apierr.ErrInvalidDateFrom)
		return
	}
	if v, err := queryTimePtr(c, "date_to"); err == nil {
		dateTo = v
	} else if strings.TrimSpace(c.Query("date_to")) != "" {
		writeJSONError(c, 400, apierr.ErrInvalidDateTo)
		return
	}

	opts := repository.OperationLogSearchOptions{
		Page:        page,
		PageSize:    pageSize,
		EntityType:  strings.TrimSpace(c.Query("entity_type")),
		EntityID:    entityID,
		Action:      strings.TrimSpace(c.Query("action")),
		ActorUserID: actorID,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
	}

	items, p, err := h.svc.SearchOperationLogs(c.Request.Context(), opts)
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}

	if strings.EqualFold(strings.TrimSpace(c.Query("export")), "csv") {
		exportOperationLogsCSV(c, opts.EntityType, 0, items)
		return
	}

	items = ensureSlice(items)
	writeJSON(c, 200, model.APIResponse[[]model.OperationLog]{
		Success:    true,
		Code:       200,
		Message:    "OK",
		Data:       items,
		Pagination: p,
	})
}

// ExportOperationLogs
// @Summary      导出操作日志
// @Tags         Admin/OperationLogs
// @Security     BearerAuth
// @Produce      text/csv
// @Param        entity_type   query  string  false  "实体类型" Enums(review,order,payment,user,player)
// @Param        entity_id     query  int     false  "实体ID"
// @Param        action        query  string  false  "动作过滤"
// @Param        actor_user_id query  int     false  "操作者用户ID"
// @Param        date_from     query  string  false  "开始日期 (YYYY-MM-DD)"
// @Param        date_to       query  string  false  "结束日期 (YYYY-MM-DD)"
// @Success      200  {file}  file
// @Router       /admin/operation-logs/export [get]
func (h *ReviewHandler) ExportOperationLogs(c *gin.Context) {
	var entityID, actorID *uint64
	if v, err := queryUint64Ptr(c, "entity_id"); err == nil {
		entityID = v
	}
	if v, err := queryUint64Ptr(c, "actor_user_id"); err == nil {
		actorID = v
	}

	var dateFrom, dateTo *time.Time
	if v, err := queryTimePtr(c, "date_from"); err == nil {
		dateFrom = v
	}
	if v, err := queryTimePtr(c, "date_to"); err == nil {
		dateTo = v
	}

	opts := repository.OperationLogSearchOptions{
		Page:        1,
		PageSize:    10000, // Export all matching records
		EntityType:  strings.TrimSpace(c.Query("entity_type")),
		EntityID:    entityID,
		Action:      strings.TrimSpace(c.Query("action")),
		ActorUserID: actorID,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
	}

	items, _, err := h.svc.SearchOperationLogs(c.Request.Context(), opts)
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}

	exportOperationLogsCSV(c, opts.EntityType, 0, items)
}
