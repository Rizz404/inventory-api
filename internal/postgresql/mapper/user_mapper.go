package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
)

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

func ToDomainUserResponse(m *model.User) domain.UserResponse {
	return domain.UserResponse{
		ID:            m.ID.String(),
		Name:          m.Name,
		Email:         m.Email,
		FullName:      m.FullName,
		Role:          m.Role,
		EmployeeID:    m.EmployeeID,
		PreferredLang: m.PreferredLang,
		IsActive:      m.IsActive,
		AvatarURL:     m.AvatarURL,
		CreatedAt:     m.CreatedAt.Format(TimeFormat),
		UpdatedAt:     m.UpdatedAt.Format(TimeFormat),
	}
}

func ToDomainUsersResponse(m []model.User) []domain.UserResponse {
	responses := make([]domain.UserResponse, len(m))
	for i, user := range m {
		responses[i] = ToDomainUserResponse(&user)
	}
	return responses
}

// * Convert domain.User directly to domain.UserResponse without going through model.User
func DomainUserToUserResponse(d *domain.User) domain.UserResponse {
	return domain.UserResponse{
		ID:            d.ID,
		Name:          d.Name,
		Email:         d.Email,
		FullName:      d.FullName,
		Role:          d.Role,
		EmployeeID:    d.EmployeeID,
		PreferredLang: d.PreferredLang,
		IsActive:      d.IsActive,
		AvatarURL:     d.AvatarURL,
		CreatedAt:     d.CreatedAt.Format(TimeFormat),
		UpdatedAt:     d.UpdatedAt.Format(TimeFormat),
	}
}

// * Convert slice of domain.User to slice of domain.UserResponse
func DomainUsersToUsersResponse(users []domain.User) []domain.UserResponse {
	responses := make([]domain.UserResponse, len(users))
	for i, user := range users {
		responses[i] = DomainUserToUserResponse(&user)
	}
	return responses
}

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
