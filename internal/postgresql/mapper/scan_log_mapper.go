package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/oklog/ulid/v2"
)

// *==================== Model conversions ====================

func ToModelScanLog(d *domain.ScanLog) model.ScanLog {
	modelScanLog := model.ScanLog{
		ScannedValue:    d.ScannedValue,
		ScanMethod:      d.ScanMethod,
		ScanTimestamp:   d.ScanTimestamp,
		ScanLocationLat: d.ScanLocationLat,
		ScanLocationLng: d.ScanLocationLng,
		ScanResult:      d.ScanResult,
	}

	if d.ID != "" {
		if parsedID, err := ulid.Parse(d.ID); err == nil {
			modelScanLog.ID = model.SQLULID(parsedID)
		}
	}

	if d.AssetID != nil && *d.AssetID != "" {
		if parsedAssetID, err := ulid.Parse(*d.AssetID); err == nil {
			modelULID := model.SQLULID(parsedAssetID)
			modelScanLog.AssetID = &modelULID
		}
	}

	if d.ScannedBy != "" {
		if parsedScannedBy, err := ulid.Parse(d.ScannedBy); err == nil {
			modelScanLog.ScannedBy = model.SQLULID(parsedScannedBy)
		}
	}

	return modelScanLog
}

func ToModelScanLogForCreate(d *domain.ScanLog) model.ScanLog {
	modelScanLog := model.ScanLog{
		ScannedValue:    d.ScannedValue,
		ScanMethod:      d.ScanMethod,
		ScanTimestamp:   d.ScanTimestamp,
		ScanLocationLat: d.ScanLocationLat,
		ScanLocationLng: d.ScanLocationLng,
		ScanResult:      d.ScanResult,
	}

	if d.AssetID != nil && *d.AssetID != "" {
		if parsedAssetID, err := ulid.Parse(*d.AssetID); err == nil {
			modelULID := model.SQLULID(parsedAssetID)
			modelScanLog.AssetID = &modelULID
		}
	}

	if d.ScannedBy != "" {
		if parsedScannedBy, err := ulid.Parse(d.ScannedBy); err == nil {
			modelScanLog.ScannedBy = model.SQLULID(parsedScannedBy)
		}
	}

	return modelScanLog
}

// *==================== Entity conversions ====================
func ToDomainScanLog(m *model.ScanLog) domain.ScanLog {
	scanLog := domain.ScanLog{
		ID:              m.ID.String(),
		ScannedValue:    m.ScannedValue,
		ScanMethod:      m.ScanMethod,
		ScannedBy:       m.ScannedBy.String(),
		ScanTimestamp:   m.ScanTimestamp,
		ScanLocationLat: m.ScanLocationLat,
		ScanLocationLng: m.ScanLocationLng,
		ScanResult:      m.ScanResult,
	}

	if m.AssetID != nil && !m.AssetID.IsZero() {
		assetIDStr := m.AssetID.String()
		scanLog.AssetID = &assetIDStr
	}

	return scanLog
}

func ToDomainScanLogs(models []model.ScanLog) []domain.ScanLog {
	scanLogs := make([]domain.ScanLog, len(models))
	for i, m := range models {
		scanLogs[i] = ToDomainScanLog(&m)
	}
	return scanLogs
}

// *==================== Entity Response conversions ====================
func ScanLogToResponse(d *domain.ScanLog) domain.ScanLogResponse {
	return domain.ScanLogResponse{
		ID:              d.ID,
		AssetID:         d.AssetID,
		ScannedValue:    d.ScannedValue,
		ScanMethod:      d.ScanMethod,
		ScannedByID:     d.ScannedBy,
		ScanTimestamp:   d.ScanTimestamp,
		ScanLocationLat: d.ScanLocationLat,
		ScanLocationLng: d.ScanLocationLng,
		ScanResult:      d.ScanResult,
	}
}

func ScanLogsToResponses(scanLogs []domain.ScanLog) []domain.ScanLogResponse {
	responses := make([]domain.ScanLogResponse, len(scanLogs))
	for i, scanLog := range scanLogs {
		responses[i] = ScanLogToResponse(&scanLog)
	}
	return responses
}

