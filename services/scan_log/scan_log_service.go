package scan_log

import (
	"context"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
)

// * Repository interface defines the contract for scan log data operations
type Repository interface {
	// * MUTATION
	CreateScanLog(ctx context.Context, payload *domain.ScanLog) (domain.ScanLog, error)
	DeleteScanLog(ctx context.Context, scanLogId string) error

	// * QUERY
	GetScanLogsPaginated(ctx context.Context, params domain.ScanLogParams) ([]domain.ScanLog, error)
	GetScanLogsCursor(ctx context.Context, params domain.ScanLogParams) ([]domain.ScanLog, error)
	GetScanLogById(ctx context.Context, scanLogId string) (domain.ScanLog, error)
	GetScanLogsByAssetId(ctx context.Context, assetId string, params domain.ScanLogParams) ([]domain.ScanLog, error)
	GetScanLogsByUserId(ctx context.Context, userId string, params domain.ScanLogParams) ([]domain.ScanLog, error)
	CheckScanLogExist(ctx context.Context, scanLogId string) (bool, error)
	CountScanLogs(ctx context.Context, params domain.ScanLogParams) (int64, error)
	GetScanLogStatistics(ctx context.Context) (domain.ScanLogStatistics, error)
	GetScanLogsForExport(ctx context.Context, params domain.ScanLogParams) ([]domain.ScanLog, error)
}

// * ScanLogService interface defines the contract for scan log business operations
type ScanLogService interface {
	// * MUTATION
	CreateScanLog(ctx context.Context, payload *domain.CreateScanLogPayload, scannedBy string) (domain.ScanLogResponse, error)
	DeleteScanLog(ctx context.Context, scanLogId string) error

	// * QUERY
	GetScanLogsPaginated(ctx context.Context, params domain.ScanLogParams) ([]domain.ScanLogResponse, int64, error)
	GetScanLogsCursor(ctx context.Context, params domain.ScanLogParams) ([]domain.ScanLogResponse, error)
	GetScanLogById(ctx context.Context, scanLogId string) (domain.ScanLogResponse, error)
	GetScanLogsByAssetId(ctx context.Context, assetId string, params domain.ScanLogParams) ([]domain.ScanLogResponse, error)
	GetScanLogsByUserId(ctx context.Context, userId string, params domain.ScanLogParams) ([]domain.ScanLogResponse, error)
	CheckScanLogExists(ctx context.Context, scanLogId string) (bool, error)
	CountScanLogs(ctx context.Context, params domain.ScanLogParams) (int64, error)
	GetScanLogStatistics(ctx context.Context) (domain.ScanLogStatisticsResponse, error)
	ExportScanLogList(ctx context.Context, payload domain.ExportScanLogListPayload, params domain.ScanLogParams, langCode string) ([]byte, string, error)
}

type Service struct {
	Repo Repository
}

// * Ensure Service implements ScanLogService interface
var _ ScanLogService = (*Service)(nil)

func NewService(r Repository) ScanLogService {
	return &Service{
		Repo: r,
	}
}

// *===========================MUTATION===========================*
func (s *Service) CreateScanLog(ctx context.Context, payload *domain.CreateScanLogPayload, scannedBy string) (domain.ScanLogResponse, error) {
	// * Prepare domain scan log
	newScanLog := domain.ScanLog{
		AssetID:         payload.AssetID,
		ScannedValue:    payload.ScannedValue,
		ScanMethod:      payload.ScanMethod,
		ScannedBy:       scannedBy,
		ScanTimestamp:   time.Now(),
		ScanLocationLat: payload.ScanLocationLat,
		ScanLocationLng: payload.ScanLocationLng,
		ScanResult:      payload.ScanResult, // Default result
	}

	createdScanLog, err := s.Repo.CreateScanLog(ctx, &newScanLog)
	if err != nil {
		return domain.ScanLogResponse{}, err
	}

	// * Convert to ScanLogResponse using mapper
	return mapper.ScanLogToResponse(&createdScanLog), nil
}

func (s *Service) DeleteScanLog(ctx context.Context, scanLogId string) error {
	// * Check if scan log exists
	exists, err := s.Repo.CheckScanLogExist(ctx, scanLogId)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrNotFound("scan log")
	}

	err = s.Repo.DeleteScanLog(ctx, scanLogId)
	if err != nil {
		return err
	}
	return nil
}

// *===========================QUERY===========================*
func (s *Service) GetScanLogsPaginated(ctx context.Context, params domain.ScanLogParams) ([]domain.ScanLogResponse, int64, error) {
	scanLogs, err := s.Repo.GetScanLogsPaginated(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// * Count total for pagination
	count, err := s.Repo.CountScanLogs(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// * Convert to ScanLogResponse using mapper
	scanLogResponses := mapper.ScanLogsToResponses(scanLogs)

	return scanLogResponses, count, nil
}

func (s *Service) GetScanLogsCursor(ctx context.Context, params domain.ScanLogParams) ([]domain.ScanLogResponse, error) {
	scanLogs, err := s.Repo.GetScanLogsCursor(ctx, params)
	if err != nil {
		return nil, err
	}

	// * Convert to ScanLogResponse using mapper
	scanLogResponses := mapper.ScanLogsToResponses(scanLogs)

	return scanLogResponses, nil
}

func (s *Service) GetScanLogById(ctx context.Context, scanLogId string) (domain.ScanLogResponse, error) {
	scanLog, err := s.Repo.GetScanLogById(ctx, scanLogId)
	if err != nil {
		return domain.ScanLogResponse{}, err
	}

	// * Convert to ScanLogResponse using mapper
	return mapper.ScanLogToResponse(&scanLog), nil
}

func (s *Service) GetScanLogsByAssetId(ctx context.Context, assetId string, params domain.ScanLogParams) ([]domain.ScanLogResponse, error) {
	scanLogs, err := s.Repo.GetScanLogsByAssetId(ctx, assetId, params)
	if err != nil {
		return nil, err
	}

	// * Convert to ScanLogResponse using mapper
	scanLogResponses := mapper.ScanLogsToResponses(scanLogs)

	return scanLogResponses, nil
}

func (s *Service) GetScanLogsByUserId(ctx context.Context, userId string, params domain.ScanLogParams) ([]domain.ScanLogResponse, error) {
	scanLogs, err := s.Repo.GetScanLogsByUserId(ctx, userId, params)
	if err != nil {
		return nil, err
	}

	// * Convert to ScanLogResponse using mapper
	scanLogResponses := mapper.ScanLogsToResponses(scanLogs)

	return scanLogResponses, nil
}

func (s *Service) CheckScanLogExists(ctx context.Context, scanLogId string) (bool, error) {
	exists, err := s.Repo.CheckScanLogExist(ctx, scanLogId)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Service) CountScanLogs(ctx context.Context, params domain.ScanLogParams) (int64, error) {
	count, err := s.Repo.CountScanLogs(ctx, params)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Service) GetScanLogStatistics(ctx context.Context) (domain.ScanLogStatisticsResponse, error) {
	stats, err := s.Repo.GetScanLogStatistics(ctx)
	if err != nil {
		return domain.ScanLogStatisticsResponse{}, err
	}

	// Convert to ScanLogStatisticsResponse using mapper
	return mapper.ScanLogStatisticsToResponse(&stats), nil
}
