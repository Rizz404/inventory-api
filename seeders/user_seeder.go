package seeders

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/services/auth"
	"github.com/Rizz404/inventory-api/services/user"
	"github.com/brianvoe/gofakeit/v6"
)

// UserSeeder handles user data seeding
type UserSeeder struct {
	authService auth.Service
	userService user.UserService
}

// NewUserSeeder creates a new user seeder
func NewUserSeeder(authService auth.Service, userService user.UserService) *UserSeeder {
	return &UserSeeder{
		authService: authService,
		userService: userService,
	}
}

// Seed creates fake users
func (us *UserSeeder) Seed(ctx context.Context, count int) error {
	// Seed random generator
	rand.Seed(time.Now().UnixNano())

	// Always create at least one admin user
	if err := us.createAdminUser(ctx); err != nil {
		return fmt.Errorf("failed to create admin user: %v", err)
	}

	// Create remaining users
	remainingCount := count - 1
	if remainingCount > 0 {
		if err := us.createRandomUsers(ctx, remainingCount); err != nil {
			return fmt.Errorf("failed to create random users: %v", err)
		}
	}

	fmt.Printf("✅ Successfully created %d users (1 admin + %d random users)\n", count, remainingCount)
	return nil
}

// createAdminUser creates a default admin user
func (us *UserSeeder) createAdminUser(ctx context.Context) error {
	adminPayload := &domain.CreateUserPayload{
		Name:          "admin",
		Email:         "admin@inventory.com",
		Password:      "admin123456",
		FullName:      "System Administrator",
		Role:          domain.RoleAdmin,
		EmployeeID:    stringPtr("ADM001"),
		PreferredLang: stringPtr("en-US"),
		IsActive:      true,
		AvatarURL:     stringPtr("https://via.placeholder.com/150/007bff/ffffff?text=AD"),
	}

	_, err := us.userService.CreateUser(ctx, adminPayload, nil)
	if err != nil {
		// If user already exists, that's fine
		return nil
	}

	fmt.Println("   ✅ Admin user created (admin@inventory.com / admin123456)")
	return nil
}

// createRandomUsers creates random users with different roles
func (us *UserSeeder) createRandomUsers(ctx context.Context, count int) error {
	languages := []string{"en-US", "id-ID"}

	successCount := 0
	for i := 0; i < count; i++ {
		// Generate fake user data
		firstName := gofakeit.FirstName()
		lastName := gofakeit.LastName()
		username := generateUsername(firstName, lastName, i)
		email := fmt.Sprintf("%s@%s", username, gofakeit.DomainName())

		// Select random role (more employees and staff than admins)
		var role domain.UserRole
		roleRand := rand.Intn(100)
		if roleRand < 10 { // 10% admin
			role = domain.RoleAdmin
		} else if roleRand < 40 { // 30% staff
			role = domain.RoleStaff
		} else { // 60% employee
			role = domain.RoleEmployee
		}

		userPayload := &domain.CreateUserPayload{
			Name:          username,
			Email:         email,
			Password:      "password123",
			FullName:      fmt.Sprintf("%s %s", firstName, lastName),
			Role:          role,
			EmployeeID:    stringPtr(generateEmployeeID(role, i)),
			PreferredLang: stringPtr(languages[rand.Intn(len(languages))]),
			IsActive:      rand.Intn(100) < 90, // 90% active
			AvatarURL:     stringPtr(generateAvatarURL(firstName, lastName)),
		}

		_, err := us.userService.CreateUser(ctx, userPayload, nil)
		if err != nil {
			fmt.Printf("   ⚠️  Failed to create user %s: %v\n", username, err)
			continue
		}

		successCount++
		if successCount%10 == 0 {
			fmt.Printf("   📊 Created %d/%d users...\n", successCount, count)
		}
	}

	if successCount < count {
		fmt.Printf("   ⚠️  Created %d out of %d requested users\n", successCount, count)
	}

	return nil
}

// generateUsername creates a unique username
func generateUsername(firstName, lastName string, index int) string {
	base := fmt.Sprintf("%s.%s",
		normalizeString(firstName),
		normalizeString(lastName))

	if index > 0 {
		return fmt.Sprintf("%s%d", base, index)
	}
	return base
}

// generateEmployeeID creates an employee ID based on role
func generateEmployeeID(role domain.UserRole, index int) string {
	var prefix string
	switch role {
	case domain.RoleAdmin:
		prefix = "ADM"
	case domain.RoleStaff:
		prefix = "STF"
	case domain.RoleEmployee:
		prefix = "EMP"
	default:
		prefix = "USR"
	}
	return fmt.Sprintf("%s%03d", prefix, index+100)
}

// generateAvatarURL creates a placeholder avatar URL
func generateAvatarURL(firstName, lastName string) string {
	initials := fmt.Sprintf("%c%c",
		normalizeString(firstName)[0],
		normalizeString(lastName)[0])

	colors := []string{"007bff", "28a745", "dc3545", "ffc107", "17a2b8", "6f42c1", "e83e8c", "fd7e14"}
	color := colors[rand.Intn(len(colors))]

	return fmt.Sprintf("https://via.placeholder.com/150/%s/ffffff?text=%s", color, initials)
}

// normalizeString converts string to lowercase and removes special characters
func normalizeString(s string) string {
	// Simple normalization - in production you might want more sophisticated handling
	result := ""
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			if r >= 'A' && r <= 'Z' {
				result += string(r + 32) // Convert to lowercase
			} else {
				result += string(r)
			}
		}
	}
	return result
}

// stringPtr returns a pointer to string
func stringPtr(s string) *string {
	return &s
}
