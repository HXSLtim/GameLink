package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// WriteError 发送错误响应
func WriteError(c *gin.Context, code int, message string) {
	c.JSON(code, Error(code, message))
}

// WriteErrorWithData 发送带数据的错误响应
func WriteErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, ErrorWithData(code, message, data))
}

// WriteValidationError 发送验证错误响应
func WriteValidationError(c *gin.Context, message string) {
	WriteError(c, http.StatusBadRequest, message)
}

// WriteNotFoundError 发送未找到错误响应
func WriteNotFoundError(c *gin.Context, message string) {
	WriteError(c, http.StatusNotFound, message)
}

// WriteUnauthorizedError 发送未授权错误响应
func WriteUnauthorizedError(c *gin.Context, message string) {
	WriteError(c, http.StatusUnauthorized, message)
}

// WriteForbiddenError 发送禁止访问错误响应
func WriteForbiddenError(c *gin.Context, message string) {
	WriteError(c, http.StatusForbidden, message)
}

// WriteInternalServerError 发送服务器内部错误响应
func WriteInternalServerError(c *gin.Context, message string) {
	WriteError(c, http.StatusInternalServerError, message)
}

// WriteBadRequestError 发送错误请求响应
func WriteBadRequestError(c *gin.Context, message string) {
	WriteError(c, http.StatusBadRequest, message)
}

// WriteConflictError 发送冲突错误响应
func WriteConflictError(c *gin.Context, message string) {
	WriteError(c, http.StatusConflict, message)
}

// WriteTooManyRequestsError 发送请求过多错误响应
func WriteTooManyRequestsError(c *gin.Context, message string) {
	WriteError(c, http.StatusTooManyRequests, message)
}

// WriteServiceUnavailableError 发送服务不可用错误响应
func WriteServiceUnavailableError(c *gin.Context, message string) {
	WriteError(c, http.StatusServiceUnavailable, message)
}

// WriteGatewayTimeoutError 发送网关超时错误响应
func WriteGatewayTimeoutError(c *gin.Context, message string) {
	WriteError(c, http.StatusGatewayTimeout, message)
}

// BindError 处理绑定错误并发送响应
func BindError(c *gin.Context, message string) {
	if message == "" {
		message = "Invalid request parameters"
	}
	WriteValidationError(c, message)
}

// DatabaseError 处理数据库错误并发送响应
func DatabaseError(c *gin.Context, message string) {
	if message == "" {
		message = "Database operation failed"
	}
	WriteInternalServerError(c, message)
}

// ServiceError 处理服务层错误并发送响应
func ServiceError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	message := err.Error()
	if message == "" {
		message = "Service operation failed"
	}
	WriteInternalServerError(c, message)
}

// ValidationErrors 发送多个验证错误响应
func ValidationErrors(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusBadRequest, ErrorWithData(http.StatusBadRequest, "Validation failed", errors))
}

// FieldError 发送字段验证错误响应
func FieldError(c *gin.Context, field, message string) {
	ValidationErrors(c, map[string]string{field: message})
}

// RequiredFieldError 发送必填字段错误响应
func RequiredFieldError(c *gin.Context, field string) {
	FieldError(c, field, "This field is required")
}

// InvalidFieldError 发送无效字段错误响应
func InvalidFieldError(c *gin.Context, field string) {
	FieldError(c, field, "This field is invalid")
}

// DuplicateError 发送重复数据错误响应
func DuplicateError(c *gin.Context, field string) {
	FieldError(c, field, "This value already exists")
}

// NotFound 发送资源未找到错误响应
func NotFound(c *gin.Context, resource string) {
	WriteNotFoundError(c, resource+" not found")
}

// Unauthorized 发送未授权访问错误响应
func Unauthorized(c *gin.Context, resource string) {
	WriteUnauthorizedError(c, "Unauthorized access to "+resource)
}

// Forbidden 发送禁止访问错误响应
func Forbidden(c *gin.Context, resource string) {
	WriteForbiddenError(c, "Access to "+resource+" is forbidden")
}
