package domain

import (
	"time"
)

// --- Enums ---

type UserRole string

const (
	RoleAdmin    UserRole = "Admin"
	RoleStaff    UserRole = "Staff"
	RoleEmployee UserRole = "Employee"
)

type UserSortField string

const (
	UserSortByName       UserSortField = "name"
	UserSortByFullName   UserSortField = "fullName"
	UserSortByEmail      UserSortField = "email"
	UserSortByRole       UserSortField = "role"
	UserSortByEmployeeID UserSortField = "employeeId"
	UserSortByIsActive   UserSortField = "isActive"
	UserSortByCreatedAt  UserSortField = "createdAt"
	UserSortByUpdatedAt  UserSortField = "updatedAt"
)

// --- Structs ---

type User struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	PasswordHash  string     `json:"passwordHash"`
	FullName      string     `json:"fullName"`
	Role          UserRole   `json:"role"`
	EmployeeID    *string    `json:"employeeId"` // ! gak usah diapa-apain dulu, soalnya belum ada
	PreferredLang string     `json:"preferredLang"`
	IsActive      bool       `json:"isActive"`
	AvatarURL     *string    `json:"avatarUrl,omitempty"`
	PhoneNumber   *string    `json:"phoneNumber,omitempty"`
	FCMToken      *string    `json:"fcmToken,omitempty"`
	LastLogin     *time.Time `json:"lastLogin,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

// ! jangan omitempty biar client nya tau
type UserResponse struct {
	ID            string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name          string     `json:"name" example:"john_doe"`
	Email         string     `json:"email" example:"john.doe@example.com"`
	FullName      string     `json:"fullName" example:"John Doe"`
	Role          UserRole   `json:"role" example:"Admin"`
	EmployeeID    *string    `json:"employeeId" example:"EMP001"`
	PreferredLang string     `json:"preferredLang" example:"en"`
	IsActive      bool       `json:"isActive" example:"true"`
	AvatarURL     *string    `json:"avatarUrl" example:"https://example.com/avatar.jpg"`
	PhoneNumber   *string    `json:"phoneNumber" example:"+6281234567890"`
	FCMToken      *string    `json:"fcmToken"`
	LastLogin     *time.Time `json:"lastLogin" example:"2023-01-01T00:00:00Z"`
	CreatedAt     time.Time  `json:"createdAt" example:"2023-01-01T00:00:00Z"`
	UpdatedAt     time.Time  `json:"updatedAt" example:"2023-01-01T00:00:00Z"`
}

type UserListResponse struct {
	ID            string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name          string     `json:"name" example:"john_doe"`
	Email         string     `json:"email" example:"john.doe@example.com"`
	FullName      string     `json:"fullName" example:"John Doe"`
	Role          UserRole   `json:"role" example:"Admin"`
	EmployeeID    *string    `json:"employeeId" example:"EMP001"`
	PreferredLang string     `json:"preferredLang" example:"en"`
	IsActive      bool       `json:"isActive" example:"true"`
	AvatarURL     *string    `json:"avatarUrl" example:"https://example.com/avatar.jpg"`
	PhoneNumber   *string    `json:"phoneNumber" example:"+6281234567890"`
	LastLogin     *time.Time `json:"lastLogin" example:"2023-01-01T00:00:00Z"`
	CreatedAt     time.Time  `json:"createdAt" example:"2023-01-01T00:00:00Z"`
	UpdatedAt     time.Time  `json:"updatedAt" example:"2023-01-01T00:00:00Z"`
}

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
	RefreshToken string       `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
}

// Alias untuk backward compatibility
type LoginResponse = AuthResponse

type BulkDeleteUsers struct {
	RequestedIDS []string `json:"requestedIds"`
	DeletedIDS   []string `json:"deletedIds"`
}

type BulkDeleteUsersResponse struct {
	RequestedIDS []string `json:"requestedIds"`
	DeletedIDS   []string `json:"deletedIds"`
}

// --- Bulk Create ---

type BulkCreateUsersPayload struct {
	Users []CreateUserPayload `json:"users" validate:"required,min=1,max=100,dive"`
}

