package auth

import (
	"context"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/utils"
)

type Repository interface {
	// * MUTATION
	CreateUser(ctx context.Context, payload *domain.User) (domain.User, error)

	// * QUERY
	GetUserById(ctx context.Context, userId string) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	CheckUserExists(ctx context.Context, userId string) (bool, error)
	CheckNameExists(ctx context.Context, name string) (bool, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
}

type Service struct {
	Repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		Repo: r,
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

// *===========================QUERY===========================*
