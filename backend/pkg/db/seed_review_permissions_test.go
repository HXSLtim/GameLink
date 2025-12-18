package db

import (
	"testing"

	"gamelink/internal/model"
	"gamelink/pkg/testutil"
)

func TestSeedReviewPermissions(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)

	// 迁移权限表
	if err := db.AutoMigrate(&model.Permission{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	// 执行种子数据函数
	if err := seedReviewPermissions(db); err != nil {
		t.Fatalf("seedReviewPermissions failed: %v", err)
	}

	// 验证权限已创建 - 使用中文分组名
	var count int64
	if err := db.Model(&model.Permission{}).Where("permissions.\"group\" = ?", "评价管理").Count(&count).Error; err != nil {
		t.Fatalf("failed to count permissions: %v", err)
	}

	if count == 0 {
		t.Fatal("expected review permissions to be created, but got 0")
	}

	t.Logf("Successfully created %d review permissions", count)

	// 验证特定权限存在 - 使用实际的权限码
	testCases := []struct {
		code        string
		method      model.HTTPMethod
		path        string
		description string
		group       string
	}{
		{"admin.reviews.list", model.HTTPMethodGET, "/api/v1/admin/reviews", "查看评价列表", "评价管理"},
		{"admin.reviews.approve.update", model.HTTPMethodPUT, "/api/v1/admin/reviews/:id/approve", "批准评价", "评价管理"},
		{"admin.reviews.delete", model.HTTPMethodDELETE, "/api/v1/admin/reviews/:id", "删除评价", "评价管理"},
		{"admin.sensitive-words.create", model.HTTPMethodPOST, "/api/v1/admin/sensitive-words", "添加敏感词", "敏感词管理"},
	}

	for _, tc := range testCases {
		var perm model.Permission
		err := db.Where("code = ? AND method = ? AND path = ?", tc.code, tc.method, tc.path).First(&perm).Error
		if err != nil {
			t.Errorf("permission not found: code=%s, method=%s, path=%s, error=%v", tc.code, tc.method, tc.path, err)
			continue
		}

		if perm.Description != tc.description {
			t.Errorf("permission description mismatch: expected=%s, got=%s", tc.description, perm.Description)
		}

		if perm.Group != tc.group {
			t.Errorf("permission group mismatch: expected=%s, got=%s", tc.group, perm.Group)
		}
	}

	// 测试幂等性 - 再次执行应该不会出错
	if err := seedReviewPermissions(db); err != nil {
		t.Fatalf("seedReviewPermissions should be idempotent, but failed on second call: %v", err)
	}
}
