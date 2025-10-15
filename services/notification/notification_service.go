package notification

import (
	"context"
	"log"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/client/fcm"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
)

// * Repository interface defines the contract for notification data operations
type Repository interface {
	// * MUTATION
	CreateNotification(ctx context.Context, payload *domain.Notification) (domain.Notification, error)
	UpdateNotification(ctx context.Context, notificationId string, payload *domain.UpdateNotificationPayload) (domain.Notification, error)
	DeleteNotification(ctx context.Context, notificationId string) error
	MarkNotificationAsRead(ctx context.Context, notificationId string, isRead bool) error
	MarkAllNotificationsAsRead(ctx context.Context, userId string) error

	// * QUERY
	GetNotificationsPaginated(ctx context.Context, params domain.NotificationParams, langCode string) ([]domain.Notification, error)
	GetNotificationsCursor(ctx context.Context, params domain.NotificationParams, langCode string) ([]domain.Notification, error)
	GetNotificationById(ctx context.Context, notificationId string) (domain.Notification, error)
	CheckNotificationExist(ctx context.Context, notificationId string) (bool, error)
	CountNotifications(ctx context.Context, params domain.NotificationParams) (int64, error)
	GetNotificationStatistics(ctx context.Context) (domain.NotificationStatistics, error)
}

// * UserRepository interface for getting user details including FCM token
type UserRepository interface {
	GetUserById(ctx context.Context, userId string) (domain.User, error)
}

// * NotificationService interface defines the contract for notification business operations
type NotificationService interface {
	// * MUTATION
	CreateNotification(ctx context.Context, payload *domain.CreateNotificationPayload) (domain.NotificationResponse, error)
	UpdateNotification(ctx context.Context, notificationId string, payload *domain.UpdateNotificationPayload) (domain.NotificationResponse, error)
	DeleteNotification(ctx context.Context, notificationId string) error
	MarkNotificationAsRead(ctx context.Context, notificationId string, isRead bool) error
	MarkAllNotificationsAsRead(ctx context.Context, userId string) error

	// * QUERY
	GetNotificationsPaginated(ctx context.Context, params domain.NotificationParams, langCode string) ([]domain.NotificationListResponse, int64, error)
	GetNotificationsCursor(ctx context.Context, params domain.NotificationParams, langCode string) ([]domain.NotificationListResponse, error)
	GetNotificationById(ctx context.Context, notificationId string, langCode string) (domain.NotificationResponse, error)
	CheckNotificationExists(ctx context.Context, notificationId string) (bool, error)
	CountNotifications(ctx context.Context, params domain.NotificationParams) (int64, error)
	GetNotificationStatistics(ctx context.Context) (domain.NotificationStatisticsResponse, error)
}

type Service struct {
	Repo      Repository
	UserRepo  UserRepository
	FCMClient *fcm.Client
}

// * Ensure Service implements NotificationService interface
var _ NotificationService = (*Service)(nil)

func NewService(r Repository, userRepo UserRepository, fcmClient *fcm.Client) NotificationService {
	return &Service{
		Repo:      r,
		UserRepo:  userRepo,
		FCMClient: fcmClient,
	}
}

// *===========================MUTATION===========================*
func (s *Service) CreateNotification(ctx context.Context, payload *domain.CreateNotificationPayload) (domain.NotificationResponse, error) {
	// * Prepare domain notification
	newNotification := domain.Notification{
		UserID:         payload.UserID,
		RelatedAssetID: payload.RelatedAssetID,
		Type:           payload.Type,
		IsRead:         false, // New notifications are always unread
		Translations:   make([]domain.NotificationTranslation, len(payload.Translations)),
	}

	// * Convert translation payloads to domain translations
	for i, translationPayload := range payload.Translations {
		newNotification.Translations[i] = domain.NotificationTranslation{
			LangCode: translationPayload.LangCode,
			Title:    translationPayload.Title,
			Message:  translationPayload.Message,
		}
	}

	createdNotification, err := s.Repo.CreateNotification(ctx, &newNotification)
	if err != nil {
		return domain.NotificationResponse{}, err
	}

	// * Send FCM notification asynchronously (non-blocking)
	go s.sendFCMNotification(ctx, &createdNotification)

	// * Convert to NotificationResponse using mapper
	return mapper.NotificationToResponse(&createdNotification, mapper.DefaultLangCode), nil
}

func (s *Service) UpdateNotification(ctx context.Context, notificationId string, payload *domain.UpdateNotificationPayload) (domain.NotificationResponse, error) {
	// * Check if notification exists
	_, err := s.Repo.GetNotificationById(ctx, notificationId)
	if err != nil {
		return domain.NotificationResponse{}, err
	}

	updatedNotification, err := s.Repo.UpdateNotification(ctx, notificationId, payload)
	if err != nil {
		return domain.NotificationResponse{}, err
	}

	// * Convert to NotificationResponse using mapper
	return mapper.NotificationToResponse(&updatedNotification, mapper.DefaultLangCode), nil
}

