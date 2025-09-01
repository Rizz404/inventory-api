package mapper

import (
	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
)

func ToDomainAssetResponse(m model.Asset, langCode string) domain.AssetResponse {
	resp := domain.AssetResponse{
		ID:            m.ID.String(),
		AssetTag:      m.AssetTag,
		QrCodeValue:   m.QrCodeValue,
		NfcTagID:      m.NfcTagID,
		AssetName:     m.AssetName,
		Brand:         m.Brand,
		Model:         m.Model,
		SerialNumber:  m.SerialNumber,
		PurchasePrice: m.PurchasePrice,
		VendorName:    m.VendorName,
		Status:        m.Status,
		Condition:     m.Condition,
		CreatedAt:     m.CreatedAt.Format(TimeFormat),
		UpdatedAt:     m.UpdatedAt.Format(TimeFormat),
	}

	if m.PurchaseDate != nil {
		resp.PurchaseDate = Ptr(m.PurchaseDate.Format(DateFormat))
	}
	if m.WarrantyEnd != nil {
		resp.WarrantyEnd = Ptr(m.WarrantyEnd.Format(DateFormat))
	}
	if m.Category.ID.String() != "" {
		resp.Category = Ptr(ToDomainCategoryResponse(&m.Category, langCode))
	}
	if m.Location != nil && m.Location.ID.String() != "" {
		resp.Location = Ptr(ToDomainLocationResponse(m.Location, langCode))
	}
	if m.User != nil && m.User.ID.String() != "" {
		resp.AssignedTo = Ptr(ToDomainUserResponse(m.User))
	}

	return resp
}

func ToDomainAssetsResponse(m []model.Asset, langCode string) []domain.AssetResponse {
	responses := make([]domain.AssetResponse, len(m))
	for i, asset := range m {
		responses[i] = ToDomainAssetResponse(asset, langCode)
	}
	return responses
}
