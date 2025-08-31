# Internationalization (i18n) Implementation Guide

## Overview

This inventory API now supports **complete internationalization (i18n)** for all error messages and success messages in three languages:
- **English (en-US)** - Default language
- **Indonesian (id-ID)** - Bahasa Indonesia
- **Japanese (ja-JP)** - 日本語

## ✅ Complete Implementation Status

All components have been updated to support i18n:

### ✅ Updated Handlers
- **Auth Handler** - Login/Register endpoints with localized messages
- **User Handler** - All user management endpoints with localized messages
- **Category Handler** - All category management endpoints with localized messages

### ✅ Updated Services
- **Auth Service** - Register/Login business logic with localized error messages
- **User Service** - User management business logic with localized error messages
- **Category Service** - Category management business logic with localized error messages

### ✅ Updated Middleware
- **Auth Middleware** - JWT validation with localized error messages
- **Role Authorization Middleware** - Permission checks with localized error messages

### ✅ Core System Components
- **Error Domain** - Complete i18n error system with message keys
- **Web Response Handler** - Automatic language detection and localized responses
- **i18n Utility** - Comprehensive message translation system

## How It Works

The i18n system automatically detects the user's preferred language from HTTP headers and returns localized messages accordingly.

### Language Detection Priority

1. **Accept-Language header** - Primary method (e.g., `Accept-Language: id-ID,id;q=0.9,en;q=0.8`)
2. **X-Language header** - Fallback method (e.g., `X-Language: ja-JP`)
3. **Default fallback** - English (en-US) if no valid language is detected

### Language Code Normalization

The system normalizes various language code formats:
- `en`, `en-us`, `en_US` → `en-US`
- `id`, `id-id`, `id_ID` → `id-ID`
- `ja`, `ja-jp`, `ja_JP` → `ja-JP`
- Unknown codes → `en-US` (default)

## Implementation Details

### 1. Message Keys System

All messages are defined using message keys in `internal/utils/i18n.go`:

```go
// Error message keys
ErrUserNotFoundKey      MessageKey = "error.user.not_found"
ErrCategoryNotFoundKey  MessageKey = "error.category.not_found"

// Success message keys
SuccessUserCreatedKey    MessageKey = "success.user.created"
SuccessCategoryCreatedKey MessageKey = "success.category.created"
```

### 2. Domain Error Updates

The `domain.AppError` struct now supports i18n:

```go
type AppError struct {
    Code       int
    Message    string
    MessageKey utils.MessageKey // For i18n support
    Params     []string         // Parameters for message formatting
    Err        error
}

// Get localized error message
func (e *AppError) GetLocalizedMessage(langCode string) string {
    if e.MessageKey != "" {
        return utils.GetLocalizedMessage(e.MessageKey, langCode, e.Params...)
    }
    return e.Message
}
```

### 3. Web Response Handler Updates

Response handlers now use message keys:

```go
// New i18n-aware success response
func Success(c *fiber.Ctx, code int, messageKey utils.MessageKey, data any) error {
    langCode := GetLanguageFromContext(c)
    message := utils.GetLocalizedMessage(messageKey, langCode)
    // ... return JSON response
}

// New i18n-aware error handling
func HandleError(c *fiber.Ctx, err error) error {
    langCode := GetLanguageFromContext(c)
    // ... handle different error types with localized messages
}
```

### 4. Handler Updates

Both category and user handlers now use message keys:

```go
// Example from category handler
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
    // ... validation and business logic
    return web.Success(c, fiber.StatusCreated, utils.SuccessCategoryCreatedKey, category)
}

// Example error handling
func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
    id := c.Params("id")
    if id == "" {
        return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrCategoryIDRequiredKey))
    }
    // ... rest of the logic
}
```

## Usage Examples

### Making API Requests with Language Headers

#### Authentication Endpoints

**Register - Indonesian Response:**
```bash
curl -X POST "http://localhost:8080/api/auth/register" \
  -H "Accept-Language: id-ID" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "john_doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "status": "success",
  "message": "Pengguna berhasil dibuat",
  "data": {...}
}
```

**Login - Japanese Response:**
```bash
curl -X POST "http://localhost:8080/api/auth/login" \
  -H "Accept-Language: ja-JP" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "status": "success",
  "message": "ログイン成功",
  "data": {
    "user": {...},
    "access_token": "...",
    "refresh_token": "..."
  }
}
```

#### Category Endpoints
```bash
curl -X GET "http://localhost:8080/api/categories" \
  -H "Accept-Language: en-US"
```

Response:
```json
{
  "status": "success",
  "message": "Categories retrieved successfully",
  "data": [...]
}
```

