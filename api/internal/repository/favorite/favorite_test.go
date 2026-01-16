package favorite

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service/integration"
)

func setupTestDB(t *testing.T) *Repository {
	t.Helper()
	db := integration.SetupTestDB(t)
	return NewRepository(db)
}

func createTestUser(t *testing.T, repo *Repository, suffix string) *model.User {
	t.Helper()
	openID := "openid_" + suffix
	user := &model.User{
		Phone:         "138" + suffix,
		Name:          "Test User " + suffix,
		Email:         "user_" + suffix + "@test.com",
		Role:          model.RoleUser,
		Status:        model.UserStatusActive,
		WeChatOpenID:  &openID,
		WeChatUnionID: "unionid_" + suffix,
	}
	err := repo.db.Create(user).Error
	require.NoError(t, err)
	return user
}

func createTestPlayer(t *testing.T, repo *Repository, userID uint64, suffix string) *model.Player {
	t.Helper()
	player := &model.Player{
		UserID:             userID,
		Nickname:           "Player " + suffix,
		Bio:                "Test bio",
		HourlyRateCents:    5000,
		VerificationStatus: model.VerificationVerified,
	}
	err := repo.db.Create(player).Error
	require.NoError(t, err)
	return player
}

// ============================================================================
// Basic CRUD Tests
// ============================================================================

func TestFavoriteRepository_Create(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	user := createTestUser(t, repo, "001")
	playerUser := createTestUser(t, repo, "002")
	player := createTestPlayer(t, repo, playerUser.ID, "001")

	fav := &model.Favorite{
		UserID:   user.ID,
		PlayerID: player.ID,
	}

	err := repo.Create(ctx, fav)
	assert.NoError(t, err)
	assert.NotZero(t, fav.ID)
}

func TestFavoriteRepository_Create_Duplicate(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	user := createTestUser(t, repo, "003")
	playerUser := createTestUser(t, repo, "004")
	player := createTestPlayer(t, repo, playerUser.ID, "002")

	fav1 := &model.Favorite{UserID: user.ID, PlayerID: player.ID}
	err := repo.Create(ctx, fav1)
	assert.NoError(t, err)

	// Try to create duplicate
	fav2 := &model.Favorite{UserID: user.ID, PlayerID: player.ID}
	err = repo.Create(ctx, fav2)
	assert.Error(t, err) // Should fail due to unique constraint
}

