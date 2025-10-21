package postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type IssueReportRepository struct {
	db *gorm.DB
}

func NewIssueReportRepository(db *gorm.DB) *IssueReportRepository {
	return &IssueReportRepository{
		db: db,
	}
}

func (r *IssueReportRepository) applyIssueReportFilters(db *gorm.DB, filters *domain.IssueReportFilterOptions) *gorm.DB {
	if filters == nil {
		return db
	}

	if filters.AssetID != nil {
		db = db.Where("ir.asset_id = ?", filters.AssetID)
	}
	if filters.ReportedBy != nil {
		db = db.Where("ir.reported_by = ?", filters.ReportedBy)
	}
	if filters.ResolvedBy != nil {
		db = db.Where("ir.resolved_by = ?", filters.ResolvedBy)
	}
	if filters.IssueType != nil {
		db = db.Where("ir.issue_type = ?", filters.IssueType)
	}
	if filters.Priority != nil {
		db = db.Where("ir.priority = ?", filters.Priority)
	}
	if filters.Status != nil {
		db = db.Where("ir.status = ?", filters.Status)
	}
	if filters.IsResolved != nil {
		if *filters.IsResolved {
			db = db.Where("ir.status IN ('Resolved', 'Closed')")
		} else {
			db = db.Where("ir.status IN ('Open', 'In Progress')")
		}
	}
	if filters.DateFrom != nil {
		db = db.Where("ir.reported_date >= ?", filters.DateFrom)
	}
	if filters.DateTo != nil {
		db = db.Where("ir.reported_date <= ?", filters.DateTo)
	}
	return db
}

func (r *IssueReportRepository) applyIssueReportSorts(db *gorm.DB, sort *domain.IssueReportSortOptions) *gorm.DB {
	if sort == nil || sort.Field == "" {
		return db.Order("ir.reported_date DESC")
	}

	// Map camelCase sort field to snake_case database column
	columnName := mapper.MapIssueReportSortFieldToColumn(sort.Field)
	orderClause := columnName

	order := "DESC"
	if sort.Order == domain.SortOrderAsc {
		order = "ASC"
	}
	return db.Order(fmt.Sprintf("%s %s", orderClause, order))
}

