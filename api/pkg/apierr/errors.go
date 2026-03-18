// Package apierr provides generic API error handling.

package apierr

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Common errors
var (
	ErrNotFound         = New(http.StatusNotFound, "资源不存在")
	ErrUnauthorized     = New(http.StatusUnauthorized, "未授权访问")
	ErrForbidden        = New(http.StatusForbidden, "权限不足")
	ErrInvalidInput     = New(http.StatusBadRequest, "输入参数无效")
	ErrInternal         = New(http.StatusInternalServerError, "服务器内部错误")
	ErrConflict         = New(http.StatusConflict, "资源冲突")
	ErrTooManyRequests  = New(http.StatusTooManyRequests, "请求过于频繁")
	ErrValidationFailed = New(http.StatusBadRequest, "验证失败")
)

// APIError represents a standardized API error response
type APIError struct {
	Code       int                    `json:"code"`
	Message    string                 `json:"message"`
	Details    string                 `json:"details,omitempty"`
	Field      string                 `json:"field,omitempty"`
	RequestID  string                 `json:"requestId,omitempty"`
	Timestamp  int64                  `json:"timestamp"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`

	// Internal fields (not exposed in JSON)
	internalError error `json:"-"` // Original error for logging
}

// Error implements the error interface
func (e APIError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// New creates a new API error
func New(code int, message string) *APIError {
	return &APIError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}
}

// WithDetails adds error details (sanitized for production)
func (e *APIError) WithDetails(details string) *APIError {
	// In production, don't expose internal details
	if isProduction() {
		// Log the details but don't expose them
		if e.internalError != nil {
			slog.Error("API error details",
				"code", e.Code,
				"message", e.Message,
				"details", details,
				"internal_error", e.internalError.Error())
		} else {
			slog.Error("API error details",
				"code", e.Code,
				"message", e.Message,
				"details", details)
		}
		// Don't set Details field in production
		return e
	}
	// In development, expose details for debugging
	e.Details = details
	return e
}

// WithInternalError stores the original error for logging (never exposed to client)
func (e *APIError) WithInternalError(err error) *APIError {
	e.internalError = err
	// Log immediately in production
	if isProduction() && err != nil {
		slog.Error("Internal error occurred",
			"code", e.Code,
			"message", e.Message,
			"error", err.Error())
	}
	return e
}

// WithField adds field information for validation errors
func (e *APIError) WithField(field string) *APIError {
	e.Field = field
	return e
}

// WithRequestID adds request ID
func (e *APIError) WithRequestID(requestID string) *APIError {
	e.RequestID = requestID
	return e
}

// WithExtension adds extension data
func (e *APIError) WithExtension(key string, value interface{}) *APIError {
	if e.Extensions == nil {
		e.Extensions = make(map[string]interface{})
	}
	e.Extensions[key] = value
	return e
}

// Sanitize returns a sanitized version of the error for production
func (e *APIError) Sanitize() *APIError {
	if !isProduction() {
		return e
	}

	// Create a clean copy without sensitive details
	sanitized := &APIError{
		Code:       e.Code,
		Message:    e.Message,
		Field:      e.Field,
		RequestID:  e.RequestID,
		Timestamp:  e.Timestamp,
		Extensions: e.Extensions,
		// Details and internalError are intentionally omitted
	}

	// Use generic messages for 5xx errors in production
	if e.Code >= 500 {
		sanitized.Message = "服务器内部错误，请稍后重试"
	}

	return sanitized
}

// isProduction checks if running in production environment
func isProduction() bool {
	env := os.Getenv("APP_ENV")
	return env == "production" || env == "staging"
}

// BadRequest creates a 400 Bad Request error
func BadRequest(message string) *APIError {
	return New(http.StatusBadRequest, message)
}

// Unauthorized creates a 401 Unauthorized error
func Unauthorized(message string) *APIError {
	return New(http.StatusUnauthorized, message)
}

// Forbidden creates a 403 Forbidden error
func Forbidden(message string) *APIError {
	return New(http.StatusForbidden, message)
}

// NotFound creates a 404 Not Found error
func NotFound(message string) *APIError {
	return New(http.StatusNotFound, message)
}

// Conflict creates a 409 Conflict error
func Conflict(message string) *APIError {
	return New(http.StatusConflict, message)
}

// InternalError creates a 500 Internal Server Error
func InternalError(message string) *APIError {
	return New(http.StatusInternalServerError, message)
}

// ServiceUnavailable creates a 503 Service Unavailable error
func ServiceUnavailable(message string) *APIError {
	return New(http.StatusServiceUnavailable, message)
}

// TooManyRequests creates a 429 Too Many Requests error
func TooManyRequests(message string) *APIError {
	return New(http.StatusTooManyRequests, message)
}

// ValidationError represents a validation error
type ValidationError struct {
	APIError
	Field string `json:"field"`
	Value string `json:"value,omitempty"`
	Tag   string `json:"tag,omitempty"`
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		APIError: APIError{
			Code:      http.StatusBadRequest,
			Message:   message,
			Timestamp: time.Now().Unix(),
		},
		Field: field,
	}
}

// DatabaseError represents a database error
type DatabaseError struct {
	APIError
	Query string `json:"query,omitempty"`
}

// NewDatabaseError creates a new database error
func NewDatabaseError(err error) *DatabaseError {
	return &DatabaseError{
		APIError: APIError{
			Code:      http.StatusInternalServerError,
			Message:   "数据库操作失败",
			Details:   err.Error(),
			Timestamp: time.Now().Unix(),
		},
	}
}

// IsNotFound checks if the error is a not found error
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	// Check if it's a direct APIError
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Code == http.StatusNotFound
	}
	// Check if it's a wrapped APIError using errors.Is
	var target *APIError
	if errors.As(err, &target) {
		return target.Code == http.StatusNotFound
	}
	return false
}

// IsUnauthorized checks if the error is an unauthorized error
func IsUnauthorized(err error) bool {
	if err == nil {
		return false
	}
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Code == http.StatusUnauthorized
	}
	var target *APIError
	if errors.As(err, &target) {
		return target.Code == http.StatusUnauthorized
	}
	return false
}

// IsForbidden checks if the error is a forbidden error
func IsForbidden(err error) bool {
	if err == nil {
		return false
	}
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Code == http.StatusForbidden
	}
	var target *APIError
	if errors.As(err, &target) {
		return target.Code == http.StatusForbidden
	}
	return false
}

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	if err == nil {
		return false
	}
	// 检查是否是 ValidationError 类型
	if _, ok := err.(*ValidationError); ok {
		return true
	}
	// 检查是否是 APIError 且状态码为 400
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Code == http.StatusBadRequest
	}
	// Check wrapped errors
	var target *ValidationError
	if errors.As(err, &target) {
		return true
	}
	var apiTarget *APIError
	if errors.As(err, &apiTarget) {
		return apiTarget.Code == http.StatusBadRequest
	}
	return false
}

// IsBadRequest checks if the error is a bad request error
func IsBadRequest(err error) bool {
	if err == nil {
		return false
	}
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Code == http.StatusBadRequest
	}
	var target *APIError
	if errors.As(err, &target) {
		return target.Code == http.StatusBadRequest
	}
	return false
}

// IsInternalError checks if the error is an internal server error
func IsInternalError(err error) bool {
	if err == nil {
		return false
	}
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Code == http.StatusInternalServerError
	}
	var target *APIError
	if errors.As(err, &target) {
		return target.Code == http.StatusInternalServerError
	}
	return false
}
