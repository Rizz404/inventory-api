# Notification Types Documentation

Dokumentasi lengkap tipe-tipe notifikasi yang tersedia di Inventory API untuk keperluan frontend development.

---

## 📋 Overview

### Supported Languages

| Code    | Language   |
| ------- | ---------- |
| `en-US` | English    |
| `id-ID` | Indonesian |
| `ja-JP` | Japanese   |

### Notification Priority Levels

| Priority | FCM Priority | Description                   |
| -------- | ------------ | ----------------------------- |
| `LOW`    | `normal`     | Informational notifications   |
| `NORMAL` | `normal`     | Standard notifications        |
| `HIGH`   | `high`       | Important notifications       |
| `URGENT` | `high`       | Critical, immediate attention |

---

## 📦 Notification Types

### 1. `MAINTENANCE` - Maintenance Notifications

Notifikasi terkait jadwal dan pelaksanaan pemeliharaan aset.

#### 1.1 Maintenance Scheduled

Notifikasi saat jadwal pemeliharaan baru dibuat.

| Language | Title                                | Message                                                                        |
| -------- | ------------------------------------ | ------------------------------------------------------------------------------ |
| `en-US`  | Maintenance Scheduled                | Maintenance for asset "{assetName}" is scheduled on {scheduledDate}.           |
| `id-ID`  | Pemeliharaan Dijadwalkan             | Pemeliharaan untuk aset "{assetName}" dijadwalkan pada {scheduledDate}.        |
| `ja-JP`  | メンテナンスがスケジュールされました | 資産 "{assetName}" のメンテナンスが {scheduledDate} にスケジュールされました。 |

**Parameters:**
- `{assetName}` - Nama aset
- `{assetTag}` - Tag/kode aset
- `{scheduledDate}` - Tanggal jadwal pemeliharaan

#### 1.2 Maintenance Due Soon

Notifikasi pengingat jadwal pemeliharaan yang akan segera jatuh tempo.

| Language | Title                            | Message                                                                                     |
| -------- | -------------------------------- | ------------------------------------------------------------------------------------------- |
| `en-US`  | Maintenance Due Soon             | Maintenance for asset "{assetName}" is due on {scheduledDate}. Please prepare.              |
| `id-ID`  | Pemeliharaan Segera Jatuh Tempo  | Pemeliharaan untuk aset "{assetName}" jatuh tempo pada {scheduledDate}. Silakan persiapkan. |
| `ja-JP`  | メンテナンス期限が近づいています | 資産 "{assetName}" のメンテナンス期限が {scheduledDate} です。準備してください。            |

**Parameters:**
- `{assetName}` - Nama aset
- `{assetTag}` - Tag/kode aset
- `{scheduledDate}` - Tanggal jadwal pemeliharaan

#### 1.3 Maintenance Overdue

Notifikasi saat pemeliharaan melewati jadwal yang ditentukan.

| Language | Title                      | Message                                                                                                 |
| -------- | -------------------------- | ------------------------------------------------------------------------------------------------------- |
| `en-US`  | Maintenance Overdue        | Maintenance for asset "{assetName}" is overdue. Scheduled date was {scheduledDate}.                     |
| `id-ID`  | Pemeliharaan Terlambat     | Pemeliharaan untuk aset "{assetName}" sudah terlambat. Tanggal yang dijadwalkan adalah {scheduledDate}. |
| `ja-JP`  | メンテナンスが期限切れです | 資産 "{assetName}" のメンテナンスが期限切れです。スケジュールされた日付は {scheduledDate} でした。      |

**Parameters:**
- `{assetName}` - Nama aset
- `{assetTag}` - Tag/kode aset
- `{scheduledDate}` - Tanggal jadwal pemeliharaan

#### 1.4 Maintenance Completed

Notifikasi saat pemeliharaan berhasil diselesaikan.

| Language | Title                      | Message                                                                   |
| -------- | -------------------------- | ------------------------------------------------------------------------- |
| `en-US`  | Maintenance Completed      | Maintenance for asset "{assetName}" has been completed. Notes: "{notes}". |
| `id-ID`  | Pemeliharaan Selesai       | Pemeliharaan untuk aset "{assetName}" telah selesai. Catatan: "{notes}".  |
| `ja-JP`  | メンテナンスが完了しました | 資産 "{assetName}" のメンテナンスが完了しました。メモ: "{notes}"。        |

