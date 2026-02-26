package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mw "gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service/sensitiveword"
)

type sensitiveWordHandlerTestRepo struct {
	getFn func(ctx context.Context, id uint64) (*model.SensitiveWord, error)
}

func (r *sensitiveWordHandlerTestRepo) Create(_ context.Context, _ *model.SensitiveWord) error {
	return nil
}

func (r *sensitiveWordHandlerTestRepo) Get(ctx context.Context, id uint64) (*model.SensitiveWord, error) {
	if r.getFn != nil {
		return r.getFn(ctx, id)
	}
	return nil, repository.ErrNotFound
}

func (r *sensitiveWordHandlerTestRepo) List(_ context.Context, _ repository.SensitiveWordListOptions) ([]model.SensitiveWord, int64, error) {
	return nil, 0, nil
}

func (r *sensitiveWordHandlerTestRepo) Update(_ context.Context, _ *model.SensitiveWord) error {
	return nil
}

func (r *sensitiveWordHandlerTestRepo) Delete(_ context.Context, _ uint64) error {
	return nil
}

func (r *sensitiveWordHandlerTestRepo) GetAll(_ context.Context) ([]model.SensitiveWord, error) {
	return nil, nil
}

func setupSensitiveWordHandlerTestRouter(repo repository.SensitiveWordRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	handler := NewSensitiveWordHandler(sensitiveword.NewSensitiveWordService(repo))
	pm := mw.NewPermissionMiddleware(nil, nil, nil)
	pm.AddToWhitelist("GET:/api/v1/admin/sensitive-words/:id")

	RegisterSensitiveWordRoutes(router, handler, pm)
	return router
}

func TestGetSensitiveWord_Success(t *testing.T) {
	repo := &sensitiveWordHandlerTestRepo{
		getFn: func(_ context.Context, id uint64) (*model.SensitiveWord, error) {
			require.Equal(t, uint64(7), id)
			return &model.SensitiveWord{
				Base: model.Base{
					ID:        7,
					CreatedAt: time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC),
					UpdatedAt: time.Date(2026, time.January, 6, 7, 8, 9, 0, time.UTC),
				},
				Word:     "badword",
				Category: model.SensitiveWordCategoryAbuse,
				Severity: model.SensitiveWordSeverityHigh,
			}, nil
		},
	}
	router := setupSensitiveWordHandlerTestRouter(repo)

	req, err := http.NewRequest(http.MethodGet, "/sensitive-words/7", nil)
	require.NoError(t, err)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, true, resp["success"])

	data, ok := resp["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(7), data["id"])
	assert.Equal(t, "badword", data["word"])
	assert.Equal(t, "abuse", data["category"])
	assert.Equal(t, "high", data["severity"])
	assert.Equal(t, "2026-01-02 03:04:05", data["createdAt"])
	assert.Equal(t, "2026-01-06 07:08:09", data["updatedAt"])
}

func TestGetSensitiveWord_NotFound(t *testing.T) {
	repo := &sensitiveWordHandlerTestRepo{
		getFn: func(_ context.Context, _ uint64) (*model.SensitiveWord, error) {
			return nil, repository.ErrNotFound
		},
	}
	router := setupSensitiveWordHandlerTestRouter(repo)

	req, err := http.NewRequest(http.MethodGet, "/sensitive-words/999", nil)
	require.NoError(t, err)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, false, resp["success"])
	assert.Equal(t, "sensitive word not found", resp["message"])
}
