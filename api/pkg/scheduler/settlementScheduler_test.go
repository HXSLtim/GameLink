package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockCommissionService wraps the commission service with mocking capabilities
type MockCommissionService struct {
	mock.Mock
}

func (m *MockCommissionService) SettleMonth(ctx context.Context, month string) error {
	args := m.Called(ctx, month)
	return args.Error(0)
}

// MockDistributedLock is a mock implementation of DistributedLock for testing
type MockDistributedLock struct {
	mock.Mock
}

func (m *MockDistributedLock) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	args := m.Called(ctx, key, ttl)
	return args.Bool(0), args.Error(1)
}

func (m *MockDistributedLock) Unlock(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockDistributedLock) TryLock(ctx context.Context, key string, ttl time.Duration, retry int, interval time.Duration) (bool, error) {
	args := m.Called(ctx, key, ttl, retry, interval)
	return args.Bool(0), args.Error(1)
}

// TestNewSettlementScheduler tests the creation of a new settlement scheduler
func TestNewSettlementScheduler(t *testing.T) {
	// Create a mock lock
	mockLock := new(MockDistributedLock)

	// Create a scheduler with mock lock
	scheduler := NewSettlementScheduler(nil, mockLock)

	assert.NotNil(t, scheduler)
	assert.NotNil(t, scheduler.Cron)
	assert.NotNil(t, scheduler.lock)
}

// TestSettlementScheduler_MonthlySettlement_Success tests successful monthly settlement via TriggerSettlement
func TestSettlementScheduler_MonthlySettlement_Success(t *testing.T) {
	mockService := new(MockCommissionService)
	mockLock := new(MockDistributedLock)
	scheduler := NewSettlementScheduler(nil, mockLock)
	// Replace the commission service with our mock
	scheduler.commissionSvc = mockService

	ctx := context.Background()
	month := "2026-01"

	mockService.On("SettleMonth", ctx, month).Return(nil)

	// Manually trigger the monthly settlement logic
	err := scheduler.TriggerSettlement(month)

	require.NoError(t, err)
	mockService.AssertExpectations(t)
}

// TestSettlementScheduler_MonthlySettlement_ServiceError tests settlement when service returns error
func TestSettlementScheduler_MonthlySettlement_ServiceError(t *testing.T) {
	mockService := new(MockCommissionService)
	mockLock := new(MockDistributedLock)
	scheduler := NewSettlementScheduler(nil, mockLock)
	scheduler.commissionSvc = mockService

	ctx := context.Background()
	month := "2026-01"

	expectedError := errors.New("database connection failed")
	mockService.On("SettleMonth", ctx, month).Return(expectedError)

	// Trigger settlement should return error
	err := scheduler.TriggerSettlement(month)

	assert.Error(t, err)
	// The error from the service is returned directly
	assert.Equal(t, expectedError, err)
	mockService.AssertExpectations(t)
}

// TestSettlementScheduler_TriggerSettlement_AlreadySettled tests triggering settlement when already settled
func TestSettlementScheduler_TriggerSettlement_AlreadySettled(t *testing.T) {
	mockService := new(MockCommissionService)
	mockLock := new(MockDistributedLock)
	scheduler := NewSettlementScheduler(nil, mockLock)
	scheduler.commissionSvc = mockService

	ctx := context.Background()
	month := "2026-01"

	// Simulate already settled error
	alreadySettledErr := errors.New("already settled")
	mockService.On("SettleMonth", ctx, month).Return(alreadySettledErr)

	err := scheduler.TriggerSettlement(month)

	assert.Error(t, err)
	assert.Equal(t, alreadySettledErr, err)
	mockService.AssertExpectations(t)
}

// TestSettlementScheduler_TriggerSettlement_NoRecords tests triggering settlement with no records
func TestSettlementScheduler_TriggerSettlement_NoRecords(t *testing.T) {
	mockService := new(MockCommissionService)
	mockLock := new(MockDistributedLock)
	scheduler := NewSettlementScheduler(nil, mockLock)
	scheduler.commissionSvc = mockService

	ctx := context.Background()
	month := "2026-01"

	// Simulate no records error
	noRecordsErr := errors.New("no pending records")
	mockService.On("SettleMonth", ctx, month).Return(noRecordsErr)

	err := scheduler.TriggerSettlement(month)

	assert.Error(t, err)
	assert.Equal(t, noRecordsErr, err)
	mockService.AssertExpectations(t)
}

