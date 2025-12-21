package asset_movement

import (
	"context"
	"log"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/notification/messages"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/Rizz404/inventory-api/internal/utils"
)

// * Repository interface defines the contract for asset movement data operations
type Repository interface {
	// * MUTATION
	CreateAssetMovement(ctx context.Context, payload *domain.AssetMovement) (domain.AssetMovement, error)
	UpdateAssetMovement(ctx context.Context, movementId string, payload *domain.UpdateAssetMovementPayload) (domain.AssetMovement, error)
	DeleteAssetMovement(ctx context.Context, movementId string) error
	BulkDeleteAssetMovements(ctx context.Context, movementIds []string) (domain.BulkDeleteAssetMovements, error)

	// * QUERY
	GetAssetMovementsPaginated(ctx context.Context, params domain.AssetMovementParams, langCode string) ([]domain.AssetMovement, error)
	GetAssetMovementsCursor(ctx context.Context, params domain.AssetMovementParams, langCode string) ([]domain.AssetMovement, error)
	GetAssetMovementById(ctx context.Context, movementId string) (domain.AssetMovement, error)
	GetAssetMovementsByAssetId(ctx context.Context, assetId string, params domain.AssetMovementParams) ([]domain.AssetMovement, error)
	CheckAssetMovementExist(ctx context.Context, movementId string) (bool, error)
	CountAssetMovements(ctx context.Context, params domain.AssetMovementParams) (int64, error)
	GetAssetMovementStatistics(ctx context.Context) (domain.AssetMovementStatistics, error)
	GetAssetMovementsForExport(ctx context.Context, params domain.AssetMovementParams, langCode string) ([]domain.AssetMovement, error)
}

// * AssetService interface for checking asset existence
type AssetService interface {
	CheckAssetExists(ctx context.Context, assetId string) (bool, error)
	GetAssetById(ctx context.Context, assetId string, langCode string) (domain.AssetResponse, error)
}

// * LocationService interface for checking location existence
type LocationService interface {
	CheckLocationExists(ctx context.Context, locationId string) (bool, error)
	GetLocationById(ctx context.Context, locationId string, langCode string) (domain.LocationResponse, error)
}

// * UserService interface for checking user existence
type UserService interface {
	CheckUserExists(ctx context.Context, userId string) (bool, error)
	GetUserById(ctx context.Context, userId string) (domain.UserResponse, error)
}

// * NotificationService interface for creating notifications
type NotificationService interface {
	CreateNotification(ctx context.Context, payload *domain.CreateNotificationPayload) (domain.NotificationResponse, error)
}
type AssetMovementService interface {
	// * MUTATION
	CreateAssetMovement(ctx context.Context, payload *domain.CreateAssetMovementPayload, movedBy string) (domain.AssetMovementResponse, error)
	UpdateAssetMovement(ctx context.Context, movementId string, payload *domain.UpdateAssetMovementPayload) (domain.AssetMovementResponse, error)
	DeleteAssetMovement(ctx context.Context, movementId string) error
	BulkDeleteAssetMovements(ctx context.Context, payload *domain.BulkDeleteAssetMovementsPayload) (domain.BulkDeleteAssetMovementsResponse, error)

	// * QUERY
	GetAssetMovementsPaginated(ctx context.Context, params domain.AssetMovementParams, langCode string) ([]domain.AssetMovementListResponse, int64, error)
	GetAssetMovementsCursor(ctx context.Context, params domain.AssetMovementParams, langCode string) ([]domain.AssetMovementListResponse, error)
	GetAssetMovementById(ctx context.Context, movementId string, langCode string) (domain.AssetMovementResponse, error)
	GetAssetMovementsByAssetId(ctx context.Context, assetId string, params domain.AssetMovementParams, langCode string) ([]domain.AssetMovementListResponse, error)
	CheckAssetMovementExists(ctx context.Context, movementId string) (bool, error)
	CountAssetMovements(ctx context.Context, params domain.AssetMovementParams) (int64, error)
	GetAssetMovementStatistics(ctx context.Context) (domain.AssetMovementStatisticsResponse, error)

	// * EXPORT
	ExportAssetMovementList(ctx context.Context, payload domain.ExportAssetMovementListPayload, params domain.AssetMovementParams, langCode string) ([]byte, string, error)
}

