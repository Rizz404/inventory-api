package mapper

import (
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/oklog/ulid/v2"
)

// *==================== Model conversions ====================

func ToModelNotification(d *domain.Notification) model.Notification {
	modelNotification := model.Notification{
		Type:      d.Type,
		Priority:  d.Priority,
		IsRead:    d.IsRead,
		ReadAt:    d.ReadAt,
		ExpiresAt: d.ExpiresAt,
	}

	if d.ID != "" {
		if parsedID, err := ulid.Parse(d.ID); err == nil {
			modelNotification.ID = model.SQLULID(parsedID)
		}
	}

	if d.UserID != "" {
		if parsedUserID, err := ulid.Parse(d.UserID); err == nil {
			modelNotification.UserID = model.SQLULID(parsedUserID)
		}
	}

	// Related entity (polymorphic)
	if d.RelatedEntityType != nil {
		modelNotification.RelatedEntityType = d.RelatedEntityType
	}

	if d.RelatedEntityID != nil && *d.RelatedEntityID != "" {
		if parsedEntityID, err := ulid.Parse(*d.RelatedEntityID); err == nil {
			sqlEntityID := model.SQLULID(parsedEntityID)
			modelNotification.RelatedEntityID = &sqlEntityID
		}
	}

	// Legacy support
	if d.RelatedAssetID != nil && *d.RelatedAssetID != "" {
		if parsedAssetID, err := ulid.Parse(*d.RelatedAssetID); err == nil {
			modelULID := model.SQLULID(parsedAssetID)
			modelNotification.RelatedAssetID = &modelULID
		}
	}

	return modelNotification
}

func ToModelNotificationForCreate(d *domain.Notification) model.Notification {
	modelNotification := model.Notification{
		Type:      d.Type,
		Priority:  d.Priority,
		IsRead:    d.IsRead,
		ReadAt:    d.ReadAt,
		ExpiresAt: d.ExpiresAt,
	}

	if d.UserID != "" {
		if parsedUserID, err := ulid.Parse(d.UserID); err == nil {
			modelNotification.UserID = model.SQLULID(parsedUserID)
		}
	}

	// Related entity (polymorphic)
	if d.RelatedEntityType != nil {
		modelNotification.RelatedEntityType = d.RelatedEntityType
	}

	if d.RelatedEntityID != nil && *d.RelatedEntityID != "" {
		if parsedEntityID, err := ulid.Parse(*d.RelatedEntityID); err == nil {
			sqlEntityID := model.SQLULID(parsedEntityID)
			modelNotification.RelatedEntityID = &sqlEntityID
		}
	}

	// Legacy support
	if d.RelatedAssetID != nil && *d.RelatedAssetID != "" {
		if parsedAssetID, err := ulid.Parse(*d.RelatedAssetID); err == nil {
			modelULID := model.SQLULID(parsedAssetID)
			modelNotification.RelatedAssetID = &modelULID
		}
	}

	return modelNotification
}

func ToModelNotificationTranslation(d *domain.NotificationTranslation) model.NotificationTranslation {
	modelTranslation := model.NotificationTranslation{
		LangCode: d.LangCode,
		Title:    d.Title,
		Message:  d.Message,
	}

	if d.ID != "" {
		if parsedID, err := ulid.Parse(d.ID); err == nil {
			modelTranslation.ID = model.SQLULID(parsedID)
		}
	}

	if d.NotificationID != "" {
		if parsedNotificationID, err := ulid.Parse(d.NotificationID); err == nil {
			modelTranslation.NotificationID = model.SQLULID(parsedNotificationID)
		}
	}

	return modelTranslation
}

func ToModelNotificationTranslationForCreate(notificationID string, d *domain.NotificationTranslation) model.NotificationTranslation {
	modelTranslation := model.NotificationTranslation{
		LangCode: d.LangCode,
		Title:    d.Title,
		Message:  d.Message,
	}

	if notificationID != "" {
		if parsedNotificationID, err := ulid.Parse(notificationID); err == nil {
			modelTranslation.NotificationID = model.SQLULID(parsedNotificationID)
		}
	}

	return modelTranslation
}

