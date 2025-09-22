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

type IssueReportRepository struct {
	db *gorm.DB
}

type IssueReportFilterOptions struct {
	AssetID    *string               `json:"assetId,omitempty"`
	ReportedBy *string               `json:"reportedBy,omitempty"`
	ResolvedBy *string               `json:"resolvedBy,omitempty"`
	IssueType  *string               `json:"issueType,omitempty"`
	Priority   *domain.IssuePriority `json:"priority,omitempty"`
	Status     *domain.IssueStatus   `json:"status,omitempty"`
	IsResolved *bool                 `json:"isResolved,omitempty"`
	DateFrom   *time.Time            `json:"dateFrom,omitempty"`
	DateTo     *time.Time            `json:"dateTo,omitempty"`
}

func NewIssueReportRepository(db *gorm.DB) *IssueReportRepository {
	return &IssueReportRepository{
		db: db,
	}
}

func (r *IssueReportRepository) applyIssueReportFilters(db *gorm.DB, filters any) *gorm.DB {
	f, ok := filters.(*IssueReportFilterOptions)
	if !ok || f == nil {
		return db
	}

	if f.AssetID != nil {
		db = db.Where("ir.asset_id = ?", f.AssetID)
	}
	if f.ReportedBy != nil {
		db = db.Where("ir.reported_by = ?", f.ReportedBy)
	}
	if f.ResolvedBy != nil {
		db = db.Where("ir.resolved_by = ?", f.ResolvedBy)
	}
	if f.IssueType != nil {
		db = db.Where("ir.issue_type = ?", f.IssueType)
	}
	if f.Priority != nil {
		db = db.Where("ir.priority = ?", f.Priority)
	}
	if f.Status != nil {
		db = db.Where("ir.status = ?", f.Status)
	}
	if f.IsResolved != nil {
		if *f.IsResolved {
			db = db.Where("ir.status IN ('Resolved', 'Closed')")
		} else {
			db = db.Where("ir.status IN ('Open', 'In Progress')")
		}
	}
	if f.DateFrom != nil {
		db = db.Where("ir.reported_date >= ?", f.DateFrom)
	}
	if f.DateTo != nil {
		db = db.Where("ir.reported_date <= ?", f.DateTo)
	}
	return db
}

