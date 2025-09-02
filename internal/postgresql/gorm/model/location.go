package model

import (
	"log"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Location struct {
	ID           SQLULID  `gorm:"primaryKey;type:varchar(26)"`
	LocationCode string   `gorm:"type:varchar(20);unique;not null"`
	Building     *string  `gorm:"type:varchar(100)"`
	Floor        *string  `gorm:"type:varchar(20)"`
	Latitude     *float64 `gorm:"type:decimal(11,8)"`
	Longitude    *float64 `gorm:"type:decimal(11,8)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Translations []LocationTranslation `gorm:"foreignKey:LocationID"`
}

func (Location) TableName() string {
	return "locations"
}

func (u *Location) BeforeCreate(tx *gorm.DB) error {
	log.Printf("🚀 Location.BeforeCreate called! Current ID: %s, IsZero: %t", u.ID.String(), u.ID.IsZero())

	if u.ID.IsZero() {
		u.ID = SQLULID(ulid.Make())
		log.Printf("🚀 Generated new ULID for Location: %s", u.ID.String())
	}

	return nil
}

type LocationTranslation struct {
	ID           SQLULID `gorm:"primaryKey;type:varchar(26)"`
	LocationID   SQLULID `gorm:"type:varchar(26);not null;uniqueIndex:idx_loc_lang"`
	LangCode     string  `gorm:"type:varchar(5);not null;uniqueIndex:idx_loc_lang"`
	LocationName string  `gorm:"type:varchar(100);not null"`
}

func (u *LocationTranslation) BeforeCreate(tx *gorm.DB) error {
	log.Printf("🚀 LocationTranslation.BeforeCreate called! Current ID: %s, IsZero: %t", u.ID.String(), u.ID.IsZero())

	if u.ID.IsZero() {
		u.ID = SQLULID(ulid.Make())
		log.Printf("🚀 Generated new ULID for LocationTranslation: %s", u.ID.String())
	}

	return nil
}
