package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"unicode"
)

// SecretValidationError represents a JWT secret validation error
type SecretValidationError struct {
	Message string
}

func (e *SecretValidationError) Error() string {
	return e.Message
}

// Common weak secrets that should never be used
var weakSecrets = []string{
	"secret",
	"password",
	"12345678",
	"test",
	"development",
	"changeme",
	"default",
	"admin",
	"jwt_secret",
	"your-secret-key",
	"your_secret_key",
	"yoursecretkey",
}

// ValidateJWTSecret validates JWT secret key strength
// Requirements:
// - Minimum 64 characters (increased from 32)
// - Must contain uppercase letters
// - Must contain lowercase letters
// - Must contain digits
// - Must contain special characters
// - Must not be a common weak secret
// - Must have sufficient entropy
func ValidateJWTSecret(secret string) error {
	// Check minimum length
	if len(secret) < 64 {
		return &SecretValidationError{
			Message: "JWT secret must be at least 64 characters long",
		}
	}

	// Check for common weak secrets (case-insensitive)
	secretLower := strings.ToLower(secret)
	for _, weak := range weakSecrets {
		if strings.Contains(secretLower, weak) {
			return &SecretValidationError{
				Message: "JWT secret contains common weak patterns",
			}
		}
	}

	// Check for repeated characters (e.g., "aaaaaaa...")
	if hasRepeatedChars(secret, 5) {
		return &SecretValidationError{
			Message: "JWT secret contains too many repeated characters",
		}
	}

	// Check character diversity
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, ch := range secret {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return &SecretValidationError{
			Message: "JWT secret must contain at least one uppercase letter",
		}
	}
	if !hasLower {
		return &SecretValidationError{
			Message: "JWT secret must contain at least one lowercase letter",
		}
	}
	if !hasDigit {
		return &SecretValidationError{
			Message: "JWT secret must contain at least one digit",
		}
	}
	if !hasSpecial {
		return &SecretValidationError{
			Message: "JWT secret must contain at least one special character",
		}
	}

	// Check entropy (simple check: unique character ratio)
	if !hasSufficientEntropy(secret) {
		return &SecretValidationError{
			Message: "JWT secret has insufficient entropy (too predictable)",
		}
	}

	return nil
}

// hasRepeatedChars checks if the string has more than maxRepeat consecutive repeated characters
func hasRepeatedChars(s string, maxRepeat int) bool {
	if len(s) < maxRepeat {
		return false
	}

	count := 1
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			count++
			if count >= maxRepeat {
				return true
			}
		} else {
			count = 1
		}
	}
	return false
}

// hasSufficientEntropy checks if the secret has sufficient randomness
// Uses a simple heuristic: unique characters should be at least 40% of total length
func hasSufficientEntropy(s string) bool {
	if len(s) == 0 {
		return false
	}

	uniqueChars := make(map[rune]bool)
	for _, ch := range s {
		uniqueChars[ch] = true
	}

	uniqueRatio := float64(len(uniqueChars)) / float64(len(s))
	return uniqueRatio >= 0.4
}

// GenerateSecretHash generates a SHA-256 hash of the secret for logging/comparison
// This allows checking if a secret has changed without exposing the actual secret
func GenerateSecretHash(secret string) string {
	hash := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(hash[:])
}

// IsSecretWeak performs a quick check if the secret is obviously weak
// This is a faster check than ValidateJWTSecret for startup validation
func IsSecretWeak(secret string) bool {
	// Too short
	if len(secret) < 32 {
		return true
	}

	// All same character (check manually instead of regex)
	if len(secret) > 0 {
		allSame := true
		firstChar := secret[0]
		for i := 1; i < len(secret); i++ {
			if secret[i] != firstChar {
				allSame = false
				break
			}
		}
		if allSame {
			return true
		}
	}

	// Common weak patterns
	secretLower := strings.ToLower(secret)
	for _, weak := range weakSecrets {
		if secretLower == weak || strings.HasPrefix(secretLower, weak) {
			return true
		}
	}

	return false
}

// ValidateJWTSecretOrPanic validates JWT secret and panics if invalid
// Use this during application startup to fail fast on weak secrets
func ValidateJWTSecretOrPanic(secret string) {
	if err := ValidateJWTSecret(secret); err != nil {
		panic("JWT secret validation failed: " + err.Error())
	}
}

// ValidateJWTSecretWithWarning validates JWT secret and returns error
// Use this for non-critical validation where you want to log warnings
func ValidateJWTSecretWithWarning(secret string) error {
	// Quick weak check
	if IsSecretWeak(secret) {
		return errors.New("JWT secret is obviously weak")
	}

	// Full validation
	return ValidateJWTSecret(secret)
}
