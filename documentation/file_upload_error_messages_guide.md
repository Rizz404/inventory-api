# File Upload Error Messages Guide

## Overview
This guide explains the detailed error messages returned when file uploads fail in the Inventory API. These improvements help clients understand exactly why their upload failed and how to fix it.

## Error Message Structure

All file upload errors now provide detailed information about:
1. **What went wrong** - The specific validation or upload issue
2. **Current values** - Size, type, or other relevant file information
3. **Expected values** - Maximum size, allowed formats, etc.
4. **How to fix** - Clear guidance on what needs to be changed

## Common Error Scenarios

### 1. File Size Too Large

**Error Message Format:**
```json
{
  "error": "avatar: File size too large (7.50 MB). Maximum allowed size is 5 MB"
}
```

**Causes:**
- Avatar files exceeding 5 MB
- Data matrix images exceeding 10 MB

**How to Fix:**
- Compress the image before uploading
- Use image optimization tools
- Reduce image resolution

**File Size Limits:**
- **Avatar images**: 5 MB maximum
- **Data matrix/QR code images**: 10 MB maximum

---

### 2. Invalid File Type

**Error Message Format:**
```json
{
  "error": "avatar: Invalid file type '.pdf'. Allowed types: JPG, JPEG, PNG, GIF, WEBP"
}
```

**Causes:**
- Uploading non-image files (PDF, DOC, etc.)
- Using unsupported image formats (BMP, TIFF, etc.)

**How to Fix:**
- Convert the file to a supported format
- Use only: JPG, JPEG, PNG, GIF, or WEBP

**Supported Formats:**
- `.jpg` / `.jpeg` - JPEG images
- `.png` - PNG images
- `.gif` - GIF images (including animated)
- `.webp` - WebP images

---

### 3. Empty File

**Error Message Format:**
```json
{
  "error": "avatar: File is empty (0 bytes)"
}
```

**Causes:**
- Uploading a corrupted file
- File transfer interrupted
- File creation failed on client side

**How to Fix:**
- Ensure the file contains data before uploading
- Try selecting the file again
- Check if the file opens correctly on your device

---

### 4. Invalid File Content

**Error Message Format:**
```json
{
  "error": "dataMatrixImage: File 'document.png' is not a valid image file. The file content does not match any supported image format"
}
```

**Causes:**
- File has image extension but contains different data
- Renamed non-image file to have image extension
- Corrupted image file

**How to Fix:**
- Use an actual image file
- Re-save the image using an image editor
- Don't just rename file extensions

**Technical Details:**
The API checks the file's magic numbers (first bytes) to verify it's actually an image:
- JPEG: `FF D8 FF`
- PNG: `89 50 4E 47`
- GIF: `47 49 46 38`
- WEBP: `52 49 46 46 ... 57 45 42 50`

---

### 5. Filename Too Long

**Error Message Format:**
```json
{
  "error": "avatar: Filename too long (287 characters). Maximum allowed is 255 characters"
}
```

**Causes:**
- Very long original filename
- File path included in filename

**How to Fix:**
- Rename the file to a shorter name before uploading
- Keep filenames under 255 characters

---

### 6. File Cannot Be Read

**Error Message Format:**
```json
{
  "error": "avatar: Cannot read file: permission denied"
}
```

**Causes:**
- File permissions issues
- File locked by another process
- Corrupted file system

**How to Fix:**
- Check file permissions
- Close any programs using the file
- Try copying the file first

---

### 7. Cloudinary Upload Failure

**Error Message Format:**
```json
{
  "error": "Failed to upload avatar: failed to upload file to cloudinary: Invalid image file"
}
```

**Common Cloudinary Errors:**

#### Network Issues
```
Failed to upload avatar: failed to upload file to cloudinary: connection timeout
```
- **Cause**: Network connectivity problems
- **Fix**: Check internet connection, try again

#### Invalid API Credentials
```
Failed to upload avatar: failed to upload file to cloudinary: Invalid API key
```
- **Cause**: Cloudinary configuration issue (server-side)
- **Fix**: Contact system administrator

#### Storage Quota Exceeded
```
Failed to upload avatar: failed to upload file to cloudinary: Quota exceeded
```
- **Cause**: Cloudinary storage limit reached (server-side)
- **Fix**: Contact system administrator

#### Rate Limit Exceeded
```
Failed to upload avatar: failed to upload file to cloudinary: Rate limit exceeded
```
- **Cause**: Too many uploads in short time
- **Fix**: Wait a few minutes and try again

