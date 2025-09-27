package rest

import (
	"strconv"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/rest/middleware"
	"github.com/Rizz404/inventory-api/internal/utils"
	"github.com/Rizz404/inventory-api/internal/web"
	"github.com/Rizz404/inventory-api/services/asset_movement"
	"github.com/gofiber/fiber/v2"
)

type AssetMovementHandler struct {
	Service asset_movement.AssetMovementService
}

func NewAssetMovementHandler(app fiber.Router, s asset_movement.AssetMovementService) {
	handler := &AssetMovementHandler{
		Service: s,
	}

	// * Bisa di group
	// ! routenya bisa tabrakan hati-hati
	movements := app.Group("/asset-movements")

	// * Create
	movements.Post("/",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin, domain.RoleManager), // ! jangan lupa uncomment pas production
		handler.CreateAssetMovement,
	)

	movements.Get("/", handler.GetAssetMovementsPaginated)
	movements.Get("/statistics", handler.GetAssetMovementStatistics)
	movements.Get("/cursor", handler.GetAssetMovementsCursor)
	movements.Get("/count", handler.CountAssetMovements)
	movements.Get("/check/:id", handler.CheckAssetMovementExists)
	movements.Get("/asset/:assetId", handler.GetAssetMovementsByAssetId)
	movements.Get("/:id", handler.GetAssetMovementById)
	movements.Patch("/:id",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin, domain.RoleManager), // ! jangan lupa uncomment pas production
		handler.UpdateAssetMovement,
	)
	movements.Delete("/:id",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin), // ! jangan lupa uncomment pas production
		handler.DeleteAssetMovement,
	)
}

func (h *AssetMovementHandler) parseAssetMovementFiltersAndSort(c *fiber.Ctx) (query.Params, error) {
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
	filters := &postgresql.AssetMovementFilterOptions{}

	// Parse asset ID filter
	if assetID := c.Query("asset_id"); assetID != "" {
		filters.AssetID = &assetID
	}

	// Parse from location ID filter
	if fromLocationID := c.Query("from_location_id"); fromLocationID != "" {
		filters.FromLocationID = &fromLocationID
	}

	// Parse to location ID filter
	if toLocationID := c.Query("to_location_id"); toLocationID != "" {
		filters.ToLocationID = &toLocationID
	}

	// Parse from user ID filter
	if fromUserID := c.Query("from_user_id"); fromUserID != "" {
		filters.FromUserID = &fromUserID
	}

	// Parse to user ID filter
	if toUserID := c.Query("to_user_id"); toUserID != "" {
		filters.ToUserID = &toUserID
	}

	// Parse moved by filter
	if movedBy := c.Query("moved_by"); movedBy != "" {
		filters.MovedBy = &movedBy
	}

	// * Parse date range filters
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		if parsedDate, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filters.DateFrom = &parsedDate
		}
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		if parsedDate, err := time.Parse("2006-01-02", dateTo); err == nil {
			filters.DateTo = &parsedDate
		}
	}

	params.Filters = filters

	return params, nil
}

// *===========================MUTATION===========================*
func (h *AssetMovementHandler) CreateAssetMovement(c *fiber.Ctx) error {
	var payload domain.CreateAssetMovementPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	// * Get user ID from context (set by auth middleware)
	userID, ok := web.GetUserIDFromContext(c)
	if !ok {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrUserIDRequiredKey))
	}

	movement, err := h.Service.CreateAssetMovement(c.Context(), &payload, userID)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusCreated, utils.SuccessAssetMovementCreatedKey, movement)
}

func (h *AssetMovementHandler) UpdateAssetMovement(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrAssetMovementIDRequiredKey))
	}

	var payload domain.UpdateAssetMovementPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	movement, err := h.Service.UpdateAssetMovement(c.Context(), id, &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessAssetMovementUpdatedKey, movement)
}

func (h *AssetMovementHandler) DeleteAssetMovement(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrAssetMovementIDRequiredKey))
	}

	err := h.Service.DeleteAssetMovement(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessAssetMovementDeletedKey, nil)
}

// *===========================QUERY===========================*
func (h *AssetMovementHandler) GetAssetMovementsPaginated(c *fiber.Ctx) error {
	params, err := h.parseAssetMovementFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	params.Pagination = &query.PaginationOptions{Limit: limit, Offset: offset}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	movements, total, err := h.Service.GetAssetMovementsPaginated(c.Context(), params, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.SuccessWithOffsetInfo(c, fiber.StatusOK, utils.SuccessAssetMovementRetrievedKey, movements, int(total), limit, (offset/limit)+1)
}

func (h *AssetMovementHandler) GetAssetMovementsCursor(c *fiber.Ctx) error {
	params, err := h.parseAssetMovementFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	cursor := c.Query("cursor")
	params.Pagination = &query.PaginationOptions{Limit: limit, Cursor: cursor}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	movements, err := h.Service.GetAssetMovementsCursor(c.Context(), params, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	var nextCursor string
	hasNextPage := len(movements) == limit
	if hasNextPage {
		nextCursor = movements[len(movements)-1].ID
	}

	return web.SuccessWithCursor(c, fiber.StatusOK, utils.SuccessAssetMovementRetrievedKey, movements, nextCursor, hasNextPage, limit)
}

func (h *AssetMovementHandler) GetAssetMovementById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrAssetMovementIDRequiredKey))
	}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	movement, err := h.Service.GetAssetMovementById(c.Context(), id, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessAssetMovementRetrievedKey, movement)
}

func (h *AssetMovementHandler) GetAssetMovementsByAssetId(c *fiber.Ctx) error {
	assetId := c.Params("assetId")
	if assetId == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrAssetIDRequiredKey))
	}

	params, err := h.parseAssetMovementFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	params.Pagination = &query.PaginationOptions{Limit: limit, Offset: offset}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	movements, err := h.Service.GetAssetMovementsByAssetId(c.Context(), assetId, params, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessAssetMovementRetrievedKey, movements)
}

func (h *AssetMovementHandler) CheckAssetMovementExists(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrAssetMovementIDRequiredKey))
	}

	exists, err := h.Service.CheckAssetMovementExists(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessAssetMovementExistenceCheckedKey, map[string]bool{"exists": exists})
}

func (h *AssetMovementHandler) CountAssetMovements(c *fiber.Ctx) error {
	params, err := h.parseAssetMovementFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	count, err := h.Service.CountAssetMovements(c.Context(), params)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessAssetMovementCountedKey, count)
}

func (h *AssetMovementHandler) GetAssetMovementStatistics(c *fiber.Ctx) error {
	stats, err := h.Service.GetAssetMovementStatistics(c.Context())
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessAssetMovementStatisticsRetrievedKey, stats)
}
