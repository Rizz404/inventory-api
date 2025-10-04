# API Filtering & Sorting Guide

Dokumentasi ini menjelaskan cara menggunakan fitur filtering dan sorting pada API untuk keperluan frontend development.

---

## Category API

### Base Endpoints
```
GET /categories
GET /categories/cursor
```

### Query Parameters

**Search & Filtering:**
- `search` - Pencarian berdasarkan category code atau category name
- `parentId` - Filter berdasarkan parent category ID
- `hasParent` - Filter kategori yang memiliki parent (`true`) atau tidak (`false`)

**Sorting:**
- `sortBy` - Field untuk sorting: `category_code`, `name`, `category_name`, `created_at`, `updated_at`
- `sortOrder` - Urutan: `asc` atau `desc` (default: `desc`)

**Pagination:**
- `limit` - Jumlah data per halaman (default: 10)
- `offset` - Offset untuk pagination (default: 0)
- `cursor` - Cursor ID untuk cursor-based pagination

### Response Format

**Offset Pagination:**
```json
{
  "success": true,
  "message": "Categories retrieved successfully",
  "data": [...],
  "pagination": {
    "total": 50,
    "limit": 10,
    "page": 1
  }
}
```

**Cursor Pagination:**
```json
{
  "success": true,
  "message": "Categories retrieved successfully",
  "data": [...],
  "pagination": {
    "nextCursor": "01HXG9876543210ZYXWVUTS",
    "hasNextPage": true,
    "limit": 10
  }
}
```

### Additional Endpoints
```
GET /categories/count?search={keyword}&hasParent={boolean}
GET /categories/statistics
GET /categories/check/{id}
GET /categories/check/code/{code}
```

---

## Location API

### Base Endpoints
```
GET /locations
GET /locations/cursor
```

### Query Parameters

**Search & Filtering:**
- `search` - Pencarian berdasarkan location code atau location name

**Sorting:**
- `sortBy` - Field untuk sorting: `location_code`, `name`, `location_name`, `building`, `floor`, `created_at`, `updated_at`
- `sortOrder` - Urutan: `asc` atau `desc` (default: `desc`)

**Pagination:**
- `limit` - Jumlah data per halaman (default: 10)
- `offset` - Offset untuk pagination (default: 0)
- `cursor` - Cursor ID untuk cursor-based pagination

### Response Format

**Offset Pagination:**
```json
{
  "success": true,
  "message": "Locations retrieved successfully",
  "data": [...],
  "pagination": {
    "total": 50,
    "limit": 10,
    "page": 1
  }
}
```

**Cursor Pagination:**
```json
{
  "success": true,
  "message": "Locations retrieved successfully",
  "data": [...],
  "pagination": {
    "nextCursor": "01HXG9876543210ZYXWVUTS",
    "hasNextPage": true,
    "limit": 10
  }
}
```

### Additional Endpoints
```
GET /locations/count?search={keyword}
GET /locations/statistics
GET /locations/check/{id}
GET /locations/check/code/{code}
```

---

## Notification API

### Base Endpoints
```
GET /notifications
GET /notifications/cursor
```

### Query Parameters

**Search & Filtering:**
- `search` - Pencarian berdasarkan title atau message
- `userId` - Filter berdasarkan user ID
- `relatedAssetId` - Filter berdasarkan related asset ID
- `type` - Filter berdasarkan tipe notifikasi: `MAINTENANCE`, `WARRANTY`, `STATUS_CHANGE`, `MOVEMENT`, `ISSUE_REPORT`
- `isRead` - Filter berdasarkan status baca (`true` atau `false`)

**Sorting:**
- `sortBy` - Field untuk sorting: `type`, `is_read`, `created_at`, `title`, `message`
- `sortOrder` - Urutan: `asc` atau `desc` (default: `desc`)

**Pagination:**
- `limit` - Jumlah data per halaman (default: 10)
- `offset` - Offset untuk pagination (default: 0)
- `cursor` - Cursor ID untuk cursor-based pagination

### Response Format

**Offset Pagination:**
```json
{
  "success": true,
  "message": "Notifications retrieved successfully",
  "data": [...],
  "pagination": {
    "total": 50,
    "limit": 10,
    "page": 1
  }
}
```