type BulkCreateUsersResponse struct {
	Users []UserResponse `json:"users"`
}

// --- Payloads ---

type LoginPayload struct {
	Email    string `json:"email" example:"john.doe@example.com" form:"email" validate:"required,email"`
	Password string `json:"password" example:"password123" form:"password" validate:"required,min=5"`
}

type RegisterPayload struct {
	Name     string `json:"name" example:"john_doe" form:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" example:"john.doe@example.com" form:"email" validate:"required,email"`
	Password string `json:"password" example:"password123" form:"password" validate:"required,min=5"`
}

type CreateUserPayload struct {
	Name          string   `json:"name" example:"john_doe" form:"name" validate:"required,min=3,max=50"`
	Email         string   `json:"email" example:"john.doe@example.com" form:"email" validate:"required,email,max=255"`
	Password      string   `json:"password" example:"password123" form:"password" validate:"required,min=8,max=100"`
	FullName      string   `json:"fullName" example:"John Doe" form:"fullName" validate:"required,min=3,max=100"`
	Role          UserRole `json:"role" example:"Admin" form:"role" validate:"required,oneof=Admin Staff Employee"`
	EmployeeID    *string  `json:"employeeId,omitempty" example:"EMP001" form:"employeeId" validate:"omitempty,max=20"` // ! gak usah diapa-apain dulu, soalnya belum ada
	PreferredLang *string  `json:"preferredLang,omitempty" example:"en" form:"preferredLang" validate:"omitempty,max=5"`
	IsActive      *bool    `json:"isActive,omitempty" example:"true" form:"isActive" validate:"omitempty"`
	AvatarURL     *string  `json:"avatarUrl,omitempty" example:"https://example.com/avatar.jpg" form:"avatarUrl" validate:"omitempty,url"`
	PhoneNumber   *string  `json:"phoneNumber,omitempty" example:"+6281234567890" form:"phoneNumber" validate:"omitempty,max=20"`
}

type UpdateUserPayload struct {
	Name          *string   `json:"name,omitempty" example:"john_doe" form:"name" validate:"omitempty,min=3,max=50"`
	Email         *string   `json:"email,omitempty" example:"john.doe@example.com" form:"email" validate:"omitempty,email,max=255"`
	FullName      *string   `json:"fullName,omitempty" example:"John Doe" form:"fullName" validate:"omitempty,min=3,max=100"`
	Role          *UserRole `json:"role,omitempty" example:"Admin" form:"role" validate:"omitempty,oneof=Admin Staff Employee"`
	EmployeeID    *string   `json:"employeeId,omitempty" example:"EMP001" form:"employeeId" validate:"omitempty,max=20"` // ! gak usah diapa-apain dulu, soalnya belum ada
	PreferredLang *string   `json:"preferredLang,omitempty" example:"en" form:"preferredLang" validate:"omitempty,max=5"`
	IsActive      *bool     `json:"isActive,omitempty" example:"true" form:"isActive" validate:"omitempty"`
	AvatarURL     *string   `json:"avatarUrl,omitempty" example:"https://example.com/avatar.jpg" form:"avatarUrl" validate:"omitempty,url"`
	PhoneNumber   *string   `json:"phoneNumber,omitempty" example:"+6281234567890" form:"phoneNumber" validate:"omitempty,max=20"`
	FCMToken      *string   `json:"fcmToken,omitempty" form:"fcmToken" validate:"omitempty"`
}

type ChangePasswordPayload struct {
	OldPassword string `json:"oldPassword" form:"oldPassword" validate:"required,min=8" example:"oldpassword123"`
	NewPassword string `json:"newPassword" form:"newPassword" validate:"required,min=8" example:"newpassword123"`
}

type RefreshTokenPayload struct {
	RefreshToken string `json:"refreshToken" form:"refreshToken" validate:"required"`
}

// --- Forgot Password Payloads ---

type ForgotPasswordPayload struct {
	Email string `json:"email" form:"email" validate:"required,email" example:"john.doe@example.com"`
}

