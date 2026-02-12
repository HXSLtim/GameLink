package resp

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service"
	"gamelink/pkg/apierr"
)

func TestError_APIError(t *testing.T) {
	c, w := setupTestContext()

	apiErr := apierr.BadRequest("invalid input").WithRequestID("req-123")
	Error(c, apiErr)

	var resp model.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, "invalid input", resp.Message)
	assert.Equal(t, "req-123", resp.TraceID)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestError_RepositoryNotFound(t *testing.T) {
	c, w := setupTestContext()

	Error(c, repository.ErrNotFound)

	var resp model.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Equal(t, "resource not found", resp.Message)
}

func TestError_ValidationError(t *testing.T) {
	c, w := setupTestContext()

	validationErr := apierr.NewValidationError("email", "invalid email format")
	Error(c, validationErr)

	var resp model.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestError_GenericError(t *testing.T) {
	c, w := setupTestContext()

	Error(c, errors.New("something went wrong"))

	var resp model.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Equal(t, "internal server error", resp.Message)
}

func TestErrorMsg(t *testing.T) {
	c, w := setupTestContext()

	ErrorMsg(c, http.StatusForbidden, "access denied")

	var resp model.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, http.StatusForbidden, resp.Code)
	assert.Equal(t, "access denied", resp.Message)
}

func TestBadRequest(t *testing.T) {
	c, w := setupTestContext()

	BadRequest(c, "invalid parameter")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUnauthorized(t *testing.T) {
	c, w := setupTestContext()

	Unauthorized(c, "token expired")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestForbidden(t *testing.T) {
	c, w := setupTestContext()

	Forbidden(c, "no permission")

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestNotFound(t *testing.T) {
	c, w := setupTestContext()

	NotFound(c, "user not found")

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestInternalError(t *testing.T) {
	c, w := setupTestContext()

	InternalError(c, "database error")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMapServiceError_Nil(t *testing.T) {
	code, resp := MapServiceError(nil)

	assert.Equal(t, http.StatusOK, code)
	assert.True(t, resp.Success)
	assert.Equal(t, "success", resp.Message)
}

func TestMapServiceError_APIError(t *testing.T) {
	apiErr := apierr.BadRequest("bad request").
		WithDetails("detail info").
		WithField("email").
		WithExtension("custom", "value")

	code, resp := MapServiceError(apiErr)

	assert.Equal(t, http.StatusBadRequest, code)
	assert.False(t, resp.Success)
	assert.Equal(t, "bad request", resp.Message)
	assert.NotNil(t, resp.Meta)
	meta := resp.Meta.(map[string]any)
	assert.Equal(t, "detail info", meta["details"])
	assert.Equal(t, "email", meta["field"])
	assert.Equal(t, "value", meta["custom"])
}

func TestMapServiceError_InvalidCredentials(t *testing.T) {
	code, resp := MapServiceError(service.ErrInvalidCredentials)

	assert.Equal(t, http.StatusUnauthorized, code)
	assert.False(t, resp.Success)
	assert.Equal(t, "用户名或密码错误", resp.Message)
}

func TestMapServiceError_UserDisabled(t *testing.T) {
	code, resp := MapServiceError(service.ErrUserDisabled)

	assert.Equal(t, http.StatusForbidden, code)
	assert.False(t, resp.Success)
	assert.Equal(t, "用户账户已禁用", resp.Message)
}

func TestMapServiceError_OrderInvalidTransition(t *testing.T) {
	code, resp := MapServiceError(service.ErrOrderInvalidTransition)

	assert.Equal(t, http.StatusBadRequest, code) // ErrOrderInvalidTransition is BadRequest
	assert.False(t, resp.Success)
}

func TestMapServiceError_UserNotFound(t *testing.T) {
	code, resp := MapServiceError(service.ErrUserNotFound)

	assert.Equal(t, http.StatusNotFound, code)
	assert.False(t, resp.Success)
}

func TestMapServiceError_NotFound(t *testing.T) {
	code, resp := MapServiceError(service.ErrNotFound)

	assert.Equal(t, http.StatusNotFound, code)
	assert.False(t, resp.Success)
	assert.Equal(t, "resource not found", resp.Message)
}

func TestMapServiceError_RepositoryNotFound(t *testing.T) {
	code, resp := MapServiceError(repository.ErrNotFound)

	assert.Equal(t, http.StatusNotFound, code)
	assert.False(t, resp.Success)
}

func TestMapServiceError_GenericError(t *testing.T) {
	code, resp := MapServiceError(errors.New("unknown error"))

	assert.Equal(t, http.StatusInternalServerError, code)
	assert.False(t, resp.Success)
	assert.Equal(t, "internal server error", resp.Message)
}

func TestMapServiceErrorWithPath_Users(t *testing.T) {
	code, resp := MapServiceErrorWithPath(service.ErrNotFound, "/api/v1/users/123")

	assert.Equal(t, http.StatusNotFound, code)
	assert.Equal(t, apierr.ErrUserNotFound, resp.Message)
}

func TestMapServiceErrorWithPath_Orders(t *testing.T) {
	code, resp := MapServiceErrorWithPath(service.ErrNotFound, "/api/v1/orders/456")

	assert.Equal(t, http.StatusNotFound, code)
	assert.Equal(t, apierr.ErrOrderNotFound, resp.Message)
}

func TestMapServiceErrorWithPath_Payments(t *testing.T) {
	code, resp := MapServiceErrorWithPath(service.ErrNotFound, "/api/v1/payments/789")

	assert.Equal(t, http.StatusNotFound, code)
	assert.Equal(t, apierr.ErrPaymentNotFound, resp.Message)
}

func TestMapServiceErrorWithPath_Players(t *testing.T) {
	code, resp := MapServiceErrorWithPath(service.ErrNotFound, "/api/v1/players/1")

	assert.Equal(t, http.StatusNotFound, code)
	assert.Equal(t, apierr.ErrPlayerNotFound, resp.Message)
}

func TestMapServiceErrorWithPath_Games(t *testing.T) {
	code, resp := MapServiceErrorWithPath(service.ErrNotFound, "/api/v1/games/2")

	assert.Equal(t, http.StatusNotFound, code)
	assert.Equal(t, apierr.ErrGameNotFound, resp.Message)
}

func TestMapServiceErrorWithPath_Unknown(t *testing.T) {
	code, resp := MapServiceErrorWithPath(service.ErrNotFound, "/api/v1/unknown/3")

	assert.Equal(t, http.StatusNotFound, code)
	assert.Equal(t, "resource not found", resp.Message)
}

func TestGetDomainNotFoundMessage(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/api/v1/users/1", apierr.ErrUserNotFound},
		{"/api/v1/orders/1", apierr.ErrOrderNotFound},
		{"/api/v1/payments/1", apierr.ErrPaymentNotFound},
		{"/api/v1/players/1", apierr.ErrPlayerNotFound},
		{"/api/v1/games/1", apierr.ErrGameNotFound},
		{"/api/v1/other/1", "resource not found"},
		{"/USERS/1", apierr.ErrUserNotFound}, // case insensitive
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := getDomainNotFoundMessage(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}
