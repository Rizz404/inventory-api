package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
)

func findNotificationTranslation(translations []model.NotificationTranslation, langCode string) (title, message string) {
	for _, t := range translations {
		if t.LangCode == langCode {
			return t.Title, t.Message
		}
	}
	for _, t := range translations {
		if t.LangCode == DefaultLangCode {
			return t.Title, t.Message
		}
	}
	if len(translations) > 0 {
		return translations[0].Title, translations[0].Message
	}
	return "", ""
}

func ToDomainNotificationResponse(m model.Notification, langCode string) domain.NotificationResponse {
	title, msg := findNotificationTranslation(m.Translations, langCode)
	resp := domain.NotificationResponse{
		ID:        m.ID.String(),
		Type:      m.Type,
		IsRead:    m.IsRead,
		CreatedAt: m.CreatedAt.Format(TimeFormat),
		Title:     title,
		Message:   msg,
	}
	if m.RelatedAssetID != nil {
		resp.RelatedAssetID = Ptr(m.RelatedAssetID.String())
	}
	return resp
}

func ToDomainNotificationsResponse(m []model.Notification, langCode string) []domain.NotificationResponse {
	responses := make([]domain.NotificationResponse, len(m))
	for i, notif := range m {
		responses[i] = ToDomainNotificationResponse(notif, langCode)
	}
	return responses
}