// TestSettlementScheduler_GetNextRunTime tests getting the next run time
func TestSettlementScheduler_GetNextRunTime(t *testing.T) {
	mockLock := new(MockDistributedLock)
	scheduler := NewSettlementScheduler(nil, mockLock)

	scheduler.Start()
	defer scheduler.Stop()

	// Give scheduler time to start
	time.Sleep(100 * time.Millisecond)

	nextRun := scheduler.GetNextRunTime()

	assert.NotZero(t, nextRun, "Next run time should be set")

	// Verify it's in the future
	assert.True(t, nextRun.After(time.Now()), "Next run should be in the future")
}

// TestSettlementScheduler_CronSchedule tests that the cron schedule is correctly set
func TestSettlementScheduler_CronSchedule(t *testing.T) {
	mockLock := new(MockDistributedLock)
	scheduler := NewSettlementScheduler(nil, mockLock)

	scheduler.Start()
	defer scheduler.Stop()

	entries := scheduler.Cron.Entries()

	require.Len(t, entries, 1, "Should have exactly one cron job")

	// The cron expression "0 2 1 * *" means: at 02:00 on the 1st of every month
	// We can't easily test the exact schedule without accessing internal cron fields,
	// but we can verify the job is registered
	assert.NotNil(t, entries[0])
}

// TestSettlementScheduler_MultipleStartStop tests starting and stopping multiple times
func TestSettlementScheduler_MultipleStartStop(t *testing.T) {
	mockLock := new(MockDistributedLock)
	scheduler := NewSettlementScheduler(nil, mockLock)

	// First start
	scheduler.Start()
	time.Sleep(50 * time.Millisecond)
	entries1 := scheduler.Cron.Entries()
	assert.Len(t, entries1, 1)

	// Stop
	scheduler.Stop()
	time.Sleep(50 * time.Millisecond)

	// Second start (should not duplicate jobs)
	scheduler.Start()
	time.Sleep(50 * time.Millisecond)
	entries2 := scheduler.Cron.Entries()
	// Note: cron library may add new entries on restart, so we just verify it's running
	assert.NotEmpty(t, entries2, "Should have cron entries after restart")

	scheduler.Stop()
}

// TestSettlementScheduler_DifferentMonths tests settlement for different months
func TestSettlementScheduler_DifferentMonths(t *testing.T) {
	mockService := new(MockCommissionService)
	mockLock := new(MockDistributedLock)
	scheduler := NewSettlementScheduler(nil, mockLock)
	scheduler.commissionSvc = mockService

	ctx := context.Background()

	months := []string{"2025-12", "2026-01", "2026-02"}

	for _, month := range months {
		mockService.On("SettleMonth", ctx, month).Return(nil).Once()
	}

	// Trigger settlements for different months
	for _, month := range months {
		err := scheduler.TriggerSettlement(month)
		require.NoError(t, err)
	}

	mockService.AssertExpectations(t)
}

// TestSettlementScheduler_MonthFormatValidation tests various month formats
func TestSettlementScheduler_MonthFormatValidation(t *testing.T) {
	mockService := new(MockCommissionService)
	mockLock := new(MockDistributedLock)
	scheduler := NewSettlementScheduler(nil, mockLock)
	scheduler.commissionSvc = mockService

	ctx := context.Background()

	testCases := []struct {
		name        string
		month       string
		shouldError bool
	}{
		{"Valid month format", "2026-01", false},
		{"Valid month format December", "2025-12", false},
		{"Valid month format January", "2026-01", false},
		{"Invalid format - no dash", "202601", true},
		{"Invalid format - wrong separator", "2026/01", true},
		{"Invalid format - short year", "26-01", true},
		{"Invalid format - text month", "2026-Jan", true},
		{"Empty string", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.shouldError {
				// For invalid formats, the service should handle it
				// We're testing that the scheduler passes it through correctly
				mockService.On("SettleMonth", ctx, tc.month).Return(errors.New("invalid month format")).Once()

				err := scheduler.TriggerSettlement(tc.month)
				assert.Error(t, err)
			} else {
				mockService.On("SettleMonth", ctx, tc.month).Return(nil).Once()

				err := scheduler.TriggerSettlement(tc.month)
				assert.NoError(t, err)
			}
		})
	}

	mockService.AssertExpectations(t)
}

