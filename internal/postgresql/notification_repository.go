package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/query"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

type NotificationFilterOptions struct {
	UserID         *string                  `json:"userId,omitempty"`
	RelatedAssetID *string                  `json:"relatedAssetId,omitempty"`
	Type           *domain.NotificationType `json:"type,omitempty"`
	IsRead         *bool                    `json:"isRead,omitempty"`
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (r *NotificationRepository) applyNotificationFilters(db *gorm.DB, filters any) *gorm.DB {
	f, ok := filters.(*NotificationFilterOptions)
	if !ok || f == nil {
		return db
	}

	if f.UserID != nil {
		db = db.Where("n.user_id = ?", f.UserID)
	}
	if f.RelatedAssetID != nil {
		db = db.Where("n.related_asset_id = ?", f.RelatedAssetID)
	}
	if f.Type != nil {
		db = db.Where("n.type = ?", f.Type)
	}
	if f.IsRead != nil {
		db = db.Where("n.is_read = ?", f.IsRead)
	}
	return db
}

func (r *NotificationRepository) applyNotificationSorts(db *gorm.DB, sort *query.SortOptions) *gorm.DB {
	if sort == nil || sort.Field == "" {
		return db.Order("n.created_at DESC")
	}

	var orderClause string
	switch strings.ToLower(sort.Field) {
	case "type", "is_read", "created_at":
		orderClause = fmt.Sprintf("n.%s", sort.Field)
	case "title":
		orderClause = "nt.title"
	case "message":
		orderClause = "nt.message"
	default:
		return db.Order("n.created_at DESC")
	}

	order := "DESC"
	if strings.ToLower(sort.Order) == "asc" {
		order = "ASC"
	}
	return db.Order(fmt.Sprintf("%s %s", orderClause, order))
}

// *===========================MUTATION===========================*
func (r *NotificationRepository) CreateNotification(ctx context.Context, payload *domain.Notification) (domain.Notification, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.Notification{}, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create notification
	modelNotification := mapper.ToModelNotificationForCreate(payload)
	if err := tx.Create(&modelNotification).Error; err != nil {
		tx.Rollback()
		return domain.Notification{}, domain.ErrInternal(err)
	}

	// Create translations
	for _, translation := range payload.Translations {
		modelTranslation := mapper.ToModelNotificationTranslationForCreate(modelNotification.ID.String(), &translation)
		if err := tx.Create(&modelTranslation).Error; err != nil {
			tx.Rollback()
			return domain.Notification{}, domain.ErrInternal(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return domain.Notification{}, domain.ErrInternal(err)
	}

	// Fetch created notification with translations
	return r.GetNotificationById(ctx, modelNotification.ID.String())
}

func (r *NotificationRepository) UpdateNotification(ctx context.Context, notificationId string, payload *domain.UpdateNotificationPayload) (domain.Notification, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.Notification{}, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update notification basic info
	updates := mapper.ToModelNotificationUpdateMap(payload)
	if len(updates) > 0 {
		if err := tx.Model(&model.Notification{}).Where("id = ?", notificationId).Updates(updates).Error; err != nil {
			tx.Rollback()
			return domain.Notification{}, domain.ErrInternal(err)
		}
	}

	// Update translations if provided
	if len(payload.Translations) > 0 {
		for _, translationPayload := range payload.Translations {
			var existingTranslation model.NotificationTranslation
			err := tx.Where("notification_id = ? AND lang_code = ?", notificationId, translationPayload.LangCode).
				First(&existingTranslation).Error

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Create new translation
					newTranslation := mapper.ToModelNotificationTranslationForCreate(notificationId, &domain.NotificationTranslation{
						LangCode: translationPayload.LangCode,
						Title:    translationPayload.Title,
						Message:  translationPayload.Message,
					})
					if err := tx.Create(&newTranslation).Error; err != nil {
						tx.Rollback()
						return domain.Notification{}, domain.ErrInternal(err)
					}
				} else {
					tx.Rollback()
					return domain.Notification{}, domain.ErrInternal(err)
				}
			} else {
				// Update existing translation
				existingTranslation.Title = translationPayload.Title
				existingTranslation.Message = translationPayload.Message
				if err := tx.Save(&existingTranslation).Error; err != nil {
					tx.Rollback()
					return domain.Notification{}, domain.ErrInternal(err)
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return domain.Notification{}, domain.ErrInternal(err)
	}

	// Fetch updated notification with translations
	return r.GetNotificationById(ctx, notificationId)
}

func (r *NotificationRepository) DeleteNotification(ctx context.Context, notificationId string) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.ErrInternal(tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete translations first (foreign key constraint)
	if err := tx.Delete(&model.NotificationTranslation{}, "notification_id = ?", notificationId).Error; err != nil {
		tx.Rollback()
		return domain.ErrInternal(err)
	}

	// Delete notification
	if err := tx.Delete(&model.Notification{}, "id = ?", notificationId).Error; err != nil {
		tx.Rollback()
		return domain.ErrInternal(err)
	}

	if err := tx.Commit().Error; err != nil {
		return domain.ErrInternal(err)
	}

	return nil
}

// *===========================QUERY===========================*
func (r *NotificationRepository) GetNotificationsPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.Notification, error) {
	var notifications []model.Notification
	db := r.db.WithContext(ctx).
		Table("notifications n").
		Preload("Translations")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN notification_translations nt ON n.id = nt.notification_id").
			Where("nt.title ILIKE ? OR nt.message ILIKE ?", searchPattern, searchPattern).
			Distinct("n.id")
	}

	// Set pagination cursor to empty for offset-based pagination
	params.Pagination.Cursor = ""
	db = query.Apply(db, params, r.applyNotificationFilters, r.applyNotificationSorts)

	if err := db.Find(&notifications).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain notifications
	return mapper.ToDomainNotifications(notifications), nil
}

func (r *NotificationRepository) GetNotificationsCursor(ctx context.Context, params query.Params, langCode string) ([]domain.Notification, error) {
	var notifications []model.Notification
	db := r.db.WithContext(ctx).
		Table("notifications n").
		Preload("Translations")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN notification_translations nt ON n.id = nt.notification_id").
			Where("nt.title ILIKE ? OR nt.message ILIKE ?", searchPattern, searchPattern).
			Distinct("n.id")
	}

	// Set offset to 0 for cursor-based pagination
	params.Pagination.Offset = 0
	db = query.Apply(db, params, r.applyNotificationFilters, r.applyNotificationSorts)

	if err := db.Find(&notifications).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain notifications
	return mapper.ToDomainNotifications(notifications), nil
}

func (r *NotificationRepository) GetNotificationById(ctx context.Context, notificationId string) (domain.Notification, error) {
	var notification model.Notification

	err := r.db.WithContext(ctx).
		Table("notifications n").
		Preload("Translations").
		First(&notification, "id = ?", notificationId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Notification{}, domain.ErrNotFound("notification not found")
		}
		return domain.Notification{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainNotification(&notification), nil
}

func (r *NotificationRepository) CheckNotificationExist(ctx context.Context, notificationId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Notification{}).Where("id = ?", notificationId).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *NotificationRepository) CountNotifications(ctx context.Context, params query.Params) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Table("notifications n")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN notification_translations nt ON n.id = nt.notification_id").
			Where("nt.title ILIKE ? OR nt.message ILIKE ?", searchPattern, searchPattern).
			Distinct("n.id")
	}

	db = query.Apply(db, query.Params{Filters: params.Filters}, r.applyNotificationFilters, nil)

	if err := db.Count(&count).Error; err != nil {
		return 0, domain.ErrInternal(err)
	}
	return count, nil
}

