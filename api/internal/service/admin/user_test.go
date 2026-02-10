package admin

import (
	"context"
	"errors"
	"testing"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/mocks"
	"gamelink/pkg/apierr"
)

// TestGetUser demonstrates how to test AdminService.GetUser with mock repository.
func TestGetUser(t *testing.T) {
	tests := []struct {
		name          string
		userID        uint64
		mockSetup     func(*mocks.MockUserRepository)
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:   "successful retrieval",
			userID: 1,
			mockSetup: func(m *mocks.MockUserRepository) {
				m.GetFunc = func(ctx context.Context, id uint64) (*model.User, error) {
					return &model.User{
						ID:    1,
						Name:  "Test User",
						Email: "test@example.com",
						Role:  model.RoleUser,
					}, nil
				}
			},
			expectedUser: &model.User{
				ID:    1,
				Name:  "Test User",
				Email: "test@example.com",
				Role:  model.RoleUser,
			},
			expectedError: nil,
		},
		{
			name:   "user not found",
			userID: 999,
			mockSetup: func(m *mocks.MockUserRepository) {
				m.GetFunc = func(ctx context.Context, id uint64) (*model.User, error) {
					return nil, repository.ErrNotFound
				}
			},
			expectedUser:  nil,
			expectedError: repository.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: Setup mock repository
			mockUserRepo := &mocks.MockUserRepository{}
			tt.mockSetup(mockUserRepo)

			svc := &AdminService{
				users: mockUserRepo,
			}

			// Act: Call the method under test
			user, err := svc.GetUser(context.Background(), tt.userID)

			// Assert: Verify results
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}

			if tt.expectedUser != nil && user != nil {
				if user.ID != tt.expectedUser.ID || user.Name != tt.expectedUser.Name {
					t.Errorf("expected user %+v, got %+v", tt.expectedUser, user)
				}
			} else if (tt.expectedUser == nil) != (user == nil) {
				t.Errorf("user mismatch: expected %+v, got %+v", tt.expectedUser, user)
			}
		})
	}
}

// TestCreateUser demonstrates testing user creation with validation.
func TestCreateUser(t *testing.T) {
	tests := []struct {
		name          string
		input         CreateUserInput
		mockSetup     func(*mocks.MockUserRepository)
		expectedError error
	}{
		{
			name: "successful creation",
			input: CreateUserInput{
				Name:     "New User",
				Email:    "new@example.com",
				Phone:    "13800138000",
				Password: "ValidPass123",
				Role:     model.RoleUser,
				Status:   model.UserStatusActive,
			},
			mockSetup: func(m *mocks.MockUserRepository) {
				// Check for existing user
				m.FindByEmailFunc = func(ctx context.Context, email string) (*model.User, error) {
					return nil, repository.ErrNotFound // No existing user
				}
				m.FindByPhoneFunc = func(ctx context.Context, phone string) (*model.User, error) {
					return nil, repository.ErrNotFound // No existing user
				}
				// Allow creation
				m.CreateFunc = func(ctx context.Context, user *model.User) error {
					user.ID = 123 // Simulate DB auto-increment
					return nil
				}
			},
			expectedError: nil,
		},
		{
			name: "validation error - empty name",
			input: CreateUserInput{
				Name:     "",
				Email:    "test@example.com",
				Password: "ValidPass123",
				Role:     model.RoleUser,
				Status:   model.UserStatusActive,
			},
			mockSetup:     func(m *mocks.MockUserRepository) {},
			expectedError: apierr.ErrValidation,
		},
		{
			name: "duplicate email",
			input: CreateUserInput{
				Name:     "User",
				Email:    "existing@example.com",
				Password: "ValidPass123",
				Role:     model.RoleUser,
				Status:   model.UserStatusActive,
			},
			mockSetup: func(m *mocks.MockUserRepository) {
				m.FindByEmailFunc = func(ctx context.Context, email string) (*model.User, error) {
					return &model.User{ID: 1, Email: email}, nil // Existing user
				}
			},
			expectedError: apierr.ErrDuplicateEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUserRepo := &mocks.MockUserRepository{}
			tt.mockSetup(mockUserRepo)

			svc := &AdminService{
				users: mockUserRepo,
			}

			// Act
			user, err := svc.CreateUser(context.Background(), tt.input)

			// Assert
			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else if !apierr.IsType(err, tt.expectedError) {
					t.Errorf("expected error type %T, got %T", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if user == nil {
					t.Error("expected user to be created, got nil")
				}
			}
		})
	}
}

// TestListUsersWithOptions demonstrates testing pagination and filtering.
func TestListUsersWithOptions(t *testing.T) {
	mockUsers := []model.User{
		{ID: 1, Name: "Alice", Role: model.RoleUser},
		{ID: 2, Name: "Bob", Role: model.RolePlayer},
	}

	tests := []struct {
		name          string
		opts          UserListOptions
		mockSetup     func(*mocks.MockUserRepository)
		expectedCount int
		expectedError error
	}{
		{
			name: "list all users",
			opts: UserListOptions{
				Page:     1,
				PageSize: 10,
			},
			mockSetup: func(m *mocks.MockUserRepository) {
				m.ListWithFiltersFunc = func(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
					return mockUsers, int64(len(mockUsers)), nil
				}
			},
			expectedCount: 2,
			expectedError: nil,
		},
		{
			name: "filter by role",
			opts: UserListOptions{
				Page:     1,
				PageSize: 10,
				Role:     []string{"user"},
			},
			mockSetup: func(m *mocks.MockUserRepository) {
				m.ListWithFiltersFunc = func(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
					filtered := []model.User{mockUsers[0]} // Only Alice
					return filtered, 1, nil
				}
			},
			expectedCount: 1,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUserRepo := &mocks.MockUserRepository{}
			tt.mockSetup(mockUserRepo)

			svc := &AdminService{
				users: mockUserRepo,
			}

			// Act
			users, pagination, err := svc.ListUsersWithOptions(context.Background(), tt.opts)

			// Assert
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}

			if len(users) != tt.expectedCount {
				t.Errorf("expected %d users, got %d", tt.expectedCount, len(users))
			}

			if pagination != nil && pagination.Total != int64(tt.expectedCount) {
				t.Errorf("expected total %d, got %d", tt.expectedCount, pagination.Total)
			}
		})
	}
}
