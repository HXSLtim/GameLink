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
	if err == nil {
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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

	log.Printf("super admin user ensured: email=%s phone=%s id=%d", email, phone, admin.ID)
	return nil
}
