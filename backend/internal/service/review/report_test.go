package review

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	orderimpl "gamelink/internal/repository/implementations"
	"gamelink/internal/repository/review"
	"gamelink/internal/repository/reviewreport"
	"gamelink/pkg/testutil"
)

func setupReportTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.User{},
		&model.Player{},
		&model.Game{},
		&model.ServiceItem{},
		&model.Order{},
		&model.Review{},
		&model.ReviewReport{},
		&model.ReviewReply{},
	)
	return db
}

func createTestReportService(db *gorm.DB) *ReviewService {
	reviewRepo := review.NewReviewRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	reportRepo := reviewreport.NewReviewReportRepository(db)

	// Create mock repositories for users and players
	userRepo := &mockUserRepository{db: db}
	playerRepo := &mockPlayerRepository{db: db}
	replyRepo := &mockReviewReplyRepository{}
	notificationRepo := &mockNotificationRepository{}
	opLogRepo := &mockOperationLogRepository{}

	return NewReviewService(reviewRepo, orderRepo, playerRepo, userRepo, replyRepo, reportRepo, notificationRepo, opLogRepo)
}

// Mock repositories
type mockUserRepository struct {
	db *gorm.DB
}

func (m *mockUserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	if err := m.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (m *mockUserRepository) List(ctx context.Context) ([]model.User, error) {
	return nil, nil
}

func (m *mockUserRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	return nil, 0, nil
}

func (m *mockUserRepository) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
	return nil, 0, nil
}

func (m *mockUserRepository) Count(ctx context.Context, opts repository.UserListOptions) (int, error) {
	return 0, nil
}

func (m *mockUserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, nil
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}

func (m *mockUserRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, nil
}

func (m *mockUserRepository) Create(ctx context.Context, user *model.User) error {
	return m.db.WithContext(ctx).Create(user).Error
}

func (m *mockUserRepository) Update(ctx context.Context, user *model.User) error {
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}

func (m *mockUserRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.User, error) {
	var users []model.User
	if len(ids) == 0 {
		return users, nil
	}
	if err := m.db.WithContext(ctx).Where("id IN ?", ids).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

type mockPlayerRepository struct {
	db *gorm.DB
}

func (m *mockPlayerRepository) Get(ctx context.Context, id uint64) (*model.Player, error) {
	var player model.Player
	if err := m.db.WithContext(ctx).First(&player, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &player, nil
}

func (m *mockPlayerRepository) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	return nil, nil
}

func (m *mockPlayerRepository) List(ctx context.Context) ([]model.Player, error) {
	return nil, nil
}

func (m *mockPlayerRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	return nil, 0, nil
}

func (m *mockPlayerRepository) Create(ctx context.Context, player *model.Player) error {
	return m.db.WithContext(ctx).Create(player).Error
}

func (m *mockPlayerRepository) Update(ctx context.Context, player *model.Player) error {
	return nil
}

func (m *mockPlayerRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}

func (m *mockPlayerRepository) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, status *model.VerificationStatus) ([]model.Player, int64, error) {
	return nil, 0, nil
}

func (m *mockPlayerRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, status model.VerificationStatus) (int64, error) {
	return 0, nil
}

func (m *mockPlayerRepository) BatchDelete(ctx context.Context, ids []uint64) (int64, error) {
	return 0, nil
}

type mockReviewReplyRepository struct{}

func (m *mockReviewReplyRepository) Create(ctx context.Context, reply *model.ReviewReply) error {
	return nil
}

func (m *mockReviewReplyRepository) Get(ctx context.Context, replyID uint64) (*model.ReviewReply, error) {
	return nil, nil
}

func (m *mockReviewReplyRepository) ListByReview(ctx context.Context, reviewID uint64) ([]model.ReviewReply, error) {
	return nil, nil
}

func (m *mockReviewReplyRepository) Update(ctx context.Context, reply *model.ReviewReply) error {
	return nil
}

func (m *mockReviewReplyRepository) Delete(ctx context.Context, replyID uint64) error {
	return nil
}

func (m *mockReviewReplyRepository) UpdateStatus(ctx context.Context, replyID uint64, status string, note string) error {
	return nil
}

type mockNotificationRepository struct{}

func (m *mockNotificationRepository) ListByUser(ctx context.Context, opts repository.NotificationListOptions) ([]model.NotificationEvent, int64, error) {
	return nil, 0, nil
}

func (m *mockNotificationRepository) MarkRead(ctx context.Context, userID uint64, ids []uint64) error {
	return nil
}

func (m *mockNotificationRepository) MarkAllRead(ctx context.Context, userID uint64) error {
	return nil
}

func (m *mockNotificationRepository) CountUnread(ctx context.Context, userID uint64) (int64, error) {
	return 0, nil
}

func (m *mockNotificationRepository) Create(ctx context.Context, event *model.NotificationEvent) error {
	return nil
}

func (m *mockNotificationRepository) Delete(ctx context.Context, userID uint64, id uint64) error {
	return nil
}

type mockOperationLogRepository struct{}

