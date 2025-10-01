package rest

import (
	"strconv"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/rest/middleware"
	"github.com/Rizz404/inventory-api/internal/utils"
	"github.com/Rizz404/inventory-api/internal/web"
	"github.com/Rizz404/inventory-api/services/issue_report"
	"github.com/gofiber/fiber/v2"
)

type IssueReportHandler struct {
	Service issue_report.IssueReportService
}

func NewIssueReportHandler(app fiber.Router, s issue_report.IssueReportService) {
	handler := &IssueReportHandler{
		Service: s,
	}

	// * Bisa di group
	// ! routenya bisa tabrakan hati-hati
	issueReports := app.Group("/issue-reports")

	// * Create
	issueReports.Post("/",
		middleware.AuthMiddleware(),
		handler.CreateIssueReport,
	)

	issueReports.Get("/",
		middleware.OptionalAuth(), // Optional auth to filter by reporter
		handler.GetIssueReportsPaginated,
	)
	issueReports.Get("/statistics", handler.GetIssueReportStatistics)
	issueReports.Get("/cursor",
		middleware.OptionalAuth(), // Optional auth to filter by reporter
		handler.GetIssueReportsCursor,
	)
	issueReports.Get("/count",
		middleware.OptionalAuth(), // Optional auth to filter by reporter
		handler.CountIssueReports,
	)
	issueReports.Get("/check/:id", handler.CheckIssueReportExists)
	issueReports.Get("/:id", handler.GetIssueReportById)

	// * Resolve/Reopen operations
	issueReports.Patch("/:id/resolve",
		middleware.AuthMiddleware(),
		handler.ResolveIssueReport,
	)
	issueReports.Patch("/:id/reopen",
		middleware.AuthMiddleware(),
		handler.ReopenIssueReport,
	)

	issueReports.Patch("/:id",
		middleware.AuthMiddleware(),
		handler.UpdateIssueReport,
	)
	issueReports.Delete("/:id",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin), // ! jangan lupa uncomment pas production
		handler.DeleteIssueReport,
	)
}

func (h *IssueReportHandler) parseIssueReportFiltersAndSort(c *fiber.Ctx) (domain.IssueReportParams, error) {
	params := domain.IssueReportParams{}

	// * Parse search query
	search := c.Query("search")
	if search != "" {
		params.SearchQuery = &search
	}

	// * Parse sorting options
	sortBy := c.Query("sortBy")
	if sortBy != "" {
		params.Sort = &domain.IssueReportSortOptions{
			Field: sortBy,
			Order: c.Query("sortOrder", "desc"),
		}
	}

	// * Parse filtering options
	filters := &domain.IssueReportFilterOptions{}

	if assetID := c.Query("assetId"); assetID != "" {
		filters.AssetID = &assetID
	}

	if reportedBy := c.Query("reportedBy"); reportedBy != "" {
		filters.ReportedBy = &reportedBy
	}

	if resolvedBy := c.Query("resolvedBy"); resolvedBy != "" {
		filters.ResolvedBy = &resolvedBy
	}

	if issueType := c.Query("issueType"); issueType != "" {
		filters.IssueType = &issueType
	}

	if priority := c.Query("priority"); priority != "" {
		if p := domain.IssuePriority(priority); p != "" {
			filters.Priority = &p
		}
	}

	if status := c.Query("status"); status != "" {
		if s := domain.IssueStatus(status); s != "" {
			filters.Status = &s
		}
	}

	if isResolvedStr := c.Query("isResolved"); isResolvedStr != "" {
		isResolved, err := strconv.ParseBool(isResolvedStr)
		if err == nil {
			filters.IsResolved = &isResolved
		}
	}

	if dateFromStr := c.Query("dateFrom"); dateFromStr != "" {
		if dateFrom, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			filters.DateFrom = &dateFrom
		}
	}

	if dateToStr := c.Query("dateTo"); dateToStr != "" {
		if dateTo, err := time.Parse("2006-01-02", dateToStr); err == nil {
			filters.DateTo = &dateTo
		}
	}

	params.Filters = filters

	return params, nil
}

