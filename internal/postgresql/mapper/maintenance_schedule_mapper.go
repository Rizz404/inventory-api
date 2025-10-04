package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/oklog/ulid/v2"
)

// *==================== Model conversions ====================

func ToModelMaintenanceSchedule(d *domain.MaintenanceSchedule) model.MaintenanceSchedule {
	modelSchedule := model.MaintenanceSchedule{
		MaintenanceType: d.MaintenanceType,
		ScheduledDate:   d.ScheduledDate,
		FrequencyMonths: d.FrequencyMonths,
		Status:          d.Status,
	}

	if d.ID != "" {
		if parsedID, err := ulid.Parse(d.ID); err == nil {
			modelSchedule.ID = model.SQLULID(parsedID)
		}
	}

	if d.AssetID != "" {
		if parsedAssetID, err := ulid.Parse(d.AssetID); err == nil {
			modelSchedule.AssetID = model.SQLULID(parsedAssetID)
		}
	}

	if d.CreatedBy != "" {
		if parsedCreatedBy, err := ulid.Parse(d.CreatedBy); err == nil {
			modelSchedule.CreatedBy = model.SQLULID(parsedCreatedBy)
		}
	}

	return modelSchedule
}

func ToModelMaintenanceScheduleForCreate(d *domain.MaintenanceSchedule) model.MaintenanceSchedule {
	modelSchedule := model.MaintenanceSchedule{
		MaintenanceType: d.MaintenanceType,
		ScheduledDate:   d.ScheduledDate,
		FrequencyMonths: d.FrequencyMonths,
		Status:          d.Status,
	}

	if d.AssetID != "" {
		if parsedAssetID, err := ulid.Parse(d.AssetID); err == nil {
			modelSchedule.AssetID = model.SQLULID(parsedAssetID)
		}
	}

	if d.CreatedBy != "" {
		if parsedCreatedBy, err := ulid.Parse(d.CreatedBy); err == nil {
			modelSchedule.CreatedBy = model.SQLULID(parsedCreatedBy)
		}
	}

	return modelSchedule
}

func ToModelMaintenanceScheduleTranslation(d *domain.MaintenanceScheduleTranslation) model.MaintenanceScheduleTranslation {
	modelTranslation := model.MaintenanceScheduleTranslation{
		LangCode:    d.LangCode,
		Title:       d.Title,
		Description: d.Description,
	}

	if d.ID != "" {
		if parsedID, err := ulid.Parse(d.ID); err == nil {
			modelTranslation.ID = model.SQLULID(parsedID)
		}
	}

	if d.ScheduleID != "" {
		if parsedScheduleID, err := ulid.Parse(d.ScheduleID); err == nil {
			modelTranslation.ScheduleID = model.SQLULID(parsedScheduleID)
		}
	}

	return modelTranslation
}

func ToModelMaintenanceScheduleTranslationForCreate(scheduleID string, d *domain.MaintenanceScheduleTranslation) model.MaintenanceScheduleTranslation {
	modelTranslation := model.MaintenanceScheduleTranslation{
		LangCode:    d.LangCode,
		Title:       d.Title,
		Description: d.Description,
	}

	if scheduleID != "" {
		if parsedScheduleID, err := ulid.Parse(scheduleID); err == nil {
			modelTranslation.ScheduleID = model.SQLULID(parsedScheduleID)
		}
	}

	return modelTranslation
}

func ToDomainMaintenanceSchedule(m *model.MaintenanceSchedule) domain.MaintenanceSchedule {
	domainSchedule := domain.MaintenanceSchedule{
		ID:              m.ID.String(),
		AssetID:         m.AssetID.String(),
		MaintenanceType: m.MaintenanceType,
		ScheduledDate:   m.ScheduledDate,
		FrequencyMonths: m.FrequencyMonths,
		Status:          m.Status,
		CreatedBy:       m.CreatedBy.String(),
		CreatedAt:       m.CreatedAt,
	}

	if len(m.Translations) > 0 {
		domainSchedule.Translations = make([]domain.MaintenanceScheduleTranslation, len(m.Translations))
		for i, translation := range m.Translations {
			domainSchedule.Translations[i] = ToDomainMaintenanceScheduleTranslation(&translation)
		}
	}

	return domainSchedule
}

func ToDomainMaintenanceScheduleTranslation(m *model.MaintenanceScheduleTranslation) domain.MaintenanceScheduleTranslation {
	return domain.MaintenanceScheduleTranslation{
		ID:          m.ID.String(),
		ScheduleID:  m.ScheduleID.String(),
		LangCode:    m.LangCode,
		Title:       m.Title,
		Description: m.Description,
	}
}

