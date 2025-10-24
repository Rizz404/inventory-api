package maintenance_schedule

import (
	"context"
	"log"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/notification/messages"
	"github.com/robfig/cron/v3"
)

// CronService manages scheduled tasks for maintenance schedules
type CronService struct {
	cron                *cron.Cron
	repo                Repository
	assetService        AssetService
	notificationService NotificationService
}

// NewCronService creates a new cron service instance
func NewCronService(repo Repository, assetService AssetService, notificationService NotificationService) *CronService {
	// Create cron instance with seconds field support
	c := cron.New(cron.WithSeconds())

	return &CronService{
		cron:                c,
		repo:                repo,
		assetService:        assetService,
		notificationService: notificationService,
	}
}

// Start begins all scheduled cron jobs
func (cs *CronService) Start() error {
	// Check maintenance due soon daily at 9:00 AM
	_, err := cs.cron.AddFunc("0 0 9 * * *", cs.checkMaintenanceDueSoon)
	if err != nil {
		return err
	}

	// Check overdue maintenance daily at 9:30 AM
	_, err = cs.cron.AddFunc("0 30 9 * * *", cs.checkOverdueMaintenance)
	if err != nil {
		return err
	}

	cs.cron.Start()
	log.Println("Maintenance schedule cron service started successfully")
	return nil
}

// Stop gracefully stops all cron jobs
func (cs *CronService) Stop() {
	ctx := cs.cron.Stop()
	<-ctx.Done()
	log.Println("Maintenance schedule cron service stopped")
}

// checkMaintenanceDueSoon checks for maintenance schedules due within 7 days
func (cs *CronService) checkMaintenanceDueSoon() {
	ctx := context.Background()
	log.Println("Running maintenance due soon check...")

	// Get schedules due within 7 days
	schedules, err := cs.repo.GetSchedulesDueSoon(ctx, 7)
	if err != nil {
		log.Printf("Failed to fetch maintenance schedules due soon: %v", err)
		return
	}

	// Send notification asynchronously for each schedule
	for _, schedule := range schedules {
		scheduleCopy := schedule // Avoid closure issue
		go cs.sendMaintenanceDueSoonNotification(context.Background(), &scheduleCopy)
	}

	log.Printf("Maintenance due soon check completed. Found %d schedules due within 7 days", len(schedules))
}

// checkOverdueMaintenance checks for overdue maintenance schedules
func (cs *CronService) checkOverdueMaintenance() {
	ctx := context.Background()
	log.Println("Running overdue maintenance check...")

	// Get overdue schedules
	schedules, err := cs.repo.GetOverdueSchedules(ctx)
	if err != nil {
		log.Printf("Failed to fetch overdue maintenance schedules: %v", err)
		return
	}

	// Send notification asynchronously for each schedule
	for _, schedule := range schedules {
		scheduleCopy := schedule // Avoid closure issue
		go cs.sendMaintenanceOverdueNotification(context.Background(), &scheduleCopy)
	}

	log.Printf("Overdue maintenance check completed. Found %d overdue schedules", len(schedules))
}

// sendMaintenanceDueSoonNotification sends notification for maintenance due soon
func (cs *CronService) sendMaintenanceDueSoonNotification(ctx context.Context, schedule *domain.MaintenanceSchedule) {
	if cs.notificationService == nil {
		log.Printf("Notification service not available, skipping maintenance due soon notification for schedule ID: %s", schedule.ID)
		return
	}

	// Get asset details
	asset, err := cs.assetService.GetAssetById(ctx, schedule.AssetID, "en-US") // Default lang
	if err != nil {
		log.Printf("Failed to get asset details for schedule ID: %s: %v", schedule.ID, err)
		return
	}

	if asset.AssignedToID == nil || *asset.AssignedToID == "" {
		return
	}

	scheduledDate := schedule.ScheduledDate.Format("2006-01-02")

	titleKey, messageKey, params := messages.MaintenanceDueSoonNotification(asset.AssetName, asset.AssetTag, scheduledDate)
	utilTranslations := messages.GetMaintenanceScheduleNotificationTranslations(titleKey, messageKey, params)

	// Convert to domain translations
	translations := make([]domain.CreateNotificationTranslationPayload, len(utilTranslations))
	for i, t := range utilTranslations {
		translations[i] = domain.CreateNotificationTranslationPayload{
			LangCode: t.LangCode,
			Title:    t.Title,
			Message:  t.Message,
		}
	}

	notificationPayload := &domain.CreateNotificationPayload{
		UserID:            *asset.AssignedToID,
		RelatedEntityType: stringPtr("maintenance_schedule"),
		RelatedEntityID:   &schedule.ID,
		RelatedAssetID:    &schedule.AssetID,
		Type:              domain.NotificationTypeMaintenance,
		Priority:          domain.NotificationPriorityNormal, // Due soon = normal priority
		Translations:      translations,
	}

	_, err = cs.notificationService.CreateNotification(ctx, notificationPayload)
	if err != nil {
		log.Printf("Failed to create maintenance due soon notification for schedule ID: %s: %v", schedule.ID, err)
	} else {
		log.Printf("Successfully created maintenance due soon notification for schedule ID: %s, user ID: %s", schedule.ID, *asset.AssignedToID)
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

// sendMaintenanceOverdueNotification sends notification for overdue maintenance
func (cs *CronService) sendMaintenanceOverdueNotification(ctx context.Context, schedule *domain.MaintenanceSchedule) {
	if cs.notificationService == nil {
		log.Printf("Notification service not available, skipping overdue maintenance notification for schedule ID: %s", schedule.ID)
		return
	}

	// Get asset details
	asset, err := cs.assetService.GetAssetById(ctx, schedule.AssetID, "en-US") // Default lang
	if err != nil {
		log.Printf("Failed to get asset details for schedule ID: %s: %v", schedule.ID, err)
		return
	}

	if asset.AssignedToID == nil || *asset.AssignedToID == "" {
		return
	}

	scheduledDate := schedule.ScheduledDate.Format("2006-01-02")

	titleKey, messageKey, params := messages.MaintenanceOverdueNotification(asset.AssetName, asset.AssetTag, scheduledDate)
	utilTranslations := messages.GetMaintenanceScheduleNotificationTranslations(titleKey, messageKey, params)

	// Convert to domain translations
	translations := make([]domain.CreateNotificationTranslationPayload, len(utilTranslations))
	for i, t := range utilTranslations {
		translations[i] = domain.CreateNotificationTranslationPayload{
			LangCode: t.LangCode,
			Title:    t.Title,
			Message:  t.Message,
		}
	}

	notificationPayload := &domain.CreateNotificationPayload{
		UserID:            *asset.AssignedToID,
		RelatedEntityType: stringPtr("maintenance_schedule"),
		RelatedEntityID:   &schedule.ID,
		RelatedAssetID:    &schedule.AssetID,
		Type:              domain.NotificationTypeMaintenance,
		Priority:          domain.NotificationPriorityHigh, // Overdue = high priority
		Translations:      translations,
	}

	_, err = cs.notificationService.CreateNotification(ctx, notificationPayload)
	if err != nil {
		log.Printf("Failed to create overdue maintenance notification for schedule ID: %s: %v", schedule.ID, err)
	} else {
		log.Printf("Successfully created overdue maintenance notification for schedule ID: %s, user ID: %s", schedule.ID, *asset.AssignedToID)
	}
}
