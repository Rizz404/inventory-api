package model

import (
	"log"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type AssetMovement struct {
	ID             SQLULID   `gorm:"primaryKey;type:varchar(26)"`
	AssetID        SQLULID   `gorm:"type:varchar(26);not null"`
	FromLocationID *SQLULID  `gorm:"type:varchar(26)"`
	ToLocationID   *SQLULID  `gorm:"type:varchar(26)"`
	FromUserID     *SQLULID  `gorm:"type:varchar(26)"`
	ToUserID       *SQLULID  `gorm:"type:varchar(26)"`
	MovementDate   time.Time `gorm:"not null"`
	MovedBy        SQLULID   `gorm:"type:varchar(26);not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Asset          Asset                      `gorm:"foreignKey:AssetID"`
	FromLocation   *Location                  `gorm:"foreignKey:FromLocationID"`
	ToLocation     *Location                  `gorm:"foreignKey:ToLocationID"`
	FromUser       *User                      `gorm:"foreignKey:FromUserID"`
	ToUser         *User                      `gorm:"foreignKey:ToUserID"`
	MovedByUser    User                       `gorm:"foreignKey:MovedBy"`
	Translations   []AssetMovementTranslation `gorm:"foreignKey:MovementID"`
}

func (AssetMovement) TableName() string {
	return "asset_movements"
}

func (u *AssetMovement) BeforeCreate(tx *gorm.DB) error {
	log.Printf("🚀 AssetMovement.BeforeCreate called! Current ID: %s, IsZero: %t", u.ID.String(), u.ID.IsZero())

	if u.ID.IsZero() {
		u.ID = SQLULID(ulid.Make())
		log.Printf("🚀 Generated new ULID for AssetMovement: %s", u.ID.String())
	}

	return nil
}

type AssetMovementTranslation struct {
	ID         SQLULID `gorm:"primaryKey;type:varchar(26)"`
	MovementID SQLULID `gorm:"type:varchar(26);not null;uniqueIndex:idx_mov_lang"`
	LangCode   string  `gorm:"type:varchar(5);not null;uniqueIndex:idx_mov_lang"`
	Notes      *string `gorm:"type:text"`
}

func (AssetMovementTranslation) TableName() string {
	return "asset_movement_translations"
}

func (u *AssetMovementTranslation) BeforeCreate(tx *gorm.DB) error {
	log.Printf("🚀 AssetMovementTranslation.BeforeCreate called! Current ID: %s, IsZero: %t", u.ID.String(), u.ID.IsZero())

	if u.ID.IsZero() {
		u.ID = SQLULID(ulid.Make())
		log.Printf("🚀 Generated new ULID for AssetMovementTranslation: %s", u.ID.String())
	}

	return nil
}
