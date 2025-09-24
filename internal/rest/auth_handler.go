// Package rest provides REST API handlers for the inventory management system.
//
//	@title			Inventory Management API
//	@version		1.0
//	@description	A comprehensive inventory management API with JWT authentication, multi-language support, and CRUD operations for assets, users, and locations.
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io
//
//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT
//
//	@host		localhost:8080
//	@BasePath	/api/v1
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
//
//	@tag.name			Authentication
//	@tag.description	Authentication related endpoints
//	@tag.name			Users
//	@tag.description	User management endpoints
//	@tag.name			Categories
//	@tag.description	Category management endpoints
//	@tag.name			Locations
//	@tag.description	Location management endpoints
//	@tag.name			Assets
//	@tag.description	Asset management endpoints
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
	users.Post("/register", handler.Register)
	users.Post("/login", handler.Login)

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
//	@Success		201				{object}	web.JSONResponse{data=domain.User}	"User registered successfully"
//	@Failure		400				{object}	web.JSONResponse{error=[]web.ValidationError}	"Validation failed"
//	@Failure		409				{object}	web.JSONResponse	"User already exists"
//	@Failure		500				{object}	web.JSONResponse	"Internal server error"
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
//	@Success		200				{object}	web.JSONResponse{data=domain.LoginResponse}	"Login successful"
//	@Failure		400				{object}	web.JSONResponse{error=[]web.ValidationError}	"Validation failed"
//	@Failure		401				{object}	web.JSONResponse	"Invalid credentials"
//	@Failure		500				{object}	web.JSONResponse	"Internal server error"
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

// *===========================QUERY===========================*
