package seeders

import (
	"context"
	"fmt"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/services/category"
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

// Seed creates categories with parent-child hierarchy from predefined data
func (cs *CategorySeeder) Seed(ctx context.Context, parentCount, childrenCount int) error {
	// Use predefined category data
	categoryData := []struct {
		Code         string
		Translations map[string]map[string]string
		Children     []struct {
			Code          string
			ItemsIncluded []string
			Translations  map[string]map[string]string
		}
	}{
		{
			Code: "FURN",
			Translations: map[string]map[string]string{
				"en-US": {"name": "Office Furniture", "desc": "General office furniture and fixtures"},
				// "id-ID": {"name": "Perabotan Kantor", "desc": "Perabotan kantor umum dan perlengkapan"},
				"ja-JP": {"name": "オフィス家具", "desc": "一般的なオフィス家具と備品"},
			},
			Children: []struct {
				Code          string
				ItemsIncluded []string
				Translations  map[string]map[string]string
			}{
				{
					Code:          "DESK",
					ItemsIncluded: []string{"Meja"},
					Translations: map[string]map[string]string{
						"en-US": {"name": "Desks & Tables", "desc": "Workstations, meeting tables, and desks"},
						// "id-ID": {"name": "Meja & Meja Kerja", "desc": "Stasiun kerja, meja rapat, dan meja"},
						"ja-JP": {"name": "机・テーブル", "desc": "ワークステーション、会議用テーブル、机"},
					},
				},
				{
					Code:          "SEAT",
					ItemsIncluded: []string{"Kursi"},
					Translations: map[string]map[string]string{
						"en-US": {"name": "Chairs & Seating", "desc": "Office chairs, sofas, and stools"},
						// "id-ID": {"name": "Kursi & Tempat Duduk", "desc": "Kursi kantor, sofa, dan bangku"},
						"ja-JP": {"name": "椅子・座席", "desc": "オフィスチェア、ソファ、スツール"},
					},
				},
				{
					Code:          "STRG",
					ItemsIncluded: []string{"Rak Nasi"},
					Translations: map[string]map[string]string{
						"en-US": {"name": "Shelves & Storage", "desc": "Cabinets, racks, and shelving units"},
						// "id-ID": {"name": "Rak & Penyimpanan", "desc": "Lemari, rak, dan unit rak"},
						"ja-JP": {"name": "棚・収納", "desc": "キャビネット、ラック、棚ユニット"},
					},
				},
			},
		},
		{
			Code: "IT",
			Translations: map[string]map[string]string{
				"en-US": {"name": "IT Equipment", "desc": "Information technology hardware and devices"},
				// "id-ID": {"name": "Peralatan IT", "desc": "Perangkat keras dan perangkat teknologi informasi"},
				"ja-JP": {"name": "IT機器", "desc": "情報技術ハードウェアおよびデバイス"},
			},
			Children: []struct {
				Code          string
				ItemsIncluded []string
				Translations  map[string]map[string]string
			}{
				{
					Code:          "COMP",
					ItemsIncluded: []string{"PC"},
					Translations: map[string]map[string]string{
						"en-US": {"name": "Computers", "desc": "Desktops, laptops, and servers"},
						// "id-ID": {"name": "Komputer", "desc": "Desktop, laptop, dan server"},
						"ja-JP": {"name": "コンピュータ", "desc": "デスクトップ、ラップトップ、サーバー"},
					},
				},
				{
					Code:          "PERI",
					ItemsIncluded: []string{"Keyboard", "Mouse"},
					Translations: map[string]map[string]string{
						"en-US": {"name": "Peripherals", "desc": "Input devices like keyboards, mice, and webcams"},
						// "id-ID": {"name": "Periferal", "desc": "Perangkat input seperti keyboard, mouse, dan webcam"},
						"ja-JP": {"name": "周辺機器", "desc": "キーボード、マウス、ウェブカメラなどの入力デバイス"},
					},
				},
				{
					Code:          "NET",
					ItemsIncluded: []string{"Router"},
					Translations: map[string]map[string]string{
						"en-US": {"name": "Networking", "desc": "Routers, switches, and modems"},
						// "id-ID": {"name": "Jaringan", "desc": "Router, switch, dan modem"},
						"ja-JP": {"name": "ネットワーク", "desc": "ルーター、スイッチ、モデム"},
					},
				},
			},
		},
		{
			Code: "ELEC",
			Translations: map[string]map[string]string{
				"en-US": {"name": "Electronics & Appliances", "desc": "General electronics and pantry appliances"},
				// "id-ID": {"name": "Elektronik & Peralatan", "desc": "Elektronik umum dan peralatan pantry"},
				"ja-JP": {"name": "家電・電子機器", "desc": "一般電子機器およびパントリー家電"},
			},
			Children: []struct {
				Code          string
				ItemsIncluded []string
				Translations  map[string]map[string]string
			}{
				{
					Code:          "PWR",
					ItemsIncluded: []string{"Colokan"},
					Translations: map[string]map[string]string{
						"en-US": {"name": "Power Management", "desc": "Power strips, extension cords, and UPS"},
						// "id-ID": {"name": "Manajemen Daya", "desc": "Strip power, kabel ekstensi, dan UPS"},
						"ja-JP": {"name": "電源管理", "desc": "電源タップ、延長コード、UPS"},
					},
				},
				{
					Code:          "KITCH",
					ItemsIncluded: []string{"Dispenser", "Teko Air Panas"},
					Translations: map[string]map[string]string{
						"en-US": {"name": "Kitchen Appliances", "desc": "Dispensers, kettles, and coffee makers"},
						// "id-ID": {"name": "Peralatan Dapur", "desc": "Dispenser, teko, dan mesin kopi"},
						"ja-JP": {"name": "キッチン家電", "desc": "ディスペンサー、ケトル、コーヒーメーカー"},
					},
				},
			},
		},
		{
			Code: "HK",
			Translations: map[string]map[string]string{
				"en-US": {"name": "Housekeeping", "desc": "Cleaning supplies and maintenance tools"},
				// "id-ID": {"name": "Kebersihan Rumah", "desc": "Persediaan pembersihan dan alat pemeliharaan"},
				"ja-JP": {"name": "清掃用品", "desc": "清掃用品およびメンテナンスツール"},
			},
			Children: []struct {
				Code          string
				ItemsIncluded []string
				Translations  map[string]map[string]string
			}{
				{
					Code:          "CLEAN",
					ItemsIncluded: []string{"Sapu", "Kain Pel", "Kemoceng", "Ember"},
					Translations: map[string]map[string]string{
						"en-US": {"name": "Cleaning Tools", "desc": "Brooms, mops, buckets, and dusters"},
						// "id-ID": {"name": "Alat Pembersih", "desc": "Sapu, kain pel, kemoceng, dan ember"},
						"ja-JP": {"name": "掃除用具", "desc": "ほうき、モップ、バケツ、ダスター"},
					},
				},
			},
		},
	}

	// Create categories from data
	var createdParentIDs []string
	for _, catData := range categoryData {
		// Create parent category
		translations := []domain.CreateCategoryTranslationPayload{}
		for langCode, trans := range catData.Translations {
			name := trans["name"]
			desc := trans["desc"]
			translations = append(translations, domain.CreateCategoryTranslationPayload{
				LangCode:     langCode,
				CategoryName: name,
				Description:  &desc,
			})
		}

		payload := &domain.CreateCategoryPayload{
			CategoryCode: catData.Code,
			Translations: translations,
		}

		created, err := cs.categoryService.CreateCategory(ctx, payload)
		if err != nil {
			fmt.Printf("   ⚠️  Failed to create parent category %s: %v\n", catData.Code, err)
			continue
		}

		createdParentIDs = append(createdParentIDs, created.ID)
		fmt.Printf("   ✅ Created parent category: %s (%s)\n", catData.Code, catData.Translations["en-US"]["name"])

		// Create children
		for _, childData := range catData.Children {
			childTranslations := []domain.CreateCategoryTranslationPayload{}
			for langCode, trans := range childData.Translations {
				name := trans["name"]
				desc := trans["desc"]
				childTranslations = append(childTranslations, domain.CreateCategoryTranslationPayload{
					LangCode:     langCode,
					CategoryName: name,
					Description:  &desc,
				})
			}

			childPayload := &domain.CreateCategoryPayload{
				ParentID:     &created.ID,
				CategoryCode: fmt.Sprintf("%s-%s", catData.Code, childData.Code),
				Translations: childTranslations,
			}

			_, err := cs.categoryService.CreateCategory(ctx, childPayload)
			if err != nil {
				fmt.Printf("   ⚠️  Failed to create child category %s: %v\n", childData.Code, err)
				continue
			}

			fmt.Printf("   ✅ Created child category: %s under %s\n", childData.Translations["en-US"]["name"], catData.Code)
		}
	}

	totalCreated := len(createdParentIDs) + len(categoryData)*len(categoryData[0].Children) // Approximate
	fmt.Printf("✅ Successfully created %d categories\n", totalCreated)
	return nil
}
