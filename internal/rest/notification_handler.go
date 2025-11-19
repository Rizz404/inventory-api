package rest

import (
	"strconv"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/rest/middleware"
	"github.com/Rizz404/inventory-api/internal/utils"
	"github.com/Rizz404/inventory-api/internal/web"
	"github.com/Rizz404/inventory-api/services/notification"
	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct {
	Service notification.NotificationService
}

func NewNotificationHandler(app fiber.Router, s notification.NotificationService) {
	handler := &NotificationHandler{
		Service: s,
	}

	// * Bisa di group
	// ! routenya bisa tabrakan hati-hati
	notifications := app.Group("/notifications")

	// * Create
	notifications.Post("/",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin), // ! jangan lupa uncomment pas production
		handler.CreateNotification,
	)

	notifications.Get("/",
		middleware.OptionalAuth(), // Optional auth to filter by user
		handler.GetNotificationsPaginated,
	)
	notifications.Get("/statistics", handler.GetNotificationStatistics)
	notifications.Get("/cursor",
		middleware.OptionalAuth(), // Optional auth to filter by user
		handler.GetNotificationsCursor,
	)
	notifications.Get("/count",
		middleware.OptionalAuth(), // Optional auth to filter by user
		handler.CountNotifications,
	)
	notifications.Get("/check/:id", handler.CheckNotificationExists)
	notifications.Get("/:id", handler.GetNotificationById)

	// * Mark operations (batch update)
	notifications.Patch("/mark-read",
		middleware.AuthMiddleware(),
		handler.MarkNotificationsAsRead,
	)
	notifications.Patch("/mark-unread",
		middleware.AuthMiddleware(),
		handler.MarkNotificationsAsUnread,
	)

	notifications.Patch("/:id",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin), // ! jangan lupa uncomment pas production
		handler.UpdateNotification,
	)
	notifications.Delete("/:id",
		middleware.AuthMiddleware(),
		// middleware.AuthorizeRole(domain.RoleAdmin), // ! jangan lupa uncomment pas production
		handler.DeleteNotification,
	)
}

func (h *NotificationHandler) parseNotificationFiltersAndSort(c *fiber.Ctx) (domain.NotificationParams, error) {
	params := domain.NotificationParams{}

	// * Parse search query
	search := c.Query("search")
	if search != "" {
		params.SearchQuery = &search
	}

	// * Parse sorting options
	sortBy := c.Query("sortBy")
	if sortBy != "" {
		sortOrder := c.Query("sortOrder", "desc")
		params.Sort = &domain.NotificationSortOptions{
			Field: domain.NotificationSortField(sortBy),
			Order: domain.SortOrder(sortOrder),
		}
	}

	// * Parse filtering options
	filters := &domain.NotificationFilterOptions{}

	if userID := c.Query("userId"); userID != "" {
		filters.UserID = &userID
	}

	if relatedEntityType := c.Query("relatedEntityType"); relatedEntityType != "" {
		filters.RelatedEntityType = &relatedEntityType
	}

	if relatedEntityID := c.Query("relatedEntityId"); relatedEntityID != "" {
		filters.RelatedEntityID = &relatedEntityID
	}

	if relatedAssetID := c.Query("relatedAssetId"); relatedAssetID != "" {
		filters.RelatedAssetID = &relatedAssetID
	}

	if notificationType := c.Query("type"); notificationType != "" {
		if nType := domain.NotificationType(notificationType); nType != "" {
			filters.Type = &nType
		}
	}

	if priorityStr := c.Query("priority"); priorityStr != "" {
		if priority := domain.NotificationPriority(priorityStr); priority != "" {
			filters.Priority = &priority
		}
	}

	if isReadStr := c.Query("isRead"); isReadStr != "" {
		isRead, err := strconv.ParseBool(isReadStr)
		if err == nil {
			filters.IsRead = &isRead
		}
	}

	params.Filters = filters

	return params, nil
}

// *===========================MUTATION===========================*
func (h *NotificationHandler) CreateNotification(c *fiber.Ctx) error {
	var payload domain.CreateNotificationPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	notification, err := h.Service.CreateNotification(c.Context(), &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusCreated, utils.SuccessNotificationCreatedKey, notification)
}

