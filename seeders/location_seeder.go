package seeders

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/services/location"
)

// LocationSeeder handles location data seeding
type LocationSeeder struct {
	locationService location.LocationService
}

// NewLocationSeeder creates a new location seeder
func NewLocationSeeder(locationService location.LocationService) *LocationSeeder {
	return &LocationSeeder{
		locationService: locationService,
	}
}

// Seed creates fake locations
func (ls *LocationSeeder) Seed(ctx context.Context, count int) error {
	// Seed random generator
	rand.Seed(time.Now().UnixNano())

	// Create predefined locations first, then random ones
	createdCount := 0

	// Create some predefined realistic locations
	predefinedLocations := ls.getPredefinedLocations()

	// Create predefined locations first
	for i := 0; i < count && i < len(predefinedLocations); i++ {
		if err := ls.createLocation(ctx, predefinedLocations[i]); err != nil {
			fmt.Printf("   ⚠️  Failed to create predefined location %s: %v\n", predefinedLocations[i].LocationCode, err)
			continue
		}
		createdCount++
	}

	// Create random locations for remaining count
	for i := len(predefinedLocations); i < count; i++ {
		locationData := ls.generateRandomLocation(i)
		if err := ls.createLocation(ctx, locationData); err != nil {
			fmt.Printf("   ⚠️  Failed to create random location %s: %v\n", locationData.LocationCode, err)
			continue
		}
		createdCount++
	}

	fmt.Printf("✅ Successfully created %d locations\n", createdCount)
	return nil
}

// createLocation creates a single location
func (ls *LocationSeeder) createLocation(ctx context.Context, locationData *domain.CreateLocationPayload) error {
	_, err := ls.locationService.CreateLocation(ctx, locationData)
	if err != nil {
		return err
	}

	fmt.Printf("   ✅ Created location: %s (%s)\n", locationData.LocationCode, locationData.Translations[0].LocationName)
	return nil
}

