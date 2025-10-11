package seeders

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/services/asset"
	"github.com/brianvoe/gofakeit/v6"
)

// AssetSeeder handles asset data seeding
type AssetSeeder struct {
	assetService asset.AssetService
}

// NewAssetSeeder creates a new asset seeder
func NewAssetSeeder(assetService asset.AssetService) *AssetSeeder {
	return &AssetSeeder{
		assetService: assetService,
	}
}

// Seed creates fake assets
func (as *AssetSeeder) Seed(ctx context.Context, count int, categoryIDs []string, locationIDs []string, userIDs []string) error {
	if len(categoryIDs) == 0 {
		return fmt.Errorf("no category IDs provided for asset seeding")
	}
	if len(locationIDs) == 0 {
		return fmt.Errorf("no location IDs provided for asset seeding")
	}

	// Seed random generator
	rand.Seed(time.Now().UnixNano())

	successCount := 0
	for i := 0; i < count; i++ {
		assetPayload := as.generateAssetPayload(categoryIDs, locationIDs, userIDs, i)

		_, err := as.assetService.CreateAsset(ctx, assetPayload, nil, "en-US")
		if err != nil {
			fmt.Printf("   ⚠️ Failed to create asset %d: %v\n", i+1, err)
			continue
		}

		successCount++
		if (i+1)%10 == 0 || i == count-1 {
			fmt.Printf("   📦 Created %d/%d assets\n", successCount, count)
		}
	}

	fmt.Printf("✅ Successfully created %d assets\n", successCount)
	return nil
}

// generateAssetPayload generates fake asset data
func (as *AssetSeeder) generateAssetPayload(categoryIDs []string, locationIDs []string, userIDs []string, index int) *domain.CreateAssetPayload {
	// Asset types and names
	assetTypes := []string{
		"Laptop", "Desktop Computer", "Monitor", "Printer", "Scanner", "Projector",
		"Office Chair", "Office Desk", "Filing Cabinet", "Whiteboard",
		"Server", "Network Switch", "Router", "UPS", "Air Conditioner",
		"Vehicle", "Forklift", "Crane", "Generator", "Industrial Machine",
	}

	brands := []string{
		"Dell", "HP", "Lenovo", "Apple", "Asus", "Acer", "Canon", "Epson",
		"Samsung", "LG", "Sony", "Panasonic", "Cisco", "APC", "Toyota", "Honda",
	}

	assetType := assetTypes[rand.Intn(len(assetTypes))]
	brand := brands[rand.Intn(len(brands))]

	// Generate asset tag with prefix and number
	assetTag := fmt.Sprintf("AST-%06d", index+1)

	// Generate purchase data
	purchaseDate := gofakeit.DateRange(time.Now().AddDate(-5, 0, 0), time.Now().AddDate(-1, 0, 0))
	warrantyEnd := purchaseDate.AddDate(rand.Intn(3)+1, 0, 0) // 1-3 years warranty
	purchaseDateStr := purchaseDate.Format("2006-01-02")
	warrantyEndStr := warrantyEnd.Format("2006-01-02")

	// Generate price based on asset type
	var priceRange [2]int
	switch {
	case contains([]string{"Laptop", "Desktop Computer", "Server"}, assetType):
		priceRange = [2]int{500, 5000}
	case contains([]string{"Monitor", "Printer", "Scanner"}, assetType):
		priceRange = [2]int{100, 1000}
	case contains([]string{"Office Chair", "Office Desk"}, assetType):
		priceRange = [2]int{50, 500}
	case contains([]string{"Vehicle", "Forklift", "Crane"}, assetType):
		priceRange = [2]int{10000, 100000}
	default:
		priceRange = [2]int{100, 2000}
	}

	purchasePrice := float64(rand.Intn(priceRange[1]-priceRange[0]+1) + priceRange[0])

	// Random status and condition
	statuses := []domain.AssetStatus{
		domain.StatusActive, domain.StatusActive, domain.StatusActive, // Higher chance of active
		domain.StatusMaintenance, domain.StatusDisposed, domain.StatusLost,
	}
	conditions := []domain.AssetCondition{
		domain.ConditionGood, domain.ConditionGood, domain.ConditionFair, // Higher chance of good
		domain.ConditionPoor, domain.ConditionDamaged,
	}

	// Randomly assign to user (50% chance)
	var assignedTo *string
	if len(userIDs) > 0 && rand.Intn(2) == 0 {
		assignedTo = &userIDs[rand.Intn(len(userIDs))]
	}

	// Random location
	locationID := locationIDs[rand.Intn(len(locationIDs))]

	// Generate asset name
	assetName := fmt.Sprintf("%s %s %s", brand, assetType, gofakeit.LetterN(3))

	status := statuses[rand.Intn(len(statuses))]
	condition := conditions[rand.Intn(len(conditions))]

	return &domain.CreateAssetPayload{
		AssetTag:      assetTag,
		AssetName:     assetName,
		CategoryID:    categoryIDs[rand.Intn(len(categoryIDs))],
		Brand:         &brand,
		Model:         stringPtr(gofakeit.CarModel()),
		SerialNumber:  stringPtr(gofakeit.LetterN(10) + gofakeit.DigitN(6)),
		PurchaseDate:  &purchaseDateStr,
		PurchasePrice: &purchasePrice,
		VendorName:    stringPtr(gofakeit.Company()),
		WarrantyEnd:   &warrantyEndStr,
		Status:        status,
		Condition:     condition,
		LocationID:    &locationID,
		AssignedTo:    assignedTo,
	}
}

// Helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
