# Swagger Documentation Guide

This document provides comprehensive information about the Swagger/OpenAPI documentation setup for this Go project.

## Overview

The project uses Swagger/OpenAPI 3.0 for API documentation with the following features:

- **Interactive API Documentation**: Swagger UI for testing endpoints
- **Code Generation**: Automatic documentation generation from Go annotations
- **Authentication Support**: JWT Bearer token authentication
- **Comprehensive Schemas**: Detailed request/response models
- **Error Handling**: Complete error response documentation

## Dependencies

The project uses the following Swagger-related packages:

- `github.com/swaggo/swag` - Swagger code generation
- `github.com/swaggo/gin-swagger` - Gin integration for Swagger
- `github.com/swaggo/files` - Static file serving for Swagger UI

## Project Structure

```
├── docs/
│   ├── swagger.yaml          # OpenAPI specification
│   ├── docs.go               # Generated docs (auto-generated)
│   ├── swagger.json          # Generated JSON (auto-generated)
│   └── SWAGGER.md            # This documentation
├── scripts/
│   ├── generate-swagger.sh   # Swagger generation script (Linux/Mac)
│   └── generate-swagger.bat # Swagger generation script (Windows)
├── handlers/
│   └── user_handler.go       # Handler with Swagger annotations
├── main.go                   # Main file with API documentation
└── routes/
    └── router.go             # Router with Swagger middleware
```

## Setup and Installation

### 1. Install Swagger Tools

```bash
# Using Makefile
make swagger-install

# Or manually
go install github.com/swaggo/swag/cmd/swag@latest
```

### 2. Generate Documentation

```bash
# Using Makefile
make swagger

# Or manually
swag init -g main.go -o ./docs
```

### 3. Start the Server

```bash
# Using Makefile
make swagger-serve

# Or manually
go run main.go
```

## Accessing Swagger UI

Once the server is running, you can access the Swagger UI at:

- **Local Development**: http://localhost:8080/swagger/index.html
- **Production**: https://your-domain.com/swagger/index.html

## API Endpoints Documentation

### Authentication Endpoints

#### POST /api/v1/register
- **Description**: Register a new user
- **Tags**: Authentication
- **Request Body**: 
  ```json
  {
    "name": "John Doe",
    "email": "john.doe@example.com",
    "password": "password123"
  }
  ```
- **Responses**:
  - `201`: User successfully registered
  - `400`: Validation error
  - `409`: User already exists

#### POST /api/v1/login
- **Description**: User login
- **Tags**: Authentication
- **Request Body**:
  ```json
  {
    "email": "john.doe@example.com",
    "password": "password123"
  }
  ```
- **Responses**:
  - `200`: Login successful (returns JWT token)
  - `400`: Validation error
  - `401`: Invalid credentials

### User Endpoints

#### GET /api/v1/profile
- **Description**: Get user profile
- **Tags**: User
- **Security**: Bearer token required
- **Headers**: `Authorization: Bearer <jwt_token>`
- **Responses**:
  - `200`: Profile retrieved successfully
  - `401`: Unauthorized

### Redis Endpoints

#### POST /api/v1/set-redis-key
- **Description**: Set Redis key-value pair
- **Tags**: Redis
- **Request Body**:
  ```json
  {
    "key": "user:session:123",
    "value": "active"
  }
  ```
- **Responses**:
  - `200`: Key set successfully
  - `400`: Validation error
  - `500`: Redis operation failed

## Authentication

The API uses JWT Bearer token authentication:

1. **Login**: POST to `/api/v1/login` with email/password
2. **Get Token**: Response includes a JWT token
3. **Use Token**: Include token in Authorization header for protected endpoints

### Example Authentication Flow

