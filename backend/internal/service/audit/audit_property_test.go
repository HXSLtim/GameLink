package audit

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"

	"gamelink/internal/model"
	"gamelink/internal/repository/permissionauditlog"
)

// propertyMockRepository is a mock repository for property testing that tracks all operations.
type propertyMockRepository struct {
	mu      sync.Mutex
	logs    []*model.PermissionAuditLog
	batches []int
}

func newPropertyMockRepository() *propertyMockRepository {
	return &propertyMockRepository{
		logs:    make([]*model.PermissionAuditLog, 0),
		batches: make([]int, 0),
	}
}

func (m *propertyMockRepository) Create(ctx context.Context, log *model.PermissionAuditLog) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logs = append(m.logs, log)
	m.batches = append(m.batches, 1)
	return nil
}

func (m *propertyMockRepository) CreateBatch(ctx context.Context, logs []*model.PermissionAuditLog) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logs = append(m.logs, logs...)
	m.batches = append(m.batches, len(logs))
	return nil
}

func (m *propertyMockRepository) List(ctx context.Context, opts permissionauditlog.ListOptions) ([]model.PermissionAuditLog, int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Apply filters
	var filtered []model.PermissionAuditLog
	for _, log := range m.logs {
		if opts.OperatorID != nil && log.OperatorID != *opts.OperatorID {
			continue
		}
		if opts.TargetType != nil && log.TargetType != *opts.TargetType {
			continue
		}
		if opts.TargetID != nil && log.TargetID != *opts.TargetID {
			continue
		}
		if opts.Action != nil && log.Action != *opts.Action {
			continue
		}
		if opts.DateFrom != nil && log.CreatedAt.Before(*opts.DateFrom) {
			continue
		}
		if opts.DateTo != nil && log.CreatedAt.After(*opts.DateTo) {
			continue
		}
		if opts.RequestID != "" && log.RequestID != opts.RequestID {
			continue
		}
		if opts.IPAddress != "" && log.IPAddress != opts.IPAddress {
			continue
		}
		filtered = append(filtered, *log)
	}

	// Apply pagination
	page := opts.Page
	if page < 1 {
		page = 1
	}
	pageSize := opts.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	total := int64(len(filtered))
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(filtered) {
		return []model.PermissionAuditLog{}, total, nil
	}
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], total, nil
}

func (m *propertyMockRepository) Get(ctx context.Context, id uint64) (*model.PermissionAuditLog, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, log := range m.logs {
		if log.ID == id {
			return log, nil
		}
	}
	return nil, nil
}

func (m *propertyMockRepository) DeleteBefore(ctx context.Context, before time.Time) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var remaining []*model.PermissionAuditLog
	var deleted int64
	for _, log := range m.logs {
		if log.CreatedAt.Before(before) {
			deleted++
		} else {
			remaining = append(remaining, log)
		}
	}
	m.logs = remaining
	return deleted, nil
}

func (m *propertyMockRepository) CountByDateRange(ctx context.Context, from, to time.Time) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var count int64
	for _, log := range m.logs {
		if !log.CreatedAt.Before(from) && !log.CreatedAt.After(to) {
			count++
		}
	}
	return count, nil
}

func (m *propertyMockRepository) getLogs() []*model.PermissionAuditLog {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]*model.PermissionAuditLog, len(m.logs))
	copy(result, m.logs)
	return result
}