// getPredefinedLocations returns a list of realistic predefined locations
func (ls *LocationSeeder) getPredefinedLocations() []*domain.CreateLocationPayload {
	return []*domain.CreateLocationPayload{
		{
			LocationCode: "HQ_LOBBY",
			Building:     stringPtr("Headquarters"),
			Floor:        stringPtr("Ground Floor"),
			Latitude:     float64Ptr(-6.2088),
			Longitude:    float64Ptr(106.8456),
			Translations: []domain.CreateLocationTranslationPayload{
				{LangCode: "en-US", LocationName: "Headquarters Lobby"},
				{LangCode: "id-ID", LocationName: "Lobi Kantor Pusat"},
				{LangCode: "ja-JP", LocationName: "本社ロビー"},
			},
		},
		{
			LocationCode: "HQ_IT_ROOM",
			Building:     stringPtr("Headquarters"),
			Floor:        stringPtr("2nd Floor"),
			Latitude:     float64Ptr(-6.2088),
			Longitude:    float64Ptr(106.8456),
			Translations: []domain.CreateLocationTranslationPayload{
				{LangCode: "en-US", LocationName: "IT Server Room"},
				{LangCode: "id-ID", LocationName: "Ruang Server IT"},
				{LangCode: "ja-JP", LocationName: "ITサーバールーム"},
			},
		},
		{
			LocationCode: "HQ_MEETING_A",
			Building:     stringPtr("Headquarters"),
			Floor:        stringPtr("3rd Floor"),
			Latitude:     float64Ptr(-6.2088),
			Longitude:    float64Ptr(106.8456),
			Translations: []domain.CreateLocationTranslationPayload{
				{LangCode: "en-US", LocationName: "Meeting Room A"},
				{LangCode: "id-ID", LocationName: "Ruang Rapat A"},
				{LangCode: "ja-JP", LocationName: "会議室A"},
			},
		},
		{
			LocationCode: "HQ_MEETING_B",
			Building:     stringPtr("Headquarters"),
			Floor:        stringPtr("3rd Floor"),
			Latitude:     float64Ptr(-6.2088),
			Longitude:    float64Ptr(106.8456),
			Translations: []domain.CreateLocationTranslationPayload{
				{LangCode: "en-US", LocationName: "Meeting Room B"},
				{LangCode: "id-ID", LocationName: "Ruang Rapat B"},
				{LangCode: "ja-JP", LocationName: "会議室B"},
			},
		},
		{
			LocationCode: "WAREHOUSE_A",
			Building:     stringPtr("Warehouse Complex"),
			Floor:        stringPtr("Ground Floor"),
			Latitude:     float64Ptr(-6.2297),
			Longitude:    float64Ptr(106.8278),
			Translations: []domain.CreateLocationTranslationPayload{
				{LangCode: "en-US", LocationName: "Warehouse A - Main Storage"},
				{LangCode: "id-ID", LocationName: "Gudang A - Penyimpanan Utama"},
				{LangCode: "ja-JP", LocationName: "倉庫A - メインストレージ"},
			},
		},
		{
			LocationCode: "WAREHOUSE_B",
			Building:     stringPtr("Warehouse Complex"),
			Floor:        stringPtr("Ground Floor"),
			Latitude:     float64Ptr(-6.2297),
			Longitude:    float64Ptr(106.8278),
			Translations: []domain.CreateLocationTranslationPayload{
				{LangCode: "en-US", LocationName: "Warehouse B - Electronics"},
				{LangCode: "id-ID", LocationName: "Gudang B - Elektronik"},
				{LangCode: "ja-JP", LocationName: "倉庫B - 電子機器"},
			},
		},
		{
			LocationCode: "OFFICE_FL1",
			Building:     stringPtr("Main Office"),
			Floor:        stringPtr("1st Floor"),
			Latitude:     float64Ptr(-6.2115),
			Longitude:    float64Ptr(106.8452),
			Translations: []domain.CreateLocationTranslationPayload{
				{LangCode: "en-US", LocationName: "Office Floor 1 - General"},
				{LangCode: "id-ID", LocationName: "Lantai Kantor 1 - Umum"},
				{LangCode: "ja-JP", LocationName: "オフィスフロア1 - 一般"},
			},
		},
		{
			LocationCode: "OFFICE_FL2",
			Building:     stringPtr("Main Office"),
			Floor:        stringPtr("2nd Floor"),
			Latitude:     float64Ptr(-6.2115),
			Longitude:    float64Ptr(106.8452),
			Translations: []domain.CreateLocationTranslationPayload{
				{LangCode: "en-US", LocationName: "Office Floor 2 - Management"},
				{LangCode: "id-ID", LocationName: "Lantai Kantor 2 - Manajemen"},
				{LangCode: "ja-JP", LocationName: "オフィスフロア2 - 管理"},
			},
		},
		{
			LocationCode: "PARKING_A",
			Building:     stringPtr("Parking Building"),
			Floor:        stringPtr("Level 1"),
			Latitude:     float64Ptr(-6.2090),
			Longitude:    float64Ptr(106.8460),
			Translations: []domain.CreateLocationTranslationPayload{
				{LangCode: "en-US", LocationName: "Parking Area A"},
				{LangCode: "id-ID", LocationName: "Area Parkir A"},
				{LangCode: "ja-JP", LocationName: "駐車場A"},
			},
		},
		{
			LocationCode: "CAFETERIA",
			Building:     stringPtr("Main Office"),
			Floor:        stringPtr("Ground Floor"),
			Latitude:     float64Ptr(-6.2115),
			Longitude:    float64Ptr(106.8452),
			Translations: []domain.CreateLocationTranslationPayload{
				{LangCode: "en-US", LocationName: "Employee Cafeteria"},
				{LangCode: "id-ID", LocationName: "Kafeteria Karyawan"},
				{LangCode: "ja-JP", LocationName: "社員食堂"},
			},
		},
	}
}

