package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"gamelink/internal/repository"
)

const (
	// Password reset token expiration time
	passwordResetTokenExpiry = 15 * time.Minute
	// Redis key prefix for password reset tokens
	passwordResetKeyPrefix = "password_reset:"
)

var (
	// ErrInvalidToken is returned when the password reset token is invalid or expired
	ErrInvalidToken = errors.New("invalid or expired reset token")
	// ErrWeakPassword is returned when the password doesn't meet security requirements
	ErrWeakPassword = errors.New("password does not meet security requirements")
)

// PasswordResetService handles password reset operations
type PasswordResetService struct {
	users  repository.UserRepository
	redis  *redis.Client
	logger *slog.Logger
}

// NewPasswordResetService creates a new password reset service
func NewPasswordResetService(users repository.UserRepository, cacheClient interface{ GetClient() *redis.Client }, logger *slog.Logger) *PasswordResetService {
	var redisClient *redis.Client
	// Try to get the underlying Redis client from cache
	if cacheClient != nil {
		redisClient = cacheClient.GetClient()
	}

	return &PasswordResetService{
		users:  users,
		redis:  redisClient,
		logger: logger,
	}
}

// RequestReset initiates a password reset request by sending a reset email
// Security note: Always returns success to prevent user enumeration attacks
func (s *PasswordResetService) RequestReset(ctx context.Context, email string) error {
	// 1. Find user by email
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		// Don't expose whether user exists - always return success
		s.logger.Info("password reset requested for non-existent email", "email", email)
		return nil
	}

	// 2. Generate secure random token (32 bytes)
	resetToken, err := generateSecureToken(32)
	if err != nil {
		s.logger.Error("failed to generate reset token", "error", err, "user_id", user.ID)
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// 3. Store token in Redis with expiration
	if s.redis != nil {
		key := passwordResetKeyPrefix + resetToken
		if err := s.redis.Set(ctx, key, user.ID, passwordResetTokenExpiry).Err(); err != nil {
			s.logger.Error("failed to store reset token in redis", "error", err, "user_id", user.ID)
			return fmt.Errorf("failed to store reset token: %w", err)
		}
	} else {
		s.logger.Warn("redis not configured, password reset will not work in production")
		return fmt.Errorf("redis not configured for password reset")
	}

	// 4. In production, send email with reset link
	// For now, log the token (REMOVE IN PRODUCTION)
	resetURL := fmt.Sprintf("https://gamelink.com/reset-password?token=%s", resetToken)
	s.logger.Info("password reset link generated",
		"user_id", user.ID,
		"email", email,
		"reset_url", resetURL,
		"token_expires_in", passwordResetTokenExpiry)

	// TODO: Integrate email service (e.g., cfg.ExternalAPI.SMS.Enabled)
	// Example:
	// if err := s.emailService.Send(email, "Password Reset", resetLink); err != nil {
	//     return fmt.Errorf("failed to send reset email: %w", err)
	// }

	return nil
}

// ResetPassword resets the user's password using a valid reset token
func (s *PasswordResetService) ResetPassword(ctx context.Context, token, newPassword string) error {
	if s.redis == nil {
		return errors.New("redis not configured for password reset")
	}

	// 1. Validate token and get user ID
	key := passwordResetKeyPrefix + token
	userIDStr, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		s.logger.Warn("password reset attempted with invalid token", "token", token[:8]+"...")
		return ErrInvalidToken
	}
	if err != nil {
		s.logger.Error("failed to retrieve reset token from redis", "error", err)
		return fmt.Errorf("failed to validate reset token: %w", err)
	}

	// 2. Parse user ID
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		s.logger.Error("invalid user ID in reset token", "user_id", userIDStr, "error", err)
		return fmt.Errorf("invalid user ID in reset token: %w", err)
	}

	// 3. Validate password strength
	if err := s.validatePasswordStrength(newPassword); err != nil {
		s.logger.Info("weak password rejected during reset", "user_id", userID, "error", err)
		return ErrWeakPassword
	}

	// 4. Update user password
	if err := s.users.UpdatePassword(ctx, userID, newPassword); err != nil {
		s.logger.Error("failed to update user password", "user_id", userID, "error", err)
		return fmt.Errorf("failed to update password: %w", err)
	}

	// 5. Delete the used token (single-use only)
	if err := s.redis.Del(ctx, key).Err(); err != nil {
		s.logger.Warn("failed to delete used reset token", "token", token[:8]+"...", "error", err)
		// Non-critical error, continue
	}

	s.logger.Info("password reset successful", "user_id", userID)
	return nil
}

// validatePasswordStrength ensures the password meets security requirements
func (s *PasswordResetService) validatePasswordStrength(password string) error {
	// Minimum length
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Maximum length (prevent DoS)
	if len(password) > 128 {
		return errors.New("password must not exceed 128 characters")
	}

	// Check for character variety
	var (
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	)

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character (!@#$%^&* etc.)")
	}

	// Check for common weak passwords
	commonPasswords := []string{
		"Password123!", "Admin123!", "Welcome123!", "Qwerty123!",
		"Password1!", "12345678", "password", "admin",
	}
	for _, weak := range commonPasswords {
		if password == weak {
			return errors.New("password is too common, please choose a stronger one")
		}
	}

	return nil
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	// Use URL-safe base64 encoding
	return base64.URLEncoding.EncodeToString(b), nil
}
