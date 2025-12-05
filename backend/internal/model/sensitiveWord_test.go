package model

import (
	"testing"
)

func TestSensitiveWordCategory_Valid(t *testing.T) {
	tests := []struct {
		name     string
		category SensitiveWordCategory
		want     bool
	}{
		{"political", SensitiveWordCategoryPolitical, true},
		{"pornographic", SensitiveWordCategoryPornographic, true},
		{"violent", SensitiveWordCategoryViolent, true},
		{"advertising", SensitiveWordCategoryAdvertising, true},
		{"other", SensitiveWordCategoryOther, true},
		{"invalid", SensitiveWordCategory("invalid"), false},
		{"empty", SensitiveWordCategory(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.category.Valid(); got != tt.want {
				t.Errorf("SensitiveWordCategory.Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSensitiveWordSeverity_Valid(t *testing.T) {
	tests := []struct {
		name     string
		severity SensitiveWordSeverity
		want     bool
	}{
		{"low", SensitiveWordSeverityLow, true},
		{"medium", SensitiveWordSeverityMedium, true},
		{"high", SensitiveWordSeverityHigh, true},
		{"invalid", SensitiveWordSeverity("invalid"), false},
		{"empty", SensitiveWordSeverity(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.severity.Valid(); got != tt.want {
				t.Errorf("SensitiveWordSeverity.Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
