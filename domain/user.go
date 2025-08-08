package domain

import "time"

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
	Username      string    `json:"username"`
	PasswordHash  string    `json:"-"` // Hide password hash from JSON responses
	FullName      string    `json:"fullName"`
	Role          UserRole  `json:"role"`
	EmployeeID    *string   `json:"employeeId"`
	PreferredLang string    `json:"preferredLang"`
	IsActive      bool      `json:"isActive"`
	CreatedAt     time.Time `json:"createdAt"`
}

type UserResponse struct {
	ID            string   `json:"id"`
	Username      string   `json:"username"`
	FullName      string   `json:"fullName"`
	Role          UserRole `json:"role"`
	EmployeeID    *string  `json:"employeeId,omitempty"`
	PreferredLang string   `json:"preferredLang"`
	IsActive      bool     `json:"isActive"`
	CreatedAt     string   `json:"createdAt"`
}

type LoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
}

// --- Payloads ---

type LoginPayload struct {
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required,min=8"`
}

type CreateUserPayload struct {
	Username      string   `json:"username" form:"username" validate:"required,min=3,max=50"`
	Password      string   `json:"password" form:"password" validate:"required,min=8,max=100"`
	FullName      string   `json:"fullName" form:"fullName" validate:"required,min=3,max=100"`
	Role          UserRole `json:"role" form:"role" validate:"required,oneof=Admin Staff Employee"`
	EmployeeID    *string  `json:"employeeId,omitempty" form:"employeeId" validate:"omitempty,max=20"`
	PreferredLang *string  `json:"preferredLang,omitempty" form:"preferredLang" validate:"omitempty,max=5"`
}

type UpdateUserPayload struct {
	Username      *string   `json:"username,omitempty" form:"username" validate:"omitempty,min=3,max=50"`
	Password      *string   `json:"password,omitempty" form:"password" validate:"omitempty,min=8,max=100"`
	FullName      *string   `json:"fullName,omitempty" form:"fullName" validate:"omitempty,min=3,max=100"`
	Role          *UserRole `json:"role,omitempty" form:"role" validate:"omitempty,oneof=Admin Staff Employee"`
	EmployeeID    *string   `json:"employeeId,omitempty" form:"employeeId" validate:"omitempty,max=20"`
	PreferredLang *string   `json:"preferredLang,omitempty" form:"preferredLang" validate:"omitempty,max=5"`
	IsActive      *bool     `json:"isActive,omitempty" form:"isActive" validate:"omitempty"`
}
