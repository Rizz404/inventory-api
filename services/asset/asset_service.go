package asset

import (
	"context"
	"mime/multipart"
	"strings"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/client/cloudinary"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/Rizz404/inventory-api/internal/utils"
	"github.com/oklog/ulid/v2"
)

// * Repository interface defines the contract for asset data operations
type Repository interface {
	// * MUTATION
	CreateAsset(ctx context.Context, payload *domain.Asset) (domain.Asset, error)
	UpdateAsset(ctx context.Context, payload *domain.Asset) (domain.Asset, error)
	UpdateAssetWithPayload(ctx context.Context, assetId string, payload *domain.UpdateAssetPayload) (domain.Asset, error)
	DeleteAsset(ctx context.Context, assetId string) error

	// * QUERY
	GetAssetsPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.Asset, error)
	GetAssetsCursor(ctx context.Context, params query.Params, langCode string) ([]domain.Asset, error)
	GetAssetById(ctx context.Context, assetId string) (domain.Asset, error)
	GetAssetByAssetTag(ctx context.Context, assetTag string) (domain.Asset, error)
	CheckAssetExists(ctx context.Context, assetId string) (bool, error)
	CheckAssetTagExists(ctx context.Context, assetTag string) (bool, error)
	CheckSerialNumberExists(ctx context.Context, serialNumber string) (bool, error)
	CheckAssetTagExistsExcluding(ctx context.Context, assetTag string, excludeAssetId string) (bool, error)
	CheckSerialNumberExistsExcluding(ctx context.Context, serialNumber string, excludeAssetId string) (bool, error)
	CountAssets(ctx context.Context, params query.Params) (int64, error)
	GetAssetStatistics(ctx context.Context) (domain.AssetStatistics, error)
}

// * AssetService interface defines the contract for asset business operations
type AssetService interface {
	// * MUTATION
	CreateAsset(ctx context.Context, payload *domain.CreateAssetPayload, dataMatrixImageFile *multipart.FileHeader, langCode string) (domain.AssetResponse, error)
	UpdateAsset(ctx context.Context, assetId string, payload *domain.UpdateAssetPayload, dataMatrixImageFile *multipart.FileHeader, langCode string) (domain.AssetResponse, error)
	DeleteAsset(ctx context.Context, assetId string) error

	// * QUERY
	GetAssetsPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.AssetResponse, int64, error)
	GetAssetsCursor(ctx context.Context, params query.Params, langCode string) ([]domain.AssetResponse, error)
	GetAssetById(ctx context.Context, assetId string, langCode string) (domain.AssetResponse, error)
	GetAssetByAssetTag(ctx context.Context, assetTag string, langCode string) (domain.AssetResponse, error)
	CheckAssetExists(ctx context.Context, assetId string) (bool, error)
	CheckAssetTagExists(ctx context.Context, assetTag string) (bool, error)
	CheckSerialNumberExists(ctx context.Context, serialNumber string) (bool, error)
	CountAssets(ctx context.Context, params query.Params) (int64, error)
	GetAssetStatistics(ctx context.Context) (domain.AssetStatisticsResponse, error)
}

type Service struct {
	Repo             Repository
	CloudinaryClient *cloudinary.Client
}

// * Ensure Service implements AssetService interface
var _ AssetService = (*Service)(nil)

func NewService(r Repository, cloudinaryClient *cloudinary.Client) AssetService {
	return &Service{
		Repo:             r,
		CloudinaryClient: cloudinaryClient,
	}
}

