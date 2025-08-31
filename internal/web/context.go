package web

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/gofiber/fiber/v2"
)

// * GetUserFromContext helper function untuk ambil user info dari context
func GetUserFromContext(c *fiber.Ctx) (idUser string, name string, email string, role domain.UserRole, isActive bool, ok bool) {
	idUser, ok1 := c.Locals("id_user").(string)
	name, ok2 := c.Locals("name").(string)
	email, ok3 := c.Locals("email").(string)
	role, ok4 := c.Locals("role").(domain.UserRole)
	isActive, ok5 := c.Locals("is_active").(bool)

	return idUser, name, email, role, isActive, ok1 && ok2 && ok3 && ok4 && ok5
}

// * GetUserIDFromContext helper function untuk ambil user ID saja
func GetUserIDFromContext(c *fiber.Ctx) (string, bool) {
	idUser, ok := c.Locals("id_user").(string)
	return idUser, ok
}
