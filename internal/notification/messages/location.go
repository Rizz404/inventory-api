package messages

// Location notification message keys
const (
	// Location Created/Updated
	NotifLocationUpdatedTitleKey   NotificationMessageKey = "notification.location.updated.title"
	NotifLocationUpdatedMessageKey NotificationMessageKey = "notification.location.updated.message"
)

// locationNotificationTranslations contains all location notification message translations
var locationNotificationTranslations = map[NotificationMessageKey]map[string]string{
	// ==================== LOCATION UPDATED ====================
	NotifLocationUpdatedTitleKey: {
		"en-US": "Location Updated",
		"ja-JP": "場所が更新されました",
	},
	NotifLocationUpdatedMessageKey: {
		"en-US": "Location \"{locationName}\" has been updated in the system.",
		"ja-JP": "場所 \"{locationName}\" がシステムで更新されました。",
	},
}

// GetLocationNotificationMessage returns the localized location notification message
func GetLocationNotificationMessage(key NotificationMessageKey, langCode string, params map[string]string) string {
	return GetNotificationMessage(key, langCode, params, locationNotificationTranslations)
}

// GetLocationNotificationTranslations returns all translations for a location notification
func GetLocationNotificationTranslations(titleKey, messageKey NotificationMessageKey, params map[string]string) []NotificationTranslation {
	return GetNotificationTranslations(titleKey, messageKey, params, locationNotificationTranslations)
}

// ==================== LOCATION NOTIFICATION HELPER FUNCTIONS ====================

// LocationUpdatedNotification creates notification for location update
func LocationUpdatedNotification(locationName string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"locationName": locationName,
	}
	return NotifLocationUpdatedTitleKey, NotifLocationUpdatedMessageKey, params
}
