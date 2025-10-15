package issue_report

import (
	"context"
	"log"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/notification/messages"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/Rizz404/inventory-api/internal/utils"
)

// * Repository interface defines the contract for issue report data operations
type Repository interface {
	// * MUTATION
	CreateIssueReport(ctx context.Context, payload *domain.IssueReport) (domain.IssueReport, error)
	UpdateIssueReport(ctx context.Context, issueReportId string, payload *domain.UpdateIssueReportPayload) (domain.IssueReport, error)
	ResolveIssueReport(ctx context.Context, issueReportId string, resolvedBy string, resolutionNotes string) (domain.IssueReport, error)
	ReopenIssueReport(ctx context.Context, issueReportId string) (domain.IssueReport, error)
	DeleteIssueReport(ctx context.Context, issueReportId string) error

	// * QUERY
	GetIssueReportsPaginated(ctx context.Context, params domain.IssueReportParams, langCode string) ([]domain.IssueReport, error)
	GetIssueReportsCursor(ctx context.Context, params domain.IssueReportParams, langCode string) ([]domain.IssueReport, error)
	GetIssueReportById(ctx context.Context, issueReportId string) (domain.IssueReport, error)
	CheckIssueReportExist(ctx context.Context, issueReportId string) (bool, error)
	CountIssueReports(ctx context.Context, params domain.IssueReportParams) (int64, error)
	GetIssueReportStatistics(ctx context.Context) (domain.IssueReportStatistics, error)
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
	UpdateIssueReport(ctx context.Context, issueReportId string, payload *domain.UpdateIssueReportPayload) (domain.IssueReportResponse, error)
	ResolveIssueReport(ctx context.Context, issueReportId string, resolvedBy string, payload *domain.ResolveIssueReportPayload) (domain.IssueReportResponse, error)
	ReopenIssueReport(ctx context.Context, issueReportId string) (domain.IssueReportResponse, error)
	DeleteIssueReport(ctx context.Context, issueReportId string) error

	// * QUERY
	GetIssueReportsPaginated(ctx context.Context, params domain.IssueReportParams, langCode string) ([]domain.IssueReportListResponse, int64, error)
	GetIssueReportsCursor(ctx context.Context, params domain.IssueReportParams, langCode string) ([]domain.IssueReportListResponse, error)
	GetIssueReportById(ctx context.Context, issueReportId string, langCode string) (domain.IssueReportResponse, error)
	CheckIssueReportExists(ctx context.Context, issueReportId string) (bool, error)
	CountIssueReports(ctx context.Context, params domain.IssueReportParams) (int64, error)
	GetIssueReportStatistics(ctx context.Context) (domain.IssueReportStatisticsResponse, error)
}

type Service struct {
	Repo                Repository
	NotificationService NotificationService
	AssetService        AssetService
	UserRepo            UserRepository
}

// * Ensure Service implements IssueReportService interface
var _ IssueReportService = (*Service)(nil)

func NewService(r Repository, notificationService NotificationService, assetService AssetService, userRepo UserRepository) IssueReportService {
	return &Service{
		Repo:                r,
		NotificationService: notificationService,
		AssetService:        assetService,
		UserRepo:            userRepo,
	}
}

