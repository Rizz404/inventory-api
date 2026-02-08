package maintenance_schedule

import (
	"context"
	"log"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/client/gtranslate"
	"github.com/Rizz404/inventory-api/internal/notification/messages"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/Rizz404/inventory-api/internal/utils"
)

// Repository defines data operations for maintenance schedules
type Repository interface {
	// Schedule mutations
	CreateSchedule(ctx context.Context, payload *domain.MaintenanceSchedule) (domain.MaintenanceSchedule, error)
	UpdateSchedule(ctx context.Context, scheduleId string, payload *domain.UpdateMaintenanceSchedulePayload) (domain.MaintenanceSchedule, error)
	DeleteSchedule(ctx context.Context, scheduleId string) error
	BulkCreateSchedules(ctx context.Context, schedules []domain.MaintenanceSchedule) ([]domain.MaintenanceSchedule, error)
	BulkDeleteSchedules(ctx context.Context, scheduleIds []string) (domain.BulkDeleteMaintenanceSchedules, error)
	AddMaintenanceScheduleTranslations(ctx context.Context, scheduleId string, translations []domain.MaintenanceScheduleTranslation) error

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
	UpdateMaintenanceSchedule(ctx context.Context, scheduleId string, payload *domain.UpdateMaintenanceSchedulePayload, langCode string) (domain.MaintenanceScheduleResponse, error)
	DeleteMaintenanceSchedule(ctx context.Context, scheduleId string) error
	BulkCreateMaintenanceSchedules(ctx context.Context, payload *domain.BulkCreateMaintenanceSchedulesPayload, createdBy string) (domain.BulkCreateMaintenanceSchedulesResponse, error)
	BulkDeleteMaintenanceSchedules(ctx context.Context, payload *domain.BulkDeleteMaintenanceSchedulesPayload) (domain.BulkDeleteMaintenanceSchedulesResponse, error)
	GetMaintenanceSchedulesPaginated(ctx context.Context, params domain.MaintenanceScheduleParams, langCode string) ([]domain.MaintenanceScheduleResponse, int64, error)
	GetMaintenanceSchedulesCursor(ctx context.Context, params domain.MaintenanceScheduleParams, langCode string) ([]domain.MaintenanceScheduleResponse, error)
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
	Translator          *gtranslate.Client
}

var _ MaintenanceScheduleService = (*Service)(nil)

func NewService(r Repository, assetSvc AssetService, userSvc UserService, notificationSvc NotificationService, translator *gtranslate.Client) MaintenanceScheduleService {
	return &Service{Repo: r, AssetService: assetSvc, UserService: userSvc, NotificationService: notificationSvc, Translator: translator}
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

	// Parse next scheduled date in UTC
	nextScheduledDate, err := time.ParseInLocation("2006-01-02", payload.NextScheduledDate, time.UTC)
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

	// * Auto-translate missing languages in background
	go s.autoTranslateCreateMaintenanceScheduleAsync(created.ID, payload.Translations)

	// Send notification asynchronously
	go s.sendMaintenanceScheduledNotification(ctx, &created)

	return mapper.MaintenanceScheduleToResponse(&created, mapper.DefaultLangCode), nil
}

func (s *Service) UpdateMaintenanceSchedule(ctx context.Context, scheduleId string, payload *domain.UpdateMaintenanceSchedulePayload, langCode string) (domain.MaintenanceScheduleResponse, error) {
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

	// * Auto-translate missing languages in background if translations updated
	if len(payload.Translations) > 0 {
		go s.autoTranslateUpdateMaintenanceScheduleAsync(scheduleId, payload.Translations, updated.Translations)
	}

	return mapper.MaintenanceScheduleToResponse(&updated, langCode), nil
}

func (s *Service) DeleteMaintenanceSchedule(ctx context.Context, scheduleId string) error {
	return s.Repo.DeleteSchedule(ctx, scheduleId)
}

