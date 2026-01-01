package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/service/content"
	"gamelink/pkg/apierr"
)

// ContentHandler 内容管理处理器
type ContentHandler struct {
	feedSvc   *content.AdminFeedService
	chatSvc   *content.ChatModerationService
	reportSvc *content.FeedReportService
	statsSvc  *content.ContentStatsService
}

// NewContentHandler 创建内容管理处理器
func NewContentHandler(
	feedSvc *content.AdminFeedService,
	chatSvc *content.ChatModerationService,
	reportSvc *content.FeedReportService,
	statsSvc *content.ContentStatsService,
) *ContentHandler {
	return &ContentHandler{
		feedSvc:   feedSvc,
		chatSvc:   chatSvc,
		reportSvc: reportSvc,
		statsSvc:  statsSvc,
	}
}

// ListFeeds 列出动态
// @Summary      列出动态
// @Description  获取动态列表，支持筛选
// @Tags         Admin - Content
// @Security     BearerAuth
// @Param        page             query     int     false  "页码"
// @Param        pageSize         query     int     false  "每页数量"
// @Param        authorId         query     int     false  "作者ID"
// @Param        categoryId       query     int     false  "分类ID"
// @Param        keyword          query     string  false  "关键词"
// @Param        moderationStatus query     string  false  "审核状态"
// @Success      200  {object}  model.APIResponse[content.AdminListFeedsResponse]
// @Router       /admin/content/feeds [get]
func (h *ContentHandler) ListFeeds(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	authorID, _ := queryUint64Ptr(c, "authorId")
	categoryID, _ := queryUint64Ptr(c, "categoryId")
	dateFrom, _ := queryTimePtr(c, "dateFrom")
	dateTo, _ := queryTimePtr(c, "dateTo")

	var moderationStatus *model.FeedModerationStatus
	if s := c.Query("moderationStatus"); s != "" {
		status := model.FeedModerationStatus(s)
		moderationStatus = &status
	}

	resp, err := h.feedSvc.ListFeeds(c.Request.Context(), content.AdminListFeedsRequest{
		Page:             page,
		PageSize:         pageSize,
		AuthorID:         authorID,
		CategoryID:       categoryID,
		Keyword:          c.Query("keyword"),
		ModerationStatus: moderationStatus,
		DateFrom:         dateFrom,
		DateTo:           dateTo,
	})
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondSuccess(c, resp)
}

// GetFeed 获取动态详情
// @Summary      获取动态详情
// @Tags         Admin - Content
// @Security     BearerAuth
// @Param        id   path      int  true  "动态ID"
// @Success      200  {object}  model.APIResponse[content.AdminFeedDTO]
// @Router       /admin/content/feeds/{id} [get]
func (h *ContentHandler) GetFeed(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	feed, err := h.feedSvc.GetFeed(c.Request.Context(), id)
	if err != nil {
		respondError(c, apierr.NotFound(err.Error()))
		return
	}

	respondSuccess(c, feed)
}

// ApproveFeed 批准动态
// @Summary      批准动态
// @Tags         Admin - Content
// @Security     BearerAuth
// @Param        id   path      int                        true  "动态ID"
// @Param        body body      content.AdminModerationRequest  false "审核备注"
// @Success      200  {object}  model.SuccessResponse
// @Router       /admin/content/feeds/{id}/approve [put]
func (h *ContentHandler) ApproveFeed(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req content.AdminModerationRequest
	_ = c.ShouldBindJSON(&req)

	if err := h.feedSvc.ApproveFeed(c.Request.Context(), id, adminID, req.Note); err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondMsg(c, "已批准")
}

// RejectFeed 拒绝动态
// @Summary      拒绝动态
// @Tags         Admin - Content
// @Security     BearerAuth
// @Param        id   path      int                        true  "动态ID"
// @Param        body body      content.AdminModerationRequest  true  "拒绝原因"
// @Success      200  {object}  model.SuccessResponse
// @Router       /admin/content/feeds/{id}/reject [put]
func (h *ContentHandler) RejectFeed(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req content.AdminModerationRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Note == "" {
		respondBadRequest(c, "拒绝原因不能为空")
		return
	}

	if err := h.feedSvc.RejectFeed(c.Request.Context(), id, adminID, req.Note); err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondMsg(c, "已拒绝")
}

