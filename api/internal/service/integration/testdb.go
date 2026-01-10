// Package integration provides integration test utilities with PostgreSQL.
package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"gamelink/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	testDB     *gorm.DB
	initOnce   sync.Once
	cleanMutex sync.Mutex
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
		&model.GameCategory{},
		&model.ServiceItem{},
		&model.Order{},
		&model.OrderItem{},
		&model.OrderPlayer{},
		&model.Payment{},
		&model.Wallet{},
		&model.Withdraw{},

		// RBAC - Permission, Role, Menu
		&model.Permission{},
		&model.RoleModel{},
		&model.RolePermission{},
		&model.UserRole{},
		&model.Menu{},
		&model.PermissionAuditLog{},

		// Operation Logs
		&model.OperationLog{},

		// Review & Dispute
		&model.Review{},
		&model.ReviewReply{},
		&model.ReviewAppeal{},
		&model.OrderDispute{},

		// Chat
		&model.ChatGroup{},
		&model.ChatGroupMember{},
		&model.ChatMessage{},
		&model.ChatReport{},

		// Notification
		&model.NotificationEvent{},
		&model.UserNotification{},

		// Sensitive Word
		&model.SensitiveWord{},

		// Statistics
		&model.PlatformStatistics{},
		&model.PlayerStatistics{},
		&model.UserStatistics{},
		&model.ServiceItemStatistics{},
		&model.GameStatistics{},

		// KPI
		&model.KPITarget{},

		// VIP & Coupon
		&model.VipLevel{},
		&model.VipConfig{},
		&model.CouponTemplate{},
		&model.Coupon{},

		// Recharge
		&model.RechargeOption{},
		&model.RechargeRecord{},

		// Activity
		&model.Activity{},
		&model.ActivityReward{},
		&model.ActivityParticipation{},
		&model.ActivityDailyStats{},

		// Team
		&model.Team{},
		&model.TeamMember{},
		&model.TeamInvite{},

		// Referral
		&model.ReferralConfig{},
		&model.ReferralCode{},
		&model.Referral{},
		&model.ReferralReward{},

		// Commission
		&model.CommissionRule{},
		&model.CommissionRecord{},
		&model.MonthlySettlement{},

		// User Block
		&model.UserBlock{},

		// Favorite
		&model.Favorite{},

		// Settlement
		&model.SettlementCompany{},
		&model.PlayerCompanyAssignment{},
		&model.SalaryPaymentRecord{},

		// Player Rank & Certification
		&model.GameRank{},
		&model.PlayerRankRecord{},
		&model.PlayerCertification{},

		// Order Timeout
		&model.OrderTimeoutConfig{},
		&model.OrderTimeoutLog{},
		&model.OrderServiceAssignment{},

		// Dispute Templates & Chat Snapshots
		&model.DisputeTemplate{},
		&model.ChatSnapshot{},

		// User Tag
		&model.UserTag{},
		&model.UserTagRelation{},

		// Ranking
		&model.PlayerRanking{},
		&model.RankingCommissionConfig{},
		&model.RankingReward{},

		// Feed/Content
		&model.Feed{},
		&model.FeedImage{},
		&model.FeedReport{},
		&model.ContentCategory{},

		// Collection Entity & Routing Rule
		&model.CollectionEntity{},
		&model.PaymentChannelConfig{},
		&model.RoutingRule{},
		&model.RoutingRuleHistory{},
		&model.CollectionEntityHistory{},
		&model.RoutingLog{},

		// System State
		&model.SystemState{},
	)
}

