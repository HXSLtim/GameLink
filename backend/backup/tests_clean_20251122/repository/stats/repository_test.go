package stats

import (
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"gamelink/internal/model"
)

func testContext() context.Context {
	return context.Background()
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.Player{}, &model.Game{}, &model.Order{}, &model.Payment{}, &model.OperationLog{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

// TestNewStatsRepository 测试构造函数。
func TestNewStatsRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	repo := NewStatsRepository(db)
	if repo == nil {
		t.Fatal("NewStatsRepository returned nil")
	}

	if _, ok := repo.(*gormStatsRepository); !ok {
		t.Errorf("expected *gormStatsRepository, got %T", repo)
	}
}

func TestStatsRepository_Dashboard(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStatsRepository(db)

	// Create test data
	now := time.Now()
	db.Create(&model.User{Phone: "13800138001", Email: "user1@test.com", Name: "User1"})
	db.Create(&model.User{Phone: "13800138002", Email: "user2@test.com", Name: "User2"})

	db.Create(&model.Player{UserID: 1, Nickname: "Player1", HourlyRateCents: 10000})

	db.Create(&model.Game{Name: "Game1", IconURL: "http://example.com/icon.jpg"})

	playerID := uint64(1)
	gameID := uint64(1)
	db.Create(&model.Order{OrderNo: "TEST_ORDER_001", UserID: 1, PlayerID: &playerID, GameID: &gameID, ItemID: 1, Title: "Order1", UnitPriceCents: 20000, TotalPriceCents: 20000, Status: model.OrderStatusPending})
	db.Create(&model.Order{OrderNo: "TEST_ORDER_002", UserID: 2, PlayerID: &playerID, GameID: &gameID, ItemID: 1, Title: "Order2", UnitPriceCents: 10000, TotalPriceCents: 10000, Status: model.OrderStatusCompleted})

	db.Create(&model.Payment{OrderID: 1, AmountCents: 20000, Status: "pending"})
	paidAt := now
	db.Create(&model.Payment{OrderID: 2, AmountCents: 10000, Status: "paid", PaidAt: &paidAt})

	dashboard, err := repo.Dashboard(testContext())
	if err != nil {
		t.Fatalf("Dashboard failed: %v", err)
	}

	if dashboard.TotalUsers != 2 {
		t.Errorf("expected 2 users, got %d", dashboard.TotalUsers)
	}
	if dashboard.TotalPlayers != 1 {
		t.Errorf("expected 1 player, got %d", dashboard.TotalPlayers)
	}
	if dashboard.TotalGames != 1 {
		t.Errorf("expected 1 game, got %d", dashboard.TotalGames)
	}
	if dashboard.TotalOrders != 2 {
		t.Errorf("expected 2 orders, got %d", dashboard.TotalOrders)
	}
	if dashboard.TotalPaidAmountCents != 10000 {
		t.Errorf("expected 10000 cents paid, got %d", dashboard.TotalPaidAmountCents)
	}
	if len(dashboard.OrdersByStatus) == 0 {
		t.Error("expected OrdersByStatus to have data")
	}
	if len(dashboard.PaymentsByStatus) == 0 {
		t.Error("expected PaymentsByStatus to have data")
	}
}

func TestStatsRepository_RevenueTrend(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStatsRepository(db)

	// Create test payments
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	db.Create(&model.Payment{OrderID: 1, AmountCents: 10000, Status: "paid", PaidAt: &yesterday})
	db.Create(&model.Payment{OrderID: 2, AmountCents: 15000, Status: "paid", PaidAt: &now})
	db.Create(&model.Payment{OrderID: 3, AmountCents: 5000, Status: "pending"}) // Not paid

	trend, err := repo.RevenueTrend(testContext(), 7)
	if err != nil {
		t.Fatalf("RevenueTrend failed: %v", err)
	}

	// Should have 1 or 2 date entries (depending on timezone)
	if len(trend) < 1 {
		t.Error("expected at least 1 date in revenue trend")
	}
}

func TestStatsRepository_UserGrowth(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStatsRepository(db)

	// Create test users with different creation dates
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	user1 := &model.User{Phone: "13800138001", Email: "user1@test.com", Name: "User1"}
	user1.CreatedAt = yesterday
	db.Create(user1)

	user2 := &model.User{Phone: "13800138002", Email: "user2@test.com", Name: "User2"}
	user2.CreatedAt = now
	db.Create(user2)

	growth, err := repo.UserGrowth(testContext(), 7)
	if err != nil {
		t.Fatalf("UserGrowth failed: %v", err)
	}

	if len(growth) < 1 {
		t.Error("expected at least 1 date in user growth")
	}
}

func TestStatsRepository_OrdersByStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStatsRepository(db)

	// Create test orders
	playerID := uint64(1)
	gameID := uint64(1)
	db.Create(&model.Order{OrderNo: "TEST_ORDER_101", UserID: 1, PlayerID: &playerID, GameID: &gameID, ItemID: 1, Title: "Order1", UnitPriceCents: 20000, TotalPriceCents: 20000, Status: model.OrderStatusPending})
	db.Create(&model.Order{OrderNo: "TEST_ORDER_102", UserID: 1, PlayerID: &playerID, GameID: &gameID, ItemID: 1, Title: "Order2", UnitPriceCents: 10000, TotalPriceCents: 10000, Status: model.OrderStatusPending})
	db.Create(&model.Order{OrderNo: "TEST_ORDER_103", UserID: 2, PlayerID: &playerID, GameID: &gameID, ItemID: 1, Title: "Order3", UnitPriceCents: 10000, TotalPriceCents: 10000, Status: model.OrderStatusCompleted})

	stats, err := repo.OrdersByStatus(testContext())
	if err != nil {
		t.Fatalf("OrdersByStatus failed: %v", err)
	}

	if len(stats) == 0 {
		t.Error("expected at least one status in stats")
	}

	if stats[string(model.OrderStatusPending)] != 2 {
		t.Errorf("expected 2 pending orders, got %d", stats[string(model.OrderStatusPending)])
	}
	if stats[string(model.OrderStatusCompleted)] != 1 {
		t.Errorf("expected 1 completed order, got %d", stats[string(model.OrderStatusCompleted)])
	}
}

