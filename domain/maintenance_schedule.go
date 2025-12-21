package domain

import "time"

// --- Enums ---

type MaintenanceScheduleType string

const (
	ScheduleTypePreventive  MaintenanceScheduleType = "Preventive"
	ScheduleTypeCorrective  MaintenanceScheduleType = "Corrective"
	ScheduleTypeInspection  MaintenanceScheduleType = "Inspection"
	ScheduleTypeCalibration MaintenanceScheduleType = "Calibration"
)

type ScheduleState string

const (
	StateActive    ScheduleState = "Active"
	StatePaused    ScheduleState = "Paused"
	StateStopped   ScheduleState = "Stopped"
	StateCompleted ScheduleState = "Completed"
)

type IntervalUnit string

const (
	IntervalMinutes IntervalUnit = "Minutes"
	IntervalHours   IntervalUnit = "Hours"
	IntervalDays    IntervalUnit = "Days"
	IntervalWeeks   IntervalUnit = "Weeks"
	IntervalMonths  IntervalUnit = "Months"
	IntervalYears   IntervalUnit = "Years"
)

type MaintenanceScheduleSortField string

const (
	MaintenanceScheduleSortByNextScheduledDate MaintenanceScheduleSortField = "nextScheduledDate"
	MaintenanceScheduleSortByMaintenanceType   MaintenanceScheduleSortField = "maintenanceType"
	MaintenanceScheduleSortByState             MaintenanceScheduleSortField = "state"
	MaintenanceScheduleSortByCreatedAt         MaintenanceScheduleSortField = "createdAt"
	MaintenanceScheduleSortByUpdatedAt         MaintenanceScheduleSortField = "updatedAt"
)

// --- Structs ---

type MaintenanceSchedule struct {
	ID                string                           `json:"id"`
	AssetID           string                           `json:"assetId"`
	MaintenanceType   MaintenanceScheduleType          `json:"maintenanceType"`
	IsRecurring       bool                             `json:"isRecurring"`
	IntervalValue     *int                             `json:"intervalValue"`
	IntervalUnit      *IntervalUnit                    `json:"intervalUnit"`
	ScheduledTime     *string                          `json:"scheduledTime"`
	NextScheduledDate time.Time                        `json:"nextScheduledDate"`
	LastExecutedDate  *time.Time                       `json:"lastExecutedDate"`
	State             ScheduleState                    `json:"state"`
	AutoComplete      bool                             `json:"autoComplete"`
	EstimatedCost     *float64                         `json:"estimatedCost"`
	CreatedBy         string                           `json:"createdBy"`
	CreatedAt         time.Time                        `json:"createdAt"`
	UpdatedAt         time.Time                        `json:"updatedAt"`
	Translations      []MaintenanceScheduleTranslation `json:"translations,omitempty"`
	// * Populated
	Asset         *Asset `json:"asset,omitempty"`
	CreatedByUser *User  `json:"createdByUser,omitempty"`
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
	ID                string                                   `json:"id"`
	AssetID           string                                   `json:"assetId"`
	MaintenanceType   MaintenanceScheduleType                  `json:"maintenanceType"`
	IsRecurring       bool                                     `json:"isRecurring"`
	IntervalValue     *int                                     `json:"intervalValue"`
	IntervalUnit      *IntervalUnit                            `json:"intervalUnit"`
	ScheduledTime     *string                                  `json:"scheduledTime"`
	NextScheduledDate time.Time                                `json:"nextScheduledDate"`
	LastExecutedDate  *time.Time                               `json:"lastExecutedDate"`
	State             ScheduleState                            `json:"state"`
	AutoComplete      bool                                     `json:"autoComplete"`
	EstimatedCost     *NullableDecimal2                        `json:"estimatedCost"`
	CreatedByID       string                                   `json:"createdById"`
	CreatedAt         time.Time                                `json:"createdAt"`
	UpdatedAt         time.Time                                `json:"updatedAt"`
	Title             string                                   `json:"title"`
	Description       *string                                  `json:"description"`
	Translations      []MaintenanceScheduleTranslationResponse `json:"translations"`
	// * Populated
	Asset     AssetResponse `json:"asset"`
	CreatedBy UserResponse  `json:"createdBy"`
}

