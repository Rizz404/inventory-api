# Decimal Formatting Guide

## Problem
When JSON encoding in Go, `float64` values without decimal parts (e.g., `3211.0`) are automatically serialized as integers (`3211`) to save space. This causes issues with client applications (like Flutter) that expect consistent decimal formatting with 2 decimal places.

Example issue:
```json
{
  "purchasePrice": 3211  // Expected: 3211.00 (as number, not string)
}
```

## Solution
Use **custom types with JSON marshaling** that always format monetary values as **numbers with exactly 2 decimal places** in Response structs.

### Changes Made

#### 1. Custom Decimal Types
Created `domain/decimal.go` with custom types that marshal as **numbers** (not strings) with 2 decimal places:

```go
// Decimal2 - Non-nullable decimal with 2 places
type Decimal2 float64

func (d Decimal2) MarshalJSON() ([]byte, error) {
	formatted := strconv.FormatFloat(float64(d), 'f', 2, 64)
	return []byte(formatted), nil  // Returns as number: 3211.00
}

// NullableDecimal2 - Nullable decimal with 2 places
type NullableDecimal2 struct {
	Value Decimal2
	Valid bool
}

func (d NullableDecimal2) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte("null"), nil
	}
	return d.Value.MarshalJSON()
}
```

#### 2. Domain Layer Updates
Updated Response structs to use custom decimal types:

**Files Modified:**
- `domain/asset.go`:
  - `AssetResponse.PurchasePrice`: `*float64` → `*NullableDecimal2`
  - `AssetListResponse.PurchasePrice`: `*float64` → `*NullableDecimal2`
  - `AssetValueStatisticsResponse`: All value fields → `*NullableDecimal2`
  - `AssetSummaryStatisticsResponse.MostExpensiveAssetValue`: `*float64` → `*NullableDecimal2`
  - `AssetSummaryStatisticsResponse.LeastExpensiveAssetValue`: `*float64` → `*NullableDecimal2`

- `domain/maintenance_record.go`:
  - `MaintenanceRecordResponse.ActualCost`: `*float64` → `*NullableDecimal2`
  - `MaintenanceRecordListResponse.ActualCost`: `*float64` → `*NullableDecimal2`
  - Statistics responses: `float64` → `Decimal2` (non-nullable), `*float64` → `*NullableDecimal2` (nullable)

#### 3. Mapper Layer Updates
Created helper functions to convert float64 to custom types:

**File:** `internal/postgresql/mapper/asset_mapper.go`
```go
// Helper function to convert *float64 to *NullableDecimal2
func toNullableDecimal2(price *float64) *domain.NullableDecimal2 {
	return domain.NewNullableDecimal2(price)
}

// Helper function to convert float64 to Decimal2
func toDecimal2(value float64) domain.Decimal2 {
	return domain.Decimal2(value)
}
```

Updated mappers:
- `AssetToResponse()` - Convert `PurchasePrice`
- `AssetToListResponse()` - Convert `PurchasePrice`
- `AssetStatisticsToResponse()` - Convert all value statistics
- `MaintenanceRecordToResponse()` - Convert `ActualCost`
- `MaintenanceRecordToListResponse()` - Convert `ActualCost`
- `MaintenanceRecordStatisticsToResponse()` - Convert all cost statistics

### Important Notes

#### Internal Domain Objects
The internal `domain.Asset` and `domain.MaintenanceRecord` structs **still use `*float64`** for:
- Database operations
- Business logic calculations
- Internal processing

Only **Response structs** use `*string` for consistent JSON formatting.

#### Client-Side Handling

**Flutter/Dart:**
```dart
// Values are already numbers, no parsing needed!
double price = asset.purchasePrice ?? 0.0;

// Handle nullable
double? price = asset.purchasePrice;

// The JSON will be: {"purchasePrice": 3211.00}
// Not: {"purchasePrice": "3211.00"}
```

**JavaScript:**
```javascript
// Values are already numbers, no conversion needed!
const price = asset.purchasePrice || 0;

// Or
const price = asset.purchasePrice ?? 0;
```### API Response Examples

**Before:**
```json
{
  "purchasePrice": 3211,
  "totalValue": 150000,
  "averageValue": 5000
}
```

**After:**
```json
{
  "purchasePrice": 3211.00,
  "totalValue": 150000.00,
  "averageValue": 5000.00
}
```

**Note:** Values are returned as **JSON numbers** (not strings), ensuring type safety and consistency.

### Benefits

1. ✅ **Consistent formatting**: Always 2 decimal places as **numbers**
2. ✅ **Type safety**: Returns as JSON number, not string
3. ✅ **Client-side compatibility**: No parsing needed, direct number usage
4. ✅ **Clear monetary values**: Explicit decimal representation (3211.00 not 3211)
5. ✅ **Database unchanged**: Internal storage still uses DECIMAL/float64
6. ✅ **Flutter/Dart friendly**: Direct double assignment without parsing

### Related Files

- Domain: `domain/asset.go`, `domain/maintenance_record.go`
- Mappers: `internal/postgresql/mapper/asset_mapper.go`, `internal/postgresql/mapper/maintenance_record_mapper.go`
- Migrations: All DECIMAL fields in migrations remain unchanged

### Future Considerations

If you need to add new monetary fields:
1. Use `*float64` in domain entity structs (internal representation)
2. Use `*NullableDecimal2` in response structs (for nullable fields)
3. Use `Decimal2` in response structs (for non-nullable fields)
4. Use `toNullableDecimal2()` or `toDecimal2()` in mapper functions
5. Document the field in this guide

### Implementation Details

**Custom Type Behavior:**
- `Decimal2(3211.5)` → JSON: `3211.50`
- `Decimal2(3211.0)` → JSON: `3211.00` (NOT `3211`)
- `NullableDecimal2{Valid: false}` → JSON: `null`
- `NullableDecimal2{Value: Decimal2(3211), Valid: true}` → JSON: `3211.00`

**Database Mapping:**
```
PostgreSQL DECIMAL(15, 2) → Go *float64 → Response *NullableDecimal2 → JSON number with 2 decimals
```
