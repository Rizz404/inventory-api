package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

type CategoryFilterOptions struct {
	ParentID  *string `json:"parentId,omitempty"`
	HasParent *bool   `json:"hasParent,omitempty"`
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) applyCategoryFilters(db *gorm.DB, filters any) *gorm.DB {
	f, ok := filters.(*CategoryFilterOptions)
	if !ok || f == nil {
		return db
	}

	if f.ParentID != nil {
		db = db.Where("c.parent_id = ?", f.ParentID)
	}
	if f.HasParent != nil {
		if *f.HasParent {
			db = db.Where("c.parent_id IS NOT NULL")
		} else {
			db = db.Where("c.parent_id IS NULL")
		}
	}
	return db
}

func (r *CategoryRepository) applyCategorySorts(db *gorm.DB, sort *query.SortOptions) *gorm.DB {
	if sort == nil || sort.Field == "" {
		return db.Order("c.created_at DESC")
	}
	var orderClause string
	switch strings.ToLower(sort.Field) {
	case "category_code", "created_at", "updated_at":
		orderClause = "c." + sort.Field
	case "name", "category_name":
		orderClause = "ct.category_name"
	default:
		return db.Order("c.created_at DESC")
	}

	order := "DESC"
	if strings.ToLower(sort.Order) == "asc" {
		order = "ASC"
	}
	return db.Order(fmt.Sprintf("%s %s", orderClause, order))
}

// *===========================MUTATION===========================*
func (r *CategoryRepository) CreateCategory(ctx context.Context, payload *domain.Category) (domain.Category, error) {
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

func (r *CategoryRepository) UpdateCategory(ctx context.Context, categoryId string, payload *domain.UpdateCategoryPayload) (domain.Category, error) {
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

func (r *CategoryRepository) DeleteCategory(ctx context.Context, categoryId string) error {
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
func (r *CategoryRepository) GetCategoriesPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.Category, error) {
	var categories []model.Category
	db := r.db.WithContext(ctx).
		Table("categories c").
		Preload("Translations")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN category_translations ct ON c.id = ct.category_id").
			Where("c.category_code ILIKE ? OR ct.category_name ILIKE ?", searchPattern, searchPattern).
			Distinct("c.id")
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

func (r *CategoryRepository) GetCategoriesCursor(ctx context.Context, params query.Params, langCode string) ([]domain.Category, error) {
	var categories []model.Category
	db := r.db.WithContext(ctx).
		Table("categories c").
		Preload("Translations")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN category_translations ct ON c.id = ct.category_id").
			Where("c.category_code ILIKE ? OR ct.category_name ILIKE ?", searchPattern, searchPattern).
			Distinct("c.id")
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

func (r *CategoryRepository) GetCategoriesResponse(ctx context.Context, params query.Params, langCode string) ([]domain.CategoryResponse, error) {
	var categories []model.Category
	db := r.db.WithContext(ctx).
		Table("categories c").
		Preload("Translations")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN category_translations ct ON c.id = ct.category_id").
			Where("c.category_code ILIKE ? OR ct.category_name ILIKE ?", searchPattern, searchPattern).
			Distinct("c.id")
	}

	db = query.Apply(db, params, r.applyCategoryFilters, r.applyCategorySorts)

	if err := db.Find(&categories).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	return mapper.ToDomainCategoriesResponse(categories, langCode), nil
}

func (r *CategoryRepository) GetCategoryById(ctx context.Context, categoryId string) (domain.Category, error) {
	var category model.Category

	err := r.db.WithContext(ctx).
		Table("categories c").
		Preload("Translations").
		First(&category, "id = ?", categoryId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Category{}, domain.ErrNotFound("category")
		}
		return domain.Category{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainCategory(&category), nil
}

func (r *CategoryRepository) GetCategoryByCode(ctx context.Context, categoryCode string) (domain.Category, error) {
	var category model.Category

	err := r.db.WithContext(ctx).
		Table("categories c").
		Preload("Translations").
		First(&category, "category_code = ?", categoryCode).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Category{}, domain.ErrNotFound("category")
		}
		return domain.Category{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainCategory(&category), nil
}

func (r *CategoryRepository) GetCategoryHierarchy(ctx context.Context, langCode string) ([]domain.CategoryResponse, error) {
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

func (r *CategoryRepository) CheckCategoryExist(ctx context.Context, categoryId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Category{}).Where("id = ?", categoryId).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *CategoryRepository) CheckCategoryCodeExist(ctx context.Context, categoryCode string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Category{}).Where("category_code = ?", categoryCode).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *CategoryRepository) CheckCategoryCodeExistExcluding(ctx context.Context, categoryCode string, excludeCategoryId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Category{}).Where("category_code = ? AND id != ?", categoryCode, excludeCategoryId).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *CategoryRepository) CountCategories(ctx context.Context, params query.Params) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Table("categories c")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN category_translations ct ON c.id = ct.category_id").
			Where("c.category_code ILIKE ? OR ct.category_name ILIKE ?", searchPattern, searchPattern).
			Distinct("c.id")
	}

	db = query.Apply(db, query.Params{Filters: params.Filters}, r.applyCategoryFilters, nil)

	if err := db.Count(&count).Error; err != nil {
		return 0, domain.ErrInternal(err)
	}
	return count, nil
}