func ToDomainMaintenanceSchedules(models []model.MaintenanceSchedule) []domain.MaintenanceSchedule {
	schedules := make([]domain.MaintenanceSchedule, len(models))
	for i, m := range models {
		schedules[i] = ToDomainMaintenanceSchedule(&m)
	}
	return schedules
}

// *==================== Entity Response conversions ====================
func MaintenanceScheduleToResponse(d *domain.MaintenanceSchedule, langCode string) domain.MaintenanceScheduleResponse {
	response := domain.MaintenanceScheduleResponse{
		ID:              d.ID,
		AssetID:         d.AssetID,
		MaintenanceType: d.MaintenanceType,
		ScheduledDate:   d.ScheduledDate,
		FrequencyMonths: d.FrequencyMonths,
		Status:          d.Status,
		CreatedByID:     d.CreatedBy,
		CreatedAt:       d.CreatedAt,
		Translations:    make([]domain.MaintenanceScheduleTranslationResponse, len(d.Translations)),
	}

	// Populate translations
	for i, translation := range d.Translations {
		response.Translations[i] = domain.MaintenanceScheduleTranslationResponse{
			LangCode:    translation.LangCode,
			Title:       translation.Title,
			Description: translation.Description,
		}
	}

	// Find translation for the requested language
	for _, translation := range d.Translations {
		if translation.LangCode == langCode {
			response.Title = translation.Title
			response.Description = translation.Description
			break
		}
	}

	// If no translation found for requested language, use first available
	if response.Title == "" && len(d.Translations) > 0 {
		response.Title = d.Translations[0].Title
		response.Description = d.Translations[0].Description
	}

	return response
}

func MaintenanceScheduleToListResponse(d *domain.MaintenanceSchedule, langCode string) domain.MaintenanceScheduleListResponse {
	response := domain.MaintenanceScheduleListResponse{
		ID:              d.ID,
		AssetID:         d.AssetID,
		MaintenanceType: d.MaintenanceType,
		ScheduledDate:   d.ScheduledDate,
		FrequencyMonths: d.FrequencyMonths,
		Status:          d.Status,
		CreatedByID:     d.CreatedBy,
		CreatedAt:       d.CreatedAt,
	}

	// Find translation for the requested language
	for _, translation := range d.Translations {
		if translation.LangCode == langCode {
			response.Title = translation.Title
			response.Description = translation.Description
			break
		}
	}

	// If no translation found for requested language, use first available
	if response.Title == "" && len(d.Translations) > 0 {
		response.Title = d.Translations[0].Title
		response.Description = d.Translations[0].Description
	}

	return response
}

func MaintenanceSchedulesToResponses(schedules []domain.MaintenanceSchedule, langCode string) []domain.MaintenanceScheduleResponse {
	responses := make([]domain.MaintenanceScheduleResponse, len(schedules))
	for i, schedule := range schedules {
		responses[i] = MaintenanceScheduleToResponse(&schedule, langCode)
	}
	return responses
}

func MaintenanceSchedulesToListResponses(schedules []domain.MaintenanceSchedule, langCode string) []domain.MaintenanceScheduleListResponse {
	responses := make([]domain.MaintenanceScheduleListResponse, len(schedules))
	for i, schedule := range schedules {
		responses[i] = MaintenanceScheduleToListResponse(&schedule, langCode)
	}
	return responses
}

// *==================== Statistics conversions ====================