func TestStatsRepository_TopPlayers(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStatsRepository(db)

	// Create test players
	db.Create(&model.Player{UserID: 1, Nickname: "TopPlayer", HourlyRateCents: 10000, RatingAverage: 4.8, RatingCount: 100})
	db.Create(&model.Player{UserID: 2, Nickname: "NewPlayer", HourlyRateCents: 8000, RatingAverage: 4.5, RatingCount: 10})
	db.Create(&model.Player{UserID: 3, Nickname: "ProPlayer", HourlyRateCents: 15000, RatingAverage: 4.9, RatingCount: 50})

	top, err := repo.TopPlayers(testContext(), 2)
	if err != nil {
		t.Fatalf("TopPlayers failed: %v", err)
	}

	if len(top) != 2 {
		t.Errorf("expected 2 top players, got %d", len(top))
	}

	// First player should have highest rating count (100)
	if len(top) > 0 && top[0].Nickname != "TopPlayer" {
		t.Errorf("expected TopPlayer to be first, got %s", top[0].Nickname)
	}
}

func TestStatsRepository_AuditOverview(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStatsRepository(db)

	// Create test operation logs
	actorID := uint64(1)
	db.Create(&model.OperationLog{EntityType: "user", EntityID: 1, Action: "create", ActorUserID: &actorID})
	db.Create(&model.OperationLog{EntityType: "user", EntityID: 2, Action: "update", ActorUserID: &actorID})
	db.Create(&model.OperationLog{EntityType: "order", EntityID: 1, Action: "create", ActorUserID: &actorID})

	byEntity, byAction, err := repo.AuditOverview(testContext(), nil, nil)
	if err != nil {
		t.Fatalf("AuditOverview failed: %v", err)
	}

	if len(byEntity) == 0 {
		t.Error("expected byEntity to have data")
	}
	if len(byAction) == 0 {
		t.Error("expected byAction to have data")
	}

	if byEntity["user"] < 2 {
		t.Errorf("expected at least 2 user entities, got %d", byEntity["user"])
	}
	if byEntity["order"] < 1 {
		t.Errorf("expected at least 1 order entity, got %d", byEntity["order"])
	}
	if byAction["create"] < 1 {
		t.Errorf("expected at least 1 create action, got %d", byAction["create"])
	}
	if byAction["update"] < 1 {
		t.Errorf("expected at least 1 update action, got %d", byAction["update"])
	}
}

func TestStatsRepository_AuditTrend(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStatsRepository(db)

	// Create test operation logs
	actorID := uint64(1)

	log1 := &model.OperationLog{EntityType: "user", EntityID: 1, Action: "create", ActorUserID: &actorID}
	log1.CreatedAt = time.Now().Add(-24 * time.Hour)
	db.Create(log1)

	db.Create(&model.OperationLog{EntityType: "user", EntityID: 2, Action: "create", ActorUserID: &actorID})
	db.Create(&model.OperationLog{EntityType: "order", EntityID: 1, Action: "create", ActorUserID: &actorID})

	t.Run("All logs", func(t *testing.T) {
		trend, err := repo.AuditTrend(testContext(), nil, nil, "", "")
		if err != nil {
			t.Fatalf("AuditTrend failed: %v", err)
		}
		if len(trend) < 1 {
			t.Error("expected at least 1 date in audit trend")
		}
	})

	t.Run("Filtered by entity", func(t *testing.T) {
		trend, err := repo.AuditTrend(testContext(), nil, nil, "user", "")
		if err != nil {
			t.Fatalf("AuditTrend failed: %v", err)
		}
		if len(trend) < 1 {
			t.Error("expected at least 1 date in filtered audit trend")
		}
	})

	t.Run("Filtered by action", func(t *testing.T) {
		trend, err := repo.AuditTrend(testContext(), nil, nil, "", "create")
		if err != nil {
			t.Fatalf("AuditTrend failed: %v", err)
		}
		if len(trend) < 1 {
			t.Error("expected at least 1 date in filtered audit trend")
		}
	})
}

func TestStatsRepository_EmptyData(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStatsRepository(db)

	t.Run("Empty Dashboard", func(t *testing.T) {
		dashboard, err := repo.Dashboard(testContext())
		if err != nil {
			t.Fatalf("Dashboard failed: %v", err)
		}
		if dashboard.TotalUsers != 0 {
			t.Errorf("expected 0 users, got %d", dashboard.TotalUsers)
		}
	})

	t.Run("Empty RevenueTrend", func(t *testing.T) {
		trend, err := repo.RevenueTrend(testContext(), 7)
		if err != nil {
			t.Fatalf("RevenueTrend failed: %v", err)
		}
		if len(trend) != 0 {
			t.Error("expected empty revenue trend")
		}
	})

	t.Run("Empty UserGrowth", func(t *testing.T) {
		growth, err := repo.UserGrowth(testContext(), 7)
		if err != nil {
			t.Fatalf("UserGrowth failed: %v", err)
		}
		if len(growth) != 0 {
			t.Error("expected empty user growth")
		}
	})
}
