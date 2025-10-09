package domain

import "time"

// --- Enums ---

type AssetStatus string

const (
	StatusActive      AssetStatus = "Active"
	StatusMaintenance AssetStatus = "Maintenance"
	StatusDisposed    AssetStatus = "Disposed"
	StatusLost        AssetStatus = "Lost"
)

type AssetCondition string

const (
	ConditionGood    AssetCondition = "Good"
	ConditionFair    AssetCondition = "Fair"
	ConditionPoor    AssetCondition = "Poor"
	ConditionDamaged AssetCondition = "Damaged"
)

type AssetSortField string

const (
	AssetSortByAssetTag      AssetSortField = "assetTag"
	AssetSortByAssetName     AssetSortField = "assetName"
	AssetSortByBrand         AssetSortField = "brand"
	AssetSortByModel         AssetSortField = "model"
	AssetSortBySerialNumber  AssetSortField = "serialNumber"
	AssetSortByPurchaseDate  AssetSortField = "purchaseDate"
	AssetSortByPurchasePrice AssetSortField = "purchasePrice"
	AssetSortByVendorName    AssetSortField = "vendorName"
	AssetSortByWarrantyEnd   AssetSortField = "warrantyEnd"
	AssetSortByStatus        AssetSortField = "status"
	AssetSortByCondition     AssetSortField = "condition"
	AssetSortByCreatedAt     AssetSortField = "createdAt"
	AssetSortByUpdatedAt     AssetSortField = "updatedAt"
)

// --- Structs ---

type Asset struct {
	ID                 string         `json:"id"`
	AssetTag           string         `json:"assetTag"`
	DataMatrixImageUrl string         `json:"dataMatrixImageUrl"`
	AssetName          string         `json:"assetName"`
	CategoryID         string         `json:"categoryId"`
	Brand              *string        `json:"brand"`
	Model              *string        `json:"model"`
	SerialNumber       *string        `json:"serialNumber"`
	PurchaseDate       *time.Time     `json:"purchaseDate"`
	PurchasePrice      *float64       `json:"purchasePrice"`
	VendorName         *string        `json:"vendorName"`
	WarrantyEnd        *time.Time     `json:"warrantyEnd"`
	Status             AssetStatus    `json:"status"`
	Condition          AssetCondition `json:"condition"`
	LocationID         *string        `json:"locationId"`
	AssignedTo         *string        `json:"assignedTo"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	// * Populated
	// Todo: Masih pake translation populated, nanti benerin diakhir
	Category *Category `json:"category"`
	Location *Location `json:"location"`
	User     *User     `json:"user"`
}

type AssetResponse struct {
	ID                 string         `json:"id"`
	AssetTag           string         `json:"assetTag"`
	DataMatrixImageUrl string         `json:"dataMatrixImageUrl"`
	AssetName          string         `json:"assetName"`
	CategoryID         string         `json:"categoryId"`
	Brand              *string        `json:"brand"`
	Model              *string        `json:"model"`
	SerialNumber       *string        `json:"serialNumber"`
	PurchaseDate       *time.Time     `json:"purchaseDate"`
	PurchasePrice      *float64       `json:"purchasePrice"`
	VendorName         *string        `json:"vendorName"`
	WarrantyEnd        *time.Time     `json:"warrantyEnd"`
	Status             AssetStatus    `json:"status"`
	Condition          AssetCondition `json:"condition"`
	LocationID         *string        `json:"locationId"`
	AssignedToID       *string        `json:"assignedToId"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	// ???
	Category   *CategoryResponse `json:"category"`
	Location   *LocationResponse `json:"location"`
	AssignedTo *UserResponse     `json:"assignedTo"`
}

