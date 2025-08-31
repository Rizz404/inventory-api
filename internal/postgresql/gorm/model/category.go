package model

import (
	"log"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Category struct {
	ID           SQLULID  `gorm:"primaryKey;type:varchar(26)"`
	ParentID     *SQLULID `gorm:"type:varchar(26)"`
	CategoryCode string   `gorm:"type:varchar(20);unique;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Parent       *Category             `gorm:"foreignKey:ParentID"`
	Children     []Category            `gorm:"foreignKey:ParentID"`
	Translations []CategoryTranslation `gorm:"foreignKey:CategoryID"`
}

func (Category) TableName() string {
	return "categories"
}

func (u *Category) BeforeCreate(tx *gorm.DB) error {
	log.Printf("🚀 Category.BeforeCreate called! Current ID: %s, IsZero: %t", u.ID.String(), u.ID.IsZero())

	if u.ID.IsZero() {
		u.ID = SQLULID(ulid.Make())
		log.Printf("🚀 Generated new ULID for Category: %s", u.ID.String())
	}

	return nil
}

type CategoryTranslation struct {
	ID           SQLULID `gorm:"primaryKey;type:varchar(26)"`
	CategoryID   SQLULID `gorm:"type:varchar(26);not null;uniqueIndex:idx_cat_lang"`
	LangCode     string  `gorm:"type:varchar(5);not null;uniqueIndex:idx_cat_lang"`
	CategoryName string  `gorm:"type:varchar(100);not null"`
	Description  *string `gorm:"type:text"`
}

func (CategoryTranslation) TableName() string {
	return "categories_translation"
}

func (u *CategoryTranslation) BeforeCreate(tx *gorm.DB) error {
	log.Printf("🚀 CategoryTranslation.BeforeCreate called! Current ID: %s, IsZero: %t", u.ID.String(), u.ID.IsZero())

	if u.ID.IsZero() {
		u.ID = SQLULID(ulid.Make())
		log.Printf("🚀 Generated new ULID for CategoryTranslation: %s", u.ID.String())
	}

	return nil
}
