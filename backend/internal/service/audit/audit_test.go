package audit

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository/permissionauditlog"
)

// mockRepository is a mock implementation of permissionauditlog.Repository for testing.
type mockRepository struct {
	mu      sync.Mutex
	logs    []*model.PermissionAuditLog
	batches []int // Track batch sizes
	err     error
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		logs:    make([]*model.PermissionAuditLog, 0),
		batches: make([]int, 0),
	}
}

func (m *mockRepository) Create(ctx context.Context, log *model.PermissionAuditLog) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.err != nil {
		return m.err
	}
	m.logs = append(m.logs, log)
	m.batches = append(m.batches, 1)
	return nil
}

func (m *mockRepository) CreateBatch(ctx context.Context, logs []*model.PermissionAuditLog) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.err != nil {
		return m.err
	}
	m.logs = append(m.logs, logs...)
	m.batches = append(m.batches, len(logs))
	return nil
}

func (m *mockRepository) List(ctx context.Context, opts permissionauditlog.ListOptions) ([]model.PermissionAuditLog, int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]model.PermissionAuditLog, len(m.logs))
	for i, log := range m.logs {
		result[i] = *log
	}
	return result, int64(len(result)), nil
}

func (m *mockRepository) Get(ctx context.Context, id uint64) (*model.PermissionAuditLog, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, log := range m.logs {
		if log.ID == id {
			return log, nil
		}
	}
	return nil, nil
}

func (m *mockRepository) DeleteBefore(ctx context.Context, before time.Time) (int64, error) {
	return 0, nil
}

func (m *mockRepository) CountByDateRange(ctx context.Context, from, to time.Time) (int64, error) {
	return int64(len(m.logs)), nil
}

func (m *mockRepository) getLogs() []*model.PermissionAuditLog {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]*model.PermissionAuditLog, len(m.logs))
	copy(result, m.logs)
	return result
}

func (m *mockRepository) getBatches() []int {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]int, len(m.batches))
	copy(result, m.batches)
	return result
}

func TestNewService(t *testing.T) {
	repo := newMockRepository()

	t.Run("with default config", func(t *testing.T) {
		svc := NewServiceWithDefaults(repo)
		assert.NotNil(t, svc)
		assert.Equal(t, DefaultBufferSize, svc.config.BufferSize)
		assert.Equal(t, DefaultBatchSize, svc.config.BatchSize)
		assert.Equal(t, DefaultFlushInterval, svc.config.FlushInterval)
	})

	t.Run("with custom config", func(t *testing.T) {
		config := Config{
			BufferSize:    500,
			BatchSize:     25,
			FlushInterval: 10 * time.Second,
		}
		svc := NewService(repo, config)
		assert.NotNil(t, svc)
		assert.Equal(t, 500, svc.config.BufferSize)
		assert.Equal(t, 25, svc.config.BatchSize)
		assert.Equal(t, 10*time.Second, svc.config.FlushInterval)
	})

	t.Run("with zero values uses defaults", func(t *testing.T) {
		config := Config{}
		svc := NewService(repo, config)
		assert.Equal(t, DefaultBufferSize, svc.config.BufferSize)
		assert.Equal(t, DefaultBatchSize, svc.config.BatchSize)
		assert.Equal(t, DefaultFlushInterval, svc.config.FlushInterval)
	})
}

func TestServiceStartStop(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, Config{
		BufferSize:    100,
		BatchSize:     10,
		FlushInterval: 100 * time.Millisecond,
	})

	// Initially not running
	assert.False(t, svc.IsRunning())

	// Start the service
	svc.Start()
	assert.True(t, svc.IsRunning())

	// Starting again should be a no-op
	svc.Start()
	assert.True(t, svc.IsRunning())

	// Stop the service
	svc.Stop()
	assert.False(t, svc.IsRunning())

	// Stopping again should be a no-op
	svc.Stop()
	assert.False(t, svc.IsRunning())
}

func TestServiceLog(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, Config{
		BufferSize:    100,
		BatchSize:     5,
		FlushInterval: 50 * time.Millisecond,
	})

	svc.Start()
	defer svc.Stop()

	// Log some entries
	for i := 0; i < 10; i++ {
		svc.Log(&model.PermissionAuditLog{
			OperatorID:   1,
			OperatorName: "admin",
			TargetType:   model.AuditTargetTypePermission,
			TargetID:     uint64(i + 1),
			TargetName:   "test_permission",
			Action:       model.AuditActionCreate,
		})
	}

	// Wait for flush
	time.Sleep(200 * time.Millisecond)

	// Verify logs were written
	logs := repo.getLogs()
	assert.Equal(t, 10, len(logs))

	// Verify batching occurred
	batches := repo.getBatches()
	assert.True(t, len(batches) >= 1, "Expected at least one batch")
}

func TestServiceLogWhenNotRunning(t *testing.T) {
	repo := newMockRepository()
	svc := NewServiceWithDefaults(repo)

	// Log without starting - should be dropped
	svc.Log(&model.PermissionAuditLog{
		OperatorID: 1,
		Action:     model.AuditActionCreate,
	})

	// No logs should be written
	logs := repo.getLogs()
	assert.Equal(t, 0, len(logs))
}

