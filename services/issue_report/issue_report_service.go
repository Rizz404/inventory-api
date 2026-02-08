package issue_report

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

// * Repository interface defines the contract for issue report data operations
type Repository interface {
	// * MUTATION
	CreateIssueReport(ctx context.Context, payload *domain.IssueReport) (domain.IssueReport, error)
	UpdateIssueReport(ctx context.Context, issueReportId string, payload *domain.UpdateIssueReportPayload) (domain.IssueReport, error)
	DeleteIssueReport(ctx context.Context, issueReportId string) error
	BulkCreateIssueReports(ctx context.Context, reports []domain.IssueReport) ([]domain.IssueReport, error)
	BulkDeleteIssueReports(ctx context.Context, reportIds []string) (domain.BulkDeleteIssueReports, error)
	AddIssueReportTranslations(ctx context.Context, issueReportId string, translations []domain.IssueReportTranslation) error

	// * QUERY
	GetIssueReportsPaginated(ctx context.Context, params domain.IssueReportParams, langCode string) ([]domain.IssueReport, error)
	GetIssueReportsCursor(ctx context.Context, params domain.IssueReportParams, langCode string) ([]domain.IssueReport, error)
	GetIssueReportById(ctx context.Context, issueReportId string) (domain.IssueReport, error)
	CheckIssueReportExist(ctx context.Context, issueReportId string) (bool, error)
	CountIssueReports(ctx context.Context, params domain.IssueReportParams) (int64, error)
	GetIssueReportStatistics(ctx context.Context) (domain.IssueReportStatistics, error)
	GetIssueReportsForExport(ctx context.Context, params domain.IssueReportParams, langCode string) ([]domain.IssueReport, error)
}

// * NotificationService interface for creating notifications
type NotificationService interface {
	CreateNotification(ctx context.Context, payload *domain.CreateNotificationPayload) (domain.NotificationResponse, error)
}

// * AssetService interface for getting asset information
type AssetService interface {
	GetAssetById(ctx context.Context, assetId string, langCode string) (domain.AssetResponse, error)
}

// * UserRepository interface for getting user details
type UserRepository interface {
	GetUsersPaginated(ctx context.Context, params domain.UserParams) ([]domain.User, error)
}

// * IssueReportService interface defines the contract for issue report business operations
type IssueReportService interface {
	// * MUTATION
	CreateIssueReport(ctx context.Context, payload *domain.CreateIssueReportPayload, reportedBy string) (domain.IssueReportResponse, error)
	UpdateIssueReport(ctx context.Context, issueReportId string, payload *domain.UpdateIssueReportPayload, langCode string) (domain.IssueReportResponse, error)
	DeleteIssueReport(ctx context.Context, issueReportId string) error
	BulkCreateIssueReports(ctx context.Context, payload *domain.BulkCreateIssueReportsPayload, reportedBy string) (domain.BulkCreateIssueReportsResponse, error)
	BulkDeleteIssueReports(ctx context.Context, payload *domain.BulkDeleteIssueReportsPayload) (domain.BulkDeleteIssueReportsResponse, error)

	// * QUERY
	GetIssueReportsPaginated(ctx context.Context, params domain.IssueReportParams, langCode string) ([]domain.IssueReportResponse, int64, error)
	GetIssueReportsCursor(ctx context.Context, params domain.IssueReportParams, langCode string) ([]domain.IssueReportResponse, error)
	GetIssueReportById(ctx context.Context, issueReportId string, langCode string) (domain.IssueReportResponse, error)
	CheckIssueReportExists(ctx context.Context, issueReportId string) (bool, error)
	CountIssueReports(ctx context.Context, params domain.IssueReportParams) (int64, error)
	GetIssueReportStatistics(ctx context.Context) (domain.IssueReportStatisticsResponse, error)
	ExportIssueReportList(ctx context.Context, payload domain.ExportIssueReportListPayload, params domain.IssueReportParams, langCode string) ([]byte, string, error)
}

type Service struct {
	Repo                Repository
	NotificationService NotificationService
	AssetService        AssetService
	UserRepo            UserRepository
	Translator          *gtranslate.Client
}

// * Ensure Service implements IssueReportService interface
var _ IssueReportService = (*Service)(nil)

func NewService(r Repository, notificationService NotificationService, assetService AssetService, userRepo UserRepository, translator *gtranslate.Client) IssueReportService {
	return &Service{
		Repo:                r,
		NotificationService: notificationService,
		AssetService:        assetService,
		UserRepo:            userRepo,
		Translator:          translator,
	}
}

