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
)

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Game{},
		&model.Player{},
		&model.PlayerGame{},
		&model.PlayerSkillTag{},
		&model.User{},
		&model.Order{},
		&model.Payment{},
		&model.Review{},
		&model.OperationLog{},
		// RBAC models
		&model.Permission{},
		&model.RoleModel{},
		&model.RolePermission{},
		&model.UserRole{},
	)
}

// runDataFixups contains data migrations that adjust existing values.
// It is idempotent and safe to run at startup.
func runDataFixups(db *gorm.DB) error {
	// Normalize order status spelling: "cancelled" -> "canceled" //nolint:misspell // legacy spelling retained for clarity
	if err := db.Exec("UPDATE orders SET status='canceled' WHERE status='cancelled'").Error; err != nil {
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
	env := strings.TrimSpace(os.Getenv("APP_ENV"))
	email := strings.TrimSpace(os.Getenv("SUPER_ADMIN_EMAIL"))
	phone := strings.TrimSpace(os.Getenv("SUPER_ADMIN_PHONE"))
	name := strings.TrimSpace(os.Getenv("SUPER_ADMIN_NAME"))
	password := os.Getenv("SUPER_ADMIN_PASSWORD")

	if name == "" {
		name = "Super Admin"
	}

	if email == "" && phone == "" {
		if env == "production" {
			return errors.New("SUPER_ADMIN_EMAIL or SUPER_ADMIN_PHONE must be set in production")
		}
		email = "superAdmin@GameLink.com"
	}

	if password == "" {
		if env == "production" {
			return errors.New("SUPER_ADMIN_PASSWORD must be set in production")
		}
		password = "admin123"
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
