# Asset Notification System - Complete Implementation Guide

## Overview
Sistem notifikasi untuk asset management yang terintegrasi penuh dengan cron job untuk automated notifications.

## 📋 Notification Types Yang Diimplementasikan

### 1. **Asset Assignment Notifications**
- **New Assignment**: Notifikasi saat asset baru di-assign ke user
- **Reassignment**: Notifikasi saat asset dipindahkan ke user lain
- **Unassignment**: Notifikasi saat asset tidak lagi di-assign ke user

### 2. **Asset Status Change Notifications**
- **General Status Change**: Notifikasi perubahan status (Active, Inactive, dll)
- **Activated**: Notifikasi khusus saat asset diaktifkan
- **Maintenance**: Notifikasi saat asset masuk maintenance
- **Disposed**: Notifikasi saat asset di-dispose
- **Lost**: Notifikasi saat asset dilaporkan hilang

### 3. **Asset Condition Change Notifications**
- **General Condition Change**: Notifikasi perubahan kondisi
- **Damaged**: Notifikasi khusus untuk kondisi rusak
- **Poor**: Notifikasi untuk kondisi buruk

### 4. **Asset Location Change Notifications**
- Notifikasi saat lokasi asset berubah

### 5. **Asset Creation/Deletion Notifications**
- **Created**: Notifikasi saat asset baru dibuat
- **Deleted**: Notifikasi saat asset dihapus
- **High Value**: Notifikasi khusus untuk asset bernilai tinggi

### 6. **Warranty Notifications (Automated via Cron)**
- **Warranty Expiring Soon**: Notifikasi 30 hari sebelum warranty expire
- **Warranty Expired**: Notifikasi saat warranty telah expire

## 🤖 Cron Job Implementation

### Cron Schedule
```go
// Check warranty expiring daily at 9:00 AM
"0 0 9 * * *"

// Check expired warranties daily at 9:30 AM
"0 30 9 * * *"
```

### Cron Service Features
- **Automatic Start**: Dimulai saat aplikasi start
- **Graceful Shutdown**: Berhenti dengan proper cleanup saat aplikasi stop
- **Error Recovery**: Menggunakan cron.Recover untuk handle panic
- **Efficient Query**: Query database hanya untuk assets dengan warranty yang relevan

### How It Works
1. **Warranty Expiring Check (9:00 AM)**
   - Query assets dengan `warranty_end` antara today dan 30 hari dari sekarang
   - Hanya assets yang `assigned_to` user tertentu
   - Kirim notifikasi ke assigned user

2. **Warranty Expired Check (9:30 AM)**
   - Query assets dengan `warranty_end` yang expired dalam 24 jam terakhir
   - Hanya assets yang `assigned_to` user tertentu
   - Kirim notifikasi ke assigned user

## 📁 File Structure

```
services/asset/
├── asset_service.go         # Main service dengan notification methods
├── cron_service.go          # Cron service untuk automated notifications
└── export_service.go        # Export functionality

internal/notification/messages/
└── asset.go                 # Message templates & translations

internal/postgresql/
└── asset_repository.go      # Repository dengan warranty query methods

domain/
└── notification.go          # Notification types definition

app/
└── main.go                  # Cron service initialization
```

## 🔧 Implementation Details

### 1. Service Layer (`asset_service.go`)

#### Notification Helper Methods:
```go
// Assignment
func (s *Service) sendAssetAssignmentNotification(...)
func (s *Service) sendAssetUnassignmentNotification(...)

// Status & Condition
func (s *Service) sendAssetStatusChangeNotification(...)
func (s *Service) sendAssetConditionChangeNotification(...)

// Location & Lifecycle
func (s *Service) sendAssetLocationChangeNotification(...)
func (s *Service) sendAssetCreatedNotification(...)
func (s *Service) sendAssetDeletedNotification(...)

// High Value
func (s *Service) sendHighValueAssetNotification(...)
```

