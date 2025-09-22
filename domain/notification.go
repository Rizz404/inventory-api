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
	UserID         string           `json:"userId"`
	RelatedAssetID *string          `json:"relatedAssetId,omitempty"`
	Type           NotificationType `json:"type"`
	IsRead         bool             `json:"isRead"`
	CreatedAt      string           `json:"createdAt"`
	Title          string           `json:"title"`
	Message        string           `json:"message"`
	// * Populated
	// ! cuma notification gak perlu populated table biar gak berat
	// User         UserResponse   `json:"user"`
	// RelatedAsset *AssetResponse `json:"relatedAsset,omitempty"`
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

type UpdateNotificationPayload struct {
	IsRead       *bool                                  `json:"isRead,omitempty"`
	Translations []CreateNotificationTranslationPayload `json:"translations,omitempty"`
}

// --- Statistics ---

// Internal statistics structs (used in repository layer)
type NotificationStatistics struct {
	Total          NotificationCountStatistics   `json:"total"`
	ByType         NotificationTypeStatistics    `json:"byType"`
	ByStatus       NotificationStatusStatistics  `json:"byStatus"`
	CreationTrends []NotificationCreationTrend   `json:"creationTrends"`
	Summary        NotificationSummaryStatistics `json:"summary"`
}

type NotificationCountStatistics struct {
	Count int `json:"count"`
}

type NotificationTypeStatistics struct {
	Maintenance  int `json:"maintenance"`
	Warranty     int `json:"warranty"`
	StatusChange int `json:"statusChange"`
	Movement     int `json:"movement"`
	IssueReport  int `json:"issueReport"`
}

type NotificationStatusStatistics struct {
	Read   int `json:"read"`
	Unread int `json:"unread"`
}

type NotificationCreationTrend struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type NotificationSummaryStatistics struct {
	TotalNotifications         int     `json:"totalNotifications"`
	ReadPercentage             float64 `json:"readPercentage"`
	UnreadPercentage           float64 `json:"unreadPercentage"`
	MostCommonType             string  `json:"mostCommonType"`
	AverageNotificationsPerDay float64 `json:"averageNotificationsPerDay"`
	LatestCreationDate         string  `json:"latestCreationDate"`
	EarliestCreationDate       string  `json:"earliestCreationDate"`
}

// Response statistics structs (used in service/handler layer)
type NotificationStatisticsResponse struct {
	Total          NotificationCountStatisticsResponse   `json:"total"`
	ByType         NotificationTypeStatisticsResponse    `json:"byType"`
	ByStatus       NotificationStatusStatisticsResponse  `json:"byStatus"`
	CreationTrends []NotificationCreationTrendResponse   `json:"creationTrends"`
	Summary        NotificationSummaryStatisticsResponse `json:"summary"`
}

type NotificationCountStatisticsResponse struct {
	Count int `json:"count"`
}

type NotificationTypeStatisticsResponse struct {
	Maintenance  int `json:"maintenance"`
	Warranty     int `json:"warranty"`
	StatusChange int `json:"statusChange"`
	Movement     int `json:"movement"`
	IssueReport  int `json:"issueReport"`
}

type NotificationStatusStatisticsResponse struct {
	Read   int `json:"read"`
	Unread int `json:"unread"`
}

type NotificationCreationTrendResponse struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type NotificationSummaryStatisticsResponse struct {
	TotalNotifications         int     `json:"totalNotifications"`
	ReadPercentage             float64 `json:"readPercentage"`
	UnreadPercentage           float64 `json:"unreadPercentage"`
	MostCommonType             string  `json:"mostCommonType"`
	AverageNotificationsPerDay float64 `json:"averageNotificationsPerDay"`
	LatestCreationDate         string  `json:"latestCreationDate"`
	EarliestCreationDate       string  `json:"earliestCreationDate"`
}
