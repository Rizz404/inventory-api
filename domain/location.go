package domain

import "time"

// --- Structs ---

type Location struct {
	ID           string                `json:"id"`
	LocationCode string                `json:"locationCode"`
	Building     *string               `json:"building"`
	Floor        *string               `json:"floor"`
	Latitude     *float64              `json:"latitude"`
	Longitude    *float64              `json:"longitude"`
	CreatedAt    time.Time             `json:"createdAt"`
	UpdatedAt    time.Time             `json:"updatedAt"`
	Translations []LocationTranslation `json:"translations,omitempty"`
}

type LocationTranslation struct {
	ID           string `json:"id"`
	LocationID   string `json:"locationId"`
	LangCode     string `json:"langCode"`
	LocationName string `json:"locationName"`
}

type LocationTranslationResponse struct {
	LangCode     string `json:"langCode"`
	LocationName string `json:"locationName"`
}

type LocationResponse struct {
	ID           string                        `json:"id"`
	LocationName string                        `json:"locationName"`
	LocationCode string                        `json:"locationCode"`
	Building     *string                       `json:"building"`
	Floor        *string                       `json:"floor"`
	Latitude     *float64                      `json:"latitude"`
	Longitude    *float64                      `json:"longitude"`
	CreatedAt    time.Time                     `json:"createdAt"`
	UpdatedAt    time.Time                     `json:"updatedAt"`
	Translations []LocationTranslationResponse `json:"translations"`
}

