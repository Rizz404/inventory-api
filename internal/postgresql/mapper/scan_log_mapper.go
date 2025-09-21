package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/oklog/ulid/v2"
)

// ===== Scan Log Mappers =====

func ToModelScanLog(d *domain.ScanLog) model.ScanLog {
	modelScanLog := model.ScanLog{
		ScanTime:    d.ScanTime,
		ScanMethod:  d.ScanMethod,
		ScanData:    d.ScanData,
		IsSucceeded: d.IsSucceeded,
	}

	if d.ID != "" {
		if parsedID, err := ulid.Parse(d.ID); err == nil {
			modelScanLog.ID = model.SQLULID(parsedID)
		}
	}

	if d.AssetID != "" {
		if parsedAssetID, err := ulid.Parse(d.AssetID); err == nil {
			modelScanLog.AssetID = model.SQLULID(parsedAssetID)
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
		ScanTime:    d.ScanTime,
		ScanMethod:  d.ScanMethod,
		ScanData:    d.ScanData,
		IsSucceeded: d.IsSucceeded,
	}

	if d.AssetID != "" {
		if parsedAssetID, err := ulid.Parse(d.AssetID); err == nil {
			modelScanLog.AssetID = model.SQLULID(parsedAssetID)
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
	return domain.ScanLog{
		ID:          m.ID.String(),
		AssetID:     m.AssetID.String(),
		ScannedBy:   m.ScannedBy.String(),
		ScanTime:    m.ScanTime,
		ScanMethod:  m.ScanMethod,
		ScanData:    m.ScanData,
		IsSucceeded: m.IsSucceeded,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
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
		ID:          d.ID,
		AssetID:     d.AssetID,
		ScannedByID: d.ScannedBy,
		ScanTime:    d.ScanTime.Format(TimeFormat),
		ScanMethod:  d.ScanMethod,
		ScanData:    d.ScanData,
		IsSucceeded: d.IsSucceeded,
		CreatedAt:   d.CreatedAt.Format(TimeFormat),
		UpdatedAt:   d.UpdatedAt.Format(TimeFormat),
	}
}

func ScanLogsToResponses(scanLogs []domain.ScanLog) []domain.ScanLogResponse {
	responses := make([]domain.ScanLogResponse, len(scanLogs))
	for i, scanLog := range scanLogs {
		responses[i] = ScanLogToResponse(&scanLog)
	}
	return responses
}
