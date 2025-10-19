package seeders

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/services/maintenance_record"
	"github.com/brianvoe/gofakeit/v6"
)

// MaintenanceRecordSeeder handles maintenance record data seeding
type MaintenanceRecordSeeder struct {
	maintenanceRecordService maintenance_record.MaintenanceRecordService
}

// NewMaintenanceRecordSeeder creates a new maintenance record seeder
func NewMaintenanceRecordSeeder(maintenanceRecordService maintenance_record.MaintenanceRecordService) *MaintenanceRecordSeeder {
	return &MaintenanceRecordSeeder{
		maintenanceRecordService: maintenanceRecordService,
	}
}

// Seed creates fake maintenance records
func (mrs *MaintenanceRecordSeeder) Seed(ctx context.Context, count int, assetIDs []string, scheduleIDs []string, userIDs []string) error {
	if len(assetIDs) == 0 {
		return fmt.Errorf("no asset IDs provided for maintenance record seeding")
	}
	if len(userIDs) == 0 {
		return fmt.Errorf("no user IDs provided for maintenance record seeding")
	}

	// Seed random generator
	rand.Seed(time.Now().UnixNano())

	successCount := 0
	for i := 0; i < count; i++ {
		recordPayload := mrs.generateMaintenanceRecordPayload(assetIDs, scheduleIDs, userIDs)

		_, err := mrs.maintenanceRecordService.CreateMaintenanceRecord(ctx, recordPayload, "en-US")
		if err != nil {
			fmt.Printf("   ⚠️ Failed to create maintenance record %d: %v\n", i+1, err)
			continue
		}

		successCount++
		if (i+1)%10 == 0 || i == count-1 {
			fmt.Printf("   📋 Created %d/%d maintenance records\n", successCount, count)
		}
	}

	fmt.Printf("✅ Successfully created %d maintenance records\n", successCount)
	return nil
}

// generateMaintenanceRecordPayload generates fake maintenance record data
func (mrs *MaintenanceRecordSeeder) generateMaintenanceRecordPayload(assetIDs []string, scheduleIDs []string, userIDs []string) *domain.CreateMaintenanceRecordPayload {
	// Select random asset
	assetID := assetIDs[rand.Intn(len(assetIDs))]

	// 70% chance to be linked to a schedule
	var scheduleID *string
	if len(scheduleIDs) > 0 && rand.Intn(10) < 7 {
		schedule := scheduleIDs[rand.Intn(len(scheduleIDs))]
		scheduleID = &schedule
	}

	// Random maintenance date in the past 2 years
	maintenanceDate := gofakeit.DateRange(time.Now().AddDate(-2, 0, 0), time.Now())

	// Determine if performed by user or vendor (60% user, 40% vendor)
	var performedByUser *string
	var performedByVendor *string

	if rand.Intn(10) < 6 {
		// Performed by internal user
		user := userIDs[rand.Intn(len(userIDs))]
		performedByUser = &user
	} else {
		// Performed by external vendor
		vendors := []string{
			"TechCorp Maintenance Services",
			"Professional IT Solutions",
			"Industrial Equipment Services",
			"ABC Technical Support",
			"Prime Maintenance Group",
			"Expert Systems Repair",
			"Global Service Partners",
			"Reliable Tech Services",
			"Advanced Equipment Care",
			"Precision Maintenance Co.",
		}
		vendor := vendors[rand.Intn(len(vendors))]
		performedByVendor = &vendor
	}

	// Generate cost based on maintenance type
	var actualCost *float64
	if performedByVendor != nil {
		// Vendor services typically cost more
		cost := float64(rand.Intn(2000) + 100) // $100-$2100
		actualCost = &cost
	} else {
		// Internal maintenance might have material costs
		if rand.Intn(2) == 0 { // 50% chance of having costs for materials
			cost := float64(rand.Intn(500) + 20) // $20-$520
			actualCost = &cost
		}
	}

	// Generate maintenance titles and detailed notes
	maintenanceTasks := []struct {
		title string
		notes string
	}{
		{
			"Routine System Cleaning",
			"Performed comprehensive cleaning of all components. Removed dust accumulation from fans, heat sinks, and air vents. Applied compressed air to clear debris. All components cleaned and tested for proper airflow.",
		},
		{
			"Software Update Installation",
			"Installed latest operating system updates and security patches. Updated all device drivers to current versions. Verified system stability after updates. Rebooted system and confirmed all applications working properly.",
		},
		{
			"Hardware Component Replacement",
			"Replaced faulty hardware component with new genuine part. Tested functionality before and after replacement. Updated asset documentation with new component information. Verified warranty coverage for replaced part.",
		},
		{
			"Performance Optimization",
			"Ran comprehensive performance diagnostics. Optimized system settings for better performance. Cleaned temporary files and defragmented storage. Measured performance improvements and documented results.",
		},
		{
			"Calibration and Testing",
			"Performed precision calibration of all measuring instruments. Tested accuracy against known standards. Adjusted settings to maintain specification compliance. Generated calibration certificate for records.",
		},
		{
			"Network Connectivity Repair",
			"Diagnosed and resolved network connectivity issues. Replaced damaged network cables and connectors. Reconfigured network settings and tested connection stability. Verified access to all required network resources.",
		},
		{
			"Power System Maintenance",
			"Inspected all power connections and electrical components. Tested power supply output voltages and stability. Replaced worn power cables and connectors. Verified proper grounding and surge protection.",
		},
		{
			"Emergency Repair Service",
			"Responded to critical system failure. Diagnosed root cause of malfunction. Implemented emergency repair procedures to restore functionality. Documented incident and recommended preventive measures.",
		},
		{
			"Preventive Maintenance Check",
			"Conducted scheduled preventive maintenance inspection. Checked all moving parts for proper lubrication. Inspected safety systems and emergency procedures. Updated maintenance log with current status.",
		},
		{
			"Firmware Update and Configuration",
			"Updated device firmware to latest stable version. Backed up existing configuration before update. Applied security patches and feature enhancements. Tested all functions after firmware update.",
		},
	}

	task := maintenanceTasks[rand.Intn(len(maintenanceTasks))]
	title := task.title
	notes := task.notes

	translations := []domain.CreateMaintenanceRecordTranslationPayload{
		{
			LangCode: "en-US",
			Title:    title,
			Notes:    &notes,
		},
		// {
		// 	LangCode: "id-ID",
		// 	Title:    title,
		// 	Notes:    &notes,
		// },
	}

	// Format maintenance date
	maintenanceDateStr := maintenanceDate.Format("2006-01-02")

	return &domain.CreateMaintenanceRecordPayload{
		ScheduleID:        scheduleID,
		AssetID:           assetID,
		MaintenanceDate:   maintenanceDateStr,
		PerformedByUser:   performedByUser,
		PerformedByVendor: performedByVendor,
		ActualCost:        actualCost,
		Translations:      translations,
	}
}