**Parameters:**
- `{assetName}` - Nama aset
- `{assetTag}` - Tag/kode aset
- `{notes}` - Catatan pemeliharaan

#### 1.5 Maintenance Failed

Notifikasi saat pemeliharaan gagal dilakukan.

| Language | Title                      | Message                                                                                    |
| -------- | -------------------------- | ------------------------------------------------------------------------------------------ |
| `en-US`  | Maintenance Failed         | Maintenance for asset "{assetName}" could not be completed. Reason: "{failureReason}".     |
| `id-ID`  | Pemeliharaan Gagal         | Pemeliharaan untuk aset "{assetName}" tidak dapat diselesaikan. Alasan: "{failureReason}". |
| `ja-JP`  | メンテナンスが失敗しました | 資産 "{assetName}" のメンテナンスが完了できませんでした。理由: "{failureReason}"。         |

**Parameters:**
- `{assetName}` - Nama aset
- `{assetTag}` - Tag/kode aset
- `{failureReason}` - Alasan kegagalan

---

### 2. `WARRANTY` - Warranty Notifications

Notifikasi terkait masa garansi aset.

#### 2.1 Warranty Expiring Soon

Notifikasi pengingat garansi yang akan segera berakhir.

| Language | Title                        | Message                                                           |
| -------- | ---------------------------- | ----------------------------------------------------------------- |
| `en-US`  | Warranty Expiring Soon       | Warranty for asset "{assetName}" will expire on {expiryDate}.     |
| `id-ID`  | Garansi Akan Berakhir        | Garansi untuk aset "{assetName}" akan berakhir pada {expiryDate}. |
| `ja-JP`  | 保証期間がまもなく終了します | 資産 "{assetName}" の保証期間が {expiryDate} に終了します。       |

**Parameters:**
- `{assetName}` - Nama aset
- `{assetTag}` - Tag/kode aset
- `{expiryDate}` - Tanggal berakhir garansi

#### 2.2 Warranty Expired

Notifikasi saat garansi telah berakhir.

| Language | Title                  | Message                                          |
| -------- | ---------------------- | ------------------------------------------------ |
| `en-US`  | Warranty Expired       | Warranty for asset "{assetName}" has expired.    |
| `id-ID`  | Garansi Berakhir       | Garansi untuk aset "{assetName}" telah berakhir. |
| `ja-JP`  | 保証期間が終了しました | 資産 "{assetName}" の保証期間が終了しました。    |

**Parameters:**
- `{assetName}` - Nama aset
- `{assetTag}` - Tag/kode aset

---

### 3. `ISSUE` - Issue Report Notifications

Notifikasi terkait laporan masalah/kerusakan aset.

#### 3.1 New Issue Reported

Notifikasi saat masalah baru dilaporkan.

| Language | Title                      | Message                                                 |
| -------- | -------------------------- | ------------------------------------------------------- |
| `en-US`  | New Issue Reported         | A new issue has been reported for asset "{assetName}".  |
| `id-ID`  | Masalah Baru Dilaporkan    | Masalah baru telah dilaporkan untuk aset "{assetName}". |
| `ja-JP`  | 新しい問題が報告されました | 資産 "{assetName}" に対して新しい問題が報告されました。 |

**Parameters:**
- `{assetName}` - Nama aset
- `{assetTag}` - Tag/kode aset

#### 3.2 Issue Updated

Notifikasi saat laporan masalah diperbarui.

| Language | Title                | Message                                                    |
| -------- | -------------------- | ---------------------------------------------------------- |
| `en-US`  | Issue Updated        | Issue report for asset "{assetName}" has been updated.     |
| `id-ID`  | Masalah Diperbarui   | Laporan masalah untuk aset "{assetName}" telah diperbarui. |
| `ja-JP`  | 問題が更新されました | 資産 "{assetName}" の問題レポートが更新されました。        |