func (m *mockOperationLogRepository) Append(ctx context.Context, log *model.OperationLog) error {
	return nil
}

func (m *mockOperationLogRepository) ListByEntity(ctx context.Context, entityType string, entityID uint64, opts repository.OperationLogListOptions) ([]model.OperationLog, int64, error) {
	return nil, 0, nil
}

func (m *mockOperationLogRepository) List(ctx context.Context, opts repository.OperationLogSearchOptions) ([]model.OperationLog, int64, error) {
	return nil, 0, nil
}

func TestReviewService_ReportReview(t *testing.T) {
	db := setupReportTestDB(t)
	defer testutil.CleanDB(t, db)

	svc := createTestReportService(db)
	ctx := context.Background()

	// Create test data
	user := &model.User{
		Name:  "Test User",
		Email: "test@example.com",
		Phone: "1234567890",
		Role:  model.RoleUser,
	}
	require.NoError(t, db.Create(user).Error)

	player := &model.Player{
		UserID:   user.ID,
		Nickname: "Test Player",
	}
	require.NoError(t, db.Create(player).Error)

	testReview := &model.Review{
		OrderID:  1001,
		UserID:   user.ID,
		PlayerID: player.ID,
		Score:    5,
		Content:  "Great service!",
		Status:   model.ReviewStatusApproved,
	}
	require.NoError(t, db.Create(testReview).Error)

	t.Run("report review successfully", func(t *testing.T) {
		req := ReportReviewRequest{
			Reason:   "评价内容包含不实信息",
			Evidence: "https://example.com/evidence.jpg",
		}

		resp, err := svc.ReportReview(ctx, testReview.ID, user.ID, req)
		require.NoError(t, err)
		assert.NotZero(t, resp.ReportID)

		// Verify review is marked as reported
		updatedReview, err := svc.reviews.Get(ctx, testReview.ID)
		require.NoError(t, err)
		assert.True(t, updatedReview.IsReported)
	})

	t.Run("report non-existent review", func(t *testing.T) {
		req := ReportReviewRequest{
			Reason: "Test reason",
		}

		_, err := svc.ReportReview(ctx, 99999, user.ID, req)
		assert.ErrorIs(t, err, ErrReviewNotFound)
	})

	t.Run("report with invalid reason", func(t *testing.T) {
		req := ReportReviewRequest{
			Reason: "", // Empty reason
		}

		_, err := svc.ReportReview(ctx, testReview.ID, user.ID, req)
		assert.Error(t, err)
	})
}

func TestReviewService_ListReports(t *testing.T) {
	db := setupReportTestDB(t)
	defer testutil.CleanDB(t, db)

	svc := createTestReportService(db)
	ctx := context.Background()

	// Create test data
	user := &model.User{
		Name:  "Test User",
		Email: "test@example.com",
		Phone: "1234567890",
		Role:  model.RoleUser,
	}
	require.NoError(t, db.Create(user).Error)

	// Create multiple reports
	for i := 0; i < 3; i++ {
		report := &model.ReviewReport{
			ReviewID:   uint64(1001 + i),
			ReporterID: user.ID,
			Reason:     "Test reason",
			Status:     model.ReviewReportStatusPending,
		}
		require.NoError(t, db.Create(report).Error)
	}

	t.Run("list all reports", func(t *testing.T) {
		req := ListReportsRequest{
			Page:     1,
			PageSize: 10,
		}

		resp, err := svc.ListReports(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, int64(3), resp.Total)
		assert.Len(t, resp.Reports, 3)
	})

	t.Run("filter by status", func(t *testing.T) {
		status := model.ReviewReportStatusPending
		req := ListReportsRequest{
			Page:     1,
			PageSize: 10,
			Status:   &status,
		}

		resp, err := svc.ListReports(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, int64(3), resp.Total)
	})

	t.Run("pagination", func(t *testing.T) {
		req := ListReportsRequest{
			Page:     1,
			PageSize: 2,
		}

		resp, err := svc.ListReports(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, int64(3), resp.Total)
		assert.Len(t, resp.Reports, 2)
	})
}

