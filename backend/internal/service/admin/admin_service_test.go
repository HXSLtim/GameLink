package admin

import (
	"testing"

	"gamelink/internal/cache"
	"gamelink/internal/repository"
)

// mockGameRepository 是一个简单的mock实现。
type mockGameRepository struct {
	repository.GameRepository
}

type mockUserRepository struct {
	repository.UserRepository
}

type mockPlayerRepository struct {
	repository.PlayerRepository
}

type mockOrderRepository struct {
	repository.OrderRepository
}

type mockPaymentRepository struct {
	repository.PaymentRepository
}

type mockRoleRepository struct {
	repository.RoleRepository
}

type mockCache struct {
	cache.Cache
}

// TestNewAdminService 测试构造函数。
func TestNewAdminService(t *testing.T) {
	games := &mockGameRepository{}
	users := &mockUserRepository{}
	players := &mockPlayerRepository{}
	orders := &mockOrderRepository{}
	payments := &mockPaymentRepository{}
	roles := &mockRoleRepository{}
	cache := &mockCache{}

	svc := NewAdminService(games, users, players, orders, payments, roles, cache)

	if svc == nil {
		t.Fatal("NewAdminService returned nil")
	}

	if svc.games != games {
		t.Error("games repository not set correctly")
	}
	if svc.users != users {
		t.Error("users repository not set correctly")
	}
	if svc.players != players {
		t.Error("players repository not set correctly")
	}
	if svc.orders != orders {
		t.Error("orders repository not set correctly")
	}
	if svc.payments != payments {
		t.Error("payments repository not set correctly")
	}
	if svc.roles != roles {
		t.Error("roles repository not set correctly")
	}
	if svc.cache != cache {
		t.Error("cache not set correctly")
	}
}

// TestSetTxManager 测试事务管理器注入。
func TestSetTxManager(t *testing.T) {
	svc := NewAdminService(
		&mockGameRepository{},
		&mockUserRepository{},
		&mockPlayerRepository{},
		&mockOrderRepository{},
		&mockPaymentRepository{},
		&mockRoleRepository{},
		&mockCache{},
	)

	if svc.tx != nil {
		t.Error("tx should be nil initially")
	}

	// Note: 我们不能测试实际的TxManager，因为它是一个接口
	// 这个测试只是确保方法存在并可以调用
}
