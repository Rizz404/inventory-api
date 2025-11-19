package domain

import "time"

// --- Enums ---

type NotificationType string

const (
	// Core notification types
	NotificationTypeMaintenance    NotificationType = "MAINTENANCE"     // All maintenance-related (due, overdue, completed)
	NotificationTypeWarranty       NotificationType = "WARRANTY"        // Warranty expiring/expired
	NotificationTypeIssue          NotificationType = "ISSUE"           // Issue reports (reported, updated, resolved)
	NotificationTypeMovement       NotificationType = "MOVEMENT"        // Asset movements
	NotificationTypeStatusChange   NotificationType = "STATUS_CHANGE"   // Asset status changes
	NotificationTypeLocationChange NotificationType = "LOCATION_CHANGE" // Location changes
	NotificationTypeCategoryChange NotificationType = "CATEGORY_CHANGE" // Category changes
)

type NotificationPriority string

const (
	NotificationPriorityLow    NotificationPriority = "LOW"
	NotificationPriorityNormal NotificationPriority = "NORMAL"
	NotificationPriorityHigh   NotificationPriority = "HIGH"
	NotificationPriorityUrgent NotificationPriority = "URGENT"
)

// ToFCMPriority maps notification priority to FCM priority
func (p NotificationPriority) ToFCMPriority() string {
	switch p {
	case NotificationPriorityUrgent, NotificationPriorityHigh:
		return "high"
	default:
		return "normal"
	}
}

type NotificationSortField string

const (
	NotificationSortByType      NotificationSortField = "type"
	NotificationSortByPriority  NotificationSortField = "priority"
	NotificationSortByTitle     NotificationSortField = "title"
	NotificationSortByMessage   NotificationSortField = "message"
	NotificationSortByIsRead    NotificationSortField = "isRead"
	NotificationSortByCreatedAt NotificationSortField = "createdAt"
)

// --- Structs ---

type Notification struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`

	// Related entity (polymorphic)
	RelatedEntityType *string `json:"relatedEntityType,omitempty"`
	RelatedEntityID   *string `json:"relatedEntityId,omitempty"`

	// Legacy support (deprecated, use RelatedEntityID instead)
	RelatedAssetID *string `json:"relatedAssetId,omitempty"`

	Type     NotificationType     `json:"type"`
	Priority NotificationPriority `json:"priority"`

	// Status
	IsRead bool       `json:"isRead"`
	ReadAt *time.Time `json:"readAt,omitempty"`

	// Expiration
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`

	CreatedAt    time.Time                 `json:"createdAt"`
	Translations []NotificationTranslation `json:"translations,omitempty"`
}

type NotificationTranslation struct {
	ID             string `json:"id"`
	NotificationID string `json:"notificationId"`
	LangCode       string `json:"langCode"`
	Title          string `json:"title"`
	Message        string `json:"message"`
}

type NotificationTranslationResponse struct {
	LangCode string `json:"langCode"`
	Title    string `json:"title"`
	Message  string `json:"message"`
}

type NotificationResponse struct {
	ID                string                            `json:"id"`
	UserID            string                            `json:"userId"`
	RelatedEntityType *string                           `json:"relatedEntityType,omitempty"`
	RelatedEntityID   *string                           `json:"relatedEntityId,omitempty"`
	RelatedAssetID    *string                           `json:"relatedAssetId,omitempty"` // deprecated
	Type              NotificationType                  `json:"type"`
	Priority          NotificationPriority              `json:"priority"`
	IsRead            bool                              `json:"isRead"`
	ReadAt            *time.Time                        `json:"readAt,omitempty"`
	ExpiresAt         *time.Time                        `json:"expiresAt,omitempty"`
	CreatedAt         time.Time                         `json:"createdAt"`
	Title             string                            `json:"title"`
	Message           string                            `json:"message"`
	Translations      []NotificationTranslationResponse `json:"translations"`
}

type NotificationListResponse struct {
	ID                string               `json:"id"`
	UserID            string               `json:"userId"`
	RelatedEntityType *string              `json:"relatedEntityType,omitempty"`
	RelatedEntityID   *string              `json:"relatedEntityId,omitempty"`
	RelatedAssetID    *string              `json:"relatedAssetId,omitempty"` // deprecated
	Type              NotificationType     `json:"type"`
	Priority          NotificationPriority `json:"priority"`
	IsRead            bool                 `json:"isRead"`
	ReadAt            *time.Time           `json:"readAt,omitempty"`
	ExpiresAt         *time.Time           `json:"expiresAt,omitempty"`
	CreatedAt         time.Time            `json:"createdAt"`
	Title             string               `json:"title"`
	Message           string               `json:"message"`
}

type BulkDeleteNotifications struct {
	RequestedIDS []string `json:"requestedIds"`
	DeletedIDS   []string `json:"deletedIds"`
}