---

## Validation Flow

The API performs validation in the following order:

### Handler Level (First Check)
1. ✅ File size validation
2. ✅ File extension validation
3. ✅ Empty file check
4. ✅ Filename length check
5. ✅ File readability check
6. ✅ File content verification (magic numbers)

### Service Level (Second Check)
7. ✅ Cloudinary configuration check
8. ✅ File upload to Cloudinary
9. ✅ URL generation

### Cloudinary Level (Final Check)
10. ✅ Additional format validation
11. ✅ Virus scanning (if enabled)
12. ✅ Storage and processing

---

## API Endpoints Affected

### User Endpoints
- `POST /api/users` - Create user with avatar
- `PATCH /api/users/:id` - Update user avatar
- `PATCH /api/users/profile` - Update current user avatar

### Asset Endpoints
- `POST /api/assets` - Create asset with data matrix image
- `PATCH /api/assets/:id` - Update asset data matrix image

---

## Example Client Implementation

### JavaScript/Fetch Example

```javascript
async function uploadAvatar(file, userId) {
  // Pre-validate on client side
  if (file.size > 5 * 1024 * 1024) {
    alert('File too large! Maximum size is 5 MB');
    return;
  }

  const allowedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
  if (!allowedTypes.includes(file.type)) {
    alert('Invalid file type! Please use JPG, PNG, GIF, or WEBP');
    return;
  }

  const formData = new FormData();
  formData.append('avatar', file);
  // ... add other fields

  try {
    const response = await fetch(`/api/users/${userId}`, {
      method: 'PATCH',
      body: formData,
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });

    if (!response.ok) {
      const error = await response.json();
      // Display the detailed error message
      alert(`Upload failed: ${error.error || error.message}`);
      return;
    }

    const result = await response.json();
    console.log('Upload successful!', result);
  } catch (error) {
    alert(`Upload failed: ${error.message}`);
  }
}
```

### React/TypeScript Example

```typescript
import { useState } from 'react';

interface UploadError {
  error?: string;
  message?: string;
}

function AvatarUpload() {
  const [error, setError] = useState<string>('');
  const [loading, setLoading] = useState(false);

  const validateFile = (file: File): string | null => {
    // Size check
    const maxSize = 5 * 1024 * 1024; // 5 MB
    if (file.size > maxSize) {
      const sizeMB = (file.size / (1024 * 1024)).toFixed(2);
      return `File too large (${sizeMB} MB). Maximum allowed is 5 MB`;
    }

    // Type check
    const allowedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
    if (!allowedTypes.includes(file.type)) {
      return 'Invalid file type. Allowed: JPG, PNG, GIF, WEBP';
    }

    // Empty file check
    if (file.size === 0) {
      return 'File is empty';
    }

    return null;
  };

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Client-side validation
    const validationError = validateFile(file);
    if (validationError) {
      setError(validationError);
      return;
    }

    setError('');
    setLoading(true);

    const formData = new FormData();
    formData.append('avatar', file);

    try {
      const response = await fetch('/api/users/profile', {
        method: 'PATCH',
        body: formData,
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });

      if (!response.ok) {
        const errorData: UploadError = await response.json();
        setError(errorData.error || errorData.message || 'Upload failed');
        return;
      }

      // Success!
      alert('Avatar updated successfully!');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Upload failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <input
        type="file"
        accept="image/jpeg,image/png,image/gif,image/webp"
        onChange={handleFileChange}
        disabled={loading}
      />
      {error && <div style={{ color: 'red' }}>{error}</div>}
      {loading && <div>Uploading...</div>}
    </div>
  );
}
```

### Flutter/Dart Example

