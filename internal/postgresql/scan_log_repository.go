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
	"gorm.io/gorm"
)

type ScanLogRepository struct {
	db *gorm.DB
}

type ScanLogFilterOptions struct {
	ScanMethod     *domain.ScanMethodType `json:"scanMethod,omitempty"`
	ScanResult     *domain.ScanResultType `json:"scanResult,omitempty"`
	ScannedBy      *string                `json:"scannedBy,omitempty"`
	AssetID        *string                `json:"assetId,omitempty"`
	DateFrom       *time.Time             `json:"dateFrom,omitempty"`
	DateTo         *time.Time             `json:"dateTo,omitempty"`
	HasCoordinates *bool                  `json:"hasCoordinates,omitempty"`
}

func NewScanLogRepository(db *gorm.DB) *ScanLogRepository {
	return &ScanLogRepository{
		db: db,
	}
}

func (r *ScanLogRepository) applyScanLogFilters(db *gorm.DB, filters any) *gorm.DB {
	f, ok := filters.(*ScanLogFilterOptions)
	if !ok || f == nil {
		return db
	}

	if f.ScanMethod != nil {
		db = db.Where("sl.scan_method = ?", *f.ScanMethod)
	}

	if f.ScanResult != nil {
		db = db.Where("sl.scan_result = ?", *f.ScanResult)
	}

	if f.ScannedBy != nil {
		db = db.Where("sl.scanned_by = ?", *f.ScannedBy)
	}

	if f.AssetID != nil {
		db = db.Where("sl.asset_id = ?", *f.AssetID)
	}

	if f.DateFrom != nil {
		db = db.Where("sl.scan_timestamp >= ?", *f.DateFrom)
	}

	if f.DateTo != nil {
		db = db.Where("sl.scan_timestamp <= ?", *f.DateTo)
	}

	if f.HasCoordinates != nil {
		if *f.HasCoordinates {
			db = db.Where("sl.scan_location_lat IS NOT NULL AND sl.scan_location_lng IS NOT NULL")
		} else {
			db = db.Where("sl.scan_location_lat IS NULL OR sl.scan_location_lng IS NULL")
		}
	}

	return db
}

func (r *ScanLogRepository) applyScanLogSorts(db *gorm.DB, sort *query.SortOptions) *gorm.DB {
	if sort == nil || sort.Field == "" {
		return db.Order("sl.scan_timestamp DESC")
	}

	var orderClause string
	switch strings.ToLower(sort.Field) {
	case "scan_timestamp", "scanned_value", "scan_method", "scan_result":
		orderClause = "sl." + sort.Field
	default:
		return db.Order("sl.scan_timestamp DESC")
	}

	order := "DESC"
	if strings.ToLower(sort.Order) == "asc" {
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

	// Fetch created scan log
	return r.GetScanLogById(ctx, modelScanLog.ID.String())
}

func (r *ScanLogRepository) DeleteScanLog(ctx context.Context, scanLogId string) error {
	if err := r.db.WithContext(ctx).Delete(&model.ScanLog{}, "id = ?", scanLogId).Error; err != nil {
		return domain.ErrInternal(err)
	}
	return nil
}

// *===========================QUERY===========================*
func (r *ScanLogRepository) GetScanLogsPaginated(ctx context.Context, params query.Params) ([]domain.ScanLog, error) {
	var scanLogs []model.ScanLog
	db := r.db.WithContext(ctx).
		Table("scan_logs sl").
		Preload("Asset").
		Preload("ScannedByUser")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("sl.scanned_value ILIKE ?", searchPattern)
	}

	// Set pagination cursor to empty for offset-based pagination
	params.Pagination.Cursor = ""
	db = query.Apply(db, params, r.applyScanLogFilters, r.applyScanLogSorts)

	if err := db.Find(&scanLogs).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain scan logs
	return mapper.ToDomainScanLogs(scanLogs), nil
}

func (r *ScanLogRepository) GetScanLogsCursor(ctx context.Context, params query.Params) ([]domain.ScanLog, error) {
	var scanLogs []model.ScanLog
	db := r.db.WithContext(ctx).
		Table("scan_logs sl").
		Preload("Asset").
		Preload("ScannedByUser")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("sl.scanned_value ILIKE ?", searchPattern)
	}

	// Set offset to 0 for cursor-based pagination
	params.Pagination.Offset = 0
	db = query.Apply(db, params, r.applyScanLogFilters, r.applyScanLogSorts)

	if err := db.Find(&scanLogs).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain scan logs
	return mapper.ToDomainScanLogs(scanLogs), nil
}

func (r *ScanLogRepository) GetScanLogById(ctx context.Context, scanLogId string) (domain.ScanLog, error) {
	var scanLog model.ScanLog

	err := r.db.WithContext(ctx).
		Preload("Asset").
		Preload("ScannedByUser").
		First(&scanLog, "id = ?", scanLogId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ScanLog{}, domain.ErrNotFound("scan log")
		}
		return domain.ScanLog{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainScanLog(&scanLog), nil
}

func (r *ScanLogRepository) GetScanLogsByAssetId(ctx context.Context, assetId string, params query.Params) ([]domain.ScanLog, error) {
	var scanLogs []model.ScanLog
	db := r.db.WithContext(ctx).
		Table("scan_logs sl").
		Where("sl.asset_id = ?", assetId).
		Preload("Asset").
		Preload("ScannedByUser")

	db = query.Apply(db, params, r.applyScanLogFilters, r.applyScanLogSorts)

	if err := db.Find(&scanLogs).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	return mapper.ToDomainScanLogs(scanLogs), nil
}

func (r *ScanLogRepository) GetScanLogsByUserId(ctx context.Context, userId string, params query.Params) ([]domain.ScanLog, error) {
	var scanLogs []model.ScanLog
	db := r.db.WithContext(ctx).
		Table("scan_logs sl").
		Where("sl.scanned_by = ?", userId).
		Preload("Asset").
		Preload("ScannedByUser")

	db = query.Apply(db, params, r.applyScanLogFilters, r.applyScanLogSorts)

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

func (r *ScanLogRepository) CountScanLogs(ctx context.Context, params query.Params) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Table("scan_logs sl")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("sl.scanned_value ILIKE ?", searchPattern)
	}

	db = query.Apply(db, query.Params{Filters: params.Filters}, r.applyScanLogFilters, nil)

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
		Date  string `json:"date"`
		Count int64  `json:"count"`
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
		// Parse the date string to time.Time
		date, _ := time.Parse("2006-01-02", st.Date)
		stats.ScanTrends[i] = domain.ScanTrend{
			Date:  date,
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
	var earliestDate, latestDate time.Time
	if err := r.db.WithContext(ctx).Model(&model.ScanLog{}).Select("MIN(scan_timestamp)").Scan(&earliestDate).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.ScanLog{}).Select("MAX(scan_timestamp)").Scan(&latestDate).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	if !earliestDate.IsZero() && !latestDate.IsZero() {
		stats.Summary.EarliestScanDate = earliestDate
		stats.Summary.LatestScanDate = latestDate

		// Calculate average scans per day
		daysDiff := latestDate.Sub(earliestDate).Hours() / 24
		if daysDiff > 0 {
			stats.Summary.AverageScansPerDay = float64(totalCount) / daysDiff
		}
	}

	return stats, nil
}
