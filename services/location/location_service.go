package location

import (
	"context"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/Rizz404/inventory-api/internal/utils"
)

// * Repository interface defines the contract for location data operations
type Repository interface {
	// * MUTATION
	CreateLocation(ctx context.Context, payload *domain.Location) (domain.Location, error)
	UpdateLocation(ctx context.Context, locationId string, payload *domain.UpdateLocationPayload) (domain.Location, error)
	DeleteLocation(ctx context.Context, locationId string) error

	// * QUERY
	GetLocationsPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.LocationListItem, error)
	GetLocationsCursor(ctx context.Context, params query.Params, langCode string) ([]domain.LocationListItem, error)
	GetLocationById(ctx context.Context, locationId string) (domain.Location, error)
	GetLocationByCode(ctx context.Context, locationCode string) (domain.Location, error)
	GetLocationHierarchy(ctx context.Context, langCode string) ([]domain.LocationResponse, error)
	CheckLocationExist(ctx context.Context, locationId string) (bool, error)
	CheckLocationCodeExist(ctx context.Context, locationCode string) (bool, error)
	CheckLocationCodeExistExcluding(ctx context.Context, locationCode string, excludeLocationId string) (bool, error)
	CountLocations(ctx context.Context, params query.Params) (int64, error)
	GetLocationStatistics(ctx context.Context) (domain.LocationStatistics, error)
}

// * LocationService interface defines the contract for location business operations
type LocationService interface {
	// * MUTATION
	CreateLocation(ctx context.Context, payload *domain.CreateLocationPayload) (domain.LocationResponse, error)
	UpdateLocation(ctx context.Context, locationId string, payload *domain.UpdateLocationPayload) (domain.LocationResponse, error)
	DeleteLocation(ctx context.Context, locationId string) error

	// * QUERY
	GetLocationsPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.LocationListItemResponse, int64, error)
	GetLocationsCursor(ctx context.Context, params query.Params, langCode string) ([]domain.LocationListItemResponse, error)
	GetLocationById(ctx context.Context, locationId string, langCode string) (domain.LocationResponse, error)
	GetLocationByCode(ctx context.Context, locationCode string, langCode string) (domain.LocationResponse, error)
	GetLocationHierarchy(ctx context.Context, langCode string) ([]domain.LocationResponse, error)
	CheckLocationExists(ctx context.Context, locationId string) (bool, error)
	CheckLocationCodeExists(ctx context.Context, locationCode string) (bool, error)
	CountLocations(ctx context.Context, params query.Params) (int64, error)
	GetLocationStatistics(ctx context.Context) (domain.LocationStatisticsResponse, error)
}

type Service struct {
	Repo Repository
}

// * Ensure Service implements LocationService interface
var _ LocationService = (*Service)(nil)

func NewService(r Repository) LocationService {
	return &Service{
		Repo: r,
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

	// * Convert to LocationResponse using mapper
	return mapper.LocationToResponse(&createdLocation, mapper.DefaultLangCode), nil
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

// *===========================QUERY===========================*
func (s *Service) GetLocationsPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.LocationListItemResponse, int64, error) {
	locations, err := s.Repo.GetLocationsPaginated(ctx, params, langCode)
	if err != nil {
		return nil, 0, err
	}

	// * Count total for pagination
	count, err := s.Repo.CountLocations(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// Convert LocationListItem to LocationListItemResponse
	responses := make([]domain.LocationListItemResponse, len(locations))
	for i, item := range locations {
		responses[i] = domain.LocationListItemResponse{
			ID:           item.ID,
			LocationCode: item.LocationCode,
			Building:     item.Building,
			Floor:        item.Floor,
			Latitude:     item.Latitude,
			Longitude:    item.Longitude,
			Name:         item.Name,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		}
	}

	return responses, count, nil
}

func (s *Service) GetLocationsCursor(ctx context.Context, params query.Params, langCode string) ([]domain.LocationListItemResponse, error) {
	locations, err := s.Repo.GetLocationsCursor(ctx, params, langCode)
	if err != nil {
		return nil, err
	}

	// Convert LocationListItem to LocationListItemResponse
	responses := make([]domain.LocationListItemResponse, len(locations))
	for i, item := range locations {
		responses[i] = domain.LocationListItemResponse{
			ID:           item.ID,
			LocationCode: item.LocationCode,
			Building:     item.Building,
			Floor:        item.Floor,
			Latitude:     item.Latitude,
			Longitude:    item.Longitude,
			Name:         item.Name,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		}
	}

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

func (s *Service) GetLocationHierarchy(ctx context.Context, langCode string) ([]domain.LocationResponse, error) {
	hierarchy, err := s.Repo.GetLocationHierarchy(ctx, langCode)
	if err != nil {
		return nil, err
	}
	return hierarchy, nil
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

func (s *Service) CountLocations(ctx context.Context, params query.Params) (int64, error) {
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
