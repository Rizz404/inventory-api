package asset

import (
	"context"
	"log"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/notification/messages"
	"github.com/robfig/cron/v3"
)

// CronService manages scheduled tasks for assets
type CronService struct {
	cron                *cron.Cron
	assetRepo           Repository
	notificationService NotificationService
}

// NewCronService creates a new cron service instance
func NewCronService(assetRepo Repository, notificationService NotificationService) *CronService {
	// Create cron instance with seconds field support
	c := cron.New(cron.WithSeconds())

	return &CronService{
		cron:                c,
		assetRepo:           assetRepo,
		notificationService: notificationService,
	}
}

// Start begins all scheduled cron jobs
func (cs *CronService) Start() error {
	// Check warranty expiring daily at 9:00 AM
	_, err := cs.cron.AddFunc("0 0 9 * * *", cs.checkWarrantyExpiring)
	if err != nil {
		return err
	}

	// Check expired warranties daily at 9:30 AM
	_, err = cs.cron.AddFunc("0 30 9 * * *", cs.checkExpiredWarranties)
	if err != nil {
		return err
	}

	cs.cron.Start()
	log.Println("Asset cron service started successfully")
	return nil
}

// Stop gracefully stops all cron jobs
func (cs *CronService) Stop() {
	ctx := cs.cron.Stop()
	<-ctx.Done()
	log.Println("Asset cron service stopped")
}

// checkWarrantyExpiring checks for assets with warranties expiring within 30 days
func (cs *CronService) checkWarrantyExpiring() {
	ctx := context.Background()
	log.Println("Running warranty expiring check...")

	// Get assets with warranties expiring within 30 days
	assets, err := cs.assetRepo.GetAssetsWithWarrantyExpiring(ctx, 30)
	if err != nil {
		log.Printf("Failed to fetch assets for warranty check: %v", err)
		return
	}

	// Send notification to each asset's assigned user
	for _, asset := range assets {
		if asset.AssignedTo != nil && *asset.AssignedTo != "" {
			cs.sendWarrantyExpiringNotification(ctx, &asset)
		}
	}

	log.Printf("Warranty expiring check completed. Found %d assets with warranties expiring within 30 days", len(assets))
}

// checkExpiredWarranties checks for assets with expired warranties
func (cs *CronService) checkExpiredWarranties() {
	ctx := context.Background()
	log.Println("Running expired warranty check...")

	// Get assets with expired warranties (today)
	assets, err := cs.assetRepo.GetAssetsWithExpiredWarranty(ctx)
	if err != nil {
		log.Printf("Failed to fetch assets for expired warranty check: %v", err)
		return
	}

	// Send notification to each asset's assigned user
	for _, asset := range assets {
		if asset.AssignedTo != nil && *asset.AssignedTo != "" {
			cs.sendWarrantyExpiredNotification(ctx, &asset)
		}
	}

	log.Printf("Expired warranty check completed. Found %d assets with expired warranties", len(assets))
}

// sendWarrantyExpiringNotification sends notification for warranty expiring soon
func (cs *CronService) sendWarrantyExpiringNotification(ctx context.Context, asset *domain.Asset) {
	if cs.notificationService == nil {
		log.Printf("Notification service not available, skipping warranty expiring notification for asset ID: %s", asset.ID)
		return
	}

	if asset.WarrantyEnd == nil {
		return
	}

	expiryDate := asset.WarrantyEnd.Format("2006-01-02")
	assetIdStr := asset.ID

	titleKey, messageKey, params := messages.AssetWarrantyExpiringNotification(asset.AssetName, asset.AssetTag, expiryDate)
	utilTranslations := messages.GetAssetNotificationTranslations(titleKey, messageKey, params)

	// Convert to domain translations
	translations := make([]domain.CreateNotificationTranslationPayload, len(utilTranslations))
	for i, t := range utilTranslations {
		translations[i] = domain.CreateNotificationTranslationPayload{
			LangCode: t.LangCode,
			Title:    t.Title,
			Message:  t.Message,
		}
	}

	entityType := "asset"
	priority := domain.NotificationPriorityNormal

	notificationPayload := &domain.CreateNotificationPayload{
		UserID:            *asset.AssignedTo,
		RelatedEntityType: &entityType,
		RelatedEntityID:   &assetIdStr,
		RelatedAssetID:    &assetIdStr,
		Type:              domain.NotificationTypeWarranty,
		Priority:          priority,
		Translations:      translations,
	}

	_, err := cs.notificationService.CreateNotification(ctx, notificationPayload)
	if err != nil {
		log.Printf("Failed to create warranty expiring notification for asset ID: %s: %v", asset.ID, err)
	} else {
		log.Printf("Successfully created warranty expiring notification for asset ID: %s, user ID: %s", asset.ID, *asset.AssignedTo)
	}
}

// sendWarrantyExpiredNotification sends notification for expired warranty
func (cs *CronService) sendWarrantyExpiredNotification(ctx context.Context, asset *domain.Asset) {
	if cs.notificationService == nil {
		log.Printf("Notification service not available, skipping warranty expired notification for asset ID: %s", asset.ID)
		return
	}

	assetIdStr := asset.ID

	titleKey, messageKey, params := messages.AssetWarrantyExpiredNotification(asset.AssetName, asset.AssetTag)
	utilTranslations := messages.GetAssetNotificationTranslations(titleKey, messageKey, params)

	// Convert to domain translations
	translations := make([]domain.CreateNotificationTranslationPayload, len(utilTranslations))
	for i, t := range utilTranslations {
		translations[i] = domain.CreateNotificationTranslationPayload{
			LangCode: t.LangCode,
			Title:    t.Title,
			Message:  t.Message,
		}
	}

	entityType := "asset"
	priority := domain.NotificationPriorityHigh

	notificationPayload := &domain.CreateNotificationPayload{
		UserID:            *asset.AssignedTo,
		RelatedEntityType: &entityType,
		RelatedEntityID:   &assetIdStr,
		RelatedAssetID:    &assetIdStr,
		Type:              domain.NotificationTypeWarranty,
		Priority:          priority,
		Translations:      translations,
	}

	_, err := cs.notificationService.CreateNotification(ctx, notificationPayload)
	if err != nil {
		log.Printf("Failed to create warranty expired notification for asset ID: %s: %v", asset.ID, err)
	} else {
		log.Printf("Successfully created warranty expired notification for asset ID: %s, user ID: %s", asset.ID, *asset.AssignedTo)
	}
}