// *===========================MUTATION===========================*
func (r *IssueReportRepository) CreateIssueReport(ctx context.Context, payload *domain.IssueReport) (domain.IssueReport, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.IssueReport{}, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create issue report
	modelIssueReport := mapper.ToModelIssueReportForCreate(payload)
	if err := tx.Create(&modelIssueReport).Error; err != nil {
		tx.Rollback()
		return domain.IssueReport{}, domain.ErrInternal(err)
	}

	// Create translations
	for _, translation := range payload.Translations {
		modelTranslation := mapper.ToModelIssueReportTranslationForCreate(modelIssueReport.ID.String(), &translation)
		if err := tx.Create(&modelTranslation).Error; err != nil {
			tx.Rollback()
			return domain.IssueReport{}, domain.ErrInternal(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return domain.IssueReport{}, domain.ErrInternal(err)
	}

	// Return created issue report with translations (no need to query again)
	// GORM has already filled the model with created data including ID and timestamps
	domainIssueReport := mapper.ToDomainIssueReport(&modelIssueReport)
	// Add translations manually since they were created separately
	for _, translation := range payload.Translations {
		domainIssueReport.Translations = append(domainIssueReport.Translations, domain.IssueReportTranslation{
			LangCode:    translation.LangCode,
			Title:       translation.Title,
			Description: translation.Description,
		})
	}
	return domainIssueReport, nil
}

func (r *IssueReportRepository) UpdateIssueReport(ctx context.Context, issueReportId string, payload *domain.UpdateIssueReportPayload) (domain.IssueReport, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.IssueReport{}, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get current issue report to check status change
	var currentReport model.IssueReport
	if err := tx.First(&currentReport, "id = ?", issueReportId).Error; err != nil {
		tx.Rollback()
		return domain.IssueReport{}, domain.ErrInternal(err)
	}

	// Update issue report basic info
	updates := mapper.ToModelIssueReportUpdateMap(payload)

	// Handle resolved_date based on status change
	if payload.Status != nil {
		if *payload.Status == domain.IssueStatusResolved {
			// If changing to resolved, set resolved_date
			now := time.Now()
			updates["resolved_date"] = &now
		} else if currentReport.Status == domain.IssueStatusResolved && *payload.Status != domain.IssueStatusResolved {
			// If changing from resolved to something else, clear resolved_date and resolved_by
			updates["resolved_date"] = nil
			updates["resolved_by"] = nil
		}
	}

	if len(updates) > 0 {
		if err := tx.Model(&model.IssueReport{}).Where("id = ?", issueReportId).Updates(updates).Error; err != nil {
			tx.Rollback()
			return domain.IssueReport{}, domain.ErrInternal(err)
		}
	}

	// Update translations if provided
	if len(payload.Translations) > 0 {
		for _, translationPayload := range payload.Translations {
			translationUpdates := mapper.ToModelIssueReportTranslationUpdateMap(&translationPayload)
			if len(translationUpdates) > 0 {
				// Try to update existing translation
				result := tx.Model(&model.IssueReportTranslation{}).
					Where("report_id = ? AND lang_code = ?", issueReportId, translationPayload.LangCode).
					Updates(translationUpdates)

				if result.Error != nil {
					tx.Rollback()
					return domain.IssueReport{}, domain.ErrInternal(result.Error)
				}

				// If no rows affected, create new translation
				if result.RowsAffected == 0 {
					newTranslation := model.IssueReportTranslation{
						LangCode:        translationPayload.LangCode,
						Title:           *translationPayload.Title,
						Description:     translationPayload.Description,
						ResolutionNotes: translationPayload.ResolutionNotes,
					}
					if parsedReportID, err := ulid.Parse(issueReportId); err == nil {
						newTranslation.ReportID = model.SQLULID(parsedReportID)
					}

					if err := tx.Create(&newTranslation).Error; err != nil {
						tx.Rollback()
						return domain.IssueReport{}, domain.ErrInternal(err)
					}
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return domain.IssueReport{}, domain.ErrInternal(err)
	}

	// Fetch updated issue report with translations and relations
	return r.GetIssueReportById(ctx, issueReportId)
}

func (r *IssueReportRepository) DeleteIssueReport(ctx context.Context, issueReportId string) error {
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
	if err := tx.Delete(&model.IssueReportTranslation{}, "report_id = ?", issueReportId).Error; err != nil {
		tx.Rollback()
		return domain.ErrInternal(err)
	}

	// Delete issue report
	if err := tx.Delete(&model.IssueReport{}, "id = ?", issueReportId).Error; err != nil {
		tx.Rollback()
		return domain.ErrInternal(err)
	}

	if err := tx.Commit().Error; err != nil {
		return domain.ErrInternal(err)
	}

	return nil
}

// *===========================QUERY===========================*
func (r *IssueReportRepository) GetIssueReportsPaginated(ctx context.Context, params domain.IssueReportParams, langCode string) ([]domain.IssueReport, error) {
	var issueReports []model.IssueReport
	db := r.db.WithContext(ctx).
		Table("issue_reports ir").
		Preload("Translations").
		Preload("Asset").
		Preload("Asset.Category").
		Preload("Asset.Category.Translations").
		Preload("Asset.Location").
		Preload("Asset.Location.Translations").
		Preload("Asset.User").
		Preload("ReportedByUser").
		Preload("ResolvedByUser")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN issue_report_translations irt ON ir.id = irt.report_id").
			Where("ir.issue_type ILIKE ? OR irt.title ILIKE ? OR irt.description ILIKE ?", searchPattern, searchPattern, searchPattern).
			Distinct("ir.id, ir.reported_date")
	}

	// Apply filters, sorts, and pagination manually
	db = r.applyIssueReportFilters(db, params.Filters)
	db = r.applyIssueReportSorts(db, params.Sort)
	if params.Pagination != nil {
		if params.Pagination.Limit > 0 {
			db = db.Limit(params.Pagination.Limit)
		}
		if params.Pagination.Offset > 0 {
			db = db.Offset(params.Pagination.Offset)
		}
	}

	if err := db.Find(&issueReports).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain issue reports
	return mapper.ToDomainIssueReports(issueReports), nil
}

func (r *IssueReportRepository) GetIssueReportsCursor(ctx context.Context, params domain.IssueReportParams, langCode string) ([]domain.IssueReport, error) {
	var issueReports []model.IssueReport
	db := r.db.WithContext(ctx).
		Table("issue_reports ir").
		Preload("Translations").
		Preload("Asset").
		Preload("Asset.Category").
		Preload("Asset.Category.Translations").
		Preload("Asset.Location").
		Preload("Asset.Location.Translations").
		Preload("Asset.User").
		Preload("ReportedByUser").
		Preload("ResolvedByUser")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN issue_report_translations irt ON ir.id = irt.report_id").
			Where("ir.issue_type ILIKE ? OR irt.title ILIKE ? OR irt.description ILIKE ?", searchPattern, searchPattern, searchPattern).
			Distinct("ir.id, ir.reported_date")
	}

	// Apply filters, sorts, and cursor pagination manually
	db = r.applyIssueReportFilters(db, params.Filters)
	db = r.applyIssueReportSorts(db, params.Sort)
	if params.Pagination != nil {
		if params.Pagination.Cursor != "" {
			db = db.Where("ir.id < ?", params.Pagination.Cursor)
		}
		if params.Pagination.Limit > 0 {
			db = db.Limit(params.Pagination.Limit)
		}
	}

	if err := db.Find(&issueReports).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain issue reports
	return mapper.ToDomainIssueReports(issueReports), nil
}

func (r *IssueReportRepository) GetIssueReportById(ctx context.Context, issueReportId string) (domain.IssueReport, error) {
	var issueReport model.IssueReport

	err := r.db.WithContext(ctx).
		Table("issue_reports ir").
		Preload("Translations").
		Preload("Asset").
		Preload("Asset.Category").
		Preload("Asset.Category.Translations").
		Preload("Asset.Location").
		Preload("Asset.Location.Translations").
		Preload("Asset.User").
		Preload("ReportedByUser").
		Preload("ResolvedByUser").
		First(&issueReport, "id = ?", issueReportId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.IssueReport{}, domain.ErrNotFound("issue report")
		}
		return domain.IssueReport{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainIssueReport(&issueReport), nil
}

func (r *IssueReportRepository) CheckIssueReportExist(ctx context.Context, issueReportId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.IssueReport{}).Where("id = ?", issueReportId).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *IssueReportRepository) CountIssueReports(ctx context.Context, params domain.IssueReportParams) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Table("issue_reports ir")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN issue_report_translations irt ON ir.id = irt.report_id").
			Where("ir.issue_type ILIKE ? OR irt.title ILIKE ? OR irt.description ILIKE ?", searchPattern, searchPattern, searchPattern).
			Distinct("ir.id, ir.reported_date")
	}

	db = r.applyIssueReportFilters(db, params.Filters)

	if err := db.Count(&count).Error; err != nil {
		return 0, domain.ErrInternal(err)
	}
	return count, nil
}

func (r *IssueReportRepository) GetIssueReportStatistics(ctx context.Context) (domain.IssueReportStatistics, error) {
	var stats domain.IssueReportStatistics

	// Get total issue report count
	var totalCount int64
	if err := r.db.WithContext(ctx).Model(&model.IssueReport{}).Count(&totalCount).Error; err != nil {
		return domain.IssueReportStatistics{}, domain.ErrInternal(err)
	}
	stats.Total.Count = int(totalCount)

	// Get priority statistics
	var priorityStats []struct {
		Priority domain.IssuePriority `json:"priority"`
		Count    int64                `json:"count"`
	}
	if err := r.db.WithContext(ctx).Model(&model.IssueReport{}).
		Select("priority, COUNT(*) as count").
		Group("priority").
		Find(&priorityStats).Error; err != nil {
		return domain.IssueReportStatistics{}, domain.ErrInternal(err)
	}

	for _, ps := range priorityStats {
		switch ps.Priority {
		case domain.PriorityLow:
			stats.ByPriority.Low = int(ps.Count)
		case domain.PriorityMedium:
			stats.ByPriority.Medium = int(ps.Count)
		case domain.PriorityHigh:
			stats.ByPriority.High = int(ps.Count)
		case domain.PriorityCritical:
			stats.ByPriority.Critical = int(ps.Count)
		}
	}

	// Get status statistics
	var statusStats []struct {
		Status domain.IssueStatus `json:"status"`
		Count  int64              `json:"count"`
	}
	if err := r.db.WithContext(ctx).Model(&model.IssueReport{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Find(&statusStats).Error; err != nil {
		return domain.IssueReportStatistics{}, domain.ErrInternal(err)
	}

	for _, ss := range statusStats {
		switch ss.Status {
		case domain.IssueStatusOpen:
			stats.ByStatus.Open = int(ss.Count)
		case domain.IssueStatusInProgress:
			stats.ByStatus.InProgress = int(ss.Count)
		case domain.IssueStatusResolved:
			stats.ByStatus.Resolved = int(ss.Count)
		case domain.IssueStatusClosed:
			stats.ByStatus.Closed = int(ss.Count)
		}
	}

	// Get type statistics
	var typeStats []struct {
		IssueType string `json:"issue_type"`
		Count     int64  `json:"count"`
	}
	if err := r.db.WithContext(ctx).Model(&model.IssueReport{}).
		Select("issue_type, COUNT(*) as count").
		Group("issue_type").
		Find(&typeStats).Error; err != nil {
		return domain.IssueReportStatistics{}, domain.ErrInternal(err)
	}

	stats.ByType.Types = make(map[string]int)
	for _, ts := range typeStats {
		stats.ByType.Types[ts.IssueType] = int(ts.Count)
	}

	// Get creation trends (last 30 days)
	var creationTrends []struct {
		Date  time.Time `json:"date"`
		Count int64     `json:"count"`
	}
	if err := r.db.WithContext(ctx).Model(&model.IssueReport{}).
		Select("DATE(reported_date) as date, COUNT(*) as count").
		Where("reported_date >= ?", time.Now().AddDate(0, 0, -30)).
		Group("DATE(reported_date)").
		Order("date DESC").
		Find(&creationTrends).Error; err != nil {
		return domain.IssueReportStatistics{}, domain.ErrInternal(err)
	}

	stats.CreationTrends = make([]domain.IssueReportCreationTrend, len(creationTrends))
	for i, ct := range creationTrends {
		stats.CreationTrends[i] = domain.IssueReportCreationTrend{
			Date:  ct.Date,
			Count: int(ct.Count),
		}
	}

	// Calculate summary statistics
	stats.Summary.TotalReports = int(totalCount)

	openCount := stats.ByStatus.Open + stats.ByStatus.InProgress
	resolvedCount := stats.ByStatus.Resolved + stats.ByStatus.Closed

	if totalCount > 0 {
		stats.Summary.OpenPercentage = float64(openCount) / float64(totalCount) * 100
		stats.Summary.ResolvedPercentage = float64(resolvedCount) / float64(totalCount) * 100
	}

	// Get critical unresolved count
	var criticalUnresolvedCount int64
	if err := r.db.WithContext(ctx).Model(&model.IssueReport{}).
		Where("priority = ? AND status IN ('Open', 'In Progress')", domain.PriorityCritical).
		Count(&criticalUnresolvedCount).Error; err != nil {
		return domain.IssueReportStatistics{}, domain.ErrInternal(err)
	}
	stats.Summary.CriticalUnresolvedCount = int(criticalUnresolvedCount)

	// Find most common priority and type
	mostCommonPriorityCount := 0
	mostCommonPriority := ""
	if stats.ByPriority.Low > mostCommonPriorityCount {
		mostCommonPriorityCount = stats.ByPriority.Low
		mostCommonPriority = "Low"
	}
	if stats.ByPriority.Medium > mostCommonPriorityCount {
		mostCommonPriorityCount = stats.ByPriority.Medium
		mostCommonPriority = "Medium"
	}
	if stats.ByPriority.High > mostCommonPriorityCount {
		mostCommonPriorityCount = stats.ByPriority.High
		mostCommonPriority = "High"
	}
	if stats.ByPriority.Critical > mostCommonPriorityCount {
		mostCommonPriority = "Critical"
	}
	stats.Summary.MostCommonPriority = mostCommonPriority

	// Find most common type
	mostCommonTypeCount := 0
	mostCommonType := ""
	for issueType, count := range stats.ByType.Types {
		if count > mostCommonTypeCount {
			mostCommonTypeCount = count
			mostCommonType = issueType
		}
	}
	stats.Summary.MostCommonType = mostCommonType

	// Calculate average resolution time
	var avgResolutionDays float64
	if err := r.db.WithContext(ctx).Model(&model.IssueReport{}).
		Select("AVG(EXTRACT(DAY FROM (resolved_date - reported_date))) as avg_days").
		Where("resolved_date IS NOT NULL").
		Row().Scan(&avgResolutionDays); err == nil {
		stats.Summary.AverageResolutionTime = avgResolutionDays
	}

	// Get earliest and latest creation dates
	var earliestDate, latestDate *time.Time
	if err := r.db.WithContext(ctx).Model(&model.IssueReport{}).
		Select("MIN(reported_date) as earliest, MAX(reported_date) as latest").
		Row().Scan(&earliestDate, &latestDate); err != nil {
		return domain.IssueReportStatistics{}, domain.ErrInternal(err)
	}

	if earliestDate != nil {
		stats.Summary.EarliestCreationDate = *earliestDate
	}
	if latestDate != nil {
		stats.Summary.LatestCreationDate = *latestDate
	}

	// Calculate average reports per day
	if earliestDate != nil && latestDate != nil {
		daysDiff := latestDate.Sub(*earliestDate).Hours() / 24
		if daysDiff > 0 {
			stats.Summary.AverageReportsPerDay = float64(totalCount) / daysDiff
		}
	}

	return stats, nil
}
