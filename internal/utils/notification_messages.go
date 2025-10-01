package utils

import (
	"fmt"
)

// * NotificationMessageKey represents a notification message key
type NotificationMessageKey string

// * Asset notification message keys
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

// * notificationMessageTranslations contains all notification message translations
var notificationMessageTranslations = map[NotificationMessageKey]map[string]string{
	// ==================== ASSET ASSIGNMENT ====================
	NotifAssetAssignedTitleKey: {
		"en":    "Asset Assigned",
		"id-ID": "Aset Ditugaskan",
	},
	NotifAssetAssignedMessageKey: {
		"en":    "Asset '{assetName}' (Tag: {assetTag}) has been assigned to you.",
		"id-ID": "Aset '{assetName}' (Tag: {assetTag}) telah ditugaskan kepada Anda.",
	},

	NotifAssetNewAssignedTitleKey: {
		"en":    "New Asset Assigned",
		"id-ID": "Aset Baru Ditugaskan",
	},
	NotifAssetNewAssignedMessageKey: {
		"en":    "New asset '{assetName}' (Tag: {assetTag}) has been assigned to you.",
		"id-ID": "Aset baru '{assetName}' (Tag: {assetTag}) telah ditugaskan kepada Anda.",
	},

	NotifAssetUnassignedTitleKey: {
		"en":    "Asset Unassigned",
		"id-ID": "Aset Tidak Ditugaskan Lagi",
	},
	NotifAssetUnassignedMessageKey: {
		"en":    "Asset '{assetName}' (Tag: {assetTag}) has been unassigned from you.",
		"id-ID": "Aset '{assetName}' (Tag: {assetTag}) telah dilepaskan dari Anda.",
	},

	// ==================== ASSET STATUS CHANGE ====================
	NotifAssetStatusChangedTitleKey: {
		"en":    "Asset Status Changed",
		"id-ID": "Status Aset Berubah",
	},
	NotifAssetStatusChangedMessageKey: {
		"en":    "Asset '{assetName}' status changed from {oldStatus} to {newStatus}.",
		"id-ID": "Status aset '{assetName}' berubah dari {oldStatus} menjadi {newStatus}.",
	},

	NotifAssetActivatedTitleKey: {
		"en":    "Asset Activated",
		"id-ID": "Aset Diaktifkan",
	},
	NotifAssetActivatedMessageKey: {
		"en":    "Asset '{assetName}' (Tag: {assetTag}) is now active and ready to use.",
		"id-ID": "Aset '{assetName}' (Tag: {assetTag}) sekarang aktif dan siap digunakan.",
	},

	NotifAssetMaintenanceTitleKey: {
		"en":    "Asset Under Maintenance",
		"id-ID": "Aset Dalam Pemeliharaan",
	},
	NotifAssetMaintenanceMessageKey: {
		"en":    "Asset '{assetName}' (Tag: {assetTag}) has been moved to maintenance status.",
		"id-ID": "Aset '{assetName}' (Tag: {assetTag}) telah dipindahkan ke status pemeliharaan.",
	},

	NotifAssetDisposedTitleKey: {
		"en":    "Asset Disposed",
		"id-ID": "Aset Dibuang",
	},
	NotifAssetDisposedMessageKey: {
		"en":    "Asset '{assetName}' (Tag: {assetTag}) has been disposed.",
		"id-ID": "Aset '{assetName}' (Tag: {assetTag}) telah dibuang.",
	},

	NotifAssetLostTitleKey: {
		"en":    "Asset Reported Lost",
		"id-ID": "Aset Dilaporkan Hilang",
	},
	NotifAssetLostMessageKey: {
		"en":    "Asset '{assetName}' (Tag: {assetTag}) has been reported as lost.",
		"id-ID": "Aset '{assetName}' (Tag: {assetTag}) telah dilaporkan hilang.",
	},

	// ==================== ASSET CONDITION CHANGE ====================
	NotifAssetConditionChangedTitleKey: {
		"en":    "Asset Condition Changed",
		"id-ID": "Kondisi Aset Berubah",
	},
	NotifAssetConditionChangedMessageKey: {
		"en":    "Asset '{assetName}' condition changed from {oldCondition} to {newCondition}.",
		"id-ID": "Kondisi aset '{assetName}' berubah dari {oldCondition} menjadi {newCondition}.",
	},

	NotifAssetConditionDamagedTitleKey: {
		"en":    "Asset Damaged",
		"id-ID": "Aset Rusak",
	},
	NotifAssetConditionDamagedMessageKey: {
		"en":    "Asset '{assetName}' (Tag: {assetTag}) has been marked as damaged. Please check immediately.",
		"id-ID": "Aset '{assetName}' (Tag: {assetTag}) telah ditandai sebagai rusak. Mohon segera dicek.",
	},

	NotifAssetConditionPoorTitleKey: {
		"en":    "Asset in Poor Condition",
		"id-ID": "Aset Kondisi Buruk",
	},
	NotifAssetConditionPoorMessageKey: {
		"en":    "Asset '{assetName}' (Tag: {assetTag}) condition has deteriorated to poor. Maintenance may be needed.",
		"id-ID": "Aset '{assetName}' (Tag: {assetTag}) kondisinya memburuk. Mungkin perlu pemeliharaan.",
	},

	// ==================== ASSET LOCATION CHANGE ====================
	NotifAssetLocationChangedTitleKey: {
		"en":    "Asset Location Changed",
		"id-ID": "Lokasi Aset Berubah",
	},
	NotifAssetLocationChangedMessageKey: {
		"en":    "Asset '{assetName}' has been moved from {oldLocation} to {newLocation}.",
		"id-ID": "Aset '{assetName}' telah dipindahkan dari {oldLocation} ke {newLocation}.",
	},

	// ==================== ASSET CREATION/DELETION ====================
	NotifAssetCreatedTitleKey: {
		"en":    "New Asset Created",
		"id-ID": "Aset Baru Dibuat",
	},
	NotifAssetCreatedMessageKey: {
		"en":    "New asset '{assetName}' (Tag: {assetTag}) has been added to the inventory.",
		"id-ID": "Aset baru '{assetName}' (Tag: {assetTag}) telah ditambahkan ke inventaris.",
	},

	NotifAssetDeletedTitleKey: {
		"en":    "Asset Deleted",
		"id-ID": "Aset Dihapus",
	},
	NotifAssetDeletedMessageKey: {
		"en":    "Asset '{assetName}' (Tag: {assetTag}) has been removed from the inventory.",
		"id-ID": "Aset '{assetName}' (Tag: {assetTag}) telah dihapus dari inventaris.",
	},

	// ==================== ASSET WARRANTY ====================
	NotifAssetWarrantyExpiringSoonTitleKey: {
		"en":    "Warranty Expiring Soon",
		"id-ID": "Garansi Akan Habis",
	},
	NotifAssetWarrantyExpiringSoonMessageKey: {
		"en":    "Warranty for asset '{assetName}' (Tag: {assetTag}) will expire on {expiryDate}.",
		"id-ID": "Garansi untuk aset '{assetName}' (Tag: {assetTag}) akan berakhir pada {expiryDate}.",
	},

	NotifAssetWarrantyExpiredTitleKey: {
		"en":    "Warranty Expired",
		"id-ID": "Garansi Telah Habis",
	},
	NotifAssetWarrantyExpiredMessageKey: {
		"en":    "Warranty for asset '{assetName}' (Tag: {assetTag}) has expired.",
		"id-ID": "Garansi untuk aset '{assetName}' (Tag: {assetTag}) telah habis.",
	},

	// ==================== ASSET VALUE/PURCHASE ====================
	NotifAssetHighValueTitleKey: {
		"en":    "High Value Asset Added",
		"id-ID": "Aset Bernilai Tinggi Ditambahkan",
	},
	NotifAssetHighValueMessageKey: {
		"en":    "High value asset '{assetName}' (Tag: {assetTag}) worth {value} has been added to your inventory.",
		"id-ID": "Aset bernilai tinggi '{assetName}' (Tag: {assetTag}) senilai {value} telah ditambahkan ke inventaris Anda.",
	},
}

