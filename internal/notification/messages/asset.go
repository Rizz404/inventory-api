package messages

// Asset notification message keys
const (
	// Asset Assignment
	NotifAssetAssignedTitleKey   NotificationMessageKey = "notification.asset.assigned.title"
	NotifAssetAssignedMessageKey NotificationMessageKey = "notification.asset.assigned.message"

	NotifAssetNewAssignedTitleKey   NotificationMessageKey = "notification.asset.new_assigned.title"
	NotifAssetNewAssignedMessageKey NotificationMessageKey = "notification.asset.new_assigned.message"

	NotifAssetUnassignedTitleKey   NotificationMessageKey = "notification.asset.unassigned.title"
	NotifAssetUnassignedMessageKey NotificationMessageKey = "notification.asset.unassigned.message"

	// Asset Status Change
	NotifAssetStatusChangedTitleKey   NotificationMessageKey = "notification.asset.status_changed.title"
	NotifAssetStatusChangedMessageKey NotificationMessageKey = "notification.asset.status_changed.message"

	NotifAssetActivatedTitleKey   NotificationMessageKey = "notification.asset.activated.title"
	NotifAssetActivatedMessageKey NotificationMessageKey = "notification.asset.activated.message"

	NotifAssetMaintenanceTitleKey   NotificationMessageKey = "notification.asset.maintenance.title"
	NotifAssetMaintenanceMessageKey NotificationMessageKey = "notification.asset.maintenance.message"

	NotifAssetDisposedTitleKey   NotificationMessageKey = "notification.asset.disposed.title"
	NotifAssetDisposedMessageKey NotificationMessageKey = "notification.asset.disposed.message"

	NotifAssetLostTitleKey   NotificationMessageKey = "notification.asset.lost.title"
	NotifAssetLostMessageKey NotificationMessageKey = "notification.asset.lost.message"

	// Asset Condition Change
	NotifAssetConditionChangedTitleKey   NotificationMessageKey = "notification.asset.condition_changed.title"
	NotifAssetConditionChangedMessageKey NotificationMessageKey = "notification.asset.condition_changed.message"

	NotifAssetConditionDamagedTitleKey   NotificationMessageKey = "notification.asset.condition_damaged.title"
	NotifAssetConditionDamagedMessageKey NotificationMessageKey = "notification.asset.condition_damaged.message"

	NotifAssetConditionPoorTitleKey   NotificationMessageKey = "notification.asset.condition_poor.title"
	NotifAssetConditionPoorMessageKey NotificationMessageKey = "notification.asset.condition_poor.message"

	// Asset Location Change
	NotifAssetLocationChangedTitleKey   NotificationMessageKey = "notification.asset.location_changed.title"
	NotifAssetLocationChangedMessageKey NotificationMessageKey = "notification.asset.location_changed.message"

	// Asset Creation/Deletion
	NotifAssetCreatedTitleKey   NotificationMessageKey = "notification.asset.created.title"
	NotifAssetCreatedMessageKey NotificationMessageKey = "notification.asset.created.message"

	NotifAssetDeletedTitleKey   NotificationMessageKey = "notification.asset.deleted.title"
	NotifAssetDeletedMessageKey NotificationMessageKey = "notification.asset.deleted.message"

	// Asset Warranty
	NotifAssetWarrantyExpiringSoonTitleKey   NotificationMessageKey = "notification.asset.warranty_expiring_soon.title"
	NotifAssetWarrantyExpiringSoonMessageKey NotificationMessageKey = "notification.asset.warranty_expiring_soon.message"

	NotifAssetWarrantyExpiredTitleKey   NotificationMessageKey = "notification.asset.warranty_expired.title"
	NotifAssetWarrantyExpiredMessageKey NotificationMessageKey = "notification.asset.warranty_expired.message"

	// Asset Value/Purchase
	NotifAssetHighValueTitleKey   NotificationMessageKey = "notification.asset.high_value.title"
	NotifAssetHighValueMessageKey NotificationMessageKey = "notification.asset.high_value.message"
)

