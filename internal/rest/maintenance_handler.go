package rest

import (
	"strconv"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/rest/middleware"
	"github.com/Rizz404/inventory-api/internal/utils"
	"github.com/Rizz404/inventory-api/internal/web"
	"github.com/Rizz404/inventory-api/services/maintenance"
	"github.com/gofiber/fiber/v2"
)

type MaintenanceHandler struct {
	Service maintenance.MaintenanceService
}

func NewMaintenanceHandler(app fiber.Router, s maintenance.MaintenanceService) {
	handler := &MaintenanceHandler{Service: s}

	m := app.Group("/maintenance")

	// SCHEDULES
	schedules := m.Group("/schedules")
	schedules.Post("/",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin, domain.RoleManager), // ! uncomment in production
		handler.CreateMaintenanceSchedule,
	)
	schedules.Get("/", handler.GetMaintenanceSchedulesPaginated)
	schedules.Get("/cursor", handler.GetMaintenanceSchedulesCursor)
	schedules.Get("/count", handler.CountMaintenanceSchedules)
	schedules.Get("/check/:id", handler.CheckMaintenanceScheduleExists)
	schedules.Get("/:id", handler.GetMaintenanceScheduleById)
	schedules.Patch("/:id",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin, domain.RoleManager), // ! uncomment in production
		handler.UpdateMaintenanceSchedule,
	)
	schedules.Delete("/:id",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin), // ! uncomment in production
		handler.DeleteMaintenanceSchedule,
	)

	// RECORDS
	records := m.Group("/records")
	records.Post("/",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin, domain.RoleManager), // ! uncomment in production
		handler.CreateMaintenanceRecord,
	)
	records.Get("/", handler.GetMaintenanceRecordsPaginated)
	records.Get("/cursor", handler.GetMaintenanceRecordsCursor)
	records.Get("/count", handler.CountMaintenanceRecords)
	records.Get("/check/:id", handler.CheckMaintenanceRecordExists)
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

	// STATISTICS
	m.Get("/statistics", handler.GetMaintenanceStatistics)
}

func (h *MaintenanceHandler) parseScheduleFiltersAndSort(c *fiber.Ctx) (query.Params, error) {
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
	filters := &postgresql.MaintenanceScheduleFilterOptions{}
	if assetID := c.Query("asset_id"); assetID != "" {
		filters.AssetID = &assetID
	}
	if mtype := c.Query("maintenance_type"); mtype != "" {
		t := domain.MaintenanceScheduleType(mtype)
		filters.MaintenanceType = &t
	}
	if status := c.Query("status"); status != "" {
		st := domain.ScheduleStatus(status)
		filters.Status = &st
	}
	if createdBy := c.Query("created_by"); createdBy != "" {
		filters.CreatedBy = &createdBy
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

func (h *MaintenanceHandler) parseRecordFiltersAndSort(c *fiber.Ctx) (query.Params, error) {
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

// ============== SCHEDULES: MUTATION ==============
func (h *MaintenanceHandler) CreateMaintenanceSchedule(c *fiber.Ctx) error {
	var payload domain.CreateMaintenanceSchedulePayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	// CreatedBy from auth
	userID, ok := web.GetUserIDFromContext(c)
	if !ok {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrUserIDRequiredKey))
	}

	schedule, err := h.Service.CreateMaintenanceSchedule(c.Context(), &payload, userID)
	if err != nil {
		return web.HandleError(c, err)
	}
	return web.Success(c, fiber.StatusCreated, utils.SuccessMaintenanceScheduleCreatedKey, schedule)
}

func (h *MaintenanceHandler) UpdateMaintenanceSchedule(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrMaintenanceScheduleIDRequiredKey))
	}
	var payload domain.CreateMaintenanceSchedulePayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}
	schedule, err := h.Service.UpdateMaintenanceSchedule(c.Context(), id, &payload)
	if err != nil {
		return web.HandleError(c, err)
	}
	return web.Success(c, fiber.StatusOK, utils.SuccessMaintenanceScheduleUpdatedKey, schedule)
}

func (h *MaintenanceHandler) DeleteMaintenanceSchedule(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrMaintenanceScheduleIDRequiredKey))
	}
	if err := h.Service.DeleteMaintenanceSchedule(c.Context(), id); err != nil {
		return web.HandleError(c, err)
	}
	return web.Success(c, fiber.StatusOK, utils.SuccessMaintenanceScheduleDeletedKey, nil)
}

