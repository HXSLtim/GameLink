package response

import (
	"net/http"

	"gamelink/internal/model"
)

// Response 定义统一的响应结构
type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 创建成功响应
func Success(data interface{}) Response {
	return Response{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
	}
}

// Error 创建错误响应
func Error(code int, message string) Response {
	return Response{
		Success: false,
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

// ErrorWithData 创建带数据的错误响应
func ErrorWithData(code int, message string, data interface{}) Response {
	return Response{
		Success: false,
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// ValidationError 创建验证错误响应
func ValidationError(message string) Response {
	return Error(http.StatusBadRequest, message)
}

// NotFoundError 创建未找到错误响应
func NotFoundError(message string) Response {
	return Error(http.StatusNotFound, message)
}

// UnauthorizedError 创建未授权错误响应
func UnauthorizedError(message string) Response {
	return Error(http.StatusUnauthorized, message)
}

// ForbiddenError 创建禁止访问错误响应
func ForbiddenError(message string) Response {
	return Error(http.StatusForbidden, message)
}

// InternalServerError 创建服务器内部错误响应
func InternalServerError(message string) Response {
	return Error(http.StatusInternalServerError, message)
}

// BadRequestError 创建错误请求响应
func BadRequestError(message string) Response {
	return Error(http.StatusBadRequest, message)
}

// ConvertFromModel 将 model.APIResponse[T] 转换为 Response
func ConvertFromModel[T any](apiResp model.APIResponse[T]) Response {
	return Response{
		Success: apiResp.Success,
		Code:    apiResp.Code,
		Message: apiResp.Message,
		Data:    apiResp.Data,
	}
}
