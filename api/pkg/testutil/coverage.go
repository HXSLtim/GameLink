// Package testutil provides utilities for testing.
package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// CoverageInfo holds coverage information.
type CoverageInfo struct {
	Total       int
	Covered     int
	CoveragePct float64
}

// AssertCoverage asserts that coverage meets the threshold.
func AssertCoverage(t *testing.T, coverage float64, threshold float64) {
	t.Helper()

	assert.GreaterOrEqual(t, coverage, threshold,
		"coverage %.2f%% is below threshold %.2f%%", coverage, threshold)
}

// AssertCoverageAtLeast asserts that coverage is at least the threshold.
func AssertCoverageAtLeast(t *testing.T, coverage float64, minCoverage float64) {
	t.Helper()

	if coverage < minCoverage {
		t.Errorf("coverage %.2f%% is below minimum %.2f%%", coverage, minCoverage)
	}
}

// AssertCoverageExactly asserts that coverage is exactly the expected value.
func AssertCoverageExactly(t *testing.T, coverage float64, expected float64) {
	t.Helper()

	assert.InDelta(t, expected, coverage, 0.01,
		"coverage should be %.2f%%, got %.2f%%", expected, coverage)
}

// AssertCoverageReport generates a coverage report for a package.
func AssertCoverageReport(t *testing.T, pkg string, expectedCoverage float64) {
	t.Helper()

	// This is a placeholder for actual coverage reporting logic
	// In a real implementation, you would run go test -cover for the package
	// and compare the actual coverage with the expected coverage.

	t.Logf("Package %s coverage report: expected %.2f%%", pkg, expectedCoverage)
}

// TrackCoverage tracks coverage for multiple packages.
func TrackCoverage(t *testing.T, coverageMap map[string]float64) {
	t.Helper()

	for pkg, coverage := range coverageMap {
		t.Logf("Package: %s, Coverage: %.2f%%", pkg, coverage)
	}
}

// CoverageThresholds defines coverage thresholds for different types of code.
var CoverageThresholds = struct {
	Model      float64
	Repository float64
	Service    float64
	Handler    float64
	Utils      float64
}{
	Model:      100.0,
	Repository: 100.0,
	Service:    100.0,
	Handler:    100.0,
	Utils:      100.0,
}
