package postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) applyCategoryFilters(db *gorm.DB, filters *domain.CategoryFilterOptions) *gorm.DB {
	if filters == nil {
		return db
	}

	if filters.ParentID != nil {
		db = db.Where("c.parent_id = ?", filters.ParentID)
	}
	if filters.HasParent != nil {
		if *filters.HasParent {
			db = db.Where("c.parent_id IS NOT NULL")
		} else {
			db = db.Where("c.parent_id IS NULL")
		}
	}
	return db
}

func (r *CategoryRepository) applyCategorySorts(db *gorm.DB, sort *domain.CategorySortOptions, langCode string) *gorm.DB {
	if sort == nil || sort.Field == "" {
		return db.Order("c.created_at DESC")
	}

	order := "DESC"
	if sort.Order == domain.SortOrderAsc {
		order = "ASC"
	}

	if sort.Field == domain.CategorySortByCategoryName {
		// Use subquery for sorting by translation
		subquery := fmt.Sprintf("(SELECT category_name FROM category_translations WHERE category_id = c.id AND lang_code = '%s' LIMIT 1)", langCode)
		return db.Order(fmt.Sprintf("%s %s", subquery, order))
	}

	// Map camelCase sort field to snake_case database column
	columnName := mapper.MapCategorySortFieldToColumn(sort.Field)
	orderClause := columnName
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

	// Return created category with translations (no need to query again)
	// GORM has already filled the model with created data including ID and timestamps
	domainCategory := mapper.ToDomainCategory(&modelCategory)
	// Add translations manually since they were created separately
	for i, translation := range payload.Translations {
		domainCategory.Translations = append(domainCategory.Translations, domain.CategoryTranslation{
			LangCode:     translation.LangCode,
			CategoryName: translation.CategoryName,
			Description:  translation.Description,
		})
		_ = i // avoid unused variable
	}
	return domainCategory, nil
}