```dart
import 'package:http/http.dart' as http;
import 'package:mime/mime.dart';

class FileUploadService {
  static const int maxAvatarSize = 5 * 1024 * 1024; // 5 MB
  static const List<String> allowedTypes = [
    'image/jpeg',
    'image/png',
    'image/gif',
    'image/webp'
  ];

  String? validateFile(File file, String fieldName) {
    // Check file size
    final fileSize = file.lengthSync();
    if (fileSize > maxAvatarSize) {
      final sizeMB = (fileSize / (1024 * 1024)).toStringAsFixed(2);
      return '$fieldName: File size too large ($sizeMB MB). Maximum allowed is 5 MB';
    }

    // Check if empty
    if (fileSize == 0) {
      return '$fieldName: File is empty (0 bytes)';
    }

    // Check file type
    final mimeType = lookupMimeType(file.path);
    if (mimeType == null || !allowedTypes.contains(mimeType)) {
      return '$fieldName: Invalid file type. Allowed types: JPG, PNG, GIF, WEBP';
    }

    return null;
  }

  Future<Map<String, dynamic>> uploadAvatar(
    File file,
    String userId,
    String token,
  ) async {
    // Validate first
    final validationError = validateFile(file, 'avatar');
    if (validationError != null) {
      throw Exception(validationError);
    }

    try {
      var request = http.MultipartRequest(
        'PATCH',
        Uri.parse('https://api.example.com/users/$userId'),
      );

      request.headers['Authorization'] = 'Bearer $token';
      request.files.add(
        await http.MultipartFile.fromPath('avatar', file.path),
      );

      var response = await request.send();
      var responseBody = await response.stream.bytesToString();

      if (response.statusCode != 200) {
        final errorData = json.decode(responseBody);
        throw Exception(errorData['error'] ?? errorData['message'] ?? 'Upload failed');
      }

      return json.decode(responseBody);
    } catch (e) {
      rethrow;
    }
  }
}
```

---

## Best Practices for Clients

### 1. **Client-Side Pre-Validation**
Always validate files on the client side before uploading to provide immediate feedback:
- Check file size
- Check file type/extension
- Check if file is empty

### 2. **Show Clear Error Messages**
Display the server's error message directly to users - it's now detailed enough to be user-friendly.

### 3. **Provide Upload Progress**
Show upload progress for large files to improve user experience.

### 4. **Handle Network Errors**
Catch and display network-related errors gracefully:
- Connection timeout
- No internet connection
- Server unreachable

### 5. **Image Preview**
Show a preview of the selected image before uploading to let users verify they selected the correct file.

### 6. **Compression Recommendations**
Suggest image compression tools if users frequently encounter size limits:
- TinyPNG
- ImageOptim
- Squoosh
- Online image compressors

---

## Testing Upload Errors

### Test Cases

1. **Upload oversized file**
   - Create a 10MB image
   - Try to upload as avatar
   - Expect: "File size too large" error

2. **Upload wrong file type**
   - Try to upload a PDF as avatar
   - Expect: "Invalid file type" error

3. **Upload empty file**
   - Create a 0-byte file
   - Try to upload
   - Expect: "File is empty" error

4. **Upload corrupted image**
   - Rename a .txt file to .png
   - Try to upload
   - Expect: "Not a valid image file" error

5. **Upload with long filename**
   - Create file with 300+ character name
   - Try to upload
   - Expect: "Filename too long" error

---

## Troubleshooting Guide

### Client Reports: "Upload keeps failing"

**Check these things:**
1. ✅ File size (must be < 5 MB for avatars, < 10 MB for data matrix images)
2. ✅ File format (must be JPG, PNG, GIF, or WEBP)
3. ✅ File is not empty
4. ✅ File is not corrupted
5. ✅ Filename length (< 255 characters)
6. ✅ Internet connection is stable

**Ask client to:**
1. Check the exact error message received
2. Try with a different, smaller image
3. Try with a simple PNG or JPG file
4. Check if the file opens correctly on their device

### Server-Side Issues

If all client-side checks pass but uploads still fail:
1. Check Cloudinary configuration
2. Verify API credentials
3. Check Cloudinary storage quota
4. Review Cloudinary logs
5. Check network connectivity from server to Cloudinary

---

## Additional Resources

- [Cloudinary Upload API Documentation](https://cloudinary.com/documentation/image_upload_api_reference)
- [Image Optimization Best Practices](https://web.dev/fast/#optimize-your-images)
- [File Upload Security Best Practices](https://owasp.org/www-community/vulnerabilities/Unrestricted_File_Upload)

---

## Change Log

### Version 1.1 (Current)
- ✅ Added detailed file size error messages with actual vs. maximum size
- ✅ Added file type validation with supported formats list
- ✅ Added empty file detection
- ✅ Added file content verification (magic numbers)
- ✅ Added filename length validation
- ✅ Improved Cloudinary error message passthrough
- ✅ Added handler-level pre-validation

### Version 1.0 (Previous)
- ❌ Generic "File upload failed" message
- ❌ No specific error details
- ❌ Difficult to debug upload issues