func MaintenanceScheduleStatisticsToResponse(stats *domain.MaintenanceScheduleStatistics) domain.MaintenanceScheduleStatisticsResponse {
	resp := domain.MaintenanceScheduleStatisticsResponse{
		Total: domain.MaintenanceScheduleCountStatisticsResponse{Count: stats.Total.Count},
		ByType: domain.MaintenanceTypeStatisticsResponse{
			Preventive: stats.ByType.Preventive,
			Corrective: stats.ByType.Corrective,
		},
		ByStatus: domain.MaintenanceScheduleStatusStatisticsResponse{
			Scheduled: stats.ByStatus.Scheduled,
			Completed: stats.ByStatus.Completed,
			Cancelled: stats.ByStatus.Cancelled,
		},
		Summary: domain.MaintenanceScheduleSummaryStatisticsResponse{
			TotalSchedules:                    stats.Summary.TotalSchedules,
			ScheduledMaintenancePercentage:    stats.Summary.ScheduledMaintenancePercentage,
			CompletedMaintenancePercentage:    stats.Summary.CompletedMaintenancePercentage,
			CancelledMaintenancePercentage:    stats.Summary.CancelledMaintenancePercentage,
			PreventiveMaintenancePercentage:   stats.Summary.PreventiveMaintenancePercentage,
			CorrectiveMaintenancePercentage:   stats.Summary.CorrectiveMaintenancePercentage,
			AverageScheduleFrequency:          stats.Summary.AverageScheduleFrequency,
			UpcomingMaintenanceCount:          stats.Summary.UpcomingMaintenanceCount,
			OverdueMaintenanceCount:           stats.Summary.OverdueMaintenanceCount,
			AssetsWithScheduledMaintenance:    stats.Summary.AssetsWithScheduledMaintenance,
			AssetsWithoutScheduledMaintenance: stats.Summary.AssetsWithoutScheduledMaintenance,
			AverageSchedulesPerDay:            stats.Summary.AverageSchedulesPerDay,
			LatestScheduleDate:                stats.Summary.LatestScheduleDate,
			EarliestScheduleDate:              stats.Summary.EarliestScheduleDate,
			TotalUniqueCreators:               stats.Summary.TotalUniqueCreators,
		},
	}

	// ByAsset schedules
	resp.ByAsset = make([]domain.AssetMaintenanceScheduleStatisticsResponse, len(stats.ByAsset))
	for i, a := range stats.ByAsset {
		resp.ByAsset[i] = domain.AssetMaintenanceScheduleStatisticsResponse{
			AssetID:         a.AssetID,
			AssetName:       a.AssetName,
			AssetTag:        a.AssetTag,
			ScheduleCount:   a.ScheduleCount,
			NextMaintenance: a.NextMaintenance,
		}
	}

	// ByCreator
	resp.ByCreator = make([]domain.UserMaintenanceScheduleStatisticsResponse, len(stats.ByCreator))
	for i, u := range stats.ByCreator {
		resp.ByCreator[i] = domain.UserMaintenanceScheduleStatisticsResponse{
			UserID:    u.UserID,
			UserName:  u.UserName,
			UserEmail: u.UserEmail,
			Count:     u.Count,
		}
	}

	// Upcoming & Overdue schedules
	resp.UpcomingSchedule = make([]domain.UpcomingMaintenanceScheduleResponse, len(stats.UpcomingSchedule))
	for i, up := range stats.UpcomingSchedule {
		resp.UpcomingSchedule[i] = domain.UpcomingMaintenanceScheduleResponse{
			ID:              up.ID,
			AssetID:         up.AssetID,
			AssetName:       up.AssetName,
			AssetTag:        up.AssetTag,
			MaintenanceType: up.MaintenanceType,
			ScheduledDate:   up.ScheduledDate,
			DaysUntilDue:    up.DaysUntilDue,
			Title:           up.Title,
			Description:     up.Description,
		}
	}

	resp.OverdueSchedule = make([]domain.OverdueMaintenanceScheduleResponse, len(stats.OverdueSchedule))
	for i, od := range stats.OverdueSchedule {
		resp.OverdueSchedule[i] = domain.OverdueMaintenanceScheduleResponse{
			ID:              od.ID,
			AssetID:         od.AssetID,
			AssetName:       od.AssetName,
			AssetTag:        od.AssetTag,
			MaintenanceType: od.MaintenanceType,
			ScheduledDate:   od.ScheduledDate,
			DaysOverdue:     od.DaysOverdue,
			Title:           od.Title,
			Description:     od.Description,
		}
	}

	// Frequency trends
	resp.FrequencyTrends = make([]domain.MaintenanceFrequencyTrendResponse, len(stats.FrequencyTrends))
	for i, ft := range stats.FrequencyTrends {
		resp.FrequencyTrends[i] = domain.MaintenanceFrequencyTrendResponse{FrequencyMonths: ft.FrequencyMonths, Count: ft.Count}
	}

	return resp
}

func MapMaintenanceScheduleSortFieldToColumn(field domain.MaintenanceScheduleSortField) string {
	columnMap := map[domain.MaintenanceScheduleSortField]string{
		domain.MaintenanceScheduleSortByScheduledDate:   "ms.scheduled_date",
		domain.MaintenanceScheduleSortByMaintenanceType: "ms.maintenance_type",
		domain.MaintenanceScheduleSortByStatus:          "ms.status",
		domain.MaintenanceScheduleSortByCreatedAt:       "ms.created_at",
		domain.MaintenanceScheduleSortByUpdatedAt:       "ms.updated_at",
	}

	if column, exists := columnMap[field]; exists {
		return column
	}
	return "ms.scheduled_date"
}
