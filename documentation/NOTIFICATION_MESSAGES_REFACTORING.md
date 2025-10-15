# Notification Messages Refactoring Summary

## Overview

Notification messages have been refactored from a single file in `internal/utils/` to a more organized structure under `internal/notification/messages/`.

## Changes

### Before (Old Structure)
```
internal/utils/
└── notification_messages.go  (All notification messages in one file)
```

### After (New Structure)
```
internal/notification/messages/
├── common.go          # Common types and helper functions
├── asset.go           # Asset notification messages
├── issue_report.go    # Issue report notification messages
└── README.md          # Documentation
```

## Benefits

1. ✅ **Better Organization** - Each domain has its own file
2. ✅ **Scalability** - Easy to add new notification types
3. ✅ **Separation of Concerns** - Notifications are not mixed with generic utils
4. ✅ **Maintainability** - Easier to find and update messages
5. ✅ **Clear Ownership** - Each team can maintain their domain's messages

## Migration Guide

### Import Changes

**Old:**
```go
import "github.com/Rizz404/inventory-api/internal/utils"
```

**New:**
```go
import "github.com/Rizz404/inventory-api/internal/notification/messages"
```

### Function Call Changes

#### Asset Notifications

**Old:**
```go
titleKey, messageKey, params := utils.AssetAssignmentNotification(assetName, assetTag, isNew)
translations := utils.GetNotificationTranslations(titleKey, messageKey, params)
message := utils.GetNotificationMessage(messageKey, "en-US", params)
```

**New:**
```go
titleKey, messageKey, params := messages.AssetAssignmentNotification(assetName, assetTag, isNew)
translations := messages.GetAssetNotificationTranslations(titleKey, messageKey, params)
message := messages.GetAssetNotificationMessage(messageKey, "en-US", params)
```

#### Constants

**Old:**
```go
utils.NotifAssetAssignedTitleKey
utils.NotifAssetAssignedMessageKey
```

**New:**
```go
messages.NotifAssetAssignedTitleKey
messages.NotifAssetAssignedMessageKey
```

## Files Updated

### Services
- ✅ `services/asset/asset_service.go` - Updated to use new import path

### Documentation
- 📝 `documentation/notification_messages_guide.md` - Update examples (if needed)

## New Features

### Issue Report Notifications

Complete notification system for issue reports:

```go
import "github.com/Rizz404/inventory-api/internal/notification/messages"

// Issue Reported
titleKey, messageKey, params := messages.IssueReportedNotification("Laptop", "ASS-001")

// Issue Updated
titleKey, messageKey, params := messages.IssueUpdatedNotification("Laptop", "ASS-001")

// Issue Resolved
titleKey, messageKey, params := messages.IssueResolvedNotification("Laptop", "ASS-001", "Fixed")

// Issue Reopened
titleKey, messageKey, params := messages.IssueReopenedNotification("Laptop", "ASS-001")

// Get all translations
translations := messages.GetIssueReportNotificationTranslations(titleKey, messageKey, params)
```

### Available Issue Report Messages

1. **Issue Reported**
   - EN: "New Issue Reported" - "A new issue has been reported for asset {assetName} ({assetTag})."
   - JP: "新しい問題が報告されました" - "資産 {assetName} ({assetTag}) に対して新しい問題が報告されました。"

2. **Issue Updated**
   - EN: "Issue Updated" - "Issue report for asset {assetName} has been updated."
   - JP: "問題が更新されました" - "資産 {assetName} の問題レポートが更新されました。"

3. **Issue Resolved**
   - EN: "Issue Resolved" - "Issue report for asset {assetName} has been resolved. Resolution: {resolutionNotes}."
   - JP: "問題が解決されました" - "資産 {assetName} の問題レポートが解決されました。解決策: {resolutionNotes}。"

4. **Issue Reopened**
   - EN: "Issue Reopened" - "Issue report for asset {assetName} has been reopened."
   - JP: "問題が再開されました" - "資産 {assetName} の問題レポートが再開されました。"

## Usage Example in Service

### Issue Report Service

