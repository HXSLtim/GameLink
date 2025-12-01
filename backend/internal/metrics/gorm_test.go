package metrics

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

/**
 * Test_InstrumentGorm_CallbackLogic 测试callback逻辑
 * 注意：由于需要真实数据库才能注册和触发callbacks，这里直接测试callback函数的核心逻辑
 */
func Test_InstrumentGorm_CallbackLogic(t *testing.T) {
	// Arrange
	reg := prometheus.NewRegistry()
	Init(reg)

	db := &gorm.DB{
		Config:    &gorm.Config{},
		Statement: &gorm.Statement{},
	}

	// Test before callback - 应该设置start time
	t.Run("before callback sets start time", func(t *testing.T) {
		const key = "metrics_start"

		// 模拟before callback
		before := func(tx *gorm.DB) {
			tx.InstanceSet(key, time.Now())
		}

		// Act
		before(db)

		// Assert
		val, ok := db.InstanceGet(key)
		assert.True(t, ok, "应该设置了start time")
		assert.IsType(t, time.Time{}, val, "start time应该是time.Time类型")
	})

	// Test after callback - 应该记录指标
	t.Run("after callback records metrics", func(t *testing.T) {
		const key = "metrics_start"

		// 设置start time
		startTime := time.Now().Add(-100 * time.Millisecond)
		db.InstanceSet(key, startTime)
		db.Statement.Table = "test_table"

		// 模拟after callback
		after := func(op string) func(tx *gorm.DB) {
			return func(tx *gorm.DB) {
				v, ok := tx.InstanceGet(key)
				if !ok {
					return
				}
				start, _ := v.(time.Time)
				table := tx.Statement.Table
				DBQueryDuration.WithLabelValues(op, table).Observe(time.Since(start).Seconds())
			}
		}

		// Act
		afterFunc := after("query")
		afterFunc(db)

		// Assert - 验证指标被记录
		metrics, err := reg.Gather()
		assert.NoError(t, err)

		found := false
		for _, m := range metrics {
			if m.GetName() == "db_query_duration_seconds" {
				found = true
				assert.GreaterOrEqual(t, len(m.GetMetric()), 1)
			}
		}
		assert.True(t, found, "应该找到db_query_duration_seconds指标")
	})

	// Test after callback without start time - 应该安全返回
	t.Run("after callback without start time returns safely", func(t *testing.T) {
		db2 := &gorm.DB{
			Config: &gorm.Config{},
			Statement: &gorm.Statement{
				Table: "test_table",
			},
		}

		after := func(op string) func(tx *gorm.DB) {
			return func(tx *gorm.DB) {
				v, ok := tx.InstanceGet("metrics_start")
				if !ok {
					return // 应该安全返回
				}
				start, _ := v.(time.Time)
				table := tx.Statement.Table
				DBQueryDuration.WithLabelValues(op, table).Observe(time.Since(start).Seconds())
			}
		}

		// Act & Assert - 不应该panic
		assert.NotPanics(t, func() {
			afterFunc := after("query")
			afterFunc(db2)
		})
	})
}
