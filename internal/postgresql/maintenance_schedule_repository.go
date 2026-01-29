package postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MaintenanceScheduleRepository struct {
	db *gorm.DB
}

func NewMaintenanceScheduleRepository(db *gorm.DB) *MaintenanceScheduleRepository {
	return &MaintenanceScheduleRepository{db: db}
}

// ===== Filters and Sorts =====

func (r *MaintenanceScheduleRepository) applyScheduleFilters(db *gorm.DB, filters *domain.MaintenanceScheduleFilterOptions) *gorm.DB {
	if filters == nil {
		return db
	}
	if filters.AssetID != nil && *filters.AssetID != "" {
		db = db.Where("maintenance_schedules.asset_id = ?", filters.AssetID)
	}
	if filters.MaintenanceType != nil && *filters.MaintenanceType != "" {
		db = db.Where("maintenance_schedules.maintenance_type = ?", filters.MaintenanceType)
	}
	if filters.State != nil && *filters.State != "" {
		db = db.Where("maintenance_schedules.state = ?", filters.State)
	}
	if filters.CreatedBy != nil && *filters.CreatedBy != "" {
		db = db.Where("maintenance_schedules.created_by = ?", filters.CreatedBy)
	}
	if filters.FromDate != nil && *filters.FromDate != "" {
		db = db.Where("maintenance_schedules.next_scheduled_date >= ?", *filters.FromDate)
	}
	if filters.ToDate != nil && *filters.ToDate != "" {
		db = db.Where("maintenance_schedules.next_scheduled_date <= ?", *filters.ToDate)
	}
	return db
}