// * GetNotificationMessage returns the localized notification message
func GetNotificationMessage(key NotificationMessageKey, langCode string, params map[string]string) string {
	translations, exists := notificationMessageTranslations[key]
	if !exists {
		return string(key)
	}

	normalizedLang := normalizeLanguageCode(langCode)
	message, exists := translations[normalizedLang]
	if !exists {
		// Fallback to English
		message, exists = translations["en"]
		if !exists {
			return string(key)
		}
	}

	// Replace placeholders with actual values
	for placeholder, value := range params {
		message = replaceAllCaseInsensitive(message, "{"+placeholder+"}", value)
	}

	return message
}

// * NotificationTranslation represents a notification translation
type NotificationTranslation struct {
	LangCode string
	Title    string
	Message  string
}

// * GetNotificationTranslations returns all translations for a notification message key
func GetNotificationTranslations(titleKey, messageKey NotificationMessageKey, params map[string]string) []NotificationTranslation {
	translations := []NotificationTranslation{}

	// Get available languages
	languages := []string{"en", "id-ID"}

	for _, lang := range languages {
		title := GetNotificationMessage(titleKey, lang, params)
		message := GetNotificationMessage(messageKey, lang, params)

		translations = append(translations, NotificationTranslation{
			LangCode: lang,
			Title:    title,
			Message:  message,
		})
	}

	return translations
}

// * Helper function to replace all occurrences of a string (case-insensitive)
func replaceAllCaseInsensitive(s, old, new string) string {
	// Simple implementation - for production, consider using regex for case-insensitive replace
	result := s
	for {
		index := -1
		for i := 0; i <= len(result)-len(old); i++ {
			if result[i:i+len(old)] == old {
				index = i
				break
			}
		}
		if index == -1 {
			break
		}
		result = result[:index] + new + result[index+len(old):]
	}
	return result
}

// * Asset notification helper functions

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
func AssetLocationChangeNotification(assetName, oldLocation, newLocation string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName":   assetName,
		"oldLocation": oldLocation,
		"newLocation": newLocation,
	}
	return NotifAssetLocationChangedTitleKey, NotifAssetLocationChangedMessageKey, params
}

// AssetCreationNotification creates notification for new asset creation
func AssetCreationNotification(assetName, assetTag string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
	params := map[string]string{
		"assetName": assetName,
		"assetTag":  assetTag,
	}
	return NotifAssetCreatedTitleKey, NotifAssetCreatedMessageKey, params
}

// AssetDeletionNotification creates notification for asset deletion
func AssetDeletionNotification(assetName, assetTag string) (NotificationMessageKey, NotificationMessageKey, map[string]string) {
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

// * Utility function to format currency
func FormatCurrency(value float64) string {
	// Simple formatting - you can enhance this based on locale
	return fmt.Sprintf("$%.2f", value)
}

// * Utility function to format date
func FormatDate(date string) string {
	// Simple formatting - you can enhance this based on locale
	return date
}
