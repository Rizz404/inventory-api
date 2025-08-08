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

type MaintenanceRecord struct {
	ID                string                         `json:"id"`
	ScheduleID        *string                        `json:"scheduleId"`
	AssetID           string                         `json:"assetId"`
	MaintenanceDate   time.Time                      `json:"maintenanceDate"`
	PerformedByUser   *string                        `json:"performedByUser"`
	PerformedByVendor *string                        `json:"performedByVendor"`
	ActualCost        *float64                       `json:"actualCost"`
	Translations      []MaintenanceRecordTranslation `json:"translations,omitempty"`
}

type MaintenanceRecordTranslation struct {
	ID       string  `json:"id"`
	RecordID string  `json:"recordId"`
	LangCode string  `json:"langCode"`
	Title    string  `json:"title"`
	Notes    *string `json:"notes"`
}

type MaintenanceScheduleResponse struct {
	ID              string                  `json:"id"`
	Asset           AssetResponse           `json:"asset"`
	MaintenanceType MaintenanceScheduleType `json:"maintenanceType"`
	ScheduledDate   string                  `json:"scheduledDate"`
	FrequencyMonths *int                    `json:"frequencyMonths,omitempty"`
	Status          ScheduleStatus          `json:"status"`
	CreatedBy       UserResponse            `json:"createdBy"`
	CreatedAt       string                  `json:"createdAt"`
	Title           string                  `json:"title"`
	Description     *string                 `json:"description,omitempty"`
}

type MaintenanceRecordResponse struct {
	ID                string                       `json:"id"`
	Schedule          *MaintenanceScheduleResponse `json:"schedule,omitempty"`
	Asset             AssetResponse                `json:"asset"`
	MaintenanceDate   string                       `json:"maintenanceDate"`
	PerformedByUser   *UserResponse                `json:"performedByUser,omitempty"`
	PerformedByVendor *string                      `json:"performedByVendor,omitempty"`
	ActualCost        *float64                     `json:"actualCost,omitempty"`
	Title             string                       `json:"title"`
	Notes             *string                      `json:"notes,omitempty"`
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
