package seeders

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/services/maintenance_schedule"
	"github.com/brianvoe/gofakeit/v6"
)

// MaintenanceScheduleSeeder handles maintenance schedule data seeding
type MaintenanceScheduleSeeder struct {
	maintenanceScheduleService maintenance_schedule.MaintenanceScheduleService
}

// NewMaintenanceScheduleSeeder creates a new maintenance schedule seeder
func NewMaintenanceScheduleSeeder(maintenanceScheduleService maintenance_schedule.MaintenanceScheduleService) *MaintenanceScheduleSeeder {
	return &MaintenanceScheduleSeeder{
		maintenanceScheduleService: maintenanceScheduleService,
	}
}

// Seed creates fake maintenance schedules
func (mss *MaintenanceScheduleSeeder) Seed(ctx context.Context, count int, assetIDs []string, userIDs []string) error {
	if len(assetIDs) == 0 {
		return fmt.Errorf("no asset IDs provided for maintenance schedule seeding")
	}
	if len(userIDs) == 0 {
		return fmt.Errorf("no user IDs provided for maintenance schedule seeding")
	}

	// Seed random generator
	rand.Seed(time.Now().UnixNano())

	successCount := 0
	for i := 0; i < count; i++ {
		// ! Select random creator user
		createdBy := userIDs[rand.Intn(len(userIDs))]

		schedulePayload := mss.generateMaintenanceSchedulePayload(assetIDs)

		_, err := mss.maintenanceScheduleService.CreateMaintenanceSchedule(ctx, schedulePayload, createdBy)
		if err != nil {
			fmt.Printf("   ⚠️ Failed to create maintenance schedule %d: %v\n", i+1, err)
			continue
		}

		successCount++
		if (i+1)%10 == 0 || i == count-1 {
			fmt.Printf("   🔧 Created %d/%d maintenance schedules\n", successCount, count)
		}
	}

	fmt.Printf("✅ Successfully created %d maintenance schedules\n", successCount)
	return nil
}

// generateMaintenanceSchedulePayload generates fake maintenance schedule data
func (mss *MaintenanceScheduleSeeder) generateMaintenanceSchedulePayload(assetIDs []string) *domain.CreateMaintenanceSchedulePayload {
	// Select random asset
	assetID := assetIDs[rand.Intn(len(assetIDs))]

	// Random maintenance type with realistic distribution
	maintenanceTypes := []domain.MaintenanceScheduleType{
		domain.ScheduleTypePreventive, domain.ScheduleTypePreventive, domain.ScheduleTypePreventive, // 75% preventive
		domain.ScheduleTypeCorrective, // 25% corrective
	}
	maintenanceType := maintenanceTypes[rand.Intn(len(maintenanceTypes))]

	// Generate next scheduled date
	var nextScheduledDate time.Time
	if maintenanceType == domain.ScheduleTypePreventive {
		// Preventive maintenance scheduled in the future (next 6 months)
		nextScheduledDate = gofakeit.DateRange(time.Now().AddDate(0, 1, 0), time.Now().AddDate(0, 6, 0))
	} else {
		// Corrective maintenance can be in past or near future
		nextScheduledDate = gofakeit.DateRange(time.Now().AddDate(0, -1, 0), time.Now().AddDate(0, 2, 0))
	}

	// Recurring settings
	isRecurring := false
	var intervalValue *int
	var intervalUnit *domain.IntervalUnit

	if maintenanceType == domain.ScheduleTypePreventive && rand.Float32() > 0.3 { // 70% of preventive are recurring
		isRecurring = true
		// Common intervals
		intervals := []struct {
			value int
			unit  domain.IntervalUnit
		}{
			{7, domain.IntervalDays},   // weekly
			{14, domain.IntervalDays},  // bi-weekly
			{1, domain.IntervalMonths}, // monthly
			{3, domain.IntervalMonths}, // quarterly
			{6, domain.IntervalMonths}, // semi-annual
			{1, domain.IntervalYears},  // annual
		}
		interval := intervals[rand.Intn(len(intervals))]
		intervalValue = &interval.value
		intervalUnit = &interval.unit
	}

	// Auto-complete for one-time maintenance
	autoComplete := false
	if !isRecurring && rand.Float32() > 0.5 {
		autoComplete = true
	}

	// Estimated cost
	var estimatedCost *float64
	if rand.Float32() > 0.2 { // 80% have estimated cost
		cost := float64(rand.Intn(500) + 50) // $50 - $550
		estimatedCost = &cost
	}

	// Generate maintenance titles and descriptions based on type
	var title, description string
	if maintenanceType == domain.ScheduleTypePreventive {
		preventiveTasks := []struct {
			title       string
			description string
		}{
			{"Regular System Cleaning", "Perform thorough cleaning of equipment components and check for dust accumulation."},
			{"Software Updates Installation", "Install latest software updates, security patches, and driver updates."},
			{"Hardware Inspection", "Inspect all hardware components for signs of wear, damage, or loose connections."},
			{"Performance Optimization", "Run performance diagnostics and optimize system settings for better performance."},
			{"Calibration Check", "Verify and calibrate equipment settings to maintain accuracy and performance."},
			{"Filter Replacement", "Replace air filters, dust filters, and other consumable filter components."},
			{"Lubrication Service", "Apply appropriate lubricants to moving parts and mechanical components."},
			{"Battery Maintenance", "Check battery health, clean terminals, and replace if necessary."},
			{"Safety Inspection", "Conduct comprehensive safety inspection and test safety features."},
			{"Backup and Recovery Test", "Test backup systems and verify data recovery procedures."},
		}
		task := preventiveTasks[rand.Intn(len(preventiveTasks))]
		title = task.title
		description = task.description
	} else {
		correctiveTasks := []struct {
			title       string
			description string
		}{
			{"Emergency Repair Service", "Address critical hardware failure requiring immediate repair or replacement."},
			{"System Recovery", "Recover system from failure state and restore normal operations."},
			{"Component Replacement", "Replace failed or damaged component with new or refurbished part."},
			{"Software Troubleshooting", "Diagnose and resolve software-related issues and conflicts."},
			{"Network Connectivity Fix", "Resolve network connectivity issues and restore communication."},
			{"Power System Repair", "Repair or replace power supply components and electrical systems."},
			{"Display System Repair", "Fix display issues, replace screens, or repair video components."},
			{"Mechanical Repair", "Repair or replace mechanical components and moving parts."},
			{"Firmware Recovery", "Restore or update firmware to resolve system instability."},
			{"Data Recovery Service", "Recover lost or corrupted data from storage devices."},
		}
		task := correctiveTasks[rand.Intn(len(correctiveTasks))]
		title = task.title
		description = task.description
	}

	translations := []domain.CreateMaintenanceScheduleTranslationPayload{
		{
			LangCode:    "en-US",
			Title:       title,
			Description: &description,
		},
		{
			LangCode:    "id-ID",
			Title:       title,
			Description: &description,
		},
		{
			LangCode:    "ja-JP",
			Title:       title,
			Description: &description,
		},
	}

	// Format next scheduled date
	nextScheduledDateStr := nextScheduledDate.Format("2006-01-02")

	return &domain.CreateMaintenanceSchedulePayload{
		AssetID:           assetID,
		MaintenanceType:   maintenanceType,
		IsRecurring:       &isRecurring,
		IntervalValue:     intervalValue,
		IntervalUnit:      intervalUnit,
		NextScheduledDate: nextScheduledDateStr,
		AutoComplete:      &autoComplete,
		EstimatedCost:     estimatedCost,
		Translations:      translations,
	}
}
