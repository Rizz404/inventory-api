package maintenance_schedule

import (
	"context"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/Rizz404/inventory-api/internal/utils"
)

// Repository defines data operations for maintenance schedules
type Repository interface {
	// Schedule mutations
	CreateSchedule(ctx context.Context, payload *domain.MaintenanceSchedule) (domain.MaintenanceSchedule, error)
	UpdateSchedule(ctx context.Context, scheduleId string, payload *domain.UpdateMaintenanceSchedulePayload) (domain.MaintenanceSchedule, error)
	DeleteSchedule(ctx context.Context, scheduleId string) error
	BulkDeleteSchedules(ctx context.Context, scheduleIds []string) (domain.BulkDeleteMaintenanceSchedules, error)

	// Schedule queries
	GetSchedulesPaginated(ctx context.Context, params domain.MaintenanceScheduleParams, langCode string) ([]domain.MaintenanceSchedule, error)
	GetSchedulesCursor(ctx context.Context, params domain.MaintenanceScheduleParams, langCode string) ([]domain.MaintenanceSchedule, error)
	GetScheduleById(ctx context.Context, scheduleId string) (domain.MaintenanceSchedule, error)
	CountSchedules(ctx context.Context, params domain.MaintenanceScheduleParams) (int64, error)
	CheckScheduleExist(ctx context.Context, scheduleId string) (bool, error)

	// Statistics
	GetMaintenanceScheduleStatistics(ctx context.Context) (domain.MaintenanceScheduleStatistics, error)

	// Export
	GetMaintenanceSchedulesForExport(ctx context.Context, params domain.MaintenanceScheduleParams, langCode string) ([]domain.MaintenanceSchedule, error)

	// Cron-related queries
	GetSchedulesDueSoon(ctx context.Context, daysFromNow int) ([]domain.MaintenanceSchedule, error)
	GetOverdueSchedules(ctx context.Context) ([]domain.MaintenanceSchedule, error)
	GetRecurringSchedulesToUpdate(ctx context.Context) ([]domain.MaintenanceSchedule, error)
	UpdateLastExecutedDate(ctx context.Context, scheduleId string, lastExecutedDate *time.Time) error
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

// NotificationService interface for creating notifications
type NotificationService interface {
	CreateNotification(ctx context.Context, payload *domain.CreateNotificationPayload) (domain.NotificationResponse, error)
}

// MaintenanceScheduleService business operations
type MaintenanceScheduleService interface {
	CreateMaintenanceSchedule(ctx context.Context, payload *domain.CreateMaintenanceSchedulePayload, createdBy string) (domain.MaintenanceScheduleResponse, error)
	UpdateMaintenanceSchedule(ctx context.Context, scheduleId string, payload *domain.UpdateMaintenanceSchedulePayload) (domain.MaintenanceScheduleResponse, error)
	DeleteMaintenanceSchedule(ctx context.Context, scheduleId string) error
	BulkDeleteMaintenanceSchedules(ctx context.Context, payload *domain.BulkDeleteMaintenanceSchedulesPayload) (domain.BulkDeleteMaintenanceSchedulesResponse, error)
	GetMaintenanceSchedulesPaginated(ctx context.Context, params domain.MaintenanceScheduleParams, langCode string) ([]domain.MaintenanceScheduleListResponse, int64, error)
	GetMaintenanceSchedulesCursor(ctx context.Context, params domain.MaintenanceScheduleParams, langCode string) ([]domain.MaintenanceScheduleListResponse, error)
	GetMaintenanceScheduleById(ctx context.Context, scheduleId string, langCode string) (domain.MaintenanceScheduleResponse, error)
	CheckMaintenanceScheduleExists(ctx context.Context, scheduleId string) (bool, error)
	CountMaintenanceSchedules(ctx context.Context, params domain.MaintenanceScheduleParams) (int64, error)
	GetMaintenanceScheduleStatistics(ctx context.Context) (domain.MaintenanceScheduleStatisticsResponse, error)
	ExportMaintenanceScheduleList(ctx context.Context, payload domain.ExportMaintenanceScheduleListPayload, params domain.MaintenanceScheduleParams, langCode string) ([]byte, string, error)
}

type Service struct {
	Repo                Repository
	AssetService        AssetService
	UserService         UserService
	NotificationService NotificationService
}

var _ MaintenanceScheduleService = (*Service)(nil)

func NewService(r Repository, assetSvc AssetService, userSvc UserService, notificationSvc NotificationService) MaintenanceScheduleService {
	return &Service{Repo: r, AssetService: assetSvc, UserService: userSvc, NotificationService: notificationSvc}
}

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

	// Parse next scheduled date
	nextScheduledDate, err := time.Parse("2006-01-02", payload.NextScheduledDate)
	if err != nil {
		return domain.MaintenanceScheduleResponse{}, domain.ErrBadRequestWithKey(utils.ErrMaintenanceScheduleDateRequiredKey)
	}

	// Set defaults
	isRecurring := false
	if payload.IsRecurring != nil {
		isRecurring = *payload.IsRecurring
	}

	autoComplete := false
	if payload.AutoComplete != nil {
		autoComplete = *payload.AutoComplete
	}

	// Validate: if recurring, must have interval
	if isRecurring && (payload.IntervalValue == nil || payload.IntervalUnit == nil) {
		return domain.MaintenanceScheduleResponse{}, domain.ErrBadRequest("recurring schedule must have interval_value and interval_unit")
	}

	// Build domain entity
	schedule := domain.MaintenanceSchedule{
		AssetID:           payload.AssetID,
		MaintenanceType:   payload.MaintenanceType,
		IsRecurring:       isRecurring,
		IntervalValue:     payload.IntervalValue,
		IntervalUnit:      payload.IntervalUnit,
		ScheduledTime:     payload.ScheduledTime,
		NextScheduledDate: nextScheduledDate,
		State:             domain.StateActive,
		AutoComplete:      autoComplete,
		EstimatedCost:     payload.EstimatedCost,
		CreatedBy:         createdBy,
		Translations:      make([]domain.MaintenanceScheduleTranslation, len(payload.Translations)),
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

func (s *Service) UpdateMaintenanceSchedule(ctx context.Context, scheduleId string, payload *domain.UpdateMaintenanceSchedulePayload) (domain.MaintenanceScheduleResponse, error) {
	// Ensure schedule exists
	if exists, err := s.Repo.CheckScheduleExist(ctx, scheduleId); err != nil {
		return domain.MaintenanceScheduleResponse{}, err
	} else if !exists {
		return domain.MaintenanceScheduleResponse{}, domain.ErrNotFoundWithKey(utils.ErrMaintenanceScheduleNotFoundKey)
	}

	updated, err := s.Repo.UpdateSchedule(ctx, scheduleId, payload)
	if err != nil {
		return domain.MaintenanceScheduleResponse{}, err
	}
	return mapper.MaintenanceScheduleToResponse(&updated, mapper.DefaultLangCode), nil
}

func (s *Service) DeleteMaintenanceSchedule(ctx context.Context, scheduleId string) error {
	return s.Repo.DeleteSchedule(ctx, scheduleId)
}

func (s *Service) BulkDeleteMaintenanceSchedules(ctx context.Context, payload *domain.BulkDeleteMaintenanceSchedulesPayload) (domain.BulkDeleteMaintenanceSchedulesResponse, error) {
	// * Validate that IDs are provided
	if len(payload.IDS) == 0 {
		return domain.BulkDeleteMaintenanceSchedulesResponse{}, domain.ErrBadRequest("maintenance schedule IDs are required")
	}

	// * Perform bulk delete operation
	result, err := s.Repo.BulkDeleteSchedules(ctx, payload.IDS)
	if err != nil {
		return domain.BulkDeleteMaintenanceSchedulesResponse{}, err
	}

	// * Convert to response
	response := domain.BulkDeleteMaintenanceSchedulesResponse{
		RequestedIDS: result.RequestedIDS,
		DeletedIDS:   result.DeletedIDS,
	}

	return response, nil
}

func (s *Service) GetMaintenanceSchedulesPaginated(ctx context.Context, params domain.MaintenanceScheduleParams, langCode string) ([]domain.MaintenanceScheduleListResponse, int64, error) {
	schedules, err := s.Repo.GetSchedulesPaginated(ctx, params, langCode)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.Repo.CountSchedules(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	schedulesResponses := mapper.MaintenanceSchedulesToListResponses(schedules, langCode)

	return schedulesResponses, count, nil
}

func (s *Service) GetMaintenanceSchedulesCursor(ctx context.Context, params domain.MaintenanceScheduleParams, langCode string) ([]domain.MaintenanceScheduleListResponse, error) {
	schedules, err := s.Repo.GetSchedulesCursor(ctx, params, langCode)
	if err != nil {
		return nil, err
	}

	schedulesResponses := mapper.MaintenanceSchedulesToListResponses(schedules, langCode)

	return schedulesResponses, nil
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

func (s *Service) CountMaintenanceSchedules(ctx context.Context, params domain.MaintenanceScheduleParams) (int64, error) {
	return s.Repo.CountSchedules(ctx, params)
}

func (s *Service) GetMaintenanceScheduleStatistics(ctx context.Context) (domain.MaintenanceScheduleStatisticsResponse, error) {
	stats, err := s.Repo.GetMaintenanceScheduleStatistics(ctx)
	if err != nil {
		return domain.MaintenanceScheduleStatisticsResponse{}, err
	}
	return mapper.MaintenanceScheduleStatisticsToResponse(&stats), nil
}
