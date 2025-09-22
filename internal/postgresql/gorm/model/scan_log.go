package model

import (
	"log"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type ScanLog struct {
	ID              SQLULID               `gorm:"primaryKey;type:varchar(26)"`
	AssetID         *SQLULID              `gorm:"type:varchar(26)"`
	ScannedValue    string                `gorm:"type:varchar(255);not null"`
	ScanMethod      domain.ScanMethodType `gorm:"type:scan_method_type;not null"`
	ScannedBy       SQLULID               `gorm:"type:varchar(26);not null"`
	ScanTimestamp   time.Time             `gorm:"default:CURRENT_TIMESTAMP"`
	ScanLocationLat *float64              `gorm:"type:decimal(11,8)"`
	ScanLocationLng *float64              `gorm:"type:decimal(11,8)"`
	ScanResult      domain.ScanResultType `gorm:"type:scan_result_type;not null"`
}

func (ScanLog) TableName() string {
	return "scan_logs"
}

func (u *ScanLog) BeforeCreate(tx *gorm.DB) error {
	log.Printf("🚀 ScanLog.BeforeCreate called! Current ID: %s, IsZero: %t", u.ID.String(), u.ID.IsZero())

	if u.ID.IsZero() {
		u.ID = SQLULID(ulid.Make())
		log.Printf("🚀 Generated new ULID for ScanLog: %s", u.ID.String())
	}

	return nil
}
