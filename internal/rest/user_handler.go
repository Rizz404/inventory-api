package rest

import (
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/rest/middleware"
	"github.com/Rizz404/inventory-api/internal/utils"
	"github.com/Rizz404/inventory-api/internal/web"
	"github.com/Rizz404/inventory-api/services/user"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	Service user.UserService
}

func NewUserHandler(app fiber.Router, s user.UserService) {
	handler := &UserHandler{
		Service: s,
	}

	// * Bisa di group
	// ! routenya bisa tabrakan hati-hati
	users := app.Group("/users")

	// * Create
	users.Post("/export/list",
		middleware.AuthMiddleware(),
		handler.ExportUserList,
	)
	users.Post("/",
		middleware.AuthMiddleware(),
		middleware.AuthorizeRole(domain.RoleAdmin),
		handler.CreateUser,
	)
	users.Post("/bulk",
		middleware.AuthMiddleware(),
		middleware.AuthorizeRole(domain.RoleAdmin),
		handler.BulkCreateUsers,
	)

	users.Get("/", handler.GetUsersPaginated)
	users.Get("/statistics", handler.GetUserStatistics)
	users.Get("/cursor", handler.GetUsersCursor)
	users.Get("/count", handler.CountUsers)
	users.Get("/profile", middleware.AuthMiddleware(), handler.GetCurrentUser)
	users.Patch("/profile", middleware.AuthMiddleware(), handler.UpdateCurrentUser)
	users.Patch("/profile/password", middleware.AuthMiddleware(), handler.ChangeCurrentUserPassword)
	users.Get("/name/:name", handler.GetUserByName)
	users.Get("/email/:email", handler.GetUserByEmail)
	users.Get("/check/name/:name", handler.CheckNameExists)
	users.Get("/check/email/:email", handler.CheckEmailExists)
	users.Get("/check/:id", handler.CheckUserExists)
	users.Get("/:id", handler.GetUserById)
	users.Patch("/:id", handler.UpdateUser)
	users.Patch("/:id/password", middleware.AuthMiddleware(), handler.ChangePassword)
	users.Delete("/:id", handler.DeleteUser)
	users.Post("/bulk-delete",
		middleware.AuthMiddleware(),
		middleware.AuthorizeRole(domain.RoleAdmin),
		handler.BulkDeleteUsers,
	)
}

func (h *UserHandler) parseUserFiltersAndSort(c *fiber.Ctx) (domain.UserParams, error) {
	params := domain.UserParams{}

	// * Parse search query
	search := c.Query("search")
	if search != "" {
		params.SearchQuery = &search
	}

	// * Parse sorting options
	sortBy := c.Query("sortBy")
	if sortBy != "" {
		sortOrder := c.Query("sortOrder", "desc")
		params.Sort = &domain.UserSortOptions{
			Field: domain.UserSortField(sortBy),
			Order: domain.SortOrder(sortOrder),
		}
	}

	// * Parse filtering options
	roleStr := c.Query("role")
	var role *domain.UserRole
	if roleStr != "" {
		r := domain.UserRole(roleStr)
		role = &r
	}

	filters := &domain.UserFilterOptions{Role: role}

	if isActiveStr := c.Query("isActive"); isActiveStr != "" {
		isActive, err := strconv.ParseBool(isActiveStr)
		if err == nil {
			filters.IsActive = &isActive
		}
	}

	if employeeID := c.Query("employeeId"); employeeID != "" {
		filters.EmployeeID = &employeeID
	}

	params.Filters = filters

	return params, nil
}

