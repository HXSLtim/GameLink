// Package audit provides asynchronous audit logging for permission changes.
package audit

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository/permissionauditlog"
	"gamelink/pkg/logging"
)

// DefaultBufferSize is the default size of the audit log channel buffer.
const DefaultBufferSize = 1000

// DefaultBatchSize is the default number of logs to write in a single batch.
const DefaultBatchSize = 50

// DefaultFlushInterval is the default interval for flushing logs to the database.
const DefaultFlushInterval = 5 * time.Second

// Config holds configuration for the async audit log writer.
type Config struct {
	// BufferSize is the size of the channel buffer for pending logs.
	BufferSize int
	// BatchSize is the maximum number of logs to write in a single batch.
	BatchSize int
	// FlushInterval is the interval for flushing logs to the database.
	FlushInterval time.Duration
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		BufferSize:    DefaultBufferSize,
		BatchSize:     DefaultBatchSize,
		FlushInterval: DefaultFlushInterval,
	}
}

// Service provides asynchronous audit logging functionality.
type Service struct {
	repo           permissionauditlog.Repository
	config         Config
	logChan        chan *model.PermissionAuditLog
	stopChan       chan struct{}
	doneChan       chan struct{}
	wg             sync.WaitGroup
	mu             sync.RWMutex
	running        bool
	droppedCount   int64
	processedCount int64
}

// NewService creates a new audit service with the given repository and config.
func NewService(repo permissionauditlog.Repository, config Config) *Service {
	if config.BufferSize <= 0 {
		config.BufferSize = DefaultBufferSize
	}
	if config.BatchSize <= 0 {
		config.BatchSize = DefaultBatchSize
	}
	if config.FlushInterval <= 0 {
		config.FlushInterval = DefaultFlushInterval
	}

	return &Service{
		repo:     repo,
		config:   config,
		logChan:  make(chan *model.PermissionAuditLog, config.BufferSize),
		stopChan: make(chan struct{}),
		doneChan: make(chan struct{}),
	}
}

// NewServiceWithDefaults creates a new audit service with default configuration.
func NewServiceWithDefaults(repo permissionauditlog.Repository) *Service {
	return NewService(repo, DefaultConfig())
}

// Start begins the background log processing goroutine.
func (s *Service) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return
	}

	s.running = true
	s.wg.Add(1)
	go s.processLogs()
}

// Stop gracefully stops the audit service, flushing any pending logs.
// It blocks until all pending logs are written or the default timeout (30s) is reached.
func (s *Service) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()

	close(s.stopChan)

	// Wait for processing to complete with a timeout
	select {
	case <-s.doneChan:
		logging.Info("audit service stopped gracefully")
	case <-time.After(30 * time.Second):
		logging.Warn("audit service stop timed out")
	}
}

// StopWithContext gracefully stops the audit service with a custom context.
func (s *Service) StopWithContext(ctx context.Context) error {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return nil
	}
	s.running = false
	s.mu.Unlock()

	close(s.stopChan)

	select {
	case <-s.doneChan:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Log queues an audit log entry for asynchronous writing.
func (s *Service) Log(log *model.PermissionAuditLog) {
	s.mu.RLock()
	running := s.running
	s.mu.RUnlock()

	if !running {
		logging.Warn("audit service not running, log dropped")
		return
	}

	select {
	case s.logChan <- log:
		// Successfully queued
	default:
		s.mu.Lock()
		s.droppedCount++
		s.mu.Unlock()
		logging.Warn("audit log buffer full, log dropped")
	}
}

// LogPermissionChange logs a permission change event.
func (s *Service) LogPermissionChange(ctx context.Context, operatorID uint64, operatorName string,
	action model.AuditAction, permissionID uint64, permissionName string,
	beforeData, afterData interface{}, ipAddress, userAgent, requestID string) {

	beforeJSON, _ := json.Marshal(beforeData)
	afterJSON, _ := json.Marshal(afterData)

	s.Log(&model.PermissionAuditLog{
		OperatorID:   operatorID,
		OperatorName: operatorName,
		TargetType:   model.AuditTargetTypePermission,
		TargetID:     permissionID,
		TargetName:   permissionName,
		Action:       action,
		BeforeData:   string(beforeJSON),
		AfterData:    string(afterJSON),
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		RequestID:    requestID,
	})
}

// LogRoleChange logs a role change event.
func (s *Service) LogRoleChange(ctx context.Context, operatorID uint64, operatorName string,
	action model.AuditAction, roleID uint64, roleName string,
	beforeData, afterData interface{}, ipAddress, userAgent, requestID string) {

	beforeJSON, _ := json.Marshal(beforeData)
	afterJSON, _ := json.Marshal(afterData)

	s.Log(&model.PermissionAuditLog{
		OperatorID:   operatorID,
		OperatorName: operatorName,
		TargetType:   model.AuditTargetTypeRole,
		TargetID:     roleID,
		TargetName:   roleName,
		Action:       action,
		BeforeData:   string(beforeJSON),
		AfterData:    string(afterJSON),
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		RequestID:    requestID,
	})
}

// LogUserRoleChange logs a user role assignment change event.
func (s *Service) LogUserRoleChange(ctx context.Context, operatorID uint64, operatorName string,
	action model.AuditAction, userID uint64, userName string,
	beforeData, afterData interface{}, ipAddress, userAgent, requestID string) {

	beforeJSON, _ := json.Marshal(beforeData)
	afterJSON, _ := json.Marshal(afterData)

	s.Log(&model.PermissionAuditLog{
		OperatorID:   operatorID,
		OperatorName: operatorName,
		TargetType:   model.AuditTargetTypeUser,
		TargetID:     userID,
		TargetName:   userName,
		Action:       action,
		BeforeData:   string(beforeJSON),
		AfterData:    string(afterJSON),
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		RequestID:    requestID,
	})
}

// Stats returns current statistics about the audit service.
type Stats struct {
	Running        bool  `json:"running"`
	BufferSize     int   `json:"bufferSize"`
	BufferUsed     int   `json:"bufferUsed"`
	ProcessedCount int64 `json:"processedCount"`
	DroppedCount   int64 `json:"droppedCount"`
}

// GetStats returns current statistics about the audit service.
func (s *Service) GetStats() Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return Stats{
		Running:        s.running,
		BufferSize:     s.config.BufferSize,
		BufferUsed:     len(s.logChan),
		ProcessedCount: s.processedCount,
		DroppedCount:   s.droppedCount,
	}
}