```bash
# 1. Register a user
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"password123"}'

# 2. Login to get token
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"password123"}'

# 3. Use token for protected endpoints
curl -X GET http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

## Swagger Annotations

### Handler Annotations

Each handler function includes Swagger annotations:

```go
// Register godoc
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body object{name=string,email=string,password=string} true "User registration data"
// @Success 201 {object} object{id=int,name=string,email=string} "User successfully registered"
// @Failure 400 {object} object{error=string} "Bad request - validation error"
// @Failure 409 {object} object{error=string} "Conflict - user already exists"
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
    // Handler implementation
}
```

### Main Package Annotations

The main package includes API-level documentation:

```go
// @title User Management API
// @version 1.0.0
// @description A comprehensive API for user management
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
```

## Available Commands

### Makefile Commands

```bash
# Install Swagger tools
make swagger-install

# Generate Swagger documentation
make swagger

# Start server with Swagger UI
make swagger-serve

# Run all tests
make test

# Install dependencies
make deps
```

### Script Commands

```bash
# Linux/Mac
./scripts/generate-swagger.sh

# Windows
scripts/generate-swagger.bat
```

## Customization

### Adding New Endpoints

1. **Add Handler**: Create handler function with Swagger annotations
2. **Add Route**: Register route in `routes/router.go`
3. **Regenerate Docs**: Run `make swagger`
4. **Test**: Access Swagger UI to test new endpoint

### Modifying Schemas

1. **Update YAML**: Modify `docs/swagger.yaml` for complex schemas
2. **Update Annotations**: Modify handler annotations for simple changes
3. **Regenerate**: Run `make swagger`

### Custom Styling

Swagger UI can be customized by modifying the HTML template or using custom CSS.

## Troubleshooting

### Common Issues

1. **Swagger UI Not Loading**:
   - Ensure server is running on correct port
   - Check that `/swagger/*any` route is registered
   - Verify `swag init` was run successfully

2. **Documentation Not Updating**:
   - Run `make swagger` to regenerate docs
   - Check for syntax errors in annotations
   - Verify file paths in `swag init` command

3. **Authentication Not Working**:
   - Verify JWT token format: `Bearer <token>`
   - Check token expiration
   - Ensure middleware is properly configured

### Debug Commands

```bash
# Check Swagger installation
swag version

# Generate docs with verbose output
swag init -g main.go -o ./docs --parseDependency --parseInternal

# Validate generated JSON
cat docs/swagger.json | jq .
```

## Best Practices

### Documentation

1. **Keep Annotations Updated**: Update Swagger annotations when changing handlers
2. **Use Descriptive Summaries**: Write clear, concise endpoint descriptions
3. **Include Examples**: Provide request/response examples
4. **Document Errors**: Include all possible error responses

### Code Organization

1. **Consistent Annotations**: Use consistent annotation format across handlers
2. **Schema Reuse**: Define reusable schemas in `swagger.yaml`
3. **Version Control**: Commit generated docs to version control
4. **CI/CD Integration**: Include Swagger generation in build pipeline

## Integration with CI/CD

### GitHub Actions Example

```yaml
name: Generate Swagger Docs
on: [push, pull_request]
jobs:
  swagger:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      - name: Install Swagger
        run: go install github.com/swaggo/swag/cmd/swag@latest
      - name: Generate Docs
        run: swag init -g main.go -o ./docs
      - name: Commit Changes
        run: |
          git config --local user.email "action@github.com"
          git config --local user.name "GitHub Action"
          git add docs/
          git commit -m "Update Swagger docs" || exit 0
          git push
```

## Future Enhancements

Potential improvements to the Swagger setup:

1. **Code Generation**: Generate client SDKs from OpenAPI spec
2. **Validation**: Add request/response validation middleware
3. **Testing**: Integrate Swagger with API testing tools
4. **Monitoring**: Add API usage analytics
5. **Versioning**: Implement API versioning strategy
6. **Rate Limiting**: Add rate limiting documentation
7. **Webhooks**: Document webhook endpoints
8. **GraphQL**: Consider GraphQL integration alongside REST
