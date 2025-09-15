package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/oklog/ulid/v2"
)

// ===== Maintenance Schedule Mappers =====

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

func ToDomainMaintenanceScheduleResponse(m *model.MaintenanceSchedule, langCode string) domain.MaintenanceScheduleResponse {
	response := domain.MaintenanceScheduleResponse{
		ID:              m.ID.String(),
		MaintenanceType: m.MaintenanceType,
		ScheduledDate:   m.ScheduledDate.Format(DateFormat),
		FrequencyMonths: m.FrequencyMonths,
		Status:          m.Status,
		CreatedAt:       m.CreatedAt.Format(TimeFormat),
	}

	// Find translation for the requested language
	for _, translation := range m.Translations {
		if translation.LangCode == langCode {
			response.Title = translation.Title
			response.Description = translation.Description
			break
		}
	}

	// If no translation found for requested language, use first available
	if response.Title == "" && len(m.Translations) > 0 {
		response.Title = m.Translations[0].Title
		response.Description = m.Translations[0].Description
	}

	// Handle related entities
	if !m.Asset.ID.IsZero() {
		assetResponse := ToDomainAssetResponse(&m.Asset, langCode)
		response.Asset = assetResponse
	}

	if !m.CreatedByUser.ID.IsZero() {
		userResponse := ToDomainUserResponse(&m.CreatedByUser)
		response.CreatedBy = userResponse
	}

	return response
}

func ToDomainMaintenanceSchedulesResponse(m []model.MaintenanceSchedule, langCode string) []domain.MaintenanceScheduleResponse {
	responses := make([]domain.MaintenanceScheduleResponse, len(m))
	for i, schedule := range m {
		responses[i] = ToDomainMaintenanceScheduleResponse(&schedule, langCode)
	}
	return responses
}

// * Convert domain.MaintenanceSchedule directly to domain.MaintenanceScheduleResponse without going through model.MaintenanceSchedule
func DomainMaintenanceScheduleToMaintenanceScheduleResponse(d *domain.MaintenanceSchedule, langCode string) domain.MaintenanceScheduleResponse {
	response := domain.MaintenanceScheduleResponse{
		ID:              d.ID,
		MaintenanceType: d.MaintenanceType,
		ScheduledDate:   d.ScheduledDate.Format(DateFormat),
		FrequencyMonths: d.FrequencyMonths,
		Status:          d.Status,
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

// * Convert slice of domain.MaintenanceSchedule to slice of domain.MaintenanceScheduleResponse
func DomainMaintenanceSchedulesToMaintenanceSchedulesResponse(schedules []domain.MaintenanceSchedule, langCode string) []domain.MaintenanceScheduleResponse {
	responses := make([]domain.MaintenanceScheduleResponse, len(schedules))
	for i, schedule := range schedules {
		responses[i] = DomainMaintenanceScheduleToMaintenanceScheduleResponse(&schedule, langCode)
	}
	return responses
}

// ===== Maintenance Record Mappers =====

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

func ToDomainMaintenanceRecordResponse(m *model.MaintenanceRecord, langCode string) domain.MaintenanceRecordResponse {
	response := domain.MaintenanceRecordResponse{
		ID:                m.ID.String(),
		MaintenanceDate:   m.MaintenanceDate.Format(DateFormat),
		PerformedByVendor: m.PerformedByVendor,
		ActualCost:        m.ActualCost,
		CreatedAt:         m.CreatedAt.Format(TimeFormat),
		UpdatedAt:         m.UpdatedAt.Format(TimeFormat),
	}

	// Find translation for the requested language
	for _, translation := range m.Translations {
		if translation.LangCode == langCode {
			response.Title = translation.Title
			response.Notes = translation.Notes
			break
		}
	}

	// If no translation found for requested language, use first available
	if response.Title == "" && len(m.Translations) > 0 {
		response.Title = m.Translations[0].Title
		response.Notes = m.Translations[0].Notes
	}

	// Handle related entities
	if m.Schedule != nil && !m.Schedule.ID.IsZero() {
		scheduleResponse := ToDomainMaintenanceScheduleResponse(m.Schedule, langCode)
		response.Schedule = &scheduleResponse
	}

	if !m.Asset.ID.IsZero() {
		assetResponse := ToDomainAssetResponse(&m.Asset, langCode)
		response.Asset = assetResponse
	}

	if m.User != nil && !m.User.ID.IsZero() {
		userResponse := ToDomainUserResponse(m.User)
		response.PerformedByUser = &userResponse
	}

	return response
}

func ToDomainMaintenanceRecordsResponse(m []model.MaintenanceRecord, langCode string) []domain.MaintenanceRecordResponse {
	responses := make([]domain.MaintenanceRecordResponse, len(m))
	for i, record := range m {
		responses[i] = ToDomainMaintenanceRecordResponse(&record, langCode)
	}
	return responses
}

// * Convert domain.MaintenanceRecord directly to domain.MaintenanceRecordResponse without going through model.MaintenanceRecord
func DomainMaintenanceRecordToMaintenanceRecordResponse(d *domain.MaintenanceRecord, langCode string) domain.MaintenanceRecordResponse {
	response := domain.MaintenanceRecordResponse{
		ID:                d.ID,
		MaintenanceDate:   d.MaintenanceDate.Format(DateFormat),
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

// * Convert slice of domain.MaintenanceRecord to slice of domain.MaintenanceRecordResponse
func DomainMaintenanceRecordsToMaintenanceRecordsResponse(records []domain.MaintenanceRecord, langCode string) []domain.MaintenanceRecordResponse {
	responses := make([]domain.MaintenanceRecordResponse, len(records))
	for i, record := range records {
		responses[i] = DomainMaintenanceRecordToMaintenanceRecordResponse(&record, langCode)
	}
	return responses
}
