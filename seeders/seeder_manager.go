package seeders

import (
	"context"
	"fmt"

	"github.com/Rizz404/inventory-api/services/auth"
	"github.com/Rizz404/inventory-api/services/category"
	"github.com/Rizz404/inventory-api/services/location"
	"github.com/Rizz404/inventory-api/services/user"
)

// SeederManager manages all seeders
type SeederManager struct {
	userSeeder     *UserSeeder
	categorySeeder *CategorySeeder
	locationSeeder *LocationSeeder
}

// NewSeederManager creates a new seeder manager
func NewSeederManager(
	authService auth.Service,
	userService user.UserService,
	categoryService category.CategoryService,
	locationService location.LocationService,
) *SeederManager {
	return &SeederManager{
		userSeeder:     NewUserSeeder(authService, userService),
		categorySeeder: NewCategorySeeder(categoryService),
		locationSeeder: NewLocationSeeder(locationService),
	}
}

// SeedUsers seeds user data
func (sm *SeederManager) SeedUsers(ctx context.Context, count int) error {
	fmt.Printf("📋 Starting user seeding (count: %d)...\n", count)
	return sm.userSeeder.Seed(ctx, count)
}

// SeedCategories seeds category data with parent-child hierarchy
func (sm *SeederManager) SeedCategories(ctx context.Context, totalCount int) error {
	// Calculate parent and children counts
	parentCount := getParentCount(totalCount)
	childrenCount := totalCount - parentCount

	fmt.Printf("📋 Starting category seeding (parents: %d, children: %d)...\n", parentCount, childrenCount)
	return sm.categorySeeder.Seed(ctx, parentCount, childrenCount)
}

// SeedLocations seeds location data
func (sm *SeederManager) SeedLocations(ctx context.Context, count int) error {
	fmt.Printf("📋 Starting location seeding (count: %d)...\n", count)
	return sm.locationSeeder.Seed(ctx, count)
}

// SeedAll seeds all data in the correct order
func (sm *SeederManager) SeedAll(ctx context.Context, count int) error {
	fmt.Println("🌱 Starting comprehensive seeding...")

	// 1. Seed users first (they might be referenced by other entities)
	fmt.Printf("\n1️⃣ Seeding users (count: %d)...\n", count)
	if err := sm.SeedUsers(ctx, count); err != nil {
		return fmt.Errorf("failed to seed users: %v", err)
	}

	// 2. Seed categories
	fmt.Printf("\n2️⃣ Seeding categories (total: %d)...\n", count)
	if err := sm.SeedCategories(ctx, count); err != nil {
		return fmt.Errorf("failed to seed categories: %v", err)
	}

	// 3. Seed locations
	fmt.Printf("\n3️⃣ Seeding locations (count: %d)...\n", count)
	if err := sm.SeedLocations(ctx, count); err != nil {
		return fmt.Errorf("failed to seed locations: %v", err)
	}

	fmt.Println("\n🎉 All seeding completed successfully!")
	return nil
}

// getParentCount returns the number of parent categories
// For categories, we want fewer parents and more children
func getParentCount(totalCount int) int {
	if totalCount <= 5 {
		return totalCount
	}
	// Approximately 25% of total count as parents, minimum 3
	parentCount := totalCount / 4
	if parentCount < 3 {
		parentCount = 3
	}
	return parentCount
}
