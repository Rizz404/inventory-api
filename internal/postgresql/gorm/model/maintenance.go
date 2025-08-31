package model

import (
	"log"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
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

func (u *MaintenanceSchedule) BeforeCreate(tx *gorm.DB) error {
	log.Printf("🚀 MaintenanceSchedule.BeforeCreate called! Current ID: %s, IsZero: %t", u.ID.String(), u.ID.IsZero())

	if u.ID.IsZero() {
		u.ID = SQLULID(ulid.Make())
		log.Printf("🚀 Generated new ULID for MaintenanceSchedule: %s", u.ID.String())
	}

	return nil
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

func (u *MaintenanceScheduleTranslation) BeforeCreate(tx *gorm.DB) error {
	log.Printf("🚀 MaintenanceScheduleTranslation.BeforeCreate called! Current ID: %s, IsZero: %t", u.ID.String(), u.ID.IsZero())

	if u.ID.IsZero() {
		u.ID = SQLULID(ulid.Make())
		log.Printf("🚀 Generated new ULID for MaintenanceScheduleTranslation: %s", u.ID.String())
	}

	return nil
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

func (u *MaintenanceRecord) BeforeCreate(tx *gorm.DB) error {
	log.Printf("🚀 MaintenanceRecord.BeforeCreate called! Current ID: %s, IsZero: %t", u.ID.String(), u.ID.IsZero())

	if u.ID.IsZero() {
		u.ID = SQLULID(ulid.Make())
		log.Printf("🚀 Generated new ULID for MaintenanceRecord: %s", u.ID.String())
	}

	return nil
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

func (u *MaintenanceRecordTranslation) BeforeCreate(tx *gorm.DB) error {
	log.Printf("🚀 MaintenanceRecordTranslation.BeforeCreate called! Current ID: %s, IsZero: %t", u.ID.String(), u.ID.IsZero())

	if u.ID.IsZero() {
		u.ID = SQLULID(ulid.Make())
		log.Printf("🚀 Generated new ULID for MaintenanceRecordTranslation: %s", u.ID.String())
	}

	return nil
}
