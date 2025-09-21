package category

import (
	"context"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/Rizz404/inventory-api/internal/utils"
)

// * Repository interface defines the contract for category data operations
type Repository interface {
	// * MUTATION
	CreateCategory(ctx context.Context, payload *domain.Category) (domain.Category, error)
	UpdateCategory(ctx context.Context, categoryId string, payload *domain.UpdateCategoryPayload) (domain.Category, error)
	DeleteCategory(ctx context.Context, categoryId string) error

	// * QUERY
	GetCategoriesPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.Category, error)
	GetCategoriesCursor(ctx context.Context, params query.Params, langCode string) ([]domain.Category, error)
	GetCategoryById(ctx context.Context, categoryId string) (domain.Category, error)
	GetCategoryByCode(ctx context.Context, categoryCode string) (domain.Category, error)
	GetCategoryHierarchy(ctx context.Context, langCode string) ([]domain.CategoryResponse, error)
	CheckCategoryExist(ctx context.Context, categoryId string) (bool, error)
	CheckCategoryCodeExist(ctx context.Context, categoryCode string) (bool, error)
	CheckCategoryCodeExistExcluding(ctx context.Context, categoryCode string, excludeCategoryId string) (bool, error)
	CountCategories(ctx context.Context, params query.Params) (int64, error)
	GetCategoryStatistics(ctx context.Context) (domain.CategoryStatistics, error)
}

// * CategoryService interface defines the contract for category business operations
type CategoryService interface {
	// * MUTATION
	CreateCategory(ctx context.Context, payload *domain.CreateCategoryPayload) (domain.CategoryResponse, error)
	UpdateCategory(ctx context.Context, categoryId string, payload *domain.UpdateCategoryPayload) (domain.CategoryResponse, error)
	DeleteCategory(ctx context.Context, categoryId string) error

	// * QUERY
	GetCategoriesPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.CategoryResponse, int64, error)
	GetCategoriesCursor(ctx context.Context, params query.Params, langCode string) ([]domain.CategoryResponse, error)
	GetCategoryById(ctx context.Context, categoryId string, langCode string) (domain.CategoryResponse, error)
	GetCategoryByCode(ctx context.Context, categoryCode string, langCode string) (domain.CategoryResponse, error)
	GetCategoryHierarchy(ctx context.Context, langCode string) ([]domain.CategoryResponse, error)
	CheckCategoryExists(ctx context.Context, categoryId string) (bool, error)
	CheckCategoryCodeExists(ctx context.Context, categoryCode string) (bool, error)
	CountCategories(ctx context.Context, params query.Params) (int64, error)
	GetCategoryStatistics(ctx context.Context) (domain.CategoryStatisticsResponse, error)
}

type Service struct {
	Repo Repository
}

// * Ensure Service implements CategoryService interface
var _ CategoryService = (*Service)(nil)

func NewService(r Repository) CategoryService {
	return &Service{
		Repo: r,
	}
}

// *===========================MUTATION===========================*
func (s *Service) CreateCategory(ctx context.Context, payload *domain.CreateCategoryPayload) (domain.CategoryResponse, error) {
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

	// * Prepare domain category
	newCategory := domain.Category{
		ParentID:     payload.ParentID,
		CategoryCode: payload.CategoryCode,
		Translations: make([]domain.CategoryTranslation, len(payload.Translations)),
	}

	// * Convert translation payloads to domain translations
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

	// * Convert to CategoryResponse using mapper
	return mapper.CategoryToResponse(&createdCategory, mapper.DefaultLangCode), nil
}

func (s *Service) UpdateCategory(ctx context.Context, categoryId string, payload *domain.UpdateCategoryPayload) (domain.CategoryResponse, error) {
	// * Check if category exists
	_, err := s.Repo.GetCategoryById(ctx, categoryId)
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

	updatedCategory, err := s.Repo.UpdateCategory(ctx, categoryId, payload)
	if err != nil {
		return domain.CategoryResponse{}, err
	}

	// * Convert to CategoryResponse using mapper
	return mapper.CategoryToResponse(&updatedCategory, mapper.DefaultLangCode), nil
}

func (s *Service) DeleteCategory(ctx context.Context, categoryId string) error {
	err := s.Repo.DeleteCategory(ctx, categoryId)
	if err != nil {
		return err
	}
	return nil
}

// *===========================QUERY===========================*
func (s *Service) GetCategoriesPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.CategoryResponse, int64, error) {
	categories, err := s.Repo.GetCategoriesPaginated(ctx, params, langCode)
	if err != nil {
		return nil, 0, err
	}

	// * Count total for pagination
	count, err := s.Repo.CountCategories(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// * Convert to CategoryResponse using mapper
	categoryResponses := mapper.CategoriesToResponses(categories, langCode)

	return categoryResponses, count, nil
}

func (s *Service) GetCategoriesCursor(ctx context.Context, params query.Params, langCode string) ([]domain.CategoryResponse, error) {
	categories, err := s.Repo.GetCategoriesCursor(ctx, params, langCode)
	if err != nil {
		return nil, err
	}

	// * Convert to CategoryResponse using mapper
	categoryResponses := mapper.CategoriesToResponses(categories, langCode)

	return categoryResponses, nil
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

func (s *Service) GetCategoryHierarchy(ctx context.Context, langCode string) ([]domain.CategoryResponse, error) {
	hierarchy, err := s.Repo.GetCategoryHierarchy(ctx, langCode)
	if err != nil {
		return nil, err
	}
	return hierarchy, nil
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

func (s *Service) CountCategories(ctx context.Context, params query.Params) (int64, error) {
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
