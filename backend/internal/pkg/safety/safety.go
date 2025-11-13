package safety

import (
	"errors"
	"strings"
	"unicode/utf8"
)

var defaultSensitiveWords = []string{
	"违法",
	"违规",
	"涉黄",
	"spam",
	"fake",
}

// ValidateText ensures text length and sensitivity requirements.
func ValidateText(content string, maxLen int) error {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return errors.New("内容不能为空")
	}
	if utf8.RuneCountInString(trimmed) > maxLen {
		return errors.New("内容长度超出限制")
	}
	if ContainsSensitiveWord(trimmed) {
		return errors.New("检测到敏感词")
	}
	return nil
}

// ContainsSensitiveWord returns true if text hits default block list.
func ContainsSensitiveWord(text string) bool {
	lower := strings.ToLower(text)
	for _, word := range defaultSensitiveWords {
		if strings.Contains(lower, strings.ToLower(word)) {
			return true
		}
	}
	return false
}
