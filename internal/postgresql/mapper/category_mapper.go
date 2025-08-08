package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
)

func findCategoryTranslation(translations []model.CategoryTranslation, langCode string) (name string, description *string) {
	for _, t := range translations {
		if t.LangCode == langCode {
			return t.CategoryName, t.Description
		}
	}
	for _, t := range translations {
		if t.LangCode == DefaultLangCode {
			return t.CategoryName, t.Description
		}
	}
	if len(translations) > 0 {
		return translations[0].CategoryName, translations[0].Description
	}
	return "", nil
}

func ToDomainCategoryResponse(m model.Category, langCode string) domain.CategoryResponse {
	name, desc := findCategoryTranslation(m.Translations, langCode)
	var children []domain.CategoryResponse
	if len(m.Children) > 0 {
		children = ToDomainCategoriesResponse(m.Children, langCode)
	}

	resp := domain.CategoryResponse{
		ID:           m.ID.String(),
		CategoryCode: m.CategoryCode,
		Name:         name,
		Description:  desc,
		Children:     children,
	}
	if m.ParentID != nil {
		resp.ParentID = Ptr(m.ParentID.String())
	}
	return resp
}

func ToDomainCategoriesResponse(m []model.Category, langCode string) []domain.CategoryResponse {
	responses := make([]domain.CategoryResponse, len(m))
	for i, category := range m {
		responses[i] = ToDomainCategoryResponse(category, langCode)
	}
	return responses
}
