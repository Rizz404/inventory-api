package rest

import (
	"strconv"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/rest/middleware"
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
	users.Post("/",
		middleware.AuthMiddleware(),
		middleware.AuthorizeRole(domain.RoleAdmin),
		handler.CreateUser,
	)

	users.Get("/", handler.GetUsersPaginated)
	users.Get("/cursor", handler.GetUsersCursor)
	users.Get("/count", handler.CountUsers)
	users.Get("/profile", middleware.AuthMiddleware(), handler.GetCurrentUser)
	users.Patch("/profile", middleware.AuthMiddleware(), handler.UpdateCurrentUser)
	users.Get("/name/:name", handler.GetUserByName)
	users.Get("/email/:email", handler.GetUserByEmail)
	users.Get("/check/name/:name", handler.CheckNameExists)
	users.Get("/check/email/:email", handler.CheckEmailExists)
	users.Get("/check/:id", handler.CheckUserExists)
	users.Get("/:id", handler.GetUserById)
	users.Patch("/:id", handler.UpdateUser)
	users.Delete("/:id", handler.DeleteUser)
}

func (h *UserHandler) parseUserFiltersAndSort(c *fiber.Ctx) (query.Params, error) {
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
	roleStr := c.Query("role")
	var role *domain.UserRole
	if roleStr != "" {
		r := domain.UserRole(roleStr)
		role = &r
	}

	filters := &postgresql.UserFilterOptions{Role: role}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive, err := strconv.ParseBool(isActiveStr)
		if err == nil {
			filters.IsActive = &isActive
		}
	}

	if employeeID := c.Query("employee_id"); employeeID != "" {
		filters.EmployeeID = &employeeID
	}

	params.Filters = filters

	return params, nil
}

// *===========================MUTATION===========================*
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var payload domain.CreateUserPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	user, err := h.Service.CreateUser(c.Context(), &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusCreated, "User created successfully", user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequest("User ID is required"))
	}

	var payload domain.UpdateUserPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	user, err := h.Service.UpdateUser(c.Context(), id, &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, "User updated successfully", user)
}

func (h *UserHandler) UpdateCurrentUser(c *fiber.Ctx) error {
	id, ok := web.GetUserIDFromContext(c)
	if !ok {
		return web.HandleError(c, domain.ErrBadRequest("User ID is required"))
	}

	var payload domain.UpdateUserPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	user, err := h.Service.UpdateUser(c.Context(), id, &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, "User updated successfully", user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequest("User ID is required"))
	}

	err := h.Service.DeleteUser(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, "User deleted successfully", nil)
}

// *===========================QUERY===========================*
func (h *UserHandler) GetUsersPaginated(c *fiber.Ctx) error {
	params, err := h.parseUserFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	params.Pagination = &query.PaginationOptions{Limit: limit, Offset: offset}

	users, total, err := h.Service.GetUsersPaginated(c.Context(), params)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.SuccessWithPageInfo(c, fiber.StatusOK, "Users retrieved successfully", users, int(total), limit, (offset/limit)+1)
}

func (h *UserHandler) GetUsersCursor(c *fiber.Ctx) error {
	params, err := h.parseUserFiltersAndSort(c)
	if err != nil {
		return web.HandleError(c, domain.ErrBadRequest(err.Error()))
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	cursor := c.Query("cursor")
	params.Pagination = &query.PaginationOptions{Limit: limit, Cursor: cursor}

	users, err := h.Service.GetUsersCursor(c.Context(), params)
	if err != nil {
		return web.HandleError(c, err)
	}

	var nextCursor string
	hasNextPage := len(users) == limit
	if hasNextPage {
		nextCursor = users[len(users)-1].ID
	}

	return web.SuccessWithCursor(c, fiber.StatusOK, "Users retrieved successfully", users, nextCursor, hasNextPage, limit)
}
func (h *UserHandler) GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequest("User ID is required"))
	}

	user, err := h.Service.GetUserById(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, "User retrieved successfully", user)
}

func (h *UserHandler) GetUserByName(c *fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return web.HandleError(c, domain.ErrBadRequest("Name is required"))
	}

	user, err := h.Service.GetUserByName(c.Context(), name)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, "User retrieved successfully by name", user)
}

func (h *UserHandler) GetUserByEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return web.HandleError(c, domain.ErrBadRequest("Email is required"))
	}

	user, err := h.Service.GetUserByEmail(c.Context(), email)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, "User retrieved successfully by email", user)
}

func (h *UserHandler) GetCurrentUser(c *fiber.Ctx) error {
	id, ok := web.GetUserIDFromContext(c)
	if !ok {
		return web.HandleError(c, domain.ErrBadRequest("User ID is required"))
	}

	user, err := h.Service.GetUserById(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, "User retrieved successfully", user)
}

func (h *UserHandler) CheckUserExists(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequest("User ID is required"))
	}

	exists, err := h.Service.CheckUserExists(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, "User existence checked successfully", map[string]bool{"exists": exists})
}

func (h *UserHandler) CheckNameExists(c *fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return web.HandleError(c, domain.ErrBadRequest("Name is required"))
	}

	exists, err := h.Service.CheckNameExists(c.Context(), name)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, "Name existence checked successfully", map[string]bool{"exists": exists})
}

func (h *UserHandler) CheckEmailExists(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return web.HandleError(c, domain.ErrBadRequest("Email is required"))
	}

	exists, err := h.Service.CheckEmailExists(c.Context(), email)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, "Email existence checked successfully", map[string]bool{"exists": exists})
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

	return web.Success(c, fiber.StatusOK, "User counted successfully", count)
}
