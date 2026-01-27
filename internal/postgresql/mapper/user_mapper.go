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
		PhoneNumber:   d.PhoneNumber,
		FCMToken:      d.FCMToken,
		LastLogin:     d.LastLogin,
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
		PhoneNumber:   d.PhoneNumber,
		FCMToken:      d.FCMToken,
		LastLogin:     d.LastLogin,
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
		PhoneNumber:   m.PhoneNumber,
		FCMToken:      m.FCMToken,
		LastLogin:     m.LastLogin,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func ToDomainUsers(models []model.User) []domain.User {
	if len(models) == 0 {
		return []domain.User{}
	}
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
		PhoneNumber:   u.PhoneNumber,
		FCMToken:      u.FCMToken,
		LastLogin:     u.LastLogin,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

func UsersToResponses(users []domain.User) []domain.UserResponse {
	if len(users) == 0 {
		return []domain.UserResponse{}
	}
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
		PhoneNumber:   u.PhoneNumber,
		LastLogin:     u.LastLogin,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

func UsersToListResponses(users []domain.User) []domain.UserListResponse {
	if len(users) == 0 {
		return []domain.UserListResponse{}
	}
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

func PersonalStatisticsToResponse(stats *domain.UserPersonalStatistics) domain.UserPersonalStatisticsResponse {
	// Convert asset items
	assetItems := make([]domain.UserPersonalAssetItemResponse, len(stats.Assets.Items))
	for i, item := range stats.Assets.Items {
		assetItems[i] = domain.UserPersonalAssetItemResponse{
			AssetID:      item.AssetID,
			AssetTag:     item.AssetTag,
			Name:         item.Name,
			Category:     item.Category,
			Condition:    item.Condition,
			Value:        domain.NewDecimal2(item.Value),
			AssignedDate: item.AssignedDate,
		}
	}

	// Convert issue report items
	issueItems := make([]domain.UserPersonalIssueReportItemResponse, len(stats.IssueReports.RecentIssues))
	for i, item := range stats.IssueReports.RecentIssues {
		issueItems[i] = domain.UserPersonalIssueReportItemResponse{
			IssueID:      item.IssueID,
			AssetID:      item.AssetID,
			AssetTag:     item.AssetTag,
			Title:        item.Title,
			Priority:     item.Priority,
			Status:       item.Status,
			ReportedDate: item.ReportedDate,
		}
	}

	return domain.UserPersonalStatisticsResponse{
		UserID:   stats.UserID,
		UserName: stats.UserName,
		Role:     stats.Role,
		Assets: domain.UserPersonalAssetStatisticsResponse{
			Total: domain.UserPersonalAssetTotalStatisticsResponse{
				Count:      stats.Assets.Total.Count,
				TotalValue: domain.NewDecimal2(stats.Assets.Total.TotalValue),
			},
			ByCondition: domain.UserPersonalAssetConditionStatisticsResponse{
				Good:    stats.Assets.ByCondition.Good,
				Fair:    stats.Assets.ByCondition.Fair,
				Poor:    stats.Assets.ByCondition.Poor,
				Damaged: stats.Assets.ByCondition.Damaged,
			},
			Items: assetItems,
		},
		IssueReports: domain.UserPersonalIssueReportStatisticsResponse{
			Total: domain.UserPersonalIssueReportTotalStatisticsResponse{
				Count: stats.IssueReports.Total.Count,
			},
			ByStatus: domain.UserPersonalIssueReportStatusStatisticsResponse{
				Open:       stats.IssueReports.ByStatus.Open,
				InProgress: stats.IssueReports.ByStatus.InProgress,
				Resolved:   stats.IssueReports.ByStatus.Resolved,
				Closed:     stats.IssueReports.ByStatus.Closed,
			},
			ByPriority: domain.UserPersonalIssueReportPriorityStatisticsResponse{
				High:   stats.IssueReports.ByPriority.High,
				Medium: stats.IssueReports.ByPriority.Medium,
				Low:    stats.IssueReports.ByPriority.Low,
			},
			RecentIssues: issueItems,
			Summary: domain.UserPersonalIssueReportSummaryStatisticsResponse{
				OpenIssuesCount:         stats.IssueReports.Summary.OpenIssuesCount,
				ResolvedIssuesCount:     stats.IssueReports.Summary.ResolvedIssuesCount,
				AverageResolutionDays:   domain.NewDecimal2(stats.IssueReports.Summary.AverageResolutionDays),
			},
		},
		Summary: domain.UserPersonalSummaryStatisticsResponse{
			AccountCreatedDate: stats.Summary.AccountCreatedDate,
			AccountAge:         stats.Summary.AccountAge,
			LastLogin:          stats.Summary.LastLogin,
			HasActiveIssues:    stats.Summary.HasActiveIssues,
			HealthScore:        stats.Summary.HealthScore,
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
