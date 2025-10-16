# Asset Notification Messages System

## Overview

Sistem notification messages yang terstruktur dan tersentralisasi untuk asset management. Menggunakan pendekatan i18n (internationalization) yang sama seperti error messages, namun khusus untuk notification.

## Architecture

```
┌────────────────────────────┐
│ Asset Service              │
│ (Business Logic)           │
└────────────┬───────────────┘
             │
             ▼
┌────────────────────────────┐
│ notification_messages.go   │
│ (Helper Functions)         │
└────────────┬───────────────┘
             │
             ▼
┌────────────────────────────┐
│ Notification Service       │
│ (Creates DB + FCM)         │
└────────────────────────────┘
```

## File Structure

```
internal/utils/
├── i18n.go                      # Error messages (existing)
└── notification_messages.go     # Notification messages (new)
```

## Features

### 1. **Centralized Message Management**
- Semua notification messages disimpan dalam satu file
- Mudah untuk maintain dan update
- Konsisten dengan error message system

### 2. **Multi-language Support**
- English (en)
- Indonesian (id-ID)
- Mudah menambahkan bahasa baru

### 3. **Type-safe Keys**
- Menggunakan custom type `NotificationMessageKey`
- Compile-time checking untuk key yang tidak valid
- Autocomplete support di IDE

### 4. **Dynamic Placeholders**
- Support untuk dynamic values: `{assetName}`, `{assetTag}`, dll
- Type-safe parameter passing

### 5. **Helper Functions**
- Pre-built functions untuk common scenarios
- Mengurangi code duplication
- Lebih mudah digunakan

## Notification Use Cases

### ✅ Implemented

#### 1. **Asset Assignment**
```go
// When asset is assigned to a user
NotifAssetAssignedTitleKey
NotifAssetAssignedMessageKey

// When new asset is created and assigned
NotifAssetNewAssignedTitleKey
NotifAssetNewAssignedMessageKey

// When asset is unassigned
NotifAssetUnassignedTitleKey
NotifAssetUnassignedMessageKey
```

#### 2. **Asset Status Change**
```go
// Generic status change
NotifAssetStatusChangedTitleKey
NotifAssetStatusChangedMessageKey

// Specific statuses
NotifAssetActivatedTitleKey          // Active
NotifAssetMaintenanceTitleKey        // Maintenance
NotifAssetDisposedTitleKey           // Disposed
NotifAssetLostTitleKey               // Lost
```

#### 3. **Asset Condition Change**
```go
// Generic condition change
NotifAssetConditionChangedTitleKey
NotifAssetConditionChangedMessageKey

// Specific conditions
NotifAssetConditionDamagedTitleKey   // Damaged
NotifAssetConditionPoorTitleKey      // Poor
```

### 📋 Planned (Not Yet Implemented)

#### 4. **Asset Location Change**
```go
NotifAssetLocationChangedTitleKey
NotifAssetLocationChangedMessageKey
```

#### 5. **Asset Creation/Deletion**
```go
NotifAssetCreatedTitleKey            // For admin/manager
NotifAssetDeletedTitleKey
```

#### 6. **Asset Warranty**
```go
NotifAssetWarrantyExpiringSoonTitleKey
NotifAssetWarrantyExpiredTitleKey
```

#### 7. **Asset High Value**
```go
NotifAssetHighValueTitleKey          // For high-value assets
```

## Usage Examples

### Example 1: Asset Assignment (Simple)