func (s *Service) DeleteNotification(ctx context.Context, notificationId string) error {
	err := s.Repo.DeleteNotification(ctx, notificationId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) MarkNotificationAsRead(ctx context.Context, notificationId string, isRead bool) error {
	// * Check if notification exists
	_, err := s.Repo.GetNotificationById(ctx, notificationId)
	if err != nil {
		return err
	}

	err = s.Repo.MarkNotificationAsRead(ctx, notificationId, isRead)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) MarkAllNotificationsAsRead(ctx context.Context, userId string) error {
	err := s.Repo.MarkAllNotificationsAsRead(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}

// *===========================QUERY===========================*
func (s *Service) GetNotificationsPaginated(ctx context.Context, params domain.NotificationParams, langCode string) ([]domain.NotificationListResponse, int64, error) {
	notifications, err := s.Repo.GetNotificationsPaginated(ctx, params, langCode)
	if err != nil {
		return nil, 0, err
	}

	// * Count total for pagination
	count, err := s.Repo.CountNotifications(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// Convert Notification to NotificationListResponse using mapper
	responses := mapper.NotificationsToListResponses(notifications, langCode)

	return responses, count, nil
}

func (s *Service) GetNotificationsCursor(ctx context.Context, params domain.NotificationParams, langCode string) ([]domain.NotificationListResponse, error) {
	notifications, err := s.Repo.GetNotificationsCursor(ctx, params, langCode)
	if err != nil {
		return nil, err
	}

	// Convert Notification to NotificationListResponse using mapper
	responses := mapper.NotificationsToListResponses(notifications, langCode)

	return responses, nil
}

func (s *Service) GetNotificationById(ctx context.Context, notificationId string, langCode string) (domain.NotificationResponse, error) {
	notification, err := s.Repo.GetNotificationById(ctx, notificationId)
	if err != nil {
		return domain.NotificationResponse{}, err
	}

	// * Convert to NotificationResponse using mapper
	return mapper.NotificationToResponse(&notification, langCode), nil
}

func (s *Service) CheckNotificationExists(ctx context.Context, notificationId string) (bool, error) {
	exists, err := s.Repo.CheckNotificationExist(ctx, notificationId)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Service) CountNotifications(ctx context.Context, params domain.NotificationParams) (int64, error) {
	count, err := s.Repo.CountNotifications(ctx, params)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Service) GetNotificationStatistics(ctx context.Context) (domain.NotificationStatisticsResponse, error) {
	stats, err := s.Repo.GetNotificationStatistics(ctx)
	if err != nil {
		return domain.NotificationStatisticsResponse{}, err
	}

	// Convert to NotificationStatisticsResponse using mapper
	return mapper.NotificationStatisticsToResponse(&stats), nil
}

// *===========================HELPER METHODS===========================*

// sendFCMNotification sends push notification via FCM to the user
func (s *Service) sendFCMNotification(ctx context.Context, notification *domain.Notification) {
	// * Skip if FCM client is not initialized
	if s.FCMClient == nil {
		log.Printf("FCM client not initialized, skipping FCM notification for notification ID: %s", notification.ID)
		return
	}

	log.Printf("Starting FCM notification send for notification ID: %s, user ID: %s", notification.ID, notification.UserID)

	// * Get user to retrieve FCM token and preferred language
	user, err := s.UserRepo.GetUserById(ctx, notification.UserID)
	if err != nil {
		// Log error but don't fail the notification creation
		log.Printf("Failed to get user for FCM notification (notification ID: %s, user ID: %s): %v", notification.ID, notification.UserID, err)
		return
	}

	// * Skip if user doesn't have FCM token
	if user.FCMToken == nil || *user.FCMToken == "" {
		log.Printf("User has no FCM token, skipping FCM notification for notification ID: %s, user ID: %s", notification.ID, notification.UserID)
		return
	}

	// * Get the appropriate translation based on user's preferred language
	var title, message string
	for _, translation := range notification.Translations {
		if translation.LangCode == user.PreferredLang {
			title = translation.Title
			message = translation.Message
			break
		}
	}

	// * Fallback to first translation if preferred language not found
	if title == "" && len(notification.Translations) > 0 {
		title = notification.Translations[0].Title
		message = notification.Translations[0].Message
	}

	// * Prepare FCM notification data
	fcmNotification := &fcm.PushNotification{
		Token: *user.FCMToken,
		Title: title,
		Body:  message,
		Data: map[string]string{
			"notification_id": notification.ID,
			"user_id":         notification.UserID,
			"type":            string(notification.Type),
			"is_read":         "false",
			"click_action":    "FLUTTER_NOTIFICATION_CLICK",
		},
	}

	// * Add related asset ID if available
	if notification.RelatedAssetID != nil {
		fcmNotification.Data["related_asset_id"] = *notification.RelatedAssetID
	}

	// * Send FCM notification
	_, err = s.FCMClient.SendToToken(ctx, fcmNotification)
	if err != nil {
		// Log error but don't fail the notification creation
		log.Printf("Failed to send FCM notification (notification ID: %s, user ID: %s): %v", notification.ID, notification.UserID, err)
	} else {
		log.Printf("Successfully sent FCM notification for notification ID: %s, user ID: %s", notification.ID, notification.UserID)
	}
}
