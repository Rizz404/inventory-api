# Example JSON Responses with Decimal2 Formatting

## Asset Response Example

### Before (Issue):
```json
{
  "id": "01K6MEDBANVRVZXRFQSMGXCM1J",
  "assetTag": "AST-000020",
  "assetName": "Dell Laptop",
  "purchasePrice": 3211,  ❌ No decimals!
  "status": "Active",
  "condition": "Good"
}
```

### After (Fixed):
```json
{
  "id": "01K6MEDBANVRVZXRFQSMGXCM1J",
  "assetTag": "AST-000020",
  "assetName": "Dell Laptop",
  "purchasePrice": 3211.00,  ✅ Always 2 decimals as number!
  "status": "Active",
  "condition": "Good"
}
```

## Statistics Response Example

### Before (Issue):
```json
{
  "valueStatistics": {
    "totalValue": 150000,      ❌ No decimals
    "averageValue": 5000,      ❌ No decimals
    "minValue": 1000,          ❌ No decimals
    "maxValue": 50000          ❌ No decimals
  }
}
```

### After (Fixed):
```json
{
  "valueStatistics": {
    "totalValue": 150000.00,   ✅ Consistent formatting
    "averageValue": 5000.00,   ✅ Consistent formatting
    "minValue": 1000.00,       ✅ Consistent formatting
    "maxValue": 50000.00       ✅ Consistent formatting
  }
}
```

## Maintenance Record Response Example

### Before (Issue):
```json
{
  "id": "01K6MED...",
  "maintenanceDate": "2025-10-03T07:08:34Z",
  "actualCost": 500,  ❌ No decimals
  "title": "Regular maintenance"
}
```

### After (Fixed):
```json
{
  "id": "01K6MED...",
  "maintenanceDate": "2025-10-03T07:08:34Z",
  "actualCost": 500.00,  ✅ Consistent formatting
  "title": "Regular maintenance"
}
```

## Null Values

```json
{
  "purchasePrice": null,  ✅ Properly handles null
  "actualCost": null      ✅ Properly handles null
}
```

## Flutter/Dart Usage

```dart
// Model class
class Asset {
  final String id;
  final double? purchasePrice;  // Can be null

  Asset.fromJson(Map<String, dynamic> json)
    : id = json['id'],
      purchasePrice = json['purchasePrice']?.toDouble();  // Direct conversion
}

// Usage - No parsing needed!
final asset = Asset.fromJson(jsonResponse);
print(asset.purchasePrice);  // 3211.0

// Formatting for display
final formattedPrice = '\$${asset.purchasePrice?.toStringAsFixed(2)}';
// Output: $3211.00
```

## JavaScript Usage

```javascript
// API Response
const response = await fetch('/api/assets/123');
const asset = await response.json();

// Direct number access - No parsing needed!
console.log(asset.purchasePrice);  // 3211.00
console.log(typeof asset.purchasePrice);  // "number"

// Formatting for display
const formatter = new Intl.NumberFormat('en-US', {
  style: 'currency',
  currency: 'USD',
});

console.log(formatter.format(asset.purchasePrice));  // "$3,211.00"
```

## Technical Details

### JSON Number Type
- Values are encoded as JSON **numbers**, not strings
- `3211.00` (number) ≠ `"3211.00"` (string)
- Type-safe and compatible with all JSON parsers

### Precision
- Always exactly 2 decimal places
- `3211` → `3211.00`
- `3211.5` → `3211.50`
- `3211.99` → `3211.99`

### Database to JSON Flow
```
PostgreSQL: DECIMAL(15, 2)
    ↓
Go Internal: *float64
    ↓
Mapper: toNullableDecimal2(*float64)
    ↓
Response: *NullableDecimal2
    ↓
JSON Marshaling: Custom MarshalJSON()
    ↓
JSON Output: 3211.00 (number with 2 decimals)
```