func TestReviewService_HandleReport(t *testing.T) {
	db := setupReportTestDB(t)
	defer testutil.CleanDB(t, db)

	svc := createTestReportService(db)
	ctx := context.Background()

	// Create test data
	user := &model.User{
		Name:  "Test User",
		Email: "test@example.com",
		Phone: "1234567890",
		Role:  model.RoleUser,
	}
	require.NoError(t, db.Create(user).Error)

	admin := &model.User{
		Name:  "Admin",
		Email: "admin@example.com",
		Phone: "0987654321",
		Role:  model.RoleAdmin,
	}
	require.NoError(t, db.Create(admin).Error)

	player := &model.Player{
		UserID:   user.ID,
		Nickname: "Test Player",
	}
	require.NoError(t, db.Create(player).Error)

	testReview := &model.Review{
		OrderID:    1001,
		UserID:     user.ID,
		PlayerID:   player.ID,
		Score:      5,
		Content:    "Great service!",
		Status:     model.ReviewStatusApproved,
		IsReported: true,
	}
	require.NoError(t, db.Create(testReview).Error)

	t.Run("delete review", func(t *testing.T) {
		report := &model.ReviewReport{
			ReviewID:   testReview.ID,
			ReporterID: user.ID,
			Reason:     "不实信息",
			Status:     model.ReviewReportStatusPending,
		}
		require.NoError(t, db.Create(report).Error)

		req := HandleReportRequest{
			Action: "delete",
			Note:   "经核实，评价内容确实存在不实信息",
		}

		resp, err := svc.HandleReport(ctx, report.ID, admin.ID, req)
		require.NoError(t, err)
		assert.Equal(t, "success", resp.Status)
		assert.Equal(t, "评价已删除", resp.Message)

		// Verify review is deleted
		updatedReview, err := svc.reviews.Get(ctx, testReview.ID)
		require.NoError(t, err)
		assert.Equal(t, model.ReviewStatusDeleted, updatedReview.Status)

		// Verify report is handled
		updatedReport, err := svc.reports.Get(ctx, report.ID)
		require.NoError(t, err)
		assert.Equal(t, model.ReviewReportStatusApproved, updatedReport.Status)
		assert.NotNil(t, updatedReport.HandledBy)
		assert.Equal(t, admin.ID, *updatedReport.HandledBy)
		assert.NotNil(t, updatedReport.HandledAt)
	})

	t.Run("warn reviewer", func(t *testing.T) {
		review2 := &model.Review{
			OrderID:    1002,
			UserID:     user.ID,
			PlayerID:   player.ID,
			Score:      4,
			Content:    "Good",
			Status:     model.ReviewStatusApproved,
			IsReported: true,
		}
		require.NoError(t, db.Create(review2).Error)

		report := &model.ReviewReport{
			ReviewID:   review2.ID,
			ReporterID: user.ID,
			Reason:     "轻微不当内容",
			Status:     model.ReviewReportStatusPending,
		}
		require.NoError(t, db.Create(report).Error)

		req := HandleReportRequest{
			Action: "warn",
			Note:   "已警告用户注意评价内容",
		}

		resp, err := svc.HandleReport(ctx, report.ID, admin.ID, req)
		require.NoError(t, err)
		assert.Equal(t, "success", resp.Status)
		assert.Equal(t, "已警告评价者", resp.Message)

		// Verify report is handled
		updatedReport, err := svc.reports.Get(ctx, report.ID)
		require.NoError(t, err)
		assert.Equal(t, model.ReviewReportStatusApproved, updatedReport.Status)
	})

	t.Run("reject report", func(t *testing.T) {
		review3 := &model.Review{
			OrderID:    1003,
			UserID:     user.ID,
			PlayerID:   player.ID,
			Score:      5,
			Content:    "Excellent",
			Status:     model.ReviewStatusApproved,
			IsReported: true,
		}
		require.NoError(t, db.Create(review3).Error)

		report := &model.ReviewReport{
			ReviewID:   review3.ID,
			ReporterID: user.ID,
			Reason:     "无理由举报",
			Status:     model.ReviewReportStatusPending,
		}
		require.NoError(t, db.Create(report).Error)

		req := HandleReportRequest{
			Action: "reject",
			Note:   "举报不成立",
		}

		resp, err := svc.HandleReport(ctx, report.ID, admin.ID, req)
		require.NoError(t, err)
		assert.Equal(t, "success", resp.Status)
		assert.Equal(t, "举报已驳回", resp.Message)

		// Verify report is rejected
		updatedReport, err := svc.reports.Get(ctx, report.ID)
		require.NoError(t, err)
		assert.Equal(t, model.ReviewReportStatusRejected, updatedReport.Status)
	})

	t.Run("handle non-existent report", func(t *testing.T) {
		req := HandleReportRequest{
			Action: "delete",
		}

		_, err := svc.HandleReport(ctx, 99999, admin.ID, req)
		assert.ErrorIs(t, err, ErrReportNotFound)
	})

	t.Run("handle already handled report", func(t *testing.T) {
		report := &model.ReviewReport{
			ReviewID:   testReview.ID,
			ReporterID: user.ID,
			Reason:     "Test",
			Status:     model.ReviewReportStatusApproved,
		}
		require.NoError(t, db.Create(report).Error)

		req := HandleReportRequest{
			Action: "delete",
		}

		_, err := svc.HandleReport(ctx, report.ID, admin.ID, req)
		assert.ErrorIs(t, err, ErrReportAlreadyHandled)
	})

	t.Run("invalid action", func(t *testing.T) {
		report := &model.ReviewReport{
			ReviewID:   testReview.ID,
			ReporterID: user.ID,
			Reason:     "Test",
			Status:     model.ReviewReportStatusPending,
		}
		require.NoError(t, db.Create(report).Error)

		req := HandleReportRequest{
			Action: "invalid",
		}

		_, err := svc.HandleReport(ctx, report.ID, admin.ID, req)
		assert.ErrorIs(t, err, ErrInvalidReportAction)
	})
}
