# Notification Messages

This package contains notification message templates and translations for the inventory API.

## Structure

```
internal/notification/messages/
├── common.go          # Common types, helpers, and shared functions
├── asset.go           # Asset-related notification messages
├── issue_report.go    # Issue report notification messages
└── README.md          # This file
```

## Overview

The notification messages system provides:
- **Multi-language support** (English, Japanese)
- **Type-safe message keys**
- **Helper functions** for creating notifications
- **Consistent message formatting** across all domains

## Supported Languages

- `en-US` - English (United States)
- `ja-JP` - Japanese (Japan)

## Usage

### 1. Asset Notifications

```go
import "github.com/Rizz404/inventory-api/internal/notification/messages"

// Create asset assignment notification
titleKey, messageKey, params := messages.AssetAssignmentNotification("Laptop Dell XPS", "ASS-LPT-0001", false)

// Get translations for all languages
translations := messages.GetAssetNotificationTranslations(titleKey, messageKey, params)

// Or get single language
message := messages.GetAssetNotificationMessage(messageKey, "en-US", params)
```

### 2. Issue Report Notifications

```go
import "github.com/Rizz404/inventory-api/internal/notification/messages"

// Create issue reported notification
titleKey, messageKey, params := messages.IssueReportedNotification("Laptop Dell XPS", "ASS-LPT-0001")

// Get translations
translations := messages.GetIssueReportNotificationTranslations(titleKey, messageKey, params)
```

## Adding New Domain Notifications

To add notifications for a new domain (e.g., maintenance, user):

### 1. Create a new file

```go
// internal/notification/messages/maintenance.go
package messages

// Define constants for message keys
const (
    NotifMaintenanceScheduledTitleKey   NotificationMessageKey = "notification.maintenance.scheduled.title"
    NotifMaintenanceScheduledMessageKey NotificationMessageKey = "notification.maintenance.scheduled.message"
)

// Define translations map
var maintenanceNotificationTranslations = map[NotificationMessageKey]map[string]string{
    NotifMaintenanceScheduledTitleKey: {
        "en-US": "Maintenance Scheduled",
        "ja-JP": "メンテナンスがスケジュールされました",
    },
    NotifMaintenanceScheduledMessageKey: {
        "en-US": "Maintenance for asset {assetName} is scheduled on {scheduledDate}.",
        "ja-JP": "資産 {assetName} のメンテナンスが {scheduledDate} にスケジュールされました。",
    },
}

// Create getter functions
func GetMaintenanceNotificationMessage(key NotificationMessageKey, langCode string, params map[string]string) string {
    return GetNotificationMessage(key, langCode, params, maintenanceNotificationTranslations)
}

func GetMaintenanceNotificationTranslations(titleKey, messageKey NotificationMessageKey, params map[string]string) []NotificationTranslation {
    return GetNotificationTranslations(titleKey, messageKey, params, maintenanceNotificationTranslations)
}

// Create helper functions
func MaintenanceScheduledNotification(assetName, scheduledDate string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
    params := map[string]string{
        "assetName":     assetName,
        "scheduledDate": scheduledDate,
    }
    return NotifMaintenanceScheduledTitleKey, NotifMaintenanceScheduledMessageKey, params
}
```

### 2. Use in service

```go
import "github.com/Rizz404/inventory-api/internal/notification/messages"

// In your service function
titleKey, messageKey, params := messages.MaintenanceScheduledNotification(
    asset.AssetName,
    schedule.ScheduledDate.Format("2006-01-02"),
)

translations := messages.GetMaintenanceNotificationTranslations(titleKey, messageKey, params)
```

## Message Key Naming Convention

Follow this pattern for consistency:

```
notification.{domain}.{action}.{type}

Examples:
- notification.asset.assigned.title
- notification.asset.assigned.message
- notification.issue_report.resolved.title
- notification.maintenance.scheduled.message
```

## Parameters in Messages

Use curly braces `{}` for placeholders:

```go
"Asset {assetName} has been assigned to {userName}."
```

Parameters are case-insensitive during replacement.

## Common Helper Functions

Located in `common.go`:

- `GetNotificationMessage()` - Get single localized message
- `GetNotificationTranslations()` - Get all language translations
- `normalizeLanguageCode()` - Normalize language code format

## Best Practices

1. **Always provide translations** for all supported languages
2. **Use helper functions** instead of building params manually
3. **Keep messages concise** but informative
4. **Include relevant context** (asset name, dates, etc.)
5. **Group related notifications** in the same section
6. **Use consistent terminology** across messages

## Examples

### Asset Assignment (New)

```go
titleKey, messageKey, params := messages.AssetAssignmentNotification("MacBook Pro", "ASS-LPT-0042", true)
// Result:
// en-US: "New Asset Assigned", "New asset "MacBook Pro" has been assigned to you."
// ja-JP: "新しい資産が割り当てられました", "新しい資産 "MacBook Pro" があなたに割り当てられました。"
```

### Issue Resolved

```go
titleKey, messageKey, params := messages.IssueResolvedNotification(
    "Dell Monitor",
    "ASS-MON-0015",
    "Screen flickering fixed by replacing cable",
)
// Result:
// en-US: "Issue Resolved", "Issue report for asset Dell Monitor has been resolved. Resolution: Screen flickering fixed by replacing cable."
// ja-JP: "問題が解決されました", "資産 Dell Monitor の問題レポートが解決されました。解決策: Screen flickering fixed by replacing cable。"
```

## Migration Notes

### Old Usage (deprecated)
```go
import "github.com/Rizz404/inventory-api/internal/utils"

titleKey, messageKey, params := utils.AssetAssignmentNotification(...)
translations := utils.GetNotificationTranslations(titleKey, messageKey, params)
```

### New Usage
```go
import "github.com/Rizz404/inventory-api/internal/notification/messages"

titleKey, messageKey, params := messages.AssetAssignmentNotification(...)
translations := messages.GetAssetNotificationTranslations(titleKey, messageKey, params)
```

## Related Documentation

- [Notification Guide](../../../documentation/notification_guide.md)
- [Notification Messages Guide](../../../documentation/notification_messages_guide.md)
- [FCM Integration Guide](../../../documentation/fcm_integration_guide.md)
