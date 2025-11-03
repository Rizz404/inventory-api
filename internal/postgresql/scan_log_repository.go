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

type ScanLogRepository struct {
	db *gorm.DB
}

func NewScanLogRepository(db *gorm.DB) *ScanLogRepository {
	return &ScanLogRepository{
		db: db,
	}
}

func (r *ScanLogRepository) applyScanLogFilters(db *gorm.DB, filters *domain.ScanLogFilterOptions) *gorm.DB {
	if filters == nil {
		return db
	}

	if filters.ScanMethod != nil {
		db = db.Where("sl.scan_method = ?", *filters.ScanMethod)
	}

	if filters.ScanResult != nil {
		db = db.Where("sl.scan_result = ?", *filters.ScanResult)
	}

	if filters.ScannedBy != nil {
		db = db.Where("sl.scanned_by = ?", *filters.ScannedBy)
	}

	if filters.AssetID != nil {
		db = db.Where("sl.asset_id = ?", *filters.AssetID)
	}

	if filters.DateFrom != nil {
		db = db.Where("sl.scan_timestamp >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		db = db.Where("sl.scan_timestamp <= ?", *filters.DateTo)
	}

	if filters.HasCoordinates != nil {
		if *filters.HasCoordinates {
			db = db.Where("sl.scan_location_lat IS NOT NULL AND sl.scan_location_lng IS NOT NULL")
		} else {
			db = db.Where("sl.scan_location_lat IS NULL OR sl.scan_location_lng IS NULL")
		}
	}

	return db
}

func (r *ScanLogRepository) applyScanLogSorts(db *gorm.DB, sort *domain.ScanLogSortOptions) *gorm.DB {
	if sort == nil || sort.Field == "" {
		return db.Order("sl.scan_timestamp DESC")
	}

	// Map camelCase sort field to snake_case database column
	columnName := mapper.MapScanLogSortFieldToColumn(sort.Field)
	orderClause := columnName

	order := "DESC"
	if sort.Order == domain.SortOrderAsc {
		order = "ASC"
	}
	return db.Order(fmt.Sprintf("%s %s", orderClause, order))
}

// *===========================MUTATION===========================*
func (r *ScanLogRepository) CreateScanLog(ctx context.Context, payload *domain.ScanLog) (domain.ScanLog, error) {
	modelScanLog := mapper.ToModelScanLogForCreate(payload)

	if err := r.db.WithContext(ctx).Create(&modelScanLog).Error; err != nil {
		return domain.ScanLog{}, domain.ErrInternal(err)
	}

	// Return created scan log (no need to query again)
	// GORM has already filled the model with created data including ID and timestamps
	return mapper.ToDomainScanLog(&modelScanLog), nil
}

func (r *ScanLogRepository) DeleteScanLog(ctx context.Context, scanLogId string) error {
	if err := r.db.WithContext(ctx).Delete(&model.ScanLog{}, "id = ?", scanLogId).Error; err != nil {
		return domain.ErrInternal(err)
	}
	return nil
}

// *===========================QUERY===========================*
func (r *ScanLogRepository) GetScanLogsPaginated(ctx context.Context, params domain.ScanLogParams) ([]domain.ScanLog, error) {
	var scanLogs []model.ScanLog
	db := r.db.WithContext(ctx).
		Table("scan_logs sl")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("sl.scanned_value ILIKE ?", searchPattern)
	}

	// Apply filters
	db = r.applyScanLogFilters(db, params.Filters)

	// Apply sorting
	db = r.applyScanLogSorts(db, params.Sort)

	// Apply pagination
	if params.Pagination != nil {
		if params.Pagination.Limit > 0 {
			db = db.Limit(params.Pagination.Limit)
		}
		if params.Pagination.Offset > 0 {
			db = db.Offset(params.Pagination.Offset)
		}
	}

	if err := db.Find(&scanLogs).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain scan logs
	return mapper.ToDomainScanLogs(scanLogs), nil
}

func (r *ScanLogRepository) GetScanLogsCursor(ctx context.Context, params domain.ScanLogParams) ([]domain.ScanLog, error) {
	var scanLogs []model.ScanLog
	db := r.db.WithContext(ctx).
		Table("scan_logs sl")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("sl.scanned_value ILIKE ?", searchPattern)
	}

	// Apply filters
	db = r.applyScanLogFilters(db, params.Filters)

	// Apply sorting
	db = r.applyScanLogSorts(db, params.Sort)

	// Apply cursor-based pagination
	if params.Pagination != nil {
		if params.Pagination.Cursor != "" {
			db = db.Where("sl.id > ?", params.Pagination.Cursor)
		}
		if params.Pagination.Limit > 0 {
			db = db.Limit(params.Pagination.Limit)
		}
	}

	if err := db.Find(&scanLogs).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain scan logs
	return mapper.ToDomainScanLogs(scanLogs), nil
}

