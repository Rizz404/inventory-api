package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/oklog/ulid/v2"
)

// *==================== Model conversions ====================

func ToModelIssueReport(d *domain.IssueReport) model.IssueReport {
	modelReport := model.IssueReport{
		ReportedDate: d.ReportedDate,
		IssueType:    d.IssueType,
		Priority:     d.Priority,
		Status:       d.Status,
		ResolvedDate: d.ResolvedDate,
	}

	if d.ID != "" {
		if parsedID, err := ulid.Parse(d.ID); err == nil {
			modelReport.ID = model.SQLULID(parsedID)
		}
	}

	if d.AssetID != "" {
		if parsedAssetID, err := ulid.Parse(d.AssetID); err == nil {
			modelReport.AssetID = model.SQLULID(parsedAssetID)
		}
	}

	if d.ReportedBy != "" {
		if parsedReportedBy, err := ulid.Parse(d.ReportedBy); err == nil {
			modelReport.ReportedBy = model.SQLULID(parsedReportedBy)
		}
	}

	if d.ResolvedBy != nil && *d.ResolvedBy != "" {
		if parsedResolvedBy, err := ulid.Parse(*d.ResolvedBy); err == nil {
			modelULID := model.SQLULID(parsedResolvedBy)
			modelReport.ResolvedBy = &modelULID
		}
	}

	return modelReport
}

func ToModelIssueReportForCreate(d *domain.IssueReport) model.IssueReport {
	modelReport := model.IssueReport{
		ReportedDate: d.ReportedDate,
		IssueType:    d.IssueType,
		Priority:     d.Priority,
		Status:       d.Status,
		ResolvedDate: d.ResolvedDate,
	}

	if d.AssetID != "" {
		if parsedAssetID, err := ulid.Parse(d.AssetID); err == nil {
			modelReport.AssetID = model.SQLULID(parsedAssetID)
		}
	}

	if d.ReportedBy != "" {
		if parsedReportedBy, err := ulid.Parse(d.ReportedBy); err == nil {
			modelReport.ReportedBy = model.SQLULID(parsedReportedBy)
		}
	}

	if d.ResolvedBy != nil && *d.ResolvedBy != "" {
		if parsedResolvedBy, err := ulid.Parse(*d.ResolvedBy); err == nil {
			modelULID := model.SQLULID(parsedResolvedBy)
			modelReport.ResolvedBy = &modelULID
		}
	}

	return modelReport
}

func ToModelIssueReportTranslation(d *domain.IssueReportTranslation) model.IssueReportTranslation {
	modelTranslation := model.IssueReportTranslation{
		LangCode:        d.LangCode,
		Title:           d.Title,
		Description:     d.Description,
		ResolutionNotes: d.ResolutionNotes,
	}

	if d.ID != "" {
		if parsedID, err := ulid.Parse(d.ID); err == nil {
			modelTranslation.ID = model.SQLULID(parsedID)
		}
	}

	if d.ReportID != "" {
		if parsedReportID, err := ulid.Parse(d.ReportID); err == nil {
			modelTranslation.ReportID = model.SQLULID(parsedReportID)
		}
	}

	return modelTranslation
}

func ToModelIssueReportTranslationForCreate(reportID string, d *domain.IssueReportTranslation) model.IssueReportTranslation {
	modelTranslation := model.IssueReportTranslation{
		LangCode:        d.LangCode,
		Title:           d.Title,
		Description:     d.Description,
		ResolutionNotes: d.ResolutionNotes,
	}

	if reportID != "" {
		if parsedReportID, err := ulid.Parse(reportID); err == nil {
			modelTranslation.ReportID = model.SQLULID(parsedReportID)
		}
	}

	return modelTranslation
}

