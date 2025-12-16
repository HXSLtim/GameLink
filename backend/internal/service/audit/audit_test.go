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

func TestServiceQuery(t *testing.T) {
	repo := newMockRepository()
	svc := NewServiceWithDefaults(repo)

	// Add some test logs directly to mock
	for i := 0; i < 5; i++ {
		repo.logs = append(repo.logs, &model.PermissionAuditLog{
			ID:           uint64(i + 1),
			OperatorID:   1,
			OperatorName: "admin",
			TargetType:   model.AuditTargetTypePermission,
			TargetID:     uint64(i + 1),
			TargetName:   "test_permission",
			Action:       model.AuditActionCreate,
		})
	}

	ctx := context.Background()

	t.Run("query with default pagination", func(t *testing.T) {
		result, err := svc.Query(ctx, QueryOptions{})
		require.NoError(t, err)
		assert.Equal(t, int64(5), result.Total)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 20, result.PageSize)
	})

	t.Run("query with custom pagination", func(t *testing.T) {
		result, err := svc.Query(ctx, QueryOptions{
			Page:     2,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, 2, result.Page)
		assert.Equal(t, 10, result.PageSize)
	})

	t.Run("query with page size limit", func(t *testing.T) {
		result, err := svc.Query(ctx, QueryOptions{
			PageSize: 200, // Exceeds max
		})
		require.NoError(t, err)
		assert.Equal(t, 100, result.PageSize) // Should be capped at 100
	})
}

