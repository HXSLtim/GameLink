package role

import (
	"testing"

	"gamelink/internal/cache"
	"gamelink/internal/repository"
)

// mockRoleRepository 是一个简单的mock实现。
type mockRoleRepository struct {
	repository.RoleRepository
}

// mockCache 是一个简单的mock实现。
type mockCache struct {
	cache.Cache
}

// TestNewRoleService 测试构造函数。
func TestNewRoleService(t *testing.T) {
	repo := &mockRoleRepository{}
	cache := &mockCache{}

	svc := NewRoleService(repo, cache)

	if svc == nil {
		t.Fatal("NewRoleService returned nil")
	}

	if svc.roles != repo {
		t.Error("roles repository not set correctly")
	}

	if svc.cache != cache {
		t.Error("cache not set correctly")
	}
}