// *===========================MUTATION===========================*
func (s *Service) CreateAsset(ctx context.Context, payload *domain.CreateAssetPayload, dataMatrixImageFile *multipart.FileHeader, langCode string) (domain.AssetResponse, error) {
	// * Check if asset tag already exists
	if tagExists, err := s.Repo.CheckAssetTagExists(ctx, payload.AssetTag); err != nil {
		return domain.AssetResponse{}, err
	} else if tagExists {
		return domain.AssetResponse{}, domain.ErrConflictWithKey(utils.ErrAssetTagExistsKey)
	}

	// * Check if serial number already exists (if provided)
	if payload.SerialNumber != nil && *payload.SerialNumber != "" {
		if serialExists, err := s.Repo.CheckSerialNumberExists(ctx, *payload.SerialNumber); err != nil {
			return domain.AssetResponse{}, err
		} else if serialExists {
			return domain.AssetResponse{}, domain.ErrConflictWithKey(utils.ErrAssetSerialNumberExistsKey)
		}
	}

	status := domain.StatusActive
	if payload.Status != nil {
		status = *payload.Status
	}

	condition := domain.ConditionGood
	if payload.Condition != nil {
		condition = *payload.Condition
	}

	// * Handle data matrix image upload if file is provided
	var dataMatrixImageURL string = ""
	if dataMatrixImageFile != nil {
		// Upload file to Cloudinary if client is available
		if s.CloudinaryClient != nil {
			// Generate temporary asset ID for image naming
			tempAssetID := "temp_" + ulid.Make().String()
			uploadConfig := cloudinary.GetDataMatrixImageUploadConfig()
			publicID := "asset_" + tempAssetID + "_datamatrix"
			uploadConfig.PublicID = &publicID

			uploadResult, err := s.CloudinaryClient.UploadSingleFile(ctx, dataMatrixImageFile, uploadConfig)
			if err != nil {
				return domain.AssetResponse{}, domain.ErrBadRequestWithKey(utils.ErrFileUploadFailedKey)
			}
			dataMatrixImageURL = uploadResult.SecureURL
		} else {
			return domain.AssetResponse{}, domain.ErrBadRequestWithKey(utils.ErrCloudinaryConfigKey)
		}
	} else if payload.DataMatrixImageUrl != nil {
		// Use provided data matrix image URL from JSON/form data
		dataMatrixImageURL = *payload.DataMatrixImageUrl
	}

	// * Parse date strings to time.Time if provided
	var purchaseDate *time.Time
	if payload.PurchaseDate != nil && *payload.PurchaseDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", *payload.PurchaseDate); err == nil {
			purchaseDate = &parsedDate
		}
	}

	var warrantyEnd *time.Time
	if payload.WarrantyEnd != nil && *payload.WarrantyEnd != "" {
		if parsedDate, err := time.Parse("2006-01-02", *payload.WarrantyEnd); err == nil {
			warrantyEnd = &parsedDate
		}
	}

	// * Prepare new asset
	newAsset := domain.Asset{
		AssetTag:           payload.AssetTag,
		DataMatrixImageUrl: dataMatrixImageURL,
		AssetName:          payload.AssetName,
		CategoryID:         payload.CategoryID,
		Brand:              payload.Brand,
		Model:              payload.Model,
		SerialNumber:       payload.SerialNumber,
		PurchaseDate:       purchaseDate,
		PurchasePrice:      payload.PurchasePrice,
		VendorName:         payload.VendorName,
		WarrantyEnd:        warrantyEnd,
		Status:             status,
		Condition:          condition,
		LocationID:         payload.LocationID,
		AssignedTo:         payload.AssignedTo,
	}

	createdAsset, err := s.Repo.CreateAsset(ctx, &newAsset)
	if err != nil {
		// * Repository already handles error translation, so return directly
		return domain.AssetResponse{}, err
	}

	// * Update data matrix image public ID with actual asset ID if file was uploaded
	if dataMatrixImageFile != nil && s.CloudinaryClient != nil && dataMatrixImageURL != "" {
		// Re-upload with correct public ID
		uploadConfig := cloudinary.GetDataMatrixImageUploadConfig()
		finalPublicID := "asset_" + createdAsset.ID + "_datamatrix"
		uploadConfig.PublicID = &finalPublicID

		uploadResult, err := s.CloudinaryClient.UploadSingleFile(ctx, dataMatrixImageFile, uploadConfig)
		if err == nil {
			// Update asset with final data matrix image URL
			updatePayload := &domain.UpdateAssetPayload{
				DataMatrixImageUrl: &uploadResult.SecureURL,
			}
			createdAsset, _ = s.Repo.UpdateAssetWithPayload(ctx, createdAsset.ID, updatePayload)
		}
		// Note: We don't return error here to avoid failing asset creation if image re-upload fails
	}

	// * Convert to AssetResponse using mapper
	return mapper.AssetToResponse(&createdAsset), nil
}

