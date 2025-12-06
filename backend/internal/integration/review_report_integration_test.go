package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	adminhandler "gamelink/internal/handler/admin"
	"gamelink/internal/model"
	"gamelink/internal/repository/order"
	"gamelink/internal/repository/reviewreport"
	"gamelink/pkg/testutil"
)

// TestReviewReportCreation 测试举报创建
// 需求: 3.1
func TestReviewReportCreation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateReviewModerationModels(t, db)

	seed := seedReviewModerationData(t, db)
	router, _ := setupReviewModerationRouter(t, db, seed.adminUserID)

	ctx := context.Background()
	reviewRepo := order.NewReviewRepository(db)
	reportRepo := reviewreport.NewReviewReportRepository(db)

	// 1. 创建一个评�?
	testReview := &model.Review{
		OrderID:  seed.orderID,
		UserID:   seed.userID,
		PlayerID: seed.playerID,
		Score:    model.Rating(3),
		Content:  "Test review for reporting",
		Status:   model.ReviewStatusApproved,
	}
	err := reviewRepo.Create(ctx, testReview)
	require.NoError(t, err)

	// 2. 创建举报
	reportPayload := map[string]interface{}{
		"reason":   "评价内容不实",
		"evidence": "https://example.com/evidence.jpg",
	}
	reportResp := doJSON(router, http.MethodPost, "/api/v1/admin/reviews/"+uintToStr(testReview.ID)+"/reports", reportPayload, "")
	require.Equal(t, http.StatusCreated, reportResp.Code, "create report should succeed")

	var reportParsed apiResp[adminhandler.CreateReviewReportResponse]
	err = json.Unmarshal(reportResp.Body.Bytes(), &reportParsed)
	require.NoError(t, err)
	assert.True(t, reportParsed.Success)
	assert.Equal(t, "created", reportParsed.Message)
	assert.Greater(t, reportParsed.Data.ReportID, uint64(0))

	reportID := reportParsed.Data.ReportID

	// 3. 验证举报记录已创�?
	report, err := reportRepo.Get(ctx, reportID)
	require.NoError(t, err)
	assert.Equal(t, testReview.ID, report.ReviewID)
	assert.Equal(t, seed.adminUserID, report.ReporterID)
	assert.Equal(t, "评价内容不实", report.Reason)
	assert.Equal(t, "https://example.com/evidence.jpg", report.Evidence)
	assert.Equal(t, model.ReviewReportStatusPending, report.Status)
	assert.Nil(t, report.HandledBy)
	assert.Nil(t, report.HandledAt)

	// 4. 验证评价被标记为已举�?
	updatedReview, err := reviewRepo.Get(ctx, testReview.ID)
	require.NoError(t, err)
	assert.True(t, updatedReview.IsReported, "review should be marked as reported")
}

// TestReviewReportCreationInvalidReview 测试举报不存在的评价
func TestReviewReportCreationInvalidReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateReviewModerationModels(t, db)

	seed := seedReviewModerationData(t, db)
	router, _ := setupReviewModerationRouter(t, db, seed.adminUserID)

	// 尝试举报不存在的评价
	reportPayload := map[string]interface{}{
		"reason": "测试举报",
	}
	reportResp := doJSON(router, http.MethodPost, "/api/v1/admin/reviews/99999/reports", reportPayload, "")
	assert.NotEqual(t, http.StatusCreated, reportResp.Code, "should fail to report non-existent review")
}