// assetNotificationTranslations contains all asset notification message translations
var assetNotificationTranslations = map[NotificationMessageKey]map[string]string{
	// ==================== ASSET ASSIGNMENT ====================
	NotifAssetAssignedTitleKey: {
		"en-US": "Asset Assigned",
		"ja-JP": "資産が割り当てられました",
	},
	NotifAssetAssignedMessageKey: {
		"en-US": "Asset \"{assetName}\" has been assigned to you.",
		"ja-JP": "資産 \"{assetName}\" があなたに割り当てられました。",
	},

	NotifAssetNewAssignedTitleKey: {
		"en-US": "New Asset Assigned",
		"ja-JP": "新しい資産が割り当てられました",
	},
	NotifAssetNewAssignedMessageKey: {
		"en-US": "New asset \"{assetName}\" has been assigned to you.",
		"ja-JP": "新しい資産 \"{assetName}\" があなたに割り当てられました。",
	},

	NotifAssetUnassignedTitleKey: {
		"en-US": "Asset Unassigned",
		"ja-JP": "資産の割り当てが解除されました",
	},
	NotifAssetUnassignedMessageKey: {
		"en-US": "Asset \"{assetName}\" has been unassigned from you.",
		"ja-JP": "資産 \"{assetName}\" の割り当てがあなたから解除されました。",
	},

	// ==================== ASSET STATUS CHANGE ====================
	NotifAssetStatusChangedTitleKey: {
		"en-US": "Asset Status Changed",
		"ja-JP": "資産ステータスが変更されました",
	},
	NotifAssetStatusChangedMessageKey: {
		"en-US": "Asset \"{assetName}\" status changed from {oldStatus} to {newStatus}.",
		"ja-JP": "資産 \"{assetName}\" のステータスが {oldStatus} から {newStatus} に変更されました。",
	},

	NotifAssetActivatedTitleKey: {
		"en-US": "Asset Activated",
		"ja-JP": "資産が有効化されました",
	},
	NotifAssetActivatedMessageKey: {
		"en-US": "Asset \"{assetName}\" is now active and ready to use.",
		"ja-JP": "資産 \"{assetName}\" が有効化され、使用準備が整いました。",
	},

	NotifAssetMaintenanceTitleKey: {
		"en-US": "Asset Under Maintenance",
		"ja-JP": "資産がメンテナンス中です",
	},
	NotifAssetMaintenanceMessageKey: {
		"en-US": "Asset \"{assetName}\" has been moved to maintenance status.",
		"ja-JP": "資産 \"{assetName}\" がメンテナンスステータスに移動されました。",
	},

	NotifAssetDisposedTitleKey: {
		"en-US": "Asset Disposed",
		"ja-JP": "資産が廃棄されました",
	},
	NotifAssetDisposedMessageKey: {
		"en-US": "Asset \"{assetName}\" has been disposed.",
		"ja-JP": "資産 \"{assetName}\" が廃棄されました。",
	},

	NotifAssetLostTitleKey: {
		"en-US": "Asset Reported Lost",
		"ja-JP": "資産が行方不明として報告されました",
	},
	NotifAssetLostMessageKey: {
		"en-US": "Asset \"{assetName}\" has been reported as lost.",
		"ja-JP": "資産 \"{assetName}\" が行方不明として報告されました。",
	},

	// ==================== ASSET CONDITION CHANGE ====================
	NotifAssetConditionChangedTitleKey: {
		"en-US": "Asset Condition Changed",
		"ja-JP": "資産の状態が変更されました",
	},
	NotifAssetConditionChangedMessageKey: {
		"en-US": "Asset \"{assetName}\" condition changed from {oldCondition} to {newCondition}.",
		"ja-JP": "資産 \"{assetName}\" の状態が {oldCondition} から {newCondition} に変更されました。",
	},

	NotifAssetConditionDamagedTitleKey: {
		"en-US": "Asset Damaged",
		"ja-JP": "資産が損傷しました",
	},
	NotifAssetConditionDamagedMessageKey: {
		"en-US": "Asset \"{assetName}\" has been marked as damaged. Please check immediately.",
		"ja-JP": "資産 \"{assetName}\" が損傷としてマークされました。すぐに確認してください。",
	},

	NotifAssetConditionPoorTitleKey: {
		"en-US": "Asset in Poor Condition",
		"ja-JP": "資産の状態が不良です",
	},
	NotifAssetConditionPoorMessageKey: {
		"en-US": "Asset \"{assetName}\" condition has deteriorated to poor. Maintenance may be needed.",
		"ja-JP": "資産 \"{assetName}\" の状態が不良に悪化しました。メンテナンスが必要かもしれません。",
	},

	// ==================== ASSET LOCATION CHANGE ====================
	NotifAssetLocationChangedTitleKey: {
		"en-US": "Asset Location Changed",
		"ja-JP": "資産の場所が変更されました",
	},
	NotifAssetLocationChangedMessageKey: {
		"en-US": "Asset \"{assetName}\" has been moved from {oldLocation} to {newLocation}.",
		"ja-JP": "資産 \"{assetName}\" が {oldLocation} から {newLocation} に移動されました。",
	},

	// ==================== ASSET CREATION/DELETION ====================
	NotifAssetCreatedTitleKey: {
		"en-US": "New Asset Created",
		"ja-JP": "新しい資産が作成されました",
	},
	NotifAssetCreatedMessageKey: {
		"en-US": "New asset \"{assetName}\" has been added to the inventory.",
		"ja-JP": "新しい資産 \"{assetName}\" が在庫に追加されました。",
	},

	NotifAssetDeletedTitleKey: {
		"en-US": "Asset Deleted",
		"ja-JP": "資産が削除されました",
	},
	NotifAssetDeletedMessageKey: {
		"en-US": "Asset \"{assetName}\" has been removed from the inventory.",
		"ja-JP": "資産 \"{assetName}\" が在庫から削除されました。",
	},

	// ==================== ASSET WARRANTY ====================
	NotifAssetWarrantyExpiringSoonTitleKey: {
		"en-US": "Warranty Expiring Soon",
		"ja-JP": "保証期間がまもなく終了します",
	},
	NotifAssetWarrantyExpiringSoonMessageKey: {
		"en-US": "Warranty for asset \"{assetName}\" will expire on {expiryDate}.",
		"ja-JP": "資産 \"{assetName}\" の保証期間が {expiryDate} に終了します。",
	},

	NotifAssetWarrantyExpiredTitleKey: {
		"en-US": "Warranty Expired",
		"ja-JP": "保証期間が終了しました",
	},
	NotifAssetWarrantyExpiredMessageKey: {
		"en-US": "Warranty for asset \"{assetName}\" has expired.",
		"ja-JP": "資産 \"{assetName}\" の保証期間が終了しました。",
	},

	// ==================== ASSET VALUE/PURCHASE ====================
	NotifAssetHighValueTitleKey: {
		"en-US": "High Value Asset Added",
		"ja-JP": "高額資産が追加されました",
	},
	NotifAssetHighValueMessageKey: {
		"en-US": "High value asset \"{assetName}\" worth {value} has been added to your inventory.",
		"ja-JP": "高額資産 \"{assetName}\" 価値 {value} が在庫に追加されました。",
	},
}

