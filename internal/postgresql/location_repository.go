package postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *LocationRepository) applyLocationFilters(db *gorm.DB, filters *domain.LocationFilterOptions) *gorm.DB {
	if filters == nil {
		return db
	}

	return db
}

func (r *LocationRepository) applyLocationSorts(db *gorm.DB, sort *domain.LocationSortOptions, langCode string) *gorm.DB {
	if sort == nil || sort.Field == "" {
		return db.Order("l.created_at DESC")
	}

	order := "DESC"
	if sort.Order == domain.SortOrderAsc {
		order = "ASC"
	}

	if sort.Field == domain.LocationSortByLocationName {
		// Use subquery for sorting by translation
		subquery := fmt.Sprintf("(SELECT location_name FROM location_translations WHERE location_id = l.id AND lang_code = '%s' LIMIT 1)", langCode)
		return db.Order(fmt.Sprintf("%s %s", subquery, order))
	}

	// Map camelCase sort field to snake_case database column
	columnName := mapper.MapLocationSortFieldToColumn(sort.Field)
	orderClause := columnName
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

	// Return created location with translations (no need to query again)
	// GORM has already filled the model with created data including ID and timestamps
	domainLocation := mapper.ToDomainLocation(&modelLocation)
	// Add translations manually since they were created separately
	for _, translation := range payload.Translations {
		domainLocation.Translations = append(domainLocation.Translations, domain.LocationTranslation{
			LangCode:     translation.LangCode,
			LocationName: translation.LocationName,
		})
	}
	return domainLocation, nil
}

