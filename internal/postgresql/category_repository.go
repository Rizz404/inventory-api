package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type gormCategoryRepository struct {
	db *gorm.DB
}

type CategoryFilterOptions struct {
	ParentID  *string `json:"parentId,omitempty"`
	HasParent *bool   `json:"hasParent,omitempty"`
}

func NewCategoryRepository(db *gorm.DB) *gormCategoryRepository {
	return &gormCategoryRepository{
		db: db,
	}
}

func (r *gormCategoryRepository) applyCategoryFilters(db *gorm.DB, filters any) *gorm.DB {
	f, ok := filters.(*CategoryFilterOptions)
	if !ok || f == nil {
		return db
	}

	if f.ParentID != nil {
		db = db.Where("categories.parent_id = ?", f.ParentID)
	}
	if f.HasParent != nil {
		if *f.HasParent {
			db = db.Where("categories.parent_id IS NOT NULL")
		} else {
			db = db.Where("categories.parent_id IS NULL")
		}
	}
	return db
}

func (r *gormCategoryRepository) applyCategorySorts(db *gorm.DB, sort *query.SortOptions) *gorm.DB {
	if sort == nil || sort.Field == "" {
		return db.Order("categories.created_at DESC")
	}
	var orderClause string
	switch strings.ToLower(sort.Field) {
	case "category_code", "created_at", "updated_at":
		orderClause = "categories." + sort.Field
	case "name", "category_name":
		orderClause = "category_translations.category_name"
	default:
		return db.Order("categories.created_at DESC")
	}

	order := "DESC"
	if strings.ToLower(sort.Order) == "asc" {
		order = "ASC"
	}
	return db.Order(fmt.Sprintf("%s %s", orderClause, order))
}

