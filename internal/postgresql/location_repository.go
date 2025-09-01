package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type LocationRepository struct {
	db *gorm.DB
}

type LocationFilterOptions struct {
}

func NewLocationRepository(db *gorm.DB) *LocationRepository {
	return &LocationRepository{
		db: db,
	}
}

func (r *LocationRepository) applyLocationFilters(db *gorm.DB, filters any) *gorm.DB {
	f, ok := filters.(*LocationFilterOptions)
	if !ok || f == nil {
		return db
	}

	return db
}

func (r *LocationRepository) applyLocationSorts(db *gorm.DB, sort *query.SortOptions) *gorm.DB {
	if sort == nil || sort.Field == "" {
		return db.Order("l.created_at DESC")
	}
	var orderClause string
	switch strings.ToLower(sort.Field) {
	case "location_code", "building", "floor", "created_at", "updated_at":
		orderClause = "l." + sort.Field
	case "name", "location_name":
		orderClause = "lt.location_name"
	default:
		return db.Order("l.created_at DESC")
	}

	order := "DESC"
	if strings.ToLower(sort.Order) == "asc" {
		order = "ASC"
	}
	return db.Order(fmt.Sprintf("%s %s", orderClause, order))
}

// *===========================MUTATION===========================*
func (r *LocationRepository) CreateLocation(ctx context.Context, payload *domain.Location) (domain.Location, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.Location{}, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create location
	modelLocation := mapper.ToModelLocationForCreate(payload)
	if err := tx.Create(&modelLocation).Error; err != nil {
		tx.Rollback()
		return domain.Location{}, domain.ErrInternal(err)
	}

	// Create translations
	for _, translation := range payload.Translations {
		modelTranslation := mapper.ToModelLocationTranslationForCreate(modelLocation.ID.String(), &translation)
		if err := tx.Create(&modelTranslation).Error; err != nil {
			tx.Rollback()
			return domain.Location{}, domain.ErrInternal(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return domain.Location{}, domain.ErrInternal(err)
	}

	// Fetch created location with translations
	return r.GetLocationById(ctx, modelLocation.ID.String())
}

func (r *LocationRepository) UpdateLocation(ctx context.Context, locationId string, payload *domain.UpdateLocationPayload) (domain.Location, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.Location{}, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update location basic info
	updates := mapper.ToModelLocationUpdateMap(payload)
	if len(updates) > 0 {
		if err := tx.Model(&model.Location{}).Where("id = ?", locationId).Updates(updates).Error; err != nil {
			tx.Rollback()
			return domain.Location{}, domain.ErrInternal(err)
		}
	}

	// Update translations if provided
	if len(payload.Translations) > 0 {
		for _, translationPayload := range payload.Translations {
			translationUpdates := mapper.ToModelLocationTranslationUpdateMap(&translationPayload)
			if len(translationUpdates) > 0 {
				// Try to update existing translation
				result := tx.Model(&model.LocationTranslation{}).
					Where("location_id = ? AND lang_code = ?", locationId, translationPayload.LangCode).
					Updates(translationUpdates)

				if result.Error != nil {
					tx.Rollback()
					return domain.Location{}, domain.ErrInternal(result.Error)
				}

				// If no rows affected, create new translation
				if result.RowsAffected == 0 {
					newTranslation := model.LocationTranslation{
						LangCode:     translationPayload.LangCode,
						LocationName: *translationPayload.LocationName,
					}
					if parsedLocationID, err := ulid.Parse(locationId); err == nil {
						newTranslation.LocationID = model.SQLULID(parsedLocationID)
					}

					if err := tx.Create(&newTranslation).Error; err != nil {
						tx.Rollback()
						return domain.Location{}, domain.ErrInternal(err)
					}
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return domain.Location{}, domain.ErrInternal(err)
	}

	// Fetch updated location with translations
	return r.GetLocationById(ctx, locationId)
}

func (r *LocationRepository) DeleteLocation(ctx context.Context, locationId string) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.ErrInternal(tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete translations first (foreign key constraint)
	if err := tx.Delete(&model.LocationTranslation{}, "location_id = ?", locationId).Error; err != nil {
		tx.Rollback()
		return domain.ErrInternal(err)
	}

	// Delete location
	if err := tx.Delete(&model.Location{}, "id = ?", locationId).Error; err != nil {
		tx.Rollback()
		return domain.ErrInternal(err)
	}

	if err := tx.Commit().Error; err != nil {
		return domain.ErrInternal(err)
	}

	return nil
}

// *===========================QUERY===========================*
func (r *LocationRepository) GetLocationsPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.Location, error) {
	var locations []model.Location
	db := r.db.WithContext(ctx).
		Table("locations l").
		Preload("Translations")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN location_translations lt ON l.id = lt.location_id").
			Where("l.location_code ILIKE ? OR lt.location_name ILIKE ?", searchPattern, searchPattern).
			Distinct("l.id")
	}

	// Set pagination cursor to empty for offset-based pagination
	params.Pagination.Cursor = ""
	db = query.Apply(db, params, r.applyLocationFilters, r.applyLocationSorts)

	if err := db.Find(&locations).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain locations
	domainLocations := make([]domain.Location, len(locations))
	for i, location := range locations {
		domainLocations[i] = mapper.ToDomainLocation(&location)
	}
	return domainLocations, nil
}

func (r *LocationRepository) GetLocationsCursor(ctx context.Context, params query.Params, langCode string) ([]domain.Location, error) {
	var locations []model.Location
	db := r.db.WithContext(ctx).
		Table("locations l").
		Preload("Translations")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN location_translations lt ON l.id = lt.location_id").
			Where("l.location_code ILIKE ? OR lt.location_name ILIKE ?", searchPattern, searchPattern).
			Distinct("l.id")
	}

	// Set offset to 0 for cursor-based pagination
	params.Pagination.Offset = 0
	db = query.Apply(db, params, r.applyLocationFilters, r.applyLocationSorts)

	if err := db.Find(&locations).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain locations
	domainLocations := make([]domain.Location, len(locations))
	for i, location := range locations {
		domainLocations[i] = mapper.ToDomainLocation(&location)
	}
	return domainLocations, nil
}

func (r *LocationRepository) GetLocationsResponse(ctx context.Context, params query.Params, langCode string) ([]domain.LocationResponse, error) {
	var locations []model.Location
	db := r.db.WithContext(ctx).
		Table("locations l").
		Preload("Translations")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN location_translations lt ON l.id = lt.location_id").
			Where("l.location_code ILIKE ? OR lt.location_name ILIKE ?", searchPattern, searchPattern).
			Distinct("l.id")
	}

	db = query.Apply(db, params, r.applyLocationFilters, r.applyLocationSorts)

	if err := db.Find(&locations).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	return mapper.ToDomainLocationsResponse(locations, langCode), nil
}

func (r *LocationRepository) GetLocationById(ctx context.Context, locationId string) (domain.Location, error) {
	var location model.Location

	err := r.db.WithContext(ctx).
		Table("locations l").
		Preload("Translations").
		First(&location, "l.id = ?", locationId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Location{}, domain.ErrNotFound("location")
		}
		return domain.Location{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainLocation(&location), nil
}

func (r *LocationRepository) GetLocationByCode(ctx context.Context, locationCode string) (domain.Location, error) {
	var location model.Location

	err := r.db.WithContext(ctx).
		Table("locations l").
		Preload("Translations").
		First(&location, "l.location_code = ?", locationCode).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Location{}, domain.ErrNotFound("location")
		}
		return domain.Location{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainLocation(&location), nil
}

func (r *LocationRepository) GetLocationHierarchy(ctx context.Context, langCode string) ([]domain.LocationResponse, error) {
	locations, err := r.GetLocationsResponse(ctx, query.Params{
		Pagination: &query.PaginationOptions{
			Limit: 1000, // Large limit to get all locations
		},
	}, langCode)
	if err != nil {
		return nil, err
	}

	return mapper.BuildLocationHierarchy(locations), nil
}

func (r *LocationRepository) CheckLocationExist(ctx context.Context, locationId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("locations l").Where("l.id = ?", locationId).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *LocationRepository) CheckLocationCodeExist(ctx context.Context, locationCode string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("locations l").Where("l.location_code = ?", locationCode).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *LocationRepository) CountLocations(ctx context.Context, params query.Params) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Table("locations l")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN location_translations lt ON l.id = lt.location_id").
			Where("l.location_code ILIKE ? OR lt.location_name ILIKE ?", searchPattern, searchPattern).
			Distinct("l.id")
	}

	db = query.Apply(db, query.Params{Filters: params.Filters}, r.applyLocationFilters, nil)

	if err := db.Count(&count).Error; err != nil {
		return 0, domain.ErrInternal(err)
	}
	return count, nil
}
