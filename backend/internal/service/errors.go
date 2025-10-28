package service

import "errors"

// Common service-level sentinel errors for cross-layer mapping.
var (
    // ErrInvalidCredentials indicates bad username or password.
    ErrInvalidCredentials = errors.New("invalid credentials")

    // ErrUserDisabled indicates the user exists but is disabled/banned.
    ErrUserDisabled = errors.New("user account is disabled")
)

