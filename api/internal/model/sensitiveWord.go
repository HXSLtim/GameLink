package model

// SensitiveWordCategory 敏感词分类
type SensitiveWordCategory string

const (
	SensitiveWordCategoryPolitics SensitiveWordCategory = "politics" // 政治
	SensitiveWordCategoryPorn     SensitiveWordCategory = "porn"     // 色情
	SensitiveWordCategoryAbuse    SensitiveWordCategory = "abuse"    // 辱骂
	SensitiveWordCategoryAd       SensitiveWordCategory = "ad"       // 广告
	SensitiveWordCategoryOther    SensitiveWordCategory = "other"    // 其他
)

// Valid 检查敏感词分类是否合法
func (c SensitiveWordCategory) Valid() bool {
	switch c {
	case SensitiveWordCategoryPolitics, SensitiveWordCategoryPorn,
		SensitiveWordCategoryAbuse, SensitiveWordCategoryAd,
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

// SensitiveWordMatchType 敏感词匹配类型
type SensitiveWordMatchType string

const (
	SensitiveWordMatchTypeExact SensitiveWordMatchType = "exact" // 精确匹配
	SensitiveWordMatchTypeFuzzy SensitiveWordMatchType = "fuzzy" // 模糊匹配
	SensitiveWordMatchTypeRegex SensitiveWordMatchType = "regex" // 正则匹配
)

// Valid 检查敏感词匹配类型是否合法
func (m SensitiveWordMatchType) Valid() bool {
	switch m {
	case SensitiveWordMatchTypeExact, SensitiveWordMatchTypeFuzzy, SensitiveWordMatchTypeRegex:
		return true
	default:
		return false
	}
}

// SensitiveWord 敏感词模型
// @Description 敏感词模型，用于过滤评价内容中的不当词汇
// @Example [{"id": 1, "word": "测试敏感词", "category": "other", "matchType": "exact", "createdAt": "2024-01-15T10:30:00Z", "updatedAt": "2024-01-15T10:30:00Z"}]
type SensitiveWord struct {
	Base
	// 敏感词内容
	// @Example 测试敏感词
	Word string `json:"word" gorm:"column:word;type:varchar(100);uniqueIndex;not null"`
	// 敏感词分类
	// @Enum politics, porn, abuse, ad, other
	// @Example other
	Category SensitiveWordCategory `json:"category" gorm:"column:category;type:varchar(20);not null;index"`
	// 匹配类型
	// @Enum exact, fuzzy, regex
	// @Example exact
	MatchType SensitiveWordMatchType `json:"matchType" gorm:"column:match_type;type:varchar(20);not null;index;default:'exact'"`
	// 严重程度
	// @Enum low, medium, high
	// @Example medium
	Severity SensitiveWordSeverity `json:"severity" gorm:"column:severity;type:varchar(20);not null;default:'medium'"`
	// 替换内容（默认 ***）
	Replacement string `json:"replacement" gorm:"column:replacement;type:varchar(100);default:'***'"`
	// 是否启用
	IsActive bool `json:"isActive" gorm:"column:is_active;default:true;index"`
	// 创建人ID
	CreatedBy uint64 `json:"createdBy" gorm:"column:created_by"`
}