```go
// In asset_service.go
func (s *Service) CreateAsset(ctx context.Context, payload *domain.CreateAssetPayload, ...) {
    // ... create asset logic ...

    // Send notification if assigned
    if payload.AssignedTo != nil && *payload.AssignedTo != "" {
        go s.sendAssetAssignmentNotification(ctx, &createdAsset, *payload.AssignedTo, true)
    }
}

// Helper method
func (s *Service) sendAssetAssignmentNotification(ctx context.Context, asset *domain.Asset, userId string, isNewAsset bool) {
    // Get message keys and params
    titleKey, messageKey, params := utils.AssetAssignmentNotification(
        asset.AssetName,
        asset.AssetTag,
        isNewAsset,
    )

    // Get translations
    utilTranslations := utils.GetNotificationTranslations(titleKey, messageKey, params)

    // Convert to domain type
    translations := make([]domain.CreateNotificationTranslationPayload, len(utilTranslations))
    for i, t := range utilTranslations {
        translations[i] = domain.CreateNotificationTranslationPayload{
            LangCode: t.LangCode,
            Title:    t.Title,
            Message:  t.Message,
        }
    }

    // Create notification
    notificationPayload := &domain.CreateNotificationPayload{
        UserID:         userId,
        RelatedAssetID: &asset.ID,
        Type:           domain.NotificationTypeStatusChange,
        Translations:   translations,
    }

    s.NotificationService.CreateNotification(ctx, notificationPayload)
}
```

### Example 2: Multiple Notifications on Update

```go
// In asset_service.go
func (s *Service) UpdateAsset(ctx context.Context, assetId string, payload *domain.UpdateAssetPayload, ...) {
    // ... update asset logic ...

    // Send all relevant notifications
    go s.sendUpdateNotifications(ctx, &existingAsset, &updatedAsset, payload)
}

func (s *Service) sendUpdateNotifications(ctx context.Context, oldAsset, newAsset *domain.Asset, payload *domain.UpdateAssetPayload) {
    // 1. Check assignment changes
    if payload.AssignedTo != nil {
        if *payload.AssignedTo != "" && (oldAsset.AssignedTo == nil || *oldAsset.AssignedTo != *payload.AssignedTo) {
            s.sendAssetAssignmentNotification(ctx, newAsset, *payload.AssignedTo, false)
        } else if *payload.AssignedTo == "" && oldAsset.AssignedTo != nil {
            s.sendAssetUnassignmentNotification(ctx, newAsset, *oldAsset.AssignedTo)
        }
    }

    // 2. Check status changes
    if payload.Status != nil && *payload.Status != oldAsset.Status {
        s.sendAssetStatusChangeNotification(ctx, newAsset, oldAsset.Status, *payload.Status)
    }

    // 3. Check condition changes
    if payload.Condition != nil && *payload.Condition != oldAsset.Condition {
        s.sendAssetConditionChangeNotification(ctx, newAsset, oldAsset.Condition, *payload.Condition)
    }
}
```

### Example 3: Custom Message with Placeholders

```go
// Get message for specific language
params := map[string]string{
    "assetName": "MacBook Pro 2023",
    "assetTag":  "LAPTOP-001",
}

titleEN := utils.GetNotificationMessage(
    utils.NotifAssetAssignedTitleKey,
    "en",
    params,
)
// Result: "Asset Assigned"

messageEN := utils.GetNotificationMessage(
    utils.NotifAssetAssignedMessageKey,
    "en",
    params,
)
// Result: "Asset 'MacBook Pro 2023' (Tag: LAPTOP-001) has been assigned to you."

messageID := utils.GetNotificationMessage(
    utils.NotifAssetAssignedMessageKey,
    "id-ID",
    params,
)
// Result: "Aset 'MacBook Pro 2023' (Tag: LAPTOP-001) telah ditugaskan kepada Anda."
```

## Adding New Notification Types

### Step 1: Add Message Keys

```go
// In notification_messages.go
const (
    // ... existing keys ...

    // New notification type
    NotifAssetMaintenanceScheduledTitleKey   NotificationMessageKey = "notification.asset.maintenance_scheduled.title"
    NotifAssetMaintenanceScheduledMessageKey NotificationMessageKey = "notification.asset.maintenance_scheduled.message"
)
```

### Step 2: Add Translations

