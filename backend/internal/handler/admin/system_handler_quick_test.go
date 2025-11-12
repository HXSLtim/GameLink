package admin

import (
    "context"
    "database/sql"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/gin-gonic/gin"

    "gamelink/internal/config"
)

type fakeCache struct{}
func (fakeCache) Get(ctx context.Context, key string) (string, bool, error) { return "", false, nil }
func (fakeCache) Set(ctx context.Context, key, value string, ttl time.Duration) error { return nil }
func (fakeCache) Delete(ctx context.Context, key string) error { return nil }
func (fakeCache) Close(ctx context.Context) error { return nil }

func TestSystemHandler_Config_Cache_Resources_Version(t *testing.T) {
    gin.SetMode(gin.TestMode)
    cfg := config.Load()
    h := NewSystemInfoHandler(cfg, &sql.DB{}, fakeCache{})

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest(http.MethodGet, "/admin/system/config", nil)
    h.Config(c)
    if w.Code != http.StatusOK { t.Fatalf("%d", w.Code) }

    w2 := httptest.NewRecorder()
    c2, _ := gin.CreateTestContext(w2)
    c2.Request = httptest.NewRequest(http.MethodGet, "/admin/system/cache", nil)
    h.CacheStatus(c2)
    if w2.Code != http.StatusOK { t.Fatalf("%d", w2.Code) }

    w3 := httptest.NewRecorder()
    c3, _ := gin.CreateTestContext(w3)
    c3.Request = httptest.NewRequest(http.MethodGet, "/admin/system/resources", nil)
    h.Resources(c3)
    if w3.Code != http.StatusOK { t.Fatalf("%d", w3.Code) }

    w4 := httptest.NewRecorder()
    c4, _ := gin.CreateTestContext(w4)
    c4.Request = httptest.NewRequest(http.MethodGet, "/admin/system/version", nil)
    h.Version(c4)
    if w4.Code != http.StatusOK { t.Fatalf("%d", w4.Code) }
}
