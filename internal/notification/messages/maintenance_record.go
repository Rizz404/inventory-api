package messages

// Maintenance Record notification message keys
const (
	// Maintenance Performed
	NotifMaintenanceCompletedTitleKey   NotificationMessageKey = "notification.maintenance_record.completed.title"
	NotifMaintenanceCompletedMessageKey NotificationMessageKey = "notification.maintenance_record.completed.message"

	// Maintenance Failed
	NotifMaintenanceFailedTitleKey   NotificationMessageKey = "notification.maintenance_record.failed.title"
	NotifMaintenanceFailedMessageKey NotificationMessageKey = "notification.maintenance_record.failed.message"
)

// maintenanceRecordNotificationTranslations contains all maintenance record notification message translations
var maintenanceRecordNotificationTranslations = map[NotificationMessageKey]map[string]string{
	// ==================== MAINTENANCE COMPLETED ====================
	NotifMaintenanceCompletedTitleKey: {
		"en-US": "Maintenance Completed",
		"ja-JP": "メンテナンスが完了しました",
	},
	NotifMaintenanceCompletedMessageKey: {
		"en-US": "Maintenance for asset \"{assetName}\" has been completed. Notes: \"{notes}\".",
		"ja-JP": "資産 \"{assetName}\" のメンテナンスが完了しました。メモ: \"{notes}\"。",
	},

	// ==================== MAINTENANCE FAILED ====================
	NotifMaintenanceFailedTitleKey: {
		"en-US": "Maintenance Failed",
		"ja-JP": "メンテナンスが失敗しました",
	},
	NotifMaintenanceFailedMessageKey: {
		"en-US": "Maintenance for asset \"{assetName}\" could not be completed. Reason: \"{failureReason}\".",
		"ja-JP": "資産 \"{assetName}\" のメンテナンスが完了できませんでした。理由: \"{failureReason}\"。",
	},
}

// GetMaintenanceRecordNotificationMessage returns the localized maintenance record notification message
func GetMaintenanceRecordNotificationMessage(key NotificationMessageKey, langCode string, params map[string]string) string {
	return GetNotificationMessage(key, langCode, params, maintenanceRecordNotificationTranslations)
}

// GetMaintenanceRecordNotificationTranslations returns all translations for a maintenance record notification
func GetMaintenanceRecordNotificationTranslations(titleKey, messageKey NotificationMessageKey, params map[string]string) []NotificationTranslation {
	return GetNotificationTranslations(titleKey, messageKey, params, maintenanceRecordNotificationTranslations)
}

// ==================== MAINTENANCE RECORD NOTIFICATION HELPER FUNCTIONS ====================

// MaintenanceCompletedNotification creates notification for completed maintenance
func MaintenanceCompletedNotification(assetName, assetTag, notes string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName": assetName,
		"assetTag":  assetTag,
		"notes":     notes,
	}
	return NotifMaintenanceCompletedTitleKey, NotifMaintenanceCompletedMessageKey, params
}

// MaintenanceFailedNotification creates notification for failed maintenance
func MaintenanceFailedNotification(assetName, assetTag, failureReason string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName":     assetName,
		"assetTag":      assetTag,
		"failureReason": failureReason,
	}
	return NotifMaintenanceFailedTitleKey, NotifMaintenanceFailedMessageKey, params
}
