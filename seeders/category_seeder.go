package seeders

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/services/category"
	"github.com/brianvoe/gofakeit/v6"
)

// CategorySeeder handles category data seeding
type CategorySeeder struct {
	categoryService category.CategoryService
}

// NewCategorySeeder creates a new category seeder
func NewCategorySeeder(categoryService category.CategoryService) *CategorySeeder {
	return &CategorySeeder{
		categoryService: categoryService,
	}
}

// Seed creates fake categories with parent-child hierarchy
func (cs *CategorySeeder) Seed(ctx context.Context, parentCount, childrenCount int) error {
	// Seed random generator
	rand.Seed(time.Now().UnixNano())

	// First, create parent categories
	parentIDs, err := cs.createParentCategories(ctx, parentCount)
	if err != nil {
		return fmt.Errorf("failed to create parent categories: %v", err)
	}

	// Then, create children categories
	if childrenCount > 0 && len(parentIDs) > 0 {
		if err := cs.createChildCategories(ctx, parentIDs, childrenCount); err != nil {
			return fmt.Errorf("failed to create child categories: %v", err)
		}
	}

	totalCreated := len(parentIDs) + childrenCount
	fmt.Printf("✅ Successfully created %d categories (%d parents + %d children)\n",
		totalCreated, len(parentIDs), childrenCount)
	return nil
}

// createParentCategories creates parent (root) categories
func (cs *CategorySeeder) createParentCategories(ctx context.Context, count int) ([]string, error) {
	// Predefined parent categories for more realistic data
	parentCategories := []struct {
		code         string
		names        map[string]string
		descriptions map[string]string
	}{
		{
			code: "ELECTRONICS",
			names: map[string]string{
				"en-US": "Electronics",
				"id-ID": "Elektronik",
			},
			descriptions: map[string]string{
				"en-US": "Electronic devices and components",
				"id-ID": "Perangkat dan komponen elektronik",
			},
		},
		{
			code: "FURNITURE",
			names: map[string]string{
				"en-US": "Furniture",
				"id-ID": "Perabotan",
			},
			descriptions: map[string]string{
				"en-US": "Office and workplace furniture",
				"id-ID": "Perabotan kantor dan tempat kerja",
			},
		},
		{
			code: "VEHICLES",
			names: map[string]string{
				"en-US": "Vehicles",
				"id-ID": "Kendaraan",
			},
			descriptions: map[string]string{
				"en-US": "Transportation vehicles and equipment",
				"id-ID": "Kendaraan transportasi dan peralatan",
			},
		},
		{
			code: "OFFICE_SUPPLIES",
			names: map[string]string{
				"en-US": "Office Supplies",
				"id-ID": "Perlengkapan Kantor",
			},
			descriptions: map[string]string{
				"en-US": "Office supplies and stationery",
				"id-ID": "Perlengkapan kantor dan alat tulis",
			},
		},
		{
			code: "TOOLS",
			names: map[string]string{
				"en-US": "Tools & Equipment",
				"id-ID": "Alat & Peralatan",
			},
			descriptions: map[string]string{
				"en-US": "Tools and maintenance equipment",
				"id-ID": "Alat dan peralatan pemeliharaan",
			},
		},
		{
			code: "SAFETY",
			names: map[string]string{
				"en-US": "Safety Equipment",
				"id-ID": "Peralatan Keselamatan",
			},
			descriptions: map[string]string{
				"en-US": "Safety and security equipment",
				"id-ID": "Peralatan keselamatan dan keamanan",
			},
		},
	}

	var createdIDs []string

	// Create predefined categories first
	predefinedCount := len(parentCategories)
	if count > predefinedCount {
		predefinedCount = count
	}

	for i := 0; i < predefinedCount && i < len(parentCategories); i++ {
		cat := parentCategories[i]

		translations := []domain.CreateCategoryTranslationPayload{}
		for langCode, name := range cat.names {
			desc := cat.descriptions[langCode]
			translations = append(translations, domain.CreateCategoryTranslationPayload{
				LangCode:     langCode,
				CategoryName: name,
				Description:  &desc,
			})
		}

		payload := &domain.CreateCategoryPayload{
			CategoryCode: cat.code,
			Translations: translations,
		}

		created, err := cs.categoryService.CreateCategory(ctx, payload)
		if err != nil {
			fmt.Printf("   ⚠️  Failed to create parent category %s: %v\n", cat.code, err)
			continue
		}

		createdIDs = append(createdIDs, created.ID)
		fmt.Printf("   ✅ Created parent category: %s (%s)\n", cat.code, cat.names["en-US"])
	}

	// Create additional random parent categories if needed
	for i := len(parentCategories); i < count; i++ {
		categoryName := gofakeit.ProductCategory()
		categoryCode := strings.ToUpper(strings.ReplaceAll(categoryName, " ", "_"))

		// Ensure unique code
		categoryCode = fmt.Sprintf("%s_%d", categoryCode, i)

		translations := []domain.CreateCategoryTranslationPayload{
			{
				LangCode:     "en-US",
				CategoryName: categoryName,
				Description:  stringPtr(fmt.Sprintf("%s category for inventory management", categoryName)),
			},
			{
				LangCode:     "id-ID",
				CategoryName: fmt.Sprintf("Kategori %s", categoryName),
				Description:  stringPtr(fmt.Sprintf("Kategori %s untuk manajemen inventori", categoryName)),
			},
		}

		payload := &domain.CreateCategoryPayload{
			CategoryCode: categoryCode,
			Translations: translations,
		}

		created, err := cs.categoryService.CreateCategory(ctx, payload)
		if err != nil {
			fmt.Printf("   ⚠️  Failed to create parent category %s: %v\n", categoryCode, err)
			continue
		}

		createdIDs = append(createdIDs, created.ID)
		fmt.Printf("   ✅ Created parent category: %s (%s)\n", categoryCode, categoryName)
	}

	return createdIDs, nil
}

