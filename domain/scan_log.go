package domain

import (
	"time"
)

// --- Enums ---

type ScanMethodType string

const (
	ScanMethodDataMatrix  ScanMethodType = "DATA_MATRIX"
	ScanMethodManualInput ScanMethodType = "MANUAL_INPUT"
)

type ScanResultType string

const (
	ScanResultSuccess       ScanResultType = "Success"
	ScanResultInvalidID     ScanResultType = "Invalid ID"
	ScanResultAssetNotFound ScanResultType = "Asset Not Found"
)

type ScanLogSortField string

const (
	ScanLogSortByScannedValue  ScanLogSortField = "scannedValue"
	ScanLogSortByScanMethod    ScanLogSortField = "scanMethod"
	ScanLogSortByScanTimestamp ScanLogSortField = "scanTimestamp"
	ScanLogSortByScanResult    ScanLogSortField = "scanResult"
)

// --- Structs ---

type ScanLog struct {
	ID              string         `json:"id"`
	AssetID         *string        `json:"assetId"`
	ScannedValue    string         `json:"scannedValue"`
	ScanMethod      ScanMethodType `json:"scanMethod"`
	ScannedBy       string         `json:"scannedBy"`
	ScanTimestamp   time.Time      `json:"scanTimestamp"`
	ScanLocationLat *float64       `json:"scanLocationLat"`
	ScanLocationLng *float64       `json:"scanLocationLng"`
	ScanResult      ScanResultType `json:"scanResult"`
}

type ScanLogResponse struct {
	ID              string         `json:"id"`
	AssetID         *string        `json:"assetId"`
	ScannedValue    string         `json:"scannedValue"`
	ScanMethod      ScanMethodType `json:"scanMethod"`
	ScannedByID     string         `json:"scannedById"`
	ScanTimestamp   time.Time      `json:"scanTimestamp"`
	ScanLocationLat *float64       `json:"scanLocationLat"`
	ScanLocationLng *float64       `json:"scanLocationLng"`
	ScanResult      ScanResultType `json:"scanResult"`
	// * Populated
	// ! cuma scan log gak perlu populated table biar gak berat
	// Asset     *AssetResponse `json:"asset,omitempty"`
	// ScannedBy UserResponse   `json:"scannedBy"`
}

type ScanLogListResponse struct {
	ID              string         `json:"id"`
	AssetID         *string        `json:"assetId"`
	ScannedValue    string         `json:"scannedValue"`
	ScanMethod      ScanMethodType `json:"scanMethod"`
	ScannedByID     string         `json:"scannedById"`
	ScanTimestamp   time.Time      `json:"scanTimestamp"`
	ScanLocationLat *float64       `json:"scanLocationLat"`
	ScanLocationLng *float64       `json:"scanLocationLng"`
	ScanResult      ScanResultType `json:"scanResult"`
}

type BulkDeleteScanLogs struct {
	RequestedIDS []string `json:"requestedIds"`
	DeletedIDS   []string `json:"deletedIds"`
}

type BulkDeleteScanLogsResponse struct {
	RequestedIDS []string `json:"requestedIds"`
	DeletedIDS   []string `json:"deletedIds"`
}

// --- Payloads ---

type CreateScanLogPayload struct {
	AssetID         *string        `json:"assetId"`
	ScannedValue    string         `json:"scannedValue" validate:"required"`
	ScanMethod      ScanMethodType `json:"scanMethod" validate:"required,oneof=DATA_MATRIX MANUAL_INPUT"`
	ScanLocationLat *float64       `json:"scanLocationLat,omitempty" validate:"omitempty,latitude"`
	ScanLocationLng *float64       `json:"scanLocationLng,omitempty" validate:"omitempty,longitude"`
	ScanResult      ScanResultType `json:"scanResult"`
}

type BulkDeleteScanLogsPayload struct {
	IDS []string `json:"ids" validate:"required,min=1,max=100,dive,required"`
}

// --- Query Parameters ---

type ScanLogFilterOptions struct {
	ScanMethod     *ScanMethodType `json:"scanMethod,omitempty"`
	ScanResult     *ScanResultType `json:"scanResult,omitempty"`
	ScannedBy      *string         `json:"scannedBy,omitempty"`
	AssetID        *string         `json:"assetId,omitempty"`
	DateFrom       *time.Time      `json:"dateFrom,omitempty"`
	DateTo         *time.Time      `json:"dateTo,omitempty"`
	HasCoordinates *bool           `json:"hasCoordinates,omitempty"`
}