// cleanTables truncates all tables to ensure test isolation.
// Uses mutex to prevent concurrent deadlocks and timeout for safety.
func cleanTables(t *testing.T, db *gorm.DB) {
	t.Helper()

	// Acquire mutex to prevent concurrent TRUNCATE operations
	cleanMutex.Lock()
	defer cleanMutex.Unlock()

	// Add timeout context to prevent indefinite blocking
	// Increased to 120 seconds to handle large number of tables
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	tables := []string{
		// Dispute Templates & Chat Snapshots
		"chat_snapshots", "dispute_templates",
		// Order Timeout
		"order_service_assignments", "order_timeout_logs", "order_timeout_configs",
		// Player Rank & Certification
		"player_certifications", "player_rank_records", "game_ranks",
		// Referral
		"referral_rewards", "referrals", "referral_codes", "referral_configs",
		// Team
		"team_invites", "team_members", "teams",
		// Activity
		"activity_daily_stats", "activity_participations", "activity_rewards", "activities",
		// Recharge
		"recharge_records", "recharge_options",
		// Coupon & VIP
		"coupons", "coupon_templates", "vip_configs", "vip_levels",
		// Commission
		"monthly_settlements", "commission_records", "commission_rules",
		// Settlement
		"player_company_assignments", "settlement_companies",
		// User Block
		"user_blocks",
		// Favorite
		"favorites",
		// User Tag
		"user_tag_relations", "user_tags",
		// Ranking
		"ranking_rewards", "ranking_commission_configs", "player_rankings",
		// Feed/Content
		"feed_images", "feed_reports", "feeds", "content_categories",
		// Collection Entity & Routing Rule
		"routing_logs", "routing_rule_histories", "collection_entity_histories",
		"payment_channel_configs", "routing_rules", "collection_entities",
		// KPI & Statistics
		"kpi_targets",
		"user_statistics", "player_statistics", "platform_statistics",
		"service_item_statistics", "game_statistics",
		// Review & Dispute
		"order_disputes", "review_appeals", "review_replies", "reviews",
		// Chat
		"chat_reports", "chat_messages", "chat_group_members", "chat_groups",
		// Notification
		"user_notifications", "notification_events",
		// Sensitive Word
		"sensitive_words",
		// RBAC & Audit Logs
		"permission_audit_logs", "operation_logs",
		"role_permissions", "user_roles", "menus", "permissions", "roles",
		// Payment & Wallet
		"withdraws", "wallets", "payments",
		// Order
		"order_players", "order_items", "orders",
		// Core
		"service_items", "games", "game_categories",
		"players", "users",
	}

	// Disable foreign key checks, truncate, then re-enable
	if err := db.WithContext(ctx).Exec("SET session_replication_role = 'replica'").Error; err != nil {
		t.Logf("Warning: failed to disable foreign key checks: %v", err)
	}

	for _, table := range tables {
		if err := db.WithContext(ctx).Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			// Log warning but don't fail - table might not exist yet
			t.Logf("Warning: failed to truncate %s: %v", table, err)
		}
	}

	if err := db.WithContext(ctx).Exec("SET session_replication_role = 'origin'").Error; err != nil {
		t.Logf("Warning: failed to re-enable foreign key checks: %v", err)
	}
}