type MaintenanceScheduleListResponse struct {
	ID                string                  `json:"id"`
	AssetID           string                  `json:"assetId"`
	MaintenanceType   MaintenanceScheduleType `json:"maintenanceType"`
	IsRecurring       bool                    `json:"isRecurring"`
	IntervalValue     *int                    `json:"intervalValue"`
	IntervalUnit      *IntervalUnit           `json:"intervalUnit"`
	ScheduledTime     *string                 `json:"scheduledTime"`
	NextScheduledDate time.Time               `json:"nextScheduledDate"`
	LastExecutedDate  *time.Time              `json:"lastExecutedDate"`
	State             ScheduleState           `json:"state"`
	AutoComplete      bool                    `json:"autoComplete"`
	EstimatedCost     *NullableDecimal2       `json:"estimatedCost"`
	CreatedByID       string                  `json:"createdById"`
	CreatedAt         time.Time               `json:"createdAt"`
	UpdatedAt         time.Time               `json:"updatedAt"`
	Title             string                  `json:"title"`
	Description       *string                 `json:"description"`
	// * Populated
	Asset     AssetResponse `json:"asset"`
	CreatedBy UserResponse  `json:"createdBy"`
}

type BulkDeleteMaintenanceSchedules struct {
	RequestedIDS []string `json:"requestedIds"`
	DeletedIDS   []string `json:"deletedIds"`
}

type BulkDeleteMaintenanceSchedulesResponse struct {
	RequestedIDS []string `json:"requestedIds"`
	DeletedIDS   []string `json:"deletedIds"`
}

// --- Bulk Create ---

type BulkCreateMaintenanceSchedulesPayload struct {
	MaintenanceSchedules []CreateMaintenanceSchedulePayload `json:"maintenanceSchedules" validate:"required,min=1,max=100,dive"`
}

type BulkCreateMaintenanceSchedulesResponse struct {
	MaintenanceSchedules []MaintenanceScheduleResponse `json:"maintenanceSchedules"`
}

// --- Payloads ---

type CreateMaintenanceSchedulePayload struct {
	AssetID           string                                        `json:"assetId" validate:"required"`
	MaintenanceType   MaintenanceScheduleType                       `json:"maintenanceType" validate:"required,oneof=Preventive Corrective Inspection Calibration"`
	IsRecurring       *bool                                         `json:"isRecurring,omitempty"`
	IntervalValue     *int                                          `json:"intervalValue,omitempty" validate:"omitempty,gt=0"`
	IntervalUnit      *IntervalUnit                                 `json:"intervalUnit,omitempty" validate:"omitempty,oneof=Minutes Hours Days Weeks Months Years"`
	ScheduledTime     *string                                       `json:"scheduledTime,omitempty" validate:"omitempty"`
	NextScheduledDate string                                        `json:"nextScheduledDate" validate:"required,datetime=2006-01-02"`
	AutoComplete      *bool                                         `json:"autoComplete,omitempty"`
	EstimatedCost     *float64                                      `json:"estimatedCost,omitempty" validate:"omitempty,gt=0"`
	Translations      []CreateMaintenanceScheduleTranslationPayload `json:"translations" validate:"required,min=1,dive"`
}

type CreateMaintenanceScheduleTranslationPayload struct {
	LangCode    string  `json:"langCode" validate:"required,max=5"`
	Title       string  `json:"title" validate:"required,max=200"`
	Description *string `json:"description,omitempty"`
}

