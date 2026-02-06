package category

import (
	"context"
	"log"
	"mime/multipart"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/client/cloudinary"
	"github.com/Rizz404/inventory-api/internal/client/gtranslate"
	"github.com/Rizz404/inventory-api/internal/notification/messages"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/Rizz404/inventory-api/internal/utils"
	"github.com/oklog/ulid/v2"
)

// * Repository interface defines the contract for category data operations
type Repository interface {
	// * MUTATION
	CreateCategory(ctx context.Context, payload *domain.Category) (domain.Category, error)
	BulkCreateCategories(ctx context.Context, categories []domain.Category) ([]domain.Category, error)
	UpdateCategory(ctx context.Context, categoryId string, payload *domain.UpdateCategoryPayload) (domain.Category, error)
	DeleteCategory(ctx context.Context, categoryId string) error
	BulkDeleteCategories(ctx context.Context, categoryIds []string) (domain.BulkDeleteCategories, error)
	AddCategoryTranslations(ctx context.Context, categoryId string, translations []domain.CategoryTranslation) error

	// * QUERY
	GetCategoriesPaginated(ctx context.Context, params domain.CategoryParams, langCode string) ([]domain.Category, error)
	GetCategoriesCursor(ctx context.Context, params domain.CategoryParams, langCode string) ([]domain.Category, error)
	GetCategoryById(ctx context.Context, categoryId string) (domain.Category, error)
	GetCategoryByCode(ctx context.Context, categoryCode string) (domain.Category, error)
	CheckCategoryExist(ctx context.Context, categoryId string) (bool, error)
	CheckCategoryCodeExist(ctx context.Context, categoryCode string) (bool, error)
	CheckCategoryCodeExistExcluding(ctx context.Context, categoryCode string, excludeCategoryId string) (bool, error)
	CountCategories(ctx context.Context, params domain.CategoryParams) (int64, error)
	GetCategoryStatistics(ctx context.Context) (domain.CategoryStatistics, error)
}

// * CategoryService interface defines the contract for category business operations
type CategoryService interface {
	// * MUTATION
	CreateCategory(ctx context.Context, payload *domain.CreateCategoryPayload, imageFile *multipart.FileHeader) (domain.CategoryResponse, error)
	BulkCreateCategories(ctx context.Context, payload *domain.BulkCreateCategoriesPayload) (domain.BulkCreateCategoriesResponse, error)
	UpdateCategory(ctx context.Context, categoryId string, payload *domain.UpdateCategoryPayload, imageFile *multipart.FileHeader, langCode string) (domain.CategoryResponse, error)
	DeleteCategory(ctx context.Context, categoryId string) error
	BulkDeleteCategories(ctx context.Context, payload *domain.BulkDeleteCategoriesPayload) (domain.BulkDeleteCategoriesResponse, error)

	// * QUERY
	GetCategoriesPaginated(ctx context.Context, params domain.CategoryParams, langCode string) ([]domain.CategoryResponse, int64, error)
	GetCategoriesCursor(ctx context.Context, params domain.CategoryParams, langCode string) ([]domain.CategoryResponse, error)
	GetCategoryById(ctx context.Context, categoryId string, langCode string) (domain.CategoryResponse, error)
	GetCategoryByCode(ctx context.Context, categoryCode string, langCode string) (domain.CategoryResponse, error)
	CheckCategoryExists(ctx context.Context, categoryId string) (bool, error)
	CheckCategoryCodeExists(ctx context.Context, categoryCode string) (bool, error)
	CountCategories(ctx context.Context, params domain.CategoryParams) (int64, error)
	GetCategoryStatistics(ctx context.Context) (domain.CategoryStatisticsResponse, error)
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
	CloudinaryClient    *cloudinary.Client
	Translator          *gtranslate.Client
}

// * Ensure Service implements CategoryService interface
var _ CategoryService = (*Service)(nil)

func NewService(r Repository, notificationService NotificationService, userRepo UserRepository, cloudinaryClient *cloudinary.Client, translator *gtranslate.Client) CategoryService {
	return &Service{
		Repo:                r,
		NotificationService: notificationService,
		UserRepo:            userRepo,
		CloudinaryClient:    cloudinaryClient,
		Translator:          translator,
	}
}