func TestFavoriteRepository_Delete(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	user := createTestUser(t, repo, "005")
	playerUser := createTestUser(t, repo, "006")
	player := createTestPlayer(t, repo, playerUser.ID, "003")

	fav := &model.Favorite{UserID: user.ID, PlayerID: player.ID}
	err := repo.Create(ctx, fav)
	require.NoError(t, err)

	err = repo.Delete(ctx, user.ID, player.ID)
	assert.NoError(t, err)

	// Verify deleted
	exists, err := repo.Exists(ctx, user.ID, player.ID)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestFavoriteRepository_Delete_NotFound(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	err := repo.Delete(ctx, 999999, 999999)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestFavoriteRepository_Exists(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	user := createTestUser(t, repo, "007")
	playerUser := createTestUser(t, repo, "008")
	player := createTestPlayer(t, repo, playerUser.ID, "004")

	// Not exists initially
	exists, err := repo.Exists(ctx, user.ID, player.ID)
	assert.NoError(t, err)
	assert.False(t, exists)

	// Create favorite
	fav := &model.Favorite{UserID: user.ID, PlayerID: player.ID}
	err = repo.Create(ctx, fav)
	require.NoError(t, err)

	// Now exists
	exists, err = repo.Exists(ctx, user.ID, player.ID)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestFavoriteRepository_ListByUserID(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	user := createTestUser(t, repo, "009")

	// Create multiple players and favorites
	for i := 0; i < 5; i++ {
		suffix := "01" + string(rune('0'+i))
		playerUser := createTestUser(t, repo, suffix)
		player := createTestPlayer(t, repo, playerUser.ID, suffix)
		fav := &model.Favorite{UserID: user.ID, PlayerID: player.ID}
		err := repo.Create(ctx, fav)
		require.NoError(t, err)
	}

	// List favorites
	favorites, total, err := repo.ListByUserID(ctx, user.ID, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, favorites, 5)
}

func TestFavoriteRepository_ListByUserID_Pagination(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	user := createTestUser(t, repo, "020")

	// Create 15 favorites
	for i := 0; i < 15; i++ {
		suffix := "02" + string(rune('0'+i/10)) + string(rune('0'+i%10))
		playerUser := createTestUser(t, repo, suffix)
		player := createTestPlayer(t, repo, playerUser.ID, suffix)
		fav := &model.Favorite{UserID: user.ID, PlayerID: player.ID}
		err := repo.Create(ctx, fav)
		require.NoError(t, err)
	}

	// Page 1
	favorites, total, err := repo.ListByUserID(ctx, user.ID, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(15), total)
	assert.Len(t, favorites, 10)

	// Page 2
	favorites, total, err = repo.ListByUserID(ctx, user.ID, 2, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(15), total)
	assert.Len(t, favorites, 5)
}

// ============================================================================
// Concurrent Tests
// ============================================================================

func TestFavoriteRepository_ConcurrentCreate(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	user := createTestUser(t, repo, "100")

	// Create multiple players with unique suffixes
	var players []*model.Player
	for i := 0; i < 10; i++ {
		suffix := "1" + string(rune('a'+i)) + "0"
		playerUser := createTestUser(t, repo, suffix)
		player := createTestPlayer(t, repo, playerUser.ID, suffix)
		players = append(players, player)
	}

	// Concurrent create favorites
	var wg sync.WaitGroup
	errors := make(chan error, len(players))

	for _, player := range players {
		wg.Add(1)
		go func(p *model.Player) {
			defer wg.Done()
			fav := &model.Favorite{UserID: user.ID, PlayerID: p.ID}
			if err := repo.Create(ctx, fav); err != nil {
				errors <- err
			}
		}(player)
	}

	wg.Wait()
	close(errors)

	// Check no errors
	for err := range errors {
		t.Errorf("Concurrent create error: %v", err)
	}

	// Verify all favorites created
	favorites, total, err := repo.ListByUserID(ctx, user.ID, 1, 100)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), total)
	assert.Len(t, favorites, 10)
}

func TestFavoriteRepository_ConcurrentCreateSame(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	user := createTestUser(t, repo, "200")
	playerUser := createTestUser(t, repo, "201")
	player := createTestPlayer(t, repo, playerUser.ID, "200")

	// Concurrent create same favorite (should only succeed once)
	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fav := &model.Favorite{UserID: user.ID, PlayerID: player.ID}
			if err := repo.Create(ctx, fav); err == nil {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Only one should succeed due to unique constraint
	assert.Equal(t, 1, successCount)

	// Verify only one favorite exists
	exists, err := repo.Exists(ctx, user.ID, player.ID)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestFavoriteRepository_ConcurrentDeleteAndCreate(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	user := createTestUser(t, repo, "300")
	playerUser := createTestUser(t, repo, "301")
	player := createTestPlayer(t, repo, playerUser.ID, "300")

	// Create initial favorite
	fav := &model.Favorite{UserID: user.ID, PlayerID: player.ID}
	err := repo.Create(ctx, fav)
	require.NoError(t, err)

	// Concurrent delete and create
	var wg sync.WaitGroup

	// Delete goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = repo.Delete(ctx, user.ID, player.ID)
	}()

	// Create goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		newFav := &model.Favorite{UserID: user.ID, PlayerID: player.ID}
		_ = repo.Create(ctx, newFav)
	}()

	wg.Wait()

	// Final state should be consistent (either exists or not)
	exists, err := repo.Exists(ctx, user.ID, player.ID)
	assert.NoError(t, err)
	// Result can be either true or false, but should not error
	t.Logf("Final exists state: %v", exists)
}
