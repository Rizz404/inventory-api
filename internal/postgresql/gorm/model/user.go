package model

import (
	"log"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type User struct {
	ID            SQLULID         `gorm:"primaryKey;type:varchar(26)"`
	Name          string          `gorm:"type:varchar(50);unique;not null"`
	Email         string          `gorm:"type:varchar(255);unique;not null"`
	PasswordHash  string          `gorm:"type:varchar(255);not null"`
	FullName      string          `gorm:"type:varchar(100);not null"`
	Role          domain.UserRole `gorm:"type:user_role;not null"`
	EmployeeID    *string         `gorm:"type:varchar(20);unique"`
	PreferredLang string          `gorm:"type:varchar(5);default:'id-ID'"`
	IsActive      bool            `gorm:"default:true"`
	AvatarURL     *string         `gorm:"type:varchar(255)"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	log.Printf("🚀 User.BeforeCreate called! Current ID: %s, IsZero: %t", u.ID.String(), u.ID.IsZero())

	if u.ID.IsZero() {
		u.ID = SQLULID(ulid.Make())
		log.Printf("🚀 Generated new ULID for User: %s", u.ID.String())
	}

	return nil
}
