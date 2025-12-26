package resp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestJSON_WithTraceID(t *testing.T) {
	c, w := setupTestContext()
	c.Set("request_id", "trace-123")

	payload := model.APIResponse[string]{
		Success: true,
		Code:    200,
		Message: "OK",
		Data:    "test",
	}

	JSON(c, http.StatusOK, payload)

	var resp model.APIResponse[string]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "trace-123", resp.TraceID)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJSON_WithExistingTraceID(t *testing.T) {
	c, w := setupTestContext()
	c.Set("request_id", "should-not-override")

	payload := model.APIResponse[string]{
		Success: true,
		Code:    200,
		Message: "OK",
		Data:    "test",
		TraceID: "existing-trace",
	}

	JSON(c, http.StatusOK, payload)

	var resp model.APIResponse[string]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "existing-trace", resp.TraceID)
}

func TestJSON_WithoutRequestID(t *testing.T) {
	c, w := setupTestContext()

	payload := model.APIResponse[string]{
		Success: true,
		Code:    200,
		Message: "OK",
		Data:    "test",
	}

	JSON(c, http.StatusOK, payload)

	var resp model.APIResponse[string]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Empty(t, resp.TraceID)
}

func TestSuccess(t *testing.T) {
	c, w := setupTestContext()

	type TestData struct {
		Name string `json:"name"`
	}
	data := TestData{Name: "test"}

	Success(c, "操作成功", data)

	var resp model.APIResponse[TestData]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "操作成功", resp.Message)
	assert.Equal(t, "test", resp.Data.Name)
}

func TestOK(t *testing.T) {
	c, w := setupTestContext()

	OK(c, map[string]int{"count": 10})

	var resp model.APIResponse[map[string]int]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "OK", resp.Message)
	assert.Equal(t, 10, resp.Data["count"])
}

func TestCreated(t *testing.T) {
	c, w := setupTestContext()

	type User struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
	}
	user := User{ID: 1, Name: "test"}

	Created(c, user)

	var resp model.APIResponse[User]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "created", resp.Message)
	assert.Equal(t, uint64(1), resp.Data.ID)
}

func TestUpdated(t *testing.T) {
	c, w := setupTestContext()

	Updated(c, "updated data")

	var resp model.APIResponse[string]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "updated", resp.Message)
	assert.Equal(t, "updated data", resp.Data)
}

func TestDeleted(t *testing.T) {
	c, w := setupTestContext()

	Deleted(c)

	var resp model.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "deleted", resp.Message)
}

func TestList_WithData(t *testing.T) {
	c, w := setupTestContext()

	items := []string{"a", "b", "c"}
	pagination := &model.Pagination{
		Page:       1,
		PageSize:   10,
		Total:      3,
		TotalPages: 1,
		HasNext:    false,
		HasPrev:    false,
	}

	List(c, items, pagination)

	var resp model.APIResponse[[]string]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "OK", resp.Message)
	assert.Len(t, resp.Data, 3)
	assert.NotNil(t, resp.Pagination)
	assert.Equal(t, 3, resp.Pagination.Total)
}

func TestList_WithNilData(t *testing.T) {
	c, w := setupTestContext()

	var items []string = nil

	List(c, items, nil)

	var resp model.APIResponse[[]string]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	// Data should be empty array, not nil (JSON serialization)
	assert.Len(t, resp.Data, 0)
}