// *===========================MUTATION===========================*
func (s *Service) CreateCategory(ctx context.Context, payload *domain.CreateCategoryPayload, imageFile *multipart.FileHeader) (domain.CategoryResponse, error) {
	// * Check if category code already exists
	if codeExists, err := s.Repo.CheckCategoryCodeExist(ctx, payload.CategoryCode); err != nil {
		return domain.CategoryResponse{}, err
	} else if codeExists {
		return domain.CategoryResponse{}, domain.ErrConflictWithKey(utils.ErrCategoryCodeExistsKey)
	}

	// * Check if parent category exists if parentId is provided
	if payload.ParentID != nil && *payload.ParentID != "" {
		if parentExists, err := s.Repo.CheckCategoryExist(ctx, *payload.ParentID); err != nil {
			return domain.CategoryResponse{}, err
		} else if !parentExists {
			return domain.CategoryResponse{}, domain.ErrNotFoundWithKey(utils.ErrCategoryNotFoundKey)
		}
	}

	// * Handle image upload if file is provided
	var imageURL *string
	if imageFile != nil {
		// Upload file to Cloudinary if client is available
		if s.CloudinaryClient != nil {
			// Generate temporary category ID for image naming
			tempCategoryID := "temp-" + ulid.Make().String()
			uploadConfig := cloudinary.GetCategoryImageUploadConfig()
			publicID := "category-" + tempCategoryID + "-image"
			uploadConfig.PublicID = &publicID

			uploadResult, err := s.CloudinaryClient.UploadSingleFile(ctx, imageFile, uploadConfig)
			if err != nil {
				// Provide detailed error message
				errorMsg := "Failed to upload category image: " + err.Error()
				return domain.CategoryResponse{}, domain.ErrBadRequest(errorMsg)
			}
			imageURL = &uploadResult.SecureURL
		} else {
			return domain.CategoryResponse{}, domain.ErrBadRequestWithKey(utils.ErrCloudinaryConfigKey)
		}
	} else if payload.ImageURL != nil {
		// Use provided image URL from JSON/form data
		imageURL = payload.ImageURL
	}

	// * Prepare domain category with user's input translations only
	newCategory := domain.Category{
		ParentID:     payload.ParentID,
		CategoryCode: payload.CategoryCode,
		ImageURL:     imageURL,
		Translations: make([]domain.CategoryTranslation, len(payload.Translations)),
	}

	// * Convert translation payloads to domain translations (only user input)
	for i, translationPayload := range payload.Translations {
		newCategory.Translations[i] = domain.CategoryTranslation{
			LangCode:     translationPayload.LangCode,
			CategoryName: translationPayload.CategoryName,
			Description:  translationPayload.Description,
		}
	}

	createdCategory, err := s.Repo.CreateCategory(ctx, &newCategory)
	if err != nil {
		return domain.CategoryResponse{}, err
	}

	// * Auto-translate missing languages in background if needed
	if len(payload.Translations) < 3 {
		go s.autoTranslateCategoryAsync(createdCategory.ID, payload.Translations)
	}

	// * Convert to CategoryResponse using mapper
	return mapper.CategoryToResponse(&createdCategory, mapper.DefaultLangCode), nil
}

