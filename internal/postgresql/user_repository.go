package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/Rizz404/inventory-api/internal/utils"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type UserFilterOptions struct {
	Role       *domain.UserRole `json:"role,omitempty"`
	IsActive   *bool            `json:"is_active,omitempty"`
	EmployeeID *string          `json:"employee_id,omitempty"`
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) applyUserFilters(db *gorm.DB, filters any) *gorm.DB {
	f, ok := filters.(*UserFilterOptions)
	if !ok || f == nil {
		return db
	}

	if f.Role != nil {
		db = db.Where("u.role = ?", f.Role)
	}
	if f.IsActive != nil {
		db = db.Where("u.is_active = ?", *f.IsActive)
	}
	if f.EmployeeID != nil {
		db = db.Where("u.employee_id = ?", *f.EmployeeID)
	}
	return db
}

func (r *UserRepository) applyUserSorts(db *gorm.DB, sort *query.SortOptions) *gorm.DB {
	if sort == nil || sort.Field == "" {
		return db.Order("u.created_at DESC")
	}
	var orderClause string
	switch strings.ToLower(sort.Field) {
	case "name", "full_name", "email", "role", "employee_id", "is_active", "created_at", "updated_at":
		orderClause = "u." + sort.Field
	default:
		return db.Order("u.created_at DESC")
	}

	order := "DESC"
	if strings.ToLower(sort.Order) == "asc" {
		order = "ASC"
	}
	return db.Order(fmt.Sprintf("%s %s", orderClause, order))
}

