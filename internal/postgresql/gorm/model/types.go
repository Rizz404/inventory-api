package model

import (
	"database/sql/driver"
	"fmt"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type SQLULID ulid.ULID

// BeforeCreate hook untuk auto-generate ULID
func (u *SQLULID) BeforeCreate(tx *gorm.DB) error {
	// Selalu generate ULID baru jika masih zero value
	if *u == SQLULID(ulid.ULID{}) {
		*u = SQLULID(ulid.Make())
	}
	return nil
}

func (u SQLULID) GormDataType() string {
	return "varchar(26)"
}

func (u SQLULID) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "varchar(26)"
}

// Scan - membaca dari database
func (u *SQLULID) Scan(value interface{}) error {
	if value == nil {
		*u = SQLULID(ulid.ULID{})
		return nil
	}

	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("SQLULID must be a string, got %T", value)
	}

	// Handle empty string - set sebagai zero value
	if s == "" {
		*u = SQLULID(ulid.ULID{})
		return nil
	}

	id, err := ulid.Parse(s)
	if err != nil {
		return fmt.Errorf("failed to parse ULID: %w", err)
	}

	*u = SQLULID(id)
	return nil
}

// Value - menulis ke database
func (u SQLULID) Value() (driver.Value, error) {
	// Jika zero value, return nil (akan jadi NULL di database)
	// BUKAN empty string yang menyebabkan duplicate key
	if u == SQLULID(ulid.ULID{}) {
		return nil, nil
	}
	return ulid.ULID(u).String(), nil
}

func (u SQLULID) String() string {
	// Return empty string untuk zero value agar JSON serialization benar
	if u == SQLULID(ulid.ULID{}) {
		return ""
	}
	return ulid.ULID(u).String()
}

// IsZero mengecek apakah ULID adalah zero value
func (u SQLULID) IsZero() bool {
	return u == SQLULID(ulid.ULID{})
}
