# Decimal Formatting Fix Summary

## Problem
Go's JSON marshaler automatically removes trailing zeros from float64 values. For example:
- `150000.00` becomes `150000` (appears as int to frontend)
- `38.095238095238095` becomes `38.095238095238095` (too many decimal places)

This causes two issues:
1. **Frontend type confusion**: Values like `150000` appear as integers instead of floats/doubles
2. **Excessive precision**: Percentages and averages display too many decimal places

## Solution
Created a custom `Decimal2` type that always formats numbers with exactly 2 decimal places during JSON marshaling.

## Changes Made

### 1. Enhanced `domain/decimal.go`
- Added `NewDecimal2()` helper function
- Ensured `Decimal2` type always marshals with exactly 2 decimal places
- Kept `NullableDecimal2` for nullable fields

### 2. Updated Response Types

#### `domain/asset.go`
Changed `AssetSummaryStatisticsResponse` fields from `float64` to `Decimal2`:
- All percentage fields (activeAssetsPercentage, maintenanceAssetsPercentage, etc.)
- AverageAssetsPerDay
- Kept MostExpensiveAssetValue and LeastExpensiveAssetValue as `*NullableDecimal2`

#### `domain/category.go`
Changed `CategorySummaryStatisticsResponse` fields from `float64` to `Decimal2`:
- TopLevelPercentage
- SubCategoriesPercentage
- AverageCategoriesPerDay

#### `domain/location.go`
Changed multiple fields to use decimal types:
- `LocationSummaryStatisticsResponse`: CoordinatesPercentage, BuildingPercentage, FloorPercentage, AverageLocationsPerDay
- `GeographicStatisticsResponse`: AverageLatitude, AverageLongitude (as `*NullableDecimal2`)

### 3. Updated Mappers

#### `internal/postgresql/mapper/asset_mapper.go`
- Wrapped all percentage and average values with `toDecimal2()`
- Already using `toNullableDecimal2()` for nullable price fields

#### `internal/postgresql/mapper/category_mapper.go`
- Wrapped percentage and average fields with `domain.NewDecimal2()`

#### `internal/postgresql/mapper/location_mapper.go`
- Wrapped percentage and average fields with `domain.NewDecimal2()`
- Wrapped nullable coordinate fields with `domain.NewNullableDecimal2()`

## Results

### Before Fix
```json
{
  "purchasePrice": 150000,  // ❌ Appears as int
  "valueStatistics": {
    "totalValue": 178943,
    "averageValue": 8521.1,
    "minValue": 63,
    "maxValue": 150000
  },
  "summary": {
    "activeAssetsPercentage": 38.095238095238095,  // ❌ Too many decimals
    "maintenanceAssetsPercentage": 33.33333333333333,
    "averageAssetsPerDay": 2.520704428056228
  }
}
```

### After Fix
```json
{
  "purchasePrice": 150000.00,  // ✅ Clearly a float
  "valueStatistics": {
    "totalValue": 178943.00,
    "averageValue": 8521.10,
    "minValue": 63.00,
    "maxValue": 150000.00
  },
  "summary": {
    "activeAssetsPercentage": 38.10,  // ✅ Clean 2 decimal places
    "maintenanceAssetsPercentage": 33.33,
    "averageAssetsPerDay": 2.52
  }
}
```

## Benefits
1. **Frontend type safety**: All numeric values are clearly floats/doubles
2. **Clean display**: Percentages and averages show exactly 2 decimal places
3. **Consistent formatting**: All monetary and statistical values use the same precision
4. **User-friendly**: Ready for direct display without frontend formatting

## Implementation Notes
- Internal domain types (non-Response) still use `float64` for calculations
- Conversion to `Decimal2` happens only in Response types during mapping
- The `Decimal2` type marshals as a JSON number (not string), maintaining type compatibility
