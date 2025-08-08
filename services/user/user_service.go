package user

import (
	"context"
	"errors"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/utils"
)

// * Sekarang service hanya mem-proxy error, karena sudah diterjemahkan oleh repo
type Repository interface {
	// * MUTATION
	CreateUser(ctx context.Context, payload *domain.User) (domain.User, error)
	UpdateUser(ctx context.Context, payload *domain.User) (domain.User, error)
	UpdateUserWithPayload(ctx context.Context, userId string, payload *domain.UpdateUserPayload) (domain.User, error)
	DeleteUser(ctx context.Context, userId string) error

	// * QUERY
	GetUsersPaginated(ctx context.Context, params query.Params) ([]domain.User, error)
	GetUsersCursor(ctx context.Context, params query.Params) ([]domain.User, error)
	GetUserById(ctx context.Context, userId string) (domain.User, error)
	GetUserByUsernameOrEmail(ctx context.Context, username string, email string) (domain.User, error)
	CheckUserExist(ctx context.Context, userId string) (bool, error)
	CountUsers(ctx context.Context, params query.Params) (int64, error)
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
func (s *Service) CreateUser(ctx context.Context, payload *domain.CreateUserPayload) (domain.User, error) {
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return domain.User{}, domain.ErrInternal(err)
	}

	// * Cek apakah username sudah ada
	_, err = s.Repo.GetUserByUsernameOrEmail(ctx, payload.Username, "")
	if err == nil {
		// * Jika tidak ada error, berarti user DITEMUKAN. Ini konflik.
		return domain.User{}, domain.ErrConflict("user with this username already exists")
	}

	var appErr *domain.AppError
	if errors.As(err, &appErr) && appErr.Code != 404 {
		// * Jika errornya bukan 404 (NotFound), maka ini adalah error internal.
		return domain.User{}, err
	}

	// Set default language if not provided
	preferredLang := "id-ID"
	if payload.PreferredLang != nil {
		preferredLang = *payload.PreferredLang
	}

	// * Siapkan user baru
	newUser := domain.User{
		Username:      payload.Username,
		PasswordHash:  hashedPassword,
		FullName:      payload.FullName,
		Role:          payload.Role,
		EmployeeID:    payload.EmployeeID,
		PreferredLang: preferredLang,
		IsActive:      true, // Default active
	}

	createdUser, err := s.Repo.CreateUser(ctx, &newUser)
	if err != nil {
		// * Repository sudah menerjemahkan error (misal: conflict), jadi langsung kembalikan
		return domain.User{}, err
	}

	// Clear password hash from response
	createdUser.PasswordHash = ""
	return createdUser, nil
}

func (s *Service) UpdateUser(ctx context.Context, userId string, payload *domain.UpdateUserPayload) (domain.User, error) {
	// Check if user exists
	_, err := s.Repo.GetUserById(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}

	// Check username uniqueness if being updated
	if payload.Username != nil {
		_, err := s.Repo.GetUserByUsernameOrEmail(ctx, *payload.Username, "")
		if err == nil {
			return domain.User{}, domain.ErrConflict("username already taken")
		}
		var appErr *domain.AppError
		if errors.As(err, &appErr) && appErr.Code != 404 {
			return domain.User{}, err
		}
	}

	// Use the new UpdateUserWithPayload method
	updatedUser, err := s.Repo.UpdateUserWithPayload(ctx, userId, payload)
	if err != nil {
		return domain.User{}, err
	}

	// Clear password hash from response
	updatedUser.PasswordHash = ""
	return updatedUser, nil
}

func (s *Service) DeleteUser(ctx context.Context, userId string) error {
	err := s.Repo.DeleteUser(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}

// *===========================QUERY===========================*
func (s *Service) GetUsersPaginated(ctx context.Context, params query.Params) ([]domain.User, int64, error) {
	users, err := s.Repo.GetUsersPaginated(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// * Hanya hitung total jika pagination-nya offset
	count, err := s.Repo.CountUsers(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// Clear password hash from all users
	for i := range users {
		users[i].PasswordHash = ""
	}

	return users, count, nil
}

func (s *Service) GetUsersCursor(ctx context.Context, params query.Params) ([]domain.User, error) {
	users, err := s.Repo.GetUsersCursor(ctx, params)
	if err != nil {
		return nil, err
	}

	// Clear password hash from all users
	for i := range users {
		users[i].PasswordHash = ""
	}

	return users, nil
}

func (s *Service) GetUserById(ctx context.Context, userId string) (domain.User, error) {
	user, err := s.Repo.GetUserById(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}

	// Clear password hash from response
	user.PasswordHash = ""
	return user, nil
}

func (s *Service) GetUserByUsernameOrEmail(ctx context.Context, username string, email string) (domain.User, error) {
	user, err := s.Repo.GetUserByUsernameOrEmail(ctx, username, email)
	if err != nil {
		return domain.User{}, err
	}

	// Clear password hash from response
	user.PasswordHash = ""
	return user, nil
}

func (s *Service) CheckUserExist(ctx context.Context, userId string) (bool, error) {
	exist, err := s.Repo.CheckUserExist(ctx, userId)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (s *Service) CountUsers(ctx context.Context, params query.Params) (int64, error) {
	count, err := s.Repo.CountUsers(ctx, params)
	if err != nil {
		return 0, err
	}
	return count, nil
}
