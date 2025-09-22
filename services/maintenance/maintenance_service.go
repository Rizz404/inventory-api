package maintenance

import (
	"context"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/Rizz404/inventory-api/internal/utils"
)

// Repository defines data operations for maintenance schedules and records
type Repository interface {
	// Schedule mutations
	CreateSchedule(ctx context.Context, payload *domain.MaintenanceSchedule) (domain.MaintenanceSchedule, error)
	UpdateSchedule(ctx context.Context, scheduleId string, payload *domain.MaintenanceSchedule) (domain.MaintenanceSchedule, error)
	DeleteSchedule(ctx context.Context, scheduleId string) error

	// Record mutations
	CreateRecord(ctx context.Context, payload *domain.MaintenanceRecord) (domain.MaintenanceRecord, error)
	UpdateRecord(ctx context.Context, recordId string, payload *domain.MaintenanceRecord) (domain.MaintenanceRecord, error)
	DeleteRecord(ctx context.Context, recordId string) error

	// Schedule queries
	GetSchedulesPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceScheduleListItem, error)
	GetSchedulesCursor(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceScheduleListItem, error)
	GetScheduleById(ctx context.Context, scheduleId string) (domain.MaintenanceSchedule, error)
	CountSchedules(ctx context.Context, params query.Params) (int64, error)
	CheckScheduleExist(ctx context.Context, scheduleId string) (bool, error)

	// Record queries
	GetRecordsPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceRecordListItem, error)
	GetRecordsCursor(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceRecordListItem, error)
	GetRecordById(ctx context.Context, recordId string) (domain.MaintenanceRecord, error)
	CountRecords(ctx context.Context, params query.Params) (int64, error)
	CheckRecordExist(ctx context.Context, recordId string) (bool, error)

	// Statistics
	GetMaintenanceStatistics(ctx context.Context) (domain.MaintenanceStatistics, error)
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

// MaintenanceService business operations
type MaintenanceService interface {
	// Schedules
	CreateMaintenanceSchedule(ctx context.Context, payload *domain.CreateMaintenanceSchedulePayload, createdBy string) (domain.MaintenanceScheduleResponse, error)
	UpdateMaintenanceSchedule(ctx context.Context, scheduleId string, payload *domain.CreateMaintenanceSchedulePayload) (domain.MaintenanceScheduleResponse, error)
	DeleteMaintenanceSchedule(ctx context.Context, scheduleId string) error
	GetMaintenanceSchedulesPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceScheduleListItem, int64, error)
	GetMaintenanceSchedulesCursor(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceScheduleListItem, error)
	GetMaintenanceScheduleById(ctx context.Context, scheduleId string, langCode string) (domain.MaintenanceScheduleResponse, error)
	CheckMaintenanceScheduleExists(ctx context.Context, scheduleId string) (bool, error)
	CountMaintenanceSchedules(ctx context.Context, params query.Params) (int64, error)

	// Records
	CreateMaintenanceRecord(ctx context.Context, payload *domain.CreateMaintenanceRecordPayload, performedBy string) (domain.MaintenanceRecordResponse, error)
	UpdateMaintenanceRecord(ctx context.Context, recordId string, payload *domain.CreateMaintenanceRecordPayload) (domain.MaintenanceRecordResponse, error)
	DeleteMaintenanceRecord(ctx context.Context, recordId string) error
	GetMaintenanceRecordsPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceRecordListItem, int64, error)
	GetMaintenanceRecordsCursor(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceRecordListItem, error)
	GetMaintenanceRecordById(ctx context.Context, recordId string, langCode string) (domain.MaintenanceRecordResponse, error)
	CheckMaintenanceRecordExists(ctx context.Context, recordId string) (bool, error)
	CountMaintenanceRecords(ctx context.Context, params query.Params) (int64, error)

	// Statistics
	GetMaintenanceStatistics(ctx context.Context) (domain.MaintenanceStatisticsResponse, error)
}

type Service struct {
	Repo         Repository
	AssetService AssetService
	UserService  UserService
}

var _ MaintenanceService = (*Service)(nil)

func NewService(r Repository, assetSvc AssetService, userSvc UserService) MaintenanceService {
	return &Service{Repo: r, AssetService: assetSvc, UserService: userSvc}
}

// =========================== SCHEDULES ===========================

func (s *Service) CreateMaintenanceSchedule(ctx context.Context, payload *domain.CreateMaintenanceSchedulePayload, createdBy string) (domain.MaintenanceScheduleResponse, error) {
	// Validate creator user exists
	if createdBy != "" {
		if exists, err := s.UserService.CheckUserExists(ctx, createdBy); err != nil {
			return domain.MaintenanceScheduleResponse{}, err
		} else if !exists {
			return domain.MaintenanceScheduleResponse{}, domain.ErrNotFoundWithKey(utils.ErrUserNotFoundKey)
		}
	}
	// Validate asset exists
	if exists, err := s.AssetService.CheckAssetExists(ctx, payload.AssetID); err != nil {
		return domain.MaintenanceScheduleResponse{}, err
	} else if !exists {
		return domain.MaintenanceScheduleResponse{}, domain.ErrNotFoundWithKey(utils.ErrAssetNotFoundKey)
	}

	// Parse scheduled date
	scheduledDate, err := time.Parse("2006-01-02", payload.ScheduledDate)
	if err != nil {
		return domain.MaintenanceScheduleResponse{}, domain.ErrBadRequestWithKey(utils.ErrMaintenanceScheduleDateRequiredKey)
	}

	// Build domain entity
	schedule := domain.MaintenanceSchedule{
		AssetID:         payload.AssetID,
		MaintenanceType: payload.MaintenanceType,
		ScheduledDate:   scheduledDate,
		FrequencyMonths: payload.FrequencyMonths,
		Status:          domain.StatusScheduled,
		CreatedBy:       createdBy,
		Translations:    make([]domain.MaintenanceScheduleTranslation, len(payload.Translations)),
	}
	for i, t := range payload.Translations {
		schedule.Translations[i] = domain.MaintenanceScheduleTranslation{
			LangCode:    t.LangCode,
			Title:       t.Title,
			Description: t.Description,
		}
	}

	created, err := s.Repo.CreateSchedule(ctx, &schedule)
	if err != nil {
		return domain.MaintenanceScheduleResponse{}, err
	}
	return mapper.MaintenanceScheduleToResponse(&created, mapper.DefaultLangCode), nil
}

