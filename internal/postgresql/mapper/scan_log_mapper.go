package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/oklog/ulid/v2"
)

// ===== Scan Log Mappers =====

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

// Domain -> Response conversions (for service layer)
func ScanLogToResponse(d *domain.ScanLog) domain.ScanLogResponse {
	return domain.ScanLogResponse{
		ID:              d.ID,
		AssetID:         d.AssetID,
		ScannedValue:    d.ScannedValue,
		ScanMethod:      d.ScanMethod,
		ScannedByID:     d.ScannedBy,
		ScanTimestamp:   d.ScanTimestamp.Format(TimeFormat),
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
