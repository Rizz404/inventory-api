package seeders

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/services/issue_report"
	"github.com/brianvoe/gofakeit/v6"
)

// IssueReportSeeder handles issue report data seeding
type IssueReportSeeder struct {
	issueReportService issue_report.IssueReportService
}

// NewIssueReportSeeder creates a new issue report seeder
func NewIssueReportSeeder(issueReportService issue_report.IssueReportService) *IssueReportSeeder {
	return &IssueReportSeeder{
		issueReportService: issueReportService,
	}
}

// Seed creates fake issue reports
func (irs *IssueReportSeeder) Seed(ctx context.Context, count int, assetIDs []string, userIDs []string) error {
	if len(assetIDs) == 0 {
		return fmt.Errorf("no asset IDs provided for issue report seeding")
	}
	if len(userIDs) == 0 {
		return fmt.Errorf("no user IDs provided for issue report seeding")
	}

	// Seed random generator
	rand.Seed(time.Now().UnixNano())

	successCount := 0
	for i := 0; i < count; i++ {
		// ! Select random reporter user
		reportedBy := userIDs[rand.Intn(len(userIDs))]

		issuePayload := irs.generateIssueReportPayload(assetIDs)

		_, err := irs.issueReportService.CreateIssueReport(ctx, issuePayload, reportedBy)
		if err != nil {
			fmt.Printf("   ⚠️ Failed to create issue report %d: %v\n", i+1, err)
			continue
		}

		successCount++
		if (i+1)%10 == 0 || i == count-1 {
			fmt.Printf("   🐛 Created %d/%d issue reports\n", successCount, count)
		}
	}

	fmt.Printf("✅ Successfully created %d issue reports\n", successCount)
	return nil
}

// generateIssueReportPayload generates fake issue report data
func (irs *IssueReportSeeder) generateIssueReportPayload(assetIDs []string) *domain.CreateIssueReportPayload {
	// Select random asset
	assetID := assetIDs[rand.Intn(len(assetIDs))]

	// Random report date in the past 6 months
	_ = gofakeit.DateRange(time.Now().AddDate(0, -6, 0), time.Now())

	// Issue types based on common asset problems
	issueTypes := []string{
		"Hardware Malfunction",
		"Software Issue",
		"Physical Damage",
		"Performance Problem",
		"Connectivity Issue",
		"Power Problem",
		"Display Issue",
		"Audio Problem",
		"Overheating",
		"Battery Issue",
		"Missing Parts",
		"Wear and Tear",
	}

	issueType := issueTypes[rand.Intn(len(issueTypes))]

	// Random priority with realistic distribution
	priorities := []domain.IssuePriority{
		domain.PriorityLow, domain.PriorityLow, domain.PriorityLow, // 30%
		domain.PriorityMedium, domain.PriorityMedium, domain.PriorityMedium, domain.PriorityMedium, // 40%
		domain.PriorityHigh, domain.PriorityHigh, // 20%
		domain.PriorityCritical, // 10%
	}
	priority := priorities[rand.Intn(len(priorities))]

	// Note: Status is managed through separate update operations after creation

	// Note: Resolved data is managed separately through update operations

	// Generate issue titles and descriptions based on type
	var title, description string
	switch issueType {
	case "Hardware Malfunction":
		title = "Hardware component not functioning properly"
		description = "The hardware component is showing signs of malfunction and needs immediate attention."
	case "Software Issue":
		title = "Software application error encountered"
		description = "Software is crashing or showing error messages during normal operation."
	case "Physical Damage":
		title = "Physical damage observed on equipment"
		description = "Equipment shows visible physical damage that may affect functionality."
	case "Performance Problem":
		title = "Equipment performance degradation"
		description = "Equipment is running slower than normal or showing performance issues."
	case "Connectivity Issue":
		title = "Network or connectivity problem"
		description = "Equipment is unable to connect to network or other devices properly."
	case "Power Problem":
		title = "Power supply or battery issue"
		description = "Equipment is experiencing power-related problems or battery drain."
	case "Display Issue":
		title = "Display or screen malfunction"
		description = "Screen is flickering, showing artifacts, or not displaying properly."
	case "Audio Problem":
		title = "Audio system malfunction"
		description = "Audio output is distorted, too quiet, or not working at all."
	case "Overheating":
		title = "Equipment overheating detected"
		description = "Equipment is running hot and may shut down due to temperature protection."
	case "Battery Issue":
		title = "Battery performance degradation"
		description = "Battery is not holding charge or draining faster than expected."
	case "Missing Parts":
		title = "Missing components or accessories"
		description = "Some parts or accessories are missing from the equipment."
	default:
		title = "General equipment issue"
		description = "Equipment is experiencing issues that need investigation."
	}

	// Note: Resolution notes are added when issues are resolved

	translations := []domain.CreateIssueReportTranslationPayload{
		{
			LangCode:    "en-US",
			Title:       title,
			Description: &description,
		},
		// {
		// 	LangCode:    "id-ID",
		// 	Title:       title,
		// 	Description: &description,
		// },
	}

	return &domain.CreateIssueReportPayload{
		AssetID:      assetID,
		IssueType:    issueType,
		Priority:     priority,
		Translations: translations,
	}
}
