package db

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/pkg/config"
)

// prepareOrdersMigration 在 autoMigrate 之前处理 orders 表的字段迁移（仅 PostgreSQL）
func prepareOrdersMigration(db *gorm.DB) error {
	if db == nil || db.Dialector == nil || db.Dialector.Name() != "postgres" {
		return nil
	}

	// 检查 orders 表是否存在
	var tableExists bool
	checkTableSQL := "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'orders')"
	if err := db.Raw(checkTableSQL).Scan(&tableExists).Error; err != nil {
		return err
	}

	if !tableExists {
		return nil // 表不存在，autoMigrate 会创建
	}

	// 检查并添加 item_id 字段（如果不存在）
	var itemIDExists bool
	checkColumnSQL := "SELECT EXISTS (SELECT FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'orders' AND column_name = 'item_id')"
	if err := db.Raw(checkColumnSQL).Scan(&itemIDExists).Error; err != nil {
		return err
	}

	if !itemIDExists {
		// 先添加字段（允许 NULL），设置默认值为 1（临时默认服务项）
		if err := db.Exec("ALTER TABLE orders ADD COLUMN item_id integer DEFAULT 1").Error; err != nil {
			return err
		}
		// 更新所有现有订单的 item_id 为 1
		if err := db.Exec("UPDATE orders SET item_id = 1 WHERE item_id IS NULL").Error; err != nil {
			return err
		}
	}

	// 检查并添加 order_no 字段（如果不存在）
	var orderNoExists bool
	checkColumnSQL = "SELECT EXISTS (SELECT FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'orders' AND column_name = 'order_no')"
	if err := db.Raw(checkColumnSQL).Scan(&orderNoExists).Error; err != nil {
		return err
	}

	if !orderNoExists {
		// 添加 order_no 字段（允许 NULL）
		if err := db.Exec("ALTER TABLE orders ADD COLUMN order_no text").Error; err != nil {
			return err
		}
	}

	// 检查并添加 unit_price_cents 字段（如果不存在）
	var unitPriceExists bool
	checkColumnSQL = "SELECT EXISTS (SELECT FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'orders' AND column_name = 'unit_price_cents')"
	if err := db.Raw(checkColumnSQL).Scan(&unitPriceExists).Error; err != nil {
		return err
	}

	if !unitPriceExists {
		// 添加字段（默认值为 0）
		if err := db.Exec("ALTER TABLE orders ADD COLUMN unit_price_cents integer DEFAULT 0").Error; err != nil {
			return err
		}
		// 如果有 price_cents 字段，从中迁移数据
		var oldPriceExists bool
		checkColumnSQL = "SELECT EXISTS (SELECT FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'orders' AND column_name = 'price_cents')"
		if err := db.Raw(checkColumnSQL).Scan(&oldPriceExists).Error; err == nil && oldPriceExists {
			if err := db.Exec("UPDATE orders SET unit_price_cents = price_cents WHERE unit_price_cents = 0").Error; err != nil {
				log.Printf("warning: failed to migrate price_cents to unit_price_cents: %v", err)
			}
		}
	}

	// 检查并添加 total_price_cents 字段（如果不存在）
	var totalPriceExists bool
	checkColumnSQL = "SELECT EXISTS (SELECT FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'orders' AND column_name = 'total_price_cents')"
	if err := db.Raw(checkColumnSQL).Scan(&totalPriceExists).Error; err != nil {
		return err
	}

	if !totalPriceExists {
		// 添加字段（默认值为 0）
		if err := db.Exec("ALTER TABLE orders ADD COLUMN total_price_cents integer DEFAULT 0").Error; err != nil {
			return err
		}
		// 如果有 price_cents 字段，从中迁移数据
		var oldPriceExists bool
		checkColumnSQL = "SELECT EXISTS (SELECT FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'orders' AND column_name = 'price_cents')"
		if err := db.Raw(checkColumnSQL).Scan(&oldPriceExists).Error; err == nil && oldPriceExists {
			if err := db.Exec("UPDATE orders SET total_price_cents = price_cents WHERE total_price_cents = 0").Error; err != nil {
				log.Printf("warning: failed to migrate price_cents to total_price_cents: %v", err)
			}
		}
	}

	return nil
}

