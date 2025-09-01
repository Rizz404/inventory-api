package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/oklog/ulid/v2"
)

func ToModelLocation(d *domain.Location) model.Location {
	modelLocation := model.Location{
		LocationCode: d.LocationCode,
		Building:     d.Building,
		Floor:        d.Floor,
	}

	if d.ID != "" {
		if parsedID, err := ulid.Parse(d.ID); err == nil {
			modelLocation.ID = model.SQLULID(parsedID)
		}
	}

	return modelLocation
}

func ToModelLocationForCreate(d *domain.Location) model.Location {
	modelLocation := model.Location{
		LocationCode: d.LocationCode,
		Building:     d.Building,
		Floor:        d.Floor,
	}

	return modelLocation
}

func ToModelLocationTranslation(d *domain.LocationTranslation) model.LocationTranslation {
	modelTranslation := model.LocationTranslation{
		LangCode:     d.LangCode,
		LocationName: d.LocationName,
	}

	if d.ID != "" {
		if parsedID, err := ulid.Parse(d.ID); err == nil {
			modelTranslation.ID = model.SQLULID(parsedID)
		}
	}

	if d.LocationID != "" {
		if parsedLocationID, err := ulid.Parse(d.LocationID); err == nil {
			modelTranslation.LocationID = model.SQLULID(parsedLocationID)
		}
	}

	return modelTranslation
}

func ToModelLocationTranslationForCreate(locationID string, d *domain.LocationTranslation) model.LocationTranslation {
	modelTranslation := model.LocationTranslation{
		LangCode:     d.LangCode,
		LocationName: d.LocationName,
	}

	if locationID != "" {
		if parsedLocationID, err := ulid.Parse(locationID); err == nil {
			modelTranslation.LocationID = model.SQLULID(parsedLocationID)
		}
	}

	return modelTranslation
}

func ToDomainLocation(m *model.Location) domain.Location {
	domainLocation := domain.Location{
		ID:           m.ID.String(),
		LocationCode: m.LocationCode,
		Building:     m.Building,
		Floor:        m.Floor,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}

	if len(m.Translations) > 0 {
		domainLocation.Translations = make([]domain.LocationTranslation, len(m.Translations))
		for i, translation := range m.Translations {
			domainLocation.Translations[i] = ToDomainLocationTranslation(&translation)
		}
	}

	return domainLocation
}

func ToDomainLocationTranslation(m *model.LocationTranslation) domain.LocationTranslation {
	return domain.LocationTranslation{
		ID:           m.ID.String(),
		LocationID:   m.LocationID.String(),
		LangCode:     m.LangCode,
		LocationName: m.LocationName,
	}
}

func ToDomainLocationResponse(m *model.Location, langCode string) domain.LocationResponse {
	response := domain.LocationResponse{
		ID:           m.ID.String(),
		LocationCode: m.LocationCode,
		Building:     m.Building,
		Floor:        m.Floor,
		CreatedAt:    m.CreatedAt.Format(TimeFormat),
		UpdatedAt:    m.UpdatedAt.Format(TimeFormat),
	}

	for _, translation := range m.Translations {
		if translation.LangCode == langCode {
			response.Name = translation.LocationName
			break
		}
	}

	if response.Name == "" && len(m.Translations) > 0 {
		response.Name = m.Translations[0].LocationName
	}

	return response
}

func ToDomainLocationsResponse(m []model.Location, langCode string) []domain.LocationResponse {
	responses := make([]domain.LocationResponse, len(m))
	for i, location := range m {
		responses[i] = ToDomainLocationResponse(&location, langCode)
	}
	return responses
}

// * Convert domain.Location directly to domain.LocationResponse without going through model.Location
func DomainLocationToLocationResponse(d *domain.Location, langCode string) domain.LocationResponse {
	response := domain.LocationResponse{
		ID:           d.ID,
		LocationCode: d.LocationCode,
		CreatedAt:    d.CreatedAt.Format(TimeFormat),
		UpdatedAt:    d.UpdatedAt.Format(TimeFormat),
	}

	// * Find translation for the requested language
	for _, translation := range d.Translations {
		if translation.LangCode == langCode {
			response.Name = translation.LocationName
			break
		}
	}

	// * If no translation found for requested language, use first available
	if response.Name == "" && len(d.Translations) > 0 {
		response.Name = d.Translations[0].LocationName
	}

	return response
}

// * Convert slice of domain.Location to slice of domain.LocationResponse
func DomainLocationsToLocationsResponse(locations []domain.Location, langCode string) []domain.LocationResponse {
	responses := make([]domain.LocationResponse, len(locations))
	for i, location := range locations {
		responses[i] = DomainLocationToLocationResponse(&location, langCode)
	}
	return responses
}

func BuildLocationHierarchy(locations []domain.LocationResponse) []domain.LocationResponse {
	locationMap := make(map[string]*domain.LocationResponse)
	var rootLocations []domain.LocationResponse

	// Gunakan slice pointer agar bisa memodifikasi item di map secara langsung
	for i := range locations {
		locationMap[locations[i].ID] = &locations[i]
	}

	return rootLocations
}

func ToModelLocationUpdateMap(payload *domain.UpdateLocationPayload) map[string]any {
	updates := make(map[string]any)

	if payload.LocationCode != nil {
		updates["location_code"] = *payload.LocationCode
	}

	return updates
}

func ToModelLocationTranslationUpdateMap(payload *domain.UpdateLocationTranslationPayload) map[string]any {
	updates := make(map[string]any)

	if payload.LocationName != nil {
		updates["location_name"] = *payload.LocationName
	}

	return updates
}
