package seeders

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/services/user"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
)

// UserSeeder handles user data seeding
type UserSeeder struct {
	userService user.UserService
}

// NewUserSeeder creates a new user seeder
func NewUserSeeder(userService user.UserService) *UserSeeder {
	return &UserSeeder{
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
		Email:         "admin@gmail.com",
		Password:      "admin123",
		FullName:      "System Administrator",
		Role:          domain.RoleAdmin,
		EmployeeID:    nil,
		PreferredLang: stringPtr("en-US"),
		IsActive:      true,
		AvatarURL:     stringPtr("https://ui-avatars.com/api/?name=System+Administrator&size=150&background=007bff&color=fff&bold=true"),
		PhoneNumber:   stringPtr("+6281234567890"),
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
	languages := []string{"en-US", "id-ID", "ja-JP"}

	successCount := 0
	for i := 0; i < count; i++ {
		// Generate fake user data
		firstName := gofakeit.FirstName()
		lastName := gofakeit.LastName()
		// Create unique username with timestamp suffix
		timestamp := time.Now().UnixNano()
		username := generateUniqueUsername(firstName, lastName, i, timestamp)
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
			EmployeeID:    nil,
			PreferredLang: stringPtr(languages[rand.Intn(len(languages))]),
			IsActive:      rand.Intn(100) < 90, // 90% active
			AvatarURL:     stringPtr(generateAvatarURL(firstName, lastName)),
			PhoneNumber:   stringPtr(gofakeit.Phone()),
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

// generateUniqueUsername creates a unique username with timestamp
func generateUniqueUsername(firstName, lastName string, index int, timestamp int64) string {
	base := fmt.Sprintf("%s.%s",
		normalizeString(firstName),
		normalizeString(lastName))

	// Add index and a portion of timestamp for uniqueness
	suffix := timestamp % 100000 // Last 5 digits
	return fmt.Sprintf("%s%d", base, index+int(suffix))
}

// generateUniqueEmployeeID creates a unique employee ID based on role using UUID
func generateUniqueEmployeeID(role domain.UserRole) string {
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
	// Generate unique ID using UUID (take first 8 chars and convert to uppercase)
	uuidStr := strings.ReplaceAll(uuid.New().String(), "-", "")
	return fmt.Sprintf("%s%s", prefix, strings.ToUpper(uuidStr[:8]))
}

// generateAvatarURL creates an avatar URL using UI Avatars API
func generateAvatarURL(firstName, lastName string) string {
	fullName := fmt.Sprintf("%s+%s", firstName, lastName)

	colors := []string{"007bff", "28a745", "dc3545", "ffc107", "17a2b8", "6f42c1", "e83e8c", "fd7e14"}
	color := colors[rand.Intn(len(colors))]

	return fmt.Sprintf("https://ui-avatars.com/api/?name=%s&size=150&background=%s&color=fff&bold=true", fullName, color)
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
