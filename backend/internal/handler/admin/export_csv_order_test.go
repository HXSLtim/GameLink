package admin

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "time"
    "github.com/gin-gonic/gin"
    "gamelink/internal/model"
)

func TestExportOperationLogsCSV_Order_EN_NoBOM_TZ(t *testing.T) {
    r := newTestEngine()
    r.GET("/export_order", func(c *gin.Context) {
        items := []model.OperationLog{
            {Base: model.Base{ID: 9, CreatedAt: time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)}, EntityType: "order", EntityID: 9, Action: "update_status", Reason: "", MetadataJSON: []byte("{\"k\":\"v\"}")},
        }
        exportOperationLogsCSV(c, "order", 9, items)
    })
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/export_order?fields=id,metadata,created_at&header_lang=en&tz=Asia/Shanghai&bom=false", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w.Code) }
    if ct := w.Header().Get("Content-Type"); ct == "" || !strings.Contains(ct, "text/csv") { t.Fatalf("expected csv content type, got %q", ct) }
    if cd := w.Header().Get("Content-Disposition"); !strings.Contains(cd, "attachment") || !strings.Contains(cd, "order_9_logs.csv") { t.Fatalf("unexpected content disposition: %q", cd) }
    body := w.Body.String()
    if strings.HasPrefix(body, "\uFEFF") { t.Fatalf("unexpected BOM present") }
    if !strings.Contains(body, "created_at") || !strings.Contains(body, "metadata") { t.Fatalf("expected headers present") }
    _ = body
}