// TestReviewReportHandlingFlow 测试举报处理流程
// 需�? 3.2, 3.3, 3.4, 3.5
func TestReviewReportHandlingFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateReviewModerationModels(t, db)

	seed := seedReviewModerationData(t, db)
	router, _ := setupReviewModerationRouter(t, db, seed.adminUserID)

	ctx := context.Background()
	reviewRepo := order.NewReviewRepository(db)
	reportRepo := reviewreport.NewReviewReportRepository(db)

	// 1. 创建三个评价用于不同的处理方式测�?
	reviews := make([]*model.Review, 3)
	for i := 0; i < 3; i++ {
		reviews[i] = &model.Review{
			OrderID:  seed.orderID,
			UserID:   seed.userID,
			PlayerID: seed.playerID,
			Score:    model.Rating(3),
			Content:  "Test review " + string(rune('A'+i)),
			Status:   model.ReviewStatusApproved,
		}
		err := reviewRepo.Create(ctx, reviews[i])
		require.NoError(t, err)
	}

	// 2. 为每个评价创建举�?
	reports := make([]*model.ReviewReport, 3)
	for i := 0; i < 3; i++ {
		reports[i] = &model.ReviewReport{
			ReviewID:   reviews[i].ID,
			ReporterID: seed.userID,
			Reason:     "举报原因 " + string(rune('A'+i)),
			Status:     model.ReviewReportStatusPending,
		}
		err := reportRepo.Create(ctx, reports[i])
		require.NoError(t, err)

		// 标记评价为已举报
		reviews[i].IsReported = true
		err = reviewRepo.Update(ctx, reviews[i])
		require.NoError(t, err)
	}

	// 3. 测试删除评价处理方式
	deletePayload := map[string]interface{}{
		"action": "delete",
		"note":   "评价内容违规，已删除",
	}
	deleteResp := doJSON(router, http.MethodPut, "/api/v1/admin/review-reports/"+uintToStr(reports[0].ID)+"/handle", deletePayload, "")
	require.Equal(t, http.StatusOK, deleteResp.Code, "handle report with delete should succeed")

	var deleteParsed apiResp[adminhandler.HandleReviewReportResponse]
	err := json.Unmarshal(deleteResp.Body.Bytes(), &deleteParsed)
	require.NoError(t, err)
	assert.True(t, deleteParsed.Success)
	assert.Equal(t, "评价已删除", deleteParsed.Data.Message)

	// 4. 验证评价已被删除
	deletedReview, err := reviewRepo.Get(ctx, reviews[0].ID)
	require.NoError(t, err)
	assert.Equal(t, model.ReviewStatusDeleted, deletedReview.Status, "review should be deleted")

	// 5. 验证举报状态已更新为已通过
	handledReport, err := reportRepo.Get(ctx, reports[0].ID)
	require.NoError(t, err)
	assert.Equal(t, model.ReviewReportStatusApproved, handledReport.Status)
	assert.NotNil(t, handledReport.HandledBy)
	assert.Equal(t, seed.adminUserID, *handledReport.HandledBy)
	assert.NotNil(t, handledReport.HandledAt)
	assert.Equal(t, "评价内容违规，已删除", handledReport.HandlingNote)

	// 6. 测试警告处理方式
	warnPayload := map[string]interface{}{
		"action": "warn",
		"note":   "已警告评价者，请注意言辞",
	}
	warnResp := doJSON(router, http.MethodPut, "/api/v1/admin/review-reports/"+uintToStr(reports[1].ID)+"/handle", warnPayload, "")
	require.Equal(t, http.StatusOK, warnResp.Code, "handle report with warn should succeed")

	var warnParsed apiResp[adminhandler.HandleReviewReportResponse]
	err = json.Unmarshal(warnResp.Body.Bytes(), &warnParsed)
	require.NoError(t, err)
	assert.True(t, warnParsed.Success)
	assert.Equal(t, "已警告评价者", warnParsed.Data.Message)

	// 7. 验证评价状态未改变（警告不删除评价）
	warnedReview, err := reviewRepo.Get(ctx, reviews[1].ID)
	require.NoError(t, err)
	assert.Equal(t, model.ReviewStatusApproved, warnedReview.Status, "review should still be approved after warning")

	// 8. 验证举报状态已更新为已通过
	warnedReport, err := reportRepo.Get(ctx, reports[1].ID)
	require.NoError(t, err)
	assert.Equal(t, model.ReviewReportStatusApproved, warnedReport.Status)
	assert.NotNil(t, warnedReport.HandledBy)
	assert.Equal(t, "已警告评价者，请注意言辞", warnedReport.HandlingNote)

	// 9. 测试驳回举报处理方式
	rejectPayload := map[string]interface{}{
		"action": "reject",
		"note":   "举报不成立，评价内容正常",
	}
	rejectResp := doJSON(router, http.MethodPut, "/api/v1/admin/review-reports/"+uintToStr(reports[2].ID)+"/handle", rejectPayload, "")
	require.Equal(t, http.StatusOK, rejectResp.Code, "handle report with reject should succeed")

	var rejectParsed apiResp[adminhandler.HandleReviewReportResponse]
	err = json.Unmarshal(rejectResp.Body.Bytes(), &rejectParsed)
	require.NoError(t, err)
	assert.True(t, rejectParsed.Success)
	assert.Equal(t, "举报已驳回", rejectParsed.Data.Message)

	// 10. 验证评价状态未改变
	rejectedReview, err := reviewRepo.Get(ctx, reviews[2].ID)
	require.NoError(t, err)
	assert.Equal(t, model.ReviewStatusApproved, rejectedReview.Status, "review should still be approved after rejection")

	// 11. 验证举报状态已更新为已驳回
	rejectedReport, err := reportRepo.Get(ctx, reports[2].ID)
	require.NoError(t, err)
	assert.Equal(t, model.ReviewReportStatusRejected, rejectedReport.Status)
	assert.NotNil(t, rejectedReport.HandledBy)
	assert.Equal(t, "举报不成立，评价内容正常", rejectedReport.HandlingNote)

	// 12. 验证驳回举报后，评价的举报标记被取消（因为没有其他待处理举报�?
	assert.False(t, rejectedReview.IsReported, "review should not be marked as reported after all reports are handled")
}

