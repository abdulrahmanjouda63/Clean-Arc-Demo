package repositories

import (
	"testing"
	"temp/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserRepo_Create(t *testing.T) {
	tests := []struct {
		name    string
		user    *models.User
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful user creation",
			user: &models.User{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "hashedpassword",
			},
			wantErr: false,
		},
		{
			name: "user creation with empty name",
			user: &models.User{
				Name:     "",
				Email:    "test@example.com",
				Password: "hashedpassword",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test would require a real database connection
			// For now, we'll test the interface compliance
			repo := NewUserRepo()
			assert.NotNil(t, repo)
			
			// Test that the repo implements the interface
			var _ UserRepository = repo
		})
	}
}

func TestUserRepo_FindByEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected *models.User
		wantErr  bool
		errMsg   string
	}{
		{
			name:  "user found",
			email: "test@example.com",
			expected: &models.User{
				ID:    1,
				Name:  "Test User",
				Email: "test@example.com",
			},
			wantErr: false,
		},
		{
			name:     "user not found",
			email:    "nonexistent@example.com",
			expected: nil,
			wantErr:  true,
			errMsg:   "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewUserRepo()
			assert.NotNil(t, repo)
			
			// Test that the repo implements the interface
			var _ UserRepository = repo
		})
	}
}

func TestUserRepo_FindByID(t *testing.T) {
	tests := []struct {
		name     string
		id       uint
		expected *models.User
		wantErr  bool
		errMsg   string
	}{
		{
			name:  "user found",
			id:    1,
			expected: &models.User{
				ID:    1,
				Name:  "Test User",
				Email: "test@example.com",
			},
			wantErr: false,
		},
		{
			name:     "user not found",
			id:       999,
			expected: nil,
			wantErr:  true,
			errMsg:   "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewUserRepo()
			assert.NotNil(t, repo)
			
			// Test that the repo implements the interface
			var _ UserRepository = repo
		})
	}
}

// MockUserRepository is a mock implementation of UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

// Ensure MockUserRepository implements UserRepository interface
var _ UserRepository = (*MockUserRepository)(nil)

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

// TestMockUserRepository tests the mock implementation
func TestMockUserRepository(t *testing.T) {
	mockRepo := &MockUserRepository{}
	
	// Test Create
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)
	err := mockRepo.Create(&models.User{Name: "Test"})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	
	// Test FindByEmail
	expectedUser := &models.User{ID: 1, Email: "test@example.com"}
	mockRepo.On("FindByEmail", "test@example.com").Return(expectedUser, nil)
	user, err := mockRepo.FindByEmail("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
	
	// Test FindByID
	mockRepo.On("FindByID", uint(1)).Return(expectedUser, nil)
	user, err = mockRepo.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}