type BulkDeleteNotificationsResponse struct {
	RequestedIDS []string `json:"requestedIds"`
	DeletedIDS   []string `json:"deletedIds"`
}

// --- Payloads ---

// Notifications are typically created by the system, not directly by users.
// Payloads might not be needed for direct API exposure but can be used internally.
type CreateNotificationPayload struct {
	UserID            string                                 `json:"userId"`
	RelatedEntityType *string                                `json:"relatedEntityType,omitempty"`
	RelatedEntityID   *string                                `json:"relatedEntityId,omitempty"`
	RelatedAssetID    *string                                `json:"relatedAssetId,omitempty"` // deprecated
	Type              NotificationType                       `json:"type"`
	Priority          NotificationPriority                   `json:"priority"`
	ExpiresAt         *time.Time                             `json:"expiresAt,omitempty"`
	Translations      []CreateNotificationTranslationPayload `json:"translations"`
}

type CreateNotificationTranslationPayload struct {
	LangCode string `json:"langCode"`
	Title    string `json:"title"`
	Message  string `json:"message"`
}

type UpdateNotificationPayload struct {
	IsRead       *bool                                  `json:"isRead,omitempty"`
	Priority     *NotificationPriority                  `json:"priority,omitempty"`
	ExpiresAt    *time.Time                             `json:"expiresAt,omitempty"`
	Translations []CreateNotificationTranslationPayload `json:"translations,omitempty"`
}

type BulkDeleteNotificationsPayload struct {
	IDS []string `json:"ids" validate:"required,min=1,max=100,dive,required"`
}

type MarkNotificationsPayload struct {
	NotificationIDs []string `json:"notificationIds" validate:"required,min=1,dive"`
}

// --- Query Parameters ---

type NotificationFilterOptions struct {
	UserID            *string               `json:"userId,omitempty"`
	RelatedEntityType *string               `json:"relatedEntityType,omitempty"`
	RelatedEntityID   *string               `json:"relatedEntityId,omitempty"`
	RelatedAssetID    *string               `json:"relatedAssetId,omitempty"` // deprecated
	Type              *NotificationType     `json:"type,omitempty"`
	Priority          *NotificationPriority `json:"priority,omitempty"`
	IsRead            *bool                 `json:"isRead,omitempty"`
}

type NotificationSortOptions struct {
	Field NotificationSortField `json:"field" example:"createdAt"`
	Order SortOrder             `json:"order" example:"desc"`
}

type NotificationParams struct {
	SearchQuery *string                    `json:"searchQuery,omitempty"`
	Filters     *NotificationFilterOptions `json:"filters,omitempty"`
	Sort        *NotificationSortOptions   `json:"sort,omitempty"`
	Pagination  *PaginationOptions         `json:"pagination,omitempty"`
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
	Maintenance    int `json:"maintenance"`
	Warranty       int `json:"warranty"`
	Issue          int `json:"issue"`
	Movement       int `json:"movement"`
	StatusChange   int `json:"statusChange"`
	LocationChange int `json:"locationChange"`
	CategoryChange int `json:"categoryChange"`
}

type NotificationStatusStatistics struct {
	Read   int `json:"read"`
	Unread int `json:"unread"`
}

type NotificationCreationTrend struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

type NotificationSummaryStatistics struct {
	TotalNotifications         int       `json:"totalNotifications"`
	ReadPercentage             float64   `json:"readPercentage"`
	UnreadPercentage           float64   `json:"unreadPercentage"`
	MostCommonType             string    `json:"mostCommonType"`
	AverageNotificationsPerDay float64   `json:"averageNotificationsPerDay"`
	LatestCreationDate         time.Time `json:"latestCreationDate"`
	EarliestCreationDate       time.Time `json:"earliestCreationDate"`
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
	Maintenance    int `json:"maintenance"`
	Warranty       int `json:"warranty"`
	Issue          int `json:"issue"`
	Movement       int `json:"movement"`
	StatusChange   int `json:"statusChange"`
	LocationChange int `json:"locationChange"`
	CategoryChange int `json:"categoryChange"`
}

type NotificationStatusStatisticsResponse struct {
	Read   int `json:"read"`
	Unread int `json:"unread"`
}

type NotificationCreationTrendResponse struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

type NotificationSummaryStatisticsResponse struct {
	TotalNotifications         int       `json:"totalNotifications"`
	ReadPercentage             Decimal2  `json:"readPercentage"`
	UnreadPercentage           Decimal2  `json:"unreadPercentage"`
	MostCommonType             string    `json:"mostCommonType"`
	AverageNotificationsPerDay Decimal2  `json:"averageNotificationsPerDay"`
	LatestCreationDate         time.Time `json:"latestCreationDate"`
	EarliestCreationDate       time.Time `json:"earliestCreationDate"`
}