// GetAssetNotificationMessage returns the localized asset notification message
func GetAssetNotificationMessage(key NotificationMessageKey, langCode string, params map[string]string) string {
	return GetNotificationMessage(key, langCode, params, assetNotificationTranslations)
}

// GetAssetNotificationTranslations returns all translations for an asset notification
func GetAssetNotificationTranslations(titleKey, messageKey NotificationMessageKey, params map[string]string) []NotificationTranslation {
	return GetNotificationTranslations(titleKey, messageKey, params, assetNotificationTranslations)
}

// ==================== ASSET NOTIFICATION HELPER FUNCTIONS ====================

// AssetAssignmentNotification creates notification for asset assignment
func AssetAssignmentNotification(assetName, assetTag string, isNew bool) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName": assetName,
		"assetTag":  assetTag,
	}

	if isNew {
		return NotifAssetNewAssignedTitleKey, NotifAssetNewAssignedMessageKey, params
	}
	return NotifAssetAssignedTitleKey, NotifAssetAssignedMessageKey, params
}

// AssetUnassignmentNotification creates notification for asset unassignment
func AssetUnassignmentNotification(assetName, assetTag string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName": assetName,
		"assetTag":  assetTag,
	}
	return NotifAssetUnassignedTitleKey, NotifAssetUnassignedMessageKey, params
}

