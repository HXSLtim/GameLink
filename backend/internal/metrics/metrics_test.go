package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestInitRegistersMetrics(t *testing.T) {
	reg := prometheus.NewRegistry()
	Init(reg)
	if HTTPRequestsTotal == nil || HTTPRequestDuration == nil || DBQueryDuration == nil {
		t.Fatal("expected metrics to be initialised")
	}
	// second init should be a no-op
	Init(reg)
}
