package services

import (
	"temp/models"
	"temp/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_Register(t *testing.T) {
	tests := []struct {
		name        string
		userName    string
		email       string
		password    string
		setupMock   func(*MockUserRepository)
		expectedErr string
	}{
		{
			name:     "successful registration",
			userName: "Test User",
			email:    "test@example.com",
			password: "password123",
			setupMock: func(mockRepo *MockUserRepository) {
				// First call to FindByEmail should return error (user doesn't exist)
				mockRepo.On("FindByEmail", "test@example.com").Return(nil, repositories.ErrUserNotFound)
				// Create should succeed
				mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)
			},
			expectedErr: "",
		},
		{
			name:     "user already exists",
			userName: "Test User",
			email:    "test@example.com",
			password: "password123",
			setupMock: func(mockRepo *MockUserRepository) {
				existingUser := &models.User{ID: 1, Email: "test@example.com"}
				mockRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil)
			},
			expectedErr: "user already exists",
		},
		{
			name:     "empty password",
			userName: "Test User",
			email:    "test@example.com",
			password: "", // Empty password - bcrypt handles this fine
			setupMock: func(mockRepo *MockUserRepository) {
				mockRepo.On("FindByEmail", "test@example.com").Return(nil, repositories.ErrUserNotFound)
				mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)
			},
			expectedErr: "", // This will succeed as bcrypt handles empty passwords
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			tt.setupMock(mockRepo)

			service := NewUserService(mockRepo, "test-secret", 24)
			user, err := service.Register(tt.userName, tt.email, tt.password)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.userName, user.Name)
				assert.Equal(t, tt.email, user.Email)
				assert.NotEmpty(t, user.Password) // Password should be hashed
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Authenticate(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		password    string
		setupMock   func(*MockUserRepository)
		expectedErr string
		expectToken bool
	}{
		{
			name:     "successful authentication",
			email:    "test@example.com",
			password: "password123",
			setupMock: func(mockRepo *MockUserRepository) {
				user := &models.User{
					ID:       1,
					Name:     "Test User",
					Email:    "test@example.com",
					Password: "$2a$10$vnz04c9pQOhKP3lc7p4LLOZYHapMZBdodhQdv5TYw/4gL3.xpGv4m", // "password123"
				}
				mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)
			},
			expectedErr: "",
			expectToken: true,
		},
		{
			name:     "user not found",
			email:    "nonexistent@example.com",
			password: "password123",
			setupMock: func(mockRepo *MockUserRepository) {
				mockRepo.On("FindByEmail", "nonexistent@example.com").Return(nil, repositories.ErrUserNotFound)
			},
			expectedErr: "invalid credentials",
			expectToken: false,
		},
		{
			name:     "invalid password",
			email:    "test@example.com",
			password: "wrongpassword",
			setupMock: func(mockRepo *MockUserRepository) {
				user := &models.User{
					ID:       1,
					Name:     "Test User",
					Email:    "test@example.com",
					Password: "$2a$10$vnz04c9pQOhKP3lc7p4LLOZYHapMZBdodhQdv5TYw/4gL3.xpGv4m", // "password123"
				}
				mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)
			},
			expectedErr: "invalid credentials",
			expectToken: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			tt.setupMock(mockRepo)

			service := NewUserService(mockRepo, "test-secret", 24)
			token, user, err := service.Authenticate(tt.email, tt.password)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Empty(t, token)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				if tt.expectToken {
					assert.NotEmpty(t, token)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestNewUserService(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo, "test-secret", 24)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
	assert.Equal(t, "test-secret", service.jwtSecret)
	assert.Equal(t, int64(24*60*60*1000000000), int64(service.expDuration)) // 24 hours in nanoseconds
}

// MockUserRepository is a mock implementation of UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

// Ensure MockUserRepository implements UserRepository interface
var _ repositories.UserRepository = (*MockUserRepository)(nil)

func (m *MockUserRepository) Migrate() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUserRepository) Create(u *models.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// MockUserService is a mock implementation of UserServiceInterface interface
type MockUserService struct {
	mock.Mock
}

// Ensure MockUserService implements UserServiceInterface interface
var _ UserServiceInterface = (*MockUserService)(nil)

func (m *MockUserService) Register(name, email, password string) (*models.User, error) {
	args := m.Called(name, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) Authenticate(email, password string) (string, *models.User, error) {
	args := m.Called(email, password)
	if args.Get(1) == nil {
		return args.String(0), nil, args.Error(2)
	}
	return args.String(0), args.Get(1).(*models.User), args.Error(2)
}

// TestMockUserService tests the mock service implementation
func TestMockUserService(t *testing.T) {
	mockService := &MockUserService{}

	// Test Register
	expectedUser := &models.User{ID: 1, Name: "Test", Email: "test@example.com"}
	mockService.On("Register", "Test", "test@example.com", "password").Return(expectedUser, nil)
	user, err := mockService.Register("Test", "test@example.com", "password")
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockService.AssertExpectations(t)

	// Test Authenticate
	token := "test-token"
	mockService.On("Authenticate", "test@example.com", "password").Return(token, expectedUser, nil)
	returnedToken, returnedUser, err := mockService.Authenticate("test@example.com", "password")
	assert.NoError(t, err)
	assert.Equal(t, token, returnedToken)
	assert.Equal(t, expectedUser, returnedUser)
	mockService.AssertExpectations(t)
}