func (r *MaintenanceScheduleRepository) applyScheduleSorts(db *gorm.DB, sort *domain.MaintenanceScheduleSortOptions) *gorm.DB {
	if sort == nil || sort.Field == "" {
		return db.Order("maintenance_schedules.created_at DESC")
	}

	// Map camelCase sort field to snake_case database column
	columnName := mapper.MapMaintenanceScheduleSortFieldToColumn(sort.Field)
	orderClause := "maintenance_schedules." + columnName

	order := "DESC"
	if sort.Order == domain.SortOrderAsc {
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

	// Return created maintenance schedule with translations (no need to query again)
	// GORM has already filled the model with created data including ID and timestamps
	domainSchedule := mapper.ToDomainMaintenanceSchedule(&m)
	// Add translations manually since they were created separately
	for _, translation := range payload.Translations {
		domainSchedule.Translations = append(domainSchedule.Translations, domain.MaintenanceScheduleTranslation{
			LangCode:    translation.LangCode,
			Title:       translation.Title,
			Description: translation.Description,
		})
	}
	return domainSchedule, nil
}

func (r *MaintenanceScheduleRepository) BulkCreateMaintenanceSchedules(ctx context.Context, schedules []domain.MaintenanceSchedule) ([]domain.MaintenanceSchedule, error) {
	if len(schedules) == 0 {
		return []domain.MaintenanceSchedule{}, nil
	}

	models := make([]*model.MaintenanceSchedule, len(schedules))
	for i := range schedules {
		m := mapper.ToModelMaintenanceScheduleForCreate(&schedules[i])
		models[i] = &m
	}

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.
		Omit(clause.Associations).
		Session(&gorm.Session{CreateBatchSize: 500}).
		Create(&models).Error; err != nil {
		tx.Rollback()
		return nil, domain.ErrInternal(err)
	}

	// Insert translations in batch
	var translations []model.MaintenanceScheduleTranslation
	for i := range models {
		s := schedules[i]
		for _, t := range s.Translations {
			mt := mapper.ToModelMaintenanceScheduleTranslationForCreate(models[i].ID.String(), &t)
			translations = append(translations, mt)
		}
	}
	if len(translations) > 0 {
		if err := tx.Session(&gorm.Session{CreateBatchSize: 500}).Create(&translations).Error; err != nil {
			tx.Rollback()
			return nil, domain.ErrInternal(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	created := make([]domain.MaintenanceSchedule, len(models))
	for i := range models {
		created[i] = mapper.ToDomainMaintenanceSchedule(models[i])
		for _, t := range schedules[i].Translations {
			created[i].Translations = append(created[i].Translations, domain.MaintenanceScheduleTranslation{
				LangCode:    t.LangCode,
				Title:       t.Title,
				Description: t.Description,
			})
		}
	}
	return created, nil
}

// BulkCreateSchedules satisfies the service repository interface by delegating
// to the existing bulk creation logic used elsewhere in the codebase.
func (r *MaintenanceScheduleRepository) BulkCreateSchedules(ctx context.Context, schedules []domain.MaintenanceSchedule) ([]domain.MaintenanceSchedule, error) {
	return r.BulkCreateMaintenanceSchedules(ctx, schedules)
}

func (r *MaintenanceScheduleRepository) UpdateSchedule(ctx context.Context, scheduleId string, payload *domain.UpdateMaintenanceSchedulePayload) (domain.MaintenanceSchedule, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.MaintenanceSchedule{}, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// Update maintenance schedule basic info
	updates := mapper.ToModelMaintenanceScheduleUpdateMap(payload)
	if len(updates) > 0 {
		if err := tx.Table("maintenance_schedules").Where("id = ?", scheduleId).Updates(updates).Error; err != nil {
			tx.Rollback()
			return domain.MaintenanceSchedule{}, domain.ErrInternal(err)
		}
	}

	// Update translations if provided
	if len(payload.Translations) > 0 {
		// Delete existing translations
		if err := tx.Where("schedule_id = ?", scheduleId).Delete(&model.MaintenanceScheduleTranslation{}).Error; err != nil {
			tx.Rollback()
			return domain.MaintenanceSchedule{}, domain.ErrInternal(err)
		}

		// Create new translations
		for _, translationPayload := range payload.Translations {
			translation := domain.MaintenanceScheduleTranslation{
				LangCode:    translationPayload.LangCode,
				Title:       *translationPayload.Title,
				Description: translationPayload.Description,
			}
			modelTranslation := mapper.ToModelMaintenanceScheduleTranslationForCreate(scheduleId, &translation)
			if err := tx.Create(&modelTranslation).Error; err != nil {
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
	if err := tx.Commit().Error; err != nil {
		return domain.ErrInternal(err)
	}
	return nil
}

func (r *MaintenanceScheduleRepository) BulkDeleteSchedules(ctx context.Context, scheduleIds []string) (domain.BulkDeleteMaintenanceSchedules, error) {
	result := domain.BulkDeleteMaintenanceSchedules{
		RequestedIDS: scheduleIds,
		DeletedIDS:   []string{},
	}

	if len(scheduleIds) == 0 {
		return result, nil
	}

	// First, find which schedules actually exist
	var existingSchedules []model.MaintenanceSchedule
	if err := r.db.WithContext(ctx).Select("id").Where("id IN ?", scheduleIds).Find(&existingSchedules).Error; err != nil {
		return result, domain.ErrInternal(err)
	}

	// Collect existing schedule IDs
	existingIds := make([]string, 0, len(existingSchedules))
	for _, schedule := range existingSchedules {
		existingIds = append(existingIds, schedule.ID.String())
	}

	// If no schedules exist, return early
	if len(existingIds) == 0 {
		return result, nil
	}

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return result, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete translations first (foreign key constraint)
	if err := tx.Delete(&model.MaintenanceScheduleTranslation{}, "schedule_id IN ?", existingIds).Error; err != nil {
		tx.Rollback()
		return result, domain.ErrInternal(err)
	}

	// Delete maintenance schedules
	if err := tx.Delete(&model.MaintenanceSchedule{}, "id IN ?", existingIds).Error; err != nil {
		tx.Rollback()
		return result, domain.ErrInternal(err)
	}

	if err := tx.Commit().Error; err != nil {
		return result, domain.ErrInternal(err)
	}

	result.DeletedIDS = existingIds
	return result, nil
}

// ===== QUERIES =====

func (r *MaintenanceScheduleRepository) GetSchedulesPaginated(ctx context.Context, params domain.MaintenanceScheduleParams, langCode string) ([]domain.MaintenanceSchedule, error) {
	var schedules []model.MaintenanceSchedule
	db := r.db.WithContext(ctx).
		Preload("Translations").
		Preload("Asset").
		Preload("Asset.Category").
		Preload("Asset.Category.Translations").
		Preload("Asset.Location").
		Preload("Asset.Location.Translations").
		Preload("Asset.User").
		Preload("CreatedByUser")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		sq := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN maintenance_schedule_translations mst ON maintenance_schedules.id = mst.schedule_id").
			Joins("LEFT JOIN assets a ON maintenance_schedules.asset_id = a.id").
			Where("mst.title ILIKE ? OR a.asset_name ILIKE ?", sq, sq).
			Group("maintenance_schedules.id")
	}

	// Apply filters, sorts, and pagination manually
	db = r.applyScheduleFilters(db, params.Filters)
	db = r.applyScheduleSorts(db, params.Sort)
	if params.Pagination != nil {
		if params.Pagination.Limit > 0 {
			db = db.Limit(params.Pagination.Limit)
		}
		if params.Pagination.Offset > 0 {
			db = db.Offset(params.Pagination.Offset)
		}
	}

	if err := db.Find(&schedules).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}
	return mapper.ToDomainMaintenanceSchedules(schedules), nil
}

func (r *MaintenanceScheduleRepository) GetSchedulesCursor(ctx context.Context, params domain.MaintenanceScheduleParams, langCode string) ([]domain.MaintenanceSchedule, error) {
	var schedules []model.MaintenanceSchedule
	db := r.db.WithContext(ctx).
		Preload("Translations").
		Preload("Asset").
		Preload("Asset.Category").
		Preload("Asset.Category.Translations").
		Preload("Asset.Location").
		Preload("Asset.Location.Translations").
		Preload("Asset.User").
		Preload("CreatedByUser")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		sq := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN maintenance_schedule_translations mst ON maintenance_schedules.id = mst.schedule_id").
			Joins("LEFT JOIN assets a ON maintenance_schedules.asset_id = a.id").
			Where("mst.title ILIKE ? OR a.asset_name ILIKE ?", sq, sq).
			Group("maintenance_schedules.id")
	}

	// Apply filters
	db = r.applyScheduleFilters(db, params.Filters)

	// Apply sorting - for cursor pagination, we need consistent ordering by ID
	if params.Sort != nil && params.Sort.Field != "" {
		db = r.applyScheduleSorts(db, params.Sort)
		// Always add secondary sort by ID DESC for consistency (ULID = newer = larger)
		db = db.Order("maintenance_schedules.id DESC")
	} else {
		// Default to ID DESC for cursor pagination (newest first)
		db = db.Order("maintenance_schedules.id DESC")
	}

	// Apply cursor-based pagination
	if params.Pagination != nil {
		if params.Pagination.Cursor != "" {
			db = db.Where("maintenance_schedules.id < ?", params.Pagination.Cursor)
		}
		if params.Pagination.Limit > 0 {
			db = db.Limit(params.Pagination.Limit)
		}
	}

	if err := db.Find(&schedules).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}
	return mapper.ToDomainMaintenanceSchedules(schedules), nil
}

func (r *MaintenanceScheduleRepository) GetScheduleById(ctx context.Context, scheduleId string) (domain.MaintenanceSchedule, error) {
	var m model.MaintenanceSchedule
	err := r.db.WithContext(ctx).
		Preload("Translations").
		Preload("Asset").
		Preload("Asset.Category").
		Preload("Asset.Category.Translations").
		Preload("Asset.Location").
		Preload("Asset.Location.Translations").
		Preload("Asset.User").
		Preload("CreatedByUser").
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

func (r *MaintenanceScheduleRepository) CountSchedules(ctx context.Context, params domain.MaintenanceScheduleParams) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(&model.MaintenanceSchedule{})
	if params.SearchQuery != nil && *params.SearchQuery != "" {
		sq := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN maintenance_schedule_translations mst ON maintenance_schedules.id = mst.schedule_id").
			Joins("LEFT JOIN assets a ON maintenance_schedules.asset_id = a.id").
			Where("mst.title ILIKE ? OR a.asset_name ILIKE ?", sq, sq).
			Group("maintenance_schedules.id")
	}
	db = r.applyScheduleFilters(db, params.Filters)
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
	var inspectionCount int64
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceSchedule{}).
		Where("maintenance_type = ?", domain.ScheduleTypeInspection).Count(&inspectionCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	var calibrationCount int64
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceSchedule{}).
		Where("maintenance_type = ?", domain.ScheduleTypeCalibration).Count(&calibrationCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.ByType.Preventive = int(preventiveCount)
	stats.ByType.Corrective = int(correctiveCount)
	stats.ByType.Inspection = int(inspectionCount)
	stats.ByType.Calibration = int(calibrationCount)

	// By status
	var activeCount, pausedCount, stoppedCount, completedCount int64
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceSchedule{}).Where("state = ?", domain.StateActive).Count(&activeCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceSchedule{}).Where("state = ?", domain.StatePaused).Count(&pausedCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceSchedule{}).Where("state = ?", domain.StateStopped).Count(&stoppedCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceSchedule{}).Where("state = ?", domain.StateCompleted).Count(&completedCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.ByStatus.Active = int(activeCount)
	stats.ByStatus.Paused = int(pausedCount)
	stats.ByStatus.Stopped = int(stoppedCount)
	stats.ByStatus.Completed = int(completedCount)

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
			"MIN(CASE WHEN ms.next_scheduled_date > NOW() AND ms.state = 'Active' THEN ms.next_scheduled_date::text ELSE NULL END) as next_maintenance").
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
		ID                string
		AssetID           string
		AssetName         string
		AssetTag          string
		MaintenanceType   domain.MaintenanceScheduleType
		NextScheduledDate string
		DaysUntilDue      int
		Title             string
		Description       *string
	}
	if err := r.db.WithContext(ctx).
		Table("maintenance_schedules ms").
		Select("ms.id, a.id as asset_id, a.asset_name, a.asset_tag, ms.maintenance_type, " +
			"ms.next_scheduled_date::text, EXTRACT(DAY FROM ms.next_scheduled_date - NOW()) as days_until_due, " +
			"mst.title, mst.description").
		Joins("LEFT JOIN assets a ON ms.asset_id = a.id").
		Joins("LEFT JOIN maintenance_schedule_translations mst ON ms.id = mst.schedule_id").
		Where("ms.next_scheduled_date > NOW() AND ms.state = 'Active'").
		Order("ms.next_scheduled_date ASC").
		Limit(10).
		Scan(&upcomingResults).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	for _, result := range upcomingResults {
		stats.UpcomingSchedule = append(stats.UpcomingSchedule, domain.UpcomingMaintenanceSchedule{
			ID:                result.ID,
			AssetID:           result.AssetID,
			AssetName:         result.AssetName,
			AssetTag:          result.AssetTag,
			MaintenanceType:   result.MaintenanceType,
			NextScheduledDate: result.NextScheduledDate,
			DaysUntilDue:      result.DaysUntilDue,
			Title:             result.Title,
			Description:       result.Description,
		})
	}

	// Overdue schedules
	var overdueResults []struct {
		ID                string
		AssetID           string
		AssetName         string
		AssetTag          string
		MaintenanceType   domain.MaintenanceScheduleType
		NextScheduledDate string
		DaysOverdue       int
		Title             string
		Description       *string
	}
	if err := r.db.WithContext(ctx).
		Table("maintenance_schedules ms").
		Select("ms.id, a.id as asset_id, a.asset_name, a.asset_tag, ms.maintenance_type, " +
			"ms.next_scheduled_date::text, EXTRACT(DAY FROM NOW() - ms.next_scheduled_date) as days_overdue, " +
			"mst.title, mst.description").
		Joins("LEFT JOIN assets a ON ms.asset_id = a.id").
		Joins("LEFT JOIN maintenance_schedule_translations mst ON ms.id = mst.schedule_id").
		Where("ms.next_scheduled_date < NOW() AND ms.state = 'Active'").
		Order("ms.next_scheduled_date ASC").
		Scan(&overdueResults).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	for _, result := range overdueResults {
		stats.OverdueSchedule = append(stats.OverdueSchedule, domain.OverdueMaintenanceSchedule{
			ID:                result.ID,
			AssetID:           result.AssetID,
			AssetName:         result.AssetName,
			AssetTag:          result.AssetTag,
			MaintenanceType:   result.MaintenanceType,
			NextScheduledDate: result.NextScheduledDate,
			DaysOverdue:       result.DaysOverdue,
			Title:             result.Title,
			Description:       result.Description,
		})
	}

	// Frequency trends
	var frequencyResults []struct {
		FrequencyMonths int
		Count           int64
	}
	if err := r.db.WithContext(ctx).
		Table("maintenance_schedules").
		Select("CASE WHEN interval_unit = 'Months' THEN interval_value ELSE interval_value * 12 END as frequency_months, COUNT(*) as count").
		Where("interval_value IS NOT NULL AND interval_unit IS NOT NULL").
		Group("CASE WHEN interval_unit = 'Months' THEN interval_value ELSE interval_value * 12 END").
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
		stats.Summary.ActiveMaintenancePercentage = float64(activeCount) / float64(total) * 100
		stats.Summary.CompletedMaintenancePercentage = float64(completedCount) / float64(total) * 100
		stats.Summary.PausedMaintenancePercentage = float64(pausedCount) / float64(total) * 100
		stats.Summary.StoppedMaintenancePercentage = float64(stoppedCount) / float64(total) * 100
		stats.Summary.PreventiveMaintenancePercentage = float64(preventiveCount) / float64(total) * 100
		stats.Summary.CorrectiveMaintenancePercentage = float64(correctiveCount) / float64(total) * 100
		stats.Summary.InspectionMaintenancePercentage = float64(inspectionCount) / float64(total) * 100
		stats.Summary.CalibrationMaintenancePercentage = float64(calibrationCount) / float64(total) * 100
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
		stats.Summary.EarliestScheduleDate = earliest
	}
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceSchedule{}).Select("MAX(created_at)").Scan(&latest).Error; err == nil {
		stats.Summary.LatestScheduleDate = latest
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

// GetSchedulesDueSoon retrieves maintenance schedules that are due within the specified number of days from now
func (r *MaintenanceScheduleRepository) GetSchedulesDueSoon(ctx context.Context, daysFromNow int) ([]domain.MaintenanceSchedule, error) {
	var models []model.MaintenanceSchedule
	now := time.Now().UTC()
	futureDate := now.AddDate(0, 0, daysFromNow)

	err := r.db.WithContext(ctx).
		Preload("Translations").
		Preload("Asset").
		Preload("Asset.User").
		Where("next_scheduled_date >= ? AND next_scheduled_date <= ? AND state = ?",
			now, futureDate, domain.StateActive).
		Find(&models).Error

	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	schedules := make([]domain.MaintenanceSchedule, len(models))
	for i, m := range models {
		schedules[i] = mapper.ToDomainMaintenanceSchedule(&m)
	}

	return schedules, nil
}

// GetOverdueSchedules retrieves maintenance schedules that are overdue (past scheduled date and still active)
func (r *MaintenanceScheduleRepository) GetOverdueSchedules(ctx context.Context) ([]domain.MaintenanceSchedule, error) {
	var models []model.MaintenanceSchedule
	now := time.Now().UTC()

	err := r.db.WithContext(ctx).
		Preload("Translations").
		Preload("Asset").
		Preload("Asset.User").
		Where("next_scheduled_date < ? AND state = ?", now, domain.StateActive).
		Find(&models).Error

	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	schedules := make([]domain.MaintenanceSchedule, len(models))
	for i, m := range models {
		schedules[i] = mapper.ToDomainMaintenanceSchedule(&m)
	}

	return schedules, nil
}

// GetRecurringSchedulesToUpdate retrieves recurring schedules that need next_scheduled_date update
func (r *MaintenanceScheduleRepository) GetRecurringSchedulesToUpdate(ctx context.Context) ([]domain.MaintenanceSchedule, error) {
	var models []model.MaintenanceSchedule
	now := time.Now().UTC()

	err := r.db.WithContext(ctx).
		Preload("Translations").
		Preload("Asset").
		Where("is_recurring = ? AND next_scheduled_date < ? AND state = ?",
			true, now, domain.StateActive).
		Find(&models).Error

	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	schedules := make([]domain.MaintenanceSchedule, len(models))
	for i, m := range models {
		schedules[i] = mapper.ToDomainMaintenanceSchedule(&m)
	}

	return schedules, nil
}

// UpdateLastExecutedDate updates the last_executed_date field for a schedule
func (r *MaintenanceScheduleRepository) UpdateLastExecutedDate(ctx context.Context, scheduleId string, lastExecutedDate *time.Time) error {
	result := r.db.WithContext(ctx).
		Model(&model.MaintenanceSchedule{}).
		Where("id = ?", scheduleId).
		Update("last_executed_date", lastExecutedDate)

	if result.Error != nil {
		return domain.ErrInternal(result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound("maintenance schedule not found")
	}

	return nil
}

func (r *MaintenanceScheduleRepository) GetMaintenanceSchedulesForExport(ctx context.Context, params domain.MaintenanceScheduleParams, langCode string) ([]domain.MaintenanceSchedule, error) {
	var schedules []model.MaintenanceSchedule
	db := r.db.WithContext(ctx).
		Preload("Translations").
		Preload("Asset").
		Preload("Asset.Category").
		Preload("Asset.Category.Translations").
		Preload("Asset.Location").
		Preload("Asset.Location.Translations").
		Preload("Asset.User").
		Preload("CreatedByUser")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		sq := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN maintenance_schedule_translations mst ON maintenance_schedules.id = mst.schedule_id").
			Joins("LEFT JOIN assets a ON maintenance_schedules.asset_id = a.id").
			Where("mst.title ILIKE ? OR a.asset_name ILIKE ?", sq, sq).
			Group("maintenance_schedules.id")
	}

	// Apply filters
	db = r.applyScheduleFilters(db, params.Filters)

	// Apply sorting
	db = r.applyScheduleSorts(db, params.Sort)

	// No pagination for export - get all matching maintenance schedules
	if err := db.Find(&schedules).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	return mapper.ToDomainMaintenanceSchedules(schedules), nil
}
