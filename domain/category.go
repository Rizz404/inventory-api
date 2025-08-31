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
	ParentID     *string            `json:"parentId,omitempty"`
	CategoryCode string             `json:"categoryCode"`
	Name         string             `json:"name"`
	Description  *string            `json:"description,omitempty"`
	Children     []CategoryResponse `json:"children,omitempty"`
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
