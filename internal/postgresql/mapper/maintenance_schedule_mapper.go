package mapper

import (
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/oklog/ulid/v2"
)

// *==================== Model conversions ====================

func ToModelMaintenanceSchedule(d *domain.MaintenanceSchedule) model.MaintenanceSchedule {
	modelSchedule := model.MaintenanceSchedule{
		MaintenanceType:   d.MaintenanceType,
		IsRecurring:       d.IsRecurring,
		IntervalValue:     d.IntervalValue,
		IntervalUnit:      d.IntervalUnit,
		ScheduledTime:     d.ScheduledTime,
		NextScheduledDate: d.NextScheduledDate,
		LastExecutedDate:  d.LastExecutedDate,
		State:             d.State,
		AutoComplete:      d.AutoComplete,
		EstimatedCost:     d.EstimatedCost,
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
		MaintenanceType:   d.MaintenanceType,
		IsRecurring:       d.IsRecurring,
		IntervalValue:     d.IntervalValue,
		IntervalUnit:      d.IntervalUnit,
		ScheduledTime:     d.ScheduledTime,
		NextScheduledDate: d.NextScheduledDate,
		LastExecutedDate:  d.LastExecutedDate,
		State:             d.State,
		AutoComplete:      d.AutoComplete,
		EstimatedCost:     d.EstimatedCost,
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
		ID:                m.ID.String(),
		AssetID:           m.AssetID.String(),
		MaintenanceType:   m.MaintenanceType,
		IsRecurring:       m.IsRecurring,
		IntervalValue:     m.IntervalValue,
		IntervalUnit:      m.IntervalUnit,
		ScheduledTime:     m.ScheduledTime,
		NextScheduledDate: m.NextScheduledDate,
		LastExecutedDate:  m.LastExecutedDate,
		State:             m.State,
		AutoComplete:      m.AutoComplete,
		EstimatedCost:     m.EstimatedCost,
		CreatedBy:         m.CreatedBy.String(),
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}

	if len(m.Translations) > 0 {
		domainSchedule.Translations = make([]domain.MaintenanceScheduleTranslation, len(m.Translations))
		for i, translation := range m.Translations {
			domainSchedule.Translations[i] = ToDomainMaintenanceScheduleTranslation(&translation)
		}
	}

	// Populate related entities if preloaded
	if !m.Asset.ID.IsZero() {
		asset := ToDomainAsset(&m.Asset)
		domainSchedule.Asset = &asset
	}

	if !m.CreatedByUser.ID.IsZero() {
		user := ToDomainUser(&m.CreatedByUser)
		domainSchedule.CreatedByUser = &user
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
	if len(models) == 0 {
		return []domain.MaintenanceSchedule{}
	}
	schedules := make([]domain.MaintenanceSchedule, len(models))
	for i, m := range models {
		schedules[i] = ToDomainMaintenanceSchedule(&m)
	}
	return schedules
}

// *==================== Entity Response conversions ====================
func MaintenanceScheduleToResponse(d *domain.MaintenanceSchedule, langCode string) domain.MaintenanceScheduleResponse {
	response := domain.MaintenanceScheduleResponse{
		ID:                d.ID,
		AssetID:           d.AssetID,
		MaintenanceType:   d.MaintenanceType,
		IsRecurring:       d.IsRecurring,
		IntervalValue:     d.IntervalValue,
		IntervalUnit:      d.IntervalUnit,
		ScheduledTime:     d.ScheduledTime,
		NextScheduledDate: d.NextScheduledDate,
		LastExecutedDate:  d.LastExecutedDate,
		State:             d.State,
		AutoComplete:      d.AutoComplete,
		EstimatedCost:     domain.NewNullableDecimal2(d.EstimatedCost),
		CreatedByID:       d.CreatedBy,
		CreatedAt:         d.CreatedAt,
		UpdatedAt:         d.UpdatedAt,
		Translations:      make([]domain.MaintenanceScheduleTranslationResponse, len(d.Translations)),
	}

	// Populate Asset if available
	if d.Asset != nil {
		response.Asset = AssetToResponse(d.Asset, langCode)
	}

	// Populate CreatedBy User if available
	if d.CreatedByUser != nil {
		response.CreatedBy = UserToResponse(d.CreatedByUser)
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
		ID:                d.ID,
		AssetID:           d.AssetID,
		MaintenanceType:   d.MaintenanceType,
		IsRecurring:       d.IsRecurring,
		IntervalValue:     d.IntervalValue,
		IntervalUnit:      d.IntervalUnit,
		ScheduledTime:     d.ScheduledTime,
		NextScheduledDate: d.NextScheduledDate,
		LastExecutedDate:  d.LastExecutedDate,
		State:             d.State,
		AutoComplete:      d.AutoComplete,
		EstimatedCost:     domain.NewNullableDecimal2(d.EstimatedCost),
		CreatedByID:       d.CreatedBy,
		CreatedAt:         d.CreatedAt,
		UpdatedAt:         d.UpdatedAt,
	}

	// Populate Asset if available
	if d.Asset != nil {
		response.Asset = AssetToResponse(d.Asset, langCode)
	}

	// Populate CreatedBy User if available
	if d.CreatedByUser != nil {
		response.CreatedBy = UserToResponse(d.CreatedByUser)
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
	if len(schedules) == 0 {
		return []domain.MaintenanceScheduleResponse{}
	}
	responses := make([]domain.MaintenanceScheduleResponse, len(schedules))
	for i, schedule := range schedules {
		responses[i] = MaintenanceScheduleToResponse(&schedule, langCode)
	}
	return responses
}

func MaintenanceSchedulesToListResponses(schedules []domain.MaintenanceSchedule, langCode string) []domain.MaintenanceScheduleListResponse {
	if len(schedules) == 0 {
		return []domain.MaintenanceScheduleListResponse{}
	}
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
			Preventive:  stats.ByType.Preventive,
			Corrective:  stats.ByType.Corrective,
			Inspection:  stats.ByType.Inspection,
			Calibration: stats.ByType.Calibration,
		},
		ByStatus: domain.MaintenanceScheduleStatusStatisticsResponse{
			Active:    stats.ByStatus.Active,
			Paused:    stats.ByStatus.Paused,
			Stopped:   stats.ByStatus.Stopped,
			Completed: stats.ByStatus.Completed,
		},
		Summary: domain.MaintenanceScheduleSummaryStatisticsResponse{
			TotalSchedules:                    stats.Summary.TotalSchedules,
			ActiveMaintenancePercentage:       domain.NewDecimal2(stats.Summary.ActiveMaintenancePercentage),
			CompletedMaintenancePercentage:    domain.NewDecimal2(stats.Summary.CompletedMaintenancePercentage),
			PausedMaintenancePercentage:       domain.NewDecimal2(stats.Summary.PausedMaintenancePercentage),
			StoppedMaintenancePercentage:      domain.NewDecimal2(stats.Summary.StoppedMaintenancePercentage),
			PreventiveMaintenancePercentage:   domain.NewDecimal2(stats.Summary.PreventiveMaintenancePercentage),
			CorrectiveMaintenancePercentage:   domain.NewDecimal2(stats.Summary.CorrectiveMaintenancePercentage),
			InspectionMaintenancePercentage:   domain.NewDecimal2(stats.Summary.InspectionMaintenancePercentage),
			CalibrationMaintenancePercentage:  domain.NewDecimal2(stats.Summary.CalibrationMaintenancePercentage),
			AverageScheduleFrequency:          domain.NewDecimal2(stats.Summary.AverageScheduleFrequency),
			UpcomingMaintenanceCount:          stats.Summary.UpcomingMaintenanceCount,
			OverdueMaintenanceCount:           stats.Summary.OverdueMaintenanceCount,
			AssetsWithScheduledMaintenance:    stats.Summary.AssetsWithScheduledMaintenance,
			AssetsWithoutScheduledMaintenance: stats.Summary.AssetsWithoutScheduledMaintenance,
			AverageSchedulesPerDay:            domain.NewDecimal2(stats.Summary.AverageSchedulesPerDay),
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
			ID:                up.ID,
			AssetID:           up.AssetID,
			AssetName:         up.AssetName,
			AssetTag:          up.AssetTag,
			MaintenanceType:   up.MaintenanceType,
			NextScheduledDate: up.NextScheduledDate,
			DaysUntilDue:      up.DaysUntilDue,
			Title:             up.Title,
			Description:       up.Description,
		}
	}

	resp.OverdueSchedule = make([]domain.OverdueMaintenanceScheduleResponse, len(stats.OverdueSchedule))
	for i, od := range stats.OverdueSchedule {
		resp.OverdueSchedule[i] = domain.OverdueMaintenanceScheduleResponse{
			ID:                od.ID,
			AssetID:           od.AssetID,
			AssetName:         od.AssetName,
			AssetTag:          od.AssetTag,
			MaintenanceType:   od.MaintenanceType,
			NextScheduledDate: od.NextScheduledDate,
			DaysOverdue:       od.DaysOverdue,
			Title:             od.Title,
			Description:       od.Description,
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
		domain.MaintenanceScheduleSortByNextScheduledDate: "ms.next_scheduled_date",
		domain.MaintenanceScheduleSortByMaintenanceType:   "ms.maintenance_type",
		domain.MaintenanceScheduleSortByState:             "ms.state",
		domain.MaintenanceScheduleSortByCreatedAt:         "ms.created_at",
		domain.MaintenanceScheduleSortByUpdatedAt:         "ms.updated_at",
	}

	if column, exists := columnMap[field]; exists {
		return column
	}
	return "ms.next_scheduled_date"
}

// *==================== Update Map conversions (Harus snake case karena untuk database) ====================
func ToModelMaintenanceScheduleUpdateMap(payload *domain.UpdateMaintenanceSchedulePayload) map[string]any {
	updates := make(map[string]any)

	if payload.MaintenanceType != nil {
		if *payload.MaintenanceType == "" {
			updates["maintenance_type"] = nil
		} else {
			updates["maintenance_type"] = *payload.MaintenanceType
		}
	}
	if payload.IsRecurring != nil {
		updates["is_recurring"] = *payload.IsRecurring
	}
	if payload.IntervalValue != nil {
		updates["interval_value"] = *payload.IntervalValue
	}
	if payload.IntervalUnit != nil {
		updates["interval_unit"] = *payload.IntervalUnit
	}
	if payload.ScheduledTime != nil {
		updates["scheduled_time"] = *payload.ScheduledTime
	}
	if payload.NextScheduledDate != nil {
		if *payload.NextScheduledDate == "" {
			updates["next_scheduled_date"] = nil
		} else {
			if parsedDate, err := time.ParseInLocation("2006-01-02", *payload.NextScheduledDate, time.UTC); err == nil {
				updates["next_scheduled_date"] = parsedDate
			}
		}
	}
	if payload.State != nil {
		updates["state"] = *payload.State
	}
	if payload.AutoComplete != nil {
		updates["auto_complete"] = *payload.AutoComplete
	}
	if payload.EstimatedCost != nil {
		updates["estimated_cost"] = *payload.EstimatedCost
	}

	return updates
}
