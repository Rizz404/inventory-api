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
		ScheduledDate:   d.ScheduledDate.Format(DateFormat),
		FrequencyMonths: d.FrequencyMonths,
		Status:          d.Status,
		CreatedByID:     d.CreatedBy,
		CreatedAt:       d.CreatedAt.Format(TimeFormat),
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
		ScheduledDate:   d.ScheduledDate.Format(DateFormat),
		FrequencyMonths: d.FrequencyMonths,
		Status:          d.Status,
		CreatedByID:     d.CreatedBy,
		CreatedAt:       d.CreatedAt.Format(TimeFormat),
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

// *==================== Entity conversions ====================

func ToModelMaintenanceRecord(d *domain.MaintenanceRecord) model.MaintenanceRecord {
	modelRecord := model.MaintenanceRecord{
		MaintenanceDate:   d.MaintenanceDate,
		PerformedByVendor: d.PerformedByVendor,
		ActualCost:        d.ActualCost,
	}

	if d.ID != "" {
		if parsedID, err := ulid.Parse(d.ID); err == nil {
			modelRecord.ID = model.SQLULID(parsedID)
		}
	}

	if d.ScheduleID != nil && *d.ScheduleID != "" {
		if parsedScheduleID, err := ulid.Parse(*d.ScheduleID); err == nil {
			modelULID := model.SQLULID(parsedScheduleID)
			modelRecord.ScheduleID = &modelULID
		}
	}

	if d.AssetID != "" {
		if parsedAssetID, err := ulid.Parse(d.AssetID); err == nil {
			modelRecord.AssetID = model.SQLULID(parsedAssetID)
		}
	}

	if d.PerformedByUser != nil && *d.PerformedByUser != "" {
		if parsedPerformedByUser, err := ulid.Parse(*d.PerformedByUser); err == nil {
			modelULID := model.SQLULID(parsedPerformedByUser)
			modelRecord.PerformedByUser = &modelULID
		}
	}

	return modelRecord
}

func ToModelMaintenanceRecordForCreate(d *domain.MaintenanceRecord) model.MaintenanceRecord {
	modelRecord := model.MaintenanceRecord{
		MaintenanceDate:   d.MaintenanceDate,
		PerformedByVendor: d.PerformedByVendor,
		ActualCost:        d.ActualCost,
	}

	if d.ScheduleID != nil && *d.ScheduleID != "" {
		if parsedScheduleID, err := ulid.Parse(*d.ScheduleID); err == nil {
			modelULID := model.SQLULID(parsedScheduleID)
			modelRecord.ScheduleID = &modelULID
		}
	}

	if d.AssetID != "" {
		if parsedAssetID, err := ulid.Parse(d.AssetID); err == nil {
			modelRecord.AssetID = model.SQLULID(parsedAssetID)
		}
	}

	if d.PerformedByUser != nil && *d.PerformedByUser != "" {
		if parsedPerformedByUser, err := ulid.Parse(*d.PerformedByUser); err == nil {
			modelULID := model.SQLULID(parsedPerformedByUser)
			modelRecord.PerformedByUser = &modelULID
		}
	}

	return modelRecord
}

func ToModelMaintenanceRecordTranslation(d *domain.MaintenanceRecordTranslation) model.MaintenanceRecordTranslation {
	modelTranslation := model.MaintenanceRecordTranslation{
		LangCode: d.LangCode,
		Title:    d.Title,
		Notes:    d.Notes,
	}

	if d.ID != "" {
		if parsedID, err := ulid.Parse(d.ID); err == nil {
			modelTranslation.ID = model.SQLULID(parsedID)
		}
	}

	if d.RecordID != "" {
		if parsedRecordID, err := ulid.Parse(d.RecordID); err == nil {
			modelTranslation.RecordID = model.SQLULID(parsedRecordID)
		}
	}

	return modelTranslation
}

func ToModelMaintenanceRecordTranslationForCreate(recordID string, d *domain.MaintenanceRecordTranslation) model.MaintenanceRecordTranslation {
	modelTranslation := model.MaintenanceRecordTranslation{
		LangCode: d.LangCode,
		Title:    d.Title,
		Notes:    d.Notes,
	}

	if recordID != "" {
		if parsedRecordID, err := ulid.Parse(recordID); err == nil {
			modelTranslation.RecordID = model.SQLULID(parsedRecordID)
		}
	}

	return modelTranslation
}

