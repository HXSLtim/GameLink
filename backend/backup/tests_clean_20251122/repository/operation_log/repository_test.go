package operationlog

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

func testContext() context.Context {
	return context.Background()
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	if err := db.AutoMigrate(&model.OperationLog{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestOperationLogRepository_Append(t *testing.T) {
	db := setupTestDB(t)
	repo := NewOperationLogRepository(db)

	actorID := uint64(1)
	metadata := json.RawMessage(`{"key":"value"}`)

	log := &model.OperationLog{
		EntityType:   "order",
		EntityID:     100,
		ActorUserID:  &actorID,
		Action:       "create",
		Reason:       "新建订单",
		TraceID:      "trace-append-001",
		MetadataJSON: metadata,
	}

	err := repo.Append(testContext(), log)
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	if log.ID == 0 {
		t.Error("expected log ID to be set")
	}

	// Verify created
	var retrieved model.OperationLog
	db.First(&retrieved, log.ID)

	if retrieved.EntityType != "order" {
		t.Errorf("expected entity type 'order', got %s", retrieved.EntityType)
	}
	if retrieved.EntityID != 100 {
		t.Errorf("expected entity ID 100, got %d", retrieved.EntityID)
	}
	if *retrieved.ActorUserID != 1 {
		t.Errorf("expected actor user ID 1, got %d", *retrieved.ActorUserID)
	}
	if retrieved.Action != "create" {
		t.Errorf("expected action 'create', got %s", retrieved.Action)
	}
	if retrieved.TraceID != "trace-append-001" {
		t.Errorf("expected trace id 'trace-append-001', got %s", retrieved.TraceID)
	}
}

func TestOperationLogRepository_AppendWithoutActor(t *testing.T) {
	db := setupTestDB(t)
	repo := NewOperationLogRepository(db)

	log := &model.OperationLog{
		EntityType: "payment",
		EntityID:   200,
		Action:     "refund",
		Reason:     "系统自动退款",
	}

	err := repo.Append(testContext(), log)
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	// Verify ActorUserID is nil
	var retrieved model.OperationLog
	db.First(&retrieved, log.ID)

	if retrieved.ActorUserID != nil {
		t.Error("expected ActorUserID to be nil")
	}
}

func TestOperationLogRepository_ListByEntity(t *testing.T) {
	db := setupTestDB(t)
	repo := NewOperationLogRepository(db)

	// Create test data
	actorID1 := uint64(1)
	actorID2 := uint64(2)

	logs := []*model.OperationLog{
		{EntityType: "order", EntityID: 1, ActorUserID: &actorID1, Action: "create"},
		{EntityType: "order", EntityID: 1, ActorUserID: &actorID1, Action: "update_status"},
		{EntityType: "order", EntityID: 1, ActorUserID: &actorID2, Action: "complete"},
		{EntityType: "order", EntityID: 2, ActorUserID: &actorID1, Action: "create"},
		{EntityType: "payment", EntityID: 1, ActorUserID: &actorID1, Action: "capture"},
	}

	for _, log := range logs {
		_ = repo.Append(testContext(), log)
	}

	t.Run("List all logs for order 1", func(t *testing.T) {
		result, total, err := repo.ListByEntity(testContext(), "order", 1, repository.OperationLogListOptions{
			Page:     1,
			PageSize: 10,
		})
		if err != nil {
			t.Fatalf("ListByEntity failed: %v", err)
		}

		if total != 3 {
			t.Errorf("expected total 3, got %d", total)
		}
		if len(result) != 3 {
			t.Errorf("expected 3 logs, got %d", len(result))
		}
	})

	t.Run("List logs for order 2", func(t *testing.T) {
		result, total, err := repo.ListByEntity(testContext(), "order", 2, repository.OperationLogListOptions{
			Page:     1,
			PageSize: 10,
		})
		if err != nil {
			t.Fatalf("ListByEntity failed: %v", err)
		}

		if total != 1 {
			t.Errorf("expected total 1, got %d", total)
		}
		if len(result) != 1 {
			t.Errorf("expected 1 log, got %d", len(result))
		}
	})

	t.Run("List logs for payment 1", func(t *testing.T) {
		result, total, err := repo.ListByEntity(testContext(), "payment", 1, repository.OperationLogListOptions{
			Page:     1,
			PageSize: 10,
		})
		if err != nil {
			t.Fatalf("ListByEntity failed: %v", err)
		}

		if total != 1 {
			t.Errorf("expected total 1, got %d", total)
		}
		if result[0].Action != "capture" {
			t.Errorf("expected action 'capture', got %s", result[0].Action)
		}
	})

	t.Run("List non-existent entity", func(t *testing.T) {
		result, total, err := repo.ListByEntity(testContext(), "order", 999, repository.OperationLogListOptions{
			Page:     1,
			PageSize: 10,
		})
		if err != nil {
			t.Fatalf("ListByEntity failed: %v", err)
		}

		if total != 0 {
			t.Errorf("expected total 0, got %d", total)
		}
		if len(result) != 0 {
			t.Errorf("expected 0 logs, got %d", len(result))
		}
	})
}

func TestOperationLogRepository_ListByEntityWithActionFilter(t *testing.T) {
	db := setupTestDB(t)
	repo := NewOperationLogRepository(db)

	actorID := uint64(1)

	logs := []*model.OperationLog{
		{EntityType: "order", EntityID: 1, ActorUserID: &actorID, Action: "create"},
		{EntityType: "order", EntityID: 1, ActorUserID: &actorID, Action: "update_status"},
		{EntityType: "order", EntityID: 1, ActorUserID: &actorID, Action: "update_status"},
		{EntityType: "order", EntityID: 1, ActorUserID: &actorID, Action: "complete"},
	}

	for _, log := range logs {
		_ = repo.Append(testContext(), log)
	}

	t.Run("Filter by action 'update_status'", func(t *testing.T) {
		result, total, err := repo.ListByEntity(testContext(), "order", 1, repository.OperationLogListOptions{
			Page:     1,
			PageSize: 10,
			Action:   "update_status",
		})
		if err != nil {
			t.Fatalf("ListByEntity failed: %v", err)
		}

		if total != 2 {
			t.Errorf("expected total 2, got %d", total)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 logs, got %d", len(result))
		}

		for _, log := range result {
			if log.Action != "update_status" {
				t.Errorf("expected action 'update_status', got %s", log.Action)
			}
		}
	})

	t.Run("Filter by action 'create'", func(t *testing.T) {
		_, total, err := repo.ListByEntity(testContext(), "order", 1, repository.OperationLogListOptions{
			Page:     1,
			PageSize: 10,
			Action:   "create",
		})
		if err != nil {
			t.Fatalf("ListByEntity failed: %v", err)
		}

		if total != 1 {
			t.Errorf("expected total 1, got %d", total)
		}
	})
}

func TestOperationLogRepository_ListByEntityWithActorFilter(t *testing.T) {
	db := setupTestDB(t)
	repo := NewOperationLogRepository(db)

	actorID1 := uint64(1)
	actorID2 := uint64(2)

	logs := []*model.OperationLog{
		{EntityType: "order", EntityID: 1, ActorUserID: &actorID1, Action: "create"},
		{EntityType: "order", EntityID: 1, ActorUserID: &actorID1, Action: "update_status"},
		{EntityType: "order", EntityID: 1, ActorUserID: &actorID2, Action: "complete"},
	}

	for _, log := range logs {
		_ = repo.Append(testContext(), log)
	}

	t.Run("Filter by actor user ID 1", func(t *testing.T) {
		actorFilter := uint64(1)
		result, total, err := repo.ListByEntity(testContext(), "order", 1, repository.OperationLogListOptions{
			Page:        1,
			PageSize:    10,
			ActorUserID: &actorFilter,
		})
		if err != nil {
			t.Fatalf("ListByEntity failed: %v", err)
		}

		if total != 2 {
			t.Errorf("expected total 2, got %d", total)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 logs, got %d", len(result))
		}

		for _, log := range result {
			if *log.ActorUserID != 1 {
				t.Errorf("expected actor user ID 1, got %d", *log.ActorUserID)
			}
		}
	})

	t.Run("Filter by actor user ID 2", func(t *testing.T) {
		actorFilter := uint64(2)
		logs, total, err := repo.ListByEntity(testContext(), "order", 1, repository.OperationLogListOptions{
			Page:        1,
			PageSize:    10,
			ActorUserID: &actorFilter,
		})
		if err != nil {
			t.Fatalf("ListByEntity failed: %v", err)
		}

		if total != 1 {
			t.Errorf("expected total 1, got %d", total)
		}
		if logs[0].Action != "complete" {
			t.Errorf("expected action 'complete', got %s", logs[0].Action)
		}
	})
}

func TestOperationLogRepository_ListByEntityWithDateFilter(t *testing.T) {
	db := setupTestDB(t)
	repo := NewOperationLogRepository(db)

	actorID := uint64(1)

	// Create logs with different timestamps
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	lastWeek := now.Add(-7 * 24 * time.Hour)

	logs := []*model.OperationLog{
		{EntityType: "order", EntityID: 1, ActorUserID: &actorID, Action: "create"},
		{EntityType: "order", EntityID: 1, ActorUserID: &actorID, Action: "update_status"},
		{EntityType: "order", EntityID: 1, ActorUserID: &actorID, Action: "complete"},
	}

	for i, log := range logs {
		_ = repo.Append(testContext(), log)
		// Manually update created_at for testing (in real scenario, timestamps would differ)
		var timestamp time.Time
		switch i {
		case 0:
			timestamp = lastWeek
		case 1:
			timestamp = yesterday
		case 2:
			timestamp = now
		}
		db.Model(&model.OperationLog{}).Where("id = ?", log.ID).Update("created_at", timestamp)
	}

	t.Run("Filter from yesterday", func(t *testing.T) {
		dateFrom := yesterday.Add(-1 * time.Hour)
		_, total, err := repo.ListByEntity(testContext(), "order", 1, repository.OperationLogListOptions{
			Page:     1,
			PageSize: 10,
			DateFrom: &dateFrom,
		})
		if err != nil {
			t.Fatalf("ListByEntity failed: %v", err)
		}

		if total < 2 {
			t.Errorf("expected at least 2 logs from yesterday, got %d", total)
		}
	})

	t.Run("Filter to yesterday", func(t *testing.T) {
		dateTo := yesterday.Add(1 * time.Hour)
		result, total, err := repo.ListByEntity(testContext(), "order", 1, repository.OperationLogListOptions{
			Page:     1,
			PageSize: 10,
			DateTo:   &dateTo,
		})
		if err != nil {
			t.Fatalf("ListByEntity failed: %v", err)
		}

		if total < 2 {
			t.Errorf("expected at least 2 logs up to yesterday, got %d", total)
		}
		if len(result) < 2 {
			t.Errorf("expected at least 2 logs, got %d", len(result))
		}
	})

	t.Run("Filter date range", func(t *testing.T) {
		dateFrom := lastWeek.Add(-1 * time.Hour)
		dateTo := yesterday.Add(1 * time.Hour)
		_, total, err := repo.ListByEntity(testContext(), "order", 1, repository.OperationLogListOptions{
			Page:     1,
			PageSize: 10,
			DateFrom: &dateFrom,
			DateTo:   &dateTo,
		})
		if err != nil {
			t.Fatalf("ListByEntity failed: %v", err)
		}

		if total < 2 {
			t.Errorf("expected at least 2 logs in date range, got %d", total)
		}
	})
}

func TestOperationLogRepository_ListByEntityPagination(t *testing.T) {
	db := setupTestDB(t)
	repo := NewOperationLogRepository(db)

	actorID := uint64(1)

	// Create 15 logs for the same order
	for i := 1; i <= 15; i++ {
		log := &model.OperationLog{
			EntityType:  "order",
			EntityID:    1,
			ActorUserID: &actorID,
			Action:      "update_status",
		}
		_ = repo.Append(testContext(), log)
	}

	t.Run("First page", func(t *testing.T) {
		result, total, err := repo.ListByEntity(testContext(), "order", 1, repository.OperationLogListOptions{
			Page:     1,
			PageSize: 10,
		})
		if err != nil {
			t.Fatalf("ListByEntity failed: %v", err)
		}

		if total != 15 {
			t.Errorf("expected total 15, got %d", total)
		}
		if len(result) != 10 {
			t.Errorf("expected 10 logs, got %d", len(result))
		}
	})

	t.Run("Second page", func(t *testing.T) {
		logs, total, err := repo.ListByEntity(testContext(), "order", 1, repository.OperationLogListOptions{
			Page:     2,
			PageSize: 10,
		})
		if err != nil {
			t.Fatalf("ListByEntity failed: %v", err)
		}

		if total != 15 {
			t.Errorf("expected total 15, got %d", total)
		}
		if len(logs) != 5 {
			t.Errorf("expected 5 logs on second page, got %d", len(logs))
		}
	})

	t.Run("Custom page size", func(t *testing.T) {
		logs, total, err := repo.ListByEntity(testContext(), "order", 1, repository.OperationLogListOptions{
			Page:     1,
			PageSize: 5,
		})
		if err != nil {
			t.Fatalf("ListByEntity failed: %v", err)
		}

		if total != 15 {
			t.Errorf("expected total 15, got %d", total)
		}
		if len(logs) != 5 {
			t.Errorf("expected 5 logs, got %d", len(logs))
		}
	})
}

func TestOperationLogRepository_CompleteWorkflow(t *testing.T) {
	db := setupTestDB(t)
	repo := NewOperationLogRepository(db)

	actorID := uint64(1)
	metadata := json.RawMessage(`{"oldStatus":"pending","newStatus":"completed"}`)

	// Append log
	log := &model.OperationLog{
		EntityType:   "order",
		EntityID:     100,
		ActorUserID:  &actorID,
		Action:       "complete",
		Reason:       "订单完成",
		MetadataJSON: metadata,
	}

	err := repo.Append(testContext(), log)
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	// List logs
	result, total, err := repo.ListByEntity(testContext(), "order", 100, repository.OperationLogListOptions{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("ListByEntity failed: %v", err)
	}

	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 log, got %d", len(result))
	}

	retrieved := result[0]
	if retrieved.EntityType != "order" {
		t.Errorf("expected entity type 'order', got %s", retrieved.EntityType)
	}
	if retrieved.Action != "complete" {
		t.Errorf("expected action 'complete', got %s", retrieved.Action)
	}
	if retrieved.Reason != "订单完成" {
		t.Errorf("expected reason '订单完成', got %s", retrieved.Reason)
	}
	if string(retrieved.MetadataJSON) != string(metadata) {
		t.Errorf("expected metadata %s, got %s", metadata, retrieved.MetadataJSON)
	}
}
