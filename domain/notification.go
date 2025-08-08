package domain

import "time"

// --- Enums ---

type NotificationType string

const (
	NotificationTypeMaintenance  NotificationType = "MAINTENANCE"
	NotificationTypeWarranty     NotificationType = "WARRANTY"
	NotificationTypeStatusChange NotificationType = "STATUS_CHANGE"
	NotificationTypeMovement     NotificationType = "MOVEMENT"
	NotificationTypeIssueReport  NotificationType = "ISSUE_REPORT"
)

// --- Structs ---

type Notification struct {
	ID             string                    `json:"id"`
	UserID         string                    `json:"userId"`
	RelatedAssetID *string                   `json:"relatedAssetId"`
	Type           NotificationType          `json:"type"`
	IsRead         bool                      `json:"isRead"`
	CreatedAt      time.Time                 `json:"createdAt"`
	Translations   []NotificationTranslation `json:"translations,omitempty"`
}

type NotificationTranslation struct {
	ID             string `json:"id"`
	NotificationID string `json:"notificationId"`
	LangCode       string `json:"langCode"`
	Title          string `json:"title"`
	Message        string `json:"message"`
}

type NotificationResponse struct {
	ID             string           `json:"id"`
	RelatedAssetID *string          `json:"relatedAssetId,omitempty"`
	Type           NotificationType `json:"type"`
	IsRead         bool             `json:"isRead"`
	CreatedAt      string           `json:"createdAt"`
	Title          string           `json:"title"`
	Message        string           `json:"message"`
}

// --- Payloads ---

// Notifications are typically created by the system, not directly by users.
// Payloads might not be needed for direct API exposure but can be used internally.
type CreateNotificationPayload struct {
	UserID         string                                 `json:"userId"`
	RelatedAssetID *string                                `json:"relatedAssetId,omitempty"`
	Type           NotificationType                       `json:"type"`
	Translations   []CreateNotificationTranslationPayload `json:"translations"`
}

type CreateNotificationTranslationPayload struct {
	LangCode string `json:"langCode"`
	Title    string `json:"title"`
	Message  string `json:"message"`
}
