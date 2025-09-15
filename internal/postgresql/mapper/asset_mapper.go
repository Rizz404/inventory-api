package mapper

import (
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/oklog/ulid/v2"
)

func ToModelAsset(d *domain.Asset) model.Asset {
	modelAsset := model.Asset{
		AssetTag:           d.AssetTag,
		DataMatrixValue:    d.DataMatrixValue,
		DataMatrixImageUrl: d.DataMatrixImageUrl,
		AssetName:          d.AssetName,
		Brand:              d.Brand,
		Model:              d.Model,
		SerialNumber:       d.SerialNumber,
		PurchaseDate:       d.PurchaseDate,
		PurchasePrice:      d.PurchasePrice,
		VendorName:         d.VendorName,
		WarrantyEnd:        d.WarrantyEnd,
		Status:             d.Status,
		Condition:          d.Condition,
	}

	if d.ID != "" {
		if parsedID, err := ulid.Parse(d.ID); err == nil {
			modelAsset.ID = model.SQLULID(parsedID)
		}
	}

	if d.CategoryID != "" {
		if parsedCategoryID, err := ulid.Parse(d.CategoryID); err == nil {
			modelAsset.CategoryID = model.SQLULID(parsedCategoryID)
		}
	}

	if d.LocationID != nil && *d.LocationID != "" {
		if parsedLocationID, err := ulid.Parse(*d.LocationID); err == nil {
			modelULID := model.SQLULID(parsedLocationID)
			modelAsset.LocationID = &modelULID
		}
	}

	if d.AssignedTo != nil && *d.AssignedTo != "" {
		if parsedAssignedTo, err := ulid.Parse(*d.AssignedTo); err == nil {
			modelULID := model.SQLULID(parsedAssignedTo)
			modelAsset.AssignedTo = &modelULID
		}
	}

	return modelAsset
}

func ToModelAssetForCreate(d *domain.Asset) model.Asset {
	modelAsset := model.Asset{
		AssetTag:           d.AssetTag,
		DataMatrixValue:    d.DataMatrixValue,
		DataMatrixImageUrl: d.DataMatrixImageUrl,
		AssetName:          d.AssetName,
		Brand:              d.Brand,
		Model:              d.Model,
		SerialNumber:       d.SerialNumber,
		PurchaseDate:       d.PurchaseDate,
		PurchasePrice:      d.PurchasePrice,
		VendorName:         d.VendorName,
		WarrantyEnd:        d.WarrantyEnd,
		Status:             d.Status,
		Condition:          d.Condition,
	}

	if d.CategoryID != "" {
		if parsedCategoryID, err := ulid.Parse(d.CategoryID); err == nil {
			modelAsset.CategoryID = model.SQLULID(parsedCategoryID)
		}
	}

	if d.LocationID != nil && *d.LocationID != "" {
		if parsedLocationID, err := ulid.Parse(*d.LocationID); err == nil {
			modelULID := model.SQLULID(parsedLocationID)
			modelAsset.LocationID = &modelULID
		}
	}

	if d.AssignedTo != nil && *d.AssignedTo != "" {
		if parsedAssignedTo, err := ulid.Parse(*d.AssignedTo); err == nil {
			modelULID := model.SQLULID(parsedAssignedTo)
			modelAsset.AssignedTo = &modelULID
		}
	}

	return modelAsset
}

func ToDomainAsset(m *model.Asset) domain.Asset {
	domainAsset := domain.Asset{
		ID:                 m.ID.String(),
		AssetTag:           m.AssetTag,
		DataMatrixValue:    m.DataMatrixValue,
		DataMatrixImageUrl: m.DataMatrixImageUrl,
		AssetName:          m.AssetName,
		CategoryID:         m.CategoryID.String(),
		Brand:              m.Brand,
		Model:              m.Model,
		SerialNumber:       m.SerialNumber,
		PurchaseDate:       m.PurchaseDate,
		PurchasePrice:      m.PurchasePrice,
		VendorName:         m.VendorName,
		WarrantyEnd:        m.WarrantyEnd,
		Status:             m.Status,
		Condition:          m.Condition,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}

	if m.LocationID != nil && !m.LocationID.IsZero() {
		locationIDStr := m.LocationID.String()
		domainAsset.LocationID = &locationIDStr
	}

	if m.AssignedTo != nil && !m.AssignedTo.IsZero() {
		assignedToStr := m.AssignedTo.String()
		domainAsset.AssignedTo = &assignedToStr
	}

	return domainAsset
}

