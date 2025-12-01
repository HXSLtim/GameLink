package model_test

import (
	"encoding/json"
	"testing"

	"gamelink/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestRatingValid(t *testing.T) {
	// 测试有效的评分
	validRatings := []model.Rating{1, 2, 3, 4, 5}
	for _, rating := range validRatings {
		assert.True(t, rating.Valid(), "Rating %d should be valid", rating)
	}

	// 测试无效的评分
	invalidRatings := []model.Rating{0, 6, 7, 10, 255}
	for _, rating := range invalidRatings {
		assert.False(t, rating.Valid(), "Rating %d should be invalid", rating)
	}
}

func TestRatingConstants(t *testing.T) {
	assert.Equal(t, model.Rating(1), model.RatingMin)
	assert.Equal(t, model.Rating(5), model.RatingMax)
}

func TestMustRating(t *testing.T) {
	// 测试有效的评分
	validValues := []uint8{1, 2, 3, 4, 5}
	for _, value := range validValues {
		rating := model.MustRating(value)
		assert.Equal(t, model.Rating(value), rating)
	}

	// 测试无效的值应该panic
	invalidValues := []uint8{0, 6, 7, 10}
	for _, value := range invalidValues {
		assert.Panics(t, func() {
			model.MustRating(value)
		}, "MustRating(%d) should panic", value)
	}
}

func TestRatingJSONSerialization(t *testing.T) {
	// 测试JSON序列化
	rating := model.Rating(4)
	data, err := json.Marshal(rating)
	assert.NoError(t, err)
	assert.Equal(t, "4", string(data))

	// 测试JSON反序列化
	var decoded model.Rating
	err = json.Unmarshal([]byte("3"), &decoded)
	assert.NoError(t, err)
	assert.Equal(t, model.Rating(3), decoded)

	// 测试无效值的反序列化
	var invalid model.Rating
	err = json.Unmarshal([]byte("6"), &invalid)
	assert.NoError(t, err) // JSON反序列化不会验证，只是设置值
	assert.Equal(t, model.Rating(6), invalid)
	assert.False(t, invalid.Valid()) // 但值是无效的
}

func TestRatingGormDataType(t *testing.T) {
	rating := model.Rating(0)
	dataType := rating.GormDataType()
	assert.Equal(t, "tinyint", dataType)
}

func TestRatingBoundaryValues(t *testing.T) {
	// 测试边界值
	minRating := model.RatingMin
	maxRating := model.RatingMax

	assert.True(t, minRating.Valid())
	assert.True(t, maxRating.Valid())

	// 测试边界外的值
	belowMin := model.Rating(0)
	aboveMax := model.Rating(6)

	assert.False(t, belowMin.Valid())
	assert.False(t, aboveMax.Valid())
}

func TestRatingZeroValue(t *testing.T) {
	var rating model.Rating
	assert.Equal(t, model.Rating(0), rating)
	assert.False(t, rating.Valid())
}

func TestRatingMaxValue(t *testing.T) {
	maxRating := model.Rating(255) // uint8的最大值
	assert.Equal(t, model.Rating(255), maxRating)
	assert.False(t, maxRating.Valid())
}

func TestRatingInStructs(t *testing.T) {
	// 测试在结构体中的使用
	type TestStruct struct {
		ID     uint64       `json:"id"`
		Rating model.Rating `json:"rating"`
	}

	testObj := TestStruct{
		ID:     1,
		Rating: model.Rating(5),
	}

	assert.Equal(t, uint64(1), testObj.ID)
	assert.Equal(t, model.Rating(5), testObj.Rating)
	assert.True(t, testObj.Rating.Valid())

	// 测试JSON序列化
	data, err := json.Marshal(testObj)
	assert.NoError(t, err)

	// 验证JSON结构
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	assert.Contains(t, result, "id")
	assert.Contains(t, result, "rating")
	assert.Equal(t, float64(1), result["id"])
	assert.Equal(t, float64(5), result["rating"])
}

func TestRatingConversions(t *testing.T) {
	// 测试从其他类型转换
	var uint8Value uint8 = 3
	rating := model.Rating(uint8Value)
	assert.Equal(t, model.Rating(3), rating)
	assert.True(t, rating.Valid())

	// 测试转换到字符串（注意：Rating是uint8，直接转换为字符串会得到ASCII字符）
	ratingStr := string(rating)
	// 由于Rating(3)对应ASCII码3，这是一个控制字符，我们验证长度即可
	assert.Len(t, ratingStr, 1)
}

func TestRatingComparison(t *testing.T) {
	rating1 := model.Rating(3)
	rating2 := model.Rating(4)
	rating3 := model.Rating(3)

	assert.True(t, rating1 < rating2)
	assert.True(t, rating2 > rating1)
	assert.True(t, rating1 == rating3)
	assert.False(t, rating1 == rating2)
}

func TestRatingArithmetic(t *testing.T) {
	rating1 := model.Rating(3)
	rating2 := model.Rating(4)

	// 加法
	sum := rating1 + rating2
	assert.Equal(t, model.Rating(7), sum)

	// 减法
	diff := rating2 - rating1
	assert.Equal(t, model.Rating(1), diff)

	// 乘法
	product := rating1 * 2
	assert.Equal(t, model.Rating(6), product)

	// 除法
	quotient := rating2 / 2
	assert.Equal(t, model.Rating(2), quotient)
}

func TestRatingPanicMessage(t *testing.T) {
	// 测试panic时的错误消息
	defer func() {
		if r := recover(); r != nil {
			assert.Contains(t, r.(string), "invalid rating: 10")
		}
	}()

	model.MustRating(10)
}

func TestRatingEdgeCases(t *testing.T) {
	// 测试所有可能的uint8值
	for i := 0; i < 256; i++ {
		rating := model.Rating(uint8(i))
		if i >= 1 && i <= 5 {
			assert.True(t, rating.Valid(), "Rating %d should be valid", i)
		} else {
			assert.False(t, rating.Valid(), "Rating %d should be invalid", i)
		}
	}
}