// generateRandomLocation creates a random location
func (ls *LocationSeeder) generateRandomLocation(index int) *domain.CreateLocationPayload {
	// Generate realistic office/building locations with shorter codes
	locationTypes := []string{"OFF", "MTG", "STG", "LAB", "WSH", "RCP", "BRK"}
	locationNames := []string{"Office", "Meeting", "Storage", "Laboratory", "Workshop", "Reception", "Break Room"}
	buildings := []string{"A", "B", "T1", "T2", "MO", "AN"}
	buildingNames := []string{"Building A", "Building B", "Tower 1", "Tower 2", "Main Office", "Annex Building"}
	floors := []string{"GF", "1F", "2F", "3F", "4F", "5F", "B1"}
	floorNames := []string{"Ground Floor", "1st Floor", "2nd Floor", "3rd Floor", "4th Floor", "5th Floor", "Basement"}

	typeIndex := rand.Intn(len(locationTypes))
	buildingIndex := rand.Intn(len(buildings))
	floorIndex := rand.Intn(len(floors))

	locationType := locationTypes[typeIndex]
	locationName := locationNames[typeIndex]
	building := buildings[buildingIndex]
	buildingName := buildingNames[buildingIndex]
	floor := floors[floorIndex]
	floorName := floorNames[floorIndex]

	// Generate short location code (max 20 chars)
	// Format: TYPE_BUILDING_FLOOR_INDEX (e.g., OFF_A_1F_001)
	locationCode := fmt.Sprintf("%s_%s_%s_%03d", locationType, building, floor, index%1000)

	// Generate coordinates around Jakarta area
	baseLat := -6.2088
	baseLng := 106.8456
	lat := baseLat + (rand.Float64()-0.5)*0.1 // ±0.05 degrees
	lng := baseLng + (rand.Float64()-0.5)*0.1 // ±0.05 degrees

	// Generate location names
	locationNameEN := fmt.Sprintf("%s - %s", buildingName, locationName)
	locationNameID := translateLocationToID(buildingName, locationName)
	locationNameJP := translateLocationToJP(buildingName, locationName)

	return &domain.CreateLocationPayload{
		LocationCode: locationCode,
		Building:     &buildingName,
		Floor:        &floorName,
		Latitude:     &lat,
		Longitude:    &lng,
		Translations: []domain.CreateLocationTranslationPayload{
			{
				LangCode:     "en-US",
				LocationName: locationNameEN,
			},
			{
				LangCode:     "id-ID",
				LocationName: locationNameID,
			},
			{
				LangCode:     "ja-JP",
				LocationName: locationNameJP,
			},
		},
	}
}

// translateLocationToID provides simple Indonesian translations
func translateLocationToID(building, locationType string) string {
	// Simple translation mapping
	buildingTranslations := map[string]string{
		"Building A":     "Gedung A",
		"Building B":     "Gedung B",
		"Tower 1":        "Menara 1",
		"Tower 2":        "Menara 2",
		"Main Office":    "Kantor Utama",
		"Annex Building": "Gedung Tambahan",
	}

	typeTranslations := map[string]string{
		"Office":     "Kantor",
		"Meeting":    "Ruang Rapat",
		"Storage":    "Penyimpanan",
		"Laboratory": "Laboratorium",
		"Workshop":   "Workshop",
		"Reception":  "Resepsi",
		"Break Room": "Ruang Istirahat",
	}

	buildingID := buildingTranslations[building]
	if buildingID == "" {
		buildingID = building
	}

	typeID := typeTranslations[locationType]
	if typeID == "" {
		typeID = locationType
	}

	return fmt.Sprintf("%s - %s", buildingID, typeID)
}

// translateLocationToJP provides simple Japanese translations
func translateLocationToJP(building, locationType string) string {
	// Simple translation mapping
	buildingTranslations := map[string]string{
		"Building A":     "ビルディングA",
		"Building B":     "ビルディングB",
		"Tower 1":        "タワー1",
		"Tower 2":        "タワー2",
		"Main Office":    "本社",
		"Annex Building": "別館",
	}

	typeTranslations := map[string]string{
		"Office":     "オフィス",
		"Meeting":    "会議室",
		"Storage":    "倉庫",
		"Laboratory": "研究室",
		"Workshop":   "作業場",
		"Reception":  "受付",
		"Break Room": "休憩室",
	}

	buildingJP := buildingTranslations[building]
	if buildingJP == "" {
		buildingJP = building
	}

	typeJP := typeTranslations[locationType]
	if typeJP == "" {
		typeJP = locationType
	}

	return fmt.Sprintf("%s - %s", buildingJP, typeJP)
}

// float64Ptr returns a pointer to float64
func float64Ptr(f float64) *float64 {
	return &f
}
