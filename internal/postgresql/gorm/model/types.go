package model

import (
	"database/sql/driver"
	"fmt"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type SQLULID ulid.ULID

func (u *SQLULID) BeforeCreate(tx *gorm.DB) (err error) {
	*u = SQLULID(ulid.Make())
	return
}

func (u SQLULID) GormDataType() string {
	return "varchar(26)"
}

func (u SQLULID) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "varchar(26)"
}

func (u *SQLULID) Scan(value interface{}) error {
	if value == nil {
		*u = SQLULID(ulid.ULID{})
		return nil
	}
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("SQLULID must be a string, got %T", value)
	}
	id, err := ulid.Parse(s)
	if err != nil {
		return err
	}
	*u = SQLULID(id)
	return nil
}

func (u SQLULID) Value() (driver.Value, error) {
	return ulid.ULID(u).String(), nil
}

func (u SQLULID) String() string {
	return ulid.ULID(u).String()
}