**Parameters:**
- `{assetName}` - Nama aset
- `{assetTag}` - Tag/kode aset

#### 3.3 Issue Resolved

Notifikasi saat masalah berhasil diselesaikan.

| Language | Title                | Message                                                                                     |
| -------- | -------------------- | ------------------------------------------------------------------------------------------- |
| `en-US`  | Issue Resolved       | Issue report for asset "{assetName}" has been resolved. Resolution: "{resolutionNotes}".    |
| `id-ID`  | Masalah Diselesaikan | Laporan masalah untuk aset "{assetName}" telah diselesaikan. Resolusi: "{resolutionNotes}". |
| `ja-JP`  | 問題が解決されました | 資産 "{assetName}" の問題レポートが解決されました。解決策: "{resolutionNotes}"。            |

**Parameters:**
- `{assetName}` - Nama aset
- `{assetTag}` - Tag/kode aset
- `{resolutionNotes}` - Catatan penyelesaian

#### 3.4 Issue Reopened

Notifikasi saat masalah dibuka kembali.

| Language | Title                  | Message                                                        |
| -------- | ---------------------- | -------------------------------------------------------------- |
| `en-US`  | Issue Reopened         | Issue report for asset "{assetName}" has been reopened.        |
| `id-ID`  | Masalah Dibuka Kembali | Laporan masalah untuk aset "{assetName}" telah dibuka kembali. |
| `ja-JP`  | 問題が再開されました   | 資産 "{assetName}" の問題レポートが再開されました。            |

**Parameters:**
- `{assetName}` - Nama aset
- `{assetTag}` - Tag/kode aset

---

### 4. `MOVEMENT` - Asset Movement Notifications

Notifikasi terkait perpindahan lokasi atau penugasan aset.

#### 4.1 Asset Location Changed

Notifikasi saat lokasi aset berubah.

| Language | Title                      | Message                                                                                    |
| -------- | -------------------------- | ------------------------------------------------------------------------------------------ |
| `en-US`  | Asset Location Changed     | Asset "{assetName}" ({assetTag}) has been moved from "{oldLocation}" to "{newLocation}".   |
| `id-ID`  | Lokasi Aset Berubah        | Aset "{assetName}" ({assetTag}) telah dipindahkan dari "{oldLocation}" ke "{newLocation}". |
| `ja-JP`  | 資産の場所が変更されました | 資産 "{assetName}" ({assetTag}) が "{oldLocation}" から "{newLocation}" に移動されました。 |

**Parameters:**
- `{assetName}` - Nama aset
- `{assetTag}` - Tag/kode aset
- `{oldLocation}` - Lokasi lama
- `{newLocation}` - Lokasi baru

#### 4.2 Asset Assigned to User

Notifikasi saat aset ditugaskan ke user baru.

| Language | Title                            | Message                                                                         |
| -------- | -------------------------------- | ------------------------------------------------------------------------------- |
| `en-US`  | Asset Assigned to You            | Asset "{assetName}" ({assetTag}) has been assigned from "{oldUser}" to you.     |
| `id-ID`  | Aset Ditugaskan kepada Anda      | Aset "{assetName}" ({assetTag}) telah ditugaskan dari "{oldUser}" kepada Anda.  |
| `ja-JP`  | 資産があなたに割り当てられました | 資産 "{assetName}" ({assetTag}) が "{oldUser}" からあなたに割り当てられました。 |

**Parameters:**
- `{assetName}` - Nama aset
- `{assetTag}` - Tag/kode aset
- `{oldUser}` - User sebelumnya
- `{newUser}` - User baru

---

### 5. `STATUS_CHANGE` - Asset Status Change Notifications

Notifikasi terkait perubahan status aset.

#### 5.1 Asset Status Changed (Generic)

| Language | Title                          | Message                                                                           |
| -------- | ------------------------------ | --------------------------------------------------------------------------------- |
| `en-US`  | Asset Status Changed           | Asset "{assetName}" status changed from {oldStatus} to {newStatus}.               |
| `id-ID`  | Status Aset Berubah            | Status aset "{assetName}" berubah dari {oldStatus} menjadi {newStatus}.           |
| `ja-JP`  | 資産ステータスが変更されました | 資産 "{assetName}" のステータスが {oldStatus} から {newStatus} に変更されました。 |

