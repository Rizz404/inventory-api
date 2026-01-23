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
	BulkCreateScanLogs(ctx context.Context, scanLogs []domain.ScanLog) ([]domain.ScanLog, error)
	BulkDeleteScanLogs(ctx context.Context, scanLogIds []string) (domain.BulkDeleteScanLogs, error)

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
	BulkCreateScanLogs(ctx context.Context, payload *domain.BulkCreateScanLogsPayload, scannedBy string) (domain.BulkCreateScanLogsResponse, error)
	BulkDeleteScanLogs(ctx context.Context, payload *domain.BulkDeleteScanLogsPayload) (domain.BulkDeleteScanLogsResponse, error)

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
	ScanTimestamp:   time.Now().UTC(),

	createdScanLog, err := s.Repo.CreateScanLog(ctx, &newScanLog)
	if err != nil {
		return domain.ScanLogResponse{}, err
	}

	// * Convert to ScanLogResponse using mapper
	return mapper.ScanLogToResponse(&createdScanLog), nil
}

func (s *Service) BulkCreateScanLogs(ctx context.Context, payload *domain.BulkCreateScanLogsPayload, scannedBy string) (domain.BulkCreateScanLogsResponse, error) {
	if payload == nil || len(payload.ScanLogs) == 0 {
		return domain.BulkCreateScanLogsResponse{}, domain.ErrBadRequest("scan logs payload is required")
	}

	// * Build domain scan logs
	scanLogs := make([]domain.ScanLog, len(payload.ScanLogs))
	for i, item := range payload.ScanLogs {
		scanLogs[i] = domain.ScanLog{
			AssetID:         item.AssetID,
			ScannedValue:    item.ScannedValue,
			ScanMethod:      item.ScanMethod,
			ScannedBy:       scannedBy,
			ScanTimestamp:   time.Now().UTC(),
			ScanLocationLat: item.ScanLocationLat,
			ScanLocationLng: item.ScanLocationLng,
			ScanResult:      item.ScanResult,
		}
	}

	// * Call repository bulk create
	created, err := s.Repo.BulkCreateScanLogs(ctx, scanLogs)
	if err != nil {
		return domain.BulkCreateScanLogsResponse{}, err
	}

	// * Convert to responses
	response := domain.BulkCreateScanLogsResponse{
		ScanLogs: mapper.ScanLogsToResponses(created),
	}
	return response, nil
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

func (s *Service) BulkDeleteScanLogs(ctx context.Context, payload *domain.BulkDeleteScanLogsPayload) (domain.BulkDeleteScanLogsResponse, error) {
	// * Validate that IDs are provided
	if len(payload.IDS) == 0 {
		return domain.BulkDeleteScanLogsResponse{}, domain.ErrBadRequest("scan log IDs are required")
	}

	// * Perform bulk delete operation
	result, err := s.Repo.BulkDeleteScanLogs(ctx, payload.IDS)
	if err != nil {
		return domain.BulkDeleteScanLogsResponse{}, err
	}

	// * Convert to response
	response := domain.BulkDeleteScanLogsResponse{
		RequestedIDS: result.RequestedIDS,
		DeletedIDS:   result.DeletedIDS,
	}

	return response, nil
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
