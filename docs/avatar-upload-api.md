# Avatar Upload API Documentation

## Overview

The avatar upload functionality is now integrated directly into user creation and update operations. Users can upload avatars along with their profile data in a single request, supporting both multipart file uploads and URL-based avatars.

## Prerequisites

### Environment Variables

You need to set up Cloudinary credentials in your `.env` file:

```env
# Option 1: Use Cloudinary URL (recommended)
CLOUDINARY_URL=cloudinary://api_key:api_secret@cloud_name

# Option 2: Use individual credentials
CLOUDINARY_CLOUD_NAME=your_cloud_name
CLOUDINARY_API_KEY=your_api_key
CLOUDINARY_API_SECRET=your_api_secret
```

## API Endpoints

### 1. Create User with Avatar

**Endpoint:** `POST /api/v1/users`
**Authentication:** Required (Bearer Token + Admin Role)
**Content-Type:** `multipart/form-data` or `application/json`

#### Option A: Multipart Form Data (File Upload)

**Request:**
```bash
curl -X POST "http://localhost:5000/api/v1/users" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -F "name=john_doe" \
  -F "email=john@example.com" \
  -F "password=password123" \
  -F "fullName=John Doe" \
  -F "role=Staff" \
  -F "isActive=true" \
  -F "avatar=@/path/to/avatar.jpg"
```

#### Option B: JSON with Avatar URL

**Request:**
```bash
curl -X POST "http://localhost:5000/api/v1/users" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "john_doe",
    "email": "john@example.com",
    "password": "password123",
    "fullName": "John Doe",
    "role": "Staff",
    "isActive": true,
    "avatarUrl": "https://example.com/avatar.jpg"
  }'
```

### 2. Update Current User Profile with Avatar

**Endpoint:** `PATCH /api/v1/users/profile`
**Authentication:** Required (Bearer Token)
**Content-Type:** `multipart/form-data` or `application/json`

#### Option A: Multipart Form Data (File Upload)

**Request:**
```bash
curl -X PATCH "http://localhost:5000/api/v1/users/profile" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "fullName=John Updated" \
  -F "avatar=@/path/to/new-avatar.jpg"
```

#### Option B: JSON to Remove Avatar

**Request:**
```bash
curl -X PATCH "http://localhost:5000/api/v1/users/profile" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "avatarUrl": null
  }'
```

#### Option C: JSON with Avatar URL

**Request:**
```bash
curl -X PATCH "http://localhost:5000/api/v1/users/profile" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "fullName": "John Updated",
    "avatarUrl": "https://example.com/new-avatar.jpg"
  }'
```

### 3. Update Specific User with Avatar (Admin Only)

**Endpoint:** `PATCH /api/v1/users/:id`
**Authentication:** Required (Bearer Token + Admin Role)
**Content-Type:** `multipart/form-data` or `application/json`

Same request format as update current user profile, but targets a specific user by ID.

## Response Format

All successful operations return the updated user object:

```json
{
  "status": "success",
  "message": "User created successfully", // or "User updated successfully"
  "data": {
    "id": "01ARZ3NDEKTSV4RRFFQ69G5FAV",
    "name": "john_doe",
    "email": "john@example.com",
    "fullName": "John Doe",
    "role": "Staff",
    "employeeId": null,
    "preferredLang": "en-US",
    "isActive": true,
    "avatarUrl": "https://res.cloudinary.com/your-cloud/image/upload/v1234567890/avatars/user_01ARZ3NDEKTSV4RRFFQ69G5FAV_avatar.jpg",
    "createdAt": "2025-01-01T00:00:00Z",
    "updatedAt": "2025-01-01T00:00:00Z"
  }
}
```

## Avatar Handling Logic

### Content Type Detection

The API automatically detects the request content type:

- **`multipart/form-data`**: Supports both form fields and file uploads
- **`application/json`** or **`application/x-www-form-urlencoded`**: Supports only text-based avatar URLs

### Avatar Priority

1. **File Upload** (multipart): If an `avatar` file is provided, it takes priority
2. **Avatar URL** (JSON/form): If no file is provided, uses the `avatarUrl` field value
3. **No Avatar**: If neither is provided, `avatarUrl` remains `null` (optional)

### Avatar Management

- **File Uploads**: Automatically uploaded to Cloudinary with consistent naming: `user_{userId}_avatar`
- **URL Updates**: Validates URL format when provided via JSON
- **Avatar Removal**: Set `avatarUrl` to `null`, `""`, or `"null"` to remove
- **Old Avatar Cleanup**: Automatically deletes old Cloudinary avatars when updated

## File Upload Constraints

### Avatar Upload Configuration

- **Allowed File Types:** JPEG, PNG, GIF, WebP
- **Maximum File Size:** 5MB
- **Maximum Files:** 1 per request
- **Cloudinary Folder:** `avatars/`
- **Public ID Format:** `user_{userId}_avatar`
- **Overwrite:** Enabled (replaces existing avatar)

## Error Responses

### Avatar-Specific Errors

**File Type Not Allowed (400)**
```json
{
  "status": "error",
  "message": "File type not allowed",
  "error": "BAD_REQUEST"
}
```

**File Size Too Large (400)**
```json
{
  "status": "error",
  "message": "File size too large",
  "error": "BAD_REQUEST"
}
```

**File Upload Failed (400)**
```json
{
  "status": "error",
  "message": "File upload failed",
  "error": "BAD_REQUEST"
}
```

**Cloudinary Not Configured (400)**
```json
{
  "status": "error",
  "message": "Cloudinary configuration error",
  "error": "BAD_REQUEST"
}
```

### Standard User Errors

**User Not Found (404)**
```json
{
  "status": "error",
  "message": "User not found",
  "error": "NOT_FOUND"
}
```

**Validation Error (400)**
```json
{
  "status": "error",
  "message": "Validation failed: email is required",
  "error": "BAD_REQUEST"
}
```

## Benefits of Integrated Approach

1. **Simplified API**: Single endpoint for user operations with optional avatar
2. **Atomic Operations**: User data and avatar are updated together
3. **Better UX**: Users can upload everything in one step
4. **Flexible Input**: Supports both file uploads and URL-based avatars
5. **Optional Avatar**: Avatar is never required, avoiding validation errors
6. **Automatic Cleanup**: Old avatars are automatically managed

## Notes

1. **Avatar is Always Optional**: You can create/update users without providing any avatar data
2. **Content-Type Detection**: The API automatically handles different request formats
3. **File Overwrite**: New avatar files automatically replace old ones with the same user ID
4. **URL Validation**: Avatar URLs are validated when provided via JSON
5. **Error Isolation**: Avatar upload failures don't prevent user operations from succeeding
6. **Internationalization**: All error messages support multiple languages (English, Indonesian, Japanese)
7. **Graceful Degradation**: If Cloudinary is not configured, file uploads fail gracefully with clear error messages