type UpdateMaintenanceSchedulePayload struct {
	MaintenanceType   *MaintenanceScheduleType                      `json:"maintenanceType,omitempty" validate:"omitempty,oneof=Preventive Corrective Inspection Calibration"`
	IsRecurring       *bool                                         `json:"isRecurring,omitempty"`
	IntervalValue     *int                                          `json:"intervalValue,omitempty" validate:"omitempty,gt=0"`
	IntervalUnit      *IntervalUnit                                 `json:"intervalUnit,omitempty" validate:"omitempty,oneof=Minutes Hours Days Weeks Months Years"`
	ScheduledTime     *string                                       `json:"scheduledTime,omitempty" validate:"omitempty"`
	NextScheduledDate *string                                       `json:"nextScheduledDate,omitempty" validate:"omitempty,datetime=2006-01-02"`
	State             *ScheduleState                                `json:"state,omitempty" validate:"omitempty,oneof=Active Paused Stopped Completed"`
	AutoComplete      *bool                                         `json:"autoComplete,omitempty"`
	EstimatedCost     *float64                                      `json:"estimatedCost,omitempty" validate:"omitempty,gt=0"`
	Translations      []UpdateMaintenanceScheduleTranslationPayload `json:"translations,omitempty" validate:"omitempty,dive"`
}

type UpdateMaintenanceScheduleTranslationPayload struct {
	LangCode    string  `json:"langCode" validate:"required,max=5"`
	Title       *string `json:"title,omitempty" validate:"omitempty,max=200"`
	Description *string `json:"description,omitempty"`
}

type BulkDeleteMaintenanceSchedulesPayload struct {
	IDS []string `json:"ids" validate:"required,min=1,max=100,dive,required"`
}

type ExportMaintenanceScheduleListPayload struct {
	Format      ExportFormat                      `json:"format" validate:"required,oneof=pdf excel"`
	SearchQuery *string                           `json:"searchQuery,omitempty"`
	Filters     *MaintenanceScheduleFilterOptions `json:"filters,omitempty"`
	Sort        *MaintenanceScheduleSortOptions   `json:"sort,omitempty"`
}

// --- Query Parameters ---

type MaintenanceScheduleFilterOptions struct {
	AssetID         *string                  `json:"assetId,omitempty"`
	MaintenanceType *MaintenanceScheduleType `json:"maintenanceType,omitempty"`
	State           *ScheduleState           `json:"state,omitempty"`
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
	Preventive  int `json:"preventive"`
	Corrective  int `json:"corrective"`
	Inspection  int `json:"inspection"`
	Calibration int `json:"calibration"`
}

type MaintenanceScheduleStatusStatistics struct {
	Active    int `json:"active"`
	Paused    int `json:"paused"`
	Stopped   int `json:"stopped"`
	Completed int `json:"completed"`
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
	ID                string                  `json:"id"`
	AssetID           string                  `json:"assetId"`
	AssetName         string                  `json:"assetName"`
	AssetTag          string                  `json:"assetTag"`
	MaintenanceType   MaintenanceScheduleType `json:"maintenanceType"`
	NextScheduledDate string                  `json:"nextScheduledDate"`
	DaysUntilDue      int                     `json:"daysUntilDue"`
	Title             string                  `json:"title"`
	Description       *string                 `json:"description,omitempty"`
}

type OverdueMaintenanceSchedule struct {
	ID                string                  `json:"id"`
	AssetID           string                  `json:"assetId"`
	AssetName         string                  `json:"assetName"`
	AssetTag          string                  `json:"assetTag"`
	MaintenanceType   MaintenanceScheduleType `json:"maintenanceType"`
	NextScheduledDate string                  `json:"nextScheduledDate"`
	DaysOverdue       int                     `json:"daysOverdue"`
	Title             string                  `json:"title"`
	Description       *string                 `json:"description,omitempty"`
}

type MaintenanceFrequencyTrend struct {
	FrequencyMonths int `json:"frequencyMonths"`
	Count           int `json:"count"`
}

