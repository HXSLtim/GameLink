package service

import (
    "errors"
    "gamelink/internal/repository"
)

// Common service-level sentinel errors for cross-layer mapping.
var (
    // ErrInvalidCredentials indicates bad username or password.
    ErrInvalidCredentials = errors.New("invalid credentials")

    // ErrUserDisabled indicates the user exists but is disabled/banned.
    ErrUserDisabled = errors.New("user account is disabled")

    // ErrValidation indicates request/domain validation failed.
    ErrValidation = errors.New("validation failed")

    // ErrNotFound is a shared alias to repository not found.
    ErrNotFound = repository.ErrNotFound

    // ErrUserNotFound indicates the requested user does not exist.
    ErrUserNotFound = errors.New("user not found")

    // ErrOrderInvalidTransition indicates an invalid status change on order.
    ErrOrderInvalidTransition = errors.New("order invalid transition")
)

