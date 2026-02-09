package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// ============================================================================
// Mock Implementations
// ============================================================================

// MockUserTagRepository is a mock implementation of UserTagRepository
type MockUserTagRepository struct {
	mock.Mock
}

func (m *MockUserTagRepository) CreateTag(ctx context.Context, tag *model.UserTag) error {
	args := m.Called(ctx, tag)
	if args.Error(0) == nil && tag.ID == 0 {
		tag.ID = 1
	}
	return args.Error(0)
}

func (m *MockUserTagRepository) GetTag(ctx context.Context, id uint64) (*model.UserTag, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserTag), args.Error(1)
}

func (m *MockUserTagRepository) ListTags(ctx context.Context) ([]model.UserTag, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.UserTag), args.Error(1)
}

func (m *MockUserTagRepository) UpdateTag(ctx context.Context, tag *model.UserTag) error {
	args := m.Called(ctx, tag)
	return args.Error(0)
}

func (m *MockUserTagRepository) DeleteTag(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserTagRepository) AddTagToUser(ctx context.Context, userID uint64, tagID uint64) error {
	args := m.Called(ctx, userID, tagID)
	return args.Error(0)
}

func (m *MockUserTagRepository) RemoveTagFromUser(ctx context.Context, userID uint64, tagID uint64) error {
	args := m.Called(ctx, userID, tagID)
	return args.Error(0)
}

func (m *MockUserTagRepository) GetUserTags(ctx context.Context, userID uint64) ([]model.UserTag, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.UserTag), args.Error(1)
}

func (m *MockUserTagRepository) BatchSetUserTags(ctx context.Context, userID uint64, tagIDs []uint64) error {
	args := m.Called(ctx, userID, tagIDs)
	return args.Error(0)
}

func (m *MockUserTagRepository) GetUsersByTag(ctx context.Context, tagID uint64, page, pageSize int) ([]model.User, int64, error) {
	args := m.Called(ctx, tagID, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

// MockTagUserRepository is a mock implementation of UserRepository for tag tests
type MockTagUserRepository struct {
	mock.Mock
}

func (m *MockTagUserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockTagUserRepository) List(ctx context.Context) ([]model.User, error) { return nil, nil }
func (m *MockTagUserRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (m *MockTagUserRepository) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (m *MockTagUserRepository) Count(ctx context.Context, opts repository.UserListOptions) (int, error) {
	return 0, nil
}
func (m *MockTagUserRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.User, error) {
	return nil, nil
}
func (m *MockTagUserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, nil
}
func (m *MockTagUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}
func (m *MockTagUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}
func (m *MockTagUserRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, nil
}
func (m *MockTagUserRepository) Create(ctx context.Context, user *model.User) error { return nil }
func (m *MockTagUserRepository) Update(ctx context.Context, user *model.User) error { return nil }
func (m *MockTagUserRepository) Delete(ctx context.Context, id uint64) error        { return nil }
func (m *MockTagUserRepository) UpdatePassword(ctx context.Context, userID uint64, newPassword string) error {
	return nil
}
func (m *MockTagUserRepository) GetByWeChatOpenID(ctx context.Context, openID string) (*model.User, error) {
	return nil, nil
}
func (m *MockTagUserRepository) GetByWeChatUnionID(ctx context.Context, unionID string) (*model.User, error) {
	return nil, nil
}

// MockTagCache is a mock implementation of cache.Cache
type MockTagCache struct {
	mock.Mock
	data map[string]string
}

func NewMockTagCache() *MockTagCache {
	return &MockTagCache{data: make(map[string]string)}
}

func (m *MockTagCache) Get(ctx context.Context, key string) (string, bool, error) {
	if val, ok := m.data[key]; ok {
		return val, true, nil
	}
	return "", false, nil
}

func (m *MockTagCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	m.data[key] = value
	return nil
}

func (m *MockTagCache) Delete(ctx context.Context, key string) error {
	delete(m.data, key)
	return nil
}

func (m *MockTagCache) Close(ctx context.Context) error {
	return nil
}

func (m *MockTagCache) GetRedisClient() interface{} {
	return nil
}

// ============================================================================
// Tests
// ============================================================================

func TestNewUserTagService(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	svc := NewUserTagService(tagRepo, userRepo, cache)

	assert.NotNil(t, svc)
	assert.Equal(t, tagRepo, svc.tagRepo)
	assert.Equal(t, userRepo, svc.userRepo)
	assert.Equal(t, cache, svc.cache)
}

func TestUserTagService_CreateTag(t *testing.T) {
	tests := []struct {
		name        string
		tagName     string
		color       string
		description string
		setupMock   func(*MockUserTagRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid tag creation",
			tagName:     "VIP",
			color:       "#FF6B6B",
			description: "VIP用户",
			setupMock: func(m *MockUserTagRepository) {
				m.On("ListTags", mock.Anything).Return([]model.UserTag{}, nil)
				m.On("CreateTag", mock.Anything, mock.AnythingOfType("*model.UserTag")).Return(nil)
			},
			expectError: false,
		},
		{
			name:        "empty tag name",
			tagName:     "",
			color:       "#FF6B6B",
			description: "Test",
			setupMock:   func(m *MockUserTagRepository) {},
			expectError: true,
			errorMsg:    "标签名称长度必须在1-64个字符之间",
		},
		{
			name:        "tag name too long",
			tagName:     string(make([]byte, 65)),
			color:       "#FF6B6B",
			description: "Test",
			setupMock:   func(m *MockUserTagRepository) {},
			expectError: true,
			errorMsg:    "标签名称长度必须在1-64个字符之间",
		},
		{
			name:        "invalid color format",
			tagName:     "Test",
			color:       "red",
			description: "Test",
			setupMock:   func(m *MockUserTagRepository) {},
			expectError: true,
			errorMsg:    "颜色格式不正确",
		},
		{
			name:        "duplicate tag name",
			tagName:     "VIP",
			color:       "#FF6B6B",
			description: "VIP用户",
			setupMock: func(m *MockUserTagRepository) {
				m.On("ListTags", mock.Anything).Return([]model.UserTag{
					{Name: "VIP"},
				}, nil)
			},
			expectError: true,
			errorMsg:    "标签名称已存在",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagRepo := &MockUserTagRepository{}
			userRepo := &MockTagUserRepository{}
			cache := NewMockTagCache()

			tt.setupMock(tagRepo)

			svc := NewUserTagService(tagRepo, userRepo, cache)
			tag, err := svc.CreateTag(context.Background(), tt.tagName, tt.color, tt.description)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Nil(t, tag)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tag)
				assert.Equal(t, tt.tagName, tag.Name)
			}
		})
	}
}

