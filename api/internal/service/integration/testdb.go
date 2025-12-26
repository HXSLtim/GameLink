// Package integration provides integration test utilities with PostgreSQL.
package integration

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"gamelink/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	testDB   *gorm.DB
	initOnce sync.Once
)

// TestDBConfig holds PostgreSQL test database configuration.
type TestDBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// DefaultTestConfig returns default test database configuration.
// Uses environment variables or defaults to local Docker PostgreSQL.
func DefaultTestConfig() TestDBConfig {
	return TestDBConfig{
		Host:     getEnvOrDefault("TEST_DB_HOST", "localhost"),
		Port:     getEnvOrDefault("TEST_DB_PORT", "5432"),
		User:     getEnvOrDefault("TEST_DB_USER", "gamelink"),
		Password: getEnvOrDefault("TEST_DB_PASSWORD", "gamelink"),
		DBName:   getEnvOrDefault("TEST_DB_NAME", "gamelink_test"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// DSN returns the PostgreSQL connection string.
func (c TestDBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		c.Host, c.Port, c.User, c.Password, c.DBName,
	)
}

// SetupTestDB initializes the test database connection.
// It creates tables and returns a clean database for each test.
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	initOnce.Do(func() {
		config := DefaultTestConfig()
		var err error
		testDB, err = gorm.Open(postgres.Open(config.DSN()), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			t.Fatalf("Failed to connect to test database: %v", err)
		}

		// Auto migrate all models
		if err := migrateModels(testDB); err != nil {
			t.Fatalf("Failed to migrate models: %v", err)
		}
	})

	if testDB == nil {
		t.Fatal("Test database not initialized")
	}

	// Clean all tables before each test
	cleanTables(t, testDB)

	return testDB
}

// migrateModels creates all database tables.
func migrateModels(db *gorm.DB) error {
	return db.AutoMigrate(
		// Core models
		&model.User{},
		&model.Player{},
		&model.Game{},
		&model.ServiceItem{},
		&model.Order{},
		&model.OrderItem{},
		&model.OrderPlayer{},
		&model.Payment{},
		&model.Wallet{},
		&model.Withdraw{},

		// Review & Dispute
		&model.Review{},
		&model.ReviewReply{},
		&model.ReviewAppeal{},
		&model.OrderDispute{},

		// Statistics
		&model.PlatformStatistics{},
		&model.PlayerStatistics{},
		&model.UserStatistics{},
		&model.ServiceItemStatistics{},
		&model.GameStatistics{},

		// KPI
		&model.KPITarget{},

		// Other models as needed
		&model.VipLevel{},
		&model.CouponTemplate{},
		&model.Coupon{},
		&model.RechargeOption{},
		&model.RechargeRecord{},
		&model.Activity{},
		&model.Team{},
		&model.TeamMember{},
	)
}

// cleanTables truncates all tables to ensure test isolation.
func cleanTables(t *testing.T, db *gorm.DB) {
	t.Helper()

	tables := []string{
		"team_members", "teams",
		"activity_participations", "activity_rewards", "activities",
		"recharge_records", "recharge_options",
		"coupons", "coupon_templates",
		"vip_levels",
		"kpi_targets",
		"user_statistics", "player_statistics", "platform_statistics",
		"order_disputes", "review_appeals", "review_replies", "reviews",
		"withdraws", "wallets", "payments",
		"order_players", "order_items", "orders",
		"service_items", "games",
		"players", "users",
	}

	// Disable foreign key checks, truncate, then re-enable
	db.Exec("SET session_replication_role = 'replica'")
	for _, table := range tables {
		db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
	}
	db.Exec("SET session_replication_role = 'origin'")
}

// CreateTestUser creates a test user and returns it.
func CreateTestUser(t *testing.T, db *gorm.DB, name string) *model.User {
	t.Helper()
	user := &model.User{
		Name:   name,
		Phone:  fmt.Sprintf("138%08d", len(name)),
		Email:  fmt.Sprintf("%s@test.com", name),
		Role:   model.RoleUser,
		Status: model.UserStatusActive,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

// CreateTestPlayer creates a test player and returns it.
func CreateTestPlayer(t *testing.T, db *gorm.DB, user *model.User) *model.Player {
	t.Helper()
	player := &model.Player{
		UserID:             user.ID,
		Nickname:           user.Name + "_player",
		VerificationStatus: model.VerificationVerified,
	}
	if err := db.Create(player).Error; err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}

	// Update user role
	db.Model(user).Update("role", model.RolePlayer)

	return player
}

// CreateTestGame creates a test game and returns it.
func CreateTestGame(t *testing.T, db *gorm.DB, name string) *model.Game {
	t.Helper()
	game := &model.Game{
		Key:      name,
		Name:     name,
		Category: "moba",
		IsActive: true,
	}
	if err := db.Create(game).Error; err != nil {
		t.Fatalf("Failed to create test game: %v", err)
	}
	return game
}

// CreateTestOrder creates a test order and returns it.
func CreateTestOrder(t *testing.T, db *gorm.DB, user *model.User, player *model.Player, status model.OrderStatus) *model.Order {
	t.Helper()
	order := &model.Order{
		OrderNo:         fmt.Sprintf("TEST%d", user.ID),
		UserID:          user.ID,
		PlayerID:        &player.ID,
		TotalPriceCents: 10000,
		Status:          status,
		Currency:        model.CurrencyCNY,
	}
	if err := db.Create(order).Error; err != nil {
		t.Fatalf("Failed to create test order: %v", err)
	}
	return order
}

// CreateTestPayment creates a test payment and returns it.
func CreateTestPayment(t *testing.T, db *gorm.DB, order *model.Order, status model.PaymentStatus) *model.Payment {
	t.Helper()
	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      order.UserID,
		AmountCents: order.TotalPriceCents,
		Method:      model.PaymentMethodWeChat,
		Status:      status,
	}
	if err := db.Create(payment).Error; err != nil {
		t.Fatalf("Failed to create test payment: %v", err)
	}
	return payment
}

// SkipIfNoTestDB skips the test if test database is not available.
func SkipIfNoTestDB(t *testing.T) {
	t.Helper()
	config := DefaultTestConfig()
	db, err := gorm.Open(postgres.Open(config.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skipf("Skipping integration test: test database not available: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.Close()
}
