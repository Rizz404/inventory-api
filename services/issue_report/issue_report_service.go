package issue_report

import (
	"context"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/Rizz404/inventory-api/internal/utils"
)

// * Repository interface defines the contract for issue report data operations
type Repository interface {
	// * MUTATION
	CreateIssueReport(ctx context.Context, payload *domain.IssueReport) (domain.IssueReport, error)
	UpdateIssueReport(ctx context.Context, issueReportId string, payload *domain.UpdateIssueReportPayload) (domain.IssueReport, error)
	ResolveIssueReport(ctx context.Context, issueReportId string, resolvedBy string, resolutionNotes string) (domain.IssueReport, error)
	ReopenIssueReport(ctx context.Context, issueReportId string) (domain.IssueReport, error)
	DeleteIssueReport(ctx context.Context, issueReportId string) error

	// * QUERY
	GetIssueReportsPaginated(ctx context.Context, params domain.IssueReportParams, langCode string) ([]domain.IssueReport, error)
	GetIssueReportsCursor(ctx context.Context, params domain.IssueReportParams, langCode string) ([]domain.IssueReport, error)
	GetIssueReportById(ctx context.Context, issueReportId string) (domain.IssueReport, error)
	CheckIssueReportExist(ctx context.Context, issueReportId string) (bool, error)
	CountIssueReports(ctx context.Context, params domain.IssueReportParams) (int64, error)
	GetIssueReportStatistics(ctx context.Context) (domain.IssueReportStatistics, error)
}

// * IssueReportService interface defines the contract for issue report business operations
type IssueReportService interface {
	// * MUTATION
	CreateIssueReport(ctx context.Context, payload *domain.CreateIssueReportPayload, reportedBy string) (domain.IssueReportResponse, error)
	UpdateIssueReport(ctx context.Context, issueReportId string, payload *domain.UpdateIssueReportPayload) (domain.IssueReportResponse, error)
	ResolveIssueReport(ctx context.Context, issueReportId string, resolvedBy string, payload *domain.ResolveIssueReportPayload) (domain.IssueReportResponse, error)
	ReopenIssueReport(ctx context.Context, issueReportId string) (domain.IssueReportResponse, error)
	DeleteIssueReport(ctx context.Context, issueReportId string) error

	// * QUERY
	GetIssueReportsPaginated(ctx context.Context, params domain.IssueReportParams, langCode string) ([]domain.IssueReportListResponse, int64, error)
	GetIssueReportsCursor(ctx context.Context, params domain.IssueReportParams, langCode string) ([]domain.IssueReportListResponse, error)
	GetIssueReportById(ctx context.Context, issueReportId string, langCode string) (domain.IssueReportResponse, error)
	CheckIssueReportExists(ctx context.Context, issueReportId string) (bool, error)
	CountIssueReports(ctx context.Context, params domain.IssueReportParams) (int64, error)
	GetIssueReportStatistics(ctx context.Context) (domain.IssueReportStatisticsResponse, error)
}

type Service struct {
	Repo Repository
}

// * Ensure Service implements IssueReportService interface
var _ IssueReportService = (*Service)(nil)

func NewService(r Repository) IssueReportService {
	return &Service{
		Repo: r,
	}
}

// *===========================MUTATION===========================*
func (s *Service) CreateIssueReport(ctx context.Context, payload *domain.CreateIssueReportPayload, reportedBy string) (domain.IssueReportResponse, error) {
	// * Prepare domain issue report
	newIssueReport := domain.IssueReport{
		AssetID:      payload.AssetID,
		ReportedBy:   reportedBy,
		ReportedDate: time.Now(),
		IssueType:    payload.IssueType,
		Priority:     payload.Priority,
		Status:       domain.IssueStatusOpen, // New reports are always open
		Translations: make([]domain.IssueReportTranslation, len(payload.Translations)),
	}

	// * Convert translation payloads to domain translations
	for i, translationPayload := range payload.Translations {
		newIssueReport.Translations[i] = domain.IssueReportTranslation{
			LangCode:    translationPayload.LangCode,
			Title:       translationPayload.Title,
			Description: translationPayload.Description,
		}
	}

	createdIssueReport, err := s.Repo.CreateIssueReport(ctx, &newIssueReport)
	if err != nil {
		return domain.IssueReportResponse{}, err
	}

	// * Convert to IssueReportResponse using mapper
	return mapper.IssueReportToResponse(&createdIssueReport, mapper.DefaultLangCode), nil
}

func (s *Service) UpdateIssueReport(ctx context.Context, issueReportId string, payload *domain.UpdateIssueReportPayload) (domain.IssueReportResponse, error) {
	// * Check if issue report exists
	existingReport, err := s.Repo.GetIssueReportById(ctx, issueReportId)
	if err != nil {
		return domain.IssueReportResponse{}, err
	}

	// * Business logic: Check if the report is already closed and trying to change status
	if existingReport.Status == domain.IssueStatusClosed && payload.Status != nil && *payload.Status != domain.IssueStatusClosed {
		return domain.IssueReportResponse{}, domain.ErrBadRequestWithKey(utils.ErrIssueReportCannotReopenKey)
	}

	updatedIssueReport, err := s.Repo.UpdateIssueReport(ctx, issueReportId, payload)
	if err != nil {
		return domain.IssueReportResponse{}, err
	}

	// * Convert to IssueReportResponse using mapper
	return mapper.IssueReportToResponse(&updatedIssueReport, mapper.DefaultLangCode), nil
}