func ToDomainMaintenanceRecord(m *model.MaintenanceRecord) domain.MaintenanceRecord {
	domainRecord := domain.MaintenanceRecord{
		ID:                m.ID.String(),
		AssetID:           m.AssetID.String(),
		MaintenanceDate:   m.MaintenanceDate,
		PerformedByVendor: m.PerformedByVendor,
		ActualCost:        m.ActualCost,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}

	if m.ScheduleID != nil && !m.ScheduleID.IsZero() {
		scheduleIDStr := m.ScheduleID.String()
		domainRecord.ScheduleID = &scheduleIDStr
	}

	if m.PerformedByUser != nil && !m.PerformedByUser.IsZero() {
		performedByUserStr := m.PerformedByUser.String()
		domainRecord.PerformedByUser = &performedByUserStr
	}

	if len(m.Translations) > 0 {
		domainRecord.Translations = make([]domain.MaintenanceRecordTranslation, len(m.Translations))
		for i, translation := range m.Translations {
			domainRecord.Translations[i] = ToDomainMaintenanceRecordTranslation(&translation)
		}
	}

	return domainRecord
}

func ToDomainMaintenanceRecordTranslation(m *model.MaintenanceRecordTranslation) domain.MaintenanceRecordTranslation {
	return domain.MaintenanceRecordTranslation{
		ID:       m.ID.String(),
		RecordID: m.RecordID.String(),
		LangCode: m.LangCode,
		Title:    m.Title,
		Notes:    m.Notes,
	}
}

func ToDomainMaintenanceRecords(models []model.MaintenanceRecord) []domain.MaintenanceRecord {
	records := make([]domain.MaintenanceRecord, len(models))
	for i, m := range models {
		records[i] = ToDomainMaintenanceRecord(&m)
	}
	return records
}

// *==================== Entity Response conversions ====================
func MaintenanceRecordToResponse(d *domain.MaintenanceRecord, langCode string) domain.MaintenanceRecordResponse {
	response := domain.MaintenanceRecordResponse{
		ID:                d.ID,
		ScheduleID:        d.ScheduleID,
		AssetID:           d.AssetID,
		MaintenanceDate:   d.MaintenanceDate.Format(DateFormat),
		PerformedByUserID: d.PerformedByUser,
		PerformedByVendor: d.PerformedByVendor,
		ActualCost:        d.ActualCost,
		CreatedAt:         d.CreatedAt.Format(TimeFormat),
		UpdatedAt:         d.UpdatedAt.Format(TimeFormat),
		Translations:      make([]domain.MaintenanceRecordTranslationResponse, len(d.Translations)),
	}

	// Populate translations
	for i, translation := range d.Translations {
		response.Translations[i] = domain.MaintenanceRecordTranslationResponse{
			LangCode: translation.LangCode,
			Title:    translation.Title,
			Notes:    translation.Notes,
		}
	}

	// Find translation for the requested language
	for _, translation := range d.Translations {
		if translation.LangCode == langCode {
			response.Title = translation.Title
			response.Notes = translation.Notes
			break
		}
	}

	// If no translation found for requested language, use first available
	if response.Title == "" && len(d.Translations) > 0 {
		response.Title = d.Translations[0].Title
		response.Notes = d.Translations[0].Notes
	}

	return response
}

func MaintenanceRecordToListResponse(d *domain.MaintenanceRecord, langCode string) domain.MaintenanceRecordListResponse {
	response := domain.MaintenanceRecordListResponse{
		ID:                d.ID,
		ScheduleID:        d.ScheduleID,
		AssetID:           d.AssetID,
		MaintenanceDate:   d.MaintenanceDate.Format(DateFormat),
		PerformedByUserID: d.PerformedByUser,
		PerformedByVendor: d.PerformedByVendor,
		ActualCost:        d.ActualCost,
		CreatedAt:         d.CreatedAt.Format(TimeFormat),
		UpdatedAt:         d.UpdatedAt.Format(TimeFormat),
	}

	// Find translation for the requested language
	for _, translation := range d.Translations {
		if translation.LangCode == langCode {
			response.Title = translation.Title
			response.Notes = translation.Notes
			break
		}
	}

	// If no translation found for requested language, use first available
	if response.Title == "" && len(d.Translations) > 0 {
		response.Title = d.Translations[0].Title
		response.Notes = d.Translations[0].Notes
	}

	return response
}

