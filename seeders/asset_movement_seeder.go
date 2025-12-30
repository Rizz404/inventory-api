package seeders

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/services/asset_movement"
)

// AssetMovementSeeder handles asset movement data seeding
type AssetMovementSeeder struct {
	assetMovementService asset_movement.AssetMovementService
}

// NewAssetMovementSeeder creates a new asset movement seeder
func NewAssetMovementSeeder(assetMovementService asset_movement.AssetMovementService) *AssetMovementSeeder {
	return &AssetMovementSeeder{
		assetMovementService: assetMovementService,
	}
}

// Seed creates fake asset movements
func (ams *AssetMovementSeeder) Seed(ctx context.Context, count int, assetIDs []string, locationIDs []string, userIDs []string) error {
	if len(assetIDs) == 0 {
		return fmt.Errorf("no asset IDs provided for asset movement seeding")
	}
	if len(locationIDs) == 0 {
		return fmt.Errorf("no location IDs provided for asset movement seeding")
	}
	if len(userIDs) == 0 {
		return fmt.Errorf("no user IDs provided for asset movement seeding")
	}

	// Seed random generator
	rand.Seed(time.Now().UnixNano())

	successCount := 0
	for i := 0; i < count; i++ {
		// ! Add small delay to avoid rapid-fire requests
		if i > 0 {
			time.Sleep(50 * time.Millisecond)
		}

		movementPayload := ams.generateAssetMovementPayload(assetIDs, locationIDs, userIDs)

		_, err := ams.assetMovementService.CreateAssetMovement(ctx, movementPayload, "en-US")
		if err != nil {
			fmt.Printf("   ⚠️ Failed to create asset movement %d: %v\n", i+1, err)
			continue
		}

		successCount++
		if (i+1)%10 == 0 || i == count-1 {
			fmt.Printf("   🚚 Created %d/%d asset movements\n", successCount, count)
		}
	}

	fmt.Printf("✅ Successfully created %d asset movements\n", successCount)
	return nil
}

// generateAssetMovementPayload generates fake asset movement data
func (ams *AssetMovementSeeder) generateAssetMovementPayload(assetIDs []string, locationIDs []string, userIDs []string) *domain.CreateAssetMovementPayload {
	// Select random asset
	assetID := assetIDs[rand.Intn(len(assetIDs))]

	// Generate movement scenario (simplified for the new payload structure)
	movementType := rand.Intn(2) // 0: move to location, 1: move to user

	var toLocationID, toUserID *string

	switch movementType {
	case 0: // Move to location
		toLoc := locationIDs[rand.Intn(len(locationIDs))]
		toLocationID = &toLoc

	case 1: // Move to user
		toUser := userIDs[rand.Intn(len(userIDs))]
		toUserID = &toUser
	}

	// Generate movement notes
	movementReasons := []string{
		"Asset relocation for office restructuring",
		"Equipment maintenance transfer",
		"Employee assignment change",
		"Department restructuring",
		"Temporary loan to another department",
		"Return from maintenance",
		"New employee assignment",
		"Office space optimization",
		"Project requirement",
		"Equipment upgrade replacement",
	}

	notes := movementReasons[rand.Intn(len(movementReasons))]

	translations := []domain.CreateAssetMovementTranslationPayload{
		{
			LangCode: "en-US",
			Notes:    notes,
		},
		{
			LangCode: "id-ID",
			Notes:    notes,
		},
		{
			LangCode: "ja-JP",
			Notes:    notes,
		},
	}

	// Determine target based on movement type
	var targetLocationID, targetUserID *string
	if toLocationID != nil {
		targetLocationID = toLocationID
	}
	if toUserID != nil {
		targetUserID = toUserID
	}

	return &domain.CreateAssetMovementPayload{
		AssetID:      assetID,
		ToLocationID: targetLocationID,
		ToUserID:     targetUserID,
		Translations: translations,
	}
}
