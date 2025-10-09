package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Rizz404/inventory-api/internal/client/cloudinary"
	"github.com/Rizz404/inventory-api/internal/postgresql"
	"github.com/Rizz404/inventory-api/seeders"
	"github.com/Rizz404/inventory-api/services/asset"
	"github.com/Rizz404/inventory-api/services/asset_movement"
	"github.com/Rizz404/inventory-api/services/auth"
	"github.com/Rizz404/inventory-api/services/category"
	"github.com/Rizz404/inventory-api/services/issue_report"
	"github.com/Rizz404/inventory-api/services/location"
	"github.com/Rizz404/inventory-api/services/maintenance_record"
	"github.com/Rizz404/inventory-api/services/maintenance_schedule"
	"github.com/Rizz404/inventory-api/services/user"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, using system environment variables")
	}
}

func main() {
	var (
		seedType = flag.String("type", "all", "Type of seed to run: users, categories, locations, assets, movements, issues, schedules, records, or all")
		count    = flag.Int("count", 20, "Number of records to create (default: 20)")
		help     = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// Validate seed type
	validTypes := []string{"users", "categories", "locations", "assets", "movements", "issues", "schedules", "records", "all"}
	if !contains(validTypes, *seedType) {
		fmt.Printf("Invalid seed type: %s\n", *seedType)
		fmt.Printf("Valid types: %s\n", strings.Join(validTypes, ", "))
		os.Exit(1)
	}

	// Validate count
	if *count <= 0 {
		fmt.Println("Count must be greater than 0")
		os.Exit(1)
	}

	// Initialize database
	db, err := initDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize services
	services := initServices(db)

	// Initialize seeder manager
	seederManager := seeders.NewSeederManager(
		services.Auth,
		services.User,
		services.Category,
		services.Location,
		services.Asset,
		services.AssetMovement,
		services.IssueReport,
		services.MaintenanceSchedule,
		services.MaintenanceRecord,
	)

	ctx := context.Background()

	// Run seeders based on type
	switch *seedType {
	case "users":
		fmt.Printf("Seeding %d users...\n", *count)
		if err := seederManager.SeedUsers(ctx, *count); err != nil {
			log.Fatalf("Failed to seed users: %v", err)
		}
		fmt.Println("✅ Users seeded successfully!")

	case "categories":
		fmt.Printf("Seeding categories (parents: %d, children per parent: %d)...\n", getParentCount(*count), getChildrenPerParent(*count))
		if err := seederManager.SeedCategories(ctx, *count); err != nil {
			log.Fatalf("Failed to seed categories: %v", err)
		}
		fmt.Println("✅ Categories seeded successfully!")

	case "locations":
		fmt.Printf("Seeding %d locations...\n", *count)
		if err := seederManager.SeedLocations(ctx, *count); err != nil {
			log.Fatalf("Failed to seed locations: %v", err)
		}
		fmt.Println("✅ Locations seeded successfully!")

	case "assets":
		fmt.Println("⚠️ Assets seeding requires existing users, categories, and locations.")
		fmt.Println("Please make sure you have run 'users', 'categories', and 'locations' seeders first.")
		fmt.Println("Example: go run cmd/seed/main.go -type=users && go run cmd/seed/main.go -type=categories && go run cmd/seed/main.go -type=locations && go run cmd/seed/main.go -type=assets")

	case "movements":
		fmt.Println("⚠️ Asset movements seeding requires existing assets, users, and locations.")
		fmt.Println("Please make sure you have run other seeders first.")

	case "issues":
		fmt.Println("⚠️ Issue reports seeding requires existing assets and users.")
		fmt.Println("Please make sure you have run other seeders first.")

	case "schedules":
		fmt.Println("⚠️ Maintenance schedules seeding requires existing assets and users.")
		fmt.Println("Please make sure you have run other seeders first.")

	case "records":
		fmt.Println("⚠️ Maintenance records seeding requires existing assets, users, and optionally schedules.")
		fmt.Println("Please make sure you have run other seeders first.")

	case "all":
		fmt.Printf("Seeding all data (count: %d)...\n", *count)
		if err := seederManager.SeedAll(ctx, *count); err != nil {
			log.Fatalf("Failed to seed all data: %v", err)
		}
		fmt.Println("✅ All data seeded successfully!")

	default:
		fmt.Printf("Unknown seed type: %s\n", *seedType)
		os.Exit(1)
	}
}

func initDatabase() (*gorm.DB, error) {
	DSN := os.Getenv("DSN")
	if DSN == "" {
		return nil, fmt.Errorf("DSN environment variable not set")
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: DSN,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to the database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic database object: %v", err)
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Printf("Successfully connected to database")
	return db, nil
}

type Services struct {
	Auth                auth.Service
	User                user.UserService
	Category            category.CategoryService
	Location            location.LocationService
	Asset               asset.AssetService
	AssetMovement       asset_movement.AssetMovementService
	IssueReport         issue_report.IssueReportService
	MaintenanceSchedule maintenance_schedule.MaintenanceScheduleService
	MaintenanceRecord   maintenance_record.MaintenanceRecordService
}

func initServices(db *gorm.DB) *Services {
	// Initialize Cloudinary client (optional for seeding)
	var cloudinaryClient *cloudinary.Client
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL != "" {
		var err error
		cloudinaryClient, err = cloudinary.NewClientFromURL(cloudinaryURL)
		if err != nil {
			log.Printf("Warning: Failed to initialize Cloudinary client: %v. Avatar URLs will be mock URLs.", err)
		}
	}

	// Initialize repositories
	userRepository := postgresql.NewUserRepository(db)
	categoryRepository := postgresql.NewCategoryRepository(db)
	locationRepository := postgresql.NewLocationRepository(db)
	assetRepository := postgresql.NewAssetRepository(db)
	assetMovementRepository := postgresql.NewAssetMovementRepository(db)
	issueReportRepository := postgresql.NewIssueReportRepository(db)
	maintenanceScheduleRepository := postgresql.NewMaintenanceScheduleRepository(db)
	maintenanceRecordRepository := postgresql.NewMaintenanceRecordRepository(db)

	// Initialize services
	authService := auth.NewService(userRepository)
	userService := user.NewService(userRepository, cloudinaryClient)
	categoryService := category.NewService(categoryRepository)
	locationService := location.NewService(locationRepository)
	assetService := asset.NewService(assetRepository, cloudinaryClient, nil, categoryService) // nil for notification service in seeder
	assetMovementService := asset_movement.NewService(assetMovementRepository, assetService, locationService, userService)
	issueReportService := issue_report.NewService(issueReportRepository)
	maintenanceScheduleService := maintenance_schedule.NewService(maintenanceScheduleRepository, assetService, userService)
	maintenanceRecordService := maintenance_record.NewService(maintenanceRecordRepository, assetService, userService)

	return &Services{
		Auth:                *authService,
		User:                userService,
		Category:            categoryService,
		Location:            locationService,
		Asset:               assetService,
		AssetMovement:       assetMovementService,
		IssueReport:         issueReportService,
		MaintenanceSchedule: maintenanceScheduleService,
		MaintenanceRecord:   maintenanceRecordService,
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// getParentCount returns the number of parent categories based on total count
// For categories, we want fewer parents and more children
func getParentCount(totalCount int) int {
	if totalCount <= 5 {
		return totalCount
	}
	// Approximately 25% of total count as parents
	return totalCount / 4
}

// getChildrenPerParent returns approximate children per parent
func getChildrenPerParent(totalCount int) int {
	parentCount := getParentCount(totalCount)
	if parentCount == 0 {
		return 0
	}
	return (totalCount - parentCount) / parentCount
}

func showHelp() {
	fmt.Println("Inventory API Seeder")
	fmt.Println("====================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/seed/main.go [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -type string")
	fmt.Println("        Type of seed to run:")
	fmt.Println("        - users: Seed user accounts")
	fmt.Println("        - categories: Seed asset categories")
	fmt.Println("        - locations: Seed location data")
	fmt.Println("        - assets: Seed assets (requires users, categories, locations)")
	fmt.Println("        - movements: Seed asset movements (requires assets, users, locations)")
	fmt.Println("        - issues: Seed issue reports (requires assets, users)")
	fmt.Println("        - schedules: Seed maintenance schedules (requires assets, users)")
	fmt.Println("        - records: Seed maintenance records (requires assets, users, schedules)")
	fmt.Println("        - all: Seed basic data (users, categories, locations only)")
	fmt.Println("        (default: all)")
	fmt.Println("  -count int")
	fmt.Println("        Number of records to create (default: 20)")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Seed basic data (users, categories, locations)")
	fmt.Println("  go run cmd/seed/main.go")
	fmt.Println()
	fmt.Println("  # Seed 50 users only")
	fmt.Println("  go run cmd/seed/main.go -type=users -count=50")
	fmt.Println()
	fmt.Println("  # Seed 30 categories (will create ~7 parents with ~3 children each)")
	fmt.Println("  go run cmd/seed/main.go -type=categories -count=30")
	fmt.Println()
	fmt.Println("  # Seed 40 locations")
	fmt.Println("  go run cmd/seed/main.go -type=locations -count=40")
	fmt.Println()
	fmt.Println("  # Complete seeding workflow:")
	fmt.Println("  go run cmd/seed/main.go -type=users -count=20")
	fmt.Println("  go run cmd/seed/main.go -type=categories -count=15")
	fmt.Println("  go run cmd/seed/main.go -type=locations -count=10")
	fmt.Println("  # Now you can seed dependent data:")
	fmt.Println("  go run cmd/seed/main.go -type=assets -count=50")
	fmt.Println("  go run cmd/seed/main.go -type=movements -count=25")
	fmt.Println("  go run cmd/seed/main.go -type=schedules -count=20")
	fmt.Println("  go run cmd/seed/main.go -type=records -count=30")
	fmt.Println("  go run cmd/seed/main.go -type=issues -count=15")
	fmt.Println()
	fmt.Println("Environment Variables Required:")
	fmt.Println("  DSN - PostgreSQL database connection string")
	fmt.Println("  CLOUDINARY_URL - Cloudinary URL (optional, for avatar uploads)")
	fmt.Println()
	fmt.Println("Note: Some seed types require existing data. Make sure to seed in the correct order:")
	fmt.Println("1. users, categories, locations (can be seeded independently)")
	fmt.Println("2. assets (requires users, categories, locations)")
	fmt.Println("3. movements, schedules, issues (requires assets and users)")
	fmt.Println("4. records (requires assets, users, and optionally schedules)")
}