func (r *ScanLogRepository) GetScanLogById(ctx context.Context, scanLogId string) (domain.ScanLog, error) {
	var scanLog model.ScanLog

	err := r.db.WithContext(ctx).
		First(&scanLog, "id = ?", scanLogId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ScanLog{}, domain.ErrNotFound("scan log")
		}
		return domain.ScanLog{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainScanLog(&scanLog), nil
}

func (r *ScanLogRepository) GetScanLogsByAssetId(ctx context.Context, assetId string, params domain.ScanLogParams) ([]domain.ScanLog, error) {
	var scanLogs []model.ScanLog
	db := r.db.WithContext(ctx).
		Table("scan_logs sl").
		Where("sl.asset_id = ?", assetId)

		// Apply filters
	db = r.applyScanLogFilters(db, params.Filters)

	// Apply sorting
	db = r.applyScanLogSorts(db, params.Sort)

	// Apply cursor-based pagination
	if params.Pagination != nil {
		if params.Pagination.Cursor != "" {
			db = db.Where("sl.id > ?", params.Pagination.Cursor)
		}
		if params.Pagination.Limit > 0 {
			db = db.Limit(params.Pagination.Limit)
		}
	}

	if err := db.Find(&scanLogs).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	return mapper.ToDomainScanLogs(scanLogs), nil
}

func (r *ScanLogRepository) GetScanLogsByUserId(ctx context.Context, userId string, params domain.ScanLogParams) ([]domain.ScanLog, error) {
	var scanLogs []model.ScanLog
	db := r.db.WithContext(ctx).
		Table("scan_logs sl").
		Where("sl.scanned_by = ?", userId)

	// Apply filters
	db = r.applyScanLogFilters(db, params.Filters)

	// Apply sorting
	db = r.applyScanLogSorts(db, params.Sort)

	// Apply cursor-based pagination
	if params.Pagination != nil {
		if params.Pagination.Cursor != "" {
			db = db.Where("sl.id > ?", params.Pagination.Cursor)
		}
		if params.Pagination.Limit > 0 {
			db = db.Limit(params.Pagination.Limit)
		}
	}

	if err := db.Find(&scanLogs).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	return mapper.ToDomainScanLogs(scanLogs), nil
}

func (r *ScanLogRepository) CheckScanLogExist(ctx context.Context, scanLogId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ScanLog{}).Where("id = ?", scanLogId).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *ScanLogRepository) CountScanLogs(ctx context.Context, params domain.ScanLogParams) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Table("scan_logs sl")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("sl.scanned_value ILIKE ?", searchPattern)
	}

	// Apply filters
	db = r.applyScanLogFilters(db, params.Filters)

	if err := db.Count(&count).Error; err != nil {
		return 0, domain.ErrInternal(err)
	}
	return count, nil
}

