# Asset Tag Generation - API Examples

## Example 1: First Asset in Category

### Request
```http
POST /api/v1/assets/generate-tag
Content-Type: application/json

{
  "categoryId": "01HXAMPLE123456789"
}
```

### Response (Category Code: LPTP, No existing assets)
```json
{
  "status": "success",
  "message": "Asset tag suggestion generated successfully",
  "data": {
    "categoryCode": "LPTP",
    "lastAssetTag": "",
    "suggestedTag": "LPTP000001",
    "nextIncrement": 1
  }
}
```

## Example 2: Subsequent Asset in Category

### Request
```http
POST /api/v1/assets/generate-tag
Content-Type: application/json

{
  "categoryId": "01HXAMPLE123456789"
}
```

### Response (Last asset: LPTP000005)
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

## Example 3: Category with Different Code

### Request
```http
POST /api/v1/assets/generate-tag
Content-Type: application/json

{
  "categoryId": "01HMONITOR987654321"
}
```

### Response (Category Code: MON, Last asset: MON000042)
```json
{
  "status": "success",
  "message": "Asset tag suggestion generated successfully",
  "data": {
    "categoryCode": "MON",
    "lastAssetTag": "MON000042",
    "suggestedTag": "MON000043",
    "nextIncrement": 43
  }
}
```

## Example 4: Complete Asset Creation Flow

### Step 1: Get Tag Suggestion
```http
POST /api/v1/assets/generate-tag
Content-Type: application/json

{
  "categoryId": "01HXAMPLE123456789"
}
```

Response:
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

### Step 2: Create Asset with Suggested Tag
```http
POST /api/v1/assets
Content-Type: application/json
Authorization: Bearer <your-token>

{
  "assetTag": "LPTP000006",
  "assetName": "Dell Inspiron 15 3000",
  "categoryId": "01HXAMPLE123456789",
  "brand": "Dell",
  "model": "Inspiron 15 3000",
  "serialNumber": "SN123456789",
  "purchaseDate": "2025-01-15",
  "purchasePrice": 15000000,
  "vendorName": "Dell Indonesia",
  "warrantyEnd": "2028-01-15",
  "status": "Active",
  "condition": "Good",
  "locationId": "01HLOCATION123456",
  "assignedTo": "01HUSER123456789"
}
```

## Example 5: Using Custom Tag (Override)

Admin decides to use custom tag instead of suggestion:

```http
POST /api/v1/assets
Content-Type: application/json
Authorization: Bearer <your-token>

{
  "assetTag": "LPTP-SPECIAL-001",
  "assetName": "MacBook Pro 16",
  "categoryId": "01HXAMPLE123456789",
  "brand": "Apple",
  "model": "MacBook Pro 16",
  "status": "Active",
  "condition": "Good"
}
```

## Example 6: Error - Category Not Found

### Request
```http
POST /api/v1/assets/generate-tag
Content-Type: application/json

{
  "categoryId": "01HINVALID000000000"
}
```

### Response
```json
{
  "status": "error",
  "message": "Category not found",
  "code": 404
}
```

## Example 7: Error - Invalid Category ID Format

### Request
```http
POST /api/v1/assets/generate-tag
Content-Type: application/json

{
  "categoryId": "invalid-id"
}
```

### Response
```json
{
  "status": "error",
  "message": "Invalid category ID format",
  "code": 400
}
```

## Using with Different Category Codes

### Desktop Computers (DSKTP)
```json
{
  "categoryCode": "DSKTP",
  "lastAssetTag": "DSKTP00015",
  "suggestedTag": "DSKTP00016",
  "nextIncrement": 16
}
```

### Monitors (MON)
```json
{
  "categoryCode": "MON",
  "lastAssetTag": "MON000099",
  "suggestedTag": "MON000100",
  "nextIncrement": 100
}
```

### Printers (PRNT)
```json
{
  "categoryCode": "PRNT",
  "lastAssetTag": "PRNT000001",
  "suggestedTag": "PRNT000002",
  "nextIncrement": 2
}
```

