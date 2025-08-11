package auth

import (
	"context"
	"errors"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/utils"
)

type Repository interface {
	// * MUTATION
	CreateUser(ctx context.Context, payload *domain.User) (domain.User, error)
	UpdateUser(ctx context.Context, payload *domain.User) (domain.User, error)

	// * QUERY
	GetUserByUsernameOrEmail(ctx context.Context, username string, email string) (domain.User, error)
	CheckUserExist(ctx context.Context, userId string) (bool, error)
	CheckUsernameExist(ctx context.Context, username string) (bool, error)
	CheckEmailExist(ctx context.Context, email string) (bool, error)
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

	// * Cek apakah username atau email sudah ada
	_, err = s.Repo.GetUserByUsernameOrEmail(ctx, payload.Username, payload.Email)
	if err == nil {
		// * Jika tidak ada error, berarti user DITEMUKAN. Ini konflik.
		return domain.User{}, domain.ErrConflict("user with this username or email already exists")
	}

	var appErr *domain.AppError
	if errors.As(err, &appErr) && appErr.Code != 404 {
		// * Jika errornya bukan 404 (NotFound), maka ini adalah error internal.
		return domain.User{}, err
	}

	// * Siapkan user baru
	newUser := domain.User{
		Username:     payload.Username,
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

func (s *Service) Login(ctx context.Context, payload *domain.LoginPayload) (domain.LoginResponse, error) {
	// Search by email
	user, err := s.Repo.GetUserByUsernameOrEmail(ctx, "", payload.Email)
	if err != nil {
		return domain.LoginResponse{}, domain.ErrUnauthorized("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return domain.LoginResponse{}, domain.ErrUnauthorized("user account is inactive")
	}

	// Verify password
	passwordIsValid := utils.CheckPasswordHash(payload.Password, user.PasswordHash)
	if !passwordIsValid {
		return domain.LoginResponse{}, domain.ErrUnauthorized("invalid email or password")
	}

	// Create JWT payload with available fields
	jwtPayload := &utils.CreateJWTPayload{
		IDUser:   user.ID,
		Username: user.Username,
		Role:     string(user.Role), // Convert UserRole to string
		IsActive: user.IsActive,
	}

	accessToken, err := utils.CreateAccessToken(jwtPayload)
	if err != nil {
		return domain.LoginResponse{}, domain.ErrInternal(err)
	}

	refreshToken, err := utils.CreateRefreshToken(jwtPayload.IDUser)
	if err != nil {
		return domain.LoginResponse{}, domain.ErrInternal(err)
	}

	// Create user response
	userResponse := domain.UserResponse{
		ID:            user.ID,
		Username:      user.Username,
		FullName:      user.FullName,
		Role:          user.Role,
		EmployeeID:    user.EmployeeID,
		PreferredLang: user.PreferredLang,
		IsActive:      user.IsActive,
		CreatedAt:     user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	loginResponse := domain.LoginResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return loginResponse, nil
}

// *===========================QUERY===========================*