type AssetListResponse struct {
	ID                 string         `json:"id"`
	AssetTag           string         `json:"assetTag"`
	DataMatrixImageUrl string         `json:"dataMatrixImageUrl"`
	AssetName          string         `json:"assetName"`
	CategoryID         string         `json:"categoryId"`
	Brand              *string        `json:"brand"`
	Model              *string        `json:"model"`
	SerialNumber       *string        `json:"serialNumber"`
	PurchaseDate       *time.Time     `json:"purchaseDate"`
	PurchasePrice      *float64       `json:"purchasePrice"`
	VendorName         *string        `json:"vendorName"`
	WarrantyEnd        *time.Time     `json:"warrantyEnd"`
	Status             AssetStatus    `json:"status"`
	Condition          AssetCondition `json:"condition"`
	LocationID         *string        `json:"locationId"`
	AssignedToID       *string        `json:"assignedToId"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	// * Populated
	Category   *CategoryResponse `json:"category"`
	Location   *LocationResponse `json:"location"`
	AssignedTo *UserResponse     `json:"assignedTo"`
}

type BulkDeleteAssets struct {
	RequestedIDS []string `json:"requestedIds"`
	DeletedIDS   []string `json:"deletedIds"`
}

type BulkDeleteAssetsResponse struct {
	RequestedIDS []string `json:"requestedIds"`
	DeletedIDS   []string `json:"deletedIds"`
}

// --- Payloads ---

type CreateAssetPayload struct {
	AssetTag           string          `json:"assetTag" validate:"required,max=50"`
	DataMatrixImageUrl *string         `json:"dataMatrixImageUrl,omitempty" validate:"omitempty,url"`
	AssetName          string          `json:"assetName" validate:"required,max=200"`
	CategoryID         string          `json:"categoryId" validate:"required"`
	Brand              *string         `json:"brand,omitempty" validate:"omitempty,max=100"`
	Model              *string         `json:"model,omitempty" validate:"omitempty,max=100"`
	SerialNumber       *string         `json:"serialNumber,omitempty" validate:"omitempty,max=100"`
	PurchaseDate       *string         `json:"purchaseDate,omitempty" validate:"omitempty,datetime=2006-01-02"`
	PurchasePrice      *float64        `json:"purchasePrice,omitempty" validate:"omitempty,gt=0"`
	VendorName         *string         `json:"vendorName,omitempty" validate:"omitempty,max=150"`
	WarrantyEnd        *string         `json:"warrantyEnd,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Status             *AssetStatus    `json:"status,omitempty" validate:"omitempty,oneof=Active Maintenance Disposed Lost"`
	Condition          *AssetCondition `json:"condition,omitempty" validate:"omitempty,oneof=Good Fair Poor Damaged"`
	LocationID         *string         `json:"locationId,omitempty" validate:"omitempty"`
	AssignedTo         *string         `json:"assignedTo,omitempty" validate:"omitempty"`
}