// *==================== Entity conversions ====================
func ToDomainNotification(m *model.Notification) domain.Notification {
	domainNotification := domain.Notification{
		ID:        m.ID.String(),
		UserID:    m.UserID.String(),
		Type:      m.Type,
		Priority:  m.Priority,
		IsRead:    m.IsRead,
		ReadAt:    m.ReadAt,
		ExpiresAt: m.ExpiresAt,
		CreatedAt: m.CreatedAt,
	}

	// Related entity (polymorphic)
	if m.RelatedEntityType != nil {
		domainNotification.RelatedEntityType = m.RelatedEntityType
	}

	if m.RelatedEntityID != nil && !m.RelatedEntityID.IsZero() {
		entityIDStr := m.RelatedEntityID.String()
		domainNotification.RelatedEntityID = &entityIDStr
	}

	// Legacy support
	if m.RelatedAssetID != nil && !m.RelatedAssetID.IsZero() {
		assetIDStr := m.RelatedAssetID.String()
		domainNotification.RelatedAssetID = &assetIDStr
	}

	if len(m.Translations) > 0 {
		domainNotification.Translations = make([]domain.NotificationTranslation, len(m.Translations))
		for i, translation := range m.Translations {
			domainNotification.Translations[i] = ToDomainNotificationTranslation(&translation)
		}
	}

	return domainNotification
}

func ToDomainNotificationTranslation(m *model.NotificationTranslation) domain.NotificationTranslation {
	return domain.NotificationTranslation{
		ID:             m.ID.String(),
		NotificationID: m.NotificationID.String(),
		LangCode:       m.LangCode,
		Title:          m.Title,
		Message:        m.Message,
	}
}

func ToDomainNotifications(models []model.Notification) []domain.Notification {
	notifications := make([]domain.Notification, len(models))
	for i, m := range models {
		notifications[i] = ToDomainNotification(&m)
	}
	return notifications
}

// *==================== Entity Response conversions ====================
func NotificationToResponse(d *domain.Notification, langCode string) domain.NotificationResponse {
	response := domain.NotificationResponse{
		ID:                d.ID,
		UserID:            d.UserID,
		RelatedEntityType: d.RelatedEntityType,
		RelatedEntityID:   d.RelatedEntityID,
		RelatedAssetID:    d.RelatedAssetID,
		Type:              d.Type,
		Priority:          d.Priority,
		IsRead:            d.IsRead,
		ReadAt:            d.ReadAt,
		ExpiresAt:         d.ExpiresAt,
		CreatedAt:         d.CreatedAt,
		Translations:      make([]domain.NotificationTranslationResponse, len(d.Translations)),
	}

	// Populate translations
	for i, translation := range d.Translations {
		response.Translations[i] = domain.NotificationTranslationResponse{
			LangCode: translation.LangCode,
			Title:    translation.Title,
			Message:  translation.Message,
		}
	}

	// Find translation for the requested language
	for _, translation := range d.Translations {
		if translation.LangCode == langCode {
			response.Title = translation.Title
			response.Message = translation.Message
			break
		}
	}

	// If no translation found for requested language, use first available
	if response.Title == "" && len(d.Translations) > 0 {
		response.Title = d.Translations[0].Title
		response.Message = d.Translations[0].Message
	}

	return response
}

func NotificationsToResponses(notifications []domain.Notification, langCode string) []domain.NotificationResponse {
	responses := make([]domain.NotificationResponse, len(notifications))
	for i, notification := range notifications {
		responses[i] = NotificationToResponse(&notification, langCode)
	}
	return responses
}

func NotificationToListResponse(d *domain.Notification, langCode string) domain.NotificationListResponse {
	response := domain.NotificationListResponse{
		ID:                d.ID,
		UserID:            d.UserID,
		RelatedEntityType: d.RelatedEntityType,
		RelatedEntityID:   d.RelatedEntityID,
		RelatedAssetID:    d.RelatedAssetID,
		Type:              d.Type,
		Priority:          d.Priority,
		IsRead:            d.IsRead,
		ReadAt:            d.ReadAt,
		ExpiresAt:         d.ExpiresAt,
		CreatedAt:         d.CreatedAt,
	}

	// Find translation for the requested language
	for _, translation := range d.Translations {
		if translation.LangCode == langCode {
			response.Title = translation.Title
			response.Message = translation.Message
			break
		}
	}

	// If no translation found for requested language, use first available
	if response.Title == "" && len(d.Translations) > 0 {
		response.Title = d.Translations[0].Title
		response.Message = d.Translations[0].Message
	}

	return response
}