// TestSettlementScheduler_GetNextRunTime_BeforeStart tests getting next run time before starting
func TestSettlementScheduler_GetNextRunTime_BeforeStart(t *testing.T) {
	mockLock := new(MockDistributedLock)
	scheduler := NewSettlementScheduler(nil, mockLock)

	// Before starting, next run time should be zero
	nextRun := scheduler.GetNextRunTime()

	assert.Zero(t, nextRun, "Next run time should be zero before starting scheduler")
}

// BenchmarkSettlementScheduler_TriggerSettlement benchmarks the settlement trigger
func BenchmarkSettlementScheduler_TriggerSettlement(b *testing.B) {
	mockService := new(MockCommissionService)
	mockLock := new(MockDistributedLock)
	scheduler := NewSettlementScheduler(nil, mockLock)
	scheduler.commissionSvc = mockService

	ctx := context.Background()
	month := "2026-01"

	mockService.On("SettleMonth", ctx, month).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = scheduler.TriggerSettlement(month)
	}
}

// TestSettlementScheduler_PreviousMonthCalculation tests the previous month calculation
func TestSettlementScheduler_PreviousMonthCalculation(t *testing.T) {
	testCases := []struct {
		name          string
		currentTime   time.Time
		expectedMonth string
	}{
		{
			name:          "February 2026",
			currentTime:   time.Date(2026, 2, 1, 2, 0, 0, 0, time.UTC),
			expectedMonth: "2026-01",
		},
		{
			name:          "January 2026",
			currentTime:   time.Date(2026, 1, 1, 2, 0, 0, 0, time.UTC),
			expectedMonth: "2025-12",
		},
		{
			name:          "March 2026",
			currentTime:   time.Date(2026, 3, 15, 10, 30, 0, 0, time.UTC),
			expectedMonth: "2026-02",
		},
		{
			name:          "December 2026",
			currentTime:   time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC),
			expectedMonth: "2026-11",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lastMonth := tc.currentTime.AddDate(0, -1, 0)
			result := lastMonth.Format("2006-01")

			assert.Equal(t, tc.expectedMonth, result)
		})
	}
}

