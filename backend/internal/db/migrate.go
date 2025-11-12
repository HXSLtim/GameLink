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

	"gamelink/internal/config"
	"gamelink/internal/model"
)

// prepareOrdersMigration 在 autoMigrate 之前处理 orders 表的字段迁移
func prepareOrdersMigration(db *gorm.DB) error {
	// 检查 orders 表是否存在
	var tableExists bool
	if err := db.Raw("SELECT COUNT(*) > 0 FROM sqlite_master WHERE type='table' AND name='orders'").Scan(&tableExists).Error; err != nil {
		return err
	}

	if !tableExists {
		return nil // 表不存在，autoMigrate 会创建
	}

	// 检查并添加 item_id 字段（如果不存在）
	var itemIDExists bool
	if err := db.Raw("SELECT COUNT(*) > 0 FROM pragma_table_info('orders') WHERE name='item_id'").Scan(&itemIDExists).Error; err != nil {
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
	if err := db.Raw("SELECT COUNT(*) > 0 FROM pragma_table_info('orders') WHERE name='order_no'").Scan(&orderNoExists).Error; err != nil {
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
	if err := db.Raw("SELECT COUNT(*) > 0 FROM pragma_table_info('orders') WHERE name='unit_price_cents'").Scan(&unitPriceExists).Error; err != nil {
		return err
	}

	if !unitPriceExists {
		// 添加字段（默认值为 0）
		if err := db.Exec("ALTER TABLE orders ADD COLUMN unit_price_cents integer DEFAULT 0").Error; err != nil {
			return err
		}
		// 如果有 price_cents 字段，从中迁移数据
		var oldPriceExists bool
		if err := db.Raw("SELECT COUNT(*) > 0 FROM pragma_table_info('orders') WHERE name='price_cents'").Scan(&oldPriceExists).Error; err == nil && oldPriceExists {
			if err := db.Exec("UPDATE orders SET unit_price_cents = price_cents WHERE unit_price_cents = 0").Error; err != nil {
				log.Printf("warning: failed to migrate price_cents to unit_price_cents: %v", err)
			}
		}
	}

	// 检查并添加 total_price_cents 字段（如果不存在）
	var totalPriceExists bool
	if err := db.Raw("SELECT COUNT(*) > 0 FROM pragma_table_info('orders') WHERE name='total_price_cents'").Scan(&totalPriceExists).Error; err != nil {
		return err
	}

	if !totalPriceExists {
		// 添加字段（默认值为 0）
		if err := db.Exec("ALTER TABLE orders ADD COLUMN total_price_cents integer DEFAULT 0").Error; err != nil {
			return err
		}
		// 如果有 price_cents 字段，从中迁移数据
		var oldPriceExists bool
		if err := db.Raw("SELECT COUNT(*) > 0 FROM pragma_table_info('orders') WHERE name='price_cents'").Scan(&oldPriceExists).Error; err == nil && oldPriceExists {
			if err := db.Exec("UPDATE orders SET total_price_cents = price_cents WHERE total_price_cents = 0").Error; err != nil {
				log.Printf("warning: failed to migrate price_cents to total_price_cents: %v", err)
			}
		}
	}

	return nil
}

func autoMigrate(db *gorm.DB) error {
	// 临时禁用外键检查（SQLite）
	db.Exec("PRAGMA foreign_keys = OFF")
	defer db.Exec("PRAGMA foreign_keys = ON")

	// 先处理 orders 表的特殊字段
	if err := prepareOrdersMigration(db); err != nil {
		return err
	}

	return db.AutoMigrate(
		&model.Game{},
		&model.Player{},
		&model.PlayerGame{},
		&model.PlayerSkillTag{},
		&model.User{},
		&model.Order{},
		&model.Payment{},
		&model.Review{},
		&model.Withdraw{},
		&model.OperationLog{},
		// Service Item (统一管理护航服务和礼物)
		&model.ServiceItem{},
		// Commission models
		&model.CommissionRule{},
		&model.CommissionRecord{},
		&model.MonthlySettlement{},
		// Ranking models
		&model.PlayerRanking{},
		&model.RankingCommissionConfig{},
		// RBAC models
		&model.Permission{},
		&model.RoleModel{},
		&model.RolePermission{},
		&model.UserRole{},
		// Upload model
		&model.Upload{},
		// Chat models
		&model.ChatGroup{},
		&model.ChatGroupMember{},
		&model.ChatMessage{},
		&model.ChatReport{},
	)
}

// runDataFixups contains data migrations that adjust existing values.
// It is idempotent and safe to run at startup.
func runDataFixups(db *gorm.DB) error {
	// Normalize order status spelling: "cancelled" -> "canceled" //nolint:misspell // legacy spelling retained for clarity
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
		// Commission indexes
		"CREATE INDEX IF NOT EXISTS idx_commission_records_player_month ON commission_records (player_id, settlement_month)",
		"CREATE INDEX IF NOT EXISTS idx_commission_records_status_month ON commission_records (settlement_status, settlement_month)",
		"CREATE INDEX IF NOT EXISTS idx_monthly_settlements_player_month ON monthly_settlements (player_id, settlement_month)",
		"CREATE INDEX IF NOT EXISTS idx_monthly_settlements_month_status ON monthly_settlements (settlement_month, status)",
		// Operation logs indexes
		"CREATE INDEX IF NOT EXISTS idx_oplogs_entity ON operation_logs (entity_type, entity_id, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_oplogs_actor ON operation_logs (actor_user_id, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_oplogs_action ON operation_logs (action, created_at DESC)",
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
		email = "admin@gamelink.local"
	}

	if password == "" {
		if env == "production" {
			return errors.New("SUPER_ADMIN_PASSWORD must be set in production")
		}
		password = "Admin@123456"
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
	switch {
	case err == nil:
		return nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		// fall through and create admin user
	default:
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &model.User{
		Name:         name,
		Email:        email,
		Phone:        phone,
		PasswordHash: string(hashed),
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
	}

	if err := db.Create(admin).Error; err != nil {
		return err
	}

	// Assign super_admin role to this user
	var superAdminRole model.RoleModel
	if err := db.Where("slug = ?", model.RoleSlugSuperAdmin).First(&superAdminRole).Error; err != nil {
		log.Printf("warning: super_admin role not found, skipping role assignment: %v", err)
	} else {
		// Check if user already has the role
		var existingUserRole model.UserRole
		err := db.Where("user_id = ? AND role_id = ?", admin.ID, superAdminRole.ID).First(&existingUserRole).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Assign role
			userRole := &model.UserRole{
				UserID: admin.ID,
				RoleID: superAdminRole.ID,
			}
			if err := db.Create(userRole).Error; err != nil {
				log.Printf("warning: failed to assign super_admin role: %v", err)
			} else {
				log.Printf("assigned super_admin role to user id=%d", admin.ID)
			}
		}
	}

	log.Printf("super admin user ensured: email=%s phone=%s id=%d", email, phone, admin.ID)
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
