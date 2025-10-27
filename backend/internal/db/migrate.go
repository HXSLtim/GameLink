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
	)
}

// runDataFixups contains data migrations that adjust existing values.
// It is idempotent and safe to run at startup.
func runDataFixups(db *gorm.DB) error {
	// Normalize order status spelling: "cancelled" -> "canceled" //nolint:misspell // legacy spelling retained for clarity
	if err := db.Exec("UPDATE orders SET status='canceled' WHERE status='cancelled'").Error; err != nil {
		return err
	}
	return nil
}
