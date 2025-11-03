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

type MaintenanceRecordRepository struct {
	db *gorm.DB
}

func NewMaintenanceRecordRepository(db *gorm.DB) *MaintenanceRecordRepository {
	return &MaintenanceRecordRepository{db: db}
}

// ===== Filters and Sorts =====

func (r *MaintenanceRecordRepository) applyRecordFilters(db *gorm.DB, filters *domain.MaintenanceRecordFilterOptions) *gorm.DB {
	if filters == nil {
		return db
	}
	if filters.AssetID != nil && *filters.AssetID != "" {
		db = db.Where("maintenance_records.asset_id = ?", filters.AssetID)
	}
	if filters.ScheduleID != nil && *filters.ScheduleID != "" {
		db = db.Where("maintenance_records.schedule_id = ?", filters.ScheduleID)
	}
	if filters.PerformedByUser != nil && *filters.PerformedByUser != "" {
		db = db.Where("maintenance_records.performed_by_user = ?", filters.PerformedByUser)
	}
	if filters.VendorName != nil && *filters.VendorName != "" {
		db = db.Where("maintenance_records.performed_by_vendor ILIKE ?", "%"+*filters.VendorName+"%")
	}
	if filters.FromDate != nil && *filters.FromDate != "" {
		db = db.Where("maintenance_records.maintenance_date >= ?", *filters.FromDate)
	}
	if filters.ToDate != nil && *filters.ToDate != "" {
		db = db.Where("maintenance_records.maintenance_date <= ?", *filters.ToDate)
	}
	return db
}

func (r *MaintenanceRecordRepository) applyRecordSorts(db *gorm.DB, sort *domain.MaintenanceRecordSortOptions) *gorm.DB {
	if sort == nil || sort.Field == "" {
		return db.Order("maintenance_records.created_at DESC")
	}

	// Map camelCase sort field to snake_case database column
	columnName := mapper.MapMaintenanceRecordSortFieldToColumn(sort.Field)
	orderClause := "maintenance_records." + columnName

	order := "DESC"
	if sort.Order == domain.SortOrderAsc {
		order = "ASC"
	}
	return db.Order(fmt.Sprintf("%s %s", orderClause, order))
}

// ===== MUTATIONS =====

func (r *MaintenanceRecordRepository) CreateRecord(ctx context.Context, payload *domain.MaintenanceRecord) (domain.MaintenanceRecord, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.MaintenanceRecord{}, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	m := mapper.ToModelMaintenanceRecordForCreate(payload)
	if err := tx.Create(&m).Error; err != nil {
		tx.Rollback()
		return domain.MaintenanceRecord{}, domain.ErrInternal(err)
	}
	for _, t := range payload.Translations {
		mt := mapper.ToModelMaintenanceRecordTranslationForCreate(m.ID.String(), &t)
		if err := tx.Create(&mt).Error; err != nil {
			tx.Rollback()
			return domain.MaintenanceRecord{}, domain.ErrInternal(err)
		}
	}
	if err := tx.Commit().Error; err != nil {
		return domain.MaintenanceRecord{}, domain.ErrInternal(err)
	}

	// Return created maintenance record with translations (no need to query again)
	// GORM has already filled the model with created data including ID and timestamps
	domainRecord := mapper.ToDomainMaintenanceRecord(&m)
	// Add translations manually since they were created separately
	for _, translation := range payload.Translations {
		domainRecord.Translations = append(domainRecord.Translations, domain.MaintenanceRecordTranslation{
			LangCode: translation.LangCode,
			Title:    translation.Title,
			Notes:    translation.Notes,
		})
	}
	return domainRecord, nil
}

