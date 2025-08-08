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
	ID            string         `json:"id"`
	AssetTag      string         `json:"assetTag"`
	QrCodeValue   *string        `json:"qrCodeValue"`
	NfcTagID      *string        `json:"nfcTagId"`
	AssetName     string         `json:"assetName"`
	CategoryID    string         `json:"categoryId"`
	Brand         *string        `json:"brand"`
	Model         *string        `json:"model"`
	SerialNumber  *string        `json:"serialNumber"`
	PurchaseDate  *time.Time     `json:"purchaseDate"`
	PurchasePrice *float64       `json:"purchasePrice"`
	VendorName    *string        `json:"vendorName"`
	WarrantyEnd   *time.Time     `json:"warrantyEnd"`
	Status        AssetStatus    `json:"status"`
	Condition     AssetCondition `json:"condition"`
	LocationID    *string        `json:"locationId"`
	AssignedTo    *string        `json:"assignedTo"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
}

type AssetResponse struct {
	ID            string            `json:"id"`
	AssetTag      string            `json:"assetTag"`
	QrCodeValue   *string           `json:"qrCodeValue,omitempty"`
	NfcTagID      *string           `json:"nfcTagId,omitempty"`
	AssetName     string            `json:"assetName"`
	Brand         *string           `json:"brand,omitempty"`
	Model         *string           `json:"model,omitempty"`
	SerialNumber  *string           `json:"serialNumber,omitempty"`
	PurchaseDate  *string           `json:"purchaseDate,omitempty"`
	PurchasePrice *float64          `json:"purchasePrice,omitempty"`
	VendorName    *string           `json:"vendorName,omitempty"`
	WarrantyEnd   *string           `json:"warrantyEnd,omitempty"`
	Status        AssetStatus       `json:"status"`
	Condition     AssetCondition    `json:"condition"`
	Category      *CategoryResponse `json:"category,omitempty"`
	Location      *LocationResponse `json:"location,omitempty"`
	AssignedTo    *UserResponse     `json:"assignedTo,omitempty"`
	CreatedAt     string            `json:"createdAt"`
	UpdatedAt     string            `json:"updatedAt"`
}

// --- Payloads ---

type CreateAssetPayload struct {
	AssetTag      string          `json:"assetTag" validate:"required,max=50"`
	QrCodeValue   *string         `json:"qrCodeValue,omitempty" validate:"omitempty,max=255"`
	NfcTagID      *string         `json:"nfcTagId,omitempty" validate:"omitempty,max=255"`
	AssetName     string          `json:"assetName" validate:"required,max=200"`
	CategoryID    string          `json:"categoryId" validate:"required"`
	Brand         *string         `json:"brand,omitempty" validate:"omitempty,max=100"`
	Model         *string         `json:"model,omitempty" validate:"omitempty,max=100"`
	SerialNumber  *string         `json:"serialNumber,omitempty" validate:"omitempty,max=100"`
	PurchaseDate  *string         `json:"purchaseDate,omitempty" validate:"omitempty,datetime=2006-01-02"`
	PurchasePrice *float64        `json:"purchasePrice,omitempty" validate:"omitempty,gt=0"`
	VendorName    *string         `json:"vendorName,omitempty" validate:"omitempty,max=150"`
	WarrantyEnd   *string         `json:"warrantyEnd,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Status        *AssetStatus    `json:"status,omitempty" validate:"omitempty,oneof=Active Maintenance Disposed Lost"`
	Condition     *AssetCondition `json:"condition,omitempty" validate:"omitempty,oneof=Good Fair Poor Damaged"`
	LocationID    *string         `json:"locationId,omitempty" validate:"omitempty"`
	AssignedTo    *string         `json:"assignedTo,omitempty" validate:"omitempty"`
}

type UpdateAssetPayload struct {
	AssetTag      *string         `json:"assetTag,omitempty" validate:"omitempty,max=50"`
	QrCodeValue   *string         `json:"qrCodeValue,omitempty" validate:"omitempty,max=255"`
	NfcTagID      *string         `json:"nfcTagId,omitempty" validate:"omitempty,max=255"`
	AssetName     *string         `json:"assetName,omitempty" validate:"omitempty,max=200"`
	CategoryID    *string         `json:"categoryId,omitempty" validate:"omitempty"`
	Brand         *string         `json:"brand,omitempty" validate:"omitempty,max=100"`
	Model         *string         `json:"model,omitempty" validate:"omitempty,max=100"`
	SerialNumber  *string         `json:"serialNumber,omitempty" validate:"omitempty,max=100"`
	PurchaseDate  *string         `json:"purchaseDate,omitempty" validate:"omitempty,datetime=2006-01-02"`
	PurchasePrice *float64        `json:"purchasePrice,omitempty" validate:"omitempty,gt=0"`
	VendorName    *string         `json:"vendorName,omitempty" validate:"omitempty,max=150"`
	WarrantyEnd   *string         `json:"warrantyEnd,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Status        *AssetStatus    `json:"status,omitempty" validate:"omitempty,oneof=Active Maintenance Disposed Lost"`
	Condition     *AssetCondition `json:"condition,omitempty" validate:"omitempty,oneof=Good Fair Poor Damaged"`
	LocationID    *string         `json:"locationId,omitempty" validate:"omitempty"`
	AssignedTo    *string         `json:"assignedTo,omitempty" validate:"omitempty"`
}