// migrateVersion 迁移版本号。修改 model 后递增此值，下次启动时触发 AutoMigrate。
const migrateVersion = "2026-02-12-v4"

// isMigrateUpToDate 检查迁移版本是否已是最新，避免每次启动都跑 AutoMigrate。
func isMigrateUpToDate(db *gorm.DB) bool {
	// 先检查核心表是否存在（首次运行时 seed_metadata 还没建）
	var exists bool
	db.Raw(`SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'users')`).Scan(&exists)
	if !exists {
		return false // 首次运行，必须执行迁移
	}

	db.Exec(`CREATE TABLE IF NOT EXISTS seed_metadata (key TEXT PRIMARY KEY, value TEXT)`)
	var val string
	if err := db.Raw(`SELECT value FROM seed_metadata WHERE key = 'migrate_version'`).Scan(&val).Error; err != nil {
		return false
	}
	return val == migrateVersion
}

func markMigrateVersion(db *gorm.DB) {
	db.Exec(`INSERT INTO seed_metadata (key, value) VALUES ('migrate_version', ?) ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value`, migrateVersion)
}

func autoMigrate(db *gorm.DB) error {
	// ── 版本检查：schema 未变则跳过 AutoMigrate ──
	if isMigrateUpToDate(db) {
		log.Printf("[startup] auto-migrate up-to-date (version=%s), skipping", migrateVersion)
		return nil
	}
	log.Printf("[startup] auto-migrate outdated, running migration (version=%s)...", migrateVersion)

	// 先处理 orders 表的特殊字段（PostgreSQL）
	if err := prepareOrdersMigration(db); err != nil {
		return err
	}

	// Phase 1: Create base tables first (tables that others depend on)
	log.Printf("Phase 1: Creating base tables (Game, User, Player, Order, Payment)...")
	if err := db.AutoMigrate(
		&model.Game{},
		&model.User{},
		&model.Player{},
		&model.Order{},
		&model.Payment{},
	); err != nil {
		log.Printf("Phase 1 failed: %v", err)
		return fmt.Errorf("phase 1 migration failed: %w", err)
	}
	log.Printf("Phase 1: Base tables created successfully")

	// Phase 2: Create all other tables (including those with foreign keys to base tables)
	log.Printf("Phase 2: Creating dependent tables...")
	if err := db.AutoMigrate(
		&model.PlayerGame{},
		&model.PlayerSkillTag{},
		&model.Wallet{}, // 用户钱包
		&model.Review{},
		&model.ReviewReport{},
		&model.ReviewDisplaySettings{}, // 评价展示设置
		&model.SensitiveWord{},         // 敏感词
		&model.Withdraw{},
		&model.OperationLog{},
		&model.OrderDispute{}, // Order disputes (must be after Order, Payment, User)
		&model.OrderGroup{},   // 主订单（用户视角）- 订单拆分功能
		// Service Item (统一管理护航服务和礼物)
		&model.ServiceItem{},
		// Order multi-player support
		&model.OrderItem{},   // 订单明细
		&model.OrderPlayer{}, // 订单陪玩师关联
		// VIP system
		&model.VipLevel{},  // VIP等级配置
		&model.VipConfig{}, // VIP系统配置
		// Coupon system
		&model.CouponTemplate{}, // 优惠券模板
		&model.Coupon{},         // 用户优惠券
		// Recharge system
		&model.RechargeOption{}, // 充值档位配置
		&model.RechargeRecord{}, // 充值记录
		// Team system
		&model.Team{},       // 陪玩师团队
		&model.TeamMember{}, // 团队成员
		&model.TeamInvite{}, // 团队邀请
		// Player rank/certification system (陪玩师等级/认证)
		&model.GameRank{},            // 游戏段位配置
		&model.PlayerRankRecord{},    // 陪玩师段位认证记录
		&model.PlayerCertification{}, // 陪玩师实名认证
		&model.PlayerService{},       // 陪玩师服务列表
		&model.PlayerSchedule{},      // 陪玩师排班
		// Referral system (预留)
		&model.ReferralConfig{}, // 推荐配置
		&model.ReferralCode{},   // 邀请码
		&model.Referral{},       // 推荐记录
		&model.ReferralReward{}, // 推荐奖励
		// Activity system
		&model.Activity{},              // 活动
		&model.ActivityReward{},        // 活动奖励配置
		&model.ActivityParticipation{}, // 活动参与记录
		&model.ActivityDailyStats{},    // 活动每日统计
		// Notification system
		&model.NotificationTemplate{},    // 通知模板
		&model.UserNotification{},        // 用户通知（扩展版）
		&model.UserNotificationSetting{}, // 用户通知设置
		&model.NotificationConfig{},      // 通知系统配置
		&model.NotificationSchedule{},    // 定时通知任务
		// Order timeout system (订单超时处理)
		&model.OrderTimeoutConfig{},     // 订单超时配置
		&model.OrderTimeoutLog{},        // 订单超时日志
		&model.OrderServiceAssignment{}, // 订单客服分配记录
		// User block system (用户拉黑)
		&model.UserBlock{}, // 用户拉黑记录
		// Commission models
		&model.CommissionRule{},
		&model.CommissionRecord{},
		&model.MonthlySettlement{},
		// Ranking models
		&model.PlayerRanking{},
		&model.RankingCommissionConfig{},
		&model.RankingReward{},
		// RBAC models
		&model.Permission{},
		&model.RoleModel{},
		&model.RolePermission{},
		&model.UserRole{},
		&model.PermissionAuditLog{}, // Permission audit logs
		// Menu models (动态路由/前端菜单)
		&model.Menu{},
		// Upload model
		&model.Upload{},
		// Chat models
		&model.ChatGroup{},
		&model.ChatGroupMember{},
		&model.ChatMessage{},
		&model.ChatReport{},
		&model.ChatSnapshot{},
		&model.Feed{},
		&model.FeedImage{},
		&model.FeedReport{},
		&model.ContentCategory{}, // 内容分类
		&model.NotificationEvent{},
		&model.ReviewReply{},
		&model.ReviewAppeal{},
		// Monitor and KPI models
		&model.Alert{},
		&model.KPITarget{},
		&model.UserActivityDaily{},
		// Settlement Company models (提现分流)
		&model.SettlementCompany{},
		&model.PlayerCompanyAssignment{},
		&model.SettlementCompanyHistory{},
		// Collection Entity models (收款分流)
		&model.CollectionEntity{},
		&model.PaymentChannelConfig{},
		&model.RoutingRule{},
		&model.RoutingRuleHistory{},
		&model.CollectionEntityHistory{},
		&model.RoutingLog{},
		// Refund records
		&model.RefundRecord{},
		// Statistics models (统计指标)
		&model.UserStatistics{},
		&model.PlayerStatistics{},
		&model.ServiceItemStatistics{},
		&model.GameStatistics{},
		&model.PlatformStatistics{},
		&model.TagThreshold{},
		// System state tracking (for initialization status)
		&model.SystemState{},
		// Game categories
		&model.GameCategory{},
		// LFG (Looking For Group) system
		&model.LFGRequest{},
		// Favorites system
		&model.Favorite{},
		// Player presence system
		&model.PlayerPresence{},
		// Banner system (首页轮播图)
		&model.Banner{},
		// User management (previously a separate migration)
		&model.UserTag{},
		&model.UserTagRelation{},
		&model.UserLoginHistory{},
		&model.UserBehavior{},
		&model.UserSettings{},
	); err != nil {
		log.Printf("Phase 2 failed: %v", err)
		return fmt.Errorf("phase 2 migration failed: %w", err)
	}
	log.Printf("Phase 2: Dependent tables created successfully")

	// 标记迁移版本
	markMigrateVersion(db)
	return nil
}

