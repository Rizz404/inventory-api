package messages

// Category notification message keys
const (
	// Category Created/Updated
	NotifCategoryUpdatedTitleKey   NotificationMessageKey = "notification.category.updated.title"
	NotifCategoryUpdatedMessageKey NotificationMessageKey = "notification.category.updated.message"
)

// categoryNotificationTranslations contains all category notification message translations
var categoryNotificationTranslations = map[NotificationMessageKey]map[string]string{
	// ==================== CATEGORY UPDATED ====================
	NotifCategoryUpdatedTitleKey: {
		"en-US": "Category Updated",
		"id-ID": "Kategori Diperbarui",
		"ja-JP": "カテゴリが更新されました",
	},
	NotifCategoryUpdatedMessageKey: {
		"en-US": "Category \"{categoryName}\" has been updated.",
		"id-ID": "Kategori \"{categoryName}\" telah diperbarui.",
		"ja-JP": "カテゴリ \"{categoryName}\" が更新されました。",
	},
}

// GetCategoryNotificationMessage returns the localized category notification message
func GetCategoryNotificationMessage(key NotificationMessageKey, langCode string, params map[string]string) string {
	return GetNotificationMessage(key, langCode, params, categoryNotificationTranslations)
}

// GetCategoryNotificationTranslations returns all translations for a category notification
func GetCategoryNotificationTranslations(titleKey, messageKey NotificationMessageKey, params map[string]string) []NotificationTranslation {
	return GetNotificationTranslations(titleKey, messageKey, params, categoryNotificationTranslations)
}

// ==================== CATEGORY NOTIFICATION HELPER FUNCTIONS ====================

// CategoryUpdatedNotification creates notification for category update
func CategoryUpdatedNotification(categoryName string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"categoryName": categoryName,
	}
	return NotifCategoryUpdatedTitleKey, NotifCategoryUpdatedMessageKey, params
}