**Cursor Pagination:**
```json
{
  "success": true,
  "message": "Notifications retrieved successfully",
  "data": [...],
  "pagination": {
    "nextCursor": "01HXG9876543210ZYXWVUTS",
    "hasNextPage": true,
    "limit": 10
  }
}
```

### Additional Endpoints
```
GET /notifications/count?userId={id}&isRead={boolean}
GET /notifications/statistics
GET /notifications/check/{id}
PATCH /notifications/{id}/read
PATCH /notifications/{id}/unread
PATCH /notifications/mark-all-read
```

### Notification Types
- `MAINTENANCE` - Notifikasi terkait maintenance/pemeliharaan
- `WARRANTY` - Notifikasi terkait warranty/garansi
- `STATUS_CHANGE` - Notifikasi perubahan status asset
- `MOVEMENT` - Notifikasi perpindahan asset
- `ISSUE_REPORT` - Notifikasi laporan masalah

---

## Scan Log API

### Base Endpoints
```
GET /scan-logs
GET /scan-logs/cursor
```

### Query Parameters

**Search & Filtering:**
- `search` - Pencarian berdasarkan scanned value
- `scanMethod` - Filter berdasarkan metode scan: `DATA_MATRIX`, `MANUAL_INPUT`
- `scanResult` - Filter berdasarkan hasil scan: `Success`, `Invalid ID`, `Asset Not Found`
- `scannedBy` - Filter berdasarkan user ID yang melakukan scan
- `assetId` - Filter berdasarkan asset ID
- `dateFrom` - Filter dari tanggal tertentu (format: `YYYY-MM-DD`)
- `dateTo` - Filter sampai tanggal tertentu (format: `YYYY-MM-DD`)
- `hasCoordinates` - Filter log yang memiliki koordinat (`true`) atau tidak (`false`)

**Sorting:**
- `sortBy` - Field untuk sorting: `scan_timestamp`, `scanned_value`, `scan_method`, `scan_result`
- `sortOrder` - Urutan: `asc` atau `desc` (default: `desc`)

**Pagination:**
- `limit` - Jumlah data per halaman (default: 10)
- `offset` - Offset untuk pagination (default: 0)
- `cursor` - Cursor ID untuk cursor-based pagination

### Response Format

**Offset Pagination:**
```json
{
  "success": true,
  "message": "Scan logs retrieved successfully",
  "data": [...],
  "pagination": {
    "total": 50,
    "limit": 10,
    "page": 1
  }
}
```

**Cursor Pagination:**
```json
{
  "success": true,
  "message": "Scan logs retrieved successfully",
  "data": [...],
  "pagination": {
    "nextCursor": "01HXG9876543210ZYXWVUTS",
    "hasNextPage": true,
    "limit": 10
  }
}
```

### Additional Endpoints
```
GET /scan-logs/count?scanMethod={method}&scanResult={result}
GET /scan-logs/statistics
GET /scan-logs/check/{id}
GET /scan-logs/user/{userId}
GET /scan-logs/asset/{assetId}
```

### Scan Methods
- `DATA_MATRIX` - Scan menggunakan QR code/data matrix
- `MANUAL_INPUT` - Input manual tanpa scan

### Scan Results
- `Success` - Scan berhasil dan asset ditemukan
- `Invalid ID` - ID yang di-scan tidak valid
- `Asset Not Found` - Asset dengan ID tersebut tidak ditemukan

---

## Asset API

### Base Endpoints
```
GET /assets
GET /assets/cursor
```

### Query Parameters

**Search & Filtering:**
- `search` - Pencarian berdasarkan asset tag, asset name, brand, model, atau serial number
- `status` - Filter berdasarkan status: `Active`, `Maintenance`, `Disposed`, `Lost`
- `condition` - Filter berdasarkan kondisi: `Good`, `Fair`, `Poor`, `Damaged`
- `categoryId` - Filter berdasarkan category ID
- `locationId` - Filter berdasarkan location ID
- `assignedTo` - Filter berdasarkan user ID yang ditugaskan
- `brand` - Filter berdasarkan brand (partial match)
- `model` - Filter berdasarkan model (partial match)