func TestUserTagService_GetTag(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	expectedTag := &model.UserTag{Name: "VIP", Color: "#FF6B6B"}
	expectedTag.ID = 1

	tagRepo.On("GetTag", mock.Anything, uint64(1)).Return(expectedTag, nil)
	tagRepo.On("GetTag", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)

	svc := NewUserTagService(tagRepo, userRepo, cache)

	// Test existing tag
	tag, err := svc.GetTag(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedTag, tag)

	// Test non-existing tag
	tag, err = svc.GetTag(context.Background(), 999)
	assert.Error(t, err)
	assert.Nil(t, tag)
}

func TestUserTagService_ListTags(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	expectedTags := []model.UserTag{
		{Name: "VIP", Color: "#FF6B6B"},
		{Name: "New", Color: "#00FF00"},
	}

	tagRepo.On("ListTags", mock.Anything).Return(expectedTags, nil)

	svc := NewUserTagService(tagRepo, userRepo, cache)

	// First call - should hit database
	tags, err := svc.ListTags(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedTags, tags)

	// Second call - should hit cache
	tags, err = svc.ListTags(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedTags, tags)
}

func TestUserTagService_UpdateTag(t *testing.T) {
	tests := []struct {
		name        string
		tagID       uint64
		tagName     string
		color       string
		description string
		setupMock   func(*MockUserTagRepository)
		expectError bool
	}{
		{
			name:        "valid update",
			tagID:       1,
			tagName:     "Updated",
			color:       "#00FF00",
			description: "Updated description",
			setupMock: func(m *MockUserTagRepository) {
				tag := &model.UserTag{Name: "VIP", Color: "#FF6B6B"}
				tag.ID = 1
				m.On("GetTag", mock.Anything, uint64(1)).Return(tag, nil)
				m.On("UpdateTag", mock.Anything, mock.AnythingOfType("*model.UserTag")).Return(nil)
			},
			expectError: false,
		},
		{
			name:        "tag not found",
			tagID:       999,
			tagName:     "Updated",
			color:       "#00FF00",
			description: "Updated description",
			setupMock: func(m *MockUserTagRepository) {
				m.On("GetTag", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
		{
			name:        "invalid color",
			tagID:       1,
			tagName:     "Updated",
			color:       "invalid",
			description: "Updated description",
			setupMock:   func(m *MockUserTagRepository) {},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagRepo := &MockUserTagRepository{}
			userRepo := &MockTagUserRepository{}
			cache := NewMockTagCache()

			tt.setupMock(tagRepo)

			svc := NewUserTagService(tagRepo, userRepo, cache)
			err := svc.UpdateTag(context.Background(), tt.tagID, tt.tagName, tt.color, tt.description)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserTagService_DeleteTag(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	tag := &model.UserTag{Name: "VIP"}
	tag.ID = 1

	tagRepo.On("GetTag", mock.Anything, uint64(1)).Return(tag, nil)
	tagRepo.On("DeleteTag", mock.Anything, uint64(1)).Return(nil)
	tagRepo.On("GetTag", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)

	svc := NewUserTagService(tagRepo, userRepo, cache)

	// Test successful delete
	err := svc.DeleteTag(context.Background(), 1)
	assert.NoError(t, err)

	// Test delete non-existing tag
	err = svc.DeleteTag(context.Background(), 999)
	assert.Error(t, err)
}

func TestUserTagService_AddTagToUser(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint64
		tagID       uint64
		setupMock   func(*MockUserTagRepository, *MockTagUserRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:   "successful add tag to user",
			userID: 1,
			tagID:  1,
			setupMock: func(tagRepo *MockUserTagRepository, userRepo *MockTagUserRepository) {
				user := &model.User{Name: "Test User"}
				user.ID = 1
				userRepo.On("Get", mock.Anything, uint64(1)).Return(user, nil)

				tag := &model.UserTag{Name: "VIP"}
				tag.ID = 1
				tagRepo.On("GetTag", mock.Anything, uint64(1)).Return(tag, nil)
				tagRepo.On("GetUserTags", mock.Anything, uint64(1)).Return([]model.UserTag{}, nil)
				tagRepo.On("AddTagToUser", mock.Anything, uint64(1), uint64(1)).Return(nil)
			},
			expectError: false,
		},
		{
			name:   "user not found",
			userID: 999,
			tagID:  1,
			setupMock: func(tagRepo *MockUserTagRepository, userRepo *MockTagUserRepository) {
				userRepo.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
			errorMsg:    "用户不存在",
		},
		{
			name:   "tag not found",
			userID: 1,
			tagID:  999,
			setupMock: func(tagRepo *MockUserTagRepository, userRepo *MockTagUserRepository) {
				user := &model.User{Name: "Test User"}
				user.ID = 1
				userRepo.On("Get", mock.Anything, uint64(1)).Return(user, nil)
				tagRepo.On("GetTag", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
			errorMsg:    "标签不存在",
		},
		{
			name:   "user already has tag",
			userID: 1,
			tagID:  1,
			setupMock: func(tagRepo *MockUserTagRepository, userRepo *MockTagUserRepository) {
				user := &model.User{Name: "Test User"}
				user.ID = 1
				userRepo.On("Get", mock.Anything, uint64(1)).Return(user, nil)

				tag := &model.UserTag{Name: "VIP"}
				tag.ID = 1
				tagRepo.On("GetTag", mock.Anything, uint64(1)).Return(tag, nil)
				tagRepo.On("GetUserTags", mock.Anything, uint64(1)).Return([]model.UserTag{*tag}, nil)
			},
			expectError: true,
			errorMsg:    "用户已拥有此标签",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagRepo := &MockUserTagRepository{}
			userRepo := &MockTagUserRepository{}
			cache := NewMockTagCache()

			tt.setupMock(tagRepo, userRepo)

			svc := NewUserTagService(tagRepo, userRepo, cache)
			err := svc.AddTagToUser(context.Background(), tt.userID, tt.tagID)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserTagService_RemoveTagFromUser(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint64
		tagID       uint64
		setupMock   func(*MockUserTagRepository)
		expectError bool
	}{
		{
			name:   "successful remove",
			userID: 1,
			tagID:  1,
			setupMock: func(m *MockUserTagRepository) {
				m.On("RemoveTagFromUser", mock.Anything, uint64(1), uint64(1)).Return(nil)
			},
			expectError: false,
		},
		{
			name:   "repository error",
			userID: 1,
			tagID:  999,
			setupMock: func(m *MockUserTagRepository) {
				m.On("RemoveTagFromUser", mock.Anything, uint64(1), uint64(999)).Return(errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagRepo := &MockUserTagRepository{}
			userRepo := &MockTagUserRepository{}
			cache := NewMockTagCache()

			tt.setupMock(tagRepo)

			svc := NewUserTagService(tagRepo, userRepo, cache)
			err := svc.RemoveTagFromUser(context.Background(), tt.userID, tt.tagID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserTagService_GetUserTags(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint64
		setupMock   func(*MockUserTagRepository, *MockTagUserRepository)
		expectError bool
		expectTags  int
	}{
		{
			name:   "successful get user tags",
			userID: 1,
			setupMock: func(tagRepo *MockUserTagRepository, userRepo *MockTagUserRepository) {
				user := &model.User{Name: "Test User"}
				user.ID = 1
				userRepo.On("Get", mock.Anything, uint64(1)).Return(user, nil)

				expectedTags := []model.UserTag{
					{Name: "VIP", Color: "#FF6B6B"},
					{Name: "New", Color: "#00FF00"},
				}
				tagRepo.On("GetUserTags", mock.Anything, uint64(1)).Return(expectedTags, nil)
			},
			expectError: false,
			expectTags:  2,
		},
		{
			name:   "user not found",
			userID: 999,
			setupMock: func(tagRepo *MockUserTagRepository, userRepo *MockTagUserRepository) {
				userRepo.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
			expectTags:  0,
		},
		{
			name:   "repository error on get tags",
			userID: 1,
			setupMock: func(tagRepo *MockUserTagRepository, userRepo *MockTagUserRepository) {
				user := &model.User{Name: "Test User"}
				user.ID = 1
				userRepo.On("Get", mock.Anything, uint64(1)).Return(user, nil)
				tagRepo.On("GetUserTags", mock.Anything, uint64(1)).Return(nil, errors.New("database error"))
			},
			expectError: true,
			expectTags:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagRepo := &MockUserTagRepository{}
			userRepo := &MockTagUserRepository{}
			cache := NewMockTagCache()

			tt.setupMock(tagRepo, userRepo)

			svc := NewUserTagService(tagRepo, userRepo, cache)
			tags, err := svc.GetUserTags(context.Background(), tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, tags)
			} else {
				assert.NoError(t, err)
				assert.Len(t, tags, tt.expectTags)
			}
		})
	}
}

func TestUserTagService_BatchSetUserTags(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint64
		tagIDs      []uint64
		setupMock   func(*MockUserTagRepository, *MockTagUserRepository)
		expectError bool
	}{
		{
			name:   "successful batch set",
			userID: 1,
			tagIDs: []uint64{1, 2},
			setupMock: func(tagRepo *MockUserTagRepository, userRepo *MockTagUserRepository) {
				user := &model.User{Name: "Test User"}
				user.ID = 1
				userRepo.On("Get", mock.Anything, uint64(1)).Return(user, nil)

				tag1 := &model.UserTag{Name: "VIP"}
				tag1.ID = 1
				tag2 := &model.UserTag{Name: "New"}
				tag2.ID = 2
				tagRepo.On("GetTag", mock.Anything, uint64(1)).Return(tag1, nil)
				tagRepo.On("GetTag", mock.Anything, uint64(2)).Return(tag2, nil)
				tagRepo.On("BatchSetUserTags", mock.Anything, uint64(1), []uint64{1, 2}).Return(nil)
			},
			expectError: false,
		},
		{
			name:   "user not found",
			userID: 999,
			tagIDs: []uint64{1},
			setupMock: func(tagRepo *MockUserTagRepository, userRepo *MockTagUserRepository) {
				userRepo.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
		{
			name:   "tag not found",
			userID: 1,
			tagIDs: []uint64{1, 999},
			setupMock: func(tagRepo *MockUserTagRepository, userRepo *MockTagUserRepository) {
				user := &model.User{Name: "Test User"}
				user.ID = 1
				userRepo.On("Get", mock.Anything, uint64(1)).Return(user, nil)

				tag1 := &model.UserTag{Name: "VIP"}
				tag1.ID = 1
				tagRepo.On("GetTag", mock.Anything, uint64(1)).Return(tag1, nil)
				tagRepo.On("GetTag", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagRepo := &MockUserTagRepository{}
			userRepo := &MockTagUserRepository{}
			cache := NewMockTagCache()

			tt.setupMock(tagRepo, userRepo)

			svc := NewUserTagService(tagRepo, userRepo, cache)
			err := svc.BatchSetUserTags(context.Background(), tt.userID, tt.tagIDs)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserTagService_GetUsersByTag(t *testing.T) {
	tests := []struct {
		name        string
		tagID       uint64
		page        int
		pageSize    int
		setupMock   func(*MockUserTagRepository)
		expectError bool
		expectUsers int
		expectTotal int64
	}{
		{
			name:     "successful get users by tag",
			tagID:    1,
			page:     1,
			pageSize: 10,
			setupMock: func(m *MockUserTagRepository) {
				tag := &model.UserTag{Name: "VIP"}
				tag.ID = 1
				m.On("GetTag", mock.Anything, uint64(1)).Return(tag, nil)

				expectedUsers := []model.User{
					{Name: "User1"},
					{Name: "User2"},
				}
				m.On("GetUsersByTag", mock.Anything, uint64(1), 1, 10).Return(expectedUsers, int64(2), nil)
			},
			expectError: false,
			expectUsers: 2,
			expectTotal: 2,
		},
		{
			name:     "tag not found",
			tagID:    999,
			page:     1,
			pageSize: 10,
			setupMock: func(m *MockUserTagRepository) {
				m.On("GetTag", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
			expectUsers: 0,
			expectTotal: 0,
		},
		{
			name:     "repository error on get users",
			tagID:    1,
			page:     1,
			pageSize: 10,
			setupMock: func(m *MockUserTagRepository) {
				tag := &model.UserTag{Name: "VIP"}
				tag.ID = 1
				m.On("GetTag", mock.Anything, uint64(1)).Return(tag, nil)
				m.On("GetUsersByTag", mock.Anything, uint64(1), 1, 10).Return(nil, int64(0), errors.New("database error"))
			},
			expectError: true,
			expectUsers: 0,
			expectTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagRepo := &MockUserTagRepository{}
			userRepo := &MockTagUserRepository{}
			cache := NewMockTagCache()

			tt.setupMock(tagRepo)

			svc := NewUserTagService(tagRepo, userRepo, cache)
			users, total, err := svc.GetUsersByTag(context.Background(), tt.tagID, tt.page, tt.pageSize)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, users)
			} else {
				assert.NoError(t, err)
				assert.Len(t, users, tt.expectUsers)
				assert.Equal(t, tt.expectTotal, total)
			}
		})
	}
}

func TestUserTagService_BatchDeleteTags(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	tag1 := &model.UserTag{Name: "VIP"}
	tag1.ID = 1
	tag2 := &model.UserTag{Name: "New"}
	tag2.ID = 2

	tagRepo.On("GetTag", mock.Anything, uint64(1)).Return(tag1, nil)
	tagRepo.On("GetTag", mock.Anything, uint64(2)).Return(tag2, nil)
	tagRepo.On("DeleteTag", mock.Anything, uint64(1)).Return(nil)
	tagRepo.On("DeleteTag", mock.Anything, uint64(2)).Return(nil)

	svc := NewUserTagService(tagRepo, userRepo, cache)

	result, err := svc.BatchDeleteTags(context.Background(), []uint64{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
}

func TestUserTagService_BatchDeleteTags_WithFailures(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	tag1 := &model.UserTag{Name: "VIP"}
	tag1.ID = 1

	tagRepo.On("GetTag", mock.Anything, uint64(1)).Return(tag1, nil)
	tagRepo.On("GetTag", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
	tagRepo.On("DeleteTag", mock.Anything, uint64(1)).Return(nil)

	svc := NewUserTagService(tagRepo, userRepo, cache)

	result, err := svc.BatchDeleteTags(context.Background(), []uint64{1, 999})
	assert.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.FailedItems, 1)
	assert.Equal(t, uint64(999), result.FailedItems[0].ID)
}

func TestUserTagService_BatchAssignTags(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	user := &model.User{Name: "Test User"}
	user.ID = 1
	userRepo.On("Get", mock.Anything, uint64(1)).Return(user, nil)

	tag := &model.UserTag{Name: "VIP"}
	tag.ID = 1
	tagRepo.On("GetTag", mock.Anything, uint64(1)).Return(tag, nil)
	tagRepo.On("GetUserTags", mock.Anything, uint64(1)).Return([]model.UserTag{}, nil)
	tagRepo.On("AddTagToUser", mock.Anything, uint64(1), uint64(1)).Return(nil)

	svc := NewUserTagService(tagRepo, userRepo, cache)

	pairs := []UserTagPair{
		{UserID: 1, TagID: 1},
	}
	result, err := svc.BatchAssignTags(context.Background(), pairs)

	assert.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
}

func TestUserTagService_BatchAssignTags_WithFailures(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	user := &model.User{Name: "Test User"}
	user.ID = 1
	userRepo.On("Get", mock.Anything, uint64(1)).Return(user, nil)
	userRepo.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)

	tag := &model.UserTag{Name: "VIP"}
	tag.ID = 1
	tagRepo.On("GetTag", mock.Anything, uint64(1)).Return(tag, nil)
	tagRepo.On("GetUserTags", mock.Anything, uint64(1)).Return([]model.UserTag{}, nil)
	tagRepo.On("AddTagToUser", mock.Anything, uint64(1), uint64(1)).Return(nil)

	svc := NewUserTagService(tagRepo, userRepo, cache)

	pairs := []UserTagPair{
		{UserID: 1, TagID: 1},
		{UserID: 999, TagID: 1}, // User not found
	}
	result, err := svc.BatchAssignTags(context.Background(), pairs)

	assert.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.FailedItems, 1)
}

func TestUserTagService_BatchAssignTagsToUsers(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	user1 := &model.User{Name: "User 1"}
	user1.ID = 1
	user2 := &model.User{Name: "User 2"}
	user2.ID = 2
	userRepo.On("Get", mock.Anything, uint64(1)).Return(user1, nil)
	userRepo.On("Get", mock.Anything, uint64(2)).Return(user2, nil)

	tag1 := &model.UserTag{Name: "VIP"}
	tag1.ID = 1
	tag2 := &model.UserTag{Name: "New"}
	tag2.ID = 2
	tagRepo.On("GetTag", mock.Anything, uint64(1)).Return(tag1, nil)
	tagRepo.On("GetTag", mock.Anything, uint64(2)).Return(tag2, nil)
	tagRepo.On("GetUserTags", mock.Anything, mock.Anything).Return([]model.UserTag{}, nil)
	tagRepo.On("AddTagToUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	svc := NewUserTagService(tagRepo, userRepo, cache)

	result, err := svc.BatchAssignTagsToUsers(context.Background(), []uint64{1, 2}, []uint64{1, 2})

	assert.NoError(t, err)
	assert.Equal(t, 4, result.TotalCount) // 2 users * 2 tags
	assert.Equal(t, 4, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
}

func TestUserTagService_BatchRemoveTagsFromUsers(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	tagRepo.On("RemoveTagFromUser", mock.Anything, uint64(1), uint64(1)).Return(nil)
	tagRepo.On("RemoveTagFromUser", mock.Anything, uint64(1), uint64(2)).Return(nil)
	tagRepo.On("RemoveTagFromUser", mock.Anything, uint64(2), uint64(1)).Return(nil)
	tagRepo.On("RemoveTagFromUser", mock.Anything, uint64(2), uint64(2)).Return(nil)

	svc := NewUserTagService(tagRepo, userRepo, cache)

	result, err := svc.BatchRemoveTagsFromUsers(context.Background(), []uint64{1, 2}, []uint64{1, 2})

	assert.NoError(t, err)
	assert.Equal(t, 4, result.TotalCount) // 2 users * 2 tags
	assert.Equal(t, 4, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
}

func TestUserTagService_BatchRemoveTagsFromUsers_WithFailures(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	tagRepo.On("RemoveTagFromUser", mock.Anything, uint64(1), uint64(1)).Return(nil)
	tagRepo.On("RemoveTagFromUser", mock.Anything, uint64(1), uint64(2)).Return(errors.New("database error"))
	tagRepo.On("RemoveTagFromUser", mock.Anything, uint64(2), uint64(1)).Return(nil)
	tagRepo.On("RemoveTagFromUser", mock.Anything, uint64(2), uint64(2)).Return(nil)

	svc := NewUserTagService(tagRepo, userRepo, cache)

	result, err := svc.BatchRemoveTagsFromUsers(context.Background(), []uint64{1, 2}, []uint64{1, 2})

	assert.NoError(t, err)
	assert.Equal(t, 4, result.TotalCount)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.FailedItems, 1)
}

func TestUserTagService_BatchRemoveTagsFromUsers_IdempotentRemove(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	// Simulate tag not assigned error (should be treated as success - idempotent)
	tagRepo.On("RemoveTagFromUser", mock.Anything, uint64(1), uint64(1)).Return(errors.New("移除标签失败: 用户标签关联不存在"))

	svc := NewUserTagService(tagRepo, userRepo, cache)

	result, err := svc.BatchRemoveTagsFromUsers(context.Background(), []uint64{1}, []uint64{1})

	assert.NoError(t, err)
	assert.Equal(t, 1, result.TotalCount)
	// This should be counted as failure since the error message doesn't match exactly
	assert.Equal(t, 1, result.FailedCount)
}

func TestBatchOperationResult_Structure(t *testing.T) {
	result := BatchOperationResult{
		TotalCount:   10,
		SuccessCount: 8,
		FailedCount:  2,
		SuccessItems: []uint64{1, 2, 3, 4, 5, 6, 7, 8},
		FailedItems: []BatchErrorItem{
			{ID: 9, Message: "error 1"},
			{ID: 10, Message: "error 2"},
		},
	}

	assert.Equal(t, 10, result.TotalCount)
	assert.Equal(t, 8, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)
	assert.Len(t, result.SuccessItems, 8)
	assert.Len(t, result.FailedItems, 2)
}

func TestUserTagPair_Structure(t *testing.T) {
	pair := UserTagPair{
		UserID: 123,
		TagID:  456,
	}

	assert.Equal(t, uint64(123), pair.UserID)
	assert.Equal(t, uint64(456), pair.TagID)
}

func TestIsValidColor(t *testing.T) {
	tests := []struct {
		color    string
		expected bool
	}{
		{"#FF6B6B", true},
		{"#000000", true},
		{"#ffffff", true},
		{"#ABCDEF", true},
		{"", true}, // empty is allowed
		{"red", false},
		{"#FFF", false},    // too short
		{"#GGGGGG", false}, // invalid hex
		{"FF6B6B", false},  // missing #
	}

	for _, tt := range tests {
		t.Run(tt.color, func(t *testing.T) {
			result := isValidColor(tt.color)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUserTagService_CreateTag_ListTagsError(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	tagRepo.On("ListTags", mock.Anything).Return(nil, errors.New("database error"))

	svc := NewUserTagService(tagRepo, userRepo, cache)
	tag, err := svc.CreateTag(context.Background(), "VIP", "#FF6B6B", "VIP用户")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "获取标签列表失败")
	assert.Nil(t, tag)
}

func TestUserTagService_CreateTag_CreateError(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	tagRepo.On("ListTags", mock.Anything).Return([]model.UserTag{}, nil)
	tagRepo.On("CreateTag", mock.Anything, mock.AnythingOfType("*model.UserTag")).Return(errors.New("database error"))

	svc := NewUserTagService(tagRepo, userRepo, cache)
	tag, err := svc.CreateTag(context.Background(), "VIP", "#FF6B6B", "VIP用户")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "创建标签失败")
	assert.Nil(t, tag)
}

func TestUserTagService_ListTags_RepositoryError(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	tagRepo.On("ListTags", mock.Anything).Return(nil, errors.New("database error"))

	svc := NewUserTagService(tagRepo, userRepo, cache)
	tags, err := svc.ListTags(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "获取标签列表失败")
	assert.Nil(t, tags)
}

func TestUserTagService_UpdateTag_NameTooLong(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	svc := NewUserTagService(tagRepo, userRepo, cache)
	err := svc.UpdateTag(context.Background(), 1, string(make([]byte, 65)), "#FF6B6B", "description")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "标签名称长度不能超过64个字符")
}

func TestUserTagService_UpdateTag_UpdateError(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	tag := &model.UserTag{Name: "VIP", Color: "#FF6B6B"}
	tag.ID = 1
	tagRepo.On("GetTag", mock.Anything, uint64(1)).Return(tag, nil)
	tagRepo.On("UpdateTag", mock.Anything, mock.AnythingOfType("*model.UserTag")).Return(errors.New("database error"))

	svc := NewUserTagService(tagRepo, userRepo, cache)
	err := svc.UpdateTag(context.Background(), 1, "Updated", "#00FF00", "description")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "更新标签失败")
}

func TestUserTagService_DeleteTag_DeleteError(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	tag := &model.UserTag{Name: "VIP"}
	tag.ID = 1
	tagRepo.On("GetTag", mock.Anything, uint64(1)).Return(tag, nil)
	tagRepo.On("DeleteTag", mock.Anything, uint64(1)).Return(errors.New("database error"))

	svc := NewUserTagService(tagRepo, userRepo, cache)
	err := svc.DeleteTag(context.Background(), 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "删除标签失败")
}

func TestUserTagService_AddTagToUser_GetUserTagsError(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	user := &model.User{Name: "Test User"}
	user.ID = 1
	userRepo.On("Get", mock.Anything, uint64(1)).Return(user, nil)

	tag := &model.UserTag{Name: "VIP"}
	tag.ID = 1
	tagRepo.On("GetTag", mock.Anything, uint64(1)).Return(tag, nil)
	tagRepo.On("GetUserTags", mock.Anything, uint64(1)).Return(nil, errors.New("database error"))

	svc := NewUserTagService(tagRepo, userRepo, cache)
	err := svc.AddTagToUser(context.Background(), 1, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "获取用户标签失败")
}

func TestUserTagService_AddTagToUser_AddError(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	user := &model.User{Name: "Test User"}
	user.ID = 1
	userRepo.On("Get", mock.Anything, uint64(1)).Return(user, nil)

	tag := &model.UserTag{Name: "VIP"}
	tag.ID = 1
	tagRepo.On("GetTag", mock.Anything, uint64(1)).Return(tag, nil)
	tagRepo.On("GetUserTags", mock.Anything, uint64(1)).Return([]model.UserTag{}, nil)
	tagRepo.On("AddTagToUser", mock.Anything, uint64(1), uint64(1)).Return(errors.New("database error"))

	svc := NewUserTagService(tagRepo, userRepo, cache)
	err := svc.AddTagToUser(context.Background(), 1, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "添加标签失败")
}

func TestUserTagService_BatchSetUserTags_BatchSetError(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	user := &model.User{Name: "Test User"}
	user.ID = 1
	userRepo.On("Get", mock.Anything, uint64(1)).Return(user, nil)

	tag := &model.UserTag{Name: "VIP"}
	tag.ID = 1
	tagRepo.On("GetTag", mock.Anything, uint64(1)).Return(tag, nil)
	tagRepo.On("BatchSetUserTags", mock.Anything, uint64(1), []uint64{1}).Return(errors.New("database error"))

	svc := NewUserTagService(tagRepo, userRepo, cache)
	err := svc.BatchSetUserTags(context.Background(), 1, []uint64{1})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "批量设置标签失败")
}

func TestUserTagService_CreateTag_NilCache(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}

	tagRepo.On("ListTags", mock.Anything).Return([]model.UserTag{}, nil)
	tagRepo.On("CreateTag", mock.Anything, mock.AnythingOfType("*model.UserTag")).Return(nil)

	// Create service with nil cache
	svc := NewUserTagService(tagRepo, userRepo, nil)
	tag, err := svc.CreateTag(context.Background(), "VIP", "#FF6B6B", "VIP用户")

	assert.NoError(t, err)
	assert.NotNil(t, tag)
}

func TestUserTagService_ListTags_NilCache(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}

	expectedTags := []model.UserTag{
		{Name: "VIP", Color: "#FF6B6B"},
	}
	tagRepo.On("ListTags", mock.Anything).Return(expectedTags, nil)

	// Create service with nil cache
	svc := NewUserTagService(tagRepo, userRepo, nil)
	tags, err := svc.ListTags(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedTags, tags)
}

func TestUserTagService_BatchAssignTagsToUsers_WithFailures(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	user1 := &model.User{Name: "User 1"}
	user1.ID = 1
	userRepo.On("Get", mock.Anything, uint64(1)).Return(user1, nil)
	userRepo.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)

	tag1 := &model.UserTag{Name: "VIP"}
	tag1.ID = 1
	tagRepo.On("GetTag", mock.Anything, uint64(1)).Return(tag1, nil)
	tagRepo.On("GetUserTags", mock.Anything, uint64(1)).Return([]model.UserTag{}, nil)
	tagRepo.On("AddTagToUser", mock.Anything, uint64(1), uint64(1)).Return(nil)

	svc := NewUserTagService(tagRepo, userRepo, cache)

	result, err := svc.BatchAssignTagsToUsers(context.Background(), []uint64{1, 999}, []uint64{1})

	assert.NoError(t, err)
	assert.Equal(t, 2, result.TotalCount) // 2 users * 1 tag
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.FailedItems, 1)
}

func TestUserTagService_BatchRemoveTagsFromUsers_WithIdempotentErrors(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	// The RemoveTagFromUser wraps the error, so the idempotent check won't match
	// This tests the failure path
	tagRepo.On("RemoveTagFromUser", mock.Anything, uint64(1), uint64(1)).Return(errors.New("some other error"))

	svc := NewUserTagService(tagRepo, userRepo, cache)

	result, err := svc.BatchRemoveTagsFromUsers(context.Background(), []uint64{1}, []uint64{1})

	assert.NoError(t, err)
	assert.Equal(t, 1, result.TotalCount)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
}

func TestUserTagService_UpdateTag_EmptyNameAndColor(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	tag := &model.UserTag{Name: "VIP", Color: "#FF6B6B", Description: "Old description"}
	tag.ID = 1
	tagRepo.On("GetTag", mock.Anything, uint64(1)).Return(tag, nil)
	tagRepo.On("UpdateTag", mock.Anything, mock.AnythingOfType("*model.UserTag")).Return(nil)

	svc := NewUserTagService(tagRepo, userRepo, cache)
	// Update with empty name and color - should keep original values
	err := svc.UpdateTag(context.Background(), 1, "", "", "New description")

	assert.NoError(t, err)
}

func TestUserTagService_ListTags_EmptyList(t *testing.T) {
	tagRepo := &MockUserTagRepository{}
	userRepo := &MockTagUserRepository{}
	cache := NewMockTagCache()

	tagRepo.On("ListTags", mock.Anything).Return([]model.UserTag{}, nil)

	svc := NewUserTagService(tagRepo, userRepo, cache)
	tags, err := svc.ListTags(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, tags)
}