// MigrateUserManagement 迁移用户管理相关数据表（保留以兼容外部调用，内部已合并到 autoMigrate）
func MigrateUserManagement(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.UserTag{},
		&model.UserTagRelation{},
		&model.UserLoginHistory{},
		&model.UserBehavior{},
	)
}

// runDataFixups contains data migrations that adjust existing values.
// It is idempotent and safe to run at startup.
func runDataFixups(db *gorm.DB) error {
	// Normalize order status spelling: "cancelled" -> "canceled" (legacy British spelling)
	//nolint:misspell // SQL contains legacy 'cancelled' string for data migration
	if err := db.Exec("UPDATE orders SET status='canceled' WHERE status='cancelled'").Error; err != nil {
		return err
	}
	// Generate OrderNo for existing orders without one
	if err := generateOrderNumbers(db); err != nil {
		return err
	}
	// Clamp player rating average to [0,5] and set negative counts to 0
	if err := db.Exec("UPDATE players SET rating_average = CASE WHEN rating_average < 0 THEN 0 WHEN rating_average > 5 THEN 5 ELSE rating_average END").Error; err != nil {
		return err
	}
	if err := db.Exec("UPDATE players SET rating_count = 0 WHERE rating_count < 0").Error; err != nil {
		return err
	}
	// Ensure RBAC default roles exist
	if err := ensureDefaultRoles(db); err != nil {
		return err
	}
	// Ensure default commission rule exists
	if err := ensureDefaultCommissionRule(db); err != nil {
		return err
	}
	return ensureSuperAdmin(db)
}