func (s *Service) BulkCreateCategories(ctx context.Context, payload *domain.BulkCreateCategoriesPayload) (domain.BulkCreateCategoriesResponse, error) {
	if payload == nil || len(payload.Categories) == 0 {
		return domain.BulkCreateCategoriesResponse{}, domain.ErrBadRequest("categories payload is required")
	}

	codeSeen := make(map[string]struct{})
	for _, catPayload := range payload.Categories {
		if _, exists := codeSeen[catPayload.CategoryCode]; exists {
			return domain.BulkCreateCategoriesResponse{}, domain.ErrBadRequest("duplicate category code: " + catPayload.CategoryCode)
		}
		codeSeen[catPayload.CategoryCode] = struct{}{}

		// Check parent if provided
		if catPayload.ParentID != nil && *catPayload.ParentID != "" {
			if parentExists, err := s.Repo.CheckCategoryExist(ctx, *catPayload.ParentID); err != nil {
				return domain.BulkCreateCategoriesResponse{}, err
			} else if !parentExists {
				return domain.BulkCreateCategoriesResponse{}, domain.ErrNotFoundWithKey(utils.ErrCategoryNotFoundKey)
			}
		}
	}

	// Check all codes against database
	for code := range codeSeen {
		exists, err := s.Repo.CheckCategoryCodeExist(ctx, code)
		if err != nil {
			return domain.BulkCreateCategoriesResponse{}, err
		}
		if exists {
			return domain.BulkCreateCategoriesResponse{}, domain.ErrConflictWithKey(utils.ErrCategoryCodeExistsKey)
		}
	}

	categories := make([]domain.Category, len(payload.Categories))
	categoriesToTranslate := []struct {
		categoryID   string
		translations []domain.CreateCategoryTranslationPayload
	}{}

	for i, catPayload := range payload.Categories {
		cat := domain.Category{
			ParentID:     catPayload.ParentID,
			CategoryCode: catPayload.CategoryCode,
			Translations: make([]domain.CategoryTranslation, len(catPayload.Translations)),
		}

		// Convert user input translations only
		for j, transPayload := range catPayload.Translations {
			cat.Translations[j] = domain.CategoryTranslation{
				LangCode:     transPayload.LangCode,
				CategoryName: transPayload.CategoryName,
				Description:  transPayload.Description,
			}
		}
		categories[i] = cat
	}

	createdCategories, err := s.Repo.BulkCreateCategories(ctx, categories)
	if err != nil {
		return domain.BulkCreateCategoriesResponse{}, err
	}

	// * Collect categories that need auto-translation
	for i, catPayload := range payload.Categories {
		if len(catPayload.Translations) < 3 {
			categoriesToTranslate = append(categoriesToTranslate, struct {
				categoryID   string
				translations []domain.CreateCategoryTranslationPayload
			}{
				categoryID:   createdCategories[i].ID,
				translations: catPayload.Translations,
			})
		}
	}

	// * Auto-translate missing languages in background
	if len(categoriesToTranslate) > 0 {
		go func() {
			for _, item := range categoriesToTranslate {
				s.autoTranslateCategoryAsync(item.categoryID, item.translations)
			}
		}()
	}

	response := domain.BulkCreateCategoriesResponse{
		Categories: mapper.CategoriesToResponses(createdCategories, mapper.DefaultLangCode),
	}

	return response, nil
}

func (s *Service) UpdateCategory(ctx context.Context, categoryId string, payload *domain.UpdateCategoryPayload, imageFile *multipart.FileHeader, langCode string) (domain.CategoryResponse, error) {
	// * Check if category exists
	existingCategory, err := s.Repo.GetCategoryById(ctx, categoryId)
	if err != nil {
		return domain.CategoryResponse{}, err
	}

	// * Check category code uniqueness if being updated
	if payload.CategoryCode != nil {
		if codeExists, err := s.Repo.CheckCategoryCodeExistExcluding(ctx, *payload.CategoryCode, categoryId); err != nil {
			return domain.CategoryResponse{}, err
		} else if codeExists {
			return domain.CategoryResponse{}, domain.ErrConflictWithKey(utils.ErrCategoryCodeExistsKey)
		}
	}

	// * Check if parent category exists if parentId is provided
	if payload.ParentID != nil && *payload.ParentID != "" {
		if parentExists, err := s.Repo.CheckCategoryExist(ctx, *payload.ParentID); err != nil {
			return domain.CategoryResponse{}, err
		} else if !parentExists {
			return domain.CategoryResponse{}, domain.ErrNotFoundWithKey(utils.ErrCategoryNotFoundKey)
		}
	}

	// * Handle image upload if file is provided
	if imageFile != nil {
		// Upload file to Cloudinary if client is available
		if s.CloudinaryClient != nil {
			uploadConfig := cloudinary.GetCategoryImageUploadConfig()
			publicID := "category-" + categoryId + "-image"
			uploadConfig.PublicID = &publicID
			uploadConfig.Overwrite = true // Overwrite existing image

			uploadResult, err := s.CloudinaryClient.UploadSingleFile(ctx, imageFile, uploadConfig)
			if err != nil {
				// Provide detailed error message
				errorMsg := "Failed to upload category image: " + err.Error()
				return domain.CategoryResponse{}, domain.ErrBadRequest(errorMsg)
			}
			newImageURL := uploadResult.SecureURL
			payload.ImageURL = &newImageURL
		} else {
			return domain.CategoryResponse{}, domain.ErrBadRequestWithKey(utils.ErrCloudinaryConfigKey)
		}
	} else if payload.ImageURL != nil && *payload.ImageURL == "" {
		// If imageUrl is explicitly set to empty string, delete the old image from Cloudinary
		if s.CloudinaryClient != nil && existingCategory.ImageURL != nil {
			// Extract public ID from existing URL and delete (optional - Cloudinary has storage limits)
			// For now, we just set it to nil in the database
		}
	}

	// * Update category with user's input translations only
	updatedCategory, err := s.Repo.UpdateCategory(ctx, categoryId, payload)
	if err != nil {
		return domain.CategoryResponse{}, err
	}

	// * Auto-translate missing languages in background if translations were updated
	if len(payload.Translations) > 0 {
		// Get current translation count after update
		currentLangCodes := make([]string, 0)
		for _, trans := range existingCategory.Translations {
			currentLangCodes = append(currentLangCodes, trans.LangCode)
		}
		for _, trans := range payload.Translations {
			found := false
			for _, code := range currentLangCodes {
				if code == trans.LangCode {
					found = true
					break
				}
			}
			if !found {
				currentLangCodes = append(currentLangCodes, trans.LangCode)
			}
		}

		// Launch background translation if incomplete
		if len(currentLangCodes) < 3 {
			go s.autoTranslateUpdateCategoryAsync(categoryId, payload.Translations, updatedCategory.Translations)
		}
	}

	// * Send notification to all admin users
	s.sendCategoryUpdatedNotificationToAdmins(ctx, &updatedCategory)

	// * Convert to CategoryResponse using mapper with requested lang code
	return mapper.CategoryToResponse(&updatedCategory, langCode), nil
}

