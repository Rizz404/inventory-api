package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/oklog/ulid/v2"
)

// *==================== Model conversions ====================

func ToModelAssetMovement(d *domain.AssetMovement) model.AssetMovement {
	modelMovement := model.AssetMovement{
		MovementDate: d.MovementDate,
	}

	if d.ID != "" {
		if parsedID, err := ulid.Parse(d.ID); err == nil {
			modelMovement.ID = model.SQLULID(parsedID)
		}
	}

	if d.AssetID != "" {
		if parsedAssetID, err := ulid.Parse(d.AssetID); err == nil {
			modelMovement.AssetID = model.SQLULID(parsedAssetID)
		}
	}

	if d.FromLocationID != nil && *d.FromLocationID != "" {
		if parsedFromLocationID, err := ulid.Parse(*d.FromLocationID); err == nil {
			modelULID := model.SQLULID(parsedFromLocationID)
			modelMovement.FromLocationID = &modelULID
		}
	}

	if d.ToLocationID != nil && *d.ToLocationID != "" {
		if parsedToLocationID, err := ulid.Parse(*d.ToLocationID); err == nil {
			modelULID := model.SQLULID(parsedToLocationID)
			modelMovement.ToLocationID = &modelULID
		}
	}

	if d.FromUserID != nil && *d.FromUserID != "" {
		if parsedFromUserID, err := ulid.Parse(*d.FromUserID); err == nil {
			modelULID := model.SQLULID(parsedFromUserID)
			modelMovement.FromUserID = &modelULID
		}
	}

	if d.ToUserID != nil && *d.ToUserID != "" {
		if parsedToUserID, err := ulid.Parse(*d.ToUserID); err == nil {
			modelULID := model.SQLULID(parsedToUserID)
			modelMovement.ToUserID = &modelULID
		}
	}

	if d.MovedBy != "" {
		if parsedMovedBy, err := ulid.Parse(d.MovedBy); err == nil {
			modelMovement.MovedBy = model.SQLULID(parsedMovedBy)
		}
	}

	return modelMovement
}

func ToModelAssetMovementForCreate(d *domain.AssetMovement) model.AssetMovement {
	modelMovement := model.AssetMovement{
		MovementDate: d.MovementDate,
	}

	if d.AssetID != "" {
		if parsedAssetID, err := ulid.Parse(d.AssetID); err == nil {
			modelMovement.AssetID = model.SQLULID(parsedAssetID)
		}
	}

	if d.FromLocationID != nil && *d.FromLocationID != "" {
		if parsedFromLocationID, err := ulid.Parse(*d.FromLocationID); err == nil {
			modelULID := model.SQLULID(parsedFromLocationID)
			modelMovement.FromLocationID = &modelULID
		}
	}

	if d.ToLocationID != nil && *d.ToLocationID != "" {
		if parsedToLocationID, err := ulid.Parse(*d.ToLocationID); err == nil {
			modelULID := model.SQLULID(parsedToLocationID)
			modelMovement.ToLocationID = &modelULID
		}
	}

	if d.FromUserID != nil && *d.FromUserID != "" {
		if parsedFromUserID, err := ulid.Parse(*d.FromUserID); err == nil {
			modelULID := model.SQLULID(parsedFromUserID)
			modelMovement.FromUserID = &modelULID
		}
	}

	if d.ToUserID != nil && *d.ToUserID != "" {
		if parsedToUserID, err := ulid.Parse(*d.ToUserID); err == nil {
			modelULID := model.SQLULID(parsedToUserID)
			modelMovement.ToUserID = &modelULID
		}
	}

	if d.MovedBy != "" {
		if parsedMovedBy, err := ulid.Parse(d.MovedBy); err == nil {
			modelMovement.MovedBy = model.SQLULID(parsedMovedBy)
		}
	}

	return modelMovement
}