// *===========================MUTATION===========================*
func (h *UserHandler) BulkCreateUsers(c *fiber.Ctx) error {
	var payload domain.BulkCreateUsersPayload

	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	users, err := h.Service.BulkCreateUsers(c.Context(), &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusCreated, utils.SuccessUsersBulkCreatedKey, users)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var payload domain.CreateUserPayload

	// Check content type to determine parsing method
	contentType := c.Get("Content-Type")
	var avatarFile *multipart.FileHeader

	if strings.Contains(contentType, "multipart/form-data") {
		// Parse multipart form data
		if err := web.ParseFormAndValidate(c, &payload); err != nil {
			return web.HandleError(c, err)
		}

		// Try to get avatar file (optional)
		file, err := c.FormFile("avatar")
		if err == nil {
			// Validate avatar file before processing (max 5MB)
			if validationErr := web.ValidateImageFile(file, "avatar", 5); validationErr != nil {
				return web.HandleError(c, domain.ErrBadRequest(web.FormatFileValidationError(validationErr)))
			}
			avatarFile = file
		}
		// Note: We don't return error if avatar file is missing since it's optional
	} else {
		// Parse JSON/form-urlencoded data
		if err := web.ParseAndValidate(c, &payload); err != nil {
			return web.HandleError(c, err)
		}
	}

	user, err := h.Service.CreateUser(c.Context(), &payload, avatarFile)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusCreated, utils.SuccessUserCreatedKey, user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrUserIDRequiredKey))
	}

	var payload domain.UpdateUserPayload

	// Check content type to determine parsing method
	contentType := c.Get("Content-Type")
	var avatarFile *multipart.FileHeader

	if strings.Contains(contentType, "multipart/form-data") {
		// Parse multipart form data
		if err := web.ParseFormAndValidate(c, &payload); err != nil {
			return web.HandleError(c, err)
		}

		// Try to get avatar file (optional)
		file, err := c.FormFile("avatar")
		if err == nil {
			// Validate avatar file before processing (max 5MB)
			if validationErr := web.ValidateImageFile(file, "avatar", 5); validationErr != nil {
				return web.HandleError(c, domain.ErrBadRequest(web.FormatFileValidationError(validationErr)))
			}
			avatarFile = file
		}
		// Note: We don't return error if avatar file is missing since it's optional
	} else {
		// Parse JSON/form-urlencoded data
		if err := web.ParseAndValidate(c, &payload); err != nil {
			return web.HandleError(c, err)
		}
	}

	user, err := h.Service.UpdateUser(c.Context(), id, &payload, avatarFile)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessUserUpdatedKey, user)
}

func (h *UserHandler) UpdateCurrentUser(c *fiber.Ctx) error {
	id, ok := web.GetUserIDFromContext(c)
	if !ok {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrUserIDRequiredKey))
	}

	var payload domain.UpdateUserPayload

	// Check content type to determine parsing method
	contentType := c.Get("Content-Type")
	var avatarFile *multipart.FileHeader

	if strings.Contains(contentType, "multipart/form-data") {
		// Parse multipart form data
		if err := web.ParseFormAndValidate(c, &payload); err != nil {
			return web.HandleError(c, err)
		}

		// Try to get avatar file (optional)
		file, err := c.FormFile("avatar")
		if err == nil {
			// Validate avatar file before processing (max 5MB)
			if validationErr := web.ValidateImageFile(file, "avatar", 5); validationErr != nil {
				return web.HandleError(c, domain.ErrBadRequest(web.FormatFileValidationError(validationErr)))
			}
			avatarFile = file
		}
		// Note: We don't return error if avatar file is missing since it's optional
	} else {
		// Parse JSON/form-urlencoded data
		if err := web.ParseAndValidate(c, &payload); err != nil {
			return web.HandleError(c, err)
		}
	}

	user, err := h.Service.UpdateUser(c.Context(), id, &payload, avatarFile)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessUserUpdatedKey, user)
}

func (h *UserHandler) ChangeCurrentUserPassword(c *fiber.Ctx) error {
	id, ok := web.GetUserIDFromContext(c)
	if !ok {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrUserIDRequiredKey))
	}

	var payload domain.ChangePasswordPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	if err := h.Service.ChangeCurrentUserPassword(c.Context(), id, &payload); err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessUpdatedKey, nil)
}

func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrUserIDRequiredKey))
	}

	var payload domain.ChangePasswordPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	// For admin changing another user's password we ignore OldPassword, but service.ChangePassword expects payload.NewPassword
	if err := h.Service.ChangePassword(c.Context(), id, &payload); err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessUpdatedKey, nil)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrUserIDRequiredKey))
	}

	err := h.Service.DeleteUser(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessUserDeletedKey, nil)
}

