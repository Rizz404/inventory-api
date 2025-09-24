package docs

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/web"
)

// Swagger model definitions for better documentation

// RegisterRequest represents the registration request payload
// swagger:model RegisterRequest
type RegisterRequest struct {
	// User's display name
	// required: true
	// min length: 3
	// max length: 50
	// example: johndoe
	Name string `json:"name" validate:"required,min=3,max=50"`

	// User's email address
	// required: true
	// format: email
	// example: john.doe@example.com
	Email string `json:"email" validate:"required,email"`

	// User's password
	// required: true
	// min length: 5
	// example: password123
	Password string `json:"password" validate:"required,min=5"`
}

// LoginRequest represents the login request payload
// swagger:model LoginRequest
type LoginRequest struct {
	// User's email address
	// required: true
	// format: email
	// example: john.doe@example.com
	Email string `json:"email" validate:"required,email"`

	// User's password
	// required: true
	// min length: 5
	// example: password123
	Password string `json:"password" validate:"required,min=5"`
}

// UserResponse represents user data in responses
// swagger:model UserResponse
type UserResponse struct {
	// User ID
	// example: 550e8400-e29b-41d4-a716-446655440000
	ID string `json:"id"`

	// User's display name
	// example: johndoe
	Name string `json:"name"`

	// User's email address
	// example: john.doe@example.com
	Email string `json:"email"`

	// User's full name
	// example: John Doe
	FullName string `json:"fullName"`

	// User's role in the system
	// enum: Admin,Staff,Employee
	// example: Employee
	Role domain.UserRole `json:"role"`

	// Employee ID (nullable)
	// example: EMP001
	EmployeeID *string `json:"employeeId"`

	// User's preferred language
	// example: en
	PreferredLang string `json:"preferredLang"`

	// Whether the user account is active
	// example: true
	IsActive bool `json:"isActive"`

	// URL to user's avatar image (nullable)
	// example: https://example.com/avatar.jpg
	AvatarURL *string `json:"avatarUrl"`

	// Account creation timestamp
	// example: 2023-01-01 12:00:00
	CreatedAt string `json:"createdAt"`

	// Last update timestamp
	// example: 2023-01-01 12:00:00
	UpdatedAt string `json:"updatedAt"`
}

// LoginResponse represents successful login response
// swagger:model LoginResponse
type LoginResponse struct {
	// User information
	User UserResponse `json:"user"`

	// JWT access token
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	AccessToken string `json:"accessToken"`

	// JWT refresh token
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	RefreshToken string `json:"refreshToken"`
}

// SuccessResponse represents successful API response
// swagger:model SuccessResponse
type SuccessResponse struct {
	// Response status
	// example: success
	Status string `json:"status"`

	// Response message
	// example: Operation completed successfully
	Message string `json:"message"`

	// Response data (varies by endpoint)
	Data interface{} `json:"data,omitempty"`

	// Pagination information (if applicable)
	PageInfo *web.PageInfo `json:"pagination,omitempty"`

	// Cursor pagination information (if applicable)
	CursorInfo *web.CursorInfo `json:"cursor,omitempty"`
}

// ErrorResponse represents error API response
// swagger:model ErrorResponse
type ErrorResponse struct {
	// Response status
	// example: error
	Status string `json:"status"`

	// Error message
	// example: Validation failed
	Message string `json:"message"`

	// Error details (varies by error type)
	Error interface{} `json:"error,omitempty"`
}

// ValidationError represents validation error details
// swagger:model ValidationError
type ValidationError struct {
	// Field name that failed validation
	// example: email
	Field string `json:"field"`

	// Validation tag that failed
	// example: required
	Tag string `json:"tag"`

	// Value that failed validation
	// example: invalid-email
	Value string `json:"value"`

	// Human-readable error message
	// example: email is required
	Message string `json:"message"`
}

// ValidationErrorResponse represents validation error response
// swagger:model ValidationErrorResponse
type ValidationErrorResponse struct {
	// Response status
	// example: error
	Status string `json:"status"`

	// Error message
	// example: Validation failed
	Message string `json:"message"`

	// Array of validation errors
	Error []ValidationError `json:"error"`
}

// PageInfo represents pagination information
// swagger:model PageInfo
type PageInfo struct {
	// Total number of items
	// example: 100
	Total int `json:"total"`

	// Number of items per page
	// example: 10
	PerPage int `json:"per_page"`

	// Current page number
	// example: 1
	CurrentPage int `json:"current_page"`

	// Total number of pages
	// example: 10
	TotalPages int `json:"total_pages"`

	// Whether there is a previous page
	// example: false
	HasPrevPage bool `json:"has_prev_page"`

	// Whether there is a next page
	// example: true
	HasNextPage bool `json:"has_next_page"`
}

// CursorInfo represents cursor-based pagination information
// swagger:model CursorInfo
type CursorInfo struct {
	// Next cursor for pagination
	// example: eyJpZCI6IjEwIiwiY3JlYXRlZEF0IjoiMjAyMy0wMS0wMSJ9
	NextCursor string `json:"next_cursor"`

	// Whether there is a next page
	// example: true
	HasNextPage bool `json:"has_next_page"`

	// Number of items per page
	// example: 10
	PerPage int `json:"per_page"`

	// Total number of items (optional)
	// example: 100
	Total int `json:"total,omitempty"`
}