// CreateTestUser creates a test user and returns it.
func CreateTestUser(t *testing.T, db *gorm.DB, name string) *model.User {
	t.Helper()
	user := &model.User{
		Base: model.Base{
			ExtJSON: "{}",
		},
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
		Base: model.Base{
			ExtJSON: "{}",
		},
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
		Base: model.Base{
			ExtJSON: "{}",
		},
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
		Base: model.Base{
			ExtJSON: "{}",
		},
		OrderNo:         fmt.Sprintf("TEST%d", user.ID),
		UserID:          user.ID,
		PlayerID:        &player.ID,
		TotalPriceCents: 10000,
		Status:          status,
		Currency:        model.CurrencyCNY,
		OrderConfig:     "{}", // JSON field
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
		Base: model.Base{
			ExtJSON: "{}",
		},
		OrderID:     order.ID,
		UserID:      order.UserID,
		AmountCents: order.TotalPriceCents,
		Method:      model.PaymentMethodWeChat,
		Status:      status,
		ProviderRaw: []byte("{}"), // JSON field - json.RawMessage
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

// CreateTestWallet creates a test wallet for a user.
func CreateTestWallet(t *testing.T, db *gorm.DB, userID uint64, balanceCents int64) *model.Wallet {
	t.Helper()
	wallet := &model.Wallet{
		Base: model.Base{
			ExtJSON: "{}",
		},
		UserID:       userID,
		BalanceCents: balanceCents,
		FrozenCents:  0,
	}
	if err := db.Create(wallet).Error; err != nil {
		t.Fatalf("Failed to create test wallet: %v", err)
	}
	return wallet
}

// CreateTestServiceItem creates a test service item.
func CreateTestServiceItem(t *testing.T, db *gorm.DB, game *model.Game, name string, priceCents int64) *model.ServiceItem {
	t.Helper()
	var gameID *uint64
	if game != nil {
		gameID = &game.ID
	}
	// Use shorter ItemCode to fit varchar(32) limit
	shortCode := fmt.Sprintf("IT%d", time.Now().UnixNano()%1000000000)
	item := &model.ServiceItem{
		ItemCode:       shortCode,
		Name:           name,
		Category:       "escort",
		GameID:         gameID,
		BasePriceCents: priceCents,
		CommissionRate: 0.20,
		IsActive:       true,
		Tags:           "[]", // JSON array field
	}
	if err := db.Create(item).Error; err != nil {
		t.Fatalf("Failed to create test service item: %v", err)
	}
	return item
}

// CreateTestReview creates a test review.
func CreateTestReview(t *testing.T, db *gorm.DB, order *model.Order, score model.Rating) *model.Review {
	t.Helper()
	review := &model.Review{
		Base: model.Base{
			ExtJSON: "{}",
		},
		OrderID:  order.ID,
		UserID:   order.UserID,
		PlayerID: *order.PlayerID,
		Score:    score,
		Content:  "Test review content",
		Status:   model.ReviewStatusApproved,
		Images:   model.StringArray{}, // JSON array field
	}
	if err := db.Create(review).Error; err != nil {
		t.Fatalf("Failed to create test review: %v", err)
	}
	return review
}

// CreateTestVipLevel creates a test VIP level.
func CreateTestVipLevel(t *testing.T, db *gorm.DB, slug string, expRequired int64) *model.VipLevel {
	t.Helper()
	level := &model.VipLevel{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Slug:          slug,
		Title:         "VIP " + slug,
		ExpRequired:   expRequired,
		OrderDiscount: 0.95,
		IsActive:      true,
		Benefits:      "{}", // JSON field
	}
	if err := db.Create(level).Error; err != nil {
		t.Fatalf("Failed to create test VIP level: %v", err)
	}
	return level
}

// CreateTestCouponTemplate creates a test coupon template.
func CreateTestCouponTemplate(t *testing.T, db *gorm.DB, name string, deductCents int64) *model.CouponTemplate {
	t.Helper()
	template := &model.CouponTemplate{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Name:              name,
		Type:              model.CouponTypeDeduct,
		Source:            model.CouponSourceManual,
		MinAmountCents:    0,
		DeductAmountCents: deductCents,
		Scope:             model.CouponScopeAll,
		ValidityType:      "days",
		ValidityDays:      30,
		TotalCount:        100,
		PerUserLimit:      1,
		IsActive:          true,
		GameIDs:           "[]", // JSON array field
		ItemIDs:           "[]", // JSON array field
	}
	if err := db.Create(template).Error; err != nil {
		t.Fatalf("Failed to create test coupon template: %v", err)
	}
	return template
}

// CreateTestCommissionRule creates a test commission rule.
func CreateTestCommissionRule(t *testing.T, db *gorm.DB, ruleType model.CommissionRuleType, rate int) *model.CommissionRule {
	t.Helper()
	rule := &model.CommissionRule{
		Name:     fmt.Sprintf("Test Rule %d%%", rate),
		Type:     ruleType,
		Rate:     rate,
		IsActive: true,
	}
	if err := db.Create(rule).Error; err != nil {
		t.Fatalf("Failed to create test commission rule: %v", err)
	}
	return rule
}

// CreateTestSettlementCompany creates a test settlement company.
func CreateTestSettlementCompany(t *testing.T, db *gorm.DB, name string) *model.SettlementCompany {
	t.Helper()
	// Create a user first for CreatedBy field
	adminUser := CreateUniqueTestUser(t, db, "admin_"+name)
	company := &model.SettlementCompany{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Name:         name,
		CreditCode:   fmt.Sprintf("91110000%010d", time.Now().UnixNano()%10000000000), // 18 chars
		BankName:     "Test Bank",
		BankAccount:  fmt.Sprintf("ACCT%d", time.Now().UnixNano()),
		ContactName:  "Test Contact",
		ContactPhone: "13800138000",
		Status:       model.CompanyStatusActive,
		CreatedBy:    adminUser.ID,
	}
	if err := db.Create(company).Error; err != nil {
		t.Fatalf("Failed to create test settlement company: %v", err)
	}
	return company
}

// CreateTestWithdraw creates a test withdraw record.
func CreateTestWithdraw(t *testing.T, db *gorm.DB, player *model.Player, amountCents int64, status model.WithdrawStatus) *model.Withdraw {
	t.Helper()
	withdraw := &model.Withdraw{
		PlayerID:    player.ID,
		AmountCents: amountCents,
		Method:      model.WithdrawMethodAlipay,
		Status:      status,
	}
	if err := db.Create(withdraw).Error; err != nil {
		t.Fatalf("Failed to create test withdraw: %v", err)
	}
	return withdraw
}

// CreateTestOrderWithDetails creates a complete test order with all related data.
func CreateTestOrderWithDetails(t *testing.T, db *gorm.DB, user *model.User, player *model.Player, game *model.Game, status model.OrderStatus, priceCents int64) *model.Order {
	t.Helper()

	// Create service item if game provided
	var itemID uint64
	if game != nil {
		item := CreateTestServiceItem(t, db, game, "Test Service", priceCents)
		itemID = item.ID
	}

	now := time.Now()
	scheduledStart := now.Add(time.Hour)
	scheduledEnd := now.Add(2 * time.Hour)

	order := &model.Order{
		Base: model.Base{
			ExtJSON: "{}",
		},
		OrderNo:           fmt.Sprintf("ORD%d%d", user.ID, time.Now().UnixNano()),
		UserID:            user.ID,
		PlayerID:          &player.ID,
		ItemID:            itemID,
		Quantity:          1,
		UnitPriceCents:    priceCents,
		TotalPriceCents:   priceCents,
		CommissionCents:   priceCents * 20 / 100,
		PlayerIncomeCents: priceCents * 80 / 100,
		Currency:          model.CurrencyCNY,
		Status:            status,
		Title:             "Test Order",
		ScheduledStart:    &scheduledStart,
		ScheduledEnd:      &scheduledEnd,
		OrderConfig:       "{}", // JSON field
	}

	if game != nil {
		order.GameID = &game.ID
	}

	if status == model.OrderStatusCompleted {
		order.CompletedAt = &now
	}

	if err := db.Create(order).Error; err != nil {
		t.Fatalf("Failed to create test order: %v", err)
	}
	return order
}

// CreateUniqueTestUser creates a test user with unique phone and email.
func CreateUniqueTestUser(t *testing.T, db *gorm.DB, prefix string) *model.User {
	t.Helper()
	ts := time.Now().UnixNano()
	user := &model.User{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Name:   fmt.Sprintf("%s_%d", prefix, ts),
		Phone:  fmt.Sprintf("138%011d", ts%100000000000),
		Email:  fmt.Sprintf("%s_%d@test.com", prefix, ts),
		Role:   model.RoleUser,
		Status: model.UserStatusActive,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

// CreateTestUserWithPassword creates a test user with password hash.
func CreateTestUserWithPassword(t *testing.T, db *gorm.DB, prefix, password string) *model.User {
	t.Helper()
	ts := time.Now().UnixNano()
	// Use bcrypt to hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	user := &model.User{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Name:         fmt.Sprintf("%s_%d", prefix, ts),
		Phone:        fmt.Sprintf("138%011d", ts%100000000000),
		Email:        fmt.Sprintf("%s_%d@test.com", prefix, ts),
		PasswordHash: string(hashedPassword),
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

// CreateTestRole creates a test role.
func CreateTestRole(t *testing.T, db *gorm.DB, slug, name string) *model.RoleModel {
	t.Helper()
	role := &model.RoleModel{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Slug:        slug,
		Name:        name,
		Description: "Test role",
		IsSystem:    false,
	}
	if err := db.Create(role).Error; err != nil {
		t.Fatalf("Failed to create test role: %v", err)
	}
	return role
}

// CreateTestPermission creates a test permission.
func CreateTestPermission(t *testing.T, db *gorm.DB, method, path, code string) *model.Permission {
	t.Helper()
	permission := &model.Permission{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Method:      model.HTTPMethod(method),
		Path:        path,
		Code:        code,
		Group:       "test",
		Description: "Test permission",
	}
	if err := db.Create(permission).Error; err != nil {
		t.Fatalf("Failed to create test permission: %v", err)
	}
	return permission
}

// CreateTestDispute creates a test dispute.
func CreateTestDispute(t *testing.T, db *gorm.DB, order *model.Order, initiatorID uint64, initiatorType model.DisputeInitiatorType) *model.OrderDispute {
	t.Helper()
	dispute := &model.OrderDispute{
		Base: model.Base{
			ExtJSON: "{}",
		},
		OrderID:       order.ID,
		InitiatorID:   initiatorID,
		InitiatorType: initiatorType,
		Type:          model.DisputeTypeServiceQuality,
		Status:        model.DisputeStatusPending,
		Reason:        "Test dispute reason",
		EvidenceURLs:  model.EvidenceURLArray{}, // JSON array field
	}
	if err := db.Create(dispute).Error; err != nil {
		t.Fatalf("Failed to create test dispute: %v", err)
	}
	return dispute
}

// CreateTestNotification creates a test notification.
func CreateTestNotification(t *testing.T, db *gorm.DB, userID uint64, title, content string) *model.UserNotification {
	t.Helper()
	notification := &model.UserNotification{
		Base: model.Base{
			ExtJSON: "{}",
		},
		UserID:  userID,
		Type:    model.NotificationTypeSystem,
		Channel: model.NotificationChannelInApp,
		Title:   title,
		Content: content,
		Status:  model.NotificationStatusPending,
	}
	if err := db.Create(notification).Error; err != nil {
		t.Fatalf("Failed to create test notification: %v", err)
	}
	return notification
}

// CreateTestChatGroup creates a test chat group.
func CreateTestChatGroup(t *testing.T, db *gorm.DB, name string, groupType model.ChatGroupType, orderID *uint64) *model.ChatGroup {
	t.Helper()
	group := &model.ChatGroup{
		Base: model.Base{
			ExtJSON: "{}",
		},
		GroupName:      name,
		GroupType:      groupType,
		RelatedOrderID: orderID,
		MaxMembers:     100,
		IsActive:       true,
		Settings:       "{}", // JSON field
	}
	if err := db.Create(group).Error; err != nil {
		t.Fatalf("Failed to create test chat group: %v", err)
	}
	return group
}

// CreateTestChatMessage creates a test chat message.
func CreateTestChatMessage(t *testing.T, db *gorm.DB, groupID, senderID uint64, content string) *model.ChatMessage {
	t.Helper()
	message := &model.ChatMessage{
		Base: model.Base{
			ExtJSON: "{}",
		},
		GroupID:     groupID,
		SenderID:    senderID,
		Content:     content,
		MessageType: model.ChatMessageTypeText,
		AuditStatus: model.ChatMessageAuditApproved,
		Metadata:    "{}", // JSON field
	}
	if err := db.Create(message).Error; err != nil {
		t.Fatalf("Failed to create test chat message: %v", err)
	}
	return message
}

// CreateTestSensitiveWord creates a test sensitive word.
func CreateTestSensitiveWord(t *testing.T, db *gorm.DB, word string, category model.SensitiveWordCategory) *model.SensitiveWord {
	t.Helper()
	sw := &model.SensitiveWord{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Word:        word,
		Category:    category,
		MatchType:   model.SensitiveWordMatchTypeExact,
		Severity:    model.SensitiveWordSeverityHigh,
		Replacement: "***",
		IsActive:    true,
	}
	if err := db.Create(sw).Error; err != nil {
		t.Fatalf("Failed to create test sensitive word: %v", err)
	}
	return sw
}

// CreateTestMenu creates a test menu.
func CreateTestMenu(t *testing.T, db *gorm.DB, name, path string, parentID *uint64) *model.Menu {
	t.Helper()
	menu := &model.Menu{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Name:     name,
		Path:     path,
		ParentID: parentID,
		Order:    0,
		Hidden:   false,
	}
	if err := db.Create(menu).Error; err != nil {
		t.Fatalf("Failed to create test menu: %v", err)
	}
	return menu
}

// CreateTestMenuWithOrder creates a test menu with specified order and hidden status.
func CreateTestMenuWithOrder(t *testing.T, db *gorm.DB, name, path string, parentID *uint64, order int, hidden bool) *model.Menu {
	t.Helper()
	menu := &model.Menu{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Name:     name,
		Path:     path,
		ParentID: parentID,
		Order:    order,
		Hidden:   hidden,
	}
	if err := db.Create(menu).Error; err != nil {
		t.Fatalf("Failed to create test menu: %v", err)
	}
	return menu
}

// CreateTestChatGroupMember creates a test chat group member.
func CreateTestChatGroupMember(t *testing.T, db *gorm.DB, groupID, userID uint64, role model.ChatMemberRole) *model.ChatGroupMember {
	t.Helper()
	now := time.Now()
	member := &model.ChatGroupMember{
		Base: model.Base{
			ExtJSON: "{}",
		},
		GroupID:  groupID,
		UserID:   userID,
		Role:     role,
		JoinedAt: now,
		IsActive: true,
	}
	if err := db.Create(member).Error; err != nil {
		t.Fatalf("Failed to create test chat group member: %v", err)
	}
	return member
}

// AssignRoleToUser assigns a role to a user via the user_roles join table.
func AssignRoleToUser(t *testing.T, db *gorm.DB, userID, roleID uint64) {
	t.Helper()
	if err := db.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", userID, roleID).Error; err != nil {
		t.Fatalf("Failed to assign role to user: %v", err)
	}
}

// AssignPermissionToRole assigns a permission to a role via the role_permissions join table.
func AssignPermissionToRole(t *testing.T, db *gorm.DB, roleID, permissionID uint64) {
	t.Helper()
	if err := db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", roleID, permissionID).Error; err != nil {
		t.Fatalf("Failed to assign permission to role: %v", err)
	}
}

// CreateTestGameRank creates a test game rank.
func CreateTestGameRank(t *testing.T, db *gorm.DB, game *model.Game, name string, level int, priceCents int64) *model.GameRank {
	t.Helper()
	rank := &model.GameRank{
		Base: model.Base{
			ExtJSON: "{}",
		},
		GameID:     game.ID,
		Name:       name,
		Level:      level,
		PriceCents: priceCents,
		IsActive:   true,
	}
	if err := db.Create(rank).Error; err != nil {
		t.Fatalf("Failed to create test game rank: %v", err)
	}
	return rank
}

// CreateTestPlayerRankRecord creates a test player rank record.
func CreateTestPlayerRankRecord(t *testing.T, db *gorm.DB, player *model.Player, game *model.Game, rank *model.GameRank, status model.PlayerRankStatus) *model.PlayerRankRecord {
	t.Helper()
	record := &model.PlayerRankRecord{
		Base: model.Base{
			ExtJSON: "{}",
		},
		PlayerID:       player.ID,
		GameID:         game.ID,
		RankID:         rank.ID,
		Status:         status,
		ScreenshotURLs: "[]",
	}
	if err := db.Create(record).Error; err != nil {
		t.Fatalf("Failed to create test player rank record: %v", err)
	}
	return record
}

// CreateTestPlayerCertification creates a test player certification.
func CreateTestPlayerCertification(t *testing.T, db *gorm.DB, player *model.Player, status model.CertificationStatus) *model.PlayerCertification {
	t.Helper()
	cert := &model.PlayerCertification{
		Base: model.Base{
			ExtJSON: "{}",
		},
		PlayerID:       player.ID,
		RealName:       "Test Name",
		IDCardNo:       "123456789012345678",
		IDCardFrontURL: "https://example.com/front.jpg",
		IDCardBackURL:  "https://example.com/back.jpg",
		Status:         status,
	}
	if err := db.Create(cert).Error; err != nil {
		t.Fatalf("Failed to create test player certification: %v", err)
	}
	return cert
}

// CreateTestOrderTimeoutConfig creates a test order timeout config.
func CreateTestOrderTimeoutConfig(t *testing.T, db *gorm.DB, key, value, description string) *model.OrderTimeoutConfig {
	t.Helper()
	config := &model.OrderTimeoutConfig{
		Base: model.Base{
			ExtJSON: "{}",
		},
		ConfigKey:   key,
		ConfigValue: value,
		Description: description,
	}
	if err := db.Create(config).Error; err != nil {
		t.Fatalf("Failed to create test order timeout config: %v", err)
	}
	return config
}

// CreateTestOrderTimeoutLog creates a test order timeout log.
func CreateTestOrderTimeoutLog(t *testing.T, db *gorm.DB, order *model.Order, timeoutType model.OrderTimeoutType, action model.OrderTimeoutAction) *model.OrderTimeoutLog {
	t.Helper()
	log := &model.OrderTimeoutLog{
		Base: model.Base{
			ExtJSON: "{}",
		},
		OrderID:     order.ID,
		TimeoutType: timeoutType,
		TimeoutAt:   time.Now(),
		Action:      action,
		Remark:      "Test timeout log",
	}
	if err := db.Create(log).Error; err != nil {
		t.Fatalf("Failed to create test order timeout log: %v", err)
	}
	return log
}

// CreateTestOrderServiceAssignment creates a test order service assignment.
func CreateTestOrderServiceAssignment(t *testing.T, db *gorm.DB, order *model.Order, serviceUser *model.User, status model.ServiceAssignmentStatus) *model.OrderServiceAssignment {
	t.Helper()
	assignment := &model.OrderServiceAssignment{
		Base: model.Base{
			ExtJSON: "{}",
		},
		OrderID:       order.ID,
		ServiceUserID: serviceUser.ID,
		Status:        status,
		AssignedAt:    time.Now(),
		AssignType:    "auto",
	}
	if err := db.Create(assignment).Error; err != nil {
		t.Fatalf("Failed to create test order service assignment: %v", err)
	}
	return assignment
}

// CreateTestChatSnapshot creates a test chat snapshot for dispute.
func CreateTestChatSnapshot(t *testing.T, db *gorm.DB, disputeID, orderID, chatGroupID uint64) *model.ChatSnapshot {
	t.Helper()
	snapshot := &model.ChatSnapshot{
		Base: model.Base{
			ExtJSON: "{}",
		},
		DisputeID:   disputeID,
		OrderID:     orderID,
		ChatGroupID: chatGroupID,
		Messages:    `[{"sender":"user1","content":"test message","time":"2025-01-01T00:00:00Z"}]`,
		SnapshotAt:  time.Now(),
	}
	if err := db.Create(snapshot).Error; err != nil {
		t.Fatalf("Failed to create test chat snapshot: %v", err)
	}
	return snapshot
}

// CreateTestDisputeTemplate creates a test dispute template.
func CreateTestDisputeTemplate(t *testing.T, db *gorm.DB, code, name string, initiatorType model.DisputeInitiatorType) *model.DisputeTemplate {
	t.Helper()
	template := &model.DisputeTemplate{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Code:          code,
		Name:          name,
		InitiatorType: initiatorType,
		Description:   "Test dispute template",
		SortOrder:     0,
		IsActive:      true,
	}
	if err := db.Create(template).Error; err != nil {
		t.Fatalf("Failed to create test dispute template: %v", err)
	}
	return template
}

// CreateTestCommissionRecord creates a test commission record.
func CreateTestCommissionRecord(t *testing.T, db *gorm.DB, orderID, playerID uint64, totalCents int64, status model.SettlementStatus) *model.CommissionRecord {
	t.Helper()
	record := &model.CommissionRecord{
		OrderID:           orderID,
		PlayerID:          playerID,
		TotalAmountCents:  totalCents,
		CommissionRate:    20,
		CommissionCents:   totalCents * 20 / 100,
		PlayerIncomeCents: totalCents * 80 / 100,
		SettlementStatus:  status,
		SettlementMonth:   time.Now().Format("2006-01"),
	}
	if err := db.Create(record).Error; err != nil {
		t.Fatalf("Failed to create test commission record: %v", err)
	}
	return record
}

// CreateTestOrderItem creates a test order item.
func CreateTestOrderItem(t *testing.T, db *gorm.DB, order *model.Order, slot int, priceCents int64, status model.OrderItemStatus) *model.OrderItem {
	t.Helper()
	item := &model.OrderItem{
		Base: model.Base{
			ExtJSON: "{}",
		},
		OrderID:        order.ID,
		ItemID:         order.ItemID,
		Slot:           slot,
		UnitPriceCents: priceCents,
		Quantity:       1,
		TotalCents:     priceCents,
		CommissionRate: 0.20,
		Status:         status,
	}
	if err := db.Create(item).Error; err != nil {
		t.Fatalf("Failed to create test order item: %v", err)
	}
	return item
}

// CreateTestOrderPlayer creates a test order player record.
func CreateTestOrderPlayer(t *testing.T, db *gorm.DB, order *model.Order, orderItem *model.OrderItem, player *model.Player, status model.OrderPlayerStatus) *model.OrderPlayer {
	t.Helper()
	record := &model.OrderPlayer{
		Base: model.Base{
			ExtJSON: "{}",
		},
		OrderID:     order.ID,
		OrderItemID: orderItem.ID,
		PlayerID:    player.ID,
		JoinedAt:    time.Now(),
		Status:      status,
	}
	if err := db.Create(record).Error; err != nil {
		t.Fatalf("Failed to create test order player: %v", err)
	}
	return record
}

// CreateTestTeam creates a test team.
func CreateTestTeam(t *testing.T, db *gorm.DB, name string, leaderID uint64) *model.Team {
	t.Helper()
	team := &model.Team{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Name:            name,
		LeaderID:        leaderID,
		Status:          model.TeamStatusActive,
		MaxMembers:      5,
		MemberCount:     1,
		IncomeShareType: "equal",
	}
	if err := db.Create(team).Error; err != nil {
		t.Fatalf("Failed to create test team: %v", err)
	}
	return team
}

// CreateTestCollectionEntity creates a test collection entity for routing rule tests.
func CreateTestCollectionEntity(t *testing.T, db *gorm.DB, name string) *model.CollectionEntity {
	t.Helper()
	adminUser := CreateUniqueTestUser(t, db, "admin_entity_"+name)
	entity := &model.CollectionEntity{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Name:       name,
		CreditCode: fmt.Sprintf("91110000%010d", time.Now().UnixNano()%10000000000),
		Status:     model.EntityStatusActive,
		CreatedBy:  adminUser.ID,
	}
	if err := db.Create(entity).Error; err != nil {
		t.Fatalf("Failed to create test collection entity: %v", err)
	}
	return entity
}

// CreateTestRoutingRule creates a test routing rule.
func CreateTestRoutingRule(t *testing.T, db *gorm.DB, entity *model.CollectionEntity, priority int) *model.RoutingRule {
	t.Helper()
	adminUser := CreateUniqueTestUser(t, db, "admin_rule")

	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rule := &model.RoutingRule{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Name:           fmt.Sprintf("Test Rule %d", priority),
		Priority:       priority,
		Conditions:     conditionsJSON,
		TargetEntityID: entity.ID,
		Status:         model.RuleStatusActive,
		Description:    "Test routing rule",
		CreatedBy:      adminUser.ID,
	}
	if err := db.Create(rule).Error; err != nil {
		t.Fatalf("Failed to create test routing rule: %v", err)
	}
	return rule
}

// CreateTestUserTag creates a test user tag.
func CreateTestUserTag(t *testing.T, db *gorm.DB, name, color, description string) *model.UserTag {
	t.Helper()
	adminUser := CreateUniqueTestUser(t, db, "admin_tag_"+name)
	tag := &model.UserTag{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Name:        name,
		Color:       color,
		Description: description,
		CreatedBy:   adminUser.ID,
	}
	if err := db.Create(tag).Error; err != nil {
		t.Fatalf("Failed to create test user tag: %v", err)
	}
	return tag
}

// AssignTagToUser assigns a tag to a user via the join table.
func AssignTagToUser(t *testing.T, db *gorm.DB, userID, tagID uint64) {
	t.Helper()
	relation := &model.UserTagRelation{
		Base: model.Base{
			ExtJSON: "{}",
		},
		UserID: userID,
		TagID:  tagID,
	}
	if err := db.Create(relation).Error; err != nil {
		t.Fatalf("Failed to assign tag to user: %v", err)
	}
}

// CreateTestContentCategory creates a test content category.
func CreateTestContentCategory(t *testing.T, db *gorm.DB, name string) *model.ContentCategory {
	t.Helper()
	category := &model.ContentCategory{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Name:        name,
		Description: "Test category description",
		SortOrder:   0,
		Status:      model.ContentCategoryStatusActive,
		IconURL:     "https://example.com/icon.png",
	}
	if err := db.Create(category).Error; err != nil {
		t.Fatalf("Failed to create test content category: %v", err)
	}
	return category
}

// CreateTestFeed creates a test feed with optional category.
func CreateTestFeed(t *testing.T, db *gorm.DB, authorID uint64, content string, status model.FeedModerationStatus, categoryID *uint64) *model.Feed {
	t.Helper()
	feed := &model.Feed{
		Base: model.Base{
			ExtJSON: "{}",
		},
		AuthorID:         authorID,
		Content:          content,
		CategoryID:       categoryID,
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: status,
		Metrics:          model.FeedMetricFields{},
	}
	if err := db.Create(feed).Error; err != nil {
		t.Fatalf("Failed to create test feed: %v", err)
	}
	return feed
}

// CreateTestFeedWithImages creates a test feed with images.
func CreateTestFeedWithImages(t *testing.T, db *gorm.DB, authorID uint64, content string, status model.FeedModerationStatus, imageUrls []string) *model.Feed {
	t.Helper()
	feed := CreateTestFeed(t, db, authorID, content, status, nil)

	// Create feed images
	for i, url := range imageUrls {
		image := &model.FeedImage{
			Base:      model.Base{ExtJSON: "{}"},
			FeedID:    feed.ID,
			URL:       url,
			Order:     i,
			Width:     800,
			Height:    600,
			SizeBytes: 102400,
		}
		if err := db.Create(image).Error; err != nil {
			t.Fatalf("Failed to create test feed image: %v", err)
		}
	}

	// Reload feed to include images
	db.Preload("Images").First(feed, feed.ID)
	return feed
}

// CreateTestFeedReport creates a test feed report.
func CreateTestFeedReport(t *testing.T, db *gorm.DB, feedID, reporterID uint64, reason string) *model.FeedReport {
	t.Helper()
	report := &model.FeedReport{
		Base: model.Base{
			ExtJSON: "{}",
		},
		FeedID:   feedID,
		Reporter: reporterID,
		Reason:   reason,
		Status:   "pending",
	}
	if err := db.Create(report).Error; err != nil {
		t.Fatalf("Failed to create test feed report: %v", err)
	}
	return report
}
