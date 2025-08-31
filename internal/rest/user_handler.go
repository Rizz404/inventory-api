package rest

import (
	"context"
	"strconv"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/rest/middleware"
	"github.com/Rizz404/inventory-api/internal/web"
	"github.com/gofiber/fiber/v2"
)

type UserService interface {
	// * MUTATION
	CreateUser(ctx context.Context, payload *domain.CreateUserPayload) (domain.User, error)
	UpdateUser(ctx context.Context, userId string, payload *domain.UpdateUserPayload) (domain.User, error)
	DeleteUser(ctx context.Context, userId string) error

	// * QUERY
	GetUsersPaginated(ctx context.Context, params query.Params) ([]domain.User, int64, error)
	GetUsersCursor(ctx context.Context, params query.Params) ([]domain.User, error)
	GetUserById(ctx context.Context, userId string) (domain.User, error)
	GetUserByNameOrEmail(ctx context.Context, name string, email string) (domain.User, error)
	CheckUserExist(ctx context.Context, userId string) (bool, error)
	CountUsers(ctx context.Context, params query.Params) (int64, error)
}

type UserHandler struct {
	Service UserService
}

func NewUserHandler(app fiber.Router, s UserService) {
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
	users.Get("/check/:id", handler.CheckUserExist)
	users.Get("/:id", handler.GetUserById)
	users.Patch("/:id", handler.UpdateUser)
	users.Delete("/:id", handler.DeleteUser)
}

func (h *UserHandler) parseUserFiltersAndSort(c *fiber.Ctx) (query.Params, error) {
	params := query.Params{}

	search := c.Query("search")
	if search != "" {
		params.SearchQuery = &search
	}

	sortBy := c.Query("sort_by")
	if sortBy != "" {
		params.Sort = &query.SortOptions{
			Field: sortBy,
			Order: c.Query("sort_order", "desc"),
		}
	}

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

func (h *UserHandler) CheckUserExist(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return web.HandleError(c, domain.ErrBadRequest("User ID is required"))
	}

	user, err := h.Service.CheckUserExist(c.Context(), id)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, "User retrieved successfully", user)
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
