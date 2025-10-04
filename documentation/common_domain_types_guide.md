# Common Domain Types Guide

## Overview
File `domain/common.go` berisi tipe-tipe yang reusable untuk semua entity di aplikasi.

## Available Types

### 1. SortOrder
Enum untuk arah sorting yang dapat digunakan di semua entity.

```go
type SortOrder string

const (
    SortOrderAsc  SortOrder = "asc"   // Ascending: A-Z, 0-9, lama ke baru
    SortOrderDesc SortOrder = "desc"  // Descending: Z-A, 9-0, baru ke lama
)
```

### 2. PaginationOptions
Struct untuk konfigurasi pagination yang mendukung **offset-based** dan **cursor-based** pagination.

```go
type PaginationOptions struct {
    Limit  int    `json:"limit"`   // Jumlah item per halaman
    Offset int    `json:"offset"`  // Untuk offset-based pagination
    Cursor string `json:"cursor"`  // Untuk cursor-based pagination
}
```

### 3. PaginationMetadata
Struct untuk metadata response pagination.

```go
type PaginationMetadata struct {
    Total       int    `json:"total,omitempty"`       // Total items (offset-based)
    Limit       int    `json:"limit"`                 // Items per page
    Offset      int    `json:"offset,omitempty"`      // Current offset
    Page        int    `json:"page,omitempty"`        // Current page number
    HasNextPage bool   `json:"hasNextPage,omitempty"` // Apakah ada halaman berikutnya
    NextCursor  string `json:"nextCursor,omitempty"`  // Cursor untuk halaman berikutnya
}
```

## Usage Pattern for New Entity

### Step 1: Define Entity-Specific Sort Fields
Setiap entity memiliki sort fields yang spesifik. Contoh untuk entity `Category`:

```go
// domain/category.go

type CategorySortField string

const (
    CategorySortByName        CategorySortField = "name"
    CategorySortByDescription CategorySortField = "description"
    CategorySortByCreatedAt   CategorySortField = "created_at"
    CategorySortByUpdatedAt   CategorySortField = "updated_at"
)
```

### Step 2: Define Sort Options Struct
Gunakan `SortOrder` yang sudah ada di `common.go`:

```go
type CategorySortOptions struct {
    Field CategorySortField `json:"field" example:"created_at"`
    Order SortOrder          `json:"order" example:"desc"`  // ← Reusable!
}
```

### Step 3: Define Filter Options (Entity-Specific)
Setiap entity punya filter yang berbeda:

```go
type CategoryFilterOptions struct {
    IsActive *bool   `json:"is_active,omitempty"`
    ParentID *string `json:"parent_id,omitempty"`
}
```

### Step 4: Define Params Struct
Gunakan `PaginationOptions` yang sudah ada di `common.go`:

```go
type CategoryParams struct {
    SearchQuery *string                `json:"searchQuery,omitempty"`
    Filters     *CategoryFilterOptions `json:"filters,omitempty"`
    Sort        *CategorySortOptions   `json:"sort,omitempty"`
    Pagination  *PaginationOptions     `json:"pagination,omitempty"`  // ← Reusable!
}
```

## Complete Example: Asset Entity

```go
// domain/asset.go

// Asset-specific sort fields
type AssetSortField string

const (
    AssetSortByName        AssetSortField = "name"
    AssetSortBySerialNo    AssetSortField = "serial_number"
    AssetSortByPurchaseDate AssetSortField = "purchase_date"
    AssetSortByValue       AssetSortField = "value"
    AssetSortByStatus      AssetSortField = "status"
    AssetSortByCreatedAt   AssetSortField = "created_at"
    AssetSortByUpdatedAt   AssetSortField = "updated_at"
)

// Asset-specific sort options (using common SortOrder)
type AssetSortOptions struct {
    Field AssetSortField `json:"field" example:"created_at"`
    Order SortOrder       `json:"order" example:"desc"`  // ← From common.go
}

// Asset-specific filters
type AssetFilterOptions struct {
    Status     *AssetStatus `json:"status,omitempty"`
    CategoryID *string      `json:"category_id,omitempty"`
    LocationID *string      `json:"location_id,omitempty"`
    MinValue   *float64     `json:"min_value,omitempty"`
    MaxValue   *float64     `json:"max_value,omitempty"`
}

// Asset query params (using common PaginationOptions)
type AssetParams struct {
    SearchQuery *string             `json:"searchQuery,omitempty"`
    Filters     *AssetFilterOptions `json:"filters,omitempty"`
    Sort        *AssetSortOptions   `json:"sort,omitempty"`
    Pagination  *PaginationOptions  `json:"pagination,omitempty"`  // ← From common.go
}
```

## Benefits

### ✅ **Consistency**
- Semua entity menggunakan pagination dan sort order yang sama
- Client-side cukup belajar satu pattern untuk semua endpoint

### ✅ **Maintainability**
- Jika ada perubahan di pagination, cukup update satu tempat
- Mengurangi duplikasi code

### ✅ **Type Safety**
- Compiler akan membantu mendeteksi error
- IDE autocomplete bekerja dengan baik

### ✅ **Scalability**
- Mudah menambah entity baru dengan pattern yang sama
- Consistent API design

## What's Reusable vs Entity-Specific

### Reusable (in `common.go`):
- ✅ `SortOrder` - Semua entity menggunakan `asc`/`desc`
- ✅ `PaginationOptions` - Semua entity menggunakan limit/offset/cursor
- ✅ `PaginationMetadata` - Format response metadata sama untuk semua

### Entity-Specific (in respective domain file):
- ❌ `SortField` - Setiap entity punya fields berbeda
- ❌ `FilterOptions` - Setiap entity punya filter berbeda
- ❌ `SortOptions` - Struct wrapper untuk sort field + order
- ❌ `Params` - Kombinasi search, filter, sort, pagination

## Migration Checklist for Existing Entities

Jika ingin migrate entity lain (Category, Location, Asset, dll):

1. [ ] Hapus `SortOrder` enum dari entity file (sudah ada di `common.go`)
2. [ ] Ganti `EntityPaginationOptions` menjadi `PaginationOptions`
3. [ ] Update semua reference di handler
4. [ ] Update semua reference di repository
5. [ ] Update semua reference di service
6. [ ] Run tests untuk memastikan tidak ada breaking changes

## API Usage Example

### Request
```http
GET /api/categories?search=office&sortBy=name&sortOrder=asc&limit=20&offset=0
```

### Response (Offset-based)
```json
{
  "data": [...],
  "pagination": {
    "total": 150,
    "limit": 20,
    "offset": 0,
    "page": 1
  }
}
```

### Request (Cursor-based)
```http
GET /api/categories?limit=20&cursor=01ARZ3NDEKTSV4RRFFQ69G5FAV
```

### Response (Cursor-based)
```json
{
  "data": [...],
  "pagination": {
    "limit": 20,
    "hasNextPage": true,
    "nextCursor": "01ARZ3NDEKTSVX9YFFQ69G5FAV"
  }
}
```