// TestReviewReportStatusUpdate 测试举报状态更�?
// 需�? 3.2, 3.3, 3.4, 3.5
func TestReviewReportStatusUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateReviewModerationModels(t, db)

	seed := seedReviewModerationData(t, db)
	router, _ := setupReviewModerationRouter(t, db, seed.adminUserID)

	ctx := context.Background()
	reviewRepo := order.NewReviewRepository(db)
	reportRepo := reviewreport.NewReviewReportRepository(db)

	// 1. 创建评价和举�?
	testReview := &model.Review{
		OrderID:    seed.orderID,
		UserID:     seed.userID,
		PlayerID:   seed.playerID,
		Score:      model.Rating(2),
		Content:    "Test review",
		Status:     model.ReviewStatusApproved,
		IsReported: true,
	}
	err := reviewRepo.Create(ctx, testReview)
	require.NoError(t, err)

	report := &model.ReviewReport{
		ReviewID:   testReview.ID,
		ReporterID: seed.userID,
		Reason:     "测试举报",
		Status:     model.ReviewReportStatusPending,
	}
	err = reportRepo.Create(ctx, report)
	require.NoError(t, err)

	// 2. 处理举报
	handlePayload := map[string]interface{}{
		"action": "delete",
		"note":   "测试处理",
	}
	handleResp := doJSON(router, http.MethodPut, "/api/v1/admin/review-reports/"+uintToStr(report.ID)+"/handle", handlePayload, "")
	require.Equal(t, http.StatusOK, handleResp.Code)

	// 3. 验证举报状态从 pending 变为 approved
	updatedReport, err := reportRepo.Get(ctx, report.ID)
	require.NoError(t, err)
	assert.Equal(t, model.ReviewReportStatusPending, report.Status, "original status should be pending")
	assert.Equal(t, model.ReviewReportStatusApproved, updatedReport.Status, "updated status should be approved")

	// 4. 验证处理时间和处理人已设�?
	assert.NotNil(t, updatedReport.HandledAt, "handled_at should be set")
	assert.NotNil(t, updatedReport.HandledBy, "handled_by should be set")
	assert.Equal(t, seed.adminUserID, *updatedReport.HandledBy)

	// 5. 验证不能重复处理已处理的举报
	reHandlePayload := map[string]interface{}{
		"action": "reject",
		"note":   "尝试重复处理",
	}
	reHandleResp := doJSON(router, http.MethodPut, "/api/v1/admin/review-reports/"+uintToStr(report.ID)+"/handle", reHandlePayload, "")
	assert.NotEqual(t, http.StatusOK, reHandleResp.Code, "should not be able to re-handle a handled report")
}

