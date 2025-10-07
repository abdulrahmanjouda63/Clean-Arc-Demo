package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"temp/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// CreateTestUser creates a test user with sample data
func CreateTestUser() *models.User {
	return &models.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
	}
}

// CreateTestUserWithPassword creates a test user with password
func CreateTestUserWithPassword() *models.User {
	return &models.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // "password"
	}
}

// CreateTestRequest creates a test HTTP request
func CreateTestRequest(method, url string, body interface{}) *http.Request {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req, _ := http.NewRequest(method, url, &buf)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// AssertJSONResponse asserts that the response matches expected JSON
func AssertJSONResponse(t assert.TestingT, expected interface{}, actual gin.H) {
	expectedJSON, _ := json.Marshal(expected)
	actualJSON, _ := json.Marshal(actual)
	assert.JSONEq(t, string(expectedJSON), string(actualJSON))
}

// CreateTestContext creates a test gin context with request and response recorder
func CreateTestContext(req *http.Request) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}

// SetUserInContext sets a user ID in the gin context for testing
func SetUserInContext(c *gin.Context, userID uint) {
	c.Set("user_id", userID)
}
