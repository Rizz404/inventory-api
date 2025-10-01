package domain

import "time"

// --- Enums ---

type IssuePriority string

const (
	PriorityLow      IssuePriority = "Low"
	PriorityMedium   IssuePriority = "Medium"
	PriorityHigh     IssuePriority = "High"
	PriorityCritical IssuePriority = "Critical"
)

type IssueStatus string

const (
	IssueStatusOpen       IssueStatus = "Open"
	IssueStatusInProgress IssueStatus = "In Progress"
	IssueStatusResolved   IssueStatus = "Resolved"
	IssueStatusClosed     IssueStatus = "Closed"
)

// --- Structs ---

type IssueReport struct {
	ID           string                   `json:"id"`
	AssetID      string                   `json:"assetId"`
	ReportedBy   string                   `json:"reportedBy"`
	ReportedDate time.Time                `json:"reportedDate"`
	IssueType    string                   `json:"issueType"`
	Priority     IssuePriority            `json:"priority"`
	Status       IssueStatus              `json:"status"`
	ResolvedDate *time.Time               `json:"resolvedDate"`
	ResolvedBy   *string                  `json:"resolvedBy"`
	Translations []IssueReportTranslation `json:"translations,omitempty"`
}

type IssueReportTranslation struct {
	ID              string  `json:"id"`
	ReportID        string  `json:"reportId"`
	LangCode        string  `json:"langCode"`
	Title           string  `json:"title"`
	Description     *string `json:"description"`
	ResolutionNotes *string `json:"resolutionNotes"`
}

type IssueReportTranslationResponse struct {
	LangCode        string  `json:"langCode"`
	Title           string  `json:"title"`
	Description     *string `json:"description"`
	ResolutionNotes *string `json:"resolutionNotes"`
}

type IssueReportResponse struct {
	ID              string                           `json:"id"`
	AssetID         string                           `json:"assetId"`
	ReportedByID    string                           `json:"reportedById"`
	ReportedDate    time.Time                        `json:"reportedDate"`
	IssueType       string                           `json:"issueType"`
	Priority        IssuePriority                    `json:"priority"`
	Status          IssueStatus                      `json:"status"`
	ResolvedDate    *time.Time                       `json:"resolvedDate,omitempty"`
	ResolvedByID    *string                          `json:"resolvedById,omitempty"`
	Title           string                           `json:"title"`
	Description     *string                          `json:"description,omitempty"`
	ResolutionNotes *string                          `json:"resolutionNotes,omitempty"`
	CreatedAt       time.Time                        `json:"createdAt"`
	UpdatedAt       time.Time                        `json:"updatedAt"`
	Translations    []IssueReportTranslationResponse `json:"translations"`
	// * Populated
	Asset      AssetResponse `json:"asset"`
	ReportedBy UserResponse  `json:"reportedBy"`
	ResolvedBy *UserResponse `json:"resolvedBy,omitempty"`
}