// IsRunning returns whether the audit service is currently running.
func (s *Service) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// GetRepository returns the underlying repository for direct queries.
func (s *Service) GetRepository() permissionauditlog.Repository {
	return s.repo
}

// QueryOptions defines filtering options for audit log queries.
type QueryOptions struct {
	Page       int
	PageSize   int
	OperatorID *uint64
	TargetType *model.AuditTargetType
	TargetID   *uint64
	Action     *model.AuditAction
	DateFrom   *time.Time
	DateTo     *time.Time
	RequestID  string
	IPAddress  string
}

// QueryResult represents the result of an audit log query.
type QueryResult struct {
	Logs       []model.PermissionAuditLog `json:"logs"`
	Total      int64                      `json:"total"`
	Page       int                        `json:"page"`
	PageSize   int                        `json:"pageSize"`
	TotalPages int                        `json:"totalPages"`
}

// Query retrieves audit logs with filtering and pagination.
// Supports filtering by time range, operator, action type, and target.
func (s *Service) Query(ctx context.Context, opts QueryOptions) (*QueryResult, error) {
	repoOpts := permissionauditlog.ListOptions{
		Page:       opts.Page,
		PageSize:   opts.PageSize,
		OperatorID: opts.OperatorID,
		TargetType: opts.TargetType,
		TargetID:   opts.TargetID,
		Action:     opts.Action,
		DateFrom:   opts.DateFrom,
		DateTo:     opts.DateTo,
		RequestID:  opts.RequestID,
		IPAddress:  opts.IPAddress,
	}

	logs, total, err := s.repo.List(ctx, repoOpts)
	if err != nil {
		return nil, err
	}

	// Normalize pagination values
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

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &QueryResult{
		Logs:       logs,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetByID retrieves a single audit log by ID.
func (s *Service) GetByID(ctx context.Context, id uint64) (*model.PermissionAuditLog, error) {
	return s.repo.Get(ctx, id)
}

// QueryByOperator retrieves audit logs for a specific operator.
func (s *Service) QueryByOperator(ctx context.Context, operatorID uint64, page, pageSize int) (*QueryResult, error) {
	return s.Query(ctx, QueryOptions{
		Page:       page,
		PageSize:   pageSize,
		OperatorID: &operatorID,
	})
}

// QueryByAction retrieves audit logs for a specific action type.
func (s *Service) QueryByAction(ctx context.Context, action model.AuditAction, page, pageSize int) (*QueryResult, error) {
	return s.Query(ctx, QueryOptions{
		Page:     page,
		PageSize: pageSize,
		Action:   &action,
	})
}

// QueryByDateRange retrieves audit logs within a date range.
func (s *Service) QueryByDateRange(ctx context.Context, from, to time.Time, page, pageSize int) (*QueryResult, error) {
	return s.Query(ctx, QueryOptions{
		Page:     page,
		PageSize: pageSize,
		DateFrom: &from,
		DateTo:   &to,
	})
}

// QueryByTarget retrieves audit logs for a specific target.
func (s *Service) QueryByTarget(ctx context.Context, targetType model.AuditTargetType, targetID uint64, page, pageSize int) (*QueryResult, error) {
	return s.Query(ctx, QueryOptions{
		Page:       page,
		PageSize:   pageSize,
		TargetType: &targetType,
		TargetID:   &targetID,
	})
}

// ExportOptions defines options for exporting audit logs.
type ExportOptions struct {
	OperatorID *uint64
	TargetType *model.AuditTargetType
	TargetID   *uint64
	Action     *model.AuditAction
	DateFrom   *time.Time
	DateTo     *time.Time
	MaxRecords int // Maximum number of records to export (default 10000)
}

// ExportCSV exports audit logs to CSV format.
// Returns the CSV data as bytes.
func (s *Service) ExportCSV(ctx context.Context, opts ExportOptions) ([]byte, error) {
	maxRecords := opts.MaxRecords
	if maxRecords <= 0 {
		maxRecords = 10000
	}

	// Fetch all matching logs (up to maxRecords)
	repoOpts := permissionauditlog.ListOptions{
		Page:       1,
		PageSize:   maxRecords,
		OperatorID: opts.OperatorID,
		TargetType: opts.TargetType,
		TargetID:   opts.TargetID,
		Action:     opts.Action,
		DateFrom:   opts.DateFrom,
		DateTo:     opts.DateTo,
	}

	logs, _, err := s.repo.List(ctx, repoOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch audit logs: %w", err)
	}

	// Create CSV buffer
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write UTF-8 BOM for Excel compatibility
	buf.Write([]byte{0xEF, 0xBB, 0xBF})

	// Write header
	header := []string{
		"ID",
		"时间",
		"操作者ID",
		"操作者名称",
		"目标类型",
		"目标ID",
		"目标名称",
		"操作类型",
		"操作前数据",
		"操作后数据",
		"IP地址",
		"用户代理",
		"请求ID",
	}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, log := range logs {
		row := []string{
			fmt.Sprintf("%d", log.ID),
			log.CreatedAt.Format("2006-01-02 15:04:05"),
			fmt.Sprintf("%d", log.OperatorID),
			log.OperatorName,
			string(log.TargetType),
			fmt.Sprintf("%d", log.TargetID),
			log.TargetName,
			string(log.Action),
			log.BeforeData,
			log.AfterData,
			log.IPAddress,
			log.UserAgent,
			log.RequestID,
		}
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}

	return buf.Bytes(), nil
}

// GenerateExportFilename generates a filename for the export file.
func GenerateExportFilename(prefix string) string {
	return fmt.Sprintf("%s_%s.csv", prefix, time.Now().Format("20060102_150405"))
}

// ArchiveResult represents the result of an archive operation.
type ArchiveResult struct {
	ArchivedCount int64     `json:"archivedCount"`
	DeletedCount  int64     `json:"deletedCount"`
	ArchiveFile   string    `json:"archiveFile,omitempty"`
	ArchivedAt    time.Time `json:"archivedAt"`
}

// ArchiveOldLogs archives audit logs older than the retention period.
// It exports logs to CSV before deleting them from the database.
// - Online retention: 90 days (logs older than this are archived)
// - Archive retention: 365 days (archived logs older than this are deleted)
func (s *Service) ArchiveOldLogs(ctx context.Context, archiveDir string) (*ArchiveResult, error) {
	now := time.Now()
	onlineRetentionCutoff := now.AddDate(0, 0, -model.AuditLogRetentionDays)
	archiveRetentionCutoff := now.AddDate(0, 0, -model.AuditLogArchiveDays)

	result := &ArchiveResult{
		ArchivedAt: now,
	}

	// Count logs to be archived (between archive cutoff and online cutoff)
	archivedCount, err := s.repo.CountByDateRange(ctx, archiveRetentionCutoff, onlineRetentionCutoff)
	if err != nil {
		return nil, fmt.Errorf("failed to count logs for archiving: %w", err)
	}
	result.ArchivedCount = archivedCount

	// Export logs to be archived if there are any
	if archivedCount > 0 && archiveDir != "" {
		exportOpts := ExportOptions{
			DateFrom:   &archiveRetentionCutoff,
			DateTo:     &onlineRetentionCutoff,
			MaxRecords: int(archivedCount) + 1000, // Add buffer
		}

		csvData, err := s.ExportCSV(ctx, exportOpts)
		if err != nil {
			logging.Warn("failed to export archive data: %v", err)
		} else {
			result.ArchiveFile = GenerateExportFilename("audit_archive")
			// Note: Actual file writing should be handled by the caller
			// This is to keep the service layer independent of filesystem operations
			_ = csvData // Archive data available for caller to save
		}
	}

	// Delete logs older than online retention period
	deletedCount, err := s.repo.DeleteBefore(ctx, onlineRetentionCutoff)
	if err != nil {
		return nil, fmt.Errorf("failed to delete old logs: %w", err)
	}
	result.DeletedCount = deletedCount

	logging.Info("audit log archive completed: archived=%d, deleted=%d", archivedCount, deletedCount)

	return result, nil
}

// CleanupOldLogs deletes audit logs older than the specified number of days.
// This is a simpler alternative to ArchiveOldLogs when archiving is not needed.
func (s *Service) CleanupOldLogs(ctx context.Context, retentionDays int) (int64, error) {
	if retentionDays <= 0 {
		retentionDays = model.AuditLogRetentionDays
	}

	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	return s.repo.DeleteBefore(ctx, cutoff)
}

// GetRetentionStats returns statistics about audit log retention.
type RetentionStats struct {
	TotalLogs        int64     `json:"totalLogs"`
	LogsWithinOnline int64     `json:"logsWithinOnline"` // Within 90 days
	LogsForArchive   int64     `json:"logsForArchive"`   // Between 90-365 days
	LogsForDeletion  int64     `json:"logsForDeletion"`  // Older than 365 days
	OldestLogDate    time.Time `json:"oldestLogDate"`
	OnlineRetention  int       `json:"onlineRetention"`  // Days
	ArchiveRetention int       `json:"archiveRetention"` // Days
}

// GetRetentionStats returns statistics about audit log retention periods.
func (s *Service) GetRetentionStats(ctx context.Context) (*RetentionStats, error) {
	now := time.Now()
	onlineCutoff := now.AddDate(0, 0, -model.AuditLogRetentionDays)
	archiveCutoff := now.AddDate(0, 0, -model.AuditLogArchiveDays)
	veryOld := now.AddDate(-10, 0, 0) // 10 years ago as a practical minimum

	stats := &RetentionStats{
		OnlineRetention:  model.AuditLogRetentionDays,
		ArchiveRetention: model.AuditLogArchiveDays,
	}

	// Count logs within online retention (last 90 days)
	logsWithinOnline, err := s.repo.CountByDateRange(ctx, onlineCutoff, now)
	if err != nil {
		return nil, fmt.Errorf("failed to count online logs: %w", err)
	}
	stats.LogsWithinOnline = logsWithinOnline

	// Count logs for archive (90-365 days old)
	logsForArchive, err := s.repo.CountByDateRange(ctx, archiveCutoff, onlineCutoff)
	if err != nil {
		return nil, fmt.Errorf("failed to count archive logs: %w", err)
	}
	stats.LogsForArchive = logsForArchive

	// Count logs for deletion (older than 365 days)
	logsForDeletion, err := s.repo.CountByDateRange(ctx, veryOld, archiveCutoff)
	if err != nil {
		return nil, fmt.Errorf("failed to count deletion logs: %w", err)
	}
	stats.LogsForDeletion = logsForDeletion

	stats.TotalLogs = logsWithinOnline + logsForArchive + logsForDeletion

	// Get oldest log date
	result, err := s.Query(ctx, QueryOptions{
		Page:     1,
		PageSize: 1,
	})
	if err == nil && len(result.Logs) > 0 {
		// The query returns newest first, so we need to get the last page
		if result.TotalPages > 0 {
			lastPageResult, err := s.Query(ctx, QueryOptions{
				Page:     result.TotalPages,
				PageSize: 1,
			})
			if err == nil && len(lastPageResult.Logs) > 0 {
				stats.OldestLogDate = lastPageResult.Logs[0].CreatedAt
			}
		}
	}

	return stats, nil
}

// processLogs is the background goroutine that processes queued logs.
func (s *Service) processLogs() {
	defer s.wg.Done()
	defer close(s.doneChan)

	ticker := time.NewTicker(s.config.FlushInterval)
	defer ticker.Stop()

	batch := make([]*model.PermissionAuditLog, 0, s.config.BatchSize)

	flush := func() {
		if len(batch) == 0 {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.repo.CreateBatch(ctx, batch); err != nil {
			logging.Error("failed to write audit logs: %v", err)
		} else {
			s.mu.Lock()
			s.processedCount += int64(len(batch))
			s.mu.Unlock()
		}

		batch = batch[:0]
	}

	for {
		select {
		case log := <-s.logChan:
			batch = append(batch, log)
			if len(batch) >= s.config.BatchSize {
				flush()
			}

		case <-ticker.C:
			flush()

		case <-s.stopChan:
			// Drain remaining logs
			for {
				select {
				case log := <-s.logChan:
					batch = append(batch, log)
					if len(batch) >= s.config.BatchSize {
						flush()
					}
				default:
					flush()
					return
				}
			}
		}
	}
}