func (r *MaintenanceRecordRepository) UpdateRecord(ctx context.Context, recordId string, payload *domain.UpdateMaintenanceRecordPayload) (domain.MaintenanceRecord, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.MaintenanceRecord{}, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// Update maintenance record basic info
	updates := mapper.ToModelMaintenanceRecordUpdateMap(payload)
	if len(updates) > 0 {
		if err := tx.Table("maintenance_records").Where("id = ?", recordId).Updates(updates).Error; err != nil {
			tx.Rollback()
			return domain.MaintenanceRecord{}, domain.ErrInternal(err)
		}
	}

	// Update translations if provided
	if len(payload.Translations) > 0 {
		// Delete existing translations
		if err := tx.Where("record_id = ?", recordId).Delete(&model.MaintenanceRecordTranslation{}).Error; err != nil {
			tx.Rollback()
			return domain.MaintenanceRecord{}, domain.ErrInternal(err)
		}

		// Create new translations
		for _, translationPayload := range payload.Translations {
			translation := domain.MaintenanceRecordTranslation{
				LangCode: translationPayload.LangCode,
				Title:    *translationPayload.Title,
				Notes:    translationPayload.Notes,
			}
			modelTranslation := mapper.ToModelMaintenanceRecordTranslationForCreate(recordId, &translation)
			if err := tx.Create(&modelTranslation).Error; err != nil {
				tx.Rollback()
				return domain.MaintenanceRecord{}, domain.ErrInternal(err)
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return domain.MaintenanceRecord{}, domain.ErrInternal(err)
	}
	return r.GetRecordById(ctx, recordId)
}

func (r *MaintenanceRecordRepository) DeleteRecord(ctx context.Context, recordId string) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.ErrInternal(tx.Error)
	}
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Delete(&model.MaintenanceRecordTranslation{}, "record_id = ?", recordId).Error; err != nil {
		tx.Rollback()
		return domain.ErrInternal(err)
	}
	if err := tx.Delete(&model.MaintenanceRecord{}, "id = ?", recordId).Error; err != nil {
		tx.Rollback()
		return domain.ErrInternal(err)
	}
	if err := tx.Commit().Error; err != nil {
		return domain.ErrInternal(err)
	}
	return nil
}

// ===== QUERIES =====

func (r *MaintenanceRecordRepository) GetRecordsPaginated(ctx context.Context, params domain.MaintenanceRecordParams, langCode string) ([]domain.MaintenanceRecord, error) {
	var records []model.MaintenanceRecord
	db := r.db.WithContext(ctx).
		Preload("Translations").
		Preload("Schedule").
		Preload("Schedule.Translations").
		Preload("Asset").
		Preload("Asset.Category").
		Preload("Asset.Category.Translations").
		Preload("Asset.Location").
		Preload("Asset.Location.Translations").
		Preload("Asset.User").
		Preload("User")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		sq := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN maintenance_record_translations mrt ON maintenance_records.id = mrt.record_id").
			Where("mrt.title ILIKE ?", sq).
			Group("maintenance_records.id")
	}

	// Apply filters, sorts, and pagination manually
	db = r.applyRecordFilters(db, params.Filters)
	db = r.applyRecordSorts(db, params.Sort)
	if params.Pagination != nil {
		if params.Pagination.Limit > 0 {
			db = db.Limit(params.Pagination.Limit)
		}
		if params.Pagination.Offset > 0 {
			db = db.Offset(params.Pagination.Offset)
		}
	}

	if err := db.Find(&records).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}
	return mapper.ToDomainMaintenanceRecords(records), nil
}

func (r *MaintenanceRecordRepository) GetRecordsCursor(ctx context.Context, params domain.MaintenanceRecordParams, langCode string) ([]domain.MaintenanceRecord, error) {
	var records []model.MaintenanceRecord
	db := r.db.WithContext(ctx).
		Preload("Translations").
		Preload("Schedule").
		Preload("Schedule.Translations").
		Preload("Asset").
		Preload("Asset.Category").
		Preload("Asset.Category.Translations").
		Preload("Asset.Location").
		Preload("Asset.Location.Translations").
		Preload("Asset.User").
		Preload("User")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		sq := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN maintenance_record_translations mrt ON maintenance_records.id = mrt.record_id").
			Where("mrt.title ILIKE ?", sq).
			Group("maintenance_records.id")
	}

	// Apply filters, sorts, and cursor pagination manually
	db = r.applyRecordFilters(db, params.Filters)
	db = r.applyRecordSorts(db, params.Sort)
	if params.Pagination != nil {
		if params.Pagination.Cursor != "" {
			db = db.Where("mr.id > ?", params.Pagination.Cursor)
		}
		if params.Pagination.Limit > 0 {
			db = db.Limit(params.Pagination.Limit)
		}
	}

	if err := db.Find(&records).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}
	return mapper.ToDomainMaintenanceRecords(records), nil
}

