package domain

import (
	"time"
)

// --- Structs ---

type AssetMovement struct {
	ID             string                     `json:"id"`
	AssetID        string                     `json:"assetId"`
	FromLocationID *string                    `json:"fromLocationId"`
	ToLocationID   *string                    `json:"toLocationId"`
	FromUserID     *string                    `json:"fromUserId"`
	ToUserID       *string                    `json:"toUserId"`
	MovementDate   time.Time                  `json:"movementDate"`
	MovedBy        string                     `json:"movedBy"`
	CreatedAt      time.Time                  `json:"createdAt"`
	UpdatedAt      time.Time                  `json:"updatedAt"`
	Translations   []AssetMovementTranslation `json:"translations,omitempty"`
}

type AssetMovementTranslation struct {
	ID         string  `json:"id"`
	MovementID string  `json:"movementId"`
	LangCode   string  `json:"langCode"`
	Notes      *string `json:"notes"`
}

type AssetMovementTranslationResponse struct {
	LangCode string  `json:"langCode"`
	Notes    *string `json:"notes"`
}

type AssetMovementResponse struct {
	ID             string                             `json:"id"`
	AssetID        string                             `json:"assetId"`
	FromLocationID *string                            `json:"fromLocationId,omitempty"`
	ToLocationID   *string                            `json:"toLocationId,omitempty"`
	FromUserID     *string                            `json:"fromUserId,omitempty"`
	ToUserID       *string                            `json:"toUserId,omitempty"`
	MovedByID      string                             `json:"movedById"`
	MovementDate   time.Time                          `json:"movementDate"`
	Notes          *string                            `json:"notes,omitempty"`
	CreatedAt      time.Time                          `json:"createdAt"`
	UpdatedAt      time.Time                          `json:"updatedAt"`
	Translations   []AssetMovementTranslationResponse `json:"translations"`
	// * Populated
	Asset        AssetResponse     `json:"asset"`
	FromLocation *LocationResponse `json:"fromLocation,omitempty"`
	ToLocation   *LocationResponse `json:"toLocation,omitempty"`
	FromUser     *UserResponse     `json:"fromUser,omitempty"`
	ToUser       *UserResponse     `json:"toUser,omitempty"`
	MovedBy      UserResponse      `json:"movedBy"`
}

