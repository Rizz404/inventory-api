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

type MaintenanceScheduleSortField string

const (
	MaintenanceScheduleSortByScheduledDate   MaintenanceScheduleSortField = "scheduledDate"
	MaintenanceScheduleSortByMaintenanceType MaintenanceScheduleSortField = "maintenanceType"
	MaintenanceScheduleSortByStatus          MaintenanceScheduleSortField = "status"
	MaintenanceScheduleSortByCreatedAt       MaintenanceScheduleSortField = "createdAt"
	MaintenanceScheduleSortByUpdatedAt       MaintenanceScheduleSortField = "updatedAt"
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

type MaintenanceScheduleResponse struct {
	ID              string                                   `json:"id"`
	AssetID         string                                   `json:"assetId"`
	MaintenanceType MaintenanceScheduleType                  `json:"maintenanceType"`
	ScheduledDate   time.Time                                `json:"scheduledDate"`
	FrequencyMonths *int                                     `json:"frequencyMonths"`
	Status          ScheduleStatus                           `json:"status"`
	CreatedByID     string                                   `json:"createdById"`
	CreatedAt       time.Time                                `json:"createdAt"`
	Title           string                                   `json:"title"`
	Description     *string                                  `json:"description"`
	Translations    []MaintenanceScheduleTranslationResponse `json:"translations"`
	// * Populated
	Asset     AssetResponse `json:"asset"`
	CreatedBy UserResponse  `json:"createdBy"`
}

type MaintenanceScheduleListResponse struct {
	ID              string                  `json:"id"`
	AssetID         string                  `json:"assetId"`
	MaintenanceType MaintenanceScheduleType `json:"maintenanceType"`
	ScheduledDate   time.Time               `json:"scheduledDate"`
	FrequencyMonths *int                    `json:"frequencyMonths"`
	Status          ScheduleStatus          `json:"status"`
	CreatedByID     string                  `json:"createdById"`
	CreatedAt       time.Time               `json:"createdAt"`
	Title           string                  `json:"title"`
	Description     *string                 `json:"description"`
	// * Populated
	Asset     AssetResponse `json:"asset"`
	CreatedBy UserResponse  `json:"createdBy"`
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

// --- Query Parameters ---

type MaintenanceScheduleFilterOptions struct {
	AssetID         *string                  `json:"assetId,omitempty"`
	MaintenanceType *MaintenanceScheduleType `json:"maintenanceType,omitempty"`
	Status          *ScheduleStatus          `json:"status,omitempty"`
	CreatedBy       *string                  `json:"createdBy,omitempty"`
	FromDate        *string                  `json:"fromDate,omitempty"` // YYYY-MM-DD
	ToDate          *string                  `json:"toDate,omitempty"`   // YYYY-MM-DD
}

type MaintenanceScheduleSortOptions struct {
	Field MaintenanceScheduleSortField `json:"field" example:"createdAt"`
	Order SortOrder                    `json:"order" example:"desc"`
}

type MaintenanceScheduleParams struct {
	SearchQuery *string                           `json:"searchQuery,omitempty"`
	Filters     *MaintenanceScheduleFilterOptions `json:"filters,omitempty"`
	Sort        *MaintenanceScheduleSortOptions   `json:"sort,omitempty"`
	Pagination  *PaginationOptions                `json:"pagination,omitempty"`
}

// --- Statistics ---

// Internal statistics structs (used in repository layer)
type MaintenanceScheduleStatistics struct {
	Total            MaintenanceScheduleCountStatistics   `json:"total"`
	ByType           MaintenanceTypeStatistics            `json:"byType"`
	ByStatus         MaintenanceScheduleStatusStatistics  `json:"byStatus"`
	ByAsset          []AssetMaintenanceScheduleStatistics `json:"byAsset"`
	ByCreator        []UserMaintenanceScheduleStatistics  `json:"byCreator"`
	UpcomingSchedule []UpcomingMaintenanceSchedule        `json:"upcomingSchedule"`
	OverdueSchedule  []OverdueMaintenanceSchedule         `json:"overdueSchedule"`
	FrequencyTrends  []MaintenanceFrequencyTrend          `json:"frequencyTrends"`
	Summary          MaintenanceScheduleSummaryStatistics `json:"summary"`
}

type MaintenanceScheduleCountStatistics struct {
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

type AssetMaintenanceScheduleStatistics struct {
	AssetID         string `json:"assetId"`
	AssetName       string `json:"assetName"`
	AssetTag        string `json:"assetTag"`
	ScheduleCount   int    `json:"scheduleCount"`
	NextMaintenance string `json:"nextMaintenance"`
}

type UserMaintenanceScheduleStatistics struct {
	UserID    string `json:"userId"`
	UserName  string `json:"userName"`
	UserEmail string `json:"userEmail"`
	Count     int    `json:"count"`
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

type MaintenanceScheduleSummaryStatistics struct {
	TotalSchedules                    int       `json:"totalSchedules"`
	ScheduledMaintenancePercentage    float64   `json:"scheduledMaintenancePercentage"`
	CompletedMaintenancePercentage    float64   `json:"completedMaintenancePercentage"`
	CancelledMaintenancePercentage    float64   `json:"cancelledMaintenancePercentage"`
	PreventiveMaintenancePercentage   float64   `json:"preventiveMaintenancePercentage"`
	CorrectiveMaintenancePercentage   float64   `json:"correctiveMaintenancePercentage"`
	AverageScheduleFrequency          float64   `json:"averageScheduleFrequency"`
	UpcomingMaintenanceCount          int       `json:"upcomingMaintenanceCount"`
	OverdueMaintenanceCount           int       `json:"overdueMaintenanceCount"`
	AssetsWithScheduledMaintenance    int       `json:"assetsWithScheduledMaintenance"`
	AssetsWithoutScheduledMaintenance int       `json:"assetsWithoutScheduledMaintenance"`
	AverageSchedulesPerDay            float64   `json:"averageSchedulesPerDay"`
	LatestScheduleDate                time.Time `json:"latestScheduleDate"`
	EarliestScheduleDate              time.Time `json:"earliestScheduleDate"`
	TotalUniqueCreators               int       `json:"totalUniqueCreators"`
}

// Response statistics structs (used in service/handler layer)
type MaintenanceScheduleStatisticsResponse struct {
	Total            MaintenanceScheduleCountStatisticsResponse   `json:"total"`
	ByType           MaintenanceTypeStatisticsResponse            `json:"byType"`
	ByStatus         MaintenanceScheduleStatusStatisticsResponse  `json:"byStatus"`
	ByAsset          []AssetMaintenanceScheduleStatisticsResponse `json:"byAsset"`
	ByCreator        []UserMaintenanceScheduleStatisticsResponse  `json:"byCreator"`
	UpcomingSchedule []UpcomingMaintenanceScheduleResponse        `json:"upcomingSchedule"`
	OverdueSchedule  []OverdueMaintenanceScheduleResponse         `json:"overdueSchedule"`
	FrequencyTrends  []MaintenanceFrequencyTrendResponse          `json:"frequencyTrends"`
	Summary          MaintenanceScheduleSummaryStatisticsResponse `json:"summary"`
}

type MaintenanceScheduleCountStatisticsResponse struct {
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

type AssetMaintenanceScheduleStatisticsResponse struct {
	AssetID         string `json:"assetId"`
	AssetName       string `json:"assetName"`
	AssetTag        string `json:"assetTag"`
	ScheduleCount   int    `json:"scheduleCount"`
	NextMaintenance string `json:"nextMaintenance"`
}

type UserMaintenanceScheduleStatisticsResponse struct {
	UserID    string `json:"userId"`
	UserName  string `json:"userName"`
	UserEmail string `json:"userEmail"`
	Count     int    `json:"count"`
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

type MaintenanceScheduleSummaryStatisticsResponse struct {
	TotalSchedules                    int       `json:"totalSchedules"`
	ScheduledMaintenancePercentage    float64   `json:"scheduledMaintenancePercentage"`
	CompletedMaintenancePercentage    float64   `json:"completedMaintenancePercentage"`
	CancelledMaintenancePercentage    float64   `json:"cancelledMaintenancePercentage"`
	PreventiveMaintenancePercentage   float64   `json:"preventiveMaintenancePercentage"`
	CorrectiveMaintenancePercentage   float64   `json:"correctiveMaintenancePercentage"`
	AverageScheduleFrequency          float64   `json:"averageScheduleFrequency"`
	UpcomingMaintenanceCount          int       `json:"upcomingMaintenanceCount"`
	OverdueMaintenanceCount           int       `json:"overdueMaintenanceCount"`
	AssetsWithScheduledMaintenance    int       `json:"assetsWithScheduledMaintenance"`
	AssetsWithoutScheduledMaintenance int       `json:"assetsWithoutScheduledMaintenance"`
	AverageSchedulesPerDay            float64   `json:"averageSchedulesPerDay"`
	LatestScheduleDate                time.Time `json:"latestScheduleDate"`
	EarliestScheduleDate              time.Time `json:"earliestScheduleDate"`
	TotalUniqueCreators               int       `json:"totalUniqueCreators"`
}
