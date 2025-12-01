package safety

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

/**
 * Test_ValidateText_ValidContent 测试有效内容
 */
func Test_ValidateText_ValidContent(t *testing.T) {
	// Arrange
	content := "这是一个正常的评论内容"
	maxLen := 100

	// Act
	err := ValidateText(content, maxLen)

	// Assert
	assert.NoError(t, err)
}

/**
 * Test_ValidateText_EmptyContent 测试空内容
 */
func Test_ValidateText_EmptyContent(t *testing.T) {
	// Arrange
	content := ""
	maxLen := 100

	// Act
	err := ValidateText(content, maxLen)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "内容不能为空", err.Error())
}

/**
 * Test_ValidateText_OnlyWhitespace 测试只有空格的内容
 */
func Test_ValidateText_OnlyWhitespace(t *testing.T) {
	// Arrange
	content := "   \n\t  "
	maxLen := 100

	// Act
	err := ValidateText(content, maxLen)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "内容不能为空", err.Error())
}

/**
 * Test_ValidateText_ExceedsMaxLength 测试超出最大长度
 */
func Test_ValidateText_ExceedsMaxLength(t *testing.T) {
	// Arrange
	content := strings.Repeat("测", 101) // 101个字符
	maxLen := 100

	// Act
	err := ValidateText(content, maxLen)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "内容长度超出限制", err.Error())
}

/**
 * Test_ValidateText_ExactMaxLength 测试刚好最大长度
 */
func Test_ValidateText_ExactMaxLength(t *testing.T) {
	// Arrange
	content := strings.Repeat("测", 100) // 刚好100个字符
	maxLen := 100

	// Act
	err := ValidateText(content, maxLen)

	// Assert
	assert.NoError(t, err)
}

/**
 * Test_ValidateText_ContainsSensitiveWord_Chinese 测试包含中文敏感词
 */
func Test_ValidateText_ContainsSensitiveWord_Chinese(t *testing.T) {
	// Arrange
	sensitiveWords := []string{"违法", "违规", "涉黄"}
	maxLen := 100

	for _, word := range sensitiveWords {
		t.Run(word, func(t *testing.T) {
			content := "这个内容包含" + word + "词汇"

			// Act
			err := ValidateText(content, maxLen)

			// Assert
			assert.Error(t, err)
			assert.Equal(t, "检测到敏感词", err.Error())
		})
	}
}

/**
 * Test_ValidateText_ContainsSensitiveWord_English 测试包含英文敏感词
 */
func Test_ValidateText_ContainsSensitiveWord_English(t *testing.T) {
	// Arrange
	content := "This is a spam message"
	maxLen := 100

	// Act
	err := ValidateText(content, maxLen)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "检测到敏感词", err.Error())
}

/**
 * Test_ValidateText_SensitiveWord_CaseInsensitive 测试敏感词不区分大小写
 */
func Test_ValidateText_SensitiveWord_CaseInsensitive(t *testing.T) {
	// Arrange
	tests := []string{
		"This is SPAM",
		"This is Spam",
		"This is spam",
		"This is SpAm",
	}
	maxLen := 100

	for _, content := range tests {
		t.Run(content, func(t *testing.T) {
			// Act
			err := ValidateText(content, maxLen)

			// Assert
			assert.Error(t, err, "应该检测到敏感词")
		})
	}
}

/**
 * Test_ValidateText_WithWhitespaceBeforeValidation 测试带空格的内容验证
 */
func Test_ValidateText_WithWhitespaceBeforeValidation(t *testing.T) {
	// Arrange
	content := "  这是一个正常的内容  "
	maxLen := 100

	// Act
	err := ValidateText(content, maxLen)

	// Assert
	assert.NoError(t, err)
}

/**
 * Test_ContainsSensitiveWord_True 测试包含敏感词
 */
