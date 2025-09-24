package maintenance_record

import (
	"context"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/Rizz404/inventory-api/internal/utils"
)

// Repository defines data operations for maintenance records
type Repository interface {
	// Record mutations
	CreateRecord(ctx context.Context, payload *domain.MaintenanceRecord) (domain.MaintenanceRecord, error)
	UpdateRecord(ctx context.Context, recordId string, payload *domain.MaintenanceRecord) (domain.MaintenanceRecord, error)
	DeleteRecord(ctx context.Context, recordId string) error

	// Record queries
	GetRecordsPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceRecord, error)
	GetRecordsCursor(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceRecord, error)
	GetRecordById(ctx context.Context, recordId string) (domain.MaintenanceRecord, error)
	CountRecords(ctx context.Context, params query.Params) (int64, error)
	CheckRecordExist(ctx context.Context, recordId string) (bool, error)

	// Statistics
	GetMaintenanceRecordStatistics(ctx context.Context) (domain.MaintenanceRecordStatistics, error)
}

// AssetService for existence checks and populating asset info
type AssetService interface {
	CheckAssetExists(ctx context.Context, assetId string) (bool, error)
	GetAssetById(ctx context.Context, assetId string, langCode string) (domain.AssetResponse, error)
}

// UserService for existence checks
type UserService interface {
	CheckUserExists(ctx context.Context, userId string) (bool, error)
	GetUserById(ctx context.Context, userId string) (domain.UserResponse, error)
}

// MaintenanceRecordService business operations
type MaintenanceRecordService interface {
	CreateMaintenanceRecord(ctx context.Context, payload *domain.CreateMaintenanceRecordPayload, performedBy string) (domain.MaintenanceRecordResponse, error)
	UpdateMaintenanceRecord(ctx context.Context, recordId string, payload *domain.CreateMaintenanceRecordPayload) (domain.MaintenanceRecordResponse, error)
	DeleteMaintenanceRecord(ctx context.Context, recordId string) error
	GetMaintenanceRecordsPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceRecordListResponse, int64, error)
	GetMaintenanceRecordsCursor(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceRecordListResponse, error)
	GetMaintenanceRecordById(ctx context.Context, recordId string, langCode string) (domain.MaintenanceRecordResponse, error)
	CheckMaintenanceRecordExists(ctx context.Context, recordId string) (bool, error)
	CountMaintenanceRecords(ctx context.Context, params query.Params) (int64, error)
	GetMaintenanceRecordStatistics(ctx context.Context) (domain.MaintenanceRecordStatisticsResponse, error)
}

type Service struct {
	Repo         Repository
	AssetService AssetService
	UserService  UserService
}

var _ MaintenanceRecordService = (*Service)(nil)

func NewService(r Repository, assetSvc AssetService, userSvc UserService) MaintenanceRecordService {
	return &Service{Repo: r, AssetService: assetSvc, UserService: userSvc}
}

func (s *Service) CreateMaintenanceRecord(ctx context.Context, payload *domain.CreateMaintenanceRecordPayload, performedBy string) (domain.MaintenanceRecordResponse, error) {
	// Validate asset exists
	if exists, err := s.AssetService.CheckAssetExists(ctx, payload.AssetID); err != nil {
		return domain.MaintenanceRecordResponse{}, err
	} else if !exists {
		return domain.MaintenanceRecordResponse{}, domain.ErrNotFoundWithKey(utils.ErrAssetNotFoundKey)
	}

	// Determine performer: use payload if provided, otherwise fallback to performedBy param
	var performerPtr *string
	if payload.PerformedByUser != nil && *payload.PerformedByUser != "" {
		performerPtr = payload.PerformedByUser
	} else if performedBy != "" {
		performerPtr = &performedBy
	}

	// Validate performer if present
	if performerPtr != nil {
		if exists, err := s.UserService.CheckUserExists(ctx, *performerPtr); err != nil {
			return domain.MaintenanceRecordResponse{}, err
		} else if !exists {
			return domain.MaintenanceRecordResponse{}, domain.ErrNotFoundWithKey(utils.ErrUserNotFoundKey)
		}
	}

	// Parse date
	maintenanceDate, err := time.Parse("2006-01-02", payload.MaintenanceDate)
	if err != nil {
		return domain.MaintenanceRecordResponse{}, domain.ErrBadRequestWithKey(utils.ErrMaintenanceRecordDateRequiredKey)
	}

	// Build domain entity
	record := domain.MaintenanceRecord{
		ScheduleID:        payload.ScheduleID,
		AssetID:           payload.AssetID,
		MaintenanceDate:   maintenanceDate,
		PerformedByUser:   performerPtr,
		PerformedByVendor: payload.PerformedByVendor,
		ActualCost:        payload.ActualCost,
		Translations:      make([]domain.MaintenanceRecordTranslation, len(payload.Translations)),
	}
	for i, t := range payload.Translations {
		record.Translations[i] = domain.MaintenanceRecordTranslation{
			LangCode: t.LangCode,
			Title:    t.Title,
			Notes:    t.Notes,
		}
	}

	created, err := s.Repo.CreateRecord(ctx, &record)
	if err != nil {
		return domain.MaintenanceRecordResponse{}, err
	}
	return mapper.MaintenanceRecordToResponse(&created, mapper.DefaultLangCode), nil
}

