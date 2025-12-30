package location

import (
	"context"
	"log"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/notification/messages"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/Rizz404/inventory-api/internal/utils"
)

// * Repository interface defines the contract for location data operations
type Repository interface {
	// * MUTATION
	CreateLocation(ctx context.Context, payload *domain.Location) (domain.Location, error)
	BulkCreateLocations(ctx context.Context, locations []domain.Location) ([]domain.Location, error)
	UpdateLocation(ctx context.Context, locationId string, payload *domain.UpdateLocationPayload) (domain.Location, error)
	DeleteLocation(ctx context.Context, locationId string) error
	BulkDeleteLocations(ctx context.Context, locationIds []string) (domain.BulkDeleteLocations, error)

	// * QUERY
	GetLocationsPaginated(ctx context.Context, params domain.LocationParams, langCode string) ([]domain.Location, error)
	GetLocationsCursor(ctx context.Context, params domain.LocationParams, langCode string) ([]domain.Location, error)
	GetLocationById(ctx context.Context, locationId string) (domain.Location, error)
	GetLocationByCode(ctx context.Context, locationCode string) (domain.Location, error)
	CheckLocationExist(ctx context.Context, locationId string) (bool, error)
	CheckLocationCodeExist(ctx context.Context, locationCode string) (bool, error)
	CheckLocationCodeExistExcluding(ctx context.Context, locationCode string, excludeLocationId string) (bool, error)
	CountLocations(ctx context.Context, params domain.LocationParams) (int64, error)
	GetLocationStatistics(ctx context.Context) (domain.LocationStatistics, error)
}

// * LocationService interface defines the contract for location business operations
type LocationService interface {
	// * MUTATION
	CreateLocation(ctx context.Context, payload *domain.CreateLocationPayload) (domain.LocationResponse, error)
	BulkCreateLocations(ctx context.Context, payload *domain.BulkCreateLocationsPayload) (domain.BulkCreateLocationsResponse, error)
	UpdateLocation(ctx context.Context, locationId string, payload *domain.UpdateLocationPayload) (domain.LocationResponse, error)
	DeleteLocation(ctx context.Context, locationId string) error
	BulkDeleteLocations(ctx context.Context, payload *domain.BulkDeleteLocationsPayload) (domain.BulkDeleteLocationsResponse, error)

	// * QUERY
	GetLocationsPaginated(ctx context.Context, params domain.LocationParams, langCode string) ([]domain.LocationListResponse, int64, error)
	GetLocationsCursor(ctx context.Context, params domain.LocationParams, langCode string) ([]domain.LocationListResponse, error)
	GetLocationById(ctx context.Context, locationId string, langCode string) (domain.LocationResponse, error)
	GetLocationByCode(ctx context.Context, locationCode string, langCode string) (domain.LocationResponse, error)
	CheckLocationExists(ctx context.Context, locationId string) (bool, error)
	CheckLocationCodeExists(ctx context.Context, locationCode string) (bool, error)
	CountLocations(ctx context.Context, params domain.LocationParams) (int64, error)
	GetLocationStatistics(ctx context.Context) (domain.LocationStatisticsResponse, error)
}

// * NotificationService interface for creating notifications
type NotificationService interface {
	CreateNotification(ctx context.Context, payload *domain.CreateNotificationPayload) (domain.NotificationResponse, error)
}

// * UserRepository interface for getting user details
type UserRepository interface {
	GetUsersPaginated(ctx context.Context, params domain.UserParams) ([]domain.User, error)
}

type Service struct {
	Repo                Repository
	NotificationService NotificationService
	UserRepo            UserRepository
}

// * Ensure Service implements LocationService interface
var _ LocationService = (*Service)(nil)

func NewService(r Repository, notificationService NotificationService, userRepo UserRepository) LocationService {
	return &Service{
		Repo:                r,
		NotificationService: notificationService,
		UserRepo:            userRepo,
	}
}

// *===========================MUTATION===========================*
func (s *Service) CreateLocation(ctx context.Context, payload *domain.CreateLocationPayload) (domain.LocationResponse, error) {
	// * Check if location code already exists
	if codeExists, err := s.Repo.CheckLocationCodeExist(ctx, payload.LocationCode); err != nil {
		return domain.LocationResponse{}, err
	} else if codeExists {
		return domain.LocationResponse{}, domain.ErrConflictWithKey(utils.ErrLocationCodeExistsKey)
	}

	// * Prepare domain location
	newLocation := domain.Location{
		LocationCode: payload.LocationCode,
		Building:     payload.Building,
		Floor:        payload.Floor,
		Latitude:     payload.Latitude,
		Longitude:    payload.Longitude,
		Translations: make([]domain.LocationTranslation, len(payload.Translations)),
	}

	// * Convert translation payloads to domain translations
	for i, translationPayload := range payload.Translations {
		newLocation.Translations[i] = domain.LocationTranslation{
			LangCode:     translationPayload.LangCode,
			LocationName: translationPayload.LocationName,
		}
	}

	createdLocation, err := s.Repo.CreateLocation(ctx, &newLocation)
	if err != nil {
		return domain.LocationResponse{}, err
	}

	// * Send notification to all admin users
	s.sendLocationUpdatedNotificationToAdmins(ctx, &createdLocation)

	// * Convert to LocationResponse using mapper
	return mapper.LocationToResponse(&createdLocation, mapper.DefaultLangCode), nil
}

