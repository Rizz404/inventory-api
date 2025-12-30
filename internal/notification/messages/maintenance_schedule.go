package messages

// Maintenance Schedule notification message keys
const (
	// Maintenance Scheduled
	NotifMaintenanceScheduledTitleKey   NotificationMessageKey = "notification.maintenance_schedule.scheduled.title"
	NotifMaintenanceScheduledMessageKey NotificationMessageKey = "notification.maintenance_schedule.scheduled.message"

	// Maintenance Due Soon
	NotifMaintenanceDueSoonTitleKey   NotificationMessageKey = "notification.maintenance_schedule.due_soon.title"
	NotifMaintenanceDueSoonMessageKey NotificationMessageKey = "notification.maintenance_schedule.due_soon.message"

	// Maintenance Overdue
	NotifMaintenanceOverdueTitleKey   NotificationMessageKey = "notification.maintenance_schedule.overdue.title"
	NotifMaintenanceOverdueMessageKey NotificationMessageKey = "notification.maintenance_schedule.overdue.message"
)

// maintenanceScheduleNotificationTranslations contains all maintenance schedule notification message translations
var maintenanceScheduleNotificationTranslations = map[NotificationMessageKey]map[string]string{
	// ==================== MAINTENANCE SCHEDULED ====================
	NotifMaintenanceScheduledTitleKey: {
		"en-US": "Maintenance Scheduled",
		"id-ID": "Pemeliharaan Dijadwalkan",
		"ja-JP": "メンテナンスがスケジュールされました",
	},
	NotifMaintenanceScheduledMessageKey: {
		"en-US": "Maintenance for asset \"{assetName}\" is scheduled on {scheduledDate}.",
		"id-ID": "Pemeliharaan untuk aset \"{assetName}\" dijadwalkan pada {scheduledDate}.",
		"ja-JP": "資産 \"{assetName}\" のメンテナンスが {scheduledDate} にスケジュールされました。",
	},

	// ==================== MAINTENANCE DUE SOON ====================
	NotifMaintenanceDueSoonTitleKey: {
		"en-US": "Maintenance Due Soon",
		"id-ID": "Pemeliharaan Segera Jatuh Tempo",
		"ja-JP": "メンテナンス期限が近づいています",
	},
	NotifMaintenanceDueSoonMessageKey: {
		"en-US": "Maintenance for asset \"{assetName}\" is due on {scheduledDate}. Please prepare.",
		"id-ID": "Pemeliharaan untuk aset \"{assetName}\" jatuh tempo pada {scheduledDate}. Silakan persiapkan.",
		"ja-JP": "資産 \"{assetName}\" のメンテナンス期限が {scheduledDate} です。準備してください。",
	},

	// ==================== MAINTENANCE OVERDUE ====================
	NotifMaintenanceOverdueTitleKey: {
		"en-US": "Maintenance Overdue",
		"id-ID": "Pemeliharaan Terlambat",
		"ja-JP": "メンテナンスが期限切れです",
	},
	NotifMaintenanceOverdueMessageKey: {
		"en-US": "Maintenance for asset \"{assetName}\" is overdue. Scheduled date was {scheduledDate}.",
		"id-ID": "Pemeliharaan untuk aset \"{assetName}\" sudah terlambat. Tanggal yang dijadwalkan adalah {scheduledDate}.",
		"ja-JP": "資産 \"{assetName}\" のメンテナンスが期限切れです。スケジュールされた日付は {scheduledDate} でした。",
	},
}

// GetMaintenanceScheduleNotificationMessage returns the localized maintenance schedule notification message
func GetMaintenanceScheduleNotificationMessage(key NotificationMessageKey, langCode string, params map[string]string) string {
	return GetNotificationMessage(key, langCode, params, maintenanceScheduleNotificationTranslations)
}

// GetMaintenanceScheduleNotificationTranslations returns all translations for a maintenance schedule notification
func GetMaintenanceScheduleNotificationTranslations(titleKey, messageKey NotificationMessageKey, params map[string]string) []NotificationTranslation {
	return GetNotificationTranslations(titleKey, messageKey, params, maintenanceScheduleNotificationTranslations)
}

// ==================== MAINTENANCE SCHEDULE NOTIFICATION HELPER FUNCTIONS ====================

// MaintenanceScheduledNotification creates notification for new maintenance schedule
func MaintenanceScheduledNotification(assetName, assetTag, scheduledDate string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName":     assetName,
		"assetTag":      assetTag,
		"scheduledDate": scheduledDate,
	}
	return NotifMaintenanceScheduledTitleKey, NotifMaintenanceScheduledMessageKey, params
}

// MaintenanceDueSoonNotification creates notification for maintenance due soon
func MaintenanceDueSoonNotification(assetName, assetTag, scheduledDate string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName":     assetName,
		"assetTag":      assetTag,
		"scheduledDate": scheduledDate,
	}
	return NotifMaintenanceDueSoonTitleKey, NotifMaintenanceDueSoonMessageKey, params
}

// MaintenanceOverdueNotification creates notification for overdue maintenance
func MaintenanceOverdueNotification(assetName, assetTag, scheduledDate string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName":     assetName,
		"assetTag":      assetTag,
		"scheduledDate": scheduledDate,
	}
	return NotifMaintenanceOverdueTitleKey, NotifMaintenanceOverdueMessageKey, params
}