func TestServiceBatchFlush(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, Config{
		BufferSize:    100,
		BatchSize:     5,
		FlushInterval: 1 * time.Second, // Long interval to test batch size trigger
	})

	svc.Start()
	defer svc.Stop()

	// Log exactly batch size entries
	for i := 0; i < 5; i++ {
		svc.Log(&model.PermissionAuditLog{
			OperatorID: 1,
			TargetID:   uint64(i + 1),
			Action:     model.AuditActionCreate,
		})
	}

	// Wait a bit for batch to be processed
	time.Sleep(100 * time.Millisecond)

	logs := repo.getLogs()
	assert.Equal(t, 5, len(logs))
}

func TestServiceGracefulShutdown(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, Config{
		BufferSize:    100,
		BatchSize:     100,              // Large batch size
		FlushInterval: 10 * time.Second, // Long interval
	})

	svc.Start()

	// Log some entries
	for i := 0; i < 7; i++ {
		svc.Log(&model.PermissionAuditLog{
			OperatorID: 1,
			TargetID:   uint64(i + 1),
			Action:     model.AuditActionCreate,
		})
	}

	// Stop should flush remaining logs
	svc.Stop()

	logs := repo.getLogs()
	assert.Equal(t, 7, len(logs))
}

func TestServiceStats(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, Config{
		BufferSize:    100,
		BatchSize:     5,
		FlushInterval: 50 * time.Millisecond,
	})

	// Check initial stats
	stats := svc.GetStats()
	assert.False(t, stats.Running)
	assert.Equal(t, 100, stats.BufferSize)
	assert.Equal(t, 0, stats.BufferUsed)
	assert.Equal(t, int64(0), stats.ProcessedCount)
	assert.Equal(t, int64(0), stats.DroppedCount)

	svc.Start()
	defer svc.Stop()

	// Log some entries
	for i := 0; i < 10; i++ {
		svc.Log(&model.PermissionAuditLog{
			OperatorID: 1,
			TargetID:   uint64(i + 1),
			Action:     model.AuditActionCreate,
		})
	}

	// Wait for processing
	time.Sleep(200 * time.Millisecond)

	stats = svc.GetStats()
	assert.True(t, stats.Running)
	assert.Equal(t, int64(10), stats.ProcessedCount)
}

func TestLogPermissionChange(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, Config{
		BufferSize:    100,
		BatchSize:     10,
		FlushInterval: 50 * time.Millisecond,
	})

	svc.Start()
	defer svc.Stop()

	ctx := context.Background()
	beforeData := map[string]string{"name": "old_name"}
	afterData := map[string]string{"name": "new_name"}

	svc.LogPermissionChange(ctx, 1, "admin", model.AuditActionUpdate,
		100, "test_permission", beforeData, afterData,
		"192.168.1.1", "Mozilla/5.0", "req-123")

	// Wait for flush
	time.Sleep(100 * time.Millisecond)

	logs := repo.getLogs()
	require.Equal(t, 1, len(logs))

	log := logs[0]
	assert.Equal(t, uint64(1), log.OperatorID)
	assert.Equal(t, "admin", log.OperatorName)
	assert.Equal(t, model.AuditTargetTypePermission, log.TargetType)
	assert.Equal(t, uint64(100), log.TargetID)
	assert.Equal(t, "test_permission", log.TargetName)
	assert.Equal(t, model.AuditActionUpdate, log.Action)
	assert.Equal(t, "192.168.1.1", log.IPAddress)
	assert.Equal(t, "Mozilla/5.0", log.UserAgent)
	assert.Equal(t, "req-123", log.RequestID)
	assert.Contains(t, log.BeforeData, "old_name")
	assert.Contains(t, log.AfterData, "new_name")
}

func TestLogRoleChange(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, Config{
		BufferSize:    100,
		BatchSize:     10,
		FlushInterval: 50 * time.Millisecond,
	})

	svc.Start()
	defer svc.Stop()

	ctx := context.Background()
	svc.LogRoleChange(ctx, 1, "admin", model.AuditActionAssign,
		200, "test_role", nil, []uint64{1, 2, 3},
		"192.168.1.1", "Mozilla/5.0", "req-456")

	// Wait for flush
	time.Sleep(100 * time.Millisecond)

	logs := repo.getLogs()
	require.Equal(t, 1, len(logs))

	log := logs[0]
	assert.Equal(t, model.AuditTargetTypeRole, log.TargetType)
	assert.Equal(t, uint64(200), log.TargetID)
	assert.Equal(t, model.AuditActionAssign, log.Action)
}

func TestLogUserRoleChange(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, Config{
		BufferSize:    100,
		BatchSize:     10,
		FlushInterval: 50 * time.Millisecond,
	})

	svc.Start()
	defer svc.Stop()

	ctx := context.Background()
	svc.LogUserRoleChange(ctx, 1, "admin", model.AuditActionBatchAssign,
		300, "test_user", []string{"role1"}, []string{"role1", "role2"},
		"192.168.1.1", "Mozilla/5.0", "req-789")

	// Wait for flush
	time.Sleep(100 * time.Millisecond)

	logs := repo.getLogs()
	require.Equal(t, 1, len(logs))

	log := logs[0]
	assert.Equal(t, model.AuditTargetTypeUser, log.TargetType)
	assert.Equal(t, uint64(300), log.TargetID)
	assert.Equal(t, model.AuditActionBatchAssign, log.Action)
}
