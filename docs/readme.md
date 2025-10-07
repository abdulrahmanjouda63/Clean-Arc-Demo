# API Documentation

This project includes comprehensive Swagger/OpenAPI documentation for all endpoints.

## Quick Start

1. **Generate Documentation**:
   ```bash
   make swagger
   ```

2. **Start the Server**:
   ```bash
   make run
   ```

3. **Access Swagger UI**: http://localhost:8080/swagger/index.html

## Available Endpoints

### Authentication
- `POST /api/v1/register` - Register a new user
- `POST /api/v1/login` - User login

### User Management
- `GET /api/v1/profile` - Get user profile (requires authentication)

### Redis Operations
- `POST /api/v1/set-redis-key` - Set Redis key-value pair

## Authentication

The API uses JWT Bearer token authentication:

1. Register a user via `/api/v1/register`
2. Login via `/api/v1/login` to get a JWT token
3. Use the token in the `Authorization` header: `Bearer <token>`

## Documentation Files

- `swagger.yaml` - OpenAPI 3.0 specification
- `swagger.json` - Generated JSON documentation
- `docs.go` - Generated Go documentation
- `SWAGGER.md` - Comprehensive Swagger guide

## Commands

```bash
# Generate Swagger documentation
make swagger

# Install Swagger tools
make swagger-install

# Start server with Swagger UI
make swagger-serve

# Run tests
make test
```

For detailed information, see [SWAGGER.md](./SWAGGER.md).
