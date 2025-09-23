package rest

import (
	"strconv"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/rest/middleware"
	"github.com/Rizz404/inventory-api/internal/utils"
	"github.com/Rizz404/inventory-api/internal/web"
	"github.com/Rizz404/inventory-api/services/category"
	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	Service category.CategoryService
}

func NewCategoryHandler(app fiber.Router, s category.CategoryService) {
	handler := &CategoryHandler{
		Service: s,
	}

	// * Bisa di group
	// ! routenya bisa tabrakan hati-hati
	categories := app.Group("/categories")

	// * Create
	categories.Post("/",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin), // ! jangan lupa uncomment pas production
		handler.CreateCategory,
	)

	categories.Get("/", handler.GetCategoriesPaginated)
	categories.Get("/statistics", handler.GetCategoryStatistics)
	categories.Get("/cursor", handler.GetCategoriesCursor)
	categories.Get("/count", handler.CountCategories)
	categories.Get("/code/:code", handler.GetCategoryByCode)
	categories.Get("/check/code/:code", handler.CheckCategoryCodeExists)
	categories.Get("/check/:id", handler.CheckCategoryExists)
	categories.Get("/:id", handler.GetCategoryById)
	categories.Patch("/:id",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin), // ! jangan lupa uncomment pas production
		handler.UpdateCategory,
	)
	categories.Delete("/:id",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin), // ! jangan lupa uncomment pas production
		handler.DeleteCategory,
	)
}

func (h *CategoryHandler) parseCategoryFiltersAndSort(c *fiber.Ctx) (query.Params, error) {
	params := query.Params{}

	// * Parse search query
	search := c.Query("search")
	if search != "" {
		params.SearchQuery = &search
	}

	// * Parse sorting options
	sortBy := c.Query("sort_by")
	if sortBy != "" {
		params.Sort = &query.SortOptions{
			Field: sortBy,
			Order: c.Query("sort_order", "desc"),
		}
	}

	// * Parse filtering options
	filters := &postgresql.CategoryFilterOptions{}

	if parentID := c.Query("parent_id"); parentID != "" {
		filters.ParentID = &parentID
	}

	if hasParentStr := c.Query("has_parent"); hasParentStr != "" {
		hasParent, err := strconv.ParseBool(hasParentStr)
		if err == nil {
			filters.HasParent = &hasParent
		}
	}

	params.Filters = filters

	return params, nil
}

// *===========================MUTATION===========================*
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var payload domain.CreateCategoryPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	category, err := h.Service.CreateCategory(c.Context(), &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusCreated, utils.SuccessCategoryCreatedKey, category)
}

func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrCategoryIDRequiredKey))
	}

	var payload domain.UpdateCategoryPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	category, err := h.Service.UpdateCategory(c.Context(), id, &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessCategoryUpdatedKey, category)
}

func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrCategoryIDRequiredKey))
	}

	err := h.Service.DeleteCategory(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessCategoryDeletedKey, nil)
}

// *===========================QUERY===========================*
func (h *CategoryHandler) GetCategoriesPaginated(c *fiber.Ctx) error {
	params, err := h.parseCategoryFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	params.Pagination = &query.PaginationOptions{Limit: limit, Offset: offset}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	categories, total, err := h.Service.GetCategoriesPaginated(c.Context(), params, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.SuccessWithPageInfo(c, fiber.StatusOK, utils.SuccessCategoryRetrievedKey, categories, int(total), limit, (offset/limit)+1)
}

func (h *CategoryHandler) GetCategoriesCursor(c *fiber.Ctx) error {
	params, err := h.parseCategoryFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	cursor := c.Query("cursor")
	params.Pagination = &query.PaginationOptions{Limit: limit, Cursor: cursor}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	categories, err := h.Service.GetCategoriesCursor(c.Context(), params, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	var nextCursor string
	hasNextPage := len(categories) == limit
	if hasNextPage {
		nextCursor = categories[len(categories)-1].ID
	}

	return web.SuccessWithCursor(c, fiber.StatusOK, utils.SuccessCategoryRetrievedKey, categories, nextCursor, hasNextPage, limit)
}

func (h *CategoryHandler) GetCategoryById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrCategoryIDRequiredKey))
	}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	category, err := h.Service.GetCategoryById(c.Context(), id, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessCategoryRetrievedKey, category)
}

func (h *CategoryHandler) GetCategoryByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrCategoryCodeRequiredKey))
	}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	category, err := h.Service.GetCategoryByCode(c.Context(), code, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessCategoryRetrievedByCodeKey, category)
}

func (h *CategoryHandler) CheckCategoryExists(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrCategoryIDRequiredKey))
	}

	exists, err := h.Service.CheckCategoryExists(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessCategoryExistenceCheckedKey, map[string]bool{"exists": exists})
}

func (h *CategoryHandler) CheckCategoryCodeExists(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrCategoryCodeRequiredKey))
	}

	exists, err := h.Service.CheckCategoryCodeExists(c.Context(), code)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessCategoryCodeExistenceCheckedKey, map[string]bool{"exists": exists})
}

func (h *CategoryHandler) CountCategories(c *fiber.Ctx) error {
	params, err := h.parseCategoryFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	count, err := h.Service.CountCategories(c.Context(), params)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessCategoryCountedKey, count)
}

func (h *CategoryHandler) GetCategoryStatistics(c *fiber.Ctx) error {
	stats, err := h.Service.GetCategoryStatistics(c.Context())
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessCategoryStatisticsRetrievedKey, stats)
}
