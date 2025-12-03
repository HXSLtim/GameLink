package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/pkg/apierr"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service"
)

// MapServiceError 将服务层错误映射为标准的API错误响应
// 这是统一错误处理的核心函数，确保所有错误都有一致的响应格式
func MapServiceError(err error) (int, *model.APIResponse[any]) {
	return MapServiceErrorWithPath(err, "")
}

// MapServiceErrorWithPath 将服务层错误映射为标准的API错误响应，支持基于路径的特定消息
func MapServiceErrorWithPath(err error, path string) (int, *model.APIResponse[any]) {
	if err == nil {
		return http.StatusOK, &model.APIResponse[any]{
			Success: true,
			Code:    http.StatusOK,
			Message: "success",
		}
	}

	// 检查是否为APIError类型
	if apiErr, ok := err.(*apierr.APIError); ok {
		// 构建meta字段，包含额外的错误信息
		meta := make(map[string]interface{})
		if apiErr.Details != "" {
			meta["details"] = apiErr.Details
		}
		if apiErr.Field != "" {
			meta["field"] = apiErr.Field
		}
		if apiErr.Timestamp != 0 {
			meta["timestamp"] = apiErr.Timestamp
		}
		if apiErr.Extensions != nil {
			for k, v := range apiErr.Extensions {
				meta[k] = v
			}
		}

		resp := &model.APIResponse[any]{
			Success: false,
			Code:    apiErr.Code,
			Message: apiErr.Message,
		}

		// 如果有meta数据，添加到响应中
		if len(meta) > 0 {
			resp.Meta = meta
		}

		return apiErr.Code, resp
	}

	// 检查服务层错误
	switch {
	case err != nil:
		return http.StatusBadRequest, &model.APIResponse[any]{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "validation failed",
		}
	case errors.Is(err, service.ErrInvalidCredentials):
		return http.StatusUnauthorized, &model.APIResponse[any]{
			Success: false,
			Code:    http.StatusUnauthorized,
			Message: "invalid credentials",
		}
	case errors.Is(err, service.ErrUserDisabled):
		return http.StatusForbidden, &model.APIResponse[any]{
			Success: false,
			Code:    http.StatusForbidden,
			Message: "user account is disabled",
		}
	case errors.Is(err, service.ErrOrderInvalidTransition):
		return http.StatusConflict, &model.APIResponse[any]{
			Success: false,
			Code:    http.StatusConflict,
			Message: apierr.ErrOrderInvalidTransition,
		}
	case errors.Is(err, service.ErrUserNotFound):
		return http.StatusNotFound, &model.APIResponse[any]{
			Success: false,
			Code:    http.StatusNotFound,
			Message: apierr.ErrUserNotFound,
		}
	case errors.Is(err, service.ErrNotFound) || errors.Is(err, repository.ErrNotFound):
		message := "resource not found"
		// 基于路径返回特定的404消息
		if path != "" {
			message = getDomainNotFoundMessage(path)
		}
		return http.StatusNotFound, &model.APIResponse[any]{
			Success: false,
			Code:    http.StatusNotFound,
			Message: message,
		}
	default:
		// 默认返回500错误
		return http.StatusInternalServerError, &model.APIResponse[any]{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
}

// getDomainNotFoundMessage returns a stable message for 404 based on route path.
func getDomainNotFoundMessage(path string) string {
	p := strings.ToLower(path)
	switch {
	case strings.Contains(p, "/users"):
		return apierr.ErrUserNotFound
	case strings.Contains(p, "/orders"):
		return apierr.ErrOrderNotFound
	case strings.Contains(p, "/payments"):
		return apierr.ErrPaymentNotFound
	case strings.Contains(p, "/players"):
		return apierr.ErrPlayerNotFound
	case strings.Contains(p, "/games"):
		return apierr.ErrGameNotFound
	default:
		return "resource not found"
	}
}

// RespondWithServiceError 使用标准格式响应服务层错误
func RespondWithServiceError(c *gin.Context, err error) {
	RespondWithServiceErrorAndPath(c, err, "")
}

// RespondWithServiceErrorAndPath 使用标准格式响应服务层错误，支持基于路径的特定消息
func RespondWithServiceErrorAndPath(c *gin.Context, err error, path string) {
	statusCode, response := MapServiceErrorWithPath(err, path)

	// 添加请求ID到响应中
	if requestID := GetRequestID(c); requestID != "" {
		response.TraceID = requestID
	}

	c.JSON(statusCode, response)
}

// ValidateAndRespond 验证请求并响应错误（如果需要）
func ValidateAndRespond(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		RespondWithServiceError(c, apierr.BadRequest(apierr.ErrInvalidJSONPayload).WithDetails(err.Error()))
		return false
	}
	return true
}

// ValidateQueryAndRespond 验证查询参数并响应错误（如果需要）
func ValidateQueryAndRespond(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		RespondWithServiceError(c, apierr.BadRequest("invalid query parameters").WithDetails(err.Error()))
		return false
	}
	return true
}

// ParseIDAndRespond 解析ID参数并响应错误（如果需要）
func ParseIDAndRespond(c *gin.Context, paramName string) (uint64, bool) {
	id, err := ParseIDParam(c, paramName)
	if err != nil {
		RespondWithServiceError(c, err)
		return 0, false
	}
	return id, true
}
