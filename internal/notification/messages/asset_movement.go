package messages

// Asset Movement notification message keys
const (
	// Asset Location Changed
	NotifAssetMovedTitleKey   NotificationMessageKey = "notification.asset_movement.moved.title"
	NotifAssetMovedMessageKey NotificationMessageKey = "notification.asset_movement.moved.message"

	// Asset User Assigned (via movement)
	NotifAssetUserAssignedTitleKey   NotificationMessageKey = "notification.asset_movement.user_assigned.title"
	NotifAssetUserAssignedMessageKey NotificationMessageKey = "notification.asset_movement.user_assigned.message"
)

// assetMovementNotificationTranslations contains all asset movement notification message translations
var assetMovementNotificationTranslations = map[NotificationMessageKey]map[string]string{
	// ==================== ASSET LOCATION MOVED ====================
	NotifAssetMovedTitleKey: {
		"en-US": "Asset Location Changed",
		"id-ID": "Lokasi Aset Berubah",
		"ja-JP": "資産の場所が変更されました",
	},
	NotifAssetMovedMessageKey: {
		"en-US": "Asset \"{assetName}\" ({assetTag}) has been moved from \"{oldLocation}\" to \"{newLocation}\".",
		"id-ID": "Aset \"{assetName}\" ({assetTag}) telah dipindahkan dari \"{oldLocation}\" ke \"{newLocation}\".",
		"ja-JP": "資産 \"{assetName}\" ({assetTag}) が \"{oldLocation}\" から \"{newLocation}\" に移動されました。",
	},

	// ==================== ASSET USER ASSIGNED ====================
	NotifAssetUserAssignedTitleKey: {
		"en-US": "Asset Assigned to You",
		"id-ID": "Aset Ditugaskan kepada Anda",
		"ja-JP": "資産があなたに割り当てられました",
	},
	NotifAssetUserAssignedMessageKey: {
		"en-US": "Asset \"{assetName}\" ({assetTag}) has been assigned from \"{oldUser}\" to you.",
		"id-ID": "Aset \"{assetName}\" ({assetTag}) telah ditugaskan dari \"{oldUser}\" kepada Anda.",
		"ja-JP": "資産 \"{assetName}\" ({assetTag}) が \"{oldUser}\" からあなたに割り当てられました。",
	},
}

// GetAssetMovementNotificationMessage returns the localized asset movement notification message
func GetAssetMovementNotificationMessage(key NotificationMessageKey, langCode string, params map[string]string) string {
	return GetNotificationMessage(key, langCode, params, assetMovementNotificationTranslations)
}

// GetAssetMovementNotificationTranslations returns all translations for an asset movement notification
func GetAssetMovementNotificationTranslations(titleKey, messageKey NotificationMessageKey, params map[string]string) []NotificationTranslation {
	return GetNotificationTranslations(titleKey, messageKey, params, assetMovementNotificationTranslations)
}

// ==================== ASSET MOVEMENT NOTIFICATION HELPER FUNCTIONS ====================

// AssetMovedNotification creates notification for asset location movement
func AssetMovedNotification(assetName, assetTag, oldLocation, newLocation string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName":   assetName,
		"assetTag":    assetTag,
		"oldLocation": oldLocation,
		"newLocation": newLocation,
	}
	return NotifAssetMovedTitleKey, NotifAssetMovedMessageKey, params
}

// AssetUserAssignedNotification creates notification for asset user assignment
func AssetUserAssignedNotification(assetName, assetTag, oldUser, newUser string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName": assetName,
		"assetTag":  assetTag,
		"oldUser":   oldUser,
		"newUser":   newUser,
	}
	return NotifAssetUserAssignedTitleKey, NotifAssetUserAssignedMessageKey, params
}
