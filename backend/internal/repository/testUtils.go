package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestDBConfig 测试数据库配置
type TestDBConfig struct {
	Models []any
}

// SetupTestDB 设置测试数据库
func SetupTestDB(t *testing.T, models []any) *gorm.DB {
	// 使用SQLite内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Skipf("sqlite driver unavailable: %v", err)
		return nil
	}

	// 迁移表结构
	if len(models) > 0 {
		err = db.AutoMigrate(models...)
		assert.NoError(t, err, "Failed to migrate test database")
	}

	return db
}

// CleanupTestDB 清理测试数据库
func CleanupTestDB(t *testing.T, db *gorm.DB) {
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

// ClearTable 清空指定表
func ClearTable(db *gorm.DB, tableName string) error {
	return db.Exec("DELETE FROM " + tableName).Error
}

// ClearTables 清空多个表
func ClearTables(db *gorm.DB, tableNames ...string) error {
	for _, tableName := range tableNames {
		if err := ClearTable(db, tableName); err != nil {
			return err
		}
	}
	return nil
}

// CreateTestContext 创建测试上下文
func CreateTestContext() context.Context {
	return context.Background()
}

// AssertRecordCount 断言记录数量
func AssertRecordCount(t *testing.T, db *gorm.DB, model any, expectedCount int64, where ...any) {
	var count int64
	err := db.Model(model).Where(where).Count(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count, "Record count mismatch")
}

// AssertRecordExists 断言记录存在
func AssertRecordExists(t *testing.T, db *gorm.DB, model any, where ...any) {
	var count int64
	err := db.Model(model).Where(where).Count(&count).Error
	assert.NoError(t, err)
	assert.Greater(t, count, int64(0), "Record should exist")
}

// AssertRecordNotExists 断言记录不存在
func AssertRecordNotExists(t *testing.T, db *gorm.DB, model any, where ...any) {
	var count int64
	err := db.Model(model).Where(where).Count(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count, "Record should not exist")
}

// PaginationTestCase 分页测试用例
type PaginationTestCase struct {
	Name         string
	Page         int
	PageSize     int
	ExpectedPage int
	ExpectedSize int
}

// GetPaginationTestCases 获取标准的分页测试用例
func GetPaginationTestCases() []PaginationTestCase {
	return []PaginationTestCase{
		{
			Name:         "negative page should use default",
			Page:         -1,
			PageSize:     10,
			ExpectedPage: 1,
			ExpectedSize: 10,
		},
		{
			Name:         "zero page should use default",
			Page:         0,
			PageSize:     10,
			ExpectedPage: 1,
			ExpectedSize: 10,
		},
		{
			Name:         "negative page size should use default",
			Page:         1,
			PageSize:     -1,
			ExpectedPage: 1,
			ExpectedSize: 20,
		},
		{
			Name:         "zero page size should use default",
			Page:         1,
			PageSize:     0,
			ExpectedPage: 1,
			ExpectedSize: 20,
		},
		{
			Name:         "page size exceeding max should be limited",
			Page:         1,
			PageSize:     150,
			ExpectedPage: 1,
			ExpectedSize: 100,
		},
		{
			Name:         "normal values should be preserved",
			Page:         2,
			PageSize:     15,
			ExpectedPage: 2,
			ExpectedSize: 15,
		},
	}
}

// ConcurrentOperation 并发操作
type ConcurrentOperation func() error

// RunConcurrentOperations 运行并发操作
func RunConcurrentOperations(t *testing.T, operations []ConcurrentOperation) (successCount, errorCount int) {
	done := make(chan bool, len(operations))
	errors := make(chan error, len(operations))

	// 启动所有操作
	for _, op := range operations {
		go func(operation ConcurrentOperation) {
			if err := operation(); err != nil {
				errors <- err
			} else {
				done <- true
			}
		}(op)
	}

	// 等待所有操作完成
	for i := 0; i < len(operations); i++ {
		select {
		case <-done:
			successCount++
		case <-errors:
			errorCount++
		}
	}

	return successCount, errorCount
}

// TestRepositorySuite 基础测试套件
type TestRepositorySuite struct {
	DB  *gorm.DB
	Ctx context.Context
}

// SetupSuite 设置测试套件
func (s *TestRepositorySuite) SetupSuite(t *testing.T, models []any) {
	s.DB = SetupTestDB(t, models)
	s.Ctx = CreateTestContext()
}

// TearDownSuite 清理测试套件
func (s *TestRepositorySuite) TearDownSuite(t *testing.T) {
	CleanupTestDB(t, s.DB)
}

// ClearTables 清空测试表
func (s *TestRepositorySuite) ClearTables(t *testing.T, tableNames ...string) {
	err := ClearTables(s.DB, tableNames...)
	assert.NoError(t, err, "Failed to clear tables")
}

// AssertRecordCount 断言记录数量
func (s *TestRepositorySuite) AssertRecordCount(t *testing.T, model any, expectedCount int64, where ...any) {
	AssertRecordCount(t, s.DB, model, expectedCount, where...)
}

// AssertRecordExists 断言记录存在
func (s *TestRepositorySuite) AssertRecordExists(t *testing.T, model any, where ...any) {
	AssertRecordExists(t, s.DB, model, where...)
}

// AssertRecordNotExists 断言记录不存在
func (s *TestRepositorySuite) AssertRecordNotExists(t *testing.T, model any, where ...any) {
	AssertRecordNotExists(t, s.DB, model, where...)
}