#### Update Flow:
```go
func (s *Service) UpdateAsset(...) {
    // ... update logic ...

    // Send notifications asynchronously
    go s.sendUpdateNotifications(ctx, oldAsset, newAsset, payload)
}

func (s *Service) sendUpdateNotifications(...) {
    // Check for assignment changes
    // Check for status changes
    // Check for condition changes
    // Check for location changes (if needed)
}
```

### 2. Cron Service (`cron_service.go`)

#### Initialization:
```go
cronService := asset.NewCronService(assetRepository, notificationService)
if err := cronService.Start(); err != nil {
    log.Fatal(err)
}
defer cronService.Stop()
```

#### Methods:
```go
// Core methods
func (cs *CronService) Start() error
func (cs *CronService) Stop()

// Scheduled tasks
func (cs *CronService) checkWarrantyExpiring()
func (cs *CronService) checkExpiredWarranties()

// Notification senders
func (cs *CronService) sendWarrantyExpiringNotification(...)
func (cs *CronService) sendWarrantyExpiredNotification(...)
```

### 3. Repository Layer (`asset_repository.go`)

#### New Methods:
```go
func (r *AssetRepository) GetAssetsWithWarrantyExpiring(ctx context.Context, daysFromNow int) ([]domain.Asset, error)
func (r *AssetRepository) GetAssetsWithExpiredWarranty(ctx context.Context) ([]domain.Asset, error)
```

#### Query Optimization:
- Filter `warranty_end IS NOT NULL`
- Filter `assigned_to IS NOT NULL`
- Only query specific date ranges
- Include preloads untuk related data

### 4. Message Templates (`internal/notification/messages/asset.go`)

#### Structure:
```go
const (
    NotifAssetAssignedTitleKey   NotificationMessageKey = "..."
    NotifAssetAssignedMessageKey NotificationMessageKey = "..."
)

var assetNotificationTranslations = map[NotificationMessageKey]map[string]string{
    NotifAssetAssignedTitleKey: {
        "en-US": "Asset Assigned",
        "ja-JP": "資産が割り当てられました",
    },
}
```

#### Helper Functions:
```go
func AssetAssignmentNotification(...) (titleKey, messageKey, params)
func AssetUnassignmentNotification(...) (titleKey, messageKey, params)
func AssetStatusChangeNotification(...) (titleKey, messageKey, params)
// ... etc
```

## 🌍 Multi-Language Support

Semua notifications mendukung multiple languages:
- English (en-US)
- Japanese (ja-JP)

Translations stored di `assetNotificationTranslations` map dan diambil menggunakan helper functions.

## 📊 Notification Types (Domain)

```go
const (
    NotificationTypeMaintenance    NotificationType = "MAINTENANCE"
    NotificationTypeWarranty       NotificationType = "WARRANTY"
    NotificationTypeStatusChange   NotificationType = "STATUS_CHANGE"
    NotificationTypeMovement       NotificationType = "MOVEMENT"
    NotificationTypeIssueReport    NotificationType = "ISSUE_REPORT"
    NotificationTypeLocationChange NotificationType = "LOCATION_CHANGE"
    NotificationTypeAssetCreated   NotificationType = "ASSET_CREATED"
    NotificationTypeAssetDeleted   NotificationType = "ASSET_DELETED"
    NotificationTypeHighValue      NotificationType = "HIGH_VALUE"
)
```

## 🚀 Usage Examples

### Manual Notification (dalam service methods):
```go
// When creating asset
if payload.AssignedTo != nil && *payload.AssignedTo != "" {
    go s.sendAssetAssignmentNotification(ctx, &createdAsset, *payload.AssignedTo, true)
}

// When updating asset
go s.sendUpdateNotifications(ctx, &existingAsset, &updatedAsset, payload)
```

### Automated Notification (via cron):
```go
// Runs automatically at scheduled times
// No manual intervention needed
// Logs all activities
```

## 📝 Logging

Semua notification activities di-log:
```
[INFO] Sending asset assignment notification for asset ID: xxx, asset tag: xxx, user ID: xxx
[INFO] Successfully created asset assignment notification for asset ID: xxx, user ID: xxx
[ERROR] Failed to create asset assignment notification for asset ID: xxx, user ID: xxx: error details
```