// *===========================MUTATION===========================*
func (r *UserRepository) CreateUser(ctx context.Context, payload *domain.User) (domain.User, error) {
	modelUser := mapper.ToModelUser(payload)

	// Create user in database
	err := r.db.WithContext(ctx).Create(&modelUser).Error
	if err != nil {
		return domain.User{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainUser(&modelUser), nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, payload *domain.User) (domain.User, error) {
	var updatedUser model.User
	userID := payload.ID

	// Update user in database
	userUpdates := model.User{
		Name:          payload.Name,
		PasswordHash:  payload.PasswordHash,
		FullName:      payload.FullName,
		Role:          payload.Role,
		EmployeeID:    payload.EmployeeID,
		PreferredLang: payload.PreferredLang,
		IsActive:      payload.IsActive,
		AvatarURL:     payload.AvatarURL,
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

func (r *UserRepository) UpdateUserWithPayload(ctx context.Context, userId string, payload *domain.UpdateUserPayload) (domain.User, error) {
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

func (r *UserRepository) DeleteUser(ctx context.Context, userId string) error {
	err := r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", userId).Error
	if err != nil {
		return domain.ErrInternal(err)
	}
	return nil
}

// *===========================QUERY===========================*
func (r *UserRepository) GetUsersPaginated(ctx context.Context, params query.Params) ([]domain.User, error) {
	var users []model.User
	db := r.db.WithContext(ctx).
		Table("users u")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("u.name ILIKE ? OR u.full_name ILIKE ?", searchPattern, searchPattern)
	}

	// * Set pagination ke nil agar query.Apply tidak memproses cursor
	params.Pagination.Cursor = ""
	db = query.Apply(db, params, r.applyUserFilters, r.applyUserSorts)

	if err := db.Find(&users).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain users
	return mapper.ToDomainUsers(users), nil
}

func (r *UserRepository) GetUsersCursor(ctx context.Context, params query.Params) ([]domain.User, error) {
	var users []model.User
	db := r.db.WithContext(ctx).
		Table("users u")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("u.name ILIKE ? OR u.full_name ILIKE ?", searchPattern, searchPattern)
	}

	// * Set offset ke 0 agar query.Apply tidak memproses offset
	params.Pagination.Offset = 0
	db = query.Apply(db, params, r.applyUserFilters, r.applyUserSorts)

	if err := db.Find(&users).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain users
	return mapper.ToDomainUsers(users), nil
}

func (r *UserRepository) GetUserById(ctx context.Context, userId string) (domain.User, error) {
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

func (r *UserRepository) GetUserByName(ctx context.Context, name string) (domain.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).Where("name = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrNotFound("user with name '" + name + "'")
		}
		return domain.User{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainUser(&user), nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrNotFound("user with email '" + email + "'")
		}
		return domain.User{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainUser(&user), nil
}

func (r *UserRepository) CheckUserExists(ctx context.Context, userId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userId).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *UserRepository) CheckNameExists(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *UserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *UserRepository) CheckNameExistsExcluding(ctx context.Context, name string, excludeUserId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("name = ? AND id != ?", name, excludeUserId).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *UserRepository) CheckEmailExistsExcluding(ctx context.Context, email string, excludeUserId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ? AND id != ?", email, excludeUserId).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *UserRepository) CountUsers(ctx context.Context, params query.Params) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Table("users u")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("u.name ILIKE ? OR u.full_name ILIKE ?", searchPattern, searchPattern)
	}

	db = query.Apply(db, query.Params{Filters: params.Filters}, r.applyUserFilters, nil)

	if err := db.Count(&count).Error; err != nil {
		return 0, domain.ErrInternal(err)
	}
	return count, nil
}

func (r *UserRepository) GetUserStatistics(ctx context.Context) (domain.UserStatistics, error) {
	var stats domain.UserStatistics

	// Get total user count
	var totalCount int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&totalCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.Total.Count = int(totalCount)

	// Get user counts by status (active/inactive)
	var activeCount, inactiveCount int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("is_active = ?", true).Count(&activeCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("is_active = ?", false).Count(&inactiveCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.ByStatus.Active = int(activeCount)
	stats.ByStatus.Inactive = int(inactiveCount)

	// Get user counts by role
	var adminCount, staffCount, employeeCount int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("role = ?", domain.RoleAdmin).Count(&adminCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("role = ?", domain.RoleStaff).Count(&staffCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("role = ?", domain.RoleEmployee).Count(&employeeCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.ByRole.Admin = int(adminCount)
	stats.ByRole.Staff = int(staffCount)
	stats.ByRole.Employee = int(employeeCount)

	// Get registration trends (last 30 days)
	var registrationTrends []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= NOW() - INTERVAL '30 days'").
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&registrationTrends).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.RegistrationTrends = make([]domain.RegistrationTrend, len(registrationTrends))
	for i, rt := range registrationTrends {
		stats.RegistrationTrends[i] = domain.RegistrationTrend{
			Date:  rt.Date,
			Count: int(rt.Count),
		}
	}

	// Calculate summary statistics
	stats.Summary.TotalUsers = int(totalCount)

	if totalCount > 0 {
		stats.Summary.ActiveUsersPercentage = float64(activeCount) / float64(totalCount) * 100
		stats.Summary.InactiveUsersPercentage = float64(inactiveCount) / float64(totalCount) * 100
		stats.Summary.AdminPercentage = float64(adminCount) / float64(totalCount) * 100
		stats.Summary.StaffPercentage = float64(staffCount) / float64(totalCount) * 100
		stats.Summary.EmployeePercentage = float64(employeeCount) / float64(totalCount) * 100
	}

	// Get earliest and latest registration dates
	var earliestDate, latestDate time.Time
	if err := r.db.WithContext(ctx).Model(&model.User{}).Select("MIN(created_at)").Scan(&earliestDate).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.User{}).Select("MAX(created_at)").Scan(&latestDate).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.Summary.EarliestRegistrationDate = earliestDate.Format("2006-01-02")
	stats.Summary.LatestRegistrationDate = latestDate.Format("2006-01-02")

	// Calculate average users per day
	if !earliestDate.IsZero() && !latestDate.IsZero() {
		daysDiff := latestDate.Sub(earliestDate).Hours() / 24
		if daysDiff > 0 {
			stats.Summary.AverageUsersPerDay = float64(totalCount) / daysDiff
		}
	}

	return stats, nil
}
