package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
)

func findMaintenanceScheduleTranslation(translations []model.MaintenanceScheduleTranslation, langCode string) (title string, description *string) {
	for _, t := range translations {
		if t.LangCode == langCode {
			return t.Title, t.Description
		}
	}
	for _, t := range translations {
		if t.LangCode == DefaultLangCode {
			return t.Title, t.Description
		}
	}
	if len(translations) > 0 {
		return translations[0].Title, translations[0].Description
	}
	return "", nil
}

func ToDomainMaintenanceScheduleResponse(m model.MaintenanceSchedule, langCode string) domain.MaintenanceScheduleResponse {
	title, desc := findMaintenanceScheduleTranslation(m.Translations, langCode)
	return domain.MaintenanceScheduleResponse{
		ID:              m.ID.String(),
		Asset:           ToDomainAssetResponse(m.Asset, langCode),
		MaintenanceType: m.MaintenanceType,
		ScheduledDate:   m.ScheduledDate.Format(DateFormat),
		FrequencyMonths: m.FrequencyMonths,
		Status:          m.Status,
		CreatedBy:       ToDomainUserResponse(&m.CreatedByUser),
		CreatedAt:       m.CreatedAt.Format(TimeFormat),
		Title:           title,
		Description:     desc,
	}
}

func findMaintenanceRecordTranslation(translations []model.MaintenanceRecordTranslation, langCode string) (title string, notes *string) {
	for _, t := range translations {
		if t.LangCode == langCode {
			return t.Title, t.Notes
		}
	}
	for _, t := range translations {
		if t.LangCode == DefaultLangCode {
			return t.Title, t.Notes
		}
	}
	if len(translations) > 0 {
		return translations[0].Title, translations[0].Notes
	}
	return "", nil
}

func ToDomainMaintenanceRecordResponse(m model.MaintenanceRecord, langCode string) domain.MaintenanceRecordResponse {
	title, notes := findMaintenanceRecordTranslation(m.Translations, langCode)
	resp := domain.MaintenanceRecordResponse{
		ID:                m.ID.String(),
		Asset:             ToDomainAssetResponse(m.Asset, langCode),
		MaintenanceDate:   m.MaintenanceDate.Format(DateFormat),
		PerformedByVendor: m.PerformedByVendor,
		ActualCost:        m.ActualCost,
		Title:             title,
		Notes:             notes,
	}
	if m.Schedule != nil {
		resp.Schedule = Ptr(ToDomainMaintenanceScheduleResponse(*m.Schedule, langCode))
	}
	if m.User != nil {
		resp.PerformedByUser = Ptr(ToDomainUserResponse(m.User))
	}
	return resp
}
