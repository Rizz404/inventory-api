package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
)

func findLocationTranslation(translations []model.LocationTranslation, langCode string) string {
	for _, t := range translations {
		if t.LangCode == langCode {
			return t.LocationName
		}
	}
	for _, t := range translations {
		if t.LangCode == DefaultLangCode {
			return t.LocationName
		}
	}
	if len(translations) > 0 {
		return translations[0].LocationName
	}
	return ""
}

func ToDomainLocationResponse(m model.Location, langCode string) domain.LocationResponse {
	return domain.LocationResponse{
		ID:           m.ID.String(),
		LocationCode: m.LocationCode,
		Building:     m.Building,
		Floor:        m.Floor,
		Name:         findLocationTranslation(m.Translations, langCode),
	}
}

func ToDomainLocationsResponse(m []model.Location, langCode string) []domain.LocationResponse {
	responses := make([]domain.LocationResponse, len(m))
	for i, location := range m {
		responses[i] = ToDomainLocationResponse(location, langCode)
	}
	return responses
}
