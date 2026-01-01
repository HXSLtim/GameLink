package middleware

import (
	"testing"
)

func TestShouldExcludePath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		exclude  string
		expected bool
	}{
		// Exact match tests
		{
			name:     "exact match - same path",
			path:     "/api/v1/health",
			exclude:  "/api/v1/health",
			expected: true,
		},
		{
			name:     "exact match - different path",
			path:     "/api/v1/healthz",
			exclude:  "/api/v1/health",
			expected: false,
		},
		{
			name:     "exact match - subpath should not match",
			path:     "/api/v1/health/check",
			exclude:  "/api/v1/health",
			expected: false,
		},
		{
			name:     "exact match - prefix should not match",
			path:     "/foo/public",
			exclude:  "/public",
			expected: false,
		},

		// Prefix wildcard tests
		{
			name:     "prefix wildcard - direct match",
			path:     "/api/v1/public",
			exclude:  "/api/v1/public/*",
			expected: true,
		},
		{
			name:     "prefix wildcard - with trailing slash",
			path:     "/api/v1/public/",
			exclude:  "/api/v1/public/*",
			expected: true,
		},
		{
			name:     "prefix wildcard - subpath matches",
			path:     "/api/v1/public/upload",
			exclude:  "/api/v1/public/*",
			expected: true,
		},
		{
			name:     "prefix wildcard - nested subpath matches",
			path:     "/api/v1/public/images/avatar.png",
			exclude:  "/api/v1/public/*",
			expected: true,
		},
		{
			name:     "prefix wildcard - different prefix does not match",
			path:     "/api/v1/private/upload",
			exclude:  "/api/v1/public/*",
			expected: false,
		},
		{
			name:     "prefix wildcard - similar but different prefix",
			path:     "/api/v1/public2/upload",
			exclude:  "/api/v1/public/*",
			expected: false,
		},

		// Catch-all tests
		{
			name:     "catch-all - matches any path",
			path:     "/any/path/here",
			exclude:  "*",
			expected: true,
		},
		{
			name:     "catch-all - matches root",
			path:     "/",
			exclude:  "*",
			expected: true,
		},

		// Suffix wildcard tests
		{
			name:     "suffix wildcard - matches ending",
			path:     "/api/v1/health",
			exclude:  "*health",
			expected: true,
		},
		{
			name:     "suffix wildcard - does not match different ending",
			path:     "/api/v1/ping",
			exclude:  "*health",
			expected: false,
		},

		// Empty/whitespace tests
		{
			name:     "empty exclude pattern",
			path:     "/api/v1/health",
			exclude:  "",
			expected: false,
		},
		{
			name:     "whitespace exclude pattern",
			path:     "/api/v1/health",
			exclude:  "  ",
			expected: false,
		},

		// Real-world scenarios
		{
			name:     "healthz should not match health",
			path:     "/api/v1/healthz",
			exclude:  "/api/v1/health",
			expected: false,
		},
		{
			name:     "auth/refresh exact match",
			path:     "/api/v1/auth/refresh",
			exclude:  "/api/v1/auth/refresh",
			expected: true,
		},
		{
			name:     "auth/refresh should not match auth/refresh_token",
			path:     "/api/v1/auth/refresh_token",
			exclude:  "/api/v1/auth/refresh",
			expected: false,
		},
		{
			name:     "wildcard public endpoint",
			path:     "/api/v1/static/css/style.css",
			exclude:  "/api/v1/static/*",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldExcludePath(tt.path, tt.exclude)
			if result != tt.expected {
				t.Errorf("shouldExcludePath(%q, %q) = %v, want %v",
					tt.path, tt.exclude, result, tt.expected)
			}
		})
	}
}
