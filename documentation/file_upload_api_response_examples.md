# File Upload API Response Examples

## Success Responses

### Upload Avatar (Success)
```http
PATCH /api/users/profile
Content-Type: multipart/form-data

avatar: [image file]
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User updated successfully",
  "data": {
    "id": "01HQXYZ1234567890ABCDEFGH",
    "name": "john_doe",
    "email": "john@example.com",
    "fullName": "John Doe",
    "role": "user",
    "avatarUrl": "https://res.cloudinary.com/demo/image/upload/v1234567890/avatars/user_01HQXYZ1234567890ABCDEFGH_avatar.jpg",
    "isActive": true,
    "createdAt": "2025-01-15T10:30:00Z",
    "updatedAt": "2025-01-15T14:20:00Z"
  }
}
```

---

## Error Responses

### 1. File Size Too Large

**Request:**
```http
POST /api/users
Content-Type: multipart/form-data

name: john_doe
email: john@example.com
password: secret123
avatar: [10MB image file]
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Bad Request",
  "error": "avatar: File size too large (10.50 MB). Maximum allowed size is 5 MB"
}
```

**Frontend Display:**
```
❌ Upload Failed
File size too large (10.50 MB). Maximum allowed size is 5 MB

💡 Suggestion: Please compress your image or use a smaller file.
```

---

### 2. Invalid File Type

**Request:**
```http
PATCH /api/users/profile
Content-Type: multipart/form-data

avatar: document.pdf
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Bad Request",
  "error": "avatar: Invalid file type '.pdf'. Allowed types: JPG, JPEG, PNG, GIF, WEBP"
}
```

**Frontend Display:**
```
❌ Upload Failed
Invalid file type '.pdf'

✅ Allowed types: JPG, JPEG, PNG, GIF, WEBP
```

---

### 3. Empty File

**Request:**
```http
POST /api/assets
Content-Type: multipart/form-data

assetTag: LAPTOP-001
assetName: Dell Laptop
categoryId: 01HQCAT1234567890ABCDEFGH
dataMatrixImage: [0 byte file]
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Bad Request",
  "error": "dataMatrixImage: File is empty (0 bytes)"
}
```

**Frontend Display:**
```
❌ Upload Failed
The selected file is empty (0 bytes)

💡 Suggestion: Please select a valid image file.
```

---

### 4. File Not Valid Image

**Request:**
```http
PATCH /api/assets/01HQAST1234567890ABCDEFGH
Content-Type: multipart/form-data

dataMatrixImage: text_renamed_to_image.png
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Bad Request",
  "error": "dataMatrixImage: File 'text_renamed_to_image.png' is not a valid image file. The file content does not match any supported image format"
}
```

**Frontend Display:**
```
❌ Upload Failed
'text_renamed_to_image.png' is not a valid image file

The file content does not match any supported image format.
Please use an actual image file (JPG, PNG, GIF, or WEBP).
```

---

### 5. Filename Too Long

**Request:**
```http
POST /api/users
Content-Type: multipart/form-data

name: john_doe
email: john@example.com
password: secret123
avatar: [file with 300 character name]
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Bad Request",
  "error": "avatar: Filename too long (312 characters). Maximum allowed is 255 characters"
}
```

**Frontend Display:**
```
❌ Upload Failed
Filename is too long (312 characters)

Maximum allowed: 255 characters
💡 Suggestion: Please rename the file to a shorter name.
```

---

### 6. File Cannot Be Read

**Request:**
```http
PATCH /api/users/profile
Content-Type: multipart/form-data

avatar: [corrupted or locked file]
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Bad Request",
  "error": "avatar: Cannot read file: permission denied"
}
```

**Frontend Display:**
```
❌ Upload Failed
Cannot read the selected file

💡 Possible causes:
- File is corrupted
- File is locked by another program
- Insufficient permissions

Please try:
1. Closing any programs using this file
2. Selecting a different file
3. Restarting your device
```

---

### 7. Cloudinary Upload Errors

#### Network Timeout

**Request:**
```http
POST /api/assets
Content-Type: multipart/form-data

assetTag: LAPTOP-001
assetName: Dell Laptop
categoryId: 01HQCAT1234567890ABCDEFGH
dataMatrixImage: [valid image]
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Bad Request",
  "error": "Failed to upload data matrix image: failed to upload file to cloudinary: connection timeout"
}
```

**Frontend Display:**
```
❌ Upload Failed
Connection timeout while uploading to cloud storage

💡 Suggestion:
- Check your internet connection
- Try again in a few moments
- If problem persists, contact support
```

---

#### Invalid Image Format (Cloudinary)