func ensureIndexes(db *gorm.DB) error {
	stmts := []string{
		// Orders composite indexes
		"CREATE INDEX IF NOT EXISTS idx_orders_status_created ON orders (status, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_orders_user_created ON orders (user_id, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_orders_player_created ON orders (player_id, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_orders_game_created ON orders (game_id, created_at DESC)",
		// Payments composite indexes
		"CREATE INDEX IF NOT EXISTS idx_payments_status_created ON payments (status, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_payments_user_created ON payments (user_id, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_payments_order_created ON payments (order_id, created_at DESC)",
		// Withdraws composite indexes
		"CREATE INDEX IF NOT EXISTS idx_withdraws_status_created ON withdraws (status, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_withdraws_player_created ON withdraws (player_id, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_withdraws_user_created ON withdraws (user_id, created_at DESC)",
		// Service Items indexes
		"CREATE INDEX IF NOT EXISTS idx_service_items_game_subcat ON service_items (game_id, sub_category)",
		"CREATE INDEX IF NOT EXISTS idx_service_items_subcat_active ON service_items (sub_category, is_active)",
		"CREATE INDEX IF NOT EXISTS idx_orders_item_created ON orders (item_id, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_orders_recipient_player ON orders (recipient_player_id, created_at DESC)",
		// Order Group indexes (订单拆分)
		"CREATE INDEX IF NOT EXISTS idx_order_groups_user_status ON order_groups (user_id, status, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_order_groups_game ON order_groups (game_id, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_orders_group_hour ON orders (group_id, hour_index)",
		"CREATE INDEX IF NOT EXISTS idx_orders_transfer ON orders (transfer_from, transfer_to)",
		// Commission indexes
		"CREATE INDEX IF NOT EXISTS idx_commission_records_player_month ON commission_records (player_id, settlement_month)",
		"CREATE INDEX IF NOT EXISTS idx_commission_records_status_month ON commission_records (settlement_status, settlement_month)",
		"CREATE INDEX IF NOT EXISTS idx_monthly_settlements_player_month ON monthly_settlements (player_id, settlement_month)",
		"CREATE INDEX IF NOT EXISTS idx_monthly_settlements_month_status ON monthly_settlements (settlement_month, status)",
		// Operation logs indexes
		"CREATE INDEX IF NOT EXISTS idx_oplogs_entity ON operation_logs (entity_type, entity_id, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_oplogs_actor ON operation_logs (actor_user_id, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_oplogs_action ON operation_logs (action, created_at DESC)",
		// Permission audit logs indexes
		"CREATE INDEX IF NOT EXISTS idx_perm_audit_operator ON permission_audit_logs (operator_id, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_perm_audit_target ON permission_audit_logs (target_type, target_id, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_perm_audit_action ON permission_audit_logs (action, created_at DESC)",
	}
	for _, s := range stmts {
		if err := db.Exec(s).Error; err != nil {
			return err
		}
	}
	return nil
}