**Parameters:**
- `{assetName}` - Nama aset
- `{oldStatus}` - Status lama
- `{newStatus}` - Status baru

#### 5.2 Asset Activated

| Language | Title                  | Message                                                 |
| -------- | ---------------------- | ------------------------------------------------------- |
| `en-US`  | Asset Activated        | Asset "{assetName}" is now active and ready to use.     |
| `id-ID`  | Aset Diaktifkan        | Aset "{assetName}" sekarang aktif dan siap digunakan.   |
| `ja-JP`  | 資産が有効化されました | 資産 "{assetName}" が有効化され、使用準備が整いました。 |

#### 5.3 Asset Under Maintenance

| Language | Title                    | Message                                                       |
| -------- | ------------------------ | ------------------------------------------------------------- |
| `en-US`  | Asset Under Maintenance  | Asset "{assetName}" has been moved to maintenance status.     |
| `id-ID`  | Aset Dalam Pemeliharaan  | Aset "{assetName}" telah dipindahkan ke status pemeliharaan.  |
| `ja-JP`  | 資産がメンテナンス中です | 資産 "{assetName}" がメンテナンスステータスに移動されました。 |

#### 5.4 Asset Disposed

| Language | Title                | Message                                |
| -------- | -------------------- | -------------------------------------- |
| `en-US`  | Asset Disposed       | Asset "{assetName}" has been disposed. |
| `id-ID`  | Aset Dibuang         | Aset "{assetName}" telah dibuang.      |
| `ja-JP`  | 資産が廃棄されました | 資産 "{assetName}" が廃棄されました。  |

#### 5.5 Asset Reported Lost

| Language | Title                              | Message                                             |
| -------- | ---------------------------------- | --------------------------------------------------- |
| `en-US`  | Asset Reported Lost                | Asset "{assetName}" has been reported as lost.      |
| `id-ID`  | Aset Dilaporkan Hilang             | Aset "{assetName}" telah dilaporkan hilang.         |
| `ja-JP`  | 資産が行方不明として報告されました | 資産 "{assetName}" が行方不明として報告されました。 |

---

### 6. `LOCATION_CHANGE` - Location Update Notifications

Notifikasi terkait perubahan data lokasi.

#### 6.1 Location Updated

| Language | Title                | Message                                                   |
| -------- | -------------------- | --------------------------------------------------------- |
| `en-US`  | Location Updated     | Location "{locationName}" has been updated in the system. |
| `id-ID`  | Lokasi Diperbarui    | Lokasi "{locationName}" telah diperbarui di sistem.       |
| `ja-JP`  | 場所が更新されました | 場所 "{locationName}" がシステムで更新されました。        |

**Parameters:**
- `{locationName}` - Nama lokasi

---

### 7. `CATEGORY_CHANGE` - Category Update Notifications

Notifikasi terkait perubahan data kategori.

#### 7.1 Category Updated

| Language | Title                    | Message                                      |
| -------- | ------------------------ | -------------------------------------------- |
| `en-US`  | Category Updated         | Category "{categoryName}" has been updated.  |
| `id-ID`  | Kategori Diperbarui      | Kategori "{categoryName}" telah diperbarui.  |
| `ja-JP`  | カテゴリが更新されました | カテゴリ "{categoryName}" が更新されました。 |

**Parameters:**
- `{categoryName}` - Nama kategori

---

## 🎨 Asset-Specific Notifications

### Asset Assignment

#### Asset Assigned

| Language | Title                    | Message                                           |
| -------- | ------------------------ | ------------------------------------------------- |
| `en-US`  | Asset Assigned           | Asset "{assetName}" has been assigned to you.     |
| `id-ID`  | Aset Ditugaskan          | Aset "{assetName}" telah ditugaskan kepada Anda.  |
| `ja-JP`  | 資産が割り当てられました | 資産 "{assetName}" があなたに割り当てられました。 |

#### New Asset Assigned