**Sorting:**
- `sortBy` - Field untuk sorting: `asset_tag`, `asset_name`, `brand`, `model`, `serial_number`, `purchase_date`, `purchase_price`, `vendor_name`, `warranty_end`, `status`, `condition_status`, `created_at`, `updated_at`
- `sortOrder` - Urutan: `asc` atau `desc` (default: `desc`)

**Pagination:**
- `limit` - Jumlah data per halaman (default: 10)
- `offset` - Offset untuk pagination (default: 0)
- `cursor` - Cursor ID untuk cursor-based pagination

### Response Format

**Offset Pagination:**
```json
{
  "success": true,
  "message": "Assets retrieved successfully",
  "data": [...],
  "pagination": {
    "total": 50,
    "limit": 10,
    "page": 1
  }
}
```

**Cursor Pagination:**
```json
{
  "success": true,
  "message": "Assets retrieved successfully",
  "data": [...],
  "pagination": {
    "nextCursor": "01HXG9876543210ZYXWVUTS",
    "hasNextPage": true,
    "limit": 10
  }
}
```

### Additional Endpoints
```
GET /assets/count?status={status}&condition={condition}
GET /assets/statistics
GET /assets/check/{id}
GET /assets/check/tag/{tag}
GET /assets/check/serial/{serial}
GET /assets/tag/{tag}
```

### Asset Status
- `Active` - Asset aktif dan dapat digunakan
- `Maintenance` - Asset sedang dalam pemeliharaan
- `Disposed` - Asset sudah dibuang/dijual
- `Lost` - Asset hilang

### Asset Condition
- `Good` - Kondisi baik
- `Fair` - Kondisi cukup baik
- `Poor` - Kondisi buruk
- `Damaged` - Rusak

---

## Asset Movement API

### Base Endpoints
```
GET /asset-movements
GET /asset-movements/cursor
```

### Query Parameters

**Search & Filtering:**
- `search` - Pencarian berdasarkan notes (dari translation)
- `assetId` - Filter berdasarkan asset ID
- `fromLocationId` - Filter berdasarkan lokasi asal
- `toLocationId` - Filter berdasarkan lokasi tujuan
- `fromUserId` - Filter berdasarkan user asal (yang melepaskan asset)
- `toUserId` - Filter berdasarkan user tujuan (yang menerima asset)
- `movedBy` - Filter berdasarkan user yang melakukan perpindahan
- `dateFrom` - Filter dari tanggal tertentu (format: `YYYY-MM-DD`)
- `dateTo` - Filter sampai tanggal tertentu (format: `YYYY-MM-DD`)

**Sorting:**
- `sortBy` - Field untuk sorting: `movement_date`, `movementdate`, `created_at`, `createdat`, `updated_at`, `updatedat`
- `sortOrder` - Urutan: `asc` atau `desc` (default: `desc`)

**Pagination:**
- `limit` - Jumlah data per halaman (default: 10)
- `offset` - Offset untuk pagination (default: 0)
- `cursor` - Cursor ID untuk cursor-based pagination

### Response Format

**Offset Pagination:**
```json
{
  "success": true,
  "message": "Asset movements retrieved successfully",
  "data": [...],
  "pagination": {
    "total": 50,
    "limit": 10,
    "page": 1
  }
}
```

**Cursor Pagination:**
```json
{
  "success": true,
  "message": "Asset movements retrieved successfully",
  "data": [...],
  "pagination": {
    "nextCursor": "01HXG9876543210ZYXWVUTS",
    "hasNextPage": true,
    "limit": 10
  }
}
```

### Additional Endpoints
```
GET /asset-movements/count?assetId={id}&dateFrom={date}
GET /asset-movements/statistics
GET /asset-movements/check/{id}
GET /asset-movements/asset/{assetId}
```

