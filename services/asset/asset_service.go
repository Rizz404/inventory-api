package asset

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/client/cloudinary"
	"github.com/Rizz404/inventory-api/internal/notification/messages"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/Rizz404/inventory-api/internal/utils"
	"github.com/oklog/ulid/v2"
)

// * Repository interface defines the contract for asset data operations
type Repository interface {
	// * MUTATION
	CreateAsset(ctx context.Context, payload *domain.Asset) (domain.Asset, error)
	BulkCreateAssets(ctx context.Context, assets []domain.Asset) ([]domain.Asset, error)
	UpdateAsset(ctx context.Context, assetId string, payload *domain.UpdateAssetPayload) (domain.Asset, error)
	DeleteAsset(ctx context.Context, assetId string) error
	BulkDeleteAssets(ctx context.Context, assetIds []string) (domain.BulkDeleteAssets, error)

	// * QUERY
	GetAssetsPaginated(ctx context.Context, params domain.AssetParams, langCode string) ([]domain.Asset, error)
	GetAssetsCursor(ctx context.Context, params domain.AssetParams, langCode string) ([]domain.Asset, error)
	GetAssetById(ctx context.Context, assetId string) (domain.Asset, error)
	GetAssetByAssetTag(ctx context.Context, assetTag string) (domain.Asset, error)
	CheckAssetExists(ctx context.Context, assetId string) (bool, error)
	CheckAssetTagExists(ctx context.Context, assetTag string) (bool, error)
	CheckSerialNumberExists(ctx context.Context, serialNumber string) (bool, error)
	CheckAssetTagExistsExcluding(ctx context.Context, assetTag string, excludeAssetId string) (bool, error)
	CheckSerialNumberExistsExcluding(ctx context.Context, serialNumber string, excludeAssetId string) (bool, error)
	CountAssets(ctx context.Context, params domain.AssetParams) (int64, error)
	GetAssetStatistics(ctx context.Context) (domain.AssetStatistics, error)
	GetLastAssetTagByCategory(ctx context.Context, categoryId string) (string, error)
	GetLastAssetTagsByCategoryBatch(ctx context.Context, categoryId string, quantity int) ([]string, error)
	GetAssetsForExport(ctx context.Context, params domain.AssetParams, langCode string) ([]domain.Asset, error)
	GetAssetsWithWarrantyExpiring(ctx context.Context, daysFromNow int) ([]domain.Asset, error)
	GetAssetsWithExpiredWarranty(ctx context.Context) ([]domain.Asset, error)

	// * IMAGE & ASSET IMAGES CRUD
	CreateImage(ctx context.Context, imageURL string, publicID *string) (domain.Image, error)
	GetImageByPublicID(ctx context.Context, publicID string) (*domain.Image, error)
	GetAvailableAssetImages(ctx context.Context, limit int, cursor string) ([]domain.Image, error)
	AttachImagesToAsset(ctx context.Context, assetID string, imageIDs []string, displayOrders []int, primaryIndex int) ([]domain.AssetImage, error)
	GetAssetImages(ctx context.Context, assetID string) ([]domain.AssetImage, error)
	DetachImageFromAsset(ctx context.Context, assetImageID string) error
	DetachAllImagesFromAsset(ctx context.Context, assetID string) error
	UpdateAssetImagePrimary(ctx context.Context, assetID string, assetImageID string) error
	DeleteUnusedImages(ctx context.Context) error
}

// * AssetService interface defines the contract for asset business operations
type AssetService interface {
	// * MUTATION
	CreateAsset(ctx context.Context, payload *domain.CreateAssetPayload, dataMatrixImageFile *multipart.FileHeader, langCode string) (domain.AssetResponse, error)
	BulkCreateAssets(ctx context.Context, payload *domain.BulkCreateAssetsPayload, langCode string) (domain.BulkCreateAssetsResponse, error)
	UpdateAsset(ctx context.Context, assetId string, payload *domain.UpdateAssetPayload, dataMatrixImageFile *multipart.FileHeader, langCode string) (domain.AssetResponse, error)
	DeleteAsset(ctx context.Context, assetId string) error
	BulkDeleteAssets(ctx context.Context, payload *domain.BulkDeleteAssetsPayload) (domain.BulkDeleteAssetsResponse, error)

	// * QUERY
	GetAssetsPaginated(ctx context.Context, params domain.AssetParams, langCode string) ([]domain.AssetResponse, int64, error)
	GetAssetsCursor(ctx context.Context, params domain.AssetParams, langCode string) ([]domain.AssetResponse, error)
	GetAssetById(ctx context.Context, assetId string, langCode string) (domain.AssetResponse, error)
	GetAssetByAssetTag(ctx context.Context, assetTag string, langCode string) (domain.AssetResponse, error)
	CheckAssetExists(ctx context.Context, assetId string) (bool, error)
	CheckAssetTagExists(ctx context.Context, assetTag string) (bool, error)
	CheckSerialNumberExists(ctx context.Context, serialNumber string) (bool, error)
	CountAssets(ctx context.Context, params domain.AssetParams) (int64, error)
	GetAssetStatistics(ctx context.Context) (domain.AssetStatisticsResponse, error)
	GenerateAssetTagSuggestion(ctx context.Context, payload *domain.GenerateAssetTagPayload) (domain.GenerateAssetTagResponse, error)
	GenerateBulkAssetTags(ctx context.Context, payload *domain.GenerateBulkAssetTagsPayload) (domain.GenerateBulkAssetTagsResponse, error)
	UploadBulkDataMatrixImages(ctx context.Context, assetTags []string, files []*multipart.FileHeader) (domain.UploadBulkDataMatrixResponse, error)
	DeleteBulkDataMatrixImages(ctx context.Context, payload *domain.DeleteBulkDataMatrixPayload) (domain.DeleteBulkDataMatrixResponse, error)

	// * TEMPLATE IMAGES (For Bulk Create)
	UploadTemplateImages(ctx context.Context, files []*multipart.FileHeader) (domain.UploadTemplateImagesResponse, error)

	// * ASSET IMAGES
	GetAvailableAssetImages(ctx context.Context, limit int, cursor string) ([]domain.ImageResponse, error)
	UploadBulkAssetImages(ctx context.Context, assetIds []string, files []*multipart.FileHeader) (domain.UploadBulkAssetImagesResponse, error)
	DeleteBulkAssetImages(ctx context.Context, payload *domain.DeleteBulkAssetImagesPayload) (domain.DeleteBulkAssetImagesResponse, error)

	// * EXPORT
	ExportAssetList(ctx context.Context, payload *domain.ExportAssetListPayload, langCode string) ([]byte, string, error)
	ExportAssetStatistics(ctx context.Context, langCode string) ([]byte, string, error)
	ExportAssetDataMatrix(ctx context.Context, payload *domain.ExportAssetDataMatrixPayload, langCode string) ([]byte, string, error)
}

// * NotificationService interface for creating notifications
type NotificationService interface {
	CreateNotification(ctx context.Context, payload *domain.CreateNotificationPayload) (domain.NotificationResponse, error)
}

// * CategoryService interface for getting category information
type CategoryService interface {
	GetCategoryById(ctx context.Context, categoryId string, langCode string) (domain.CategoryResponse, error)
}

// * UserRepository interface for getting user details
type UserRepository interface {
	GetUsersPaginated(ctx context.Context, params domain.UserParams) ([]domain.User, error)
}

type Service struct {
	Repo                Repository
	CloudinaryClient    *cloudinary.Client
	NotificationService NotificationService
	CategoryService     CategoryService
	UserRepo            UserRepository
}