func (s *Service) BulkCreateMaintenanceSchedules(ctx context.Context, payload *domain.BulkCreateMaintenanceSchedulesPayload, createdBy string) (domain.BulkCreateMaintenanceSchedulesResponse, error) {
	if payload == nil || len(payload.MaintenanceSchedules) == 0 {
		return domain.BulkCreateMaintenanceSchedulesResponse{}, domain.ErrBadRequest("maintenance schedules payload is required")
	}

	// * Validate creator user exists
	if createdBy != "" {
		if exists, err := s.UserService.CheckUserExists(ctx, createdBy); err != nil {
			return domain.BulkCreateMaintenanceSchedulesResponse{}, err
		} else if !exists {
			return domain.BulkCreateMaintenanceSchedulesResponse{}, domain.ErrNotFoundWithKey(utils.ErrUserNotFoundKey)
		}
	}

	// * Validate all assets exist
	assetMap := make(map[string]struct{})
	for _, item := range payload.MaintenanceSchedules {
		if _, exists := assetMap[item.AssetID]; !exists {
			assetMap[item.AssetID] = struct{}{}
			if exists, err := s.AssetService.CheckAssetExists(ctx, item.AssetID); err != nil {
				return domain.BulkCreateMaintenanceSchedulesResponse{}, err
			} else if !exists {
				return domain.BulkCreateMaintenanceSchedulesResponse{}, domain.ErrNotFoundWithKey(utils.ErrAssetNotFoundKey)
			}
		}
	}

	// * Build domain schedules
	schedules := make([]domain.MaintenanceSchedule, len(payload.MaintenanceSchedules))
	for i, item := range payload.MaintenanceSchedules {
		nextScheduledDate, err := time.ParseInLocation("2006-01-02", item.NextScheduledDate, time.UTC)
		if err != nil {
			return domain.BulkCreateMaintenanceSchedulesResponse{}, domain.ErrBadRequestWithKey(utils.ErrMaintenanceScheduleDateRequiredKey)
		}

		isRecurring := false
		if item.IsRecurring != nil {
			isRecurring = *item.IsRecurring
		}

		autoComplete := false
		if item.AutoComplete != nil {
			autoComplete = *item.AutoComplete
		}

		// * Validate: if recurring, must have interval
		if isRecurring && (item.IntervalValue == nil || item.IntervalUnit == nil) {
			return domain.BulkCreateMaintenanceSchedulesResponse{}, domain.ErrBadRequest("recurring schedule must have interval_value and interval_unit")
		}

		schedules[i] = domain.MaintenanceSchedule{
			AssetID:           item.AssetID,
			MaintenanceType:   item.MaintenanceType,
			IsRecurring:       isRecurring,
			IntervalValue:     item.IntervalValue,
			IntervalUnit:      item.IntervalUnit,
			ScheduledTime:     item.ScheduledTime,
			NextScheduledDate: nextScheduledDate,
			State:             domain.StateActive,
			AutoComplete:      autoComplete,
			EstimatedCost:     item.EstimatedCost,
			CreatedBy:         createdBy,
			Translations:      make([]domain.MaintenanceScheduleTranslation, len(item.Translations)),
		}

		// * Convert translations
		for j, t := range item.Translations {
			schedules[i].Translations[j] = domain.MaintenanceScheduleTranslation{
				LangCode:    t.LangCode,
				Title:       t.Title,
				Description: t.Description,
			}
		}
	}

	// * Call repository bulk create
	created, err := s.Repo.BulkCreateSchedules(ctx, schedules)
	if err != nil {
		return domain.BulkCreateMaintenanceSchedulesResponse{}, err
	}

	// Send notifications for all created schedules
	for i := range created {
		go s.sendMaintenanceScheduledNotification(ctx, &created[i])
	}

	// * Convert to responses
	response := domain.BulkCreateMaintenanceSchedulesResponse{
		MaintenanceSchedules: mapper.MaintenanceSchedulesToResponses(created, mapper.DefaultLangCode),
	}
	return response, nil
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

func (s *Service) GetMaintenanceSchedulesPaginated(ctx context.Context, params domain.MaintenanceScheduleParams, langCode string) ([]domain.MaintenanceScheduleResponse, int64, error) {
	schedules, err := s.Repo.GetSchedulesPaginated(ctx, params, langCode)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.Repo.CountSchedules(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	schedulesResponses := mapper.MaintenanceSchedulesToResponses(schedules, langCode)

	return schedulesResponses, count, nil
}

func (s *Service) GetMaintenanceSchedulesCursor(ctx context.Context, params domain.MaintenanceScheduleParams, langCode string) ([]domain.MaintenanceScheduleResponse, error) {
	schedules, err := s.Repo.GetSchedulesCursor(ctx, params, langCode)
	if err != nil {
		return nil, err
	}

	schedulesResponses := mapper.MaintenanceSchedulesToResponses(schedules, langCode)

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

// *===========================HELPER METHODS===========================*

// sendMaintenanceScheduledNotification sends notification for new maintenance schedule
func (s *Service) sendMaintenanceScheduledNotification(ctx context.Context, schedule *domain.MaintenanceSchedule) {
	if s.NotificationService == nil {
		log.Printf("Notification service not available, skipping maintenance scheduled notification for schedule ID: %s", schedule.ID)
		return
	}

	// Get asset information
	asset, err := s.AssetService.GetAssetById(ctx, schedule.AssetID, mapper.DefaultLangCode)
	if err != nil {
		log.Printf("Failed to get asset for maintenance scheduled notification: %v", err)
		return
	}

	// Format scheduled date
	scheduledDate := schedule.NextScheduledDate.Format("2006-01-02")

	titleKey, messageKey, params := messages.MaintenanceScheduledNotification(asset.AssetName, asset.AssetTag, scheduledDate)
	utilTranslations := messages.GetMaintenanceScheduleNotificationTranslations(titleKey, messageKey, params)

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
		UserID:            schedule.CreatedBy,
		RelatedEntityType: utils.StringPtr("maintenance_schedule"),
		RelatedEntityID:   utils.StringPtr(schedule.ID),
		Type:              domain.NotificationTypeStatusChange,
		Translations:      translations,
	}

	_, err = s.NotificationService.CreateNotification(ctx, notificationPayload)
	if err != nil {
		log.Printf("Failed to create maintenance scheduled notification for schedule ID: %s: %v", schedule.ID, err)
	} else {
		log.Printf("Successfully created maintenance scheduled notification for schedule ID: %s", schedule.ID)
	}
}

// *===========================ASYNC TRANSLATION===========================*

// autoTranslateCreateMaintenanceScheduleAsync translates maintenance schedule to missing languages in background
func (s *Service) autoTranslateCreateMaintenanceScheduleAsync(scheduleID string, userTranslations []domain.CreateMaintenanceScheduleTranslationPayload) {
	if len(userTranslations) >= 3 {
		return // All languages provided, no need to translate
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Convert domain types to utils types
	utilsTranslations := make([]utils.MaintenanceScheduleCreateTranslation, len(userTranslations))
	for i, t := range userTranslations {
		utilsTranslations[i] = utils.MaintenanceScheduleCreateTranslation{
			LangCode:    t.LangCode,
			Title:       t.Title,
			Description: t.Description,
		}
	}

	// Translate using utils helper
	translatedPayloads, err := utils.AutoTranslateMaintenanceScheduleCreate(ctx, s.Translator, utilsTranslations)
	if err != nil {
		log.Printf("Failed to auto-translate maintenance schedule ID %s: %v", scheduleID, err)
		return
	}

	// Extract only new translations
	newTranslations := make([]domain.MaintenanceScheduleTranslation, 0)
	for _, translated := range translatedPayloads {
		// Skip user-provided translations
		isUserProvided := false
		for _, userTrans := range userTranslations {
			if userTrans.LangCode == translated.LangCode {
				isUserProvided = true
				break
			}
		}

		if !isUserProvided {
			newTranslations = append(newTranslations, domain.MaintenanceScheduleTranslation{
				LangCode:    translated.LangCode,
				Title:       translated.Title,
				Description: translated.Description,
			})
		}
	}

	// Save auto-translated translations to database
	if len(newTranslations) > 0 {
		if err := s.Repo.AddMaintenanceScheduleTranslations(ctx, scheduleID, newTranslations); err != nil {
			log.Printf("Failed to save auto-translated maintenance schedule translations for ID %s: %v", scheduleID, err)
		}
	}
}

// autoTranslateUpdateMaintenanceScheduleAsync translates maintenance schedule updates to missing languages in background
func (s *Service) autoTranslateUpdateMaintenanceScheduleAsync(scheduleID string, userUpdates []domain.UpdateMaintenanceScheduleTranslationPayload, existingTranslations []domain.MaintenanceScheduleTranslation) {
	if len(userUpdates) == 0 {
		return // No updates provided
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// If user updated all 3 languages, no need to auto-translate
	updatedLangCodes := make([]string, len(userUpdates))
	for i, t := range userUpdates {
		updatedLangCodes[i] = t.LangCode
	}

	if len(updatedLangCodes) >= 3 {
		return
	}

	// Convert domain types to utils types
	utilsUpdates := make([]utils.MaintenanceScheduleUpdateTranslation, len(userUpdates))
	for i, t := range userUpdates {
		utilsUpdates[i] = utils.MaintenanceScheduleUpdateTranslation{
			LangCode:    t.LangCode,
			Title:       t.Title,
			Description: t.Description,
		}
	}

	utilsExisting := make([]utils.MaintenanceScheduleExistingTranslation, len(existingTranslations))
	for i, t := range existingTranslations {
		utilsExisting[i] = utils.MaintenanceScheduleExistingTranslation{
			LangCode:    t.LangCode,
			Title:       t.Title,
			Description: t.Description,
		}
	}

	// Translate using utils helper
	translatedPayloads, err := utils.AutoTranslateMaintenanceScheduleUpdate(ctx, s.Translator, utilsUpdates, utilsExisting)
	if err != nil {
		log.Printf("Failed to auto-translate updated maintenance schedule ID %s: %v", scheduleID, err)
		return
	}

	// Extract only new translations (not in userUpdates)
	newTranslations := make([]domain.MaintenanceScheduleTranslation, 0)
	for _, translated := range translatedPayloads {
		// Skip user-updated translations
		isUserUpdated := false
		for _, userUpdate := range userUpdates {
			if userUpdate.LangCode == translated.LangCode {
				isUserUpdated = true
				break
			}
		}

		if !isUserUpdated && (translated.Title != nil || translated.Description != nil) {
			// Build full translation from update + existing
			var existing *domain.MaintenanceScheduleTranslation
			for _, e := range existingTranslations {
				if e.LangCode == translated.LangCode {
					existing = &e
					break
				}
			}

			finalTitle := ""
			if translated.Title != nil {
				finalTitle = *translated.Title
			} else if existing != nil {
				finalTitle = existing.Title
			}

			var finalDescription *string
			if translated.Description != nil {
				finalDescription = translated.Description
			} else if existing != nil && existing.Description != nil {
				finalDescription = existing.Description
			}

			if finalTitle != "" {
				newTranslations = append(newTranslations, domain.MaintenanceScheduleTranslation{
					LangCode:    translated.LangCode,
					Title:       finalTitle,
					Description: finalDescription,
				})
			}
		}
	}

	// Save auto-translated translations to database
	if len(newTranslations) > 0 {
		if err := s.Repo.AddMaintenanceScheduleTranslations(ctx, scheduleID, newTranslations); err != nil {
			log.Printf("Failed to save auto-translated maintenance schedule update translations for ID %s: %v", scheduleID, err)
		}
	}
}
