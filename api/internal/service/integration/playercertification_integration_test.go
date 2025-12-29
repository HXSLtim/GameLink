// Package integration provides integration tests for services.
package integration

import (
	"context"
	"testing"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/playercertification"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
		playerUser := CreateUniqueTestUser(t, db, "pc_paged")
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
	err := repo.Update(ctx, cert)
	require.NoError(t, err)

	// Verify
	updated, err := repo.Get(ctx, cert.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.RealName)
	assert.Equal(t, model.CertificationStatusVerified, updated.Status)
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

func TestPlayerCertificationRepository_CountByStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := playercertification.NewPlayerCertificationRepository(db)

	// Create certifications with different statuses
	for i := 0; i < 3; i++ {
		playerUser := CreateUniqueTestUser(t, db, "pc_count_pending")
		player := CreateTestPlayer(t, db, playerUser)
		CreateTestPlayerCertification(t, db, player, model.CertificationStatusPending)
	}

	for i := 0; i < 2; i++ {
		playerUser := CreateUniqueTestUser(t, db, "pc_count_verified")
		player := CreateTestPlayer(t, db, playerUser)
		CreateTestPlayerCertification(t, db, player, model.CertificationStatusVerified)
	}

	counts, err := repo.CountByStatus(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, counts[model.CertificationStatusPending], int64(3))
	assert.GreaterOrEqual(t, counts[model.CertificationStatusVerified], int64(2))
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