### Network Equipment (NET)
```json
{
  "categoryCode": "NET",
  "lastAssetTag": "",
  "suggestedTag": "NET000001",
  "nextIncrement": 1
}
```

## Frontend Integration Example (JavaScript)

```javascript
// React Hook Example
const useAssetTagGeneration = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const generateTag = async (categoryId) => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch('/api/v1/assets/generate-tag', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ categoryId }),
      });

      if (!response.ok) {
        throw new Error('Failed to generate tag');
      }

      const result = await response.json();
      return result.data;
    } catch (err) {
      setError(err.message);
      return null;
    } finally {
      setLoading(false);
    }
  };

  return { generateTag, loading, error };
};

// Usage in Component
const AssetForm = () => {
  const [categoryId, setCategoryId] = useState('');
  const [assetTag, setAssetTag] = useState('');
  const [tagSuggestion, setTagSuggestion] = useState(null);
  const { generateTag, loading } = useAssetTagGeneration();

  const handleCategoryChange = async (newCategoryId) => {
    setCategoryId(newCategoryId);

    // Auto-generate tag suggestion
    const suggestion = await generateTag(newCategoryId);
    if (suggestion) {
      setTagSuggestion(suggestion);
      setAssetTag(suggestion.suggestedTag); // Pre-fill input
    }
  };

  return (
    <form>
      <select onChange={(e) => handleCategoryChange(e.target.value)}>
        <option value="">Select Category</option>
        {/* Category options */}
      </select>

      <input
        type="text"
        value={assetTag}
        onChange={(e) => setAssetTag(e.target.value)}
        placeholder="Asset Tag"
      />

      {tagSuggestion && (
        <div className="suggestion-info">
          <p>Last tag: {tagSuggestion.lastAssetTag || 'None'}</p>
          <p>Suggested: {tagSuggestion.suggestedTag}</p>
          <button
            type="button"
            onClick={() => setAssetTag(tagSuggestion.suggestedTag)}
          >
            Use Suggested Tag
          </button>
        </div>
      )}

      {/* Other form fields */}
    </form>
  );
};
```

## cURL Testing Commands

### Generate Tag for Laptop Category
```bash
curl -X POST http://localhost:8080/api/v1/assets/generate-tag \
  -H "Content-Type: application/json" \
  -d '{
    "categoryId": "01HXAMPLE123456789"
  }'
```

### Complete Flow: Generate + Create Asset
```bash
# Step 1: Generate tag
SUGGESTED_TAG=$(curl -s -X POST http://localhost:8080/api/v1/assets/generate-tag \
  -H "Content-Type: application/json" \
  -d '{"categoryId": "01HXAMPLE123456789"}' \
  | jq -r '.data.suggestedTag')

# Step 2: Create asset with suggested tag
curl -X POST http://localhost:8080/api/v1/assets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d "{
    \"assetTag\": \"$SUGGESTED_TAG\",
    \"assetName\": \"Dell Laptop\",
    \"categoryId\": \"01HXAMPLE123456789\",
    \"status\": \"Active\",
    \"condition\": \"Good\"
  }"
```

## Postman Collection Structure

```json
{
  "info": {
    "name": "Asset Tag Generation",
    "_postman_id": "asset-tag-gen-123",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Generate Asset Tag",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"categoryId\": \"{{categoryId}}\"\n}"
        },
        "url": {
          "raw": "{{baseUrl}}/api/v1/assets/generate-tag",
          "host": ["{{baseUrl}}"],
          "path": ["api", "v1", "assets", "generate-tag"]
        }
      }
    },
    {
      "name": "Create Asset with Generated Tag",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          },
          {
            "key": "Authorization",
            "value": "Bearer {{authToken}}"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"assetTag\": \"{{generatedTag}}\",\n  \"assetName\": \"Sample Asset\",\n  \"categoryId\": \"{{categoryId}}\",\n  \"status\": \"Active\",\n  \"condition\": \"Good\"\n}"
        },
        "url": {
          "raw": "{{baseUrl}}/api/v1/assets",
          "host": ["{{baseUrl}}"],
          "path": ["api", "v1", "assets"]
        }
      }
    }
  ]
}
```
