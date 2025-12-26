package model

import "fmt"

// Rating 表示 1-5 的打分区间。
type Rating uint8

// Rating bounds define the valid inclusive range.
const (
	RatingMin Rating = 1
	RatingMax Rating = 5
)

// Valid 检查评分是否在合法区间内。
func (r Rating) Valid() bool {
	return r >= RatingMin && r <= RatingMax
}

// MustRating 返回合法评分，越界时 panic，用于初始化常量或测试。
func MustRating(value uint8) Rating {
	r := Rating(value)
	if !r.Valid() {
		panic(fmt.Sprintf("invalid rating: %d", value))
	}
	return r
}

// GormDataType 指定评分的默认列类型。
func (Rating) GormDataType() string {
	return "tinyint"
}
