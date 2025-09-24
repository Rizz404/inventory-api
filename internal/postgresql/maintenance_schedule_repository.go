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
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type MaintenanceScheduleRepository struct {
	db *gorm.DB
}

func NewMaintenanceScheduleRepository(db *gorm.DB) *MaintenanceScheduleRepository {
	return &MaintenanceScheduleRepository{db: db}
}

// ===== Filters and Sorts =====

type MaintenanceScheduleFilterOptions struct {
	AssetID         *string                         `json:"assetId,omitempty"`
	MaintenanceType *domain.MaintenanceScheduleType `json:"maintenanceType,omitempty"`
	Status          *domain.ScheduleStatus          `json:"status,omitempty"`
	CreatedBy       *string                         `json:"createdBy,omitempty"`
	FromDate        *string                         `json:"fromDate,omitempty"` // YYYY-MM-DD
	ToDate          *string                         `json:"toDate,omitempty"`   // YYYY-MM-DD
}

func (r *MaintenanceScheduleRepository) applyScheduleFilters(db *gorm.DB, filters any) *gorm.DB {
	f, ok := filters.(*MaintenanceScheduleFilterOptions)
	if !ok || f == nil {
		return db
	}
	if f.AssetID != nil && *f.AssetID != "" {
		db = db.Where("ms.asset_id = ?", f.AssetID)
	}
	if f.MaintenanceType != nil && *f.MaintenanceType != "" {
		db = db.Where("ms.maintenance_type = ?", f.MaintenanceType)
	}
	if f.Status != nil && *f.Status != "" {
		db = db.Where("ms.status = ?", f.Status)
	}
	if f.CreatedBy != nil && *f.CreatedBy != "" {
		db = db.Where("ms.created_by = ?", f.CreatedBy)
	}
	if f.FromDate != nil && *f.FromDate != "" {
		db = db.Where("ms.scheduled_date >= ?", *f.FromDate)
	}
	if f.ToDate != nil && *f.ToDate != "" {
		db = db.Where("ms.scheduled_date <= ?", *f.ToDate)
	}
	return db
}

func (r *MaintenanceScheduleRepository) applyScheduleSorts(db *gorm.DB, sort *query.SortOptions) *gorm.DB {
	if sort == nil || sort.Field == "" {
		return db.Order("ms.created_at DESC")
	}
	var orderClause string
	switch strings.ToLower(sort.Field) {
	case "scheduled_date", "created_at":
		orderClause = "ms." + sort.Field
	case "title":
		orderClause = "mst.title"
	default:
		return db.Order("ms.created_at DESC")
	}
	order := "DESC"
	if strings.ToLower(sort.Order) == "asc" {
		order = "ASC"
	}
	return db.Order(fmt.Sprintf("%s %s", orderClause, order))
}

// ===== MUTATIONS =====

