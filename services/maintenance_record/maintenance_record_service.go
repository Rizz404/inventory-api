package maintenance_record

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/client/gtranslate"
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
	BulkCreateRecords(ctx context.Context, records []domain.MaintenanceRecord) ([]domain.MaintenanceRecord, error)
	BulkDeleteRecords(ctx context.Context, recordIds []string) (domain.BulkDeleteMaintenanceRecords, error)
	AddMaintenanceRecordTranslations(ctx context.Context, recordId string, translations []domain.MaintenanceRecordTranslation) error

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
	UpdateMaintenanceRecord(ctx context.Context, recordId string, payload *domain.UpdateMaintenanceRecordPayload, langCode string) (domain.MaintenanceRecordResponse, error)
	DeleteMaintenanceRecord(ctx context.Context, recordId string) error
	BulkCreateMaintenanceRecords(ctx context.Context, payload *domain.BulkCreateMaintenanceRecordsPayload, performedBy string) (domain.BulkCreateMaintenanceRecordsResponse, error)
	BulkDeleteMaintenanceRecords(ctx context.Context, payload *domain.BulkDeleteMaintenanceRecordsPayload) (domain.BulkDeleteMaintenanceRecordsResponse, error)
	GetMaintenanceRecordsPaginated(ctx context.Context, params domain.MaintenanceRecordParams, langCode string) ([]domain.MaintenanceRecordResponse, int64, error)
	GetMaintenanceRecordsCursor(ctx context.Context, params domain.MaintenanceRecordParams, langCode string) ([]domain.MaintenanceRecordResponse, error)
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
	Translator          *gtranslate.Client
}

var _ MaintenanceRecordService = (*Service)(nil)

