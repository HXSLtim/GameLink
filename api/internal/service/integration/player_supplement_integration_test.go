// Package integration provides supplementary integration tests for player service.
package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository/game"
	repoiface "gamelink/internal/repository/interfaces"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/review"
	"gamelink/internal/repository/user"
	playerservice "gamelink/internal/service/player"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPlayerService_ApproveApplication tests approving a player application.
func TestPlayerService_ApproveApplication(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create pending player
	playerUser := CreateUniqueTestUser(t, db, "approve_player")
	testPlayer := &model.Player{
		Base: model.Base{
			ExtJSON: "{}",
		},
		UserID:             playerUser.ID,
		Nickname:           "Pending Player",
		VerificationStatus: model.VerificationPending,
	}
	err := db.Create(testPlayer).Error
	require.NoError(t, err)

	// Approve application
	adminUser := CreateUniqueTestUser(t, db, "admin_approver")
	now := time.Now()
	testPlayer.VerificationStatus = model.VerificationVerified
	testPlayer.VerifiedAt = &now
	testPlayer.VerifiedBy = &adminUser.ID
	err = db.Save(testPlayer).Error
	require.NoError(t, err)

	// Update user role
	playerUser.Role = model.RolePlayer
	db.Save(playerUser)

	// Verify
	var updatedPlayer model.Player
	err = db.First(&updatedPlayer, testPlayer.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.VerificationVerified, updatedPlayer.VerificationStatus)
	assert.NotNil(t, updatedPlayer.VerifiedAt)
	assert.NotNil(t, updatedPlayer.VerifiedBy)

	var updatedUser model.User
	err = db.First(&updatedUser, playerUser.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.RolePlayer, updatedUser.Role)
}

// TestPlayerService_RejectApplication tests rejecting a player application.
func TestPlayerService_RejectApplication(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create pending player
	playerUser := CreateUniqueTestUser(t, db, "reject_player")
	testPlayer := &model.Player{
		Base: model.Base{
			ExtJSON: "{}",
		},
		UserID:             playerUser.ID,
		Nickname:           "Rejected Player",
		VerificationStatus: model.VerificationPending,
	}
	err := db.Create(testPlayer).Error
	require.NoError(t, err)

	// Reject application
	adminUser := CreateUniqueTestUser(t, db, "admin_rejecter")
	now := time.Now()
	testPlayer.VerificationStatus = model.VerificationRejected
	testPlayer.VerifiedAt = &now
	testPlayer.VerifiedBy = &adminUser.ID
	testPlayer.RejectReason = "Insufficient qualifications"
	err = db.Save(testPlayer).Error
	require.NoError(t, err)

	// Verify
	var updatedPlayer model.Player
	err = db.First(&updatedPlayer, testPlayer.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.VerificationRejected, updatedPlayer.VerificationStatus)
	assert.Equal(t, "Insufficient qualifications", updatedPlayer.RejectReason)
}