func ToDomainAssetResponse(m *model.Asset, langCode string) domain.AssetResponse {
	response := domain.AssetResponse{
		ID:                 m.ID.String(),
		AssetTag:           m.AssetTag,
		DataMatrixValue:    m.DataMatrixValue,
		DataMatrixImageUrl: m.DataMatrixImageUrl,
		AssetName:          m.AssetName,
		Brand:              m.Brand,
		Model:              m.Model,
		SerialNumber:       m.SerialNumber,
		PurchasePrice:      m.PurchasePrice,
		VendorName:         m.VendorName,
		Status:             m.Status,
		Condition:          m.Condition,
		CreatedAt:          m.CreatedAt.Format(TimeFormat),
		UpdatedAt:          m.UpdatedAt.Format(TimeFormat),
	}

	// Handle date formatting
	if m.PurchaseDate != nil {
		purchaseDateStr := m.PurchaseDate.Format("2006-01-02")
		response.PurchaseDate = &purchaseDateStr
	}

	if m.WarrantyEnd != nil {
		warrantyEndStr := m.WarrantyEnd.Format("2006-01-02")
		response.WarrantyEnd = &warrantyEndStr
	}

	// Handle related entities
	if !m.Category.ID.IsZero() {
		categoryResponse := ToDomainCategoryResponse(&m.Category, langCode)
		response.Category = &categoryResponse
	}

	if m.Location != nil && !m.Location.ID.IsZero() {
		locationResponse := ToDomainLocationResponse(m.Location, langCode)
		response.Location = &locationResponse
	}

	if m.User != nil && !m.User.ID.IsZero() {
		userResponse := ToDomainUserResponse(m.User)
		response.AssignedTo = &userResponse
	}

	return response
}

func ToDomainAssetsResponse(m []model.Asset, langCode string) []domain.AssetResponse {
	responses := make([]domain.AssetResponse, len(m))
	for i, asset := range m {
		responses[i] = ToDomainAssetResponse(&asset, langCode)
	}
	return responses
}

// * Convert domain.Asset directly to domain.AssetResponse without going through model.Asset
func DomainAssetToAssetResponse(d *domain.Asset) domain.AssetResponse {
	response := domain.AssetResponse{
		ID:                 d.ID,
		AssetTag:           d.AssetTag,
		DataMatrixValue:    d.DataMatrixValue,
		DataMatrixImageUrl: d.DataMatrixImageUrl,
		AssetName:          d.AssetName,
		Brand:              d.Brand,
		Model:              d.Model,
		SerialNumber:       d.SerialNumber,
		PurchasePrice:      d.PurchasePrice,
		VendorName:         d.VendorName,
		Status:             d.Status,
		Condition:          d.Condition,
		CreatedAt:          d.CreatedAt.Format(TimeFormat),
		UpdatedAt:          d.UpdatedAt.Format(TimeFormat),
	}

	// Handle date formatting
	if d.PurchaseDate != nil {
		purchaseDateStr := d.PurchaseDate.Format("2006-01-02")
		response.PurchaseDate = &purchaseDateStr
	}

	if d.WarrantyEnd != nil {
		warrantyEndStr := d.WarrantyEnd.Format("2006-01-02")
		response.WarrantyEnd = &warrantyEndStr
	}

	return response
}

// * Convert slice of domain.Asset to slice of domain.AssetResponse
func DomainAssetsToAssetsResponse(assets []domain.Asset) []domain.AssetResponse {
	responses := make([]domain.AssetResponse, len(assets))
	for i, asset := range assets {
		responses[i] = DomainAssetToAssetResponse(&asset)
	}
	return responses
}

func ToModelAssetUpdateMap(payload *domain.UpdateAssetPayload) map[string]any {
	updates := make(map[string]any)

	if payload.AssetTag != nil {
		updates["asset_tag"] = *payload.AssetTag
	}
	if payload.DataMatrixValue != nil {
		updates["data_matrix_value"] = *payload.DataMatrixValue
	}
	if payload.DataMatrixImageUrl != nil {
		updates["data_matrix_image_url"] = *payload.DataMatrixImageUrl
	}
	if payload.AssetName != nil {
		updates["asset_name"] = *payload.AssetName
	}
	if payload.CategoryID != nil {
		updates["category_id"] = *payload.CategoryID
	}
	if payload.Brand != nil {
		updates["brand"] = payload.Brand
	}
	if payload.Model != nil {
		updates["model"] = payload.Model
	}
	if payload.SerialNumber != nil {
		updates["serial_number"] = payload.SerialNumber
	}
	if payload.PurchaseDate != nil {
		if *payload.PurchaseDate == "" {
			updates["purchase_date"] = nil
		} else {
			if parsedDate, err := time.Parse("2006-01-02", *payload.PurchaseDate); err == nil {
				updates["purchase_date"] = parsedDate
			}
		}
	}
	if payload.PurchasePrice != nil {
		updates["purchase_price"] = payload.PurchasePrice
	}
	if payload.VendorName != nil {
		updates["vendor_name"] = payload.VendorName
	}
	if payload.WarrantyEnd != nil {
		if *payload.WarrantyEnd == "" {
			updates["warranty_end"] = nil
		} else {
			if parsedDate, err := time.Parse("2006-01-02", *payload.WarrantyEnd); err == nil {
				updates["warranty_end"] = parsedDate
			}
		}
	}
	if payload.Status != nil {
		updates["status"] = *payload.Status
	}
	if payload.Condition != nil {
		updates["condition_status"] = *payload.Condition
	}
	if payload.LocationID != nil {
		if *payload.LocationID == "" {
			updates["location_id"] = nil
		} else {
			updates["location_id"] = *payload.LocationID
		}
	}
	if payload.AssignedTo != nil {
		if *payload.AssignedTo == "" {
			updates["assigned_to"] = nil
		} else {
			updates["assigned_to"] = *payload.AssignedTo
		}
	}

	return updates
}