func ScanLogToListResponse(d *domain.ScanLog) domain.ScanLogListResponse {
	return domain.ScanLogListResponse{
		ID:              d.ID,
		AssetID:         d.AssetID,
		ScannedValue:    d.ScannedValue,
		ScanMethod:      d.ScanMethod,
		ScannedByID:     d.ScannedBy,
		ScanTimestamp:   d.ScanTimestamp,
		ScanLocationLat: d.ScanLocationLat,
		ScanLocationLng: d.ScanLocationLng,
		ScanResult:      d.ScanResult,
	}
}

func ScanLogsToListResponses(scanLogs []domain.ScanLog) []domain.ScanLogListResponse {
	responses := make([]domain.ScanLogListResponse, len(scanLogs))
	for i, scanLog := range scanLogs {
		responses[i] = ScanLogToListResponse(&scanLog)
	}
	return responses
}

// *==================== Statistics conversions ====================
func ScanLogStatisticsToResponse(stats *domain.ScanLogStatistics) domain.ScanLogStatisticsResponse {
	response := domain.ScanLogStatisticsResponse{
		Total: domain.ScanLogCountStatisticsResponse{
			Count: stats.Total.Count,
		},
		Geographic: domain.ScanGeographicStatisticsResponse{
			WithCoordinates:    stats.Geographic.WithCoordinates,
			WithoutCoordinates: stats.Geographic.WithoutCoordinates,
		},
		Summary: domain.ScanLogSummaryStatisticsResponse{
			TotalScans:            stats.Summary.TotalScans,
			SuccessRate:           domain.NewDecimal2(stats.Summary.SuccessRate),
			ScansWithCoordinates:  stats.Summary.ScansWithCoordinates,
			CoordinatesPercentage: domain.NewDecimal2(stats.Summary.CoordinatesPercentage),
			AverageScansPerDay:    domain.NewDecimal2(stats.Summary.AverageScansPerDay),
			LatestScanDate:        stats.Summary.LatestScanDate,
			EarliestScanDate:      stats.Summary.EarliestScanDate,
		},
	}

	// Convert method statistics
	response.ByMethod = make([]domain.ScanMethodStatisticsResponse, len(stats.ByMethod))
	for i, method := range stats.ByMethod {
		response.ByMethod[i] = domain.ScanMethodStatisticsResponse{
			Method: method.Method,
			Count:  method.Count,
		}
	}

	// Convert result statistics
	response.ByResult = make([]domain.ScanResultStatisticsResponse, len(stats.ByResult))
	for i, result := range stats.ByResult {
		response.ByResult[i] = domain.ScanResultStatisticsResponse{
			Result: result.Result,
			Count:  result.Count,
		}
	}

	// Convert scan trends
	response.ScanTrends = make([]domain.ScanTrendResponse, len(stats.ScanTrends))
	for i, trend := range stats.ScanTrends {
		response.ScanTrends[i] = domain.ScanTrendResponse{
			Date:  trend.Date,
			Count: trend.Count,
		}
	}

	// Convert top scanners
	response.TopScanners = make([]domain.ScannerStatisticsResponse, len(stats.TopScanners))
	for i, scanner := range stats.TopScanners {
		response.TopScanners[i] = domain.ScannerStatisticsResponse{
			UserID: scanner.UserID,
			Count:  scanner.Count,
		}
	}

	return response
}

func MapScanLogSortFieldToColumn(field domain.ScanLogSortField) string {
	columnMap := map[domain.ScanLogSortField]string{
		domain.ScanLogSortByScannedValue:  "scanned_value",
		domain.ScanLogSortByScanMethod:    "scan_method",
		domain.ScanLogSortByScanTimestamp: "scan_timestamp",
		domain.ScanLogSortByScanResult:    "scan_result",
	}

	if column, exists := columnMap[field]; exists {
		return column
	}
	return "scan_timestamp"
}
