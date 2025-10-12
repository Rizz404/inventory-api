# Asset Export API Examples

## Prerequisites
- Server running on http://localhost:8080
- Valid authentication token

## 1. Export Asset List to PDF

### Request
```http
POST http://localhost:8080/api/v1/assets/export/list
Authorization: Bearer YOUR_TOKEN_HERE
Content-Type: application/json
Accept-Language: en

{
  "format": "pdf",
  "searchQuery": "",
  "filters": {
    "status": "Active",
    "condition": "Good"
  },
  "sort": {
    "field": "assetTag",
    "order": "asc"
  },
  "includeDataMatrixImage": false
}
```

### Using cURL
```bash
curl -X POST "http://localhost:8080/api/v1/assets/export/list" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -H "Accept-Language: en" \
  -d '{
    "format": "pdf",
    "filters": {
      "status": "Active"
    },
    "sort": {
      "field": "assetTag",
      "order": "asc"
    }
  }' \
  --output asset_list.pdf
```

## 2. Export Asset List to Excel

### Request
```http
POST http://localhost:8080/api/v1/assets/export/list
Authorization: Bearer YOUR_TOKEN_HERE
Content-Type: application/json
Accept-Language: id

{
  "format": "excel",
  "searchQuery": "laptop",
  "filters": {
    "categoryId": "01JK8H7WNAE7BDYQ5C0M8HQPQT",
    "status": "Active"
  },
  "sort": {
    "field": "assetName",
    "order": "asc"
  }
}
```

### Using cURL
```bash
curl -X POST "http://localhost:8080/api/v1/assets/export/list" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -H "Accept-Language: id" \
  -d '{
    "format": "excel",
    "filters": {
      "categoryId": "01JK8H7WNAE7BDYQ5C0M8HQPQT"
    }
  }' \
  --output asset_list.xlsx
```

## 3. Export All Assets (No Filters)

### Request
```http
POST http://localhost:8080/api/v1/assets/export/list
Authorization: Bearer YOUR_TOKEN_HERE
Content-Type: application/json

{
  "format": "pdf"
}
```

### Using cURL
```bash
curl -X POST "http://localhost:8080/api/v1/assets/export/list" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{"format": "pdf"}' \
  --output all_assets.pdf
```

## 4. Export Asset Statistics

### Request
```http
GET http://localhost:8080/api/v1/assets/export/statistics
Authorization: Bearer YOUR_TOKEN_HERE
Accept-Language: en
```

### Using cURL
```bash
curl -X GET "http://localhost:8080/api/v1/assets/export/statistics" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Accept-Language: en" \
  --output asset_statistics.pdf
```

## 5. Export with Multiple Filters

### Request
```http
POST http://localhost:8080/api/v1/assets/export/list
Authorization: Bearer YOUR_TOKEN_HERE
Content-Type: application/json

{
  "format": "excel",
  "searchQuery": "computer",
  "filters": {
    "status": "Active",
    "condition": "Good",
    "categoryId": "01JK8H7WNAE7BDYQ5C0M8HQPQT",
    "locationId": "01JK8H7WNAE7BDYQ5C0M8HQPQW",
    "brand": "Dell"
  },
  "sort": {
    "field": "purchaseDate",
    "order": "desc"
  }
}
```

## Response Examples

### Success Response (Binary File)
```
HTTP/1.1 200 OK
Content-Type: application/pdf (or application/vnd.openxmlformats-officedocument.spreadsheetml.sheet)
Content-Disposition: attachment; filename=asset_list.pdf (or asset_list.xlsx)
Content-Length: [file size]

[Binary file content]
```

### Error Response
```json
{
  "success": false,
  "message": "Invalid export format",
  "data": null
}
```

## Filter Options

### Status Filter
```json
{
  "filters": {
    "status": "Active"  // Active, Maintenance, Disposed, Lost
  }
}
```

### Condition Filter
```json
{
  "filters": {
    "condition": "Good"  // Good, Fair, Poor, Damaged
  }
}
```

### Category Filter
```json
{
  "filters": {
    "categoryId": "01JK8H7WNAE7BDYQ5C0M8HQPQT"
  }
}
```

### Location Filter
```json
{
  "filters": {
    "locationId": "01JK8H7WNAE7BDYQ5C0M8HQPQW"
  }
}
```

### Assigned User Filter
```json
{
  "filters": {
    "assignedTo": "01JK8H7WNAE7BDYQ5C0M8HQPQX"
  }
}
```

### Brand/Model Filter
```json
{
  "filters": {
    "brand": "Dell",
    "model": "Latitude"
  }
}
```

## Sort Options

### Available Sort Fields
- `assetTag`
- `assetName`
- `brand`
- `model`
- `serialNumber`
- `purchaseDate`
- `purchasePrice`
- `vendorName`
- `warrantyEnd`
- `status`
- `condition`
- `createdAt`
- `updatedAt`

### Sort Order
- `asc` - Ascending
- `desc` - Descending (default)

### Sort Example
```json
{
  "sort": {
    "field": "assetName",
    "order": "asc"
  }
}
```

## Testing with Postman

1. **Import Collection**: Create a new Postman collection
2. **Set Variables**:
   - `base_url`: http://localhost:8080
   - `token`: Your authentication token
3. **Create Requests**:
   - Add requests from examples above
   - Use `{{base_url}}` and `{{token}}` variables
4. **Send & Download**: Click "Send and Download" to save the file

## Testing with HTTPie

### Export to PDF
```bash
http POST http://localhost:8080/api/v1/assets/export/list \
  Authorization:"Bearer YOUR_TOKEN" \
  format=pdf \
  --download
```

### Export to Excel
```bash
http POST http://localhost:8080/api/v1/assets/export/list \
  Authorization:"Bearer YOUR_TOKEN" \
  format=excel \
  --download
```

### Export Statistics
```bash
http GET http://localhost:8080/api/v1/assets/export/statistics \
  Authorization:"Bearer YOUR_TOKEN" \
  --download
```
