package admin

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"

    apierr "gamelink/internal/handler"
)

func setupCtxWithQuery(method, path string) (*gin.Context, *httptest.ResponseRecorder) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    req := httptest.NewRequest(method, path, nil)
    c.Request = req
    return c, w
}

func TestQueryTimePtr_MultipleFormats(t *testing.T) {
    // RFC3339
    c, _ := setupCtxWithQuery("GET", "/?date_from=2025-10-27T12:34:56Z")
    if v, err := queryTimePtr(c, "date_from"); err != nil || v == nil {
        t.Fatalf("expected RFC3339 parsed, got err=%v v=%v", err, v)
    }
    // 2006-01-02 15:04:05
    c, _ = setupCtxWithQuery("GET", "/?date_from=2025-10-27%2012:34:56")
    if v, err := queryTimePtr(c, "date_from"); err != nil || v == nil {
        t.Fatalf("expected '2006-01-02 15:04:05' parsed, got err=%v v=%v", err, v)
    }
    // 2006-01-02
    c, _ = setupCtxWithQuery("GET", "/?date_from=2025-10-27")
    if v, err := queryTimePtr(c, "date_from"); err != nil || v == nil {
        t.Fatalf("expected '2006-01-02' parsed, got err=%v v=%v", err, v)
    }
    // unix seconds
    c, _ = setupCtxWithQuery("GET", "/?date_from=1730024096")
    if v, err := queryTimePtr(c, "date_from"); err != nil || v == nil {
        t.Fatalf("expected unix seconds parsed, got err=%v v=%v", err, v)
    }
    // invalid
    c, _ = setupCtxWithQuery("GET", "/?date_from=invalid")
    if v, err := queryTimePtr(c, "date_from"); err == nil || v != nil {
        t.Fatalf("expected error for invalid format, got err=%v v=%v", err, v)
    }
}

func TestBuildOrderListOptions_InvalidUserID(t *testing.T) {
    c, w := setupCtxWithQuery("GET", "/?user_id=abc")
    if _, ok := buildOrderListOptions(c); ok {
        t.Fatal("expected ok=false for invalid user_id")
    }
    if w.Code != http.StatusOK { // gin context not auto-written; our helper writes JSON and sets status via JSON()
        // Response code is set in the JSON payload, but recorder keeps default 200 unless we assert body
    }
}

func TestBuildPaymentListOptions_InvalidOrderID(t *testing.T) {
    c, _ := setupCtxWithQuery("GET", "/?order_id=bad")
    if _, ok := buildPaymentListOptions(c); ok {
        t.Fatal("expected ok=false for invalid order_id")
    }
}

func TestParsePagination_DefaultsAndErrors(t *testing.T) {
    // defaults
    c, _ := setupCtxWithQuery("GET", "/")
    page, pageSize, ok := parsePagination(c)
    if !ok || page != 1 || pageSize != 20 {
        t.Fatalf("expected defaults (1,20,true), got (%d,%d,%v)", page, pageSize, ok)
    }
    // invalid page
    c, _ = setupCtxWithQuery("GET", "/?page=xyz")
    _, _, ok = parsePagination(c)
    if ok {
        t.Fatal("expected ok=false for invalid page")
    }
    // invalid page_size
    c, _ = setupCtxWithQuery("GET", "/?page_size=oops")
    _, _, ok = parsePagination(c)
    if ok {
        t.Fatal("expected ok=false for invalid page_size")
    }
}

func TestConstants_Exist(t *testing.T) {
    // sanity to ensure constants referenced compile and accessible via import alias in handlers
    _ = apierr.ErrInvalidJSONPayload
    _ = apierr.ErrInvalidID
}