func (r *MaintenanceRecordRepository) GetRecordById(ctx context.Context, recordId string) (domain.MaintenanceRecord, error) {
	var m model.MaintenanceRecord
	err := r.db.WithContext(ctx).
		Preload("Translations").
		Preload("Schedule").
		Preload("Schedule.Translations").
		Preload("Asset").
		Preload("Asset.Category").
		Preload("Asset.Category.Translations").
		Preload("Asset.Location").
		Preload("Asset.Location.Translations").
		Preload("Asset.User").
		Preload("User").
		First(&m, "id = ?", recordId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.MaintenanceRecord{}, domain.ErrNotFound("maintenance_record")
		}
		return domain.MaintenanceRecord{}, domain.ErrInternal(err)
	}
	return mapper.ToDomainMaintenanceRecord(&m), nil
}

// ===== Checks and Counts =====

func (r *MaintenanceRecordRepository) CheckRecordExist(ctx context.Context, recordId string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceRecord{}).Where("id = ?", recordId).Count(&count).Error; err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *MaintenanceRecordRepository) CountRecords(ctx context.Context, params domain.MaintenanceRecordParams) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(&model.MaintenanceRecord{})
	if params.SearchQuery != nil && *params.SearchQuery != "" {
		sq := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN maintenance_record_translations mrt ON maintenance_records.id = mrt.record_id").
			Where("mrt.title ILIKE ?", sq).
			Group("maintenance_records.id")
	}
	db = r.applyRecordFilters(db, params.Filters)
	if err := db.Count(&count).Error; err != nil {
		return 0, domain.ErrInternal(err)
	}
	return count, nil
}

// ===== Statistics =====

