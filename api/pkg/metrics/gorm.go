package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

// InstrumentGorm registers callbacks to measure operation durations.
func InstrumentGorm(db *gorm.DB) error {
	// Ensure metrics are initialized; use default registerer
	Init(prometheus.DefaultRegisterer)

	// we store start time in instance settings
	const key = "metrics_start"

	before := func(tx *gorm.DB) {
		tx.InstanceSet(key, time.Now())
	}
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

	_ = db.Callback().Query().Before("gorm:query").Register("metrics:before_query", before)
	_ = db.Callback().Query().After("gorm:after_query").Register("metrics:after_query", after("query"))

	_ = db.Callback().Create().Before("gorm:create").Register("metrics:before_create", before)
	_ = db.Callback().Create().After("gorm:after_create").Register("metrics:after_create", after("create"))

	_ = db.Callback().Update().Before("gorm:update").Register("metrics:before_update", before)
	_ = db.Callback().Update().After("gorm:after_update").Register("metrics:after_update", after("update"))

	_ = db.Callback().Delete().Before("gorm:delete").Register("metrics:before_delete", before)
	_ = db.Callback().Delete().After("gorm:after_delete").Register("metrics:after_delete", after("delete"))

	return nil
}
