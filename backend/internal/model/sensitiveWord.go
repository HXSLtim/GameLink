package model

// SensitiveWordCategory 敏感词分类
type SensitiveWordCategory string

const (
	SensitiveWordCategoryPolitical    SensitiveWordCategory = "political"    // 政治
	SensitiveWordCategoryPornographic SensitiveWordCategory = "pornographic" // 色情
	SensitiveWordCategoryViolent      SensitiveWordCategory = "violent"      // 暴力
	SensitiveWordCategoryAdvertising  SensitiveWordCategory = "advertising"  // 广告
	SensitiveWordCategoryOther        SensitiveWordCategory = "other"        // 其他
)

// Valid 检查敏感词分类是否合法
func (c SensitiveWordCategory) Valid() bool {
	switch c {
	case SensitiveWordCategoryPolitical, SensitiveWordCategoryPornographic,
		SensitiveWordCategoryViolent, SensitiveWordCategoryAdvertising,
		SensitiveWordCategoryOther:
		return true
	default:
		return false
	}
}

// SensitiveWordSeverity 敏感词严重程度
type SensitiveWordSeverity string

const (
	SensitiveWordSeverityLow    SensitiveWordSeverity = "low"    // 低
	SensitiveWordSeverityMedium SensitiveWordSeverity = "medium" // 中
	SensitiveWordSeverityHigh   SensitiveWordSeverity = "high"   // 高
)

// Valid 检查敏感词严重程度是否合法
func (s SensitiveWordSeverity) Valid() bool {
	switch s {
	case SensitiveWordSeverityLow, SensitiveWordSeverityMedium, SensitiveWordSeverityHigh:
		return true
	default:
		return false
	}
}

// SensitiveWord 敏感词模型
// @Description 敏感词模型，用于过滤评价内容中的不当词汇
// @Example [{"id": 1, "word": "测试敏感词", "category": "other", "severity": "low", "createdAt": "2024-01-15T10:30:00Z", "updatedAt": "2024-01-15T10:30:00Z"}]
type SensitiveWord struct {
	Base
	// 敏感词内容
	// @Example 测试敏感词
	Word string `json:"word" gorm:"column:word;type:varchar(100);uniqueIndex;not null"`
	// 敏感词分类
	// @Enum political, pornographic, violent, advertising, other
	// @Example other
	Category SensitiveWordCategory `json:"category" gorm:"column:category;type:varchar(20);not null;index"`
	// 严重程度
	// @Enum low, medium, high
	// @Example low
	Severity SensitiveWordSeverity `json:"severity" gorm:"column:severity;type:varchar(20);not null;index"`
}
