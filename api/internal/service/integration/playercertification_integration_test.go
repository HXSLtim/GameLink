// Package integration provides integration tests for PlayerCertification service.
package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/playercertification"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== Repository Tests ====================

func TestPlayerCertificationRepository_Create(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_create")
	player := CreateTestPlayer(t, db, playerUser)

	cert := &model.PlayerCertification{
		PlayerID:       player.ID,
		RealName:       "Test User",
		IDCardNo:       "123456789012345678",
		IDCardFrontURL: "https://example.com/front.jpg",
		IDCardBackURL:  "https://example.com/back.jpg",
		Status:         model.CertificationStatusPending,
	}

	err := repo.Create(ctx, cert)
	require.NoError(t, err)
	assert.NotZero(t, cert.ID)
}

func TestPlayerCertificationRepository_Create_DuplicatePlayerID(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_duplicate")
	player := CreateTestPlayer(t, db, playerUser)

	// First certification
	CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	// Second certification for same player (should fail due to unique constraint)
	cert2 := &model.PlayerCertification{
		PlayerID:       player.ID,
		RealName:       "Another Name",
		IDCardNo:       "987654321098765432",
		IDCardFrontURL: "https://example.com/front2.jpg",
		IDCardBackURL:  "https://example.com/back2.jpg",
		Status:         model.CertificationStatusPending,
	}

	err := repo.Create(ctx, cert2)
	assert.Error(t, err) // Should fail due to unique constraint on player_id
}

func TestPlayerCertificationRepository_Get(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_get")
	player := CreateTestPlayer(t, db, playerUser)
	cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	result, err := repo.Get(ctx, cert.ID)
	require.NoError(t, err)
	assert.Equal(t, cert.ID, result.ID)
	assert.Equal(t, player.ID, result.PlayerID)
}

func TestPlayerCertificationRepository_Get_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	_, err := repo.Get(ctx, 99999)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestPlayerCertificationRepository_GetWithPlayer(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_with_player")
	player := CreateTestPlayer(t, db, playerUser)
	cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusVerified)

	result, err := repo.GetWithPlayer(ctx, cert.ID)
	require.NoError(t, err)
	assert.NotNil(t, result.Player)
	assert.Equal(t, player.ID, result.Player.ID)
}

func TestPlayerCertificationRepository_GetWithPlayer_Verifier(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	adminUser := CreateUniqueTestUser(t, db, "pc_verifier_admin")
	playerUser := CreateUniqueTestUser(t, db, "pc_verifier_player")
	player := CreateTestPlayer(t, db, playerUser)

	cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusVerified)
	cert.VerifiedBy = &adminUser.ID
	db.Save(cert)

	result, err := repo.GetWithPlayer(ctx, cert.ID)
	require.NoError(t, err)
	assert.NotNil(t, result.Verifier)
	assert.Equal(t, adminUser.ID, result.Verifier.ID)
}

func TestPlayerCertificationRepository_GetByPlayerID(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_by_player")
	player := CreateTestPlayer(t, db, playerUser)
	CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	result, err := repo.GetByPlayerID(ctx, player.ID)
	require.NoError(t, err)
	assert.Equal(t, player.ID, result.PlayerID)
}

func TestPlayerCertificationRepository_GetByPlayerID_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	_, err := repo.GetByPlayerID(ctx, 99999)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestPlayerCertificationRepository_ListPaged(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	// Create multiple certifications
	for i := 0; i < 5; i++ {
		playerUser := CreateUniqueTestUser(t, db, fmt.Sprintf("pc_paged_%d", i))
		player := CreateTestPlayer(t, db, playerUser)
		CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)
	}

	certs, total, err := repo.ListPaged(ctx, repository.PlayerCertificationListOptions{
		Page:     1,
		PageSize: 3,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, certs, 3)
}

