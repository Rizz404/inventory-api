package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
)

func findIssueReportTranslation(translations []model.IssueReportTranslation, langCode string) (title string, description, resolutionNotes *string) {
	for _, t := range translations {
		if t.LangCode == langCode {
			return t.Title, t.Description, t.ResolutionNotes
		}
	}
	for _, t := range translations {
		if t.LangCode == DefaultLangCode {
			return t.Title, t.Description, t.ResolutionNotes
		}
	}
	if len(translations) > 0 {
		return translations[0].Title, translations[0].Description, translations[0].ResolutionNotes
	}
	return "", nil, nil
}

func ToDomainIssueReportResponse(m model.IssueReport, langCode string) domain.IssueReportResponse {
	title, desc, resNotes := findIssueReportTranslation(m.Translations, langCode)
	resp := domain.IssueReportResponse{
		ID:              m.ID.String(),
		Asset:           ToDomainAssetResponse(m.Asset, langCode),
		ReportedBy:      ToDomainUserResponse(&m.ReportedByUser),
		ReportedDate:    m.ReportedDate.Format(TimeFormat),
		IssueType:       m.IssueType,
		Priority:        m.Priority,
		Status:          m.Status,
		Title:           title,
		Description:     desc,
		ResolutionNotes: resNotes,
	}
	if m.ResolvedDate != nil {
		resp.ResolvedDate = Ptr(m.ResolvedDate.Format(TimeFormat))
	}
	if m.ResolvedByUser != nil {
		resp.ResolvedBy = Ptr(ToDomainUserResponse(m.ResolvedByUser))
	}
	return resp
}
