package model

import (
	"time"

	"github.com/Rizz404/inventory-api/domain"
)

type Notification struct {
	ID             SQLULID                 `gorm:"primaryKey;type:varchar(26)"`
	UserID         SQLULID                 `gorm:"type:varchar(26);not null"`
	RelatedAssetID *SQLULID                `gorm:"type:varchar(26)"`
	Type           domain.NotificationType `gorm:"type:notification_type;not null"`
	IsRead         bool                    `gorm:"default:false"`
	CreatedAt      time.Time
	User           User                      `gorm:"foreignKey:UserID"`
	Asset          *Asset                    `gorm:"foreignKey:RelatedAssetID"`
	Translations   []NotificationTranslation `gorm:"foreignKey:NotificationID"`
}

func (Notification) TableName() string {
	return "notifications"
}

type NotificationTranslation struct {
	ID             SQLULID `gorm:"primaryKey;type:varchar(26)"`
	NotificationID SQLULID `gorm:"type:varchar(26);not null;uniqueIndex:idx_notif_lang"`
	LangCode       string  `gorm:"type:varchar(5);not null;uniqueIndex:idx_notif_lang"`
	Title          string  `gorm:"type:varchar(200);not null"`
	Message        string  `gorm:"type:text;not null"`
}

func (NotificationTranslation) TableName() string {
	return "notifications_translation"
}
