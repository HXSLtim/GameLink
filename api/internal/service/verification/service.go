package verification

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"gamelink/internal/service/sms"
	"gamelink/pkg/cache"
)

const (
	CodeLength      = 6
	CodeExpiry      = 5 * time.Minute
	MaxAttempts     = 5
	RateLimitWindow = 60 * time.Second
	MasterCode      = "666666"
)

var (
	ErrCodeExpired     = errors.New("verification code expired")
	ErrCodeInvalid     = errors.New("invalid verification code")
	ErrTooManyAttempts = errors.New("too many verification attempts")
	ErrRateLimited     = errors.New("please wait before requesting another code")
	ErrCodeNotFound    = errors.New("verification code not found")
)

type CodeType string

const (
	CodeTypePhone CodeType = "phone"
	CodeTypeEmail CodeType = "email"
)

type Service struct {
	cache       cache.Cache
	masterCode  string
	environment string
	smsSvc      *sms.Service
}

func NewService(c cache.Cache, smsSvc *sms.Service) *Service {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	masterCode := os.Getenv("VERIFICATION_MASTER_CODE")
	if masterCode == "" {
		masterCode = MasterCode
	}
	return &Service{
		cache:       c,
		masterCode:  masterCode,
		environment: env,
		smsSvc:      smsSvc,
	}
}

func (s *Service) GenerateCode(ctx context.Context, target string, codeType CodeType) (string, error) {
	rateLimitKey := s.rateLimitKey(target, codeType)
	if _, ok, _ := s.cache.Get(ctx, rateLimitKey); ok {
		return "", ErrRateLimited
	}

	code, err := generateRandomCode(CodeLength)
	if err != nil {
		return "", fmt.Errorf("failed to generate code: %w", err)
	}

	codeKey := s.codeKey(target, codeType)
	codeData := codeRecord{
		Code:      code,
		Attempts:  0,
		CreatedAt: time.Now(),
	}

	jsonData, err := json.Marshal(codeData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal code data: %w", err)
	}

	if err := s.cache.Set(ctx, codeKey, string(jsonData), CodeExpiry); err != nil {
		return "", fmt.Errorf("failed to store code: %w", err)
	}

	if err := s.cache.Set(ctx, rateLimitKey, "1", RateLimitWindow); err != nil {
		fmt.Printf("[Verification] Failed to set rate limit: %v\n", err)
	}

	return code, nil
}

// SendCode sends verification code via SMS or email
func (s *Service) SendCode(ctx context.Context, target string, codeType CodeType) error {
	code, err := s.GenerateCode(ctx, target, codeType)
	if err != nil {
		return err
	}

	switch codeType {
	case CodeTypePhone:
		if s.smsSvc != nil {
			return s.smsSvc.SendCode(ctx, target, code)
		}
		fmt.Printf("[Verification] SMS not configured, code for %s: %s\n", target, code)
	case CodeTypeEmail:
		fmt.Printf("[Verification] Email not implemented, code for %s: %s\n", target, code)
	}

	return nil
}

func (s *Service) VerifyCode(ctx context.Context, target, code string, codeType CodeType) error {
	if s.isMasterCodeEnabled() && code == s.masterCode {
		fmt.Printf("[Verification] Master code used for %s: %s\n", codeType, target)
		return nil
	}

	codeKey := s.codeKey(target, codeType)

	jsonData, ok, err := s.cache.Get(ctx, codeKey)
	if err != nil {
		return fmt.Errorf("failed to get code: %w", err)
	}
	if !ok {
		return ErrCodeNotFound
	}

	var record codeRecord
	if err := json.Unmarshal([]byte(jsonData), &record); err != nil {
		return fmt.Errorf("failed to unmarshal code data: %w", err)
	}

	if record.Attempts >= MaxAttempts {
		_ = s.cache.Delete(ctx, codeKey)
		return ErrTooManyAttempts
	}

	if record.Code != code {
		record.Attempts++
		remaining := CodeExpiry - time.Since(record.CreatedAt)
		if remaining > 0 {
			newData, _ := json.Marshal(record)
			_ = s.cache.Set(ctx, codeKey, string(newData), remaining)
		}
		return ErrCodeInvalid
	}

	_ = s.cache.Delete(ctx, codeKey)
	return nil
}

func (s *Service) IsMasterCode(code string) bool {
	return s.isMasterCodeEnabled() && code == s.masterCode
}

func (s *Service) GetMasterCode() string {
	if s.isMasterCodeEnabled() {
		return s.masterCode
	}
	return ""
}

func (s *Service) isMasterCodeEnabled() bool {
	return s.environment != "production"
}

func (s *Service) codeKey(target string, codeType CodeType) string {
	return fmt.Sprintf("verification:%s:%s", codeType, target)
}

func (s *Service) rateLimitKey(target string, codeType CodeType) string {
	return fmt.Sprintf("verification:ratelimit:%s:%s", codeType, target)
}

type codeRecord struct {
	Code      string    `json:"code"`
	Attempts  int       `json:"attempts"`
	CreatedAt time.Time `json:"created_at"`
}

func generateRandomCode(length int) (string, error) {
	const digits = "0123456789"
	code := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code[i] = digits[n.Int64()]
	}
	return string(code), nil
}