**Request:**
```http
PATCH /api/users/profile
Content-Type: multipart/form-data

avatar: [unsupported image format that passed initial validation]
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Bad Request",
  "error": "Failed to upload avatar: failed to upload file to cloudinary: Invalid image file"
}
```

**Frontend Display:**
```
❌ Upload Failed
Invalid image file format

💡 Suggestion:
- Try converting your image to JPG or PNG format
- Use an image editing tool to re-save the image
- Select a different image
```

---

#### Quota Exceeded

**Request:**
```http
POST /api/assets
Content-Type: multipart/form-data

[...fields...]
dataMatrixImage: [valid image]
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Bad Request",
  "error": "Failed to upload data matrix image: failed to upload file to cloudinary: Quota exceeded"
}
```

**Frontend Display:**
```
❌ Upload Failed
Storage quota exceeded

💡 This is a server issue. Please contact the system administrator.
```

---

#### Rate Limit Exceeded

**Request:**
```http
PATCH /api/users/profile
Content-Type: multipart/form-data

avatar: [valid image]
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Bad Request",
  "error": "Failed to upload avatar: failed to upload file to cloudinary: Rate limit exceeded"
}
```

**Frontend Display:**
```
❌ Upload Failed
Too many uploads in a short time

💡 Suggestion: Please wait a few minutes and try again.
```

---

## Implementation Examples

### JavaScript/React Error Handler

```jsx
function UserProfileUpload() {
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);

  const handleUpload = async (file) => {
    setError(null);
    setLoading(true);

    const formData = new FormData();
    formData.append('avatar', file);

    try {
      const response = await fetch('/api/users/profile', {
        method: 'PATCH',
        body: formData,
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      const data = await response.json();

      if (!response.ok) {
        // Display the detailed error message
        setError(parseUploadError(data.error));
        return;
      }

      // Success
      alert('Avatar updated successfully!');
    } catch (err) {
      setError({
        title: 'Network Error',
        message: 'Could not connect to server. Please check your internet connection.',
        suggestion: 'Try again in a few moments.'
      });
    } finally {
      setLoading(false);
    }
  };

  const parseUploadError = (errorMsg) => {
    // Parse different error types
    if (errorMsg.includes('File size too large')) {
      const match = errorMsg.match(/\(([0-9.]+) MB\)/);
      const actualSize = match ? match[1] : 'unknown';
      return {
        title: 'File Too Large',
        message: `Your file (${actualSize} MB) exceeds the maximum allowed size of 5 MB.`,
        suggestion: 'Please compress your image or use a smaller file.',
        icon: '📦'
      };
    }

    if (errorMsg.includes('Invalid file type')) {
      return {
        title: 'Invalid File Type',
        message: 'The selected file is not a supported image format.',
        suggestion: 'Please use JPG, PNG, GIF, or WEBP format.',
        icon: '🖼️'
      };
    }

    if (errorMsg.includes('File is empty')) {
      return {
        title: 'Empty File',
        message: 'The selected file is empty (0 bytes).',
        suggestion: 'Please select a valid image file.',
        icon: '📄'
      };
    }

    if (errorMsg.includes('not a valid image file')) {
      return {
        title: 'Invalid Image',
        message: 'The file content does not match any supported image format.',
        suggestion: 'Make sure you are uploading an actual image file, not a renamed document.',
        icon: '⚠️'
      };
    }

    if (errorMsg.includes('Filename too long')) {
      return {
        title: 'Filename Too Long',
        message: 'The filename exceeds 255 characters.',
        suggestion: 'Please rename the file to a shorter name before uploading.',
        icon: '📝'
      };
    }

    if (errorMsg.includes('connection timeout')) {
      return {
        title: 'Connection Timeout',
        message: 'Upload timed out while connecting to cloud storage.',
        suggestion: 'Please check your internet connection and try again.',
        icon: '🌐'
      };
    }

    if (errorMsg.includes('Rate limit exceeded')) {
      return {
        title: 'Too Many Uploads',
        message: 'You have uploaded too many files in a short time.',
        suggestion: 'Please wait a few minutes and try again.',
        icon: '⏱️'
      };
    }

    if (errorMsg.includes('Quota exceeded')) {
      return {
        title: 'Storage Full',
        message: 'Server storage quota has been exceeded.',
        suggestion: 'Please contact the system administrator.',
        icon: '💾'
      };
    }

    // Generic error
    return {
      title: 'Upload Failed',
      message: errorMsg,
      suggestion: 'Please try again or contact support if the problem persists.',
      icon: '❌'
    };
  };

  return (
    <div>
      <input
        type="file"
        accept="image/*"
        onChange={(e) => handleUpload(e.target.files[0])}
        disabled={loading}
      />

      {loading && <div>Uploading...</div>}

      {error && (
        <div className="error-message">
          <h4>{error.icon} {error.title}</h4>
          <p>{error.message}</p>
          <small>💡 {error.suggestion}</small>
        </div>
      )}
    </div>
  );
}
```

