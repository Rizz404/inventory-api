package model

import (
	"log"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Notification struct {
	ID     SQLULID `gorm:"primaryKey;type:varchar(26)"`
	UserID SQLULID `gorm:"type:varchar(26);not null"`

	// Related entity (polymorphic)
	RelatedEntityType *string  `gorm:"type:varchar(50)"`
	RelatedEntityID   *SQLULID `gorm:"type:varchar(26)"`

	// Legacy support (deprecated)
	RelatedAssetID *SQLULID `gorm:"type:varchar(26)"`

	Type     domain.NotificationType     `gorm:"type:notification_type;not null"`
	Priority domain.NotificationPriority `gorm:"type:notification_priority;default:'NORMAL'"`

	// Status
	IsRead bool       `gorm:"default:false"`
	ReadAt *time.Time `gorm:"type:timestamp with time zone"`

	// Expiration
	ExpiresAt *time.Time `gorm:"type:timestamp with time zone"`

	CreatedAt    time.Time
	Translations []NotificationTranslation `gorm:"foreignKey:NotificationID"`
}

func (Notification) TableName() string {
	return "notifications"
}

func (u *Notification) BeforeCreate(tx *gorm.DB) error {
	log.Printf("🚀 Notification.BeforeCreate called! Current ID: %s, IsZero: %t", u.ID.String(), u.ID.IsZero())

	if u.ID.IsZero() {
		u.ID = SQLULID(ulid.Make())
		log.Printf("🚀 Generated new ULID for Notification: %s", u.ID.String())
	}

	return nil
}

type NotificationTranslation struct {
	ID             SQLULID `gorm:"primaryKey;type:varchar(26)"`
	NotificationID SQLULID `gorm:"type:varchar(26);not null;uniqueIndex:idx_notif_lang"`
	LangCode       string  `gorm:"type:varchar(5);not null;uniqueIndex:idx_notif_lang"`
	Title          string  `gorm:"type:varchar(200);not null"`
	Message        string  `gorm:"type:text;not null"`
}

func (NotificationTranslation) TableName() string {
	return "notification_translations"
}

func (u *NotificationTranslation) BeforeCreate(tx *gorm.DB) error {
	log.Printf("🚀 NotificationTranslation.BeforeCreate called! Current ID: %s, IsZero: %t", u.ID.String(), u.ID.IsZero())

	if u.ID.IsZero() {
		u.ID = SQLULID(ulid.Make())
		log.Printf("🚀 Generated new ULID for NotificationTranslation: %s", u.ID.String())
	}

	return nil
}