// *==================== Entity conversions ====================
func ToDomainIssueReport(m *model.IssueReport) domain.IssueReport {
	domainReport := domain.IssueReport{
		ID:           m.ID.String(),
		AssetID:      m.AssetID.String(),
		ReportedBy:   m.ReportedBy.String(),
		ReportedDate: m.ReportedDate,
		IssueType:    m.IssueType,
		Priority:     m.Priority,
		Status:       m.Status,
		ResolvedDate: m.ResolvedDate,
	}

	if m.ResolvedBy != nil && !m.ResolvedBy.IsZero() {
		resolvedByStr := m.ResolvedBy.String()
		domainReport.ResolvedBy = &resolvedByStr
	}

	if len(m.Translations) > 0 {
		domainReport.Translations = make([]domain.IssueReportTranslation, len(m.Translations))
		for i, translation := range m.Translations {
			domainReport.Translations[i] = ToDomainIssueReportTranslation(&translation)
		}
	}

	return domainReport
}

func ToDomainIssueReportTranslation(m *model.IssueReportTranslation) domain.IssueReportTranslation {
	return domain.IssueReportTranslation{
		ID:              m.ID.String(),
		ReportID:        m.ReportID.String(),
		LangCode:        m.LangCode,
		Title:           m.Title,
		Description:     m.Description,
		ResolutionNotes: m.ResolutionNotes,
	}
}

func ToDomainIssueReports(models []model.IssueReport) []domain.IssueReport {
	reports := make([]domain.IssueReport, len(models))
	for i, m := range models {
		reports[i] = ToDomainIssueReport(&m)
	}
	return reports
}

// *==================== Entity Response conversions ====================
func IssueReportToResponse(d *domain.IssueReport, langCode string) domain.IssueReportResponse {
	response := domain.IssueReportResponse{
		ID:           d.ID,
		AssetID:      d.AssetID,
		ReportedByID: d.ReportedBy,
		ReportedDate: d.ReportedDate,
		IssueType:    d.IssueType,
		Priority:     d.Priority,
		Status:       d.Status,
		ResolvedDate: d.ResolvedDate,
		ResolvedByID: d.ResolvedBy,
		CreatedAt:    d.ReportedDate, // Use ReportedDate as CreatedAt since domain doesn't have CreatedAt
		UpdatedAt:    d.ReportedDate, // Use ReportedDate as UpdatedAt since domain doesn't have UpdatedAt
		Translations: make([]domain.IssueReportTranslationResponse, len(d.Translations)),
	}

	// Populate translations
	for i, translation := range d.Translations {
		response.Translations[i] = domain.IssueReportTranslationResponse{
			LangCode:        translation.LangCode,
			Title:           translation.Title,
			Description:     translation.Description,
			ResolutionNotes: translation.ResolutionNotes,
		}
	}

	// Find translation for the requested language
	for _, translation := range d.Translations {
		if translation.LangCode == langCode {
			response.Title = translation.Title
			response.Description = translation.Description
			response.ResolutionNotes = translation.ResolutionNotes
			break
		}
	}

	// If no translation found for requested language, use first available
	if response.Title == "" && len(d.Translations) > 0 {
		response.Title = d.Translations[0].Title
		response.Description = d.Translations[0].Description
		response.ResolutionNotes = d.Translations[0].ResolutionNotes
	}

	return response
}

func IssueReportsToResponses(reports []domain.IssueReport, langCode string) []domain.IssueReportResponse {
	responses := make([]domain.IssueReportResponse, len(reports))
	for i, report := range reports {
		responses[i] = IssueReportToResponse(&report, langCode)
	}
	return responses
}

func IssueReportToListResponse(d *domain.IssueReport, langCode string) domain.IssueReportListResponse {
	response := domain.IssueReportListResponse{
		ID:           d.ID,
		AssetID:      d.AssetID,
		ReportedByID: d.ReportedBy,
		ReportedDate: d.ReportedDate,
		IssueType:    d.IssueType,
		Priority:     d.Priority,
		Status:       d.Status,
		ResolvedDate: d.ResolvedDate,
		ResolvedByID: d.ResolvedBy,
		CreatedAt:    d.ReportedDate, // Use ReportedDate as CreatedAt since domain doesn't have CreatedAt
		UpdatedAt:    d.ReportedDate, // Use ReportedDate as UpdatedAt since domain doesn't have UpdatedAt
	}

	// Find translation for the requested language
	for _, translation := range d.Translations {
		if translation.LangCode == langCode {
			response.Title = translation.Title
			response.Description = translation.Description
			response.ResolutionNotes = translation.ResolutionNotes
			break
		}
	}

	// If no translation found for requested language, use first available
	if response.Title == "" && len(d.Translations) > 0 {
		response.Title = d.Translations[0].Title
		response.Description = d.Translations[0].Description
		response.ResolutionNotes = d.Translations[0].ResolutionNotes
	}

	return response
}