type Service struct {
	Repo                Repository
	AssetService        AssetService
	LocationService     LocationService
	UserService         UserService
	NotificationService NotificationService
}

// * Ensure Service implements AssetMovementService interface
var _ AssetMovementService = (*Service)(nil)

func NewService(r Repository, assetService AssetService, locationService LocationService, userService UserService, notificationService NotificationService) AssetMovementService {
	return &Service{
		Repo:                r,
		AssetService:        assetService,
		LocationService:     locationService,
		UserService:         userService,
		NotificationService: notificationService,
	}
}

// *===========================MUTATION===========================*
func (s *Service) CreateAssetMovement(ctx context.Context, payload *domain.CreateAssetMovementPayload, movedBy string) (domain.AssetMovementResponse, error) {
	// * Check if asset exists
	if assetExists, err := s.AssetService.CheckAssetExists(ctx, payload.AssetID); err != nil {
		return domain.AssetMovementResponse{}, err
	} else if !assetExists {
		return domain.AssetMovementResponse{}, domain.ErrNotFoundWithKey(utils.ErrAssetNotFoundKey)
	}

	// * Get current asset information to determine from location/user
	asset, err := s.AssetService.GetAssetById(ctx, payload.AssetID, "en-US")
	if err != nil {
		return domain.AssetMovementResponse{}, err
	}

	// * Validate destination (must have at least one: ToLocationID or ToUserID)
	if payload.ToLocationID == nil && payload.ToUserID == nil {
		return domain.AssetMovementResponse{}, domain.ErrBadRequestWithKey(utils.ErrAssetMovementNoChangeKey)
	}

	// * Check if destination location exists
	if payload.ToLocationID != nil {
		if locationExists, err := s.LocationService.CheckLocationExists(ctx, *payload.ToLocationID); err != nil {
			return domain.AssetMovementResponse{}, err
		} else if !locationExists {
			return domain.AssetMovementResponse{}, domain.ErrNotFoundWithKey(utils.ErrLocationNotFoundKey)
		}

		// * Check if asset is already in the same location
		if asset.LocationID != nil && *asset.LocationID == *payload.ToLocationID {
			return domain.AssetMovementResponse{}, domain.ErrBadRequestWithKey(utils.ErrAssetMovementSameLocationKey)
		}
	}

	// * Check if destination user exists
	if payload.ToUserID != nil {
		if userExists, err := s.UserService.CheckUserExists(ctx, *payload.ToUserID); err != nil {
			return domain.AssetMovementResponse{}, err
		} else if !userExists {
			return domain.AssetMovementResponse{}, domain.ErrNotFoundWithKey(utils.ErrUserNotFoundKey)
		}

		// * Check if asset is already assigned to the same user
		if asset.AssignedToID != nil && *asset.AssignedToID == *payload.ToUserID {
			return domain.AssetMovementResponse{}, domain.ErrBadRequestWithKey(utils.ErrAssetMovementNoChangeKey)
		}
	}

	// * Check if moved by user exists
	if movedByExists, err := s.UserService.CheckUserExists(ctx, movedBy); err != nil {
		return domain.AssetMovementResponse{}, err
	} else if !movedByExists {
		return domain.AssetMovementResponse{}, domain.ErrNotFoundWithKey(utils.ErrUserNotFoundKey)
	}

	// * Prepare domain asset movement
	newMovement := domain.AssetMovement{
		AssetID:        payload.AssetID,
		FromLocationID: asset.LocationID,
		ToLocationID:   payload.ToLocationID,
		FromUserID:     asset.AssignedToID,
		ToUserID:       payload.ToUserID,
		MovementDate:   time.Now(),
		MovedBy:        movedBy,
		Translations:   make([]domain.AssetMovementTranslation, len(payload.Translations)),
	}

	// * Convert translation payloads to domain translations
	for i, translationPayload := range payload.Translations {
		newMovement.Translations[i] = domain.AssetMovementTranslation{
			LangCode: translationPayload.LangCode,
			Notes:    &translationPayload.Notes,
		}
	}

	createdMovement, err := s.Repo.CreateAssetMovement(ctx, &newMovement)
	if err != nil {
		return domain.AssetMovementResponse{}, err
	}

	// * Send notifications asynchronously based on what changed
	locationChanged := payload.ToLocationID != nil && *payload.ToLocationID != "" &&
		(asset.LocationID == nil || *asset.LocationID != *payload.ToLocationID)
	userChanged := payload.ToUserID != nil && *payload.ToUserID != "" &&
		(asset.AssignedToID == nil || *asset.AssignedToID != *payload.ToUserID)

	if locationChanged {
		go s.sendLocationChangeNotification(context.Background(), &createdMovement, &asset)
	}

	if userChanged {
		go s.sendUserAssignmentNotification(context.Background(), &createdMovement, &asset)
	}

	// * Convert to AssetMovementResponse using mapper
	return mapper.AssetMovementToResponse(&createdMovement, mapper.DefaultLangCode), nil
}