---

### Flutter/Dart Error Handler

```dart
class UploadErrorParser {
  static UploadError parse(String errorMessage) {
    if (errorMessage.contains('File size too large')) {
      final regex = RegExp(r'\(([0-9.]+) MB\)');
      final match = regex.firstMatch(errorMessage);
      final actualSize = match?.group(1) ?? 'unknown';

      return UploadError(
        title: 'File Too Large',
        message: 'Your file ($actualSize MB) exceeds the maximum allowed size of 5 MB.',
        suggestion: 'Please compress your image or use a smaller file.',
        icon: '📦',
      );
    }

    if (errorMessage.contains('Invalid file type')) {
      return UploadError(
        title: 'Invalid File Type',
        message: 'The selected file is not a supported image format.',
        suggestion: 'Please use JPG, PNG, GIF, or WEBP format.',
        icon: '🖼️',
      );
    }

    if (errorMessage.contains('File is empty')) {
      return UploadError(
        title: 'Empty File',
        message: 'The selected file is empty (0 bytes).',
        suggestion: 'Please select a valid image file.',
        icon: '📄',
      );
    }

    if (errorMessage.contains('not a valid image file')) {
      return UploadError(
        title: 'Invalid Image',
        message: 'The file content does not match any supported image format.',
        suggestion: 'Make sure you are uploading an actual image file.',
        icon: '⚠️',
      );
    }

    if (errorMessage.contains('Filename too long')) {
      return UploadError(
        title: 'Filename Too Long',
        message: 'The filename exceeds 255 characters.',
        suggestion: 'Please rename the file to a shorter name.',
        icon: '📝',
      );
    }

    if (errorMessage.contains('connection timeout')) {
      return UploadError(
        title: 'Connection Timeout',
        message: 'Upload timed out while connecting to cloud storage.',
        suggestion: 'Please check your internet connection and try again.',
        icon: '🌐',
      );
    }

    if (errorMessage.contains('Rate limit exceeded')) {
      return UploadError(
        title: 'Too Many Uploads',
        message: 'You have uploaded too many files in a short time.',
        suggestion: 'Please wait a few minutes and try again.',
        icon: '⏱️',
      );
    }

    if (errorMessage.contains('Quota exceeded')) {
      return UploadError(
        title: 'Storage Full',
        message: 'Server storage quota has been exceeded.',
        suggestion: 'Please contact the system administrator.',
        icon: '💾',
      );
    }

    return UploadError(
      title: 'Upload Failed',
      message: errorMessage,
      suggestion: 'Please try again or contact support.',
      icon: '❌',
    );
  }
}

class UploadError {
  final String title;
  final String message;
  final String suggestion;
  final String icon;

  UploadError({
    required this.title,
    required this.message,
    required this.suggestion,
    required this.icon,
  });
}

// Usage in widget
Widget buildErrorDisplay(UploadError error) {
  return Container(
    padding: EdgeInsets.all(16),
    decoration: BoxDecoration(
      color: Colors.red.shade50,
      border: Border.all(color: Colors.red.shade200),
      borderRadius: BorderRadius.circular(8),
    ),
    child: Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          '${error.icon} ${error.title}',
          style: TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.bold,
            color: Colors.red.shade900,
          ),
        ),
        SizedBox(height: 8),
        Text(error.message),
        SizedBox(height: 8),
        Container(
          padding: EdgeInsets.all(8),
          decoration: BoxDecoration(
            color: Colors.blue.shade50,
            borderRadius: BorderRadius.circular(4),
          ),
          child: Row(
            children: [
              Text('💡 ', style: TextStyle(fontSize: 16)),
              Expanded(
                child: Text(
                  error.suggestion,
                  style: TextStyle(
                    fontStyle: FontStyle.italic,
                    color: Colors.blue.shade900,
                  ),
                ),
              ),
            ],
          ),
        ),
      ],
    ),
  );
}
```

---

## Summary

### Error Response Structure
```json
{
  "success": false,
  "message": "Bad Request",
  "error": "[detailed error message]"
}
```

### Error Message Format
```
[field]: [specific error] ([details]). [requirement/limit]
```

### Examples
- `avatar: File size too large (7.50 MB). Maximum allowed size is 5 MB`
- `dataMatrixImage: Invalid file type '.pdf'. Allowed types: JPG, JPEG, PNG, GIF, WEBP`
- `avatar: File is empty (0 bytes)`
- `Failed to upload avatar: failed to upload file to cloudinary: connection timeout`

### Best Practice
Always display the full error message to users - it's now designed to be user-friendly and actionable!
