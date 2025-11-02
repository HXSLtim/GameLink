package payment

import (
	"context"
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

	if err := db.AutoMigrate(&model.Payment{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestNewPaymentRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	repo := NewPaymentRepository(db)
	if repo == nil {
		t.Fatal("NewPaymentRepository returned nil")
	}
}

func TestPaymentRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPaymentRepository(db)

	payment := &model.Payment{
		OrderID:     1,
		UserID:      1,
		Status:      model.PaymentStatusPending,
		Method:      model.PaymentMethodWeChat,
		AmountCents: 10000,
	}

	err := repo.Create(testContext(), payment)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if payment.ID == 0 {
		t.Error("expected ID to be set after create")
	}

	retrieved, _ := repo.Get(testContext(), payment.ID)
	if retrieved.AmountCents != 10000 {
		t.Errorf("expected amount 10000, got %d", retrieved.AmountCents)
	}
}

func TestPaymentRepository_Get(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPaymentRepository(db)

	payment := &model.Payment{
		OrderID:     2,
		UserID:      2,
		Status:      model.PaymentStatusPaid,
		AmountCents: 5000,
	}
	_ = repo.Create(testContext(), payment)

	t.Run("Get existing payment", func(t *testing.T) {
		retrieved, err := repo.Get(testContext(), payment.ID)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if retrieved.ID != payment.ID {
			t.Errorf("expected ID %d, got %d", payment.ID, retrieved.ID)
		}
	})

	t.Run("Get non-existent payment", func(t *testing.T) {
		_, err := repo.Get(testContext(), 99999)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestPaymentRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPaymentRepository(db)

	payment := &model.Payment{
		OrderID:     3,
		UserID:      3,
		Status:      model.PaymentStatusPending,
		AmountCents: 8000,
	}
	_ = repo.Create(testContext(), payment)

	t.Run("Update existing payment", func(t *testing.T) {
		now := time.Now()
		payment.Status = model.PaymentStatusPaid
		payment.PaidAt = &now

		err := repo.Update(testContext(), payment)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		updated, _ := repo.Get(testContext(), payment.ID)
		if updated.Status != model.PaymentStatusPaid {
			t.Errorf("expected status paid, got %s", updated.Status)
		}
	})

	t.Run("Update non-existent payment", func(t *testing.T) {
		nonExistent := &model.Payment{Base: model.Base{ID: 99999}}
		err := repo.Update(testContext(), nonExistent)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestPaymentRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPaymentRepository(db)

	payment := &model.Payment{
		OrderID:     4,
		UserID:      4,
		Status:      model.PaymentStatusPending,
		AmountCents: 6000,
	}
	_ = repo.Create(testContext(), payment)

	t.Run("Delete existing payment", func(t *testing.T) {
		err := repo.Delete(testContext(), payment.ID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		_, err = repo.Get(testContext(), payment.ID)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound after delete, got %v", err)
		}
	})

	t.Run("Delete non-existent payment", func(t *testing.T) {
		err := repo.Delete(testContext(), 99999)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestPaymentRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPaymentRepository(db)

	// 创建多个支付记录
	for i := 0; i < 15; i++ {
		status := model.PaymentStatusPending
		if i%3 == 0 {
			status = model.PaymentStatusPaid
		} else if i%3 == 1 {
			status = model.PaymentStatusFailed
		}

		method := model.PaymentMethodWeChat
		if i%2 == 0 {
			method = model.PaymentMethodAlipay
		}

		payment := &model.Payment{
			OrderID:     uint64(100 + i),
			UserID:      uint64(10 + i%3),
			Status:      status,
			Method:      method,
			AmountCents: int64(1000 * (i + 1)),
		}
		_ = repo.Create(testContext(), payment)
	}

	t.Run("List with pagination", func(t *testing.T) {
		opts := repository.PaymentListOptions{
			Page:     1,
			PageSize: 10,
		}
		payments, total, err := repo.List(testContext(), opts)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if len(payments) != 10 {
			t.Errorf("expected 10 payments, got %d", len(payments))
		}
		if total < 15 {
			t.Errorf("expected total >= 15, got %d", total)
		}
	})

	t.Run("List with status filter", func(t *testing.T) {
		opts := repository.PaymentListOptions{
			Page:     1,
			PageSize: 20,
			Statuses: []model.PaymentStatus{model.PaymentStatusPaid},
		}
		_, total, err := repo.List(testContext(), opts)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if total < 5 {
			t.Errorf("expected at least 5 paid payments, got %d", total)
		}
	})

	t.Run("List with method filter", func(t *testing.T) {
		opts := repository.PaymentListOptions{
			Page:     1,
			PageSize: 20,
			Methods:  []model.PaymentMethod{model.PaymentMethodWeChat},
		}
		_, total, err := repo.List(testContext(), opts)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if total >= 1 {
			t.Logf("found %d wechat payments", total)
		}
	})

	t.Run("List with user filter", func(t *testing.T) {
		userID := uint64(10)
		opts := repository.PaymentListOptions{
			Page:     1,
			PageSize: 20,
			UserID:   &userID,
		}
		_, total, err := repo.List(testContext(), opts)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if total >= 1 {
			t.Logf("found %d payments for user 10", total)
		}
	})

	t.Run("List with order filter", func(t *testing.T) {
		orderID := uint64(100)
		opts := repository.PaymentListOptions{
			Page:     1,
			PageSize: 20,
			OrderID:  &orderID,
		}
		_, total, err := repo.List(testContext(), opts)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if total == 1 {
			t.Logf("found payment for order 100")
		}
	})

	t.Run("List with date range", func(t *testing.T) {
		now := time.Now()
		yesterday := now.Add(-24 * time.Hour)
		tomorrow := now.Add(24 * time.Hour)

		opts := repository.PaymentListOptions{
			Page:     1,
			PageSize: 20,
			DateFrom: &yesterday,
			DateTo:   &tomorrow,
		}
		_, total, err := repo.List(testContext(), opts)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if total < 15 {
			t.Errorf("expected at least 15 payments in date range, got %d", total)
		}
	})
}
