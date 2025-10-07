# Testing Guide

This document provides comprehensive information about the testing setup and how to run tests in this Go project.

## Overview

The project uses a layered architecture with comprehensive unit testing using mocks. The testing setup includes:

- **Unit Tests**: Test individual components in isolation using mocks
- **Mock Objects**: Generated mocks for repository and service interfaces
- **Test Utilities**: Helper functions for creating test data and assertions
- **Coverage Reports**: Detailed coverage analysis

## Testing Dependencies

The project uses the following testing libraries:

- `github.com/stretchr/testify` - Testing toolkit with assertions and mocks
- Built-in Go testing package

## Project Structure

```
├── handlers/
│   ├── user_handler.go
│   └── user_handler_test.go
├── services/
│   ├── interfaces.go
│   ├── user_service.go
│   └── user_service_test.go
├── repositories/
│   ├── interfaces.go
│   ├── user_repo.go
│   └── user_repo_test.go
├── utils/
│   ├── hash.go
│   └── hash_test.go
├── mocks/
│   ├── user_repository_mock.go
│   └── user_service_mock.go
├── testutils/
│   └── test_helpers.go
├── test_config.yaml
├── Makefile
└── TESTING.md
```

## Running Tests

### Using Makefile (Recommended)

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage report
make test-coverage

# Run specific test suites
make test-handlers
make test-services
make test-repositories
make test-utils

# Run tests with race detection
make test-race

# Clean test artifacts
make clean
```

### Using Go Commands Directly

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run tests for specific package
go test ./handlers/...
go test ./services/...
go test ./repositories/...
go test ./utils/...

# Run tests with race detection
go test -race ./...
```

## Test Categories

### 1. Unit Tests

Unit tests test individual components in isolation using mocks:

- **Handler Tests**: Test HTTP request/response handling
- **Service Tests**: Test business logic
- **Repository Tests**: Test data access layer (with mocks)
- **Utility Tests**: Test helper functions

### 2. Mock Objects

The project includes manually created mock objects:

- `MockUserRepository`: Mocks the user repository interface
- `MockUserService`: Mocks the user service interface

### 3. Test Utilities

Located in `testutils/test_helpers.go`:

- `CreateTestUser()`: Creates test user data
- `CreateTestRequest()`: Creates HTTP test requests
- `CreateTestContext()`: Creates Gin test context
- `AssertJSONResponse()`: Asserts JSON response equality
- `SetUserInContext()`: Sets user ID in context for testing

## Writing Tests

### Handler Tests

```go
func TestUserHandler_Register(t *testing.T) {
    mockService := &mocks.MockUserService{}
    mockService.On("Register", "Test User", "test@example.com", "password123").Return(expectedUser, nil)
    
    handler := NewUserHandler(mockService)
    req := testutils.CreateTestRequest("POST", "/register", requestBody)
    c, w := testutils.CreateTestContext(req)
    
    handler.Register(c)
    
    assert.Equal(t, http.StatusCreated, w.Code)
    mockService.AssertExpectations(t)
}
```

### Service Tests

```go
func TestUserService_Register(t *testing.T) {
    mockRepo := &mocks.MockUserRepository{}
    mockRepo.On("FindByEmail", "test@example.com").Return(nil, repositories.ErrUserNotFound)
    mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)
    
    service := NewUserService(mockRepo, "test-secret", 24)
    user, err := service.Register("Test User", "test@example.com", "password123")
    
    assert.NoError(t, err)
    assert.NotNil(t, user)
    mockRepo.AssertExpectations(t)
}
```

### Repository Tests

```go
func TestMockUserRepository(t *testing.T) {
    mockRepo := &mocks.MockUserRepository{}
    mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)
    
    err := mockRepo.Create(&models.User{Name: "Test"})
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

## Test Data

### Test Users

The test utilities provide functions to create test users:

```go
// Create a basic test user
user := testutils.CreateTestUser()

// Create a test user with password hash
user := testutils.CreateTestUserWithPassword()
```

### Test Requests

```go
// Create a JSON request
req := testutils.CreateTestRequest("POST", "/register", gin.H{
    "name":     "Test User",
    "email":    "test@example.com",
    "password": "password123",
})
```

## Coverage Reports

Generate and view coverage reports:

```bash
# Generate coverage report
make test-coverage

# View HTML coverage report
open coverage.html
```

The coverage report shows:
- Overall coverage percentage
- Line-by-line coverage analysis
- Coverage by package and file

## Best Practices

### 1. Test Structure

Follow the AAA pattern:
- **Arrange**: Set up test data and mocks
- **Act**: Execute the code under test
- **Assert**: Verify the results

### 2. Mock Usage

- Use mocks to isolate components
- Set up mock expectations before calling the code under test
- Always call `AssertExpectations(t)` to verify mock calls

### 3. Test Naming

- Use descriptive test names: `TestUserService_Register_Success`
- Include the scenario being tested

### 4. Test Data

- Use test utilities for consistent test data
- Avoid hardcoded values in tests
- Use meaningful test data

### 5. Assertions

- Use specific assertions: `assert.Equal(t, expected, actual)`
- Test both success and error cases
- Verify all important return values

## Continuous Integration

The tests are designed to run in CI environments:

```bash
# CI test command
make test-coverage
```

## Troubleshooting

### Common Issues

1. **Import Errors**: Ensure all dependencies are installed with `make deps`
2. **Mock Failures**: Check that mock expectations match actual calls
3. **Test Failures**: Run tests with `-v` flag for verbose output

### Debugging Tests

```bash
# Run specific test with verbose output
go test -v -run TestUserHandler_Register ./handlers/

# Run tests with race detection
go test -race ./...

# Run tests with coverage for specific package
go test -coverprofile=coverage.out ./handlers/
go tool cover -html=coverage.out
```

## Future Enhancements

Potential improvements to the testing setup:

1. **Integration Tests**: Add tests that use real database connections
2. **Performance Tests**: Add benchmark tests for critical paths
3. **Contract Tests**: Add tests to verify interface contracts
4. **Test Data Builders**: Create fluent builders for test data
5. **Custom Matchers**: Create custom assertion matchers for domain objects
