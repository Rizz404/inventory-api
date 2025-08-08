package model

import (
	"time"

	"github.com/Rizz404/inventory-api/domain"
)

type User struct {
	ID            SQLULID         `gorm:"primaryKey;type:varchar(26)"`
	Username      string          `gorm:"type:varchar(50);unique;not null"`
	Email         string          `gorm:"type:varchar(255);unique;not null"`
	PasswordHash  string          `gorm:"type:varchar(255);not null"`
	FullName      string          `gorm:"type:varchar(100);not null"`
	Role          domain.UserRole `gorm:"type:user_role;not null"`
	EmployeeID    *string         `gorm:"type:varchar(20);unique"`
	PreferredLang string          `gorm:"type:varchar(5);default:'id-ID'"`
	IsActive      bool            `gorm:"default:true"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (User) TableName() string {
	return "users"
}
