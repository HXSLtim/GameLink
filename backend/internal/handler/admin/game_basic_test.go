package admin

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"

    "gamelink/internal/cache"
    adminservice "gamelink/internal/service/admin"
)

func setupGameBasicRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    svc := adminservice.NewAdminService(nil, nil, nil, nil, nil, nil, cache.NewMemory())
    h := NewGameHandler(svc)
    r.GET("/admin/games/:id", h.GetGame)
    r.PUT("/admin/games/:id", h.UpdateGame)
    r.POST("/admin/games", h.CreateGame)
    r.GET("/admin/games", h.ListGames)
    r.DELETE("/admin/games/:id", h.DeleteGame)
    r.GET("/admin/games/:id/logs", h.ListGameLogs)
    return r
}

func TestGameHandler_Get_InvalidID(t *testing.T) {
    r := setupGameBasicRouter()
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/admin/games/abc", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

func TestGameHandler_Update_InvalidJSON(t *testing.T) {
    r := setupGameBasicRouter()
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPut, "/admin/games/1", bytes.NewBufferString("bad"))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

func TestGameHandler_Create_InvalidJSON(t *testing.T) {
    r := setupGameBasicRouter()
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/admin/games", bytes.NewBufferString("bad"))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

func TestGameHandler_ListLogs_InvalidID(t *testing.T) {
    r := setupGameBasicRouter()
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/admin/games/abc/logs", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