type UpdateAssetPayload struct {
	AssetTag           *string         `json:"assetTag,omitempty" validate:"omitempty,max=50"`
	DataMatrixImageUrl *string         `json:"dataMatrixImageUrl,omitempty" validate:"omitempty,url"`
	AssetName          *string         `json:"assetName,omitempty" validate:"omitempty,max=200"`
	CategoryID         *string         `json:"categoryId,omitempty" validate:"omitempty"`
	Brand              *string         `json:"brand,omitempty" validate:"omitempty,max=100"`
	Model              *string         `json:"model,omitempty" validate:"omitempty,max=100"`
	SerialNumber       *string         `json:"serialNumber,omitempty" validate:"omitempty,max=100"`
	PurchaseDate       *string         `json:"purchaseDate,omitempty" validate:"omitempty,datetime=2006-01-02"`
	PurchasePrice      *float64        `json:"purchasePrice,omitempty" validate:"omitempty,gt=0"`
	VendorName         *string         `json:"vendorName,omitempty" validate:"omitempty,max=150"`
	WarrantyEnd        *string         `json:"warrantyEnd,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Status             *AssetStatus    `json:"status,omitempty" validate:"omitempty,oneof=Active Maintenance Disposed Lost"`
	Condition          *AssetCondition `json:"condition,omitempty" validate:"omitempty,oneof=Good Fair Poor Damaged"`
	LocationID         *string         `json:"locationId,omitempty" validate:"omitempty"`
	AssignedTo         *string         `json:"assignedTo,omitempty" validate:"omitempty"`
}

type BulkDeleteAssetsPayload struct {
	IDS []string `json:"ids" validate:"required,min=1,max=100,dive,required"`
}

// --- Query Parameters ---

type AssetFilterOptions struct {
	Status     *AssetStatus    `json:"status,omitempty"`
	Condition  *AssetCondition `json:"condition,omitempty"`
	CategoryID *string         `json:"categoryId,omitempty"`
	LocationID *string         `json:"locationId,omitempty"`
	AssignedTo *string         `json:"assignedTo,omitempty"`
	Brand      *string         `json:"brand,omitempty"`
	Model      *string         `json:"model,omitempty"`
}

type AssetSortOptions struct {
	Field AssetSortField `json:"field" example:"createdAt"`
	Order SortOrder      `json:"order" example:"desc"`
}

type AssetParams struct {
	SearchQuery *string             `json:"searchQuery,omitempty"`
	Filters     *AssetFilterOptions `json:"filters,omitempty"`
	Sort        *AssetSortOptions   `json:"sort,omitempty"`
	Pagination  *PaginationOptions  `json:"pagination,omitempty"`
}

// --- Statistics ---

// Internal statistics structs (used in repository layer)
type AssetStatistics struct {
	Total              AssetCountStatistics      `json:"total"`
	ByStatus           AssetStatusStatistics     `json:"byStatus"`
	ByCondition        AssetConditionStatistics  `json:"byCondition"`
	ByCategory         []CategoryStatistics      `json:"byCategory"`
	ByLocation         []LocationStatistics      `json:"byLocation"`
	ByAssignment       AssetAssignmentStatistics `json:"byAssignment"`
	ValueStatistics    AssetValueStatistics      `json:"valueStatistics"`
	WarrantyStatistics AssetWarrantyStatistics   `json:"warrantyStatistics"`
	CreationTrends     []AssetCreationTrend      `json:"creationTrends"`
	Summary            AssetSummaryStatistics    `json:"summary"`
}

type AssetCountStatistics struct {
	Count int `json:"count"`
}

type AssetStatusStatistics struct {
	Active      int `json:"active"`
	Maintenance int `json:"maintenance"`
	Disposed    int `json:"disposed"`
	Lost        int `json:"lost"`
}

type AssetConditionStatistics struct {
	Good    int `json:"good"`
	Fair    int `json:"fair"`
	Poor    int `json:"poor"`
	Damaged int `json:"damaged"`
}

type AssetAssignmentStatistics struct {
	Assigned   int `json:"assigned"`
	Unassigned int `json:"unassigned"`
}

type AssetValueStatistics struct {
	TotalValue         *float64 `json:"totalValue"`
	AverageValue       *float64 `json:"averageValue"`
	MinValue           *float64 `json:"minValue"`
	MaxValue           *float64 `json:"maxValue"`
	AssetsWithValue    int      `json:"assetsWithValue"`
	AssetsWithoutValue int      `json:"assetsWithoutValue"`
}

type AssetWarrantyStatistics struct {
	ActiveWarranties  int `json:"activeWarranties"`
	ExpiredWarranties int `json:"expiredWarranties"`
	NoWarrantyInfo    int `json:"noWarrantyInfo"`
}

type AssetCreationTrend struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

type AssetSummaryStatistics struct {
	TotalAssets                 int       `json:"totalAssets"`
	ActiveAssetsPercentage      float64   `json:"activeAssetsPercentage"`
	MaintenanceAssetsPercentage float64   `json:"maintenanceAssetsPercentage"`
	DisposedAssetsPercentage    float64   `json:"disposedAssetsPercentage"`
	LostAssetsPercentage        float64   `json:"lostAssetsPercentage"`
	GoodConditionPercentage     float64   `json:"goodConditionPercentage"`
	FairConditionPercentage     float64   `json:"fairConditionPercentage"`
	PoorConditionPercentage     float64   `json:"poorConditionPercentage"`
	DamagedConditionPercentage  float64   `json:"damagedConditionPercentage"`
	AssignedAssetsPercentage    float64   `json:"assignedAssetsPercentage"`
	UnassignedAssetsPercentage  float64   `json:"unassignedAssetsPercentage"`
	AssetsWithPurchasePrice     int       `json:"assetsWithPurchasePrice"`
	PurchasePricePercentage     float64   `json:"purchasePricePercentage"`
	AssetsWithDataMatrix        int       `json:"assetsWithDataMatrix"`
	DataMatrixPercentage        float64   `json:"dataMatrixPercentage"`
	AssetsWithWarranty          int       `json:"assetsWithWarranty"`
	WarrantyPercentage          float64   `json:"warrantyPercentage"`
	TotalCategories             int       `json:"totalCategories"`
	TotalLocations              int       `json:"totalLocations"`
	AverageAssetsPerDay         float64   `json:"averageAssetsPerDay"`
	LatestCreationDate          time.Time `json:"latestCreationDate"`
	EarliestCreationDate        time.Time `json:"earliestCreationDate"`
	MostExpensiveAssetValue     *float64  `json:"mostExpensiveAssetValue"`
	LeastExpensiveAssetValue    *float64  `json:"leastExpensiveAssetValue"`
}

// Response statistics structs (used in service/handler layer)
type AssetStatisticsResponse struct {
	Total              AssetCountStatisticsResponse      `json:"total"`
	ByStatus           AssetStatusStatisticsResponse     `json:"byStatus"`
	ByCondition        AssetConditionStatisticsResponse  `json:"byCondition"`
	ByCategory         []CategoryStatisticsResponse      `json:"byCategory"`
	ByLocation         []LocationStatisticsResponse      `json:"byLocation"`
	ByAssignment       AssetAssignmentStatisticsResponse `json:"byAssignment"`
	ValueStatistics    AssetValueStatisticsResponse      `json:"valueStatistics"`
	WarrantyStatistics AssetWarrantyStatisticsResponse   `json:"warrantyStatistics"`
	CreationTrends     []AssetCreationTrendResponse      `json:"creationTrends"`
	Summary            AssetSummaryStatisticsResponse    `json:"summary"`
}

type AssetCountStatisticsResponse struct {
	Count int `json:"count"`
}

type AssetStatusStatisticsResponse struct {
	Active      int `json:"active"`
	Maintenance int `json:"maintenance"`
	Disposed    int `json:"disposed"`
	Lost        int `json:"lost"`
}

type AssetConditionStatisticsResponse struct {
	Good    int `json:"good"`
	Fair    int `json:"fair"`
	Poor    int `json:"poor"`
	Damaged int `json:"damaged"`
}

type AssetAssignmentStatisticsResponse struct {
	Assigned   int `json:"assigned"`
	Unassigned int `json:"unassigned"`
}

type AssetValueStatisticsResponse struct {
	TotalValue         *float64 `json:"totalValue"`
	AverageValue       *float64 `json:"averageValue"`
	MinValue           *float64 `json:"minValue"`
	MaxValue           *float64 `json:"maxValue"`
	AssetsWithValue    int      `json:"assetsWithValue"`
	AssetsWithoutValue int      `json:"assetsWithoutValue"`
}

type AssetWarrantyStatisticsResponse struct {
	ActiveWarranties  int `json:"activeWarranties"`
	ExpiredWarranties int `json:"expiredWarranties"`
	NoWarrantyInfo    int `json:"noWarrantyInfo"`
}

type AssetCreationTrendResponse struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

type AssetSummaryStatisticsResponse struct {
	TotalAssets                 int       `json:"totalAssets"`
	ActiveAssetsPercentage      float64   `json:"activeAssetsPercentage"`
	MaintenanceAssetsPercentage float64   `json:"maintenanceAssetsPercentage"`
	DisposedAssetsPercentage    float64   `json:"disposedAssetsPercentage"`
	LostAssetsPercentage        float64   `json:"lostAssetsPercentage"`
	GoodConditionPercentage     float64   `json:"goodConditionPercentage"`
	FairConditionPercentage     float64   `json:"fairConditionPercentage"`
	PoorConditionPercentage     float64   `json:"poorConditionPercentage"`
	DamagedConditionPercentage  float64   `json:"damagedConditionPercentage"`
	AssignedAssetsPercentage    float64   `json:"assignedAssetsPercentage"`
	UnassignedAssetsPercentage  float64   `json:"unassignedAssetsPercentage"`
	AssetsWithPurchasePrice     int       `json:"assetsWithPurchasePrice"`
	PurchasePricePercentage     float64   `json:"purchasePricePercentage"`
	AssetsWithDataMatrix        int       `json:"assetsWithDataMatrix"`
	DataMatrixPercentage        float64   `json:"dataMatrixPercentage"`
	AssetsWithWarranty          int       `json:"assetsWithWarranty"`
	WarrantyPercentage          float64   `json:"warrantyPercentage"`
	TotalCategories             int       `json:"totalCategories"`
	TotalLocations              int       `json:"totalLocations"`
	AverageAssetsPerDay         float64   `json:"averageAssetsPerDay"`
	LatestCreationDate          time.Time `json:"latestCreationDate"`
	EarliestCreationDate        time.Time `json:"earliestCreationDate"`
	MostExpensiveAssetValue     *float64  `json:"mostExpensiveAssetValue"`
	LeastExpensiveAssetValue    *float64  `json:"leastExpensiveAssetValue"`
}
