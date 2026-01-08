# Bulk Asset Creation Workflow

## 📋 Overview
Fitur copy asset untuk membuat banyak aset sekaligus dengan data yang mirip. Berguna saat membeli item dalam jumlah banyak dengan spesifikasi sama (contoh: 25 kursi kantor).

## 🔄 Workflow - 3 Step Process

### **Step 1: Generate Bulk Asset Tags**
Dapatkan sequential asset tags untuk semua aset yang akan dibuat.

**Endpoint:** `POST /assets/generate-bulk-tags`

**Request:**
```json
{
  "categoryId": "01JKPT8XXXXXXXXXXX",
  "quantity": 25
}
```

**Response:**
```json
{
  "success": true,
  "message": "Bulk asset tags generated successfully",
  "data": {
    "categoryCode": "FURN",
    "lastAssetTag": "FURN-00100",
    "startTag": "FURN-00101",
    "endTag": "FURN-00125",
    "tags": [
      "FURN-00101",
      "FURN-00102",
      "...truncated...",
      "FURN-00125"
    ],
    "quantity": 25,
    "startIncrement": 101,
    "endIncrement": 125
  }
}
```

---

### **Step 2: Upload Bulk Data Matrix Images**
Upload semua gambar data matrix yang sudah di-generate di mobile.

**Endpoint:** `POST /assets/upload/bulk-datamatrix`

**Content-Type:** `multipart/form-data`

**Form Fields:**
- `assetTags[]` - Array of asset tags (string array)
- `dataMatrixImages` - Multiple image files

**Example using JavaScript (FormData):**
```javascript
const formData = new FormData();

// Add asset tags
assetTags.forEach(tag => {
  formData.append('assetTags', tag);
});

// Add image files
imageFiles.forEach(file => {
  formData.append('dataMatrixImages', file);
});

fetch('/assets/upload/bulk-datamatrix', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer YOUR_TOKEN'
  },
  body: formData
});
```

**Response:**
```json
{
  "success": true,
  "message": "Bulk data matrix images uploaded successfully",
  "data": {
    "urls": [
      "https://res.cloudinary.com/.../FURN-00101.png",
      "https://res.cloudinary.com/.../FURN-00102.png",
      "...truncated...",
      "https://res.cloudinary.com/.../FURN-00125.png"
    ],
    "count": 25,
    "assetTags": ["FURN-00101", "FURN-00102", "...", "FURN-00125"]
  }
}
```

---

### **Step 3: Bulk Create Assets**
Create semua aset dengan tags dan URLs yang sudah didapat.

**Endpoint:** `POST /assets/bulk`

**Request:**
```json
{
  "assets": [
    {
      "assetTag": "FURN-00101",
      "dataMatrixImageUrl": "https://res.cloudinary.com/.../FURN-00101.png",
      "assetName": "Kursi Kantor Executive",
      "categoryId": "01JKPT8XXXXXXXXXXX",
      "brand": "Ergotec",
      "model": "EX-500",
      "serialNumber": "SN-001",
      "purchaseDate": "2026-01-08",
      "purchasePrice": 2500000,
      "vendorName": "PT Furniture Indonesia",
      "warrantyEnd": "2027-01-08",
      "status": "Active",
      "condition": "Good",
      "locationId": "01JKPT9XXXXXXXXXXX",
      "assignedTo": null
    },
    {
      "assetTag": "FURN-00102",
      "dataMatrixImageUrl": "https://res.cloudinary.com/.../FURN-00102.png",
      "assetName": "Kursi Kantor Executive",
      "categoryId": "01JKPT8XXXXXXXXXXX",
      "brand": "Ergotec",
      "model": "EX-500",
      "serialNumber": "SN-002",
      "purchaseDate": "2026-01-08",
      "purchasePrice": 2500000,
      "vendorName": "PT Furniture Indonesia",
      "warrantyEnd": "2027-01-08",
      "status": "Active",
      "condition": "Good",
      "locationId": "01JKPT9XXXXXXXXXXX",
      "assignedTo": null
    }
    // ... 23 more assets
  ]
}
```

**Response:**
```json
{
  "success": true,
  "message": "Assets bulk created successfully",
  "data": {
    "assets": [
      {
        "id": "01JKPTA1XXXXXXXXXX",
        "assetTag": "FURN-00101",
        "dataMatrixImageUrl": "https://res.cloudinary.com/.../FURN-00101.png",
        "assetName": "Kursi Kantor Executive",
        // ... full asset data
      }
      // ... 24 more assets
    ]
  }
}
```

---

## 🎯 Key Points

### ✅ **Validations:**
- Maximum 100 assets per bulk operation
- Asset tags must be unique (checked against DB)
- Serial numbers must be unique (if provided)
- Number of images must match number of asset tags
- Each image max 10MB

### 🔄 **Data yang Berbeda:**
- `assetTag` - **WAJIB** beda untuk setiap asset
- `dataMatrixImageUrl` - **WAJIB** beda untuk setiap asset
- `serialNumber` - **OPTIONAL**, bisa null atau beda untuk setiap asset
- `createdAt`, `updatedAt` - Auto-generated oleh sistem

### 🔁 **Data yang Sama:**
- `assetName`, `categoryId`, `brand`, `model`, `purchaseDate`, `purchasePrice`, `vendorName`, `warrantyEnd`, `status`, `condition`, `locationId`, `assignedTo` - Semua bisa sama untuk semua aset

---

## 📱 Mobile Implementation Example

```dart
// 1. Generate tags
final tagsResponse = await http.post(
  Uri.parse('$baseUrl/assets/generate-bulk-tags'),
  body: json.encode({
    'categoryId': categoryId,
    'quantity': copyCount,
  }),
);
final tags = tagsResponse.data['tags'] as List<String>;

// 2. Generate QR images for each tag
List<File> qrImages = [];
for (String tag in tags) {
  final qrImage = await generateDataMatrixImage(tag);
  qrImages.add(qrImage);
}

// 3. Upload images
final formData = FormData();
for (String tag in tags) {
  formData.fields.add(MapEntry('assetTags', tag));
}
for (File image in qrImages) {
  formData.files.add(
    MapEntry('dataMatrixImages', await MultipartFile.fromFile(image.path))
  );
}
final uploadResponse = await dio.post(
  '$baseUrl/assets/upload/bulk-datamatrix',
  data: formData,
);
final imageUrls = uploadResponse.data['urls'] as List<String>;

// 4. Bulk create assets
final assets = tags.asMap().entries.map((entry) {
  int index = entry.key;
  String tag = entry.value;

  return {
    'assetTag': tag,
    'dataMatrixImageUrl': imageUrls[index],
    'serialNumber': originalAsset.serialNumber != null
        ? '${originalAsset.serialNumber}-${index + 1}'
        : null,
    // Copy all other fields from original asset
    ...originalAsset.toJson(),
  };
}).toList();

final createResponse = await http.post(
  Uri.parse('$baseUrl/assets/bulk'),
  body: json.encode({'assets': assets}),
);
```

---

## ⚠️ Error Handling

### Common Errors:
1. **Asset tag already exists** - Tag sudah terpakai di database
2. **Serial number exists** - Serial number sudah terpakai
3. **File count mismatch** - Jumlah file tidak sama dengan jumlah tag
4. **File too large** - Image lebih dari 10MB
5. **Cloudinary upload failed** - Masalah saat upload image

### Best Practice:
- Selalu validasi response setelah step 1
- Handle partial upload failures di step 2 (beberapa image mungkin gagal)
- Gunakan transaction/rollback jika bulk create gagal di tengah