func ToModelAssetMovementTranslation(d *domain.AssetMovementTranslation) model.AssetMovementTranslation {
	modelTranslation := model.AssetMovementTranslation{
		LangCode: d.LangCode,
		Notes:    d.Notes,
	}

	if d.ID != "" {
		if parsedID, err := ulid.Parse(d.ID); err == nil {
			modelTranslation.ID = model.SQLULID(parsedID)
		}
	}

	if d.MovementID != "" {
		if parsedMovementID, err := ulid.Parse(d.MovementID); err == nil {
			modelTranslation.MovementID = model.SQLULID(parsedMovementID)
		}
	}

	return modelTranslation
}

func ToModelAssetMovementTranslationForCreate(movementID string, d *domain.AssetMovementTranslation) model.AssetMovementTranslation {
	modelTranslation := model.AssetMovementTranslation{
		LangCode: d.LangCode,
		Notes:    d.Notes,
	}

	if movementID != "" {
		if parsedMovementID, err := ulid.Parse(movementID); err == nil {
			modelTranslation.MovementID = model.SQLULID(parsedMovementID)
		}
	}

	return modelTranslation
}

func ToModelAssetMovementUpdateMap(payload *domain.UpdateAssetMovementPayload) map[string]interface{} {
	updates := make(map[string]interface{})

	if payload.ToLocationID != nil {
		if *payload.ToLocationID == "" {
			updates["to_location_id"] = nil
		} else {
			if parsedToLocationID, err := ulid.Parse(*payload.ToLocationID); err == nil {
				modelULID := model.SQLULID(parsedToLocationID)
				updates["to_location_id"] = &modelULID
			}
		}
	}

	if payload.ToUserID != nil {
		if *payload.ToUserID == "" {
			updates["to_user_id"] = nil
		} else {
			if parsedToUserID, err := ulid.Parse(*payload.ToUserID); err == nil {
				modelULID := model.SQLULID(parsedToUserID)
				updates["to_user_id"] = &modelULID
			}
		}
	}

	return updates
}

// *==================== Entity conversions ====================
func ToDomainAssetMovement(m *model.AssetMovement) domain.AssetMovement {
	domainMovement := domain.AssetMovement{
		ID:           m.ID.String(),
		AssetID:      m.AssetID.String(),
		MovementDate: m.MovementDate,
		MovedBy:      m.MovedBy.String(),
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}

	if m.FromLocationID != nil && !m.FromLocationID.IsZero() {
		fromLocationIDStr := m.FromLocationID.String()
		domainMovement.FromLocationID = &fromLocationIDStr
	}

	if m.ToLocationID != nil && !m.ToLocationID.IsZero() {
		toLocationIDStr := m.ToLocationID.String()
		domainMovement.ToLocationID = &toLocationIDStr
	}

	if m.FromUserID != nil && !m.FromUserID.IsZero() {
		fromUserIDStr := m.FromUserID.String()
		domainMovement.FromUserID = &fromUserIDStr
	}

	if m.ToUserID != nil && !m.ToUserID.IsZero() {
		toUserIDStr := m.ToUserID.String()
		domainMovement.ToUserID = &toUserIDStr
	}

	if len(m.Translations) > 0 {
		domainMovement.Translations = make([]domain.AssetMovementTranslation, len(m.Translations))
		for i, translation := range m.Translations {
			domainMovement.Translations[i] = ToDomainAssetMovementTranslation(&translation)
		}
	}

	return domainMovement
}

func ToDomainAssetMovements(models []model.AssetMovement) []domain.AssetMovement {
	movements := make([]domain.AssetMovement, len(models))
	for i, m := range models {
		movements[i] = ToDomainAssetMovement(&m)
	}
	return movements
}

func ToDomainAssetMovementTranslation(m *model.AssetMovementTranslation) domain.AssetMovementTranslation {
	return domain.AssetMovementTranslation{
		ID:         m.ID.String(),
		MovementID: m.MovementID.String(),
		LangCode:   m.LangCode,
		Notes:      m.Notes,
	}
}