func (s *Service) UpdateAssetMovement(ctx context.Context, movementId string, payload *domain.UpdateAssetMovementPayload) (domain.AssetMovementResponse, error) {
	// * Check if asset movement exists
	_, err := s.Repo.GetAssetMovementById(ctx, movementId)
	if err != nil {
		return domain.AssetMovementResponse{}, err
	}

	// * Check if destination location exists if being updated
	if payload.ToLocationID != nil && *payload.ToLocationID != "" {
		if locationExists, err := s.LocationService.CheckLocationExists(ctx, *payload.ToLocationID); err != nil {
			return domain.AssetMovementResponse{}, err
		} else if !locationExists {
			return domain.AssetMovementResponse{}, domain.ErrNotFoundWithKey(utils.ErrLocationNotFoundKey)
		}
	}

	// * Check if destination user exists if being updated
	if payload.ToUserID != nil && *payload.ToUserID != "" {
		if userExists, err := s.UserService.CheckUserExists(ctx, *payload.ToUserID); err != nil {
			return domain.AssetMovementResponse{}, err
		} else if !userExists {
			return domain.AssetMovementResponse{}, domain.ErrNotFoundWithKey(utils.ErrUserNotFoundKey)
		}
	}

	updatedMovement, err := s.Repo.UpdateAssetMovement(ctx, movementId, payload)
	if err != nil {
		return domain.AssetMovementResponse{}, err
	}

	// * Convert to AssetMovementResponse using mapper
	return mapper.AssetMovementToResponse(&updatedMovement, mapper.DefaultLangCode), nil
}

