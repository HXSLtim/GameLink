package resp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
)

func TestGetUserID_Found(t *testing.T) {
	c, _ := setupTestContext()
	c.Set("user_id", uint64(123))

	userID := GetUserID(c)

	assert.Equal(t, uint64(123), userID)
}

func TestGetUserID_NotFound(t *testing.T) {
	c, _ := setupTestContext()

	userID := GetUserID(c)

	assert.Equal(t, uint64(0), userID)
}

func TestGetUserID_WrongType(t *testing.T) {
	c, _ := setupTestContext()
	c.Set("user_id", "not-a-uint64")

	userID := GetUserID(c)

	assert.Equal(t, uint64(0), userID)
}

func TestGetUserIDOrFail_Success(t *testing.T) {
	c, _ := setupTestContext()
	c.Set("user_id", uint64(456))

	userID, ok := GetUserIDOrFail(c)

	assert.True(t, ok)
	assert.Equal(t, uint64(456), userID)
}

func TestGetUserIDOrFail_NotFound(t *testing.T) {
	c, w := setupTestContext()

	userID, ok := GetUserIDOrFail(c)

	assert.False(t, ok)
	assert.Equal(t, uint64(0), userID)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestParseUintParam_Valid(t *testing.T) {
	c, _ := setupTestContext()
	c.Params = gin.Params{{Key: "id", Value: "123"}}

	id, err := ParseUintParam(c, "id")

	assert.NoError(t, err)
	assert.Equal(t, uint64(123), id)
}

func TestParseUintParam_Invalid(t *testing.T) {
	c, _ := setupTestContext()
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	_, err := ParseUintParam(c, "id")

	assert.Error(t, err)
}

func TestParseIDOrFail_Success(t *testing.T) {
	c, _ := setupTestContext()
	c.Params = gin.Params{{Key: "id", Value: "789"}}

	id, ok := ParseIDOrFail(c, "id")

	assert.True(t, ok)
	assert.Equal(t, uint64(789), id)
}

func TestParseIDOrFail_Empty(t *testing.T) {
	c, w := setupTestContext()
	c.Params = gin.Params{{Key: "id", Value: ""}}

	id, ok := ParseIDOrFail(c, "id")

	assert.False(t, ok)
	assert.Equal(t, uint64(0), id)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp model.APIResponse[any]
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp.Message, "缺少id参数")
}

func TestParseIDOrFail_Invalid(t *testing.T) {
	c, w := setupTestContext()
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	id, ok := ParseIDOrFail(c, "id")

	assert.False(t, ok)
	assert.Equal(t, uint64(0), id)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp model.APIResponse[any]
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp.Message, "无效的id格式")
}

func TestParseIDOrFail_Zero(t *testing.T) {
	c, w := setupTestContext()
	c.Params = gin.Params{{Key: "id", Value: "0"}}

	id, ok := ParseIDOrFail(c, "id")

	assert.False(t, ok)
	assert.Equal(t, uint64(0), id)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp model.APIResponse[any]
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp.Message, "id不能为0")
}

func TestBindJSONOrFail_Success(t *testing.T) {
	c, _ := setupTestContext()
	c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"test"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	type Req struct {
		Name string `json:"name"`
	}
	var req Req

	ok := BindJSONOrFail(c, &req)

	assert.True(t, ok)
	assert.Equal(t, "test", req.Name)
}

func TestBindJSONOrFail_InvalidJSON(t *testing.T) {
	c, w := setupTestContext()
	c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{invalid}`))
	c.Request.Header.Set("Content-Type", "application/json")

	type Req struct {
		Name string `json:"name"`
	}
	var req Req

	ok := BindJSONOrFail(c, &req)

	assert.False(t, ok)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBindQueryOrFail_Success(t *testing.T) {
	c, _ := setupTestContext()
	c.Request = httptest.NewRequest(http.MethodGet, "/?page=1&size=10", nil)

	type Query struct {
		Page int `form:"page"`
		Size int `form:"size"`
	}
	var q Query

	ok := BindQueryOrFail(c, &q)

	assert.True(t, ok)
	assert.Equal(t, 1, q.Page)
	assert.Equal(t, 10, q.Size)
}

func TestBindQueryOrFail_Invalid(t *testing.T) {
	c, w := setupTestContext()
	c.Request = httptest.NewRequest(http.MethodGet, "/?page=abc", nil)

	type Query struct {
		Page int `form:"page" binding:"required"`
	}
	var q Query

	ok := BindQueryOrFail(c, &q)

	assert.False(t, ok)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
