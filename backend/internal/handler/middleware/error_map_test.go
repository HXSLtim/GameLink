package middleware

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"

    apierr "gamelink/internal/handler"
    "gamelink/internal/model"
    "gamelink/internal/service"
)

func TestErrorMap_ValidationMapsTo400(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.Use(ErrorMap())
    r.GET("/test", func(c *gin.Context) {
        _ = c.Error(service.ErrValidation)
    })

    w := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/test", nil)
    r.ServeHTTP(w, req)

    if w.Code != http.StatusBadRequest {
        t.Fatalf("expected 400, got %d", w.Code)
    }
}

func TestErrorMap_NotFoundMapsTo404(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.Use(ErrorMap())
    r.GET("/test", func(c *gin.Context) {
        _ = c.Error(service.ErrNotFound)
    })

    w := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/test", nil)
    r.ServeHTTP(w, req)

    if w.Code != http.StatusNotFound {
        t.Fatalf("expected 404, got %d", w.Code)
    }
}

// ensure envelope shape compatible with APIResponse
func expectEnvelope(t *testing.T, body []byte) {
    var _ model.APIResponse[any]
}

func TestErrorMap_NotFound_UsersMessage(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.Use(ErrorMap())
    r.GET("/api/v1/admin/users/:id", func(c *gin.Context) {
        _ = c.Error(service.ErrNotFound)
    })

    w := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/api/v1/admin/users/1", nil)
    r.ServeHTTP(w, req)

    if w.Code != http.StatusNotFound {
        t.Fatalf("expected 404, got %d", w.Code)
    }
    if got := w.Body.String(); !containsAll(got, apierr.ErrUserNotFound) {
        t.Fatalf("expected body to contain %q, got: %s", apierr.ErrUserNotFound, got)
    }
}

func containsAll(s string, subs ...string) bool {
    for _, sub := range subs {
        if !contains(s, sub) { return false }
    }
    return true
}
func contains(s, sub string) bool { return len(sub) == 0 || (len(s) >= len(sub) && (stringIndex(s, sub) >= 0)) }
func stringIndex(s, sub string) int { return findIndex(s, sub) }
func findIndex(s, sub string) int {
    // simple fallback since we avoid extra imports
    return indexOf(s, sub)
}
func indexOf(s, sub string) int {
    // naive search
    n, m := len(s), len(sub)
    if m == 0 { return 0 }
    for i := 0; i+m <= n; i++ {
        if s[i:i+m] == sub { return i }
    }
    return -1
}
