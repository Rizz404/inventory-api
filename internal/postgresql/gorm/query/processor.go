package query

import (
	"gorm.io/gorm"
)

type FilterApplier func(db *gorm.DB, filters any) *gorm.DB
type SortApplier func(db *gorm.DB, sort *SortOptions) *gorm.DB

// * Apply menerapkan semua parameter query ke *gorm.DB.
func Apply(
	db *gorm.DB,
	params Params,
	applyFilters FilterApplier,
	applySorts SortApplier,
) *gorm.DB {
	if params.Filters != nil && applyFilters != nil {
		db = applyFilters(db, params.Filters)
	}

	if params.Sort != nil && applySorts != nil {
		db = applySorts(db, params.Sort)
	} else {
		db = db.Order("id DESC")
	}

	if params.Pagination != nil {
		if params.Pagination.Cursor != "" {
			// * Mengasumsikan sorting DESC by ID untuk cursor
			// * Bisa dibuat lebih kompleks jika diperlukan
			db = db.Where("id < ?", params.Pagination.Cursor)
		} else {
			db = db.Offset(params.Pagination.Offset)
		}
		db = db.Limit(params.Pagination.Limit)
	}

	return db
}
