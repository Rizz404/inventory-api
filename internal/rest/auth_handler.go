package rest

import (
	"context"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/web"
	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	// * MUTATION
	// Register(ctx context.Context, payload *domain.RegisterPayload) (domain.User, error)
	Login(ctx context.Context, payload *domain.LoginPayload) (domain.LoginResponse, error)

	// * QUERY
}

type AuthHandler struct {
	Service AuthService
}

func NewAuthHandler(app fiber.Router, s AuthService) {
	handler := &AuthHandler{
		Service: s,
	}

	// * Bisa di group
	// ! routenya bisa tabrakan hati-hati
	users := app.Group("/auth")

	// * Create
	// users.Post("/register", handler.Register)
	users.Post("/login", handler.Login)

}

// *===========================MUTATION===========================*
// func (h *AuthHandler) Register(c *fiber.Ctx) error {

// 	var payload domain.RegisterPayload
// 	if err := web.ParseAndValidate(c, &payload); err != nil {
// 		return web.HandleError(c, err)
// 	}

// 	user, err := h.Service.Register(c.Context(), &payload)
// 	if err != nil {
// 		return web.HandleError(c, err)
// 	}

// 	return web.Success(c, fiber.StatusCreated, "register successfully", user)
// }

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var payload domain.LoginPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	user, err := h.Service.Login(c.Context(), &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, "login successfully", user)
}

// *===========================QUERY===========================*
