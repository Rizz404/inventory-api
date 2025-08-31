package web

import (
	"strings"

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

// * GetLanguageFromContext helper function untuk ambil bahasa dari header Accept-Language
func GetLanguageFromContext(c *fiber.Ctx) string {
	// * Check Accept-Language header
	acceptLang := c.Get("Accept-Language")
	if acceptLang != "" {
		// * Parse Accept-Language header (e.g., "en-US,en;q=0.9,id;q=0.8")
		languages := strings.Split(acceptLang, ",")
		if len(languages) > 0 {
			// * Get first language and remove quality factor
			firstLang := strings.TrimSpace(languages[0])
			if idx := strings.Index(firstLang, ";"); idx != -1 {
				firstLang = firstLang[:idx]
			}
			// * Normalize language code
			firstLang = strings.ToLower(firstLang)
			if strings.HasPrefix(firstLang, "en") {
				return "en-US"
			} else if strings.HasPrefix(firstLang, "id") {
				return "id-ID"
			} else if strings.HasPrefix(firstLang, "ja") {
				return "ja-JP"
			}
			return firstLang
		}
	}

	// * Check X-Language header as fallback
	xLang := c.Get("X-Language")
	if xLang != "" {
		return xLang
	}

	// * Default language
	return "en-US"
}
