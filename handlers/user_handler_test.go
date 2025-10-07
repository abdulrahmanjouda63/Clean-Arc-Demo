package handlers

import (
	"encoding/json"
	"net/http"
	"testing"
	"temp/models"
	"temp/services"
	"temp/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockUserService)
		expectedStatus int
		expectedBody   gin.H
	}{
		{
			name: "successful registration",
			requestBody: gin.H{
				"name":     "Test User",
				"email":    "test@example.com",
				"password": "password123",
			},
			setupMock: func(mockService *MockUserService) {
				expectedUser := &models.User{
					ID:    1,
					Name:  "Test User",
					Email: "test@example.com",
				}
				mockService.On("Register", "Test User", "test@example.com", "password123").Return(expectedUser, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: gin.H{
				"id":    float64(1),
				"name":  "Test User",
				"email": "test@example.com",
			},
		},
		{
			name: "invalid request body",
			requestBody: gin.H{
				"name": "Test User",
				// Missing email and password
			},
			setupMock: func(mockService *MockUserService) {
				// No mock expectations as validation should fail before service call
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: gin.H{
				"error": mock.AnythingOfType("string"),
			},
		},
		{
			name: "service error",
			requestBody: gin.H{
				"name":     "Test User",
				"email":    "test@example.com",
				"password": "password123",
			},
			setupMock: func(mockService *MockUserService) {
				mockService.On("Register", "Test User", "test@example.com", "password123").Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: gin.H{
				"error": mock.AnythingOfType("string"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockUserService{}
			tt.setupMock(mockService)

			handler := NewUserHandler(mockService)
			req := testutils.CreateTestRequest("POST", "/register", tt.requestBody)
			c, w := testutils.CreateTestContext(req)

			handler.Register(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			
			var response gin.H
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			
			if tt.name == "successful registration" {
				assert.Equal(t, tt.expectedBody["id"], response["id"])
				assert.Equal(t, tt.expectedBody["name"], response["name"])
				assert.Equal(t, tt.expectedBody["email"], response["email"])
			} else {
				assert.Contains(t, response, "error")
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockUserService)
		expectedStatus int
		expectedBody   gin.H
	}{
		{
			name: "successful login",
			requestBody: gin.H{
				"email":    "test@example.com",
				"password": "password123",
			},
			setupMock: func(mockService *MockUserService) {
				expectedUser := &models.User{
					ID:    1,
					Name:  "Test User",
					Email: "test@example.com",
				}
				mockService.On("Authenticate", "test@example.com", "password123").Return("test-token", expectedUser, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: gin.H{
				"token": "test-token",
				"user": gin.H{
					"id":    float64(1),
					"email": "test@example.com",
					"name":  "Test User",
				},
			},
		},
		{
			name: "invalid credentials",
			requestBody: gin.H{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			setupMock: func(mockService *MockUserService) {
				mockService.On("Authenticate", "test@example.com", "wrongpassword").Return("", nil, assert.AnError)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: gin.H{
				"error": "invalid credentials",
			},
		},
		{
			name: "invalid request body",
			requestBody: gin.H{
				"email": "test@example.com",
				// Missing password
			},
			setupMock: func(mockService *MockUserService) {
				// No mock expectations as validation should fail before service call
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: gin.H{
				"error": mock.AnythingOfType("string"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockUserService{}
			tt.setupMock(mockService)

			handler := NewUserHandler(mockService)
			req := testutils.CreateTestRequest("POST", "/login", tt.requestBody)
			c, w := testutils.CreateTestContext(req)

			handler.Login(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			
			var response gin.H
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			
			if tt.name == "successful login" {
				assert.Equal(t, tt.expectedBody["token"], response["token"])
				assert.NotNil(t, response["user"])
			} else {
				assert.Contains(t, response, "error")
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Profile(t *testing.T) {
	tests := []struct {
		name           string
		userID         interface{}
		expectedStatus int
		expectedBody   gin.H
	}{
		{
			name:           "successful profile access",
			userID:         uint(1),
			expectedStatus: http.StatusOK,
			expectedBody: gin.H{
				"message":  "this is a protected route",
				"user_id":  float64(1),
			},
		},
		{
			name:           "no user in context",
			userID:         nil,
			expectedStatus: http.StatusUnauthorized,
			expectedBody: gin.H{
				"error": "no user in context",
			},
		},
		{
			name:           "user ID as string",
			userID:         "123",
			expectedStatus: http.StatusOK,
			expectedBody: gin.H{
				"message":  "this is a protected route",
				"user_id":  float64(123),
			},
		},
		{
			name:           "user ID as float64",
			userID:         float64(456),
			expectedStatus: http.StatusOK,
			expectedBody: gin.H{
				"message":  "this is a protected route",
				"user_id":  float64(456),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockUserService{}
			handler := NewUserHandler(mockService)
			
			req := testutils.CreateTestRequest("GET", "/profile", nil)
			c, w := testutils.CreateTestContext(req)

			if tt.userID != nil {
				c.Set("user_id", tt.userID)
			}

			handler.Profile(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			
			var response gin.H
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			
			if tt.name == "successful profile access" || tt.name == "user ID as string" || tt.name == "user ID as float64" {
				assert.Equal(t, tt.expectedBody["message"], response["message"])
				assert.Equal(t, tt.expectedBody["user_id"], response["user_id"])
			} else {
				assert.Equal(t, tt.expectedBody["error"], response["error"])
			}
		})
	}
}

// MockUserService is a mock implementation of UserServiceInterface interface
type MockUserService struct {
	mock.Mock
}

// Ensure MockUserService implements UserServiceInterface interface
var _ services.UserServiceInterface = (*MockUserService)(nil)

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

func TestNewUserHandler(t *testing.T) {
	mockService := &MockUserService{}
	handler := NewUserHandler(mockService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
}
