package asset

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
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
	GetAssetsForExport(ctx context.Context, params domain.AssetParams, langCode string) ([]domain.Asset, error)
	GetAssetsWithWarrantyExpiring(ctx context.Context, daysFromNow int) ([]domain.Asset, error)
	GetAssetsWithExpiredWarranty(ctx context.Context) ([]domain.Asset, error)
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
			// Generate temporary asset ID for image naming
			tempAssetID := "temp_" + ulid.Make().String()
			uploadConfig := cloudinary.GetDataMatrixImageUploadConfig()
			publicID := "asset_" + tempAssetID + "_datamatrix"
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

	// * Send notification if asset is assigned to a user
	if payload.AssignedTo != nil && *payload.AssignedTo != "" {
		go s.sendAssetAssignmentNotification(context.Background(), &createdAsset, *payload.AssignedTo, true)
	}

	// * Send notification if asset is high-value
	if payload.PurchasePrice != nil && *payload.PurchasePrice > 10000 {
		go s.sendHighValueAssetNotificationToAdmins(context.Background(), &createdAsset)
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
			createdAsset, _ = s.Repo.UpdateAsset(ctx, createdAsset.ID, updatePayload)
		}
		// Note: We don't return error here to avoid failing asset creation if image re-upload fails
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
			if parsedDate, err := time.Parse("2006-01-02", *assetPayload.PurchaseDate); err == nil {
				purchaseDate = &parsedDate
			}
		}

		var warrantyEnd *time.Time
		if assetPayload.WarrantyEnd != nil && *assetPayload.WarrantyEnd != "" {
			if parsedDate, err := time.Parse("2006-01-02", *assetPayload.WarrantyEnd); err == nil {
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

	for i := range createdAssets {
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
				// Provide detailed error message
				errorMsg := "Failed to upload data matrix image: " + err.Error()
				return domain.AssetResponse{}, domain.ErrBadRequest(errorMsg)
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

	// Use the UpdateAsset method
	_, err = s.Repo.UpdateAsset(ctx, assetId, payload)
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
	updatedAsset, err := s.Repo.UpdateAsset(ctx, assetId, payload)
	if err != nil {
		return domain.AssetResponse{}, err
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
		// Format expected: CATEGORYCODE-00001
		dashIndex := strings.Index(lastAssetTag, "-")
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

// *===========================HELPER METHODS===========================*

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
