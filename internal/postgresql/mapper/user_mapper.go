package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
)

// Model <-> Domain conversions (for repository layer)
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
	}
}

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

// Domain -> Response conversions (for service layer)
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
		CreatedAt:     u.CreatedAt.Format(TimeFormat),
		UpdatedAt:     u.UpdatedAt.Format(TimeFormat),
	}
}

func UsersToResponses(users []domain.User) []domain.UserResponse {
	responses := make([]domain.UserResponse, len(users))
	for i, user := range users {
		responses[i] = UserToResponse(&user)
	}
	return responses
}

// Statistics conversions
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
			ActiveUsersPercentage:    stats.Summary.ActiveUsersPercentage,
			InactiveUsersPercentage:  stats.Summary.InactiveUsersPercentage,
			AdminPercentage:          stats.Summary.AdminPercentage,
			StaffPercentage:          stats.Summary.StaffPercentage,
			EmployeePercentage:       stats.Summary.EmployeePercentage,
			AverageUsersPerDay:       stats.Summary.AverageUsersPerDay,
			LatestRegistrationDate:   stats.Summary.LatestRegistrationDate,
			EarliestRegistrationDate: stats.Summary.EarliestRegistrationDate,
		},
	}
}

// Update payload mapping
func ToModelUserUpdateMap(payload *domain.UpdateUserPayload) map[string]any {
	updates := make(map[string]any)

	if payload.Name != nil {
		updates["name"] = *payload.Name
	}
	if payload.Email != nil {
		updates["email"] = *payload.Email
	}
	if payload.Password != nil {
		updates["password_hash"] = *payload.Password
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

	return updates
}
