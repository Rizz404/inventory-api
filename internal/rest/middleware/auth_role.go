package middleware

import (
	"slices"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/web"
	"github.com/gofiber/fiber/v2"
)

func AuthorizeRole(allowedRoles ...domain.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(domain.UserRole)
		if !ok {
			return web.Error(c, fiber.StatusUnauthorized, "user role not found", nil)
		}

		if slices.Contains(allowedRoles, userRole) {
			return c.Next()
		}

		return web.Error(c, fiber.StatusForbidden, "you don't have permission to access this resource", nil)
	}
}
