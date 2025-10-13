package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
)

// *==================== Model conversions ====================
func ToModelUser(d *domain.User) model.User {
	return model.User{
		Name:          d.Name,
		Email:         d.Email,
		PasswordHash:  d.PasswordHash,
		FullName:      d.FullName,
		Role:          d.Role,
		EmployeeID:    d.EmployeeID,
		PreferredLang: d.PreferredLang,
		IsActive:      d.IsActive,
		AvatarURL:     d.AvatarURL,
		FCMToken:      d.FCMToken,
	}
}

func ToModelUserForCreate(d *domain.User) model.User {
	return model.User{
		Name:          d.Name,
		Email:         d.Email,
		PasswordHash:  d.PasswordHash,
		FullName:      d.FullName,
		Role:          d.Role,
		EmployeeID:    d.EmployeeID,
		PreferredLang: d.PreferredLang,
		IsActive:      d.IsActive,
		AvatarURL:     d.AvatarURL,
		FCMToken:      d.FCMToken,
	}
}

// *==================== Entity conversions ====================
func ToDomainUser(m *model.User) domain.User {
	return domain.User{
		ID:            m.ID.String(),
		Name:          m.Name,
		Email:         m.Email,
		PasswordHash:  m.PasswordHash,
		FullName:      m.FullName,
		Role:          m.Role,
		EmployeeID:    m.EmployeeID,
		PreferredLang: m.PreferredLang,
		IsActive:      m.IsActive,
		AvatarURL:     m.AvatarURL,
		FCMToken:      m.FCMToken,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func ToDomainUsers(models []model.User) []domain.User {
	users := make([]domain.User, len(models))
	for i, m := range models {
		users[i] = ToDomainUser(&m)
	}
	return users
}

// *==================== Entity Response conversions ====================
func UserToResponse(u *domain.User) domain.UserResponse {
	return domain.UserResponse{
		ID:            u.ID,
		Name:          u.Name,
		Email:         u.Email,
		FullName:      u.FullName,
		Role:          u.Role,
		EmployeeID:    u.EmployeeID,
		PreferredLang: u.PreferredLang,
		IsActive:      u.IsActive,
		AvatarURL:     u.AvatarURL,
		FCMToken:      u.FCMToken,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

func UsersToResponses(users []domain.User) []domain.UserResponse {
	responses := make([]domain.UserResponse, len(users))
	for i, user := range users {
		responses[i] = UserToResponse(&user)
	}
	return responses
}

func UserToListResponse(u *domain.User) domain.UserListResponse {
	return domain.UserListResponse{
		ID:            u.ID,
		Name:          u.Name,
		Email:         u.Email,
		FullName:      u.FullName,
		Role:          u.Role,
		EmployeeID:    u.EmployeeID,
		PreferredLang: u.PreferredLang,
		IsActive:      u.IsActive,
		AvatarURL:     u.AvatarURL,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

func UsersToListResponses(users []domain.User) []domain.UserListResponse {
	responses := make([]domain.UserListResponse, len(users))
	for i, user := range users {
		responses[i] = UserToListResponse(&user)
	}
	return responses
}

// *==================== Statistics conversions ====================
func StatisticsToResponse(stats *domain.UserStatistics) domain.UserStatisticsResponse {
	trends := make([]domain.RegistrationTrendResponse, len(stats.RegistrationTrends))
	for i, trend := range stats.RegistrationTrends {
		trends[i] = domain.RegistrationTrendResponse{
			Date:  trend.Date,
			Count: trend.Count,
		}
	}

	return domain.UserStatisticsResponse{
		Total: domain.UserCountStatisticsResponse{
			Count: stats.Total.Count,
		},
		ByStatus: domain.UserStatusStatisticsResponse{
			Active:   stats.ByStatus.Active,
			Inactive: stats.ByStatus.Inactive,
		},
		ByRole: domain.UserRoleStatisticsResponse{
			Admin:    stats.ByRole.Admin,
			Staff:    stats.ByRole.Staff,
			Employee: stats.ByRole.Employee,
		},
		RegistrationTrends: trends,
		Summary: domain.UserSummaryStatisticsResponse{
			TotalUsers:               stats.Summary.TotalUsers,
			ActiveUsersPercentage:    domain.NewDecimal2(stats.Summary.ActiveUsersPercentage),
			InactiveUsersPercentage:  domain.NewDecimal2(stats.Summary.InactiveUsersPercentage),
			AdminPercentage:          domain.NewDecimal2(stats.Summary.AdminPercentage),
			StaffPercentage:          domain.NewDecimal2(stats.Summary.StaffPercentage),
			EmployeePercentage:       domain.NewDecimal2(stats.Summary.EmployeePercentage),
			AverageUsersPerDay:       domain.NewDecimal2(stats.Summary.AverageUsersPerDay),
			LatestRegistrationDate:   stats.Summary.LatestRegistrationDate,
			EarliestRegistrationDate: stats.Summary.EarliestRegistrationDate,
		},
	}
}

// *==================== Update Map conversions (Harus snake case karena untuk database) ====================
func ToModelUserUpdateMap(payload *domain.UpdateUserPayload) map[string]any {
	updates := make(map[string]any)

	if payload.Name != nil {
		updates["name"] = *payload.Name
	}
	if payload.Email != nil {
		updates["email"] = *payload.Email
	}
	if payload.FullName != nil {
		updates["full_name"] = *payload.FullName
	}
	if payload.Role != nil {
		updates["role"] = *payload.Role
	}
	if payload.EmployeeID != nil {
		updates["employee_id"] = payload.EmployeeID
	}
	if payload.PreferredLang != nil {
		updates["preferred_lang"] = *payload.PreferredLang
	}
	if payload.IsActive != nil {
		updates["is_active"] = *payload.IsActive
	}
	if payload.AvatarURL != nil {
		updates["avatar_url"] = payload.AvatarURL
	}
	if payload.FCMToken != nil {
		updates["fcm_token"] = payload.FCMToken
	}

	return updates
}

func MapUserSortFieldToColumn(field domain.UserSortField) string {
	columnMap := map[domain.UserSortField]string{
		domain.UserSortByName:       "name",
		domain.UserSortByFullName:   "full_name",
		domain.UserSortByEmail:      "email",
		domain.UserSortByRole:       "role",
		domain.UserSortByEmployeeID: "employee_id",
		domain.UserSortByIsActive:   "is_active",
		domain.UserSortByCreatedAt:  "created_at",
		domain.UserSortByUpdatedAt:  "updated_at",
	}

	if column, exists := columnMap[field]; exists {
		return column
	}
	return "created_at"
}
