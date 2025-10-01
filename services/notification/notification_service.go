package notification

import (
	"context"

	"github.com/Rizz404/inventory-api/domain"
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
	Repo Repository
}

// * Ensure Service implements NotificationService interface
var _ NotificationService = (*Service)(nil)

func NewService(r Repository) NotificationService {
	return &Service{
		Repo: r,
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
