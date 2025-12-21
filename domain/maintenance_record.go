package domain

import "time"

// --- Enums ---

type MaintenanceRecordSortField string

const (
	MaintenanceRecordSortByMaintenanceDate MaintenanceRecordSortField = "maintenanceDate"
	MaintenanceRecordSortByActualCost      MaintenanceRecordSortField = "actualCost"
	MaintenanceRecordSortByCreatedAt       MaintenanceRecordSortField = "createdAt"
	MaintenanceRecordSortByUpdatedAt       MaintenanceRecordSortField = "updatedAt"
)

type MaintenanceResult string

const (
	ResultSuccess     MaintenanceResult = "Success"
	ResultPartial     MaintenanceResult = "Partial"
	ResultFailed      MaintenanceResult = "Failed"
	ResultRescheduled MaintenanceResult = "Rescheduled"
)

type MaintenanceRecord struct {
	ID                string                         `json:"id"`
	ScheduleID        *string                        `json:"scheduleId"`
	AssetID           string                         `json:"assetId"`
	MaintenanceDate   time.Time                      `json:"maintenanceDate"`
	CompletionDate    *time.Time                     `json:"completionDate"`
	DurationMinutes   *int                           `json:"durationMinutes"`
	PerformedByUser   *string                        `json:"performedByUser"`
	PerformedByVendor *string                        `json:"performedByVendor"`
	Result            MaintenanceResult              `json:"result"`
	ActualCost        *float64                       `json:"actualCost"`
	CreatedAt         time.Time                      `json:"createdAt"`
	UpdatedAt         time.Time                      `json:"updatedAt"`
	Translations      []MaintenanceRecordTranslation `json:"translations,omitempty"`
	// * Preloaded relationships
	Schedule *MaintenanceSchedule `json:"schedule,omitempty"`
	Asset    *Asset               `json:"asset,omitempty"`
	User     *User                `json:"user,omitempty"`
}

type MaintenanceRecordTranslation struct {
	ID       string  `json:"id"`
	RecordID string  `json:"recordId"`
	LangCode string  `json:"langCode"`
	Title    string  `json:"title"`
	Notes    *string `json:"notes"`
}

type MaintenanceRecordTranslationResponse struct {
	LangCode string  `json:"langCode"`
	Title    string  `json:"title"`
	Notes    *string `json:"notes"`
}

type MaintenanceRecordResponse struct {
	ID                string                                 `json:"id"`
	ScheduleID        *string                                `json:"scheduleId"`
	AssetID           string                                 `json:"assetId"`
	MaintenanceDate   time.Time                              `json:"maintenanceDate"`
	CompletionDate    *time.Time                             `json:"completionDate"`
	DurationMinutes   *int                                   `json:"durationMinutes"`
	PerformedByUserID *string                                `json:"performedByUserId"`
	PerformedByVendor *string                                `json:"performedByVendor"`
	Result            MaintenanceResult                      `json:"result"`
	ActualCost        *NullableDecimal2                      `json:"actualCost"` // Custom type to ensure 2 decimal places as number
	Title             string                                 `json:"title"`
	Notes             *string                                `json:"notes"`
	CreatedAt         time.Time                              `json:"createdAt"`
	UpdatedAt         time.Time                              `json:"updatedAt"`
	Translations      []MaintenanceRecordTranslationResponse `json:"translations"`
	// * Populated
	Schedule        *MaintenanceScheduleResponse `json:"schedule"`
	Asset           AssetResponse                `json:"asset"`
	PerformedByUser *UserResponse                `json:"performedByUser"`
}

type MaintenanceRecordListResponse struct {
	ID                string            `json:"id"`
	ScheduleID        *string           `json:"scheduleId"`
	AssetID           string            `json:"assetId"`
	MaintenanceDate   time.Time         `json:"maintenanceDate"`
	CompletionDate    *time.Time        `json:"completionDate"`
	DurationMinutes   *int              `json:"durationMinutes"`
	PerformedByUserID *string           `json:"performedByUserId"`
	PerformedByVendor *string           `json:"performedByVendor"`
	Result            MaintenanceResult `json:"result"`
	ActualCost        *NullableDecimal2 `json:"actualCost"` // Custom type to ensure 2 decimal places as number
	Title             string            `json:"title"`
	Notes             *string           `json:"notes"`
	CreatedAt         time.Time         `json:"createdAt"`
	UpdatedAt         time.Time         `json:"updatedAt"`
	// * Populated
	Schedule        *MaintenanceScheduleListResponse `json:"schedule"`
	Asset           AssetResponse                    `json:"asset"`
	PerformedByUser *UserResponse                    `json:"performedByUser"`
}

