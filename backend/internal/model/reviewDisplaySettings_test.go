package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReviewSortBy_Valid(t *testing.T) {
	tests := []struct {
		name     string
		sortBy   ReviewSortBy
		expected bool
	}{
		{"time is valid", ReviewSortByTime, true},
		{"score is valid", ReviewSortByScore, true},
		{"likes is valid", ReviewSortByLikes, true},
		{"empty is invalid", "", false},
		{"random is invalid", "random", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.sortBy.Valid())
		})
	}
}

func TestReviewDisplaySettings_Validate(t *testing.T) {
	tests := []struct {
		name        string
		settings    *ReviewDisplaySettings
		expectError bool
	}{
		{
			name: "valid settings",
			settings: &ReviewDisplaySettings{
				SortBy:        ReviewSortByTime,
				MinScore:      1,
				ShowAnonymous: true,
				PageSize:      20,
			},
			expectError: false,
		},
		{
			name: "valid settings with max values",
			settings: &ReviewDisplaySettings{
				SortBy:        ReviewSortByLikes,
				MinScore:      5,
				ShowAnonymous: false,
				PageSize:      100,
			},
			expectError: false,
		},
		{
			name: "invalid sortBy",
			settings: &ReviewDisplaySettings{
				SortBy:        "invalid",
				MinScore:      1,
				ShowAnonymous: true,
				PageSize:      20,
			},
			expectError: true,
		},
		{
			name: "minScore too low",
			settings: &ReviewDisplaySettings{
				SortBy:        ReviewSortByTime,
				MinScore:      0,
				ShowAnonymous: true,
				PageSize:      20,
			},
			expectError: true,
		},
		{
			name: "minScore too high",
			settings: &ReviewDisplaySettings{
				SortBy:        ReviewSortByTime,
				MinScore:      6,
				ShowAnonymous: true,
				PageSize:      20,
			},
			expectError: true,
		},
		{
			name: "pageSize too low",
			settings: &ReviewDisplaySettings{
				SortBy:        ReviewSortByTime,
				MinScore:      1,
				ShowAnonymous: true,
				PageSize:      0,
			},
			expectError: true,
		},
		{
			name: "pageSize too high",
			settings: &ReviewDisplaySettings{
				SortBy:        ReviewSortByTime,
				MinScore:      1,
				ShowAnonymous: true,
				PageSize:      101,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.settings.Validate()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDefaultReviewDisplaySettings(t *testing.T) {
	settings := DefaultReviewDisplaySettings()

	assert.Equal(t, uint64(1), settings.ID)
	assert.Equal(t, ReviewSortByTime, settings.SortBy)
	assert.Equal(t, 1, settings.MinScore)
	assert.True(t, settings.ShowAnonymous)
	assert.Equal(t, 20, settings.PageSize)

	// Default settings should be valid
	assert.NoError(t, settings.Validate())
}

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{
		Field:   "testField",
		Message: "test error message",
	}

	assert.Equal(t, "test error message", err.Error())
}