func (s *Service) DeleteCategory(ctx context.Context, categoryId string) error {
	err := s.Repo.DeleteCategory(ctx, categoryId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) BulkDeleteCategories(ctx context.Context, payload *domain.BulkDeleteCategoriesPayload) (domain.BulkDeleteCategoriesResponse, error) {
	// * Validate that IDs are provided
	if len(payload.IDS) == 0 {
		return domain.BulkDeleteCategoriesResponse{}, domain.ErrBadRequestWithKey(utils.ErrCategoryIDRequiredKey)
	}

	// * Perform bulk delete operation
	result, err := s.Repo.BulkDeleteCategories(ctx, payload.IDS)
	if err != nil {
		return domain.BulkDeleteCategoriesResponse{}, err
	}

	// * Convert to response
	response := domain.BulkDeleteCategoriesResponse{
		RequestedIDS: result.RequestedIDS,
		DeletedIDS:   result.DeletedIDS,
	}

	return response, nil
}

// *===========================QUERY===========================*
func (s *Service) GetCategoriesPaginated(ctx context.Context, params domain.CategoryParams, langCode string) ([]domain.CategoryResponse, int64, error) {
	categories, err := s.Repo.GetCategoriesPaginated(ctx, params, langCode)
	if err != nil {
		return nil, 0, err
	}

	// * Count total for pagination
	count, err := s.Repo.CountCategories(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// Convert Category to CategoryResponse using mapper (includes translations)
	responses := mapper.CategoriesToResponses(categories, langCode)

	return responses, count, nil
}

func (s *Service) GetCategoriesCursor(ctx context.Context, params domain.CategoryParams, langCode string) ([]domain.CategoryResponse, error) {
	categories, err := s.Repo.GetCategoriesCursor(ctx, params, langCode)
	if err != nil {
		return nil, err
	}

	// Convert Category to CategoryResponse using mapper (includes translations)
	responses := mapper.CategoriesToResponses(categories, langCode)

	return responses, nil
}

func (s *Service) GetCategoryById(ctx context.Context, categoryId string, langCode string) (domain.CategoryResponse, error) {
	category, err := s.Repo.GetCategoryById(ctx, categoryId)
	if err != nil {
		return domain.CategoryResponse{}, err
	}

	// * Convert to CategoryResponse using mapper
	return mapper.CategoryToResponse(&category, langCode), nil
}

func (s *Service) GetCategoryByCode(ctx context.Context, categoryCode string, langCode string) (domain.CategoryResponse, error) {
	category, err := s.Repo.GetCategoryByCode(ctx, categoryCode)
	if err != nil {
		return domain.CategoryResponse{}, err
	}

	// * Convert to CategoryResponse using mapper
	return mapper.CategoryToResponse(&category, langCode), nil
}

func (s *Service) CheckCategoryExists(ctx context.Context, categoryId string) (bool, error) {
	exists, err := s.Repo.CheckCategoryExist(ctx, categoryId)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Service) CheckCategoryCodeExists(ctx context.Context, categoryCode string) (bool, error) {
	exists, err := s.Repo.CheckCategoryCodeExist(ctx, categoryCode)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Service) CountCategories(ctx context.Context, params domain.CategoryParams) (int64, error) {
	count, err := s.Repo.CountCategories(ctx, params)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Service) GetCategoryStatistics(ctx context.Context) (domain.CategoryStatisticsResponse, error) {
	stats, err := s.Repo.GetCategoryStatistics(ctx)
	if err != nil {
		return domain.CategoryStatisticsResponse{}, err
	}

	// Convert to CategoryStatisticsResponse using mapper
	return mapper.CategoryStatisticsToResponse(&stats), nil
}

// *===========================HELPER METHODS===========================*

// sendCategoryUpdatedNotificationToAdmins sends notification for category update to all admin users
func (s *Service) sendCategoryUpdatedNotificationToAdmins(ctx context.Context, category *domain.Category) {
	if s.NotificationService == nil {
		log.Printf("Notification service not available, skipping category updated notification for category ID: %s", category.ID)
		return
	}

	if s.UserRepo == nil {
		log.Printf("User repository not available, skipping category updated notification for category ID: %s", category.ID)
		return
	}

	log.Printf("Sending category updated notification to admins for category ID: %s, category code: %s", category.ID, category.CategoryCode)

	// Get category name in default language
	categoryName := ""
	for _, translation := range category.Translations {
		if translation.LangCode == "en-US" {
			categoryName = translation.CategoryName
			break
		}
	}
	if categoryName == "" && len(category.Translations) > 0 {
		categoryName = category.Translations[0].CategoryName
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
		log.Printf("Failed to get admin users for category updated notification: %v", err)
		return
	}

	if len(admins) == 0 {
		log.Printf("No admin users found, skipping category updated notification for category ID: %s", category.ID)
		return
	}

	// Send notification to each admin
	for _, admin := range admins {
		s.sendCategoryUpdatedNotification(ctx, category.ID, categoryName, admin.ID)
	}

	log.Printf("Successfully sent category updated notification to %d admin(s) for category ID: %s", len(admins), category.ID)
}

// sendCategoryUpdatedNotification sends notification for category update to a specific user
func (s *Service) sendCategoryUpdatedNotification(ctx context.Context, categoryID, categoryName, userID string) {
	titleKey, messageKey, params := messages.CategoryUpdatedNotification(categoryName)
	utilTranslations := messages.GetCategoryNotificationTranslations(titleKey, messageKey, params)

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
		RelatedEntityType: utils.StringPtr("category"),
		RelatedEntityID:   utils.StringPtr(categoryID),
		Type:              domain.NotificationTypeCategoryChange,
		Translations:      translations,
	}

	_, err := s.NotificationService.CreateNotification(ctx, notificationPayload)
	if err != nil {
		log.Printf("Failed to create category updated notification for user ID: %s: %v", userID, err)
	} else {
		log.Printf("Successfully created category updated notification for user ID: %s", userID)
	}
}

// *===========================BACKGROUND TRANSLATION HELPERS===========================*

// autoTranslateCategoryAsync translates missing category translations in background
// Timeout: 30 seconds
func (s *Service) autoTranslateCategoryAsync(categoryID string, userTranslations []domain.CreateCategoryTranslationPayload) {
	// Create context with 30-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("Starting background translation for category ID: %s", categoryID)

	// Get missing language codes
	existingLangCodes := make([]string, len(userTranslations))
	for i, t := range userTranslations {
		existingLangCodes[i] = t.LangCode
	}
	missingLangCodes := utils.GetMissingTranslationLangCodes(existingLangCodes)

	if len(missingLangCodes) == 0 {
		log.Printf("No missing translations for category ID: %s", categoryID)
		return
	}

	// Convert domain types to utils types
	utilsTranslations := make([]utils.CreateTranslationPayload, len(userTranslations))
	for i, t := range userTranslations {
		utilsTranslations[i] = utils.CreateTranslationPayload{
			LangCode:     t.LangCode,
			CategoryName: t.CategoryName,
			Description:  t.Description,
		}
	}

	// Translate using utils helper
	translatedPayloads, err := utils.AutoTranslateCategoryCreate(ctx, s.Translator, utilsTranslations)
	if err != nil {
		log.Printf("Failed to auto-translate category ID %s: %v", categoryID, err)
		return
	}

	// Extract only new translations
	newTranslations := make([]domain.CategoryTranslation, 0)
	for _, translated := range translatedPayloads {
		// Skip user-provided translations
		isUserProvided := false
		for _, userTrans := range userTranslations {
			if userTrans.LangCode == translated.LangCode {
				isUserProvided = true
				break
			}
		}

		if !isUserProvided {
			newTranslations = append(newTranslations, domain.CategoryTranslation{
				LangCode:     translated.LangCode,
				CategoryName: translated.CategoryName,
				Description:  translated.Description,
			})
		}
	}

	// Add translations to database
	if len(newTranslations) > 0 {
		err = s.Repo.AddCategoryTranslations(ctx, categoryID, newTranslations)
		if err != nil {
			log.Printf("Failed to save auto-translated translations for category ID %s: %v", categoryID, err)
		} else {
			log.Printf("Successfully saved %d auto-translated translations for category ID: %s", len(newTranslations), categoryID)
		}
	}
}

// autoTranslateUpdateCategoryAsync translates missing category update translations in background
// Timeout: 30 seconds
func (s *Service) autoTranslateUpdateCategoryAsync(categoryID string, userUpdates []domain.UpdateCategoryTranslationPayload, existingTranslations []domain.CategoryTranslation) {
	// Create context with 30-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("Starting background translation for updated category ID: %s", categoryID)

	// Get missing language codes
	updatedLangCodes := make([]string, len(userUpdates))
	for i, t := range userUpdates {
		updatedLangCodes[i] = t.LangCode
	}

	existingLangCodes := make([]string, len(existingTranslations))
	for i, t := range existingTranslations {
		existingLangCodes[i] = t.LangCode
	}

	// Combine updated and existing
	allLangCodes := make(map[string]bool)
	for _, code := range updatedLangCodes {
		allLangCodes[code] = true
	}
	for _, code := range existingLangCodes {
		allLangCodes[code] = true
	}

	currentCodes := make([]string, 0, len(allLangCodes))
	for code := range allLangCodes {
		currentCodes = append(currentCodes, code)
	}

	missingLangCodes := utils.GetMissingTranslationLangCodes(currentCodes)
	if len(missingLangCodes) == 0 {
		log.Printf("No missing translations for updated category ID: %s", categoryID)
		return
	}

	// Convert domain types to utils types
	utilsUpdates := make([]utils.UpdateTranslationPayload, len(userUpdates))
	for i, t := range userUpdates {
		utilsUpdates[i] = utils.UpdateTranslationPayload{
			LangCode:     t.LangCode,
			CategoryName: t.CategoryName,
			Description:  t.Description,
		}
	}

	utilsExisting := make([]utils.ExistingTranslation, len(existingTranslations))
	for i, t := range existingTranslations {
		utilsExisting[i] = utils.ExistingTranslation{
			LangCode:     t.LangCode,
			CategoryName: t.CategoryName,
			Description:  t.Description,
		}
	}

	// Translate using utils helper
	translatedPayloads, err := utils.AutoTranslateCategoryUpdate(ctx, s.Translator, utilsUpdates, utilsExisting)
	if err != nil {
		log.Printf("Failed to auto-translate updated category ID %s: %v", categoryID, err)
		return
	}

	// Extract only new translations (not in userUpdates)
	newTranslations := make([]domain.CategoryTranslation, 0)
	for _, translated := range translatedPayloads {
		// Skip user-updated translations
		isUserUpdated := false
		for _, userUpdate := range userUpdates {
			if userUpdate.LangCode == translated.LangCode {
				isUserUpdated = true
				break
			}
		}

		if !isUserUpdated && translated.CategoryName != nil {
			newTranslations = append(newTranslations, domain.CategoryTranslation{
				LangCode:     translated.LangCode,
				CategoryName: *translated.CategoryName,
				Description:  translated.Description,
			})
		}
	}

	// Add translations to database
	if len(newTranslations) > 0 {
		err = s.Repo.AddCategoryTranslations(ctx, categoryID, newTranslations)
		if err != nil {
			log.Printf("Failed to save auto-translated translations for updated category ID %s: %v", categoryID, err)
		} else {
			log.Printf("Successfully saved %d auto-translated translations for updated category ID: %s", len(newTranslations), categoryID)
		}
	}
}
