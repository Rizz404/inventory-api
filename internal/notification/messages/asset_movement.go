package messages

// Asset Movement notification message keys
const (
	// Asset Moved
	NotifAssetMovedTitleKey   NotificationMessageKey = "notification.asset_movement.moved.title"
	NotifAssetMovedMessageKey NotificationMessageKey = "notification.asset_movement.moved.message"
)

// assetMovementNotificationTranslations contains all asset movement notification message translations
var assetMovementNotificationTranslations = map[NotificationMessageKey]map[string]string{
	// ==================== ASSET MOVED ====================
	NotifAssetMovedTitleKey: {
		"en-US": "Asset Moved",
		"ja-JP": "資産が移動されました",
	},
	NotifAssetMovedMessageKey: {
		"en-US": "Asset \"{assetName}\" (\"{assetTag}\") has been moved from \"{oldLocation}\" to \"{newLocation}\".",
		"ja-JP": "資産 \"{assetName}\" (\"{assetTag}\") が \"{oldLocation}\" から \"{newLocation}\" に移動されました。",
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

// AssetMovedNotification creates notification for asset movement
func AssetMovedNotification(assetName, assetTag, oldLocation, newLocation string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName":   assetName,
		"assetTag":    assetTag,
		"oldLocation": oldLocation,
		"newLocation": newLocation,
	}
	return NotifAssetMovedTitleKey, NotifAssetMovedMessageKey, params
}
