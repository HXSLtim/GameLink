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

func setupPlayerBasicRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    svc := adminservice.NewAdminService(nil, nil, nil, nil, nil, nil, cache.NewMemory())
    h := NewPlayerHandler(svc)
    r.PUT("/admin/players/:id/verification", h.UpdatePlayerVerification)
    r.PUT("/admin/players/:id/games", h.UpdatePlayerGames)
    r.PUT("/admin/players/:id/skill-tags", h.UpdatePlayerSkillTags)
    r.GET("/admin/players/:id/logs", h.ListPlayerLogs)
    return r
}

func TestPlayerHandler_UpdateVerification_InvalidID(t *testing.T) {
    r := setupPlayerBasicRouter()
    w := httptest.NewRecorder()
    body := []byte(`{"verification_status":"approved"}`)
    req := httptest.NewRequest(http.MethodPut, "/admin/players/abc/verification", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

func TestPlayerHandler_UpdateVerification_InvalidJSON(t *testing.T) {
    r := setupPlayerBasicRouter()
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPut, "/admin/players/1/verification", bytes.NewBufferString("bad"))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

func TestPlayerHandler_UpdateGames_InvalidID(t *testing.T) {
    r := setupPlayerBasicRouter()
    w := httptest.NewRecorder()
    body := []byte(`{"game_ids":[1,2]}`)
    req := httptest.NewRequest(http.MethodPut, "/admin/players/abc/games", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

func TestPlayerHandler_UpdateGames_InvalidJSON(t *testing.T) {
    r := setupPlayerBasicRouter()
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPut, "/admin/players/1/games", bytes.NewBufferString("bad"))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

func TestPlayerHandler_UpdateSkillTags_InvalidID(t *testing.T) {
    r := setupPlayerBasicRouter()
    w := httptest.NewRecorder()
    body := []byte(`{"tags":["a","b"]}`)
    req := httptest.NewRequest(http.MethodPut, "/admin/players/abc/skill-tags", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

func TestPlayerHandler_UpdateSkillTags_InvalidJSON(t *testing.T) {
    r := setupPlayerBasicRouter()
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPut, "/admin/players/1/skill-tags", bytes.NewBufferString("bad"))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

func TestPlayerHandler_ListLogs_InvalidDate(t *testing.T) {
    r := setupPlayerBasicRouter()
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/admin/players/1/logs?date_from=bad&date_to=bad", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

