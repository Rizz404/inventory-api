package domain

import "time"

// --- Enums ---

type ScanMethodType string

const (
	ScanMethodDataMatrix  ScanMethodType = "DATA_MATRIX"
	ScanMethodManualInput ScanMethodType = "MANUAL_INPUT"
)

type ScanResultType string

const (
	ScanResultSuccess       ScanResultType = "Success"
	ScanResultInvalidID     ScanResultType = "Invalid ID"
	ScanResultAssetNotFound ScanResultType = "Asset Not Found"
)

// --- Structs ---

type ScanLog struct {
	ID              string         `json:"id"`
	AssetID         *string        `json:"assetId"`
	ScannedValue    string         `json:"scannedValue"`
	ScanMethod      ScanMethodType `json:"scanMethod"`
	ScannedBy       string         `json:"scannedBy"`
	ScanTimestamp   time.Time      `json:"scanTimestamp"`
	ScanLocationLat *float64       `json:"scanLocationLat"`
	ScanLocationLng *float64       `json:"scanLocationLng"`
	ScanResult      ScanResultType `json:"scanResult"`
}

type ScanLogResponse struct {
	ID              string         `json:"id"`
	Asset           *AssetResponse `json:"asset,omitempty"`
	ScannedValue    string         `json:"scannedValue"`
	ScanMethod      ScanMethodType `json:"scanMethod"`
	ScannedBy       UserResponse   `json:"scannedBy"`
	ScanTimestamp   string         `json:"scanTimestamp"`
	ScanLocationLat *float64       `json:"scanLocationLat,omitempty"`
	ScanLocationLng *float64       `json:"scanLocationLng,omitempty"`
	ScanResult      ScanResultType `json:"scanResult"`
}

// --- Payloads ---

type CreateScanLogPayload struct {
	ScannedValue    string         `json:"scannedValue" validate:"required"`
	ScanMethod      ScanMethodType `json:"scanMethod" validate:"required,oneof=DATA_MATRIX MANUAL_INPUT"`
	ScanLocationLat *float64       `json:"scanLocationLat,omitempty" validate:"omitempty,latitude"`
	ScanLocationLng *float64       `json:"scanLocationLng,omitempty" validate:"omitempty,longitude"`
}
