package db

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
)

func setupSeedTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(
		&model.Game{},
		&model.User{},
		&model.Player{},
		&model.Order{},
		&model.Payment{},
		&model.Review{},
		&model.ServiceItem{},
	)
	require.NoError(t, err)

	err = db.Create(&model.ServiceItem{
		ItemCode:       "ESCORT_SEED",
		Name:           "Escort Seed",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 1000,
		ServiceHours:   1,
		CommissionRate: 0.2,
		IsActive:       true,
	}).Error
	require.NoError(t, err)

	return db
}

func TestApplySeedsPopulatesData(t *testing.T) {
	t.Parallel()
	db := setupSeedTestDB(t)

	require.NoError(t, applySeeds(db))
	// apply a second time to ensure idempotency code paths are also covered
	require.NoError(t, applySeeds(db))

	var gameCount, userCount, playerCount, orderCount, paymentCount, reviewCount int64

	require.NoError(t, db.Model(&model.Game{}).Count(&gameCount).Error)
	require.NoError(t, db.Model(&model.User{}).Count(&userCount).Error)
	require.NoError(t, db.Model(&model.Player{}).Count(&playerCount).Error)
	require.NoError(t, db.Model(&model.Order{}).Count(&orderCount).Error)
	require.NoError(t, db.Model(&model.Payment{}).Count(&paymentCount).Error)
	require.NoError(t, db.Model(&model.Review{}).Count(&reviewCount).Error)

	require.GreaterOrEqual(t, gameCount, int64(10))
	require.GreaterOrEqual(t, userCount, int64(10))
	require.GreaterOrEqual(t, playerCount, int64(1))
	require.GreaterOrEqual(t, orderCount, int64(1))
	require.GreaterOrEqual(t, paymentCount, int64(1))
	require.GreaterOrEqual(t, reviewCount, int64(1))
}

func TestSeedPointerHelpers(t *testing.T) {
	now := time.Now()
	d := time.Hour

	require.Equal(t, now, *ptrTime(now))
	require.Equal(t, d, *ptrDuration(d))
	require.Nil(t, ptrTimeWithOffset(now, nil))

	offset := ptrDuration(30 * time.Minute)
	require.WithinDuration(t, now.Add(*offset), *ptrTimeWithOffset(now, offset), time.Millisecond)
}