type MaintenanceScheduleSummaryStatistics struct {
	TotalSchedules                    int       `json:"totalSchedules"`
	ActiveMaintenancePercentage       float64   `json:"activeMaintenancePercentage"`
	CompletedMaintenancePercentage    float64   `json:"completedMaintenancePercentage"`
	PausedMaintenancePercentage       float64   `json:"pausedMaintenancePercentage"`
	StoppedMaintenancePercentage      float64   `json:"stoppedMaintenancePercentage"`
	PreventiveMaintenancePercentage   float64   `json:"preventiveMaintenancePercentage"`
	CorrectiveMaintenancePercentage   float64   `json:"correctiveMaintenancePercentage"`
	InspectionMaintenancePercentage   float64   `json:"inspectionMaintenancePercentage"`
	CalibrationMaintenancePercentage  float64   `json:"calibrationMaintenancePercentage"`
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
	Preventive  int `json:"preventive"`
	Corrective  int `json:"corrective"`
	Inspection  int `json:"inspection"`
	Calibration int `json:"calibration"`
}

type MaintenanceScheduleStatusStatisticsResponse struct {
	Active    int `json:"active"`
	Paused    int `json:"paused"`
	Stopped   int `json:"stopped"`
	Completed int `json:"completed"`
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
	ID                string                  `json:"id"`
	AssetID           string                  `json:"assetId"`
	AssetName         string                  `json:"assetName"`
	AssetTag          string                  `json:"assetTag"`
	MaintenanceType   MaintenanceScheduleType `json:"maintenanceType"`
	NextScheduledDate string                  `json:"nextScheduledDate"`
	DaysUntilDue      int                     `json:"daysUntilDue"`
	Title             string                  `json:"title"`
	Description       *string                 `json:"description,omitempty"`
}

type OverdueMaintenanceScheduleResponse struct {
	ID                string                  `json:"id"`
	AssetID           string                  `json:"assetId"`
	AssetName         string                  `json:"assetName"`
	AssetTag          string                  `json:"assetTag"`
	MaintenanceType   MaintenanceScheduleType `json:"maintenanceType"`
	NextScheduledDate string                  `json:"nextScheduledDate"`
	DaysOverdue       int                     `json:"daysOverdue"`
	Title             string                  `json:"title"`
	Description       *string                 `json:"description,omitempty"`
}

type MaintenanceFrequencyTrendResponse struct {
	FrequencyMonths int `json:"frequencyMonths"`
	Count           int `json:"count"`
}

type MaintenanceScheduleSummaryStatisticsResponse struct {
	TotalSchedules                    int       `json:"totalSchedules"`
	ActiveMaintenancePercentage       Decimal2  `json:"activeMaintenancePercentage"`
	CompletedMaintenancePercentage    Decimal2  `json:"completedMaintenancePercentage"`
	PausedMaintenancePercentage       Decimal2  `json:"pausedMaintenancePercentage"`
	StoppedMaintenancePercentage      Decimal2  `json:"stoppedMaintenancePercentage"`
	PreventiveMaintenancePercentage   Decimal2  `json:"preventiveMaintenancePercentage"`
	CorrectiveMaintenancePercentage   Decimal2  `json:"correctiveMaintenancePercentage"`
	InspectionMaintenancePercentage   Decimal2  `json:"inspectionMaintenancePercentage"`
	CalibrationMaintenancePercentage  Decimal2  `json:"calibrationMaintenancePercentage"`
	AverageScheduleFrequency          Decimal2  `json:"averageScheduleFrequency"`
	UpcomingMaintenanceCount          int       `json:"upcomingMaintenanceCount"`
	OverdueMaintenanceCount           int       `json:"overdueMaintenanceCount"`
	AssetsWithScheduledMaintenance    int       `json:"assetsWithScheduledMaintenance"`
	AssetsWithoutScheduledMaintenance int       `json:"assetsWithoutScheduledMaintenance"`
	AverageSchedulesPerDay            Decimal2  `json:"averageSchedulesPerDay"`
	LatestScheduleDate                time.Time `json:"latestScheduleDate"`
	EarliestScheduleDate              time.Time `json:"earliestScheduleDate"`
	TotalUniqueCreators               int       `json:"totalUniqueCreators"`
}
