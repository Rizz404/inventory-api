package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/Rizz404/inventory-api/internal/utils"
	"gorm.io/gorm"
)

type gormUserRepository struct {
	db *gorm.DB
}

type UserFilterOptions struct {
	Role     *domain.UserRole `json:"role,omitempty"`
	IsActive *bool            `json:"is_active,omitempty"`
}

func NewUserRepository(db *gorm.DB) *gormUserRepository {
	return &gormUserRepository{
		db: db,
	}
}

func (r *gormUserRepository) applyUserFilters(db *gorm.DB, filters any) *gorm.DB {
	f, ok := filters.(*UserFilterOptions)
	if !ok || f == nil {
		return db
	}

	if f.Role != nil {
		db = db.Where("users.role = ?", f.Role)
	}
	if f.IsActive != nil {
		db = db.Where("users.is_active = ?", *f.IsActive)
	}
	return db
}

func (r *gormUserRepository) applyUserSorts(db *gorm.DB, sort *query.SortOptions) *gorm.DB {
	if sort == nil || sort.Field == "" {
		return db.Order("users.created_at DESC")
	}
	var orderClause string
	switch strings.ToLower(sort.Field) {
	case "username", "full_name", "created_at", "updated_at", "role":
		orderClause = "users." + sort.Field
	default:
		return db.Order("users.created_at DESC")
	}

	order := "DESC"
	if strings.ToLower(sort.Order) == "asc" {
		order = "ASC"
	}
	return db.Order(fmt.Sprintf("%s %s", orderClause, order))
}

// *===========================MUTATION===========================*
func (r *gormUserRepository) CreateUser(ctx context.Context, payload *domain.User) (domain.User, error) {
	modelUser := mapper.ToModelUser(payload)

	// Create user in database
	err := r.db.WithContext(ctx).Create(&modelUser).Error
	if err != nil {
		return domain.User{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainUser(&modelUser), nil
}

func (r *gormUserRepository) UpdateUser(ctx context.Context, payload *domain.User) (domain.User, error) {
	var updatedUser model.User
	userID := payload.ID

	// Update user in database
	userUpdates := model.User{
		Username:      payload.Username,
		PasswordHash:  payload.PasswordHash,
		FullName:      payload.FullName,
		Role:          payload.Role,
		EmployeeID:    payload.EmployeeID,
		PreferredLang: payload.PreferredLang,
		IsActive:      payload.IsActive,
	}

	err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Updates(userUpdates).Error
	if err != nil {
		return domain.User{}, domain.ErrInternal(err)
	}

	// Get updated user
	err = r.db.WithContext(ctx).First(&updatedUser, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrNotFound("user")
		}
		return domain.User{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainUser(&updatedUser), nil
}

func (r *gormUserRepository) UpdateUserWithPayload(ctx context.Context, userId string, payload *domain.UpdateUserPayload) (domain.User, error) {
	var updatedUser model.User

	// Build update map from payload
	updates := mapper.ToModelUserUpdateMap(payload)

	// If password is provided, hash it
	if payload.Password != nil {
		hashedPassword, err := utils.HashPassword(*payload.Password)
		if err != nil {
			return domain.User{}, domain.ErrInternal(err)
		}
		updates["password_hash"] = hashedPassword
	}

	// Perform update
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userId).Updates(updates).Error
	if err != nil {
		return domain.User{}, domain.ErrInternal(err)
	}

	// Get updated user
	err = r.db.WithContext(ctx).First(&updatedUser, "id = ?", userId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrNotFound("user")
		}
		return domain.User{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainUser(&updatedUser), nil
}

func (r *gormUserRepository) DeleteUser(ctx context.Context, userId string) error {
	err := r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", userId).Error
	if err != nil {
		return domain.ErrInternal(err)
	}
	return nil
}

// *===========================QUERY===========================*
func (r *gormUserRepository) GetUsersPaginated(ctx context.Context, params query.Params) ([]domain.User, error) {
	var users []model.User
	db := r.db.WithContext(ctx)

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("users.username ILIKE ? OR users.full_name ILIKE ?", searchPattern, searchPattern)
	}

	// * Set pagination ke nil agar query.Apply tidak memproses cursor
	params.Pagination.Cursor = ""
	db = query.Apply(db, params, r.applyUserFilters, r.applyUserSorts)

	if err := db.Find(&users).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain users
	domainUsers := make([]domain.User, len(users))
	for i, user := range users {
		domainUsers[i] = mapper.ToDomainUser(&user)
	}
	return domainUsers, nil
}

func (r *gormUserRepository) GetUsersCursor(ctx context.Context, params query.Params) ([]domain.User, error) {
	var users []model.User
	db := r.db.WithContext(ctx)

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("users.username ILIKE ? OR users.full_name ILIKE ?", searchPattern, searchPattern)
	}

	// * Set offset ke 0 agar query.Apply tidak memproses offset
	params.Pagination.Offset = 0
	db = query.Apply(db, params, r.applyUserFilters, r.applyUserSorts)

	if err := db.Find(&users).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain users
	domainUsers := make([]domain.User, len(users))
	for i, user := range users {
		domainUsers[i] = mapper.ToDomainUser(&user)
	}
	return domainUsers, nil
}

func (r *gormUserRepository) GetUserById(ctx context.Context, userId string) (domain.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).First(&user, "id = ?", userId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrNotFound("user")
		}
		return domain.User{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainUser(&user), nil
}

func (r *gormUserRepository) GetUserByUsernameOrEmail(ctx context.Context, username string, email string) (domain.User, error) {
	var user model.User

	// Since we don't have email field in current domain, just search by username
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrNotFound("user")
		}
		return domain.User{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainUser(&user), nil
}

func (r *gormUserRepository) CheckUserExist(ctx context.Context, userId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userId).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *gormUserRepository) CheckUsernameExist(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *gormUserRepository) CheckEmailExist(ctx context.Context, email string) (bool, error) {
	// Since we don't have email field in current domain, return false
	// This method exists for interface compatibility
	return false, nil
}

func (r *gormUserRepository) CountUsers(ctx context.Context, params query.Params) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(&model.User{})

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("users.username ILIKE ? OR users.full_name ILIKE ?", searchPattern, searchPattern)
	}

	db = query.Apply(db, query.Params{Filters: params.Filters}, r.applyUserFilters, nil)

	if err := db.Count(&count).Error; err != nil {
		return 0, domain.ErrInternal(err)
	}
	return count, nil
}