### Use Cases
- **Track asset location history** - Filter by `assetId` untuk melihat riwayat perpindahan asset
- **Monitor location transfers** - Filter by `fromLocationId` atau `toLocationId` untuk melihat perpindahan dari/ke lokasi tertentu
- **Audit user assignments** - Filter by `fromUserId` atau `toUserId` untuk melihat asset yang diserahterimakan
- **Date range reports** - Gunakan `dateFrom` dan `dateTo` untuk laporan periode tertentu
- **Track who moved what** - Filter by `movedBy` untuk melihat siapa yang melakukan perpindahan

---

## Issue Report API

### Base Endpoints
```
GET /issue-reports
GET /issue-reports/cursor
```

### Query Parameters

**Search & Filtering:**
- `search` - Pencarian berdasarkan title atau description
- `assetId` - Filter berdasarkan asset ID
- `reportedBy` - Filter berdasarkan user ID yang melaporkan
- `resolvedBy` - Filter berdasarkan user ID yang menyelesaikan
- `issueType` - Filter berdasarkan tipe issue
- `priority` - Filter berdasarkan prioritas: `Low`, `Medium`, `High`, `Critical`
- `status` - Filter berdasarkan status: `Open`, `In Progress`, `Resolved`, `Closed`
- `isResolved` - Filter berdasarkan apakah sudah diselesaikan (`true`) atau belum (`false`)
- `dateFrom` - Filter dari tanggal tertentu (format: `YYYY-MM-DD`)
- `dateTo` - Filter sampai tanggal tertentu (format: `YYYY-MM-DD`)

**Sorting:**
- `sortBy` - Field untuk sorting: `reported_date`, `resolved_date`, `issue_type`, `priority`, `status`, `title`, `description`
- `sortOrder` - Urutan: `asc` atau `desc` (default: `desc`)

**Pagination:**
- `limit` - Jumlah data per halaman (default: 10)
- `offset` - Offset untuk pagination (default: 0)
- `cursor` - Cursor ID untuk cursor-based pagination

### Response Format

**Offset Pagination:**
```json
{
  "success": true,
  "message": "Issue reports retrieved successfully",
  "data": [...],
  "pagination": {
    "total": 50,
    "limit": 10,
    "page": 1
  }
}
```

**Cursor Pagination:**
```json
{
  "success": true,
  "message": "Issue reports retrieved successfully",
  "data": [...],
  "pagination": {
    "nextCursor": "01HXG9876543210ZYXWVUTS",
    "hasNextPage": true,
    "limit": 10
  }
}
```

### Additional Endpoints
```
GET /issue-reports/count?status={status}&priority={priority}
GET /issue-reports/statistics
GET /issue-reports/check/{id}
PATCH /issue-reports/{id}/resolve
PATCH /issue-reports/{id}/reopen
```

### Priority Levels
- `Low` - Prioritas rendah, tidak urgent
- `Medium` - Prioritas sedang
- `High` - Prioritas tinggi, perlu segera ditangani
- `Critical` - Prioritas kritis, harus segera ditangani

### Status Types
- `Open` - Issue baru dilaporkan, belum ditangani
- `In Progress` - Issue sedang dalam proses penanganan
- `Resolved` - Issue sudah diselesaikan
- `Closed` - Issue ditutup (selesai atau dibatalkan)

### Use Cases
- **Track unresolved issues** - Filter `isResolved=false` untuk melihat issue yang belum selesai
- **Priority management** - Filter by `priority` untuk fokus pada issue dengan prioritas tertentu
- **Asset issue history** - Filter by `assetId` untuk melihat riwayat masalah pada asset tertentu
- **User performance tracking** - Filter by `resolvedBy` untuk melihat issue yang diselesaikan oleh user tertentu
- **Periodic reports** - Gunakan `dateFrom` dan `dateTo` untuk laporan periode tertentu

---

## User API

### Base Endpoints
```
GET /users
GET /users/cursor
```

### Query Parameters

**Search & Filtering:**
- `search` - Pencarian berdasarkan name, full name, atau email
- `role` - Filter berdasarkan role: `Admin`, `Staff`, `Employee`
- `isActive` - Filter berdasarkan status aktif (`true`) atau tidak aktif (`false`)
- `employeeId` - Filter berdasarkan employee ID

