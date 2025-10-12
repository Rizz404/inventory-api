package postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"gorm.io/gorm"
)

type AssetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) *AssetRepository {
	return &AssetRepository{
		db: db,
	}
}

func (r *AssetRepository) applyAssetFilters(db *gorm.DB, filters *domain.AssetFilterOptions) *gorm.DB {
	if filters == nil {
		return db
	}

	if filters.Status != nil {
		db = db.Where("a.status = ?", filters.Status)
	}
	if filters.Condition != nil {
		db = db.Where("a.condition_status = ?", filters.Condition)
	}
	if filters.CategoryID != nil {
		db = db.Where("a.category_id = ?", filters.CategoryID)
	}
	if filters.LocationID != nil {
		db = db.Where("a.location_id = ?", filters.LocationID)
	}
	if filters.AssignedTo != nil {
		db = db.Where("a.assigned_to = ?", filters.AssignedTo)
	}
	if filters.Brand != nil {
		db = db.Where("a.brand ILIKE ?", "%"+*filters.Brand+"%")
	}
	if filters.Model != nil {
		db = db.Where("a.model ILIKE ?", "%"+*filters.Model+"%")
	}
	return db
}

func (r *AssetRepository) applyAssetSorts(db *gorm.DB, sort *domain.AssetSortOptions) *gorm.DB {
	if sort == nil || sort.Field == "" {
		return db.Order("a.created_at DESC")
	}

	// Map camelCase sort field to snake_case database column
	columnName := mapper.MapAssetSortFieldToColumn(sort.Field)
	orderClause := columnName

	order := "DESC"
	if sort.Order == domain.SortOrderAsc {
		order = "ASC"
	}
	return db.Order(fmt.Sprintf("%s %s", orderClause, order))
}

