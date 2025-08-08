package model

import (
	"time"

	"github.com/Rizz404/inventory-api/domain"
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
	Asset           *Asset                `gorm:"foreignKey:AssetID"`
	ScannedByUser   User                  `gorm:"foreignKey:ScannedBy"`
}

func (ScanLog) TableName() string {
	return "scan_logs"
}