func (s *Service) DeleteAssetMovement(ctx context.Context, movementId string) error {
	err := s.Repo.DeleteAssetMovement(ctx, movementId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) BulkDeleteAssetMovements(ctx context.Context, payload *domain.BulkDeleteAssetMovementsPayload) (domain.BulkDeleteAssetMovementsResponse, error) {
	// * Validate that IDs are provided
	if len(payload.IDS) == 0 {
		return domain.BulkDeleteAssetMovementsResponse{}, domain.ErrBadRequest("asset movement IDs are required")
	}

	// * Perform bulk delete operation
	result, err := s.Repo.BulkDeleteAssetMovements(ctx, payload.IDS)
	if err != nil {
		return domain.BulkDeleteAssetMovementsResponse{}, err
	}

	// * Convert to response
	response := domain.BulkDeleteAssetMovementsResponse{
		RequestedIDS: result.RequestedIDS,
		DeletedIDS:   result.DeletedIDS,
	}

	return response, nil
}

// *===========================QUERY===========================*
func (s *Service) GetAssetMovementsPaginated(ctx context.Context, params domain.AssetMovementParams, langCode string) ([]domain.AssetMovementListResponse, int64, error) {
	listItems, err := s.Repo.GetAssetMovementsPaginated(ctx, params, langCode)
	if err != nil {
		return nil, 0, err
	}

	// * Count total for pagination
	count, err := s.Repo.CountAssetMovements(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// Convert AssetMovement to AssetMovementListResponse using mapper
	responses := mapper.AssetMovementsToListResponses(listItems, langCode)

	return responses, count, nil
}

func (s *Service) GetAssetMovementsCursor(ctx context.Context, params domain.AssetMovementParams, langCode string) ([]domain.AssetMovementListResponse, error) {
	listItems, err := s.Repo.GetAssetMovementsCursor(ctx, params, langCode)
	if err != nil {
		return nil, err
	}

	// Convert AssetMovement to AssetMovementListResponse using mapper
	responses := mapper.AssetMovementsToListResponses(listItems, langCode)

	return responses, nil
}

func (s *Service) GetAssetMovementById(ctx context.Context, movementId string, langCode string) (domain.AssetMovementResponse, error) {
	movement, err := s.Repo.GetAssetMovementById(ctx, movementId)
	if err != nil {
		return domain.AssetMovementResponse{}, err
	}

	// * Convert to AssetMovementResponse using mapper
	return mapper.AssetMovementToResponse(&movement, langCode), nil
}

func (s *Service) GetAssetMovementsByAssetId(ctx context.Context, assetId string, params domain.AssetMovementParams, langCode string) ([]domain.AssetMovementListResponse, error) {
	// * Check if asset exists
	if assetExists, err := s.AssetService.CheckAssetExists(ctx, assetId); err != nil {
		return nil, err
	} else if !assetExists {
		return nil, domain.ErrNotFoundWithKey(utils.ErrAssetNotFoundKey)
	}

	movements, err := s.Repo.GetAssetMovementsByAssetId(ctx, assetId, params)
	if err != nil {
		return nil, err
	}

	// * Convert to AssetMovementListResponse using mapper
	movementResponses := mapper.AssetMovementsToListResponses(movements, langCode)

	return movementResponses, nil
}

func (s *Service) CheckAssetMovementExists(ctx context.Context, movementId string) (bool, error) {
	exists, err := s.Repo.CheckAssetMovementExist(ctx, movementId)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Service) CountAssetMovements(ctx context.Context, params domain.AssetMovementParams) (int64, error) {
	count, err := s.Repo.CountAssetMovements(ctx, params)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Service) GetAssetMovementStatistics(ctx context.Context) (domain.AssetMovementStatisticsResponse, error) {
	stats, err := s.Repo.GetAssetMovementStatistics(ctx)
	if err != nil {
		return domain.AssetMovementStatisticsResponse{}, err
	}

	// Convert to AssetMovementStatisticsResponse using mapper
	return mapper.AssetMovementStatisticsToResponse(&stats), nil
}

// *===========================HELPER METHODS===========================*

// sendLocationChangeNotification sends notification when asset location changes
func (s *Service) sendLocationChangeNotification(ctx context.Context, movement *domain.AssetMovement, asset *domain.AssetResponse) {
	if s.NotificationService == nil {
		log.Printf("Notification service not available, skipping location change notification for asset ID: %s", asset.ID)
		return
	}

	// Get old location name
	oldLocation := "Unassigned"
	if movement.FromLocationID != nil && *movement.FromLocationID != "" {
		if loc, err := s.LocationService.GetLocationById(ctx, *movement.FromLocationID, "en-US"); err == nil {
			oldLocation = loc.LocationName
		}
	}

	// Get new location name
	newLocation := "Unknown"
	if movement.ToLocationID != nil && *movement.ToLocationID != "" {
		if loc, err := s.LocationService.GetLocationById(ctx, *movement.ToLocationID, "en-US"); err == nil {
			newLocation = loc.LocationName
		}
	}

	titleKey, messageKey, params := messages.AssetMovedNotification(asset.AssetName, asset.AssetTag, oldLocation, newLocation)
	utilTranslations := messages.GetAssetMovementNotificationTranslations(titleKey, messageKey, params)

	// Convert to domain translations
	translations := make([]domain.CreateNotificationTranslationPayload, len(utilTranslations))
	for i, t := range utilTranslations {
		translations[i] = domain.CreateNotificationTranslationPayload{
			LangCode: t.LangCode,
			Title:    t.Title,
			Message:  t.Message,
		}
	}

	entityType := "asset_movement"
	priority := domain.NotificationPriorityNormal

	// Determine recipient: new user if assigned, old user if unassigning, or admins
	var recipientUserID string
	if movement.ToUserID != nil && *movement.ToUserID != "" {
		recipientUserID = *movement.ToUserID
	} else if movement.FromUserID != nil && *movement.FromUserID != "" {
		recipientUserID = *movement.FromUserID
	} else {
		// Skip notification if no user is involved
		log.Printf("No user involved in location change, skipping notification for asset ID: %s", asset.ID)
		return
	}

	notificationPayload := &domain.CreateNotificationPayload{
		UserID:            recipientUserID,
		RelatedEntityType: &entityType,
		RelatedEntityID:   &movement.ID,
		RelatedAssetID:    &asset.ID,
		Type:              domain.NotificationTypeLocationChange,
		Priority:          priority,
		Translations:      translations,
	}

	_, err := s.NotificationService.CreateNotification(ctx, notificationPayload)
	if err != nil {
		log.Printf("Failed to create location change notification for asset ID: %s: %v", asset.ID, err)
	} else {
		log.Printf("Successfully created location change notification for asset ID: %s, user ID: %s", asset.ID, recipientUserID)
	}
}

// sendUserAssignmentNotification sends notification when asset is assigned to a user
func (s *Service) sendUserAssignmentNotification(ctx context.Context, movement *domain.AssetMovement, asset *domain.AssetResponse) {
	if s.NotificationService == nil {
		log.Printf("Notification service not available, skipping user assignment notification for asset ID: %s", asset.ID)
		return
	}

	if movement.ToUserID == nil || *movement.ToUserID == "" {
		log.Printf("No target user for assignment, skipping notification for asset ID: %s", asset.ID)
		return
	}

	// Get old user name
	oldUser := "Unassigned"
	if movement.FromUserID != nil && *movement.FromUserID != "" {
		if user, err := s.UserService.GetUserById(ctx, *movement.FromUserID); err == nil {
			oldUser = user.Name
		}
	}

	// Get new user name
	newUserName := "Unknown"
	if user, err := s.UserService.GetUserById(ctx, *movement.ToUserID); err == nil {
		newUserName = user.Name
	}

	titleKey, messageKey, params := messages.AssetUserAssignedNotification(asset.AssetName, asset.AssetTag, oldUser, newUserName)
	utilTranslations := messages.GetAssetMovementNotificationTranslations(titleKey, messageKey, params)

	// Convert to domain translations
	translations := make([]domain.CreateNotificationTranslationPayload, len(utilTranslations))
	for i, t := range utilTranslations {
		translations[i] = domain.CreateNotificationTranslationPayload{
			LangCode: t.LangCode,
			Title:    t.Title,
			Message:  t.Message,
		}
	}

	entityType := "asset_movement"
	priority := domain.NotificationPriorityNormal

	notificationPayload := &domain.CreateNotificationPayload{
		UserID:            *movement.ToUserID,
		RelatedEntityType: &entityType,
		RelatedEntityID:   &movement.ID,
		RelatedAssetID:    &asset.ID,
		Type:              domain.NotificationTypeMovement,
		Priority:          priority,
		Translations:      translations,
	}

	_, err := s.NotificationService.CreateNotification(ctx, notificationPayload)
	if err != nil {
		log.Printf("Failed to create user assignment notification for asset ID: %s: %v", asset.ID, err)
	} else {
		log.Printf("Successfully created user assignment notification for asset ID: %s, user ID: %s", asset.ID, *movement.ToUserID)
	}
}
