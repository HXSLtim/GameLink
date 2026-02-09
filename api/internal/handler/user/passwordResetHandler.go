package user

import (
	"context"
	"errors"

	"log/slog"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	_ "gamelink/internal/model" // Imported for Swagger annotations
	"gamelink/pkg/apierr"
)

// PasswordResetSuccessResponse swagger response type
type PasswordResetSuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var (
	// ErrInvalidToken is returned when the password reset token is invalid or expired
	ErrInvalidToken = errors.New("invalid or expired reset token")
	// ErrWeakPassword is returned when the password doesn't meet security requirements
	ErrWeakPassword = errors.New("password does not meet security requirements")
)

// PasswordResetHandler handles password reset requests
type PasswordResetHandler struct {
	passwordResetService interface {
		RequestReset(ctx context.Context, email string) error
		ResetPassword(ctx context.Context, token, newPassword string) error
	}
}

// NewPasswordResetHandler creates a new password reset handler
func NewPasswordResetHandler(passwordResetService interface {
	RequestReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}) *PasswordResetHandler {
	return &PasswordResetHandler{
		passwordResetService: passwordResetService,
	}
}

// RequestPasswordResetRequest represents the request body for password reset initiation
type RequestPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents the request body for password reset confirmation
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

// RequestPasswordReset initiates a password reset by sending a reset email
// @Summary Request password reset
// @Description Sends a password reset link to the user's email
// @Tags user-auth
// @Accept json
// @Produce json
// @Param request body RequestPasswordResetRequest true "Email address"
// @Success 200 {object} PasswordResetSuccessResponse "If the email is registered, a reset link will be sent"
// @Router /api/v1/auth/password-reset/request [post]
func (h *PasswordResetHandler) RequestPasswordReset(c *gin.Context) {
	var req RequestPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("invalid email format"))
		return
	}

	ctx := c.Request.Context()

	// Always return success to prevent user enumeration
	if err := h.passwordResetService.RequestReset(ctx, req.Email); err != nil {
		// Log the error but don't expose it to the client
		slog.Default().Warn("password reset request failed",
			"email", req.Email,
			"error", err)
	}

	// Return success regardless of whether the email exists
	resp.Success[any](c, "If the email is registered, a password reset link will be sent", nil)
}

// ResetPassword resets the user's password using a valid reset token
// @Summary Reset password
// @Description Resets the user's password using a reset token from email
// @Tags user-auth
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Reset token and new password"
// @Success 200 {object} PasswordResetSuccessResponse "Password reset successful"
// @Failure 400 {object} apierr.APIError "Invalid or expired token"
// @Router /api/v1/auth/password-reset/confirm [post]
func (h *PasswordResetHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("invalid request format"))
		return
	}

	err := h.passwordResetService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			resp.Error(c, apierr.BadRequest("invalid or expired reset token"))
			return
		}
		if errors.Is(err, ErrWeakPassword) {
			resp.Error(c, apierr.BadRequest("password does not meet security requirements: must be 8+ characters with uppercase, lowercase, number, and special character"))
			return
		}
		resp.Error(c, apierr.InternalError("failed to reset password"))
		return
	}

	resp.Success[any](c, "Password reset successfully. Please login with your new password.", nil)
}
