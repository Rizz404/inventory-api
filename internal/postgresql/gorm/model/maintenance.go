package model

import (
	"time"

	"github.com/Rizz404/inventory-api/domain"
)

type MaintenanceSchedule struct {
	ID              SQLULID                        `gorm:"primaryKey;type:varchar(26)"`
	AssetID         SQLULID                        `gorm:"type:varchar(26);not null"`
	MaintenanceType domain.MaintenanceScheduleType `gorm:"type:maintenance_schedule_type;not null"`
	ScheduledDate   time.Time                      `gorm:"type:date;not null"`
	FrequencyMonths *int
	Status          domain.ScheduleStatus `gorm:"type:schedule_status;default:'Scheduled'"`
	CreatedBy       SQLULID               `gorm:"type:varchar(26);not null"`
	CreatedAt       time.Time
	Asset           Asset                            `gorm:"foreignKey:AssetID"`
	CreatedByUser   User                             `gorm:"foreignKey:CreatedBy"`
	Translations    []MaintenanceScheduleTranslation `gorm:"foreignKey:ScheduleID"`
}

func (MaintenanceSchedule) TableName() string {
	return "maintenance_schedules"
}

type MaintenanceScheduleTranslation struct {
	ID          SQLULID `gorm:"primaryKey;type:varchar(26)"`
	ScheduleID  SQLULID `gorm:"type:varchar(26);not null;uniqueIndex:idx_sch_lang"`
	LangCode    string  `gorm:"type:varchar(5);not null;uniqueIndex:idx_sch_lang"`
	Title       string  `gorm:"type:varchar(200);not null"`
	Description *string `gorm:"type:text"`
}

func (MaintenanceScheduleTranslation) TableName() string {
	return "maintenance_schedules_translation"
}

type MaintenanceRecord struct {
	ID                SQLULID   `gorm:"primaryKey;type:varchar(26)"`
	ScheduleID        *SQLULID  `gorm:"type:varchar(26)"`
	AssetID           SQLULID   `gorm:"type:varchar(26);not null"`
	MaintenanceDate   time.Time `gorm:"type:date;not null"`
	PerformedByUser   *SQLULID  `gorm:"type:varchar(26)"`
	PerformedByVendor *string   `gorm:"type:varchar(150)"`
	ActualCost        *float64  `gorm:"type:decimal(12,2)"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Schedule          *MaintenanceSchedule           `gorm:"foreignKey:ScheduleID"`
	Asset             Asset                          `gorm:"foreignKey:AssetID"`
	User              *User                          `gorm:"foreignKey:PerformedByUser"`
	Translations      []MaintenanceRecordTranslation `gorm:"foreignKey:RecordID"`
}

func (MaintenanceRecord) TableName() string {
	return "maintenance_records"
}

type MaintenanceRecordTranslation struct {
	ID       SQLULID `gorm:"primaryKey;type:varchar(26)"`
	RecordID SQLULID `gorm:"type:varchar(26);not null;uniqueIndex:idx_rec_lang"`
	LangCode string  `gorm:"type:varchar(5);not null;uniqueIndex:idx_rec_lang"`
	Title    string  `gorm:"type:varchar(200);not null"`
	Notes    *string `gorm:"type:text"`
}

func (MaintenanceRecordTranslation) TableName() string {
	return "maintenance_records_translation"
}