// createChildCategories creates child categories under existing parents
func (cs *CategorySeeder) createChildCategories(ctx context.Context, parentIDs []string, totalChildren int) error {
	if len(parentIDs) == 0 {
		return fmt.Errorf("no parent categories available")
	}

	// Predefined child categories for common parents
	childTemplates := map[string][]struct {
		code         string
		names        map[string]string
		descriptions map[string]string
	}{
		"ELECTRONICS": {
			{code: "COMPUTERS", names: map[string]string{"en-US": "Computers", "id-ID": "Komputer"}, descriptions: map[string]string{"en-US": "Desktop and laptop computers", "id-ID": "Komputer desktop dan laptop"}},
			{code: "PHONES", names: map[string]string{"en-US": "Phones", "id-ID": "Telepon"}, descriptions: map[string]string{"en-US": "Mobile phones and landlines", "id-ID": "Telepon seluler dan telepon rumah"}},
			{code: "PRINTERS", names: map[string]string{"en-US": "Printers", "id-ID": "Printer"}, descriptions: map[string]string{"en-US": "Printing devices and scanners", "id-ID": "Perangkat cetak dan scanner"}},
		},
		"FURNITURE": {
			{code: "CHAIRS", names: map[string]string{"en-US": "Chairs", "id-ID": "Kursi"}, descriptions: map[string]string{"en-US": "Office chairs and seating", "id-ID": "Kursi kantor dan tempat duduk"}},
			{code: "DESKS", names: map[string]string{"en-US": "Desks", "id-ID": "Meja"}, descriptions: map[string]string{"en-US": "Office desks and tables", "id-ID": "Meja kantor dan meja kerja"}},
			{code: "STORAGE", names: map[string]string{"en-US": "Storage", "id-ID": "Penyimpanan"}, descriptions: map[string]string{"en-US": "Cabinets and storage furniture", "id-ID": "Lemari dan perabotan penyimpanan"}},
		},
		"VEHICLES": {
			{code: "CARS", names: map[string]string{"en-US": "Cars", "id-ID": "Mobil"}, descriptions: map[string]string{"en-US": "Company cars and vehicles", "id-ID": "Mobil dan kendaraan perusahaan"}},
			{code: "MOTORCYCLES", names: map[string]string{"en-US": "Motorcycles", "id-ID": "Motor"}, descriptions: map[string]string{"en-US": "Motorcycles and scooters", "id-ID": "Sepeda motor dan skuter"}},
		},
	}

	childrenPerParent := totalChildren / len(parentIDs)
	if childrenPerParent == 0 {
		childrenPerParent = 1
	}

	createdCount := 0
	for i, parentID := range parentIDs {
		// Get parent category to determine predefined children
		parentCategory, err := cs.categoryService.GetCategoryById(ctx, parentID, "en-US")
		if err != nil {
			fmt.Printf("   ⚠️  Failed to get parent category %s: %v\n", parentID, err)
			continue
		}

		// Determine how many children to create for this parent
		remainingChildren := totalChildren - createdCount
		currentChildrenCount := childrenPerParent
		if i == len(parentIDs)-1 { // Last parent gets all remaining children
			currentChildrenCount = remainingChildren
		}
		if currentChildrenCount > remainingChildren {
			currentChildrenCount = remainingChildren
		}

		// Create children for this parent
		childrenCreated := cs.createChildrenForParent(ctx, parentID, parentCategory.CategoryCode, currentChildrenCount, childTemplates)
		createdCount += childrenCreated

		if createdCount >= totalChildren {
			break
		}
	}

	return nil
}