// *==================== Entity Response conversions ====================
func AssetMovementToResponse(d *domain.AssetMovement, langCode string) domain.AssetMovementResponse {
	response := domain.AssetMovementResponse{
		ID:             d.ID,
		AssetID:        d.AssetID,
		FromLocationID: d.FromLocationID,
		ToLocationID:   d.ToLocationID,
		FromUserID:     d.FromUserID,
		ToUserID:       d.ToUserID,
		MovedByID:      d.MovedBy,
		MovementDate:   d.MovementDate.Format(TimeFormat),
		CreatedAt:      d.CreatedAt.Format(TimeFormat),
		UpdatedAt:      d.UpdatedAt.Format(TimeFormat),
		Translations:   make([]domain.AssetMovementTranslationResponse, len(d.Translations)),
	}

	// Populate translations
	for i, translation := range d.Translations {
		response.Translations[i] = domain.AssetMovementTranslationResponse{
			LangCode: translation.LangCode,
			Notes:    translation.Notes,
		}
	}

	// Find translation for the requested language
	for _, translation := range d.Translations {
		if translation.LangCode == langCode {
			response.Notes = translation.Notes
			break
		}
	}

	// If no translation found for requested language, use first available
	if response.Notes == nil && len(d.Translations) > 0 {
		response.Notes = d.Translations[0].Notes
	}

	return response
}

func AssetMovementsToResponses(movements []domain.AssetMovement, langCode string) []domain.AssetMovementResponse {
	responses := make([]domain.AssetMovementResponse, len(movements))
	for i, movement := range movements {
		responses[i] = AssetMovementToResponse(&movement, langCode)
	}
	return responses
}

func AssetMovementToListResponse(d *domain.AssetMovement, langCode string) domain.AssetMovementListResponse {
	response := domain.AssetMovementListResponse{
		ID:             d.ID,
		AssetID:        d.AssetID,
		FromLocationID: d.FromLocationID,
		ToLocationID:   d.ToLocationID,
		FromUserID:     d.FromUserID,
		ToUserID:       d.ToUserID,
		MovedByID:      d.MovedBy,
		MovementDate:   d.MovementDate.Format(TimeFormat),
		CreatedAt:      d.CreatedAt.Format(TimeFormat),
		UpdatedAt:      d.UpdatedAt.Format(TimeFormat),
	}

	// Find translation for the requested language
	for _, translation := range d.Translations {
		if translation.LangCode == langCode {
			response.Notes = translation.Notes
			break
		}
	}

	// If no translation found for requested language, use first available
	if response.Notes == nil && len(d.Translations) > 0 {
		response.Notes = d.Translations[0].Notes
	}

	return response
}

func AssetMovementsToListResponses(movements []domain.AssetMovement, langCode string) []domain.AssetMovementListResponse {
	responses := make([]domain.AssetMovementListResponse, len(movements))
	for i, movement := range movements {
		responses[i] = AssetMovementToListResponse(&movement, langCode)
	}
	return responses
}

// *==================== Statistics conversions ====================
func AssetMovementStatisticsToResponse(stats *domain.AssetMovementStatistics) domain.AssetMovementStatisticsResponse {
	return domain.AssetMovementStatisticsResponse{
		Total: domain.AssetMovementCountStatisticsResponse{
			Count: stats.Total.Count,
		},
		ByAsset:         AssetMovementByAssetStatsToResponses(stats.ByAsset),
		ByLocation:      AssetMovementByLocationStatsToResponses(stats.ByLocation),
		ByUser:          AssetMovementByUserStatsToResponses(stats.ByUser),
		ByMovementType:  AssetMovementTypeStatisticsToResponse(stats.ByMovementType),
		RecentMovements: AssetMovementRecentStatsToResponses(stats.RecentMovements),
		MovementTrends:  AssetMovementTrendsToResponses(stats.MovementTrends),
		Summary:         AssetMovementSummaryStatisticsToResponse(stats.Summary),
	}
}

func AssetMovementByAssetStatsToResponses(stats []domain.AssetMovementByAssetStats) []domain.AssetMovementByAssetStatsResponse {
	responses := make([]domain.AssetMovementByAssetStatsResponse, len(stats))
	for i, stat := range stats {
		responses[i] = domain.AssetMovementByAssetStatsResponse{
			AssetID:       stat.AssetID,
			AssetTag:      stat.AssetTag,
			AssetName:     stat.AssetName,
			MovementCount: stat.MovementCount,
		}
	}
	return responses
}

