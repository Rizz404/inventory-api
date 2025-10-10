# File Upload Error Handling - Implementation Summary

## Overview
Implementasi pesan error yang lebih detail dan informatif untuk upload file (avatar dan data matrix image) di Inventory API.

## Masalah yang Diselesaikan
Client kesulitan men-debug upload file yang gagal karena hanya menerima pesan generic "File upload failed" tanpa informasi spesifik tentang apa yang salah.

## Perubahan yang Dilakukan

### 1. **File Validator Baru** (`internal/rest/file_validator.go`)
File baru yang berisi fungsi validasi komprehensif untuk file upload.

**Fitur:**
- ✅ Validasi ukuran file dengan pesan detail (MB actual vs max)
- ✅ Validasi tipe file (extensi dan MIME type)
- ✅ Deteksi file kosong (0 bytes)
- ✅ Validasi panjang nama file (max 255 karakter)
- ✅ Verifikasi file bisa dibaca
- ✅ Verifikasi konten file menggunakan magic numbers (byte signature)

**Magic Numbers yang Divalidasi:**
- JPEG: `FF D8 FF`
- PNG: `89 50 4E 47`
- GIF: `47 49 46 38`
- WEBP: `52 49 46 46 ... 57 45 42 50`

**Format yang Didukung:**
- `.jpg` / `.jpeg`
- `.png`
- `.gif`
- `.webp`

**Batasan Ukuran:**
- Avatar: 5 MB max
- Data Matrix Image: 10 MB max

---

### 2. **User Service Updates** (`services/user/user_service.go`)

#### CreateUser Method
**Sebelum:**
```go
uploadResult, err := s.CloudinaryClient.UploadSingleFile(ctx, avatarFile, uploadConfig)
if err != nil {
    return domain.UserResponse{}, domain.ErrBadRequestWithKey(utils.ErrFileUploadFailedKey)
}
```

**Sesudah:**
```go
uploadResult, err := s.CloudinaryClient.UploadSingleFile(ctx, avatarFile, uploadConfig)
if err != nil {
    // Provide detailed error message
    errorMsg := "Failed to upload avatar: " + err.Error()
    return domain.UserResponse{}, domain.ErrBadRequest(errorMsg)
}
```

#### UpdateUser Method
Perubahan yang sama diterapkan untuk method UpdateUser.

**Benefit:**
- Error dari Cloudinary sekarang di-pass langsung ke client
- Client dapat melihat detail error (network, quota, invalid file, dll)

---

### 3. **Asset Service Updates** (`services/asset/asset_service.go`)

#### CreateAsset Method
**Sebelum:**
```go
uploadResult, err := s.CloudinaryClient.UploadSingleFile(ctx, dataMatrixImageFile, uploadConfig)
if err != nil {
    return domain.AssetResponse{}, domain.ErrBadRequestWithKey(utils.ErrFileUploadFailedKey)
}
```

**Sesudah:**
```go
uploadResult, err := s.CloudinaryClient.UploadSingleFile(ctx, dataMatrixImageFile, uploadConfig)
if err != nil {
    // Provide detailed error message
    errorMsg := "Failed to upload data matrix image: " + err.Error()
    return domain.AssetResponse{}, domain.ErrBadRequest(errorMsg)
}
```

#### UpdateAsset Method
Perubahan yang sama diterapkan untuk method UpdateAsset.

---

### 4. **User Handler Updates** (`internal/rest/user_handler.go`)

Ditambahkan validasi file di handler level sebelum dikirim ke service:

#### CreateUser Handler
```go
file, err := c.FormFile("avatar")
if err == nil {
    // Validate avatar file before processing (max 5MB)
    if validationErr := ValidateImageFile(file, "avatar", 5); validationErr != nil {
        return web.HandleError(c, domain.ErrBadRequest(FormatFileValidationError(validationErr)))
    }
    avatarFile = file
}
```

#### UpdateUser Handler
Validasi yang sama diterapkan.

#### UpdateCurrentUser Handler
Validasi yang sama diterapkan.

**Benefit:**
- Validasi terjadi sebelum upload ke Cloudinary (menghemat bandwidth & waktu)
- Error langsung dikembalikan ke client
- Mengurangi beban server Cloudinary

---

### 5. **Asset Handler Updates** (`internal/rest/asset_handler.go`)

#### CreateAsset Handler
```go
file, err := c.FormFile("dataMatrixImage")
if err == nil {
    // Validate data matrix image file before processing (max 10MB for QR/barcode images)
    if validationErr := ValidateImageFile(file, "dataMatrixImage", 10); validationErr != nil {
        return web.HandleError(c, domain.ErrBadRequest(FormatFileValidationError(validationErr)))
    }
    dataMatrixImageFile = file
}
```

#### UpdateAsset Handler
Validasi yang sama diterapkan.

---

### 6. **Dokumentasi Lengkap** (`documentation/file_upload_error_messages_guide.md`)

Dokumen komprehensif yang mencakup:
- Semua jenis error yang mungkin terjadi
- Format error message
- Penyebab dan solusi untuk setiap error
- Contoh implementasi di berbagai platform (JavaScript, React/TypeScript, Flutter/Dart)
- Best practices untuk client
- Test cases
- Troubleshooting guide

---

## Contoh Error Messages Baru

### 1. File Terlalu Besar
```json
{
  "error": "avatar: File size too large (7.50 MB). Maximum allowed size is 5 MB"
}
```

### 2. Tipe File Tidak Valid
```json
{
  "error": "avatar: Invalid file type '.pdf'. Allowed types: JPG, JPEG, PNG, GIF, WEBP"
}
```

