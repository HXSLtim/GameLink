package middleware

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/gin-gonic/gin"

    "gamelink/internal/auth"
)

func TestRequireAuth_MissingHeader(t *testing.T) {
    gin.SetMode(gin.TestMode)
    jwt := auth.NewJWTManager("s", 1)
    m := NewPermissionMiddleware(jwt, nil, nil)
    r := gin.New()
    r.GET("/p", m.RequireAuth(), func(c *gin.Context){ c.JSON(200, gin.H{"ok":true}) })
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/p", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusUnauthorized { t.Fatalf("%d", w.Code) }
}

func TestRequireAuth_ValidToken(t *testing.T) {
    gin.SetMode(gin.TestMode)
    jwt := auth.NewJWTManager("s", 60*time.Second)
    tok, _ := jwt.GenerateToken(42, "admin")
    m := NewPermissionMiddleware(jwt, nil, nil)
    r := gin.New()
    r.GET("/p", m.RequireAuth(), func(c *gin.Context){ c.JSON(200, gin.H{"ok":true}) })
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/p", nil)
    req.Header.Set("Authorization", "Bearer "+tok)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK { t.Fatalf("%d", w.Code) }
}
