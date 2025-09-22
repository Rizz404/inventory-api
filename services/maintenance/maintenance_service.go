package maintenance

import (
	"context"
	"fmt"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/utils"
	"github.com/Rizz404/inventory-api/internal/web"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// * Repository interface defines the contract for maintenance data operations
type Repository interface {
	// * SCHEDULE MUTATIONS
	CreateMaintenanceSchedule(ctx context.Context, payload *domain.MaintenanceSchedule) (domain.MaintenanceSchedule, error)
	UpdateMaintenanceSchedule(ctx context.Context, scheduleId string, payload *domain.UpdateMaintenanceSchedulePayload) (domain.MaintenanceSchedule, error)
	DeleteMaintenanceSchedule(ctx context.Context, scheduleId string) error

	// * RECORD MUTATIONS
	CreateMaintenanceRecord(ctx context.Context, payload *domain.MaintenanceRecord) (domain.MaintenanceRecord, error)
	UpdateMaintenanceRecord(ctx context.Context, recordId string, payload *domain.UpdateMaintenanceRecordPayload) (domain.MaintenanceRecord, error)
	DeleteMaintenanceRecord(ctx context.Context, recordId string) error

	// * SCHEDULE QUERIES
	GetMaintenanceSchedulesPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceSchedule, error)
	GetMaintenanceSchedulesCursor(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceSchedule, error)
	GetMaintenanceScheduleById(ctx context.Context, scheduleId string) (domain.MaintenanceSchedule, error)
	GetMaintenanceSchedulesByAssetId(ctx context.Context, assetId string, params query.Params) ([]domain.MaintenanceSchedule, error)
	CheckMaintenanceScheduleExist(ctx context.Context, scheduleId string) (bool, error)
	CountMaintenanceSchedules(ctx context.Context, params query.Params) (int64, error)

	// * RECORD QUERIES
	GetMaintenanceRecordsPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceRecord, error)
	GetMaintenanceRecordsCursor(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceRecord, error)
	GetMaintenanceRecordById(ctx context.Context, recordId string) (domain.MaintenanceRecord, error)
	GetMaintenanceRecordsByAssetId(ctx context.Context, assetId string, params query.Params) ([]domain.MaintenanceRecord, error)
	GetMaintenanceRecordsByScheduleId(ctx context.Context, scheduleId string, params query.Params) ([]domain.MaintenanceRecord, error)
	CheckMaintenanceRecordExist(ctx context.Context, recordId string) (bool, error)
	CountMaintenanceRecords(ctx context.Context, params query.Params) (int64, error)

	// * STATISTICS
	GetMaintenanceStatistics(ctx context.Context) (domain.MaintenanceStatistics, error)
}

// * AssetService interface for checking asset existence
type AssetService interface {
	CheckAssetExist(ctx context.Context, assetId string) (bool, error)
}

// * UserService interface for checking user existence
type UserService interface {
	CheckUserExist(ctx context.Context, userId string) (bool, error)
}

type MaintenanceService struct {
	repository   Repository
	assetService AssetService
	userService  UserService
}

func NewMaintenanceService(
	repository Repository,
	assetService AssetService,
	userService UserService,
) *MaintenanceService {
	return &MaintenanceService{
		repository:   repository,
		assetService: assetService,
		userService:  userService,
	}
}