func (h *UserHandler) BulkDeleteUsers(c *fiber.Ctx) error {
	var payload domain.BulkDeleteUsersPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	result, err := h.Service.BulkDeleteUsers(c.Context(), &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessUsersBulkDeletedKey, result)
}

// *===========================QUERY===========================*
func (h *UserHandler) GetUsersPaginated(c *fiber.Ctx) error {
	params, err := h.parseUserFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	params.Pagination = &domain.PaginationOptions{Limit: limit, Offset: offset}

	users, total, err := h.Service.GetUsersPaginated(c.Context(), params)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.SuccessWithOffsetInfo(c, fiber.StatusOK, utils.SuccessUserRetrievedKey, users, int(total), limit, (offset/limit)+1)
}

func (h *UserHandler) GetUsersCursor(c *fiber.Ctx) error {
	params, err := h.parseUserFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	cursor := c.Query("cursor")
	params.Pagination = &domain.PaginationOptions{Limit: limit, Cursor: cursor}

	users, err := h.Service.GetUsersCursor(c.Context(), params)
	if err != nil {
		return web.HandleError(c, err)
	}

	var nextCursor string
	hasNextPage := len(users) == limit
	if hasNextPage {
		nextCursor = users[len(users)-1].ID
	}

	return web.SuccessWithCursor(c, fiber.StatusOK, utils.SuccessUserRetrievedKey, users, nextCursor, hasNextPage, limit)
}
func (h *UserHandler) GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrUserIDRequiredKey))
	}

	user, err := h.Service.GetUserById(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessUserRetrievedKey, user)
}

func (h *UserHandler) GetUserByName(c *fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrUserNameRequiredKey))
	}

	user, err := h.Service.GetUserByName(c.Context(), name)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessUserRetrievedByNameKey, user)
}

func (h *UserHandler) GetUserByEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrUserEmailRequiredKey))
	}

	user, err := h.Service.GetUserByEmail(c.Context(), email)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessUserRetrievedByEmailKey, user)
}

func (h *UserHandler) GetCurrentUser(c *fiber.Ctx) error {
	id, ok := web.GetUserIDFromContext(c)
	if !ok {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrUserIDRequiredKey))
	}

	user, err := h.Service.GetUserById(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessUserRetrievedKey, user)
}

func (h *UserHandler) CheckUserExists(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrUserIDRequiredKey))
	}

	exists, err := h.Service.CheckUserExists(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessUserExistenceCheckedKey, map[string]bool{"exists": exists})
}

func (h *UserHandler) CheckNameExists(c *fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrUserNameRequiredKey))
	}

	exists, err := h.Service.CheckNameExists(c.Context(), name)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessUserNameExistenceCheckedKey, map[string]bool{"exists": exists})
}

func (h *UserHandler) CheckEmailExists(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return web.HandleError(c, domain.ErrBadRequestWithKey(utils.ErrUserEmailRequiredKey))
	}

	exists, err := h.Service.CheckEmailExists(c.Context(), email)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessUserEmailExistenceCheckedKey, map[string]bool{"exists": exists})
}

func (h *UserHandler) CountUsers(c *fiber.Ctx) error {
	params, err := h.parseUserFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	count, err := h.Service.CountUsers(c.Context(), params)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessUserCountedKey, count)
}

func (h *UserHandler) GetUserStatistics(c *fiber.Ctx) error {
	stats, err := h.Service.GetUserStatistics(c.Context())
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessUserStatisticsRetrievedKey, stats)
}

// *===========================EXPORT===========================*
func (h *UserHandler) ExportUserList(c *fiber.Ctx) error {
	var payload domain.ExportUserListPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	params, err := h.parseUserFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	// * Get language from headers
	langCode := web.GetLanguageFromContext(c)

	fileBytes, filename, err := h.Service.ExportUserList(c.Context(), payload, params, langCode)
	if err != nil {
		return web.HandleError(c, err)
	}

	// Set appropriate content type and headers
	var contentType string
	switch payload.Format {
	case domain.ExportFormatPDF:
		contentType = "application/pdf"
	case domain.ExportFormatExcel:
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", "attachment; filename="+filename)

	return c.Send(fileBytes)
}
