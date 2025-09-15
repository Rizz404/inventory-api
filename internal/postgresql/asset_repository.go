package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type AssetRepository struct {
	db *gorm.DB
}

type AssetFilterOptions struct {
	Status     *domain.AssetStatus    `json:"status,omitempty"`
	Condition  *domain.AssetCondition `json:"condition,omitempty"`
	CategoryID *string                `json:"category_id,omitempty"`
	LocationID *string                `json:"location_id,omitempty"`
	AssignedTo *string                `json:"assigned_to,omitempty"`
	Brand      *string                `json:"brand,omitempty"`
	Model      *string                `json:"model,omitempty"`
}

func NewAssetRepository(db *gorm.DB) *AssetRepository {
	return &AssetRepository{
		db: db,
	}
}

func (r *AssetRepository) applyAssetFilters(db *gorm.DB, filters any) *gorm.DB {
	f, ok := filters.(*AssetFilterOptions)
	if !ok || f == nil {
		return db
	}

	if f.Status != nil {
		db = db.Where("a.status = ?", f.Status)
	}
	if f.Condition != nil {
		db = db.Where("a.condition_status = ?", f.Condition)
	}
	if f.CategoryID != nil {
		db = db.Where("a.category_id = ?", f.CategoryID)
	}
	if f.LocationID != nil {
		db = db.Where("a.location_id = ?", f.LocationID)
	}
	if f.AssignedTo != nil {
		db = db.Where("a.assigned_to = ?", f.AssignedTo)
	}
	if f.Brand != nil {
		db = db.Where("a.brand ILIKE ?", "%"+*f.Brand+"%")
	}
	if f.Model != nil {
		db = db.Where("a.model ILIKE ?", "%"+*f.Model+"%")
	}
	return db
}

