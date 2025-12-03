// Package logging provides generic logging utilities.
package logging

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

/**
 * Test_WithRequestID_StoresID 测试存储请求ID
 */
func Test_WithRequestID_StoresID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	requestID := "req-123456"

	// Act
	ctx = WithRequestID(ctx, requestID)

	// Assert
	id, ok := RequestIDFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, requestID, id)
}

/**
 * Test_RequestIDFromContext_NotPresent 测试请求ID不存在
 */
func Test_RequestIDFromContext_NotPresent(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	id, ok := RequestIDFromContext(ctx)

	// Assert
	assert.False(t, ok)
	assert.Empty(t, id)
}

/**
 * Test_RequestIDFromContext_WrongType 测试错误的类型
 */
func Test_RequestIDFromContext_WrongType(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ctx = context.WithValue(ctx, keyRequestID, 12345) // 错误的类型

	// Act
	id, ok := RequestIDFromContext(ctx)

	// Assert
	assert.False(t, ok)
	assert.Empty(t, id)
}

/**
 * Test_WithActorUserID_StoresID 测试存储用户ID
 */
func Test_WithActorUserID_StoresID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uint64(12345)

	// Act
	ctx = WithActorUserID(ctx, userID)

	// Assert
	id, ok := ActorUserIDFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, userID, id)
}

/**
 * Test_ActorUserIDFromContext_NotPresent 测试用户ID不存在
 */
func Test_ActorUserIDFromContext_NotPresent(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	id, ok := ActorUserIDFromContext(ctx)

	// Assert
	assert.False(t, ok)
	assert.Equal(t, uint64(0), id)
}

/**
 * Test_ActorUserIDFromContext_WrongType 测试错误的类型
 */
func Test_ActorUserIDFromContext_WrongType(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ctx = context.WithValue(ctx, keyActorUserID, "wrong-type") // 错误的类型

	// Act
	id, ok := ActorUserIDFromContext(ctx)

	// Assert
	assert.False(t, ok)
	assert.Equal(t, uint64(0), id)
}

/**
 * Test_Context_BothValues 测试同时存储两个值
 */
func Test_Context_BothValues(t *testing.T) {
	// Arrange
	ctx := context.Background()
	requestID := "req-789"
	userID := uint64(999)

	// Act
	ctx = WithRequestID(ctx, requestID)
	ctx = WithActorUserID(ctx, userID)

	// Assert
	reqID, ok1 := RequestIDFromContext(ctx)
	assert.True(t, ok1)
	assert.Equal(t, requestID, reqID)

	actorID, ok2 := ActorUserIDFromContext(ctx)
	assert.True(t, ok2)
	assert.Equal(t, userID, actorID)
}

/**
 * Test_WithRequestID_UpdateValue 测试更新请求ID
 */
func Test_WithRequestID_UpdateValue(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ctx = WithRequestID(ctx, "old-id")

	// Act
	ctx = WithRequestID(ctx, "new-id")

	// Assert
	id, ok := RequestIDFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, "new-id", id)
}

/**
 * Test_WithActorUserID_ZeroValue 测试零值用户ID
 */
func Test_WithActorUserID_ZeroValue(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	ctx = WithActorUserID(ctx, 0)

	// Assert
	id, ok := ActorUserIDFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, uint64(0), id)
}
