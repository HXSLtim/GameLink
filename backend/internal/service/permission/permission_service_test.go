package permission

import (
	"testing"

	"gamelink/internal/cache"
	"gamelink/internal/repository"
)

// mockPermissionRepository 是一个简单的mock实现。
type mockPermissionRepository struct {
	repository.PermissionRepository
}

// mockCache 是一个简单的mock实现。
type mockCache struct {
	cache.Cache
}

// TestNewPermissionService 测试构造函数。
func TestNewPermissionService(t *testing.T) {
	permRepo := &mockPermissionRepository{}
	cache := &mockCache{}

	svc := NewPermissionService(permRepo, cache)

	if svc == nil {
		t.Fatal("NewPermissionService returned nil")
	}

	if svc.permissions != permRepo {
		t.Error("permissions repository not set correctly")
	}

	if svc.cache != cache {
		t.Error("cache not set correctly")
	}
}