type IssueReportListResponse struct {
	ID              string        `json:"id"`
	AssetID         string        `json:"assetId"`
	ReportedByID    string        `json:"reportedById"`
	ReportedDate    time.Time     `json:"reportedDate"`
	IssueType       string        `json:"issueType"`
	Priority        IssuePriority `json:"priority"`
	Status          IssueStatus   `json:"status"`
	ResolvedDate    *time.Time    `json:"resolvedDate,omitempty"`
	ResolvedByID    *string       `json:"resolvedById,omitempty"`
	Title           string        `json:"title"`
	Description     *string       `json:"description,omitempty"`
	ResolutionNotes *string       `json:"resolutionNotes,omitempty"`
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`
	// * Populated
	Asset      AssetResponse `json:"asset"`
	ReportedBy UserResponse  `json:"reportedBy"`
	ResolvedBy *UserResponse `json:"resolvedBy,omitempty"`
}

// --- Payloads ---

type CreateIssueReportPayload struct {
	AssetID      string                                `json:"assetId" validate:"required"`
	IssueType    string                                `json:"issueType" validate:"required,max=50"`
	Priority     IssuePriority                         `json:"priority" validate:"required,oneof=Low Medium High Critical"`
	Translations []CreateIssueReportTranslationPayload `json:"translations" validate:"required,min=1,dive"`
}

type CreateIssueReportTranslationPayload struct {
	LangCode    string  `json:"langCode" validate:"required,max=5"`
	Title       string  `json:"title" validate:"required,max=200"`
	Description *string `json:"description,omitempty"`
}

type UpdateIssueReportPayload struct {
	Priority        *IssuePriority `json:"priority,omitempty" validate:"omitempty,oneof=Low Medium High Critical"`
	Status          *IssueStatus   `json:"status,omitempty" validate:"omitempty,oneof=Open 'In Progress' Resolved Closed"`
	ResolutionNotes *string        `json:"resolutionNotes,omitempty"` // This should be updated via a specific action
}

type ResolveIssueReportPayload struct {
	ResolutionNotes string `json:"resolutionNotes" validate:"required,max=1000"`
}

type UpdateIssueReportTranslationPayload struct {
	LangCode        string  `json:"langCode" validate:"required,max=5"`
	Title           *string `json:"title,omitempty" validate:"omitempty,max=200"`
	Description     *string `json:"description,omitempty"`
	ResolutionNotes *string `json:"resolutionNotes,omitempty"`
}

// --- Query Parameters ---

type IssueReportFilterOptions struct {
	AssetID    *string        `json:"assetId,omitempty"`
	ReportedBy *string        `json:"reportedBy,omitempty"`
	ResolvedBy *string        `json:"resolvedBy,omitempty"`
	IssueType  *string        `json:"issueType,omitempty"`
	Priority   *IssuePriority `json:"priority,omitempty"`
	Status     *IssueStatus   `json:"status,omitempty"`
	IsResolved *bool          `json:"isResolved,omitempty"`
	DateFrom   *time.Time     `json:"dateFrom,omitempty"`
	DateTo     *time.Time     `json:"dateTo,omitempty"`
}

type IssueReportSortOptions struct {
	Field string `json:"field,omitempty"`
	Order string `json:"order,omitempty"`
}

type IssueReportPaginationOptions struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

type IssueReportParams struct {
	SearchQuery *string                       `json:"searchQuery,omitempty"`
	Filters     *IssueReportFilterOptions     `json:"filters,omitempty"`
	Sort        *IssueReportSortOptions       `json:"sort,omitempty"`
	Pagination  *IssueReportPaginationOptions `json:"pagination,omitempty"`
}

// --- Statistics ---

// Internal statistics structs (used in repository layer)
type IssueReportStatistics struct {
	Total          IssueReportCountStatistics    `json:"total"`
	ByPriority     IssueReportPriorityStatistics `json:"byPriority"`
	ByStatus       IssueReportStatusStatistics   `json:"byStatus"`
	ByType         IssueReportTypeStatistics     `json:"byType"`
	CreationTrends []IssueReportCreationTrend    `json:"creationTrends"`
	Summary        IssueReportSummaryStatistics  `json:"summary"`
}

type IssueReportCountStatistics struct {
	Count int `json:"count"`
}

type IssueReportPriorityStatistics struct {
	Low      int `json:"low"`
	Medium   int `json:"medium"`
	High     int `json:"high"`
	Critical int `json:"critical"`
}

type IssueReportStatusStatistics struct {
	Open       int `json:"open"`
	InProgress int `json:"inProgress"`
	Resolved   int `json:"resolved"`
	Closed     int `json:"closed"`
}

type IssueReportTypeStatistics struct {
	Types map[string]int `json:"types"`
}

type IssueReportCreationTrend struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

type IssueReportSummaryStatistics struct {
	TotalReports            int       `json:"totalReports"`
	OpenPercentage          float64   `json:"openPercentage"`
	ResolvedPercentage      float64   `json:"resolvedPercentage"`
	AverageResolutionTime   float64   `json:"averageResolutionTimeInDays"`
	MostCommonPriority      string    `json:"mostCommonPriority"`
	MostCommonType          string    `json:"mostCommonType"`
	CriticalUnresolvedCount int       `json:"criticalUnresolvedCount"`
	AverageReportsPerDay    float64   `json:"averageReportsPerDay"`
	LatestCreationDate      time.Time `json:"latestCreationDate"`
	EarliestCreationDate    time.Time `json:"earliestCreationDate"`
}

// Response statistics structs (used in service/handler layer)
type IssueReportStatisticsResponse struct {
	Total          IssueReportCountStatisticsResponse    `json:"total"`
	ByPriority     IssueReportPriorityStatisticsResponse `json:"byPriority"`
	ByStatus       IssueReportStatusStatisticsResponse   `json:"byStatus"`
	ByType         IssueReportTypeStatisticsResponse     `json:"byType"`
	CreationTrends []IssueReportCreationTrendResponse    `json:"creationTrends"`
	Summary        IssueReportSummaryStatisticsResponse  `json:"summary"`
}

type IssueReportCountStatisticsResponse struct {
	Count int `json:"count"`
}

type IssueReportPriorityStatisticsResponse struct {
	Low      int `json:"low"`
	Medium   int `json:"medium"`
	High     int `json:"high"`
	Critical int `json:"critical"`
}

type IssueReportStatusStatisticsResponse struct {
	Open       int `json:"open"`
	InProgress int `json:"inProgress"`
	Resolved   int `json:"resolved"`
	Closed     int `json:"closed"`
}

type IssueReportTypeStatisticsResponse struct {
	Types map[string]int `json:"types"`
}

type IssueReportCreationTrendResponse struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

type IssueReportSummaryStatisticsResponse struct {
	TotalReports            int       `json:"totalReports"`
	OpenPercentage          float64   `json:"openPercentage"`
	ResolvedPercentage      float64   `json:"resolvedPercentage"`
	AverageResolutionTime   float64   `json:"averageResolutionTimeInDays"`
	MostCommonPriority      string    `json:"mostCommonPriority"`
	MostCommonType          string    `json:"mostCommonType"`
	CriticalUnresolvedCount int       `json:"criticalUnresolvedCount"`
	AverageReportsPerDay    float64   `json:"averageReportsPerDay"`
	LatestCreationDate      time.Time `json:"latestCreationDate"`
	EarliestCreationDate    time.Time `json:"earliestCreationDate"`
}
