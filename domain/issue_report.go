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

type IssueReportResponse struct {
	ID              string        `json:"id"`
	Asset           AssetResponse `json:"asset"`
	ReportedBy      UserResponse  `json:"reportedBy"`
	ReportedDate    string        `json:"reportedDate"`
	IssueType       string        `json:"issueType"`
	Priority        IssuePriority `json:"priority"`
	Status          IssueStatus   `json:"status"`
	ResolvedDate    *string       `json:"resolvedDate,omitempty"`
	ResolvedBy      *UserResponse `json:"resolvedBy,omitempty"`
	Title           string        `json:"title"`
	Description     *string       `json:"description,omitempty"`
	ResolutionNotes *string       `json:"resolutionNotes,omitempty"`
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