// *===========================MUTATION===========================*
func (s *Service) CreateIssueReport(ctx context.Context, payload *domain.CreateIssueReportPayload, reportedBy string) (domain.IssueReportResponse, error) {
	// * Prepare domain issue report
	newIssueReport := domain.IssueReport{
		AssetID:      payload.AssetID,
		ReportedBy:   reportedBy,
		ReportedDate: time.Now().UTC(),
		IssueType:    payload.IssueType,
		Priority:     payload.Priority,
		Status:       domain.IssueStatusOpen, // New reports are always open
		Translations: make([]domain.IssueReportTranslation, len(payload.Translations)),
	}

	// * Convert translation payloads to domain translations
	for i, translationPayload := range payload.Translations {
		newIssueReport.Translations[i] = domain.IssueReportTranslation{
			LangCode:    translationPayload.LangCode,
			Title:       translationPayload.Title,
			Description: translationPayload.Description,
		}
	}

	createdIssueReport, err := s.Repo.CreateIssueReport(ctx, &newIssueReport)
	if err != nil {
		return domain.IssueReportResponse{}, err
	}

	// * Auto-translate missing languages in background
	go s.autoTranslateCreateIssueReportAsync(createdIssueReport.ID, payload.Translations)

	// * Send notification asynchronously
	go s.sendIssueReportedNotification(context.Background(), &createdIssueReport)

	// * Convert to IssueReportResponse using mapper
	return mapper.IssueReportToResponse(&createdIssueReport, mapper.DefaultLangCode), nil
}

func (s *Service) UpdateIssueReport(ctx context.Context, issueReportId string, payload *domain.UpdateIssueReportPayload, langCode string) (domain.IssueReportResponse, error) {
	// * Check if issue report exists
	_, err := s.Repo.GetIssueReportById(ctx, issueReportId)
	if err != nil {
		return domain.IssueReportResponse{}, err
	}

	updatedIssueReport, err := s.Repo.UpdateIssueReport(ctx, issueReportId, payload)
	if err != nil {
		return domain.IssueReportResponse{}, err
	}

	// * Auto-translate missing languages in background if translations updated
	if len(payload.Translations) > 0 {
		go s.autoTranslateUpdateIssueReportAsync(issueReportId, payload.Translations, updatedIssueReport.Translations)
	}

	// * Send notification asynchronously
	go s.sendIssueUpdatedNotification(context.Background(), &updatedIssueReport)

	// * Convert to IssueReportResponse using mapper with requested lang code
	return mapper.IssueReportToResponse(&updatedIssueReport, langCode), nil
}