func (h *NotificationHandler) UpdateNotification(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrNotificationIDRequiredKey))
	}

	var payload domain.UpdateNotificationPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	notification, err := h.Service.UpdateNotification(c.Context(), id, &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessNotificationUpdatedKey, notification)
}

func (h *NotificationHandler) DeleteNotification(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrNotificationIDRequiredKey))
	}

	err := h.Service.DeleteNotification(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessNotificationDeletedKey, nil)
}

func (h *NotificationHandler) MarkNotificationsAsRead(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return web.HandleError(c, domain.ErrUnauthorized("user ID not found"))
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return web.HandleError(c, domain.ErrUnauthorized("invalid user ID"))
	}

	var payload domain.MarkNotificationsPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	err := h.Service.MarkNotifications(c.Context(), userIDStr, payload.NotificationIDs, true)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessNotificationMarkedAsReadKey, nil)
}

func (h *NotificationHandler) MarkNotificationsAsUnread(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return web.HandleError(c, domain.ErrUnauthorized("user ID not found"))
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return web.HandleError(c, domain.ErrUnauthorized("invalid user ID"))
	}

	var payload domain.MarkNotificationsPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	err := h.Service.MarkNotifications(c.Context(), userIDStr, payload.NotificationIDs, false)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessNotificationMarkedAsUnreadKey, nil)
}

// *===========================QUERY===========================*
func (h *NotificationHandler) GetNotificationsPaginated(c *fiber.Ctx) error {
	params, err := h.parseNotificationFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	// If user is authenticated, automatically filter by their notifications unless explicitly filtering by userId
	if userID := c.Locals("userID"); userID != nil && c.Query("userId") == "" {
		if userIDStr, ok := userID.(string); ok {
			if params.Filters != nil {
				params.Filters.UserID = &userIDStr
			}
		}
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	params.Pagination = &domain.PaginationOptions{Limit: limit, Offset: offset}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	notifications, total, err := h.Service.GetNotificationsPaginated(c.Context(), params, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.SuccessWithOffsetInfo(c, fiber.StatusOK, utils.SuccessNotificationRetrievedKey, notifications, int(total), limit, (offset/limit)+1)
}

func (h *NotificationHandler) GetNotificationsCursor(c *fiber.Ctx) error {
	params, err := h.parseNotificationFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	// If user is authenticated, automatically filter by their notifications unless explicitly filtering by userId
	if userID := c.Locals("userID"); userID != nil && c.Query("userId") == "" {
		if userIDStr, ok := userID.(string); ok {
			if params.Filters != nil {
				params.Filters.UserID = &userIDStr
			}
		}
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	cursor := c.Query("cursor")
	params.Pagination = &domain.PaginationOptions{Limit: limit, Cursor: cursor}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	notifications, err := h.Service.GetNotificationsCursor(c.Context(), params, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	var nextCursor string
	hasNextPage := len(notifications) == limit
	if hasNextPage {
		nextCursor = notifications[len(notifications)-1].ID
	}

	return web.SuccessWithCursor(c, fiber.StatusOK, utils.SuccessNotificationRetrievedKey, notifications, nextCursor, hasNextPage, limit)
}

func (h *NotificationHandler) GetNotificationById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrNotificationIDRequiredKey))
	}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	notification, err := h.Service.GetNotificationById(c.Context(), id, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessNotificationRetrievedKey, notification)
}

func (h *NotificationHandler) CheckNotificationExists(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrNotificationIDRequiredKey))
	}

	exists, err := h.Service.CheckNotificationExists(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessNotificationExistenceCheckedKey, map[string]bool{"exists": exists})
}

func (h *NotificationHandler) CountNotifications(c *fiber.Ctx) error {
	params, err := h.parseNotificationFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	// If user is authenticated, automatically filter by their notifications unless explicitly filtering by userId
	if userID := c.Locals("userID"); userID != nil && c.Query("userId") == "" {
		if userIDStr, ok := userID.(string); ok {
			if params.Filters != nil {
				params.Filters.UserID = &userIDStr
			}
		}
	}

	count, err := h.Service.CountNotifications(c.Context(), params)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessNotificationCountedKey, count)
}

func (h *NotificationHandler) GetNotificationStatistics(c *fiber.Ctx) error {
	stats, err := h.Service.GetNotificationStatistics(c.Context())
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessNotificationStatisticsRetrievedKey, stats)
}
