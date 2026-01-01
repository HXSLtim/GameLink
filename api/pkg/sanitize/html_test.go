package sanitize

import (
	"fmt"
	"strings"
	"testing"
)

func TestEscapeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string // Check if output contains these substrings
	}{
		{
			name:     "simple text",
			input:    "hello world",
			contains: []string{"hello world"},
		},
		{
			name:     "script tag",
			input:    `<script>alert('XSS')</script>`,
			contains: []string{"&lt;script&gt;", "&lt;/script&gt;"},
		},
		{
			name:     "HTML tag",
			input:    `<div>content</div>`,
			contains: []string{"&lt;div&gt;", "content", "&lt;/div&gt;"},
		},
		{
			name:     "special characters",
			input:    `&<>"'`,
			contains: []string{"&amp;", "&lt;", "&gt;", "&#34;", "&#39;"},
		},
		{
			name:     "mixed content",
			input:    `<img src=x onerror=alert(1)>`,
			contains: []string{"&lt;img", "onerror=alert(1)", "&gt;"},
		},
		{
			name:     "empty string",
			input:    "",
			contains: []string{""},
		},
		{
			name:     "Chinese characters",
			input:    "你好<script>alert('XSS')</script>世界",
			contains: []string{"你好", "&lt;script&gt;", "世界"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeString(tt.input)
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("EscapeString() = %v, want to contain %v", result, expected)
				}
			}
		})
	}
}

func TestUnescapeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "HTML entities",
			input:    "&lt;div&gt;content&lt;/div&gt;",
			expected: "<div>content</div>",
		},
		{
			name:     "special characters",
			input:    "&amp;&lt;&gt;&quot;&#39;",
			expected: `&<>"'`,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UnescapeString(tt.input)
			if result != tt.expected {
				t.Errorf("UnescapeString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSanitizeHTML(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		shouldNotHave  []string // These patterns should be removed
		mustContain    []string // These must still exist
	}{
		{
			name:          "script tag",
			input:         `<script>alert('XSS')</script>`,
			shouldNotHave: []string{"<script", "</script>", "script"},
			mustContain:   []string{},
		},
		{
			name:          "event handler",
			input:         `<div onclick="alert('XSS')">content</div>`,
			shouldNotHave: []string{"onclick", "alert"},
			mustContain:   []string{"<div", "content", "</div>"},
		},
		{
			name:          "javascript protocol",
			input:         `<a href="javascript:alert('XSS')">link</a>`,
			shouldNotHave: []string{"javascript:"},
			mustContain:   []string{"<a", "link", "</a>"},
		},
		{
			name:          "data URI",
			input:         `<iframe src="data:text/html,<script>alert('XSS')</script>"></iframe>`,
			shouldNotHave: []string{"data:", "text/html"},
			mustContain:   []string{"<iframe"},
		},
		{
			name:          "clean text",
			input:         "Hello World",
			shouldNotHave: []string{},
			mustContain:   []string{"Hello World"},
		},
		{
			name:          "XSS vector",
			input:         `<IMG SRC=j&#97;v&#97;script:alert('XSS')>`,
			shouldNotHave: []string{"javascript"},
			mustContain:   []string{"<IMG"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeHTML(tt.input)
			for _, pattern := range tt.shouldNotHave {
				if strings.Contains(strings.ToLower(result), strings.ToLower(pattern)) {
					// For case-insensitive check
					lowerPattern := strings.ToLower(pattern)
					lowerResult := strings.ToLower(result)
					if strings.Contains(lowerResult, lowerPattern) {
						t.Errorf("SanitizeHTML() should not contain %v, got %v", pattern, result)
					}
				}
			}
			for _, pattern := range tt.mustContain {
				if !strings.Contains(result, pattern) {
					t.Errorf("SanitizeHTML() should contain %v, got %v", pattern, result)
				}
			}
		})
	}
}

func TestContainsXSSPatterns(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "script tag",
			input:    `<script>alert('XSS')</script>`,
			expected: true,
		},
		{
			name:     "event handler",
			input:    `<div onclick="alert('XSS')">content</div>`,
			expected: true,
		},
		{
			name:     "javascript protocol",
			input:    `<a href="javascript:alert('XSS')">link</a>`,
			expected: true,
		},
		{
			name:     "data URI",
			input:    `data:text/html,<script>alert('XSS')</script>`,
			expected: true,
		},
		{
			name:     "clean text",
			input:    "Hello World",
			expected: false,
		},
		{
			name:     "HTML without JS",
			input:    `<div>content</div>`,
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "case insensitive script",
			input:    `<SCRIPT>alert('XSS')</SCRIPT>`,
			expected: true,
		},
		{
			name:     "case insensitive javascript",
			input:    `JAVASCRIPT:alert('XSS')`,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsXSSPatterns(tt.input)
			if result != tt.expected {
				t.Errorf("ContainsXSSPatterns() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestStripTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple HTML",
			input:    `<div>content</div>`,
			expected: "content",
		},
		{
			name:     "nested tags",
			input:    `<div><p>nested</p></div>`,
			expected: "nested",
		},
		{
			name:     "HTML entities",
			input:    "Hello&nbsp;World&lt;3",
			expected: "Hello World<3",
		},
		{
			name:     "mixed content",
			input:    `<p>Hello</p> <strong>World</strong>`,
			expected: "Hello World",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no tags",
			input:    "Plain text",
			expected: "Plain text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripTags(tt.input)
			if result != tt.expected {
				t.Errorf("StripTags() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLength int
		maxLen    int // Expected max length (either <= maxLength or maxLength+3 for ellipsis)
	}{
		{
			name:      "short string",
			input:     "hello",
			maxLength: 10,
			maxLen:    10,
		},
		{
			name:      "exact length",
			input:     "hello",
			maxLength: 5,
			maxLen:    5,
		},
		{
			name:      "long string",
			input:     "hello world",
			maxLength: 5,
			maxLen:    8, // 5 + "..."
		},
		{
			name:      "Chinese characters",
			input:     "你好世界",
			maxLength: 2,
			maxLen:    9, // 2 Chinese chars = 6 bytes + 3 for "..."
		},
		{
			name:      "mixed UTF-8",
			input:     "Hi你好",
			maxLength: 3,
			maxLen:    9, // "H"=1, "i"=1, "你"=3, total=5 runes + 3 = 8 bytes for "..."
		},
		{
			name:      "empty string",
			input:     "",
			maxLength: 10,
			maxLen:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateString(tt.input, tt.maxLength)
			if len(result) > tt.maxLen {
				t.Errorf("TruncateString() length = %v, want <= %v", len(result), tt.maxLen)
			}
		})
	}
}

func TestValidateUTF8(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid ASCII",
			input:    "Hello World",
			expected: true,
		},
		{
			name:     "valid UTF-8 Chinese",
			input:    "你好世界",
			expected: true,
		},
		{
			name:     "valid UTF-8 emoji",
			input:    "😀🎉",
			expected: true,
		},
		{
			name:     "empty string",
			input:    "",
			expected: true,
		},
		{
			name:     "invalid UTF-8 sequence",
			input:    string([]byte{0xFF, 0xFE, 0xFD}),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateUTF8(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateUTF8() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSanitizeNickname(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		contains []string
	}{
		{
			name:     "simple nickname",
			input:    "PlayerOne",
			maxLen:   32,
			contains: []string{"PlayerOne"},
		},
		{
			name:     "nickname with script",
			input:    `<script>alert('XSS')</script>Player`,
			maxLen:   32,
			contains: []string{"Player"},
		},
		{
			name:     "nickname with HTML tags",
			input:    "<div>Player</div>",
			maxLen:   32,
			contains: []string{"Player"},
		},
		{
			name:     "long nickname",
			input:    "ThisIsAVeryLongNicknameThatExceedsThirtyTwoCharacters",
			maxLen:   35, // 32 + "..."
			contains: []string{"ThisIsAVeryLongNicknameThatExcee"},
		},
		{
			name:     "nickname with spaces",
			input:    "  Player One  ",
			maxLen:   32,
			contains: []string{"Player One"},
		},
		{
			name:     "nickname with special chars",
			input:    "Player<>&\"'",
			maxLen:   32,
			contains: []string{"Player"},
		},
		{
			name:     "Chinese nickname",
			input:    "游戏玩家",
			maxLen:   32,
			contains: []string{"游戏玩家"},
		},
		{
			name:     "empty string",
			input:    "",
			maxLen:   0,
			contains: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeNickname(tt.input)
			if len(result) > tt.maxLen {
				t.Errorf("SanitizeNickname() length = %v, want <= %v", len(result), tt.maxLen)
			}
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("SanitizeNickname() = %v, want to contain %v", result, expected)
				}
			}
		})
	}
}

func TestSanitizeMessage(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		shouldNotHave []string
		mustContain   []string
	}{
		{
			name:          "simple message",
			input:         "Hello, how are you?",
			shouldNotHave: []string{"&lt;", "&gt;"},
			mustContain:   []string{"Hello, how are you?"},
		},
		{
			name:          "message with script",
			input:         `<script>alert('XSS')</script>Hello`,
			shouldNotHave: []string{"<script", "</script>", "script"},
			mustContain:   []string{"Hello"},
		},
		{
			name:          "message with event handler",
			input:         `<div onclick="alert('XSS')">Click me</div>`,
			shouldNotHave: []string{"onclick", "alert"},
			mustContain:   []string{"Click me"},
		},
		{
			name:          "message with javascript protocol",
			input:         `<a href="javascript:alert('XSS')">link</a>`,
			shouldNotHave: []string{"javascript:"},
			mustContain:   []string{"link"},
		},
		{
			name:          "Chinese message",
			input:         "你好<script>alert('XSS')</script>世界",
			shouldNotHave: []string{"script", "<"},
			mustContain:   []string{"你好", "世界"},
		},
		{
			name:          "empty string",
			input:         "",
			shouldNotHave: []string{},
			mustContain:   []string{""},
		},
		{
			name:          "message with special chars",
			input:         "Hello <>&\"' World",
			shouldNotHave: []string{},
			mustContain:   []string{"Hello", "World"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeMessage(tt.input)
			for _, pattern := range tt.shouldNotHave {
				if strings.Contains(result, pattern) {
					t.Errorf("SanitizeMessage() should not contain %v, got %v", pattern, result)
				}
			}
			for _, expected := range tt.mustContain {
				if !strings.Contains(result, expected) {
					t.Errorf("SanitizeMessage() = %v, want to contain %v", result, expected)
				}
			}
		})
	}
}

