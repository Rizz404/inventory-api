package middleware

import (
	"os"
	"strings"

	"github.com/Rizz404/inventory-api/domain"
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
			return web.HandleError(c, domain.ErrUnauthorizedWithKey(utils.ErrTokenInvalidKey))
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			auth := c.Get("Authorization")
			if auth == "" {
				return web.HandleError(c, domain.ErrUnauthorizedWithKey(utils.ErrUnauthorizedKey))
			}

			tokenString := strings.TrimPrefix(auth, "Bearer ")
			if tokenString == auth {
				return web.HandleError(c, domain.ErrUnauthorizedWithKey(utils.ErrTokenInvalidKey))
			}

			claims, err := utils.ValidateToken(tokenString, accessTokenSecret)
			if err != nil {
				return web.HandleError(c, domain.ErrUnauthorizedWithKey(utils.ErrTokenInvalidKey))
			}

			// Set user info ke context untuk digunakan di handler selanjutnya
			c.Locals("id_user", claims.IDUser)
			if claims.Name != nil {
				c.Locals("name", *claims.Name)
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