func (s *Service) UpdateAsset(ctx context.Context, assetId string, payload *domain.UpdateAssetPayload, dataMatrixImageFile *multipart.FileHeader, langCode string) (domain.AssetResponse, error) {
	// Check if asset exists
	existingAsset, err := s.Repo.GetAssetById(ctx, assetId)
	if err != nil {
		return domain.AssetResponse{}, err
	}

	// * Check asset tag uniqueness if being updated
	if payload.AssetTag != nil {
		if tagExists, err := s.Repo.CheckAssetTagExistsExcluding(ctx, *payload.AssetTag, assetId); err != nil {
			return domain.AssetResponse{}, err
		} else if tagExists {
			return domain.AssetResponse{}, domain.ErrConflictWithKey(utils.ErrAssetTagExistsKey)
		}
	}

	// * Check serial number uniqueness if being updated
	if payload.SerialNumber != nil && *payload.SerialNumber != "" {
		if serialExists, err := s.Repo.CheckSerialNumberExistsExcluding(ctx, *payload.SerialNumber, assetId); err != nil {
			return domain.AssetResponse{}, err
		} else if serialExists {
			return domain.AssetResponse{}, domain.ErrConflictWithKey(utils.ErrAssetSerialNumberExistsKey)
		}
	}

	// * Handle data matrix image update
	var shouldDeleteOldImage bool
	oldImagePublicID := "asset_" + assetId + "_datamatrix"

	if dataMatrixImageFile != nil {
		// Upload new data matrix image file
		if s.CloudinaryClient != nil {
			uploadConfig := cloudinary.GetDataMatrixImageUploadConfig()
			publicID := "asset_" + assetId + "_datamatrix"
			uploadConfig.PublicID = &publicID

			uploadResult, err := s.CloudinaryClient.UploadSingleFile(ctx, dataMatrixImageFile, uploadConfig)
			if err != nil {
				return domain.AssetResponse{}, domain.ErrBadRequestWithKey(utils.ErrFileUploadFailedKey)
			}

			// Set new data matrix image URL in payload
			payload.DataMatrixImageUrl = &uploadResult.SecureURL
			// Note: Cloudinary will automatically overwrite old image due to same public ID
		} else {
			return domain.AssetResponse{}, domain.ErrBadRequestWithKey(utils.ErrCloudinaryConfigKey)
		}
	} else if payload.DataMatrixImageUrl != nil {
		// Handle data matrix image URL changes from JSON/form data
		if *payload.DataMatrixImageUrl == "" || *payload.DataMatrixImageUrl == "null" {
			// User wants to remove data matrix image
			payload.DataMatrixImageUrl = nil
			shouldDeleteOldImage = true
		}
		// If payload.DataMatrixImageUrl has a valid URL, it will be used as-is
	}

	// Use the UpdateAssetWithPayload method
	_, err = s.Repo.UpdateAssetWithPayload(ctx, assetId, payload)
	if err != nil {
		return domain.AssetResponse{}, err
	}

	// * Delete old data matrix image from Cloudinary if needed
	if shouldDeleteOldImage && s.CloudinaryClient != nil && existingAsset.DataMatrixImageUrl != "" {
		// Only delete if the old image was stored in Cloudinary (contains our public ID pattern)
		if strings.Contains(existingAsset.DataMatrixImageUrl, "asset_"+assetId+"_datamatrix") {
			_ = s.CloudinaryClient.DeleteFile(ctx, oldImagePublicID)
			// Note: We don't return error here to avoid failing asset update if image deletion fails
		}
	}

	// * Update asset and convert to AssetResponse using mapper
	updatedAsset, err := s.Repo.UpdateAssetWithPayload(ctx, assetId, payload)
	if err != nil {
		return domain.AssetResponse{}, err
	}

	return mapper.AssetToResponse(&updatedAsset), nil
}

func (s *Service) DeleteAsset(ctx context.Context, assetId string) error {
	err := s.Repo.DeleteAsset(ctx, assetId)
	if err != nil {
		return err
	}
	return nil
}

// *===========================QUERY===========================*
func (s *Service) GetAssetsPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.AssetResponse, int64, error) {
	assets, err := s.Repo.GetAssetsPaginated(ctx, params, langCode)
	if err != nil {
		return nil, 0, err
	}

	// * Only count total if pagination is offset-based
	count, err := s.Repo.CountAssets(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	return mapper.AssetsToResponses(assets), count, nil
}

func (s *Service) GetAssetsCursor(ctx context.Context, params query.Params, langCode string) ([]domain.AssetResponse, error) {
	assets, err := s.Repo.GetAssetsCursor(ctx, params, langCode)
	if err != nil {
		return nil, err
	}

	return mapper.AssetsToResponses(assets), nil
}

func (s *Service) GetAssetById(ctx context.Context, assetId string, langCode string) (domain.AssetResponse, error) {
	asset, err := s.Repo.GetAssetById(ctx, assetId)
	if err != nil {
		return domain.AssetResponse{}, err
	}

	return mapper.AssetToResponse(&asset), nil
}

func (s *Service) GetAssetByAssetTag(ctx context.Context, assetTag string, langCode string) (domain.AssetResponse, error) {
	asset, err := s.Repo.GetAssetByAssetTag(ctx, assetTag)
	if err != nil {
		return domain.AssetResponse{}, err
	}

	return mapper.AssetToResponse(&asset), nil
}

func (s *Service) CheckAssetExists(ctx context.Context, assetId string) (bool, error) {
	exists, err := s.Repo.CheckAssetExists(ctx, assetId)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Service) CheckAssetTagExists(ctx context.Context, assetTag string) (bool, error) {
	exists, err := s.Repo.CheckAssetTagExists(ctx, assetTag)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Service) CheckSerialNumberExists(ctx context.Context, serialNumber string) (bool, error) {
	exists, err := s.Repo.CheckSerialNumberExists(ctx, serialNumber)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Service) CountAssets(ctx context.Context, params query.Params) (int64, error) {
	count, err := s.Repo.CountAssets(ctx, params)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Service) GetAssetStatistics(ctx context.Context) (domain.AssetStatisticsResponse, error) {
	stats, err := s.Repo.GetAssetStatistics(ctx)
	if err != nil {
		return domain.AssetStatisticsResponse{}, err
	}
	return mapper.AssetStatisticsToResponse(&stats), nil
}