type VerifyResetCodePayload struct {
	Email string `json:"email" form:"email" validate:"required,email" example:"john.doe@example.com"`
	Code  string `json:"code" form:"code" validate:"required,len=6" example:"123456"`
}

type ResetPasswordPayload struct {
	Email       string `json:"email" form:"email" validate:"required,email" example:"john.doe@example.com"`
	Code        string `json:"code" form:"code" validate:"required,len=6" example:"123456"`
	NewPassword string `json:"newPassword" form:"newPassword" validate:"required,min=8" example:"newpassword123"`
}

// --- Forgot Password Responses ---

type VerifyResetCodeResponse struct {
	Valid bool `json:"valid" example:"true"`
}

type BulkDeleteUsersPayload struct {
	IDS []string `json:"ids" validate:"required,min=1,max=100,dive,required"`
}

type ExportUserListPayload struct {
	Format      ExportFormat       `json:"format" validate:"required,oneof=pdf excel"`
	SearchQuery *string            `json:"searchQuery,omitempty"`
	Filters     *UserFilterOptions `json:"filters,omitempty"`
	Sort        *UserSortOptions   `json:"sort,omitempty"`
}

// --- Query Parameters ---

type UserFilterOptions struct {
	Role       *UserRole `json:"role,omitempty"`
	IsActive   *bool     `json:"is_active,omitempty"`
	EmployeeID *string   `json:"employee_id,omitempty"`
}

type UserSortOptions struct {
	Field UserSortField `json:"field" example:"created_at"`
	Order SortOrder     `json:"order" example:"desc"`
}

type UserParams struct {
	SearchQuery *string            `json:"searchQuery,omitempty"`
	Filters     *UserFilterOptions `json:"filters,omitempty"`
	Sort        *UserSortOptions   `json:"sort,omitempty"`
	Pagination  *PaginationOptions `json:"pagination,omitempty"`
}

// --- Statistics ---

// Internal statistics structs (used in repository layer)
type UserStatistics struct {
	Total              UserCountStatistics   `json:"total"`
	ByStatus           UserStatusStatistics  `json:"byStatus"`
	ByRole             UserRoleStatistics    `json:"byRole"`
	RegistrationTrends []RegistrationTrend   `json:"registrationTrends"`
	Summary            UserSummaryStatistics `json:"summary"`
}

type UserCountStatistics struct {
	Count int `json:"count"`
}

type UserStatusStatistics struct {
	Active   int `json:"active"`
	Inactive int `json:"inactive"`
}

type UserRoleStatistics struct {
	Admin    int `json:"admin"`
	Staff    int `json:"staff"`
	Employee int `json:"employee"`
}

type RegistrationTrend struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

type UserSummaryStatistics struct {
	TotalUsers               int       `json:"totalUsers"`
	ActiveUsersPercentage    float64   `json:"activeUsersPercentage"`
	InactiveUsersPercentage  float64   `json:"inactiveUsersPercentage"`
	AdminPercentage          float64   `json:"adminPercentage"`
	StaffPercentage          float64   `json:"staffPercentage"`
	EmployeePercentage       float64   `json:"employeePercentage"`
	AverageUsersPerDay       float64   `json:"averageUsersPerDay"`
	LatestRegistrationDate   time.Time `json:"latestRegistrationDate"`
	EarliestRegistrationDate time.Time `json:"earliestRegistrationDate"`
}

// Response statistics structs (used in service/handler layer)
type UserStatisticsResponse struct {
	Total              UserCountStatisticsResponse   `json:"total"`
	ByStatus           UserStatusStatisticsResponse  `json:"byStatus"`
	ByRole             UserRoleStatisticsResponse    `json:"byRole"`
	RegistrationTrends []RegistrationTrendResponse   `json:"registrationTrends"`
	Summary            UserSummaryStatisticsResponse `json:"summary"`
}

type UserCountStatisticsResponse struct {
	Count int `json:"count"`
}

type UserStatusStatisticsResponse struct {
	Active   int `json:"active"`
	Inactive int `json:"inactive"`
}