// DeleteFeed 删除动态
// @Summary      删除动态
// @Tags         Admin - Content
// @Security     BearerAuth
// @Param        id   path      int  true  "动态ID"
// @Success      200  {object}  model.SuccessResponse
// @Router       /admin/content/feeds/{id} [delete]
func (h *ContentHandler) DeleteFeed(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req content.AdminModerationRequest
	_ = c.ShouldBindJSON(&req)

	if err := h.feedSvc.DeleteFeed(c.Request.Context(), id, adminID, req.Note); err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondDeleted(c)
}

// BatchApproveFeed 批量批准动态
// @Summary      批量批准动态
// @Tags         Admin - Content
// @Security     BearerAuth
// @Param        body body      content.AdminBatchModerationRequest  true  "批量审核请求"
// @Success      200  {object}  model.SuccessResponse
// @Router       /admin/content/feeds/batch-approve [post]
func (h *ContentHandler) BatchApproveFeed(c *gin.Context) {
	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req content.AdminBatchModerationRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	if err := h.feedSvc.BatchApproveFeed(c.Request.Context(), req.FeedIDs, adminID, req.Note); err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondMsg(c, "批量批准成功")
}

// BatchRejectFeed 批量拒绝动态
// @Summary      批量拒绝动态
// @Tags         Admin - Content
// @Security     BearerAuth
// @Param        body body      content.AdminBatchModerationRequest  true  "批量审核请求"
// @Success      200  {object}  model.SuccessResponse
// @Router       /admin/content/feeds/batch-reject [post]
func (h *ContentHandler) BatchRejectFeed(c *gin.Context) {
	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req content.AdminBatchModerationRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Note == "" {
		respondBadRequest(c, "拒绝原因不能为空")
		return
	}

	if err := h.feedSvc.BatchRejectFeed(c.Request.Context(), req.FeedIDs, adminID, req.Note); err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondMsg(c, "批量拒绝成功")
}

// GetContentStats 获取内容统计
// @Summary      获取内容统计
// @Tags         Admin - Content
// @Security     BearerAuth
// @Param        days query     int  false  "趋势天数"
// @Success      200  {object}  model.APIResponse[content.ContentStatsDTO]
// @Router       /admin/content/stats [get]
func (h *ContentHandler) GetContentStats(c *gin.Context) {
	days, _ := queryIntDefault(c, "days", 30)

	stats, err := h.statsSvc.GetStats(c.Request.Context(), days)
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondSuccess(c, stats)
}

// ExportContentStats 导出内容统计
// @Summary      导出内容统计Excel
// @Tags         Admin - Content
// @Security     BearerAuth
// @Param        days query     int  false  "趋势天数"
// @Produce      application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Success      200  {file}    binary
// @Router       /admin/content/stats/export [get]
func (h *ContentHandler) ExportContentStats(c *gin.Context) {
	days, _ := queryIntDefault(c, "days", 30)

	buf, filename, err := h.statsSvc.ExportStats(c.Request.Context(), days)
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Header("Content-Length", fmt.Sprintf("%d", buf.Len()))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

// ListChatMessages 列出聊天消息
// @Summary      列出聊天消息
// @Tags         Admin - Content
// @Security     BearerAuth
// @Param        page        query     int     false  "页码"
// @Param        pageSize    query     int     false  "每页数量"
// @Param        groupId     query     int     false  "群组ID"
// @Param        senderId    query     int     false  "发送者ID"
// @Param        auditStatus query     string  false  "审核状态"
// @Success      200  {object}  model.APIResponse[content.ListMessagesResponse]
// @Router       /admin/content/chat/messages [get]
func (h *ContentHandler) ListChatMessages(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	groupID, _ := queryUint64Ptr(c, "groupId")
	senderID, _ := queryUint64Ptr(c, "senderId")
	dateFrom, _ := queryTimePtr(c, "dateFrom")
	dateTo, _ := queryTimePtr(c, "dateTo")

	var auditStatus *model.ChatMessageAuditStatus
	if s := c.Query("auditStatus"); s != "" {
		status := model.ChatMessageAuditStatus(s)
		auditStatus = &status
	}

	resp, err := h.chatSvc.ListMessages(c.Request.Context(), content.ListMessagesRequest{
		Page:        page,
		PageSize:    pageSize,
		GroupID:     groupID,
		SenderID:    senderID,
		AuditStatus: auditStatus,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
	})
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondSuccess(c, resp)
}

// DeleteChatMessage 删除聊天消息
// @Summary      删除聊天消息
// @Tags         Admin - Content
// @Security     BearerAuth
// @Param        id   path      int  true  "消息ID"
// @Success      200  {object}  model.SuccessResponse
// @Router       /admin/content/chat/messages/{id} [delete]
func (h *ContentHandler) DeleteChatMessage(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&req)

	if err := h.chatSvc.DeleteMessage(c.Request.Context(), id, adminID, req.Reason); err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondDeleted(c)
}

