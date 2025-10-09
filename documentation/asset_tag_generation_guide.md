# Asset Tag Generation Guide

## Overview

This feature allows admins to generate suggested asset tags based on category codes. The asset tag follows the format: `CATEGORYCODE000001`, where the numeric part is automatically incremented based on existing assets in that category.

## Why This Approach?

Instead of auto-generating asset tags in the backend automatically, this approach gives admins more control:
- Admins can **preview** the suggested tag before creating an asset
- Admins can **manually override** the suggested tag if needed
- Admins maintain **flexibility** in their asset naming conventions

## API Endpoint

### Generate Asset Tag Suggestion

**Endpoint:** `POST /api/v1/assets/generate-tag`

**Request Body:**
```json
{
  "categoryId": "01HXAMPLE123456789"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Asset tag suggestion generated successfully",
  "data": {
    "categoryCode": "LPTP",
    "lastAssetTag": "LPTP000005",
    "suggestedTag": "LPTP000006",
    "nextIncrement": 6
  }
}
```

## Response Fields

| Field           | Type    | Description                                                         |
| --------------- | ------- | ------------------------------------------------------------------- |
| `categoryCode`  | string  | The category code from the selected category                        |
| `lastAssetTag`  | string  | The last asset tag used in this category (empty if no assets exist) |
| `suggestedTag`  | string  | The suggested next asset tag (formatted with 6-digit padding)       |
| `nextIncrement` | integer | The next increment number                                           |

## Example Use Cases

### 1. First Asset in Category
If there are no assets in the "Laptop" category (code: `LPTP`):
```json
{
  "categoryCode": "LPTP",
  "lastAssetTag": "",
  "suggestedTag": "LPTP000001",
  "nextIncrement": 1
}
```

### 2. Existing Assets
If the last asset tag in "Laptop" category is `LPTP000005`:
```json
{
  "categoryCode": "LPTP",
  "lastAssetTag": "LPTP000005",
  "suggestedTag": "LPTP000006",
  "nextIncrement": 6
}
```

## Frontend Integration

### Recommended Workflow

1. **Category Selection:**
   - User selects a category when creating a new asset

2. **Auto-fetch Suggestion:**
   - Frontend automatically calls `/api/v1/assets/generate-tag` with the selected category ID
   - Display the suggested tag in the asset tag input field

3. **User Options:**
   - Accept the suggestion (use as-is)
   - Modify the suggestion (edit the field)
   - Enter completely custom tag

### Example Frontend Code (React)

```javascript
const [assetTag, setAssetTag] = useState('');
const [suggestedTag, setSuggestedTag] = useState(null);

// Called when category changes
const handleCategoryChange = async (categoryId) => {
  try {
    const response = await fetch('/api/v1/assets/generate-tag', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ categoryId })
    });

    const result = await response.json();
    setSuggestedTag(result.data);
    setAssetTag(result.data.suggestedTag); // Pre-fill the input
  } catch (error) {
    console.error('Failed to generate tag suggestion:', error);
  }
};
```

## Asset Tag Format

### Format Structure
```
CATEGORYCODE + 6-digit number
```

### Examples
- `LPTP000001` - First laptop
- `DSKTP00025` - 25th desktop (assuming category code is DSKTP)
- `MON000100` - 100th monitor

### Padding
- Numbers are always padded to **6 digits**
- Supports up to 999,999 assets per category
- Leading zeros are preserved (e.g., `000001`, `000042`)

## Implementation Details

### Database Layer
**Repository Method:** `GetLastAssetTagByCategory`
```go
// Returns the most recent asset tag for a given category
// Ordered by asset_tag DESC to get the highest increment
func (r *AssetRepository) GetLastAssetTagByCategory(ctx context.Context, categoryId string) (string, error)
```

### Service Layer
**Service Method:** `GenerateAssetTagSuggestion`
```go
// 1. Fetches category to get CategoryCode
// 2. Gets last asset tag for that category
// 3. Extracts numeric part and increments by 1
// 4. Formats new tag with 6-digit padding
func (s *Service) GenerateAssetTagSuggestion(ctx context.Context, payload *domain.GenerateAssetTagPayload) (domain.GenerateAssetTagResponse, error)
```

### Handler Layer
**Handler Method:** `GenerateAssetTagSuggestion`
```go
// Validates request payload
// Calls service to generate suggestion
// Returns formatted response
func (h *AssetHandler) GenerateAssetTagSuggestion(c *fiber.Ctx) error
```

## Error Handling

### Possible Errors

1. **Category Not Found:**
```json
{
  "status": "error",
  "message": "Category not found",
  "code": 404
}
```

2. **Invalid Category ID:**
```json
{
  "status": "error",
  "message": "Invalid category ID format",
  "code": 400
}
```

3. **Database Error:**
```json
{
  "status": "error",
  "message": "Internal server error",
  "code": 500
}
```

## Best Practices

### For Admins
1. **Use Descriptive Category Codes:**
   - Keep codes short but meaningful (3-5 characters)
   - Examples: `LPTP` (Laptop), `DSKTP` (Desktop), `MON` (Monitor)

2. **Follow Suggestions When Possible:**
   - Helps maintain consistency
   - Makes asset tracking easier
   - Simplifies reporting and analytics

3. **Custom Tags:**
   - Only use custom tags when necessary
   - Ensure uniqueness before saving
   - Document reasoning for deviations

### For Developers
1. **Always Validate:**
   - Check if asset tag already exists before creating
   - Validate format if enforcing patterns

2. **Handle Edge Cases:**
   - Empty categories (no assets yet)
   - Deleted assets (gaps in sequence)
   - Manual overrides (non-standard formats)

3. **Performance:**
   - Cache category codes if needed
   - Consider indexing `asset_tag` column
   - Use appropriate database queries

## Related Endpoints

- `POST /api/v1/assets` - Create new asset (uses the generated or custom tag)
- `GET /api/v1/assets/check/tag/:tag` - Check if asset tag exists
- `GET /api/v1/categories/:id` - Get category details

## Testing

### Manual Testing
```bash
# Generate suggestion for a category
curl -X POST http://localhost:8080/api/v1/assets/generate-tag \
  -H "Content-Type: application/json" \
  -d '{"categoryId": "01HXAMPLE123456789"}'

# Create asset with suggested tag
curl -X POST http://localhost:8080/api/v1/assets \
  -H "Content-Type: application/json" \
  -d '{
    "assetTag": "LPTP000001",
    "assetName": "Dell Laptop",
    "categoryId": "01HXAMPLE123456789",
    "status": "Active",
    "condition": "Good"
  }'
```

### Test Scenarios
1. ✅ First asset in category
2. ✅ Subsequent assets in category
3. ✅ Multiple categories with different codes
4. ✅ Custom asset tag override
5. ✅ Non-existent category (error case)

## Future Enhancements

Potential improvements for future versions:

1. **Custom Format Templates:**
   - Allow admins to define tag format patterns
   - Support different padding lengths
   - Include date/location codes

2. **Bulk Generation:**
   - Generate multiple sequential tags at once
   - Useful for bulk asset imports

3. **Tag Validation Rules:**
   - Enforce format patterns
   - Validate against organization standards
   - Warn about duplicates in real-time

4. **Analytics:**
   - Track tag usage patterns
   - Identify gaps in sequences
   - Report on tag compliance

## Changelog

### Version 1.0.0 (Current)
- Initial implementation
- Basic category-based tag generation
- 6-digit numeric padding
- REST API endpoint