type UserRoleStatisticsResponse struct {
	Admin    int `json:"admin"`
	Staff    int `json:"staff"`
	Employee int `json:"employee"`
}

type RegistrationTrendResponse struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

type UserSummaryStatisticsResponse struct {
	TotalUsers               int       `json:"totalUsers"`
	ActiveUsersPercentage    Decimal2  `json:"activeUsersPercentage"`
	InactiveUsersPercentage  Decimal2  `json:"inactiveUsersPercentage"`
	AdminPercentage          Decimal2  `json:"adminPercentage"`
	StaffPercentage          Decimal2  `json:"staffPercentage"`
	EmployeePercentage       Decimal2  `json:"employeePercentage"`
	AverageUsersPerDay       Decimal2  `json:"averageUsersPerDay"`
	LatestRegistrationDate   time.Time `json:"latestRegistrationDate"`
	EarliestRegistrationDate time.Time `json:"earliestRegistrationDate"`
}

// --- Personal Statistics (Employee only) ---

// Internal personal statistics structs (used in repository layer)
type UserPersonalStatistics struct {
	UserID       string                                `json:"userId"`
	UserName     string                                `json:"userName"`
	Role         UserRole                              `json:"role"`
	Assets       UserPersonalAssetStatistics           `json:"assets"`
	IssueReports UserPersonalIssueReportStatistics     `json:"issueReports"`
	Summary      UserPersonalSummaryStatistics         `json:"summary"`
}

type UserPersonalAssetStatistics struct {
	Total       UserPersonalAssetTotalStatistics      `json:"total"`
	ByCondition UserPersonalAssetConditionStatistics   `json:"byCondition"`
	Items       []UserPersonalAssetItem                `json:"items"`
}

type UserPersonalAssetTotalStatistics struct {
	Count      int     `json:"count"`
	TotalValue float64 `json:"totalValue"`
}

type UserPersonalAssetConditionStatistics struct {
	Good    int `json:"good"`
	Fair    int `json:"fair"`
	Poor    int `json:"poor"`
	Damaged int `json:"damaged"`
}

type UserPersonalAssetItem struct {
	AssetID      string    `json:"assetId"`
	AssetTag     string    `json:"assetTag"`
	Name         string    `json:"name"`
	Category     string    `json:"category"`
	Condition    string    `json:"condition"`
	Value        float64   `json:"value"`
	AssignedDate time.Time `json:"assignedDate"`
}

type UserPersonalIssueReportStatistics struct {
	Total         UserPersonalIssueReportTotalStatistics    `json:"total"`
	ByStatus      UserPersonalIssueReportStatusStatistics   `json:"byStatus"`
	ByPriority    UserPersonalIssueReportPriorityStatistics `json:"byPriority"`
	RecentIssues  []UserPersonalIssueReportItem             `json:"recentIssues"`
	Summary       UserPersonalIssueReportSummaryStatistics  `json:"summary"`
}

type UserPersonalIssueReportTotalStatistics struct {
	Count int `json:"count"`
}

type UserPersonalIssueReportStatusStatistics struct {
	Open       int `json:"open"`
	InProgress int `json:"inProgress"`
	Resolved   int `json:"resolved"`
	Closed     int `json:"closed"`
}

type UserPersonalIssueReportPriorityStatistics struct {
	High   int `json:"high"`
	Medium int `json:"medium"`
	Low    int `json:"low"`
}

type UserPersonalIssueReportItem struct {
	IssueID      string    `json:"issueId"`
	AssetID      *string   `json:"assetId"`
	AssetTag     *string   `json:"assetTag"`
	Title        string    `json:"title"`
	Priority     string    `json:"priority"`
	Status       string    `json:"status"`
	ReportedDate time.Time `json:"reportedDate"`
}

type UserPersonalIssueReportSummaryStatistics struct {
	OpenIssuesCount         int     `json:"openIssuesCount"`
	ResolvedIssuesCount     int     `json:"resolvedIssuesCount"`
	AverageResolutionDays   float64 `json:"averageResolutionDays"`
}

