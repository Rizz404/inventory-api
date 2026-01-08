package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/client/smtp"
	"github.com/Rizz404/inventory-api/internal/utils"
)

type Repository interface {
	// * MUTATION
	CreateUser(ctx context.Context, payload *domain.User) (domain.User, error)
	UpdateUserPassword(ctx context.Context, email, passwordHash string) error
	UpdateLastLogin(ctx context.Context, userId string) error

	// * QUERY
	GetUserById(ctx context.Context, userId string) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	CheckUserExists(ctx context.Context, userId string) (bool, error)
	CheckNameExists(ctx context.Context, name string) (bool, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
}

// resetCodeEntry stores reset code with expiration
type resetCodeEntry struct {
	Code      string
	Email     string
	ExpiresAt time.Time
}

type Service struct {
	Repo       Repository
	SMTPClient *smtp.Client
	resetCodes sync.Map // email -> resetCodeEntry
}

func NewService(r Repository, smtpClient *smtp.Client) *Service {
	return &Service{
		Repo:       r,
		SMTPClient: smtpClient,
	}
}

// *===========================MUTATION===========================*

func (s *Service) Register(ctx context.Context, payload *domain.RegisterPayload) (domain.User, error) {
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return domain.User{}, domain.ErrInternal(err)
	}

	// * Check if name or email already exists
	if nameExists, err := s.Repo.CheckNameExists(ctx, payload.Name); err != nil {
		return domain.User{}, err
	} else if nameExists {
		return domain.User{}, domain.ErrConflictWithKey(utils.ErrUserNameExistsKey)
	}

	if emailExists, err := s.Repo.CheckEmailExists(ctx, payload.Email); err != nil {
		return domain.User{}, err
	} else if emailExists {
		return domain.User{}, domain.ErrConflictWithKey(utils.ErrUserEmailExistsKey)
	}

	// * Siapkan user baru
	newUser := domain.User{
		Name:         payload.Name,
		Email:        payload.Email,
		PasswordHash: hashedPassword,
		Role:         domain.RoleEmployee,
	}

	createdUser, err := s.Repo.CreateUser(ctx, &newUser)
	if err != nil {
		return domain.User{}, err
	}

	return createdUser, nil
}

func (s *Service) Login(ctx context.Context, payload *domain.LoginPayload) (domain.AuthResponse, error) {
	// Search by email
	user, err := s.Repo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		return domain.AuthResponse{}, domain.ErrUnauthorizedWithKey(utils.ErrInvalidCredentialsKey)
	}

	// Check if user is active
	if !user.IsActive {
		return domain.AuthResponse{}, domain.ErrUnauthorizedWithKey(utils.ErrUserNotFoundKey) // * Use generic message for security
	}

	// Verify password
	passwordIsValid := utils.CheckPasswordHash(payload.Password, user.PasswordHash)
	if !passwordIsValid {
		return domain.AuthResponse{}, domain.ErrUnauthorizedWithKey(utils.ErrInvalidCredentialsKey)
	}

	// Create JWT payload with available fields
	jwtPayload := &utils.CreateJWTPayload{
		IDUser:   user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Role:     string(user.Role), // Convert UserRole to string
		IsActive: user.IsActive,
	}

	accessToken, err := utils.CreateAccessToken(jwtPayload)
	if err != nil {
		return domain.AuthResponse{}, domain.ErrInternal(err)
	}

	refreshToken, err := utils.CreateRefreshToken(jwtPayload.IDUser)
	if err != nil {
		return domain.AuthResponse{}, domain.ErrInternal(err)
	}

	// Update last login timestamp (fire and forget, don't block login on failure)
	go func() {
		_ = s.Repo.UpdateLastLogin(context.Background(), user.ID)
	}()

	// Get current time for last login in response
	now := time.Now()

	// Create user response
	userResponse := domain.UserResponse{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		FullName:      user.FullName,
		Role:          user.Role,
		EmployeeID:    user.EmployeeID,
		PreferredLang: user.PreferredLang,
		IsActive:      user.IsActive,
		AvatarURL:     user.AvatarURL,
		PhoneNumber:   user.PhoneNumber,
		FCMToken:      user.FCMToken,
		LastLogin:     &now,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}

	authResponse := domain.AuthResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return authResponse, nil
}

