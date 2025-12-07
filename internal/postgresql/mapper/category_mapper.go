package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/oklog/ulid/v2"
)

// *==================== Model conversions ====================
func ToModelCategory(d *domain.Category) model.Category {
	modelCategory := model.Category{
		CategoryCode: d.CategoryCode,
	}

	if d.ID != "" {
		if parsedID, err := ulid.Parse(d.ID); err == nil {
			modelCategory.ID = model.SQLULID(parsedID)
		}
	}

	if d.ParentID != nil && *d.ParentID != "" {
		if parsedID, err := ulid.Parse(*d.ParentID); err == nil {
			modelULID := model.SQLULID(parsedID)
			modelCategory.ParentID = &modelULID
		}
	}

	return modelCategory
}

func ToModelCategoryForCreate(d *domain.Category) model.Category {
	modelCategory := model.Category{
		CategoryCode: d.CategoryCode,
	}

	if d.ParentID != nil && *d.ParentID != "" {
		if parsedID, err := ulid.Parse(*d.ParentID); err == nil {
			modelULID := model.SQLULID(parsedID)
			modelCategory.ParentID = &modelULID
		}
	}

	return modelCategory
}

func ToModelCategoryTranslation(d *domain.CategoryTranslation) model.CategoryTranslation {
	modelTranslation := model.CategoryTranslation{
		LangCode:     d.LangCode,
		CategoryName: d.CategoryName,
		Description:  d.Description,
	}

	if d.ID != "" {
		if parsedID, err := ulid.Parse(d.ID); err == nil {
			modelTranslation.ID = model.SQLULID(parsedID)
		}
	}

	if d.CategoryID != "" {
		if parsedCategoryID, err := ulid.Parse(d.CategoryID); err == nil {
			modelTranslation.CategoryID = model.SQLULID(parsedCategoryID)
		}
	}

	return modelTranslation
}

func ToModelCategoryTranslationForCreate(categoryID string, d *domain.CategoryTranslation) model.CategoryTranslation {
	modelTranslation := model.CategoryTranslation{
		LangCode:     d.LangCode,
		CategoryName: d.CategoryName,
		Description:  d.Description,
	}

	if categoryID != "" {
		if parsedCategoryID, err := ulid.Parse(categoryID); err == nil {
			modelTranslation.CategoryID = model.SQLULID(parsedCategoryID)
		}
	}

	return modelTranslation
}

// *==================== Entity conversions ====================
func ToDomainCategory(m *model.Category) domain.Category {
	domainCategory := domain.Category{
		ID:           m.ID.String(),
		CategoryCode: m.CategoryCode,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}

	if m.ParentID != nil && !m.ParentID.IsZero() {
		parentIDStr := m.ParentID.String()
		domainCategory.ParentID = &parentIDStr
	}

	if len(m.Translations) > 0 {
		domainCategory.Translations = make([]domain.CategoryTranslation, len(m.Translations))
		for i, translation := range m.Translations {
			domainCategory.Translations[i] = ToDomainCategoryTranslation(&translation)
		}
	}

	if m.Parent != nil {
		parentDomain := ToDomainCategory(m.Parent)
		domainCategory.Parent = &parentDomain
	}

	return domainCategory
}

func ToDomainCategories(models []model.Category) []domain.Category {
	if len(models) == 0 {
		return []domain.Category{}
	}
	categories := make([]domain.Category, len(models))
	for i, m := range models {
		categories[i] = ToDomainCategory(&m)
	}
	return categories
}

func ToDomainCategoryTranslation(m *model.CategoryTranslation) domain.CategoryTranslation {
	return domain.CategoryTranslation{
		ID:           m.ID.String(),
		CategoryID:   m.CategoryID.String(),
		LangCode:     m.LangCode,
		CategoryName: m.CategoryName,
		Description:  m.Description,
	}
}

// *==================== Entity Response conversions ====================
func CategoryToResponse(d *domain.Category, langCode string) domain.CategoryResponse {
	response := domain.CategoryResponse{
		ID:           d.ID,
		ParentID:     d.ParentID,
		CategoryCode: d.CategoryCode,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
		Translations: make([]domain.CategoryTranslationResponse, len(d.Translations)),
	}

	// Populate translations
	for i, translation := range d.Translations {
		response.Translations[i] = domain.CategoryTranslationResponse{
			LangCode:     translation.LangCode,
			CategoryName: translation.CategoryName,
			Description:  translation.Description,
		}
	}

	// Find translation for the requested language
	for _, translation := range d.Translations {
		if translation.LangCode == langCode {
			response.CategoryName = translation.CategoryName
			response.Description = translation.Description
			break
		}
	}

	// If no translation found for requested language, use first available
	if response.CategoryName == "" && len(d.Translations) > 0 {
		response.CategoryName = d.Translations[0].CategoryName
		response.Description = d.Translations[0].Description
	}

	// Populate parent if exists
	if d.Parent != nil {
		parentResponse := CategoryToResponse(d.Parent, langCode)
		response.Parent = &parentResponse
	}

	return response
}