func NotificationsToListResponses(notifications []domain.Notification, langCode string) []domain.NotificationListResponse {
	responses := make([]domain.NotificationListResponse, len(notifications))
	for i, notification := range notifications {
		responses[i] = NotificationToListResponse(&notification, langCode)
	}
	return responses
}

// ToModelNotificationUpdateMap converts UpdateNotificationPayload to map for database updates
func ToModelNotificationUpdateMap(payload *domain.UpdateNotificationPayload) map[string]interface{} {
	updates := make(map[string]interface{})

	if payload.IsRead != nil {
		updates["is_read"] = *payload.IsRead
		if *payload.IsRead {
			now := time.Now()
			updates["read_at"] = now
		} else {
			updates["read_at"] = nil
		}
	}

	if payload.Priority != nil {
		updates["priority"] = *payload.Priority
	}

	if payload.ExpiresAt != nil {
		updates["expires_at"] = *payload.ExpiresAt
	}

	return updates
}

// *==================== Statistics conversions ====================
// NotificationStatisticsToResponse converts NotificationStatistics to NotificationStatisticsResponse
func NotificationStatisticsToResponse(stats *domain.NotificationStatistics) domain.NotificationStatisticsResponse {
	response := domain.NotificationStatisticsResponse{
		Total: domain.NotificationCountStatisticsResponse{
			Count: stats.Total.Count,
		},
		ByType: domain.NotificationTypeStatisticsResponse{
			Maintenance:    stats.ByType.Maintenance,
			Warranty:       stats.ByType.Warranty,
			Issue:          stats.ByType.Issue,
			Movement:       stats.ByType.Movement,
			StatusChange:   stats.ByType.StatusChange,
			LocationChange: stats.ByType.LocationChange,
			CategoryChange: stats.ByType.CategoryChange,
		},
		ByStatus: domain.NotificationStatusStatisticsResponse{
			Read:   stats.ByStatus.Read,
			Unread: stats.ByStatus.Unread,
		},
		Summary: domain.NotificationSummaryStatisticsResponse{
			TotalNotifications:         stats.Summary.TotalNotifications,
			ReadPercentage:             domain.NewDecimal2(stats.Summary.ReadPercentage),
			UnreadPercentage:           domain.NewDecimal2(stats.Summary.UnreadPercentage),
			MostCommonType:             stats.Summary.MostCommonType,
			AverageNotificationsPerDay: domain.NewDecimal2(stats.Summary.AverageNotificationsPerDay),
			LatestCreationDate:         stats.Summary.LatestCreationDate,
			EarliestCreationDate:       stats.Summary.EarliestCreationDate,
		},
	}

	// Convert creation trends
	response.CreationTrends = make([]domain.NotificationCreationTrendResponse, len(stats.CreationTrends))
	for i, trend := range stats.CreationTrends {
		response.CreationTrends[i] = domain.NotificationCreationTrendResponse{
			Date:  trend.Date,
			Count: trend.Count,
		}
	}

	return response
}

// MapNotificationSortFieldToColumn maps the NotificationSortField to the corresponding database column
func MapNotificationSortFieldToColumn(field domain.NotificationSortField) string {
	columnMap := map[domain.NotificationSortField]string{
		domain.NotificationSortByType:      "n.type",
		domain.NotificationSortByPriority:  "n.priority",
		domain.NotificationSortByTitle:     "nt.title",
		domain.NotificationSortByMessage:   "nt.message",
		domain.NotificationSortByIsRead:    "n.is_read",
		domain.NotificationSortByCreatedAt: "n.created_at",
	}

	if column, exists := columnMap[field]; exists {
		return column
	}
	return "n.created_at"
}