func (s *Service) BulkCreateLocations(ctx context.Context, payload *domain.BulkCreateLocationsPayload) (domain.BulkCreateLocationsResponse, error) {
	if payload == nil || len(payload.Locations) == 0 {
		return domain.BulkCreateLocationsResponse{}, domain.ErrBadRequest("locations payload is required")
	}

	codeSeen := make(map[string]struct{})
	for _, locPayload := range payload.Locations {
		if _, exists := codeSeen[locPayload.LocationCode]; exists {
			return domain.BulkCreateLocationsResponse{}, domain.ErrBadRequest("duplicate location code: " + locPayload.LocationCode)
		}
		codeSeen[locPayload.LocationCode] = struct{}{}
	}

	// Check all codes against database
	for code := range codeSeen {
		exists, err := s.Repo.CheckLocationCodeExist(ctx, code)
		if err != nil {
			return domain.BulkCreateLocationsResponse{}, err
		}
		if exists {
			return domain.BulkCreateLocationsResponse{}, domain.ErrConflictWithKey(utils.ErrLocationCodeExistsKey)
		}
	}

	locations := make([]domain.Location, len(payload.Locations))
	for i, locPayload := range payload.Locations {
		loc := domain.Location{
			LocationCode: locPayload.LocationCode,
			Building:     locPayload.Building,
			Floor:        locPayload.Floor,
			Latitude:     locPayload.Latitude,
			Longitude:    locPayload.Longitude,
			Translations: make([]domain.LocationTranslation, len(locPayload.Translations)),
		}

		for j, transPayload := range locPayload.Translations {
			loc.Translations[j] = domain.LocationTranslation{
				LangCode:     transPayload.LangCode,
				LocationName: transPayload.LocationName,
			}
		}
		locations[i] = loc
	}

	createdLocations, err := s.Repo.BulkCreateLocations(ctx, locations)
	if err != nil {
		return domain.BulkCreateLocationsResponse{}, err
	}

	// Send notifications for all created locations
	for i := range createdLocations {
		s.sendLocationUpdatedNotificationToAdmins(ctx, &createdLocations[i])
	}

	response := domain.BulkCreateLocationsResponse{
		Locations: mapper.LocationsToResponses(createdLocations, mapper.DefaultLangCode),
	}

	return response, nil
}

func (s *Service) UpdateLocation(ctx context.Context, locationId string, payload *domain.UpdateLocationPayload) (domain.LocationResponse, error) {
	// * Check if location exists
	_, err := s.Repo.GetLocationById(ctx, locationId)
	if err != nil {
		return domain.LocationResponse{}, err
	}

	// * Check location code uniqueness if being updated
	if payload.LocationCode != nil {
		if codeExists, err := s.Repo.CheckLocationCodeExistExcluding(ctx, *payload.LocationCode, locationId); err != nil {
			return domain.LocationResponse{}, err
		} else if codeExists {
			return domain.LocationResponse{}, domain.ErrConflictWithKey(utils.ErrLocationCodeExistsKey)
		}
	}

	updatedLocation, err := s.Repo.UpdateLocation(ctx, locationId, payload)
	if err != nil {
		return domain.LocationResponse{}, err
	}

	// * Send notification to all admin users
	s.sendLocationUpdatedNotificationToAdmins(ctx, &updatedLocation)

	// * Convert to LocationResponse using mapper
	return mapper.LocationToResponse(&updatedLocation, mapper.DefaultLangCode), nil
}