func (s *Service) UpdateMaintenanceSchedule(ctx context.Context, scheduleId string, payload *domain.CreateMaintenanceSchedulePayload) (domain.MaintenanceScheduleResponse, error) {
	// Ensure schedule exists
	if exists, err := s.Repo.CheckScheduleExist(ctx, scheduleId); err != nil {
		return domain.MaintenanceScheduleResponse{}, err
	} else if !exists {
		return domain.MaintenanceScheduleResponse{}, domain.ErrNotFoundWithKey(utils.ErrMaintenanceScheduleNotFoundKey)
	}

	// Validate asset if provided
	if payload.AssetID != "" {
		if exists, err := s.AssetService.CheckAssetExists(ctx, payload.AssetID); err != nil {
			return domain.MaintenanceScheduleResponse{}, err
		} else if !exists {
			return domain.MaintenanceScheduleResponse{}, domain.ErrNotFoundWithKey(utils.ErrAssetNotFoundKey)
		}
	}

	// Parse scheduled date
	var scheduledDate time.Time
	if payload.ScheduledDate != "" {
		d, err := time.Parse("2006-01-02", payload.ScheduledDate)
		if err != nil {
			return domain.MaintenanceScheduleResponse{}, domain.ErrBadRequestWithKey(utils.ErrMaintenanceScheduleDateRequiredKey)
		}
		scheduledDate = d
	}

	// Build partial domain for update
	schedule := domain.MaintenanceSchedule{
		AssetID:         payload.AssetID,
		MaintenanceType: payload.MaintenanceType,
		ScheduledDate:   scheduledDate,
		FrequencyMonths: payload.FrequencyMonths,
		Translations:    make([]domain.MaintenanceScheduleTranslation, len(payload.Translations)),
	}
	for i, t := range payload.Translations {
		schedule.Translations[i] = domain.MaintenanceScheduleTranslation{
			LangCode:    t.LangCode,
			Title:       t.Title,
			Description: t.Description,
		}
	}

	updated, err := s.Repo.UpdateSchedule(ctx, scheduleId, &schedule)
	if err != nil {
		return domain.MaintenanceScheduleResponse{}, err
	}
	return mapper.MaintenanceScheduleToResponse(&updated, mapper.DefaultLangCode), nil
}

func (s *Service) DeleteMaintenanceSchedule(ctx context.Context, scheduleId string) error {
	return s.Repo.DeleteSchedule(ctx, scheduleId)
}

func (s *Service) GetMaintenanceSchedulesPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceScheduleListItem, int64, error) {
	schedules, err := s.Repo.GetSchedulesPaginated(ctx, params, langCode)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.Repo.CountSchedules(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	return schedules, count, nil
}

func (s *Service) GetMaintenanceSchedulesCursor(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceScheduleListItem, error) {
	schedules, err := s.Repo.GetSchedulesCursor(ctx, params, langCode)
	if err != nil {
		return nil, err
	}
	return schedules, nil
}

func (s *Service) GetMaintenanceScheduleById(ctx context.Context, scheduleId string, langCode string) (domain.MaintenanceScheduleResponse, error) {
	schedule, err := s.Repo.GetScheduleById(ctx, scheduleId)
	if err != nil {
		return domain.MaintenanceScheduleResponse{}, err
	}
	return mapper.MaintenanceScheduleToResponse(&schedule, langCode), nil
}

func (s *Service) CheckMaintenanceScheduleExists(ctx context.Context, scheduleId string) (bool, error) {
	return s.Repo.CheckScheduleExist(ctx, scheduleId)
}

func (s *Service) CountMaintenanceSchedules(ctx context.Context, params query.Params) (int64, error) {
	return s.Repo.CountSchedules(ctx, params)
}

// =========================== RECORDS ===========================

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

func (s *Service) GetMaintenanceRecordsPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceRecordListItem, int64, error) {
	records, err := s.Repo.GetRecordsPaginated(ctx, params, langCode)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.Repo.CountRecords(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	return records, count, nil
}

func (s *Service) GetMaintenanceRecordsCursor(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceRecordListItem, error) {
	records, err := s.Repo.GetRecordsCursor(ctx, params, langCode)
	if err != nil {
		return nil, err
	}
	return records, nil
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

// =========================== STATISTICS ===========================

func (s *Service) GetMaintenanceStatistics(ctx context.Context) (domain.MaintenanceStatisticsResponse, error) {
	stats, err := s.Repo.GetMaintenanceStatistics(ctx)
	if err != nil {
		return domain.MaintenanceStatisticsResponse{}, err
	}
	return mapper.MaintenanceStatisticsToResponse(&stats), nil
}