Cron job activities juga di-log:
```
[INFO] Asset cron service started successfully
[INFO] Running warranty expiring check...
[INFO] Warranty expiring check completed. Found 5 assets with warranties expiring within 30 days
[INFO] Running expired warranty check...
[INFO] Expired warranty check completed. Found 2 assets with expired warranties
```

## ⚙️ Configuration

### Cron Schedule Customization:
Edit di `cron_service.go`:
```go
// Format: second minute hour day month weekday
"0 0 9 * * *"  // Daily at 9:00 AM
"0 30 9 * * *" // Daily at 9:30 AM
```

### Warranty Check Window:
```go
// Default: 30 days
assets, err := cs.assetRepo.GetAssetsWithWarrantyExpiring(ctx, 30)

// Ubah angka untuk window yang berbeda (e.g., 60 untuk 60 hari)
```

## 🔍 Testing Recommendations

### 1. Unit Tests:
- Test individual notification methods
- Mock notification service
- Verify correct message keys dan params

### 2. Integration Tests:
- Test notification creation flow
- Verify database queries
- Check multi-language support

### 3. Cron Job Tests:
- Test manual cron trigger
- Verify date calculations
- Check notification sending

### 4. Manual Testing:
```go
// Trigger cron manually for testing
cronService.checkWarrantyExpiring()
cronService.checkExpiredWarranties()
```

## 🎯 Best Practices

1. **Async Notifications**: Selalu gunakan `go` untuk send notifications agar tidak blocking
2. **Error Handling**: Log errors tapi jangan fail main operation
3. **Null Checks**: Selalu check `AssignedTo`, `WarrantyEnd`, dll sebelum digunakan
4. **Performance**: Query hanya data yang diperlukan dengan filter yang tepat
5. **Logging**: Log semua notification activities untuk debugging
6. **Graceful Shutdown**: Ensure cron service stops properly saat app shutdown

## 🐛 Troubleshooting

### Notifications tidak terkirim:
1. Check NotificationService != nil
2. Check AssignedTo tidak nil/empty
3. Check logs untuk error details
4. Verify FCM configuration

### Cron job tidak running:
1. Check cron service started di main.go
2. Verify cron schedule format
3. Check logs untuk startup errors
4. Ensure proper timezone configuration

### Duplicate notifications:
1. Verify cron schedule tidak overlap
2. Check warranty date filters
3. Ensure single cron instance running

## 📚 Dependencies

```go
// Cron library
"github.com/robfig/cron/v3"

// Internal packages
"github.com/Rizz404/inventory-api/domain"
"github.com/Rizz404/inventory-api/internal/notification/messages"
"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
```

## 🔄 Future Enhancements

1. **Notification Preferences**: Allow users to customize notification types
2. **Batch Notifications**: Group multiple notifications into digest
3. **Custom Schedules**: Allow dynamic cron schedule configuration
4. **Notification History**: Track sent notifications
5. **Retry Logic**: Implement retry for failed notifications
6. **Email Integration**: Send critical notifications via email
7. **SMS Integration**: SMS for urgent notifications
8. **Webhook Support**: Call external webhooks on events

## ✅ Checklist Implementation

- [x] Asset assignment notifications
- [x] Asset status change notifications
- [x] Asset condition change notifications
- [x] Asset location change notifications
- [x] Asset creation/deletion notifications
- [x] High value asset notifications
- [x] Warranty expiring notifications (automated)
- [x] Warranty expired notifications (automated)
- [x] Multi-language support
- [x] Cron job implementation
- [x] Repository methods for warranty queries
- [x] Proper error handling dan logging
- [x] Graceful shutdown
- [x] Documentation

## 📞 Support

Untuk pertanyaan atau issues terkait notification system, silakan check:
1. Log files untuk error details
2. Documentation di `documentation/notification_guide.md`
3. Message templates di `internal/notification/messages/asset.go`
4. Cron service implementation di `services/asset/cron_service.go`
