# Cloudinary Image Optimization Guide

## Overview

Semua image yang di-upload ke Cloudinary **otomatis dioptimasi** menggunakan **incoming transformations**. Original file langsung di-transform SEBELUM disimpan, jadi yang tersimpan adalah versi optimized. Ini menghemat storage & bandwidth di free tier.

## Bagaimana Cara Kerjanya?

### 1. Incoming Transformations (Bukan Eager!)
- Transformasi dilakukan **saat upload** di server Cloudinary
- Original file di-transform SEBELUM disimpan
- **Hanya versi optimized yang disimpan**, original dibuang
- **Tidak ada delay** untuk client
- **Hemat storage quota** karena tidak ada duplikasi file

### ❌ Perbedaan dengan Eager Transformations
**Eager** (yang lama):
- Original disimpan UTUH (4.73 MB)
- Derived assets dibuat (Banner 31KB, Banner wide 24KB)
- Total storage: 4.73 MB + 31 KB + 24 KB = ~4.8 MB per image
- Boros untuk ratusan records

**Incoming** (yang sekarang):
- Original di-transform DULU jadi WebP + resize
- Hanya versi optimized yang disimpan (~500 KB)
- Total storage: 500 KB per image
- **Hemat 90% storage!**

### 2. Format Conversion
**WebP Format:**
- 25-35% lebih kecil dari JPEG
- 26% lebih kecil dari PNG
- Didukung 95%+ browser modern

**Fallback:**
- Browser lama otomatis dapat format asli
- Cloudinary handle secara otomatis

### 3. Auto Quality (`q_auto`)
- Cloudinary AI analyze setiap image
- Compress tanpa menurunkan kualitas visual
- Hemat 40-80% file size

## Konfigurasi Per Type

### Avatar Images
```go
Transformation: "w_500,c_limit/f_webp,q_auto"
```
- Resize max width 500px (maintain aspect ratio)
- Convert ke WebP
- Auto quality
- Hasil: Original ~3MB → Optimized ~150KB (~95% smaller)

### Category Images
```go
Transformation: "w_800,c_limit/f_webp,q_auto"
```
- Resize max width 800px (maintain aspect ratio)
- Convert ke WebP
- Auto quality
- Hasil: Original ~2MB → Optimized ~200KB (~90% smaller)

### Asset Images
```go
Transformation: "w_1920,c_limit/f_webp,q_auto"
```
- Resize max width 1920px (Full HD)
- Convert ke WebP
- Auto quality
- Hasil: Original ~5MB → Optimized ~500KB (~90% smaller)

### QR/Barcode Images
```go
Transformation: "f_png,q_auto:best"
```
- Keep PNG (better for barcodes)
- Best quality compression
- Maintain scanability
- Hasil: Original ~500KB → Optimized ~200KB (~60% smaller)

## Free Tier Compatibility

### ✅ Supported
- WebP conversion
- Auto quality
- Resize/crop
- Format conversion
- Incoming transformations (transform original)

### 📊 Limits
- 25 monthly credits
- 1 credit = ~1000 transformations
- **Incoming transformations count 1x** (hanya sekali saat upload)
- **Eager transformations count 2x** (upload + derived asset) - TIDAK DIPAKAI

### 💡 Tips Menghemat
1. ✅ Gunakan incoming transformations (transform sekali, simpan optimized)
2. ❌ Hindari eager transformations (bikin file duplikat)
3. ❌ Hindari on-the-fly transformations di URL (count setiap request)
4. ✅ Resize original sebelum disimpan
5. ✅ Cache delivery URLs di frontend

## Performance Impact

### Upload Time
- **+100-300ms** untuk incoming transformations
- Transform dilakukan di server (tidak block client)
- Client dapat response dengan URL final

### Storage
**Before (Eager):**
- 100 images × 3MB (original) = 300MB
- 100 images × 500KB (derived) = 50MB
- **Total: 350MB**

**After (Incoming):**
- 100 images × 500KB (optimized only) = 50MB
- **Savings: 85% storage reduction**
**After:** 1000 views × 500KB = 500MB/month
**Savings:** 83% bandwidth reduction

## Real World Example

### Scenario: Category Image Upload

**Input Image:**
- Format: JPEG
- Size: 3.2 MB
- Dimensions: 3000×2000px

**Incoming Transformation:**
```go
Transformation: "w_800,c_limit/f_webp,q_auto"
```

**❌ SEBELUM (Eager - Yang Lama):**
- Original: 3.2 MB (JPEG, 3000×2000px) - disimpan
- Derived: 200 KB (WebP, 800×533px) - dibuat
- **Total storage: 3.4 MB**
- App masih deliver original 3.2 MB kalau akses URL biasa

**✅ SESUDAH (Incoming - Sekarang):**
- Original: DIBUANG
- Optimized: 180 KB (WebP, 800×533px) - disimpan sebagai original
- **Total storage: 180 KB (95% savings!)**
- App selalu deliver optimized 180 KB
- Quality: Visually identical

## Delivery URLs

Result delivery URL otomatis optimized:

**URL dari upload result:**
```
https://res.cloudinary.com/.../sigma-asset/categories/cat-123.webp
```

Yang tersimpan sudah optimized (800px, WebP, ~180KB), bukan original 3.2MB!

## Code Usage (Tidak Berubah!)

Upload sudah otomatis optimized, code tetap sama:

```go
// Category image - auto resized + WebP (incoming transformation)
uploadResult, err := cloudinaryClient.UploadSingleFile(
    ctx,
    file,
    cloudinary.GetCategoryImageUploadConfig(),
)
// uploadResult.SecureURL sudah pointing ke optimized file!
```

## Monitoring

Cek usage di Cloudinary Dashboard:
- Storage: Berapa GB terpakai
- Bandwidth: Berapa GB delivered
- Transformations: Berapa transformasi terpakai
- Credits: Sisa monthly credits

## Best Practices

1. **Upload original quality** - Let Cloudinary optimize
2. **Set max dimensions** - Hindari serve 4K untuk thumbnail
3. **Use eager transformations** - Transform sekali, deliver berkali-kali
4. **Monitor usage** - Check dashboard monthly
5. **Delete unused assets** - Clean up storage

## FAQ

**Q: Apakah WebP didukung semua browser?**
A: 95%+ browser modern. Cloudinary auto fallback untuk browser lama.

**Q: Berapa lama waktu convert?**
A: +100-300ms saat upload, tidak terasa karena proses upload sendiri butuh waktu.

**Q: Apakah free tier cukup?**
A: 25 credits = ~25,000 transformations/month. Cukup untuk MVP.

**Q: Bagaimana kalau limit habis?**
A: Upgrade ke Plus ($89/month) atau kurangi upload frequency.

**Q: Kalau butuh original yang besar gimana?**
A: Tidak bisa, original sudah dibuang. Incoming transformation untuk user-generated content yang ukuran tidak penting. Kalau butuh original, jangan pakai incoming transformation.

**Q: Apakah bisa diganti ke eager transformation?**
A: Bisa, tapi boros storage. Original 3MB tetap disimpan + derived 200KB = 3.2MB total.

**Q: Kenapa tidak pakai f_auto saja?**
A: f_auto di URL count sebagai transformation setiap request. Incoming transformation cuma count sekali saat upload.

**Q: Banner (tall) dan Banner (wide) itu apa?**
A: Itu default eager transformations dari Cloudinary (preset lama). Bisa dihapus kalau tidak terpakai.
- [Free Tier Limits](https://cloudinary.com/pricing)
