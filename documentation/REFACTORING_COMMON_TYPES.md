# Refactoring Summary: Common Domain Types

## Date
October 4, 2025

## Changes Made

### 1. Created `domain/common.go`
File baru yang berisi tipe-tipe reusable untuk semua entity:
- `SortOrder` enum (asc/desc)
- `PaginationOptions` struct (limit, offset, cursor)
- `PaginationMetadata` struct (untuk response metadata)

### 2. Updated `domain/user.go`
- ❌ Removed: `SortOrder` enum (moved to common.go)
- ❌ Removed: `UserPaginationOptions` struct (replaced with common PaginationOptions)
- ✅ Updated: `UserParams` now uses `*PaginationOptions` instead of `*UserPaginationOptions`
- ✅ Kept: `UserSortField` (entity-specific, as expected)

### 3. Updated `internal/rest/user_handler.go`
- Changed `domain.UserPaginationOptions` → `domain.PaginationOptions` (2 occurrences)
  - In `GetUsersPaginated()` method
  - In `GetUsersCursor()` method

### 4. Updated `seeders/seeder_manager.go`
- Changed `domain.UserPaginationOptions` → `domain.PaginationOptions`
  - In `getUserIDs()` method

### 5. Updated Documentation
- ✅ Created: `documentation/common_domain_types_guide.md`
  - Complete guide untuk menggunakan common types
  - Pattern untuk entity baru
  - Migration checklist untuk entity existing
  - Examples dan best practices

- ✅ Updated: `documentation/user_sort_guide.md`
  - Added note bahwa `SortOrder` adalah common type
  - Added section tentang `PaginationOptions`
  - Cross-reference ke common domain types guide

## Benefits

### 🎯 Consistency
- Semua entity akan menggunakan pagination dan sort order yang sama
- Client hanya perlu belajar satu pattern

### 🔧 Maintainability
- Single source of truth untuk common types
- Perubahan di satu tempat berlaku untuk semua entity

### 🚀 Scalability
- Pattern yang jelas untuk menambah entity baru
- Reusable components mengurangi boilerplate

### 💡 Developer Experience
- Type safety dengan enums
- IDE autocomplete untuk valid values
- Clear documentation

## What's Reusable

✅ **In `common.go`:**
- `SortOrder` - Direction (asc/desc)
- `PaginationOptions` - Limit, offset, cursor
- `PaginationMetadata` - Response metadata

❌ **Still Entity-Specific:**
- `EntitySortField` - Setiap entity punya fields berbeda
- `EntityFilterOptions` - Setiap entity punya filter berbeda
- `EntitySortOptions` - Wrapper untuk sort field + order
- `EntityParams` - Kombinasi search, filter, sort, pagination

## Next Steps for Other Entities

Untuk migrate entity lain (Category, Location, Asset, dll):

1. Import `PaginationOptions` dari common instead of creating own
2. Remove `SortOrder` enum from entity file
3. Update all references in handler/service/repository
4. Follow pattern documented in `common_domain_types_guide.md`

## Testing

✅ All files compile without errors
✅ No breaking changes to existing functionality
✅ Pattern is ready to be applied to other entities

## Files Modified

```
✨ New:
- domain/common.go
- documentation/common_domain_types_guide.md

📝 Modified:
- domain/user.go
- internal/rest/user_handler.go
- seeders/seeder_manager.go
- documentation/user_sort_guide.md
```

## Example for Future Entity

```go
// domain/category.go

type CategorySortField string
const (
    CategorySortByName      CategorySortField = "name"
    CategorySortByCreatedAt CategorySortField = "created_at"
)

type CategorySortOptions struct {
    Field CategorySortField `json:"field"`
    Order SortOrder          `json:"order"` // ← From common.go
}

type CategoryParams struct {
    SearchQuery *string                `json:"searchQuery,omitempty"`
    Filters     *CategoryFilterOptions `json:"filters,omitempty"`
    Sort        *CategorySortOptions   `json:"sort,omitempty"`
    Pagination  *PaginationOptions     `json:"pagination,omitempty"` // ← From common.go
}
```