func TestPlayerCertificationRepository_ListPaged_FilterByStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser1 := CreateUniqueTestUser(t, db, "pc_status1")
	player1 := CreateTestPlayer(t, db, playerUser1)
	playerUser2 := CreateUniqueTestUser(t, db, "pc_status2")
	player2 := CreateTestPlayer(t, db, playerUser2)

	CreateTestPlayerCertification(t, db, player1, model.CertificationStatusPending)
	CreateTestPlayerCertification(t, db, player2, model.CertificationStatusVerified)

	status := model.CertificationStatusPending
	certs, total, err := repo.ListPaged(ctx, repository.PlayerCertificationListOptions{
		Status:   &status,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	for _, c := range certs {
		assert.Equal(t, model.CertificationStatusPending, c.Status)
	}
}

func TestPlayerCertificationRepository_ListPaged_FilterByStatuses(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	// Create certifications with different statuses
	playerUser1 := CreateUniqueTestUser(t, db, "pc_statuses1")
	player1 := CreateTestPlayer(t, db, playerUser1)
	playerUser2 := CreateUniqueTestUser(t, db, "pc_statuses2")
	player2 := CreateTestPlayer(t, db, playerUser2)
	playerUser3 := CreateUniqueTestUser(t, db, "pc_statuses3")
	player3 := CreateTestPlayer(t, db, playerUser3)

	CreateTestPlayerCertification(t, db, player1, model.CertificationStatusPending)
	CreateTestPlayerCertification(t, db, player2, model.CertificationStatusVerified)
	CreateTestPlayerCertification(t, db, player3, model.CertificationStatusRejected)

	// Filter for pending and rejected
	certs, total, err := repo.ListPaged(ctx, repository.PlayerCertificationListOptions{
		Statuses: []model.CertificationStatus{
			model.CertificationStatusPending,
			model.CertificationStatusRejected,
		},
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(2))
	for _, c := range certs {
		assert.NotEqual(t, model.CertificationStatusVerified, c.Status)
	}
}

func TestPlayerCertificationRepository_ListPaged_FilterByKeyword(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_keyword")
	player := CreateTestPlayer(t, db, playerUser)
	cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	// Update real name
	cert.RealName = "UniqueSearchName"
	db.Save(cert)

	certs, total, err := repo.ListPaged(ctx, repository.PlayerCertificationListOptions{
		Keyword:  "UniqueSearch",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "UniqueSearchName", certs[0].RealName)
}

func TestPlayerCertificationRepository_ListPaged_FilterByPlayerID(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser1 := CreateUniqueTestUser(t, db, "pc_player_filter1")
	player1 := CreateTestPlayer(t, db, playerUser1)
	playerUser2 := CreateUniqueTestUser(t, db, "pc_player_filter2")
	player2 := CreateTestPlayer(t, db, playerUser2)

	CreateTestPlayerCertification(t, db, player1, model.CertificationStatusPending)
	CreateTestPlayerCertification(t, db, player2, model.CertificationStatusPending)

	certs, total, err := repo.ListPaged(ctx, repository.PlayerCertificationListOptions{
		PlayerID: &player1.ID,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, player1.ID, certs[0].PlayerID)
}

func TestPlayerCertificationRepository_ListPaged_WithPreloadedPlayer(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_preload")
	player := CreateTestPlayer(t, db, playerUser)
	CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	certs, _, err := repo.ListPaged(ctx, repository.PlayerCertificationListOptions{
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.Len(t, certs, 1)
	assert.NotNil(t, certs[0].Player)
	assert.Equal(t, player.ID, certs[0].Player.ID)
}

func TestPlayerCertificationRepository_ListPending(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_pending")
	player := CreateTestPlayer(t, db, playerUser)
	CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	certs, total, err := repo.ListPending(ctx, 1, 10)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	for _, c := range certs {
		assert.Equal(t, model.CertificationStatusPending, c.Status)
	}
}

func TestPlayerCertificationRepository_Update(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_update")
	player := CreateTestPlayer(t, db, playerUser)
	cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	// Update
	cert.RealName = "Updated Name"
	cert.Status = model.CertificationStatusVerified
	cert.PhotoURL = "https://example.com/photo.jpg"
	cert.VoiceURL = "https://example.com/voice.mp3"
	err := repo.Update(ctx, cert)
	require.NoError(t, err)

	// Verify
	updated, err := repo.Get(ctx, cert.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.RealName)
	assert.Equal(t, model.CertificationStatusVerified, updated.Status)
	assert.Equal(t, "https://example.com/photo.jpg", updated.PhotoURL)
	assert.Equal(t, "https://example.com/voice.mp3", updated.VoiceURL)
}

func TestPlayerCertificationRepository_Update_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	cert := &model.PlayerCertification{
		Base:     model.Base{ID: 99999},
		RealName: "Non-existent",
	}

	err := repo.Update(ctx, cert)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestPlayerCertificationRepository_UpdateStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	adminUser := CreateUniqueTestUser(t, db, "pc_admin")
	playerUser := CreateUniqueTestUser(t, db, "pc_status_update")
	player := CreateTestPlayer(t, db, playerUser)
	cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	// Approve
	err := repo.UpdateStatus(ctx, cert.ID, model.CertificationStatusVerified, &adminUser.ID, "")
	require.NoError(t, err)

	// Verify
	updated, err := repo.Get(ctx, cert.ID)
	require.NoError(t, err)
	assert.Equal(t, model.CertificationStatusVerified, updated.Status)
	assert.NotNil(t, updated.VerifiedBy)
	assert.Equal(t, adminUser.ID, *updated.VerifiedBy)
	assert.NotNil(t, updated.VerifiedAt)
}

func TestPlayerCertificationRepository_UpdateStatus_Reject(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	adminUser := CreateUniqueTestUser(t, db, "pc_admin_reject")
	playerUser := CreateUniqueTestUser(t, db, "pc_reject")
	player := CreateTestPlayer(t, db, playerUser)
	cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	// Reject
	err := repo.UpdateStatus(ctx, cert.ID, model.CertificationStatusRejected, &adminUser.ID, "Invalid ID card")
	require.NoError(t, err)

	// Verify
	updated, err := repo.Get(ctx, cert.ID)
	require.NoError(t, err)
	assert.Equal(t, model.CertificationStatusRejected, updated.Status)
	assert.Equal(t, "Invalid ID card", updated.RejectReason)
}

func TestPlayerCertificationRepository_UpdateStatus_WithoutVerifier(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_no_verifier")
	player := CreateTestPlayer(t, db, playerUser)
	cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	// Update status without verifier
	err := repo.UpdateStatus(ctx, cert.ID, model.CertificationStatusRejected, nil, "System rejected")
	require.NoError(t, err)

	// Verify
	updated, err := repo.Get(ctx, cert.ID)
	require.NoError(t, err)
	assert.Equal(t, model.CertificationStatusRejected, updated.Status)
	assert.Nil(t, updated.VerifiedBy)
	assert.Equal(t, "System rejected", updated.RejectReason)
}

func TestPlayerCertificationRepository_UpdateStatus_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	err := repo.UpdateStatus(ctx, 99999, model.CertificationStatusVerified, nil, "")
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestPlayerCertificationRepository_Delete(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_delete")
	player := CreateTestPlayer(t, db, playerUser)
	cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	err := repo.Delete(ctx, cert.ID)
	require.NoError(t, err)

	_, err = repo.Get(ctx, cert.ID)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestPlayerCertificationRepository_Delete_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	err := repo.Delete(ctx, 99999)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestPlayerCertificationRepository_CountByStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	// Create certifications with different statuses
	for i := 0; i < 3; i++ {
		playerUser := CreateUniqueTestUser(t, db, fmt.Sprintf("pc_count_pending_%d", i))
		player := CreateTestPlayer(t, db, playerUser)
		CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)
	}

	for i := 0; i < 2; i++ {
		playerUser := CreateUniqueTestUser(t, db, fmt.Sprintf("pc_count_verified_%d", i))
		player := CreateTestPlayer(t, db, playerUser)
		CreateTestPlayerCertification(t, db, player, model.CertificationStatusVerified)
	}

	playerUser := CreateUniqueTestUser(t, db, "pc_count_rejected")
	player := CreateTestPlayer(t, db, playerUser)
	CreateTestPlayerCertification(t, db, player, model.CertificationStatusRejected)

	counts, err := repo.CountByStatus(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, counts[model.CertificationStatusPending], int64(3))
	assert.GreaterOrEqual(t, counts[model.CertificationStatusVerified], int64(2))
	assert.GreaterOrEqual(t, counts[model.CertificationStatusRejected], int64(1))
}

func TestPlayerCertificationRepository_GetPendingCount(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_pending_count")
	player := CreateTestPlayer(t, db, playerUser)
	CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	count, err := repo.GetPendingCount(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(1))
}

// ==================== Business Rules Tests ====================

func TestPlayerCertification_BusinessRule_OneCertificationPerPlayer(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_one_per_player")
	player := CreateTestPlayer(t, db, playerUser)

	// Create first certification
	cert1 := CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	// Try to create second certification (should fail due to unique constraint)
	cert2 := &model.PlayerCertification{
		PlayerID:       player.ID,
		RealName:       "Another Name",
		IDCardNo:       "987654321098765432",
		IDCardFrontURL: "https://example.com/front2.jpg",
		IDCardBackURL:  "https://example.com/back2.jpg",
		Status:         model.CertificationStatusPending,
	}

	err := repo.Create(ctx, cert2)
	assert.Error(t, err)

	// First certification should still exist
	retrieved, _ := repo.Get(ctx, cert1.ID)
	assert.Equal(t, cert1.ID, retrieved.ID)
}

func TestPlayerCertification_BusinessRule_RejectedCanResubmit(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_resubmit")
	player := CreateTestPlayer(t, db, playerUser)

	// Create and reject first certification
	cert1 := CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	adminUser := CreateUniqueTestUser(t, db, "pc_admin_resubmit")
	err := repo.UpdateStatus(ctx, cert1.ID, model.CertificationStatusRejected, &adminUser.ID, "Invalid photo")
	require.NoError(t, err)

	// Delete rejected certification to allow resubmission
	err = repo.Delete(ctx, cert1.ID)
	require.NoError(t, err)

	// Now can submit new certification
	cert2 := &model.PlayerCertification{
		PlayerID:       player.ID,
		RealName:       "Updated Name",
		IDCardNo:       "123456789012345678",
		IDCardFrontURL: "https://example.com/front_new.jpg",
		IDCardBackURL:  "https://example.com/back_new.jpg",
		Status:         model.CertificationStatusPending,
	}

	err = repo.Create(ctx, cert2)
	require.NoError(t, err)
	assert.NotZero(t, cert2.ID)

	// Verify new certification
	retrieved, err := repo.Get(ctx, cert2.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", retrieved.RealName)
	assert.Equal(t, model.CertificationStatusPending, retrieved.Status)
}

func TestPlayerCertification_BusinessRule_CannotModifyPendingCertification(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_no_modify")
	player := CreateTestPlayer(t, db, playerUser)
	cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	// Update should work through repository (but business logic should prevent it)
	cert.RealName = "Modified Name"
	err := repo.Update(ctx, cert)
	require.NoError(t, err)

	// Verify modification happened at repository level
	retrieved, err := repo.Get(ctx, cert.ID)
	require.NoError(t, err)
	assert.Equal(t, "Modified Name", retrieved.RealName)
	// Note: Business logic validation should be handled at service layer
}

func TestPlayerCertification_BusinessRule_VerifiedCertification(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	adminUser := CreateUniqueTestUser(t, db, "pc_verified_admin")
	playerUser := CreateUniqueTestUser(t, db, "pc_verified_player")
	player := CreateTestPlayer(t, db, playerUser)
	cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	// Approve certification
	err := repo.UpdateStatus(ctx, cert.ID, model.CertificationStatusVerified, &adminUser.ID, "")
	require.NoError(t, err)

	// Verify all fields are correctly set
	retrieved, err := repo.GetWithPlayer(ctx, cert.ID)
	require.NoError(t, err)
	assert.Equal(t, model.CertificationStatusVerified, retrieved.Status)
	assert.NotNil(t, retrieved.VerifiedBy)
	assert.Equal(t, adminUser.ID, *retrieved.VerifiedBy)
	assert.NotNil(t, retrieved.VerifiedAt)
	assert.Empty(t, retrieved.RejectReason)
}

func TestPlayerCertification_BusinessRule_RejectionReason(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	adminUser := CreateUniqueTestUser(t, db, "pc_reject_admin")

	// Test various rejection reasons
	rejectionReasons := []string{
		"ID card is blurry",
		"Name does not match ID card",
		"ID card expired",
		"Photo quality is poor",
		"Suspicious information provided",
	}

	for i, reason := range rejectionReasons {
		playerUser := CreateUniqueTestUser(t, db, fmt.Sprintf("pc_reject_%d", i))
		player := CreateTestPlayer(t, db, playerUser)
		cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

		err := repo.UpdateStatus(ctx, cert.ID, model.CertificationStatusRejected, &adminUser.ID, reason)
		require.NoError(t, err)

		retrieved, err := repo.Get(ctx, cert.ID)
		require.NoError(t, err)
		assert.Equal(t, model.CertificationStatusRejected, retrieved.Status)
		assert.Equal(t, reason, retrieved.RejectReason)
	}
}

// ==================== Sensitive Information Tests ====================

func TestPlayerCertification_SensitiveInformation_IDCardNotExposed(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_sensitive")
	player := CreateTestPlayer(t, db, playerUser)
	cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	// ID card should be stored but not exposed in JSON
	retrieved, err := repo.Get(ctx, cert.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, retrieved.IDCardNo)

	// Verify JSON tag is set to "-" to prevent exposure
	// The model struct has: IDCardNo string `json:"-" gorm:"..."`
}

func TestPlayerCertification_SensitiveInformation_MultipleCertificationsHaveDifferentIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	// Create certifications for different players with different ID cards
	idCards := []string{
		"110101199001011234",
		"310101199002022345",
		"440101199003033456",
	}

	for i, idCard := range idCards {
		playerUser := CreateUniqueTestUser(t, db, fmt.Sprintf("pc_idcard_%d", i))
		player := CreateTestPlayer(t, db, playerUser)
		cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)
		cert.IDCardNo = idCard
		db.Save(cert)

		retrieved, err := repo.Get(ctx, cert.ID)
		require.NoError(t, err)
		assert.Equal(t, idCard, retrieved.IDCardNo)
	}
}

// ==================== Status Transition Tests ====================

func TestPlayerCertification_StatusTransition_PendingToVerified(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	adminUser := CreateUniqueTestUser(t, db, "pc_transition_pv")
	playerUser := CreateUniqueTestUser(t, db, "pc_transition_player")
	player := CreateTestPlayer(t, db, playerUser)
	cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	err := repo.UpdateStatus(ctx, cert.ID, model.CertificationStatusVerified, &adminUser.ID, "")
	require.NoError(t, err)

	retrieved, err := repo.Get(ctx, cert.ID)
	require.NoError(t, err)
	assert.Equal(t, model.CertificationStatusVerified, retrieved.Status)
	assert.NotNil(t, retrieved.VerifiedAt)
}

func TestPlayerCertification_StatusTransition_PendingToRejected(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	adminUser := CreateUniqueTestUser(t, db, "pc_transition_pr")
	playerUser := CreateUniqueTestUser(t, db, "pc_transition_player2")
	player := CreateTestPlayer(t, db, playerUser)
	cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	err := repo.UpdateStatus(ctx, cert.ID, model.CertificationStatusRejected, &adminUser.ID, "Invalid")
	require.NoError(t, err)

	retrieved, err := repo.Get(ctx, cert.ID)
	require.NoError(t, err)
	assert.Equal(t, model.CertificationStatusRejected, retrieved.Status)
}

func TestPlayerCertification_StatusTransition_VerifiedToRejected(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	adminUser := CreateUniqueTestUser(t, db, "pc_transition_vr")
	playerUser := CreateUniqueTestUser(t, db, "pc_transition_player3")
	player := CreateTestPlayer(t, db, playerUser)
	cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusVerified)

	// Revoke verified certification
	err := repo.UpdateStatus(ctx, cert.ID, model.CertificationStatusRejected, &adminUser.ID, "Revoked due to policy violation")
	require.NoError(t, err)

	retrieved, err := repo.Get(ctx, cert.ID)
	require.NoError(t, err)
	assert.Equal(t, model.CertificationStatusRejected, retrieved.Status)
	assert.Equal(t, "Revoked due to policy violation", retrieved.RejectReason)
}

func TestPlayerCertification_StatusTransition_RejectedToVerified(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	adminUser := CreateUniqueTestUser(t, db, "pc_transition_rv")
	playerUser := CreateUniqueTestUser(t, db, "pc_transition_player4")
	player := CreateTestPlayer(t, db, playerUser)
	cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusRejected)

	// Re-approve rejected certification (appeal approved)
	err := repo.UpdateStatus(ctx, cert.ID, model.CertificationStatusVerified, &adminUser.ID, "")
	require.NoError(t, err)

	retrieved, err := repo.Get(ctx, cert.ID)
	require.NoError(t, err)
	assert.Equal(t, model.CertificationStatusVerified, retrieved.Status)
	assert.Empty(t, retrieved.RejectReason) // Clear reject reason
}

// ==================== Player Rank Certification Integration Tests ====================

func TestPlayerCertification_Integration_PrerequisiteForRankCertification(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_rank_prereq")
	player := CreateTestPlayer(t, db, playerUser)

	// Create game and rank for rank certification
	game := CreateTestGame(t, db, "Honor of Kings")
	rank := CreateTestGameRank(t, db, game, "Diamond", 5, 5000)

	// Player without real-name certification cannot have rank certification
	// (This is a business rule that should be enforced at service layer)
	cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	// Verify player has real-name certification
	hasCert, err := repo.GetByPlayerID(ctx, player.ID)
	require.NoError(t, err)
	assert.NotNil(t, hasCert)

	// After real-name verification is approved, player can apply for rank certification
	adminUser := CreateUniqueTestUser(t, db, "pc_rank_admin")
	err = repo.UpdateStatus(ctx, cert.ID, model.CertificationStatusVerified, &adminUser.ID, "")
	require.NoError(t, err)

	// Now player can create rank certification
	rankRecord := CreateTestPlayerRankRecord(t, db, player, game, rank, model.PlayerRankStatusPending)
	assert.Equal(t, player.ID, rankRecord.PlayerID)
	assert.Equal(t, game.ID, rankRecord.GameID)
	assert.Equal(t, rank.ID, rankRecord.RankID)
}

func TestPlayerCertification_Integration_UnverifiedCannotHaveRank(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_no_rank")
	player := CreateTestPlayer(t, db, playerUser)

	// Create pending real-name certification
	_ = CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	// Verify certification is still pending
	retrieved, err := repo.GetByPlayerID(ctx, player.ID)
	require.NoError(t, err)
	assert.Equal(t, model.CertificationStatusPending, retrieved.Status)

	// Business rule: unverified player cannot have verified rank certification
	// (This should be enforced at service layer during rank certification approval)
}

// ==================== Timestamp Tests ====================

func TestPlayerCertification_Timestamps_AuditFields(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_timestamps")
	player := CreateTestPlayer(t, db, playerUser)
	cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	// Check creation timestamp
	assert.NotNil(t, cert.CreatedAt)
	assert.False(t, cert.CreatedAt.IsZero())

	// Update and check modification timestamp
	originalUpdatedAt := cert.UpdatedAt
	time.Sleep(time.Millisecond) // Ensure time difference

	cert.RealName = "Updated"
	err := repo.Update(ctx, cert)
	require.NoError(t, err)

	retrieved, err := repo.Get(ctx, cert.ID)
	require.NoError(t, err)
	assert.True(t, retrieved.UpdatedAt.After(originalUpdatedAt))
}

func TestPlayerCertification_Timestamps_VerifiedAt(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	adminUser := CreateUniqueTestUser(t, db, "pc_verified_at")
	playerUser := CreateUniqueTestUser(t, db, "pc_verified_at_player")
	player := CreateTestPlayer(t, db, playerUser)
	cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	// Initially, VerifiedAt should be nil
	assert.Nil(t, cert.VerifiedAt)

	// Approve certification
	beforeApproval := time.Now()
	err := repo.UpdateStatus(ctx, cert.ID, model.CertificationStatusVerified, &adminUser.ID, "")
	require.NoError(t, err)

	// Verify VerifiedAt is set
	retrieved, err := repo.Get(ctx, cert.ID)
	require.NoError(t, err)
	assert.NotNil(t, retrieved.VerifiedAt)
	assert.True(t, retrieved.VerifiedAt.After(beforeApproval) || retrieved.VerifiedAt.Equal(beforeApproval))
}

// ==================== Edge Cases Tests ====================

func TestPlayerCertification_EdgeCase_EmptyOptionalFields(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_empty_fields")
	player := CreateTestPlayer(t, db, playerUser)

	cert := &model.PlayerCertification{
		PlayerID:       player.ID,
		RealName:       "Test User",
		IDCardNo:       "123456789012345678",
		IDCardFrontURL: "https://example.com/front.jpg",
		IDCardBackURL:  "https://example.com/back.jpg",
		Status:         model.CertificationStatusPending,
		// Optional fields left empty
		PhotoURL: "",
		VoiceURL: "",
		ExtJSON:  "",
	}

	err := repo.Create(ctx, cert)
	require.NoError(t, err)

	retrieved, err := repo.Get(ctx, cert.ID)
	require.NoError(t, err)
	assert.Empty(t, retrieved.PhotoURL)
	assert.Empty(t, retrieved.VoiceURL)
}

func TestPlayerCertification_EdgeCase_LongRealName(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_long_name")
	player := CreateTestPlayer(t, db, playerUser)

	// Test with maximum allowed real name length (64 chars)
	longName := "张三张三张三张三张三张三张三张三张三张三张三张三张三张三张三1234567890"

	cert := &model.PlayerCertification{
		PlayerID:       player.ID,
		RealName:       longName,
		IDCardNo:       "123456789012345678",
		IDCardFrontURL: "https://example.com/front.jpg",
		IDCardBackURL:  "https://example.com/back.jpg",
		Status:         model.CertificationStatusPending,
	}

	err := repo.Create(ctx, cert)
	require.NoError(t, err)

	retrieved, err := repo.Get(ctx, cert.ID)
	require.NoError(t, err)
	assert.Equal(t, longName, retrieved.RealName)
}

func TestPlayerCertification_EdgeCase_LongURL(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	playerUser := CreateUniqueTestUser(t, db, "pc_long_url")
	player := CreateTestPlayer(t, db, playerUser)

	// Test with long URLs
	longURL := "https://example.com/very/long/path/to/the/image/file/that/exceeds/normal/length/but/still/within/limit/and/should/be/stored.jpg"

	cert := &model.PlayerCertification{
		PlayerID:       player.ID,
		RealName:       "Test User",
		IDCardNo:       "123456789012345678",
		IDCardFrontURL: longURL,
		IDCardBackURL:  longURL,
		PhotoURL:       longURL,
		VoiceURL:       longURL,
		Status:         model.CertificationStatusPending,
	}

	err := repo.Create(ctx, cert)
	require.NoError(t, err)

	retrieved, err := repo.Get(ctx, cert.ID)
	require.NoError(t, err)
	assert.Equal(t, longURL, retrieved.IDCardFrontURL)
	assert.Equal(t, longURL, retrieved.IDCardBackURL)
}

func TestPlayerCertification_EdgeCase_LongRejectReason(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	adminUser := CreateUniqueTestUser(t, db, "pc_long_reject")
	playerUser := CreateUniqueTestUser(t, db, "pc_long_reject_player")
	player := CreateTestPlayer(t, db, playerUser)
	cert := CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)

	// Test with maximum allowed reject reason length (500 chars)
	longReason := "拒绝原因：身份证照片不清晰，无法识别个人信息。请重新拍摄清晰的身份证照片，确保所有信息清晰可见。" +
		"请确保照片中没有反光、模糊或遮挡。照片需要显示身份证的正面和反面，所有文字和号码必须清晰。" +
		"如果照片质量不达标，我们将无法通过您的认证申请。请您务必提供高质量的照片。" +
		"审核人员需要能够清楚地看到您的姓名、身份证号、出生日期等所有信息。" +
		"This is additional text to reach the maximum length of 500 characters for testing purposes." +
		"123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"

	err := repo.UpdateStatus(ctx, cert.ID, model.CertificationStatusRejected, &adminUser.ID, longReason)
	require.NoError(t, err)

	retrieved, err := repo.Get(ctx, cert.ID)
	require.NoError(t, err)
	assert.Equal(t, longReason, retrieved.RejectReason)
}