// MuteUser 禁言用户
// @Summary      禁言用户
// @Tags         Admin - Content
// @Security     BearerAuth
// @Param        body body      content.MuteUserRequest  true  "禁言请求"
// @Success      200  {object}  model.SuccessResponse
// @Router       /admin/content/chat/mute [post]
func (h *ContentHandler) MuteUser(c *gin.Context) {
	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req content.MuteUserRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	if err := h.chatSvc.MuteUser(c.Request.Context(), req, adminID); err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondMsg(c, "已禁言")
}

// UnmuteUser 解除禁言
// @Summary      解除禁言
// @Tags         Admin - Content
// @Security     BearerAuth
// @Param        groupId query     int  true  "群组ID"
// @Param        userId  query     int  true  "用户ID"
// @Success      200  {object}  model.SuccessResponse
// @Router       /admin/content/chat/unmute [post]
func (h *ContentHandler) UnmuteUser(c *gin.Context) {
	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	groupID, ok := QueryUint64PtrAndRespond(c, "groupId", "groupId is required")
	if !ok || groupID == nil {
		if ok {
			respondBadRequest(c, "groupId is required")
		}
		return
	}

	userID, ok := QueryUint64PtrAndRespond(c, "userId", "userId is required")
	if !ok || userID == nil {
		if ok {
			respondBadRequest(c, "userId is required")
		}
		return
	}

	if err := h.chatSvc.UnmuteUser(c.Request.Context(), *groupID, *userID, adminID); err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondMsg(c, "已解除禁言")
}

// ListFeedReports 列出动态举报
// @Summary      列出动态举报
// @Tags         Admin - Content
// @Security     BearerAuth
// @Param        page       query     int     false  "页码"
// @Param        pageSize   query     int     false  "每页数量"
// @Param        feedId     query     int     false  "动态ID"
// @Param        status     query     string  false  "状态"
// @Success      200  {object}  model.APIResponse[content.ListFeedReportsResponse]
// @Router       /admin/content/reports [get]
func (h *ContentHandler) ListFeedReports(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	feedID, _ := queryUint64Ptr(c, "feedId")
	reporterID, _ := queryUint64Ptr(c, "reporterId")
	dateFrom, _ := queryTimePtr(c, "dateFrom")
	dateTo, _ := queryTimePtr(c, "dateTo")

	var status *string
	if s := c.Query("status"); s != "" {
		status = &s
	}

	resp, err := h.reportSvc.ListFeedReports(c.Request.Context(), content.ListFeedReportsRequest{
		Page:       page,
		PageSize:   pageSize,
		FeedID:     feedID,
		ReporterID: reporterID,
		Status:     status,
		DateFrom:   dateFrom,
		DateTo:     dateTo,
	})
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondSuccess(c, resp)
}

// GetFeedReport 获取举报详情
// @Summary      获取举报详情
// @Tags         Admin - Content
// @Security     BearerAuth
// @Param        id   path      int  true  "举报ID"
// @Success      200  {object}  model.APIResponse[content.FeedReportDTO]
// @Router       /admin/content/reports/{id} [get]
func (h *ContentHandler) GetFeedReport(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	report, err := h.reportSvc.GetFeedReport(c.Request.Context(), id)
	if err != nil {
		respondError(c, apierr.NotFound(err.Error()))
		return
	}

	respondSuccess(c, report)
}

// ProcessFeedReport 处理举报
// @Summary      处理举报
// @Tags         Admin - Content
// @Security     BearerAuth
// @Param        id   path      int                        true  "举报ID"
// @Param        body body      content.ProcessReportRequest  true  "处理请求"
// @Success      200  {object}  model.SuccessResponse
// @Router       /admin/content/reports/{id}/process [post]
func (h *ContentHandler) ProcessFeedReport(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var req content.ProcessReportRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	if err := h.reportSvc.ProcessFeedReport(c.Request.Context(), id, req, adminID); err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	respondMsg(c, "处理成功")
}