func (r *MaintenanceScheduleRepository) CreateSchedule(ctx context.Context, payload *domain.MaintenanceSchedule) (domain.MaintenanceSchedule, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.MaintenanceSchedule{}, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	m := mapper.ToModelMaintenanceScheduleForCreate(payload)
	if err := tx.Create(&m).Error; err != nil {
		tx.Rollback()
		return domain.MaintenanceSchedule{}, domain.ErrInternal(err)
	}

	for _, t := range payload.Translations {
		mt := mapper.ToModelMaintenanceScheduleTranslationForCreate(m.ID.String(), &t)
		if err := tx.Create(&mt).Error; err != nil {
			tx.Rollback()
			return domain.MaintenanceSchedule{}, domain.ErrInternal(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return domain.MaintenanceSchedule{}, domain.ErrInternal(err)
	}
	return r.GetScheduleById(ctx, m.ID.String())
}

func (r *MaintenanceScheduleRepository) UpdateSchedule(ctx context.Context, scheduleId string, payload *domain.MaintenanceSchedule) (domain.MaintenanceSchedule, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.MaintenanceSchedule{}, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// Build updates
	updates := map[string]any{}
	if payload.MaintenanceType != "" {
		updates["maintenance_type"] = payload.MaintenanceType
	}
	if !payload.ScheduledDate.IsZero() {
		updates["scheduled_date"] = payload.ScheduledDate
	}
	if payload.FrequencyMonths != nil {
		updates["frequency_months"] = payload.FrequencyMonths
	}
	if payload.Status != "" {
		updates["status"] = payload.Status
	}
	if payload.AssetID != "" {
		if parsed, err := ulid.Parse(payload.AssetID); err == nil {
			updates["asset_id"] = model.SQLULID(parsed)
		}
	}
	if payload.CreatedBy != "" {
		if parsed, err := ulid.Parse(payload.CreatedBy); err == nil {
			updates["created_by"] = model.SQLULID(parsed)
		}
	}

	if len(updates) > 0 {
		if err := tx.Model(&model.MaintenanceSchedule{}).Where("id = ?", scheduleId).Updates(updates).Error; err != nil {
			tx.Rollback()
			return domain.MaintenanceSchedule{}, domain.ErrInternal(err)
		}
	}

	// Replace translations if provided
	if len(payload.Translations) > 0 {
		if err := tx.Delete(&model.MaintenanceScheduleTranslation{}, "schedule_id = ?", scheduleId).Error; err != nil {
			tx.Rollback()
			return domain.MaintenanceSchedule{}, domain.ErrInternal(err)
		}
		for _, t := range payload.Translations {
			mt := mapper.ToModelMaintenanceScheduleTranslationForCreate(scheduleId, &t)
			if err := tx.Create(&mt).Error; err != nil {
				tx.Rollback()
				return domain.MaintenanceSchedule{}, domain.ErrInternal(err)
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return domain.MaintenanceSchedule{}, domain.ErrInternal(err)
	}
	return r.GetScheduleById(ctx, scheduleId)
}

func (r *MaintenanceScheduleRepository) DeleteSchedule(ctx context.Context, scheduleId string) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.ErrInternal(tx.Error)
	}
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Delete(&model.MaintenanceScheduleTranslation{}, "schedule_id = ?", scheduleId).Error; err != nil {
		tx.Rollback()
		return domain.ErrInternal(err)
	}
	if err := tx.Delete(&model.MaintenanceSchedule{}, "id = ?", scheduleId).Error; err != nil {
		tx.Rollback()
		return domain.ErrInternal(err)
	}
	return tx.Commit().Error
}

// ===== QUERIES =====

func (r *MaintenanceScheduleRepository) GetSchedulesPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceSchedule, error) {
	var schedules []model.MaintenanceSchedule
	db := r.db.WithContext(ctx).
		Table("maintenance_schedules ms").
		Preload("Translations").
		Preload("Asset").
		Preload("CreatedByUser")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		sq := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN maintenance_schedule_translations mst ON ms.id = mst.schedule_id").
			Where("mst.title ILIKE ?", sq).
			Distinct("ms.id")
	}

	params.Pagination.Cursor = ""
	db = query.Apply(db, params, r.applyScheduleFilters, r.applyScheduleSorts)

	if err := db.Find(&schedules).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}
	return mapper.ToDomainMaintenanceSchedules(schedules), nil
}

func (r *MaintenanceScheduleRepository) GetSchedulesCursor(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceSchedule, error) {
	var schedules []model.MaintenanceSchedule
	db := r.db.WithContext(ctx).
		Table("maintenance_schedules ms").
		Preload("Translations").
		Preload("Asset").
		Preload("CreatedByUser")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		sq := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN maintenance_schedule_translations mst ON ms.id = mst.schedule_id").
			Where("mst.title ILIKE ?", sq).
			Distinct("ms.id")
	}

	params.Pagination.Offset = 0
	db = query.Apply(db, params, r.applyScheduleFilters, r.applyScheduleSorts)

	if err := db.Find(&schedules).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}
	return mapper.ToDomainMaintenanceSchedules(schedules), nil
}

func (r *MaintenanceScheduleRepository) GetScheduleById(ctx context.Context, scheduleId string) (domain.MaintenanceSchedule, error) {
	var m model.MaintenanceSchedule
	err := r.db.WithContext(ctx).
		Table("maintenance_schedules ms").
		Preload("Translations").
		First(&m, "id = ?", scheduleId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.MaintenanceSchedule{}, domain.ErrNotFound("maintenance_schedule")
		}
		return domain.MaintenanceSchedule{}, domain.ErrInternal(err)
	}
	return mapper.ToDomainMaintenanceSchedule(&m), nil
}

// ===== Checks and Counts =====

func (r *MaintenanceScheduleRepository) CheckScheduleExist(ctx context.Context, scheduleId string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceSchedule{}).Where("id = ?", scheduleId).Count(&count).Error; err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *MaintenanceScheduleRepository) CountSchedules(ctx context.Context, params query.Params) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Table("maintenance_schedules ms")
	if params.SearchQuery != nil && *params.SearchQuery != "" {
		sq := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN maintenance_schedule_translations mst ON ms.id = mst.schedule_id").
			Where("mst.title ILIKE ?", sq).
			Distinct("ms.id")
	}
	db = query.Apply(db, query.Params{Filters: params.Filters}, r.applyScheduleFilters, nil)
	if err := db.Count(&count).Error; err != nil {
		return 0, domain.ErrInternal(err)
	}
	return count, nil
}

// ===== Statistics =====