func (r *IssueReportRepository) applyIssueReportSorts(db *gorm.DB, sort *query.SortOptions) *gorm.DB {
	if sort == nil || sort.Field == "" {
		return db.Order("ir.reported_date DESC")
	}

	var orderClause string
	switch strings.ToLower(sort.Field) {
	case "reported_date", "resolved_date", "issue_type", "priority", "status":
		orderClause = fmt.Sprintf("ir.%s", sort.Field)
	case "title":
		orderClause = "irt.title"
	case "description":
		orderClause = "irt.description"
	default:
		return db.Order("ir.reported_date DESC")
	}

	order := "DESC"
	if strings.ToLower(sort.Order) == "asc" {
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

	// Fetch created issue report with translations and relations
	return r.GetIssueReportById(ctx, modelIssueReport.ID.String())
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

	// Update issue report basic info
	updates := mapper.ToModelIssueReportUpdateMap(payload)
	if len(updates) > 0 {
		if err := tx.Model(&model.IssueReport{}).Where("id = ?", issueReportId).Updates(updates).Error; err != nil {
			tx.Rollback()
			return domain.IssueReport{}, domain.ErrInternal(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return domain.IssueReport{}, domain.ErrInternal(err)
	}

	// Fetch updated issue report with translations and relations
	return r.GetIssueReportById(ctx, issueReportId)
}

func (r *IssueReportRepository) ResolveIssueReport(ctx context.Context, issueReportId string, resolvedBy string, resolutionNotes string) (domain.IssueReport, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.IssueReport{}, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()
	updates := map[string]interface{}{
		"status":        domain.IssueStatusResolved,
		"resolved_date": &now,
		"resolved_by":   resolvedBy,
	}

	if err := tx.Model(&model.IssueReport{}).Where("id = ?", issueReportId).Updates(updates).Error; err != nil {
		tx.Rollback()
		return domain.IssueReport{}, domain.ErrInternal(err)
	}

	// Update resolution notes in all translations
	if err := tx.Model(&model.IssueReportTranslation{}).
		Where("report_id = ?", issueReportId).
		Update("resolution_notes", resolutionNotes).Error; err != nil {
		tx.Rollback()
		return domain.IssueReport{}, domain.ErrInternal(err)
	}

	if err := tx.Commit().Error; err != nil {
		return domain.IssueReport{}, domain.ErrInternal(err)
	}

	return r.GetIssueReportById(ctx, issueReportId)
}

func (r *IssueReportRepository) ReopenIssueReport(ctx context.Context, issueReportId string) (domain.IssueReport, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.IssueReport{}, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	updates := map[string]interface{}{
		"status":        domain.IssueStatusOpen,
		"resolved_date": nil,
		"resolved_by":   nil,
	}

	if err := tx.Model(&model.IssueReport{}).Where("id = ?", issueReportId).Updates(updates).Error; err != nil {
		tx.Rollback()
		return domain.IssueReport{}, domain.ErrInternal(err)
	}

	// Clear resolution notes in all translations
	if err := tx.Model(&model.IssueReportTranslation{}).
		Where("report_id = ?", issueReportId).
		Update("resolution_notes", nil).Error; err != nil {
		tx.Rollback()
		return domain.IssueReport{}, domain.ErrInternal(err)
	}

	if err := tx.Commit().Error; err != nil {
		return domain.IssueReport{}, domain.ErrInternal(err)
	}

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
func (r *IssueReportRepository) GetIssueReportsPaginated(ctx context.Context, params query.Params, langCode string) ([]domain.IssueReportListItem, error) {
	var issueReports []model.IssueReport
	db := r.db.WithContext(ctx).
		Table("issue_reports ir").
		Preload("Translations").
		Preload("Asset").
		Preload("ReportedByUser").
		Preload("ResolvedByUser")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN issue_report_translations irt ON ir.id = irt.report_id").
			Where("ir.issue_type ILIKE ? OR irt.title ILIKE ? OR irt.description ILIKE ?", searchPattern, searchPattern, searchPattern).
			Distinct("ir.id")
	}

	// Set pagination cursor to empty for offset-based pagination
	params.Pagination.Cursor = ""
	db = query.Apply(db, params, r.applyIssueReportFilters, r.applyIssueReportSorts)

	if err := db.Find(&issueReports).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain issue reports first, then to list items
	domainIssueReports := mapper.ToDomainIssueReports(issueReports)
	return mapper.IssueReportsToListItems(domainIssueReports, langCode), nil
}

func (r *IssueReportRepository) GetIssueReportsCursor(ctx context.Context, params query.Params, langCode string) ([]domain.IssueReportListItem, error) {
	var issueReports []model.IssueReport
	db := r.db.WithContext(ctx).
		Table("issue_reports ir").
		Preload("Translations").
		Preload("Asset").
		Preload("ReportedByUser").
		Preload("ResolvedByUser")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN issue_report_translations irt ON ir.id = irt.report_id").
			Where("ir.issue_type ILIKE ? OR irt.title ILIKE ? OR irt.description ILIKE ?", searchPattern, searchPattern, searchPattern).
			Distinct("ir.id")
	}

	// Set offset to 0 for cursor-based pagination
	params.Pagination.Offset = 0
	db = query.Apply(db, params, r.applyIssueReportFilters, r.applyIssueReportSorts)

	if err := db.Find(&issueReports).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain issue reports
	domainIssueReports := mapper.ToDomainIssueReports(issueReports)
	return mapper.IssueReportsToListItems(domainIssueReports, langCode), nil
}

func (r *IssueReportRepository) GetIssueReportById(ctx context.Context, issueReportId string) (domain.IssueReport, error) {
	var issueReport model.IssueReport

	err := r.db.WithContext(ctx).
		Table("issue_reports ir").
		Preload("Translations").
		Preload("Asset").
		Preload("ReportedByUser").
		Preload("ResolvedByUser").
		First(&issueReport, "id = ?", issueReportId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.IssueReport{}, domain.ErrNotFound("issue report not found")
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

func (r *IssueReportRepository) CountIssueReports(ctx context.Context, params query.Params) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Table("issue_reports ir")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Joins("LEFT JOIN issue_report_translations irt ON ir.id = irt.report_id").
			Where("ir.issue_type ILIKE ? OR irt.title ILIKE ? OR irt.description ILIKE ?", searchPattern, searchPattern, searchPattern).
			Distinct("ir.id")
	}

	db = query.Apply(db, query.Params{Filters: params.Filters}, r.applyIssueReportFilters, nil)

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
		Date  string `json:"date"`
		Count int64  `json:"count"`
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
	var earliestDate, latestDate time.Time
	if err := r.db.WithContext(ctx).Model(&model.IssueReport{}).
		Select("MIN(reported_date) as earliest, MAX(reported_date) as latest").
		Row().Scan(&earliestDate, &latestDate); err != nil {
		return domain.IssueReportStatistics{}, domain.ErrInternal(err)
	}

	stats.Summary.EarliestCreationDate = earliestDate.Format("2006-01-02")
	stats.Summary.LatestCreationDate = latestDate.Format("2006-01-02")

	// Calculate average reports per day
	if !earliestDate.IsZero() && !latestDate.IsZero() {
		daysDiff := latestDate.Sub(earliestDate).Hours() / 24
		if daysDiff > 0 {
			stats.Summary.AverageReportsPerDay = float64(totalCount) / daysDiff
		}
	}

	return stats, nil
}