// *===========================MUTATION===========================*
func (s *Service) CreateIssueReport(ctx context.Context, payload *domain.CreateIssueReportPayload, reportedBy string) (domain.IssueReportResponse, error) {
	// * Prepare domain issue report
	newIssueReport := domain.IssueReport{
		AssetID:      payload.AssetID,
		ReportedBy:   reportedBy,
		ReportedDate: time.Now(),
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

	// * Send notification asynchronously
	go s.sendIssueReportedNotification(ctx, &createdIssueReport)

	// * Convert to IssueReportResponse using mapper
	return mapper.IssueReportToResponse(&createdIssueReport, mapper.DefaultLangCode), nil
}

func (s *Service) UpdateIssueReport(ctx context.Context, issueReportId string, payload *domain.UpdateIssueReportPayload) (domain.IssueReportResponse, error) {
	// * Check if issue report exists
	existingReport, err := s.Repo.GetIssueReportById(ctx, issueReportId)
	if err != nil {
		return domain.IssueReportResponse{}, err
	}

	// * Business logic: Check if the report is already closed and trying to change status
	if existingReport.Status == domain.IssueStatusClosed && payload.Status != nil && *payload.Status != domain.IssueStatusClosed {
		return domain.IssueReportResponse{}, domain.ErrBadRequestWithKey(utils.ErrIssueReportCannotReopenKey)
	}

	updatedIssueReport, err := s.Repo.UpdateIssueReport(ctx, issueReportId, payload)
	if err != nil {
		return domain.IssueReportResponse{}, err
	}

	// * Send notification asynchronously
	go s.sendIssueUpdatedNotification(ctx, &updatedIssueReport)

	// * Convert to IssueReportResponse using mapper
	return mapper.IssueReportToResponse(&updatedIssueReport, mapper.DefaultLangCode), nil
}

func (s *Service) ResolveIssueReport(ctx context.Context, issueReportId string, resolvedBy string, payload *domain.ResolveIssueReportPayload) (domain.IssueReportResponse, error) {
	// * Check if issue report exists
	existingReport, err := s.Repo.GetIssueReportById(ctx, issueReportId)
	if err != nil {
		return domain.IssueReportResponse{}, err
	}

	// * Business logic: Check if the report is already resolved
	if existingReport.Status == domain.IssueStatusResolved || existingReport.Status == domain.IssueStatusClosed {
		return domain.IssueReportResponse{}, domain.ErrBadRequestWithKey(utils.ErrIssueReportAlreadyResolvedKey)
	}

	resolvedIssueReport, err := s.Repo.ResolveIssueReport(ctx, issueReportId, resolvedBy, payload.ResolutionNotes)
	if err != nil {
		return domain.IssueReportResponse{}, err
	}

	// * Send notification asynchronously
	go s.sendIssueResolvedNotification(ctx, &resolvedIssueReport)

	// * Convert to IssueReportResponse using mapper
	return mapper.IssueReportToResponse(&resolvedIssueReport, mapper.DefaultLangCode), nil
}

func (s *Service) ReopenIssueReport(ctx context.Context, issueReportId string) (domain.IssueReportResponse, error) {
	// * Check if issue report exists
	existingReport, err := s.Repo.GetIssueReportById(ctx, issueReportId)
	if err != nil {
		return domain.IssueReportResponse{}, err
	}

	// * Business logic: Check if the report can be reopened
	if existingReport.Status == domain.IssueStatusClosed {
		return domain.IssueReportResponse{}, domain.ErrBadRequestWithKey(utils.ErrIssueReportCannotReopenKey)
	}

	if existingReport.Status == domain.IssueStatusOpen || existingReport.Status == domain.IssueStatusInProgress {
		return domain.IssueReportResponse{}, domain.ErrBadRequest("issue report is already open")
	}

	reopenedIssueReport, err := s.Repo.ReopenIssueReport(ctx, issueReportId)
	if err != nil {
		return domain.IssueReportResponse{}, err
	}

	// * Send notification asynchronously
	go s.sendIssueReopenedNotification(ctx, &reopenedIssueReport)

	// * Convert to IssueReportResponse using mapper
	return mapper.IssueReportToResponse(&reopenedIssueReport, mapper.DefaultLangCode), nil
}

func (s *Service) DeleteIssueReport(ctx context.Context, issueReportId string) error {
	err := s.Repo.DeleteIssueReport(ctx, issueReportId)
	if err != nil {
		return err
	}
	return nil
}

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
			UserID:         userId,
			RelatedAssetID: &issueReport.AssetID,
			Type:           domain.NotificationTypeIssueReport,
			Translations:   translations,
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
			UserID:         userId,
			RelatedAssetID: &issueReport.AssetID,
			Type:           domain.NotificationTypeIssueReport,
			Translations:   translations,
		}

		_, err = s.NotificationService.CreateNotification(ctx, notificationPayload)
		if err != nil {
			log.Printf("Failed to create issue updated notification for user ID: %s, issue report ID: %s: %v", userId, issueReport.ID, err)
		}
	}
}

