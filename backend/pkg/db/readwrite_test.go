package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"gamelink/pkg/testutil"
)

func newTestDB(t *testing.T) *gorm.DB {
	return testutil.NewMemoryDB(t)
}

func TestNewReadWriteDB(t *testing.T) {
	writer := newTestDB(t)
	reader1 := newTestDB(t)
	reader2 := newTestDB(t)

	rw := NewReadWriteDB(writer, reader1, reader2)

	assert.NotNil(t, rw)
	assert.Equal(t, writer, rw.Writer())
	assert.Equal(t, 2, rw.ReaderCount())
}

func TestNewReadWriteDB_NoReaders(t *testing.T) {
	writer := newTestDB(t)

	rw := NewReadWriteDB(writer)

	assert.NotNil(t, rw)
	assert.Equal(t, writer, rw.Writer())
	assert.Equal(t, writer, rw.Reader()) // 没有从库时使用主库
	assert.Equal(t, 1, rw.ReaderCount())
}

func TestReadWriteDB_RoundRobin(t *testing.T) {
	writer := newTestDB(t)
	reader1 := newTestDB(t)
	reader2 := newTestDB(t)

	rw := NewReadWriteDB(writer, reader1, reader2)
	rw.SetStrategy(RoundRobin)

	// 轮询应该交替返回不同的 reader
	readers := make(map[*gorm.DB]int)
	for i := 0; i < 10; i++ {
		r := rw.Reader()
		readers[r]++
	}

	// 两个 reader 都应该被使用
	assert.Equal(t, 2, len(readers))
	assert.Equal(t, 5, readers[reader1])
	assert.Equal(t, 5, readers[reader2])
}

func TestReadWriteDB_WithRead(t *testing.T) {
	writer := newTestDB(t)
	reader := newTestDB(t)

	rw := NewReadWriteDB(writer, reader)

	called := false
	err := rw.WithRead(func(db *gorm.DB) error {
		called = true
		assert.Equal(t, reader, db)
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, called)
}

func TestReadWriteDB_WithWrite(t *testing.T) {
	writer := newTestDB(t)
	reader := newTestDB(t)

	rw := NewReadWriteDB(writer, reader)

	called := false
	err := rw.WithWrite(func(db *gorm.DB) error {
		called = true
		assert.Equal(t, writer, db)
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, called)
}

func TestReadWriteDB_AddReader(t *testing.T) {
	writer := newTestDB(t)
	reader1 := newTestDB(t)

	rw := NewReadWriteDB(writer, reader1)
	assert.Equal(t, 1, rw.ReaderCount())

	reader2 := newTestDB(t)
	rw.AddReader(reader2)
	assert.Equal(t, 2, rw.ReaderCount())
}

func TestReadWriteDB_RemoveReader(t *testing.T) {
	writer := newTestDB(t)
	reader1 := newTestDB(t)
	reader2 := newTestDB(t)

	rw := NewReadWriteDB(writer, reader1, reader2)
	assert.Equal(t, 2, rw.ReaderCount())

	err := rw.RemoveReader(0)
	assert.NoError(t, err)
	assert.Equal(t, 1, rw.ReaderCount())

	// 无效索引
	err = rw.RemoveReader(10)
	assert.Error(t, err)
}

func TestReadWriteDB_HealthCheck(t *testing.T) {
	writer := newTestDB(t)
	reader := newTestDB(t)

	rw := NewReadWriteDB(writer, reader)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	results := rw.HealthCheck(ctx)

	assert.Contains(t, results, "writer")
	assert.Nil(t, results["writer"])
}

func TestReadWriteDB_Stats(t *testing.T) {
	writer := newTestDB(t)
	reader := newTestDB(t)

	rw := NewReadWriteDB(writer, reader)

	stats := rw.Stats()

	assert.Contains(t, stats, "writer")
	writerStats := stats["writer"].(map[string]interface{})
	assert.Contains(t, writerStats, "max_open")
	assert.Contains(t, writerStats, "open")
	assert.Contains(t, writerStats, "in_use")
	assert.Contains(t, writerStats, "idle")
}

func TestReadWriteDB_Transaction(t *testing.T) {
	writer := newTestDB(t)

	// 创建测试表
	type TestModel struct {
		ID   uint
		Name string
	}
	writer.AutoMigrate(&TestModel{})

	rw := NewReadWriteDB(writer)

	err := rw.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&TestModel{Name: "test"}).Error
	})

	assert.NoError(t, err)

	var count int64
	writer.Model(&TestModel{}).Count(&count)
	assert.Equal(t, int64(1), count)
}
