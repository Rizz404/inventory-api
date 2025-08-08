package model

import "time"

type Location struct {
	ID           SQLULID `gorm:"primaryKey;type:varchar(26)"`
	LocationCode string  `gorm:"type:varchar(20);unique;not null"`
	Building     *string `gorm:"type:varchar(100)"`
	Floor        *string `gorm:"type:varchar(20)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Translations []LocationTranslation `gorm:"foreignKey:LocationID"`
}

func (Location) TableName() string {
	return "locations"
}

type LocationTranslation struct {
	ID           SQLULID `gorm:"primaryKey;type:varchar(26)"`
	LocationID   SQLULID `gorm:"type:varchar(26);not null;uniqueIndex:idx_loc_lang"`
	LangCode     string  `gorm:"type:varchar(5);not null;uniqueIndex:idx_loc_lang"`
	LocationName string  `gorm:"type:varchar(100);not null"`
}

func (LocationTranslation) TableName() string {
	return "locations_translation"
}
