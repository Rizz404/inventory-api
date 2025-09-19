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
}

type AssetResponse struct {
	ID                 string            `json:"id"`
	AssetTag           string            `json:"assetTag"`
	DataMatrixImageUrl string            `json:"dataMatrixImageUrl"`
	AssetName          string            `json:"assetName"`
	Brand              *string           `json:"brand,omitempty"`
	Model              *string           `json:"model,omitempty"`
	SerialNumber       *string           `json:"serialNumber,omitempty"`
	PurchaseDate       *string           `json:"purchaseDate,omitempty"`
	PurchasePrice      *float64          `json:"purchasePrice,omitempty"`
	VendorName         *string           `json:"vendorName,omitempty"`
	WarrantyEnd        *string           `json:"warrantyEnd,omitempty"`
	Status             AssetStatus       `json:"status"`
	Condition          AssetCondition    `json:"condition"`
	Category           *CategoryResponse `json:"category,omitempty"`
	Location           *LocationResponse `json:"location,omitempty"`
	AssignedTo         *UserResponse     `json:"assignedTo,omitempty"`
	CreatedAt          string            `json:"createdAt"`
	UpdatedAt          string            `json:"updatedAt"`
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

// --- Statistics ---

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
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type AssetSummaryStatistics struct {
	TotalAssets                 int      `json:"totalAssets"`
	ActiveAssetsPercentage      float64  `json:"activeAssetsPercentage"`
	MaintenanceAssetsPercentage float64  `json:"maintenanceAssetsPercentage"`
	DisposedAssetsPercentage    float64  `json:"disposedAssetsPercentage"`
	LostAssetsPercentage        float64  `json:"lostAssetsPercentage"`
	GoodConditionPercentage     float64  `json:"goodConditionPercentage"`
	FairConditionPercentage     float64  `json:"fairConditionPercentage"`
	PoorConditionPercentage     float64  `json:"poorConditionPercentage"`
	DamagedConditionPercentage  float64  `json:"damagedConditionPercentage"`
	AssignedAssetsPercentage    float64  `json:"assignedAssetsPercentage"`
	UnassignedAssetsPercentage  float64  `json:"unassignedAssetsPercentage"`
	AssetsWithPurchasePrice     int      `json:"assetsWithPurchasePrice"`
	PurchasePricePercentage     float64  `json:"purchasePricePercentage"`
	AssetsWithDataMatrix        int      `json:"assetsWithDataMatrix"`
	DataMatrixPercentage        float64  `json:"dataMatrixPercentage"`
	AssetsWithWarranty          int      `json:"assetsWithWarranty"`
	WarrantyPercentage          float64  `json:"warrantyPercentage"`
	TotalCategories             int      `json:"totalCategories"`
	TotalLocations              int      `json:"totalLocations"`
	AverageAssetsPerDay         float64  `json:"averageAssetsPerDay"`
	LatestCreationDate          string   `json:"latestCreationDate"`
	EarliestCreationDate        string   `json:"earliestCreationDate"`
	MostExpensiveAssetValue     *float64 `json:"mostExpensiveAssetValue"`
	LeastExpensiveAssetValue    *float64 `json:"leastExpensiveAssetValue"`
}
