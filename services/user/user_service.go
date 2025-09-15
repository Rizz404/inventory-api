package user

import (
	"context"
	"mime/multipart"
	"strings"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/client/cloudinary"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/Rizz404/inventory-api/internal/utils"
	"github.com/oklog/ulid/v2"
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
	CheckNameExistsExcluding(ctx context.Context, name string, excludeUserId string) (bool, error)
	CheckEmailExistsExcluding(ctx context.Context, email string, excludeUserId string) (bool, error)
	CountUsers(ctx context.Context, params query.Params) (int64, error)
	GetUserStatistics(ctx context.Context) (domain.UserStatistics, error)
}

// * UserService interface defines the contract for user business operations
type UserService interface {
	// * MUTATION
	CreateUser(ctx context.Context, payload *domain.CreateUserPayload, avatarFile *multipart.FileHeader) (domain.UserResponse, error)
	UpdateUser(ctx context.Context, userId string, payload *domain.UpdateUserPayload, avatarFile *multipart.FileHeader) (domain.UserResponse, error)
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
	GetUserStatistics(ctx context.Context) (domain.UserStatistics, error)
}

type Service struct {
	Repo             Repository
	CloudinaryClient *cloudinary.Client
}

// * Ensure Service implements UserService interface
var _ UserService = (*Service)(nil)

func NewService(r Repository, cloudinaryClient *cloudinary.Client) UserService {
	return &Service{
		Repo:             r,
		CloudinaryClient: cloudinaryClient,
	}
}

// *===========================MUTATION===========================*
func (s *Service) CreateUser(ctx context.Context, payload *domain.CreateUserPayload, avatarFile *multipart.FileHeader) (domain.UserResponse, error) {
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return domain.UserResponse{}, domain.ErrInternal(err)
	}

	// * Check if name or email already exists
	if nameExists, err := s.Repo.CheckNameExists(ctx, payload.Name); err != nil {
		return domain.UserResponse{}, err
	} else if nameExists {
		return domain.UserResponse{}, domain.ErrConflictWithKey(utils.ErrUserNameExistsKey)
	}

	if emailExists, err := s.Repo.CheckEmailExists(ctx, payload.Email); err != nil {
		return domain.UserResponse{}, err
	} else if emailExists {
		return domain.UserResponse{}, domain.ErrConflictWithKey(utils.ErrUserEmailExistsKey)
	}

	// Set default language if not provided
	preferredLang := "id-ID"
	if payload.PreferredLang != nil {
		preferredLang = *payload.PreferredLang
	}

	// * Handle avatar upload if file is provided
	var avatarURL *string
	if avatarFile != nil {
		// Upload file to Cloudinary if client is available
		if s.CloudinaryClient != nil {
			// Generate temporary user ID for avatar naming
			tempUserID := "temp_" + ulid.Make().String()
			uploadConfig := cloudinary.GetAvatarUploadConfig()
			publicID := "user_" + tempUserID + "_avatar"
			uploadConfig.PublicID = &publicID

			uploadResult, err := s.CloudinaryClient.UploadSingleFile(ctx, avatarFile, uploadConfig)
			if err != nil {
				return domain.UserResponse{}, domain.ErrBadRequestWithKey(utils.ErrFileUploadFailedKey)
			}
			avatarURL = &uploadResult.SecureURL
		} else {
			return domain.UserResponse{}, domain.ErrBadRequestWithKey(utils.ErrCloudinaryConfigKey)
		}
	} else if payload.AvatarURL != nil {
		// Use provided avatar URL from JSON/form data
		avatarURL = payload.AvatarURL
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
		AvatarURL:     avatarURL,
	}

	createdUser, err := s.Repo.CreateUser(ctx, &newUser)
	if err != nil {
		// * Repository sudah menerjemahkan error (misal: conflict), jadi langsung kembalikan
		return domain.UserResponse{}, err
	}

	// * Update avatar public ID with actual user ID if file was uploaded
	if avatarFile != nil && s.CloudinaryClient != nil && avatarURL != nil {
		// Re-upload with correct public ID
		uploadConfig := cloudinary.GetAvatarUploadConfig()
		finalPublicID := "user_" + createdUser.ID + "_avatar"
		uploadConfig.PublicID = &finalPublicID

		uploadResult, err := s.CloudinaryClient.UploadSingleFile(ctx, avatarFile, uploadConfig)
		if err == nil {
			// Update user with final avatar URL
			updatePayload := &domain.UpdateUserPayload{
				AvatarURL: &uploadResult.SecureURL,
			}
			createdUser, _ = s.Repo.UpdateUserWithPayload(ctx, createdUser.ID, updatePayload)
		}
		// Note: We don't return error here to avoid failing user creation if avatar re-upload fails
	}

	// * Convert to UserResponse using direct mapper
	return mapper.DomainUserToUserResponse(&createdUser), nil
}