func (r *LocationRepository) BulkCreateLocations(ctx context.Context, locations []domain.Location) ([]domain.Location, error) {
	if len(locations) == 0 {
		return []domain.Location{}, nil
	}

	models := make([]*model.Location, len(locations))
	for i := range locations {
		m := mapper.ToModelLocationForCreate(&locations[i])
		models[i] = &m
	}

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.
		Omit(clause.Associations).
		Session(&gorm.Session{CreateBatchSize: 500}).
		Create(&models).Error; err != nil {
		tx.Rollback()
		return nil, domain.ErrInternal(err)
	}

	var translations []model.LocationTranslation
	for i := range models {
		l := locations[i]
		for _, t := range l.Translations {
			mt := mapper.ToModelLocationTranslationForCreate(models[i].ID.String(), &t)
			translations = append(translations, mt)
		}
	}
	if len(translations) > 0 {
		if err := tx.Session(&gorm.Session{CreateBatchSize: 500}).Create(&translations).Error; err != nil {
			tx.Rollback()
			return nil, domain.ErrInternal(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	created := make([]domain.Location, len(models))
	for i := range models {
		created[i] = mapper.ToDomainLocation(models[i])
		for _, t := range locations[i].Translations {
			created[i].Translations = append(created[i].Translations, domain.LocationTranslation{
				LangCode:     t.LangCode,
				LocationName: t.LocationName,
			})
		}
	}
	return created, nil
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
	var updatedLocation model.Location
	err := r.db.WithContext(ctx).
		Table("locations l").
		Preload("Translations").
		First(&updatedLocation, "l.id = ?", locationId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Location{}, domain.ErrNotFound("location")
		}
		return domain.Location{}, domain.ErrInternal(err)
	}
	return mapper.ToDomainLocation(&updatedLocation), nil
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

func (r *LocationRepository) BulkDeleteLocations(ctx context.Context, locationIds []string) (domain.BulkDeleteLocations, error) {
	result := domain.BulkDeleteLocations{
		RequestedIDS: locationIds,
		DeletedIDS:   []string{},
	}

	if len(locationIds) == 0 {
		return result, nil
	}

	// First, find which locations actually exist
	var existingLocations []model.Location
	if err := r.db.WithContext(ctx).Select("id").Where("id IN ?", locationIds).Find(&existingLocations).Error; err != nil {
		return result, domain.ErrInternal(err)
	}

	// Collect existing location IDs
	existingIds := make([]string, 0, len(existingLocations))
	for _, location := range existingLocations {
		existingIds = append(existingIds, location.ID.String())
	}

	// If no locations exist, return early
	if len(existingIds) == 0 {
		return result, nil
	}

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return result, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete translations first (foreign key constraint)
	if err := tx.Delete(&model.LocationTranslation{}, "location_id IN ?", existingIds).Error; err != nil {
		tx.Rollback()
		return result, domain.ErrInternal(err)
	}

	// Delete locations
	if err := tx.Delete(&model.Location{}, "id IN ?", existingIds).Error; err != nil {
		tx.Rollback()
		return result, domain.ErrInternal(err)
	}

	if err := tx.Commit().Error; err != nil {
		return result, domain.ErrInternal(err)
	}

	result.DeletedIDS = existingIds
	return result, nil
}

// *===========================QUERY===========================*
func (r *LocationRepository) GetLocationsPaginated(ctx context.Context, params domain.LocationParams, langCode string) ([]domain.Location, error) {
	var locations []model.Location
	db := r.db.WithContext(ctx).
		Table("locations l").
		Preload("Translations")

	needsJoin := params.SearchQuery != nil && *params.SearchQuery != "" ||
		(params.Sort != nil && params.Sort.Field == domain.LocationSortByLocationName)

	if needsJoin {
		db = db.Select("l.id, l.location_code, l.building, l.floor, l.latitude, l.longitude, l.created_at, l.updated_at").
			Joins("LEFT JOIN location_translations lt ON l.id = lt.location_id")
		if params.SearchQuery != nil && *params.SearchQuery != "" {
			searchPattern := "%" + *params.SearchQuery + "%"
			db = db.Where("l.location_code ILIKE ? OR lt.location_name ILIKE ?", searchPattern, searchPattern).
				Distinct("l.id, l.created_at")
		}
	}

	// Apply filters
	db = r.applyLocationFilters(db, params.Filters)

	// Apply sorting
	db = r.applyLocationSorts(db, params.Sort, langCode)

	// Apply pagination
	if params.Pagination != nil {
		if params.Pagination.Limit > 0 {
			db = db.Limit(params.Pagination.Limit)
		}
		if params.Pagination.Offset > 0 {
			db = db.Offset(params.Pagination.Offset)
		}
	}

	if err := db.Find(&locations).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain locations
	return mapper.ToDomainLocations(locations), nil
}

func (r *LocationRepository) GetLocationsCursor(ctx context.Context, params domain.LocationParams, langCode string) ([]domain.Location, error) {
	var locations []model.Location
	db := r.db.WithContext(ctx).
		Table("locations l").
		Preload("Translations")

	needsJoin := params.SearchQuery != nil && *params.SearchQuery != "" ||
		(params.Sort != nil && params.Sort.Field == domain.LocationSortByLocationName)

	if needsJoin {
		db = db.Select("l.id, l.location_code, l.building, l.floor, l.latitude, l.longitude, l.created_at, l.updated_at").
			Joins("LEFT JOIN location_translations lt ON l.id = lt.location_id")
		if params.SearchQuery != nil && *params.SearchQuery != "" {
			searchPattern := "%" + *params.SearchQuery + "%"
			db = db.Where("l.location_code ILIKE ? OR lt.location_name ILIKE ?", searchPattern, searchPattern).
				Distinct("l.id, l.created_at")
		}
	}

	// Apply filters
	db = r.applyLocationFilters(db, params.Filters)

	// Apply sorting
	db = r.applyLocationSorts(db, params.Sort, langCode)

	// Apply cursor-based pagination
	if params.Pagination != nil {
		if params.Pagination.Limit > 0 {
			db = db.Limit(params.Pagination.Limit)
		}
		if params.Pagination.Cursor != "" {
			// Assuming sorting DESC by ID for cursor
			db = db.Where("l.id < ?", params.Pagination.Cursor)
		}
	}

	if err := db.Find(&locations).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain locations
	return mapper.ToDomainLocations(locations), nil
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

func (r *LocationRepository) CheckLocationCodeExistExcluding(ctx context.Context, locationCode string, excludeLocationId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("locations l").Where("l.location_code = ? AND l.id != ?", locationCode, excludeLocationId).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *LocationRepository) CountLocations(ctx context.Context, params domain.LocationParams) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Table("locations l")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN location_translations lt ON l.id = lt.location_id").
			Where("l.location_code ILIKE ? OR lt.location_name ILIKE ?", searchPattern, searchPattern).
			Distinct("l.id")
	}

	// Apply filters
	db = r.applyLocationFilters(db, params.Filters)

	if err := db.Count(&count).Error; err != nil {
		return 0, domain.ErrInternal(err)
	}
	return count, nil
}

func (r *LocationRepository) GetLocationStatistics(ctx context.Context) (domain.LocationStatistics, error) {
	var stats domain.LocationStatistics

	// Get total location count
	var totalCount int64
	if err := r.db.WithContext(ctx).Model(&model.Location{}).Count(&totalCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.Total.Count = int(totalCount)

	// Get location counts by building
	var buildingStats []struct {
		Building string `json:"building"`
		Count    int64  `json:"count"`
	}
	if err := r.db.WithContext(ctx).Model(&model.Location{}).
		Select("COALESCE(building, 'No Building') as building, COUNT(*) as count").
		Group("building").
		Order("count DESC").
		Scan(&buildingStats).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.ByBuilding = make([]domain.BuildingStatistics, len(buildingStats))
	for i, bs := range buildingStats {
		stats.ByBuilding[i] = domain.BuildingStatistics{
			Building: bs.Building,
			Count:    int(bs.Count),
		}
	}

	// Get location counts by floor
	var floorStats []struct {
		Floor string `json:"floor"`
		Count int64  `json:"count"`
	}
	if err := r.db.WithContext(ctx).Model(&model.Location{}).
		Select("COALESCE(floor, 'No Floor') as floor, COUNT(*) as count").
		Group("floor").
		Order("count DESC").
		Scan(&floorStats).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.ByFloor = make([]domain.FloorStatistics, len(floorStats))
	for i, fs := range floorStats {
		stats.ByFloor[i] = domain.FloorStatistics{
			Floor: fs.Floor,
			Count: int(fs.Count),
		}
	}

	// Get geographic statistics
	var withCoordinates, withoutCoordinates int64
	if err := r.db.WithContext(ctx).Model(&model.Location{}).
		Where("latitude IS NOT NULL AND longitude IS NOT NULL").
		Count(&withCoordinates).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.Location{}).
		Where("latitude IS NULL OR longitude IS NULL").
		Count(&withoutCoordinates).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.Geographic.WithCoordinates = int(withCoordinates)
	stats.Geographic.WithoutCoordinates = int(withoutCoordinates)

	// Get average coordinates if there are locations with coordinates
	if withCoordinates > 0 {
		var avgLat, avgLng float64
		if err := r.db.WithContext(ctx).Model(&model.Location{}).
			Select("AVG(latitude)").
			Where("latitude IS NOT NULL").
			Scan(&avgLat).Error; err != nil {
			return stats, domain.ErrInternal(err)
		}
		if err := r.db.WithContext(ctx).Model(&model.Location{}).
			Select("AVG(longitude)").
			Where("longitude IS NOT NULL").
			Scan(&avgLng).Error; err != nil {
			return stats, domain.ErrInternal(err)
		}
		stats.Geographic.AverageLatitude = &avgLat
		stats.Geographic.AverageLongitude = &avgLng
	}

	// Get creation trends (last 30 days)
	var creationTrends []struct {
		Date  time.Time `json:"date"`
		Count int64     `json:"count"`
	}
	if err := r.db.WithContext(ctx).Model(&model.Location{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= NOW() - INTERVAL '30 days'").
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&creationTrends).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.CreationTrends = make([]domain.LocationCreationTrend, len(creationTrends))
	for i, ct := range creationTrends {
		stats.CreationTrends[i] = domain.LocationCreationTrend{
			Date:  ct.Date,
			Count: int(ct.Count),
		}
	}

	// Calculate summary statistics
	stats.Summary.TotalLocations = int(totalCount)

	// Get locations with/without building
	var withBuilding, withoutBuilding int64
	if err := r.db.WithContext(ctx).Model(&model.Location{}).
		Where("building IS NOT NULL AND building != ''").
		Count(&withBuilding).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.Location{}).
		Where("building IS NULL OR building = ''").
		Count(&withoutBuilding).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.Summary.LocationsWithBuilding = int(withBuilding)
	stats.Summary.LocationsWithoutBuilding = int(withoutBuilding)

	// Get locations with/without floor
	var withFloor, withoutFloor int64
	if err := r.db.WithContext(ctx).Model(&model.Location{}).
		Where("floor IS NOT NULL AND floor != ''").
		Count(&withFloor).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.Location{}).
		Where("floor IS NULL OR floor = ''").
		Count(&withoutFloor).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.Summary.LocationsWithFloor = int(withFloor)
	stats.Summary.LocationsWithoutFloor = int(withoutFloor)

	stats.Summary.LocationsWithCoordinates = int(withCoordinates)

	// Calculate percentages
	if totalCount > 0 {
		stats.Summary.CoordinatesPercentage = float64(withCoordinates) / float64(totalCount) * 100
		stats.Summary.BuildingPercentage = float64(withBuilding) / float64(totalCount) * 100
		stats.Summary.FloorPercentage = float64(withFloor) / float64(totalCount) * 100
	}

	// Get total unique buildings and floors
	var totalBuildings, totalFloors int64
	if err := r.db.WithContext(ctx).Model(&model.Location{}).
		Select("COUNT(DISTINCT building)").
		Where("building IS NOT NULL AND building != ''").
		Scan(&totalBuildings).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.Location{}).
		Select("COUNT(DISTINCT floor)").
		Where("floor IS NOT NULL AND floor != ''").
		Scan(&totalFloors).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.Summary.TotalBuildings = int(totalBuildings)
	stats.Summary.TotalFloors = int(totalFloors)

	// Get earliest and latest creation dates
	var earliestDate, latestDate time.Time
	if err := r.db.WithContext(ctx).Model(&model.Location{}).Select("MIN(created_at)").Scan(&earliestDate).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.Location{}).Select("MAX(created_at)").Scan(&latestDate).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.Summary.EarliestCreationDate = earliestDate
	stats.Summary.LatestCreationDate = latestDate

	// Calculate average locations per day
	if !earliestDate.IsZero() && !latestDate.IsZero() {
		daysDiff := latestDate.Sub(earliestDate).Hours() / 24
		if daysDiff > 0 {
			stats.Summary.AverageLocationsPerDay = float64(totalCount) / daysDiff
		}
	}

	return stats, nil
}