type UserPersonalSummaryStatistics struct {
	AccountCreatedDate time.Time  `json:"accountCreatedDate"`
	AccountAge         string     `json:"accountAge"`
	LastLogin          *time.Time `json:"lastLogin"`
	HasActiveIssues    bool       `json:"hasActiveIssues"`
	HealthScore        int        `json:"healthScore"`
}

// Response personal statistics structs (used in service/handler layer)
type UserPersonalStatisticsResponse struct {
	UserID       string                                        `json:"userId"`
	UserName     string                                        `json:"userName"`
	Role         UserRole                                      `json:"role"`
	Assets       UserPersonalAssetStatisticsResponse           `json:"assets"`
	IssueReports UserPersonalIssueReportStatisticsResponse     `json:"issueReports"`
	Summary      UserPersonalSummaryStatisticsResponse         `json:"summary"`
}

type UserPersonalAssetStatisticsResponse struct {
	Total       UserPersonalAssetTotalStatisticsResponse      `json:"total"`
	ByCondition UserPersonalAssetConditionStatisticsResponse   `json:"byCondition"`
	Items       []UserPersonalAssetItemResponse                `json:"items"`
}

type UserPersonalAssetTotalStatisticsResponse struct {
	Count      int      `json:"count"`
	TotalValue Decimal2 `json:"totalValue"`
}

type UserPersonalAssetConditionStatisticsResponse struct {
	Good    int `json:"good"`
	Fair    int `json:"fair"`
	Poor    int `json:"poor"`
	Damaged int `json:"damaged"`
}

type UserPersonalAssetItemResponse struct {
	AssetID      string    `json:"assetId"`
	AssetTag     string    `json:"assetTag"`
	Name         string    `json:"name"`
	Category     string    `json:"category"`
	Condition    string    `json:"condition"`
	Value        Decimal2  `json:"value"`
	AssignedDate time.Time `json:"assignedDate"`
}

type UserPersonalIssueReportStatisticsResponse struct {
	Total         UserPersonalIssueReportTotalStatisticsResponse    `json:"total"`
	ByStatus      UserPersonalIssueReportStatusStatisticsResponse   `json:"byStatus"`
	ByPriority    UserPersonalIssueReportPriorityStatisticsResponse `json:"byPriority"`
	RecentIssues  []UserPersonalIssueReportItemResponse             `json:"recentIssues"`
	Summary       UserPersonalIssueReportSummaryStatisticsResponse  `json:"summary"`
}

type UserPersonalIssueReportTotalStatisticsResponse struct {
	Count int `json:"count"`
}

type UserPersonalIssueReportStatusStatisticsResponse struct {
	Open       int `json:"open"`
	InProgress int `json:"inProgress"`
	Resolved   int `json:"resolved"`
	Closed     int `json:"closed"`
}

type UserPersonalIssueReportPriorityStatisticsResponse struct {
	High   int `json:"high"`
	Medium int `json:"medium"`
	Low    int `json:"low"`
}

type UserPersonalIssueReportItemResponse struct {
	IssueID      string    `json:"issueId"`
	AssetID      *string   `json:"assetId"`
	AssetTag     *string   `json:"assetTag"`
	Title        string    `json:"title"`
	Priority     string    `json:"priority"`
	Status       string    `json:"status"`
	ReportedDate time.Time `json:"reportedDate"`
}

type UserPersonalIssueReportSummaryStatisticsResponse struct {
	OpenIssuesCount         int      `json:"openIssuesCount"`
	ResolvedIssuesCount     int      `json:"resolvedIssuesCount"`
	AverageResolutionDays   Decimal2 `json:"averageResolutionDays"`
}

type UserPersonalSummaryStatisticsResponse struct {
	AccountCreatedDate time.Time  `json:"accountCreatedDate"`
	AccountAge         string     `json:"accountAge"`
	LastLogin          *time.Time `json:"lastLogin"`
	HasActiveIssues    bool       `json:"hasActiveIssues"`
	HealthScore        int        `json:"healthScore"`
}