type BulkDeleteMaintenanceRecords struct {
	RequestedIDS []string `json:"requestedIds"`
	DeletedIDS   []string `json:"deletedIds"`
}

type BulkDeleteMaintenanceRecordsResponse struct {
	RequestedIDS []string `json:"requestedIds"`
	DeletedIDS   []string `json:"deletedIds"`
}

// --- Bulk Create ---

type BulkCreateMaintenanceRecordsPayload struct {
	MaintenanceRecords []CreateMaintenanceRecordPayload `json:"maintenanceRecords" validate:"required,min=1,max=100,dive"`
}

type BulkCreateMaintenanceRecordsResponse struct {
	MaintenanceRecords []MaintenanceRecordResponse `json:"maintenanceRecords"`
}

// --- Payloads ---

type CreateMaintenanceRecordPayload struct {
	ScheduleID        *string                                     `json:"scheduleId,omitempty"`
	AssetID           string                                      `json:"assetId" validate:"required"`
	MaintenanceDate   string                                      `json:"maintenanceDate" validate:"required,datetime=2006-01-02"`
	CompletionDate    *string                                     `json:"completionDate,omitempty" validate:"omitempty,datetime=2006-01-02"`
	DurationMinutes   *int                                        `json:"durationMinutes,omitempty" validate:"omitempty,gt=0"`
	PerformedByUser   *string                                     `json:"performedByUser,omitempty" validate:"omitempty"`
	PerformedByVendor *string                                     `json:"performedByVendor,omitempty" validate:"omitempty,max=150"`
	Result            MaintenanceResult                           `json:"result" validate:"required,oneof=Success Partial Failed Rescheduled"`
	ActualCost        *float64                                    `json:"actualCost,omitempty" validate:"omitempty,gt=0"`
	Translations      []CreateMaintenanceRecordTranslationPayload `json:"translations" validate:"required,min=1,dive"`
}

type CreateMaintenanceRecordTranslationPayload struct {
	LangCode string  `json:"langCode" validate:"required,max=5"`
	Title    string  `json:"title" validate:"required,max=200"`
	Notes    *string `json:"notes,omitempty"`
}

type UpdateMaintenanceRecordPayload struct {
	ScheduleID        *string                                     `json:"scheduleId,omitempty"`
	MaintenanceDate   *string                                     `json:"maintenanceDate,omitempty" validate:"omitempty,datetime=2006-01-02"`
	CompletionDate    *string                                     `json:"completionDate,omitempty" validate:"omitempty,datetime=2006-01-02"`
	DurationMinutes   *int                                        `json:"durationMinutes,omitempty" validate:"omitempty,gt=0"`
	PerformedByUser   *string                                     `json:"performedByUser,omitempty" validate:"omitempty"`
	PerformedByVendor *string                                     `json:"performedByVendor,omitempty" validate:"omitempty,max=150"`
	Result            *MaintenanceResult                          `json:"result,omitempty" validate:"omitempty,oneof=Success Partial Failed Rescheduled"`
	ActualCost        *float64                                    `json:"actualCost,omitempty" validate:"omitempty,gt=0"`
	Translations      []UpdateMaintenanceRecordTranslationPayload `json:"translations,omitempty" validate:"omitempty,dive"`
}

type UpdateMaintenanceRecordTranslationPayload struct {
	LangCode string  `json:"langCode" validate:"required,max=5"`
	Title    *string `json:"title,omitempty" validate:"omitempty,max=200"`
	Notes    *string `json:"notes,omitempty"`
}

type BulkDeleteMaintenanceRecordsPayload struct {
	IDS []string `json:"ids" validate:"required,min=1,max=100,dive,required"`
}

type ExportMaintenanceRecordListPayload struct {
	Format      ExportFormat                    `json:"format" validate:"required,oneof=pdf excel"`
	SearchQuery *string                         `json:"searchQuery,omitempty"`
	Filters     *MaintenanceRecordFilterOptions `json:"filters,omitempty"`
	Sort        *MaintenanceRecordSortOptions   `json:"sort,omitempty"`
}

// --- Query Parameters ---

