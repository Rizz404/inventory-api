package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/Rizz404/inventory-api/internal/utils"
	"github.com/Rizz404/inventory-api/internal/web"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

var accessTokenSecret = []byte(os.Getenv("JWT_ACCESS_SECRET"))

func AuthMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: accessTokenSecret,
		},
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Add debugging
			fmt.Printf("JWT Error: %v\n", err)
			fmt.Printf("Authorization Header: %s\n", c.Get("Authorization"))
			return web.Error(c, fiber.StatusUnauthorized, "invalid or missing JWT token", err)
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			auth := c.Get("Authorization")
			if auth == "" {
				return web.Error(c, fiber.StatusUnauthorized, "missing Authorization header", nil)
			}

			tokenString := strings.TrimPrefix(auth, "Bearer ")
			if tokenString == auth {
				return web.Error(c, fiber.StatusUnauthorized, "invalid token format", nil)
			}

			claims, err := utils.ValidateToken(tokenString, accessTokenSecret)
			if err != nil {
				return web.Error(c, fiber.StatusUnauthorized, "invalid token", nil)
			}

			// Set user info ke context untuk digunakan di handler selanjutnya
			c.Locals("id_user", claims.IDUser)
			if claims.Username != nil {
				c.Locals("username", *claims.Username)
			}

			if claims.Email != nil {
				c.Locals("email", *claims.Email)
			}

			if claims.Role != nil {
				c.Locals("role", *claims.Role)
			}

			if claims.IsActive != nil {
				c.Locals("is_active", *claims.IsActive)
			}

			return c.Next()
		},
	})

}
