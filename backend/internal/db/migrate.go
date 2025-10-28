package db

import (
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
	return nil
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
	}
	for _, s := range stmts {
		if err := db.Exec(s).Error; err != nil {
			return err
		}
	}
	return nil
}
