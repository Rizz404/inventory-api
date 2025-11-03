package maintenance_record

import (
	"context"
	"strings"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/notification/messages"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/Rizz404/inventory-api/internal/utils"
)

// Repository defines data operations for maintenance records
type Repository interface {
	// Record mutations
	CreateRecord(ctx context.Context, payload *domain.MaintenanceRecord) (domain.MaintenanceRecord, error)
	UpdateRecord(ctx context.Context, recordId string, payload *domain.UpdateMaintenanceRecordPayload) (domain.MaintenanceRecord, error)
	DeleteRecord(ctx context.Context, recordId string) error

	// Record queries
	GetRecordsPaginated(ctx context.Context, params domain.MaintenanceRecordParams, langCode string) ([]domain.MaintenanceRecord, error)
	GetRecordsCursor(ctx context.Context, params domain.MaintenanceRecordParams, langCode string) ([]domain.MaintenanceRecord, error)
	GetRecordById(ctx context.Context, recordId string) (domain.MaintenanceRecord, error)
	CountRecords(ctx context.Context, params domain.MaintenanceRecordParams) (int64, error)
	CheckRecordExist(ctx context.Context, recordId string) (bool, error)
	GetMaintenanceRecordsForExport(ctx context.Context, params domain.MaintenanceRecordParams, langCode string) ([]domain.MaintenanceRecord, error)

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

// NotificationService interface for creating notifications
type NotificationService interface {
	CreateNotification(ctx context.Context, payload *domain.CreateNotificationPayload) (domain.NotificationResponse, error)
}

// MaintenanceRecordService business operations
type MaintenanceRecordService interface {
	CreateMaintenanceRecord(ctx context.Context, payload *domain.CreateMaintenanceRecordPayload, performedBy string) (domain.MaintenanceRecordResponse, error)
	UpdateMaintenanceRecord(ctx context.Context, recordId string, payload *domain.UpdateMaintenanceRecordPayload) (domain.MaintenanceRecordResponse, error)
	DeleteMaintenanceRecord(ctx context.Context, recordId string) error
	GetMaintenanceRecordsPaginated(ctx context.Context, params domain.MaintenanceRecordParams, langCode string) ([]domain.MaintenanceRecordListResponse, int64, error)
	GetMaintenanceRecordsCursor(ctx context.Context, params domain.MaintenanceRecordParams, langCode string) ([]domain.MaintenanceRecordListResponse, error)
	GetMaintenanceRecordById(ctx context.Context, recordId string, langCode string) (domain.MaintenanceRecordResponse, error)
	CheckMaintenanceRecordExists(ctx context.Context, recordId string) (bool, error)
	CountMaintenanceRecords(ctx context.Context, params domain.MaintenanceRecordParams) (int64, error)
	GetMaintenanceRecordStatistics(ctx context.Context) (domain.MaintenanceRecordStatisticsResponse, error)
	ExportMaintenanceRecordList(ctx context.Context, payload domain.ExportMaintenanceRecordListPayload, params domain.MaintenanceRecordParams, langCode string) ([]byte, string, error)
}

type Service struct {
	Repo                Repository
	AssetService        AssetService
	UserService         UserService
	NotificationService NotificationService
}

var _ MaintenanceRecordService = (*Service)(nil)

func NewService(r Repository, assetSvc AssetService, userSvc UserService, notificationSvc NotificationService) MaintenanceRecordService {
	return &Service{Repo: r, AssetService: assetSvc, UserService: userSvc, NotificationService: notificationSvc}
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

	// Parse completion date if provided
	var completionDate *time.Time
	if payload.CompletionDate != nil && *payload.CompletionDate != "" {
		parsed, err := time.Parse("2006-01-02", *payload.CompletionDate)
		if err != nil {
			return domain.MaintenanceRecordResponse{}, domain.ErrBadRequest("invalid completion date format")
		}
		completionDate = &parsed
	}

	// Build domain entity
	record := domain.MaintenanceRecord{
		ScheduleID:        payload.ScheduleID,
		AssetID:           payload.AssetID,
		MaintenanceDate:   maintenanceDate,
		CompletionDate:    completionDate,
		DurationMinutes:   payload.DurationMinutes,
		PerformedByUser:   performerPtr,
		PerformedByVendor: payload.PerformedByVendor,
		Result:            payload.Result,
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

	// Send notification for completed maintenance
	s.sendMaintenanceCompletedNotification(ctx, &created)

	return mapper.MaintenanceRecordToResponse(&created, mapper.DefaultLangCode), nil
}

func (s *Service) UpdateMaintenanceRecord(ctx context.Context, recordId string, payload *domain.UpdateMaintenanceRecordPayload) (domain.MaintenanceRecordResponse, error) {
	// Ensure record exists
	if exists, err := s.Repo.CheckRecordExist(ctx, recordId); err != nil {
		return domain.MaintenanceRecordResponse{}, err
	} else if !exists {
		return domain.MaintenanceRecordResponse{}, domain.ErrNotFoundWithKey(utils.ErrMaintenanceRecordNotFoundKey)
	}

	// Validate performer if being updated
	if payload.PerformedByUser != nil && *payload.PerformedByUser != "" {
		if exists, err := s.UserService.CheckUserExists(ctx, *payload.PerformedByUser); err != nil {
			return domain.MaintenanceRecordResponse{}, err
		} else if !exists {
			return domain.MaintenanceRecordResponse{}, domain.ErrNotFoundWithKey(utils.ErrUserNotFoundKey)
		}
	}

	updated, err := s.Repo.UpdateRecord(ctx, recordId, payload)
	if err != nil {
		return domain.MaintenanceRecordResponse{}, err
	}

	// Check if this update indicates a failed maintenance (e.g., notes contain "failed")
	failureReason := ""
	for _, t := range payload.Translations {
		if t.Notes != nil && (strings.Contains(strings.ToLower(*t.Notes), "failed") || strings.Contains(strings.ToLower(*t.Notes), "error")) {
			failureReason = *t.Notes
			break
		}
	}
	if failureReason != "" {
		s.sendMaintenanceFailedNotification(ctx, &updated, failureReason)
	}

	return mapper.MaintenanceRecordToResponse(&updated, mapper.DefaultLangCode), nil
}

func (s *Service) DeleteMaintenanceRecord(ctx context.Context, recordId string) error {
	return s.Repo.DeleteRecord(ctx, recordId)
}

func (s *Service) GetMaintenanceRecordsPaginated(ctx context.Context, params domain.MaintenanceRecordParams, langCode string) ([]domain.MaintenanceRecordListResponse, int64, error) {
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

func (s *Service) GetMaintenanceRecordsCursor(ctx context.Context, params domain.MaintenanceRecordParams, langCode string) ([]domain.MaintenanceRecordListResponse, error) {
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

func (s *Service) CountMaintenanceRecords(ctx context.Context, params domain.MaintenanceRecordParams) (int64, error) {
	return s.Repo.CountRecords(ctx, params)
}

func (s *Service) GetMaintenanceRecordStatistics(ctx context.Context) (domain.MaintenanceRecordStatisticsResponse, error) {
	stats, err := s.Repo.GetMaintenanceRecordStatistics(ctx)
	if err != nil {
		return domain.MaintenanceRecordStatisticsResponse{}, err
	}
	return mapper.MaintenanceRecordStatisticsToResponse(&stats), nil
}

// sendMaintenanceCompletedNotification sends notification for completed maintenance
func (s *Service) sendMaintenanceCompletedNotification(ctx context.Context, record *domain.MaintenanceRecord) {
	if s.NotificationService == nil {
		return
	}

	// Get asset details
	asset, err := s.AssetService.GetAssetById(ctx, record.AssetID, "en-US") // Default lang
	if err != nil {
		return
	}

	if asset.AssignedToID == nil || *asset.AssignedToID == "" {
		return
	}

	// Get notes from translations (default to first or empty)
	notes := ""
	for _, t := range record.Translations {
		if t.LangCode == "en-US" && t.Notes != nil {
			notes = *t.Notes
			break
		}
	}
	if notes == "" && len(record.Translations) > 0 && record.Translations[0].Notes != nil {
		notes = *record.Translations[0].Notes
	}

	titleKey, messageKey, params := messages.MaintenanceCompletedNotification(asset.AssetName, asset.AssetTag, notes)
	utilTranslations := messages.GetMaintenanceRecordNotificationTranslations(titleKey, messageKey, params)

	// Convert to domain translations
	translations := make([]domain.CreateNotificationTranslationPayload, len(utilTranslations))
	for i, t := range utilTranslations {
		translations[i] = domain.CreateNotificationTranslationPayload{
			LangCode: t.LangCode,
			Title:    t.Title,
			Message:  t.Message,
		}
	}

	notificationPayload := &domain.CreateNotificationPayload{
		UserID:            *asset.AssignedToID,
		RelatedEntityType: stringPtr("maintenance_record"),
		RelatedEntityID:   &record.ID,
		RelatedAssetID:    &record.AssetID,
		Type:              domain.NotificationTypeMaintenance,
		Priority:          domain.NotificationPriorityNormal, // Completed = normal priority
		Translations:      translations,
	}

	_, err = s.NotificationService.CreateNotification(ctx, notificationPayload)
	if err != nil {
		// Log error but don't fail the operation
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

// sendMaintenanceFailedNotification sends notification for failed maintenance
func (s *Service) sendMaintenanceFailedNotification(ctx context.Context, record *domain.MaintenanceRecord, failureReason string) {
	if s.NotificationService == nil {
		return
	}

	// Get asset details
	asset, err := s.AssetService.GetAssetById(ctx, record.AssetID, "en-US") // Default lang
	if err != nil {
		return
	}

	if asset.AssignedToID == nil || *asset.AssignedToID == "" {
		return
	}

	titleKey, messageKey, params := messages.MaintenanceFailedNotification(asset.AssetName, asset.AssetTag, failureReason)
	utilTranslations := messages.GetMaintenanceRecordNotificationTranslations(titleKey, messageKey, params)

	// Convert to domain translations
	translations := make([]domain.CreateNotificationTranslationPayload, len(utilTranslations))
	for i, t := range utilTranslations {
		translations[i] = domain.CreateNotificationTranslationPayload{
			LangCode: t.LangCode,
			Title:    t.Title,
			Message:  t.Message,
		}
	}

	notificationPayload := &domain.CreateNotificationPayload{
		UserID:            *asset.AssignedToID,
		RelatedEntityType: stringPtr("maintenance_record"),
		RelatedEntityID:   &record.ID,
		RelatedAssetID:    &record.AssetID,
		Type:              domain.NotificationTypeMaintenance,
		Priority:          domain.NotificationPriorityHigh, // Failed = high priority
		Translations:      translations,
	}

	_, err = s.NotificationService.CreateNotification(ctx, notificationPayload)
	if err != nil {
		// Log error but don't fail the operation
	}
}
