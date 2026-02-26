package router

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestSyncAPIPermissionsSkipList_UsesRootMetricsPath(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve current test file path")
	}

	routerFile := filepath.Join(filepath.Dir(filename), "router.go")
	content, err := os.ReadFile(routerFile)
	if err != nil {
		t.Fatalf("failed to read %s: %v", routerFile, err)
	}

	source := string(content)
	if !strings.Contains(source, "\"/metrics\"") {
		t.Fatalf("expected sync skip list to include /metrics in %s", routerFile)
	}
	if strings.Contains(source, "\"/api/v1/metrics\"") {
		t.Fatalf("expected legacy /api/v1/metrics reference to be removed from %s", routerFile)
	}
}