func (r *CategoryRepository) BulkCreateCategories(ctx context.Context, categories []domain.Category) ([]domain.Category, error) {
	if len(categories) == 0 {
		return []domain.Category{}, nil
	}

	models := make([]*model.Category, len(categories))
	for i := range categories {
		m := mapper.ToModelCategoryForCreate(&categories[i])
		models[i] = &m
	}

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.
		Omit(clause.Associations).
		Session(&gorm.Session{CreateBatchSize: 500}).
		Create(&models).Error; err != nil {
		tx.Rollback()
		return nil, domain.ErrInternal(err)
	}

	var translations []model.CategoryTranslation
	for i := range models {
		c := categories[i]
		for _, t := range c.Translations {
			mt := mapper.ToModelCategoryTranslationForCreate(models[i].ID.String(), &t)
			translations = append(translations, mt)
		}
	}
	if len(translations) > 0 {
		if err := tx.Session(&gorm.Session{CreateBatchSize: 500}).Create(&translations).Error; err != nil {
			tx.Rollback()
			return nil, domain.ErrInternal(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	created := make([]domain.Category, len(models))
	for i := range models {
		created[i] = mapper.ToDomainCategory(models[i])
		for _, t := range categories[i].Translations {
			created[i].Translations = append(created[i].Translations, domain.CategoryTranslation{
				LangCode:     t.LangCode,
				CategoryName: t.CategoryName,
				Description:  t.Description,
			})
		}
	}
	return created, nil
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
	var updatedCategory model.Category
	err := r.db.WithContext(ctx).
		Table("categories c").
		Preload("Translations").
		First(&updatedCategory, "id = ?", categoryId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Category{}, domain.ErrNotFound("category")
		}
		return domain.Category{}, domain.ErrInternal(err)
	}
	return mapper.ToDomainCategory(&updatedCategory), nil
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

func (r *CategoryRepository) BulkDeleteCategories(ctx context.Context, categoryIds []string) (domain.BulkDeleteCategories, error) {
	result := domain.BulkDeleteCategories{
		RequestedIDS: categoryIds,
		DeletedIDS:   []string{},
	}

	if len(categoryIds) == 0 {
		return result, nil
	}

	// First, find which categories actually exist
	var existingCategories []model.Category
	if err := r.db.WithContext(ctx).Select("id").Where("id IN ?", categoryIds).Find(&existingCategories).Error; err != nil {
		return result, domain.ErrInternal(err)
	}

	// Collect existing category IDs
	existingIds := make([]string, 0, len(existingCategories))
	for _, cat := range existingCategories {
		existingIds = append(existingIds, cat.ID.String())
	}

	// If no categories exist, return early
	if len(existingIds) == 0 {
		return result, nil
	}

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return result, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete translations first (foreign key constraint)
	if err := tx.Delete(&model.CategoryTranslation{}, "category_id IN ?", existingIds).Error; err != nil {
		tx.Rollback()
		return result, domain.ErrInternal(err)
	}

	// Delete categories
	if err := tx.Delete(&model.Category{}, "id IN ?", existingIds).Error; err != nil {
		tx.Rollback()
		return result, domain.ErrInternal(err)
	}

	if err := tx.Commit().Error; err != nil {
		return result, domain.ErrInternal(err)
	}

	result.DeletedIDS = existingIds
	return result, nil
}

// *===========================QUERY===========================*
func (r *CategoryRepository) GetCategoriesPaginated(ctx context.Context, params domain.CategoryParams, langCode string) ([]domain.Category, error) {
	var categories []model.Category
	db := r.db.WithContext(ctx).
		Table("categories c").
		Preload("Translations").
		Preload("Parent").
		Preload("Parent.Translations")

	needsJoin := params.SearchQuery != nil && *params.SearchQuery != "" ||
		(params.Sort != nil && params.Sort.Field == domain.CategorySortByCategoryName)

	if needsJoin {
		db = db.Select("c.id, c.category_code, c.created_at, c.updated_at").
			Joins("LEFT JOIN category_translations ct ON c.id = ct.category_id")
		if params.SearchQuery != nil && *params.SearchQuery != "" {
			searchPattern := "%" + *params.SearchQuery + "%"
			db = db.Where("c.category_code ILIKE ? OR ct.category_name ILIKE ?", searchPattern, searchPattern).
				Distinct("c.id, c.created_at")
		}
	}

	// Apply filters
	db = r.applyCategoryFilters(db, params.Filters)

	// Apply sorting
	db = r.applyCategorySorts(db, params.Sort, langCode)

	// Apply pagination
	if params.Pagination != nil {
		if params.Pagination.Limit > 0 {
			db = db.Limit(params.Pagination.Limit)
		}
		if params.Pagination.Offset > 0 {
			db = db.Offset(params.Pagination.Offset)
		}
	}

	if err := db.Find(&categories).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain categories
	return mapper.ToDomainCategories(categories), nil
}

func (r *CategoryRepository) GetCategoriesCursor(ctx context.Context, params domain.CategoryParams, langCode string) ([]domain.Category, error) {
	var categories []model.Category
	db := r.db.WithContext(ctx).
		Table("categories c").
		Preload("Translations").
		Preload("Parent").
		Preload("Parent.Translations")

	needsJoin := params.SearchQuery != nil && *params.SearchQuery != "" ||
		(params.Sort != nil && params.Sort.Field == domain.CategorySortByCategoryName)

	if needsJoin {
		db = db.Select("c.id, c.category_code, c.created_at, c.updated_at").
			Joins("LEFT JOIN category_translations ct ON c.id = ct.category_id")
		if params.SearchQuery != nil && *params.SearchQuery != "" {
			searchPattern := "%" + *params.SearchQuery + "%"
			db = db.Where("c.category_code ILIKE ? OR ct.category_name ILIKE ?", searchPattern, searchPattern).
				Distinct("c.id, c.created_at")
		}
	}

	// Apply filters
	db = r.applyCategoryFilters(db, params.Filters)

	// Apply sorting - for cursor pagination, we need consistent ordering by ID
	if params.Sort != nil && params.Sort.Field != "" {
		db = r.applyCategorySorts(db, params.Sort, langCode)
		// Always add secondary sort by ID DESC for consistency (ULID = newer = larger)
		db = db.Order("c.id DESC")
	} else {
		// Default to ID DESC for cursor pagination (newest first)
		db = db.Order("c.id DESC")
	}

	// Apply cursor-based pagination
	if params.Pagination != nil {
		if params.Pagination.Cursor != "" {
			db = db.Where("c.id < ?", params.Pagination.Cursor)
		}
		if params.Pagination.Limit > 0 {
			db = db.Limit(params.Pagination.Limit)
		}
	}

	if err := db.Find(&categories).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain categories
	return mapper.ToDomainCategories(categories), nil
}

func (r *CategoryRepository) GetCategoryById(ctx context.Context, categoryId string) (domain.Category, error) {
	var category model.Category

	err := r.db.WithContext(ctx).
		Table("categories c").
		Preload("Translations").
		Preload("Parent").
		Preload("Parent.Translations").
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
		Preload("Parent").
		Preload("Parent.Translations").
		First(&category, "category_code = ?", categoryCode).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Category{}, domain.ErrNotFound("category")
		}
		return domain.Category{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainCategory(&category), nil
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

func (r *CategoryRepository) CountCategories(ctx context.Context, params domain.CategoryParams) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Table("categories c")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN category_translations ct ON c.id = ct.category_id").
			Where("c.category_code ILIKE ? OR ct.category_name ILIKE ?", searchPattern, searchPattern).
			Distinct("c.id")
	}

	// Apply filters
	db = r.applyCategoryFilters(db, params.Filters)

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
		Date  time.Time `json:"date"`
		Count int64     `json:"count"`
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

	stats.Summary.EarliestCreationDate = earliestDate
	stats.Summary.LatestCreationDate = latestDate

	// Calculate average categories per day
	if !earliestDate.IsZero() && !latestDate.IsZero() {
		daysDiff := latestDate.Sub(earliestDate).Hours() / 24
		if daysDiff > 0 {
			stats.Summary.AverageCategoriesPerDay = float64(totalCount) / daysDiff
		}
	}

	return stats, nil
}
