package admin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gamelink/internal/repository"
	stats "gamelink/internal/service/stats"

	"github.com/gin-gonic/gin"
)

type fakeStatsRepo struct{}

func (fakeStatsRepo) Dashboard(_ context.Context) (repository.Dashboard, error) {
	return repository.Dashboard{TotalUsers: 1}, nil
}
func (fakeStatsRepo) RevenueTrend(_ context.Context, days int) ([]repository.DateValue, error) {
	return []repository.DateValue{{Date: "2025-01-01", Value: int64(days)}}, nil
}
func (fakeStatsRepo) UserGrowth(_ context.Context, days int) ([]repository.DateValue, error) {
	return []repository.DateValue{{Date: "2025-01-01", Value: int64(days)}}, nil
}
func (fakeStatsRepo) OrdersByStatus(_ context.Context) (map[string]int64, error) {
	return map[string]int64{"pending": 1}, nil
}
func (fakeStatsRepo) TopPlayers(_ context.Context, limit int) ([]repository.PlayerTop, error) {
	return []repository.PlayerTop{{PlayerID: 1, Nickname: "p"}}, nil
}
func (fakeStatsRepo) AuditOverview(_ context.Context, _ *time.Time, _ *time.Time) (map[string]int64, map[string]int64, error) {
	return map[string]int64{"order": 1}, map[string]int64{"create": 1}, nil
}
func (fakeStatsRepo) AuditTrend(_ context.Context, _ *time.Time, _ *time.Time, _ string, _ string) ([]repository.DateValue, error) {
	return []repository.DateValue{{Date: "2025-01-01", Value: 1}}, nil
}

type errorStatsRepo struct{}

func (errorStatsRepo) Dashboard(_ context.Context) (repository.Dashboard, error) {
	return repository.Dashboard{}, repository.ErrNotFound
}
func (errorStatsRepo) RevenueTrend(_ context.Context, days int) ([]repository.DateValue, error) {
	return nil, repository.ErrNotFound
}
func (errorStatsRepo) UserGrowth(_ context.Context, days int) ([]repository.DateValue, error) {
	return nil, repository.ErrNotFound
}
func (errorStatsRepo) OrdersByStatus(_ context.Context) (map[string]int64, error) {
	return nil, repository.ErrNotFound
}
func (errorStatsRepo) TopPlayers(_ context.Context, limit int) ([]repository.PlayerTop, error) {
	return nil, repository.ErrNotFound
}
func (errorStatsRepo) AuditOverview(_ context.Context, _ *time.Time, _ *time.Time) (map[string]int64, map[string]int64, error) {
	return nil, nil, repository.ErrNotFound
}
func (errorStatsRepo) AuditTrend(_ context.Context, _ *time.Time, _ *time.Time, _ string, _ string) ([]repository.DateValue, error) {
	return nil, repository.ErrNotFound
}

func TestStatsHandler_AllEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewStatsHandler(stats.NewStatsService(fakeStatsRepo{}))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/dashboard", nil)
	h.Dashboard(c)
	if w.Code == 0 {
		t.Fatalf("no response")
	}

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/revenue-trend?days=3", nil)
	h.RevenueTrend(c2)
	if w2.Code == 0 {
		t.Fatalf("no response")
	}

	w3 := httptest.NewRecorder()
	c3, _ := gin.CreateTestContext(w3)
	c3.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/user-growth?days=5", nil)
	h.UserGrowth(c3)
	if w3.Code == 0 {
		t.Fatalf("no response")
	}

	w4 := httptest.NewRecorder()
	c4, _ := gin.CreateTestContext(w4)
	c4.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/orders", nil)
	h.OrdersSummary(c4)
	if w4.Code == 0 {
		t.Fatalf("no response")
	}

	w5 := httptest.NewRecorder()
	c5, _ := gin.CreateTestContext(w5)
	c5.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/top-players?limit=2", nil)
	h.TopPlayers(c5)
	if w5.Code == 0 {
		t.Fatalf("no response")
	}

	w6 := httptest.NewRecorder()
	c6, _ := gin.CreateTestContext(w6)
	c6.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/audit/overview?from=2025-01-01T00:00:00Z&to=2025-01-02T00:00:00Z", nil)
	h.AuditOverview(c6)
	if w6.Code == 0 {
		t.Fatalf("no response")
	}

	w7 := httptest.NewRecorder()
	c7, _ := gin.CreateTestContext(w7)
	c7.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/audit/trend?from=2025-01-01T00:00:00Z&to=2025-01-02T00:00:00Z&entity=order&action=create", nil)
	h.AuditTrend(c7)
	if w7.Code == 0 {
		t.Fatalf("no response")
	}
}

func TestStatsHandler_ErrorEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewStatsHandler(stats.NewStatsService(errorStatsRepo{}))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/dashboard", nil)
	h.Dashboard(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d", w.Code)
	}

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/revenue-trend?days=3", nil)
	h.RevenueTrend(c2)
	if w2.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d", w2.Code)
	}

	w3 := httptest.NewRecorder()
	c3, _ := gin.CreateTestContext(w3)
	c3.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/user-growth?days=5", nil)
	h.UserGrowth(c3)
	if w3.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d", w3.Code)
	}

	w4 := httptest.NewRecorder()
	c4, _ := gin.CreateTestContext(w4)
	c4.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/orders", nil)
	h.OrdersSummary(c4)
	if w4.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d", w4.Code)
	}

	w5 := httptest.NewRecorder()
	c5, _ := gin.CreateTestContext(w5)
	c5.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/top-players?limit=2", nil)
	h.TopPlayers(c5)
	if w5.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d", w5.Code)
	}

	w6 := httptest.NewRecorder()
	c6, _ := gin.CreateTestContext(w6)
	c6.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/audit/overview?from=2025-01-01T00:00:00Z&to=2025-01-02T00:00:00Z", nil)
	h.AuditOverview(c6)
	if w6.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d", w6.Code)
	}

	w7 := httptest.NewRecorder()
	c7, _ := gin.CreateTestContext(w7)
	c7.Request = httptest.NewRequest(http.MethodGet, "/admin/stats/audit/trend?from=2025-01-01T00:00:00Z&to=2025-01-02T00:00:00Z&entity=order&action=create", nil)
	h.AuditTrend(c7)
	if w7.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d", w7.Code)
	}
}
