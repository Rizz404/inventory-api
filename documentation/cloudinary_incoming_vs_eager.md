# Cloudinary: Incoming vs Eager Transformations

## TL;DR

✅ **USE: Incoming Transformations** - Transform original BEFORE saving
❌ **AVOID: Eager Transformations** - Create additional derived assets

## Visual Comparison

### ❌ Eager Transformations (Boros!)

```
Upload 3MB JPEG → Cloudinary Server
                    ↓
              [Transform]
                    ↓
        ┌───────────┴───────────┐
        ↓                       ↓
  Original 3MB           Derived 200KB
  (disimpan)             (disimpan)
        ↓                       ↓
   TOTAL: 3.2 MB storage
```

**Problem:**
- Original besar masih disimpan
- Boros 90% storage
- App bisa akses original (besar) atau derived (kecil)
- Perlu logic untuk pilih URL mana yang dipakai

---

### ✅ Incoming Transformations (Efisien!)

```
Upload 3MB JPEG → Cloudinary Server
                    ↓
              [Transform]
                    ↓
              Optimized 200KB
              (disimpan sebagai original)
                    ↓
   TOTAL: 200KB storage only!
```

**Benefits:**
- Original dibuang setelah transform
- Hemat 95% storage
- App selalu dapat versi optimized
- Tidak perlu logic tambahan

## Code Comparison

### Eager (Old Way)
```go
type UploadConfig struct {
    EagerTransformations string // Creates additional file
}

uploadParams.Eager = "f_webp,q_auto"
// Result: Original + Derived = 2 files
```

### Incoming (New Way)
```go
type UploadConfig struct {
    Transformation string // Transforms original before saving
}

uploadParams.Transformation = "w_800,c_limit/f_webp,q_auto"
// Result: Only optimized = 1 file
```

## Real Example

### Upload Category Image (3000×2000px, 3.2MB JPEG)

**Eager:**
```
Storage: 3.2 MB (original) + 200 KB (derived) = 3.4 MB
URLs:
  - Original: https://.../cat-123.jpg (3.2 MB) ❌ Masih besar!
  - Derived:  https://.../cat-123.webp (200 KB) ✅ Optimized
```

**Incoming:**
```
Storage: 200 KB (optimized only)
URL:
  - Original: https://.../cat-123.webp (200 KB) ✅ Selalu optimized!
```

## When to Use Each

### ✅ Incoming Transformation
- **User-generated content** (profile pics, product images)
- **Tidak butuh original** dalam ukuran penuh
- **Hemat storage & bandwidth** (free tier!)
- **Enforce size/quality limits**

### ❌ Eager Transformation
- **Butuh original** dan derived versions
- **Professional content** (gallery, portfolio)
- **Multiple sizes** untuk responsive (thumbnail, medium, large)
- **Tidak masalah** dengan storage cost

### ⚠️ On-the-Fly Transformation (URL)
```
https://.../upload/w_800,c_limit/f_webp/cat-123.jpg
```
- Count sebagai **transformation SETIAP request**
- Boros quota kalau high traffic
- Hanya untuk testing atau low traffic

## Migration dari Eager ke Incoming

### Before (Eager)
```go
func GetCategoryImageUploadConfig() UploadConfig {
    return UploadConfig{
        EagerTransformations: "w_800,c_limit/f_webp,q_auto",
    }
}
```

### After (Incoming)
```go
func GetCategoryImageUploadConfig() UploadConfig {
    return UploadConfig{
        Transformation: "w_800,c_limit/f_webp,q_auto",
    }
}
```

### Impact
- **Storage:** -90% (3.4 MB → 200 KB per image)
- **Bandwidth:** -95% (always deliver optimized)
- **Code:** Tidak perlu ubah logic delivery URL
- **Free tier:** Lebih awet (hemat storage quota)

## Best Practices

### ✅ DO
1. Use incoming untuk user uploads
2. Set reasonable max dimensions (800px categories, 1920px assets)
3. Use WebP untuk photos
4. Use PNG untuk barcodes/QR codes
5. Monitor storage usage di dashboard

### ❌ DON'T
1. Use eager untuk user-generated content
2. Simpan original kalau tidak dibutuhkan
3. Transform on-the-fly di URL untuk production
4. Upload tanpa size limits
5. Lupa monitor quota usage

## Troubleshooting

### "File masih besar di Cloudinary!"
- Cek apakah pakai `Transformation` (incoming) bukan `Eager`
- Cek upload response, format harusnya WebP bukan JPEG
- Cek dimensi, harusnya sudah resize

### "Quality jelek setelah transform"
- Pakai `q_auto:best` instead of `q_auto`
- Atau set quality manual: `q_85`
- Trade-off: file size vs quality

### "Butuh original size"
- Jangan pakai incoming transformation
- Simpan di storage lain (S3, local)
- Atau pakai eager transformation (tapi boros)

## Conclusion

**Incoming transformation** = Transform → Save optimized → Discard original
**Eager transformation** = Save original → Transform → Save both

Untuk app dengan **ratusan/ribuan user-generated images**, incoming transformation adalah pilihan terbaik untuk **hemat storage & bandwidth** di free tier Cloudinary.
