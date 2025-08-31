package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/oklog/ulid/v2"
)

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

func ToDomainCategory(m *model.Category) domain.Category {
	domainCategory := domain.Category{
		ID:           m.ID.String(),
		CategoryCode: m.CategoryCode,
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

	return domainCategory
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

func ToDomainCategoryResponse(m *model.Category, langCode string) domain.CategoryResponse {
	response := domain.CategoryResponse{
		ID:           m.ID.String(),
		CategoryCode: m.CategoryCode,
	}

	if m.ParentID != nil && !m.ParentID.IsZero() {
		parentIDStr := m.ParentID.String()
		response.ParentID = &parentIDStr
	}

	for _, translation := range m.Translations {
		if translation.LangCode == langCode {
			response.Name = translation.CategoryName
			response.Description = translation.Description
			break
		}
	}

	if response.Name == "" && len(m.Translations) > 0 {
		response.Name = m.Translations[0].CategoryName
		response.Description = m.Translations[0].Description
	}

	return response
}

func ToDomainCategoriesResponse(m []model.Category, langCode string) []domain.CategoryResponse {
	responses := make([]domain.CategoryResponse, len(m))
	for i, category := range m {
		responses[i] = ToDomainCategoryResponse(&category, langCode)
	}
	return responses
}

func BuildCategoryHierarchy(categories []domain.CategoryResponse) []domain.CategoryResponse {
	categoryMap := make(map[string]*domain.CategoryResponse)
	var rootCategories []domain.CategoryResponse

	// Gunakan slice pointer agar bisa memodifikasi item di map secara langsung
	for i := range categories {
		categoryMap[categories[i].ID] = &categories[i]
	}

	for _, category := range categories {
		// ParentID bisa nil, jadi cek dulu
		if category.ParentID == nil || *category.ParentID == "" {
			rootCategories = append(rootCategories, category)
		} else {
			if parent, exists := categoryMap[*category.ParentID]; exists {
				if parent.Children == nil {
					parent.Children = make([]domain.CategoryResponse, 0)
				}
				parent.Children = append(parent.Children, category)
			}
		}
	}

	return rootCategories
}

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
