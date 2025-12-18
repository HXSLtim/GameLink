package db

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"gorm.io/gorm"
)

// ReadWriteDB 读写分离数据库管理器
type ReadWriteDB struct {
	writer   *gorm.DB   // 主库（写）
	readers  []*gorm.DB // 从库（读）
	current  uint64     // 轮询计数器
	mu       sync.RWMutex
	strategy LoadBalanceStrategy
}

// LoadBalanceStrategy 负载均衡策略
type LoadBalanceStrategy int

const (
	// RoundRobin 轮询策略
	RoundRobin LoadBalanceStrategy = iota
	// Random 随机策略
	Random
	// LeastConn 最少连接（简化实现）
	LeastConn
)

// ReadWriteConfig 读写分离配置
type ReadWriteConfig struct {
	WriterDSN  string              // 主库 DSN
	ReaderDSNs []string            // 从库 DSN 列表
	Strategy   LoadBalanceStrategy // 负载均衡策略
	MaxRetries int                 // 最大重试次数
	RetryDelay time.Duration       // 重试延迟
}

// NewReadWriteDB 创建读写分离数据库管理器
func NewReadWriteDB(writer *gorm.DB, readers ...*gorm.DB) *ReadWriteDB {
	rw := &ReadWriteDB{
		writer:   writer,
		readers:  readers,
		strategy: RoundRobin,
	}

	// 如果没有从库，使用主库作为读库
	if len(readers) == 0 {
		rw.readers = []*gorm.DB{writer}
	}

	return rw
}

// Writer 获取写库连接
func (rw *ReadWriteDB) Writer() *gorm.DB {
	return rw.writer
}

// Reader 获取读库连接（负载均衡）
func (rw *ReadWriteDB) Reader() *gorm.DB {
	rw.mu.RLock()
	defer rw.mu.RUnlock()

	if len(rw.readers) == 0 {
		return rw.writer
	}

	if len(rw.readers) == 1 {
		return rw.readers[0]
	}

	switch rw.strategy {
	case RoundRobin:
		idx := atomic.AddUint64(&rw.current, 1) % uint64(len(rw.readers))
		return rw.readers[idx]
	case Random:
		idx := time.Now().UnixNano() % int64(len(rw.readers))
		return rw.readers[idx]
	default:
		return rw.readers[0]
	}
}

// DB 获取默认数据库连接（写库）
func (rw *ReadWriteDB) DB() *gorm.DB {
	return rw.writer
}

// WithRead 在读库上执行查询
func (rw *ReadWriteDB) WithRead(fn func(*gorm.DB) error) error {
	return fn(rw.Reader())
}

// WithWrite 在写库上执行操作
func (rw *ReadWriteDB) WithWrite(fn func(*gorm.DB) error) error {
	return fn(rw.Writer())
}

// Transaction 在写库上执行事务
func (rw *ReadWriteDB) Transaction(fn func(*gorm.DB) error) error {
	return rw.writer.Transaction(fn)
}

// AddReader 动态添加从库
func (rw *ReadWriteDB) AddReader(reader *gorm.DB) {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	rw.readers = append(rw.readers, reader)
}

// RemoveReader 动态移除从库
func (rw *ReadWriteDB) RemoveReader(idx int) error {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	if idx < 0 || idx >= len(rw.readers) {
		return errors.New("invalid reader index")
	}

	rw.readers = append(rw.readers[:idx], rw.readers[idx+1:]...)
	return nil
}

// SetStrategy 设置负载均衡策略
func (rw *ReadWriteDB) SetStrategy(strategy LoadBalanceStrategy) {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	rw.strategy = strategy
}

// ReaderCount 获取从库数量
func (rw *ReadWriteDB) ReaderCount() int {
	rw.mu.RLock()
	defer rw.mu.RUnlock()
	return len(rw.readers)
}

// HealthCheck 健康检查所有数据库连接
func (rw *ReadWriteDB) HealthCheck(ctx context.Context) map[string]error {
	results := make(map[string]error)

	// 检查主库
	if sqlDB, err := rw.writer.DB(); err != nil {
		results["writer"] = err
	} else if err := sqlDB.PingContext(ctx); err != nil {
		results["writer"] = err
	} else {
		results["writer"] = nil
	}

	// 检查从库
	rw.mu.RLock()
	defer rw.mu.RUnlock()

	for i, reader := range rw.readers {
		key := "reader_" + string(rune('0'+i))
		if sqlDB, err := reader.DB(); err != nil {
			results[key] = err
		} else if err := sqlDB.PingContext(ctx); err != nil {
			results[key] = err
		} else {
			results[key] = nil
		}
	}

	return results
}

// Close 关闭所有数据库连接
func (rw *ReadWriteDB) Close() error {
	var errs []error

	// 关闭主库
	if sqlDB, err := rw.writer.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	// 关闭从库
	rw.mu.Lock()
	defer rw.mu.Unlock()

	for _, reader := range rw.readers {
		if reader == rw.writer {
			continue // 跳过与主库相同的连接
		}
		if sqlDB, err := reader.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// Stats 获取连接池统计信息
func (rw *ReadWriteDB) Stats() map[string]interface{} {
	stats := make(map[string]interface{})

	// 主库统计
	if sqlDB, err := rw.writer.DB(); err == nil {
		s := sqlDB.Stats()
		stats["writer"] = map[string]interface{}{
			"max_open":   s.MaxOpenConnections,
			"open":       s.OpenConnections,
			"in_use":     s.InUse,
			"idle":       s.Idle,
			"wait_count": s.WaitCount,
			"wait_time":  s.WaitDuration.String(),
		}
	}

	// 从库统计
	rw.mu.RLock()
	defer rw.mu.RUnlock()

	readerStats := make([]map[string]interface{}, 0, len(rw.readers))
	for _, reader := range rw.readers {
		if reader == rw.writer {
			continue
		}
		if sqlDB, err := reader.DB(); err == nil {
			s := sqlDB.Stats()
			readerStats = append(readerStats, map[string]interface{}{
				"max_open":   s.MaxOpenConnections,
				"open":       s.OpenConnections,
				"in_use":     s.InUse,
				"idle":       s.Idle,
				"wait_count": s.WaitCount,
				"wait_time":  s.WaitDuration.String(),
			})
		}
	}
	stats["readers"] = readerStats

	return stats
}
