package middleware

import (
	"os"
	"strings"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// * Ambil token dari Authorization header
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Next()
		}

		// * Parse Bearer token
		tokenString := strings.TrimPrefix(auth, "Bearer ")
		if tokenString == auth {
			return c.Next()
		}

		// * Parse JWT token
		claims, err := utils.ValidateToken(tokenString, []byte(os.Getenv("JWT_ACCESS_SECRET")))
		if err != nil {
			return c.Next()
		}

		// * Set user info jika token valid
		c.Locals("id_user", claims.IDUser)
		if claims.Name != nil {
			c.Locals("name", *claims.Name)
		}
		if claims.Email != nil {
			c.Locals("email", *claims.Email)
		}
		if claims.Role != nil {
			c.Locals("role", domain.UserRole(*claims.Role))
		}
		if claims.IsActive != nil {
			c.Locals("is_active", *claims.IsActive)
		}

		return c.Next()
	}
}