func (r *NotificationRepository) GetNotificationStatistics(ctx context.Context) (domain.NotificationStatistics, error) {
	var stats domain.NotificationStatistics

	// Get total notification count
	var totalCount int64
	if err := r.db.WithContext(ctx).Model(&model.Notification{}).Count(&totalCount).Error; err != nil {
		return domain.NotificationStatistics{}, domain.ErrInternal(err)
	}
	stats.Total.Count = int(totalCount)

	// Get type statistics
	var typeStats []struct {
		Type  domain.NotificationType `json:"type"`
		Count int64                   `json:"count"`
	}
	if err := r.db.WithContext(ctx).Model(&model.Notification{}).
		Select("type, COUNT(*) as count").
		Group("type").
		Find(&typeStats).Error; err != nil {
		return domain.NotificationStatistics{}, domain.ErrInternal(err)
	}

	for _, ts := range typeStats {
		switch ts.Type {
		case domain.NotificationTypeMaintenance:
			stats.ByType.Maintenance = int(ts.Count)
		case domain.NotificationTypeWarranty:
			stats.ByType.Warranty = int(ts.Count)
		case domain.NotificationTypeStatusChange:
			stats.ByType.StatusChange = int(ts.Count)
		case domain.NotificationTypeMovement:
			stats.ByType.Movement = int(ts.Count)
		case domain.NotificationTypeIssueReport:
			stats.ByType.IssueReport = int(ts.Count)
		}
	}

	// Get status statistics
	var readCount, unreadCount int64
	if err := r.db.WithContext(ctx).Model(&model.Notification{}).Where("is_read = ?", true).Count(&readCount).Error; err != nil {
		return domain.NotificationStatistics{}, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.Notification{}).Where("is_read = ?", false).Count(&unreadCount).Error; err != nil {
		return domain.NotificationStatistics{}, domain.ErrInternal(err)
	}

	stats.ByStatus.Read = int(readCount)
	stats.ByStatus.Unread = int(unreadCount)

	// Get creation trends (last 30 days)
	var creationTrends []struct {
		Date  time.Time `json:"date"`
		Count int64     `json:"count"`
	}
	if err := r.db.WithContext(ctx).Model(&model.Notification{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", time.Now().AddDate(0, 0, -30)).
		Group("DATE(created_at)").
		Order("date DESC").
		Find(&creationTrends).Error; err != nil {
		return domain.NotificationStatistics{}, domain.ErrInternal(err)
	}

	stats.CreationTrends = make([]domain.NotificationCreationTrend, len(creationTrends))
	for i, ct := range creationTrends {
		stats.CreationTrends[i] = domain.NotificationCreationTrend{
			Date:  ct.Date,
			Count: int(ct.Count),
		}
	}

	// Calculate summary statistics
	stats.Summary.TotalNotifications = int(totalCount)

	if totalCount > 0 {
		stats.Summary.ReadPercentage = float64(readCount) / float64(totalCount) * 100
		stats.Summary.UnreadPercentage = float64(unreadCount) / float64(totalCount) * 100
	}

	// Find most common type
	mostCommonCount := 0
	mostCommonType := ""
	if stats.ByType.Maintenance > mostCommonCount {
		mostCommonCount = stats.ByType.Maintenance
		mostCommonType = "MAINTENANCE"
	}
	if stats.ByType.Warranty > mostCommonCount {
		mostCommonCount = stats.ByType.Warranty
		mostCommonType = "WARRANTY"
	}
	if stats.ByType.StatusChange > mostCommonCount {
		mostCommonCount = stats.ByType.StatusChange
		mostCommonType = "STATUS_CHANGE"
	}
	if stats.ByType.Movement > mostCommonCount {
		mostCommonCount = stats.ByType.Movement
		mostCommonType = "MOVEMENT"
	}
	if stats.ByType.IssueReport > mostCommonCount {
		mostCommonType = "ISSUE_REPORT"
	}
	stats.Summary.MostCommonType = mostCommonType

	// Get earliest and latest creation dates
	var earliestDate, latestDate time.Time
	if err := r.db.WithContext(ctx).Model(&model.Notification{}).
		Select("MIN(created_at) as earliest, MAX(created_at) as latest").
		Row().Scan(&earliestDate, &latestDate); err != nil {
		return domain.NotificationStatistics{}, domain.ErrInternal(err)
	}

	stats.Summary.EarliestCreationDate = earliestDate
	stats.Summary.LatestCreationDate = latestDate

	// Calculate average notifications per day
	if !earliestDate.IsZero() && !latestDate.IsZero() {
		daysDiff := latestDate.Sub(earliestDate).Hours() / 24
		if daysDiff > 0 {
			stats.Summary.AverageNotificationsPerDay = float64(totalCount) / daysDiff
		}
	}

	return stats, nil
}

// Mark notification as read/unread
func (r *NotificationRepository) MarkNotificationAsRead(ctx context.Context, notificationId string, isRead bool) error {
	err := r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("id = ?", notificationId).
		Update("is_read", isRead).Error
	if err != nil {
		return domain.ErrInternal(err)
	}
	return nil
}

// Mark all notifications for a user as read
func (r *NotificationRepository) MarkAllNotificationsAsRead(ctx context.Context, userId string) error {
	err := r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userId, false).
		Update("is_read", true).Error
	if err != nil {
		return domain.ErrInternal(err)
	}
	return nil
}
