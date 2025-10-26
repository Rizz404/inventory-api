package model

import (
	"log"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type MaintenanceSchedule struct {
	ID                SQLULID                        `gorm:"primaryKey;type:varchar(26)"`
	AssetID           SQLULID                        `gorm:"type:varchar(26);not null"`
	MaintenanceType   domain.MaintenanceScheduleType `gorm:"type:maintenance_type;not null"`
	IsRecurring       bool                           `gorm:"default:false"`
	IntervalValue     *int                           `gorm:"type:int"`
	IntervalUnit      *domain.IntervalUnit           `gorm:"type:interval_unit"`
	ScheduledTime     *string                        `gorm:"type:time"`
	NextScheduledDate time.Time                      `gorm:"type:timestamp with time zone;not null"`
	LastExecutedDate  *time.Time                     `gorm:"type:timestamp with time zone"`
	State             domain.ScheduleState           `gorm:"type:schedule_state;default:'Active'"`
	AutoComplete      bool                           `gorm:"default:false"`
	EstimatedCost     *float64                       `gorm:"type:decimal(12,2)"`
	CreatedBy         SQLULID                        `gorm:"type:varchar(26);not null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Asset             Asset                            `gorm:"foreignKey:AssetID"`
	CreatedByUser     User                             `gorm:"foreignKey:CreatedBy"`
	Translations      []MaintenanceScheduleTranslation `gorm:"foreignKey:ScheduleID"`
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
	return "maintenance_schedule_translations"
}

func (u *MaintenanceScheduleTranslation) BeforeCreate(tx *gorm.DB) error {
	log.Printf("🚀 MaintenanceScheduleTranslation.BeforeCreate called! Current ID: %s, IsZero: %t", u.ID.String(), u.ID.IsZero())

	if u.ID.IsZero() {
		u.ID = SQLULID(ulid.Make())
		log.Printf("🚀 Generated new ULID for MaintenanceScheduleTranslation: %s", u.ID.String())
	}

	return nil
}
