package rest

import (
	"strconv"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/rest/middleware"
	"github.com/Rizz404/inventory-api/internal/utils"
	"github.com/Rizz404/inventory-api/internal/web"
	"github.com/Rizz404/inventory-api/services/scan_log"
	"github.com/gofiber/fiber/v2"
)

type ScanLogHandler struct {
	Service scan_log.ScanLogService
}

func NewScanLogHandler(app fiber.Router, s scan_log.ScanLogService) {
	handler := &ScanLogHandler{
		Service: s,
	}

	// * Bisa di group
	// ! routenya bisa tabrakan hati-hati
	scanLogs := app.Group("/scan-logs")

	// * Create
	scanLogs.Post("/",
		middleware.AuthMiddleware(),
		handler.CreateScanLog,
	)

	scanLogs.Get("/", handler.GetScanLogsPaginated)
	scanLogs.Get("/statistics", handler.GetScanLogStatistics)
	scanLogs.Get("/cursor", handler.GetScanLogsCursor)
	scanLogs.Get("/count", handler.CountScanLogs)
	scanLogs.Get("/user/:userId", handler.GetScanLogsByUserId)
	scanLogs.Get("/asset/:assetId", handler.GetScanLogsByAssetId)
	scanLogs.Get("/check/:id", handler.CheckScanLogExists)
	scanLogs.Get("/:id", handler.GetScanLogById)
	scanLogs.Delete("/:id",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin), // ! jangan lupa uncomment pas production
		handler.DeleteScanLog,
	)
}

func (h *ScanLogHandler) parseScanLogFiltersAndSort(c *fiber.Ctx) (domain.ScanLogParams, error) {
	params := domain.ScanLogParams{}

	// * Parse search query
	search := c.Query("search")
	if search != "" {
		params.SearchQuery = &search
	}

	// * Parse sorting options
	sortBy := c.Query("sortBy")
	if sortBy != "" {
		params.Sort = &domain.ScanLogSortOptions{
			Field: sortBy,
			Order: c.Query("sortOrder", "desc"),
		}
	}

	// * Parse filtering options
	filters := &domain.ScanLogFilterOptions{}

	// * Parse scan method filter
	if scanMethod := c.Query("scanMethod"); scanMethod != "" {
		method := domain.ScanMethodType(scanMethod)
		filters.ScanMethod = &method
	}

	// * Parse scan result filter
	if scanResult := c.Query("scanResult"); scanResult != "" {
		result := domain.ScanResultType(scanResult)
		filters.ScanResult = &result
	}

	// * Parse scanned by filter
	if scannedBy := c.Query("scannedBy"); scannedBy != "" {
		filters.ScannedBy = &scannedBy
	}

	// * Parse asset ID filter
	if assetId := c.Query("assetId"); assetId != "" {
		filters.AssetID = &assetId
	}

	// * Parse date range filters
	if dateFrom := c.Query("dateFrom"); dateFrom != "" {
		if parsedDate, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filters.DateFrom = &parsedDate
		}
	}

	if dateTo := c.Query("dateTo"); dateTo != "" {
		if parsedDate, err := time.Parse("2006-01-02", dateTo); err == nil {
			filters.DateTo = &parsedDate
		}
	}

	// * Parse coordinates filter
	if hasCoordinates := c.Query("hasCoordinates"); hasCoordinates != "" {
		if hasCoords, err := strconv.ParseBool(hasCoordinates); err == nil {
			filters.HasCoordinates = &hasCoords
		}
	}

	params.Filters = filters

	return params, nil
}

// *===========================MUTATION===========================*
func (h *ScanLogHandler) CreateScanLog(c *fiber.Ctx) error {
	var payload domain.CreateScanLogPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	// * Get user ID from auth context
	userId, ok := web.GetUserIDFromContext(c)
	if !ok {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrUserIDRequiredKey))
	}

	scanLog, err := h.Service.CreateScanLog(c.Context(), &payload, userId)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusCreated, utils.SuccessScanLogCreatedKey, scanLog)
}

func (h *ScanLogHandler) DeleteScanLog(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrScanLogIDRequiredKey))
	}

	err := h.Service.DeleteScanLog(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessScanLogDeletedKey, nil)
}

// *===========================QUERY===========================*
func (h *ScanLogHandler) GetScanLogsPaginated(c *fiber.Ctx) error {
	params, err := h.parseScanLogFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	params.Pagination = &domain.ScanLogPaginationOptions{Limit: limit, Offset: offset}

	scanLogs, total, err := h.Service.GetScanLogsPaginated(c.Context(), params)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.SuccessWithOffsetInfo(c, fiber.StatusOK, utils.SuccessScanLogRetrievedKey, scanLogs, int(total), limit, (offset/limit)+1)
}

func (h *ScanLogHandler) GetScanLogsCursor(c *fiber.Ctx) error {
	params, err := h.parseScanLogFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	cursor := c.Query("cursor")
	params.Pagination = &domain.ScanLogPaginationOptions{Limit: limit, Cursor: cursor}

	scanLogs, err := h.Service.GetScanLogsCursor(c.Context(), params)
	if err != nil {
		return web.HandleError(c, err)
	}

	var nextCursor string
	hasNextPage := len(scanLogs) == limit
	if hasNextPage {
		nextCursor = scanLogs[len(scanLogs)-1].ID
	}

	return web.SuccessWithCursor(c, fiber.StatusOK, utils.SuccessScanLogRetrievedKey, scanLogs, nextCursor, hasNextPage, limit)
}

func (h *ScanLogHandler) GetScanLogById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrScanLogIDRequiredKey))
	}

	scanLog, err := h.Service.GetScanLogById(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessScanLogRetrievedKey, scanLog)
}

func (h *ScanLogHandler) GetScanLogsByAssetId(c *fiber.Ctx) error {
	assetId := c.Params("assetId")
	if assetId == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrAssetIDRequiredKey))
	}

	params, err := h.parseScanLogFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	// * Apply pagination for asset scan logs
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	params.Pagination = &domain.ScanLogPaginationOptions{Limit: limit, Offset: offset}

	scanLogs, err := h.Service.GetScanLogsByAssetId(c.Context(), assetId, params)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessScanLogRetrievedKey, scanLogs)
}

func (h *ScanLogHandler) GetScanLogsByUserId(c *fiber.Ctx) error {
	userId := c.Params("userId")
	if userId == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrUserIDRequiredKey))
	}

	params, err := h.parseScanLogFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	// * Apply pagination for user scan logs
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	params.Pagination = &domain.ScanLogPaginationOptions{Limit: limit, Offset: offset}

	scanLogs, err := h.Service.GetScanLogsByUserId(c.Context(), userId, params)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessScanLogRetrievedKey, scanLogs)
}

func (h *ScanLogHandler) CheckScanLogExists(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrScanLogIDRequiredKey))
	}

	exists, err := h.Service.CheckScanLogExists(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessScanLogExistenceCheckedKey, map[string]bool{"exists": exists})
}

func (h *ScanLogHandler) CountScanLogs(c *fiber.Ctx) error {
	params, err := h.parseScanLogFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	count, err := h.Service.CountScanLogs(c.Context(), params)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessScanLogCountedKey, count)
}

func (h *ScanLogHandler) GetScanLogStatistics(c *fiber.Ctx) error {
	stats, err := h.Service.GetScanLogStatistics(c.Context())
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessScanLogStatisticsRetrievedKey, stats)
}
