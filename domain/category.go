package domain

import "time"

// --- Structs ---

type Category struct {
	ID           string                `json:"id"`
	ParentID     *string               `json:"parentId"`
	CategoryCode string                `json:"categoryCode"`
	CreatedAt    time.Time             `json:"createdAt"`
	UpdatedAt    time.Time             `json:"updatedAt"`
	Translations []CategoryTranslation `json:"translations,omitempty"`
}

type CategoryTranslation struct {
	ID           string  `json:"id"`
	CategoryID   string  `json:"categoryId"`
	LangCode     string  `json:"langCode"`
	CategoryName string  `json:"categoryName"`
	Description  *string `json:"description"`
}

type CategoryResponse struct {
	ID           string             `json:"id"`
	ParentID     *string            `json:"parentId"`
	CategoryCode string             `json:"categoryCode"`
	Name         string             `json:"name"`
	Description  *string            `json:"description"`
	Children     []CategoryResponse `json:"children"`
	CreatedAt    string             `json:"createdAt"`
	UpdatedAt    string             `json:"updatedAt"`
}

// --- Payloads ---

type CreateCategoryPayload struct {
	ParentID     *string                            `json:"parentId,omitempty" validate:"omitempty"`
	CategoryCode string                             `json:"categoryCode" validate:"required,max=20"`
	Translations []CreateCategoryTranslationPayload `json:"translations" validate:"required,min=1,dive"`
}

type CreateCategoryTranslationPayload struct {
	LangCode     string  `json:"langCode" validate:"required,max=5"`
	CategoryName string  `json:"categoryName" validate:"required,max=100"`
	Description  *string `json:"description,omitempty"`
}

type UpdateCategoryPayload struct {
	ParentID     *string                            `json:"parentId,omitempty" validate:"omitempty"`
	CategoryCode *string                            `json:"categoryCode,omitempty" validate:"omitempty,max=20"`
	Translations []UpdateCategoryTranslationPayload `json:"translations,omitempty" validate:"omitempty,dive"`
}

type UpdateCategoryTranslationPayload struct {
	LangCode     string  `json:"langCode" validate:"required,max=5"`
	CategoryName *string `json:"categoryName,omitempty" validate:"omitempty,max=100"`
	Description  *string `json:"description,omitempty"`
}

// --- Statistics ---

type CategoryStatistics struct {
	Total          CategoryCountStatistics     `json:"total"`
	ByHierarchy    CategoryHierarchyStatistics `json:"byHierarchy"`
	CreationTrends []CategoryCreationTrend     `json:"creationTrends"`
	Summary        CategorySummaryStatistics   `json:"summary"`
}

type CategoryCountStatistics struct {
	Count int `json:"count"`
}

type CategoryHierarchyStatistics struct {
	TopLevel     int `json:"topLevel"`
	WithChildren int `json:"withChildren"`
	WithParent   int `json:"withParent"`
}

type CategoryCreationTrend struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type CategorySummaryStatistics struct {
	TotalCategories                int     `json:"totalCategories"`
	TopLevelPercentage             float64 `json:"topLevelPercentage"`
	SubCategoriesPercentage        float64 `json:"subCategoriesPercentage"`
	CategoriesWithChildrenCount    int     `json:"categoriesWithChildrenCount"`
	CategoriesWithoutChildrenCount int     `json:"categoriesWithoutChildrenCount"`
	MaxDepthLevel                  int     `json:"maxDepthLevel"`
	AverageCategoriesPerDay        float64 `json:"averageCategoriesPerDay"`
	LatestCreationDate             string  `json:"latestCreationDate"`
	EarliestCreationDate           string  `json:"earliestCreationDate"`
}