```go
// In notificationMessageTranslations map
var notificationMessageTranslations = map[NotificationMessageKey]map[string]string{
    // ... existing translations ...

    NotifAssetMaintenanceScheduledTitleKey: {
        "en":    "Maintenance Scheduled",
        "id-ID": "Pemeliharaan Dijadwalkan",
    },
    NotifAssetMaintenanceScheduledMessageKey: {
        "en":    "Maintenance for asset '{assetName}' scheduled on {date}.",
        "id-ID": "Pemeliharaan untuk aset '{assetName}' dijadwalkan pada {date}.",
    },
}
```

### Step 3: Add Helper Function

```go
// In notification_messages.go
func AssetMaintenanceScheduledNotification(assetName, date string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
    params := map[string]string{
        "assetName": assetName,
        "date":      date,
    }
    return NotifAssetMaintenanceScheduledTitleKey, NotifAssetMaintenanceScheduledMessageKey, params
}
```

### Step 4: Use in Service

```go
// In asset_service.go or maintenance_service.go
func (s *Service) ScheduleMaintenance(ctx context.Context, ...) {
    // ... schedule maintenance logic ...

    titleKey, messageKey, params := utils.AssetMaintenanceScheduledNotification(
        asset.AssetName,
        schedule.Date.Format("2006-01-02"),
    )

    utilTranslations := utils.GetNotificationTranslations(titleKey, messageKey, params)

    // ... create notification ...
}
```

## Best Practices

### 1. **Always Use Helper Functions**
```go
// ✅ Good - Using helper function
titleKey, messageKey, params := utils.AssetAssignmentNotification(assetName, assetTag, isNew)

// ❌ Bad - Hardcoding keys and params
titleKey := utils.NotifAssetAssignedTitleKey
params := map[string]string{"assetName": assetName, "assetTag": assetTag}
```

### 2. **Send Notifications Asynchronously**
```go
// ✅ Good - Non-blocking
go s.sendAssetAssignmentNotification(ctx, asset, userId, true)

// ❌ Bad - Blocking
s.sendAssetAssignmentNotification(ctx, asset, userId, true)
```

### 3. **Always Check Service Availability**
```go
// ✅ Good - Graceful degradation
func (s *Service) sendNotification(...) {
    if s.NotificationService == nil {
        return
    }
    // ... send notification ...
}

// ❌ Bad - Will panic if nil
func (s *Service) sendNotification(...) {
    s.NotificationService.CreateNotification(...)
}
```

### 4. **Use Appropriate Notification Types**
```go
// ✅ Good - Using correct type
Type: domain.NotificationTypeStatusChange  // for status/assignment changes
Type: domain.NotificationTypeMaintenance   // for maintenance alerts
Type: domain.NotificationTypeWarranty      // for warranty alerts

// ❌ Bad - Using generic type for everything
Type: domain.NotificationTypeStatusChange  // for everything
```

### 5. **Include Relevant Asset ID**
```go
// ✅ Good - Including related asset
RelatedAssetID: &asset.ID

// ❌ Bad - No context
RelatedAssetID: nil
```

## Notification Triggers Matrix

| Trigger Event               | Recipient             | Notification Type | Implemented |
| --------------------------- | --------------------- | ----------------- | ----------- |
| Asset Created               | Admin/Manager         | Asset Created     | ❌           |
| Asset Assigned (New)        | Assigned User         | Asset Assignment  | ✅           |
| Asset Assigned (Transfer)   | New Assigned User     | Asset Assignment  | ✅           |
| Asset Unassigned            | Previous User         | Asset Assignment  | ✅           |
| Status → Active             | Assigned User         | Status Change     | ✅           |
| Status → Maintenance        | Assigned User         | Status Change     | ✅           |
| Status → Disposed           | Assigned User         | Status Change     | ✅           |
| Status → Lost               | Admin + Assigned User | Status Change     | ✅           |
| Condition → Damaged         | Assigned User + Admin | Condition Change  | ✅           |
| Condition → Poor            | Assigned User         | Condition Change  | ✅           |
| Location Changed            | Assigned User         | Location Change   | ❌           |
| Warranty Expiring (30 days) | Admin                 | Warranty Alert    | ❌           |
| Warranty Expired            | Admin                 | Warranty Alert    | ❌           |
| High Value Asset Created    | Admin                 | Asset Created     | ✅           |
| Asset Deleted               | Admin                 | Asset Deleted     | ❌           |