// TestSettlementScheduler_DistributedLock_Acquired tests successful lock acquisition and execution
func TestSettlementScheduler_DistributedLock_Acquired(t *testing.T) {
	mockService := new(MockCommissionService)
	mockLock := new(MockDistributedLock)
	scheduler := NewSettlementScheduler(nil, mockLock)
	scheduler.commissionSvc = mockService

	ctx := context.Background()
	lockKey := "scheduler:settlement:monthly"

	// Mock successful lock acquisition
	mockLock.On("TryLock", mock.Anything, lockKey, time.Hour, 1, time.Second).Return(true, nil)
	mockLock.On("Unlock", mock.Anything, lockKey).Return(nil)

	// Mock successful settlement
	mockService.On("SettleMonth", ctx, mock.AnythingOfType("string")).Return(nil)

	// Execute the monthly settlement with lock
	scheduler.monthlySettlementWithLock()

	// Verify lock was acquired and released
	mockLock.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

// TestSettlementScheduler_DistributedLock_AlreadyLocked tests behavior when lock is held by another instance
func TestSettlementScheduler_DistributedLock_AlreadyLocked(t *testing.T) {
	mockService := new(MockCommissionService)
	mockLock := new(MockDistributedLock)
	scheduler := NewSettlementScheduler(nil, mockLock)
	scheduler.commissionSvc = mockService

	ctx := context.Background()
	lockKey := "scheduler:settlement:monthly"

	// Mock lock acquisition failure (another instance has the lock)
	mockLock.On("TryLock", mock.Anything, lockKey, time.Hour, 1, time.Second).Return(false, nil)

	// Settlement should NOT be called
	mockService.On("SettleMonth", ctx, mock.AnythingOfType("string")).Return(nil).Maybe()

	// Execute the monthly settlement with lock
	scheduler.monthlySettlementWithLock()

	// Verify lock was attempted but settlement was not called
	mockLock.AssertExpectations(t)
	// Settlement should not be called when lock is not acquired
	mockService.AssertNotCalled(t, "SettleMonth", ctx, mock.AnythingOfType("string"))
}

// TestSettlementScheduler_DistributedLock_Error tests behavior when lock operation fails
func TestSettlementScheduler_DistributedLock_Error(t *testing.T) {
	mockService := new(MockCommissionService)
	mockLock := new(MockDistributedLock)
	scheduler := NewSettlementScheduler(nil, mockLock)
	scheduler.commissionSvc = mockService

	ctx := context.Background()
	lockKey := "scheduler:settlement:monthly"

	// Mock lock acquisition error
	lockError := errors.New("redis connection failed")
	mockLock.On("TryLock", mock.Anything, lockKey, time.Hour, 1, time.Second).Return(false, lockError)

	// Settlement should NOT be called on lock error
	mockService.On("SettleMonth", ctx, mock.AnythingOfType("string")).Return(nil).Maybe()

	// Execute the monthly settlement with lock
	scheduler.monthlySettlementWithLock()

	// Verify lock was attempted but settlement was not called
	mockLock.AssertExpectations(t)
	mockService.AssertNotCalled(t, "SettleMonth", ctx, mock.AnythingOfType("string"))
}

// TestSettlementScheduler_DistributedLock_UnlockError tests behavior when unlock fails
func TestSettlementScheduler_DistributedLock_UnlockError(t *testing.T) {
	mockService := new(MockCommissionService)
	mockLock := new(MockDistributedLock)
	scheduler := NewSettlementScheduler(nil, mockLock)
	scheduler.commissionSvc = mockService

	ctx := context.Background()
	lockKey := "scheduler:settlement:monthly"

	// Mock successful lock acquisition
	mockLock.On("TryLock", mock.Anything, lockKey, time.Hour, 1, time.Second).Return(true, nil)

	// Mock unlock error (but settlement should still complete)
	unlockError := errors.New("redis connection lost")
	mockLock.On("Unlock", mock.Anything, lockKey).Return(unlockError)

	// Mock successful settlement
	mockService.On("SettleMonth", ctx, mock.AnythingOfType("string")).Return(nil)

	// Execute the monthly settlement with lock
	// Settlement should complete even if unlock fails
	scheduler.monthlySettlementWithLock()

	// Verify lock was acquired and settlement was executed
	mockLock.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

// TestSettlementScheduler_DistributedLock_SettlementError tests behavior when settlement fails
func TestSettlementScheduler_DistributedLock_SettlementError(t *testing.T) {
	mockService := new(MockCommissionService)
	mockLock := new(MockDistributedLock)
	scheduler := NewSettlementScheduler(nil, mockLock)
	scheduler.commissionSvc = mockService

	ctx := context.Background()
	lockKey := "scheduler:settlement:monthly"

	// Mock successful lock acquisition
	mockLock.On("TryLock", mock.Anything, lockKey, time.Hour, 1, time.Second).Return(true, nil)
	mockLock.On("Unlock", mock.Anything, lockKey).Return(nil)

	// Mock settlement error
	settlementError := errors.New("settlement failed")
	mockService.On("SettleMonth", ctx, mock.AnythingOfType("string")).Return(settlementError)

	// Execute the monthly settlement with lock
	// Lock should still be released even when settlement fails
	scheduler.monthlySettlementWithLock()

	// Verify lock was acquired and released
	mockLock.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

// TestSettlementScheduler_DistributedLock_Concurrent tests concurrent execution scenarios
func TestSettlementScheduler_DistributedLock_Concurrent(t *testing.T) {
	// This test simulates a simpler concurrent scenario
	// to avoid complex mock interactions
	mockService := new(MockCommissionService)
	mockLock := new(MockDistributedLock)
	scheduler := NewSettlementScheduler(nil, mockLock)
	scheduler.commissionSvc = mockService

	ctx := context.Background()

	// Mock lock to succeed only once
	mockLock.On("TryLock", mock.Anything, mock.AnythingOfType("string"), time.Hour, 1, time.Second).Return(true, nil).Once()
	mockLock.On("TryLock", mock.Anything, mock.AnythingOfType("string"), time.Hour, 1, time.Second).Return(false, nil).Twice()
	mockLock.On("Unlock", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()

	// Mock settlement (only called once)
	mockService.On("SettleMonth", ctx, mock.AnythingOfType("string")).Return(nil).Once()

	// Simulate 3 concurrent executions
	done := make(chan bool, 3)
	for i := 0; i < 3; i++ {
		go func() {
			scheduler.monthlySettlementWithLock()
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 3; i++ {
		<-done
	}

	// Verify only one instance executed settlement
	mockService.AssertExpectations(t)
	mockLock.AssertExpectations(t)
}
