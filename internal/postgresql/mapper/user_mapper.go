package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
)

func ToModelUser(d *domain.User) model.User {
	var employeeID *string
	if d.EmployeeID != nil {
		employeeID = d.EmployeeID
	}

	return model.User{
		Username:      d.Username,
		Email:         d.Email,
		PasswordHash:  d.PasswordHash,
		FullName:      d.FullName,
		Role:          d.Role,
		EmployeeID:    employeeID,
		PreferredLang: d.PreferredLang,
		IsActive:      d.IsActive,
	}
}

func ToDomainUser(m *model.User) domain.User {
	return domain.User{
		ID:            m.ID.String(),
		Username:      m.Username,
		Email:         m.Email,
		PasswordHash:  m.PasswordHash,
		FullName:      m.FullName,
		Role:          m.Role,
		EmployeeID:    m.EmployeeID,
		PreferredLang: m.PreferredLang,
		IsActive:      m.IsActive,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func ToDomainUserResponse(m *model.User) domain.UserResponse {
	return domain.UserResponse{
		ID:            m.ID.String(),
		Username:      m.Username,
		Email:         m.Email,
		FullName:      m.FullName,
		Role:          m.Role,
		EmployeeID:    m.EmployeeID,
		PreferredLang: m.PreferredLang,
		IsActive:      m.IsActive,
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

func ToModelUserUpdateMap(payload *domain.UpdateUserPayload) map[string]interface{} {
	updates := make(map[string]interface{})

	if payload.Username != nil {
		updates["username"] = *payload.Username
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

	return updates
}