// TestListReviewReports 测试举报列表查询
// 需�? 3.1
func TestListReviewReports(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateReviewModerationModels(t, db)

	seed := seedReviewModerationData(t, db)
	router, _ := setupReviewModerationRouter(t, db, seed.adminUserID)

	ctx := context.Background()
	reviewRepo := order.NewReviewRepository(db)
	reportRepo := reviewreport.NewReviewReportRepository(db)

	// 1. 创建多个评价和举�?
	for i := 0; i < 5; i++ {
		testReview := &model.Review{
			OrderID:    seed.orderID,
			UserID:     seed.userID,
			PlayerID:   seed.playerID,
			Score:      model.Rating(3),
			Content:    "Test review " + string(rune('A'+i)),
			Status:     model.ReviewStatusApproved,
			IsReported: true,
		}
		err := reviewRepo.Create(ctx, testReview)
		require.NoError(t, err)

		status := model.ReviewReportStatusPending
		if i >= 3 {
			status = model.ReviewReportStatusApproved
		}

		report := &model.ReviewReport{
			ReviewID:   testReview.ID,
			ReporterID: seed.userID,
			Reason:     "举报原因 " + string(rune('A'+i)),
			Status:     status,
		}
		err = reportRepo.Create(ctx, report)
		require.NoError(t, err)
	}

	// 2. 获取所有举报列�?
	listResp := doJSON(router, http.MethodGet, "/api/v1/admin/review-reports?page=1&pageSize=20", nil, "")
	require.Equal(t, http.StatusOK, listResp.Code, "list reports should succeed")

	var listParsed apiResp[[]adminhandler.ReviewReportDTO]
	err := json.Unmarshal(listResp.Body.Bytes(), &listParsed)
	require.NoError(t, err)
	assert.True(t, listParsed.Success)
	assert.GreaterOrEqual(t, len(listParsed.Data), 5, "should have at least 5 reports")

	// 3. 按状态筛�?- 只获取待处理的举�?
	pendingResp := doJSON(router, http.MethodGet, "/api/v1/admin/review-reports?page=1&pageSize=20&status=pending", nil, "")
	require.Equal(t, http.StatusOK, pendingResp.Code)

	var pendingParsed apiResp[[]adminhandler.ReviewReportDTO]
	err = json.Unmarshal(pendingResp.Body.Bytes(), &pendingParsed)
	require.NoError(t, err)
	assert.True(t, pendingParsed.Success)
	assert.GreaterOrEqual(t, len(pendingParsed.Data), 3, "should have at least 3 pending reports")

	// 验证所有返回的举报都是待处理状�?
	for _, report := range pendingParsed.Data {
		assert.Equal(t, model.ReviewReportStatusPending, report.Status)
	}

	// 4. 按状态筛�?- 只获取已处理的举�?
	approvedResp := doJSON(router, http.MethodGet, "/api/v1/admin/review-reports?page=1&pageSize=20&status=approved", nil, "")
	require.Equal(t, http.StatusOK, approvedResp.Code)

	var approvedParsed apiResp[[]adminhandler.ReviewReportDTO]
	err = json.Unmarshal(approvedResp.Body.Bytes(), &approvedParsed)
	require.NoError(t, err)
	assert.True(t, approvedParsed.Success)
	assert.GreaterOrEqual(t, len(approvedParsed.Data), 2, "should have at least 2 approved reports")

	// 验证所有返回的举报都是已通过状�?
	for _, report := range approvedParsed.Data {
		assert.Equal(t, model.ReviewReportStatusApproved, report.Status)
	}
}

// TestReviewReportInvalidAction 测试无效的处理动�?
func TestReviewReportInvalidAction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateReviewModerationModels(t, db)

	seed := seedReviewModerationData(t, db)
	router, _ := setupReviewModerationRouter(t, db, seed.adminUserID)

	ctx := context.Background()
	reviewRepo := order.NewReviewRepository(db)
	reportRepo := reviewreport.NewReviewReportRepository(db)

	// 创建评价和举�?
	testReview := &model.Review{
		OrderID:  seed.orderID,
		UserID:   seed.userID,
		PlayerID: seed.playerID,
		Score:    model.Rating(3),
		Content:  "Test review",
		Status:   model.ReviewStatusApproved,
	}
	err := reviewRepo.Create(ctx, testReview)
	require.NoError(t, err)

	report := &model.ReviewReport{
		ReviewID:   testReview.ID,
		ReporterID: seed.userID,
		Reason:     "测试举报",
		Status:     model.ReviewReportStatusPending,
	}
	err = reportRepo.Create(ctx, report)
	require.NoError(t, err)

	// 尝试使用无效的处理动�?
	invalidPayload := map[string]interface{}{
		"action": "invalid_action",
		"note":   "测试",
	}
	invalidResp := doJSON(router, http.MethodPut, "/api/v1/admin/review-reports/"+uintToStr(report.ID)+"/handle", invalidPayload, "")
	assert.NotEqual(t, http.StatusOK, invalidResp.Code, "should fail with invalid action")

	// 验证举报状态未改变
	unchangedReport, err := reportRepo.Get(ctx, report.ID)
	require.NoError(t, err)
	assert.Equal(t, model.ReviewReportStatusPending, unchangedReport.Status, "report status should remain pending")
}