func (s *Service) UpdateMaintenanceRecord(ctx context.Context, recordId string, payload *domain.CreateMaintenanceRecordPayload) (domain.MaintenanceRecordResponse, error) {
	// Ensure record exists
	if exists, err := s.Repo.CheckRecordExist(ctx, recordId); err != nil {
		return domain.MaintenanceRecordResponse{}, err
	} else if !exists {
		return domain.MaintenanceRecordResponse{}, domain.ErrNotFoundWithKey(utils.ErrMaintenanceRecordNotFoundKey)
	}

	// Validate asset if provided
	if payload.AssetID != "" {
		if exists, err := s.AssetService.CheckAssetExists(ctx, payload.AssetID); err != nil {
			return domain.MaintenanceRecordResponse{}, err
		} else if !exists {
			return domain.MaintenanceRecordResponse{}, domain.ErrNotFoundWithKey(utils.ErrAssetNotFoundKey)
		}
	}

	// Validate performer if provided
	if payload.PerformedByUser != nil && *payload.PerformedByUser != "" {
		if exists, err := s.UserService.CheckUserExists(ctx, *payload.PerformedByUser); err != nil {
			return domain.MaintenanceRecordResponse{}, err
		} else if !exists {
			return domain.MaintenanceRecordResponse{}, domain.ErrNotFoundWithKey(utils.ErrUserNotFoundKey)
		}
	}

	// Parse date if provided
	var maintenanceDate time.Time
	if payload.MaintenanceDate != "" {
		d, err := time.Parse("2006-01-02", payload.MaintenanceDate)
		if err != nil {
			return domain.MaintenanceRecordResponse{}, domain.ErrBadRequestWithKey(utils.ErrMaintenanceRecordDateRequiredKey)
		}
		maintenanceDate = d
	}

	record := domain.MaintenanceRecord{
		ScheduleID:        payload.ScheduleID,
		AssetID:           payload.AssetID,
		MaintenanceDate:   maintenanceDate,
		PerformedByUser:   payload.PerformedByUser,
		PerformedByVendor: payload.PerformedByVendor,
		ActualCost:        payload.ActualCost,
		Translations:      make([]domain.MaintenanceRecordTranslation, len(payload.Translations)),
	}
	for i, t := range payload.Translations {
		record.Translations[i] = domain.MaintenanceRecordTranslation{
			LangCode: t.LangCode,
			Title:    t.Title,
			Notes:    t.Notes,
		}
	}

	updated, err := s.Repo.UpdateRecord(ctx, recordId, &record)
	if err != nil {
		return domain.MaintenanceRecordResponse{}, err
	}
	return mapper.MaintenanceRecordToResponse(&updated, mapper.DefaultLangCode), nil
}

func (s *Service) DeleteMaintenanceRecord(ctx context.Context, recordId string) error {
	return s.Repo.DeleteRecord(ctx, recordId)
}

func (s *Service) GetMaintenanceRecordsPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceRecordListResponse, int64, error) {
	records, err := s.Repo.GetRecordsPaginated(ctx, params, langCode)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.Repo.CountRecords(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	recordResponses := mapper.MaintenanceRecordsToListResponses(records, langCode)

	return recordResponses, count, nil
}

func (s *Service) GetMaintenanceRecordsCursor(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceRecordListResponse, error) {
	records, err := s.Repo.GetRecordsCursor(ctx, params, langCode)
	if err != nil {
		return nil, err
	}

	recordResponses := mapper.MaintenanceRecordsToListResponses(records, langCode)

	return recordResponses, nil
}

func (s *Service) GetMaintenanceRecordById(ctx context.Context, recordId string, langCode string) (domain.MaintenanceRecordResponse, error) {
	record, err := s.Repo.GetRecordById(ctx, recordId)
	if err != nil {
		return domain.MaintenanceRecordResponse{}, err
	}
	return mapper.MaintenanceRecordToResponse(&record, langCode), nil
}

func (s *Service) CheckMaintenanceRecordExists(ctx context.Context, recordId string) (bool, error) {
	return s.Repo.CheckRecordExist(ctx, recordId)
}

func (s *Service) CountMaintenanceRecords(ctx context.Context, params query.Params) (int64, error) {
	return s.Repo.CountRecords(ctx, params)
}

func (s *Service) GetMaintenanceRecordStatistics(ctx context.Context) (domain.MaintenanceRecordStatisticsResponse, error) {
	stats, err := s.Repo.GetMaintenanceRecordStatistics(ctx)
	if err != nil {
		return domain.MaintenanceRecordStatisticsResponse{}, err
	}
	return mapper.MaintenanceRecordStatisticsToResponse(&stats), nil
}
