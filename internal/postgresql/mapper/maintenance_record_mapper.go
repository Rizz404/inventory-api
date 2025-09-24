package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/oklog/ulid/v2"
)

// *==================== Model conversions ====================

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

func MaintenanceRecordStatisticsToResponse(stats *domain.MaintenanceRecordStatistics) domain.MaintenanceRecordStatisticsResponse {
	resp := domain.MaintenanceRecordStatisticsResponse{
		Total: domain.MaintenanceRecordCountStatisticsResponse{Count: stats.Total.Count},
		CostStatistics: domain.MaintenanceRecordCostStatisticsResponse{
			TotalCost:          stats.CostStatistics.TotalCost,
			AverageCost:        stats.CostStatistics.AverageCost,
			MinCost:            stats.CostStatistics.MinCost,
			MaxCost:            stats.CostStatistics.MaxCost,
			RecordsWithCost:    stats.CostStatistics.RecordsWithCost,
			RecordsWithoutCost: stats.CostStatistics.RecordsWithoutCost,
		},
		Summary: domain.MaintenanceRecordSummaryStatisticsResponse{
			TotalRecords:                  stats.Summary.TotalRecords,
			RecordsWithCostInfo:           stats.Summary.RecordsWithCostInfo,
			CostInfoPercentage:            stats.Summary.CostInfoPercentage,
			TotalUniqueVendors:            stats.Summary.TotalUniqueVendors,
			TotalUniquePerformers:         stats.Summary.TotalUniquePerformers,
			AverageRecordsPerDay:          stats.Summary.AverageRecordsPerDay,
			LatestRecordDate:              stats.Summary.LatestRecordDate,
			EarliestRecordDate:            stats.Summary.EarliestRecordDate,
			MostExpensiveMaintenanceCost:  stats.Summary.MostExpensiveMaintenanceCost,
			LeastExpensiveMaintenanceCost: stats.Summary.LeastExpensiveMaintenanceCost,
			AssetsWithMaintenance:         stats.Summary.AssetsWithMaintenance,
			AverageMaintenancePerAsset:    stats.Summary.AverageMaintenancePerAsset,
		},
	}

	// ByAsset records
	resp.ByAsset = make([]domain.AssetMaintenanceRecordStatisticsResponse, len(stats.ByAsset))
	for i, a := range stats.ByAsset {
		resp.ByAsset[i] = domain.AssetMaintenanceRecordStatisticsResponse{
			AssetID:         a.AssetID,
			AssetName:       a.AssetName,
			AssetTag:        a.AssetTag,
			RecordCount:     a.RecordCount,
			LastMaintenance: a.LastMaintenance,
			TotalCost:       a.TotalCost,
			AverageCost:     a.AverageCost,
		}
	}

	// ByPerformer
	resp.ByPerformer = make([]domain.UserMaintenanceRecordStatisticsResponse, len(stats.ByPerformer))
	for i, u := range stats.ByPerformer {
		resp.ByPerformer[i] = domain.UserMaintenanceRecordStatisticsResponse{
			UserID:      u.UserID,
			UserName:    u.UserName,
			UserEmail:   u.UserEmail,
			Count:       u.Count,
			TotalCost:   u.TotalCost,
			AverageCost: u.AverageCost,
		}
	}

	// ByVendor
	resp.ByVendor = make([]domain.VendorMaintenanceRecordStatisticsResponse, len(stats.ByVendor))
	for i, v := range stats.ByVendor {
		resp.ByVendor[i] = domain.VendorMaintenanceRecordStatisticsResponse{
			VendorName:  v.VendorName,
			Count:       v.Count,
			TotalCost:   v.TotalCost,
			AverageCost: v.AverageCost,
		}
	}

	// Completion trends
	resp.CompletionTrend = make([]domain.MaintenanceRecordCompletionTrendResponse, len(stats.CompletionTrend))
	for i, ct := range stats.CompletionTrend {
		resp.CompletionTrend[i] = domain.MaintenanceRecordCompletionTrendResponse{Date: ct.Date, Count: ct.Count}
	}

	// Monthly trends
	resp.MonthlyTrends = make([]domain.MaintenanceRecordMonthlyTrendResponse, len(stats.MonthlyTrends))
	for i, mt := range stats.MonthlyTrends {
		resp.MonthlyTrends[i] = domain.MaintenanceRecordMonthlyTrendResponse{
			Month:       mt.Month,
			RecordCount: mt.RecordCount,
			TotalCost:   mt.TotalCost,
		}
	}

	return resp
}