// TestMultipleReportsOnSameReview 测试同一评价的多个举�?
// 需�? 3.1, 3.5
func TestMultipleReportsOnSameReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateReviewModerationModels(t, db)

	seed := seedReviewModerationData(t, db)
	router, _ := setupReviewModerationRouter(t, db, seed.adminUserID)

	ctx := context.Background()
	reviewRepo := order.NewReviewRepository(db)
	reportRepo := reviewreport.NewReviewReportRepository(db)

	// 1. 创建一个评�?
	testReview := &model.Review{
		OrderID:  seed.orderID,
		UserID:   seed.userID,
		PlayerID: seed.playerID,
		Score:    model.Rating(2),
		Content:  "Controversial review",
		Status:   model.ReviewStatusApproved,
	}
	err := reviewRepo.Create(ctx, testReview)
	require.NoError(t, err)

	// 2. 创建多个举报
	reports := make([]*model.ReviewReport, 3)
	for i := 0; i < 3; i++ {
		reports[i] = &model.ReviewReport{
			ReviewID:   testReview.ID,
			ReporterID: seed.userID,
			Reason:     "举报原因 " + string(rune('A'+i)),
			Status:     model.ReviewReportStatusPending,
		}
		err = reportRepo.Create(ctx, reports[i])
		require.NoError(t, err)
	}

	// 标记评价为已举报
	testReview.IsReported = true
	err = reviewRepo.Update(ctx, testReview)
	require.NoError(t, err)

	// 3. 驳回第一个举报
	rejectPayload := map[string]interface{}{
		"action": "reject",
		"note":   "举报不成立",
	}
	rejectResp := doJSON(router, http.MethodPut, "/api/v1/admin/review-reports/"+uintToStr(reports[0].ID)+"/handle", rejectPayload, "")
	require.Equal(t, http.StatusOK, rejectResp.Code)

	// 4. 验证评价仍然标记为已举报（因为还有其他待处理举报�?
	reviewAfterFirstReject, err := reviewRepo.Get(ctx, testReview.ID)
	require.NoError(t, err)
	assert.True(t, reviewAfterFirstReject.IsReported, "review should still be marked as reported")

	// 5. 驳回第二个举�?
	rejectResp2 := doJSON(router, http.MethodPut, "/api/v1/admin/review-reports/"+uintToStr(reports[1].ID)+"/handle", rejectPayload, "")
	require.Equal(t, http.StatusOK, rejectResp2.Code)

	// 6. 验证评价仍然标记为已举报
	reviewAfterSecondReject, err := reviewRepo.Get(ctx, testReview.ID)
	require.NoError(t, err)
	assert.True(t, reviewAfterSecondReject.IsReported, "review should still be marked as reported")

	// 7. 驳回第三个举报（最后一个）
	rejectResp3 := doJSON(router, http.MethodPut, "/api/v1/admin/review-reports/"+uintToStr(reports[2].ID)+"/handle", rejectPayload, "")
	require.Equal(t, http.StatusOK, rejectResp3.Code)

	// 8. 验证评价的举报标记被取消（所有举报都已处理）
	reviewAfterAllRejected, err := reviewRepo.Get(ctx, testReview.ID)
	require.NoError(t, err)
	assert.False(t, reviewAfterAllRejected.IsReported, "review should not be marked as reported after all reports are rejected")
}