type AssetMovementListResponse struct {
	ID             string    `json:"id"`
	AssetID        string    `json:"assetId"`
	FromLocationID *string   `json:"fromLocationId,omitempty"`
	ToLocationID   *string   `json:"toLocationId,omitempty"`
	FromUserID     *string   `json:"fromUserId,omitempty"`
	ToUserID       *string   `json:"toUserId,omitempty"`
	MovedByID      string    `json:"movedById"`
	MovementDate   time.Time `json:"movementDate"`
	Notes          *string   `json:"notes,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	// * Populated
	Asset        AssetResponse     `json:"asset"`
	FromLocation *LocationResponse `json:"fromLocation,omitempty"`
	ToLocation   *LocationResponse `json:"toLocation,omitempty"`
	FromUser     *UserResponse     `json:"fromUser,omitempty"`
	ToUser       *UserResponse     `json:"toUser,omitempty"`
	MovedBy      UserResponse      `json:"movedBy"`
}

// --- Payloads ---

type CreateAssetMovementPayload struct {
	AssetID      string                                  `json:"assetId" validate:"required"`
	ToLocationID *string                                 `json:"toLocationId,omitempty" validate:"omitempty"`
	ToUserID     *string                                 `json:"toUserId,omitempty" validate:"omitempty"`
	Translations []CreateAssetMovementTranslationPayload `json:"translations,omitempty" validate:"omitempty,dive"`
}

type CreateAssetMovementTranslationPayload struct {
	LangCode string `json:"langCode" validate:"required,max=5"`
	Notes    string `json:"notes" validate:"required"`
}

type UpdateAssetMovementPayload struct {
	ToLocationID *string                                 `json:"toLocationId,omitempty" validate:"omitempty"`
	ToUserID     *string                                 `json:"toUserId,omitempty" validate:"omitempty"`
	Translations []UpdateAssetMovementTranslationPayload `json:"translations,omitempty" validate:"omitempty,dive"`
}

type UpdateAssetMovementTranslationPayload struct {
	LangCode string  `json:"langCode" validate:"required,max=5"`
	Notes    *string `json:"notes,omitempty" validate:"omitempty"`
}

// --- Statistics ---

// Internal statistics structs (used in repository layer)
type AssetMovementStatistics struct {
	Total           AssetMovementCountStatistics   `json:"total"`
	ByAsset         []AssetMovementByAssetStats    `json:"byAsset"`
	ByLocation      []AssetMovementByLocationStats `json:"byLocation"`
	ByUser          []AssetMovementByUserStats     `json:"byUser"`
	ByMovementType  AssetMovementTypeStatistics    `json:"byMovementType"`
	RecentMovements []AssetMovementRecentStats     `json:"recentMovements"`
	MovementTrends  []AssetMovementTrend           `json:"movementTrends"`
	Summary         AssetMovementSummaryStatistics `json:"summary"`
}

type AssetMovementCountStatistics struct {
	Count int `json:"count"`
}

type AssetMovementByAssetStats struct {
	AssetID       string `json:"assetId"`
	AssetTag      string `json:"assetTag"`
	AssetName     string `json:"assetName"`
	MovementCount int    `json:"movementCount"`
}

type AssetMovementByLocationStats struct {
	LocationID    string `json:"locationId"`
	LocationCode  string `json:"locationCode"`
	LocationName  string `json:"locationName"`
	IncomingCount int    `json:"incomingCount"`
	OutgoingCount int    `json:"outgoingCount"`
	NetMovement   int    `json:"netMovement"`
}

type AssetMovementByUserStats struct {
	UserID        string `json:"userId"`
	UserName      string `json:"userName"`
	MovementCount int    `json:"movementCount"`
}

type AssetMovementTypeStatistics struct {
	LocationToLocation int `json:"locationToLocation"`
	LocationToUser     int `json:"locationToUser"`
	UserToLocation     int `json:"userToLocation"`
	UserToUser         int `json:"userToUser"`
	NewAsset           int `json:"newAsset"`
}

type AssetMovementRecentStats struct {
	ID           string  `json:"id"`
	AssetTag     string  `json:"assetTag"`
	AssetName    string  `json:"assetName"`
	FromLocation *string `json:"fromLocation"`
	ToLocation   *string `json:"toLocation"`
	FromUser     *string `json:"fromUser"`
	ToUser       *string `json:"toUser"`
	MovedBy      string  `json:"movedBy"`
	MovementDate string  `json:"movementDate"`
	MovementType string  `json:"movementType"`
}

type AssetMovementTrend struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type AssetMovementSummaryStatistics struct {
	TotalMovements            int     `json:"totalMovements"`
	MovementsToday            int     `json:"movementsToday"`
	MovementsThisWeek         int     `json:"movementsThisWeek"`
	MovementsThisMonth        int     `json:"movementsThisMonth"`
	MostActiveAsset           string  `json:"mostActiveAsset"`
	MostActiveLocation        string  `json:"mostActiveLocation"`
	MostActiveUser            string  `json:"mostActiveUser"`
	AverageMovementsPerDay    float64 `json:"averageMovementsPerDay"`
	AverageMovementsPerAsset  float64 `json:"averageMovementsPerAsset"`
	LatestMovementDate        string  `json:"latestMovementDate"`
	EarliestMovementDate      string  `json:"earliestMovementDate"`
	UniqueAssetsWithMovements int     `json:"uniqueAssetsWithMovements"`
	UniqueLocationsInvolved   int     `json:"uniqueLocationsInvolved"`
	UniqueUsersInvolved       int     `json:"uniqueUsersInvolved"`
}

// Response statistics structs (used in service/handler layer)
type AssetMovementStatisticsResponse struct {
	Total           AssetMovementCountStatisticsResponse   `json:"total"`
	ByAsset         []AssetMovementByAssetStatsResponse    `json:"byAsset"`
	ByLocation      []AssetMovementByLocationStatsResponse `json:"byLocation"`
	ByUser          []AssetMovementByUserStatsResponse     `json:"byUser"`
	ByMovementType  AssetMovementTypeStatisticsResponse    `json:"byMovementType"`
	RecentMovements []AssetMovementRecentStatsResponse     `json:"recentMovements"`
	MovementTrends  []AssetMovementTrendResponse           `json:"movementTrends"`
	Summary         AssetMovementSummaryStatisticsResponse `json:"summary"`
}

type AssetMovementCountStatisticsResponse struct {
	Count int `json:"count"`
}

type AssetMovementByAssetStatsResponse struct {
	AssetID       string `json:"assetId"`
	AssetTag      string `json:"assetTag"`
	AssetName     string `json:"assetName"`
	MovementCount int    `json:"movementCount"`
}

type AssetMovementByLocationStatsResponse struct {
	LocationID    string `json:"locationId"`
	LocationCode  string `json:"locationCode"`
	LocationName  string `json:"locationName"`
	IncomingCount int    `json:"incomingCount"`
	OutgoingCount int    `json:"outgoingCount"`
	NetMovement   int    `json:"netMovement"`
}

type AssetMovementByUserStatsResponse struct {
	UserID        string `json:"userId"`
	UserName      string `json:"userName"`
	MovementCount int    `json:"movementCount"`
}

type AssetMovementTypeStatisticsResponse struct {
	LocationToLocation int `json:"locationToLocation"`
	LocationToUser     int `json:"locationToUser"`
	UserToLocation     int `json:"userToLocation"`
	UserToUser         int `json:"userToUser"`
	NewAsset           int `json:"newAsset"`
}

type AssetMovementRecentStatsResponse struct {
	ID           string  `json:"id"`
	AssetTag     string  `json:"assetTag"`
	AssetName    string  `json:"assetName"`
	FromLocation *string `json:"fromLocation"`
	ToLocation   *string `json:"toLocation"`
	FromUser     *string `json:"fromUser"`
	ToUser       *string `json:"toUser"`
	MovedBy      string  `json:"movedBy"`
	MovementDate string  `json:"movementDate"`
	MovementType string  `json:"movementType"`
}

type AssetMovementTrendResponse struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type AssetMovementSummaryStatisticsResponse struct {
	TotalMovements            int     `json:"totalMovements"`
	MovementsToday            int     `json:"movementsToday"`
	MovementsThisWeek         int     `json:"movementsThisWeek"`
	MovementsThisMonth        int     `json:"movementsThisMonth"`
	MostActiveAsset           string  `json:"mostActiveAsset"`
	MostActiveLocation        string  `json:"mostActiveLocation"`
	MostActiveUser            string  `json:"mostActiveUser"`
	AverageMovementsPerDay    float64 `json:"averageMovementsPerDay"`
	AverageMovementsPerAsset  float64 `json:"averageMovementsPerAsset"`
	LatestMovementDate        string  `json:"latestMovementDate"`
	EarliestMovementDate      string  `json:"earliestMovementDate"`
	UniqueAssetsWithMovements int     `json:"uniqueAssetsWithMovements"`
	UniqueLocationsInvolved   int     `json:"uniqueLocationsInvolved"`
	UniqueUsersInvolved       int     `json:"uniqueUsersInvolved"`
}
