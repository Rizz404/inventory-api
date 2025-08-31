package category

import (
	"context"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
)

type Repository interface {
	// * MUTATION
	CreateCategory(ctx context.Context, payload *domain.Category) (domain.Category, error)
	UpdateCategory(ctx context.Context, payload *domain.Category) (domain.Category, error)
	UpdateCategoryWithPayload(ctx context.Context, userId string, payload *domain.UpdateCategoryPayload) (domain.Category, error)
	DeleteCategory(ctx context.Context, userId string) error

	// * QUERY
	GetCategoriesPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.Category, error)
	GetCategoriesCursor(ctx context.Context, params query.Params, langCode string) ([]domain.Category, error)
	GetCategoryById(ctx context.Context, userId string) (domain.Category, error)
	GetCategoryByCategorynameOrEmail(ctx context.Context, name string, email string) (domain.Category, error)
	CheckCategoryExist(ctx context.Context, userId string) (bool, error)
	CountCategories(ctx context.Context, params query.Params) (int64, error)
}

type Service struct {
	Repo Repository
}