// * Ensure Service implements AssetService interface
var _ AssetService = (*Service)(nil)

func NewService(r Repository, cloudinaryClient *cloudinary.Client, notificationService NotificationService, categoryService CategoryService, userRepo UserRepository) AssetService {
	return &Service{
		Repo:                r,
		CloudinaryClient:    cloudinaryClient,
		NotificationService: notificationService,
		CategoryService:     categoryService,
		UserRepo:            userRepo,
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
	if payload.Status != "" {
		status = payload.Status
	}

	condition := domain.ConditionGood
	if payload.Condition != "" {
		condition = payload.Condition
	}

	// * Handle data matrix image upload if file is provided
	var dataMatrixImageURL string = ""
	if dataMatrixImageFile != nil {
		// Upload file to Cloudinary if client is available
		if s.CloudinaryClient != nil {
			// Generate ULID for unique filename
			ulidStr := ulid.Make().String()
			uploadConfig := cloudinary.GetDataMatrixImageUploadConfig()
			// Naming pattern: {assetTag}_{ulid}
			publicID := fmt.Sprintf("%s_%s", payload.AssetTag, ulidStr)
			uploadConfig.PublicID = &publicID

			uploadResult, err := s.CloudinaryClient.UploadSingleFile(ctx, dataMatrixImageFile, uploadConfig)
			if err != nil {
				// Provide detailed error message
				errorMsg := "Failed to upload data matrix image: " + err.Error()
				return domain.AssetResponse{}, domain.ErrBadRequest(errorMsg)
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
		if parsedDate, err := time.ParseInLocation("2006-01-02", *payload.PurchaseDate, time.UTC); err == nil {
			purchaseDate = &parsedDate
		}
	}

	var warrantyEnd *time.Time
	if payload.WarrantyEnd != nil && *payload.WarrantyEnd != "" {
		if parsedDate, err := time.ParseInLocation("2006-01-02", *payload.WarrantyEnd, time.UTC); err == nil {
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

	// * Attach images if imageUrls provided
	if len(payload.ImageUrls) > 0 {
		err = s.attachImageUrlsToAsset(ctx, createdAsset.ID, payload.ImageUrls)
		if err != nil {
			log.Printf("Failed to attach images to asset %s: %v", createdAsset.AssetTag, err)
		}
	}

	// * Send notification if asset is assigned to a user
	if payload.AssignedTo != nil && *payload.AssignedTo != "" {
		go s.sendAssetAssignmentNotification(context.Background(), &createdAsset, *payload.AssignedTo, true)
	}

	// * Send notification if asset is high-value
	if payload.PurchasePrice != nil && *payload.PurchasePrice > 10000 {
		go s.sendHighValueAssetNotificationToAdmins(context.Background(), &createdAsset)
	}

	// * Convert to AssetResponse using mapper
	return mapper.AssetToResponse(&createdAsset, langCode), nil
}

func (s *Service) BulkCreateAssets(ctx context.Context, payload *domain.BulkCreateAssetsPayload, langCode string) (domain.BulkCreateAssetsResponse, error) {
	if payload == nil || len(payload.Assets) == 0 {
		return domain.BulkCreateAssetsResponse{}, domain.ErrBadRequest("assets payload is required")
	}

	assetTagSeen := make(map[string]struct{})
	serialSeen := make(map[string]struct{})
	for _, assetPayload := range payload.Assets {
		if _, exists := assetTagSeen[assetPayload.AssetTag]; exists {
			return domain.BulkCreateAssetsResponse{}, domain.ErrConflictWithKey(utils.ErrAssetTagExistsKey)
		}
		assetTagSeen[assetPayload.AssetTag] = struct{}{}

		if assetPayload.SerialNumber != nil && *assetPayload.SerialNumber != "" {
			if _, exists := serialSeen[*assetPayload.SerialNumber]; exists {
				return domain.BulkCreateAssetsResponse{}, domain.ErrConflictWithKey(utils.ErrAssetSerialNumberExistsKey)
			}
			serialSeen[*assetPayload.SerialNumber] = struct{}{}
		}
	}

	for tag := range assetTagSeen {
		exists, err := s.Repo.CheckAssetTagExists(ctx, tag)
		if err != nil {
			return domain.BulkCreateAssetsResponse{}, err
		}
		if exists {
			return domain.BulkCreateAssetsResponse{}, domain.ErrConflictWithKey(utils.ErrAssetTagExistsKey)
		}
	}

	for serial := range serialSeen {
		exists, err := s.Repo.CheckSerialNumberExists(ctx, serial)
		if err != nil {
			return domain.BulkCreateAssetsResponse{}, err
		}
		if exists {
			return domain.BulkCreateAssetsResponse{}, domain.ErrConflictWithKey(utils.ErrAssetSerialNumberExistsKey)
		}
	}

	assets := make([]domain.Asset, len(payload.Assets))
	for i, assetPayload := range payload.Assets {
		status := domain.StatusActive
		if assetPayload.Status != "" {
			status = assetPayload.Status
		}

		condition := domain.ConditionGood
		if assetPayload.Condition != "" {
			condition = assetPayload.Condition
		}

		var purchaseDate *time.Time
		if assetPayload.PurchaseDate != nil && *assetPayload.PurchaseDate != "" {
			if parsedDate, err := time.ParseInLocation("2006-01-02", *assetPayload.PurchaseDate, time.UTC); err == nil {
				purchaseDate = &parsedDate
			}
		}

		var warrantyEnd *time.Time
		if assetPayload.WarrantyEnd != nil && *assetPayload.WarrantyEnd != "" {
			if parsedDate, err := time.ParseInLocation("2006-01-02", *assetPayload.WarrantyEnd, time.UTC); err == nil {
				warrantyEnd = &parsedDate
			}
		}

		dataMatrixImageURL := ""
		if assetPayload.DataMatrixImageUrl != nil {
			dataMatrixImageURL = *assetPayload.DataMatrixImageUrl
		}

		assets[i] = domain.Asset{
			AssetTag:           assetPayload.AssetTag,
			DataMatrixImageUrl: dataMatrixImageURL,
			AssetName:          assetPayload.AssetName,
			CategoryID:         assetPayload.CategoryID,
			Brand:              assetPayload.Brand,
			Model:              assetPayload.Model,
			SerialNumber:       assetPayload.SerialNumber,
			PurchaseDate:       purchaseDate,
			PurchasePrice:      assetPayload.PurchasePrice,
			VendorName:         assetPayload.VendorName,
			WarrantyEnd:        warrantyEnd,
			Status:             status,
			Condition:          condition,
			LocationID:         assetPayload.LocationID,
			AssignedTo:         assetPayload.AssignedTo,
		}
	}

	createdAssets, err := s.Repo.BulkCreateAssets(ctx, assets)
	if err != nil {
		return domain.BulkCreateAssetsResponse{}, err
	}

	// Attach images to assets if imageUrls provided
	for i := range createdAssets {
		if len(payload.Assets[i].ImageUrls) > 0 {
			if err := s.attachImageUrlsToAsset(ctx, createdAssets[i].ID, payload.Assets[i].ImageUrls); err != nil {
				log.Printf("Failed to attach images to asset %s: %v", createdAssets[i].AssetTag, err)
				// Continue with other assets, don't fail entire bulk operation
			}
		}

		if payload.Assets[i].AssignedTo != nil && *payload.Assets[i].AssignedTo != "" {
			go s.sendAssetAssignmentNotification(context.Background(), &createdAssets[i], *payload.Assets[i].AssignedTo, true)
		}

		if payload.Assets[i].PurchasePrice != nil && *payload.Assets[i].PurchasePrice > 10000 {
			go s.sendHighValueAssetNotificationToAdmins(context.Background(), &createdAssets[i])
		}
	}

	response := domain.BulkCreateAssetsResponse{
		Assets: mapper.AssetsToResponses(createdAssets, langCode),
	}

	return response, nil
}

func (s *Service) UpdateAsset(ctx context.Context, assetId string, payload *domain.UpdateAssetPayload, dataMatrixImageFile *multipart.FileHeader, langCode string) (domain.AssetResponse, error) {
	// Check if asset exists
	existingAsset, err := s.Repo.GetAssetById(ctx, assetId)
	if err != nil {
		return domain.AssetResponse{}, err
	}

	// * Auto-regenerate asset tag if category is being changed
	var categoryChanged bool
	if payload.CategoryID != nil && *payload.CategoryID != existingAsset.CategoryID {
		categoryChanged = true

		// Get new category info
		newCategory, err := s.CategoryService.GetCategoryById(ctx, *payload.CategoryID, langCode)
		if err != nil {
			return domain.AssetResponse{}, err
		}

		// Get last asset tag for new category
		lastTag, err := s.Repo.GetLastAssetTagByCategory(ctx, *payload.CategoryID)
		if err != nil {
			return domain.AssetResponse{}, err
		}

		// Calculate next increment
		nextIncrement := 1
		if lastTag != "" {
			// Extract number from tag format: PREFIX-00001
			parts := strings.Split(lastTag, "-")
			if len(parts) > 0 {
				numberStr := parts[len(parts)-1]
				if lastNum, err := strconv.Atoi(numberStr); err == nil {
					nextIncrement = lastNum + 1
				}
			}
		}

		// Generate new asset tag with dash and 5-digit padding
		newAssetTag := fmt.Sprintf("%s-%05d", newCategory.CategoryCode, nextIncrement)
		payload.AssetTag = &newAssetTag

		log.Printf("Category changed for asset %s: %s -> %s, regenerated tag: %s -> %s",
			assetId, existingAsset.CategoryID, *payload.CategoryID, existingAsset.AssetTag, newAssetTag)
	}

	// * Check asset tag uniqueness if being updated (and not auto-generated from category change)
	if payload.AssetTag != nil && !categoryChanged {
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
	// Extract old public ID from URL if exists
	oldImagePublicID := ""
	if existingAsset.DataMatrixImageUrl != "" && strings.Contains(existingAsset.DataMatrixImageUrl, "sigma-asset/datamatrix/") {
		// Extract public ID from Cloudinary URL
		// Format: https://res.cloudinary.com/.../sigma-asset/datamatrix/{assetTag}_{ulid}.ext
		parts := strings.Split(existingAsset.DataMatrixImageUrl, "/")
		for i, part := range parts {
			if part == "datamatrix" && i+1 < len(parts) {
				// Get filename with extension, remove extension
				fileWithExt := parts[i+1]
				lastDot := strings.LastIndex(fileWithExt, ".")
				if lastDot > 0 {
					oldImagePublicID = "sigma-asset/datamatrix/" + fileWithExt[:lastDot]
				}
				break
			}
		}
	}

	// If category changed, clear old data matrix image (tag changed, QR code invalid)
	if categoryChanged && existingAsset.DataMatrixImageUrl != "" {
		emptyString := ""
		payload.DataMatrixImageUrl = &emptyString
		shouldDeleteOldImage = true
		log.Printf("Category changed for asset %s, clearing old data matrix image", assetId)
	}

	if dataMatrixImageFile != nil {
		// Upload new data matrix image file
		if s.CloudinaryClient != nil {
			// Generate ULID for unique filename
			ulidStr := ulid.Make().String()
			uploadConfig := cloudinary.GetDataMatrixImageUploadConfig()
			// Use asset tag from existing asset or payload
			assetTag := existingAsset.AssetTag
			if payload.AssetTag != nil {
				assetTag = *payload.AssetTag
			}
			// Naming pattern: {assetTag}_{ulid}
			publicID := fmt.Sprintf("%s_%s", assetTag, ulidStr)
			uploadConfig.PublicID = &publicID

			uploadResult, err := s.CloudinaryClient.UploadSingleFile(ctx, dataMatrixImageFile, uploadConfig)
			if err != nil {
				return domain.AssetResponse{}, domain.ErrBadRequest("Failed to upload data matrix image: " + err.Error())
			}

			// Set new data matrix image URL in payload
			payload.DataMatrixImageUrl = &uploadResult.SecureURL
			shouldDeleteOldImage = true
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

	// * Update asset and convert to AssetResponse using mapper
	updatedAsset, err := s.Repo.UpdateAsset(ctx, assetId, payload)
	if err != nil {
		return domain.AssetResponse{}, err
	}

	// * Delete old data matrix image from Cloudinary if needed
	if shouldDeleteOldImage && s.CloudinaryClient != nil && oldImagePublicID != "" {
		err = s.CloudinaryClient.DeleteFile(ctx, oldImagePublicID)
		if err != nil {
			log.Printf("Failed to delete old data matrix image: %v", err)
		}
		// Note: We don't return error here to avoid failing asset update if image deletion fails
	}

	// * Send notifications for changes
	go s.sendUpdateNotifications(context.Background(), &existingAsset, &updatedAsset, payload)

	return mapper.AssetToResponse(&updatedAsset, langCode), nil
}

func (s *Service) DeleteAsset(ctx context.Context, assetId string) error {
	err := s.Repo.DeleteAsset(ctx, assetId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) BulkDeleteAssets(ctx context.Context, payload *domain.BulkDeleteAssetsPayload) (domain.BulkDeleteAssetsResponse, error) {
	// * Validate that IDs are provided
	if len(payload.IDS) == 0 {
		return domain.BulkDeleteAssetsResponse{}, domain.ErrBadRequestWithKey(utils.ErrAssetIDRequiredKey)
	}

	// * Perform bulk delete operation
	result, err := s.Repo.BulkDeleteAssets(ctx, payload.IDS)
	if err != nil {
		return domain.BulkDeleteAssetsResponse{}, err
	}

	// * Convert to response
	response := domain.BulkDeleteAssetsResponse{
		RequestedIDS: result.RequestedIDS,
		DeletedIDS:   result.DeletedIDS,
	}

	return response, nil
}

// *===========================QUERY===========================*
func (s *Service) GetAssetsPaginated(ctx context.Context, params domain.AssetParams, langCode string) ([]domain.AssetResponse, int64, error) {
	assets, err := s.Repo.GetAssetsPaginated(ctx, params, langCode)
	if err != nil {
		return nil, 0, err
	}

	// * Only count total if pagination is offset-based
	count, err := s.Repo.CountAssets(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	return mapper.AssetsToResponses(assets, langCode), count, nil
}

func (s *Service) GetAssetsCursor(ctx context.Context, params domain.AssetParams, langCode string) ([]domain.AssetResponse, error) {
	assets, err := s.Repo.GetAssetsCursor(ctx, params, langCode)
	if err != nil {
		return nil, err
	}

	return mapper.AssetsToResponses(assets, langCode), nil
}

func (s *Service) GetAssetById(ctx context.Context, assetId string, langCode string) (domain.AssetResponse, error) {
	asset, err := s.Repo.GetAssetById(ctx, assetId)
	if err != nil {
		return domain.AssetResponse{}, err
	}

	return mapper.AssetToResponse(&asset, langCode), nil
}

func (s *Service) GetAssetByAssetTag(ctx context.Context, assetTag string, langCode string) (domain.AssetResponse, error) {
	asset, err := s.Repo.GetAssetByAssetTag(ctx, assetTag)
	if err != nil {
		return domain.AssetResponse{}, err
	}

	return mapper.AssetToResponse(&asset, langCode), nil
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

func (s *Service) CountAssets(ctx context.Context, params domain.AssetParams) (int64, error) {
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

// GenerateAssetTagSuggestion generates a suggested asset tag based on category code
func (s *Service) GenerateAssetTagSuggestion(ctx context.Context, payload *domain.GenerateAssetTagPayload) (domain.GenerateAssetTagResponse, error) {
	// * Get category to retrieve CategoryCode
	category, err := s.CategoryService.GetCategoryById(ctx, payload.CategoryID, "en")
	if err != nil {
		return domain.GenerateAssetTagResponse{}, err
	}

	// * Get the last asset tag for this category
	lastAssetTag, err := s.Repo.GetLastAssetTagByCategory(ctx, payload.CategoryID)
	if err != nil {
		return domain.GenerateAssetTagResponse{}, err
	}

	// * Calculate next increment
	nextIncrement := 1
	if lastAssetTag != "" {
		// Extract the numeric part from the last asset tag
		// Format expected: CATEGORYCODE-00001 (use LastIndex to handle multi-dash codes like FURN-SEAT)
		dashIndex := strings.LastIndex(lastAssetTag, "-")
		if dashIndex != -1 && dashIndex < len(lastAssetTag)-1 {
			numericPart := lastAssetTag[dashIndex+1:]
			// Try to parse the numeric part
			var parsedNum int
			_, err := fmt.Sscanf(numericPart, "%d", &parsedNum)
			if err == nil {
				nextIncrement = parsedNum + 1
			}
		}
	}

	// * Generate suggested tag with dash and 5-digit padding
	suggestedTag := fmt.Sprintf("%s-%05d", category.CategoryCode, nextIncrement)

	return domain.GenerateAssetTagResponse{
		CategoryCode:  category.CategoryCode,
		LastAssetTag:  lastAssetTag,
		SuggestedTag:  suggestedTag,
		NextIncrement: nextIncrement,
	}, nil
}

// GenerateBulkAssetTags generates multiple sequential asset tags for bulk operations
func (s *Service) GenerateBulkAssetTags(ctx context.Context, payload *domain.GenerateBulkAssetTagsPayload) (domain.GenerateBulkAssetTagsResponse, error) {
	// * Validate quantity
	if payload.Quantity < 1 || payload.Quantity > 100 {
		return domain.GenerateBulkAssetTagsResponse{}, domain.ErrBadRequest("quantity must be between 1 and 100")
	}

	// * Get category to retrieve CategoryCode
	category, err := s.CategoryService.GetCategoryById(ctx, payload.CategoryID, "en")
	if err != nil {
		return domain.GenerateBulkAssetTagsResponse{}, err
	}

	// * Get the last asset tag for this category
	lastAssetTag, err := s.Repo.GetLastAssetTagByCategory(ctx, payload.CategoryID)
	if err != nil {
		return domain.GenerateBulkAssetTagsResponse{}, err
	}

	// * Calculate starting increment
	startIncrement := 1
	if lastAssetTag != "" {
		// Use LastIndex to handle multi-dash category codes like FURN-SEAT
		dashIndex := strings.LastIndex(lastAssetTag, "-")
		if dashIndex != -1 && dashIndex < len(lastAssetTag)-1 {
			numericPart := lastAssetTag[dashIndex+1:]
			var parsedNum int
			_, err := fmt.Sscanf(numericPart, "%d", &parsedNum)
			if err == nil {
				startIncrement = parsedNum + 1
			}
		}
	}

	// * Generate sequential tags
	tags := make([]string, payload.Quantity)
	for i := 0; i < payload.Quantity; i++ {
		tags[i] = fmt.Sprintf("%s-%05d", category.CategoryCode, startIncrement+i)
	}

	endIncrement := startIncrement + payload.Quantity - 1

	return domain.GenerateBulkAssetTagsResponse{
		CategoryCode:   category.CategoryCode,
		LastAssetTag:   lastAssetTag,
		StartTag:       tags[0],
		EndTag:         tags[len(tags)-1],
		Tags:           tags,
		Quantity:       payload.Quantity,
		StartIncrement: startIncrement,
		EndIncrement:   endIncrement,
	}, nil
}

// UploadBulkDataMatrixImages uploads multiple data matrix images to Cloudinary
func (s *Service) UploadBulkDataMatrixImages(ctx context.Context, assetTags []string, files []*multipart.FileHeader) (domain.UploadBulkDataMatrixResponse, error) {
	if len(files) == 0 {
		return domain.UploadBulkDataMatrixResponse{}, domain.ErrBadRequest("at least one file is required")
	}

	if len(files) > 100 {
		return domain.UploadBulkDataMatrixResponse{}, domain.ErrBadRequest("maximum 100 files allowed")
	}

	if len(assetTags) != len(files) {
		return domain.UploadBulkDataMatrixResponse{}, domain.ErrBadRequest("number of asset tags must match number of files")
	}

	if s.CloudinaryClient == nil {
		return domain.UploadBulkDataMatrixResponse{}, domain.ErrInternal(fmt.Errorf("cloudinary client not configured"))
	}

	// Prepare public IDs for each file
	// Naming pattern: {assetTag}-{ulid}
	publicIDs := make([]string, len(assetTags))
	for i, tag := range assetTags {
		ulidStr := ulid.Make().String()
		publicIDs[i] = fmt.Sprintf("%s-%s", tag, ulidStr)
	}

	// Get bulk upload config for data matrix images
	baseConfig := cloudinary.GetBulkDataMatrixImageUploadConfig()

	log.Printf("Starting bulk upload of %d files to Cloudinary", len(files))

	// Upload all files using efficient bulk upload method
	uploadResult, err := s.CloudinaryClient.UploadMultipleFilesWithPublicIDs(ctx, files, publicIDs, baseConfig)
	if err != nil {
		log.Printf("ERROR: Bulk upload to Cloudinary failed: %v", err)
		return domain.UploadBulkDataMatrixResponse{}, domain.ErrInternal(err)
	}

	log.Printf("Cloudinary upload completed: %d succeeded, %d failed", len(uploadResult.Results), len(uploadResult.Failed))

	// Log detailed failures
	for _, failure := range uploadResult.Failed {
		log.Printf("Upload failed for file '%s': %s", failure.FileName, failure.Error)
	}

	// Extract URLs from successful uploads using index-based mapping
	// Since uploads are sequential, Results array order matches files array order
	urls := make([]string, len(files))
	uploadedCount := len(uploadResult.Results)

	// Map results by index (most reliable method)
	for i, result := range uploadResult.Results {
		if i < len(urls) {
			urls[i] = result.SecureURL
			log.Printf("SUCCESS: Mapped file[%d] '%s' (tag: %s) to URL: %s", i, files[i].Filename, assetTags[i], result.SecureURL)
		}
	}

	// Check for any failed uploads and log them
	if len(uploadResult.Failed) > 0 {
		for i, failure := range uploadResult.Failed {
			log.Printf("FAILED: File[%d] '%s' upload failed: %s", i, failure.FileName, failure.Error)
		}
	}

	// Return error if all uploads failed
	if uploadedCount == 0 {
		return domain.UploadBulkDataMatrixResponse{}, domain.ErrInternal(fmt.Errorf("all image uploads failed"))
	}

	return domain.UploadBulkDataMatrixResponse{
		Urls:      urls,
		Count:     uploadedCount,
		AssetTags: assetTags,
	}, nil
}

// DeleteBulkDataMatrixImages deletes data matrix images from Cloudinary and nullifies DB field
func (s *Service) DeleteBulkDataMatrixImages(ctx context.Context, payload *domain.DeleteBulkDataMatrixPayload) (domain.DeleteBulkDataMatrixResponse, error) {
	if len(payload.AssetTags) == 0 {
		return domain.DeleteBulkDataMatrixResponse{}, domain.ErrBadRequest("at least one asset tag is required")
	}

	if len(payload.AssetTags) > 100 {
		return domain.DeleteBulkDataMatrixResponse{}, domain.ErrBadRequest("maximum 100 asset tags allowed")
	}

	if s.CloudinaryClient == nil {
		return domain.DeleteBulkDataMatrixResponse{}, domain.ErrInternal(fmt.Errorf("cloudinary client not configured"))
	}

	deletedCount := 0
	failedTags := []string{}

	// Process each asset tag
	for _, tag := range payload.AssetTags {
		// Get asset by tag
		asset, err := s.Repo.GetAssetByAssetTag(ctx, tag)
		if err != nil {
			failedTags = append(failedTags, tag)
			log.Printf("Failed to get asset with tag %s: %v", tag, err)
			continue
		}

		// Skip if asset has no data matrix image
		if asset.DataMatrixImageUrl == "" {
			failedTags = append(failedTags, tag)
			log.Printf("Asset with tag %s has no data matrix image", tag)
			continue
		}

		// Extract public ID from Cloudinary URL
		publicID := cloudinary.ExtractPublicIDFromURL(asset.DataMatrixImageUrl)
		if publicID == "" {
			failedTags = append(failedTags, tag)
			log.Printf("Failed to extract public ID from URL for tag %s", tag)
			continue
		}

		// Delete from Cloudinary
		if err := s.CloudinaryClient.DeleteFile(ctx, publicID); err != nil {
			log.Printf("Failed to delete file from Cloudinary for tag %s: %v", tag, err)
			// Continue anyway to nullify DB field
		}

		// Update asset to remove data matrix image URL from DB
		emptyString := ""
		updatePayload := &domain.UpdateAssetPayload{
			DataMatrixImageUrl: &emptyString,
		}

		_, err = s.Repo.UpdateAsset(ctx, asset.ID, updatePayload)
		if err != nil {
			failedTags = append(failedTags, tag)
			log.Printf("Failed to update asset with tag %s: %v", tag, err)
			continue
		}

		deletedCount++
	}

	return domain.DeleteBulkDataMatrixResponse{
		DeletedCount: deletedCount,
		FailedTags:   failedTags,
		AssetTags:    payload.AssetTags,
	}, nil
}

// *===========================HELPER METHODS===========================*

// uploadAndAttachAssetImages uploads images to Cloudinary and attaches them to an asset
// Uses many-to-many relationship with image reusability via public_id deduplication
func (s *Service) uploadAndAttachAssetImages(ctx context.Context, assetID string, imageFiles []*multipart.FileHeader) error {
	if len(imageFiles) == 0 {
		return nil
	}

	if s.CloudinaryClient == nil {
		return fmt.Errorf("cloudinary client not configured")
	}

	// Get upload config for asset images
	uploadConfig := cloudinary.GetAssetImageUploadConfig()

	log.Printf("Starting upload of %d asset images to Cloudinary", len(imageFiles))

	// Upload all files to Cloudinary (auto-generated publicIDs for potential reuse)
	uploadResult, err := s.CloudinaryClient.UploadMultipleFiles(ctx, imageFiles, uploadConfig)
	if err != nil {
		log.Printf("ERROR: Asset images upload to Cloudinary failed: %v", err)
		return fmt.Errorf("failed to upload images to cloudinary: %w", err)
	}

	log.Printf("Cloudinary upload completed: %d succeeded, %d failed", len(uploadResult.Results), len(uploadResult.Failed))

	// Log detailed failures
	for _, failure := range uploadResult.Failed {
		log.Printf("Upload failed for file '%s': %s", failure.FileName, failure.Error)
	}

	// Return error if all uploads failed
	if len(uploadResult.Results) == 0 {
		return fmt.Errorf("all %d image uploads failed", len(imageFiles))
	}

	// Process successful uploads and attach to asset
	var imageIDs []string
	var displayOrders []int

	for i, result := range uploadResult.Results {
		// Check if image with this public_id already exists (deduplication)
		existingImage, err := s.Repo.GetImageByPublicID(ctx, result.PublicID)
		if err != nil {
			log.Printf("Error checking existing image: %v", err)
		}

		var imageID string
		if existingImage != nil {
			// Image already exists in pool, reuse it!
			imageID = existingImage.ID
			log.Printf("✓ Reusing existing image: %s (public_id: %s)", imageID, result.PublicID)
		} else {
			// Create new image record in pool
			newImage, err := s.Repo.CreateImage(ctx, result.SecureURL, &result.PublicID)
			if err != nil {
				log.Printf("Failed to create image record for %s: %v", result.PublicID, err)
				continue
			}
			imageID = newImage.ID
			log.Printf("✓ Created new image: %s (public_id: %s)", imageID, result.PublicID)
		}

		imageIDs = append(imageIDs, imageID)
		displayOrders = append(displayOrders, i)
	}

	// Attach images to asset via junction table (first image is primary by default)
	if len(imageIDs) > 0 {
		_, err := s.Repo.AttachImagesToAsset(ctx, assetID, imageIDs, displayOrders, 0)
		if err != nil {
			return fmt.Errorf("failed to attach images to asset: %w", err)
		}
		log.Printf("✓ Attached %d images to asset %s", len(imageIDs), assetID)
	}

	return nil
}

// *===========================TEMPLATE IMAGES (FOR BULK CREATE)===========================*

// UploadTemplateImages uploads images to Cloudinary for later reuse in bulk asset creation
// These images are not attached to any asset yet, just stored in images table
func (s *Service) UploadTemplateImages(ctx context.Context, files []*multipart.FileHeader) (domain.UploadTemplateImagesResponse, error) {
	if len(files) == 0 {
		return domain.UploadTemplateImagesResponse{}, domain.ErrBadRequest("at least one image file is required")
	}

	if len(files) > 10 {
		return domain.UploadTemplateImagesResponse{}, domain.ErrBadRequest("maximum 10 template images per request")
	}

	if s.CloudinaryClient == nil {
		return domain.UploadTemplateImagesResponse{}, domain.ErrInternalWithMessage("cloudinary client not configured")
	}

	// Get upload config for template images (same folder as regular asset images)
	uploadConfig := cloudinary.GetAssetImageUploadConfig()
	// No need to override folder, use default sigma-asset/assets

	log.Printf("Starting upload of %d template images to Cloudinary", len(files))

	// Upload all files to Cloudinary (auto-generated publicIDs for reusability)
	uploadResult, err := s.CloudinaryClient.UploadMultipleFiles(ctx, files, uploadConfig)
	if err != nil {
		return domain.UploadTemplateImagesResponse{}, domain.ErrInternalWithMessage(fmt.Sprintf("failed to upload template images: %v", err))
	}

	log.Printf("Cloudinary upload completed: %d succeeded, %d failed", len(uploadResult.Results), len(uploadResult.Failed))

	// Log detailed failures
	for _, failure := range uploadResult.Failed {
		log.Printf("Failed to upload template image %s: %v", failure.FileName, failure.Error)
	}

	// Return error if all uploads failed
	if len(uploadResult.Results) == 0 {
		return domain.UploadTemplateImagesResponse{}, domain.ErrInternalWithMessage("all template image uploads failed")
	}

	// Process successful uploads and save to images table
	imageUrls := make([]string, 0, len(uploadResult.Results))
	for _, result := range uploadResult.Results {
		// Check if image with this public_id already exists (deduplication)
		existingImage, err := s.Repo.GetImageByPublicID(ctx, result.PublicID)
		if err == nil && existingImage != nil {
			// Image already exists, reuse it
			log.Printf("Template image already exists with public_id %s, reusing", result.PublicID)
			imageUrls = append(imageUrls, existingImage.ImageURL)
			continue
		}

		// Save new image to database for reusability
		publicID := result.PublicID
		newImage, err := s.Repo.CreateImage(ctx, result.SecureURL, &publicID)
		if err != nil {
			log.Printf("Failed to save template image to DB (public_id: %s): %v", result.PublicID, err)
			continue
		}
		imageUrls = append(imageUrls, newImage.ImageURL)
	}

	if len(imageUrls) == 0 {
		return domain.UploadTemplateImagesResponse{}, domain.ErrInternalWithMessage("no template images saved successfully")
	}

	return domain.UploadTemplateImagesResponse{
		ImageUrls: imageUrls,
		Count:     len(imageUrls),
	}, nil
}

// attachImageUrlsToAsset attaches existing images (by URL) to an asset via junction table
// Used for reusing uploaded template images across multiple assets
func (s *Service) attachImageUrlsToAsset(ctx context.Context, assetID string, imageUrls []string) error {
	if len(imageUrls) == 0 {
		return nil
	}

	// Find existing images by URLs
	imageIDs := make([]string, 0, len(imageUrls))
	displayOrders := make([]int, 0, len(imageUrls))

	for i, imageUrl := range imageUrls {
		// Extract public_id from Cloudinary URL for lookup
		publicID := cloudinary.ExtractPublicIDFromURL(imageUrl)
		if publicID == "" {
			log.Printf("Failed to extract public_id from URL: %s", imageUrl)
			continue
		}

		// Find image by public_id
		image, err := s.Repo.GetImageByPublicID(ctx, publicID)
		if err != nil {
			log.Printf("Image not found for public_id %s: %v", publicID, err)
			continue
		}

		imageIDs = append(imageIDs, image.ID)
		displayOrders = append(displayOrders, i+1)
	}

	if len(imageIDs) == 0 {
		return domain.ErrBadRequest("no valid images found for provided URLs")
	}

	// Attach images to asset (first image is primary)
	_, err := s.Repo.AttachImagesToAsset(ctx, assetID, imageIDs, displayOrders, 0)
	if err != nil {
		return domain.ErrInternalWithMessage(fmt.Sprintf("failed to attach images to asset: %v", err))
	}

	return nil
}

// *===========================ASSET IMAGES (INDEPENDENT OPERATIONS)===========================*

// GetAvailableAssetImages retrieves available asset images (from sigma-asset/assets folder only) that can be reused
// Uses cursor-based pagination for efficient scrolling through large image pools
func (s *Service) GetAvailableAssetImages(ctx context.Context, limit int, cursor string) ([]domain.ImageResponse, error) {
	images, err := s.Repo.GetAvailableAssetImages(ctx, limit, cursor)
	if err != nil {
		return nil, err
	}

	// Convert to response
	responses := make([]domain.ImageResponse, len(images))
	for i, img := range images {
		responses[i] = domain.ImageResponse{
			ID:        img.ID,
			ImageURL:  img.ImageURL,
			CreatedAt: img.CreatedAt,
			UpdatedAt: img.UpdatedAt,
		}
	}

	return responses, nil
}

// UploadBulkAssetImages uploads images to Cloudinary and attaches them to their respective assets
// This is an independent operation that can be called separately from asset creation/update
func (s *Service) UploadBulkAssetImages(ctx context.Context, assetIds []string, files []*multipart.FileHeader) (domain.UploadBulkAssetImagesResponse, error) {
	// Validate inputs
	if len(assetIds) == 0 || len(files) == 0 {
		return domain.UploadBulkAssetImagesResponse{}, domain.ErrBadRequest("assetIds and files are required")
	}

	if len(assetIds) != len(files) {
		return domain.UploadBulkAssetImagesResponse{}, domain.ErrBadRequest("number of assetIds must match number of files")
	}

	if len(assetIds) > 100 {
		return domain.UploadBulkAssetImagesResponse{}, domain.ErrBadRequest("maximum 100 images per request")
	}

	// Check Cloudinary client availability
	if s.CloudinaryClient == nil {
		return domain.UploadBulkAssetImagesResponse{}, domain.ErrInternalWithMessage("cloudinary client not configured")
	}

	// Verify all assets exist before uploading
	for _, assetID := range assetIds {
		exists, err := s.Repo.CheckAssetExists(ctx, assetID)
		if err != nil {
			return domain.UploadBulkAssetImagesResponse{}, fmt.Errorf("failed to verify asset %s: %w", assetID, err)
		}
		if !exists {
			return domain.UploadBulkAssetImagesResponse{}, domain.ErrNotFound("asset not found: " + assetID)
		}
	}

	// Get upload config
	uploadConfig := cloudinary.GetAssetImageUploadConfig()

	// Upload all files to Cloudinary
	uploadResult, err := s.CloudinaryClient.UploadMultipleFiles(ctx, files, uploadConfig)
	if err != nil {
		log.Printf("ERROR: Bulk asset images upload failed: %v", err)
		return domain.UploadBulkAssetImagesResponse{}, fmt.Errorf("failed to upload images to cloudinary: %w", err)
	}

	// Process results
	results := make([]domain.AssetImageUploadResult, len(assetIds))
	successCount := 0

	// Process successful uploads
	for i, result := range uploadResult.Results {
		assetID := assetIds[i]

		// Check if image already exists (deduplication)
		existingImage, _ := s.Repo.GetImageByPublicID(ctx, result.PublicID)

		var imageID string
		if existingImage != nil {
			imageID = existingImage.ID
			log.Printf("✓ Reusing existing image for asset %s (public_id: %s)", assetID, result.PublicID)
		} else {
			// Create new image record
			newImage, err := s.Repo.CreateImage(ctx, result.SecureURL, &result.PublicID)
			if err != nil {
				results[i] = domain.AssetImageUploadResult{
					AssetID: assetID,
					Success: false,
					Error:   fmt.Sprintf("failed to create image record: %v", err),
				}
				continue
			}
			imageID = newImage.ID
			log.Printf("✓ Created new image for asset %s (public_id: %s)", assetID, result.PublicID)
		}

		// Attach image to asset (as non-primary by default to avoid conflict)
		_, err = s.Repo.AttachImagesToAsset(ctx, assetID, []string{imageID}, []int{0}, -1)
		if err != nil {
			results[i] = domain.AssetImageUploadResult{
				AssetID: assetID,
				Success: false,
				Error:   fmt.Sprintf("failed to attach image: %v", err),
			}
			continue
		}

		results[i] = domain.AssetImageUploadResult{
			AssetID:  assetID,
			ImageURL: result.SecureURL,
			Success:  true,
		}
		successCount++
	}

	// Process failed uploads
	for _, failure := range uploadResult.Failed {
		// Find corresponding assetID by index
		idx := -1
		for i, file := range files {
			if file.Filename == failure.FileName {
				idx = i
				break
			}
		}

		if idx >= 0 && idx < len(assetIds) {
			results[idx] = domain.AssetImageUploadResult{
				AssetID: assetIds[idx],
				Success: false,
				Error:   failure.Error,
			}
		}
	}

	return domain.UploadBulkAssetImagesResponse{
		Results: results,
		Count:   successCount,
	}, nil
}

// DeleteBulkAssetImages removes asset_images junction records and cleans up orphaned images
// This is an independent operation that can be called separately from asset update
func (s *Service) DeleteBulkAssetImages(ctx context.Context, payload *domain.DeleteBulkAssetImagesPayload) (domain.DeleteBulkAssetImagesResponse, error) {
	if len(payload.AssetImageIDs) == 0 {
		return domain.DeleteBulkAssetImagesResponse{}, domain.ErrBadRequest("assetImageIds are required")
	}

	if len(payload.AssetImageIDs) > 100 {
		return domain.DeleteBulkAssetImagesResponse{}, domain.ErrBadRequest("maximum 100 asset images per request")
	}

	deletedCount := 0
	failedIDs := []string{}

	// Delete each asset_image junction record
	for _, assetImageID := range payload.AssetImageIDs {
		err := s.Repo.DetachImageFromAsset(ctx, assetImageID)
		if err != nil {
			log.Printf("Failed to detach asset image %s: %v", assetImageID, err)
			failedIDs = append(failedIDs, assetImageID)
			continue
		}
		deletedCount++
	}

	// Clean up orphaned images (images no longer attached to any asset)
	orphanedCleaned := 0
	if err := s.Repo.DeleteUnusedImages(ctx); err != nil {
		log.Printf("Warning: Failed to clean up orphaned images: %v", err)
	} else {
		// Note: We don't have a count of deleted orphaned images from the query
		// This is just a cleanup operation
		log.Printf("✓ Cleaned up orphaned images")
	}

	return domain.DeleteBulkAssetImagesResponse{
		DeletedCount:    deletedCount,
		FailedIDs:       failedIDs,
		AssetImageIDs:   payload.AssetImageIDs,
		OrphanedCleaned: orphanedCleaned,
	}, nil
}

// sendUpdateNotifications sends all relevant notifications when asset is updated
func (s *Service) sendUpdateNotifications(ctx context.Context, oldAsset, newAsset *domain.Asset, payload *domain.UpdateAssetPayload) {
	// Skip if notification service is not available
	if s.NotificationService == nil {
		return
	}

	// 1. Check for assignment changes
	if payload.AssignedTo != nil {
		if *payload.AssignedTo != "" && (oldAsset.AssignedTo == nil || *oldAsset.AssignedTo != *payload.AssignedTo) {
			// Asset was assigned to a new user
			s.sendAssetAssignmentNotification(ctx, newAsset, *payload.AssignedTo, false)
		} else if *payload.AssignedTo == "" && oldAsset.AssignedTo != nil && *oldAsset.AssignedTo != "" {
			// Asset was unassigned
			s.sendAssetUnassignmentNotification(ctx, newAsset, *oldAsset.AssignedTo)
		}
	}

	// 2. Check for status changes
	if payload.Status != nil && *payload.Status != oldAsset.Status {
		s.sendAssetStatusChangeNotification(ctx, newAsset, oldAsset.Status, *payload.Status)
	}

	// 3. Check for condition changes
	if payload.Condition != nil && *payload.Condition != oldAsset.Condition {
		s.sendAssetConditionChangeNotification(ctx, newAsset, oldAsset.Condition, *payload.Condition)
	}

	// 4. Check for location changes (if needed, you might want to track this)
	// This would require fetching location names, so skipping for now
}

// sendAssetAssignmentNotification sends notification when asset is assigned to a user
func (s *Service) sendAssetAssignmentNotification(ctx context.Context, asset *domain.Asset, userId string, isNewAsset bool) {
	// Skip if notification service is not available
	if s.NotificationService == nil {
		log.Printf("Notification service not available, skipping asset assignment notification for asset ID: %s, user ID: %s", asset.ID, userId)
		return
	}

	log.Printf("Sending asset assignment notification for asset ID: %s, asset tag: %s, user ID: %s, is new asset: %t", asset.ID, asset.AssetTag, userId, isNewAsset)

	assetIdStr := asset.ID

	// Get notification message keys and params using helper function
	titleKey, messageKey, params := messages.AssetAssignmentNotification(asset.AssetName, asset.AssetTag, isNewAsset)

	// Get translations for all supported languages
	msgTranslations := messages.GetAssetNotificationTranslations(titleKey, messageKey, params)

	// Convert to domain translations
	translations := make([]domain.CreateNotificationTranslationPayload, len(msgTranslations))
	for i, t := range msgTranslations {
		translations[i] = domain.CreateNotificationTranslationPayload{
			LangCode: t.LangCode,
			Title:    t.Title,
			Message:  t.Message,
		}
	}

	entityType := "asset"
	priority := domain.NotificationPriorityNormal

	notificationPayload := &domain.CreateNotificationPayload{
		UserID:            userId,
		RelatedEntityType: &entityType,
		RelatedEntityID:   &assetIdStr,
		RelatedAssetID:    &assetIdStr,
		Type:              domain.NotificationTypeStatusChange,
		Priority:          priority,
		Translations:      translations,
	}

	// Send notification (errors are logged internally, won't block asset creation/update)
	_, err := s.NotificationService.CreateNotification(ctx, notificationPayload)
	if err != nil {
		log.Printf("Failed to create asset assignment notification for asset ID: %s, user ID: %s: %v", asset.ID, userId, err)
	} else {
		log.Printf("Successfully created asset assignment notification for asset ID: %s, user ID: %s", asset.ID, userId)
	}
}

// sendAssetUnassignmentNotification sends notification when asset is unassigned from a user
func (s *Service) sendAssetUnassignmentNotification(ctx context.Context, asset *domain.Asset, userId string) {
	if s.NotificationService == nil {
		log.Printf("Notification service not available, skipping asset unassignment notification for asset ID: %s, user ID: %s", asset.ID, userId)
		return
	}

	log.Printf("Sending asset unassignment notification for asset ID: %s, asset tag: %s, user ID: %s", asset.ID, asset.AssetTag, userId)

	assetIdStr := asset.ID
	titleKey, messageKey, params := messages.AssetUnassignmentNotification(asset.AssetName, asset.AssetTag)
	utilTranslations := messages.GetAssetNotificationTranslations(titleKey, messageKey, params)

	// Convert to domain translations
	translations := make([]domain.CreateNotificationTranslationPayload, len(utilTranslations))
	for i, t := range utilTranslations {
		translations[i] = domain.CreateNotificationTranslationPayload{
			LangCode: t.LangCode,
			Title:    t.Title,
			Message:  t.Message,
		}
	}

	entityType := "asset"
	priority := domain.NotificationPriorityNormal

	notificationPayload := &domain.CreateNotificationPayload{
		UserID:            userId,
		RelatedEntityType: &entityType,
		RelatedEntityID:   &assetIdStr,
		RelatedAssetID:    &assetIdStr,
		Type:              domain.NotificationTypeStatusChange,
		Priority:          priority,
		Translations:      translations,
	}

	_, err := s.NotificationService.CreateNotification(ctx, notificationPayload)
	if err != nil {
		log.Printf("Failed to create asset unassignment notification for asset ID: %s, user ID: %s: %v", asset.ID, userId, err)
	} else {
		log.Printf("Successfully created asset unassignment notification for asset ID: %s, user ID: %s", asset.ID, userId)
	}
}

// sendAssetStatusChangeNotification sends notification when asset status changes
func (s *Service) sendAssetStatusChangeNotification(ctx context.Context, asset *domain.Asset, oldStatus, newStatus domain.AssetStatus) {
	if s.NotificationService == nil {
		log.Printf("Notification service not available, skipping asset status change notification for asset ID: %s", asset.ID)
		return
	}

	// Notify assigned user if exists
	if asset.AssignedTo == nil || *asset.AssignedTo == "" {
		log.Printf("Asset not assigned to any user, skipping status change notification for asset ID: %s", asset.ID)
		return
	}

	log.Printf("Sending asset status change notification for asset ID: %s, asset tag: %s, old status: %s, new status: %s, user ID: %s", asset.ID, asset.AssetTag, oldStatus, newStatus, *asset.AssignedTo)

	assetIdStr := asset.ID
	titleKey, messageKey, params := messages.AssetStatusChangeNotification(asset.AssetName, string(oldStatus), string(newStatus))

	// Add assetTag to params if not already there
	params["assetTag"] = asset.AssetTag

	utilTranslations := messages.GetAssetNotificationTranslations(titleKey, messageKey, params)

	// Convert to domain translations
	translations := make([]domain.CreateNotificationTranslationPayload, len(utilTranslations))
	for i, t := range utilTranslations {
		translations[i] = domain.CreateNotificationTranslationPayload{
			LangCode: t.LangCode,
			Title:    t.Title,
			Message:  t.Message,
		}
	}

	entityType := "asset"
	priority := domain.NotificationPriorityNormal

	notificationPayload := &domain.CreateNotificationPayload{
		UserID:            *asset.AssignedTo,
		RelatedEntityType: &entityType,
		RelatedEntityID:   &assetIdStr,
		RelatedAssetID:    &assetIdStr,
		Type:              domain.NotificationTypeStatusChange,
		Priority:          priority,
		Translations:      translations,
	}

	_, err := s.NotificationService.CreateNotification(ctx, notificationPayload)
	if err != nil {
		log.Printf("Failed to create asset status change notification for asset ID: %s, user ID: %s: %v", asset.ID, *asset.AssignedTo, err)
	} else {
		log.Printf("Successfully created asset status change notification for asset ID: %s, user ID: %s", asset.ID, *asset.AssignedTo)
	}
}

// sendAssetConditionChangeNotification sends notification when asset condition changes
func (s *Service) sendAssetConditionChangeNotification(ctx context.Context, asset *domain.Asset, oldCondition, newCondition domain.AssetCondition) {
	if s.NotificationService == nil {
		log.Printf("Notification service not available, skipping asset condition change notification for asset ID: %s", asset.ID)
		return
	}

	// Notify assigned user if exists
	if asset.AssignedTo == nil || *asset.AssignedTo == "" {
		log.Printf("Asset not assigned to any user, skipping condition change notification for asset ID: %s", asset.ID)
		return
	}

	log.Printf("Sending asset condition change notification for asset ID: %s, asset tag: %s, old condition: %s, new condition: %s, user ID: %s", asset.ID, asset.AssetTag, oldCondition, newCondition, *asset.AssignedTo)

	assetIdStr := asset.ID
	titleKey, messageKey, params := messages.AssetConditionChangeNotification(asset.AssetName, asset.AssetTag, string(oldCondition), string(newCondition))
	utilTranslations := messages.GetAssetNotificationTranslations(titleKey, messageKey, params)

	// Convert to domain translations
	translations := make([]domain.CreateNotificationTranslationPayload, len(utilTranslations))
	for i, t := range utilTranslations {
		translations[i] = domain.CreateNotificationTranslationPayload{
			LangCode: t.LangCode,
			Title:    t.Title,
			Message:  t.Message,
		}
	}

	entityType := "asset"
	priority := domain.NotificationPriorityNormal

	notificationPayload := &domain.CreateNotificationPayload{
		UserID:            *asset.AssignedTo,
		RelatedEntityType: &entityType,
		RelatedEntityID:   &assetIdStr,
		RelatedAssetID:    &assetIdStr,
		Type:              domain.NotificationTypeStatusChange,
		Priority:          priority,
		Translations:      translations,
	}

	_, err := s.NotificationService.CreateNotification(ctx, notificationPayload)
	if err != nil {
		log.Printf("Failed to create asset condition change notification for asset ID: %s, user ID: %s: %v", asset.ID, *asset.AssignedTo, err)
	} else {
		log.Printf("Successfully created asset condition change notification for asset ID: %s, user ID: %s", asset.ID, *asset.AssignedTo)
	}
}

// sendHighValueAssetNotification sends notification for high-value asset creation
func (s *Service) sendHighValueAssetNotification(ctx context.Context, asset *domain.Asset, recipientUserId string) {
	if s.NotificationService == nil {
		log.Printf("Notification service not available, skipping high-value asset notification for asset ID: %s", asset.ID)
		return
	}

	log.Printf("Sending high-value asset notification for asset ID: %s, asset tag: %s, user ID: %s",
		asset.ID, asset.AssetTag, recipientUserId)

	assetIdStr := asset.ID

	// Format purchase price
	value := "N/A"
	if asset.PurchasePrice != nil {
		value = fmt.Sprintf("%.2f", *asset.PurchasePrice)
	}

	titleKey, messageKey, params := messages.AssetHighValueNotification(asset.AssetName, asset.AssetTag, value)
	utilTranslations := messages.GetAssetNotificationTranslations(titleKey, messageKey, params)

	// Convert to domain translations
	translations := make([]domain.CreateNotificationTranslationPayload, len(utilTranslations))
	for i, t := range utilTranslations {
		translations[i] = domain.CreateNotificationTranslationPayload{
			LangCode: t.LangCode,
			Title:    t.Title,
			Message:  t.Message,
		}
	}

	entityType := "asset"
	priority := domain.NotificationPriorityHigh

	notificationPayload := &domain.CreateNotificationPayload{
		UserID:            recipientUserId,
		RelatedEntityType: &entityType,
		RelatedEntityID:   &assetIdStr,
		RelatedAssetID:    &assetIdStr,
		Type:              domain.NotificationTypeStatusChange,
		Priority:          priority,
		Translations:      translations,
	}

	_, err := s.NotificationService.CreateNotification(ctx, notificationPayload)
	if err != nil {
		log.Printf("Failed to create high-value asset notification for asset ID: %s, user ID: %s: %v", asset.ID, recipientUserId, err)
	} else {
		log.Printf("Successfully created high-value asset notification for asset ID: %s, user ID: %s", asset.ID, recipientUserId)
	}
}

// sendHighValueAssetNotificationToAdmins sends notification for high-value asset creation to all admin users
func (s *Service) sendHighValueAssetNotificationToAdmins(ctx context.Context, asset *domain.Asset) {
	if s.NotificationService == nil {
		log.Printf("Notification service not available, skipping high-value asset notification for asset ID: %s", asset.ID)
		return
	}

	if s.UserRepo == nil {
		log.Printf("User repository not available, skipping high-value asset notification for asset ID: %s", asset.ID)
		return
	}

	log.Printf("Sending high-value asset notification to admins for asset ID: %s, asset tag: %s", asset.ID, asset.AssetTag)

	// Get all admin users
	adminRole := domain.RoleAdmin
	userParams := domain.UserParams{
		Filters: &domain.UserFilterOptions{
			Role: &adminRole,
		},
	}
	admins, err := s.UserRepo.GetUsersPaginated(ctx, userParams)
	if err != nil {
		log.Printf("Failed to get admin users for high-value asset notification: %v", err)
		return
	}

	if len(admins) == 0 {
		log.Printf("No admin users found, skipping high-value asset notification for asset ID: %s", asset.ID)
		return
	}

	// Send notification to each admin
	for _, admin := range admins {
		s.sendHighValueAssetNotification(ctx, asset, admin.ID)
	}

	log.Printf("Successfully sent high-value asset notification to %d admin(s) for asset ID: %s", len(admins), asset.ID)
}
