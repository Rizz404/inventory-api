# Decimal Formatting Guide

## Problem
When JSON encoding in Go, `float64` values without decimal parts (e.g., `3211.0`) are automatically serialized as integers (`3211`) to save space. This causes issues with client applications (like Flutter) that expect consistent decimal formatting.

Example issue:
```json
{
  "purchasePrice": 3211  // Expected: "3211.00" or 3211.00
}
```

## Solution
Convert all monetary values (prices, costs, etc.) to **strings with 2 decimal places** in the Response structs.

### Changes Made

#### 1. Domain Layer Updates
Updated Response structs to use `*string` instead of `*float64` for monetary fields:

**Files Modified:**
- `domain/asset.go`:
  - `AssetResponse.PurchasePrice`: `*float64` → `*string`
  - `AssetListResponse.PurchasePrice`: `*float64` → `*string`
  - `AssetValueStatisticsResponse`: All value fields → `*string`
  - `AssetSummaryStatisticsResponse.MostExpensiveAssetValue`: `*float64` → `*string`
  - `AssetSummaryStatisticsResponse.LeastExpensiveAssetValue`: `*float64` → `*string`

- `domain/maintenance_record.go`:
  - `MaintenanceRecordResponse.ActualCost`: `*float64` → `*string`
  - `MaintenanceRecordListResponse.ActualCost`: `*float64` → `*string`
  - All statistics response cost fields → `string` or `*string`

#### 2. Mapper Layer Updates
Created helper functions to format float64 values to strings with 2 decimal places:

**File:** `internal/postgresql/mapper/asset_mapper.go`
```go
// Helper function to format float64 pointer to string pointer with 2 decimal places
func formatPriceToString(price *float64) *string {
	if price == nil {
		return nil
	}
	formatted := fmt.Sprintf("%.2f", *price)
	return &formatted
}

// Helper function to format float64 to string with 2 decimal places
func formatFloat64ToString(value float64) string {
	return fmt.Sprintf("%.2f", value)
}
```

Updated mappers:
- `AssetToResponse()` - Format `PurchasePrice`
- `AssetToListResponse()` - Format `PurchasePrice`
- `AssetStatisticsToResponse()` - Format all value statistics
- `MaintenanceRecordToResponse()` - Format `ActualCost`
- `MaintenanceRecordToListResponse()` - Format `ActualCost`
- `MaintenanceRecordStatisticsToResponse()` - Format all cost statistics

#### 3. Custom Decimal Type (Optional - Not Used)
Created `domain/decimal.go` with custom types for more advanced use cases:
- `DecimalPrice`: Always formats as decimal string
- `DecimalValue`: Nullable decimal value

This is available for future use if needed.

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
// Parse string to double
double price = double.parse(asset.purchasePrice ?? "0.0");

// Or handle nullable
double? price = asset.purchasePrice != null
    ? double.parse(asset.purchasePrice!)
    : null;
```

**JavaScript:**
```javascript
// Parse string to number
const price = parseFloat(asset.purchasePrice || "0");

// Or
const price = Number(asset.purchasePrice);
```

### API Response Examples

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
  "purchasePrice": "3211.00",
  "totalValue": "150000.00",
  "averageValue": "5000.00"
}
```

### Benefits

1. ✅ **Consistent formatting**: Always 2 decimal places
2. ✅ **Client-side compatibility**: No type confusion
3. ✅ **Clear monetary values**: Explicit decimal representation
4. ✅ **Backward compatible**: Clients can easily parse strings to numbers
5. ✅ **Database unchanged**: Internal storage still uses DECIMAL/float64

### Related Files

- Domain: `domain/asset.go`, `domain/maintenance_record.go`
- Mappers: `internal/postgresql/mapper/asset_mapper.go`, `internal/postgresql/mapper/maintenance_record_mapper.go`
- Migrations: All DECIMAL fields in migrations remain unchanged

### Future Considerations

If you need to add new monetary fields:
1. Use `*float64` in domain entity structs
2. Use `*string` in response structs
3. Use `formatPriceToString()` in mapper functions
4. Document the field in this guide