func AssetMovementByLocationStatsToResponses(stats []domain.AssetMovementByLocationStats) []domain.AssetMovementByLocationStatsResponse {
	responses := make([]domain.AssetMovementByLocationStatsResponse, len(stats))
	for i, stat := range stats {
		responses[i] = domain.AssetMovementByLocationStatsResponse{
			LocationID:    stat.LocationID,
			LocationCode:  stat.LocationCode,
			LocationName:  stat.LocationName,
			IncomingCount: stat.IncomingCount,
			OutgoingCount: stat.OutgoingCount,
			NetMovement:   stat.NetMovement,
		}
	}
	return responses
}

func AssetMovementByUserStatsToResponses(stats []domain.AssetMovementByUserStats) []domain.AssetMovementByUserStatsResponse {
	responses := make([]domain.AssetMovementByUserStatsResponse, len(stats))
	for i, stat := range stats {
		responses[i] = domain.AssetMovementByUserStatsResponse{
			UserID:        stat.UserID,
			UserName:      stat.UserName,
			MovementCount: stat.MovementCount,
		}
	}
	return responses
}

func AssetMovementTypeStatisticsToResponse(stats domain.AssetMovementTypeStatistics) domain.AssetMovementTypeStatisticsResponse {
	return domain.AssetMovementTypeStatisticsResponse{
		LocationToLocation: stats.LocationToLocation,
		LocationToUser:     stats.LocationToUser,
		UserToLocation:     stats.UserToLocation,
		UserToUser:         stats.UserToUser,
		NewAsset:           stats.NewAsset,
	}
}

func AssetMovementRecentStatsToResponses(stats []domain.AssetMovementRecentStats) []domain.AssetMovementRecentStatsResponse {
	responses := make([]domain.AssetMovementRecentStatsResponse, len(stats))
	for i, stat := range stats {
		responses[i] = domain.AssetMovementRecentStatsResponse{
			ID:           stat.ID,
			AssetTag:     stat.AssetTag,
			AssetName:    stat.AssetName,
			FromLocation: stat.FromLocation,
			ToLocation:   stat.ToLocation,
			FromUser:     stat.FromUser,
			ToUser:       stat.ToUser,
			MovedBy:      stat.MovedBy,
			MovementDate: stat.MovementDate,
			MovementType: stat.MovementType,
		}
	}
	return responses
}

func AssetMovementTrendsToResponses(trends []domain.AssetMovementTrend) []domain.AssetMovementTrendResponse {
	responses := make([]domain.AssetMovementTrendResponse, len(trends))
	for i, trend := range trends {
		responses[i] = domain.AssetMovementTrendResponse{
			Date:  trend.Date,
			Count: trend.Count,
		}
	}
	return responses
}

func AssetMovementSummaryStatisticsToResponse(summary domain.AssetMovementSummaryStatistics) domain.AssetMovementSummaryStatisticsResponse {
	return domain.AssetMovementSummaryStatisticsResponse{
		TotalMovements:            summary.TotalMovements,
		MovementsToday:            summary.MovementsToday,
		MovementsThisWeek:         summary.MovementsThisWeek,
		MovementsThisMonth:        summary.MovementsThisMonth,
		MostActiveAsset:           summary.MostActiveAsset,
		MostActiveLocation:        summary.MostActiveLocation,
		MostActiveUser:            summary.MostActiveUser,
		AverageMovementsPerDay:    summary.AverageMovementsPerDay,
		AverageMovementsPerAsset:  summary.AverageMovementsPerAsset,
		LatestMovementDate:        summary.LatestMovementDate,
		EarliestMovementDate:      summary.EarliestMovementDate,
		UniqueAssetsWithMovements: summary.UniqueAssetsWithMovements,
		UniqueLocationsInvolved:   summary.UniqueLocationsInvolved,
		UniqueUsersInvolved:       summary.UniqueUsersInvolved,
	}
}