// *===========================MUTATION===========================*
func (r *gormCategoryRepository) CreateCategory(ctx context.Context, payload *domain.Category) (domain.Category, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.Category{}, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create category
	modelCategory := mapper.ToModelCategoryForCreate(payload)
	if err := tx.Create(&modelCategory).Error; err != nil {
		tx.Rollback()
		return domain.Category{}, domain.ErrInternal(err)
	}

	// Create translations
	for _, translation := range payload.Translations {
		modelTranslation := mapper.ToModelCategoryTranslationForCreate(modelCategory.ID.String(), &translation)
		if err := tx.Create(&modelTranslation).Error; err != nil {
			tx.Rollback()
			return domain.Category{}, domain.ErrInternal(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return domain.Category{}, domain.ErrInternal(err)
	}

	// Fetch created category with translations
	return r.GetCategoryById(ctx, modelCategory.ID.String())
}

func (r *gormCategoryRepository) UpdateCategory(ctx context.Context, categoryId string, payload *domain.UpdateCategoryPayload) (domain.Category, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.Category{}, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update category basic info
	updates := mapper.ToModelCategoryUpdateMap(payload)
	if len(updates) > 0 {
		if err := tx.Model(&model.Category{}).Where("id = ?", categoryId).Updates(updates).Error; err != nil {
			tx.Rollback()
			return domain.Category{}, domain.ErrInternal(err)
		}
	}

	// Update translations if provided
	if len(payload.Translations) > 0 {
		for _, translationPayload := range payload.Translations {
			translationUpdates := mapper.ToModelCategoryTranslationUpdateMap(&translationPayload)
			if len(translationUpdates) > 0 {
				// Try to update existing translation
				result := tx.Model(&model.CategoryTranslation{}).
					Where("category_id = ? AND lang_code = ?", categoryId, translationPayload.LangCode).
					Updates(translationUpdates)

				if result.Error != nil {
					tx.Rollback()
					return domain.Category{}, domain.ErrInternal(result.Error)
				}

				// If no rows affected, create new translation
				if result.RowsAffected == 0 {
					newTranslation := model.CategoryTranslation{
						LangCode:     translationPayload.LangCode,
						CategoryName: *translationPayload.CategoryName,
						Description:  translationPayload.Description,
					}
					if parsedCategoryID, err := ulid.Parse(categoryId); err == nil {
						newTranslation.CategoryID = model.SQLULID(parsedCategoryID)
					}

					if err := tx.Create(&newTranslation).Error; err != nil {
						tx.Rollback()
						return domain.Category{}, domain.ErrInternal(err)
					}
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return domain.Category{}, domain.ErrInternal(err)
	}

	// Fetch updated category with translations
	return r.GetCategoryById(ctx, categoryId)
}

func (r *gormCategoryRepository) DeleteCategory(ctx context.Context, categoryId string) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.ErrInternal(tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete translations first (foreign key constraint)
	if err := tx.Delete(&model.CategoryTranslation{}, "category_id = ?", categoryId).Error; err != nil {
		tx.Rollback()
		return domain.ErrInternal(err)
	}

	// Delete category
	if err := tx.Delete(&model.Category{}, "id = ?", categoryId).Error; err != nil {
		tx.Rollback()
		return domain.ErrInternal(err)
	}

	if err := tx.Commit().Error; err != nil {
		return domain.ErrInternal(err)
	}

	return nil
}

// *===========================QUERY===========================*
func (r *gormCategoryRepository) GetCategoriesPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.Category, error) {
	var categories []model.Category
	db := r.db.WithContext(ctx).Preload("Translations")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN category_translations ON categories.id = category_translations.category_id").
			Where("categories.category_code ILIKE ? OR category_translations.category_name ILIKE ?", searchPattern, searchPattern).
			Distinct("categories.id")
	}

	// Set pagination cursor to empty for offset-based pagination
	params.Pagination.Cursor = ""
	db = query.Apply(db, params, r.applyCategoryFilters, r.applyCategorySorts)

	if err := db.Find(&categories).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain categories
	domainCategories := make([]domain.Category, len(categories))
	for i, category := range categories {
		domainCategories[i] = mapper.ToDomainCategory(&category)
	}
	return domainCategories, nil
}

func (r *gormCategoryRepository) GetCategoriesCursor(ctx context.Context, params query.Params, langCode string) ([]domain.Category, error) {
	var categories []model.Category
	db := r.db.WithContext(ctx).Preload("Translations")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN category_translations ON categories.id = category_translations.category_id").
			Where("categories.category_code ILIKE ? OR category_translations.category_name ILIKE ?", searchPattern, searchPattern).
			Distinct("categories.id")
	}

	// Set offset to 0 for cursor-based pagination
	params.Pagination.Offset = 0
	db = query.Apply(db, params, r.applyCategoryFilters, r.applyCategorySorts)

	if err := db.Find(&categories).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain categories
	domainCategories := make([]domain.Category, len(categories))
	for i, category := range categories {
		domainCategories[i] = mapper.ToDomainCategory(&category)
	}
	return domainCategories, nil
}

func (r *gormCategoryRepository) GetCategoriesResponse(ctx context.Context, params query.Params, langCode string) ([]domain.CategoryResponse, error) {
	var categories []model.Category
	db := r.db.WithContext(ctx).Preload("Translations")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN category_translations ON categories.id = category_translations.category_id").
			Where("categories.category_code ILIKE ? OR category_translations.category_name ILIKE ?", searchPattern, searchPattern).
			Distinct("categories.id")
	}

	db = query.Apply(db, params, r.applyCategoryFilters, r.applyCategorySorts)

	if err := db.Find(&categories).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	return mapper.ToDomainCategoriesResponse(categories, langCode), nil
}

func (r *gormCategoryRepository) GetCategoryById(ctx context.Context, categoryId string) (domain.Category, error) {
	var category model.Category

	err := r.db.WithContext(ctx).Preload("Translations").First(&category, "id = ?", categoryId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Category{}, domain.ErrNotFound("category")
		}
		return domain.Category{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainCategory(&category), nil
}

func (r *gormCategoryRepository) GetCategoryByCode(ctx context.Context, categoryCode string) (domain.Category, error) {
	var category model.Category

	err := r.db.WithContext(ctx).Preload("Translations").First(&category, "category_code = ?", categoryCode).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Category{}, domain.ErrNotFound("category")
		}
		return domain.Category{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainCategory(&category), nil
}

func (r *gormCategoryRepository) GetCategoryHierarchy(ctx context.Context, langCode string) ([]domain.CategoryResponse, error) {
	categories, err := r.GetCategoriesResponse(ctx, query.Params{
		Pagination: &query.PaginationOptions{
			Limit: 1000, // Large limit to get all categories
		},
	}, langCode)
	if err != nil {
		return nil, err
	}

	return mapper.BuildCategoryHierarchy(categories), nil
}

func (r *gormCategoryRepository) CheckCategoryExist(ctx context.Context, categoryId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Category{}).Where("id = ?", categoryId).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *gormCategoryRepository) CheckCategoryCodeExist(ctx context.Context, categoryCode string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Category{}).Where("category_code = ?", categoryCode).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *gormCategoryRepository) CountCategories(ctx context.Context, params query.Params) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(&model.Category{})

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN category_translations ON categories.id = category_translations.category_id").
			Where("categories.category_code ILIKE ? OR category_translations.category_name ILIKE ?", searchPattern, searchPattern).
			Distinct("categories.id")
	}

	db = query.Apply(db, query.Params{Filters: params.Filters}, r.applyCategoryFilters, nil)

	if err := db.Count(&count).Error; err != nil {
		return 0, domain.ErrInternal(err)
	}
	return count, nil
}