type LocationListResponse struct {
	ID           string    `json:"id"`
	LocationName string    `json:"locationName"`
	LocationCode string    `json:"locationCode"`
	Building     *string   `json:"building"`
	Floor        *string   `json:"floor"`
	Latitude     *float64  `json:"latitude"`
	Longitude    *float64  `json:"longitude"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// --- Payloads ---

type CreateLocationPayload struct {
	LocationCode string                             `json:"locationCode" validate:"required,max=20"`
	Building     *string                            `json:"building,omitempty" validate:"omitempty,max=100"`
	Floor        *string                            `json:"floor,omitempty" validate:"omitempty,max=20"`
	Latitude     *float64                           `json:"latitude,omitempty" validate:"omitempty,latitude"`
	Longitude    *float64                           `json:"longitude,omitempty" validate:"omitempty,longitude"`
	Translations []CreateLocationTranslationPayload `json:"translations" validate:"required,min=1,dive"`
}

type CreateLocationTranslationPayload struct {
	LangCode     string `json:"langCode" validate:"required,max=5"`
	LocationName string `json:"locationName" validate:"required,max=100"`
}

type UpdateLocationPayload struct {
	LocationCode *string                            `json:"locationCode,omitempty" validate:"omitempty,max=20"`
	Building     *string                            `json:"building,omitempty" validate:"omitempty,max=100"`
	Floor        *string                            `json:"floor,omitempty" validate:"omitempty,max=20"`
	Latitude     *float64                           `json:"latitude,omitempty" validate:"omitempty,latitude"`
	Longitude    *float64                           `json:"longitude,omitempty" validate:"omitempty,longitude"`
	Translations []UpdateLocationTranslationPayload `json:"translations,omitempty" validate:"omitempty,dive"`
}

type UpdateLocationTranslationPayload struct {
	LangCode     string  `json:"langCode" validate:"required,max=5"`
	LocationName *string `json:"locationName,omitempty" validate:"omitempty,max=100"`
}

// --- Statistics ---

// Internal statistics structs (used in repository layer)
type LocationStatistics struct {
	Total          LocationCountStatistics   `json:"total"`
	ByBuilding     []BuildingStatistics      `json:"byBuilding"`
	ByFloor        []FloorStatistics         `json:"byFloor"`
	Geographic     GeographicStatistics      `json:"geographic"`
	CreationTrends []LocationCreationTrend   `json:"creationTrends"`
	Summary        LocationSummaryStatistics `json:"summary"`
}

type LocationCountStatistics struct {
	Count int `json:"count"`
}

type BuildingStatistics struct {
	Building string `json:"building"`
	Count    int    `json:"count"`
}

type FloorStatistics struct {
	Floor string `json:"floor"`
	Count int    `json:"count"`
}

type GeographicStatistics struct {
	WithCoordinates    int      `json:"withCoordinates"`
	WithoutCoordinates int      `json:"withoutCoordinates"`
	AverageLatitude    *float64 `json:"averageLatitude"`
	AverageLongitude   *float64 `json:"averageLongitude"`
}

type LocationCreationTrend struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

type LocationSummaryStatistics struct {
	TotalLocations           int       `json:"totalLocations"`
	LocationsWithBuilding    int       `json:"locationsWithBuilding"`
	LocationsWithoutBuilding int       `json:"locationsWithoutBuilding"`
	LocationsWithFloor       int       `json:"locationsWithFloor"`
	LocationsWithoutFloor    int       `json:"locationsWithoutFloor"`
	LocationsWithCoordinates int       `json:"locationsWithCoordinates"`
	CoordinatesPercentage    float64   `json:"coordinatesPercentage"`
	BuildingPercentage       float64   `json:"buildingPercentage"`
	FloorPercentage          float64   `json:"floorPercentage"`
	TotalBuildings           int       `json:"totalBuildings"`
	TotalFloors              int       `json:"totalFloors"`
	AverageLocationsPerDay   float64   `json:"averageLocationsPerDay"`
	LatestCreationDate       time.Time `json:"latestCreationDate"`
	EarliestCreationDate     time.Time `json:"earliestCreationDate"`
}

// Response statistics structs (used in service/handler layer)
type LocationStatisticsResponse struct {
	Total          LocationCountStatisticsResponse   `json:"total"`
	ByBuilding     []BuildingStatisticsResponse      `json:"byBuilding"`
	ByFloor        []FloorStatisticsResponse         `json:"byFloor"`
	Geographic     GeographicStatisticsResponse      `json:"geographic"`
	CreationTrends []LocationCreationTrendResponse   `json:"creationTrends"`
	Summary        LocationSummaryStatisticsResponse `json:"summary"`
}

type LocationCountStatisticsResponse struct {
	Count int `json:"count"`
}

type BuildingStatisticsResponse struct {
	Building string `json:"building"`
	Count    int    `json:"count"`
}

type FloorStatisticsResponse struct {
	Floor string `json:"floor"`
	Count int    `json:"count"`
}

type GeographicStatisticsResponse struct {
	WithCoordinates    int      `json:"withCoordinates"`
	WithoutCoordinates int      `json:"withoutCoordinates"`
	AverageLatitude    *float64 `json:"averageLatitude"`
	AverageLongitude   *float64 `json:"averageLongitude"`
}

type LocationCreationTrendResponse struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

type LocationSummaryStatisticsResponse struct {
	TotalLocations           int       `json:"totalLocations"`
	LocationsWithBuilding    int       `json:"locationsWithBuilding"`
	LocationsWithoutBuilding int       `json:"locationsWithoutBuilding"`
	LocationsWithFloor       int       `json:"locationsWithFloor"`
	LocationsWithoutFloor    int       `json:"locationsWithoutFloor"`
	LocationsWithCoordinates int       `json:"locationsWithCoordinates"`
	CoordinatesPercentage    float64   `json:"coordinatesPercentage"`
	BuildingPercentage       float64   `json:"buildingPercentage"`
	FloorPercentage          float64   `json:"floorPercentage"`
	TotalBuildings           int       `json:"totalBuildings"`
	TotalFloors              int       `json:"totalFloors"`
	AverageLocationsPerDay   float64   `json:"averageLocationsPerDay"`
	LatestCreationDate       time.Time `json:"latestCreationDate"`
	EarliestCreationDate     time.Time `json:"earliestCreationDate"`
}
