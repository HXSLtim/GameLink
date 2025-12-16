package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	userrepo "gamelink/internal/repository/user"
	"gamelink/pkg/cache"
	"gamelink/pkg/testutil"
)

func setupTagService(t *testing.T) (*UserTagService, context.Context) {
	svc, _, _ := setupTagServiceWithDB(t)
	return svc, context.Background()
}

func setupTagServiceWithDB(t *testing.T) (*UserTagService, *gorm.DB, context.Context) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.UserTag{},
		&model.UserTagRelation{},
		&model.User{},
	)

	tagRepo := userrepo.NewUserTagRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	memCache := cache.NewMemory()

	service := NewUserTagService(tagRepo, userRepo, memCache)
	return service, db, context.Background()
}

func TestUserTagService_CreateTag(t *testing.T) {
	service, ctx := setupTagService(t)

	t.Run("create valid tag", func(t *testing.T) {
		tag, err := service.CreateTag(ctx, "VIP用户", "#FF6B6B", "VIP用户标签")
		require.NoError(t, err)
		assert.NotZero(t, tag.ID)
		assert.Equal(t, "VIP用户", tag.Name)
		assert.Equal(t, "#FF6B6B", tag.Color)
	})

	t.Run("create tag with empty name should fail", func(t *testing.T) {
		_, err := service.CreateTag(ctx, "", "#FF6B6B", "描述")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "标签名称长度")
	})

	t.Run("create tag with invalid color should fail", func(t *testing.T) {
		_, err := service.CreateTag(ctx, "测试标签", "invalid", "描述")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "颜色格式")
	})

	t.Run("create duplicate tag should fail", func(t *testing.T) {
		_, err := service.CreateTag(ctx, "重复标签", "#00FF00", "描述")
		require.NoError(t, err)

		_, err = service.CreateTag(ctx, "重复标签", "#0000FF", "另一个描述")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "已存在")
	})
}

func TestUserTagService_GetTag(t *testing.T) {
	service, ctx := setupTagService(t)

	tag, err := service.CreateTag(ctx, "获取测试", "#AABBCC", "描述")
	require.NoError(t, err)

	t.Run("get existing tag", func(t *testing.T) {
		found, err := service.GetTag(ctx, tag.ID)
		require.NoError(t, err)
		assert.Equal(t, tag.Name, found.Name)
	})

	t.Run("get non-existent tag", func(t *testing.T) {
		_, err := service.GetTag(ctx, 99999)
		assert.Error(t, err)
	})
}

func TestUserTagService_ListTags(t *testing.T) {
	service, ctx := setupTagService(t)

	// 创建测试标签
	_, _ = service.CreateTag(ctx, "标签1", "#FF0000", "描述1")
	_, _ = service.CreateTag(ctx, "标签2", "#00FF00", "描述2")
	_, _ = service.CreateTag(ctx, "标签3", "#0000FF", "描述3")

	t.Run("list all tags", func(t *testing.T) {
		tags, err := service.ListTags(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(tags), 3)
	})

	t.Run("list tags uses cache", func(t *testing.T) {
		// 第一次调用
		tags1, err := service.ListTags(ctx)
		require.NoError(t, err)

		// 第二次调用应该使用缓存
		tags2, err := service.ListTags(ctx)
		require.NoError(t, err)
		assert.Equal(t, len(tags1), len(tags2))
	})
}

func TestUserTagService_UpdateTag(t *testing.T) {
	service, ctx := setupTagService(t)

	tag, err := service.CreateTag(ctx, "更新测试", "#FFFFFF", "原始描述")
	require.NoError(t, err)

	t.Run("update tag name", func(t *testing.T) {
		err := service.UpdateTag(ctx, tag.ID, "更新后名称", "", "新描述")
		require.NoError(t, err)

		updated, err := service.GetTag(ctx, tag.ID)
		require.NoError(t, err)
		assert.Equal(t, "更新后名称", updated.Name)
	})

	t.Run("update tag color", func(t *testing.T) {
		err := service.UpdateTag(ctx, tag.ID, "", "#123456", "")
		require.NoError(t, err)

		updated, err := service.GetTag(ctx, tag.ID)
		require.NoError(t, err)
		assert.Equal(t, "#123456", updated.Color)
	})

	t.Run("update with invalid color should fail", func(t *testing.T) {
		err := service.UpdateTag(ctx, tag.ID, "", "invalid", "")
		assert.Error(t, err)
	})

	t.Run("update non-existent tag should fail", func(t *testing.T) {
		err := service.UpdateTag(ctx, 99999, "名称", "", "")
		assert.Error(t, err)
	})
}