#### Indonesian Response
```bash
curl -X GET "http://localhost:8080/api/categories" \
  -H "Accept-Language: id-ID"
```

Response:
```json
{
  "status": "success",
  "message": "Kategori berhasil diambil",
  "data": [...]
}
```

#### Japanese Response
```bash
curl -X GET "http://localhost:8080/api/categories" \
  -H "Accept-Language: ja-JP"
```

Response:
```json
{
  "status": "success",
  "message": "カテゴリが正常に取得されました",
  "data": [...]
}
```

#### Alternative Header Format
```bash
curl -X GET "http://localhost:8080/api/categories" \
  -H "X-Language: id-ID"
```

### Error Responses

#### English Error
```bash
curl -X GET "http://localhost:8080/api/categories/invalid-id" \
  -H "Accept-Language: en-US"
```

Response:
```json
{
  "status": "error",
  "message": "Category not found",
  "error": null
}
```

#### Indonesian Error
```bash
curl -X GET "http://localhost:8080/api/categories/invalid-id" \
  -H "Accept-Language: id-ID"
```

Response:
```json
{
  "status": "error",
  "message": "Kategori tidak ditemukan",
  "error": null
}
```

## Available Message Keys

### Error Messages

#### Common Errors
- `ErrBadRequestKey` - Bad request
- `ErrUnauthorizedKey` - Unauthorized access
- `ErrForbiddenKey` - Access forbidden
- `ErrNotFoundKey` - Resource not found
- `ErrConflictKey` - Resource conflict
- `ErrInternalKey` - Internal server error
- `ErrValidationKey` - Validation failed

#### User-Specific Errors
- `ErrUserNotFoundKey` - User not found
- `ErrUserNameExistsKey` - Username already exists
- `ErrUserEmailExistsKey` - Email already exists
- `ErrUserIDRequiredKey` - User ID is required
- `ErrUserNameRequiredKey` - Username is required
- `ErrUserEmailRequiredKey` - Email is required

#### Category-Specific Errors
- `ErrCategoryNotFoundKey` - Category not found
- `ErrCategoryCodeExistsKey` - Category code already exists
- `ErrCategoryIDRequiredKey` - Category ID is required
- `ErrCategoryCodeRequiredKey` - Category code is required
- `ErrCategoryNameRequiredKey` - Category name is required

### Success Messages

#### User Operations
- `SuccessUserCreatedKey` - User created successfully
- `SuccessUserUpdatedKey` - User updated successfully
- `SuccessUserDeletedKey` - User deleted successfully
- `SuccessUserRetrievedKey` - User retrieved successfully
- `SuccessUserRetrievedByNameKey` - User retrieved successfully by name
- `SuccessUserRetrievedByEmailKey` - User retrieved successfully by email
- `SuccessUserCountedKey` - Users counted successfully
- `SuccessUserExistenceCheckedKey` - User existence checked successfully

#### Category Operations
- `SuccessCategoryCreatedKey` - Category created successfully
- `SuccessCategoryUpdatedKey` - Category updated successfully
- `SuccessCategoryDeletedKey` - Category deleted successfully
- `SuccessCategoryRetrievedKey` - Categories retrieved successfully
- `SuccessCategoryRetrievedByCodeKey` - Category retrieved successfully by code
- `SuccessCategoryHierarchyRetrievedKey` - Category hierarchy retrieved successfully
- `SuccessCategoryCountedKey` - Categories counted successfully
- `SuccessCategoryExistenceCheckedKey` - Category existence checked successfully

## Backward Compatibility

The implementation maintains backward compatibility by providing both new i18n-aware functions and legacy functions:

- `web.Success()` - New i18n-aware version
- `web.SuccessWithMessage()` - Legacy version with custom message
- `web.HandleError()` - New i18n-aware version
- `web.HandleErrorWithMessage()` - Legacy version

## Adding New Languages

To add support for new languages:

1. Add translations to the `messageTranslations` map in `internal/utils/i18n.go`
2. Update the `normalizeLanguageCode()` function to handle the new language code
3. Add the new language code to `GetAvailableLanguages()` function

Example for French (fr-FR):
```go
ErrUserNotFoundKey: {
    "en-US": "User not found",
    "id-ID": "Pengguna tidak ditemukan",
    "ja-JP": "ユーザーが見つかりません",
    "fr-FR": "Utilisateur non trouvé", // Add French translation
},
```

## Testing

Use the provided test file `examples/i18n_test.go` to test the i18n functionality:

```bash
go run examples/i18n_test.go
```

This will output examples of messages in all supported languages and demonstrate language code normalization.