// AssetStatusChangeNotification creates notification for asset status change
func AssetStatusChangeNotification(assetName string, oldStatus, newStatus string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName": assetName,
		"oldStatus": oldStatus,
		"newStatus": newStatus,
	}

	// Check for specific status changes
	switch newStatus {
	case "Active":
		params["assetTag"] = assetName // You may need to pass this separately
		return NotifAssetActivatedTitleKey, NotifAssetActivatedMessageKey, params
	case "Maintenance":
		params["assetTag"] = assetName
		return NotifAssetMaintenanceTitleKey, NotifAssetMaintenanceMessageKey, params
	case "Disposed":
		params["assetTag"] = assetName
		return NotifAssetDisposedTitleKey, NotifAssetDisposedMessageKey, params
	case "Lost":
		params["assetTag"] = assetName
		return NotifAssetLostTitleKey, NotifAssetLostMessageKey, params
	default:
		return NotifAssetStatusChangedTitleKey, NotifAssetStatusChangedMessageKey, params
	}
}

// AssetConditionChangeNotification creates notification for asset condition change
func AssetConditionChangeNotification(assetName, assetTag string, oldCondition, newCondition string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName":    assetName,
		"assetTag":     assetTag,
		"oldCondition": oldCondition,
		"newCondition": newCondition,
	}

	// Check for specific condition changes
	switch newCondition {
	case "Damaged":
		return NotifAssetConditionDamagedTitleKey, NotifAssetConditionDamagedMessageKey, params
	case "Poor":
		return NotifAssetConditionPoorTitleKey, NotifAssetConditionPoorMessageKey, params
	default:
		return NotifAssetConditionChangedTitleKey, NotifAssetConditionChangedMessageKey, params
	}
}

// AssetLocationChangeNotification creates notification for asset location change
func AssetLocationChangeNotification(assetName, assetTag, oldLocation, newLocation string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName":   assetName,
		"assetTag":    assetTag,
		"oldLocation": oldLocation,
		"newLocation": newLocation,
	}
	return NotifAssetLocationChangedTitleKey, NotifAssetLocationChangedMessageKey, params
}

// AssetCreatedNotification creates notification for new asset creation
func AssetCreatedNotification(assetName, assetTag string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName": assetName,
		"assetTag":  assetTag,
	}
	return NotifAssetCreatedTitleKey, NotifAssetCreatedMessageKey, params
}

// AssetDeletedNotification creates notification for asset deletion
func AssetDeletedNotification(assetName, assetTag string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName": assetName,
		"assetTag":  assetTag,
	}
	return NotifAssetDeletedTitleKey, NotifAssetDeletedMessageKey, params
}

// AssetWarrantyExpiringNotification creates notification for warranty expiring soon
func AssetWarrantyExpiringNotification(assetName, assetTag, expiryDate string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName":  assetName,
		"assetTag":   assetTag,
		"expiryDate": expiryDate,
	}
	return NotifAssetWarrantyExpiringSoonTitleKey, NotifAssetWarrantyExpiringSoonMessageKey, params
}

// AssetWarrantyExpiredNotification creates notification for expired warranty
func AssetWarrantyExpiredNotification(assetName, assetTag string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName": assetName,
		"assetTag":  assetTag,
	}
	return NotifAssetWarrantyExpiredTitleKey, NotifAssetWarrantyExpiredMessageKey, params
}

// AssetHighValueNotification creates notification for high value asset
func AssetHighValueNotification(assetName, assetTag, value string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName": assetName,
		"assetTag":  assetTag,
		"value":     value,
	}
	return NotifAssetHighValueTitleKey, NotifAssetHighValueMessageKey, params
}