// *===========================MAINTENANCE SCHEDULE MUTATIONS===========================*
func (s *MaintenanceService) CreateMaintenanceSchedule(ctx context.Context, payload *domain.MaintenanceSchedule) (*web.SuccessResponse, error) {
	// Validate required fields
	if payload.AssetID == "" {
		return nil, domain.ErrAssetIDRequired
	}
	if payload.MaintenanceType == "" {
		return nil, domain.ErrMaintenanceTypeRequired
	}
	if payload.ScheduledDate.IsZero() {
		return nil, domain.ErrScheduledDateRequired
	}
	if len(payload.Translations) == 0 {
		return nil, domain.ErrTranslationsRequired
	}

	// Validate asset exists
	exists, err := s.assetService.CheckAssetExist(ctx, payload.AssetID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrAssetNotFound
	}

	// Validate created by user exists
	if payload.CreatedBy != "" {
		exists, err := s.userService.CheckUserExist(ctx, payload.CreatedBy)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, domain.ErrUserNotFound
		}
	}

	// Validate maintenance type
	if !s.isValidMaintenanceType(payload.MaintenanceType) {
		return nil, domain.ErrInvalidMaintenanceType
	}

	// Validate status
	if payload.Status != "" && !s.isValidScheduleStatus(payload.Status) {
		return nil, domain.ErrInvalidStatus
	}

	// Set default status
	if payload.Status == "" {
		payload.Status = domain.ScheduleStatusScheduled
	}

	// Set ID and timestamps
	payload.ID = ulid.Make().String()
	payload.CreatedAt = time.Now()
	payload.UpdatedAt = time.Now()

	// Validate and set translations
	if err := s.validateScheduleTranslations(payload.Translations); err != nil {
		return nil, err
	}

	// Validate scheduled date is not in the past for new schedules
	if payload.ScheduledDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, domain.ErrScheduledDateInPast
	}

	// Create maintenance schedule
	schedule, err := s.repository.CreateMaintenanceSchedule(ctx, payload)
	if err != nil {
		return nil, err
	}

	return &web.SuccessResponse{
		Code:    201,
		Message: utils.GetLocalizedMessage(ctx, "maintenance_schedule_created_successfully"),
		Data:    schedule,
	}, nil
}

func (s *MaintenanceService) UpdateMaintenanceSchedule(ctx context.Context, scheduleId string, payload *domain.UpdateMaintenanceSchedulePayload) (*web.SuccessResponse, error) {
	// Validate schedule exists
	exists, err := s.repository.CheckMaintenanceScheduleExist(ctx, scheduleId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrMaintenanceScheduleNotFound
	}

	// Validate asset if provided
	if payload.AssetID != nil && *payload.AssetID != "" {
		exists, err := s.assetService.CheckAssetExist(ctx, *payload.AssetID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, domain.ErrAssetNotFound
		}
	}

	// Validate maintenance type if provided
	if payload.MaintenanceType != nil && !s.isValidMaintenanceType(*payload.MaintenanceType) {
		return nil, domain.ErrInvalidMaintenanceType
	}

	// Validate status if provided
	if payload.Status != nil && !s.isValidScheduleStatus(*payload.Status) {
		return nil, domain.ErrInvalidStatus
	}

	// Validate scheduled date if provided
	if payload.ScheduledDate != nil && payload.ScheduledDate.Before(time.Now().Truncate(24*time.Hour)) {
		// Allow past dates for status updates (e.g., marking as completed)
		// Only restrict for scheduled status
		if payload.Status == nil || *payload.Status == domain.ScheduleStatusScheduled {
			return nil, domain.ErrScheduledDateInPast
		}
	}

	// Validate translations if provided
	if len(payload.Translations) > 0 {
		if err := s.validateScheduleTranslations(payload.Translations); err != nil {
			return nil, err
		}
	}

	// Set updated timestamp
	now := time.Now()
	payload.UpdatedAt = &now

	// Update maintenance schedule
	schedule, err := s.repository.UpdateMaintenanceSchedule(ctx, scheduleId, payload)
	if err != nil {
		return nil, err
	}

	return &web.SuccessResponse{
		Code:    200,
		Message: utils.GetLocalizedMessage(ctx, "maintenance_schedule_updated_successfully"),
		Data:    schedule,
	}, nil
}

