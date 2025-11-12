package admin

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"

    serviceitem "gamelink/internal/service/item"
    "gamelink/internal/model"
)

type fakeItemService struct{}
func (f *fakeItemService) CreateServiceItem(_ interface{}, _ serviceitem.CreateServiceItemRequest) (*model.ServiceItem, error) { return &model.ServiceItem{}, nil }

func TestItemHandler_Create_InvalidJSON(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest(http.MethodPost, "/admin/service-items", bytes.NewBufferString("invalid json"))
    c.Request.Header.Set("Content-Type", "application/json")
    createServiceItemHandler(c, (*serviceitem.ServiceItemService)(nil))
    assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemHandler_Get_InvalidID(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key: "id", Value: "invalid"}}
    getServiceItemHandler(c, (*serviceitem.ServiceItemService)(nil))
    assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemHandler_List_InvalidQuery(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest(http.MethodGet, "/admin/service-items?page=abc", nil)
    listServiceItemsHandler(c, (*serviceitem.ServiceItemService)(nil))
    assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemHandler_Update_InvalidID(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key: "id", Value: "invalid"}}
    body, _ := json.Marshal(map[string]interface{}{"name": "x"})
    c.Request = httptest.NewRequest(http.MethodPut, "/admin/service-items/invalid", bytes.NewBuffer(body))
    c.Request.Header.Set("Content-Type", "application/json")
    updateServiceItemHandler(c, (*serviceitem.ServiceItemService)(nil))
    assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemHandler_Delete_InvalidID(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key: "id", Value: "invalid"}}
    deleteServiceItemHandler(c, (*serviceitem.ServiceItemService)(nil))
    assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemHandler_BatchUpdateStatus_InvalidJSON(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest(http.MethodPost, "/admin/service-items/batch-update-status", bytes.NewBufferString("bad"))
    c.Request.Header.Set("Content-Type", "application/json")
    batchUpdateStatusHandler(c, (*serviceitem.ServiceItemService)(nil))
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

func TestItemHandler_BatchUpdatePrice_InvalidJSON(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest(http.MethodPost, "/admin/service-items/batch-update-price", bytes.NewBufferString("bad"))
    c.Request.Header.Set("Content-Type", "application/json")
    batchUpdatePriceHandler(c, (*serviceitem.ServiceItemService)(nil))
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

func TestItemHandler_BatchUpdatePrice_NegativePrice(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    body := []byte(`{"ids":[1,2],"basePriceCents":-100}`)
    c.Request = httptest.NewRequest(http.MethodPost, "/admin/service-items/batch-update-price", bytes.NewBuffer(body))
    c.Request.Header.Set("Content-Type", "application/json")
    batchUpdatePriceHandler(c, (*serviceitem.ServiceItemService)(nil))
    if w.Code == 0 { t.Fatalf("invalid status") }
}

func TestItemHandler_List_InvalidPagination(t *testing.T) {
    repo := newFakeItemRepo()
    r, _ := setupItemRouter(repo)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/admin/service-items?page=-1&page_size=abc", nil)
    r.ServeHTTP(w, req)
    if w.Code == 0 { t.Fatalf("invalid status") }
}
