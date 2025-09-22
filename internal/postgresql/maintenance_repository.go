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

type MaintenanceRepository struct {
	db *gorm.DB
}

func NewMaintenanceRepository(db *gorm.DB) *MaintenanceRepository {
	return &MaintenanceRepository{db: db}
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

type MaintenanceRecordFilterOptions struct {
	AssetID         *string `json:"assetId,omitempty"`
	ScheduleID      *string `json:"scheduleId,omitempty"`
	PerformedByUser *string `json:"performedByUser,omitempty"`
	VendorName      *string `json:"vendorName,omitempty"`
	FromDate        *string `json:"fromDate,omitempty"` // YYYY-MM-DD
	ToDate          *string `json:"toDate,omitempty"`   // YYYY-MM-DD
}

func (r *MaintenanceRepository) applyScheduleFilters(db *gorm.DB, filters any) *gorm.DB {
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

func (r *MaintenanceRepository) applyScheduleSorts(db *gorm.DB, sort *query.SortOptions) *gorm.DB {
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

func (r *MaintenanceRepository) applyRecordFilters(db *gorm.DB, filters any) *gorm.DB {
	f, ok := filters.(*MaintenanceRecordFilterOptions)
	if !ok || f == nil {
		return db
	}
	if f.AssetID != nil && *f.AssetID != "" {
		db = db.Where("mr.asset_id = ?", f.AssetID)
	}
	if f.ScheduleID != nil && *f.ScheduleID != "" {
		db = db.Where("mr.schedule_id = ?", f.ScheduleID)
	}
	if f.PerformedByUser != nil && *f.PerformedByUser != "" {
		db = db.Where("mr.performed_by_user = ?", f.PerformedByUser)
	}
	if f.VendorName != nil && *f.VendorName != "" {
		db = db.Where("mr.performed_by_vendor ILIKE ?", "%"+*f.VendorName+"%")
	}
	if f.FromDate != nil && *f.FromDate != "" {
		db = db.Where("mr.maintenance_date >= ?", *f.FromDate)
	}
	if f.ToDate != nil && *f.ToDate != "" {
		db = db.Where("mr.maintenance_date <= ?", *f.ToDate)
	}
	return db
}

func (r *MaintenanceRepository) applyRecordSorts(db *gorm.DB, sort *query.SortOptions) *gorm.DB {
	if sort == nil || sort.Field == "" {
		return db.Order("mr.created_at DESC")
	}
	var orderClause string
	switch strings.ToLower(sort.Field) {
	case "maintenance_date", "created_at", "updated_at":
		orderClause = "mr." + sort.Field
	case "title":
		orderClause = "mrt.title"
	default:
		return db.Order("mr.created_at DESC")
	}
	order := "DESC"
	if strings.ToLower(sort.Order) == "asc" {
		order = "ASC"
	}
	return db.Order(fmt.Sprintf("%s %s", orderClause, order))
}

// ===== MUTATIONS: Schedule =====

func (r *MaintenanceRepository) CreateSchedule(ctx context.Context, payload *domain.MaintenanceSchedule) (domain.MaintenanceSchedule, error) {
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

func (r *MaintenanceRepository) UpdateSchedule(ctx context.Context, scheduleId string, payload *domain.MaintenanceSchedule) (domain.MaintenanceSchedule, error) {
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

func (r *MaintenanceRepository) DeleteSchedule(ctx context.Context, scheduleId string) error {
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

// ===== MUTATIONS: Record =====

func (r *MaintenanceRepository) CreateRecord(ctx context.Context, payload *domain.MaintenanceRecord) (domain.MaintenanceRecord, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.MaintenanceRecord{}, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	m := mapper.ToModelMaintenanceRecordForCreate(payload)
	if err := tx.Create(&m).Error; err != nil {
		tx.Rollback()
		return domain.MaintenanceRecord{}, domain.ErrInternal(err)
	}
	for _, t := range payload.Translations {
		mt := mapper.ToModelMaintenanceRecordTranslationForCreate(m.ID.String(), &t)
		if err := tx.Create(&mt).Error; err != nil {
			tx.Rollback()
			return domain.MaintenanceRecord{}, domain.ErrInternal(err)
		}
	}
	if err := tx.Commit().Error; err != nil {
		return domain.MaintenanceRecord{}, domain.ErrInternal(err)
	}
	return r.GetRecordById(ctx, m.ID.String())
}

func (r *MaintenanceRepository) UpdateRecord(ctx context.Context, recordId string, payload *domain.MaintenanceRecord) (domain.MaintenanceRecord, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.MaintenanceRecord{}, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	updates := map[string]any{}
	if payload.AssetID != "" {
		if parsed, err := ulid.Parse(payload.AssetID); err == nil {
			updates["asset_id"] = model.SQLULID(parsed)
		}
	}
	if payload.ScheduleID != nil {
		if *payload.ScheduleID == "" {
			updates["schedule_id"] = nil
		} else if parsed, err := ulid.Parse(*payload.ScheduleID); err == nil {
			v := model.SQLULID(parsed)
			updates["schedule_id"] = &v
		}
	}
	if !payload.MaintenanceDate.IsZero() {
		updates["maintenance_date"] = payload.MaintenanceDate
	}
	if payload.PerformedByUser != nil {
		if *payload.PerformedByUser == "" {
			updates["performed_by_user"] = nil
		} else if parsed, err := ulid.Parse(*payload.PerformedByUser); err == nil {
			v := model.SQLULID(parsed)
			updates["performed_by_user"] = &v
		}
	}
	if payload.PerformedByVendor != nil {
		updates["performed_by_vendor"] = payload.PerformedByVendor
	}
	if payload.ActualCost != nil {
		updates["actual_cost"] = payload.ActualCost
	}

	if len(updates) > 0 {
		if err := tx.Model(&model.MaintenanceRecord{}).Where("id = ?", recordId).Updates(updates).Error; err != nil {
			tx.Rollback()
			return domain.MaintenanceRecord{}, domain.ErrInternal(err)
		}
	}

	if len(payload.Translations) > 0 {
		if err := tx.Delete(&model.MaintenanceRecordTranslation{}, "record_id = ?", recordId).Error; err != nil {
			tx.Rollback()
			return domain.MaintenanceRecord{}, domain.ErrInternal(err)
		}
		for _, t := range payload.Translations {
			mt := mapper.ToModelMaintenanceRecordTranslationForCreate(recordId, &t)
			if err := tx.Create(&mt).Error; err != nil {
				tx.Rollback()
				return domain.MaintenanceRecord{}, domain.ErrInternal(err)
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return domain.MaintenanceRecord{}, domain.ErrInternal(err)
	}
	return r.GetRecordById(ctx, recordId)
}

func (r *MaintenanceRepository) DeleteRecord(ctx context.Context, recordId string) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.ErrInternal(tx.Error)
	}
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Delete(&model.MaintenanceRecordTranslation{}, "record_id = ?", recordId).Error; err != nil {
		tx.Rollback()
		return domain.ErrInternal(err)
	}
	if err := tx.Delete(&model.MaintenanceRecord{}, "id = ?", recordId).Error; err != nil {
		tx.Rollback()
		return domain.ErrInternal(err)
	}
	return tx.Commit().Error
}

// ===== QUERIES =====

func (r *MaintenanceRepository) GetSchedulesPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceScheduleListItem, error) {
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
	return mapper.MaintenanceSchedulesToListItems(mapper.ToDomainMaintenanceSchedules(schedules), langCode), nil
}

func (r *MaintenanceRepository) GetSchedulesCursor(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceScheduleListItem, error) {
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
	return mapper.MaintenanceSchedulesToListItems(mapper.ToDomainMaintenanceSchedules(schedules), langCode), nil
}

func (r *MaintenanceRepository) GetScheduleById(ctx context.Context, scheduleId string) (domain.MaintenanceSchedule, error) {
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

func (r *MaintenanceRepository) GetRecordsPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceRecordListItem, error) {
	var records []model.MaintenanceRecord
	db := r.db.WithContext(ctx).
		Table("maintenance_records mr").
		Preload("Translations").
		Preload("Asset").
		Preload("PerformedByUser")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		sq := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN maintenance_record_translations mrt ON mr.id = mrt.record_id").
			Where("mrt.title ILIKE ?", sq).
			Distinct("mr.id")
	}

	params.Pagination.Cursor = ""
	db = query.Apply(db, params, r.applyRecordFilters, r.applyRecordSorts)
	if err := db.Find(&records).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}
	return mapper.MaintenanceRecordsToListItems(mapper.ToDomainMaintenanceRecords(records), langCode), nil
}

func (r *MaintenanceRepository) GetRecordsCursor(ctx context.Context, params query.Params, langCode string) ([]domain.MaintenanceRecordListItem, error) {
	var records []model.MaintenanceRecord
	db := r.db.WithContext(ctx).
		Table("maintenance_records mr").
		Preload("Translations").
		Preload("Asset").
		Preload("PerformedByUser")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		sq := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN maintenance_record_translations mrt ON mr.id = mrt.record_id").
			Where("mrt.title ILIKE ?", sq).
			Distinct("mr.id")
	}

	params.Pagination.Offset = 0
	db = query.Apply(db, params, r.applyRecordFilters, r.applyRecordSorts)
	if err := db.Find(&records).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}
	return mapper.MaintenanceRecordsToListItems(mapper.ToDomainMaintenanceRecords(records), langCode), nil
}

func (r *MaintenanceRepository) GetRecordById(ctx context.Context, recordId string) (domain.MaintenanceRecord, error) {
	var m model.MaintenanceRecord
	err := r.db.WithContext(ctx).
		Table("maintenance_records mr").
		Preload("Translations").
		First(&m, "id = ?", recordId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.MaintenanceRecord{}, domain.ErrNotFound("maintenance_record")
		}
		return domain.MaintenanceRecord{}, domain.ErrInternal(err)
	}
	return mapper.ToDomainMaintenanceRecord(&m), nil
}

