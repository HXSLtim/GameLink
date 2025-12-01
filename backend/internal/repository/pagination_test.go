package repository_test

import (
	"testing"

	"gamelink/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestNormalizePage(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{
			name:     "negative value should return default",
			input:    -1,
			expected: 1,
		},
		{
			name:     "zero should return default",
			input:    0,
			expected: 1,
		},
		{
			name:     "positive value should return as is",
			input:    5,
			expected: 5,
		},
		{
			name:     "large positive value should return as is",
			input:    1000,
			expected: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repository.NormalizePage(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizePageSize(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{
			name:     "negative value should return default",
			input:    -1,
			expected: 20,
		},
		{
			name:     "zero should return default",
			input:    0,
			expected: 20,
		},
		{
			name:     "positive value within limit should return as is",
			input:    10,
			expected: 10,
		},
		{
			name:     "value equal to max should return max",
			input:    100,
			expected: 100,
		},
		{
			name:     "value exceeding max should return max",
			input:    150,
			expected: 100,
		},
		{
			name:     "large value should return max",
			input:    1000,
			expected: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repository.NormalizePageSize(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPaginationConstants(t *testing.T) {
	// Test that constants have expected values
	assert.Equal(t, 1, repository.NormalizePage(1))
	assert.Equal(t, 20, repository.NormalizePageSize(20))

	// Test edge cases around max page size
	assert.Equal(t, 100, repository.NormalizePageSize(100))
	assert.Equal(t, 100, repository.NormalizePageSize(101))
	assert.Equal(t, 100, repository.NormalizePageSize(1000))
}
