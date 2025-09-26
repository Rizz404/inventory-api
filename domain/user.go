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
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// ! jangan omitempty biar client nya tau
type UserResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	FullName      string    `json:"fullName"`
	Role          UserRole  `json:"role"`
	EmployeeID    *string   `json:"employeeId"` // ! gak usah diapa-apain dulu, soalnya belum ada
	PreferredLang string    `json:"preferredLang"`
	IsActive      bool      `json:"isActive"`
	AvatarURL     *string   `json:"avatarUrl"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type UserListResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	FullName      string    `json:"fullName"`
	Role          UserRole  `json:"role"`
	EmployeeID    *string   `json:"employeeId"`
	PreferredLang string    `json:"preferredLang"`
	IsActive      bool      `json:"isActive"`
	AvatarURL     *string   `json:"avatarUrl"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type LoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
}

// --- Payloads ---

type LoginPayload struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=5"`
}

type RegisterPayload struct {
	Name     string `json:"name" form:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=5"`
}

type CreateUserPayload struct {
	Name          string   `json:"name" form:"name" validate:"required,min=3,max=50"`
	Email         string   `json:"email" form:"email" validate:"required,email,max=255"`
	Password      string   `json:"password" form:"password" validate:"required,min=8,max=100"`
	FullName      string   `json:"fullName" form:"fullName" validate:"required,min=3,max=100"`
	Role          UserRole `json:"role" form:"role" validate:"required,oneof=Admin Staff Employee"`
	EmployeeID    *string  `json:"employeeId,omitempty" form:"employeeId" validate:"omitempty,max=20"` // ! gak usah diapa-apain dulu, soalnya belum ada
	PreferredLang *string  `json:"preferredLang,omitempty" form:"preferredLang" validate:"omitempty,max=5"`
	IsActive      bool     `json:"isActive" form:"isActive" validate:"required"`
	AvatarURL     *string  `json:"avatarUrl,omitempty" form:"avatarUrl" validate:"omitempty,url"`
}

type UpdateUserPayload struct {
	Name          *string   `json:"name,omitempty" form:"name" validate:"omitempty,min=3,max=50"`
	Email         *string   `json:"email,omitempty" form:"email" validate:"omitempty,email,max=255"`
	Password      *string   `json:"password,omitempty" form:"password" validate:"omitempty,min=8,max=100"`
	FullName      *string   `json:"fullName,omitempty" form:"fullName" validate:"omitempty,min=3,max=100"`
	Role          *UserRole `json:"role,omitempty" form:"role" validate:"omitempty,oneof=Admin Staff Employee"`
	EmployeeID    *string   `json:"employeeId,omitempty" form:"employeeId" validate:"omitempty,max=20"` // ! gak usah diapa-apain dulu, soalnya belum ada
	PreferredLang *string   `json:"preferredLang,omitempty" form:"preferredLang" validate:"omitempty,max=5"`
	IsActive      *bool     `json:"isActive,omitempty" form:"isActive" validate:"omitempty"`
	AvatarURL     *string   `json:"avatarUrl,omitempty" form:"avatarUrl" validate:"omitempty,url"`
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
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type UserSummaryStatistics struct {
	TotalUsers               int     `json:"totalUsers"`
	ActiveUsersPercentage    float64 `json:"activeUsersPercentage"`
	InactiveUsersPercentage  float64 `json:"inactiveUsersPercentage"`
	AdminPercentage          float64 `json:"adminPercentage"`
	StaffPercentage          float64 `json:"staffPercentage"`
	EmployeePercentage       float64 `json:"employeePercentage"`
	AverageUsersPerDay       float64 `json:"averageUsersPerDay"`
	LatestRegistrationDate   string  `json:"latestRegistrationDate"`
	EarliestRegistrationDate string  `json:"earliestRegistrationDate"`
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
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type UserSummaryStatisticsResponse struct {
	TotalUsers               int     `json:"totalUsers"`
	ActiveUsersPercentage    float64 `json:"activeUsersPercentage"`
	InactiveUsersPercentage  float64 `json:"inactiveUsersPercentage"`
	AdminPercentage          float64 `json:"adminPercentage"`
	StaffPercentage          float64 `json:"staffPercentage"`
	EmployeePercentage       float64 `json:"employeePercentage"`
	AverageUsersPerDay       float64 `json:"averageUsersPerDay"`
	LatestRegistrationDate   string  `json:"latestRegistrationDate"`
	EarliestRegistrationDate string  `json:"earliestRegistrationDate"`
}