func (s *Service) DeleteLocation(ctx context.Context, locationId string) error {
	err := s.Repo.DeleteLocation(ctx, locationId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) BulkDeleteLocations(ctx context.Context, payload *domain.BulkDeleteLocationsPayload) (domain.BulkDeleteLocationsResponse, error) {
	// * Validate that IDs are provided
	if len(payload.IDS) == 0 {
		return domain.BulkDeleteLocationsResponse{}, domain.ErrBadRequestWithKey(utils.ErrLocationIDRequiredKey)
	}

	// * Perform bulk delete operation
	result, err := s.Repo.BulkDeleteLocations(ctx, payload.IDS)
	if err != nil {
		return domain.BulkDeleteLocationsResponse{}, err
	}

	// * Convert to response
	response := domain.BulkDeleteLocationsResponse{
		RequestedIDS: result.RequestedIDS,
		DeletedIDS:   result.DeletedIDS,
	}

	return response, nil
}

// *===========================QUERY===========================*
func (s *Service) GetLocationsPaginated(ctx context.Context, params domain.LocationParams, langCode string) ([]domain.LocationListResponse, int64, error) {
	locations, err := s.Repo.GetLocationsPaginated(ctx, params, langCode)
	if err != nil {
		return nil, 0, err
	}

	// * Count total for pagination
	count, err := s.Repo.CountLocations(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// Convert Location to LocationListResponse using mapper
	responses := mapper.LocationsToListResponses(locations, langCode)

	return responses, count, nil
}

func (s *Service) GetLocationsCursor(ctx context.Context, params domain.LocationParams, langCode string) ([]domain.LocationListResponse, error) {
	locations, err := s.Repo.GetLocationsCursor(ctx, params, langCode)
	if err != nil {
		return nil, err
	}

	// Convert Location to LocationListResponse using mapper
	responses := mapper.LocationsToListResponses(locations, langCode)

	return responses, nil
}

func (s *Service) GetLocationById(ctx context.Context, locationId string, langCode string) (domain.LocationResponse, error) {
	location, err := s.Repo.GetLocationById(ctx, locationId)
	if err != nil {
		return domain.LocationResponse{}, err
	}

	// * Convert to LocationResponse using mapper
	return mapper.LocationToResponse(&location, langCode), nil
}

func (s *Service) GetLocationByCode(ctx context.Context, locationCode string, langCode string) (domain.LocationResponse, error) {
	location, err := s.Repo.GetLocationByCode(ctx, locationCode)
	if err != nil {
		return domain.LocationResponse{}, err
	}

	// * Convert to LocationResponse using mapper
	return mapper.LocationToResponse(&location, langCode), nil
}

func (s *Service) CheckLocationExists(ctx context.Context, locationId string) (bool, error) {
	exists, err := s.Repo.CheckLocationExist(ctx, locationId)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Service) CheckLocationCodeExists(ctx context.Context, locationCode string) (bool, error) {
	exists, err := s.Repo.CheckLocationCodeExist(ctx, locationCode)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Service) CountLocations(ctx context.Context, params domain.LocationParams) (int64, error) {
	count, err := s.Repo.CountLocations(ctx, params)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Service) GetLocationStatistics(ctx context.Context) (domain.LocationStatisticsResponse, error) {
	stats, err := s.Repo.GetLocationStatistics(ctx)
	if err != nil {
		return domain.LocationStatisticsResponse{}, err
	}

	// Convert to LocationStatisticsResponse using mapper
	return mapper.LocationStatisticsToResponse(&stats), nil
}

// *===========================HELPER METHODS===========================*

// sendLocationUpdatedNotificationToAdmins sends notification for location update to all admin users
func (s *Service) sendLocationUpdatedNotificationToAdmins(ctx context.Context, location *domain.Location) {
	if s.NotificationService == nil {
		log.Printf("Notification service not available, skipping location updated notification for location ID: %s", location.ID)
		return
	}

	if s.UserRepo == nil {
		log.Printf("User repository not available, skipping location updated notification for location ID: %s", location.ID)
		return
	}

	log.Printf("Sending location updated notification to admins for location ID: %s, location code: %s", location.ID, location.LocationCode)

	// Get location name in default language
	locationName := ""
	for _, translation := range location.Translations {
		if translation.LangCode == "en-US" {
			locationName = translation.LocationName
			break
		}
	}
	if locationName == "" && len(location.Translations) > 0 {
		locationName = location.Translations[0].LocationName
	}

	// Get all admin users
	adminRole := domain.RoleAdmin
	userParams := domain.UserParams{
		Filters: &domain.UserFilterOptions{
			Role: &adminRole,
		},
	}
	admins, err := s.UserRepo.GetUsersPaginated(ctx, userParams)
	if err != nil {
		log.Printf("Failed to get admin users for location updated notification: %v", err)
		return
	}

	if len(admins) == 0 {
		log.Printf("No admin users found, skipping location updated notification for location ID: %s", location.ID)
		return
	}

	// Send notification to each admin
	for _, admin := range admins {
		s.sendLocationUpdatedNotification(ctx, location.ID, locationName, admin.ID)
	}

	log.Printf("Successfully sent location updated notification to %d admin(s) for location ID: %s", len(admins), location.ID)
}

// sendLocationUpdatedNotification sends notification for location update to a specific user
func (s *Service) sendLocationUpdatedNotification(ctx context.Context, locationID, locationName, userID string) {
	titleKey, messageKey, params := messages.LocationUpdatedNotification(locationName)
	utilTranslations := messages.GetLocationNotificationTranslations(titleKey, messageKey, params)

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
		UserID:            userID,
		RelatedEntityType: utils.StringPtr("location"),
		RelatedEntityID:   utils.StringPtr(locationID),
		Type:              domain.NotificationTypeLocationChange,
		Priority:          domain.NotificationPriorityLow, // Location change = low priority
		Translations:      translations,
	}

	_, err := s.NotificationService.CreateNotification(ctx, notificationPayload)
	if err != nil {
		log.Printf("Failed to create location updated notification for user ID: %s: %v", userID, err)
	} else {
		log.Printf("Successfully created location updated notification for user ID: %s", userID)
	}
}