func (r *CategoryRepository) GetCategoryStatistics(ctx context.Context) (domain.CategoryStatistics, error) {
	var stats domain.CategoryStatistics

	// Get total category count
	var totalCount int64
	if err := r.db.WithContext(ctx).Model(&model.Category{}).Count(&totalCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.Total.Count = int(totalCount)

	// Get hierarchy statistics
	var topLevelCount, withChildrenCount, withParentCount int64

	// Top level categories (no parent)
	if err := r.db.WithContext(ctx).Model(&model.Category{}).Where("parent_id IS NULL").Count(&topLevelCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	// Categories with children
	if err := r.db.WithContext(ctx).Model(&model.Category{}).
		Where("id IN (SELECT DISTINCT parent_id FROM categories WHERE parent_id IS NOT NULL)").
		Count(&withChildrenCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	// Categories with parent
	if err := r.db.WithContext(ctx).Model(&model.Category{}).Where("parent_id IS NOT NULL").Count(&withParentCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.ByHierarchy.TopLevel = int(topLevelCount)
	stats.ByHierarchy.WithChildren = int(withChildrenCount)
	stats.ByHierarchy.WithParent = int(withParentCount)

	// Get creation trends (last 30 days)
	var creationTrends []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	if err := r.db.WithContext(ctx).Model(&model.Category{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= NOW() - INTERVAL '30 days'").
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&creationTrends).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.CreationTrends = make([]domain.CategoryCreationTrend, len(creationTrends))
	for i, ct := range creationTrends {
		stats.CreationTrends[i] = domain.CategoryCreationTrend{
			Date:  ct.Date,
			Count: int(ct.Count),
		}
	}

	// Calculate summary statistics
	stats.Summary.TotalCategories = int(totalCount)

	if totalCount > 0 {
		stats.Summary.TopLevelPercentage = float64(topLevelCount) / float64(totalCount) * 100
		stats.Summary.SubCategoriesPercentage = float64(withParentCount) / float64(totalCount) * 100
	}

	stats.Summary.CategoriesWithChildrenCount = int(withChildrenCount)
	stats.Summary.CategoriesWithoutChildrenCount = int(totalCount - withChildrenCount)

	// Calculate max depth level using recursive CTE
	var maxDepth int
	if err := r.db.WithContext(ctx).Raw(`
		WITH RECURSIVE category_depth AS (
			SELECT id, parent_id, 1 as depth
			FROM categories
			WHERE parent_id IS NULL

			UNION ALL

			SELECT c.id, c.parent_id, cd.depth + 1
			FROM categories c
			INNER JOIN category_depth cd ON c.parent_id = cd.id
		)
		SELECT COALESCE(MAX(depth), 0) FROM category_depth
	`).Scan(&maxDepth).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.Summary.MaxDepthLevel = maxDepth

	// Get earliest and latest creation dates
	var earliestDate, latestDate time.Time
	if err := r.db.WithContext(ctx).Model(&model.Category{}).Select("MIN(created_at)").Scan(&earliestDate).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.Category{}).Select("MAX(created_at)").Scan(&latestDate).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.Summary.EarliestCreationDate = earliestDate.Format("2006-01-02")
	stats.Summary.LatestCreationDate = latestDate.Format("2006-01-02")

	// Calculate average categories per day
	if !earliestDate.IsZero() && !latestDate.IsZero() {
		daysDiff := latestDate.Sub(earliestDate).Hours() / 24
		if daysDiff > 0 {
			stats.Summary.AverageCategoriesPerDay = float64(totalCount) / daysDiff
		}
	}

	return stats, nil
}