func (s *Service) RefreshToken(ctx context.Context, payload *domain.RefreshTokenPayload) (domain.AuthResponse, error) {
	// Validate refresh token
	claims, err := utils.ValidateRefreshToken(payload.RefreshToken)
	if err != nil {
		return domain.AuthResponse{}, domain.ErrUnauthorizedWithKey(utils.ErrTokenInvalidKey)
	}

	// Check if user still exists and is active
	exists, err := s.Repo.CheckUserExists(ctx, claims.IDUser)
	if err != nil {
		return domain.AuthResponse{}, err
	}
	if !exists {
		return domain.AuthResponse{}, domain.ErrUnauthorizedWithKey(utils.ErrUserNotFoundKey)
	}

	// Get fresh user data by ID (refresh token only contains user ID)
	user, err := s.Repo.GetUserById(ctx, claims.IDUser)
	if err != nil {
		return domain.AuthResponse{}, domain.ErrUnauthorizedWithKey(utils.ErrUserNotFoundKey)
	}

	// Check if user is still active
	if !user.IsActive {
		return domain.AuthResponse{}, domain.ErrUnauthorizedWithKey(utils.ErrUserNotFoundKey)
	}

	// Create new tokens with fresh user data
	jwtPayload := &utils.CreateJWTPayload{
		IDUser:   user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Role:     string(user.Role),
		IsActive: user.IsActive,
	}

	accessToken, err := utils.CreateAccessToken(jwtPayload)
	if err != nil {
		return domain.AuthResponse{}, domain.ErrInternal(err)
	}

	refreshToken, err := utils.CreateRefreshToken(user.ID)
	if err != nil {
		return domain.AuthResponse{}, domain.ErrInternal(err)
	}

	// Create user response
	userResponse := domain.UserResponse{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		FullName:      user.FullName,
		Role:          user.Role,
		EmployeeID:    user.EmployeeID,
		PreferredLang: user.PreferredLang,
		IsActive:      user.IsActive,
		AvatarURL:     user.AvatarURL,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}

	authResponse := domain.AuthResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return authResponse, nil
}

// ForgotPassword generates and sends a password reset code to the user's email
func (s *Service) ForgotPassword(ctx context.Context, payload *domain.ForgotPasswordPayload) (string, error) {
	// Check if SMTP is enabled
	if s.SMTPClient == nil || !s.SMTPClient.IsEnabled() {
		return "", domain.ErrInternalWithMessage("Email service is not available")
	}

	// Check if user exists
	user, err := s.Repo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		// Return success even if email doesn't exist (security: don't reveal if email exists)
		return "If the email exists, a reset code will be sent", nil
	}

	// Generate 6-digit code
	code, err := s.generateResetCode()
	if err != nil {
		return "", domain.ErrInternal(err)
	}

	// Store code with 15 minutes expiration
	s.resetCodes.Store(payload.Email, resetCodeEntry{
		Code:      code,
		Email:     payload.Email,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	})

	// Send email
	userName := user.Name
	if user.FullName != "" {
		userName = user.FullName
	}

	if err := s.SMTPClient.SendPasswordResetEmail(ctx, payload.Email, code, userName); err != nil {
		return "", domain.ErrInternalWithKey(utils.ErrEmailSendFailedKey)
	}

	return "Reset code sent to your email", nil
}

// VerifyResetCode verifies the password reset code
func (s *Service) VerifyResetCode(ctx context.Context, payload *domain.VerifyResetCodePayload) (domain.VerifyResetCodeResponse, error) {
	entry, ok := s.resetCodes.Load(payload.Email)
	if !ok {
		return domain.VerifyResetCodeResponse{Valid: false}, nil
	}

	resetEntry := entry.(resetCodeEntry)

	// Check if code expired
	if time.Now().After(resetEntry.ExpiresAt) {
		s.resetCodes.Delete(payload.Email)
		return domain.VerifyResetCodeResponse{Valid: false}, nil
	}

	// Check if code matches
	if resetEntry.Code != payload.Code {
		return domain.VerifyResetCodeResponse{Valid: false}, nil
	}

	return domain.VerifyResetCodeResponse{Valid: true}, nil
}

// ResetPassword resets the user's password using the verification code
func (s *Service) ResetPassword(ctx context.Context, payload *domain.ResetPasswordPayload) (string, error) {
	// Verify code first
	entry, ok := s.resetCodes.Load(payload.Email)
	if !ok {
		return "", domain.ErrBadRequestWithKey(utils.ErrResetCodeNotFoundKey)
	}

	resetEntry := entry.(resetCodeEntry)

	// Check if code expired
	if time.Now().After(resetEntry.ExpiresAt) {
		s.resetCodes.Delete(payload.Email)
		return "", domain.ErrBadRequestWithKey(utils.ErrResetCodeExpiredKey)
	}

	// Check if code matches
	if resetEntry.Code != payload.Code {
		return "", domain.ErrBadRequestWithKey(utils.ErrResetCodeInvalidKey)
	}

	// Check if user exists
	_, err := s.Repo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		return "", domain.ErrNotFoundWithKey(utils.ErrUserNotFoundKey)
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(payload.NewPassword)
	if err != nil {
		return "", domain.ErrInternal(err)
	}

	// Update password
	if err := s.Repo.UpdateUserPassword(ctx, payload.Email, hashedPassword); err != nil {
		return "", err
	}

	// Delete used reset code
	s.resetCodes.Delete(payload.Email)

	return "Password reset successfully", nil
}

// generateResetCode generates a random 6-digit code
func (s *Service) generateResetCode() (string, error) {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// Convert to 6-digit number
	code := int(b[0])*10000 + int(b[1])*100 + int(b[2])
	code = code % 1000000
	return fmt.Sprintf("%06d", code), nil
}

// *===========================QUERY===========================*
