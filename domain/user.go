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

type UserResponse struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Email         string   `json:"email"`
	FullName      string   `json:"fullName"`
	Role          UserRole `json:"role"`
	EmployeeID    *string  `json:"employeeId,omitempty"` // ! gak usah diapa-apain dulu, soalnya belum ada
	PreferredLang string   `json:"preferredLang"`
	IsActive      bool     `json:"isActive"`
	AvatarURL     *string  `json:"avatarUrl,omitempty"`
	CreatedAt     string   `json:"createdAt"`
	UpdatedAt     string   `json:"updatedAt"`
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