// ensureDefaultRoles creates system predefined roles if they don't exist.
func ensureDefaultRoles(db *gorm.DB) error {
	roles := []model.RoleModel{
		{
			Slug:        string(model.RoleSlugSuperAdmin),
			Name:        "超级管理员",
			Description: "拥有系统所有权限，不可删除",
			IsSystem:    true,
		},
		{
			Slug:        string(model.RoleSlugAdmin),
			Name:        "管理员",
			Description: "后台管理权限",
			IsSystem:    true,
		},
		{
			Slug:        string(model.RoleSlugPlayer),
			Name:        "陪玩师",
			Description: "提供陪玩服务的用户",
			IsSystem:    true,
		},
		{
			Slug:        string(model.RoleSlugUser),
			Name:        "普通用户",
			Description: "平台普通用户",
			IsSystem:    true,
		},
	}

	for i := range roles {
		role := &roles[i]
		var existing model.RoleModel
		err := db.Where("slug = ?", role.Slug).First(&existing).Error
		if err == nil {
			// Role exists, update description if needed
			if existing.Name != role.Name || existing.Description != role.Description {
				db.Model(&existing).Updates(map[string]interface{}{
					"name":        role.Name,
					"description": role.Description,
					"is_system":   true,
				})
			}
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		// Create new role
		if err := db.Create(role).Error; err != nil {
			return err
		}
		log.Printf("created system role: %s (id=%d)", role.Slug, role.ID)
	}

	return nil
}

func ensureSuperAdmin(db *gorm.DB) error {
	cfg := config.Load()
	env := strings.TrimSpace(os.Getenv("APP_ENV"))
	email := strings.TrimSpace(cfg.SuperAdmin.Email)
	phone := strings.TrimSpace(cfg.SuperAdmin.Phone)
	name := strings.TrimSpace(cfg.SuperAdmin.Name)
	password := cfg.SuperAdmin.Password

	if name == "" {
		name = "Super Admin"
	}

	if email == "" && phone == "" {
		if env == "production" {
			return errors.New("SUPER_ADMIN_EMAIL or SUPER_ADMIN_PHONE must be set in production")
		}
		log.Printf("⚠️  警告：未配置超级管理员邮箱或手机号，使用示例邮箱 'admin@gamelink.local'")
		log.Printf("   开发环境可以在配置文件中设置 super_admin.email 和 super_admin.password")
		log.Printf("   生产环境必须设置 SUPER_ADMIN_EMAIL、SUPER_ADMIN_PASSWORD 环境变量")
		email = "admin@gamelink.local"
	}

	switch password {
	case "":
		if env == "production" {
			return errors.New("SUPER_ADMIN_PASSWORD must be set in production")
		}
		log.Printf("⚠️  警告：未配置超级管理员密码，使用示例密码")
		log.Printf("   请务必在开发完成后修改默认密码！")
		password = "Admin@123456"
	case "Admin@123456", "123456":
		log.Printf("⚠️  警告：超级管理员正在使用弱密码 '%s'", password)
		log.Printf("   建议修改为包含大小写字母、数字和特殊符号的强密码")
	}

	// Avoid unique constraint conflicts when phone 为空且已有空手机号行
	if phone == "" {
		phone = fmt.Sprintf("superadmin-%d", time.Now().UnixNano())
	}

	lookup := db.Model(&model.User{})
	if email != "" {
		lookup = lookup.Where("email = ?", email)
	} else {
		lookup = lookup.Where("phone = ?", phone)
	}

	var existing model.User
	err := lookup.First(&existing).Error

	var admin *model.User
	userExists := false

	switch {
	case err == nil:
		// User already exists, use existing user
		admin = &existing
		userExists = true
		log.Printf("super admin user already exists: email=%s id=%d", admin.Email, admin.ID)
	case errors.Is(err, gorm.ErrRecordNotFound):
		// Create new admin user
		hashed, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if hashErr != nil {
			return hashErr
		}

		admin = &model.User{
			Name:         name,
			Email:        email,
			Phone:        phone,
			PasswordHash: string(hashed),
			Role:         model.RoleAdmin,
			Status:       model.UserStatusActive,
		}

		if createErr := db.Create(admin).Error; createErr != nil {
			return createErr
		}
		log.Printf("created super admin user: email=%s id=%d", admin.Email, admin.ID)
	default:
		return err
	}

	// Assign super_admin role to this user (whether new or existing)
	var superAdminRole model.RoleModel
	if err := db.Where("slug = ?", model.RoleSlugSuperAdmin).First(&superAdminRole).Error; err != nil {
		return fmt.Errorf("super_admin role not found: %w", err)
	}

	// Check if user already has the role
	var existingUserRole model.UserRole
	err = db.Where("user_id = ? AND role_id = ?", admin.ID, superAdminRole.ID).First(&existingUserRole).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Assign role
		userRole := &model.UserRole{
			UserID: admin.ID,
			RoleID: superAdminRole.ID,
		}
		if err := db.Create(userRole).Error; err != nil {
			return fmt.Errorf("failed to assign super_admin role to user %d: %w", admin.ID, err)
		}
		log.Printf("assigned super_admin role to user id=%d", admin.ID)
	} else if err != nil {
		return fmt.Errorf("failed to check user role: %w", err)
	} else {
		log.Printf("user id=%d already has super_admin role", admin.ID)
	}

	if userExists {
		log.Printf("super admin user ensured (existing): email=%s phone=%s id=%d", email, phone, admin.ID)
	} else {
		log.Printf("super admin user ensured (created): email=%s phone=%s id=%d", email, phone, admin.ID)
	}
	return nil
}

