package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/oklog/ulid/v2"
)

// ===== Notification Mappers =====

func ToModelNotification(d *domain.Notification) model.Notification {
	modelNotification := model.Notification{
		Type:   d.Type,
		IsRead: d.IsRead,
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
		Type:   d.Type,
		IsRead: d.IsRead,
	}

	if d.UserID != "" {
		if parsedUserID, err := ulid.Parse(d.UserID); err == nil {
			modelNotification.UserID = model.SQLULID(parsedUserID)
		}
	}

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

func ToDomainNotification(m *model.Notification) domain.Notification {
	domainNotification := domain.Notification{
		ID:        m.ID.String(),
		UserID:    m.UserID.String(),
		Type:      m.Type,
		IsRead:    m.IsRead,
		CreatedAt: m.CreatedAt,
	}

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

// Domain -> Response conversions (for service layer)
func NotificationToResponse(d *domain.Notification, langCode string) domain.NotificationResponse {
	response := domain.NotificationResponse{
		ID:             d.ID,
		UserID:         d.UserID,
		RelatedAssetID: d.RelatedAssetID,
		Type:           d.Type,
		IsRead:         d.IsRead,
		CreatedAt:      d.CreatedAt.Format(TimeFormat),
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