func (s *MaintenanceService) DeleteMaintenanceSchedule(ctx context.Context, scheduleId string) (*web.SuccessResponse, error) {
	// Check if maintenance schedule exists
	exists, err := s.repository.CheckMaintenanceScheduleExist(ctx, scheduleId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrMaintenanceScheduleNotFound
	}

	// Delete maintenance schedule
	if err := s.repository.DeleteMaintenanceSchedule(ctx, scheduleId); err != nil {
		return nil, err
	}

	return &web.SuccessResponse{
		Code:    200,
		Message: utils.GetLocalizedMessage(ctx, "maintenance_schedule_deleted_successfully"),
		Data:    nil,
	}, nil
}

// *===========================MAINTENANCE RECORD MUTATIONS===========================*
func (s *MaintenanceService) CreateMaintenanceRecord(ctx context.Context, payload *domain.MaintenanceRecord) (*web.SuccessResponse, error) {
	// Validate required fields
	if payload.AssetID == "" {
		return nil, domain.ErrAssetIDRequired
	}
	if payload.MaintenanceDate.IsZero() {
		return nil, domain.ErrMaintenanceDateRequired
	}
	if len(payload.Translations) == 0 {
		return nil, domain.ErrTranslationsRequired
	}

	// Validate at least one performer is specified
	if payload.PerformedByUser == nil && (payload.PerformedByVendor == nil || *payload.PerformedByVendor == "") {
		return nil, domain.ErrPerformerRequired
	}

	// Validate asset exists
	exists, err := s.assetService.CheckAssetExist(ctx, payload.AssetID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrAssetNotFound
	}

	// Validate schedule if provided
	if payload.ScheduleID != nil && *payload.ScheduleID != "" {
		exists, err := s.repository.CheckMaintenanceScheduleExist(ctx, *payload.ScheduleID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, domain.ErrMaintenanceScheduleNotFound
		}
	}

	// Validate performed by user if provided
	if payload.PerformedByUser != nil && *payload.PerformedByUser != "" {
		exists, err := s.userRepo.CheckUserExist(ctx, *payload.PerformedByUser)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, domain.ErrUserNotFound
		}
	}

	// Set ID and timestamps
	payload.ID = ulid.Make().String()
	payload.CreatedAt = time.Now()
	payload.UpdatedAt = time.Now()

	// Validate and set translations
	if err := s.validateRecordTranslations(payload.Translations); err != nil {
		return nil, err
	}

	// Validate maintenance date is not in the future
	if payload.MaintenanceDate.After(time.Now()) {
		return nil, domain.ErrMaintenanceDateInFuture
	}

	// Validate cost if provided
	if payload.ActualCost != nil && *payload.ActualCost < 0 {
		return nil, domain.ErrInvalidCost
	}

	// Create maintenance record
	record, err := s.repository.CreateMaintenanceRecord(ctx, payload)
	if err != nil {
		return nil, err
	}

	return &web.SuccessResponse{
		Code:    201,
		Message: utils.GetLocalizedMessage(ctx, "maintenance_record_created_successfully"),
		Data:    record,
	}, nil
}

func (s *MaintenanceService) UpdateMaintenanceRecord(ctx context.Context, recordId string, payload *domain.UpdateMaintenanceRecordPayload) (*web.SuccessResponse, error) {
	// Validate record exists
	exists, err := s.repository.CheckMaintenanceRecordExist(ctx, recordId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrMaintenanceRecordNotFound
	}

	// Validate asset if provided
	if payload.AssetID != nil && *payload.AssetID != "" {
		exists, err := s.assetService.CheckAssetExist(ctx, *payload.AssetID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, domain.ErrAssetNotFound
		}
	}

	// Validate schedule if provided
	if payload.ScheduleID != nil && *payload.ScheduleID != "" {
		exists, err := s.repository.CheckMaintenanceScheduleExist(ctx, *payload.ScheduleID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, domain.ErrMaintenanceScheduleNotFound
		}
	}

	// Validate performed by user if provided
	if payload.PerformedByUser != nil && *payload.PerformedByUser != "" {
		exists, err := s.userRepo.CheckUserExist(ctx, *payload.PerformedByUser)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, domain.ErrUserNotFound
		}
	}

	// Validate maintenance date if provided
	if payload.MaintenanceDate != nil && payload.MaintenanceDate.After(time.Now()) {
		return nil, domain.ErrMaintenanceDateInFuture
	}

	// Validate cost if provided
	if payload.ActualCost != nil && *payload.ActualCost < 0 {
		return nil, domain.ErrInvalidCost
	}

	// Validate translations if provided
	if len(payload.Translations) > 0 {
		if err := s.validateRecordTranslations(payload.Translations); err != nil {
			return nil, err
		}
	}

	// Set updated timestamp
	now := time.Now()
	payload.UpdatedAt = &now

	// Update maintenance record
	record, err := s.repository.UpdateMaintenanceRecord(ctx, recordId, payload)
	if err != nil {
		return nil, err
	}

	return &web.SuccessResponse{
		Code:    200,
		Message: utils.GetLocalizedMessage(ctx, "maintenance_record_updated_successfully"),
		Data:    record,
	}, nil
}