| Language | Title                          | Message                                                 |
| -------- | ------------------------------ | ------------------------------------------------------- |
| `en-US`  | New Asset Assigned             | New asset "{assetName}" has been assigned to you.       |
| `id-ID`  | Aset Baru Ditugaskan           | Aset baru "{assetName}" telah ditugaskan kepada Anda.   |
| `ja-JP`  | 新しい資産が割り当てられました | 新しい資産 "{assetName}" があなたに割り当てられました。 |

#### Asset Unassigned

| Language | Title                          | Message                                                   |
| -------- | ------------------------------ | --------------------------------------------------------- |
| `en-US`  | Asset Unassigned               | Asset "{assetName}" has been unassigned from you.         |
| `id-ID`  | Aset Dibatalkan                | Aset "{assetName}" telah dibatalkan dari Anda.            |
| `ja-JP`  | 資産の割り当てが解除されました | 資産 "{assetName}" の割り当てがあなたから解除されました。 |

### Asset Condition

#### Condition Changed (Generic)

| Language | Title                      | Message                                                                           |
| -------- | -------------------------- | --------------------------------------------------------------------------------- |
| `en-US`  | Asset Condition Changed    | Asset "{assetName}" condition changed from {oldCondition} to {newCondition}.      |
| `id-ID`  | Kondisi Aset Berubah       | Kondisi aset "{assetName}" berubah dari {oldCondition} menjadi {newCondition}.    |
| `ja-JP`  | 資産の状態が変更されました | 資産 "{assetName}" の状態が {oldCondition} から {newCondition} に変更されました。 |

#### Asset Damaged

| Language | Title              | Message                                                                   |
| -------- | ------------------ | ------------------------------------------------------------------------- |
| `en-US`  | Asset Damaged      | Asset "{assetName}" has been marked as damaged. Please check immediately. |
| `id-ID`  | Aset Rusak         | Aset "{assetName}" telah ditandai sebagai rusak. Silakan periksa segera.  |
| `ja-JP`  | 資産が損傷しました | 資産 "{assetName}" が損傷としてマークされました。すぐに確認してください。 |

#### Asset in Poor Condition

| Language | Title                    | Message                                                                            |
| -------- | ------------------------ | ---------------------------------------------------------------------------------- |
| `en-US`  | Asset in Poor Condition  | Asset "{assetName}" condition has deteriorated to poor. Maintenance may be needed. |
| `id-ID`  | Aset Dalam Kondisi Buruk | Kondisi aset "{assetName}" telah memburuk. Pemeliharaan mungkin diperlukan.        |
| `ja-JP`  | 資産の状態が不良です     | 資産 "{assetName}" の状態が不良に悪化しました。メンテナンスが必要かもしれません。  |

### High Value Asset

| Language | Title                            | Message                                                                                  |
| -------- | -------------------------------- | ---------------------------------------------------------------------------------------- |
| `en-US`  | High Value Asset Added           | High value asset "{assetName}" worth {value} has been added to your inventory.           |
| `id-ID`  | Aset Bernilai Tinggi Ditambahkan | Aset bernilai tinggi "{assetName}" senilai {value} telah ditambahkan ke inventaris Anda. |
| `ja-JP`  | 高額資産が追加されました         | 高額資産 "{assetName}" 価値 {value} が在庫に追加されました。                             |

**Parameters:**
- `{assetName}` - Nama aset
- `{assetTag}` - Tag/kode aset
- `{value}` - Nilai aset

---

## 📡 API Response Structure

### Single Notification Response

```json
{
  "id": "uuid-string",
  "userId": "uuid-string",
  "relatedEntityType": "asset|maintenance_schedule|issue_report|...",
  "relatedEntityId": "uuid-string",
  "relatedAssetId": "uuid-string",
  "type": "MAINTENANCE|WARRANTY|ISSUE|MOVEMENT|STATUS_CHANGE|LOCATION_CHANGE|CATEGORY_CHANGE",
  "priority": "LOW|NORMAL|HIGH|URGENT",
  "isRead": false,
  "readAt": "2025-01-01T00:00:00Z",
  "expiresAt": "2025-01-01T00:00:00Z",
  "createdAt": "2025-01-01T00:00:00Z",
  "title": "Localized title based on Accept-Language header",
  "message": "Localized message based on Accept-Language header",
  "translations": [
    {
      "langCode": "en-US",
      "title": "English Title",
      "message": "English Message"
    },
    {
      "langCode": "id-ID",
      "title": "Judul Indonesia",
      "message": "Pesan Indonesia"
    },
    {
      "langCode": "ja-JP",
      "title": "日本語タイトル",
      "message": "日本語メッセージ"
    }
  ]
}
```