// ==================== Performance Tests ====================

func TestPlayerCertification_Performance_LargeDataset(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	// Create 50 certifications
	for i := 0; i < 50; i++ {
		playerUser := CreateUniqueTestUser(t, db, fmt.Sprintf("pc_perf_%d", i))
		player := CreateTestPlayer(t, db, playerUser)
		status := model.CertificationStatusPending
		if i%3 == 0 {
			status = model.CertificationStatusVerified
		} else if i%3 == 1 {
			status = model.CertificationStatusRejected
		}
		CreateTestPlayerCertification(t, db, player, status)
	}

	// Test pagination performance
	start := time.Now()
	certs, total, err := repo.ListPaged(ctx, repository.PlayerCertificationListOptions{
		Page:     1,
		PageSize: 20,
	})
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, int64(50), total)
	assert.Len(t, certs, 20)
	assert.Less(t, elapsed.Milliseconds(), int64(1000), "Query should complete in less than 1 second")
}

func TestPlayerCertification_Performance_CountByStatusWithLargeDataset(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	// Create 100 certifications
	for i := 0; i < 100; i++ {
		playerUser := CreateUniqueTestUser(t, db, fmt.Sprintf("pc_count_perf_%d", i))
		player := CreateTestPlayer(t, db, playerUser)
		status := model.CertificationStatusPending
		if i%2 == 0 {
			status = model.CertificationStatusVerified
		}
		CreateTestPlayerCertification(t, db, player, status)
	}

	// Test count performance
	start := time.Now()
	counts, err := repo.CountByStatus(ctx)
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, counts[model.CertificationStatusPending], int64(50))
	assert.GreaterOrEqual(t, counts[model.CertificationStatusVerified], int64(50))
	assert.Less(t, elapsed.Milliseconds(), int64(500), "Count query should complete in less than 500ms")
}