// TestReviewReportWithDeletedReview 测试处理已删除评价的举报
func TestReviewReportWithDeletedReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateReviewModerationModels(t, db)

	seed := seedReviewModerationData(t, db)
	router, _ := setupReviewModerationRouter(t, db, seed.adminUserID)

	ctx := context.Background()
	reviewRepo := order.NewReviewRepository(db)
	reportRepo := reviewreport.NewReviewReportRepository(db)

	// 1. 创建评价
	testReview := &model.Review{
		OrderID:  seed.orderID,
		UserID:   seed.userID,
		PlayerID: seed.playerID,
		Score:    model.Rating(3),
		Content:  "Test review",
		Status:   model.ReviewStatusApproved,
	}
	err := reviewRepo.Create(ctx, testReview)
	require.NoError(t, err)

	// 2. 创建举报
	report := &model.ReviewReport{
		ReviewID:   testReview.ID,
		ReporterID: seed.userID,
		Reason:     "测试举报",
		Status:     model.ReviewReportStatusPending,
	}
	err = reportRepo.Create(ctx, report)
	require.NoError(t, err)

	// 3. 先删除评�?
	testReview.Status = model.ReviewStatusDeleted
	err = reviewRepo.Update(ctx, testReview)
	require.NoError(t, err)

	// 4. 尝试处理举报（评价已删除，但举报仍可处理�?
	handlePayload := map[string]interface{}{
		"action": "reject",
		"note":   "评价已删除，驳回举报",
	}
	handleResp := doJSON(router, http.MethodPut, "/api/v1/admin/review-reports/"+uintToStr(report.ID)+"/handle", handlePayload, "")
	require.Equal(t, http.StatusOK, handleResp.Code, "should be able to handle report even if review is deleted")

	// 5. 验证举报已处�?
	handledReport, err := reportRepo.Get(ctx, report.ID)
	require.NoError(t, err)
	assert.Equal(t, model.ReviewReportStatusRejected, handledReport.Status)
}

// TestReviewReportPagination 测试举报列表分页
func TestReviewReportPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateReviewModerationModels(t, db)

	seed := seedReviewModerationData(t, db)
	router, _ := setupReviewModerationRouter(t, db, seed.adminUserID)

	ctx := context.Background()
	reviewRepo := order.NewReviewRepository(db)
	reportRepo := reviewreport.NewReviewReportRepository(db)

	// 创建25个举�?
	for i := 0; i < 25; i++ {
		testReview := &model.Review{
			OrderID:  seed.orderID,
			UserID:   seed.userID,
			PlayerID: seed.playerID,
			Score:    model.Rating(3),
			Content:  "Test review " + string(rune('A'+i%26)),
			Status:   model.ReviewStatusApproved,
		}
		err := reviewRepo.Create(ctx, testReview)
		require.NoError(t, err)

		report := &model.ReviewReport{
			ReviewID:   testReview.ID,
			ReporterID: seed.userID,
			Reason:     "举报原因",
			Status:     model.ReviewReportStatusPending,
		}
		err = reportRepo.Create(ctx, report)
		require.NoError(t, err)
	}

	// 获取第一页（20条）
	page1Resp := doJSON(router, http.MethodGet, "/api/v1/admin/review-reports?page=1&pageSize=20", nil, "")
	require.Equal(t, http.StatusOK, page1Resp.Code)

	var page1Parsed apiResp[[]adminhandler.ReviewReportDTO]
	err := json.Unmarshal(page1Resp.Body.Bytes(), &page1Parsed)
	require.NoError(t, err)
	assert.True(t, page1Parsed.Success)
	assert.GreaterOrEqual(t, len(page1Parsed.Data), 20, "first page should have at least 20 reports")

	// 获取第二页（剩余的）
	page2Resp := doJSON(router, http.MethodGet, "/api/v1/admin/review-reports?page=2&pageSize=20", nil, "")
	require.Equal(t, http.StatusOK, page2Resp.Code)

	var page2Parsed apiResp[[]adminhandler.ReviewReportDTO]
	err = json.Unmarshal(page2Resp.Body.Bytes(), &page2Parsed)
	require.NoError(t, err)
	assert.True(t, page2Parsed.Success)
	assert.GreaterOrEqual(t, len(page2Parsed.Data), 5, "second page should have at least 5 reports")
}