**Sorting:**
- `sortBy` - Field untuk sorting: `name`, `full_name`, `email`, `role`, `employee_id`, `is_active`, `created_at`, `updated_at`
- `sortOrder` - Urutan: `asc` atau `desc` (default: `desc`)

**Pagination:**
- `limit` - Jumlah data per halaman (default: 10)
- `offset` - Offset untuk pagination (default: 0)
- `cursor` - Cursor ID untuk cursor-based pagination

### Response Format

**Offset Pagination:**
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [...],
  "pagination": {
    "total": 50,
    "limit": 10,
    "page": 1
  }
}
```

**Cursor Pagination:**
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [...],
  "pagination": {
    "nextCursor": "01HXG9876543210ZYXWVUTS",
    "hasNextPage": true,
    "limit": 10
  }
}
```

### Additional Endpoints
```
GET /users/count?role={role}&isActive={boolean}
GET /users/statistics
GET /users/check/{id}
GET /users/check/name/{name}
GET /users/check/email/{email}
GET /users/profile
GET /users/name/{name}
GET /users/email/{email}
```

### User Roles
- `Admin` - Administrator dengan akses penuh
- `Staff` - Staff dengan akses terbatas
- `Employee` - Employee dengan akses dasar

### Use Cases
- **Filter active users** - Filter `isActive=true` untuk melihat user yang aktif
- **Role-based management** - Filter by `role` untuk mengelola user berdasarkan peran
- **Search users** - Gunakan `search` untuk mencari user berdasarkan nama atau email
- **Employee lookup** - Filter by `employeeId` untuk mencari berdasarkan ID karyawan

---

## Maintenance Schedule API

### Base Endpoints
```
GET /maintenance/schedules
GET /maintenance/schedules/cursor
```

### Query Parameters

**Search & Filtering:**
- `search` - Pencarian berdasarkan title atau description
- `assetId` - Filter berdasarkan asset ID
- `maintenanceType` - Filter berdasarkan tipe maintenance: `Preventive`, `Corrective`
- `status` - Filter berdasarkan status: `Scheduled`, `Completed`, `Cancelled`
- `createdBy` - Filter berdasarkan user ID yang membuat schedule
- `fromDate` - Filter dari tanggal tertentu (format: `YYYY-MM-DD`)
- `toDate` - Filter sampai tanggal tertentu (format: `YYYY-MM-DD`)

**Sorting:**
- `sortBy` - Field untuk sorting: `scheduled_date`, `created_at`, `title`
- `sortOrder` - Urutan: `asc` atau `desc` (default: `desc`)

**Pagination:**
- `limit` - Jumlah data per halaman (default: 10)
- `offset` - Offset untuk pagination (default: 0)
- `cursor` - Cursor ID untuk cursor-based pagination

### Response Format

**Offset Pagination:**
```json
{
  "success": true,
  "message": "Maintenance schedules retrieved successfully",
  "data": [...],
  "pagination": {
    "total": 50,
    "limit": 10,
    "page": 1
  }
}
```

**Cursor Pagination:**
```json
{
  "success": true,
  "message": "Maintenance schedules retrieved successfully",
  "data": [...],
  "pagination": {
    "nextCursor": "01HXG9876543210ZYXWVUTS",
    "hasNextPage": true,
    "limit": 10
  }
}
```

### Additional Endpoints
```
GET /maintenance/schedules/count?status={status}&maintenanceType={type}
GET /maintenance/schedules/statistics
GET /maintenance/schedules/check/{id}
```

### Maintenance Types
- `Preventive` - Maintenance pencegahan yang terjadwal secara berkala
- `Corrective` - Maintenance perbaikan untuk mengatasi masalah

### Status Types
- `Scheduled` - Maintenance yang sudah dijadwalkan
- `Completed` - Maintenance yang sudah selesai dilakukan
- `Cancelled` - Maintenance yang dibatalkan

### Use Cases
- **View upcoming maintenance** - Filter `status=Scheduled` dan `fromDate=today` untuk melihat maintenance yang akan datang
- **Track preventive maintenance** - Filter `maintenanceType=Preventive` untuk fokus pada maintenance berkala
- **Asset maintenance history** - Filter by `assetId` untuk melihat riwayat maintenance suatu asset
- **Periodic reports** - Gunakan `fromDate` dan `toDate` untuk laporan periode tertentu
- **Monitor completion rate** - Filter by `status` untuk melihat maintenance yang completed vs scheduled