```go
package issue_report

import (
    "context"
    "github.com/Rizz404/inventory-api/domain"
    "github.com/Rizz404/inventory-api/internal/notification/messages"
)

func (s *issueReportServiceImpl) CreateIssueReport(ctx context.Context, payload *domain.CreateIssueReportPayload) (domain.IssueReport, error) {
    // Create issue report
    issueReport, err := s.repository.CreateIssueReport(ctx, payload)
    if err != nil {
        return domain.IssueReport{}, err
    }

    // Get asset info
    asset, _ := s.assetRepo.GetAssetById(ctx, issueReport.AssetID)

    // Send notification
    titleKey, messageKey, params := messages.IssueReportedNotification(
        asset.AssetName,
        asset.AssetTag,
    )

    translations := messages.GetIssueReportNotificationTranslations(titleKey, messageKey, params)

    // Convert to domain and send
    notifTranslations := make([]domain.CreateNotificationTranslationPayload, len(translations))
    for i, t := range translations {
        notifTranslations[i] = domain.CreateNotificationTranslationPayload{
            LangCode: t.LangCode,
            Title:    t.Title,
            Message:  t.Message,
        }
    }

    notifPayload := &domain.CreateNotificationPayload{
        UserID:       issueReport.ReportedByID,
        Type:         "ISSUE_REPORT",
        RelatedID:    &issueReport.ID,
        Translations: notifTranslations,
    }

    s.notificationService.CreateNotification(ctx, notifPayload)

    return issueReport, nil
}

func (s *issueReportServiceImpl) ResolveIssueReport(ctx context.Context, issueId string, resolutionNotes string) (domain.IssueReport, error) {
    // Resolve issue
    issueReport, err := s.repository.ResolveIssueReport(ctx, issueId, resolutionNotes)
    if err != nil {
        return domain.IssueReport{}, err
    }

    // Get asset info
    asset, _ := s.assetRepo.GetAssetById(ctx, issueReport.AssetID)

    // Send resolved notification
    titleKey, messageKey, params := messages.IssueResolvedNotification(
        asset.AssetName,
        asset.AssetTag,
        resolutionNotes,
    )

    translations := messages.GetIssueReportNotificationTranslations(titleKey, messageKey, params)

    // Send to reporter
    // ... create and send notification

    return issueReport, nil
}
```

## Adding New Domain Notifications

See [README.md](../internal/notification/messages/README.md) for detailed instructions on adding notifications for new domains.

Quick example:

```go
// internal/notification/messages/maintenance.go
package messages

const (
    NotifMaintenanceScheduledTitleKey   NotificationMessageKey = "notification.maintenance.scheduled.title"
    NotifMaintenanceScheduledMessageKey NotificationMessageKey = "notification.maintenance.scheduled.message"
)

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

func GetMaintenanceNotificationMessage(key NotificationMessageKey, langCode string, params map[string]string) string {
    return GetNotificationMessage(key, langCode, params, maintenanceNotificationTranslations)
}

func GetMaintenanceNotificationTranslations(titleKey, messageKey NotificationMessageKey, params map[string]string) []NotificationTranslation {
    return GetNotificationTranslations(titleKey, messageKey, params, maintenanceNotificationTranslations)
}

func MaintenanceScheduledNotification(assetName, scheduledDate string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
    params := map[string]string{
        "assetName":     assetName,
        "scheduledDate": scheduledDate,
    }
    return NotifMaintenanceScheduledTitleKey, NotifMaintenanceScheduledMessageKey, params
}
```

## Testing

Make sure to test:

1. ✅ All services using asset notifications still work
2. ✅ New issue report notifications function correctly
3. ✅ All translations are returned properly
4. ✅ Placeholders are replaced correctly
5. ✅ Fallback to English works when unsupported language is requested

## Next Steps

To complete the refactoring:

1. ✅ Create `maintenance.go` for maintenance notifications
2. ✅ Create `user.go` for user notifications
3. ✅ Create `category.go` and `location.go` for admin notifications
4. ✅ Update all services to use the new import paths
5. ✅ Remove old `internal/utils/notification_messages.go` file
6. ✅ Update documentation and examples

## Related Files

- `internal/notification/messages/common.go` - Common helpers
- `internal/notification/messages/asset.go` - Asset notifications
- `internal/notification/messages/issue_report.go` - Issue report notifications
- `internal/notification/messages/README.md` - Detailed documentation
- `services/asset/asset_service.go` - Example usage

## Date

Created: 2025-10-15
