// Package sanitize provides input sanitization functions to prevent XSS attacks.
package sanitize

import (
	"html"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// HTML special characters that need to be escaped
var (
	// Script tag pattern to detect potential XSS
	scriptPattern = regexp.MustCompile(`(?i)<\s*script\b[^>]*>.*?<\s*/\s*script\s*>`)
	// Event handler pattern (onclick, onload, etc.)
	eventHandlerPattern = regexp.MustCompile(`(?i)\s+on\w+\s*=\s*[^>]*`)
	// javascript: protocol pattern
	jsProtocolPattern = regexp.MustCompile(`(?i)javascript\s*:`)
	// Data URI with base64 pattern (can be used for XSS)
	dataURIPattern = regexp.MustCompile(`(?i)data\s*:\s*text/html`)
)

// EscapeString escapes HTML special characters to prevent XSS attacks.
// It converts &, <, >, ", ' to their HTML entity equivalents.
//
// Example:
//
//	input := `<script>alert('XSS')</script>`
//	output := sanitize.EscapeString(input)
//	// Output: &lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;
func EscapeString(s string) string {
	if s == "" {
		return s
	}
	return html.EscapeString(s)
}

// UnescapeString unescapes HTML entities.
// Use with caution - only on trusted data.
func UnescapeString(s string) string {
	if s == "" {
		return s
	}
	return html.UnescapeString(s)
}

// SanitizeHTML removes potentially dangerous HTML content.
// This is more aggressive than EscapeString and completely removes script tags,
// event handlers, and dangerous protocols.
//
// Use this for user-generated content that should be displayed as plain text.
func SanitizeHTML(s string) string {
	if s == "" {
		return s
	}

	// Remove script tags
	s = scriptPattern.ReplaceAllString(s, "")
	// Remove event handlers
	s = eventHandlerPattern.ReplaceAllString(s, "")
	// Remove javascript: protocol
	s = jsProtocolPattern.ReplaceAllString(s, "")
	// Remove data:text/html URIs
	s = dataURIPattern.ReplaceAllString(s, "")

	return s
}

// EscapeAll escapes all string fields in a map.
// Useful for escaping request bodies before processing.
func EscapeAll(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(data))
	for key, value := range data {
		switch v := value.(type) {
		case string:
			result[key] = EscapeString(v)
		case map[string]interface{}:
			result[key] = EscapeAll(v)
		case []interface{}:
			result[key] = escapeArray(v)
		default:
			result[key] = value
		}
	}
	return result
}

func escapeArray(arr []interface{}) []interface{} {
	result := make([]interface{}, len(arr))
	for i, value := range arr {
		switch v := value.(type) {
		case string:
			result[i] = EscapeString(v)
		case map[string]interface{}:
			result[i] = EscapeAll(v)
		case []interface{}:
			result[i] = escapeArray(v)
		default:
			result[i] = value
		}
	}
	return result
}

// ContainsXSSPatterns checks if a string contains potential XSS patterns.
// Returns true if suspicious patterns are detected.
func ContainsXSSPatterns(s string) bool {
	if s == "" {
		return false
	}

	// Check for script tags
	if scriptPattern.MatchString(s) {
		return true
	}

	// Check for event handlers
	if eventHandlerPattern.MatchString(s) {
		return true
	}

	// Check for javascript: protocol
	if jsProtocolPattern.MatchString(s) {
		return true
	}

	// Check for data:text/html URIs
	if dataURIPattern.MatchString(s) {
		return true
	}

	return false
}

// StripTags removes all HTML tags from a string.
// This is useful when you want to display only the text content.
func StripTags(s string) string {
	if s == "" {
		return s
	}

	// Simple tag removal - for production, consider using a proper HTML parser
	tagPattern := regexp.MustCompile(`<[^>]*>`)
	s = tagPattern.ReplaceAllString(s, "")

	// Also remove common HTML entities
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")
	s = strings.ReplaceAll(s, "&apos;", "'")

	return s
}

// TruncateString truncates a string to a maximum length, adding ellipsis if needed.
// This helps prevent extremely long strings that could be used for DoS attacks.
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}

	// Truncate by runes (not bytes) to handle UTF-8 correctly
	runes := []rune(s)
	if len(runes) <= maxLength {
		return s
	}

	return string(runes[:maxLength]) + "..."
}

// ValidateUTF8 checks if a string contains only valid UTF-8 characters.
// Invalid UTF-8 sequences can be used to bypass security checks.
func ValidateUTF8(s string) bool {
	for i, r := range s {
		if r == utf8.RuneError {
			// Check if it's a real error or just a valid rune encoding issue
			if !utf8.ValidRune(r) {
				return false
			}
			// Check if the byte sequence is invalid
			if s[i] >= 0x80 && s[i] <= 0xBF {
				// Continuation byte without a starter
				return false
			}
		}
	}
	return utf8.ValidString(s)
}

// SanitizeNickname sanitizes user nicknames specifically.
// Nicknames have special requirements: they should be short and contain no HTML.
func SanitizeNickname(nickname string) string {
	if nickname == "" {
		return nickname
	}

	// Remove dangerous patterns first (before StripTags)
	nickname = SanitizeHTML(nickname)

	// Remove HTML tags
	nickname = StripTags(nickname)

	// Trim whitespace
	nickname = strings.TrimSpace(nickname)

	// Limit length (typically 32 characters for nicknames)
	nickname = TruncateString(nickname, 32)

	// Escape any remaining special characters
	return EscapeString(nickname)
}

// SanitizeMessage sanitizes chat messages and similar content.
// Allows basic formatting but removes dangerous content.
func SanitizeMessage(content string) string {
	if content == "" {
		return content
	}

	// Remove dangerous HTML/JS patterns
	content = SanitizeHTML(content)

	// Validate UTF-8
	if !ValidateUTF8(content) {
		// If invalid, strip all non-printable characters
		content = strings.Map(func(r rune) rune {
			if unicode.IsPrint(r) || unicode.IsSpace(r) {
				return r
			}
			return -1
		}, content)
	}

	// Limit length (typically 5000 characters for messages)
	content = TruncateString(content, 5000)

	// Escape special characters
	return EscapeString(content)
}

// SanitizeReview sanitizes review and feedback content.
func SanitizeReview(content string) string {
	return SanitizeMessage(content) // Same rules as messages
}

// SanitizeReport sanitizes report/dispute content.
func SanitizeReport(content string) string {
	if content == "" {
		return content
	}

	// Remove dangerous patterns
	content = SanitizeHTML(content)

	// Validate UTF-8
	if !ValidateUTF8(content) {
		content = strings.Map(func(r rune) rune {
			if unicode.IsPrint(r) || unicode.IsSpace(r) {
				return r
			}
			return -1
		}, content)
	}

	// Limit length (typically 1000 characters for reports)
	content = TruncateString(content, 1000)

	return EscapeString(content)
}
