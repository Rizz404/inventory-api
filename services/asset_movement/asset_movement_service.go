package asset_movement

import (
	"context"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/Rizz404/inventory-api/internal/utils"
)

// * Repository interface defines the contract for asset movement data operations
type Repository interface {
	// * MUTATION
	CreateAssetMovement(ctx context.Context, payload *domain.AssetMovement) (domain.AssetMovement, error)
	UpdateAssetMovement(ctx context.Context, movementId string, payload *domain.UpdateAssetMovementPayload) (domain.AssetMovement, error)
	DeleteAssetMovement(ctx context.Context, movementId string) error

	// * QUERY
	GetAssetMovementsPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.AssetMovementListItem, error)
	GetAssetMovementsCursor(ctx context.Context, params query.Params, langCode string) ([]domain.AssetMovementListItem, error)
	GetAssetMovementById(ctx context.Context, movementId string) (domain.AssetMovement, error)
	GetAssetMovementsByAssetId(ctx context.Context, assetId string, params query.Params) ([]domain.AssetMovement, error)
	CheckAssetMovementExist(ctx context.Context, movementId string) (bool, error)
	CountAssetMovements(ctx context.Context, params query.Params) (int64, error)
	GetAssetMovementStatistics(ctx context.Context) (domain.AssetMovementStatistics, error)
}

// * AssetService interface for checking asset existence
type AssetService interface {
	CheckAssetExists(ctx context.Context, assetId string) (bool, error)
	GetAssetById(ctx context.Context, assetId string, langCode string) (domain.AssetResponse, error)
}

// * LocationService interface for checking location existence
type LocationService interface {
	CheckLocationExists(ctx context.Context, locationId string) (bool, error)
}

// * UserService interface for checking user existence
type UserService interface {
	CheckUserExists(ctx context.Context, userId string) (bool, error)
}

// * AssetMovementService interface defines the contract for asset movement business operations
type AssetMovementService interface {
	// * MUTATION
	CreateAssetMovement(ctx context.Context, payload *domain.CreateAssetMovementPayload, movedBy string) (domain.AssetMovementResponse, error)
	UpdateAssetMovement(ctx context.Context, movementId string, payload *domain.UpdateAssetMovementPayload) (domain.AssetMovementResponse, error)
	DeleteAssetMovement(ctx context.Context, movementId string) error

	// * QUERY
	GetAssetMovementsPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.AssetMovementListItemResponse, int64, error)
	GetAssetMovementsCursor(ctx context.Context, params query.Params, langCode string) ([]domain.AssetMovementListItemResponse, error)
	GetAssetMovementById(ctx context.Context, movementId string, langCode string) (domain.AssetMovementResponse, error)
	GetAssetMovementsByAssetId(ctx context.Context, assetId string, params query.Params, langCode string) ([]domain.AssetMovementResponse, error)
	CheckAssetMovementExists(ctx context.Context, movementId string) (bool, error)
	CountAssetMovements(ctx context.Context, params query.Params) (int64, error)
	GetAssetMovementStatistics(ctx context.Context) (domain.AssetMovementStatisticsResponse, error)
}

type Service struct {
	Repo            Repository
	AssetService    AssetService
	LocationService LocationService
	UserService     UserService
}

// * Ensure Service implements AssetMovementService interface
var _ AssetMovementService = (*Service)(nil)

func NewService(r Repository, assetService AssetService, locationService LocationService, userService UserService) AssetMovementService {
	return &Service{
		Repo:            r,
		AssetService:    assetService,
		LocationService: locationService,
		UserService:     userService,
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

	// * Validate destination (must have either ToLocationID or ToUserID, but not both)
	if payload.ToLocationID == nil && payload.ToUserID == nil {
		return domain.AssetMovementResponse{}, domain.ErrBadRequestWithKey(utils.ErrAssetMovementNoChangeKey)
	}

	if payload.ToLocationID != nil && payload.ToUserID != nil {
		return domain.AssetMovementResponse{}, domain.ErrBadRequestWithKey(utils.ErrAssetMovementInvalidLocationKey)
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

	// * Convert to AssetMovementResponse using mapper
	return mapper.AssetMovementToResponse(&createdMovement, mapper.DefaultLangCode), nil
}

func (s *Service) UpdateAssetMovement(ctx context.Context, movementId string, payload *domain.UpdateAssetMovementPayload) (domain.AssetMovementResponse, error) {
	// * Check if asset movement exists
	_, err := s.Repo.GetAssetMovementById(ctx, movementId)
	if err != nil {
		return domain.AssetMovementResponse{}, err
	}

	// * Validate destination if being updated
	if payload.ToLocationID != nil && payload.ToUserID != nil {
		return domain.AssetMovementResponse{}, domain.ErrBadRequestWithKey(utils.ErrAssetMovementInvalidLocationKey)
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

// *===========================QUERY===========================*
func (s *Service) GetAssetMovementsPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.AssetMovementListItemResponse, int64, error) {
	listItems, err := s.Repo.GetAssetMovementsPaginated(ctx, params, langCode)
	if err != nil {
		return nil, 0, err
	}

	// * Count total for pagination
	count, err := s.Repo.CountAssetMovements(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// Convert AssetMovementListItem to AssetMovementListItemResponse
	responses := make([]domain.AssetMovementListItemResponse, len(listItems))
	for i, item := range listItems {
		responses[i] = domain.AssetMovementListItemResponse{
			ID:             item.ID,
			AssetID:        item.AssetID,
			FromLocationID: item.FromLocationID,
			ToLocationID:   item.ToLocationID,
			FromUserID:     item.FromUserID,
			ToUserID:       item.ToUserID,
			MovedByID:      item.MovedByID,
			MovementDate:   item.MovementDate,
			Notes:          item.Notes,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		}
	}

	return responses, count, nil
}

func (s *Service) GetAssetMovementsCursor(ctx context.Context, params query.Params, langCode string) ([]domain.AssetMovementListItemResponse, error) {
	listItems, err := s.Repo.GetAssetMovementsCursor(ctx, params, langCode)
	if err != nil {
		return nil, err
	}

	// Convert AssetMovementListItem to AssetMovementListItemResponse
	responses := make([]domain.AssetMovementListItemResponse, len(listItems))
	for i, item := range listItems {
		responses[i] = domain.AssetMovementListItemResponse{
			ID:             item.ID,
			AssetID:        item.AssetID,
			FromLocationID: item.FromLocationID,
			ToLocationID:   item.ToLocationID,
			FromUserID:     item.FromUserID,
			ToUserID:       item.ToUserID,
			MovedByID:      item.MovedByID,
			MovementDate:   item.MovementDate,
			Notes:          item.Notes,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		}
	}

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

func (s *Service) GetAssetMovementsByAssetId(ctx context.Context, assetId string, params query.Params, langCode string) ([]domain.AssetMovementResponse, error) {
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

	// * Convert to AssetMovementResponse using mapper
	movementResponses := mapper.AssetMovementsToResponses(movements, langCode)

	return movementResponses, nil
}

func (s *Service) CheckAssetMovementExists(ctx context.Context, movementId string) (bool, error) {
	exists, err := s.Repo.CheckAssetMovementExist(ctx, movementId)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Service) CountAssetMovements(ctx context.Context, params query.Params) (int64, error) {
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