### Notification List Response

```json
{
  "id": "uuid-string",
  "userId": "uuid-string",
  "relatedEntityType": "asset",
  "relatedEntityId": "uuid-string",
  "relatedAssetId": "uuid-string",
  "type": "MAINTENANCE",
  "priority": "HIGH",
  "isRead": false,
  "readAt": null,
  "expiresAt": null,
  "createdAt": "2025-01-01T00:00:00Z",
  "title": "Maintenance Due Soon",
  "message": "Maintenance for asset \"Laptop Dell XPS\" is due on 2025-01-05. Please prepare."
}
```

---

## 🔧 Frontend Implementation Notes

### 1. Language Detection

Kirim header `Accept-Language` untuk mendapatkan notifikasi dalam bahasa yang diinginkan:

```http
Accept-Language: id-ID
```

### 2. Priority-based Styling

| Priority | Suggested Color | Icon            |
| -------- | --------------- | --------------- |
| `LOW`    | Gray/Blue       | Info icon       |
| `NORMAL` | Blue            | Bell icon       |
| `HIGH`   | Orange/Yellow   | Warning icon    |
| `URGENT` | Red             | Alert/Fire icon |

### 3. Type-based Categorization

Gunakan field `type` untuk mengelompokkan atau memfilter notifikasi:

- **MAINTENANCE** - 🔧 Maintenance icon
- **WARRANTY** - 📋 Document/shield icon
- **ISSUE** - ⚠️ Warning icon
- **MOVEMENT** - 📍 Location/arrow icon
- **STATUS_CHANGE** - 🔄 Refresh/sync icon
- **LOCATION_CHANGE** - 🏢 Building/map icon
- **CATEGORY_CHANGE** - 🏷️ Tag/label icon

### 4. Related Entity Navigation

Gunakan `relatedEntityType` dan `relatedEntityId` untuk navigasi ke detail entity terkait:

```javascript
// Example
if (notification.relatedEntityType === 'asset' && notification.relatedEntityId) {
  navigateTo(`/assets/${notification.relatedEntityId}`);
}
```

---

## 📊 Statistics Response

```json
{
  "total": {
    "count": 150
  },
  "byType": {
    "maintenance": 45,
    "warranty": 12,
    "issue": 23,
    "movement": 30,
    "statusChange": 20,
    "locationChange": 10,
    "categoryChange": 10
  },
  "byStatus": {
    "read": 100,
    "unread": 50
  },
  "creationTrends": [
    {
      "date": "2025-01-01T00:00:00Z",
      "count": 10
    }
  ],
  "summary": {
    "totalNotifications": 150,
    "readPercentage": 66.67,
    "unreadPercentage": 33.33,
    "mostCommonType": "maintenance",
    "averageNotificationsPerDay": 5.0,
    "latestCreationDate": "2025-01-01T00:00:00Z",
    "earliestCreationDate": "2024-12-01T00:00:00Z"
  }
}
```

---

## 🔔 FCM Push Notification Format

Notifikasi juga dikirim via Firebase Cloud Messaging (FCM) dengan format:

```json
{
  "to": "user-fcm-token",
  "notification": {
    "title": "Localized Title",
    "body": "Localized Message"
  },
  "data": {
    "notificationId": "uuid-string",
    "type": "MAINTENANCE",
    "priority": "HIGH",
    "relatedEntityType": "asset",
    "relatedEntityId": "uuid-string"
  },
  "android": {
    "priority": "high"
  },
  "apns": {
    "headers": {
      "apns-priority": "10"
    }
  }
}
```

---

*Last updated: December 30, 2025*