func (s *MaintenanceService) DeleteMaintenanceRecord(ctx context.Context, recordId string) (*web.SuccessResponse, error) {
	// Check if maintenance record exists
	exists, err := s.repository.CheckMaintenanceRecordExist(ctx, recordId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrMaintenanceRecordNotFound
	}

	// Delete maintenance record
	if err := s.repository.DeleteMaintenanceRecord(ctx, recordId); err != nil {
		return nil, err
	}

	return &web.SuccessResponse{
		Code:    200,
		Message: utils.GetLocalizedMessage(ctx, "maintenance_record_deleted_successfully"),
		Data:    nil,
	}, nil
}

// *===========================MAINTENANCE SCHEDULE QUERIES===========================*
func (s *MaintenanceService) GetMaintenanceSchedulesPaginated(ctx context.Context, params query.Params) (*web.PaginatedResponse, error) {
	langCode := utils.GetLanguageCode(ctx)

	// Get schedules
	schedules, err := s.repository.GetMaintenanceSchedulesPaginated(ctx, params, langCode)
	if err != nil {
		return nil, err
	}

	// Get total count for pagination
	totalCount, err := s.repository.CountMaintenanceSchedules(ctx, params)
	if err != nil {
		return nil, err
	}

	// Calculate pagination info
	totalPages := (totalCount + int64(params.Pagination.Limit) - 1) / int64(params.Pagination.Limit)
	currentPage := (params.Pagination.Offset / params.Pagination.Limit) + 1

	return &web.PaginatedResponse{
		Code:        200,
		Message:     utils.GetLocalizedMessage(ctx, "maintenance_schedules_retrieved_successfully"),
		Data:        schedules,
		CurrentPage: currentPage,
		TotalPages:  int(totalPages),
		TotalItems:  int(totalCount),
		ItemsPerPage: params.Pagination.Limit,
	}, nil
}

func (s *MaintenanceService) GetMaintenanceSchedulesCursor(ctx context.Context, params query.Params) (*web.CursorResponse, error) {
	langCode := utils.GetLanguageCode(ctx)

	// Get schedules
	schedules, err := s.repository.GetMaintenanceSchedulesCursor(ctx, params, langCode)
	if err != nil {
		return nil, err
	}

	// Determine next cursor
	var nextCursor *string
	if len(schedules) == params.Pagination.Limit {
		if len(schedules) > 0 {
			nextCursor = &schedules[len(schedules)-1].ID
		}
	}

	return &web.CursorResponse{
		Code:       200,
		Message:    utils.GetLocalizedMessage(ctx, "maintenance_schedules_retrieved_successfully"),
		Data:       schedules,
		NextCursor: nextCursor,
		HasMore:    nextCursor != nil,
	}, nil
}