func CategoriesToResponses(categories []domain.Category, langCode string) []domain.CategoryResponse {
	if len(categories) == 0 {
		return []domain.CategoryResponse{}
	}
	responses := make([]domain.CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = CategoryToResponse(&category, langCode)
	}
	return responses
}

func CategoryToListResponse(d *domain.Category, langCode string) domain.CategoryListResponse {
	response := domain.CategoryListResponse{
		ID:           d.ID,
		ParentID:     d.ParentID,
		CategoryCode: d.CategoryCode,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}

	// Find translation for the requested language
	for _, translation := range d.Translations {
		if translation.LangCode == langCode {
			response.CategoryName = translation.CategoryName
			response.Description = translation.Description
			break
		}
	}

	// If no translation found for requested language, use first available
	if response.CategoryName == "" && len(d.Translations) > 0 {
		response.CategoryName = d.Translations[0].CategoryName
		response.Description = d.Translations[0].Description
	}

	// Populate parent if exists
	if d.Parent != nil {
		parentResponse := CategoryToListResponse(d.Parent, langCode)
		response.Parent = &parentResponse
	}

	return response
}

func CategoriesToListResponses(categories []domain.Category, langCode string) []domain.CategoryListResponse {
	if len(categories) == 0 {
		return []domain.CategoryListResponse{}
	}
	responses := make([]domain.CategoryListResponse, len(categories))
	for i, category := range categories {
		responses[i] = CategoryToListResponse(&category, langCode)
	}
	return responses
}

func CategoryStatisticsToResponse(stats *domain.CategoryStatistics) domain.CategoryStatisticsResponse {
	trends := make([]domain.CategoryCreationTrendResponse, len(stats.CreationTrends))
	for i, trend := range stats.CreationTrends {
		trends[i] = domain.CategoryCreationTrendResponse{
			Date:  trend.Date,
			Count: trend.Count,
		}
	}

	return domain.CategoryStatisticsResponse{
		Total: domain.CategoryCountStatisticsResponse{
			Count: stats.Total.Count,
		},
		ByHierarchy: domain.CategoryHierarchyStatisticsResponse{
			TopLevel:     stats.ByHierarchy.TopLevel,
			WithChildren: stats.ByHierarchy.WithChildren,
			WithParent:   stats.ByHierarchy.WithParent,
		},
		CreationTrends: trends,
		Summary: domain.CategorySummaryStatisticsResponse{
			TotalCategories:                stats.Summary.TotalCategories,
			TopLevelPercentage:             domain.NewDecimal2(stats.Summary.TopLevelPercentage),
			SubCategoriesPercentage:        domain.NewDecimal2(stats.Summary.SubCategoriesPercentage),
			CategoriesWithChildrenCount:    stats.Summary.CategoriesWithChildrenCount,
			CategoriesWithoutChildrenCount: stats.Summary.CategoriesWithoutChildrenCount,
			MaxDepthLevel:                  stats.Summary.MaxDepthLevel,
			AverageCategoriesPerDay:        domain.NewDecimal2(stats.Summary.AverageCategoriesPerDay),
			LatestCreationDate:             stats.Summary.LatestCreationDate,
			EarliestCreationDate:           stats.Summary.EarliestCreationDate,
		},
	}
}

// *==================== Update Map conversions (Harus snake case karena untuk database) ====================
func ToModelCategoryUpdateMap(payload *domain.UpdateCategoryPayload) map[string]any {
	updates := make(map[string]any)

	// Mengizinkan untuk set ParentID menjadi NULL
	if payload.ParentID != nil {
		if *payload.ParentID == "" {
			updates["parent_id"] = nil
		} else {
			updates["parent_id"] = *payload.ParentID
		}
	}
	if payload.CategoryCode != nil {
		updates["category_code"] = *payload.CategoryCode
	}

	return updates
}

func ToModelCategoryTranslationUpdateMap(payload *domain.UpdateCategoryTranslationPayload) map[string]any {
	updates := make(map[string]any)

	if payload.CategoryName != nil {
		updates["category_name"] = *payload.CategoryName
	}
	if payload.Description != nil {
		updates["description"] = payload.Description
	}

	return updates
}

func MapCategorySortFieldToColumn(field domain.CategorySortField) string {
	columnMap := map[domain.CategorySortField]string{
		domain.CategorySortByCategoryCode: "category_code",
		domain.CategorySortByCategoryName: "ct.category_name",
		domain.CategorySortByCreatedAt:    "c.created_at",
		domain.CategorySortByUpdatedAt:    "c.updated_at",
	}

	if column, exists := columnMap[field]; exists {
		return column
	}
	return "c.created_at"
}
