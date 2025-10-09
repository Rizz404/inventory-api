package postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"gorm.io/gorm"
)

type AssetMovementRepository struct {
	db *gorm.DB
}

func NewAssetMovementRepository(db *gorm.DB) *AssetMovementRepository {
	return &AssetMovementRepository{
		db: db,
	}
}

func (r *AssetMovementRepository) applyAssetMovementFilters(db *gorm.DB, filters *domain.AssetMovementFilterOptions) *gorm.DB {
	if filters == nil {
		return db
	}

	if filters.AssetID != nil && *filters.AssetID != "" {
		db = db.Where("am.asset_id = ?", *filters.AssetID)
	}

	if filters.FromLocationID != nil && *filters.FromLocationID != "" {
		db = db.Where("am.from_location_id = ?", *filters.FromLocationID)
	}

	if filters.ToLocationID != nil && *filters.ToLocationID != "" {
		db = db.Where("am.to_location_id = ?", *filters.ToLocationID)
	}

	if filters.FromUserID != nil && *filters.FromUserID != "" {
		db = db.Where("am.from_user_id = ?", *filters.FromUserID)
	}

	if filters.ToUserID != nil && *filters.ToUserID != "" {
		db = db.Where("am.to_user_id = ?", *filters.ToUserID)
	}

	if filters.MovedBy != nil && *filters.MovedBy != "" {
		db = db.Where("am.moved_by = ?", *filters.MovedBy)
	}

	if filters.DateFrom != nil {
		db = db.Where("am.movement_date >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		db = db.Where("am.movement_date <= ?", *filters.DateTo)
	}

	return db
}

func (r *AssetMovementRepository) applyAssetMovementSorts(db *gorm.DB, sort *domain.AssetMovementSortOptions) *gorm.DB {
	if sort == nil || sort.Field == "" {
		return db.Order("am.movement_date DESC")
	}

	// Map camelCase sort field to snake_case database column
	columnName := mapper.MapAssetMovementSortFieldToColumn(sort.Field)
	orderClause := columnName

	order := "DESC"
	if sort.Order == domain.SortOrderAsc {
		order = "ASC"
	}
	return db.Order(fmt.Sprintf("%s %s", orderClause, order))
}

// *===========================MUTATION===========================*
func (r *AssetMovementRepository) CreateAssetMovement(ctx context.Context, payload *domain.AssetMovement) (domain.AssetMovement, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.AssetMovement{}, tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// Create asset movement
	modelMovement := mapper.ToModelAssetMovementForCreate(payload)
	if err := tx.Create(&modelMovement).Error; err != nil {
		tx.Rollback()
		return domain.AssetMovement{}, err
	}

	// Create translations
	for _, translation := range payload.Translations {
		modelTranslation := mapper.ToModelAssetMovementTranslationForCreate(modelMovement.ID.String(), &translation)
		if err := tx.Create(&modelTranslation).Error; err != nil {
			tx.Rollback()
			return domain.AssetMovement{}, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return domain.AssetMovement{}, err
	}

	// Return created asset movement with translations (no need to query again)
	// GORM has already filled the model with created data including ID and timestamps
	domainMovement := mapper.ToDomainAssetMovement(&modelMovement)
	// Add translations manually since they were created separately
	for _, translation := range payload.Translations {
		domainMovement.Translations = append(domainMovement.Translations, domain.AssetMovementTranslation{
			LangCode: translation.LangCode,
			Notes:    translation.Notes,
		})
	}
	return domainMovement, nil
}

func (r *AssetMovementRepository) UpdateAssetMovement(ctx context.Context, movementId string, payload *domain.UpdateAssetMovementPayload) (domain.AssetMovement, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.AssetMovement{}, domain.ErrInternal(tx.Error)
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// Update asset movement basic info
	updates := mapper.ToModelAssetMovementUpdateMap(payload)
	if len(updates) > 0 {
		if err := tx.Table("asset_movements").Where("id = ?", movementId).Updates(updates).Error; err != nil {
			tx.Rollback()
			return domain.AssetMovement{}, domain.ErrInternal(err)
		}
	}

	// Update translations if provided
	if len(payload.Translations) > 0 {
		// Delete existing translations
		if err := tx.Where("movement_id = ?", movementId).Delete(&model.AssetMovementTranslation{}).Error; err != nil {
			tx.Rollback()
			return domain.AssetMovement{}, domain.ErrInternal(err)
		}

		// Create new translations
		for _, translationPayload := range payload.Translations {
			translation := domain.AssetMovementTranslation{
				LangCode: translationPayload.LangCode,
				Notes:    translationPayload.Notes,
			}
			modelTranslation := mapper.ToModelAssetMovementTranslationForCreate(movementId, &translation)
			if err := tx.Create(&modelTranslation).Error; err != nil {
				tx.Rollback()
				return domain.AssetMovement{}, domain.ErrInternal(err)
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return domain.AssetMovement{}, domain.ErrInternal(err)
	}

	// Fetch updated asset movement with translations
	return r.GetAssetMovementById(ctx, movementId)
}

func (r *AssetMovementRepository) DeleteAssetMovement(ctx context.Context, movementId string) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// Delete translations first (foreign key constraint)
	if err := tx.Where("movement_id = ?", movementId).Delete(&model.AssetMovementTranslation{}).Error; err != nil {
		tx.Rollback()
		return domain.ErrInternal(err)
	}

	// Delete asset movement
	result := tx.Where("id = ?", movementId).Delete(&model.AssetMovement{})
	if result.Error != nil {
		tx.Rollback()
		return domain.ErrInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return domain.ErrNotFound("asset movement")
	}

	if err := tx.Commit().Error; err != nil {
		return domain.ErrInternal(err)
	}
	return nil
}

// *===========================QUERY===========================*
func (r *AssetMovementRepository) GetAssetMovementsPaginated(ctx context.Context, params domain.AssetMovementParams, langCode string) ([]domain.AssetMovement, error) {
	var movements []model.AssetMovement
	db := r.db.WithContext(ctx).
		Table("asset_movements am").
		Preload("Translations").
		Preload("Asset").
		Preload("Asset.Category").
		Preload("Asset.Category.Translations").
		Preload("Asset.Location").
		Preload("Asset.Location.Translations").
		Preload("Asset.User").
		Preload("FromLocation").
		Preload("FromLocation.Translations").
		Preload("ToLocation").
		Preload("ToLocation.Translations").
		Preload("FromUser").
		Preload("ToUser").
		Preload("MovedByUser")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		// Join with assets table for search in asset tag/name
		db = db.Joins("LEFT JOIN assets a ON am.asset_id = a.id").
			Where("a.asset_tag ILIKE ? OR a.serial_number ILIKE ?",
				"%"+*params.SearchQuery+"%", "%"+*params.SearchQuery+"%")
	}

	// Apply filters, sorts, and pagination manually
	db = r.applyAssetMovementFilters(db, params.Filters)
	db = r.applyAssetMovementSorts(db, params.Sort)
	if params.Pagination != nil {
		if params.Pagination.Limit > 0 {
			db = db.Limit(params.Pagination.Limit)
		}
		if params.Pagination.Offset > 0 {
			db = db.Offset(params.Pagination.Offset)
		}
	}

	if err := db.Find(&movements).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain asset movements
	return mapper.ToDomainAssetMovements(movements), nil
}

func (r *AssetMovementRepository) GetAssetMovementsCursor(ctx context.Context, params domain.AssetMovementParams, langCode string) ([]domain.AssetMovement, error) {
	var movements []model.AssetMovement
	db := r.db.WithContext(ctx).
		Table("asset_movements am").
		Preload("Translations").
		Preload("Asset").
		Preload("Asset.Category").
		Preload("Asset.Category.Translations").
		Preload("Asset.Location").
		Preload("Asset.Location.Translations").
		Preload("Asset.User").
		Preload("FromLocation").
		Preload("FromLocation.Translations").
		Preload("ToLocation").
		Preload("ToLocation.Translations").
		Preload("FromUser").
		Preload("ToUser").
		Preload("MovedByUser")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		// Join with assets table for search in asset tag/name
		db = db.Joins("LEFT JOIN assets a ON am.asset_id = a.id").
			Where("a.asset_tag ILIKE ? OR a.serial_number ILIKE ?",
				"%"+*params.SearchQuery+"%", "%"+*params.SearchQuery+"%")
	}

	// Apply filters, sorts, and cursor pagination manually
	db = r.applyAssetMovementFilters(db, params.Filters)
	db = r.applyAssetMovementSorts(db, params.Sort)
	if params.Pagination != nil {
		if params.Pagination.Cursor != "" {
			db = db.Where("am.id > ?", params.Pagination.Cursor)
		}
		if params.Pagination.Limit > 0 {
			db = db.Limit(params.Pagination.Limit)
		}
	}

	if err := db.Find(&movements).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain asset movements
	return mapper.ToDomainAssetMovements(movements), nil
}

func (r *AssetMovementRepository) GetAssetMovementById(ctx context.Context, movementId string) (domain.AssetMovement, error) {
	var movement model.AssetMovement

	err := r.db.WithContext(ctx).
		Table("asset_movements am").
		Preload("Translations").
		Preload("Asset").
		Preload("Asset.Category").
		Preload("Asset.Category.Translations").
		Preload("Asset.Location").
		Preload("Asset.Location.Translations").
		Preload("Asset.User").
		Preload("FromLocation").
		Preload("FromLocation.Translations").
		Preload("ToLocation").
		Preload("ToLocation.Translations").
		Preload("FromUser").
		Preload("ToUser").
		Preload("MovedByUser").
		First(&movement, "am.id = ?", movementId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.AssetMovement{}, domain.ErrNotFound("asset movement")
		}
		return domain.AssetMovement{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainAssetMovement(&movement), nil
}

func (r *AssetMovementRepository) GetAssetMovementsByAssetId(ctx context.Context, assetId string, params domain.AssetMovementParams) ([]domain.AssetMovement, error) {
	var movements []model.AssetMovement
	db := r.db.WithContext(ctx).
		Table("asset_movements am").
		Preload("Translations").
		Preload("Asset").
		Preload("FromLocation").
		Preload("ToLocation").
		Preload("FromUser").
		Preload("ToUser").
		Preload("MovedByUser").
		Where("am.asset_id = ?", assetId)

	db = r.applyAssetMovementFilters(db, params.Filters)
	db = r.applyAssetMovementSorts(db, params.Sort)
	if params.Pagination != nil {
		if params.Pagination.Limit > 0 {
			db = db.Limit(params.Pagination.Limit)
		}
		if params.Pagination.Offset > 0 {
			db = db.Offset(params.Pagination.Offset)
		}
	}

	if err := db.Find(&movements).Error; err != nil {
		return nil, err
	}

	return mapper.ToDomainAssetMovements(movements), nil
}

func (r *AssetMovementRepository) CheckAssetMovementExist(ctx context.Context, movementId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("asset_movements am").Where("am.id = ?", movementId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *AssetMovementRepository) CountAssetMovements(ctx context.Context, params domain.AssetMovementParams) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Table("asset_movements am")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		// Join with assets table for search in asset tag/name
		db = db.Joins("LEFT JOIN assets a ON am.asset_id = a.id").
			Where("a.asset_tag ILIKE ? OR a.serial_number ILIKE ?",
				"%"+*params.SearchQuery+"%", "%"+*params.SearchQuery+"%")
	}

	db = r.applyAssetMovementFilters(db, params.Filters)

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *AssetMovementRepository) GetAssetMovementStatistics(ctx context.Context) (domain.AssetMovementStatistics, error) {
	var stats domain.AssetMovementStatistics

	// Get total asset movement count
	var totalCount int64
	if err := r.db.WithContext(ctx).Table("asset_movements").Count(&totalCount).Error; err != nil {
		return stats, err
	}
	stats.Total.Count = int(totalCount)

	// Get asset movement counts by asset (top 10)
	var assetStats []struct {
		AssetID       string `json:"asset_id"`
		AssetTag      string `json:"asset_tag"`
		AssetName     string `json:"asset_name"`
		MovementCount int64  `json:"movement_count"`
	}
	if err := r.db.WithContext(ctx).
		Table("asset_movements am").
		Select("am.asset_id, a.asset_tag, COALESCE(at.asset_name, '') as asset_name, COUNT(*) as movement_count").
		Joins("LEFT JOIN assets a ON am.asset_id = a.id").
		Joins("LEFT JOIN asset_translations at ON a.id = at.asset_id AND at.lang_code = 'en-US'").
		Group("am.asset_id, a.asset_tag, at.asset_name").
		Order("movement_count DESC").
		Limit(10).
		Find(&assetStats).Error; err != nil {
		return stats, err
	}

	stats.ByAsset = make([]domain.AssetMovementByAssetStats, len(assetStats))
	for i, as := range assetStats {
		stats.ByAsset[i] = domain.AssetMovementByAssetStats{
			AssetID:       as.AssetID,
			AssetTag:      as.AssetTag,
			AssetName:     as.AssetName,
			MovementCount: int(as.MovementCount),
		}
	}

	// Get asset movement counts by location (top 10)
	var locationStats []struct {
		LocationID    string `json:"location_id"`
		LocationCode  string `json:"location_code"`
		LocationName  string `json:"location_name"`
		IncomingCount int64  `json:"incoming_count"`
		OutgoingCount int64  `json:"outgoing_count"`
	}
	if err := r.db.WithContext(ctx).Raw(`
		SELECT
			l.id as location_id,
			l.location_code,
			COALESCE(lt.location_name, '') as location_name,
			COALESCE(incoming.count, 0) as incoming_count,
			COALESCE(outgoing.count, 0) as outgoing_count
		FROM locations l
		LEFT JOIN location_translations lt ON l.id = lt.location_id AND lt.lang_code = 'en-US'
		LEFT JOIN (
			SELECT to_location_id as location_id, COUNT(*) as count
			FROM asset_movements
			WHERE to_location_id IS NOT NULL
			GROUP BY to_location_id
		) incoming ON l.id = incoming.location_id
		LEFT JOIN (
			SELECT from_location_id as location_id, COUNT(*) as count
			FROM asset_movements
			WHERE from_location_id IS NOT NULL
			GROUP BY from_location_id
		) outgoing ON l.id = outgoing.location_id
		WHERE incoming.count > 0 OR outgoing.count > 0
		ORDER BY (COALESCE(incoming.count, 0) + COALESCE(outgoing.count, 0)) DESC
		LIMIT 10
	`).Find(&locationStats).Error; err != nil {
		return stats, err
	}

	stats.ByLocation = make([]domain.AssetMovementByLocationStats, len(locationStats))
	for i, ls := range locationStats {
		stats.ByLocation[i] = domain.AssetMovementByLocationStats{
			LocationID:    ls.LocationID,
			LocationCode:  ls.LocationCode,
			LocationName:  ls.LocationName,
			IncomingCount: int(ls.IncomingCount),
			OutgoingCount: int(ls.OutgoingCount),
			NetMovement:   int(ls.IncomingCount - ls.OutgoingCount),
		}
	}

	// Get asset movement counts by user (top 10)
	var userStats []struct {
		UserID        string `json:"user_id"`
		UserName      string `json:"user_name"`
		MovementCount int64  `json:"movement_count"`
	}
	if err := r.db.WithContext(ctx).
		Table("asset_movements am").
		Select("am.moved_by as user_id, u.name as user_name, COUNT(*) as movement_count").
		Joins("LEFT JOIN users u ON am.moved_by = u.id").
		Group("am.moved_by, u.name").
		Order("movement_count DESC").
		Limit(10).
		Find(&userStats).Error; err != nil {
		return stats, err
	}

	stats.ByUser = make([]domain.AssetMovementByUserStats, len(userStats))
	for i, us := range userStats {
		stats.ByUser[i] = domain.AssetMovementByUserStats{
			UserID:        us.UserID,
			UserName:      us.UserName,
			MovementCount: int(us.MovementCount),
		}
	}

	// Get movement type statistics
	var typeStats struct {
		LocationToLocation int64 `json:"location_to_location"`
		LocationToUser     int64 `json:"location_to_user"`
		UserToLocation     int64 `json:"user_to_location"`
		UserToUser         int64 `json:"user_to_user"`
		NewAsset           int64 `json:"new_asset"`
	}

	// Location to Location
	r.db.WithContext(ctx).Table("asset_movements").
		Where("from_location_id IS NOT NULL AND to_location_id IS NOT NULL").
		Count(&typeStats.LocationToLocation)

	// Location to User
	r.db.WithContext(ctx).Table("asset_movements").
		Where("from_location_id IS NOT NULL AND to_user_id IS NOT NULL").
		Count(&typeStats.LocationToUser)

	// User to Location
	r.db.WithContext(ctx).Table("asset_movements").
		Where("from_user_id IS NOT NULL AND to_location_id IS NOT NULL").
		Count(&typeStats.UserToLocation)

	// User to User
	r.db.WithContext(ctx).Table("asset_movements").
		Where("from_user_id IS NOT NULL AND to_user_id IS NOT NULL").
		Count(&typeStats.UserToUser)

	// New Asset (no from location or user)
	r.db.WithContext(ctx).Table("asset_movements").
		Where("from_location_id IS NULL AND from_user_id IS NULL").
		Count(&typeStats.NewAsset)

	stats.ByMovementType = domain.AssetMovementTypeStatistics{
		LocationToLocation: int(typeStats.LocationToLocation),
		LocationToUser:     int(typeStats.LocationToUser),
		UserToLocation:     int(typeStats.UserToLocation),
		UserToUser:         int(typeStats.UserToUser),
		NewAsset:           int(typeStats.NewAsset),
	}

	// Get recent movements (last 10)
	var recentMovements []struct {
		ID           string    `json:"id"`
		AssetTag     string    `json:"asset_tag"`
		AssetName    string    `json:"asset_name"`
		FromLocation *string   `json:"from_location"`
		ToLocation   *string   `json:"to_location"`
		FromUser     *string   `json:"from_user"`
		ToUser       *string   `json:"to_user"`
		MovedBy      string    `json:"moved_by"`
		MovementDate time.Time `json:"movement_date"`
	}
	if err := r.db.WithContext(ctx).Raw(`
		SELECT
			am.id,
			a.asset_tag,
			COALESCE(at.asset_name, '') as asset_name,
			COALESCE(fl.location_code, '') as from_location,
			COALESCE(tl.location_code, '') as to_location,
			COALESCE(fu.name, '') as from_user,
			COALESCE(tu.name, '') as to_user,
			mb.name as moved_by,
			am.movement_date
		FROM asset_movements am
		LEFT JOIN assets a ON am.asset_id = a.id
		LEFT JOIN asset_translations at ON a.id = at.asset_id AND at.lang_code = 'en-US'
		LEFT JOIN locations fl ON am.from_location_id = fl.id
		LEFT JOIN locations tl ON am.to_location_id = tl.id
		LEFT JOIN users fu ON am.from_user_id = fu.id
		LEFT JOIN users tu ON am.to_user_id = tu.id
		LEFT JOIN users mb ON am.moved_by = mb.id
		ORDER BY am.movement_date DESC
		LIMIT 10
	`).Find(&recentMovements).Error; err != nil {
		return stats, err
	}

	stats.RecentMovements = make([]domain.AssetMovementRecentStats, len(recentMovements))
	for i, rm := range recentMovements {
		movementType := "New Asset"
		if rm.FromLocation != nil && *rm.FromLocation != "" && rm.ToLocation != nil && *rm.ToLocation != "" {
			movementType = "Location to Location"
		} else if rm.FromLocation != nil && *rm.FromLocation != "" && rm.ToUser != nil && *rm.ToUser != "" {
			movementType = "Location to User"
		} else if rm.FromUser != nil && *rm.FromUser != "" && rm.ToLocation != nil && *rm.ToLocation != "" {
			movementType = "User to Location"
		} else if rm.FromUser != nil && *rm.FromUser != "" && rm.ToUser != nil && *rm.ToUser != "" {
			movementType = "User to User"
		}

		stats.RecentMovements[i] = domain.AssetMovementRecentStats{
			ID:           rm.ID,
			AssetTag:     rm.AssetTag,
			AssetName:    rm.AssetName,
			FromLocation: rm.FromLocation,
			ToLocation:   rm.ToLocation,
			FromUser:     rm.FromUser,
			ToUser:       rm.ToUser,
			MovedBy:      rm.MovedBy,
			MovementDate: rm.MovementDate,
			MovementType: movementType,
		}
	}

	// Get movement trends (last 30 days)
	var movementTrends []struct {
		Date  time.Time `json:"date"`
		Count int64     `json:"count"`
	}
	if err := r.db.WithContext(ctx).Raw(`
		SELECT
			DATE(movement_date) as date,
			COUNT(*) as count
		FROM asset_movements
		WHERE movement_date >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY DATE(movement_date)
		ORDER BY date DESC
	`).Find(&movementTrends).Error; err != nil {
		return stats, err
	}

	stats.MovementTrends = make([]domain.AssetMovementTrend, len(movementTrends))
	for i, mt := range movementTrends {
		stats.MovementTrends[i] = domain.AssetMovementTrend{
			Date:  mt.Date,
			Count: int(mt.Count),
		}
	}

	// Calculate summary statistics
	stats.Summary.TotalMovements = int(totalCount)

	// Get movements today
	var movementsToday int64
	r.db.WithContext(ctx).Table("asset_movements").
		Where("DATE(movement_date) = CURRENT_DATE").
		Count(&movementsToday)
	stats.Summary.MovementsToday = int(movementsToday)

	// Get movements this week
	var movementsThisWeek int64
	r.db.WithContext(ctx).Table("asset_movements").
		Where("movement_date >= DATE_TRUNC('week', CURRENT_DATE)").
		Count(&movementsThisWeek)
	stats.Summary.MovementsThisWeek = int(movementsThisWeek)

	// Get movements this month
	var movementsThisMonth int64
	r.db.WithContext(ctx).Table("asset_movements").
		Where("movement_date >= DATE_TRUNC('month', CURRENT_DATE)").
		Count(&movementsThisMonth)
	stats.Summary.MovementsThisMonth = int(movementsThisMonth)

	// Get most active asset, location, user
	if len(stats.ByAsset) > 0 {
		stats.Summary.MostActiveAsset = stats.ByAsset[0].AssetTag
	}
	if len(stats.ByLocation) > 0 {
		stats.Summary.MostActiveLocation = stats.ByLocation[0].LocationCode
	}
	if len(stats.ByUser) > 0 {
		stats.Summary.MostActiveUser = stats.ByUser[0].UserName
	}

	// Get unique counts
	var uniqueAssets, uniqueLocations, uniqueUsers int64
	r.db.WithContext(ctx).Table("asset_movements").Select("COUNT(DISTINCT asset_id)").Row().Scan(&uniqueAssets)
	r.db.WithContext(ctx).Raw("SELECT COUNT(DISTINCT location_id) FROM (SELECT from_location_id as location_id FROM asset_movements WHERE from_location_id IS NOT NULL UNION SELECT to_location_id as location_id FROM asset_movements WHERE to_location_id IS NOT NULL) t").Row().Scan(&uniqueLocations)
	r.db.WithContext(ctx).Raw("SELECT COUNT(DISTINCT user_id) FROM (SELECT from_user_id as user_id FROM asset_movements WHERE from_user_id IS NOT NULL UNION SELECT to_user_id as user_id FROM asset_movements WHERE to_user_id IS NOT NULL UNION SELECT moved_by as user_id FROM asset_movements) t").Row().Scan(&uniqueUsers)

	stats.Summary.UniqueAssetsWithMovements = int(uniqueAssets)
	stats.Summary.UniqueLocationsInvolved = int(uniqueLocations)
	stats.Summary.UniqueUsersInvolved = int(uniqueUsers)

	// Get earliest and latest movement dates
	var earliestDate, latestDate *time.Time
	r.db.WithContext(ctx).Table("asset_movements").Select("MIN(movement_date)").Row().Scan(&earliestDate)
	r.db.WithContext(ctx).Table("asset_movements").Select("MAX(movement_date)").Row().Scan(&latestDate)

	if earliestDate != nil {
		stats.Summary.EarliestMovementDate = *earliestDate
	}
	if latestDate != nil {
		stats.Summary.LatestMovementDate = *latestDate
	}

	// Calculate average movements per day and per asset
	if earliestDate != nil && latestDate != nil {
		daysDiff := latestDate.Sub(*earliestDate).Hours() / 24
		if daysDiff > 0 {
			stats.Summary.AverageMovementsPerDay = float64(totalCount) / daysDiff
		}
	}

	if uniqueAssets > 0 {
		stats.Summary.AverageMovementsPerAsset = float64(totalCount) / float64(uniqueAssets)
	}

	return stats, nil
}