func TestSanitizeReview(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:     "simple review",
			input:    "Great service!",
			contains: []string{"Great service!"},
		},
		{
			name:     "review with XSS",
			input:    `<script>alert('XSS')</script>Great service!`,
			contains: []string{"Great service!"},
		},
		{
			name:     "empty string",
			input:    "",
			contains: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeReview(tt.input)
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("SanitizeReview() = %v, want to contain %v", result, expected)
				}
			}
		})
	}
}

func TestSanitizeReport(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
		maxLen   int
	}{
		{
			name:     "simple report",
			input:    "User violated rules",
			contains: []string{"User violated rules"},
			maxLen:   1000,
		},
		{
			name:     "report with XSS",
			input:    `<script>alert('XSS')</script>User violated rules`,
			contains: []string{"User violated rules"},
			maxLen:   1000,
		},
		{
			name:     "empty string",
			input:    "",
			contains: []string{""},
			maxLen:   0,
		},
		{
			name:     "long report",
			input:    string(make([]byte, 2000)),
			contains: []string{},
			maxLen:   1003, // 1000 + "..."
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeReport(tt.input)
			if len(result) > tt.maxLen {
				t.Errorf("SanitizeReport() length = %v, want <= %v", len(result), tt.maxLen)
			}
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("SanitizeReport() = %v, want to contain %v", result, expected)
				}
			}
		})
	}
}

func TestEscapeAll(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		checks   map[string]string // key -> expected substring
	}{
		{
			name: "simple map",
			input: map[string]interface{}{
				"name":  "John<script>",
				"email": "john@example.com",
			},
			checks: map[string]string{
				"name":  "&lt;script&gt;",
				"email": "john@example.com",
			},
		},
		{
			name: "nested map",
			input: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "John<script>",
				},
			},
			checks: map[string]string{
				"user.name": "&lt;script&gt;",
			},
		},
		{
			name: "array",
			input: map[string]interface{}{
				"items": []interface{}{"item1<script>", "item2"},
			},
			checks: map[string]string{
				"items.0": "&lt;script&gt;",
			},
		},
		{
			name: "mixed types",
			input: map[string]interface{}{
				"name":  "John",
				"age":   30,
				"admin": true,
			},
			checks: map[string]string{
				"name": "John",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeAll(tt.input)
			for keyPath, expectedValue := range tt.checks {
				// Navigate nested structure using keyPath (e.g., "user.name" or "items.0")
				parts := strings.Split(keyPath, ".")
				var current interface{} = result
				for _, part := range parts {
					switch v := current.(type) {
					case map[string]interface{}:
						current = v[part]
					case []interface{}:
						idx := 0
						fmt.Sscanf(part, "%d", &idx)
						if idx >= 0 && idx < len(v) {
							current = v[idx]
						} else {
							current = nil
						}
					}
				}
				if resultStr, ok := current.(string); ok {
					if !strings.Contains(resultStr, expectedValue) {
						t.Errorf("EscapeAll()[%v] = %v, want to contain %v", keyPath, resultStr, expectedValue)
					}
				}
			}
		})
	}
}

// Benchmark tests
func BenchmarkEscapeString(b *testing.B) {
	input := `<script>alert('XSS')</script>Hello World`
	for i := 0; i < b.N; i++ {
		EscapeString(input)
	}
}

func BenchmarkSanitizeHTML(b *testing.B) {
	input := `<script>alert('XSS')</script><div onclick="alert(1)">Content</div>`
	for i := 0; i < b.N; i++ {
		SanitizeHTML(input)
	}
}

func BenchmarkContainsXSSPatterns(b *testing.B) {
	input := `<div onclick="alert('XSS')">Content</div>`
	for i := 0; i < b.N; i++ {
		ContainsXSSPatterns(input)
	}
}

func BenchmarkSanitizeMessage(b *testing.B) {
	input := `<script>alert('XSS')</script>Hello World`
	for i := 0; i < b.N; i++ {
		SanitizeMessage(input)
	}
}