func (s *Service) DeleteIssueReport(ctx context.Context, issueReportId string) error {
	err := s.Repo.DeleteIssueReport(ctx, issueReportId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) BulkCreateIssueReports(ctx context.Context, payload *domain.BulkCreateIssueReportsPayload, reportedBy string) (domain.BulkCreateIssueReportsResponse, error) {
	if payload == nil || len(payload.IssueReports) == 0 {
		return domain.BulkCreateIssueReportsResponse{}, domain.ErrBadRequest("issue reports payload is required")
	}

	// * Validate no duplicates in payload
	seenAssets := make(map[string]struct{})
	for _, item := range payload.IssueReports {
		if _, exists := seenAssets[item.AssetID]; exists {
			return domain.BulkCreateIssueReportsResponse{}, domain.ErrBadRequest("duplicate asset ID: " + item.AssetID)
		}
		seenAssets[item.AssetID] = struct{}{}
	}

	// * Validate all assets exist
	for assetID := range seenAssets {
		_, err := s.AssetService.GetAssetById(ctx, assetID, mapper.DefaultLangCode)
		if err != nil {
			return domain.BulkCreateIssueReportsResponse{}, domain.ErrNotFound("asset")
		}
	}

	// * Build domain issue reports
	issueReports := make([]domain.IssueReport, len(payload.IssueReports))
	for i, item := range payload.IssueReports {
		issueReports[i] = domain.IssueReport{
			AssetID:      item.AssetID,
			ReportedBy:   reportedBy,
			ReportedDate: time.Now().UTC(),
			IssueType:    item.IssueType,
			Priority:     item.Priority,
			Status:       domain.IssueStatusOpen,
			Translations: make([]domain.IssueReportTranslation, len(item.Translations)),
		}

		// * Convert translation payloads
		for j, translationPayload := range item.Translations {
			issueReports[i].Translations[j] = domain.IssueReportTranslation{
				LangCode:        translationPayload.LangCode,
				Title:           translationPayload.Title,
				Description:     translationPayload.Description,
				ResolutionNotes: nil,
			}
		}
	}

	// * Call repository bulk create
	created, err := s.Repo.BulkCreateIssueReports(ctx, issueReports)
	if err != nil {
		return domain.BulkCreateIssueReportsResponse{}, err
	}

	// * Send notifications asynchronously
	for i := range created {
		go s.sendIssueReportedNotification(context.Background(), &created[i])
	}

	// * Convert to responses
	response := domain.BulkCreateIssueReportsResponse{
		IssueReports: mapper.IssueReportsToResponses(created, mapper.DefaultLangCode),
	}
	return response, nil
}

func (s *Service) BulkDeleteIssueReports(ctx context.Context, payload *domain.BulkDeleteIssueReportsPayload) (domain.BulkDeleteIssueReportsResponse, error) {
	// * Validate that IDs are provided
	if len(payload.IDS) == 0 {
		return domain.BulkDeleteIssueReportsResponse{}, domain.ErrBadRequest("issue report IDs are required")
	}

	// * Perform bulk delete operation
	result, err := s.Repo.BulkDeleteIssueReports(ctx, payload.IDS)
	if err != nil {
		return domain.BulkDeleteIssueReportsResponse{}, err
	}

	// * Convert to response
	response := domain.BulkDeleteIssueReportsResponse{
		RequestedIDS: result.RequestedIDS,
		DeletedIDS:   result.DeletedIDS,
	}

	return response, nil
}

// *===========================HELPER METHODS===========================*

// sendIssueReportedNotification sends notification when a new issue report is created
func (s *Service) sendIssueReportedNotification(ctx context.Context, issueReport *domain.IssueReport) {
	// Skip if notification service is not available
	if s.NotificationService == nil {
		log.Printf("Notification service not available, skipping issue reported notification for issue report ID: %s", issueReport.ID)
		return
	}

	// Get asset information
	asset, err := s.AssetService.GetAssetById(ctx, issueReport.AssetID, mapper.DefaultLangCode)
	if err != nil {
		log.Printf("Failed to get asset for issue report notification (issue report ID: %s, asset ID: %s): %v", issueReport.ID, issueReport.AssetID, err)
		return
	}

	log.Printf("Sending issue reported notification for issue report ID: %s, asset ID: %s, asset tag: %s", issueReport.ID, asset.ID, asset.AssetTag)

	// Determine notification recipients
	var userIds []string
	if asset.AssignedToID != nil && *asset.AssignedToID != "" {
		// Notify assigned user
		userIds = []string{*asset.AssignedToID}
	} else {
		// Notify all admins
		adminRole := domain.RoleAdmin
		userParams := domain.UserParams{
			Filters: &domain.UserFilterOptions{
				Role: &adminRole,
			},
		}
		admins, err := s.UserRepo.GetUsersPaginated(ctx, userParams)
		if err != nil {
			log.Printf("Failed to get admin users for issue report notification: %v", err)
			return
		}
		for _, admin := range admins {
			userIds = append(userIds, admin.ID)
		}
	}

	// Get notification message keys and params
	titleKey, messageKey, params := messages.IssueReportedNotification(asset.AssetName, asset.AssetTag)

	// Get translations for all supported languages
	msgTranslations := messages.GetIssueReportNotificationTranslations(titleKey, messageKey, params)

	// Send notification to each recipient
	for _, userId := range userIds {
		// Convert to domain translations
		translations := make([]domain.CreateNotificationTranslationPayload, len(msgTranslations))
		for i, t := range msgTranslations {
			translations[i] = domain.CreateNotificationTranslationPayload{
				LangCode: t.LangCode,
				Title:    t.Title,
				Message:  t.Message,
			}
		}

		notificationPayload := &domain.CreateNotificationPayload{
			UserID:            userId,
			RelatedEntityType: utils.StringPtr("issue_report"),
			RelatedEntityID:   &issueReport.ID,
			RelatedAssetID:    &issueReport.AssetID,
			Type:              domain.NotificationTypeIssue,
			Priority:          determinePriorityFromIssue(issueReport),
			Translations:      translations,
		}

		_, err = s.NotificationService.CreateNotification(ctx, notificationPayload)
		if err != nil {
			log.Printf("Failed to create issue reported notification for user ID: %s, issue report ID: %s: %v", userId, issueReport.ID, err)
		}
	}
}

// sendIssueUpdatedNotification sends notification when an issue report is updated
func (s *Service) sendIssueUpdatedNotification(ctx context.Context, issueReport *domain.IssueReport) {
	// Skip if notification service is not available
	if s.NotificationService == nil {
		log.Printf("Notification service not available, skipping issue updated notification for issue report ID: %s", issueReport.ID)
		return
	}

	// Get asset information
	asset, err := s.AssetService.GetAssetById(ctx, issueReport.AssetID, mapper.DefaultLangCode)
	if err != nil {
		log.Printf("Failed to get asset for issue update notification (issue report ID: %s, asset ID: %s): %v", issueReport.ID, issueReport.AssetID, err)
		return
	}

	log.Printf("Sending issue updated notification for issue report ID: %s, asset ID: %s, asset tag: %s", issueReport.ID, asset.ID, asset.AssetTag)

	// Determine notification recipients: reporter and assigned user if different
	userIds := []string{issueReport.ReportedBy}
	if asset.AssignedToID != nil && *asset.AssignedToID != "" && *asset.AssignedToID != issueReport.ReportedBy {
		userIds = append(userIds, *asset.AssignedToID)
	}

	// Get notification message keys and params
	titleKey, messageKey, params := messages.IssueUpdatedNotification(asset.AssetName, asset.AssetTag)

	// Get translations for all supported languages
	msgTranslations := messages.GetIssueReportNotificationTranslations(titleKey, messageKey, params)

	// Send notification to each recipient
	for _, userId := range userIds {
		// Convert to domain translations
		translations := make([]domain.CreateNotificationTranslationPayload, len(msgTranslations))
		for i, t := range msgTranslations {
			translations[i] = domain.CreateNotificationTranslationPayload{
				LangCode: t.LangCode,
				Title:    t.Title,
				Message:  t.Message,
			}
		}

		notificationPayload := &domain.CreateNotificationPayload{
			UserID:            userId,
			RelatedEntityType: utils.StringPtr("issue_report"),
			RelatedEntityID:   &issueReport.ID,
			RelatedAssetID:    &issueReport.AssetID,
			Type:              domain.NotificationTypeIssue,
			Priority:          determinePriorityFromIssue(issueReport),
			Translations:      translations,
		}

		_, err = s.NotificationService.CreateNotification(ctx, notificationPayload)
		if err != nil {
			log.Printf("Failed to create issue updated notification for user ID: %s, issue report ID: %s: %v", userId, issueReport.ID, err)
		}
	}
}

// Helper function to determine notification priority based on issue priority
func determinePriorityFromIssue(issue *domain.IssueReport) domain.NotificationPriority {
	switch issue.Priority {
	case domain.PriorityCritical:
		return domain.NotificationPriorityUrgent
	case domain.PriorityHigh:
		return domain.NotificationPriorityHigh
	case domain.PriorityMedium:
		return domain.NotificationPriorityNormal
	case domain.PriorityLow:
		return domain.NotificationPriorityLow
	default:
		return domain.NotificationPriorityNormal
	}
}

// *===========================QUERY===========================*
func (s *Service) GetIssueReportsPaginated(ctx context.Context, params domain.IssueReportParams, langCode string) ([]domain.IssueReportResponse, int64, error) {
	listItems, err := s.Repo.GetIssueReportsPaginated(ctx, params, langCode)
	if err != nil {
		return nil, 0, err
	}

	// * Count total for pagination
	count, err := s.Repo.CountIssueReports(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// Convert IssueReport to IssueReportResponse using mapper (includes translations)
	responses := mapper.IssueReportsToResponses(listItems, langCode)

	return responses, count, nil
}

func (s *Service) GetIssueReportsCursor(ctx context.Context, params domain.IssueReportParams, langCode string) ([]domain.IssueReportResponse, error) {
	listItems, err := s.Repo.GetIssueReportsCursor(ctx, params, langCode)
	if err != nil {
		return nil, err
	}

	// Convert IssueReport to IssueReportResponse using mapper (includes translations)
	responses := mapper.IssueReportsToResponses(listItems, langCode)

	return responses, nil
}

func (s *Service) GetIssueReportById(ctx context.Context, issueReportId string, langCode string) (domain.IssueReportResponse, error) {
	issueReport, err := s.Repo.GetIssueReportById(ctx, issueReportId)
	if err != nil {
		return domain.IssueReportResponse{}, err
	}

	// * Convert to IssueReportResponse using mapper
	return mapper.IssueReportToResponse(&issueReport, langCode), nil
}

func (s *Service) CheckIssueReportExists(ctx context.Context, issueReportId string) (bool, error) {
	exists, err := s.Repo.CheckIssueReportExist(ctx, issueReportId)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Service) CountIssueReports(ctx context.Context, params domain.IssueReportParams) (int64, error) {
	count, err := s.Repo.CountIssueReports(ctx, params)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Service) GetIssueReportStatistics(ctx context.Context) (domain.IssueReportStatisticsResponse, error) {
	stats, err := s.Repo.GetIssueReportStatistics(ctx)
	if err != nil {
		return domain.IssueReportStatisticsResponse{}, err
	}

	// Convert to IssueReportStatisticsResponse using mapper
	return mapper.IssueReportStatisticsToResponse(&stats), nil
}

// *===========================ASYNC TRANSLATION===========================*

// autoTranslateCreateIssueReportAsync translates issue report to missing languages in background
func (s *Service) autoTranslateCreateIssueReportAsync(issueReportID string, userTranslations []domain.CreateIssueReportTranslationPayload) {
	if len(userTranslations) >= 3 {
		return // All languages provided, no need to translate
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Convert domain types to utils types
	utilsTranslations := make([]utils.IssueReportCreateTranslation, len(userTranslations))
	for i, t := range userTranslations {
		utilsTranslations[i] = utils.IssueReportCreateTranslation{
			LangCode:        t.LangCode,
			Title:           t.Title,
			Description:     t.Description,
			ResolutionNotes: nil,
		}
	}

	// Translate using utils helper
	translatedPayloads, err := utils.AutoTranslateIssueReportCreate(ctx, s.Translator, utilsTranslations)
	if err != nil {
		log.Printf("Failed to auto-translate issue report ID %s: %v", issueReportID, err)
		return
	}

	// Extract only new translations
	newTranslations := make([]domain.IssueReportTranslation, 0)
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
			newTranslations = append(newTranslations, domain.IssueReportTranslation{
				LangCode:        translated.LangCode,
				Title:           translated.Title,
				Description:     translated.Description,
				ResolutionNotes: translated.ResolutionNotes,
			})
		}
	}

	// Save auto-translated translations to database
	if len(newTranslations) > 0 {
		if err := s.Repo.AddIssueReportTranslations(ctx, issueReportID, newTranslations); err != nil {
			log.Printf("Failed to save auto-translated issue report translations for ID %s: %v", issueReportID, err)
		}
	}
}

// autoTranslateUpdateIssueReportAsync translates issue report updates to missing languages in background
func (s *Service) autoTranslateUpdateIssueReportAsync(issueReportID string, userUpdates []domain.UpdateIssueReportTranslationPayload, existingTranslations []domain.IssueReportTranslation) {
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
	utilsUpdates := make([]utils.IssueReportUpdateTranslation, len(userUpdates))
	for i, t := range userUpdates {
		utilsUpdates[i] = utils.IssueReportUpdateTranslation{
			LangCode:        t.LangCode,
			Title:           t.Title,
			Description:     t.Description,
			ResolutionNotes: t.ResolutionNotes,
		}
	}

	utilsExisting := make([]utils.IssueReportExistingTranslation, len(existingTranslations))
	for i, t := range existingTranslations {
		utilsExisting[i] = utils.IssueReportExistingTranslation{
			LangCode:        t.LangCode,
			Title:           t.Title,
			Description:     t.Description,
			ResolutionNotes: t.ResolutionNotes,
		}
	}

	// Translate using utils helper
	translatedPayloads, err := utils.AutoTranslateIssueReportUpdate(ctx, s.Translator, utilsUpdates, utilsExisting)
	if err != nil {
		log.Printf("Failed to auto-translate updated issue report ID %s: %v", issueReportID, err)
		return
	}

	// Extract only new translations (not in userUpdates)
	newTranslations := make([]domain.IssueReportTranslation, 0)
	for _, translated := range translatedPayloads {
		// Skip user-updated translations
		isUserUpdated := false
		for _, userUpdate := range userUpdates {
			if userUpdate.LangCode == translated.LangCode {
				isUserUpdated = true
				break
			}
		}

		if !isUserUpdated && (translated.Title != nil || translated.Description != nil || translated.ResolutionNotes != nil) {
			// Build full translation from update + existing
			var existing *domain.IssueReportTranslation
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

			var finalResolutionNotes *string
			if translated.ResolutionNotes != nil {
				finalResolutionNotes = translated.ResolutionNotes
			} else if existing != nil {
				finalResolutionNotes = existing.ResolutionNotes
			}

			if finalTitle != "" {
				newTranslations = append(newTranslations, domain.IssueReportTranslation{
					LangCode:        translated.LangCode,
					Title:           finalTitle,
					Description:     finalDescription,
					ResolutionNotes: finalResolutionNotes,
				})
			}
		}
	}

	// Save auto-translated translations to database
	if len(newTranslations) > 0 {
		if err := s.Repo.AddIssueReportTranslations(ctx, issueReportID, newTranslations); err != nil {
			log.Printf("Failed to save auto-translated issue report update translations for ID %s: %v", issueReportID, err)
		}
	}
}