func MaintenanceRecordsToResponses(records []domain.MaintenanceRecord, langCode string) []domain.MaintenanceRecordResponse {
	responses := make([]domain.MaintenanceRecordResponse, len(records))
	for i, record := range records {
		responses[i] = MaintenanceRecordToResponse(&record, langCode)
	}
	return responses
}

func MaintenanceRecordsToListResponses(records []domain.MaintenanceRecord, langCode string) []domain.MaintenanceRecordListResponse {
	responses := make([]domain.MaintenanceRecordListResponse, len(records))
	for i, record := range records {
		responses[i] = MaintenanceRecordToListResponse(&record, langCode)
	}
	return responses
}

// *==================== Statistics conversions ====================

func MaintenanceStatisticsToResponse(stats *domain.MaintenanceStatistics) domain.MaintenanceStatisticsResponse {
	resp := domain.MaintenanceStatisticsResponse{
		Schedules: domain.MaintenanceScheduleStatisticsResponse{
			Total: domain.MaintenanceCountStatisticsResponse{Count: stats.Schedules.Total.Count},
			ByType: domain.MaintenanceTypeStatisticsResponse{
				Preventive: stats.Schedules.ByType.Preventive,
				Corrective: stats.Schedules.ByType.Corrective,
			},
			ByStatus: domain.MaintenanceScheduleStatusStatisticsResponse{
				Scheduled: stats.Schedules.ByStatus.Scheduled,
				Completed: stats.Schedules.ByStatus.Completed,
				Cancelled: stats.Schedules.ByStatus.Cancelled,
			},
		},
		Records: domain.MaintenanceRecordStatisticsResponse{
			Total: domain.MaintenanceCountStatisticsResponse{Count: stats.Records.Total.Count},
		},
		Summary: domain.MaintenanceSummaryStatisticsResponse{
			TotalSchedules:                  stats.Summary.TotalSchedules,
			TotalRecords:                    stats.Summary.TotalRecords,
			ScheduledMaintenancePercentage:  stats.Summary.ScheduledMaintenancePercentage,
			CompletedMaintenancePercentage:  stats.Summary.CompletedMaintenancePercentage,
			CancelledMaintenancePercentage:  stats.Summary.CancelledMaintenancePercentage,
			PreventiveMaintenancePercentage: stats.Summary.PreventiveMaintenancePercentage,
			CorrectiveMaintenancePercentage: stats.Summary.CorrectiveMaintenancePercentage,
			AverageMaintenancePerAsset:      stats.Summary.AverageMaintenancePerAsset,
			AssetsWithMaintenance:           stats.Summary.AssetsWithMaintenance,
			AssetsWithoutMaintenance:        stats.Summary.AssetsWithoutMaintenance,
			MaintenanceComplianceRate:       stats.Summary.MaintenanceComplianceRate,
			AverageMaintenanceFrequency:     stats.Summary.AverageMaintenanceFrequency,
			UpcomingMaintenanceCount:        stats.Summary.UpcomingMaintenanceCount,
			OverdueMaintenanceCount:         stats.Summary.OverdueMaintenanceCount,
			RecordsWithCostInfo:             stats.Summary.RecordsWithCostInfo,
			CostInfoPercentage:              stats.Summary.CostInfoPercentage,
			TotalUniqueVendors:              stats.Summary.TotalUniqueVendors,
			TotalUniquePerformers:           stats.Summary.TotalUniquePerformers,
			AverageRecordsPerDay:            stats.Summary.AverageRecordsPerDay,
			LatestRecordDate:                stats.Summary.LatestRecordDate,
			EarliestRecordDate:              stats.Summary.EarliestRecordDate,
			MostExpensiveMaintenanceCost:    stats.Summary.MostExpensiveMaintenanceCost,
			LeastExpensiveMaintenanceCost:   stats.Summary.LeastExpensiveMaintenanceCost,
		},
	}

	// ByAsset schedules and records
	resp.Schedules.ByAsset = make([]domain.AssetMaintenanceStatisticsResponse, len(stats.Schedules.ByAsset))
	for i, a := range stats.Schedules.ByAsset {
		resp.Schedules.ByAsset[i] = domain.AssetMaintenanceStatisticsResponse{
			AssetID:         a.AssetID,
			AssetName:       a.AssetName,
			AssetTag:        a.AssetTag,
			ScheduleCount:   a.ScheduleCount,
			RecordCount:     a.RecordCount,
			LastMaintenance: a.LastMaintenance,
			NextMaintenance: a.NextMaintenance,
		}
	}

	resp.Records.ByAsset = make([]domain.AssetMaintenanceStatisticsResponse, len(stats.Records.ByAsset))
	for i, a := range stats.Records.ByAsset {
		resp.Records.ByAsset[i] = domain.AssetMaintenanceStatisticsResponse{
			AssetID:         a.AssetID,
			AssetName:       a.AssetName,
			AssetTag:        a.AssetTag,
			ScheduleCount:   a.ScheduleCount,
			RecordCount:     a.RecordCount,
			LastMaintenance: a.LastMaintenance,
			NextMaintenance: a.NextMaintenance,
		}
	}

	// ByCreator / ByPerformer
	resp.Schedules.ByCreator = make([]domain.UserMaintenanceStatisticsResponse, len(stats.Schedules.ByCreator))
	for i, u := range stats.Schedules.ByCreator {
		resp.Schedules.ByCreator[i] = domain.UserMaintenanceStatisticsResponse{
			UserID:      u.UserID,
			UserName:    u.UserName,
			UserEmail:   u.UserEmail,
			Count:       u.Count,
			TotalCost:   u.TotalCost,
			AverageCost: u.AverageCost,
		}
	}

	resp.Records.ByPerformer = make([]domain.UserMaintenanceStatisticsResponse, len(stats.Records.ByPerformer))
	for i, u := range stats.Records.ByPerformer {
		resp.Records.ByPerformer[i] = domain.UserMaintenanceStatisticsResponse{
			UserID:      u.UserID,
			UserName:    u.UserName,
			UserEmail:   u.UserEmail,
			Count:       u.Count,
			TotalCost:   u.TotalCost,
			AverageCost: u.AverageCost,
		}
	}

	// ByVendor
	resp.Records.ByVendor = make([]domain.VendorMaintenanceStatisticsResponse, len(stats.Records.ByVendor))
	for i, v := range stats.Records.ByVendor {
		resp.Records.ByVendor[i] = domain.VendorMaintenanceStatisticsResponse{
			VendorName:  v.VendorName,
			Count:       v.Count,
			TotalCost:   v.TotalCost,
			AverageCost: v.AverageCost,
		}
	}

	// Upcoming & Overdue schedules
	resp.Schedules.UpcomingSchedule = make([]domain.UpcomingMaintenanceScheduleResponse, len(stats.Schedules.UpcomingSchedule))
	for i, up := range stats.Schedules.UpcomingSchedule {
		resp.Schedules.UpcomingSchedule[i] = domain.UpcomingMaintenanceScheduleResponse{
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

	resp.Schedules.OverdueSchedule = make([]domain.OverdueMaintenanceScheduleResponse, len(stats.Schedules.OverdueSchedule))
	for i, od := range stats.Schedules.OverdueSchedule {
		resp.Schedules.OverdueSchedule[i] = domain.OverdueMaintenanceScheduleResponse{
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
	resp.Schedules.FrequencyTrends = make([]domain.MaintenanceFrequencyTrendResponse, len(stats.Schedules.FrequencyTrends))
	for i, ft := range stats.Schedules.FrequencyTrends {
		resp.Schedules.FrequencyTrends[i] = domain.MaintenanceFrequencyTrendResponse{FrequencyMonths: ft.FrequencyMonths, Count: ft.Count}
	}

	// Completion trends
	resp.Records.CompletionTrend = make([]domain.MaintenanceCompletionTrendResponse, len(stats.Records.CompletionTrend))
	for i, ct := range stats.Records.CompletionTrend {
		resp.Records.CompletionTrend[i] = domain.MaintenanceCompletionTrendResponse{Date: ct.Date, Count: ct.Count}
	}

	// Monthly trends
	resp.Records.MonthlyTrends = make([]domain.MaintenanceMonthlyTrendResponse, len(stats.Records.MonthlyTrends))
	for i, mt := range stats.Records.MonthlyTrends {
		resp.Records.MonthlyTrends[i] = domain.MaintenanceMonthlyTrendResponse{
			Month:           mt.Month,
			ScheduleCount:   mt.ScheduleCount,
			RecordCount:     mt.RecordCount,
			TotalCost:       mt.TotalCost,
			PreventiveCount: mt.PreventiveCount,
			CorrectiveCount: mt.CorrectiveCount,
		}
	}

	return resp
}
