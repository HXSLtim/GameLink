package metrics

import (
	"sync"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

func resetMetricsForTest() {
	if HTTPRequestsTotal != nil {
		prometheus.Unregister(HTTPRequestsTotal)
	}
	if HTTPRequestDuration != nil {
		prometheus.Unregister(HTTPRequestDuration)
	}
	if DBQueryDuration != nil {
		prometheus.Unregister(DBQueryDuration)
	}
	HTTPRequestsTotal = nil
	HTTPRequestDuration = nil
	DBQueryDuration = nil
	once = sync.Once{}
}

func TestInitRegistersMetrics(t *testing.T) {
	resetMetricsForTest()
	reg := prometheus.NewRegistry()
	Init(reg)
	if HTTPRequestsTotal == nil || HTTPRequestDuration == nil || DBQueryDuration == nil {
		t.Fatal("expected metrics to be initialised")
	}
	// second init should be a no-op
	Init(reg)
}

func TestInstrumentGormRecordsDurations(t *testing.T) {
	resetMetricsForTest()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	type instrumentUser struct {
		ID   uint
		Name string
	}
	if err := db.AutoMigrate(&instrumentUser{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	if err := InstrumentGorm(db); err != nil {
		t.Fatalf("instrument gorm: %v", err)
	}

	// create, query, update, delete to trigger callbacks
	u := &instrumentUser{Name: "alice"}
	if err := db.Create(u).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := db.First(&instrumentUser{}, u.ID).Error; err != nil {
		t.Fatalf("query: %v", err)
	}
	if err := db.Model(u).Update("Name", "bob").Error; err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := db.Delete(&instrumentUser{}, u.ID).Error; err != nil {
		t.Fatalf("delete: %v", err)
	}

	mfs, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		t.Fatalf("gather: %v", err)
	}
	found := false
	for _, mf := range mfs {
		if mf.GetName() != "db_query_duration_seconds" {
			continue
		}
		for _, metric := range mf.GetMetric() {
			if metric.GetHistogram().GetSampleCount() > 0 {
				found = true
				break
			}
		}
	}
	if !found {
		t.Fatal("expected db_query_duration_seconds to have samples")
	}
}
