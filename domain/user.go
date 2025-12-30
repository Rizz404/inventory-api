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
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"passwordHash"`
	FullName      string    `json:"fullName"`
	Role          UserRole  `json:"role"`
	EmployeeID    *string   `json:"employeeId"` // ! gak usah diapa-apain dulu, soalnya belum ada
	PreferredLang string    `json:"preferredLang"`
	IsActive      bool      `json:"isActive"`
	AvatarURL     *string   `json:"avatarUrl,omitempty"`
	FCMToken      *string   `json:"fcmToken,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// ! jangan omitempty biar client nya tau
type UserResponse struct {
	ID            string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name          string    `json:"name" example:"john_doe"`
	Email         string    `json:"email" example:"john.doe@example.com"`
	FullName      string    `json:"fullName" example:"John Doe"`
	Role          UserRole  `json:"role" example:"Admin"`
	EmployeeID    *string   `json:"employeeId" example:"EMP001"`
	PreferredLang string    `json:"preferredLang" example:"en"`
	IsActive      bool      `json:"isActive" example:"true"`
	AvatarURL     *string   `json:"avatarUrl" example:"https://example.com/avatar.jpg"`
	FCMToken      *string   `json:"fcmToken"`
	CreatedAt     time.Time `json:"createdAt" example:"2023-01-01T00:00:00Z"`
	UpdatedAt     time.Time `json:"updatedAt" example:"2023-01-01T00:00:00Z"`
}

type UserListResponse struct {
	ID            string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name          string    `json:"name" example:"john_doe"`
	Email         string    `json:"email" example:"john.doe@example.com"`
	FullName      string    `json:"fullName" example:"John Doe"`
	Role          UserRole  `json:"role" example:"Admin"`
	EmployeeID    *string   `json:"employeeId" example:"EMP001"`
	PreferredLang string    `json:"preferredLang" example:"en"`
	IsActive      bool      `json:"isActive" example:"true"`
	AvatarURL     *string   `json:"avatarUrl" example:"https://example.com/avatar.jpg"`
	CreatedAt     time.Time `json:"createdAt" example:"2023-01-01T00:00:00Z"`
	UpdatedAt     time.Time `json:"updatedAt" example:"2023-01-01T00:00:00Z"`
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
	IsActive      bool     `json:"isActive" example:"true" form:"isActive" validate:"required"`
	AvatarURL     *string  `json:"avatarUrl,omitempty" example:"https://example.com/avatar.jpg" form:"avatarUrl" validate:"omitempty,url"`
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

type ForgotPasswordResponse struct {
	Message string `json:"message" example:"Reset code sent to your email"`
}

type VerifyResetCodeResponse struct {
	Valid bool `json:"valid" example:"true"`
}

type ResetPasswordResponse struct {
	Message string `json:"message" example:"Password reset successfully"`
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
