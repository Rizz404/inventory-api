package domain

import "time"

// --- Enums ---

type LocationSortField string

const (
	LocationSortByLocationCode LocationSortField = "locationCode"
	LocationSortByLocationName LocationSortField = "locationName"
	LocationSortByBuilding     LocationSortField = "building"
	LocationSortByFloor        LocationSortField = "floor"
	LocationSortByCreatedAt    LocationSortField = "createdAt"
	LocationSortByUpdatedAt    LocationSortField = "updatedAt"
)

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

type BulkDeleteLocations struct {
	RequestedIDS []string `json:"requestedIds"`
	DeletedIDS   []string `json:"deletedIds"`
}

type BulkDeleteLocationsResponse struct {
	RequestedIDS []string `json:"requestedIds"`
	DeletedIDS   []string `json:"deletedIds"`
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

type BulkDeleteLocationsPayload struct {
	IDS []string `json:"ids" validate:"required,min=1,max=100,dive,required"`
}

// --- Query Parameters ---

type LocationFilterOptions struct {
}

type LocationSortOptions struct {
	Field LocationSortField `json:"field" example:"createdAt"`
	Order SortOrder         `json:"order" example:"desc"`
}

type LocationParams struct {
	SearchQuery *string                `json:"searchQuery,omitempty"`
	Filters     *LocationFilterOptions `json:"filters,omitempty"`
	Sort        *LocationSortOptions   `json:"sort,omitempty"`
	Pagination  *PaginationOptions     `json:"pagination,omitempty"`
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
	AverageLatitude    *float64 `json:"averageLatitude"`  // Always 2 decimal places or null
	AverageLongitude   *float64 `json:"averageLongitude"` // Always 2 decimal places or null
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
	CoordinatesPercentage    Decimal2  `json:"coordinatesPercentage"` // Always 2 decimal places
	BuildingPercentage       Decimal2  `json:"buildingPercentage"`    // Always 2 decimal places
	FloorPercentage          Decimal2  `json:"floorPercentage"`       // Always 2 decimal places
	TotalBuildings           int       `json:"totalBuildings"`
	TotalFloors              int       `json:"totalFloors"`
	AverageLocationsPerDay   Decimal2  `json:"averageLocationsPerDay"` // Always 2 decimal places
	LatestCreationDate       time.Time `json:"latestCreationDate"`
	EarliestCreationDate     time.Time `json:"earliestCreationDate"`
}
