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
		Latitude:     d.Latitude,
		Longitude:    d.Longitude,
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
		Latitude:     d.Latitude,
		Longitude:    d.Longitude,
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
		Latitude:     m.Latitude,
		Longitude:    m.Longitude,
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

func ToDomainLocations(models []model.Location) []domain.Location {
	locations := make([]domain.Location, len(models))
	for i, m := range models {
		locations[i] = ToDomainLocation(&m)
	}
	return locations
}

func ToDomainLocationTranslation(m *model.LocationTranslation) domain.LocationTranslation {
	return domain.LocationTranslation{
		ID:           m.ID.String(),
		LocationID:   m.LocationID.String(),
		LangCode:     m.LangCode,
		LocationName: m.LocationName,
	}
}

// Domain -> Response conversions (for service layer)
func LocationToResponse(d *domain.Location, langCode string) domain.LocationResponse {
	response := domain.LocationResponse{
		ID:           d.ID,
		LocationCode: d.LocationCode,
		Building:     d.Building,
		Floor:        d.Floor,
		Latitude:     d.Latitude,
		Longitude:    d.Longitude,
		CreatedAt:    d.CreatedAt.Format(TimeFormat),
		UpdatedAt:    d.UpdatedAt.Format(TimeFormat),
		Translations: make([]domain.LocationTranslationResponse, len(d.Translations)),
	}

	// Populate translations
	for i, translation := range d.Translations {
		response.Translations[i] = domain.LocationTranslationResponse{
			LangCode:     translation.LangCode,
			LocationName: translation.LocationName,
		}
	}

	// Find translation for the requested language
	for _, translation := range d.Translations {
		if translation.LangCode == langCode {
			response.LocationName = translation.LocationName
			break
		}
	}

	// If no translation found for requested language, use first available
	if response.LocationName == "" && len(d.Translations) > 0 {
		response.LocationName = d.Translations[0].LocationName
	}

	return response
}

func LocationsToResponses(locations []domain.Location, langCode string) []domain.LocationResponse {
	responses := make([]domain.LocationResponse, len(locations))
	for i, location := range locations {
		responses[i] = LocationToResponse(&location, langCode)
	}
	return responses
}

// Statistics conversions
func LocationStatisticsToResponse(stats *domain.LocationStatistics) domain.LocationStatisticsResponse {
	buildingStats := make([]domain.BuildingStatisticsResponse, len(stats.ByBuilding))
	for i, building := range stats.ByBuilding {
		buildingStats[i] = domain.BuildingStatisticsResponse{
			Building: building.Building,
			Count:    building.Count,
		}
	}

	floorStats := make([]domain.FloorStatisticsResponse, len(stats.ByFloor))
	for i, floor := range stats.ByFloor {
		floorStats[i] = domain.FloorStatisticsResponse{
			Floor: floor.Floor,
			Count: floor.Count,
		}
	}

	trends := make([]domain.LocationCreationTrendResponse, len(stats.CreationTrends))
	for i, trend := range stats.CreationTrends {
		trends[i] = domain.LocationCreationTrendResponse{
			Date:  trend.Date,
			Count: trend.Count,
		}
	}

	return domain.LocationStatisticsResponse{
		Total: domain.LocationCountStatisticsResponse{
			Count: stats.Total.Count,
		},
		ByBuilding: buildingStats,
		ByFloor:    floorStats,
		Geographic: domain.GeographicStatisticsResponse{
			WithCoordinates:    stats.Geographic.WithCoordinates,
			WithoutCoordinates: stats.Geographic.WithoutCoordinates,
			AverageLatitude:    stats.Geographic.AverageLatitude,
			AverageLongitude:   stats.Geographic.AverageLongitude,
		},
		CreationTrends: trends,
		Summary: domain.LocationSummaryStatisticsResponse{
			TotalLocations:           stats.Summary.TotalLocations,
			LocationsWithBuilding:    stats.Summary.LocationsWithBuilding,
			LocationsWithoutBuilding: stats.Summary.LocationsWithoutBuilding,
			LocationsWithFloor:       stats.Summary.LocationsWithFloor,
			LocationsWithoutFloor:    stats.Summary.LocationsWithoutFloor,
			LocationsWithCoordinates: stats.Summary.LocationsWithCoordinates,
			CoordinatesPercentage:    stats.Summary.CoordinatesPercentage,
			BuildingPercentage:       stats.Summary.BuildingPercentage,
			FloorPercentage:          stats.Summary.FloorPercentage,
			TotalBuildings:           stats.Summary.TotalBuildings,
			TotalFloors:              stats.Summary.TotalFloors,
			AverageLocationsPerDay:   stats.Summary.AverageLocationsPerDay,
			LatestCreationDate:       stats.Summary.LatestCreationDate,
			EarliestCreationDate:     stats.Summary.EarliestCreationDate,
		},
	}
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
	if payload.Building != nil {
		updates["building"] = payload.Building
	}
	if payload.Floor != nil {
		updates["floor"] = payload.Floor
	}
	if payload.Latitude != nil {
		updates["latitude"] = payload.Latitude
	}
	if payload.Longitude != nil {
		updates["longitude"] = payload.Longitude
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

func LocationToListResponse(d *domain.Location, langCode string) domain.LocationListResponse {
	response := domain.LocationListResponse{
		ID:           d.ID,
		LocationCode: d.LocationCode,
		Building:     d.Building,
		Floor:        d.Floor,
		Latitude:     d.Latitude,
		Longitude:    d.Longitude,
		CreatedAt:    d.CreatedAt.Format(TimeFormat),
		UpdatedAt:    d.UpdatedAt.Format(TimeFormat),
	}

	// Find translation for the requested language
	for _, translation := range d.Translations {
		if translation.LangCode == langCode {
			response.LocationName = translation.LocationName
			break
		}
	}

	// If no translation found for requested language, use first available
	if response.LocationName == "" && len(d.Translations) > 0 {
		response.LocationName = d.Translations[0].LocationName
	}

	return response
}

func LocationsToListResponses(locations []domain.Location, langCode string) []domain.LocationListResponse {
	responses := make([]domain.LocationListResponse, len(locations))
	for i, location := range locations {
		responses[i] = LocationToListResponse(&location, langCode)
	}
	return responses
}