func Test_ContainsSensitiveWord_True(t *testing.T) {
	// Arrange
	tests := []string{
		"这个内容违法",
		"这是违规的",
		"涉黄内容",
		"This is spam",
		"fake news",
		"违法fake spam", // 多个敏感词
	}

	for _, text := range tests {
		t.Run(text, func(t *testing.T) {
			// Act
			result := ContainsSensitiveWord(text)

			// Assert
			assert.True(t, result)
		})
	}
}

/**
 * Test_ContainsSensitiveWord_False 测试不包含敏感词
 */
func Test_ContainsSensitiveWord_False(t *testing.T) {
	// Arrange
	tests := []string{
		"这是一个正常的内容",
		"This is a normal message",
		"游戏陪玩服务",
		"订单处理中",
		"",
	}

	for _, text := range tests {
		t.Run(text, func(t *testing.T) {
			// Act
			result := ContainsSensitiveWord(text)

			// Assert
			assert.False(t, result)
		})
	}
}

/**
 * Test_ContainsSensitiveWord_EmptyText 测试空文本
 */
func Test_ContainsSensitiveWord_EmptyText(t *testing.T) {
	// Act
	result := ContainsSensitiveWord("")

	// Assert
	assert.False(t, result)
}

/**
 * Test_ContainsSensitiveWord_PartialMatch 测试部分匹配
 */
func Test_ContainsSensitiveWord_PartialMatch(t *testing.T) {
	// Arrange
	tests := []string{
		"这是一个违法行为", // 包含"违法"
		"spam邮件",   // 包含"spam"
		"fake消息",   // 包含"fake"
		"内容违规处理",   // 包含"违规"
		"不涉黄内容",    // 包含"涉黄"
	}

	for _, text := range tests {
		t.Run(text, func(t *testing.T) {
			// Act
			result := ContainsSensitiveWord(text)

			// Assert
			assert.True(t, result, "应该检测到部分匹配的敏感词")
		})
	}
}

/**
 * Test_ContainsSensitiveWord_CaseVariations 测试大小写变化
 */
func Test_ContainsSensitiveWord_CaseVariations(t *testing.T) {
	// Arrange
	tests := []string{
		"SPAM",
		"Spam",
		"sPaM",
		"FAKE",
		"Fake",
		"FaKe",
	}

	for _, text := range tests {
		t.Run(text, func(t *testing.T) {
			// Act
			result := ContainsSensitiveWord(text)

			// Assert
			assert.True(t, result, "应该不区分大小写检测敏感词")
		})
	}
}

/**
 * Test_ValidateText_UTF8Characters 测试UTF-8字符计数
 */
func Test_ValidateText_UTF8Characters(t *testing.T) {
	// Arrange - 使用emoji和中文字符
	content := "测试😀内容🎉"
	maxLen := 10

	// Act
	err := ValidateText(content, maxLen)

	// Assert
	assert.NoError(t, err, "应该正确计数UTF-8字符")
}

/**
 * Test_ValidateText_MultipleErrors 测试多个错误情况优先级
 */
func Test_ValidateText_MultipleErrors(t *testing.T) {
	// Arrange - 空内容优先于其他错误
	content := ""
	maxLen := 10

	// Act
	err := ValidateText(content, maxLen)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "内容不能为空", err.Error())
}

/**
 * Test_ValidateText_LongTextWithSensitiveWord 测试超长且包含敏感词
 */
func Test_ValidateText_LongTextWithSensitiveWord(t *testing.T) {
	// Arrange - 先检查长度，再检查敏感词
	content := strings.Repeat("测", 101) + "违法"
	maxLen := 100

	// Act
	err := ValidateText(content, maxLen)

	// Assert
	assert.Error(t, err)
	// 应该先检测到长度问题
	assert.Equal(t, "内容长度超出限制", err.Error())
}

/**
 * Test_ValidateText_EdgeCases 测试边界情况
 */
func Test_ValidateText_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		content string
		maxLen  int
		wantErr bool
	}{
		{"单个字符", "a", 1, false},
		{"零长度限制", "a", 0, true},
		{"负数长度限制", "a", -1, true},
		{"单个中文", "测", 1, false},
		{"单个emoji", "😀", 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidateText(tt.content, tt.maxLen)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