func (s *Service) ResolveIssueReport(ctx context.Context, issueReportId string, resolvedBy string, payload *domain.ResolveIssueReportPayload) (domain.IssueReportResponse, error) {
	// * Check if issue report exists
	existingReport, err := s.Repo.GetIssueReportById(ctx, issueReportId)
	if err != nil {
		return domain.IssueReportResponse{}, err
	}

	// * Business logic: Check if the report is already resolved
	if existingReport.Status == domain.IssueStatusResolved || existingReport.Status == domain.IssueStatusClosed {
		return domain.IssueReportResponse{}, domain.ErrBadRequestWithKey(utils.ErrIssueReportAlreadyResolvedKey)
	}

	resolvedIssueReport, err := s.Repo.ResolveIssueReport(ctx, issueReportId, resolvedBy, payload.ResolutionNotes)
	if err != nil {
		return domain.IssueReportResponse{}, err
	}

	// * Convert to IssueReportResponse using mapper
	return mapper.IssueReportToResponse(&resolvedIssueReport, mapper.DefaultLangCode), nil
}

func (s *Service) ReopenIssueReport(ctx context.Context, issueReportId string) (domain.IssueReportResponse, error) {
	// * Check if issue report exists
	existingReport, err := s.Repo.GetIssueReportById(ctx, issueReportId)
	if err != nil {
		return domain.IssueReportResponse{}, err
	}

	// * Business logic: Check if the report can be reopened
	if existingReport.Status == domain.IssueStatusClosed {
		return domain.IssueReportResponse{}, domain.ErrBadRequestWithKey(utils.ErrIssueReportCannotReopenKey)
	}

	if existingReport.Status == domain.IssueStatusOpen || existingReport.Status == domain.IssueStatusInProgress {
		return domain.IssueReportResponse{}, domain.ErrBadRequest("issue report is already open")
	}

	reopenedIssueReport, err := s.Repo.ReopenIssueReport(ctx, issueReportId)
	if err != nil {
		return domain.IssueReportResponse{}, err
	}

	// * Convert to IssueReportResponse using mapper
	return mapper.IssueReportToResponse(&reopenedIssueReport, mapper.DefaultLangCode), nil
}

func (s *Service) DeleteIssueReport(ctx context.Context, issueReportId string) error {
	err := s.Repo.DeleteIssueReport(ctx, issueReportId)
	if err != nil {
		return err
	}
	return nil
}

// *===========================QUERY===========================*
func (s *Service) GetIssueReportsPaginated(ctx context.Context, params domain.IssueReportParams, langCode string) ([]domain.IssueReportListResponse, int64, error) {
	listItems, err := s.Repo.GetIssueReportsPaginated(ctx, params, langCode)
	if err != nil {
		return nil, 0, err
	}

	// * Count total for pagination
	count, err := s.Repo.CountIssueReports(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// Convert IssueReport to IssueReportListResponse using mapper
	responses := mapper.IssueReportsToListResponses(listItems, langCode)

	return responses, count, nil
}

func (s *Service) GetIssueReportsCursor(ctx context.Context, params domain.IssueReportParams, langCode string) ([]domain.IssueReportListResponse, error) {
	listItems, err := s.Repo.GetIssueReportsCursor(ctx, params, langCode)
	if err != nil {
		return nil, err
	}

	// Convert IssueReport to IssueReportListResponse using mapper
	responses := mapper.IssueReportsToListResponses(listItems, langCode)

	return responses, nil
}

func (s *Service) GetIssueReportById(ctx context.Context, issueReportId string, langCode string) (domain.IssueReportResponse, error) {
	issueReport, err := s.Repo.GetIssueReportById(ctx, issueReportId)
	if err != nil {
		return domain.IssueReportResponse{}, err
	}

	// * Convert to IssueReportResponse using mapper
	return mapper.IssueReportToResponse(&issueReport, langCode), nil
}

func (s *Service) CheckIssueReportExists(ctx context.Context, issueReportId string) (bool, error) {
	exists, err := s.Repo.CheckIssueReportExist(ctx, issueReportId)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Service) CountIssueReports(ctx context.Context, params domain.IssueReportParams) (int64, error) {
	count, err := s.Repo.CountIssueReports(ctx, params)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Service) GetIssueReportStatistics(ctx context.Context) (domain.IssueReportStatisticsResponse, error) {
	stats, err := s.Repo.GetIssueReportStatistics(ctx)
	if err != nil {
		return domain.IssueReportStatisticsResponse{}, err
	}

	// Convert to IssueReportStatisticsResponse using mapper
	return mapper.IssueReportStatisticsToResponse(&stats), nil
}