func (r *MaintenanceRecordRepository) GetMaintenanceRecordStatistics(ctx context.Context) (domain.MaintenanceRecordStatistics, error) {
	var stats domain.MaintenanceRecordStatistics

	// Total records
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceRecord{}).Count(&total).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.Total.Count = int(total)

	// By performer
	var byPerformerResults []struct {
		UserID      string
		UserName    string
		UserEmail   string
		Count       int64
		TotalCost   float64
		AverageCost float64
	}
	if err := r.db.WithContext(ctx).
		Table("maintenance_records mr").
		Select("u.id as user_id, u.full_name as user_name, u.email as user_email, " +
			"COUNT(*) as count, COALESCE(SUM(mr.actual_cost), 0) as total_cost, " +
			"COALESCE(AVG(mr.actual_cost), 0) as average_cost").
		Joins("LEFT JOIN users u ON mr.performed_by_user = u.id").
		Where("mr.performed_by_user IS NOT NULL").
		Group("u.id, u.full_name, u.email").
		Scan(&byPerformerResults).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	for _, result := range byPerformerResults {
		stats.ByPerformer = append(stats.ByPerformer, domain.UserMaintenanceRecordStatistics{
			UserID:      result.UserID,
			UserName:    result.UserName,
			UserEmail:   result.UserEmail,
			Count:       int(result.Count),
			TotalCost:   result.TotalCost,
			AverageCost: result.AverageCost,
		})
	}

	// By vendor
	var byVendorResults []struct {
		VendorName  string
		Count       int64
		TotalCost   float64
		AverageCost float64
	}
	if err := r.db.WithContext(ctx).
		Table("maintenance_records").
		Select("performed_by_vendor as vendor_name, COUNT(*) as count, " +
			"COALESCE(SUM(actual_cost), 0) as total_cost, " +
			"COALESCE(AVG(actual_cost), 0) as average_cost").
		Where("performed_by_vendor IS NOT NULL AND performed_by_vendor != ''").
		Group("performed_by_vendor").
		Scan(&byVendorResults).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	for _, result := range byVendorResults {
		stats.ByVendor = append(stats.ByVendor, domain.VendorMaintenanceRecordStatistics{
			VendorName:  result.VendorName,
			Count:       int(result.Count),
			TotalCost:   result.TotalCost,
			AverageCost: result.AverageCost,
		})
	}

	// By asset
	var byAssetResults []struct {
		AssetID         string
		AssetName       string
		AssetTag        string
		RecordCount     int64
		LastMaintenance string
		TotalCost       float64
		AverageCost     float64
	}
	if err := r.db.WithContext(ctx).
		Table("maintenance_records mr").
		Select("a.id as asset_id, a.asset_name, a.asset_tag, COUNT(*) as record_count, " +
			"MAX(mr.maintenance_date)::text as last_maintenance, " +
			"COALESCE(SUM(mr.actual_cost), 0) as total_cost, " +
			"COALESCE(AVG(mr.actual_cost), 0) as average_cost").
		Joins("LEFT JOIN assets a ON mr.asset_id = a.id").
		Group("a.id, a.asset_name, a.asset_tag").
		Scan(&byAssetResults).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	for _, result := range byAssetResults {
		stats.ByAsset = append(stats.ByAsset, domain.AssetMaintenanceRecordStatistics{
			AssetID:         result.AssetID,
			AssetName:       result.AssetName,
			AssetTag:        result.AssetTag,
			RecordCount:     int(result.RecordCount),
			LastMaintenance: result.LastMaintenance,
			TotalCost:       result.TotalCost,
			AverageCost:     result.AverageCost,
		})
	}

	// Cost statistics
	var costStats struct {
		TotalCost   *float64
		AverageCost *float64
		MinCost     *float64
		MaxCost     *float64
	}
	var recordsWithCost, recordsWithoutCost int64

	if err := r.db.WithContext(ctx).Model(&model.MaintenanceRecord{}).
		Select("SUM(actual_cost) as total_cost, AVG(actual_cost) as average_cost, MIN(actual_cost) as min_cost, MAX(actual_cost) as max_cost").
		Where("actual_cost IS NOT NULL").
		Scan(&costStats).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	if err := r.db.WithContext(ctx).Model(&model.MaintenanceRecord{}).
		Where("actual_cost IS NOT NULL").Count(&recordsWithCost).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	if err := r.db.WithContext(ctx).Model(&model.MaintenanceRecord{}).
		Where("actual_cost IS NULL").Count(&recordsWithoutCost).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.CostStatistics = domain.MaintenanceRecordCostStatistics{
		TotalCost:          costStats.TotalCost,
		AverageCost:        costStats.AverageCost,
		MinCost:            costStats.MinCost,
		MaxCost:            costStats.MaxCost,
		RecordsWithCost:    int(recordsWithCost),
		RecordsWithoutCost: int(recordsWithoutCost),
	}

	// Completion trend (last 30 days)
	var completionTrends []struct {
		Date  time.Time
		Count int64
	}
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceRecord{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= NOW() - INTERVAL '30 days'").
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&completionTrends).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	for _, ct := range completionTrends {
		stats.CompletionTrend = append(stats.CompletionTrend, domain.MaintenanceRecordCompletionTrend{
			Date:  ct.Date,
			Count: int(ct.Count),
		})
	}

	// Monthly trends (last 12 months)
	var monthlyTrends []struct {
		Month       string
		RecordCount int64
		TotalCost   float64
	}
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceRecord{}).
		Select("TO_CHAR(created_at, 'YYYY-MM') as month, COUNT(*) as record_count, " +
			"COALESCE(SUM(actual_cost), 0) as total_cost").
		Where("created_at >= NOW() - INTERVAL '12 months'").
		Group("TO_CHAR(created_at, 'YYYY-MM')").
		Order("month ASC").
		Scan(&monthlyTrends).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	for _, mt := range monthlyTrends {
		stats.MonthlyTrends = append(stats.MonthlyTrends, domain.MaintenanceRecordMonthlyTrend{
			Month:       mt.Month,
			RecordCount: int(mt.RecordCount),
			TotalCost:   mt.TotalCost,
		})
	}

	// Summary statistics
	stats.Summary.TotalRecords = int(total)
	stats.Summary.RecordsWithCostInfo = int(recordsWithCost)
	if total > 0 {
		stats.Summary.CostInfoPercentage = float64(recordsWithCost) / float64(total) * 100
	}

	// Unique vendors and performers
	var uniqueVendors int64
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceRecord{}).
		Select("COUNT(DISTINCT performed_by_vendor)").
		Where("performed_by_vendor IS NOT NULL AND performed_by_vendor != ''").
		Scan(&uniqueVendors).Error; err == nil {
		stats.Summary.TotalUniqueVendors = int(uniqueVendors)
	}

	var uniquePerformers int64
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceRecord{}).
		Select("COUNT(DISTINCT performed_by_user)").
		Where("performed_by_user IS NOT NULL").
		Scan(&uniquePerformers).Error; err == nil {
		stats.Summary.TotalUniquePerformers = int(uniquePerformers)
	}

	// Date ranges and averages
	var earliest, latest time.Time
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceRecord{}).Select("MIN(created_at)").Scan(&earliest).Error; err == nil {
		stats.Summary.EarliestRecordDate = earliest.Format("2006-01-02")
	}
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceRecord{}).Select("MAX(created_at)").Scan(&latest).Error; err == nil {
		stats.Summary.LatestRecordDate = latest.Format("2006-01-02")
	}
	if !earliest.IsZero() && !latest.IsZero() {
		days := latest.Sub(earliest).Hours() / 24
		if days > 0 {
			stats.Summary.AverageRecordsPerDay = float64(total) / days
		}
	}

	// Most and least expensive maintenance
	stats.Summary.MostExpensiveMaintenanceCost = costStats.MaxCost
	stats.Summary.LeastExpensiveMaintenanceCost = costStats.MinCost

	// Assets with maintenance
	var assetsWithMaintenance int64
	if err := r.db.WithContext(ctx).
		Table("assets a").
		Joins("INNER JOIN maintenance_records mr ON a.id = mr.asset_id").
		Select("COUNT(DISTINCT a.id)").
		Scan(&assetsWithMaintenance).Error; err == nil {
		stats.Summary.AssetsWithMaintenance = int(assetsWithMaintenance)
	}

	if assetsWithMaintenance > 0 {
		stats.Summary.AverageMaintenancePerAsset = float64(total) / float64(assetsWithMaintenance)
	}

	return stats, nil
}

func (r *MaintenanceRecordRepository) GetMaintenanceRecordsForExport(ctx context.Context, params domain.MaintenanceRecordParams, langCode string) ([]domain.MaintenanceRecord, error) {
	var records []model.MaintenanceRecord
	db := r.db.WithContext(ctx).
		Preload("Translations").
		Preload("Schedule").
		Preload("Schedule.Translations").
		Preload("Asset").
		Preload("Asset.Category").
		Preload("Asset.Category.Translations").
		Preload("Asset.Location").
		Preload("Asset.Location.Translations").
		Preload("Asset.User").
		Preload("User")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		sq := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN maintenance_record_translations mrt ON maintenance_records.id = mrt.record_id").
			Where("mrt.title ILIKE ?", sq).
			Group("maintenance_records.id")
	}

	// Apply filters
	db = r.applyRecordFilters(db, params.Filters)

	// Apply sorting
	db = r.applyRecordSorts(db, params.Sort)

	// No pagination for export - get all matching maintenance records
	if err := db.Find(&records).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	return mapper.ToDomainMaintenanceRecords(records), nil
}