func (s *Service) UpdateUser(ctx context.Context, userId string, payload *domain.UpdateUserPayload, avatarFile *multipart.FileHeader) (domain.UserResponse, error) {
	// Check if user exists
	existingUser, err := s.Repo.GetUserById(ctx, userId)
	if err != nil {
		return domain.UserResponse{}, err
	}

	// * Check name/email uniqueness if being updated
	if payload.Name != nil {
		if nameExists, err := s.Repo.CheckNameExistsExcluding(ctx, *payload.Name, userId); err != nil {
			return domain.UserResponse{}, err
		} else if nameExists {
			return domain.UserResponse{}, domain.ErrConflictWithKey(utils.ErrUserNameExistsKey)
		}
	}

	if payload.Email != nil {
		if emailExists, err := s.Repo.CheckEmailExistsExcluding(ctx, *payload.Email, userId); err != nil {
			return domain.UserResponse{}, err
		} else if emailExists {
			return domain.UserResponse{}, domain.ErrConflictWithKey(utils.ErrUserEmailExistsKey)
		}
	}

	// * Handle avatar update
	var shouldDeleteOldAvatar bool
	oldAvatarPublicID := "user_" + userId + "_avatar"

	if avatarFile != nil {
		// Upload new avatar file
		if s.CloudinaryClient != nil {
			uploadConfig := cloudinary.GetAvatarUploadConfig()
			publicID := "user_" + userId + "_avatar"
			uploadConfig.PublicID = &publicID

			uploadResult, err := s.CloudinaryClient.UploadSingleFile(ctx, avatarFile, uploadConfig)
			if err != nil {
				return domain.UserResponse{}, domain.ErrBadRequestWithKey(utils.ErrFileUploadFailedKey)
			}

			// Set new avatar URL in payload
			payload.AvatarURL = &uploadResult.SecureURL
			// Note: Cloudinary will automatically overwrite old avatar due to same public ID
		} else {
			return domain.UserResponse{}, domain.ErrBadRequestWithKey(utils.ErrCloudinaryConfigKey)
		}
	} else if payload.AvatarURL != nil {
		// Handle avatar URL changes from JSON/form data
		if *payload.AvatarURL == "" || *payload.AvatarURL == "null" {
			// User wants to remove avatar
			payload.AvatarURL = nil
			shouldDeleteOldAvatar = true
		}
		// If payload.AvatarURL has a valid URL, it will be used as-is
	}

	// Use the UpdateUserWithPayload method
	updatedUser, err := s.Repo.UpdateUserWithPayload(ctx, userId, payload)
	if err != nil {
		return domain.UserResponse{}, err
	}

	// * Delete old avatar from Cloudinary if needed
	if shouldDeleteOldAvatar && s.CloudinaryClient != nil && existingUser.AvatarURL != nil && *existingUser.AvatarURL != "" {
		// Only delete if the old avatar was stored in Cloudinary (contains our public ID pattern)
		if strings.Contains(*existingUser.AvatarURL, "user_"+userId+"_avatar") {
			_ = s.CloudinaryClient.DeleteFile(ctx, oldAvatarPublicID)
			// Note: We don't return error here to avoid failing user update if avatar deletion fails
		}
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

func (s *Service) GetUserStatistics(ctx context.Context) (domain.UserStatistics, error) {
	stats, err := s.Repo.GetUserStatistics(ctx)
	if err != nil {
		return domain.UserStatistics{}, err
	}
	return stats, nil
}
