package middleware

import (
	"os"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/utils"
	"github.com/Rizz404/inventory-api/internal/web"
	"github.com/gofiber/fiber/v2"
)

const (
	APIKeyHeader = "X-API-Key"
)

// APIKeyMiddleware validates API key for mobile client access
func APIKeyMiddleware() fiber.Handler {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		panic("API_KEY environment variable not set")
	}

	return func(c *fiber.Ctx) error {
		clientKey := c.Get(APIKeyHeader)

		if clientKey == "" {
			return web.HandleError(c, domain.ErrUnauthorizedWithKey(utils.ErrAPIKeyMissingKey))
		}

		if clientKey != apiKey {
			return web.HandleError(c, domain.ErrUnauthorizedWithKey(utils.ErrAPIKeyInvalidKey))
		}

		return c.Next()
	}
}