// ===== Checks and Counts =====

func (r *MaintenanceRepository) CheckScheduleExist(ctx context.Context, scheduleId string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceSchedule{}).Where("id = ?", scheduleId).Count(&count).Error; err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *MaintenanceRepository) CheckRecordExist(ctx context.Context, recordId string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceRecord{}).Where("id = ?", recordId).Count(&count).Error; err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *MaintenanceRepository) CountSchedules(ctx context.Context, params query.Params) (int64, error) {
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

func (r *MaintenanceRepository) CountRecords(ctx context.Context, params query.Params) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Table("maintenance_records mr")
	if params.SearchQuery != nil && *params.SearchQuery != "" {
		sq := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN maintenance_record_translations mrt ON mr.id = mrt.record_id").
			Where("mrt.title ILIKE ?", sq).
			Distinct("mr.id")
	}
	db = query.Apply(db, query.Params{Filters: params.Filters}, r.applyRecordFilters, nil)
	if err := db.Count(&count).Error; err != nil {
		return 0, domain.ErrInternal(err)
	}
	return count, nil
}

// ===== Statistics =====
// Implementing a comprehensive statistics function similar to assets/categories.
func (r *MaintenanceRepository) GetMaintenanceStatistics(ctx context.Context) (domain.MaintenanceStatistics, error) {
	var stats domain.MaintenanceStatistics

	// Schedules: total
	var schTotal int64
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceSchedule{}).Count(&schTotal).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.Schedules.Total.Count = int(schTotal)

	// Schedules by type
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
	stats.Schedules.ByType.Preventive = int(preventiveCount)
	stats.Schedules.ByType.Corrective = int(correctiveCount)

	// Schedules by status
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
	stats.Schedules.ByStatus.Scheduled = int(scheduledCount)
	stats.Schedules.ByStatus.Completed = int(completedCount)
	stats.Schedules.ByStatus.Cancelled = int(cancelledCount)

	// Records total
	var recTotal int64
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceRecord{}).Count(&recTotal).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.Records.Total.Count = int(recTotal)

	// Completion trend (last 30 days)
	var completionTrends []struct {
		Date  string
		Count int64
	}
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceRecord{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= NOW() - INTERVAL '30 days'").
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&completionTrends).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.Records.CompletionTrend = make([]domain.MaintenanceCompletionTrend, len(completionTrends))
	for i, ct := range completionTrends {
		stats.Records.CompletionTrend[i] = domain.MaintenanceCompletionTrend{Date: ct.Date, Count: int(ct.Count)}
	}

	// Summary basic
	stats.Summary.TotalSchedules = int(schTotal)
	stats.Summary.TotalRecords = int(recTotal)

	// Earliest/latest record
	var earliest, latest time.Time
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceRecord{}).Select("MIN(created_at)").Scan(&earliest).Error; err == nil {
		stats.Summary.EarliestRecordDate = earliest.Format("2006-01-02")
	}
	if err := r.db.WithContext(ctx).Model(&model.MaintenanceRecord{}).Select("MAX(created_at)").Scan(&latest).Error; err == nil {
		stats.Summary.LatestRecordDate = latest.Format("2006-01-02")
	}
	if !earliest.IsZero() && !latest.IsZero() {
		days := latest.Sub(earliest).Hours() / 24
		if days > 0 {
			stats.Summary.AverageRecordsPerDay = float64(recTotal) / days
		}
	}

	return stats, nil
}
