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
	ForgotPassword(ctx context.Context, payload *domain.ForgotPasswordPayload) (domain.ForgotPasswordResponse, error)
	VerifyResetCode(ctx context.Context, payload *domain.VerifyResetCodePayload) (domain.VerifyResetCodeResponse, error)
	ResetPassword(ctx context.Context, payload *domain.ResetPasswordPayload) (domain.ResetPasswordResponse, error)

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
	users.Post("/forgot-password", handler.ForgotPassword)
	users.Post("/verify-reset-code", handler.VerifyResetCode)
	users.Post("/reset-password", handler.ResetPassword)

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

// ForgotPassword godoc
//
//	@Summary		Request password reset
//	@Description	Send a password reset code to the user's email
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			forgotPasswordPayload	body		domain.ForgotPasswordPayload	true	"Email for password reset"
//	@Success		200						{object}	web.SuccessResponse{data=domain.ForgotPasswordResponse}	"Reset code sent"
//	@Failure		400						{object}	web.ErrorResponse{error=web.ValidationErrors}	"Validation failed"
//	@Failure		500						{object}	web.ErrorResponse	"Internal server error"
//	@Router			/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var payload domain.ForgotPasswordPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	response, err := h.Service.ForgotPassword(c.Context(), &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessResetCodeSentKey, response)
}

// VerifyResetCode godoc
//
//	@Summary		Verify password reset code
//	@Description	Verify if the password reset code is valid
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			verifyResetCodePayload	body		domain.VerifyResetCodePayload	true	"Email and reset code"
//	@Success		200						{object}	web.SuccessResponse{data=domain.VerifyResetCodeResponse}	"Code verification result"
//	@Failure		400						{object}	web.ErrorResponse{error=web.ValidationErrors}	"Validation failed"
//	@Failure		500						{object}	web.ErrorResponse	"Internal server error"
//	@Router			/auth/verify-reset-code [post]
func (h *AuthHandler) VerifyResetCode(c *fiber.Ctx) error {
	var payload domain.VerifyResetCodePayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	response, err := h.Service.VerifyResetCode(c.Context(), &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessResetCodeVerifiedKey, response)
}

// ResetPassword godoc
//
//	@Summary		Reset password
//	@Description	Reset user password using the verification code
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			resetPasswordPayload	body		domain.ResetPasswordPayload	true	"Email, code, and new password"
//	@Success		200						{object}	web.SuccessResponse{data=domain.ResetPasswordResponse}	"Password reset successfully"
//	@Failure		400						{object}	web.ErrorResponse{error=web.ValidationErrors}	"Validation failed or invalid code"
//	@Failure		404						{object}	web.ErrorResponse	"User not found"
//	@Failure		500						{object}	web.ErrorResponse	"Internal server error"
//	@Router			/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var payload domain.ResetPasswordPayload
	if err := web.ParseAndValidate(c, &payload); err != nil {
		return web.HandleError(c, err)
	}

	response, err := h.Service.ResetPassword(c.Context(), &payload)
	if err != nil {
		return web.HandleError(c, err)
	}

	return web.Success(c, fiber.StatusOK, utils.SuccessPasswordResetKey, response)
}

// *===========================QUERY===========================*