func (s *MaintenanceService) GetMaintenanceScheduleById(ctx context.Context, scheduleId string) (*web.SuccessResponse, error) {
	// Validate schedule exists
	exists, err := s.repository.CheckMaintenanceScheduleExist(ctx, scheduleId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrMaintenanceScheduleNotFound
	}

	// Get schedule
	schedule, err := s.repository.GetMaintenanceScheduleById(ctx, scheduleId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrMaintenanceScheduleNotFound
		}
		return nil, err
	}

	return &web.SuccessResponse{
		Code:    200,
		Message: utils.GetLocalizedMessage(ctx, "maintenance_schedule_retrieved_successfully"),
		Data:    schedule,
	}, nil
}

func (s *MaintenanceService) GetMaintenanceSchedulesByAssetId(ctx context.Context, assetId string, params query.Params) (*web.SuccessResponse, error) {
	// Validate asset exists
	exists, err := s.assetService.CheckAssetExist(ctx, assetId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrAssetNotFound
	}

	// Get schedules
	schedules, err := s.repository.GetMaintenanceSchedulesByAssetId(ctx, assetId, params)
	if err != nil {
		return nil, err
	}

	return &web.SuccessResponse{
		Code:    200,
		Message: utils.GetLocalizedMessage(ctx, "asset_maintenance_schedules_retrieved_successfully"),
		Data:    schedules,
	}, nil
}

// *===========================MAINTENANCE RECORD QUERIES===========================*
func (s *MaintenanceService) GetMaintenanceRecordsPaginated(ctx context.Context, params query.Params) (*web.PaginatedResponse, error) {
	langCode := utils.GetLanguageCode(ctx)

	// Get records
	records, err := s.repository.GetMaintenanceRecordsPaginated(ctx, params, langCode)
	if err != nil {
		return nil, err
	}

	// Get total count for pagination
	totalCount, err := s.repository.CountMaintenanceRecords(ctx, params)
	if err != nil {
		return nil, err
	}

	// Calculate pagination info
	totalPages := (totalCount + int64(params.Pagination.Limit) - 1) / int64(params.Pagination.Limit)
	currentPage := (params.Pagination.Offset / params.Pagination.Limit) + 1

	return &web.PaginatedResponse{
		Code:        200,
		Message:     utils.GetLocalizedMessage(ctx, "maintenance_records_retrieved_successfully"),
		Data:        records,
		CurrentPage: currentPage,
		TotalPages:  int(totalPages),
		TotalItems:  int(totalCount),
		ItemsPerPage: params.Pagination.Limit,
	}, nil
}

func (s *MaintenanceService) GetMaintenanceRecordsCursor(ctx context.Context, params query.Params) (*web.CursorResponse, error) {
	langCode := utils.GetLanguageCode(ctx)

	// Get records
	records, err := s.repository.GetMaintenanceRecordsCursor(ctx, params, langCode)
	if err != nil {
		return nil, err
	}

	// Determine next cursor
	var nextCursor *string
	if len(records) == params.Pagination.Limit {
		if len(records) > 0 {
			nextCursor = &records[len(records)-1].ID
		}
	}

	return &web.CursorResponse{
		Code:       200,
		Message:    utils.GetLocalizedMessage(ctx, "maintenance_records_retrieved_successfully"),
		Data:       records,
		NextCursor: nextCursor,
		HasMore:    nextCursor != nil,
	}, nil
}

func (s *MaintenanceService) GetMaintenanceRecordById(ctx context.Context, recordId string) (*web.SuccessResponse, error) {
	// Validate record exists
	exists, err := s.repository.CheckMaintenanceRecordExist(ctx, recordId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrMaintenanceRecordNotFound
	}

	// Get record
	record, err := s.repository.GetMaintenanceRecordById(ctx, recordId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrMaintenanceRecordNotFound
		}
		return nil, err
	}

	return &web.SuccessResponse{
		Code:    200,
		Message: utils.GetLocalizedMessage(ctx, "maintenance_record_retrieved_successfully"),
		Data:    record,
	}, nil
}

