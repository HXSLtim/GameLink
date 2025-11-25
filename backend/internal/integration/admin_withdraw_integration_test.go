package integration

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	adminhandler "gamelink/internal/handler/admin"
	"gamelink/internal/model"
	withdrawrepo "gamelink/internal/repository/withdraw"
	"gamelink/internal/testutil"
)

// 管理端提现审批流：批准->完成、拒绝，状态与处理人应正确更新
func TestAdminWithdrawApprovalFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateWithdrawModels(t, db)

	repo := withdrawrepo.NewWithdrawRepository(db)
	router := gin.New()
	admin := router.Group("/api/v1/admin")
	admin.Use(setUserID(999)) // 设置 admin user_id
	adminhandler.RegisterWithdrawRoutes(admin, repo)

	// seed pending withdraw
	pending := &model.Withdraw{
		PlayerID:    1,
		UserID:      1,
		AmountCents: 10000,
		Method:      model.WithdrawMethodAlipay,
		AccountInfo: "ali-acc",
		Status:      model.WithdrawStatusPending,
	}
	if err := repo.Create(ctx(), pending); err != nil {
		t.Fatalf("seed withdraw: %v", err)
	}

	// approve
	approvePayload := map[string]interface{}{"remark": "ok"}
	approveResp := doJSON(router, http.MethodPost, "/api/v1/admin/withdraws/"+uintToStr(pending.ID)+"/approve", approvePayload, "")
	if approveResp.Code != http.StatusOK {
		t.Fatalf("approve status=%d body=%s", approveResp.Code, approveResp.Body.String())
	}
	approved, _ := repo.Get(ctx(), pending.ID)
	if approved.Status != model.WithdrawStatusApproved || approved.ProcessedBy == nil || approved.ProcessedAt == nil {
		t.Fatalf("unexpected approved state: %+v", approved)
	}

	// complete
	completeResp := doJSON(router, http.MethodPost, "/api/v1/admin/withdraws/"+uintToStr(pending.ID)+"/complete", nil, "")
	if completeResp.Code != http.StatusOK {
		t.Fatalf("complete status=%d body=%s", completeResp.Code, completeResp.Body.String())
	}
	completed, _ := repo.Get(ctx(), pending.ID)
	if completed.Status != model.WithdrawStatusCompleted || completed.CompletedAt == nil {
		t.Fatalf("unexpected completed state: %+v", completed)
	}

	// seed another pending and reject
	rejectW := &model.Withdraw{
		PlayerID:    2,
		UserID:      2,
		AmountCents: 20000,
		Method:      model.WithdrawMethodWeChat,
		AccountInfo: "wx-acc",
		Status:      model.WithdrawStatusPending,
		CreatedAt:   time.Now(),
	}
	_ = repo.Create(ctx(), rejectW)

	rejectPayload := map[string]interface{}{"reason": "risk"}
	rejectResp := doJSON(router, http.MethodPost, "/api/v1/admin/withdraws/"+uintToStr(rejectW.ID)+"/reject", rejectPayload, "")
	if rejectResp.Code != http.StatusOK {
		t.Fatalf("reject status=%d body=%s", rejectResp.Code, rejectResp.Body.String())
	}
	rejected, _ := repo.Get(ctx(), rejectW.ID)
	if rejected.Status != model.WithdrawStatusRejected || rejected.RejectReason != "risk" {
		t.Fatalf("unexpected rejected state: %+v", rejected)
	}

	// ensure list filtered by status returns only expected counts
	listResp := doJSON(router, http.MethodGet, "/api/v1/admin/withdraws?status=completed", nil, "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("list status=%d body=%s", listResp.Code, listResp.Body.String())
	}
	var listParsed apiResp[struct {
		Withdraws []model.Withdraw `json:"withdraws"`
		Total     int64            `json:"total"`
	}]
	_ = json.Unmarshal(listResp.Body.Bytes(), &listParsed)
	if listParsed.Data.Total != 1 || len(listParsed.Data.Withdraws) != 1 || listParsed.Data.Withdraws[0].ID != pending.ID {
		t.Fatalf("unexpected list result: %+v", listParsed.Data)
	}
}

// 非法状态流转：重复审批/完成、错误状态的完成或拒绝
func TestAdminWithdrawInvalidTransitions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateWithdrawModels(t, db)

	repo := withdrawrepo.NewWithdrawRepository(db)
	router := gin.New()
	admin := router.Group("/api/v1/admin")
	admin.Use(setUserID(999))
	adminhandler.RegisterWithdrawRoutes(admin, repo)

	// seed approved withdraw
	approved := &model.Withdraw{
		PlayerID:    3,
		UserID:      3,
		AmountCents: 30000,
		Method:      model.WithdrawMethodAlipay,
		AccountInfo: "ali-acc",
		Status:      model.WithdrawStatusApproved,
	}
	_ = repo.Create(ctx(), approved)

	// 重复 approve -> 400
	dupApprove := doJSON(router, http.MethodPost, "/api/v1/admin/withdraws/"+uintToStr(approved.ID)+"/approve", nil, "")
	if dupApprove.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 on duplicate approve, got %d body=%s", dupApprove.Code, dupApprove.Body.String())
	}

	// reject 已批准 -> 400
	rejectApproved := doJSON(router, http.MethodPost, "/api/v1/admin/withdraws/"+uintToStr(approved.ID)+"/reject", map[string]interface{}{"reason": "no"}, "")
	if rejectApproved.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 on reject approved, got %d body=%s", rejectApproved.Code, rejectApproved.Body.String())
	}

	// seed pending withdraw
	pending := &model.Withdraw{
		PlayerID:    4,
		UserID:      4,
		AmountCents: 15000,
		Method:      model.WithdrawMethodWeChat,
		AccountInfo: "wx-acc",
		Status:      model.WithdrawStatusPending,
	}
	_ = repo.Create(ctx(), pending)

	// complete pending -> 400
	completePending := doJSON(router, http.MethodPost, "/api/v1/admin/withdraws/"+uintToStr(pending.ID)+"/complete", nil, "")
	if completePending.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 on complete pending, got %d body=%s", completePending.Code, completePending.Body.String())
	}
}