func TestServiceQueryByOperator(t *testing.T) {
	repo := newMockRepository()
	svc := NewServiceWithDefaults(repo)

	ctx := context.Background()
	result, err := svc.QueryByOperator(ctx, 1, 1, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestServiceQueryByAction(t *testing.T) {
	repo := newMockRepository()
	svc := NewServiceWithDefaults(repo)

	ctx := context.Background()
	result, err := svc.QueryByAction(ctx, model.AuditActionCreate, 1, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestServiceQueryByDateRange(t *testing.T) {
	repo := newMockRepository()
	svc := NewServiceWithDefaults(repo)

	ctx := context.Background()
	from := time.Now().Add(-24 * time.Hour)
	to := time.Now()
	result, err := svc.QueryByDateRange(ctx, from, to, 1, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestServiceQueryByTarget(t *testing.T) {
	repo := newMockRepository()
	svc := NewServiceWithDefaults(repo)

	ctx := context.Background()
	result, err := svc.QueryByTarget(ctx, model.AuditTargetTypePermission, 1, 1, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestServiceGetByID(t *testing.T) {
	repo := newMockRepository()
	repo.logs = append(repo.logs, &model.PermissionAuditLog{
		ID:         1,
		OperatorID: 1,
		Action:     model.AuditActionCreate,
	})
	svc := NewServiceWithDefaults(repo)

	ctx := context.Background()

	t.Run("get existing log", func(t *testing.T) {
		log, err := svc.GetByID(ctx, 1)
		require.NoError(t, err)
		assert.NotNil(t, log)
		assert.Equal(t, uint64(1), log.ID)
	})

	t.Run("get non-existing log", func(t *testing.T) {
		log, err := svc.GetByID(ctx, 999)
		require.NoError(t, err)
		assert.Nil(t, log)
	})
}

func TestServiceExportCSV(t *testing.T) {
	repo := newMockRepository()
	now := time.Now()
	repo.logs = append(repo.logs, &model.PermissionAuditLog{
		ID:           1,
		OperatorID:   1,
		OperatorName: "admin",
		TargetType:   model.AuditTargetTypePermission,
		TargetID:     100,
		TargetName:   "test_permission",
		Action:       model.AuditActionCreate,
		BeforeData:   "{}",
		AfterData:    `{"name":"test"}`,
		IPAddress:    "192.168.1.1",
		UserAgent:    "Mozilla/5.0",
		RequestID:    "req-123",
		CreatedAt:    now,
	})
	svc := NewServiceWithDefaults(repo)

	ctx := context.Background()

	t.Run("export with default options", func(t *testing.T) {
		data, err := svc.ExportCSV(ctx, ExportOptions{})
		require.NoError(t, err)
		assert.NotEmpty(t, data)
		// Check UTF-8 BOM
		assert.Equal(t, byte(0xEF), data[0])
		assert.Equal(t, byte(0xBB), data[1])
		assert.Equal(t, byte(0xBF), data[2])
		// Check header exists
		assert.Contains(t, string(data), "ID")
		assert.Contains(t, string(data), "操作者ID")
	})

	t.Run("export with max records", func(t *testing.T) {
		data, err := svc.ExportCSV(ctx, ExportOptions{MaxRecords: 100})
		require.NoError(t, err)
		assert.NotEmpty(t, data)
	})
}

func TestGenerateExportFilename(t *testing.T) {
	filename := GenerateExportFilename("audit_export")
	assert.Contains(t, filename, "audit_export_")
	assert.Contains(t, filename, ".csv")
}

func TestServiceGetRepository(t *testing.T) {
	repo := newMockRepository()
	svc := NewServiceWithDefaults(repo)

	gotRepo := svc.GetRepository()
	assert.Equal(t, repo, gotRepo)
}

func TestServiceStopWithContext(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, Config{
		BufferSize:    100,
		BatchSize:     10,
		FlushInterval: 50 * time.Millisecond,
	})

	svc.Start()

	// Log some entries
	for i := 0; i < 3; i++ {
		svc.Log(&model.PermissionAuditLog{
			OperatorID: 1,
			TargetID:   uint64(i + 1),
			Action:     model.AuditActionCreate,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := svc.StopWithContext(ctx)
	require.NoError(t, err)
	assert.False(t, svc.IsRunning())

	// Verify logs were flushed
	logs := repo.getLogs()
	assert.Equal(t, 3, len(logs))
}

func TestServiceStopWithContextTimeout(t *testing.T) {
	repo := newMockRepository()
	svc := NewServiceWithDefaults(repo)

	// Stop without starting - should return nil immediately
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	err := svc.StopWithContext(ctx)
	require.NoError(t, err)
}

func TestServiceBufferFull(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo, Config{
		BufferSize:    2, // Very small buffer
		BatchSize:     100,
		FlushInterval: 10 * time.Second, // Long interval
	})

	svc.Start()
	defer svc.Stop()

	// Fill the buffer and overflow
	for i := 0; i < 10; i++ {
		svc.Log(&model.PermissionAuditLog{
			OperatorID: 1,
			TargetID:   uint64(i + 1),
			Action:     model.AuditActionCreate,
		})
	}

	// Wait a bit
	time.Sleep(50 * time.Millisecond)

	stats := svc.GetStats()
	assert.True(t, stats.DroppedCount > 0, "Expected some logs to be dropped")
}

func TestServiceCleanupOldLogs(t *testing.T) {
	repo := newMockRepository()
	svc := NewServiceWithDefaults(repo)

	ctx := context.Background()

	t.Run("cleanup with default retention", func(t *testing.T) {
		count, err := svc.CleanupOldLogs(ctx, 0)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("cleanup with custom retention", func(t *testing.T) {
		count, err := svc.CleanupOldLogs(ctx, 30)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestServiceArchiveOldLogs(t *testing.T) {
	repo := newMockRepository()
	svc := NewServiceWithDefaults(repo)

	ctx := context.Background()

	result, err := svc.ArchiveOldLogs(ctx, "")
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotZero(t, result.ArchivedAt)
}

func TestServiceGetRetentionStats(t *testing.T) {
	repo := newMockRepository()
	svc := NewServiceWithDefaults(repo)

	ctx := context.Background()

	stats, err := svc.GetRetentionStats(ctx)
	require.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, model.AuditLogRetentionDays, stats.OnlineRetention)
	assert.Equal(t, model.AuditLogArchiveDays, stats.ArchiveRetention)
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	assert.Equal(t, DefaultBufferSize, config.BufferSize)
	assert.Equal(t, DefaultBatchSize, config.BatchSize)
	assert.Equal(t, DefaultFlushInterval, config.FlushInterval)
}
