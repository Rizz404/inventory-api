package rest

import (
	"strconv"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/rest/middleware"
	"github.com/Rizz404/inventory-api/internal/utils"
	"github.com/Rizz404/inventory-api/internal/web"
	"github.com/Rizz404/inventory-api/services/maintenance_schedule"
	"github.com/gofiber/fiber/v2"
)

type MaintenanceScheduleHandler struct {
	Service maintenance_schedule.MaintenanceScheduleService
}

func NewMaintenanceScheduleHandler(app fiber.Router, s maintenance_schedule.MaintenanceScheduleService) {
	handler := &MaintenanceScheduleHandler{Service: s}

	schedules := app.Group("/maintenance/schedules")
	schedules.Post("/",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin, domain.RoleManager), // ! uncomment in production
		handler.CreateMaintenanceSchedule,
	)
	schedules.Get("/", handler.GetMaintenanceSchedulesPaginated)
	schedules.Get("/cursor", handler.GetMaintenanceSchedulesCursor)
	schedules.Get("/count", handler.CountMaintenanceSchedules)
	schedules.Get("/check/:id", handler.CheckMaintenanceScheduleExists)
	schedules.Get("/statistics", handler.GetMaintenanceScheduleStatistics)
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
}

func (h *MaintenanceScheduleHandler) parseFiltersAndSort(c *fiber.Ctx) (domain.MaintenanceScheduleParams, error) {
	params := domain.MaintenanceScheduleParams{}

	// Search
	if s := c.Query("search"); s != "" {
		params.SearchQuery = &s
	}

	// Sort
	if sortBy := c.Query("sortBy"); sortBy != "" {
		sortOrder := c.Query("sortOrder", "desc")
		params.Sort = &domain.MaintenanceScheduleSortOptions{
			Field: domain.MaintenanceScheduleSortField(sortBy),
			Order: domain.SortOrder(sortOrder),
		}
	}

	// Filters
	filters := &domain.MaintenanceScheduleFilterOptions{}
	if assetID := c.Query("assetId"); assetID != "" {
		filters.AssetID = &assetID
	}
	if mtype := c.Query("maintenanceType"); mtype != "" {
		t := domain.MaintenanceScheduleType(mtype)
		filters.MaintenanceType = &t
	}
	if state := c.Query("state"); state != "" {
		st := domain.ScheduleState(state)
		filters.State = &st
	}
	if createdBy := c.Query("createdBy"); createdBy != "" {
		filters.CreatedBy = &createdBy
	}
	if fromDate := c.Query("fromDate"); fromDate != "" {
		// keep as string YYYY-MM-DD to match repo expectations
		fd := fromDate
		filters.FromDate = &fd
	}
	if toDate := c.Query("toDate"); toDate != "" {
		td := toDate
		filters.ToDate = &td
	}
	params.Filters = filters
	return params, nil
}

// ============== SCHEDULES: MUTATION ==============
func (h *MaintenanceScheduleHandler) CreateMaintenanceSchedule(c *fiber.Ctx) error {
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

func (h *MaintenanceScheduleHandler) UpdateMaintenanceSchedule(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrMaintenanceScheduleIDRequiredKey))
	}

	var payload domain.UpdateMaintenanceSchedulePayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	schedule, err := h.Service.UpdateMaintenanceSchedule(c.Context(), id, &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessMaintenanceScheduleUpdatedKey, schedule)
}

func (h *MaintenanceScheduleHandler) DeleteMaintenanceSchedule(c *fiber.Ctx) error {
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
func (h *MaintenanceScheduleHandler) GetMaintenanceSchedulesPaginated(c *fiber.Ctx) error {
	params, err := h.parseFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	params.Pagination = &domain.PaginationOptions{Limit: limit, Offset: offset}

	langCode := web.GetLanguageFromContext(c)
	schedules, total, err := h.Service.GetMaintenanceSchedulesPaginated(c.Context(), params, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}
	return web.SuccessWithOffsetInfo(c, fiber.StatusOK, utils.SuccessMaintenanceScheduleRetrievedKey, schedules, int(total), limit, (offset/limit)+1)
}

func (h *MaintenanceScheduleHandler) GetMaintenanceSchedulesCursor(c *fiber.Ctx) error {
	params, err := h.parseFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	cursor := c.Query("cursor")
	params.Pagination = &domain.PaginationOptions{Limit: limit, Cursor: cursor}

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

func (h *MaintenanceScheduleHandler) GetMaintenanceScheduleById(c *fiber.Ctx) error {
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

func (h *MaintenanceScheduleHandler) CheckMaintenanceScheduleExists(c *fiber.Ctx) error {
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

func (h *MaintenanceScheduleHandler) CountMaintenanceSchedules(c *fiber.Ctx) error {
	params, err := h.parseFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}
	count, err := h.Service.CountMaintenanceSchedules(c.Context(), params)
	if err != nil {
		return web.HandleError(c, err)
	}
	return web.Success(c, fiber.StatusOK, utils.SuccessMaintenanceScheduleCountedKey, count)
}

// ============== STATISTICS ==============
func (h *MaintenanceScheduleHandler) GetMaintenanceScheduleStatistics(c *fiber.Ctx) error {
	stats, err := h.Service.GetMaintenanceScheduleStatistics(c.Context())
	if err != nil {
		return web.HandleError(c, err)
	}
	return web.Success(c, fiber.StatusOK, utils.SuccessMaintenanceScheduleStatisticsRetrievedKey, stats)
}