### 3. File Kosong
```json
{
  "error": "dataMatrixImage: File is empty (0 bytes)"
}
```

### 4. Konten File Tidak Valid
```json
{
  "error": "avatar: File 'document.png' is not a valid image file. The file content does not match any supported image format"
}
```

### 5. Nama File Terlalu Panjang
```json
{
  "error": "avatar: Filename too long (287 characters). Maximum allowed is 255 characters"
}
```

### 6. Error dari Cloudinary
```json
{
  "error": "Failed to upload avatar: failed to upload file to cloudinary: connection timeout"
}
```

---

## Flow Validasi

### Handler Level (First Line of Defense)
1. File size check → Return error immediately
2. File extension check → Return error immediately
3. Empty file check → Return error immediately
4. Filename length check → Return error immediately
5. File readability check → Return error immediately
6. File content verification → Return error immediately

### Service Level (Second Line)
7. Cloudinary configuration check
8. Upload attempt to Cloudinary
9. Handle Cloudinary errors with detail

### Cloudinary Level (Final)
10. Additional format validation
11. Storage and processing

---

## Benefits untuk Client

### Sebelum
❌ "File upload failed" → Client bingung kenapa
❌ Harus trial & error
❌ Sulit debugging
❌ Banyak support request

### Sesudah
✅ "File size too large (7.50 MB). Maximum allowed is 5 MB" → Client tahu harus compress file
✅ "Invalid file type '.pdf'" → Client tahu harus pakai gambar
✅ "Failed to upload avatar: connection timeout" → Client tahu ada masalah network
✅ Error messages actionable
✅ Reduced support requests

---

## Testing

### Test Cases yang Harus Dijalankan

1. **Upload file > 5MB sebagai avatar**
   - Expected: Error dengan size detail

2. **Upload file PDF sebagai avatar**
   - Expected: Error dengan list format yang diperbolehkan

3. **Upload file 0 bytes**
   - Expected: Error "File is empty"

4. **Upload file .txt yang di-rename jadi .png**
   - Expected: Error "not a valid image file"

5. **Upload file dengan nama 300+ karakter**
   - Expected: Error "Filename too long"

6. **Test dengan Cloudinary offline**
   - Expected: Error dengan detail dari Cloudinary

---

## Endpoints yang Terpengaruh

### User Endpoints
- `POST /api/users` - Create user dengan avatar
- `PATCH /api/users/:id` - Update user avatar
- `PATCH /api/users/profile` - Update current user avatar

### Asset Endpoints
- `POST /api/assets` - Create asset dengan data matrix image
- `PATCH /api/assets/:id` - Update asset data matrix image

---

## Compatibility

### Backward Compatible
✅ Ya, perubahan ini backward compatible:
- Endpoints sama
- Request format sama
- Response structure sama (hanya error messages yang lebih detail)
- Client lama tetap bisa bekerja (hanya error messages yang lebih informatif)

### Breaking Changes
❌ Tidak ada breaking changes

---

## Recommendations untuk Client Developers

### 1. Client-Side Pre-Validation
Tambahkan validasi di client sebelum upload:
```javascript
if (file.size > 5 * 1024 * 1024) {
  alert('File too large! Max 5 MB');
  return;
}
```

### 2. Display Server Error Messages
Error messages dari server sudah user-friendly, langsung tampilkan:
```javascript
const error = await response.json();
alert(error.error); // Already detailed and helpful
```

### 3. Show Upload Progress
Untuk file besar, tampilkan progress:
```javascript
xhr.upload.onprogress = (e) => {
  const percent = (e.loaded / e.total) * 100;
  console.log(`Uploaded: ${percent}%`);
};
```

### 4. Implement Retry Logic
Untuk network errors, implementasikan retry:
```javascript
if (error.includes('timeout') || error.includes('connection')) {
  // Retry upload
}
```

---

## Monitoring & Logging

### Metrics to Track
1. Upload success rate
2. Most common error types
3. Average file sizes
4. Most common file types uploaded
5. Cloudinary error frequency

### Logging
Service layer sekarang log error details:
```go
errorMsg := "Failed to upload avatar: " + err.Error()
// This includes full Cloudinary error details
```

---

## Future Improvements

### Potential Enhancements
1. ⭐ Image format conversion (auto-convert BMP to PNG, etc.)
2. ⭐ Automatic image compression on server-side
3. ⭐ Image dimension validation
4. ⭐ Malware scanning integration
5. ⭐ Multiple file upload support
6. ⭐ Upload resume capability
7. ⭐ CDN caching optimization

### Client SDK
Consider creating a client SDK with built-in validation:
```javascript
import { InventoryAPIClient } from '@inventory/client-sdk';

const client = new InventoryAPIClient(apiKey);
// SDK automatically validates before upload
await client.users.uploadAvatar(userId, file);
```

---

## Conclusion

Perubahan ini membuat debugging upload issues jauh lebih mudah dengan:
- ✅ Error messages yang spesifik dan actionable
- ✅ Validasi di multiple levels (handler, service, cloudinary)
- ✅ Detail lengkap tentang apa yang salah dan bagaimana memperbaikinya
- ✅ Dokumentasi komprehensif untuk client developers
- ✅ Backward compatible
- ✅ Better user experience

Client sekarang dapat:
- 🎯 Memahami dengan jelas kenapa upload gagal
- 🎯 Tahu persis apa yang harus diperbaiki
- 🎯 Mengurangi trial & error
- 🎯 Mendapat feedback yang cepat

Result: **Fewer support requests, happier users, faster debugging!** 🚀
