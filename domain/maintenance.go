package domain

import "time"

// --- Enums ---

type MaintenanceScheduleType string

const (
	ScheduleTypePreventive MaintenanceScheduleType = "Preventive"
	ScheduleTypeCorrective MaintenanceScheduleType = "Corrective"
)

type ScheduleStatus string

const (
	StatusScheduled ScheduleStatus = "Scheduled"
	StatusCompleted ScheduleStatus = "Completed"
	StatusCancelled ScheduleStatus = "Cancelled"
)

// --- Structs ---

type MaintenanceSchedule struct {
	ID              string                           `json:"id"`
	AssetID         string                           `json:"assetId"`
	MaintenanceType MaintenanceScheduleType          `json:"maintenanceType"`
	ScheduledDate   time.Time                        `json:"scheduledDate"`
	FrequencyMonths *int                             `json:"frequencyMonths"`
	Status          ScheduleStatus                   `json:"status"`
	CreatedBy       string                           `json:"createdBy"`
	CreatedAt       time.Time                        `json:"createdAt"`
	Translations    []MaintenanceScheduleTranslation `json:"translations,omitempty"`
}

type MaintenanceScheduleTranslation struct {
	ID          string  `json:"id"`
	ScheduleID  string  `json:"scheduleId"`
	LangCode    string  `json:"langCode"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type MaintenanceScheduleTranslationResponse struct {
	LangCode    string  `json:"langCode"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type MaintenanceRecord struct {
	ID                string                         `json:"id"`
	ScheduleID        *string                        `json:"scheduleId"`
	AssetID           string                         `json:"assetId"`
	MaintenanceDate   time.Time                      `json:"maintenanceDate"`
	PerformedByUser   *string                        `json:"performedByUser"`
	PerformedByVendor *string                        `json:"performedByVendor"`
	ActualCost        *float64                       `json:"actualCost"`
	CreatedAt         time.Time                      `json:"createdAt"`
	UpdatedAt         time.Time                      `json:"updatedAt"`
	Translations      []MaintenanceRecordTranslation `json:"translations,omitempty"`
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

type MaintenanceScheduleResponse struct {
	ID              string                                   `json:"id"`
	AssetID         string                                   `json:"assetId"`
	MaintenanceType MaintenanceScheduleType                  `json:"maintenanceType"`
	ScheduledDate   string                                   `json:"scheduledDate"`
	FrequencyMonths *int                                     `json:"frequencyMonths,omitempty"`
	Status          ScheduleStatus                           `json:"status"`
	CreatedByID     string                                   `json:"createdById"`
	CreatedAt       string                                   `json:"createdAt"`
	Title           string                                   `json:"title"`
	Description     *string                                  `json:"description,omitempty"`
	Translations    []MaintenanceScheduleTranslationResponse `json:"translations"`
	// * Populated
	Asset     AssetResponse `json:"asset"`
	CreatedBy UserResponse  `json:"createdBy"`
}

type MaintenanceScheduleListResponse struct {
	ID              string                  `json:"id"`
	AssetID         string                  `json:"assetId"`
	MaintenanceType MaintenanceScheduleType `json:"maintenanceType"`
	ScheduledDate   string                  `json:"scheduledDate"`
	FrequencyMonths *int                    `json:"frequencyMonths,omitempty"`
	Status          ScheduleStatus          `json:"status"`
	CreatedByID     string                  `json:"createdById"`
	CreatedAt       string                  `json:"createdAt"`
	Title           string                  `json:"title"`
	Description     *string                 `json:"description,omitempty"`
	// * Populated
	Asset     AssetResponse `json:"asset"`
	CreatedBy UserResponse  `json:"createdBy"`
}

type MaintenanceRecordResponse struct {
	ID                string                                 `json:"id"`
	ScheduleID        *string                                `json:"scheduleId,omitempty"`
	AssetID           string                                 `json:"assetId"`
	MaintenanceDate   string                                 `json:"maintenanceDate"`
	PerformedByUserID *string                                `json:"performedByUserId,omitempty"`
	PerformedByVendor *string                                `json:"performedByVendor,omitempty"`
	ActualCost        *float64                               `json:"actualCost,omitempty"`
	Title             string                                 `json:"title"`
	Notes             *string                                `json:"notes,omitempty"`
	CreatedAt         string                                 `json:"createdAt"`
	UpdatedAt         string                                 `json:"updatedAt"`
	Translations      []MaintenanceRecordTranslationResponse `json:"translations"`
	// * Populated
	Schedule        *MaintenanceScheduleResponse `json:"schedule,omitempty"`
	Asset           AssetResponse                `json:"asset"`
	PerformedByUser *UserResponse                `json:"performedByUser,omitempty"`
}

type MaintenanceRecordListResponse struct {
	ID                string   `json:"id"`
	ScheduleID        *string  `json:"scheduleId,omitempty"`
	AssetID           string   `json:"assetId"`
	MaintenanceDate   string   `json:"maintenanceDate"`
	PerformedByUserID *string  `json:"performedByUserId,omitempty"`
	PerformedByVendor *string  `json:"performedByVendor,omitempty"`
	ActualCost        *float64 `json:"actualCost,omitempty"`
	Title             string   `json:"title"`
	Notes             *string  `json:"notes,omitempty"`
	CreatedAt         string   `json:"createdAt"`
	UpdatedAt         string   `json:"updatedAt"`
	// * Populated
	Schedule        *MaintenanceScheduleResponse `json:"schedule,omitempty"`
	Asset           AssetResponse                `json:"asset"`
	PerformedByUser *UserResponse                `json:"performedByUser,omitempty"`
}

// --- Payloads ---

type CreateMaintenanceSchedulePayload struct {
	AssetID         string                                        `json:"assetId" validate:"required"`
	MaintenanceType MaintenanceScheduleType                       `json:"maintenanceType" validate:"required,oneof=Preventive Corrective"`
	ScheduledDate   string                                        `json:"scheduledDate" validate:"required,datetime=2006-01-02"`
	FrequencyMonths *int                                          `json:"frequencyMonths,omitempty" validate:"omitempty,gt=0"`
	Translations    []CreateMaintenanceScheduleTranslationPayload `json:"translations" validate:"required,min=1,dive"`
}

type CreateMaintenanceScheduleTranslationPayload struct {
	LangCode    string  `json:"langCode" validate:"required,max=5"`
	Title       string  `json:"title" validate:"required,max=200"`
	Description *string `json:"description,omitempty"`
}

type CreateMaintenanceRecordPayload struct {
	ScheduleID        *string                                     `json:"scheduleId,omitempty"`
	AssetID           string                                      `json:"assetId" validate:"required"`
	MaintenanceDate   string                                      `json:"maintenanceDate" validate:"required,datetime=2006-01-02"`
	PerformedByUser   *string                                     `json:"performedByUser,omitempty" validate:"omitempty"`
	PerformedByVendor *string                                     `json:"performedByVendor,omitempty" validate:"omitempty,max=150"`
	ActualCost        *float64                                    `json:"actualCost,omitempty" validate:"omitempty,gt=0"`
	Translations      []CreateMaintenanceRecordTranslationPayload `json:"translations" validate:"required,min=1,dive"`
}

type CreateMaintenanceRecordTranslationPayload struct {
	LangCode string  `json:"langCode" validate:"required,max=5"`
	Title    string  `json:"title" validate:"required,max=200"`
	Notes    *string `json:"notes,omitempty"`
}

// --- Statistics ---

// Internal statistics structs (used in repository layer)
type MaintenanceStatistics struct {
	Schedules MaintenanceScheduleStatistics `json:"schedules"`
	Records   MaintenanceRecordStatistics   `json:"records"`
	Summary   MaintenanceSummaryStatistics  `json:"summary"`
}

type MaintenanceScheduleStatistics struct {
	Total            MaintenanceCountStatistics          `json:"total"`
	ByType           MaintenanceTypeStatistics           `json:"byType"`
	ByStatus         MaintenanceScheduleStatusStatistics `json:"byStatus"`
	ByAsset          []AssetMaintenanceStatistics        `json:"byAsset"`
	ByCreator        []UserMaintenanceStatistics         `json:"byCreator"`
	UpcomingSchedule []UpcomingMaintenanceSchedule       `json:"upcomingSchedule"`
	OverdueSchedule  []OverdueMaintenanceSchedule        `json:"overdueSchedule"`
	FrequencyTrends  []MaintenanceFrequencyTrend         `json:"frequencyTrends"`
}

type MaintenanceRecordStatistics struct {
	Total           MaintenanceCountStatistics    `json:"total"`
	ByPerformer     []UserMaintenanceStatistics   `json:"byPerformer"`
	ByVendor        []VendorMaintenanceStatistics `json:"byVendor"`
	ByAsset         []AssetMaintenanceStatistics  `json:"byAsset"`
	CostStatistics  MaintenanceCostStatistics     `json:"costStatistics"`
	CompletionTrend []MaintenanceCompletionTrend  `json:"completionTrend"`
	MonthlyTrends   []MaintenanceMonthlyTrend     `json:"monthlyTrends"`
}

type MaintenanceCountStatistics struct {
	Count int `json:"count"`
}

type MaintenanceTypeStatistics struct {
	Preventive int `json:"preventive"`
	Corrective int `json:"corrective"`
}

type MaintenanceScheduleStatusStatistics struct {
	Scheduled int `json:"scheduled"`
	Completed int `json:"completed"`
	Cancelled int `json:"cancelled"`
}

type AssetMaintenanceStatistics struct {
	AssetID         string `json:"assetId"`
	AssetName       string `json:"assetName"`
	AssetTag        string `json:"assetTag"`
	ScheduleCount   int    `json:"scheduleCount"`
	RecordCount     int    `json:"recordCount"`
	LastMaintenance string `json:"lastMaintenance"`
	NextMaintenance string `json:"nextMaintenance"`
}

type UserMaintenanceStatistics struct {
	UserID      string  `json:"userId"`
	UserName    string  `json:"userName"`
	UserEmail   string  `json:"userEmail"`
	Count       int     `json:"count"`
	TotalCost   float64 `json:"totalCost"`
	AverageCost float64 `json:"averageCost"`
}

type VendorMaintenanceStatistics struct {
	VendorName  string  `json:"vendorName"`
	Count       int     `json:"count"`
	TotalCost   float64 `json:"totalCost"`
	AverageCost float64 `json:"averageCost"`
}

type UpcomingMaintenanceSchedule struct {
	ID              string                  `json:"id"`
	AssetID         string                  `json:"assetId"`
	AssetName       string                  `json:"assetName"`
	AssetTag        string                  `json:"assetTag"`
	MaintenanceType MaintenanceScheduleType `json:"maintenanceType"`
	ScheduledDate   string                  `json:"scheduledDate"`
	DaysUntilDue    int                     `json:"daysUntilDue"`
	Title           string                  `json:"title"`
	Description     *string                 `json:"description,omitempty"`
}

type OverdueMaintenanceSchedule struct {
	ID              string                  `json:"id"`
	AssetID         string                  `json:"assetId"`
	AssetName       string                  `json:"assetName"`
	AssetTag        string                  `json:"assetTag"`
	MaintenanceType MaintenanceScheduleType `json:"maintenanceType"`
	ScheduledDate   string                  `json:"scheduledDate"`
	DaysOverdue     int                     `json:"daysOverdue"`
	Title           string                  `json:"title"`
	Description     *string                 `json:"description,omitempty"`
}

type MaintenanceFrequencyTrend struct {
	FrequencyMonths int `json:"frequencyMonths"`
	Count           int `json:"count"`
}

type MaintenanceCostStatistics struct {
	TotalCost          *float64 `json:"totalCost"`
	AverageCost        *float64 `json:"averageCost"`
	MinCost            *float64 `json:"minCost"`
	MaxCost            *float64 `json:"maxCost"`
	RecordsWithCost    int      `json:"recordsWithCost"`
	RecordsWithoutCost int      `json:"recordsWithoutCost"`
}

type MaintenanceCompletionTrend struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type MaintenanceMonthlyTrend struct {
	Month           string  `json:"month"`
	ScheduleCount   int     `json:"scheduleCount"`
	RecordCount     int     `json:"recordCount"`
	TotalCost       float64 `json:"totalCost"`
	PreventiveCount int     `json:"preventiveCount"`
	CorrectiveCount int     `json:"correctiveCount"`
}

type MaintenanceSummaryStatistics struct {
	TotalSchedules                  int      `json:"totalSchedules"`
	TotalRecords                    int      `json:"totalRecords"`
	ScheduledMaintenancePercentage  float64  `json:"scheduledMaintenancePercentage"`
	CompletedMaintenancePercentage  float64  `json:"completedMaintenancePercentage"`
	CancelledMaintenancePercentage  float64  `json:"cancelledMaintenancePercentage"`
	PreventiveMaintenancePercentage float64  `json:"preventiveMaintenancePercentage"`
	CorrectiveMaintenancePercentage float64  `json:"correctiveMaintenancePercentage"`
	AverageMaintenancePerAsset      float64  `json:"averageMaintenancePerAsset"`
	AssetsWithMaintenance           int      `json:"assetsWithMaintenance"`
	AssetsWithoutMaintenance        int      `json:"assetsWithoutMaintenance"`
	MaintenanceComplianceRate       float64  `json:"maintenanceComplianceRate"`
	AverageMaintenanceFrequency     float64  `json:"averageMaintenanceFrequency"`
	UpcomingMaintenanceCount        int      `json:"upcomingMaintenanceCount"`
	OverdueMaintenanceCount         int      `json:"overdueMaintenanceCount"`
	RecordsWithCostInfo             int      `json:"recordsWithCostInfo"`
	CostInfoPercentage              float64  `json:"costInfoPercentage"`
	TotalUniqueVendors              int      `json:"totalUniqueVendors"`
	TotalUniquePerformers           int      `json:"totalUniquePerformers"`
	AverageRecordsPerDay            float64  `json:"averageRecordsPerDay"`
	LatestRecordDate                string   `json:"latestRecordDate"`
	EarliestRecordDate              string   `json:"earliestRecordDate"`
	MostExpensiveMaintenanceCost    *float64 `json:"mostExpensiveMaintenanceCost"`
	LeastExpensiveMaintenanceCost   *float64 `json:"leastExpensiveMaintenanceCost"`
}

// Response statistics structs (used in service/handler layer)
type MaintenanceStatisticsResponse struct {
	Schedules MaintenanceScheduleStatisticsResponse `json:"schedules"`
	Records   MaintenanceRecordStatisticsResponse   `json:"records"`
	Summary   MaintenanceSummaryStatisticsResponse  `json:"summary"`
}

type MaintenanceScheduleStatisticsResponse struct {
	Total            MaintenanceCountStatisticsResponse          `json:"total"`
	ByType           MaintenanceTypeStatisticsResponse           `json:"byType"`
	ByStatus         MaintenanceScheduleStatusStatisticsResponse `json:"byStatus"`
	ByAsset          []AssetMaintenanceStatisticsResponse        `json:"byAsset"`
	ByCreator        []UserMaintenanceStatisticsResponse         `json:"byCreator"`
	UpcomingSchedule []UpcomingMaintenanceScheduleResponse       `json:"upcomingSchedule"`
	OverdueSchedule  []OverdueMaintenanceScheduleResponse        `json:"overdueSchedule"`
	FrequencyTrends  []MaintenanceFrequencyTrendResponse         `json:"frequencyTrends"`
}

type MaintenanceRecordStatisticsResponse struct {
	Total           MaintenanceCountStatisticsResponse    `json:"total"`
	ByPerformer     []UserMaintenanceStatisticsResponse   `json:"byPerformer"`
	ByVendor        []VendorMaintenanceStatisticsResponse `json:"byVendor"`
	ByAsset         []AssetMaintenanceStatisticsResponse  `json:"byAsset"`
	CostStatistics  MaintenanceCostStatisticsResponse     `json:"costStatistics"`
	CompletionTrend []MaintenanceCompletionTrendResponse  `json:"completionTrend"`
	MonthlyTrends   []MaintenanceMonthlyTrendResponse     `json:"monthlyTrends"`
}

type MaintenanceCountStatisticsResponse struct {
	Count int `json:"count"`
}

type MaintenanceTypeStatisticsResponse struct {
	Preventive int `json:"preventive"`
	Corrective int `json:"corrective"`
}

type MaintenanceScheduleStatusStatisticsResponse struct {
	Scheduled int `json:"scheduled"`
	Completed int `json:"completed"`
	Cancelled int `json:"cancelled"`
}

type AssetMaintenanceStatisticsResponse struct {
	AssetID         string `json:"assetId"`
	AssetName       string `json:"assetName"`
	AssetTag        string `json:"assetTag"`
	ScheduleCount   int    `json:"scheduleCount"`
	RecordCount     int    `json:"recordCount"`
	LastMaintenance string `json:"lastMaintenance"`
	NextMaintenance string `json:"nextMaintenance"`
}

type UserMaintenanceStatisticsResponse struct {
	UserID      string  `json:"userId"`
	UserName    string  `json:"userName"`
	UserEmail   string  `json:"userEmail"`
	Count       int     `json:"count"`
	TotalCost   float64 `json:"totalCost"`
	AverageCost float64 `json:"averageCost"`
}

type VendorMaintenanceStatisticsResponse struct {
	VendorName  string  `json:"vendorName"`
	Count       int     `json:"count"`
	TotalCost   float64 `json:"totalCost"`
	AverageCost float64 `json:"averageCost"`
}

type UpcomingMaintenanceScheduleResponse struct {
	ID              string                  `json:"id"`
	AssetID         string                  `json:"assetId"`
	AssetName       string                  `json:"assetName"`
	AssetTag        string                  `json:"assetTag"`
	MaintenanceType MaintenanceScheduleType `json:"maintenanceType"`
	ScheduledDate   string                  `json:"scheduledDate"`
	DaysUntilDue    int                     `json:"daysUntilDue"`
	Title           string                  `json:"title"`
	Description     *string                 `json:"description,omitempty"`
}

type OverdueMaintenanceScheduleResponse struct {
	ID              string                  `json:"id"`
	AssetID         string                  `json:"assetId"`
	AssetName       string                  `json:"assetName"`
	AssetTag        string                  `json:"assetTag"`
	MaintenanceType MaintenanceScheduleType `json:"maintenanceType"`
	ScheduledDate   string                  `json:"scheduledDate"`
	DaysOverdue     int                     `json:"daysOverdue"`
	Title           string                  `json:"title"`
	Description     *string                 `json:"description,omitempty"`
}

type MaintenanceFrequencyTrendResponse struct {
	FrequencyMonths int `json:"frequencyMonths"`
	Count           int `json:"count"`
}

type MaintenanceCostStatisticsResponse struct {
	TotalCost          *float64 `json:"totalCost"`
	AverageCost        *float64 `json:"averageCost"`
	MinCost            *float64 `json:"minCost"`
	MaxCost            *float64 `json:"maxCost"`
	RecordsWithCost    int      `json:"recordsWithCost"`
	RecordsWithoutCost int      `json:"recordsWithoutCost"`
}

type MaintenanceCompletionTrendResponse struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type MaintenanceMonthlyTrendResponse struct {
	Month           string  `json:"month"`
	ScheduleCount   int     `json:"scheduleCount"`
	RecordCount     int     `json:"recordCount"`
	TotalCost       float64 `json:"totalCost"`
	PreventiveCount int     `json:"preventiveCount"`
	CorrectiveCount int     `json:"correctiveCount"`
}

type MaintenanceSummaryStatisticsResponse struct {
	TotalSchedules                  int      `json:"totalSchedules"`
	TotalRecords                    int      `json:"totalRecords"`
	ScheduledMaintenancePercentage  float64  `json:"scheduledMaintenancePercentage"`
	CompletedMaintenancePercentage  float64  `json:"completedMaintenancePercentage"`
	CancelledMaintenancePercentage  float64  `json:"cancelledMaintenancePercentage"`
	PreventiveMaintenancePercentage float64  `json:"preventiveMaintenancePercentage"`
	CorrectiveMaintenancePercentage float64  `json:"correctiveMaintenancePercentage"`
	AverageMaintenancePerAsset      float64  `json:"averageMaintenancePerAsset"`
	AssetsWithMaintenance           int      `json:"assetsWithMaintenance"`
	AssetsWithoutMaintenance        int      `json:"assetsWithoutMaintenance"`
	MaintenanceComplianceRate       float64  `json:"maintenanceComplianceRate"`
	AverageMaintenanceFrequency     float64  `json:"averageMaintenanceFrequency"`
	UpcomingMaintenanceCount        int      `json:"upcomingMaintenanceCount"`
	OverdueMaintenanceCount         int      `json:"overdueMaintenanceCount"`
	RecordsWithCostInfo             int      `json:"recordsWithCostInfo"`
	CostInfoPercentage              float64  `json:"costInfoPercentage"`
	TotalUniqueVendors              int      `json:"totalUniqueVendors"`
	TotalUniquePerformers           int      `json:"totalUniquePerformers"`
	AverageRecordsPerDay            float64  `json:"averageRecordsPerDay"`
	LatestRecordDate                string   `json:"latestRecordDate"`
	EarliestRecordDate              string   `json:"earliestRecordDate"`
	MostExpensiveMaintenanceCost    *float64 `json:"mostExpensiveMaintenanceCost"`
	LeastExpensiveMaintenanceCost   *float64 `json:"leastExpensiveMaintenanceCost"`
}