---

## Maintenance Record API

### Base Endpoints
```
GET /maintenance/records
GET /maintenance/records/cursor
```

### Query Parameters

**Search & Filtering:**
- `search` - Pencarian berdasarkan title atau notes
- `assetId` - Filter berdasarkan asset ID
- `scheduleId` - Filter berdasarkan maintenance schedule ID
- `performedByUser` - Filter berdasarkan user ID yang melakukan maintenance
- `vendorName` - Filter berdasarkan nama vendor (partial match)
- `fromDate` - Filter dari tanggal tertentu (format: `YYYY-MM-DD`)
- `toDate` - Filter sampai tanggal tertentu (format: `YYYY-MM-DD`)

**Sorting:**
- `sortBy` - Field untuk sorting: `maintenance_date`, `created_at`, `updated_at`, `title`
- `sortOrder` - Urutan: `asc` atau `desc` (default: `desc`)

**Pagination:**
- `limit` - Jumlah data per halaman (default: 10)
- `offset` - Offset untuk pagination (default: 0)
- `cursor` - Cursor ID untuk cursor-based pagination

### Response Format

**Offset Pagination:**
```json
{
  "success": true,
  "message": "Maintenance records retrieved successfully",
  "data": [...],
  "pagination": {
    "total": 50,
    "limit": 10,
    "page": 1
  }
}
```

**Cursor Pagination:**
```json
{
  "success": true,
  "message": "Maintenance records retrieved successfully",
  "data": [...],
  "pagination": {
    "nextCursor": "01HXG9876543210ZYXWVUTS",
    "hasNextPage": true,
    "limit": 10
  }
}
```

### Additional Endpoints
```
GET /maintenance/records/count?assetId={id}&performedByUser={userId}
GET /maintenance/records/statistics
GET /maintenance/records/check/{id}
```

### Use Cases
- **Asset maintenance history** - Filter by `assetId` untuk melihat riwayat maintenance suatu asset
- **Track user performance** - Filter by `performedByUser` untuk melihat maintenance yang dilakukan user tertentu
- **Vendor performance tracking** - Filter by `vendorName` untuk melihat maintenance yang dilakukan vendor tertentu
- **Cost analysis** - Gunakan dengan date range untuk analisis biaya maintenance
- **Periodic reports** - Gunakan `fromDate` dan `toDate` untuk laporan periode tertentu
- **Schedule completion tracking** - Filter by `scheduleId` untuk melihat record dari schedule tertentu

---

## Language Support

Semua endpoint mendukung multi-bahasa. Gunakan header `Accept-Language` untuk mendapatkan data dalam bahasa yang diinginkan.

```
Accept-Language: id
Accept-Language: en
```

---

## Notes

- Search adalah case-insensitive
- Default sorting: `reported_date DESC` untuk Issue Report, `movement_date DESC` untuk Asset Movement, `scan_timestamp DESC` untuk Scan Log, `created_at DESC` untuk yang lain
- Gunakan cursor-based pagination untuk dataset besar atau infinite scroll
- Gunakan offset-based pagination untuk pagination tradisional dengan nomor halaman
- Semua parameter adalah optional
- Format tanggal untuk `dateFrom` dan `dateTo`: `YYYY-MM-DD` (contoh: `2025-01-15`)
- Filter `brand` dan `model` di Asset API menggunakan partial match (ILIKE)
- Asset Movement track perpindahan lokasi dan assignment user
- Filter `isResolved` di Issue Report akan mengecek status: Resolved/Closed = true, Open/In Progress = false
- Search di User API mencari pada field: name, full_name, dan email
- Maintenance Schedule menggunakan `fromDate` dan `toDate` untuk filter scheduled_date
- Maintenance Record menggunakan `fromDate` dan `toDate` untuk filter maintenance_date
- Filter `vendorName` di Maintenance Record menggunakan partial match (ILIKE)