func (r *AssetRepository) applyAssetSorts(db *gorm.DB, sort *query.SortOptions) *gorm.DB {
	if sort == nil || sort.Field == "" {
		return db.Order("a.created_at DESC")
	}
	var orderClause string
	switch strings.ToLower(sort.Field) {
	case "asset_tag", "asset_name", "brand", "model", "serial_number", "purchase_date", "purchase_price", "vendor_name", "warranty_end", "status", "condition_status", "created_at", "updated_at":
		orderClause = "a." + sort.Field
	default:
		return db.Order("a.created_at DESC")
	}

	order := "DESC"
	if strings.ToLower(sort.Order) == "asc" {
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

func (r *AssetRepository) UpdateAsset(ctx context.Context, payload *domain.Asset) (domain.Asset, error) {
	var updatedAsset model.Asset
	assetID := payload.ID

	// Update asset in database
	assetUpdates := model.Asset{
		AssetTag:           payload.AssetTag,
		DataMatrixValue:    payload.DataMatrixValue,
		DataMatrixImageUrl: payload.DataMatrixImageUrl,
		AssetName:          payload.AssetName,
		Brand:              payload.Brand,
		Model:              payload.Model,
		SerialNumber:       payload.SerialNumber,
		PurchaseDate:       payload.PurchaseDate,
		PurchasePrice:      payload.PurchasePrice,
		VendorName:         payload.VendorName,
		WarrantyEnd:        payload.WarrantyEnd,
		Status:             payload.Status,
		Condition:          payload.Condition,
	}

	// Handle foreign keys
	if payload.CategoryID != "" {
		if parsedCategoryID, err := ulid.Parse(payload.CategoryID); err == nil {
			assetUpdates.CategoryID = model.SQLULID(parsedCategoryID)
		}
	}

	if payload.LocationID != nil && *payload.LocationID != "" {
		if parsedLocationID, err := ulid.Parse(*payload.LocationID); err == nil {
			modelULID := model.SQLULID(parsedLocationID)
			assetUpdates.LocationID = &modelULID
		}
	}

	if payload.AssignedTo != nil && *payload.AssignedTo != "" {
		if parsedAssignedTo, err := ulid.Parse(*payload.AssignedTo); err == nil {
			modelULID := model.SQLULID(parsedAssignedTo)
			assetUpdates.AssignedTo = &modelULID
		}
	}

	err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("id = ?", assetID).Updates(assetUpdates).Error
	if err != nil {
		return domain.Asset{}, domain.ErrInternal(err)
	}

	// Get updated asset
	err = r.db.WithContext(ctx).Preload("Category").Preload("Category.Translations").Preload("Location").Preload("Location.Translations").Preload("User").First(&updatedAsset, "id = ?", assetID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Asset{}, domain.ErrNotFound("asset")
		}
		return domain.Asset{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainAsset(&updatedAsset), nil
}

func (r *AssetRepository) UpdateAssetWithPayload(ctx context.Context, assetId string, payload *domain.UpdateAssetPayload) (domain.Asset, error) {
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
func (r *AssetRepository) GetAssetsPaginated(ctx context.Context, params query.Params) ([]domain.Asset, error) {
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

	// * Set pagination ke nil agar query.Apply tidak memproses cursor
	params.Pagination.Cursor = ""
	db = query.Apply(db, params, r.applyAssetFilters, r.applyAssetSorts)

	if err := db.Find(&assets).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain assets
	domainAssets := make([]domain.Asset, len(assets))
	for i, asset := range assets {
		domainAssets[i] = mapper.ToDomainAsset(&asset)
	}
	return domainAssets, nil
}

func (r *AssetRepository) GetAssetsCursor(ctx context.Context, params query.Params) ([]domain.Asset, error) {
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

	// * Set offset ke 0 agar query.Apply tidak memproses offset
	params.Pagination.Offset = 0
	db = query.Apply(db, params, r.applyAssetFilters, r.applyAssetSorts)

	if err := db.Find(&assets).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain assets
	domainAssets := make([]domain.Asset, len(assets))
	for i, asset := range assets {
		domainAssets[i] = mapper.ToDomainAsset(&asset)
	}
	return domainAssets, nil
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

func (r *AssetRepository) GetAssetByDataMatrixValue(ctx context.Context, dataMatrixValue string) (domain.Asset, error) {
	var asset model.Asset

	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Category.Translations").
		Preload("Location").
		Preload("Location.Translations").
		Preload("User").
		Where("data_matrix_value = ?", dataMatrixValue).First(&asset).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Asset{}, domain.ErrNotFound("asset with data matrix value '" + dataMatrixValue + "'")
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

func (r *AssetRepository) CheckDataMatrixValueExists(ctx context.Context, dataMatrixValue string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("data_matrix_value = ?", dataMatrixValue).Count(&count).Error
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

func (r *AssetRepository) GetAssetsResponse(ctx context.Context, params query.Params, langCode string) ([]domain.AssetResponse, error) {
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

	// * Set pagination ke nil agar query.Apply tidak memproses cursor
	params.Pagination.Cursor = ""
	db = query.Apply(db, params, r.applyAssetFilters, r.applyAssetSorts)

	if err := db.Find(&assets).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to AssetResponse using mapper
	assetResponses := mapper.ToDomainAssetsResponse(assets, langCode)
	return assetResponses, nil
}

func (r *AssetRepository) GetAssetResponseById(ctx context.Context, assetId string, langCode string) (domain.AssetResponse, error) {
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
			return domain.AssetResponse{}, domain.ErrNotFound("asset")
		}
		return domain.AssetResponse{}, domain.ErrInternal(err)
	}

	// Convert to AssetResponse using mapper
	return mapper.ToDomainAssetResponse(&asset, langCode), nil
}

func (r *AssetRepository) CheckAssetTagExistsExcluding(ctx context.Context, assetTag string, excludeAssetId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("asset_tag = ? AND id != ?", assetTag, excludeAssetId).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *AssetRepository) CheckDataMatrixValueExistsExcluding(ctx context.Context, dataMatrixValue string, excludeAssetId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Asset{}).Where("data_matrix_value = ? AND id != ?", dataMatrixValue, excludeAssetId).Count(&count).Error
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

func (r *AssetRepository) CountAssets(ctx context.Context, params query.Params) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Table("assets a")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("a.asset_tag ILIKE ? OR a.asset_name ILIKE ? OR a.brand ILIKE ? OR a.model ILIKE ? OR a.serial_number ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	db = query.Apply(db, query.Params{Filters: params.Filters}, r.applyAssetFilters, nil)

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
		Date  string `json:"date"`
		Count int64  `json:"count"`
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

	stats.Summary.EarliestCreationDate = earliestDate.Format("2006-01-02")
	stats.Summary.LatestCreationDate = latestDate.Format("2006-01-02")

	// Calculate average assets per day
	if !earliestDate.IsZero() && !latestDate.IsZero() {
		daysDiff := latestDate.Sub(earliestDate).Hours() / 24
		if daysDiff > 0 {
			stats.Summary.AverageAssetsPerDay = float64(totalCount) / daysDiff
		}
	}

	return stats, nil
}