func (s *MaintenanceService) GetMaintenanceRecordsByAssetId(ctx context.Context, assetId string, params query.Params) (*web.SuccessResponse, error) {
	// Validate asset exists
	exists, err := s.assetService.CheckAssetExist(ctx, assetId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrAssetNotFound
	}

	// Get records
	records, err := s.repository.GetMaintenanceRecordsByAssetId(ctx, assetId, params)
	if err != nil {
		return nil, err
	}

	return &web.SuccessResponse{
		Code:    200,
		Message: utils.GetLocalizedMessage(ctx, "asset_maintenance_records_retrieved_successfully"),
		Data:    records,
	}, nil
}

func (s *MaintenanceService) GetMaintenanceRecordsByScheduleId(ctx context.Context, scheduleId string, params query.Params) (*web.SuccessResponse, error) {
	// Validate schedule exists
	exists, err := s.repository.CheckMaintenanceScheduleExist(ctx, scheduleId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrMaintenanceScheduleNotFound
	}

	// Get records
	records, err := s.repository.GetMaintenanceRecordsByScheduleId(ctx, scheduleId, params)
	if err != nil {
		return nil, err
	}

	return &web.SuccessResponse{
		Code:    200,
		Message: utils.GetLocalizedMessage(ctx, "schedule_maintenance_records_retrieved_successfully"),
		Data:    records,
	}, nil
}

// *===========================MAINTENANCE STATISTICS===========================*
func (s *MaintenanceService) GetMaintenanceStatistics(ctx context.Context) (*web.SuccessResponse, error) {
	stats, err := s.repository.GetMaintenanceStatistics(ctx)
	if err != nil {
		return nil, err
	}

	return &web.SuccessResponse{
		Code:    200,
		Message: utils.GetLocalizedMessage(ctx, "maintenance_statistics_retrieved_successfully"),
		Data:    stats,
	}, nil
}

// *===========================VALIDATION HELPERS===========================*
func (s *MaintenanceService) isValidMaintenanceType(maintenanceType domain.MaintenanceScheduleType) bool {
	switch maintenanceType {
	case domain.MaintenanceTypePreventive, domain.MaintenanceTypeCorrective, domain.MaintenanceTypePredictive, domain.MaintenanceTypeEmergency:
		return true
	default:
		return false
	}
}

func (s *MaintenanceService) isValidScheduleStatus(status domain.ScheduleStatus) bool {
	switch status {
	case domain.ScheduleStatusScheduled, domain.ScheduleStatusInProgress, domain.ScheduleStatusCompleted, domain.ScheduleStatusCancelled:
		return true
	default:
		return false
	}
}

func (s *MaintenanceService) validateScheduleTranslations(translations []domain.MaintenanceScheduleTranslation) error {
	if len(translations) == 0 {
		return domain.ErrTranslationsRequired
	}

	// Check for required English translation
	hasEnglish := false
	for _, translation := range translations {
		if translation.LangCode == "en-US" {
			hasEnglish = true
			if translation.Title == "" {
				return domain.ErrScheduleTitleRequired
			}
			break
		}
	}

	if !hasEnglish {
		return domain.ErrEnglishTranslationRequired
	}

	// Validate each translation
	for _, translation := range translations {
		if translation.LangCode == "" {
			return domain.ErrTranslationLangCodeRequired
		}
		if translation.Title == "" {
			return domain.ErrScheduleTitleRequired
		}
	}

	return nil
}

func (s *MaintenanceService) validateRecordTranslations(translations []domain.MaintenanceRecordTranslation) error {
	if len(translations) == 0 {
		return domain.ErrTranslationsRequired
	}

	// Check for required English translation
	hasEnglish := false
	for _, translation := range translations {
		if translation.LangCode == "en-US" {
			hasEnglish = true
			if translation.Title == "" {
				return domain.ErrRecordTitleRequired
			}
			break
		}
	}

	if !hasEnglish {
		return domain.ErrEnglishTranslationRequired
	}

	// Validate each translation
	for _, translation := range translations {
		if translation.LangCode == "" {
			return domain.ErrTranslationLangCodeRequired
		}
		if translation.Title == "" {
			return domain.ErrRecordTitleRequired
		}
	}

	return nil
}
