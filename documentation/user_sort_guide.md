# User Sort Options Guide

## Overview
Panduan ini menjelaskan cara menggunakan sort options untuk endpoint User API.

## Available Sort Fields

Sort fields yang tersedia didefinisikan sebagai enum `UserSortField`:

| Field                  | Value         | Description                                  |
| ---------------------- | ------------- | -------------------------------------------- |
| `UserSortByName`       | `name`        | Sort berdasarkan username                    |
| `UserSortByFullName`   | `full_name`   | Sort berdasarkan nama lengkap                |
| `UserSortByEmail`      | `email`       | Sort berdasarkan email                       |
| `UserSortByRole`       | `role`        | Sort berdasarkan role (Admin/Staff/Employee) |
| `UserSortByEmployeeID` | `employee_id` | Sort berdasarkan employee ID                 |
| `UserSortByIsActive`   | `is_active`   | Sort berdasarkan status aktif                |
| `UserSortByCreatedAt`  | `created_at`  | Sort berdasarkan tanggal dibuat              |
| `UserSortByUpdatedAt`  | `updated_at`  | Sort berdasarkan tanggal update              |

## Sort Order

Sort order yang tersedia didefinisikan sebagai enum `SortOrder` (di `domain/common.go` - reusable untuk semua entity):

| Order           | Value  | Description                         |
| --------------- | ------ | ----------------------------------- |
| `SortOrderAsc`  | `asc`  | Ascending (A-Z, 0-9, lama ke baru)  |
| `SortOrderDesc` | `desc` | Descending (Z-A, 9-0, baru ke lama) |

> **Note**: `SortOrder` adalah tipe reusable yang didefinisikan di `domain/common.go` dan dapat digunakan oleh semua entity.

## Usage Examples

### 1. Sort by Created Date (Default)
```
GET /api/users?sortBy=created_at&sortOrder=desc
```
Response akan menampilkan user terbaru di atas.

### 2. Sort by Name Ascending
```
GET /api/users?sortBy=name&sortOrder=asc
```
Response akan menampilkan user berdasarkan nama A-Z.

### 3. Sort by Role Descending
```
GET /api/users?sortBy=role&sortOrder=desc
```
Response akan menampilkan user berdasarkan role secara descending.

### 4. Sort with Filters
```
GET /api/users?sortBy=full_name&sortOrder=asc&role=Admin&isActive=true
```
Response akan menampilkan user Admin yang aktif, diurutkan berdasarkan nama lengkap A-Z.

### 5. Sort with Search
```
GET /api/users?search=john&sortBy=created_at&sortOrder=desc
```
Response akan menampilkan user yang mengandung "john" di name/full_name, diurutkan dari yang terbaru.

## Default Behavior

Jika `sortBy` tidak diberikan, sistem akan menggunakan default sorting:
- **Field**: `created_at`
- **Order**: `desc` (terbaru di atas)

## Type Safety

Dengan menggunakan enum, Anda mendapatkan:
1. **Autocompletion** - IDE akan menyarankan field yang valid
2. **Type Safety** - Compiler akan mendeteksi jika menggunakan field yang tidak valid
3. **Clear Documentation** - Semua field sort tersedia di satu tempat

## Code Example (Go Client)

```go
// Contoh penggunaan di code
params := domain.UserParams{
    Sort: &domain.UserSortOptions{
        Field: domain.UserSortByFullName,
        Order: domain.SortOrderAsc,
    },
}
```

## Pagination Options

`PaginationOptions` juga merupakan tipe reusable yang didefinisikan di `domain/common.go`:

```go
type PaginationOptions struct {
    Limit  int    `json:"limit"`   // Items per page
    Offset int    `json:"offset"`  // For offset-based pagination
    Cursor string `json:"cursor"`  // For cursor-based pagination
}
```

Lihat [Common Domain Types Guide](./common_domain_types_guide.md) untuk detail lebih lanjut.

## Notes

- Sort field yang tidak valid akan di-fallback ke `created_at desc`
- Sort order yang tidak valid akan di-fallback ke `desc`
- Kombinasi sort dengan filter dan pagination tetap didukung
- `SortOrder` dan `PaginationOptions` adalah common types yang reusable untuk semua entity
