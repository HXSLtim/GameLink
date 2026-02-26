package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	reconservice "gamelink/internal/service/reconciliation"
	"gamelink/pkg/auth"
)

type handlerFakeRepo struct {
	listFn    func(ctx context.Context, opts repository.ReconciliationListOptions) ([]model.Reconciliation, int64, error)
	getFn     func(ctx context.Context, id uint64, withDetails bool) (*model.Reconciliation, error)
	createFn  func(ctx context.Context, rec *model.Reconciliation) error
	executeFn func(ctx context.Context, id uint64, opts repository.ReconciliationExecuteOptions) (*model.Reconciliation, error)
}

func (f *handlerFakeRepo) List(ctx context.Context, opts repository.ReconciliationListOptions) ([]model.Reconciliation, int64, error) {
	if f.listFn != nil {
		return f.listFn(ctx, opts)
	}
	return nil, 0, nil
}

func (f *handlerFakeRepo) Get(ctx context.Context, id uint64, withDetails bool) (*model.Reconciliation, error) {
	if f.getFn != nil {
		return f.getFn(ctx, id, withDetails)
	}
	return nil, repository.ErrNotFound
}

func (f *handlerFakeRepo) Create(ctx context.Context, rec *model.Reconciliation) error {
	if f.createFn != nil {
		return f.createFn(ctx, rec)
	}
	return nil
}

func (f *handlerFakeRepo) Execute(ctx context.Context, id uint64, opts repository.ReconciliationExecuteOptions) (*model.Reconciliation, error) {
	if f.executeFn != nil {
		return f.executeFn(ctx, id, opts)
	}
	return nil, nil
}

func setupReconciliationRouter(t *testing.T, repo repository.ReconciliationRepository) (*gin.Engine, string) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	jwtMgr := auth.NewJWTManager("test-secret", time.Hour)
	token, err := jwtMgr.GenerateToken(99, "admin")
	require.NoError(t, err)

	pm := middleware.NewPermissionMiddleware(jwtMgr, nil, nil)
	pm.AddPathsToWhitelist([]string{
		"GET:/api/v1/admin/reconciliations",
		"GET:/api/v1/admin/reconciliations/:id",
		"POST:/api/v1/admin/reconciliations",
		"POST:/api/v1/admin/reconciliations/:id/execute",
	})

	router := gin.New()
	svc := reconservice.NewService(repo)
	RegisterReconciliationRoutes(router.Group("/api/v1/admin"), svc, pm)
	return router, token
}

func TestReconciliationHandler_List(t *testing.T) {
	now := time.Now().UTC()
	repo := &handlerFakeRepo{
		listFn: func(_ context.Context, opts repository.ReconciliationListOptions) ([]model.Reconciliation, int64, error) {
			require.Equal(t, model.ReconciliationStatusPending, *opts.Status)
			return []model.Reconciliation{
				{
					ReconciliationNo:   "RCN-H-001",
					ReconciliationDate: now,
					Type:               model.ReconciliationTypeManual,
					Status:             model.ReconciliationStatusPending,
					PeriodStart:        now.Add(-24 * time.Hour),
					PeriodEnd:          now,
				},
			}, 1, nil
		},
	}
	router, token := setupReconciliationRouter(t, repo)

	req, err := http.NewRequest(http.MethodGet, "/api/v1/admin/reconciliations?status=pending", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, true, body["success"])
	data, ok := body["data"].([]any)
	require.True(t, ok)
	require.Len(t, data, 1)
}

func TestReconciliationHandler_Create(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339)
	repo := &handlerFakeRepo{
		createFn: func(_ context.Context, rec *model.Reconciliation) error {
			rec.ID = 7
			return nil
		},
	}
	router, token := setupReconciliationRouter(t, repo)

	payload := map[string]any{
		"reconciliationDate": now,
		"type":               "manual",
		"periodStart":        now,
		"periodEnd":          now,
		"abstract":           "manual run",
		"details": []map[string]any{
			{
				"externalType":   "payment",
				"externalNo":     "EXT-1",
				"externalAmount": 1000,
				"externalDate":   now,
				"internalType":   "payment",
				"internalNo":     "INT-1",
				"internalAmount": 1000,
				"internalDate":   now,
			},
		},
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/admin/reconciliations", bytes.NewReader(raw))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, true, body["success"])
	data, ok := body["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(7), data["id"])
}

func TestReconciliationHandler_ExecuteUsesAdminUserFromToken(t *testing.T) {
	repo := &handlerFakeRepo{
		executeFn: func(_ context.Context, id uint64, opts repository.ReconciliationExecuteOptions) (*model.Reconciliation, error) {
			require.Equal(t, uint64(8), id)
			require.Equal(t, uint64(99), opts.ProcessedBy)
			return &model.Reconciliation{
				Base:   model.Base{ID: 8},
				Status: model.ReconciliationStatusSuccess,
			}, nil
		},
	}
	router, token := setupReconciliationRouter(t, repo)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/admin/reconciliations/8/execute", bytes.NewReader([]byte(`{}`)))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, true, body["success"])
	data, ok := body["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "success", data["status"])
}

func TestReconciliationHandler_ListRejectsInvalidPagination(t *testing.T) {
	router, token := setupReconciliationRouter(t, &handlerFakeRepo{})

	req, err := http.NewRequest(http.MethodGet, "/api/v1/admin/reconciliations?page=bad", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestReconciliationHandler_CreateRejectsInvalidPayload(t *testing.T) {
	router, token := setupReconciliationRouter(t, &handlerFakeRepo{})

	req, err := http.NewRequest(http.MethodPost, "/api/v1/admin/reconciliations", bytes.NewReader([]byte(`{}`)))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestReconciliationHandler_ExecuteRejectsInvalidStatus(t *testing.T) {
	repo := &handlerFakeRepo{
		executeFn: func(_ context.Context, _ uint64, _ repository.ReconciliationExecuteOptions) (*model.Reconciliation, error) {
			t.Fatalf("execute should not be called when status validation fails")
			return nil, nil
		},
	}
	router, token := setupReconciliationRouter(t, repo)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/admin/reconciliations/8/execute", bytes.NewReader([]byte(`{"status":"pending"}`)))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}
