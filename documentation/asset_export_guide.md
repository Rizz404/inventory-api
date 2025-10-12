# Asset Export Feature Documentation

## Overview
Sistem export untuk daftar asset dan statistik asset dengan dukungan format PDF dan Excel.

## Features

### 1. Export Asset List
Export daftar asset ke dalam format PDF atau Excel dengan filtering dan sorting.

**Endpoint:**
```
POST /api/v1/assets/export/list
```

**Headers:**
- `Authorization`: Bearer token (required)
- `Accept-Language`: `en` atau `id` (optional, default: `en`)

**Request Body:**
```json
{
  "format": "pdf",  // "pdf" atau "excel"
  "searchQuery": "laptop",  // optional
  "filters": {  // optional
    "status": "Active",
    "condition": "Good",
    "categoryId": "category-id",
    "locationId": "location-id",
    "assignedTo": "user-id",
    "brand": "Dell",
    "model": "Latitude"
  },
  "sort": {  // optional
    "field": "assetTag",
    "order": "asc"
  },
  "includeDataMatrixImage": false  // optional, hanya untuk PDF
}
```

**Response:**
- Content-Type:
  - `application/pdf` untuk format PDF
  - `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet` untuk Excel
- Content-Disposition: `attachment; filename=asset_list.pdf` atau `asset_list.xlsx`
- Body: File binary

#### PDF Output
- **Layout**: Landscape A4
- **Content**:
  - Header dengan judul dan timestamp
  - Tabel dengan kolom: Asset Tag, Asset Name, Category, Brand, Model, Status, Condition, Location
  - Footer dengan total jumlah asset

#### Excel Output
- **Sheet**: "Assets"
- **Columns**:
  - Asset Tag
  - Asset Name
  - Category
  - Brand
  - Model
  - Serial Number
  - Purchase Date
  - Purchase Price
  - Vendor
  - Warranty End
  - Status
  - Condition
  - Location
  - Assigned To
- **Styling**:
  - Header dengan background biru dan teks bold
  - Auto-fit column width

### 2. Export Asset Statistics
Export statistik asset dalam format PDF dengan grafik visual.

**Endpoint:**
```
GET /api/v1/assets/export/statistics
```

**Headers:**
- `Authorization`: Bearer token (required)
- `Accept-Language`: `en` atau `id` (optional, default: `en`)

**Response:**
- Content-Type: `application/pdf`
- Content-Disposition: `attachment; filename=asset_statistics.pdf`
- Body: File binary

#### PDF Output
- **Layout**: Portrait A4
- **Content**:
  1. **Title Page**: "Asset Statistics Report" dengan timestamp

  2. **Status Distribution Chart**: Pie chart menampilkan distribusi status asset
     - Active
     - Maintenance
     - Disposed
     - Lost

  3. **Condition Distribution Chart**: Pie chart menampilkan distribusi kondisi asset
     - Good
     - Fair
     - Poor
     - Damaged

  4. **Summary Statistics Table**:
     - Total Assets
     - Total Categories
     - Total Locations
     - Active Assets (dengan persentase)
     - Assigned Assets (dengan persentase)
     - Total Value (jika ada)
     - Average Value (jika ada)

## Libraries Used

### 1. Excelize v2
- **Purpose**: Generate Excel files (.xlsx)
- **GitHub**: https://github.com/xuri/excelize
- **Features Used**:
  - Create new workbook
  - Add sheets
  - Set cell values
  - Apply cell styles
  - Auto-fit columns

### 2. Maroto v2
- **Purpose**: Generate PDF files
- **GitHub**: https://github.com/johnfercher/maroto
- **Features Used**:
  - Create PDF documents
  - Add rows and columns
  - Add text with styling
  - Add images (for charts)
  - Configure page layout (A4, margins, orientation)

### 3. go-echarts v2
- **Purpose**: Generate charts
- **GitHub**: https://github.com/go-echarts/go-echarts
- **Features Used**:
  - Pie charts
  - Chart rendering to HTML

## Implementation Details

### Service Layer
File: `services/asset/asset_export.go`

**Key Functions:**
1. `ExportAssetList`: Main handler untuk export daftar asset
2. `exportAssetListToPDF`: Generate PDF untuk daftar asset
3. `exportAssetListToExcel`: Generate Excel untuk daftar asset
4. `ExportAssetStatistics`: Main handler untuk export statistik
5. `exportAssetStatisticsToPDF`: Generate PDF untuk statistik dengan grafik
6. `generateStatusChart`: Generate pie chart untuk distribusi status
7. `generateConditionChart`: Generate pie chart untuk distribusi kondisi

### Repository Layer
File: `internal/postgresql/asset_repository.go`

**New Method:**
- `GetAssetsForExport`: Fetch all assets tanpa pagination untuk export

### Handler Layer
File: `internal/rest/asset_handler.go`

**New Endpoints:**
- `ExportAssetList`: POST /assets/export/list
- `ExportAssetStatistics`: GET /assets/export/statistics

### Domain Layer
File: `domain/asset.go`

**New Types:**
- `ExportFormat`: Enum untuk format export (pdf, excel)
- `ExportAssetListPayload`: Payload untuk export daftar asset
- `ExportAssetStatisticsPayload`: Payload untuk export statistik

## Usage Examples

### Export Asset List to PDF
```bash
curl -X POST http://localhost:8080/api/v1/assets/export/list \
  -H "Authorization: Bearer YOUR_TOKEN" \
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

### Export Asset List to Excel
```bash
curl -X POST http://localhost:8080/api/v1/assets/export/list \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "format": "excel",
    "filters": {
      "categoryId": "cat-123"
    }
  }' \
  --output asset_list.xlsx
```

### Export Asset Statistics
```bash
curl -X GET http://localhost:8080/api/v1/assets/export/statistics \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Accept-Language: id" \
  --output asset_statistics.pdf
```

## Error Handling

**Possible Errors:**
- `400 Bad Request`: Invalid format atau payload validation error
- `401 Unauthorized`: Missing or invalid auth token
- `500 Internal Server Error`: Error generating PDF/Excel atau database error

## Performance Considerations

1. **Export daftar asset**: Tidak ada pagination, akan fetch semua matching assets
   - Untuk dataset besar, consider menambah limit di backend

2. **Chart generation**: Charts di-generate sebagai temp HTML files
   - Files otomatis di-cleanup setelah PDF generation

3. **Memory usage**: Large exports may consume significant memory
   - Monitor untuk datasets > 10k assets

## Future Enhancements

1. ~~Add support untuk export berdasarkan date range~~
2. ~~Add more chart types (bar chart, line chart)~~
3. ~~Add email delivery option~~
4. ~~Add scheduling untuk regular exports~~
5. ~~Add compression untuk large files~~
6. ~~Add custom template support~~