// *===========================MUTATION===========================*
func (h *IssueReportHandler) CreateIssueReport(c *fiber.Ctx) error {
	var payload domain.CreateIssueReportPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	// Get reporter ID from context (set by auth middleware)
	reporterID := c.Locals("userID")
	if reporterID == nil {
		return web.HandleError(c, domain.ErrUnauthorized("user ID not found"))
	}

	reporterIDStr, ok := reporterID.(string)
	if !ok {
		return web.HandleError(c, domain.ErrUnauthorized("invalid user ID"))
	}

	issueReport, err := h.Service.CreateIssueReport(c.Context(), &payload, reporterIDStr)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusCreated, utils.SuccessIssueReportCreatedKey, issueReport)
}

func (h *IssueReportHandler) UpdateIssueReport(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrIssueReportIDRequiredKey))
	}

	var payload domain.UpdateIssueReportPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	issueReport, err := h.Service.UpdateIssueReport(c.Context(), id, &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessIssueReportUpdatedKey, issueReport)
}

func (h *IssueReportHandler) ResolveIssueReport(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrIssueReportIDRequiredKey))
	}

	var payload domain.ResolveIssueReportPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	// Get resolver ID from context (set by auth middleware)
	resolverID := c.Locals("userID")
	if resolverID == nil {
		return web.HandleError(c, domain.ErrUnauthorized("user ID not found"))
	}

	resolverIDStr, ok := resolverID.(string)
	if !ok {
		return web.HandleError(c, domain.ErrUnauthorized("invalid user ID"))
	}

	issueReport, err := h.Service.ResolveIssueReport(c.Context(), id, resolverIDStr, &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessIssueReportResolvedKey, issueReport)
}

func (h *IssueReportHandler) ReopenIssueReport(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrIssueReportIDRequiredKey))
	}

	issueReport, err := h.Service.ReopenIssueReport(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessIssueReportReopenedKey, issueReport)
}

func (h *IssueReportHandler) DeleteIssueReport(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrIssueReportIDRequiredKey))
	}

	err := h.Service.DeleteIssueReport(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessIssueReportDeletedKey, nil)
}

// *===========================QUERY===========================*
func (h *IssueReportHandler) GetIssueReportsPaginated(c *fiber.Ctx) error {
	params, err := h.parseIssueReportFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	// If user is authenticated, they can see all reports, but we can add role-based filtering later
	// For now, we'll let them see all reports if they're authenticated

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	params.Pagination = &domain.IssueReportPaginationOptions{Limit: limit, Offset: offset}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	issueReports, total, err := h.Service.GetIssueReportsPaginated(c.Context(), params, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.SuccessWithOffsetInfo(c, fiber.StatusOK, utils.SuccessIssueReportRetrievedKey, issueReports, int(total), limit, (offset/limit)+1)
}

func (h *IssueReportHandler) GetIssueReportsCursor(c *fiber.Ctx) error {
	params, err := h.parseIssueReportFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	cursor := c.Query("cursor")
	params.Pagination = &domain.IssueReportPaginationOptions{Limit: limit, Cursor: cursor}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	issueReports, err := h.Service.GetIssueReportsCursor(c.Context(), params, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	var nextCursor string
	hasNextPage := len(issueReports) == limit
	if hasNextPage {
		nextCursor = issueReports[len(issueReports)-1].ID
	}

	return web.SuccessWithCursor(c, fiber.StatusOK, utils.SuccessIssueReportRetrievedKey, issueReports, nextCursor, hasNextPage, limit)
}

func (h *IssueReportHandler) GetIssueReportById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrIssueReportIDRequiredKey))
	}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	issueReport, err := h.Service.GetIssueReportById(c.Context(), id, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessIssueReportRetrievedKey, issueReport)
}

func (h *IssueReportHandler) CheckIssueReportExists(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrIssueReportIDRequiredKey))
	}

	exists, err := h.Service.CheckIssueReportExists(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessIssueReportExistenceCheckedKey, map[string]bool{"exists": exists})
}

func (h *IssueReportHandler) CountIssueReports(c *fiber.Ctx) error {
	params, err := h.parseIssueReportFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	count, err := h.Service.CountIssueReports(c.Context(), params)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessIssueReportCountedKey, count)
}

func (h *IssueReportHandler) GetIssueReportStatistics(c *fiber.Ctx) error {
	stats, err := h.Service.GetIssueReportStatistics(c.Context())
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessIssueReportStatisticsRetrievedKey, stats)
}