func IssueReportsToListResponses(reports []domain.IssueReport, langCode string) []domain.IssueReportListResponse {
	responses := make([]domain.IssueReportListResponse, len(reports))
	for i, report := range reports {
		responses[i] = IssueReportToListResponse(&report, langCode)
	}
	return responses
}

// ToModelIssueReportUpdateMap converts UpdateIssueReportPayload to map for database updates
func ToModelIssueReportUpdateMap(payload *domain.UpdateIssueReportPayload) map[string]interface{} {
	updates := make(map[string]interface{})

	if payload.Priority != nil {
		updates["priority"] = *payload.Priority
	}

	if payload.Status != nil {
		updates["status"] = *payload.Status
	}

	if payload.ResolutionNotes != nil {
		updates["resolution_notes"] = *payload.ResolutionNotes
	}

	return updates
}

// *==================== Statistics conversions ====================
// IssueReportStatisticsToResponse converts IssueReportStatistics to IssueReportStatisticsResponse
func IssueReportStatisticsToResponse(stats *domain.IssueReportStatistics) domain.IssueReportStatisticsResponse {
	response := domain.IssueReportStatisticsResponse{
		Total: domain.IssueReportCountStatisticsResponse{
			Count: stats.Total.Count,
		},
		ByPriority: domain.IssueReportPriorityStatisticsResponse{
			Low:      stats.ByPriority.Low,
			Medium:   stats.ByPriority.Medium,
			High:     stats.ByPriority.High,
			Critical: stats.ByPriority.Critical,
		},
		ByStatus: domain.IssueReportStatusStatisticsResponse{
			Open:       stats.ByStatus.Open,
			InProgress: stats.ByStatus.InProgress,
			Resolved:   stats.ByStatus.Resolved,
			Closed:     stats.ByStatus.Closed,
		},
		ByType: domain.IssueReportTypeStatisticsResponse{
			Types: stats.ByType.Types,
		},
		Summary: domain.IssueReportSummaryStatisticsResponse{
			TotalReports:            stats.Summary.TotalReports,
			OpenPercentage:          stats.Summary.OpenPercentage,
			ResolvedPercentage:      stats.Summary.ResolvedPercentage,
			AverageResolutionTime:   stats.Summary.AverageResolutionTime,
			MostCommonPriority:      stats.Summary.MostCommonPriority,
			MostCommonType:          stats.Summary.MostCommonType,
			CriticalUnresolvedCount: stats.Summary.CriticalUnresolvedCount,
			AverageReportsPerDay:    stats.Summary.AverageReportsPerDay,
			LatestCreationDate:      stats.Summary.LatestCreationDate,
			EarliestCreationDate:    stats.Summary.EarliestCreationDate,
		},
	}

	// Convert creation trends
	response.CreationTrends = make([]domain.IssueReportCreationTrendResponse, len(stats.CreationTrends))
	for i, trend := range stats.CreationTrends {
		response.CreationTrends[i] = domain.IssueReportCreationTrendResponse{
			Date:  trend.Date,
			Count: trend.Count,
		}
	}

	return response
}

func MapIssueReportSortFieldToColumn(field domain.IssueReportSortField) string {
	columnMap := map[domain.IssueReportSortField]string{
		domain.IssueReportSortByReportedDate: "ir.reported_date",
		domain.IssueReportSortByPriority:     "ir.priority",
		domain.IssueReportSortByStatus:       "ir.status",
		domain.IssueReportSortByCreatedAt:    "ir.created_at",
		domain.IssueReportSortByUpdatedAt:    "ir.updated_at",
	}

	if column, exists := columnMap[field]; exists {
		return column
	}
	return "ir.reported_date"
}
