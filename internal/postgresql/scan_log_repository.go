package postgresql

import "gorm.io/gorm"

type ScanLogRepository struct {
	db *gorm.DB
}

type ScanLogFilterOptions struct {
}
