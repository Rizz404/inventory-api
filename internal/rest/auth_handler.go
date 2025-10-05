// Package rest provides REST API handlers for authentication endpoints.
package rest

import (
	"context"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/utils"
	"github.com/Rizz404/inventory-api/internal/web"
	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	// * MUTATION
	Register(ctx context.Context, payload *domain.RegisterPayload) (domain.User, error)
	Login(ctx context.Context, payload *domain.LoginPayload) (domain.AuthResponse, error)
	RefreshToken(ctx context.Context, payload *domain.RefreshTokenPayload) (domain.AuthResponse, error)

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
	users.Post("/register", handler.Register)
	users.Post("/login", handler.Login)
	users.Post("/refresh", handler.RefreshToken)

}

// *===========================MUTATION===========================*

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Register a new user account with name, email, and password
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			registerPayload	body		domain.RegisterPayload	true	"User registration data"
//	@Success		201				{object}	web.SuccessResponse{data=domain.UserResponse}	"User registered successfully"
//	@Failure		400				{object}	web.ErrorResponse{error=web.ValidationErrors}	"Validation failed"
//	@Failure		409				{object}	web.ErrorResponse	"User already exists"
//	@Failure		500				{object}	web.ErrorResponse	"Internal server error"
//	@Router			/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {

	var payload domain.RegisterPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	user, err := h.Service.Register(c.Context(), &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusCreated, utils.SuccessUserCreatedKey, user)
}

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate user with email and password, returns JWT tokens
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			loginPayload	body		domain.LoginPayload	true	"User login credentials"
//	@Success		200				{object}	web.SuccessResponse{data=domain.AuthResponse}	"Login successful"
//	@Failure		400				{object}	web.ErrorResponse{error=web.ValidationErrors}	"Validation failed"
//	@Failure		401				{object}	web.ErrorResponse	"Invalid credentials"
//	@Failure		500				{object}	web.ErrorResponse	"Internal server error"
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var payload domain.LoginPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	user, err := h.Service.Login(c.Context(), &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessLoginKey, user)
}

// RefreshToken godoc
//
//	@Summary		Refresh access token
//	@Description	Get new access and refresh tokens using a valid refresh token
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			refreshTokenPayload	body		domain.RefreshTokenPayload	true	"Refresh token data"
//	@Success		200					{object}	web.SuccessResponse{data=domain.AuthResponse}	"Token refreshed successfully"
//	@Failure		400					{object}	web.ErrorResponse{error=web.ValidationErrors}	"Validation failed"
//	@Failure		401					{object}	web.ErrorResponse	"Invalid or expired refresh token"
//	@Failure		500					{object}	web.ErrorResponse	"Internal server error"
//	@Router			/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var payload domain.RefreshTokenPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	authResponse, err := h.Service.RefreshToken(c.Context(), &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessTokenRefreshedKey, authResponse)
}

// *===========================QUERY===========================*