func (r *ScanLogRepository) GetScanLogStatistics(ctx context.Context) (domain.ScanLogStatistics, error) {
	var stats domain.ScanLogStatistics

	// Get total scan log count
	var totalCount int64
	if err := r.db.WithContext(ctx).Model(&model.ScanLog{}).Count(&totalCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.Total.Count = int(totalCount)

	// Get scan counts by method
	var methodStats []struct {
		ScanMethod string `json:"scanMethod"`
		Count      int64  `json:"count"`
	}
	if err := r.db.WithContext(ctx).Model(&model.ScanLog{}).
		Select("scan_method, COUNT(*) as count").
		Group("scan_method").
		Order("count DESC").
		Scan(&methodStats).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.ByMethod = make([]domain.ScanMethodStatistics, len(methodStats))
	for i, ms := range methodStats {
		stats.ByMethod[i] = domain.ScanMethodStatistics{
			Method: domain.ScanMethodType(ms.ScanMethod),
			Count:  int(ms.Count),
		}
	}

	// Get scan counts by result
	var resultStats []struct {
		ScanResult string `json:"scanResult"`
		Count      int64  `json:"count"`
	}
	if err := r.db.WithContext(ctx).Model(&model.ScanLog{}).
		Select("scan_result, COUNT(*) as count").
		Group("scan_result").
		Order("count DESC").
		Scan(&resultStats).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.ByResult = make([]domain.ScanResultStatistics, len(resultStats))
	for i, rs := range resultStats {
		stats.ByResult[i] = domain.ScanResultStatistics{
			Result: domain.ScanResultType(rs.ScanResult),
			Count:  int(rs.Count),
		}
	}

	// Get scan trends (last 30 days)
	var scanTrends []struct {
		Date  time.Time `json:"date"`
		Count int64     `json:"count"`
	}
	if err := r.db.WithContext(ctx).Model(&model.ScanLog{}).
		Select("DATE(scan_timestamp) as date, COUNT(*) as count").
		Where("scan_timestamp >= NOW() - INTERVAL '30 days'").
		Group("DATE(scan_timestamp)").
		Order("date ASC").
		Scan(&scanTrends).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.ScanTrends = make([]domain.ScanTrend, len(scanTrends))
	for i, st := range scanTrends {
		stats.ScanTrends[i] = domain.ScanTrend{
			Date:  st.Date,
			Count: int(st.Count),
		}
	}

	// Get geographic statistics
	var withCoordinates, withoutCoordinates int64
	if err := r.db.WithContext(ctx).Model(&model.ScanLog{}).
		Where("scan_location_lat IS NOT NULL AND scan_location_lng IS NOT NULL").
		Count(&withCoordinates).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.ScanLog{}).
		Where("scan_location_lat IS NULL OR scan_location_lng IS NULL").
		Count(&withoutCoordinates).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.Geographic.WithCoordinates = int(withCoordinates)
	stats.Geographic.WithoutCoordinates = int(withoutCoordinates)

	// Get summary statistics
	stats.Summary.TotalScans = int(totalCount)

	// Get success rate
	var successCount int64
	if err := r.db.WithContext(ctx).Model(&model.ScanLog{}).
		Where("scan_result = ?", domain.ScanResultSuccess).
		Count(&successCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	if totalCount > 0 {
		stats.Summary.SuccessRate = float64(successCount) / float64(totalCount) * 100
		stats.Summary.CoordinatesPercentage = float64(withCoordinates) / float64(totalCount) * 100
	}

	stats.Summary.ScansWithCoordinates = int(withCoordinates)

	// Get top scanners (most active users)
	var topScanners []struct {
		ScannedBy string `json:"scannedBy"`
		Count     int64  `json:"count"`
	}
	if err := r.db.WithContext(ctx).Model(&model.ScanLog{}).
		Select("scanned_by, COUNT(*) as count").
		Group("scanned_by").
		Order("count DESC").
		Limit(10).
		Scan(&topScanners).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.TopScanners = make([]domain.ScannerStatistics, len(topScanners))
	for i, ts := range topScanners {
		stats.TopScanners[i] = domain.ScannerStatistics{
			UserID: ts.ScannedBy,
			Count:  int(ts.Count),
		}
	}

	// Get earliest and latest scan dates
	var earliestDate, latestDate *time.Time
	if err := r.db.WithContext(ctx).Model(&model.ScanLog{}).Select("MIN(scan_timestamp)").Row().Scan(&earliestDate); err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.ScanLog{}).Select("MAX(scan_timestamp)").Row().Scan(&latestDate); err != nil {
		return stats, domain.ErrInternal(err)
	}

	if earliestDate != nil && latestDate != nil {
		stats.Summary.EarliestScanDate = *earliestDate
		stats.Summary.LatestScanDate = *latestDate

		// Calculate average scans per day
		daysDiff := latestDate.Sub(*earliestDate).Hours() / 24
		if daysDiff > 0 {
			stats.Summary.AverageScansPerDay = float64(totalCount) / daysDiff
		}
	}

	return stats, nil
}

func (r *ScanLogRepository) GetScanLogsForExport(ctx context.Context, params domain.ScanLogParams) ([]domain.ScanLog, error) {
	var scanLogs []model.ScanLog
	db := r.db.WithContext(ctx).
		Table("scan_logs sl")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("sl.scanned_value ILIKE ?", searchPattern)
	}

	// Apply filters
	db = r.applyScanLogFilters(db, params.Filters)

	// Apply sorting
	db = r.applyScanLogSorts(db, params.Sort)

	// No pagination for export - get all matching scan logs
	if err := db.Find(&scanLogs).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain scan logs
	return mapper.ToDomainScanLogs(scanLogs), nil
}
