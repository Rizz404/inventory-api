package seeders

import (
	"context"
	"fmt"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/services/asset"
	"github.com/Rizz404/inventory-api/services/asset_movement"
	"github.com/Rizz404/inventory-api/services/auth"
	"github.com/Rizz404/inventory-api/services/category"
	"github.com/Rizz404/inventory-api/services/issue_report"
	"github.com/Rizz404/inventory-api/services/location"
	"github.com/Rizz404/inventory-api/services/maintenance_record"
	"github.com/Rizz404/inventory-api/services/maintenance_schedule"
	"github.com/Rizz404/inventory-api/services/user"
)

// SeederManager manages all seeders
type SeederManager struct {
	userSeeder                *UserSeeder
	categorySeeder            *CategorySeeder
	locationSeeder            *LocationSeeder
	assetSeeder               *AssetSeeder
	assetMovementSeeder       *AssetMovementSeeder
	issueReportSeeder         *IssueReportSeeder
	maintenanceScheduleSeeder *MaintenanceScheduleSeeder
	maintenanceRecordSeeder   *MaintenanceRecordSeeder
}

// NewSeederManager creates a new seeder manager
func NewSeederManager(
	authService auth.Service,
	userService user.UserService,
	categoryService category.CategoryService,
	locationService location.LocationService,
	assetService asset.AssetService,
	assetMovementService asset_movement.AssetMovementService,
	issueReportService issue_report.IssueReportService,
	maintenanceScheduleService maintenance_schedule.MaintenanceScheduleService,
	maintenanceRecordService maintenance_record.MaintenanceRecordService,
) *SeederManager {
	return &SeederManager{
		userSeeder:                NewUserSeeder(authService, userService),
		categorySeeder:            NewCategorySeeder(categoryService),
		locationSeeder:            NewLocationSeeder(locationService),
		assetSeeder:               NewAssetSeeder(assetService),
		assetMovementSeeder:       NewAssetMovementSeeder(assetMovementService),
		issueReportSeeder:         NewIssueReportSeeder(issueReportService),
		maintenanceScheduleSeeder: NewMaintenanceScheduleSeeder(maintenanceScheduleService),
		maintenanceRecordSeeder:   NewMaintenanceRecordSeeder(maintenanceRecordService),
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

// SeedAssets seeds asset data
func (sm *SeederManager) SeedAssets(ctx context.Context, count int, categoryIDs []string, locationIDs []string, userIDs []string) error {
	fmt.Printf("📋 Starting asset seeding (count: %d)...\n", count)
	return sm.assetSeeder.Seed(ctx, count, categoryIDs, locationIDs, userIDs)
}

// SeedAssetMovements seeds asset movement data
func (sm *SeederManager) SeedAssetMovements(ctx context.Context, count int, assetIDs []string, locationIDs []string, userIDs []string) error {
	fmt.Printf("📋 Starting asset movement seeding (count: %d)...\n", count)
	return sm.assetMovementSeeder.Seed(ctx, count, assetIDs, locationIDs, userIDs)
}

// SeedIssueReports seeds issue report data
func (sm *SeederManager) SeedIssueReports(ctx context.Context, count int, assetIDs []string, userIDs []string) error {
	fmt.Printf("📋 Starting issue report seeding (count: %d)...\n", count)
	return sm.issueReportSeeder.Seed(ctx, count, assetIDs, userIDs)
}

// SeedMaintenanceSchedules seeds maintenance schedule data
func (sm *SeederManager) SeedMaintenanceSchedules(ctx context.Context, count int, assetIDs []string, userIDs []string) error {
	fmt.Printf("📋 Starting maintenance schedule seeding (count: %d)...\n", count)
	return sm.maintenanceScheduleSeeder.Seed(ctx, count, assetIDs, userIDs)
}

// SeedMaintenanceRecords seeds maintenance record data
func (sm *SeederManager) SeedMaintenanceRecords(ctx context.Context, count int, assetIDs []string, scheduleIDs []string, userIDs []string) error {
	fmt.Printf("📋 Starting maintenance record seeding (count: %d)...\n", count)
	return sm.maintenanceRecordSeeder.Seed(ctx, count, assetIDs, scheduleIDs, userIDs)
}

// SeedAll seeds all data in the correct order with proper dependencies
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

	// Get IDs of seeded basic data for dependent seeding
	userIDs, err := sm.getUserIDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user IDs: %v", err)
	}

	categoryIDs, err := sm.getCategoryIDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get category IDs: %v", err)
	}

	locationIDs, err := sm.getLocationIDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get location IDs: %v", err)
	}

	// Validate we have required data
	if len(userIDs) == 0 {
		return fmt.Errorf("no users found for dependent seeding")
	}
	if len(categoryIDs) == 0 {
		return fmt.Errorf("no categories found for dependent seeding")
	}
	if len(locationIDs) == 0 {
		return fmt.Errorf("no locations found for dependent seeding")
	}

	// 4. Seed assets (requires users, categories, and locations)
	fmt.Printf("\n4️⃣ Seeding assets (count: %d)...\n", count)
	if err := sm.SeedAssets(ctx, count, categoryIDs, locationIDs, userIDs); err != nil {
		return fmt.Errorf("failed to seed assets: %v", err)
	}

	// Get asset IDs for dependent seeding
	assetIDs, err := sm.getAssetIDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get asset IDs: %v", err)
	}

	if len(assetIDs) == 0 {
		return fmt.Errorf("no assets found for dependent seeding")
	}

	// 5. Seed asset movements (requires assets, locations, and users)
	movementCount := count / 2 // Fewer movements than assets
	if movementCount < 5 {
		movementCount = 5
	}
	fmt.Printf("\n5️⃣ Seeding asset movements (count: %d)...\n", movementCount)
	if err := sm.SeedAssetMovements(ctx, movementCount, assetIDs, locationIDs, userIDs); err != nil {
		return fmt.Errorf("failed to seed asset movements: %v", err)
	}

	// 6. Seed maintenance schedules (requires assets and users)
	scheduleCount := count / 3 // Fewer schedules than assets
	if scheduleCount < 5 {
		scheduleCount = 5
	}
	fmt.Printf("\n6️⃣ Seeding maintenance schedules (count: %d)...\n", scheduleCount)
	if err := sm.SeedMaintenanceSchedules(ctx, scheduleCount, assetIDs, userIDs); err != nil {
		return fmt.Errorf("failed to seed maintenance schedules: %v", err)
	}

	// Get schedule IDs for maintenance records
	scheduleIDs, err := sm.getMaintenanceScheduleIDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get maintenance schedule IDs: %v", err)
	}

	// 7. Seed maintenance records (requires assets, optionally schedules, and users)
	recordCount := scheduleCount * 2 // More records than schedules
	fmt.Printf("\n7️⃣ Seeding maintenance records (count: %d)...\n", recordCount)
	if err := sm.SeedMaintenanceRecords(ctx, recordCount, assetIDs, scheduleIDs, userIDs); err != nil {
		return fmt.Errorf("failed to seed maintenance records: %v", err)
	}

	// 8. Seed issue reports (requires assets and users)
	issueCount := count / 4 // Fewer issues than assets
	if issueCount < 3 {
		issueCount = 3
	}
	fmt.Printf("\n8️⃣ Seeding issue reports (count: %d)...\n", issueCount)
	if err := sm.SeedIssueReports(ctx, issueCount, assetIDs, userIDs); err != nil {
		return fmt.Errorf("failed to seed issue reports: %v", err)
	}

	fmt.Printf("\n🎉 Comprehensive seeding completed successfully!")
	fmt.Printf("\n📊 Summary:")
	fmt.Printf("\n   - Users: %d", len(userIDs))
	fmt.Printf("\n   - Categories: %d", len(categoryIDs))
	fmt.Printf("\n   - Locations: %d", len(locationIDs))
	fmt.Printf("\n   - Assets: %d", len(assetIDs))
	fmt.Printf("\n   - Asset Movements: %d", movementCount)
	fmt.Printf("\n   - Maintenance Schedules: %d", scheduleCount)
	fmt.Printf("\n   - Maintenance Records: %d", recordCount)
	fmt.Printf("\n   - Issue Reports: %d\n", issueCount)

	return nil
}

