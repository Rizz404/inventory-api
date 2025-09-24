# Swagger Documentation

This directory contains the generated Swagger/OpenAPI documentation for the Inventory Management API.

## 📚 Generated Files

- `docs.go` - Go code with embedded Swagger spec (auto-generated)
- `swagger.json` - OpenAPI 3.0 specification in JSON format
- `swagger.yaml` - OpenAPI 3.0 specification in YAML format
- `models.go` - Custom Swagger model definitions for better documentation

## 🚀 Quick Start

### 1. Generate Documentation
```bash
make swagger-gen
```

### 2. Run the Application
```bash
make dev  # Generates docs + builds + runs
# OR
make run  # Just builds and runs
```

### 3. Access Documentation
- **Swagger UI**: http://localhost:8080/docs/
- **JSON Spec**: http://localhost:8080/docs/doc.json
- **API Base URL**: http://localhost:8080/api/v1/

## 📖 Available Endpoints

### Authentication
- `POST /auth/register` - Register a new user account
- `POST /auth/login` - Authenticate user and get JWT tokens

### Users (Coming Soon)
- `GET /users` - List all users
- `POST /users` - Create a new user
- `GET /users/{id}` - Get user by ID
- `PUT /users/{id}` - Update user
- `DELETE /users/{id}` - Delete user

### Assets (Coming Soon)
- `GET /assets` - List all assets
- `POST /assets` - Create a new asset
- `GET /assets/{id}` - Get asset by ID
- `PUT /assets/{id}` - Update asset
- `DELETE /assets/{id}` - Delete asset

### Categories (Coming Soon)
- `GET /categories` - List all categories
- `POST /categories` - Create a new category
- `GET /categories/{id}` - Get category by ID
- `PUT /categories/{id}` - Update category
- `DELETE /categories/{id}` - Delete category

### Locations (Coming Soon)
- `GET /locations` - List all locations
- `POST /locations` - Create a new location
- `GET /locations/{id}` - Get location by ID
- `PUT /locations/{id}` - Update location
- `DELETE /locations/{id}` - Delete location

## 🔐 Authentication

The API uses JWT (JSON Web Tokens) for authentication. After successful login, include the token in requests:

```bash
Authorization: Bearer <your-jwt-token>
```

## 📝 Request/Response Examples

### Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "johndoe",
    "email": "john.doe@example.com",
    "password": "password123"
  }'
```

**Response:**
```json
{
  "status": "success",
  "message": "User created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "johndoe",
    "email": "john.doe@example.com",
    "fullName": "",
    "role": "Employee",
    "employeeId": null,
    "preferredLang": "en",
    "isActive": true,
    "avatarUrl": null,
    "createdAt": "2023-01-01 12:00:00",
    "updatedAt": "2023-01-01 12:00:00"
  }
}
```

### Login User
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "password123"
  }'
```

**Response:**
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "johndoe",
      "email": "john.doe@example.com",
      "fullName": "",
      "role": "Employee",
      "employeeId": null,
      "preferredLang": "en",
      "isActive": true,
      "avatarUrl": null,
      "createdAt": "2023-01-01 12:00:00",
      "updatedAt": "2023-01-01 12:00:00"
    },
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

## 🚨 Error Responses

### Validation Error
```json
{
  "status": "error",
  "message": "Validation failed",
  "error": [
    {
      "field": "email",
      "tag": "required",
      "value": "",
      "message": "email is required"
    }
  ]
}
```

### Authentication Error
```json
{
  "status": "error",
  "message": "Invalid credentials"
}
```

### Internal Server Error
```json
{
  "status": "error",
  "message": "An unexpected error occurred"
}
```

## 🔧 Swagger Annotations Guide

### Handler Annotations
```go
// Register godoc
//	@Summary		Register a new user
//	@Description	Register a new user account with name, email, and password
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			registerPayload	body		domain.RegisterPayload	true	"User registration data"
//	@Success		201				{object}	web.JSONResponse{data=domain.User}	"User registered successfully"
//	@Failure		400				{object}	web.JSONResponse{error=[]web.ValidationError}	"Validation failed"
//	@Failure		409				{object}	web.JSONResponse	"User already exists"
//	@Failure		500				{object}	web.JSONResponse	"Internal server error"
//	@Router			/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
    // Implementation...
}
```

### Model Annotations (in struct tags)
```go
type RegisterPayload struct {
    Name     string `json:"name" form:"name" validate:"required,min=3,max=50" example:"johndoe"`
    Email    string `json:"email" form:"email" validate:"required,email" example:"john.doe@example.com"`
    Password string `json:"password" form:"password" validate:"required,min=5" example:"password123"`
}
```

## 🛠️ Development Commands

```bash
# Generate Swagger docs
make swagger-gen

# Install Swagger CLI (if not installed)
make install-swagger

# Clean generated docs
make clean

# Development workflow
make dev  # Generate docs + build + run
```

## 📊 Supported Features

- ✅ **Complete Auth endpoints** (register, login)
- ✅ **JWT Authentication** with Bearer tokens
- ✅ **Validation error handling** with detailed messages
- ✅ **Multi-language support** (i18n)
- ✅ **Pagination support** (both offset and cursor-based)
- ✅ **Comprehensive error responses**
- ✅ **Request/Response examples**
- 🔄 **User management endpoints** (coming soon)
- 🔄 **Asset management endpoints** (coming soon)
- 🔄 **Category management endpoints** (coming soon)
- 🔄 **Location management endpoints** (coming soon)

## 🎯 Best Practices

1. **Always regenerate docs** after API changes: `make swagger-gen`
2. **Use descriptive summaries** and descriptions in annotations
3. **Include all possible response codes** (2xx, 4xx, 5xx)
4. **Add examples** to request/response models
5. **Group related endpoints** using Tags
6. **Document authentication requirements** for protected endpoints
7. **Use consistent naming** for parameters and models

## 🚀 Next Steps

1. Add Swagger annotations to remaining handlers (users, assets, categories, locations)
2. Add authentication middleware documentation
3. Include file upload endpoints documentation
4. Add query parameter documentation for filtering and sorting
5. Document webhook endpoints (if any)

---

For more information about Swagger/OpenAPI annotations, visit:
- [Swaggo Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://spec.openapis.org/oas/v3.0.3/)
