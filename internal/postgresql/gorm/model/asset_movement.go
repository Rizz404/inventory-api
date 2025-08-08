package model

import "time"

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

type AssetMovementTranslation struct {
	ID         SQLULID `gorm:"primaryKey;type:varchar(26)"`
	MovementID SQLULID `gorm:"type:varchar(26);not null;uniqueIndex:idx_mov_lang"`
	LangCode   string  `gorm:"type:varchar(5);not null;uniqueIndex:idx_mov_lang"`
	Notes      *string `gorm:"type:text"`
}

func (AssetMovementTranslation) TableName() string {
	return "asset_movements_translation"
}