func NewService(r Repository, assetSvc AssetService, userSvc UserService, notificationSvc NotificationService, translator *gtranslate.Client) MaintenanceRecordService {
	return &Service{Repo: r, AssetService: assetSvc, UserService: userSvc, NotificationService: notificationSvc, Translator: translator}
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
	maintenanceDate, err := time.ParseInLocation("2006-01-02", payload.MaintenanceDate, time.UTC)
	if err != nil {
		return domain.MaintenanceRecordResponse{}, domain.ErrBadRequestWithKey(utils.ErrMaintenanceRecordDateRequiredKey)
	}

	// Parse completion date if provided
	var completionDate *time.Time
	if payload.CompletionDate != nil && *payload.CompletionDate != "" {
		parsed, err := time.ParseInLocation("2006-01-02", *payload.CompletionDate, time.UTC)
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

	// * Auto-translate missing languages in background
	go s.autoTranslateCreateMaintenanceRecordAsync(created.ID, payload.Translations)

	// Send notification for completed maintenance
	s.sendMaintenanceCompletedNotification(ctx, &created)

	return mapper.MaintenanceRecordToResponse(&created, mapper.DefaultLangCode), nil
}

func (s *Service) UpdateMaintenanceRecord(ctx context.Context, recordId string, payload *domain.UpdateMaintenanceRecordPayload, langCode string) (domain.MaintenanceRecordResponse, error) {
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

	// * Auto-translate missing languages in background if translations updated
	if len(payload.Translations) > 0 {
		go s.autoTranslateUpdateMaintenanceRecordAsync(recordId, payload.Translations, updated.Translations)
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

	return mapper.MaintenanceRecordToResponse(&updated, langCode), nil
}

func (s *Service) DeleteMaintenanceRecord(ctx context.Context, recordId string) error {
	return s.Repo.DeleteRecord(ctx, recordId)
}

func (s *Service) BulkCreateMaintenanceRecords(ctx context.Context, payload *domain.BulkCreateMaintenanceRecordsPayload, performedBy string) (domain.BulkCreateMaintenanceRecordsResponse, error) {
	if payload == nil || len(payload.MaintenanceRecords) == 0 {
		return domain.BulkCreateMaintenanceRecordsResponse{}, domain.ErrBadRequest("maintenance records payload is required")
	}

	// * Validate all assets exist
	assetMap := make(map[string]struct{})
	for _, item := range payload.MaintenanceRecords {
		if _, exists := assetMap[item.AssetID]; !exists {
			assetMap[item.AssetID] = struct{}{}
			if exists, err := s.AssetService.CheckAssetExists(ctx, item.AssetID); err != nil {
				return domain.BulkCreateMaintenanceRecordsResponse{}, err
			} else if !exists {
				return domain.BulkCreateMaintenanceRecordsResponse{}, domain.ErrNotFoundWithKey(utils.ErrAssetNotFoundKey)
			}
		}
	}

	// * Build domain records
	records := make([]domain.MaintenanceRecord, len(payload.MaintenanceRecords))
	for i, item := range payload.MaintenanceRecords {
		maintenanceDate, err := time.ParseInLocation("2006-01-02", item.MaintenanceDate, time.UTC)
		if err != nil {
			return domain.BulkCreateMaintenanceRecordsResponse{}, domain.ErrBadRequest("invalid maintenance_date format")
		}

		var completionDate *time.Time
		if item.CompletionDate != nil {
			parsedDate, err := time.ParseInLocation("2006-01-02", *item.CompletionDate, time.UTC)
			if err != nil {
				return domain.BulkCreateMaintenanceRecordsResponse{}, domain.ErrBadRequest("invalid completion_date format")
			}
			completionDate = &parsedDate
		}

		// * Determine performer
		var performerPtr *string
		if item.PerformedByUser != nil && *item.PerformedByUser != "" {
			performerPtr = item.PerformedByUser
		} else if performedBy != "" {
			performerPtr = &performedBy
		}

		// * Validate performer if present
		if performerPtr != nil {
			if exists, err := s.UserService.CheckUserExists(ctx, *performerPtr); err != nil {
				return domain.BulkCreateMaintenanceRecordsResponse{}, err
			} else if !exists {
				return domain.BulkCreateMaintenanceRecordsResponse{}, domain.ErrNotFoundWithKey(utils.ErrUserNotFoundKey)
			}
		}

		records[i] = domain.MaintenanceRecord{
			ScheduleID:        item.ScheduleID,
			AssetID:           item.AssetID,
			MaintenanceDate:   maintenanceDate,
			CompletionDate:    completionDate,
			DurationMinutes:   item.DurationMinutes,
			PerformedByUser:   performerPtr,
			PerformedByVendor: item.PerformedByVendor,
			Result:            item.Result,
			ActualCost:        item.ActualCost,
			Translations:      make([]domain.MaintenanceRecordTranslation, len(item.Translations)),
		}

		// * Convert translations
		for j, t := range item.Translations {
			records[i].Translations[j] = domain.MaintenanceRecordTranslation{
				LangCode: t.LangCode,
				Title:    t.Title,
				Notes:    t.Notes,
			}
		}
	}

	// * Call repository bulk create
	created, err := s.Repo.BulkCreateRecords(ctx, records)
	if err != nil {
		return domain.BulkCreateMaintenanceRecordsResponse{}, err
	}

	// * Send notifications asynchronously
	for i := range created {
		go s.sendMaintenanceCompletedNotification(ctx, &created[i])
	}

	// * Convert to responses
	response := domain.BulkCreateMaintenanceRecordsResponse{
		MaintenanceRecords: mapper.MaintenanceRecordsToResponses(created, mapper.DefaultLangCode),
	}
	return response, nil
}

func (s *Service) BulkDeleteMaintenanceRecords(ctx context.Context, payload *domain.BulkDeleteMaintenanceRecordsPayload) (domain.BulkDeleteMaintenanceRecordsResponse, error) {
	// * Validate that IDs are provided
	if len(payload.IDS) == 0 {
		return domain.BulkDeleteMaintenanceRecordsResponse{}, domain.ErrBadRequest("maintenance record IDs are required")
	}

	// * Perform bulk delete operation
	result, err := s.Repo.BulkDeleteRecords(ctx, payload.IDS)
	if err != nil {
		return domain.BulkDeleteMaintenanceRecordsResponse{}, err
	}

	// * Convert to response
	response := domain.BulkDeleteMaintenanceRecordsResponse{
		RequestedIDS: result.RequestedIDS,
		DeletedIDS:   result.DeletedIDS,
	}

	return response, nil
}

func (s *Service) GetMaintenanceRecordsPaginated(ctx context.Context, params domain.MaintenanceRecordParams, langCode string) ([]domain.MaintenanceRecordResponse, int64, error) {
	records, err := s.Repo.GetRecordsPaginated(ctx, params, langCode)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.Repo.CountRecords(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	recordResponses := mapper.MaintenanceRecordsToResponses(records, langCode)

	return recordResponses, count, nil
}

func (s *Service) GetMaintenanceRecordsCursor(ctx context.Context, params domain.MaintenanceRecordParams, langCode string) ([]domain.MaintenanceRecordResponse, error) {
	records, err := s.Repo.GetRecordsCursor(ctx, params, langCode)
	if err != nil {
		return nil, err
	}

	recordResponses := mapper.MaintenanceRecordsToResponses(records, langCode)

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
		RelatedEntityType: utils.StringPtr("maintenance_record"),
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
		RelatedEntityType: utils.StringPtr("maintenance_record"),
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

// *===========================ASYNC TRANSLATION===========================*

// autoTranslateCreateMaintenanceRecordAsync translates maintenance record to missing languages in background
func (s *Service) autoTranslateCreateMaintenanceRecordAsync(recordID string, userTranslations []domain.CreateMaintenanceRecordTranslationPayload) {
	if len(userTranslations) >= 3 {
		return // All languages provided, no need to translate
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Convert domain types to utils types
	utilsTranslations := make([]utils.MaintenanceRecordCreateTranslation, len(userTranslations))
	for i, t := range userTranslations {
		utilsTranslations[i] = utils.MaintenanceRecordCreateTranslation{
			LangCode: t.LangCode,
			Title:    t.Title,
			Notes:    t.Notes,
		}
	}

	// Translate using utils helper
	translatedPayloads, err := utils.AutoTranslateMaintenanceRecordCreate(ctx, s.Translator, utilsTranslations)
	if err != nil {
		log.Printf("Failed to auto-translate maintenance record ID %s: %v", recordID, err)
		return
	}

	// Extract only new translations
	newTranslations := make([]domain.MaintenanceRecordTranslation, 0)
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
			newTranslations = append(newTranslations, domain.MaintenanceRecordTranslation{
				LangCode: translated.LangCode,
				Title:    translated.Title,
				Notes:    translated.Notes,
			})
		}
	}

	// Save auto-translated translations to database
	if len(newTranslations) > 0 {
		if err := s.Repo.AddMaintenanceRecordTranslations(ctx, recordID, newTranslations); err != nil {
			log.Printf("Failed to save auto-translated maintenance record translations for ID %s: %v", recordID, err)
		}
	}
}

// autoTranslateUpdateMaintenanceRecordAsync translates maintenance record updates to missing languages in background
func (s *Service) autoTranslateUpdateMaintenanceRecordAsync(recordID string, userUpdates []domain.UpdateMaintenanceRecordTranslationPayload, existingTranslations []domain.MaintenanceRecordTranslation) {
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
	utilsUpdates := make([]utils.MaintenanceRecordUpdateTranslation, len(userUpdates))
	for i, t := range userUpdates {
		utilsUpdates[i] = utils.MaintenanceRecordUpdateTranslation{
			LangCode: t.LangCode,
			Title:    t.Title,
			Notes:    t.Notes,
		}
	}

	utilsExisting := make([]utils.MaintenanceRecordExistingTranslation, len(existingTranslations))
	for i, t := range existingTranslations {
		utilsExisting[i] = utils.MaintenanceRecordExistingTranslation{
			LangCode: t.LangCode,
			Title:    t.Title,
			Notes:    t.Notes,
		}
	}

	// Translate using utils helper
	translatedPayloads, err := utils.AutoTranslateMaintenanceRecordUpdate(ctx, s.Translator, utilsUpdates, utilsExisting)
	if err != nil {
		log.Printf("Failed to auto-translate updated maintenance record ID %s: %v", recordID, err)
		return
	}

	// Extract only new translations (not in userUpdates)
	newTranslations := make([]domain.MaintenanceRecordTranslation, 0)
	for _, translated := range translatedPayloads {
		// Skip user-updated translations
		isUserUpdated := false
		for _, userUpdate := range userUpdates {
			if userUpdate.LangCode == translated.LangCode {
				isUserUpdated = true
				break
			}
		}

		if !isUserUpdated && (translated.Title != nil || translated.Notes != nil) {
			// Build full translation from update + existing
			var existing *domain.MaintenanceRecordTranslation
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

			var finalNotes *string
			if translated.Notes != nil {
				finalNotes = translated.Notes
			} else if existing != nil {
				finalNotes = existing.Notes
			}

			if finalTitle != "" {
				newTranslations = append(newTranslations, domain.MaintenanceRecordTranslation{
					LangCode: translated.LangCode,
					Title:    finalTitle,
					Notes:    finalNotes,
				})
			}
		}
	}

	// Save auto-translated translations to database
	if len(newTranslations) > 0 {
		if err := s.Repo.AddMaintenanceRecordTranslations(ctx, recordID, newTranslations); err != nil {
			log.Printf("Failed to save auto-translated maintenance record update translations for ID %s: %v", recordID, err)
		}
	}
}