type MaintenanceRecordFilterOptions struct {
	AssetID         *string `json:"assetId,omitempty"`
	ScheduleID      *string `json:"scheduleId,omitempty"`
	PerformedByUser *string `json:"performedByUser,omitempty"`
	VendorName      *string `json:"vendorName,omitempty"`
	FromDate        *string `json:"fromDate,omitempty"` // YYYY-MM-DD
	ToDate          *string `json:"toDate,omitempty"`   // YYYY-MM-DD
}

type MaintenanceRecordSortOptions struct {
	Field MaintenanceRecordSortField `json:"field" example:"createdAt"`
	Order SortOrder                  `json:"order" example:"desc"`
}

type MaintenanceRecordParams struct {
	SearchQuery *string                         `json:"searchQuery,omitempty"`
	Filters     *MaintenanceRecordFilterOptions `json:"filters,omitempty"`
	Sort        *MaintenanceRecordSortOptions   `json:"sort,omitempty"`
	Pagination  *PaginationOptions              `json:"pagination,omitempty"`
}

// --- Statistics ---

// Internal statistics structs (used in repository layer)
type MaintenanceRecordStatistics struct {
	Total           MaintenanceRecordCountStatistics    `json:"total"`
	ByPerformer     []UserMaintenanceRecordStatistics   `json:"byPerformer"`
	ByVendor        []VendorMaintenanceRecordStatistics `json:"byVendor"`
	ByAsset         []AssetMaintenanceRecordStatistics  `json:"byAsset"`
	CostStatistics  MaintenanceRecordCostStatistics     `json:"costStatistics"`
	CompletionTrend []MaintenanceRecordCompletionTrend  `json:"completionTrend"`
	MonthlyTrends   []MaintenanceRecordMonthlyTrend     `json:"monthlyTrends"`
	Summary         MaintenanceRecordSummaryStatistics  `json:"summary"`
}

type MaintenanceRecordCountStatistics struct {
	Count int `json:"count"`
}

type UserMaintenanceRecordStatistics struct {
	UserID      string  `json:"userId"`
	UserName    string  `json:"userName"`
	UserEmail   string  `json:"userEmail"`
	Count       int     `json:"count"`
	TotalCost   float64 `json:"totalCost"`
	AverageCost float64 `json:"averageCost"`
}

type VendorMaintenanceRecordStatistics struct {
	VendorName  string  `json:"vendorName"`
	Count       int     `json:"count"`
	TotalCost   float64 `json:"totalCost"`
	AverageCost float64 `json:"averageCost"`
}

type AssetMaintenanceRecordStatistics struct {
	AssetID         string  `json:"assetId"`
	AssetName       string  `json:"assetName"`
	AssetTag        string  `json:"assetTag"`
	RecordCount     int     `json:"recordCount"`
	LastMaintenance string  `json:"lastMaintenance"`
	TotalCost       float64 `json:"totalCost"`
	AverageCost     float64 `json:"averageCost"`
}

type MaintenanceRecordCostStatistics struct {
	TotalCost          *float64 `json:"totalCost"`
	AverageCost        *float64 `json:"averageCost"`
	MinCost            *float64 `json:"minCost"`
	MaxCost            *float64 `json:"maxCost"`
	RecordsWithCost    int      `json:"recordsWithCost"`
	RecordsWithoutCost int      `json:"recordsWithoutCost"`
}

type MaintenanceRecordCompletionTrend struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

type MaintenanceRecordMonthlyTrend struct {
	Month       string  `json:"month"`
	RecordCount int     `json:"recordCount"`
	TotalCost   float64 `json:"totalCost"`
}

type MaintenanceRecordSummaryStatistics struct {
	TotalRecords                  int      `json:"totalRecords"`
	RecordsWithCostInfo           int      `json:"recordsWithCostInfo"`
	CostInfoPercentage            float64  `json:"costInfoPercentage"`
	TotalUniqueVendors            int      `json:"totalUniqueVendors"`
	TotalUniquePerformers         int      `json:"totalUniquePerformers"`
	AverageRecordsPerDay          float64  `json:"averageRecordsPerDay"`
	LatestRecordDate              string   `json:"latestRecordDate"`
	EarliestRecordDate            string   `json:"earliestRecordDate"`
	MostExpensiveMaintenanceCost  *float64 `json:"mostExpensiveMaintenanceCost"`
	LeastExpensiveMaintenanceCost *float64 `json:"leastExpensiveMaintenanceCost"`
	AssetsWithMaintenance         int      `json:"assetsWithMaintenance"`
	AverageMaintenancePerAsset    float64  `json:"averageMaintenancePerAsset"`
}