func (r *MaintenanceScheduleRepository) GetMaintenanceScheduleStatistics(ctx context.Context) (domain.MaintenanceScheduleStatistics, error) {
	var stats domain.MaintenanceScheduleStatistics

	// Total schedules
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceSchedule{}).Count(&total).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.Total.Count = int(total)

	// By type
	var preventiveCount int64
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceSchedule{}).
		Where("maintenance_type = ?", domain.ScheduleTypePreventive).Count(&preventiveCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	var correctiveCount int64
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceSchedule{}).
		Where("maintenance_type = ?", domain.ScheduleTypeCorrective).Count(&correctiveCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.ByType.Preventive = int(preventiveCount)
	stats.ByType.Corrective = int(correctiveCount)

	// By status
	var scheduledCount, completedCount, cancelledCount int64
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceSchedule{}).Where("status = ?", domain.StatusScheduled).Count(&scheduledCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceSchedule{}).Where("status = ?", domain.StatusCompleted).Count(&completedCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceSchedule{}).Where("status = ?", domain.StatusCancelled).Count(&cancelledCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.ByStatus.Scheduled = int(scheduledCount)
	stats.ByStatus.Completed = int(completedCount)
	stats.ByStatus.Cancelled = int(cancelledCount)

	// By asset
	var byAssetResults []struct {
		AssetID         string
		AssetName       string
		AssetTag        string
		ScheduleCount   int64
		NextMaintenance string
	}
	if err := r.db.WithContext(ctx).
		Table("maintenance_schedules ms").
		Select("a.id as asset_id, a.asset_name, a.asset_tag, COUNT(*) as schedule_count, " +
			"MIN(CASE WHEN ms.scheduled_date > NOW() AND ms.status = 'Scheduled' THEN ms.scheduled_date::text ELSE NULL END) as next_maintenance").
		Joins("LEFT JOIN assets a ON ms.asset_id = a.id").
		Group("a.id, a.asset_name, a.asset_tag").
		Scan(&byAssetResults).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	for _, result := range byAssetResults {
		stats.ByAsset = append(stats.ByAsset, domain.AssetMaintenanceScheduleStatistics{
			AssetID:         result.AssetID,
			AssetName:       result.AssetName,
			AssetTag:        result.AssetTag,
			ScheduleCount:   int(result.ScheduleCount),
			NextMaintenance: result.NextMaintenance,
		})
	}

	// By creator
	var byCreatorResults []struct {
		UserID    string
		UserName  string
		UserEmail string
		Count     int64
	}
	if err := r.db.WithContext(ctx).
		Table("maintenance_schedules ms").
		Select("u.id as user_id, u.full_name as user_name, u.email as user_email, COUNT(*) as count").
		Joins("LEFT JOIN users u ON ms.created_by = u.id").
		Group("u.id, u.full_name, u.email").
		Scan(&byCreatorResults).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	for _, result := range byCreatorResults {
		stats.ByCreator = append(stats.ByCreator, domain.UserMaintenanceScheduleStatistics{
			UserID:    result.UserID,
			UserName:  result.UserName,
			UserEmail: result.UserEmail,
			Count:     int(result.Count),
		})
	}

	// Upcoming schedules
	var upcomingResults []struct {
		ID              string
		AssetID         string
		AssetName       string
		AssetTag        string
		MaintenanceType domain.MaintenanceScheduleType
		ScheduledDate   string
		DaysUntilDue    int
		Title           string
		Description     *string
	}
	if err := r.db.WithContext(ctx).
		Table("maintenance_schedules ms").
		Select("ms.id, a.id as asset_id, a.asset_name, a.asset_tag, ms.maintenance_type, " +
			"ms.scheduled_date::text, EXTRACT(DAY FROM ms.scheduled_date - NOW()) as days_until_due, " +
			"mst.title, mst.description").
		Joins("LEFT JOIN assets a ON ms.asset_id = a.id").
		Joins("LEFT JOIN maintenance_schedule_translations mst ON ms.id = mst.schedule_id").
		Where("ms.scheduled_date > NOW() AND ms.status = 'Scheduled'").
		Order("ms.scheduled_date ASC").
		Limit(10).
		Scan(&upcomingResults).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	for _, result := range upcomingResults {
		stats.UpcomingSchedule = append(stats.UpcomingSchedule, domain.UpcomingMaintenanceSchedule{
			ID:              result.ID,
			AssetID:         result.AssetID,
			AssetName:       result.AssetName,
			AssetTag:        result.AssetTag,
			MaintenanceType: result.MaintenanceType,
			ScheduledDate:   result.ScheduledDate,
			DaysUntilDue:    result.DaysUntilDue,
			Title:           result.Title,
			Description:     result.Description,
		})
	}

	// Overdue schedules
	var overdueResults []struct {
		ID              string
		AssetID         string
		AssetName       string
		AssetTag        string
		MaintenanceType domain.MaintenanceScheduleType
		ScheduledDate   string
		DaysOverdue     int
		Title           string
		Description     *string
	}
	if err := r.db.WithContext(ctx).
		Table("maintenance_schedules ms").
		Select("ms.id, a.id as asset_id, a.asset_name, a.asset_tag, ms.maintenance_type, " +
			"ms.scheduled_date::text, EXTRACT(DAY FROM NOW() - ms.scheduled_date) as days_overdue, " +
			"mst.title, mst.description").
		Joins("LEFT JOIN assets a ON ms.asset_id = a.id").
		Joins("LEFT JOIN maintenance_schedule_translations mst ON ms.id = mst.schedule_id").
		Where("ms.scheduled_date < NOW() AND ms.status = 'Scheduled'").
		Order("ms.scheduled_date ASC").
		Scan(&overdueResults).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	for _, result := range overdueResults {
		stats.OverdueSchedule = append(stats.OverdueSchedule, domain.OverdueMaintenanceSchedule{
			ID:              result.ID,
			AssetID:         result.AssetID,
			AssetName:       result.AssetName,
			AssetTag:        result.AssetTag,
			MaintenanceType: result.MaintenanceType,
			ScheduledDate:   result.ScheduledDate,
			DaysOverdue:     result.DaysOverdue,
			Title:           result.Title,
			Description:     result.Description,
		})
	}

	// Frequency trends
	var frequencyResults []struct {
		FrequencyMonths int
		Count           int64
	}
	if err := r.db.WithContext(ctx).
		Table("maintenance_schedules").
		Select("frequency_months, COUNT(*) as count").
		Where("frequency_months IS NOT NULL").
		Group("frequency_months").
		Order("frequency_months ASC").
		Scan(&frequencyResults).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	for _, result := range frequencyResults {
		stats.FrequencyTrends = append(stats.FrequencyTrends, domain.MaintenanceFrequencyTrend{
			FrequencyMonths: result.FrequencyMonths,
			Count:           int(result.Count),
		})
	}

	// Summary statistics
	stats.Summary.TotalSchedules = int(total)
	if total > 0 {
		stats.Summary.ScheduledMaintenancePercentage = float64(scheduledCount) / float64(total) * 100
		stats.Summary.CompletedMaintenancePercentage = float64(completedCount) / float64(total) * 100
		stats.Summary.CancelledMaintenancePercentage = float64(cancelledCount) / float64(total) * 100
		stats.Summary.PreventiveMaintenancePercentage = float64(preventiveCount) / float64(total) * 100
		stats.Summary.CorrectiveMaintenancePercentage = float64(correctiveCount) / float64(total) * 100
	}

	// Average schedule frequency
	var avgFreq float64
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceSchedule{}).
		Select("AVG(frequency_months)").Where("frequency_months IS NOT NULL").Scan(&avgFreq).Error; err == nil {
		stats.Summary.AverageScheduleFrequency = avgFreq
	}

	stats.Summary.UpcomingMaintenanceCount = len(stats.UpcomingSchedule)
	stats.Summary.OverdueMaintenanceCount = len(stats.OverdueSchedule)

	// Assets with/without scheduled maintenance
	var assetsWithMaintenance int64
	if err := r.db.WithContext(ctx).
		Table("assets a").
		Joins("INNER JOIN maintenance_schedules ms ON a.id = ms.asset_id").
		Where("ms.status = 'Scheduled'").
		Count(&assetsWithMaintenance).Error; err == nil {
		stats.Summary.AssetsWithScheduledMaintenance = int(assetsWithMaintenance)
	}

	var totalAssets int64
	if err := r.db.WithContext(ctx).Model(&model.Asset{}).Count(&totalAssets).Error; err == nil {
		stats.Summary.AssetsWithoutScheduledMaintenance = int(totalAssets - assetsWithMaintenance)
	}

	// Average schedules per day
	var earliest, latest time.Time
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceSchedule{}).Select("MIN(created_at)").Scan(&earliest).Error; err == nil {
		stats.Summary.EarliestScheduleDate = earliest.Format("2006-01-02")
	}
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceSchedule{}).Select("MAX(created_at)").Scan(&latest).Error; err == nil {
		stats.Summary.LatestScheduleDate = latest.Format("2006-01-02")
	}
	if !earliest.IsZero() && !latest.IsZero() {
		days := latest.Sub(earliest).Hours() / 24
		if days > 0 {
			stats.Summary.AverageSchedulesPerDay = float64(total) / days
		}
	}

	// Total unique creators
	var uniqueCreators int64
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceSchedule{}).
		Select("COUNT(DISTINCT created_by)").Scan(&uniqueCreators).Error; err == nil {
		stats.Summary.TotalUniqueCreators = int(uniqueCreators)
	}

	return stats, nil
}