// sendIssueResolvedNotification sends notification when an issue report is resolved
func (s *Service) sendIssueResolvedNotification(ctx context.Context, issueReport *domain.IssueReport) {
	// Skip if notification service is not available
	if s.NotificationService == nil {
		log.Printf("Notification service not available, skipping issue resolved notification for issue report ID: %s", issueReport.ID)
		return
	}

	// Get asset information
	asset, err := s.AssetService.GetAssetById(ctx, issueReport.AssetID, mapper.DefaultLangCode)
	if err != nil {
		log.Printf("Failed to get asset for issue resolved notification (issue report ID: %s, asset ID: %s): %v", issueReport.ID, issueReport.AssetID, err)
		return
	}

	log.Printf("Sending issue resolved notification for issue report ID: %s, asset ID: %s, asset tag: %s", issueReport.ID, asset.ID, asset.AssetTag)

	// Notify reporter
	userIds := []string{issueReport.ReportedBy}

	// Get resolution notes from the issue report translations
	var resolutionNotes string
	if issueReport.Translations != nil && len(issueReport.Translations) > 0 {
		for _, translation := range issueReport.Translations {
			if translation.ResolutionNotes != nil && *translation.ResolutionNotes != "" {
				resolutionNotes = *translation.ResolutionNotes
				break
			}
		}
	}

	// Get notification message keys and params
	titleKey, messageKey, params := messages.IssueResolvedNotification(asset.AssetName, asset.AssetTag, resolutionNotes)

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
			UserID:         userId,
			RelatedAssetID: &issueReport.AssetID,
			Type:           domain.NotificationTypeIssueReport,
			Translations:   translations,
		}

		_, err = s.NotificationService.CreateNotification(ctx, notificationPayload)
		if err != nil {
			log.Printf("Failed to create issue resolved notification for user ID: %s, issue report ID: %s: %v", userId, issueReport.ID, err)
		}
	}
}

// sendIssueReopenedNotification sends notification when an issue report is reopened
func (s *Service) sendIssueReopenedNotification(ctx context.Context, issueReport *domain.IssueReport) {
	// Skip if notification service is not available
	if s.NotificationService == nil {
		log.Printf("Notification service not available, skipping issue reopened notification for issue report ID: %s", issueReport.ID)
		return
	}

	// Get asset information
	asset, err := s.AssetService.GetAssetById(ctx, issueReport.AssetID, mapper.DefaultLangCode)
	if err != nil {
		log.Printf("Failed to get asset for issue reopened notification (issue report ID: %s, asset ID: %s): %v", issueReport.ID, issueReport.AssetID, err)
		return
	}

	log.Printf("Sending issue reopened notification for issue report ID: %s, asset ID: %s, asset tag: %s", issueReport.ID, asset.ID, asset.AssetTag)

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
			log.Printf("Failed to get admin users for issue reopened notification: %v", err)
			return
		}
		for _, admin := range admins {
			userIds = append(userIds, admin.ID)
		}
	}

	// Get notification message keys and params
	titleKey, messageKey, params := messages.IssueReopenedNotification(asset.AssetName, asset.AssetTag)

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
			UserID:         userId,
			RelatedAssetID: &issueReport.AssetID,
			Type:           domain.NotificationTypeIssueReport,
			Translations:   translations,
		}

		_, err = s.NotificationService.CreateNotification(ctx, notificationPayload)
		if err != nil {
			log.Printf("Failed to create issue reopened notification for user ID: %s, issue report ID: %s: %v", userId, issueReport.ID, err)
		}
	}
}

// *===========================QUERY===========================*
func (s *Service) GetIssueReportsPaginated(ctx context.Context, params domain.IssueReportParams, langCode string) ([]domain.IssueReportListResponse, int64, error) {
	listItems, err := s.Repo.GetIssueReportsPaginated(ctx, params, langCode)
	if err != nil {
		return nil, 0, err
	}

	// * Count total for pagination
	count, err := s.Repo.CountIssueReports(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// Convert IssueReport to IssueReportListResponse using mapper
	responses := mapper.IssueReportsToListResponses(listItems, langCode)

	return responses, count, nil
}

func (s *Service) GetIssueReportsCursor(ctx context.Context, params domain.IssueReportParams, langCode string) ([]domain.IssueReportListResponse, error) {
	listItems, err := s.Repo.GetIssueReportsCursor(ctx, params, langCode)
	if err != nil {
		return nil, err
	}

	// Convert IssueReport to IssueReportListResponse using mapper
	responses := mapper.IssueReportsToListResponses(listItems, langCode)

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