// Response statistics structs (used in service/handler layer)
type MaintenanceRecordStatisticsResponse struct {
	Total           MaintenanceRecordCountStatisticsResponse    `json:"total"`
	ByPerformer     []UserMaintenanceRecordStatisticsResponse   `json:"byPerformer"`
	ByVendor        []VendorMaintenanceRecordStatisticsResponse `json:"byVendor"`
	ByAsset         []AssetMaintenanceRecordStatisticsResponse  `json:"byAsset"`
	CostStatistics  MaintenanceRecordCostStatisticsResponse     `json:"costStatistics"`
	CompletionTrend []MaintenanceRecordCompletionTrendResponse  `json:"completionTrend"`
	MonthlyTrends   []MaintenanceRecordMonthlyTrendResponse     `json:"monthlyTrends"`
	Summary         MaintenanceRecordSummaryStatisticsResponse  `json:"summary"`
}

type MaintenanceRecordCountStatisticsResponse struct {
	Count int `json:"count"`
}

type UserMaintenanceRecordStatisticsResponse struct {
	UserID      string   `json:"userId"`
	UserName    string   `json:"userName"`
	UserEmail   string   `json:"userEmail"`
	Count       int      `json:"count"`
	TotalCost   Decimal2 `json:"totalCost"`   // Custom type to ensure 2 decimal places as number
	AverageCost Decimal2 `json:"averageCost"` // Custom type to ensure 2 decimal places as number
}

type VendorMaintenanceRecordStatisticsResponse struct {
	VendorName  string   `json:"vendorName"`
	Count       int      `json:"count"`
	TotalCost   Decimal2 `json:"totalCost"`   // Custom type to ensure 2 decimal places as number
	AverageCost Decimal2 `json:"averageCost"` // Custom type to ensure 2 decimal places as number
}

type AssetMaintenanceRecordStatisticsResponse struct {
	AssetID         string   `json:"assetId"`
	AssetName       string   `json:"assetName"`
	AssetTag        string   `json:"assetTag"`
	RecordCount     int      `json:"recordCount"`
	LastMaintenance string   `json:"lastMaintenance"`
	TotalCost       Decimal2 `json:"totalCost"`   // Custom type to ensure 2 decimal places as number
	AverageCost     Decimal2 `json:"averageCost"` // Custom type to ensure 2 decimal places as number
}

type MaintenanceRecordCostStatisticsResponse struct {
	TotalCost          *NullableDecimal2 `json:"totalCost"`   // Custom type to ensure 2 decimal places as number
	AverageCost        *NullableDecimal2 `json:"averageCost"` // Custom type to ensure 2 decimal places as number
	MinCost            *NullableDecimal2 `json:"minCost"`     // Custom type to ensure 2 decimal places as number
	MaxCost            *NullableDecimal2 `json:"maxCost"`     // Custom type to ensure 2 decimal places as number
	RecordsWithCost    int               `json:"recordsWithCost"`
	RecordsWithoutCost int               `json:"recordsWithoutCost"`
}

type MaintenanceRecordCompletionTrendResponse struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

type MaintenanceRecordMonthlyTrendResponse struct {
	Month       string   `json:"month"`
	RecordCount int      `json:"recordCount"`
	TotalCost   Decimal2 `json:"totalCost"` // Custom type to ensure 2 decimal places as number
}

type MaintenanceRecordSummaryStatisticsResponse struct {
	TotalRecords                  int               `json:"totalRecords"`
	RecordsWithCostInfo           int               `json:"recordsWithCostInfo"`
	CostInfoPercentage            Decimal2          `json:"costInfoPercentage"`
	TotalUniqueVendors            int               `json:"totalUniqueVendors"`
	TotalUniquePerformers         int               `json:"totalUniquePerformers"`
	AverageRecordsPerDay          Decimal2          `json:"averageRecordsPerDay"`
	LatestRecordDate              string            `json:"latestRecordDate"`
	EarliestRecordDate            string            `json:"earliestRecordDate"`
	MostExpensiveMaintenanceCost  *NullableDecimal2 `json:"mostExpensiveMaintenanceCost"`  // Custom type to ensure 2 decimal places as number
	LeastExpensiveMaintenanceCost *NullableDecimal2 `json:"leastExpensiveMaintenanceCost"` // Custom type to ensure 2 decimal places as number
	AssetsWithMaintenance         int               `json:"assetsWithMaintenance"`
	AverageMaintenancePerAsset    Decimal2          `json:"averageMaintenancePerAsset"`
}
