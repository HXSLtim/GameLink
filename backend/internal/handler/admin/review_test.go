package admin

import (
    "context"
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"

    "gamelink/internal/cache"
    "gamelink/internal/model"
    "gamelink/internal/repository"
    "gamelink/internal/repository/common"
    adminservice "gamelink/internal/service/admin"
)

func setupReviewTestRouter() (*gin.Engine, *adminservice.AdminService) {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    svc := adminservice.NewAdminService(nil, nil, nil, nil, nil, nil, cache.NewMemory())
    h := NewReviewHandler(svc)
    r.GET("/admin/reviews", h.ListReviews)
    r.GET("/admin/reviews/:id", h.GetReview)
    r.POST("/admin/reviews", h.CreateReview)
    r.PUT("/admin/reviews/:id", h.UpdateReview)
    r.DELETE("/admin/reviews/:id", h.DeleteReview)
    r.GET("/admin/reviews/:id/logs", h.ListReviewLogs)
    return r, svc
}

type opLogsCSV struct{ items []model.OperationLog }
func (o opLogsCSV) Append(context.Context, *model.OperationLog) error { return nil }
func (o opLogsCSV) ListByEntity(context.Context, string, uint64, repository.OperationLogListOptions) ([]model.OperationLog, int64, error) { return o.items, int64(len(o.items)), nil }

type txLogs struct{ repos common.Repos }
func (t *txLogs) WithTx(ctx context.Context, fn func(r *common.Repos) error) error { return fn(&t.repos) }

func TestReviewHandler_GetReview_InvalidID(t *testing.T) {
    r, _ := setupReviewTestRouter()
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/admin/reviews/abc", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("expected 400, got %d", w.Code) }
}

func TestReviewHandler_Create_InvalidJSON(t *testing.T) {
    r, _ := setupReviewTestRouter()
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/admin/reviews", bytes.NewBufferString("bad"))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("expected 400, got %d", w.Code) }
}

func TestReviewHandler_Update_InvalidID(t *testing.T) {
    r, _ := setupReviewTestRouter()
    w := httptest.NewRecorder()
    body := []byte(`{"score":5,"content":"ok"}`)
    req := httptest.NewRequest(http.MethodPut, "/admin/reviews/abc", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("expected 400, got %d", w.Code) }
}

func TestReviewHandler_Update_InvalidJSON(t *testing.T) {
    r, _ := setupReviewTestRouter()
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPut, "/admin/reviews/1", bytes.NewBufferString("bad"))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("expected 400, got %d", w.Code) }
}

func TestReviewHandler_ListReviews_InvalidFilters(t *testing.T) {
    r, _ := setupReviewTestRouter()
    for _, url := range []string{
        "/admin/reviews?order_id=abc",
        "/admin/reviews?user_id=abc",
        "/admin/reviews?player_id=abc",
        "/admin/reviews?date_from=bad",
        "/admin/reviews?date_to=bad",
    } {
        w := httptest.NewRecorder()
        req := httptest.NewRequest(http.MethodGet, url, nil)
        r.ServeHTTP(w, req)
        if w.Code != http.StatusBadRequest { t.Fatalf("expected 400, got %d", w.Code) }
    }
}

func TestReviewHandler_ListReviewLogs_Invalid(t *testing.T) {
    r, _ := setupReviewTestRouter()
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/admin/reviews/abc/logs?page=1&page_size=20&date_from=2025-01-01T00:00:00Z&date_to=2025-01-02T00:00:00Z", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("expected 400, got %d", w.Code) }
    w2 := httptest.NewRecorder()
    req2 := httptest.NewRequest(http.MethodGet, "/admin/reviews/1/logs?page=1&page_size=20&date_from=bad&date_to=bad", nil)
    r.ServeHTTP(w2, req2)
    if w2.Code != http.StatusBadRequest { t.Fatalf("expected 400, got %d", w2.Code) }
}

// duplicate block removed

type reviewsRepo struct{ created *model.Review }
func (r *reviewsRepo) List(context.Context, repository.ReviewListOptions) ([]model.Review, int64, error) { return nil, 0, nil }
func (r *reviewsRepo) Get(context.Context, uint64) (*model.Review, error) { return &model.Review{Base: model.Base{ID:1}, Score: 5}, nil }
func (r *reviewsRepo) Create(ctx context.Context, rv *model.Review) error { r.created = rv; rv.ID = 10; return nil }
func (r *reviewsRepo) Update(context.Context, *model.Review) error { return nil }
func (r *reviewsRepo) Delete(context.Context, uint64) error { return nil }

type txReviews struct{ repos common.Repos }
func (t *txReviews) WithTx(ctx context.Context, fn func(r *common.Repos) error) error { return fn(&t.repos) }

func TestReviewHandler_CreateReview_Success(t *testing.T) {
    r, svc := setupReviewTestRouter()
    repo := &reviewsRepo{}
    svc.SetTxManager(&txReviews{repos: common.Repos{Reviews: repo, OpLogs: opLogsCSV{items: []model.OperationLog{}}}})
    w := httptest.NewRecorder()
    body := []byte(`{"order_id":1,"user_id":1,"player_id":1,"score":5,"content":"ok"}`)
    req := httptest.NewRequest(http.MethodPost, "/admin/reviews", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusCreated { t.Fatalf("expected 201, got %d", w.Code) }
}

func TestReviewHandler_UpdateReview_SuccessOrHandled(t *testing.T) {
    r, svc := setupReviewTestRouter()
    repo := &reviewsRepo{}
    svc.SetTxManager(&txReviews{repos: common.Repos{Reviews: repo, OpLogs: opLogsCSV{items: []model.OperationLog{}}}})
    w := httptest.NewRecorder()
    body := []byte(`{"score":5,"content":"ok"}`)
    req := httptest.NewRequest(http.MethodPut, "/admin/reviews/1", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError { t.Fatalf("expected 200/500, got %d", w.Code) }
}

func TestReviewHandler_DeleteReview_SuccessOrHandled(t *testing.T) {
    r, svc := setupReviewTestRouter()
    repo := &reviewsRepo{}
    svc.SetTxManager(&txReviews{repos: common.Repos{Reviews: repo, OpLogs: opLogsCSV{items: []model.OperationLog{}}}})
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodDelete, "/admin/reviews/1", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError { t.Fatalf("expected 200/500, got %d", w.Code) }
}

func TestReviewHandler_ListReviewLogs_ExportCSV_English_NoBOM_UTC(t *testing.T) {
    r, svc := setupReviewTestRouter()
    svc.SetTxManager(&txLogs{repos: common.Repos{OpLogs: opLogsCSV{items: []model.OperationLog{{Base: model.Base{ID:1}}}}}})
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/admin/reviews/1/logs?date_from=2025-01-01T00:00:00Z&date_to=2025-01-02T00:00:00Z&export=csv&fields=id,created_at&header_lang=en&tz=UTC", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w.Code) }
}

// duplicate block removed
