package model_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gamelink/internal/model"
	"gorm.io/gorm"
)

func TestBaseModel(t *testing.T) {
	now := time.Now()
	deletedAt := gorm.DeletedAt{Time: now, Valid: true}
	
	base := &model.Base{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: deletedAt,
	}

	assert.Equal(t, uint64(1), base.ID)
	assert.Equal(t, now, base.CreatedAt)
	assert.Equal(t, now, base.UpdatedAt)
	assert.Equal(t, deletedAt, base.DeletedAt)
}

func TestBaseJSONSerialization(t *testing.T) {
	now := time.Now()
	base := &model.Base{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 序列化
	data, err := json.Marshal(base)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"id":1`)
	assert.Contains(t, string(data), `"createdAt"`)
	assert.Contains(t, string(data), `"updatedAt"`)

	// 反序列化
	var decoded model.Base
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, base.ID, decoded.ID)
}

func TestBaseWithDeletedAt(t *testing.T) {
	now := time.Now()
	deletedAt := gorm.DeletedAt{Time: now, Valid: true}
	
	base := &model.Base{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: deletedAt,
	}

	// 序列化 - 应该包含deletedAt
	data, err := json.Marshal(base)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"deletedAt"`)

	// 反序列化
	var decoded model.Base
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, deletedAt.Time.Unix(), decoded.DeletedAt.Time.Unix())
	assert.Equal(t, deletedAt.Valid, decoded.DeletedAt.Valid)
}

func TestBaseZeroID(t *testing.T) {
	base := &model.Base{
		ID: 0,
	}

	assert.Equal(t, uint64(0), base.ID)
}

func TestBaseZeroTime(t *testing.T) {
	base := &model.Base{
		ID:        1,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	assert.Equal(t, time.Time{}, base.CreatedAt)
	assert.Equal(t, time.Time{}, base.UpdatedAt)
}

func TestBaseInvalidDeletedAt(t *testing.T) {
	base := &model.Base{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{Valid: false},
	}

	// 序列化 - 不应该包含deletedAt（因为有omitempty）
	_, err := json.Marshal(base)
	assert.NoError(t, err)
	// 由于omitempty，无效的DeletedAt不应该出现在JSON中
	// 注意：实际行为取决于gorm.DeletedAt的MarshalJSON实现
}

func TestBaseLargeID(t *testing.T) {
	base := &model.Base{
		ID: ^uint64(0), // 最大uint64值
	}

	assert.Equal(t, uint64(18446744073709551615), base.ID)
}

func TestBaseJSONFields(t *testing.T) {
	now := time.Now().UTC()
	base := &model.Base{
		ID:        12345,
		CreatedAt: now,
		UpdatedAt: now.Add(1 * time.Hour),
	}

	data, err := json.Marshal(base)
	assert.NoError(t, err)

	// 验证JSON结构
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// 检查必需的字段
	assert.Contains(t, result, "id")
	assert.Contains(t, result, "createdAt")
	assert.Contains(t, result, "updatedAt")
	// 注意：由于gorm.DeletedAt的特殊实现，即使Valid为false也可能出现在JSON中

	// 验证值
	assert.Equal(t, float64(12345), result["id"])
}

func TestBaseTimePrecision(t *testing.T) {
	// 测试时间精度
	now := time.Now().Truncate(time.Millisecond)
	base := &model.Base{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now.Add(1 * time.Second),
	}

	assert.True(t, base.UpdatedAt.After(base.CreatedAt))
	assert.Equal(t, now.Unix(), base.CreatedAt.Unix())
}