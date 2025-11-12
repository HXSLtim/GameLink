package db

import (
    "strings"
    "testing"

    "github.com/glebarez/sqlite"
    "github.com/stretchr/testify/require"
    "gorm.io/gorm"

    "gamelink/internal/model"
)

func newMemDB(t *testing.T) *gorm.DB {
    t.Helper()
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)
    return db
}

func migrateCore(t *testing.T, db *gorm.DB) {
    t.Helper()
    err := db.AutoMigrate(
        &model.User{},
        &model.RoleModel{},
        &model.UserRole{},
        &model.Order{},
    )
    require.NoError(t, err)
}


func TestEnsureDefaultRoles_UpdateExisting(t *testing.T) {
    db := newMemDB(t)
    migrateCore(t, db)

    r := &model.RoleModel{
        Slug:        string(model.RoleSlugSuperAdmin),
        Name:        "超级管理员",
        Description: "old",
        IsSystem:    false,
    }
    require.NoError(t, db.Create(r).Error)

    require.NoError(t, ensureDefaultRoles(db))

    var got model.RoleModel
    require.NoError(t, db.Where("slug = ?", model.RoleSlugSuperAdmin).First(&got).Error)
    if got.Description == "old" || !got.IsSystem {
        t.Fatalf("expected role updated, got desc=%q is_system=%v", got.Description, got.IsSystem)
    }
}

func TestGenerateOrderNumbers_AssignsForEmpty(t *testing.T) {
    db := newMemDB(t)
    migrateCore(t, db)

    o1 := &model.Order{Title: "t1", OrderNo: ""}
    o2 := &model.Order{Title: "t2"}
    require.NoError(t, db.Create([]*model.Order{o1, o2}).Error)

    require.NoError(t, generateOrderNumbers(db))

    var a, b model.Order
    require.NoError(t, db.First(&a, o1.ID).Error)
    require.NoError(t, db.First(&b, o2.ID).Error)
    if strings.TrimSpace(a.OrderNo) == "" {
        t.Fatal("expected o1 to have generated order_no")
    }
    if strings.TrimSpace(b.OrderNo) == "" {
        t.Fatal("expected o2 to have generated order_no")
    }
}

func TestPrepareOrdersMigration_AddsMissingColumns(t *testing.T) {
    db := newMemDB(t)
    require.NoError(t, db.Exec("CREATE TABLE orders (id integer primary key, user_id integer, status text)").Error)

    require.NoError(t, prepareOrdersMigration(db))

    var exists bool
    require.NoError(t, db.Raw("SELECT COUNT(*) > 0 FROM pragma_table_info('orders') WHERE name='item_id'").Scan(&exists).Error)
    require.True(t, exists)
    require.NoError(t, db.Raw("SELECT COUNT(*) > 0 FROM pragma_table_info('orders') WHERE name='order_no'").Scan(&exists).Error)
    require.True(t, exists)
    require.NoError(t, db.Raw("SELECT COUNT(*) > 0 FROM pragma_table_info('orders') WHERE name='unit_price_cents'").Scan(&exists).Error)
    require.True(t, exists)
    require.NoError(t, db.Raw("SELECT COUNT(*) > 0 FROM pragma_table_info('orders') WHERE name='total_price_cents'").Scan(&exists).Error)
    require.True(t, exists)
}

func TestPrepareOrdersMigration_Idempotent(t *testing.T) {
    db := newMemDB(t)
    // initial create with price_cents for migration path
    require.NoError(t, db.Exec("CREATE TABLE orders (id integer primary key, status text, created_at datetime, price_cents integer)").Error)
    // first run adds new columns and migrates
    require.NoError(t, prepareOrdersMigration(db))
    // insert a row with empty order_no
    require.NoError(t, db.Exec("INSERT INTO orders (id, status, created_at, price_cents) VALUES (1, 'pending', CURRENT_TIMESTAMP, 123)").Error)
    // second run should be no-op (no duplicate errors)
    require.NoError(t, prepareOrdersMigration(db))
    // verify columns exist
    var exists bool
    require.NoError(t, db.Raw("SELECT COUNT(*) > 0 FROM pragma_table_info('orders') WHERE name='unit_price_cents'").Scan(&exists).Error)
    require.True(t, exists)
    require.NoError(t, db.Raw("SELECT COUNT(*) > 0 FROM pragma_table_info('orders') WHERE name='total_price_cents'").Scan(&exists).Error)
    require.True(t, exists)
}


func TestGenerateOrderNumbers_Unique(t *testing.T) {
    db := newMemDB(t)
    require.NoError(t, db.AutoMigrate(&model.Order{}))
    // two orders without numbers
    require.NoError(t, db.Create(&model.Order{Title: "a"}).Error)
    require.NoError(t, db.Create(&model.Order{Title: "b"}).Error)

    require.NoError(t, generateOrderNumbers(db))

    var orders []model.Order
    require.NoError(t, db.Find(&orders).Error)
    require.Len(t, orders, 2)
    if strings.TrimSpace(orders[0].OrderNo) == "" || strings.TrimSpace(orders[1].OrderNo) == "" {
        t.Fatal("expected non-empty order_no for all")
    }
    if orders[0].OrderNo == orders[1].OrderNo {
        t.Fatal("expected unique order_no values")
    }
}

func TestEnsureDefaultCommissionRule_AlreadyExists(t *testing.T) {
    db := newMemDB(t)
    require.NoError(t, db.AutoMigrate(&model.CommissionRule{}))

    r := &model.CommissionRule{
        Name:        "默认抽成规则",
        Description: "20%",
        Type:        "default",
        Rate:        20,
        IsActive:    true,
    }
    require.NoError(t, db.Create(r).Error)

    require.NoError(t, ensureDefaultCommissionRule(db))

    var cnt int64
    require.NoError(t, db.Model(&model.CommissionRule{}).
        Where("type = ? AND is_active = ?", "default", true).
        Where("game_id IS NULL AND player_id IS NULL AND service_type IS NULL").
        Count(&cnt).Error)
    require.Equal(t, int64(1), cnt)
}

func TestEnsureIndexesCreatesExpectedIndexes(t *testing.T) {
    db := newMemDB(t)
    // create minimal tables with referenced columns
    stmts := []string{
        "CREATE TABLE orders (id integer primary key, status text, created_at datetime, user_id integer, player_id integer, game_id integer, item_id integer, recipient_player_id integer)",
        "CREATE TABLE payments (id integer primary key, status text, created_at datetime, user_id integer, order_id integer)",
        "CREATE TABLE withdraws (id integer primary key, status text, created_at datetime, player_id integer, user_id integer)",
        "CREATE TABLE service_items (id integer primary key, game_id integer, sub_category text, is_active boolean)",
        "CREATE TABLE commission_records (id integer primary key, player_id integer, settlement_month text, settlement_status text)",
        "CREATE TABLE monthly_settlements (id integer primary key, player_id integer, settlement_month text, status text)",
        "CREATE TABLE operation_logs (id integer primary key, entity_type text, entity_id integer, actor_user_id integer, action text, created_at datetime)",
    }
    for _, s := range stmts { require.NoError(t, db.Exec(s).Error) }

    require.NoError(t, ensureIndexes(db))

    // sample of expected indexes
    var count int64
    require.NoError(t, db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name='idx_orders_status_created'").Scan(&count).Error)
    require.Equal(t, int64(1), count)
    require.NoError(t, db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name='idx_payments_status_created'").Scan(&count).Error)
    require.Equal(t, int64(1), count)
    require.NoError(t, db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name='idx_oplogs_entity'").Scan(&count).Error)
    require.Equal(t, int64(1), count)
}
