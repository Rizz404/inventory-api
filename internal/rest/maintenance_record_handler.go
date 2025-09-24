package rest

import (
	"strconv"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/rest/middleware"
	"github.com/Rizz404/inventory-api/internal/utils"
	"github.com/Rizz404/inventory-api/internal/web"
	"github.com/Rizz404/inventory-api/services/maintenance_record"
	"github.com/gofiber/fiber/v2"
)

type MaintenanceRecordHandler struct {
	Service maintenance_record.MaintenanceRecordService
}

func NewMaintenanceRecordHandler(app fiber.Router, s maintenance_record.MaintenanceRecordService) {
	handler := &MaintenanceRecordHandler{Service: s}

	records := app.Group("/maintenance/records")
	records.Post("/",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin, domain.RoleManager), // ! uncomment in production
		handler.CreateMaintenanceRecord,
	)
	records.Get("/", handler.GetMaintenanceRecordsPaginated)
	records.Get("/cursor", handler.GetMaintenanceRecordsCursor)
	records.Get("/count", handler.CountMaintenanceRecords)
	records.Get("/check/:id", handler.CheckMaintenanceRecordExists)
	records.Get("/statistics", handler.GetMaintenanceRecordStatistics)
	records.Get("/:id", handler.GetMaintenanceRecordById)
	records.Patch("/:id",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin, domain.RoleManager), // ! uncomment in production
		handler.UpdateMaintenanceRecord,
	)
	records.Delete("/:id",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin), // ! uncomment in production
		handler.DeleteMaintenanceRecord,
	)
}

func (h *MaintenanceRecordHandler) parseFiltersAndSort(c *fiber.Ctx) (query.Params, error) {
	params := query.Params{}

	// Search
	if s := c.Query("search"); s != "" {
		params.SearchQuery = &s
	}

	// Sort
	if sortBy := c.Query("sort_by"); sortBy != "" {
		params.Sort = &query.SortOptions{Field: sortBy, Order: c.Query("sort_order", "desc")}
	}

	// Filters
	filters := &postgresql.MaintenanceRecordFilterOptions{}
	if assetID := c.Query("asset_id"); assetID != "" {
		filters.AssetID = &assetID
	}
	if scheduleID := c.Query("schedule_id"); scheduleID != "" {
		filters.ScheduleID = &scheduleID
	}
	if performedByUser := c.Query("performed_by_user"); performedByUser != "" {
		filters.PerformedByUser = &performedByUser
	}
	if vendorName := c.Query("vendor_name"); vendorName != "" {
		filters.VendorName = &vendorName
	}
	if fromDate := c.Query("from_date"); fromDate != "" {
		// keep as string YYYY-MM-DD to match repo expectations
		fd := fromDate
		filters.FromDate = &fd
	}
	if toDate := c.Query("to_date"); toDate != "" {
		td := toDate
		filters.ToDate = &td
	}
	params.Filters = filters
	return params, nil
}

// ============== RECORDS: MUTATION ==============
func (h *MaintenanceRecordHandler) CreateMaintenanceRecord(c *fiber.Ctx) error {
	var payload domain.CreateMaintenanceRecordPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	// User performing maintenance (if authenticated)
	userID, _ := web.GetUserIDFromContext(c)
	record, err := h.Service.CreateMaintenanceRecord(c.Context(), &payload, userID)
	if err != nil {
		return web.HandleError(c, err)
	}
	return web.Success(c, fiber.StatusCreated, utils.SuccessMaintenanceRecordCreatedKey, record)
}

func (h *MaintenanceRecordHandler) UpdateMaintenanceRecord(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrMaintenanceRecordIDRequiredKey))
	}
	var payload domain.CreateMaintenanceRecordPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}
	record, err := h.Service.UpdateMaintenanceRecord(c.Context(), id, &payload)
	if err != nil {
		return web.HandleError(c, err)
	}
	return web.Success(c, fiber.StatusOK, utils.SuccessMaintenanceRecordUpdatedKey, record)
}

func (h *MaintenanceRecordHandler) DeleteMaintenanceRecord(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrMaintenanceRecordIDRequiredKey))
	}
	if err := h.Service.DeleteMaintenanceRecord(c.Context(), id); err != nil {
		return web.HandleError(c, err)
	}
	return web.Success(c, fiber.StatusOK, utils.SuccessMaintenanceRecordDeletedKey, nil)
}

// ============== RECORDS: QUERY ==============
func (h *MaintenanceRecordHandler) GetMaintenanceRecordsPaginated(c *fiber.Ctx) error {
	params, err := h.parseFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	params.Pagination = &query.PaginationOptions{Limit: limit, Offset: offset}

	langCode := web.GetLanguageFromContext(c)
	records, total, err := h.Service.GetMaintenanceRecordsPaginated(c.Context(), params, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}
	return web.SuccessWithPageInfo(c, fiber.StatusOK, utils.SuccessMaintenanceRecordRetrievedKey, records, int(total), limit, (offset/limit)+1)
}

func (h *MaintenanceRecordHandler) GetMaintenanceRecordsCursor(c *fiber.Ctx) error {
	params, err := h.parseFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	cursor := c.Query("cursor")
	params.Pagination = &query.PaginationOptions{Limit: limit, Cursor: cursor}

	langCode := web.GetLanguageFromContext(c)
	records, err := h.Service.GetMaintenanceRecordsCursor(c.Context(), params, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}
	var nextCursor string
	hasNextPage := len(records) == limit
	if hasNextPage {
		nextCursor = records[len(records)-1].ID
	}
	return web.SuccessWithCursor(c, fiber.StatusOK, utils.SuccessMaintenanceRecordRetrievedKey, records, nextCursor, hasNextPage, limit)
}

func (h *MaintenanceRecordHandler) GetMaintenanceRecordById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrMaintenanceRecordIDRequiredKey))
	}
	langCode := web.GetLanguageFromContext(c)
	record, err := h.Service.GetMaintenanceRecordById(c.Context(), id, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}
	return web.Success(c, fiber.StatusOK, utils.SuccessMaintenanceRecordRetrievedKey, record)
}

func (h *MaintenanceRecordHandler) CheckMaintenanceRecordExists(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrMaintenanceRecordIDRequiredKey))
	}
	exists, err := h.Service.CheckMaintenanceRecordExists(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}
	return web.Success(c, fiber.StatusOK, utils.SuccessMaintenanceRecordCountedKey, map[string]bool{"exists": exists})
}

func (h *MaintenanceRecordHandler) CountMaintenanceRecords(c *fiber.Ctx) error {
	params, err := h.parseFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}
	count, err := h.Service.CountMaintenanceRecords(c.Context(), params)
	if err != nil {
		return web.HandleError(c, err)
	}
	return web.Success(c, fiber.StatusOK, utils.SuccessMaintenanceRecordCountedKey, count)
}

// ============== STATISTICS ==============
func (h *MaintenanceRecordHandler) GetMaintenanceRecordStatistics(c *fiber.Ctx) error {
	stats, err := h.Service.GetMaintenanceRecordStatistics(c.Context())
	if err != nil {
		return web.HandleError(c, err)
	}
	return web.Success(c, fiber.StatusOK, utils.SuccessMaintenanceRecordStatisticsRetrievedKey, stats)
}