// ============== SCHEDULES: QUERY ==============
func (h *MaintenanceHandler) GetMaintenanceSchedulesPaginated(c *fiber.Ctx) error {
	params, err := h.parseScheduleFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	params.Pagination = &query.PaginationOptions{Limit: limit, Offset: offset}

	langCode := web.GetLanguageFromContext(c)
	schedules, total, err := h.Service.GetMaintenanceSchedulesPaginated(c.Context(), params, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}
	return web.SuccessWithPageInfo(c, fiber.StatusOK, utils.SuccessMaintenanceScheduleRetrievedKey, schedules, int(total), limit, (offset/limit)+1)
}

func (h *MaintenanceHandler) GetMaintenanceSchedulesCursor(c *fiber.Ctx) error {
	params, err := h.parseScheduleFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	cursor := c.Query("cursor")
	params.Pagination = &query.PaginationOptions{Limit: limit, Cursor: cursor}

	langCode := web.GetLanguageFromContext(c)
	schedules, err := h.Service.GetMaintenanceSchedulesCursor(c.Context(), params, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}
	var nextCursor string
	hasNextPage := len(schedules) == limit
	if hasNextPage {
		nextCursor = schedules[len(schedules)-1].ID
	}
	return web.SuccessWithCursor(c, fiber.StatusOK, utils.SuccessMaintenanceScheduleRetrievedKey, schedules, nextCursor, hasNextPage, limit)
}

func (h *MaintenanceHandler) GetMaintenanceScheduleById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrMaintenanceScheduleIDRequiredKey))
	}
	langCode := web.GetLanguageFromContext(c)
	schedule, err := h.Service.GetMaintenanceScheduleById(c.Context(), id, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}
	return web.Success(c, fiber.StatusOK, utils.SuccessMaintenanceScheduleRetrievedKey, schedule)
}

func (h *MaintenanceHandler) CheckMaintenanceScheduleExists(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrMaintenanceScheduleIDRequiredKey))
	}
	exists, err := h.Service.CheckMaintenanceScheduleExists(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}
	return web.Success(c, fiber.StatusOK, utils.SuccessMaintenanceScheduleCountedKey, map[string]bool{"exists": exists})
}

func (h *MaintenanceHandler) CountMaintenanceSchedules(c *fiber.Ctx) error {
	params, err := h.parseScheduleFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}
	count, err := h.Service.CountMaintenanceSchedules(c.Context(), params)
	if err != nil {
		return web.HandleError(c, err)
	}
	return web.Success(c, fiber.StatusOK, utils.SuccessMaintenanceScheduleCountedKey, count)
}

// ============== RECORDS: MUTATION ==============
func (h *MaintenanceHandler) CreateMaintenanceRecord(c *fiber.Ctx) error {
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

func (h *MaintenanceHandler) UpdateMaintenanceRecord(c *fiber.Ctx) error {
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

func (h *MaintenanceHandler) DeleteMaintenanceRecord(c *fiber.Ctx) error {
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
func (h *MaintenanceHandler) GetMaintenanceRecordsPaginated(c *fiber.Ctx) error {
	params, err := h.parseRecordFiltersAndSort(c)
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

func (h *MaintenanceHandler) GetMaintenanceRecordsCursor(c *fiber.Ctx) error {
	params, err := h.parseRecordFiltersAndSort(c)
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

func (h *MaintenanceHandler) GetMaintenanceRecordById(c *fiber.Ctx) error {
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

func (h *MaintenanceHandler) CheckMaintenanceRecordExists(c *fiber.Ctx) error {
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

func (h *MaintenanceHandler) CountMaintenanceRecords(c *fiber.Ctx) error {
	params, err := h.parseRecordFiltersAndSort(c)
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
func (h *MaintenanceHandler) GetMaintenanceStatistics(c *fiber.Ctx) error {
	stats, err := h.Service.GetMaintenanceStatistics(c.Context())
	if err != nil {
		return web.HandleError(c, err)
	}
	return web.Success(c, fiber.StatusOK, utils.SuccessMaintenanceScheduleStatisticsRetrievedKey, stats)
}