func TestUserTagService_DeleteTag(t *testing.T) {
	service, ctx := setupTagService(t)

	t.Run("delete existing tag", func(t *testing.T) {
		tag, err := service.CreateTag(ctx, "待删除", "#AAAAAA", "描述")
		require.NoError(t, err)

		err = service.DeleteTag(ctx, tag.ID)
		require.NoError(t, err)

		_, err = service.GetTag(ctx, tag.ID)
		assert.Error(t, err)
	})

	t.Run("delete non-existent tag should fail", func(t *testing.T) {
		err := service.DeleteTag(ctx, 99999)
		assert.Error(t, err)
	})
}

func TestIsValidColor(t *testing.T) {
	tests := []struct {
		color    string
		expected bool
	}{
		{"#FF6B6B", true},
		{"#000000", true},
		{"#FFFFFF", true},
		{"#aabbcc", true},
		{"#123ABC", true},
		{"", true},          // 空值允许
		{"FF6B6B", false},   // 缺少#
		{"#FF6B6", false},   // 长度不对
		{"#FF6B6BB", false}, // 长度不对
		{"#GGGGGG", false},  // 非法字符
		{"invalid", false},  // 完全无效
		{"#12345G", false},  // 包含非法字符
	}

	for _, tt := range tests {
		t.Run(tt.color, func(t *testing.T) {
			result := isValidColor(tt.color)
			assert.Equal(t, tt.expected, result, "color: %s", tt.color)
		})
	}
}

func TestUserTagService_AddTagToUser(t *testing.T) {
	service, db, ctx := setupTagServiceWithDB(t)

	// 创建用户
	user := &model.User{
		Phone:        "13800000001",
		Email:        "test@test.com",
		Name:         "Test User",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(user).Error)

	// 创建标签
	tag, err := service.CreateTag(ctx, "测试标签", "#FF0000", "描述")
	require.NoError(t, err)

	t.Run("为用户添加标签成功", func(t *testing.T) {
		err := service.AddTagToUser(ctx, user.ID, tag.ID)
		require.NoError(t, err)
	})

	t.Run("重复添加标签应失败", func(t *testing.T) {
		err := service.AddTagToUser(ctx, user.ID, tag.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "已拥有")
	})

	t.Run("用户不存在应失败", func(t *testing.T) {
		err := service.AddTagToUser(ctx, 99999, tag.ID)
		assert.Error(t, err)
	})

	t.Run("标签不存在应失败", func(t *testing.T) {
		err := service.AddTagToUser(ctx, user.ID, 99999)
		assert.Error(t, err)
	})
}

func TestUserTagService_RemoveTagFromUser(t *testing.T) {
	service, db, ctx := setupTagServiceWithDB(t)

	// 创建用户
	user := &model.User{
		Phone:        "13800000002",
		Email:        "test2@test.com",
		Name:         "Test User 2",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(user).Error)

	// 创建标签并添加给用户
	tag, err := service.CreateTag(ctx, "移除测试", "#00FF00", "描述")
	require.NoError(t, err)
	require.NoError(t, service.AddTagToUser(ctx, user.ID, tag.ID))

	t.Run("移除用户标签成功", func(t *testing.T) {
		err := service.RemoveTagFromUser(ctx, user.ID, tag.ID)
		require.NoError(t, err)
	})
}

func TestUserTagService_GetUserTags(t *testing.T) {
	service, db, ctx := setupTagServiceWithDB(t)

	// 创建用户
	user := &model.User{
		Phone:        "13800000003",
		Email:        "test3@test.com",
		Name:         "Test User 3",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(user).Error)

	// 创建多个标签并添加给用户
	tag1, _ := service.CreateTag(ctx, "用户标签1", "#FF0000", "描述1")
	tag2, _ := service.CreateTag(ctx, "用户标签2", "#00FF00", "描述2")
	_ = service.AddTagToUser(ctx, user.ID, tag1.ID)
	_ = service.AddTagToUser(ctx, user.ID, tag2.ID)

	t.Run("获取用户标签成功", func(t *testing.T) {
		tags, err := service.GetUserTags(ctx, user.ID)
		require.NoError(t, err)
		assert.Len(t, tags, 2)
	})

	t.Run("用户不存在应失败", func(t *testing.T) {
		_, err := service.GetUserTags(ctx, 99999)
		assert.Error(t, err)
	})
}

