package user

import (
	"context"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/Rizz404/inventory-api/internal/utils"
)

// * Repository interface defines the contract for user data operations
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
	GetUserByName(ctx context.Context, name string) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	CheckUserExists(ctx context.Context, userId string) (bool, error)
	CheckNameExists(ctx context.Context, name string) (bool, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CountUsers(ctx context.Context, params query.Params) (int64, error)
}

// * UserService interface defines the contract for user business operations
type UserService interface {
	// * MUTATION
	CreateUser(ctx context.Context, payload *domain.CreateUserPayload) (domain.UserResponse, error)
	UpdateUser(ctx context.Context, userId string, payload *domain.UpdateUserPayload) (domain.UserResponse, error)
	DeleteUser(ctx context.Context, userId string) error

	// * QUERY
	GetUsersPaginated(ctx context.Context, params query.Params) ([]domain.UserResponse, int64, error)
	GetUsersCursor(ctx context.Context, params query.Params) ([]domain.UserResponse, error)
	GetUserById(ctx context.Context, userId string) (domain.UserResponse, error)
	GetUserByName(ctx context.Context, name string) (domain.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (domain.UserResponse, error)
	CheckUserExists(ctx context.Context, userId string) (bool, error)
	CheckNameExists(ctx context.Context, name string) (bool, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CountUsers(ctx context.Context, params query.Params) (int64, error)
}

type Service struct {
	Repo Repository
}

// * Ensure Service implements UserService interface
var _ UserService = (*Service)(nil)

func NewService(r Repository) UserService {
	return &Service{
		Repo: r,
	}
}

// *===========================MUTATION===========================*
func (s *Service) CreateUser(ctx context.Context, payload *domain.CreateUserPayload) (domain.UserResponse, error) {
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return domain.UserResponse{}, domain.ErrInternal(err)
	}

	// * Check if name or email already exists
	if nameExists, err := s.Repo.CheckNameExists(ctx, payload.Name); err != nil {
		return domain.UserResponse{}, err
	} else if nameExists {
		return domain.UserResponse{}, domain.ErrConflict("user with name '" + payload.Name + "' already exists")
	}

	if emailExists, err := s.Repo.CheckEmailExists(ctx, payload.Email); err != nil {
		return domain.UserResponse{}, err
	} else if emailExists {
		return domain.UserResponse{}, domain.ErrConflict("user with email '" + payload.Email + "' already exists")
	}

	// Set default language if not provided
	preferredLang := "id-ID"
	if payload.PreferredLang != nil {
		preferredLang = *payload.PreferredLang
	}

	// * Siapkan user baru
	newUser := domain.User{
		Name:          payload.Name,
		Email:         payload.Email,
		PasswordHash:  hashedPassword,
		FullName:      payload.FullName,
		Role:          payload.Role,
		EmployeeID:    payload.EmployeeID,
		PreferredLang: preferredLang,
		IsActive:      true, // Default active
		AvatarURL:     payload.AvatarURL,
	}

	createdUser, err := s.Repo.CreateUser(ctx, &newUser)
	if err != nil {
		// * Repository sudah menerjemahkan error (misal: conflict), jadi langsung kembalikan
		return domain.UserResponse{}, err
	}

	// * Convert to UserResponse using direct mapper
	return mapper.DomainUserToUserResponse(&createdUser), nil
}

func (s *Service) UpdateUser(ctx context.Context, userId string, payload *domain.UpdateUserPayload) (domain.UserResponse, error) {
	// Check if user exists
	_, err := s.Repo.GetUserById(ctx, userId)
	if err != nil {
		return domain.UserResponse{}, err
	}

	// * Check name/email uniqueness if being updated
	if payload.Name != nil {
		if nameExists, err := s.Repo.CheckNameExists(ctx, *payload.Name); err != nil {
			return domain.UserResponse{}, err
		} else if nameExists {
			return domain.UserResponse{}, domain.ErrConflict("name '" + *payload.Name + "' is already taken")
		}
	}

	if payload.Email != nil {
		if emailExists, err := s.Repo.CheckEmailExists(ctx, *payload.Email); err != nil {
			return domain.UserResponse{}, err
		} else if emailExists {
			return domain.UserResponse{}, domain.ErrConflict("email '" + *payload.Email + "' is already taken")
		}
	}

	// Use the new UpdateUserWithPayload method
	updatedUser, err := s.Repo.UpdateUserWithPayload(ctx, userId, payload)
	if err != nil {
		return domain.UserResponse{}, err
	}

	// * Convert to UserResponse using direct mapper
	return mapper.DomainUserToUserResponse(&updatedUser), nil
}

func (s *Service) DeleteUser(ctx context.Context, userId string) error {
	err := s.Repo.DeleteUser(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}

// *===========================QUERY===========================*
func (s *Service) GetUsersPaginated(ctx context.Context, params query.Params) ([]domain.UserResponse, int64, error) {
	users, err := s.Repo.GetUsersPaginated(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// * Hanya hitung total jika pagination-nya offset
	count, err := s.Repo.CountUsers(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// * Convert to UserResponse using direct mapper
	userResponses := mapper.DomainUsersToUsersResponse(users)

	return userResponses, count, nil
}

func (s *Service) GetUsersCursor(ctx context.Context, params query.Params) ([]domain.UserResponse, error) {
	users, err := s.Repo.GetUsersCursor(ctx, params)
	if err != nil {
		return nil, err
	}

	// * Convert to UserResponse using direct mapper
	userResponses := mapper.DomainUsersToUsersResponse(users)

	return userResponses, nil
}

func (s *Service) GetUserById(ctx context.Context, userId string) (domain.UserResponse, error) {
	user, err := s.Repo.GetUserById(ctx, userId)
	if err != nil {
		return domain.UserResponse{}, err
	}

	// * Convert to UserResponse using direct mapper
	return mapper.DomainUserToUserResponse(&user), nil
}

func (s *Service) GetUserByName(ctx context.Context, name string) (domain.UserResponse, error) {
	user, err := s.Repo.GetUserByName(ctx, name)
	if err != nil {
		return domain.UserResponse{}, err
	}

	// * Convert to UserResponse using direct mapper
	return mapper.DomainUserToUserResponse(&user), nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (domain.UserResponse, error) {
	user, err := s.Repo.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.UserResponse{}, err
	}

	// * Convert to UserResponse using direct mapper
	return mapper.DomainUserToUserResponse(&user), nil
}

func (s *Service) CheckUserExists(ctx context.Context, userId string) (bool, error) {
	exists, err := s.Repo.CheckUserExists(ctx, userId)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Service) CheckNameExists(ctx context.Context, name string) (bool, error) {
	exists, err := s.Repo.CheckNameExists(ctx, name)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Service) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	exists, err := s.Repo.CheckEmailExists(ctx, email)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Service) CountUsers(ctx context.Context, params query.Params) (int64, error) {
	count, err := s.Repo.CountUsers(ctx, params)
	if err != nil {
		return 0, err
	}
	return count, nil
}