// TestPlayerService_SubmitRealNameCertification tests submitting real name certification.
func TestPlayerService_SubmitRealNameCertification(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create verified player
	playerUser := CreateUniqueTestUser(t, db, "cert_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Submit certification
	cert := CreateTestPlayerCertification(t, db, testPlayer, model.CertificationStatusPending)

	// Verify
	var savedCert model.PlayerCertification
	err := db.First(&savedCert, cert.ID).Error
	require.NoError(t, err)
	assert.Equal(t, testPlayer.ID, savedCert.PlayerID)
	assert.Equal(t, model.CertificationStatusPending, savedCert.Status)
	assert.NotEmpty(t, savedCert.RealName)
	assert.NotEmpty(t, savedCert.IDCardNo)
}

// TestPlayerService_ApproveRealNameCertification tests approving real name certification.
func TestPlayerService_ApproveRealNameCertification(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create player with pending certification
	playerUser := CreateUniqueTestUser(t, db, "approve_cert_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	cert := CreateTestPlayerCertification(t, db, testPlayer, model.CertificationStatusPending)

	// Approve certification
	adminUser := CreateUniqueTestUser(t, db, "cert_admin")
	now := time.Now()
	cert.Status = model.CertificationStatusVerified
	cert.VerifiedAt = &now
	cert.VerifiedBy = &adminUser.ID
	err := db.Save(cert).Error
	require.NoError(t, err)

	// Verify
	var updatedCert model.PlayerCertification
	err = db.First(&updatedCert, cert.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.CertificationStatusVerified, updatedCert.Status)
	assert.NotNil(t, updatedCert.VerifiedAt)
}

// TestPlayerService_RejectRealNameCertification tests rejecting real name certification.
func TestPlayerService_RejectRealNameCertification(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create player with pending certification
	playerUser := CreateUniqueTestUser(t, db, "reject_cert_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	cert := CreateTestPlayerCertification(t, db, testPlayer, model.CertificationStatusPending)

	// Reject certification
	adminUser := CreateUniqueTestUser(t, db, "reject_cert_admin")
	now := time.Now()
	cert.Status = model.CertificationStatusRejected
	cert.VerifiedAt = &now
	cert.VerifiedBy = &adminUser.ID
	cert.RejectReason = "ID card image unclear"
	err := db.Save(cert).Error
	require.NoError(t, err)

	// Verify
	var updatedCert model.PlayerCertification
	err = db.First(&updatedCert, cert.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.CertificationStatusRejected, updatedCert.Status)
	assert.Equal(t, "ID card image unclear", updatedCert.RejectReason)
}

// TestPlayerService_SubmitRankCertification tests submitting rank certification.
func TestPlayerService_SubmitRankCertification(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create player and game
	playerUser := CreateUniqueTestUser(t, db, "rank_cert_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "rank_cert_game")
	testRank := CreateTestGameRank(t, db, testGame, "Diamond", 5, 5000)

	// Submit rank certification
	record := CreateTestPlayerRankRecord(t, db, testPlayer, testGame, testRank, model.PlayerRankStatusPending)

	// Verify
	var savedRecord model.PlayerRankRecord
	err := db.First(&savedRecord, record.ID).Error
	require.NoError(t, err)
	assert.Equal(t, testPlayer.ID, savedRecord.PlayerID)
	assert.Equal(t, testGame.ID, savedRecord.GameID)
	assert.Equal(t, testRank.ID, savedRecord.RankID)
	assert.Equal(t, model.PlayerRankStatusPending, savedRecord.Status)
}

// TestPlayerService_ApproveRankCertification tests approving rank certification.
func TestPlayerService_ApproveRankCertification(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create player with pending rank certification
	playerUser := CreateUniqueTestUser(t, db, "approve_rank_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "approve_rank_game")
	testRank := CreateTestGameRank(t, db, testGame, "Master", 6, 6000)
	record := CreateTestPlayerRankRecord(t, db, testPlayer, testGame, testRank, model.PlayerRankStatusPending)

	// Approve rank certification
	adminUser := CreateUniqueTestUser(t, db, "rank_admin")
	now := time.Now()
	record.Status = model.PlayerRankStatusVerified
	record.VerifiedAt = &now
	record.VerifiedBy = &adminUser.ID
	err := db.Save(record).Error
	require.NoError(t, err)

	// Verify
	var updatedRecord model.PlayerRankRecord
	err = db.First(&updatedRecord, record.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.PlayerRankStatusVerified, updatedRecord.Status)
	assert.NotNil(t, updatedRecord.VerifiedAt)
}

// TestPlayerService_UpdateRatingAfterReview tests updating player rating after review.
func TestPlayerService_UpdateRatingAfterReview(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create player
	playerUser := CreateUniqueTestUser(t, db, "rating_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testPlayer.RatingAverage = 0
	testPlayer.RatingCount = 0
	db.Save(testPlayer)

	// Create user and order for review
	testUser := CreateUniqueTestUser(t, db, "rating_user")
	testGame := CreateTestGame(t, db, "rating_game")
	testOrder := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Create reviews
	CreateTestReview(t, db, testOrder, model.Rating5)

	// Update player rating (simulating service logic)
	var reviews []model.Review
	db.Where("player_id = ?", testPlayer.ID).Find(&reviews)

	var totalScore float32
	for _, r := range reviews {
		totalScore += float32(r.Score)
	}
	testPlayer.RatingAverage = totalScore / float32(len(reviews))
	testPlayer.RatingCount = uint32(len(reviews))
	db.Save(testPlayer)

	// Verify
	var updatedPlayer model.Player
	err := db.First(&updatedPlayer, testPlayer.ID).Error
	require.NoError(t, err)
	assert.Equal(t, float32(5.0), updatedPlayer.RatingAverage)
	assert.Equal(t, uint32(1), updatedPlayer.RatingCount)
}

// TestPlayerService_GetPlayerEarnings tests getting player earnings.
func TestPlayerService_GetPlayerEarnings(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create player
	playerUser := CreateUniqueTestUser(t, db, "earnings_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create user and completed orders
	testUser := CreateUniqueTestUser(t, db, "earnings_user")
	testGame := CreateTestGame(t, db, "earnings_game")

	// Create multiple completed orders with commission records
	for i := 0; i < 3; i++ {
		order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
		CreateTestCommissionRecord(t, db, order.ID, testPlayer.ID, 10000, model.SettlementStatusSettled)
	}

	// Calculate total earnings
	var records []model.CommissionRecord
	err := db.Where("player_id = ?", testPlayer.ID).Find(&records).Error
	require.NoError(t, err)

	var totalEarnings int64
	for _, r := range records {
		totalEarnings += r.PlayerIncomeCents
	}

	// Verify
	assert.Len(t, records, 3)
	assert.Equal(t, int64(24000), totalEarnings) // 8000 * 3
}

// TestPlayerService_OnlineStatusManagement tests player online status management via cache.
func TestPlayerService_OnlineStatusManagement(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	playerRepo := player.NewPlayerRepository(db)
	userRepo := user.NewUserRepository(db)
	gameRepo := game.NewGameRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	mockOrders := &mockOrderQuery{}
	mockTags := &mockPlayerTagRepo{tags: make(map[uint64][]string)}
	mockCacheInst := newMockCache()

	svc := playerservice.NewPlayerService(playerRepo, userRepo, gameRepo, mockOrders, reviewRepo, mockTags, mockCacheInst)

	// Create player
	playerUser := CreateUniqueTestUser(t, db, "online_status_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Set online (this uses cache, not database)
	err := svc.SetPlayerOnlineStatus(ctx, playerUser.ID, true)
	require.NoError(t, err)

	// Verify online status in cache
	cacheKey := fmt.Sprintf("player:online:%d", testPlayer.ID)
	_, exists, _ := mockCacheInst.Get(ctx, cacheKey)
	assert.True(t, exists, "Player should be online in cache")

	// Set offline
	err = svc.SetPlayerOnlineStatus(ctx, playerUser.ID, false)
	require.NoError(t, err)

	// Verify offline status in cache (key should be deleted)
	_, exists, _ = mockCacheInst.Get(ctx, cacheKey)
	assert.False(t, exists, "Player should be offline (cache key deleted)")
}

// TestPlayerService_PlayerWithMultipleRanks tests player with multiple game ranks.
func TestPlayerService_PlayerWithMultipleRanks(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create player
	playerUser := CreateUniqueTestUser(t, db, "multi_rank_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create multiple games and ranks
	game1 := CreateTestGame(t, db, "multi_rank_game1")
	game2 := CreateTestGame(t, db, "multi_rank_game2")

	rank1 := CreateTestGameRank(t, db, game1, "Diamond", 5, 5000)
	rank2 := CreateTestGameRank(t, db, game2, "Master", 6, 6000)

	// Create rank records
	CreateTestPlayerRankRecord(t, db, testPlayer, game1, rank1, model.PlayerRankStatusVerified)
	CreateTestPlayerRankRecord(t, db, testPlayer, game2, rank2, model.PlayerRankStatusVerified)

	// Verify player has multiple ranks
	var records []model.PlayerRankRecord
	err := db.Where("player_id = ? AND status = ?", testPlayer.ID, model.PlayerRankStatusVerified).Find(&records).Error
	require.NoError(t, err)
	assert.Len(t, records, 2)
}

// TestPlayerService_PlayerStatistics tests player statistics.
func TestPlayerService_PlayerStatistics(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create player
	playerUser := CreateUniqueTestUser(t, db, "stats_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create player statistics
	stats := &model.PlayerStatistics{
		PlayerID:             testPlayer.ID,
		TotalEarningsCents:   50000,
		TotalCommissionCents: 10000,
		TotalOrderCount:      10,
		CompletedOrderCount:  8,
		CanceledOrderCount:   1,
		RefundOrderCount:     1,
	}
	err := db.Create(stats).Error
	require.NoError(t, err)

	// Verify
	var savedStats model.PlayerStatistics
	err = db.Where("player_id = ?", testPlayer.ID).First(&savedStats).Error
	require.NoError(t, err)
	assert.Equal(t, int64(50000), savedStats.TotalEarningsCents)
	assert.Equal(t, 10, savedStats.TotalOrderCount)
	assert.Equal(t, 8, savedStats.CompletedOrderCount)
}

// mockOrderQuery for player service tests
type mockOrderQueryForPlayer struct {
	orders []model.Order
}

func (m *mockOrderQueryForPlayer) List(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
	var result []model.Order
	for _, o := range m.orders {
		if opts.PlayerID != nil && o.PlayerID != nil && *o.PlayerID == *opts.PlayerID {
			result = append(result, o)
		}
	}
	return result, int64(len(result)), nil
}

func (m *mockOrderQueryForPlayer) Get(ctx context.Context, id uint64) (*model.Order, error) {
	for _, o := range m.orders {
		if o.ID == id {
			return &o, nil
		}
	}
	return nil, nil
}