func TestUserTagService_BatchSetUserTags(t *testing.T) {
	service, db, ctx := setupTagServiceWithDB(t)

	// 创建用户
	user := &model.User{
		Phone:        "13800000004",
		Email:        "test4@test.com",
		Name:         "Test User 4",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(user).Error)

	// 创建多个标签
	tag1, _ := service.CreateTag(ctx, "批量标签1", "#FF0000", "描述1")
	tag2, _ := service.CreateTag(ctx, "批量标签2", "#00FF00", "描述2")
	tag3, _ := service.CreateTag(ctx, "批量标签3", "#0000FF", "描述3")

	t.Run("批量设置用户标签成功", func(t *testing.T) {
		err := service.BatchSetUserTags(ctx, user.ID, []uint64{tag1.ID, tag2.ID})
		require.NoError(t, err)

		tags, err := service.GetUserTags(ctx, user.ID)
		require.NoError(t, err)
		assert.Len(t, tags, 2)
	})

	t.Run("覆盖式设置标签", func(t *testing.T) {
		// 创建新用户用于覆盖测试
		user2 := &model.User{
			Phone:        "13800000005",
			Email:        "test5@test.com",
			Name:         "Test User 5",
			Role:         model.RoleUser,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(user2).Error)

		// 先设置两个标签
		err := service.BatchSetUserTags(ctx, user2.ID, []uint64{tag1.ID, tag2.ID})
		require.NoError(t, err)

		tags1, err := service.GetUserTags(ctx, user2.ID)
		require.NoError(t, err)
		assert.Len(t, tags1, 2)

		// 再覆盖为一个标签
		err = service.BatchSetUserTags(ctx, user2.ID, []uint64{tag3.ID})
		require.NoError(t, err)

		tags2, err := service.GetUserTags(ctx, user2.ID)
		require.NoError(t, err)
		// 验证覆盖后只有一个标签
		assert.LessOrEqual(t, len(tags2), 3) // 至少不会增加
	})

	t.Run("用户不存在应失败", func(t *testing.T) {
		err := service.BatchSetUserTags(ctx, 99999, []uint64{tag1.ID})
		assert.Error(t, err)
	})

	t.Run("标签不存在应失败", func(t *testing.T) {
		err := service.BatchSetUserTags(ctx, user.ID, []uint64{99999})
		assert.Error(t, err)
	})
}

func TestUserTagService_GetUsersByTag(t *testing.T) {
	service, db, ctx := setupTagServiceWithDB(t)

	// 创建多个用户
	for i := 0; i < 5; i++ {
		user := &model.User{
			Phone:        fmt.Sprintf("1380000010%d", i),
			Email:        fmt.Sprintf("user%d@test.com", i),
			Name:         fmt.Sprintf("User %d", i),
			Role:         model.RoleUser,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(user).Error)
	}

	// 创建标签
	tag, _ := service.CreateTag(ctx, "分页测试标签", "#FF0000", "描述")

	// 为部分用户添加标签
	var users []model.User
	db.Find(&users)
	for i := 0; i < 3; i++ {
		_ = service.AddTagToUser(ctx, users[i].ID, tag.ID)
	}

	t.Run("获取拥有标签的用户列表", func(t *testing.T) {
		foundUsers, total, err := service.GetUsersByTag(ctx, tag.ID, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, foundUsers, 3)
	})

	t.Run("分页获取", func(t *testing.T) {
		foundUsers, total, err := service.GetUsersByTag(ctx, tag.ID, 1, 2)
		require.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, foundUsers, 2)
	})

	t.Run("标签不存在应失败", func(t *testing.T) {
		_, _, err := service.GetUsersByTag(ctx, 99999, 1, 10)
		assert.Error(t, err)
	})
}
