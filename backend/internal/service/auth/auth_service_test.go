package auth

import (
	"testing"
	"time"

	"gamelink/internal/auth"
	"gamelink/internal/repository"
)

// mockUserRepository 是一个简单的mock实现。
type mockUserRepository struct {
	repository.UserRepository
}

// TestNewAuthService 测试构造函数。
func TestNewAuthService(t *testing.T) {
	userRepo := &mockUserRepository{}
	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)

	svc := NewAuthService(userRepo, jwtManager)

	if svc == nil {
		t.Fatal("NewAuthService returned nil")
	}

	if svc.userRepo != userRepo {
		t.Error("userRepo not set correctly")
	}

	if svc.jwtManager != jwtManager {
		t.Error("jwtManager not set correctly")
	}
}

// TestAuthServiceWithNilInputs 测试nil输入处理。
func TestAuthServiceWithNilInputs(t *testing.T) {
	// 虽然不推荐传nil，但我们测试service是否会panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("NewAuthService panicked with nil inputs: %v", r)
		}
	}()

	_ = NewAuthService(nil, nil)
}