type ScanLogSortOptions struct {
	Field ScanLogSortField `json:"field" example:"scanTimestamp"`
	Order SortOrder        `json:"order" example:"desc"`
}

type ScanLogParams struct {
	SearchQuery *string               `json:"searchQuery,omitempty"`
	Filters     *ScanLogFilterOptions `json:"filters,omitempty"`
	Sort        *ScanLogSortOptions   `json:"sort,omitempty"`
	Pagination  *PaginationOptions    `json:"pagination,omitempty"`
}

// --- Statistics ---

// Internal statistics structs (used in repository layer)
type ScanLogStatistics struct {
	Total       ScanLogCountStatistics   `json:"total"`
	ByMethod    []ScanMethodStatistics   `json:"byMethod"`
	ByResult    []ScanResultStatistics   `json:"byResult"`
	Geographic  ScanGeographicStatistics `json:"geographic"`
	ScanTrends  []ScanTrend              `json:"scanTrends"`
	TopScanners []ScannerStatistics      `json:"topScanners"`
	Summary     ScanLogSummaryStatistics `json:"summary"`
}

type ScanLogCountStatistics struct {
	Count int `json:"count"`
}

type ScanMethodStatistics struct {
	Method ScanMethodType `json:"method"`
	Count  int            `json:"count"`
}

type ScanResultStatistics struct {
	Result ScanResultType `json:"result"`
	Count  int            `json:"count"`
}

type ScanGeographicStatistics struct {
	WithCoordinates    int `json:"withCoordinates"`
	WithoutCoordinates int `json:"withoutCoordinates"`
}

type ScanTrend struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

type ScannerStatistics struct {
	UserID string `json:"userId"`
	Count  int    `json:"count"`
}

type ScanLogSummaryStatistics struct {
	TotalScans            int       `json:"totalScans"`
	SuccessRate           float64   `json:"successRate"`
	ScansWithCoordinates  int       `json:"scansWithCoordinates"`
	CoordinatesPercentage float64   `json:"coordinatesPercentage"`
	AverageScansPerDay    float64   `json:"averageScansPerDay"`
	LatestScanDate        time.Time `json:"latestScanDate"`
	EarliestScanDate      time.Time `json:"earliestScanDate"`
}

// Response statistics structs (used in service/handler layer)
type ScanLogStatisticsResponse struct {
	Total       ScanLogCountStatisticsResponse   `json:"total"`
	ByMethod    []ScanMethodStatisticsResponse   `json:"byMethod"`
	ByResult    []ScanResultStatisticsResponse   `json:"byResult"`
	Geographic  ScanGeographicStatisticsResponse `json:"geographic"`
	ScanTrends  []ScanTrendResponse              `json:"scanTrends"`
	TopScanners []ScannerStatisticsResponse      `json:"topScanners"`
	Summary     ScanLogSummaryStatisticsResponse `json:"summary"`
}

type ScanLogCountStatisticsResponse struct {
	Count int `json:"count"`
}

type ScanMethodStatisticsResponse struct {
	Method ScanMethodType `json:"method"`
	Count  int            `json:"count"`
}

type ScanResultStatisticsResponse struct {
	Result ScanResultType `json:"result"`
	Count  int            `json:"count"`
}

type ScanGeographicStatisticsResponse struct {
	WithCoordinates    int `json:"withCoordinates"`
	WithoutCoordinates int `json:"withoutCoordinates"`
}

type ScanTrendResponse struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

type ScannerStatisticsResponse struct {
	UserID string `json:"userId"`
	Count  int    `json:"count"`
}

type ScanLogSummaryStatisticsResponse struct {
	TotalScans            int       `json:"totalScans"`
	SuccessRate           float64   `json:"successRate"`
	ScansWithCoordinates  int       `json:"scansWithCoordinates"`
	CoordinatesPercentage float64   `json:"coordinatesPercentage"`
	AverageScansPerDay    float64   `json:"averageScansPerDay"`
	LatestScanDate        time.Time `json:"latestScanDate"`
	EarliestScanDate      time.Time `json:"earliestScanDate"`
}
