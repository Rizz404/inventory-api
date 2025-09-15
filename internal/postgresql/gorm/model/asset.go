package model

import (
	"log"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Asset struct {
	ID                 SQLULID               `gorm:"primaryKey;type:varchar(26)"`
	AssetTag           string                `gorm:"type:varchar(50);unique;not null"`
	DataMatrixValue    string                `gorm:"type:varchar(255);unique;not null"`
	DataMatrixImageUrl string                `gorm:"type:varchar(255);not null"`
	AssetName          string                `gorm:"type:varchar(200);not null"`
	CategoryID         SQLULID               `gorm:"type:varchar(26);not null"`
	Brand              *string               `gorm:"type:varchar(100)"`
	Model              *string               `gorm:"type:varchar(100)"`
	SerialNumber       *string               `gorm:"type:varchar(100);unique"`
	PurchaseDate       *time.Time            `gorm:"type:date"`
	PurchasePrice      *float64              `gorm:"type:decimal(15,2)"`
	VendorName         *string               `gorm:"type:varchar(150)"`
	WarrantyEnd        *time.Time            `gorm:"type:date"`
	Status             domain.AssetStatus    `gorm:"type:asset_status;default:'Active'"`
	Condition          domain.AssetCondition `gorm:"type:asset_condition;default:'Good';column:condition_status"`
	LocationID         *SQLULID              `gorm:"type:varchar(26)"`
	AssignedTo         *SQLULID              `gorm:"type:varchar(26)"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Category           Category  `gorm:"foreignKey:CategoryID"`
	Location           *Location `gorm:"foreignKey:LocationID"`
	User               *User     `gorm:"foreignKey:AssignedTo"`
}

func (Asset) TableName() string {
	return "assets"
}

func (u *Asset) BeforeCreate(tx *gorm.DB) error {
	log.Printf("🚀 Asset.BeforeCreate called! Current ID: %s, IsZero: %t", u.ID.String(), u.ID.IsZero())

	if u.ID.IsZero() {
		u.ID = SQLULID(ulid.Make())
		log.Printf("🚀 Generated new ULID for Asset: %s", u.ID.String())
	}

	return nil
}