## Testing

### Unit Test Example

```go
func TestAssetAssignmentNotification(t *testing.T) {
    titleKey, messageKey, params := utils.AssetAssignmentNotification(
        "MacBook Pro",
        "LAPTOP-001",
        true,
    )

    // Test keys
    assert.Equal(t, utils.NotifAssetNewAssignedTitleKey, titleKey)
    assert.Equal(t, utils.NotifAssetNewAssignedMessageKey, messageKey)

    // Test params
    assert.Equal(t, "MacBook Pro", params["assetName"])
    assert.Equal(t, "LAPTOP-001", params["assetTag"])

    // Test English translation
    titleEN := utils.GetNotificationMessage(titleKey, "en", params)
    assert.Equal(t, "New Asset Assigned", titleEN)

    messageEN := utils.GetNotificationMessage(messageKey, "en", params)
    assert.Contains(t, messageEN, "MacBook Pro")
    assert.Contains(t, messageEN, "LAPTOP-001")

    // Test Indonesian translation
    titleID := utils.GetNotificationMessage(titleKey, "id-ID", params)
    assert.Equal(t, "Aset Baru Ditugaskan", titleID)
}
```

### Integration Test Example

```go
func TestAssetServiceSendsNotificationOnAssignment(t *testing.T) {
    // Setup mocks
    mockRepo := &MockAssetRepository{}
    mockNotifService := &MockNotificationService{}
    service := asset.NewService(mockRepo, nil, mockNotifService)

    // Create asset with assignment
    payload := &domain.CreateAssetPayload{
        AssetTag:   "TEST-001",
        AssetName:  "Test Asset",
        AssignedTo: &userId,
    }

    _, err := service.CreateAsset(context.Background(), payload, nil, "en")
    assert.NoError(t, err)

    // Wait for async notification
    time.Sleep(100 * time.Millisecond)

    // Verify notification was created
    assert.Equal(t, 1, mockNotifService.CreateNotificationCallCount)

    // Verify notification content
    notif := mockNotifService.LastNotification
    assert.Equal(t, userId, notif.UserID)
    assert.Equal(t, 2, len(notif.Translations)) // EN + ID
    assert.Contains(t, notif.Translations[0].Title, "Asset")
}
```

## Future Enhancements

1. **Notification Templates with Rich Content**
   - HTML templates for email
   - Push notification with images
   - Action buttons

2. **Notification Batching**
   - Combine multiple notifications
   - Digest emails
   - Summary notifications

3. **Notification Preferences**
   - User can choose which notifications to receive
   - Preferred delivery methods
   - Quiet hours

4. **Notification Analytics**
   - Track open rates
   - Click-through rates
   - User engagement

5. **More Languages**
   - Spanish
   - Chinese
   - Japanese
   - etc.

6. **Context-aware Messages**
   - Different messages for different roles
   - Urgency levels
   - Priority indicators

## Troubleshooting

### Issue: Notification not received

**Check:**
1. Is notification created in database?
   ```sql
   SELECT * FROM notifications WHERE user_id = 'xxx' ORDER BY created_at DESC;
   ```
2. Is FCM token set for user?
   ```sql
   SELECT fcm_token FROM users WHERE id = 'xxx';
   ```
3. Check server logs for errors
4. Verify NotificationService is not nil

### Issue: Wrong language in notification

**Check:**
1. User's `preferred_lang` in database
2. Available translations in `notificationMessageTranslations`
3. Fallback to English if language not found

### Issue: Placeholders not replaced

**Check:**
1. Placeholder names match exactly: `{assetName}` not `{asset_name}`
2. Params map contains all required placeholders
3. Check `replaceAllCaseInsensitive` function

## License

MIT