// *===========================MUTATION===========================*
func (r *AssetRepository) CreateAsset(ctx context.Context, payload *domain.Asset) (domain.Asset, error) {
	modelAsset := mapper.ToModelAssetForCreate(payload)

	// Create asset in database
	err := r.db.WithContext(ctx).Create(&modelAsset).Error
	if err != nil {
		return domain.Asset{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainAsset(&modelAsset), nil
}

func (r *AssetRepository) UpdateAsset(ctx context.Context, assetId string, payload *domain.UpdateAssetPayload) (domain.Asset, error) {
	var updatedAsset model.Asset

	// Build update map from payload
	updates := mapper.ToModelAssetUpdateMap(payload)

	// Perform update
	err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("id = ?", assetId).Updates(updates).Error
	if err != nil {
		return domain.Asset{}, domain.ErrInternal(err)
	}

	// Get updated asset
	err = r.db.WithContext(ctx).Preload("Category").Preload("Category.Translations").Preload("Location").Preload("Location.Translations").Preload("User").First(&updatedAsset, "id = ?", assetId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Asset{}, domain.ErrNotFound("asset")
		}
		return domain.Asset{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainAsset(&updatedAsset), nil
}

func (r *AssetRepository) DeleteAsset(ctx context.Context, assetId string) error {
	err := r.db.WithContext(ctx).Delete(&model.Asset{}, "id = ?", assetId).Error
	if err != nil {
		return domain.ErrInternal(err)
	}
	return nil
}

// *===========================QUERY===========================*
func (r *AssetRepository) GetAssetsPaginated(ctx context.Context, params domain.AssetParams, langCode string) ([]domain.Asset, error) {
	var assets []model.Asset
	db := r.db.WithContext(ctx).
		Table("assets a").
		Preload("Category").
		Preload("Category.Translations").
		Preload("Location").
		Preload("Location.Translations").
		Preload("User")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("a.asset_tag ILIKE ? OR a.asset_name ILIKE ? OR a.brand ILIKE ? OR a.model ILIKE ? OR a.serial_number ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Apply filters
	db = r.applyAssetFilters(db, params.Filters)

	// Apply sorting
	db = r.applyAssetSorts(db, params.Sort)

	// Apply pagination
	if params.Pagination != nil {
		if params.Pagination.Limit > 0 {
			db = db.Limit(params.Pagination.Limit)
		}
		if params.Pagination.Offset > 0 {
			db = db.Offset(params.Pagination.Offset)
		}
	}

	if err := db.Find(&assets).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain assets using helper function
	return mapper.ToDomainAssets(assets), nil
}

func (r *AssetRepository) GetAssetsCursor(ctx context.Context, params domain.AssetParams, langCode string) ([]domain.Asset, error) {
	var assets []model.Asset
	db := r.db.WithContext(ctx).
		Table("assets a").
		Preload("Category").
		Preload("Category.Translations").
		Preload("Location").
		Preload("Location.Translations").
		Preload("User")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("a.asset_tag ILIKE ? OR a.asset_name ILIKE ? OR a.brand ILIKE ? OR a.model ILIKE ? OR a.serial_number ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Apply filters
	db = r.applyAssetFilters(db, params.Filters)

	// Apply sorting
	db = r.applyAssetSorts(db, params.Sort)

	// Apply cursor-based pagination
	if params.Pagination != nil {
		if params.Pagination.Limit > 0 {
			db = db.Limit(params.Pagination.Limit)
		}
		if params.Pagination.Cursor != "" {
			// Assuming sorting DESC by ID for cursor
			db = db.Where("a.id < ?", params.Pagination.Cursor)
		}
	}

	if err := db.Find(&assets).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain assets using helper function
	return mapper.ToDomainAssets(assets), nil
}

func (r *AssetRepository) GetAssetById(ctx context.Context, assetId string) (domain.Asset, error) {
	var asset model.Asset

	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Category.Translations").
		Preload("Location").
		Preload("Location.Translations").
		Preload("User").
		First(&asset, "id = ?", assetId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Asset{}, domain.ErrNotFound("asset")
		}
		return domain.Asset{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainAsset(&asset), nil
}

func (r *AssetRepository) GetAssetByAssetTag(ctx context.Context, assetTag string) (domain.Asset, error) {
	var asset model.Asset

	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Category.Translations").
		Preload("Location").
		Preload("Location.Translations").
		Preload("User").
		Where("asset_tag = ?", assetTag).First(&asset).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Asset{}, domain.ErrNotFound("asset with tag '" + assetTag + "'")
		}
		return domain.Asset{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainAsset(&asset), nil
}

func (r *AssetRepository) CheckAssetExists(ctx context.Context, assetId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("id = ?", assetId).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *AssetRepository) CheckAssetTagExists(ctx context.Context, assetTag string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("asset_tag = ?", assetTag).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *AssetRepository) CheckSerialNumberExists(ctx context.Context, serialNumber string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("serial_number = ?", serialNumber).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *AssetRepository) CheckAssetTagExistsExcluding(ctx context.Context, assetTag string, excludeAssetId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("asset_tag = ? AND id != ?", assetTag, excludeAssetId).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *AssetRepository) CheckSerialNumberExistsExcluding(ctx context.Context, serialNumber string, excludeAssetId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("serial_number = ? AND id != ?", serialNumber, excludeAssetId).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *AssetRepository) CountAssets(ctx context.Context, params domain.AssetParams) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Table("assets a")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("a.asset_tag ILIKE ? OR a.asset_name ILIKE ? OR a.brand ILIKE ? OR a.model ILIKE ? OR a.serial_number ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Apply filters
	db = r.applyAssetFilters(db, params.Filters)

	if err := db.Count(&count).Error; err != nil {
		return 0, domain.ErrInternal(err)
	}
	return count, nil
}

func (r *AssetRepository) GetAssetStatistics(ctx context.Context) (domain.AssetStatistics, error) {
	var stats domain.AssetStatistics

	// Get total asset count
	var totalCount int64
	if err := r.db.WithContext(ctx).Model(&model.Asset{}).Count(&totalCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.Total.Count = int(totalCount)

	// Get asset counts by status
	var activeCount, maintenanceCount, disposedCount, lostCount int64
	if err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("status = ?", domain.StatusActive).Count(&activeCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("status = ?", domain.StatusMaintenance).Count(&maintenanceCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("status = ?", domain.StatusDisposed).Count(&disposedCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("status = ?", domain.StatusLost).Count(&lostCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.ByStatus.Active = int(activeCount)
	stats.ByStatus.Maintenance = int(maintenanceCount)
	stats.ByStatus.Disposed = int(disposedCount)
	stats.ByStatus.Lost = int(lostCount)

	// Get asset counts by condition
	var goodCount, fairCount, poorCount, damagedCount int64
	if err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("condition_status = ?", domain.ConditionGood).Count(&goodCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("condition_status = ?", domain.ConditionFair).Count(&fairCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("condition_status = ?", domain.ConditionPoor).Count(&poorCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("condition_status = ?", domain.ConditionDamaged).Count(&damagedCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.ByCondition.Good = int(goodCount)
	stats.ByCondition.Fair = int(fairCount)
	stats.ByCondition.Poor = int(poorCount)
	stats.ByCondition.Damaged = int(damagedCount)

	// Get assignment statistics
	var assignedCount, unassignedCount int64
	if err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("assigned_to IS NOT NULL").Count(&assignedCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("assigned_to IS NULL").Count(&unassignedCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.ByAssignment.Assigned = int(assignedCount)
	stats.ByAssignment.Unassigned = int(unassignedCount)

	// Get value statistics
	var valueStats struct {
		TotalValue         *float64 `json:"total_value"`
		AverageValue       *float64 `json:"average_value"`
		MinValue           *float64 `json:"min_value"`
		MaxValue           *float64 `json:"max_value"`
		AssetsWithValue    int64    `json:"assets_with_value"`
		AssetsWithoutValue int64    `json:"assets_without_value"`
	}

	if err := r.db.WithContext(ctx).Model(&model.Asset{}).
		Select("SUM(purchase_price) as total_value, AVG(purchase_price) as average_value, MIN(purchase_price) as min_value, MAX(purchase_price) as max_value").
		Where("purchase_price IS NOT NULL").
		Scan(&valueStats).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	if err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("purchase_price IS NOT NULL").Count(&valueStats.AssetsWithValue).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("purchase_price IS NULL").Count(&valueStats.AssetsWithoutValue).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.ValueStatistics.TotalValue = valueStats.TotalValue
	stats.ValueStatistics.AverageValue = valueStats.AverageValue
	stats.ValueStatistics.MinValue = valueStats.MinValue
	stats.ValueStatistics.MaxValue = valueStats.MaxValue
	stats.ValueStatistics.AssetsWithValue = int(valueStats.AssetsWithValue)
	stats.ValueStatistics.AssetsWithoutValue = int(valueStats.AssetsWithoutValue)

	// Get warranty statistics
	var activeWarranties, expiredWarranties, noWarrantyInfo int64
	currentTime := time.Now()
	if err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("warranty_end IS NOT NULL AND warranty_end > ?", currentTime).Count(&activeWarranties).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("warranty_end IS NOT NULL AND warranty_end <= ?", currentTime).Count(&expiredWarranties).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("warranty_end IS NULL").Count(&noWarrantyInfo).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.WarrantyStatistics.ActiveWarranties = int(activeWarranties)
	stats.WarrantyStatistics.ExpiredWarranties = int(expiredWarranties)
	stats.WarrantyStatistics.NoWarrantyInfo = int(noWarrantyInfo)

	// Get creation trends (last 30 days)
	var creationTrends []struct {
		Date  time.Time `json:"date"`
		Count int64     `json:"count"`
	}
	if err := r.db.WithContext(ctx).Model(&model.Asset{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= NOW() - INTERVAL '30 days'").
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&creationTrends).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.CreationTrends = make([]domain.AssetCreationTrend, len(creationTrends))
	for i, ct := range creationTrends {
		stats.CreationTrends[i] = domain.AssetCreationTrend{
			Date:  ct.Date,
			Count: int(ct.Count),
		}
	}

	// Get asset counts by category
	var categoryStats []struct {
		CategoryID   string `json:"category_id"`
		CategoryCode string `json:"category_code"`
		CategoryName string `json:"category_name"`
		AssetCount   int64  `json:"asset_count"`
	}
	if err := r.db.WithContext(ctx).
		Table("assets a").
		Select("c.id as category_id, c.category_code, COALESCE(ct.category_name, c.category_code) as category_name, COUNT(a.id) as asset_count").
		Joins("INNER JOIN categories c ON a.category_id = c.id").
		Joins("LEFT JOIN category_translations ct ON c.id = ct.category_id AND ct.lang_code = 'en'").
		Group("c.id, c.category_code, ct.category_name").
		Order("asset_count DESC").
		Scan(&categoryStats).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.ByCategory = make([]domain.AssetByCategoryStatistics, len(categoryStats))
	for i, cs := range categoryStats {
		percentage := float64(0)
		if totalCount > 0 {
			percentage = float64(cs.AssetCount) / float64(totalCount) * 100
		}
		stats.ByCategory[i] = domain.AssetByCategoryStatistics{
			CategoryID:   cs.CategoryID,
			CategoryName: cs.CategoryName,
			CategoryCode: cs.CategoryCode,
			AssetCount:   int(cs.AssetCount),
			Percentage:   percentage,
		}
	}

	// Get asset counts by location
	var locationStats []struct {
		LocationID   string `json:"location_id"`
		LocationCode string `json:"location_code"`
		LocationName string `json:"location_name"`
		AssetCount   int64  `json:"asset_count"`
	}
	if err := r.db.WithContext(ctx).
		Table("assets a").
		Select("l.id as location_id, l.location_code, COALESCE(lt.location_name, l.location_code) as location_name, COUNT(a.id) as asset_count").
		Joins("INNER JOIN locations l ON a.location_id = l.id").
		Joins("LEFT JOIN location_translations lt ON l.id = lt.location_id AND lt.lang_code = 'en'").
		Where("a.location_id IS NOT NULL").
		Group("l.id, l.location_code, lt.location_name").
		Order("asset_count DESC").
		Scan(&locationStats).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.ByLocation = make([]domain.AssetByLocationStatistics, len(locationStats))
	for i, ls := range locationStats {
		percentage := float64(0)
		if totalCount > 0 {
			percentage = float64(ls.AssetCount) / float64(totalCount) * 100
		}
		stats.ByLocation[i] = domain.AssetByLocationStatistics{
			LocationID:   ls.LocationID,
			LocationName: ls.LocationName,
			LocationCode: ls.LocationCode,
			AssetCount:   int(ls.AssetCount),
			Percentage:   percentage,
		}
	}

	// Count unique categories and locations for summary
	var uniqueCategories, uniqueLocations int64
	if err := r.db.WithContext(ctx).Model(&model.Asset{}).Distinct("category_id").Count(&uniqueCategories).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("location_id IS NOT NULL").Distinct("location_id").Count(&uniqueLocations).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.Summary.TotalCategories = int(uniqueCategories)
	stats.Summary.TotalLocations = int(uniqueLocations)

	// Calculate summary statistics
	stats.Summary.TotalAssets = int(totalCount)

	if totalCount > 0 {
		stats.Summary.ActiveAssetsPercentage = float64(activeCount) / float64(totalCount) * 100
		stats.Summary.MaintenanceAssetsPercentage = float64(maintenanceCount) / float64(totalCount) * 100
		stats.Summary.DisposedAssetsPercentage = float64(disposedCount) / float64(totalCount) * 100
		stats.Summary.LostAssetsPercentage = float64(lostCount) / float64(totalCount) * 100
		stats.Summary.GoodConditionPercentage = float64(goodCount) / float64(totalCount) * 100
		stats.Summary.FairConditionPercentage = float64(fairCount) / float64(totalCount) * 100
		stats.Summary.PoorConditionPercentage = float64(poorCount) / float64(totalCount) * 100
		stats.Summary.DamagedConditionPercentage = float64(damagedCount) / float64(totalCount) * 100
		stats.Summary.AssignedAssetsPercentage = float64(assignedCount) / float64(totalCount) * 100
		stats.Summary.UnassignedAssetsPercentage = float64(unassignedCount) / float64(totalCount) * 100
		stats.Summary.PurchasePricePercentage = float64(valueStats.AssetsWithValue) / float64(totalCount) * 100
		stats.Summary.WarrantyPercentage = float64(activeWarranties+expiredWarranties) / float64(totalCount) * 100
	}

	// Additional summary fields
	stats.Summary.AssetsWithPurchasePrice = int(valueStats.AssetsWithValue)
	stats.Summary.AssetsWithWarranty = int(activeWarranties + expiredWarranties)
	stats.Summary.MostExpensiveAssetValue = valueStats.MaxValue
	stats.Summary.LeastExpensiveAssetValue = valueStats.MinValue

	// Get earliest and latest creation dates
	var earliestDate, latestDate time.Time
	if err := r.db.WithContext(ctx).Model(&model.Asset{}).Select("MIN(created_at)").Scan(&earliestDate).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.Asset{}).Select("MAX(created_at)").Scan(&latestDate).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.Summary.EarliestCreationDate = earliestDate
	stats.Summary.LatestCreationDate = latestDate

	// Calculate average assets per day
	if !earliestDate.IsZero() && !latestDate.IsZero() {
		daysDiff := latestDate.Sub(earliestDate).Hours() / 24
		if daysDiff > 0 {
			stats.Summary.AverageAssetsPerDay = float64(totalCount) / daysDiff
		}
	}

	return stats, nil
}

func (r *AssetRepository) GetLastAssetTagByCategory(ctx context.Context, categoryId string) (string, error) {
	var asset model.Asset

	// Get the last asset tag for the given category, ordered by asset_tag descending
	err := r.db.WithContext(ctx).
		Where("category_id = ?", categoryId).
		Order("asset_tag DESC").
		Limit(1).
		First(&asset).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return empty string if no asset found for this category
			return "", nil
		}
		return "", domain.ErrInternal(err)
	}

	return asset.AssetTag, nil
}

func (r *AssetRepository) GetAssetsForExport(ctx context.Context, params domain.AssetParams, langCode string) ([]domain.Asset, error) {
	var assets []model.Asset
	db := r.db.WithContext(ctx).
		Table("assets a").
		Preload("Category").
		Preload("Category.Translations").
		Preload("Location").
		Preload("Location.Translations").
		Preload("User")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("a.asset_tag ILIKE ? OR a.asset_name ILIKE ? OR a.brand ILIKE ? OR a.model ILIKE ? OR a.serial_number ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Apply filters
	db = r.applyAssetFilters(db, params.Filters)

	// Apply sorting
	db = r.applyAssetSorts(db, params.Sort)

	// No pagination for export - get all matching assets
	if err := db.Find(&assets).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain assets using helper function
	return mapper.ToDomainAssets(assets), nil
}