// Helper methods to get IDs of seeded data using simple database queries
func (sm *SeederManager) getUserIDs(ctx context.Context) ([]string, error) {
	// Create query params to get all users (large limit)
	params := domain.UserParams{
		Pagination: &domain.UserPaginationOptions{
			Limit:  1000, // Large enough to get all seeded users
			Offset: 0,
		},
	}

	users, _, err := sm.userSeeder.userService.GetUsersPaginated(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}

	var userIDs []string
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}
	return userIDs, nil
}

func (sm *SeederManager) getCategoryIDs(ctx context.Context) ([]string, error) {
	params := domain.CategoryParams{
		Pagination: &domain.CategoryPaginationOptions{
			Limit:  1000,
			Offset: 0,
		},
	}

	categories, _, err := sm.categorySeeder.categoryService.GetCategoriesPaginated(ctx, params, "en-US")
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %v", err)
	}

	var categoryIDs []string
	for _, category := range categories {
		categoryIDs = append(categoryIDs, category.ID)
	}
	return categoryIDs, nil
}

func (sm *SeederManager) getLocationIDs(ctx context.Context) ([]string, error) {
	params := domain.LocationParams{
		Pagination: &domain.LocationPaginationOptions{
			Limit:  1000,
			Offset: 0,
		},
	}

	locations, _, err := sm.locationSeeder.locationService.GetLocationsPaginated(ctx, params, "en-US")
	if err != nil {
		return nil, fmt.Errorf("failed to get locations: %v", err)
	}

	var locationIDs []string
	for _, location := range locations {
		locationIDs = append(locationIDs, location.ID)
	}
	return locationIDs, nil
}

func (sm *SeederManager) getAssetIDs(ctx context.Context) ([]string, error) {
	params := domain.AssetParams{
		Pagination: &domain.AssetPaginationOptions{
			Limit:  1000,
			Offset: 0,
		},
	}

	assets, _, err := sm.assetSeeder.assetService.GetAssetsPaginated(ctx, params, "en-US")
	if err != nil {
		return nil, fmt.Errorf("failed to get assets: %v", err)
	}

	var assetIDs []string
	for _, asset := range assets {
		assetIDs = append(assetIDs, asset.ID)
	}
	return assetIDs, nil
}

func (sm *SeederManager) getMaintenanceScheduleIDs(ctx context.Context) ([]string, error) {
	params := domain.MaintenanceScheduleParams{
		Pagination: &domain.MaintenanceSchedulePaginationOptions{
			Limit:  1000,
			Offset: 0,
		},
	}

	schedules, _, err := sm.maintenanceScheduleSeeder.maintenanceScheduleService.GetMaintenanceSchedulesPaginated(ctx, params, "en-US")
	if err != nil {
		return nil, fmt.Errorf("failed to get maintenance schedules: %v", err)
	}

	var scheduleIDs []string
	for _, schedule := range schedules {
		scheduleIDs = append(scheduleIDs, schedule.ID)
	}
	return scheduleIDs, nil
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