// createChildrenForParent creates children for a specific parent
func (cs *CategorySeeder) createChildrenForParent(ctx context.Context, parentID, parentCode string, count int, childTemplates map[string][]struct {
	code         string
	names        map[string]string
	descriptions map[string]string
}) int {

	var createdCount int

	// Try to use predefined children first
	if templates, exists := childTemplates[parentCode]; exists {
		for i := 0; i < count && i < len(templates); i++ {
			template := templates[i]

			translations := []domain.CreateCategoryTranslationPayload{}
			for langCode, name := range template.names {
				desc := template.descriptions[langCode]
				translations = append(translations, domain.CreateCategoryTranslationPayload{
					LangCode:     langCode,
					CategoryName: name,
					Description:  &desc,
				})
			}

			payload := &domain.CreateCategoryPayload{
				ParentID:     &parentID,
				CategoryCode: fmt.Sprintf("%s_%s", parentCode, template.code),
				Translations: translations,
			}

			_, err := cs.categoryService.CreateCategory(ctx, payload)
			if err != nil {
				fmt.Printf("   ⚠️  Failed to create child category %s: %v\n", template.code, err)
				continue
			}

			createdCount++
			fmt.Printf("   ✅ Created child category: %s under %s\n", template.names["en-US"], parentCode)
		}

		// Create additional random children if needed
		for i := len(templates); i < count; i++ {
			if createdCount >= count {
				break
			}
			createdCount += cs.createRandomChild(ctx, parentID, parentCode, i)
		}
	} else {
		// Create all random children
		for i := 0; i < count; i++ {
			createdCount += cs.createRandomChild(ctx, parentID, parentCode, i)
		}
	}

	return createdCount
}

// createRandomChild creates a random child category
func (cs *CategorySeeder) createRandomChild(ctx context.Context, parentID, parentCode string, index int) int {
	productName := gofakeit.ProductName()
	childCode := strings.ToUpper(strings.ReplaceAll(productName, " ", "_"))
	childCode = fmt.Sprintf("%s_%s_%d", parentCode, childCode, index)

	translations := []domain.CreateCategoryTranslationPayload{
		{
			LangCode:     "en-US",
			CategoryName: productName,
			Description:  stringPtr(fmt.Sprintf("%s subcategory", productName)),
		},
		{
			LangCode:     "id-ID",
			CategoryName: fmt.Sprintf("Subkategori %s", productName),
			Description:  stringPtr(fmt.Sprintf("Subkategori %s", productName)),
		},
	}

	payload := &domain.CreateCategoryPayload{
		ParentID:     &parentID,
		CategoryCode: childCode,
		Translations: translations,
	}

	_, err := cs.categoryService.CreateCategory(ctx, payload)
	if err != nil {
		fmt.Printf("   ⚠️  Failed to create child category %s: %v\n", childCode, err)
		return 0
	}

	fmt.Printf("   ✅ Created child category: %s under %s\n", productName, parentCode)
	return 1
}
