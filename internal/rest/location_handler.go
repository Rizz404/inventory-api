package rest

import (
	"strconv"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/rest/middleware"
	"github.com/Rizz404/inventory-api/internal/utils"
	"github.com/Rizz404/inventory-api/internal/web"
	"github.com/Rizz404/inventory-api/services/location"
	"github.com/gofiber/fiber/v2"
)

type LocationHandler struct {
	Service location.LocationService
}

func NewLocationHandler(app fiber.Router, s location.LocationService) {
	handler := &LocationHandler{
		Service: s,
	}

	// * Bisa di group
	// ! routenya bisa tabrakan hati-hati
	locations := app.Group("/locations")

	// * Create
	locations.Post("/",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin), // ! jangan lupa uncomment pas production
		handler.CreateLocation,
	)

	locations.Get("/", handler.GetLocationsPaginated)
	locations.Get("/statistics", handler.GetLocationStatistics)
	locations.Get("/cursor", handler.GetLocationsCursor)
	locations.Get("/count", handler.CountLocations)
	locations.Get("/code/:code", handler.GetLocationByCode)
	locations.Get("/check/code/:code", handler.CheckLocationCodeExists)
	locations.Get("/check/:id", handler.CheckLocationExists)
	locations.Get("/:id", handler.GetLocationById)
	locations.Patch("/:id",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin), // ! jangan lupa uncomment pas production
		handler.UpdateLocation,
	)
	locations.Delete("/:id",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin), // ! jangan lupa uncomment pas production
		handler.DeleteLocation,
	)
}

func (h *LocationHandler) parseLocationFiltersAndSort(c *fiber.Ctx) (domain.LocationParams, error) {
	params := domain.LocationParams{}

	// * Parse search query
	search := c.Query("search")
	if search != "" {
		params.SearchQuery = &search
	}

	// * Parse sorting options
	sortBy := c.Query("sortBy")
	if sortBy != "" {
		sortOrder := c.Query("sortOrder", "desc")
		params.Sort = &domain.LocationSortOptions{
			Field: domain.LocationSortField(sortBy),
			Order: domain.SortOrder(sortOrder),
		}
	}

	// * Parse filtering options
	filters := &domain.LocationFilterOptions{}

	params.Filters = filters

	return params, nil
}

// *===========================MUTATION===========================*
func (h *LocationHandler) CreateLocation(c *fiber.Ctx) error {
	var payload domain.CreateLocationPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	location, err := h.Service.CreateLocation(c.Context(), &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusCreated, utils.SuccessLocationCreatedKey, location)
}

func (h *LocationHandler) UpdateLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrLocationIDRequiredKey))
	}

	var payload domain.UpdateLocationPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	location, err := h.Service.UpdateLocation(c.Context(), id, &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessLocationUpdatedKey, location)
}

func (h *LocationHandler) DeleteLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrLocationIDRequiredKey))
	}

	err := h.Service.DeleteLocation(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessLocationDeletedKey, nil)
}

// *===========================QUERY===========================*
func (h *LocationHandler) GetLocationsPaginated(c *fiber.Ctx) error {
	params, err := h.parseLocationFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	params.Pagination = &domain.PaginationOptions{Limit: limit, Offset: offset}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	locations, total, err := h.Service.GetLocationsPaginated(c.Context(), params, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.SuccessWithOffsetInfo(c, fiber.StatusOK, utils.SuccessLocationRetrievedKey, locations, int(total), limit, (offset/limit)+1)
}

func (h *LocationHandler) GetLocationsCursor(c *fiber.Ctx) error {
	params, err := h.parseLocationFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	cursor := c.Query("cursor")
	params.Pagination = &domain.PaginationOptions{Limit: limit, Cursor: cursor}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	locations, err := h.Service.GetLocationsCursor(c.Context(), params, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	var nextCursor string
	hasNextPage := len(locations) == limit
	if hasNextPage {
		nextCursor = locations[len(locations)-1].ID
	}

	return web.SuccessWithCursor(c, fiber.StatusOK, utils.SuccessLocationRetrievedKey, locations, nextCursor, hasNextPage, limit)
}

func (h *LocationHandler) GetLocationById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrLocationIDRequiredKey))
	}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	location, err := h.Service.GetLocationById(c.Context(), id, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessLocationRetrievedKey, location)
}

func (h *LocationHandler) GetLocationByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrLocationCodeRequiredKey))
	}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	location, err := h.Service.GetLocationByCode(c.Context(), code, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessLocationRetrievedByCodeKey, location)
}

func (h *LocationHandler) CheckLocationExists(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrLocationIDRequiredKey))
	}

	exists, err := h.Service.CheckLocationExists(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessLocationExistenceCheckedKey, map[string]bool{"exists": exists})
}

func (h *LocationHandler) CheckLocationCodeExists(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrLocationCodeRequiredKey))
	}

	exists, err := h.Service.CheckLocationCodeExists(c.Context(), code)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessLocationCodeExistenceCheckedKey, map[string]bool{"exists": exists})
}

func (h *LocationHandler) CountLocations(c *fiber.Ctx) error {
	params, err := h.parseLocationFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	count, err := h.Service.CountLocations(c.Context(), params)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessLocationCountedKey, count)
}

func (h *LocationHandler) GetLocationStatistics(c *fiber.Ctx) error {
	stats, err := h.Service.GetLocationStatistics(c.Context())
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessLocationStatisticsRetrievedKey, stats)
}