// TestAuditLogCompleteness tests Property 11: Audit Log Completeness
// For any permission or role modification operation, an audit log entry should be created
// containing the operator ID, timestamp, target information, and the before/after values.
// **Feature: rbac-button-level-permission, Property 11: 审计日志完整性**
// **Validates: Requirements 6.1, 6.2**
func TestAuditLogCompleteness(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property 1: Permission change logs should contain all required fields
	properties.Property("permission change logs should contain all required fields", prop.ForAll(
		func(operatorID uint64, operatorName, permissionName, ipAddress, userAgent, requestID string) bool {
			repo := newPropertyMockRepository()
			svc := NewService(repo, Config{
				BufferSize:    100,
				BatchSize:     10,
				FlushInterval: 50 * time.Millisecond,
			})
			svc.Start()
			defer svc.Stop()

			ctx := context.Background()
			beforeData := map[string]string{"name": "old"}
			afterData := map[string]string{"name": "new"}

			svc.LogPermissionChange(ctx, operatorID, operatorName, model.AuditActionUpdate,
				100, permissionName, beforeData, afterData, ipAddress, userAgent, requestID)

			// Wait for flush
			time.Sleep(150 * time.Millisecond)

			logs := repo.getLogs()
			if len(logs) != 1 {
				return false
			}

			log := logs[0]
			// Verify all required fields are present
			return log.OperatorID == operatorID &&
				log.OperatorName == operatorName &&
				log.TargetType == model.AuditTargetTypePermission &&
				log.TargetID == 100 &&
				log.TargetName == permissionName &&
				log.Action == model.AuditActionUpdate &&
				log.IPAddress == ipAddress &&
				log.UserAgent == userAgent &&
				log.RequestID == requestID &&
				strings.Contains(log.BeforeData, "old") &&
				strings.Contains(log.AfterData, "new")
		},
		gen.UInt64(),
		genNonEmptyString(),
		genNonEmptyString(),
		genIPAddress(),
		genUserAgent(),
		genRequestID(),
	))

	// Property 2: Role change logs should contain all required fields
	properties.Property("role change logs should contain all required fields", prop.ForAll(
		func(operatorID uint64, operatorName, roleName, ipAddress, userAgent, requestID string) bool {
			repo := newPropertyMockRepository()
			svc := NewService(repo, Config{
				BufferSize:    100,
				BatchSize:     10,
				FlushInterval: 50 * time.Millisecond,
			})
			svc.Start()
			defer svc.Stop()

			ctx := context.Background()
			beforeData := []uint64{1, 2}
			afterData := []uint64{1, 2, 3}

			svc.LogRoleChange(ctx, operatorID, operatorName, model.AuditActionAssign,
				200, roleName, beforeData, afterData, ipAddress, userAgent, requestID)

			// Wait for flush
			time.Sleep(150 * time.Millisecond)

			logs := repo.getLogs()
			if len(logs) != 1 {
				return false
			}

			log := logs[0]
			return log.OperatorID == operatorID &&
				log.OperatorName == operatorName &&
				log.TargetType == model.AuditTargetTypeRole &&
				log.TargetID == 200 &&
				log.TargetName == roleName &&
				log.Action == model.AuditActionAssign &&
				log.IPAddress == ipAddress &&
				log.UserAgent == userAgent &&
				log.RequestID == requestID &&
				log.BeforeData != "" &&
				log.AfterData != ""
		},
		gen.UInt64(),
		genNonEmptyString(),
		genNonEmptyString(),
		genIPAddress(),
		genUserAgent(),
		genRequestID(),
	))

	// Property 3: User role change logs should contain all required fields
	properties.Property("user role change logs should contain all required fields", prop.ForAll(
		func(operatorID uint64, operatorName, userName, ipAddress, userAgent, requestID string) bool {
			repo := newPropertyMockRepository()
			svc := NewService(repo, Config{
				BufferSize:    100,
				BatchSize:     10,
				FlushInterval: 50 * time.Millisecond,
			})
			svc.Start()
			defer svc.Stop()

			ctx := context.Background()
			beforeData := []string{"role1"}
			afterData := []string{"role1", "role2"}

			svc.LogUserRoleChange(ctx, operatorID, operatorName, model.AuditActionBatchAssign,
				300, userName, beforeData, afterData, ipAddress, userAgent, requestID)

			// Wait for flush
			time.Sleep(150 * time.Millisecond)

			logs := repo.getLogs()
			if len(logs) != 1 {
				return false
			}

			log := logs[0]
			return log.OperatorID == operatorID &&
				log.OperatorName == operatorName &&
				log.TargetType == model.AuditTargetTypeUser &&
				log.TargetID == 300 &&
				log.TargetName == userName &&
				log.Action == model.AuditActionBatchAssign &&
				log.IPAddress == ipAddress &&
				log.UserAgent == userAgent &&
				log.RequestID == requestID &&
				log.BeforeData != "" &&
				log.AfterData != ""
		},
		gen.UInt64(),
		genNonEmptyString(),
		genNonEmptyString(),
		genIPAddress(),
		genUserAgent(),
		genRequestID(),
	))

	// Property 4: All audit actions should be logged with correct action type
	properties.Property("all audit actions should be logged with correct action type", prop.ForAll(
		func(action model.AuditAction) bool {
			repo := newPropertyMockRepository()
			svc := NewService(repo, Config{
				BufferSize:    100,
				BatchSize:     10,
				FlushInterval: 50 * time.Millisecond,
			})
			svc.Start()
			defer svc.Stop()

			ctx := context.Background()
			svc.LogPermissionChange(ctx, 1, "admin", action, 100, "test", nil, nil, "", "", "")

			// Wait for flush
			time.Sleep(150 * time.Millisecond)

			logs := repo.getLogs()
			if len(logs) != 1 {
				return false
			}

			return logs[0].Action == action
		},
		genAuditAction(),
	))

	// Property 5: Multiple logs should all be persisted
	properties.Property("multiple logs should all be persisted", prop.ForAll(
		func(count int) bool {
			if count < 1 || count > 50 {
				count = 10
			}

			repo := newPropertyMockRepository()
			svc := NewService(repo, Config{
				BufferSize:    100,
				BatchSize:     10,
				FlushInterval: 50 * time.Millisecond,
			})
			svc.Start()
			defer svc.Stop()

			ctx := context.Background()
			for i := 0; i < count; i++ {
				svc.LogPermissionChange(ctx, uint64(i+1), "admin", model.AuditActionCreate,
					uint64(i+100), "perm", nil, nil, "", "", "")
			}

			// Wait for flush
			time.Sleep(300 * time.Millisecond)

			logs := repo.getLogs()
			return len(logs) == count
		},
		gen.IntRange(1, 50),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// genNonEmptyString generates a non-empty string.
func genNonEmptyString() gopter.Gen {
	return gen.AlphaString().SuchThat(func(s string) bool {
		return len(s) > 0
	}).Map(func(s string) string {
		if len(s) == 0 {
			return "default"
		}
		return s
	})
}

// genIPAddress generates a valid IP address string.
func genIPAddress() gopter.Gen {
	ipAddresses := []string{
		"192.168.1.1",
		"10.0.0.1",
		"172.16.0.1",
		"127.0.0.1",
		"8.8.8.8",
	}
	return gen.OneConstOf(ipAddresses[0], ipAddresses[1], ipAddresses[2], ipAddresses[3], ipAddresses[4])
}

// genUserAgent generates a user agent string.
func genUserAgent() gopter.Gen {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
		"Mozilla/5.0 (X11; Linux x86_64)",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_0)",
		"Mozilla/5.0 (Android 11; Mobile)",
	}
	return gen.OneConstOf(userAgents[0], userAgents[1], userAgents[2], userAgents[3], userAgents[4])
}

// genRequestID generates a request ID string.
func genRequestID() gopter.Gen {
	return gen.AlphaString().Map(func(s string) string {
		if len(s) < 8 {
			return "req-" + s + "12345678"[:8-len(s)]
		}
		return "req-" + s[:8]
	})
}

// genAuditAction generates a valid audit action.
func genAuditAction() gopter.Gen {
	return gen.OneConstOf(
		model.AuditActionCreate,
		model.AuditActionUpdate,
		model.AuditActionDelete,
		model.AuditActionAssign,
		model.AuditActionRevoke,
		model.AuditActionBatchAssign,
		model.AuditActionBatchRevoke,
		model.AuditActionInheritChange,
	)
}
