package service

import (
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
)

// Common service-level sentinel errors for cross-layer mapping.
var (
	// ErrInvalidCredentials indicates bad username or password.
	ErrInvalidCredentials = apierr.Unauthorized("用户名或密码错误")

	// ErrUserDisabled indicates the user exists but is disabled/banned.
	ErrUserDisabled = apierr.Forbidden("用户账户已禁用")

	// ErrValidation indicates request/domain validation failed.
	ErrValidation = apierr.BadRequest("验证失败")

	// ErrNotFound is a shared alias to repository not found.
	ErrNotFound = repository.ErrNotFound

	// ErrUserNotFound indicates the requested user does not exist.
	ErrUserNotFound = apierr.NotFound("用户不存在")

	// ErrOrderInvalidTransition indicates an invalid status change on order.
	ErrOrderInvalidTransition = apierr.BadRequest("订单状态转换无效")
)