// generateOrderNumbers 为没有订单号的订单生成订单号
func generateOrderNumbers(db *gorm.DB) error {
	// 查询所有没有订单号的订单（空字符串或NULL）
	var orders []model.Order
	if err := db.Where("order_no = ? OR order_no IS NULL OR TRIM(order_no) = ''", "").Find(&orders).Error; err != nil {
		return err
	}

	if len(orders) == 0 {
		return nil
	}

	log.Printf("generating order numbers for %d orders", len(orders))

	// 为每个订单生成唯一订单号
	timestamp := time.Now().Unix()
	for i := range orders {
		// 格式: ORD + 时间戳 + 订单ID (确保唯一性)
		orderNo := fmt.Sprintf("ORD%d%08d", timestamp, orders[i].ID)
		if err := db.Model(&orders[i]).Update("order_no", orderNo).Error; err != nil {
			log.Printf("warning: failed to update order %d: %v", orders[i].ID, err)
			continue // 继续处理其他订单，不要中断整个流程
		}
	}

	log.Printf("successfully generated order numbers for %d orders", len(orders))
	return nil
}

// ensureDefaultCommissionRule 确保默认抽成规则存在
func ensureDefaultCommissionRule(db *gorm.DB) error {
	var existing model.CommissionRule
	err := db.Where("type = ? AND is_active = ?", "default", true).
		Where("game_id IS NULL AND player_id IS NULL AND service_type IS NULL").
		First(&existing).Error

	if err == nil {
		// 默认规则已存在
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 创建默认抽成规则：20%
	defaultRule := &model.CommissionRule{
		Name:        "默认抽成规则",
		Description: "平台默认抽成比例为20%",
		Type:        "default",
		Rate:        20,
		IsActive:    true,
	}

	if err := db.Create(defaultRule).Error; err != nil {
		return err
	}

	log.Printf("created default commission rule: 20%% (id=%d)", defaultRule.ID)
	return nil
}
