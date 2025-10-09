package mapper

import (
	"fmt"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/oklog/ulid/v2"
)

// Helper function to format float64 pointer to string pointer with 2 decimal places
func formatPriceToString(price *float64) *string {
	if price == nil {
		return nil
	}
	formatted := fmt.Sprintf("%.2f", *price)
	return &formatted
}

// Helper function to format float64 to string with 2 decimal places
func formatFloat64ToString(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

// *==================== Model conversions ====================
func ToModelAsset(d *domain.Asset) model.Asset {
	modelAsset := model.Asset{
		AssetTag:           d.AssetTag,
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

// *==================== Entity conversions ====================
func ToDomainAsset(m *model.Asset) domain.Asset {
	domainAsset := domain.Asset{
		ID:                 m.ID.String(),
		AssetTag:           m.AssetTag,
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

	// Populate related entities if preloaded
	if !m.Category.ID.IsZero() {
		category := ToDomainCategory(&m.Category)
		domainAsset.Category = &category
	}

	if m.Location != nil && !m.Location.ID.IsZero() {
		location := ToDomainLocation(m.Location)
		domainAsset.Location = &location
	}

	if m.User != nil && !m.User.ID.IsZero() {
		user := ToDomainUser(m.User)
		domainAsset.User = &user
	}

	return domainAsset
}

func ToDomainAssets(models []model.Asset) []domain.Asset {
	assets := make([]domain.Asset, len(models))
	for i, model := range models {
		assets[i] = ToDomainAsset(&model)
	}
	return assets
}

// *==================== Entity Response conversions ====================
func AssetToResponse(d *domain.Asset, langCode string) domain.AssetResponse {
	response := domain.AssetResponse{
		ID:                 d.ID,
		AssetTag:           d.AssetTag,
		DataMatrixImageUrl: d.DataMatrixImageUrl,
		AssetName:          d.AssetName,
		CategoryID:         d.CategoryID,
		Brand:              d.Brand,
		Model:              d.Model,
		SerialNumber:       d.SerialNumber,
		PurchaseDate:       d.PurchaseDate,
		PurchasePrice:      formatPriceToString(d.PurchasePrice), // Format price as string with 2 decimals
		VendorName:         d.VendorName,
		WarrantyEnd:        d.WarrantyEnd,
		Status:             d.Status,
		Condition:          d.Condition,
		LocationID:         d.LocationID,
		AssignedToID:       d.AssignedTo,
		CreatedAt:          d.CreatedAt,
		UpdatedAt:          d.UpdatedAt,
	}

	// Populate related entities
	if d.Category != nil {
		categoryResponse := CategoryToResponse(d.Category, langCode)
		response.Category = &categoryResponse
	}

	if d.Location != nil {
		locationResponse := LocationToResponse(d.Location, langCode)
		response.Location = &locationResponse
	}

	if d.User != nil {
		userResponse := UserToResponse(d.User)
		response.AssignedTo = &userResponse
	}

	return response
}

func AssetsToResponses(assets []domain.Asset, langCode string) []domain.AssetResponse {
	responses := make([]domain.AssetResponse, len(assets))
	for i, asset := range assets {
		responses[i] = AssetToResponse(&asset, langCode)
	}
	return responses
}

func AssetToListResponse(d *domain.Asset, langCode string) domain.AssetListResponse {
	response := domain.AssetListResponse{
		ID:                 d.ID,
		AssetTag:           d.AssetTag,
		DataMatrixImageUrl: d.DataMatrixImageUrl,
		AssetName:          d.AssetName,
		CategoryID:         d.CategoryID,
		Brand:              d.Brand,
		Model:              d.Model,
		SerialNumber:       d.SerialNumber,
		PurchaseDate:       d.PurchaseDate,
		PurchasePrice:      formatPriceToString(d.PurchasePrice), // Format price as string with 2 decimals
		VendorName:         d.VendorName,
		WarrantyEnd:        d.WarrantyEnd,
		Status:             d.Status,
		Condition:          d.Condition,
		LocationID:         d.LocationID,
		AssignedToID:       d.AssignedTo,
		CreatedAt:          d.CreatedAt,
		UpdatedAt:          d.UpdatedAt,
	}

	// Populate related entities
	if d.Category != nil {
		categoryResponse := CategoryToResponse(d.Category, langCode)
		response.Category = &categoryResponse
	}

	if d.Location != nil {
		locationResponse := LocationToResponse(d.Location, langCode)
		response.Location = &locationResponse
	}

	if d.User != nil {
		userResponse := UserToResponse(d.User)
		response.AssignedTo = &userResponse
	}

	return response
}

func AssetsToListResponses(assets []domain.Asset, langCode string) []domain.AssetListResponse {
	responses := make([]domain.AssetListResponse, len(assets))
	for i, asset := range assets {
		responses[i] = AssetToListResponse(&asset, langCode)
	}
	return responses
}

func AssetStatisticsToResponse(stats *domain.AssetStatistics) domain.AssetStatisticsResponse {
	response := domain.AssetStatisticsResponse{
		Total: domain.AssetCountStatisticsResponse{
			Count: stats.Total.Count,
		},
		ByStatus: domain.AssetStatusStatisticsResponse{
			Active:      stats.ByStatus.Active,
			Maintenance: stats.ByStatus.Maintenance,
			Disposed:    stats.ByStatus.Disposed,
			Lost:        stats.ByStatus.Lost,
		},
		ByCondition: domain.AssetConditionStatisticsResponse{
			Good:    stats.ByCondition.Good,
			Fair:    stats.ByCondition.Fair,
			Poor:    stats.ByCondition.Poor,
			Damaged: stats.ByCondition.Damaged,
		},
		ByAssignment: domain.AssetAssignmentStatisticsResponse{
			Assigned:   stats.ByAssignment.Assigned,
			Unassigned: stats.ByAssignment.Unassigned,
		},
		ValueStatistics: domain.AssetValueStatisticsResponse{
			TotalValue:         formatPriceToString(stats.ValueStatistics.TotalValue),
			AverageValue:       formatPriceToString(stats.ValueStatistics.AverageValue),
			MinValue:           formatPriceToString(stats.ValueStatistics.MinValue),
			MaxValue:           formatPriceToString(stats.ValueStatistics.MaxValue),
			AssetsWithValue:    stats.ValueStatistics.AssetsWithValue,
			AssetsWithoutValue: stats.ValueStatistics.AssetsWithoutValue,
		},
		WarrantyStatistics: domain.AssetWarrantyStatisticsResponse{
			ActiveWarranties:  stats.WarrantyStatistics.ActiveWarranties,
			ExpiredWarranties: stats.WarrantyStatistics.ExpiredWarranties,
			NoWarrantyInfo:    stats.WarrantyStatistics.NoWarrantyInfo,
		},
		Summary: domain.AssetSummaryStatisticsResponse{
			TotalAssets:                 stats.Summary.TotalAssets,
			ActiveAssetsPercentage:      stats.Summary.ActiveAssetsPercentage,
			MaintenanceAssetsPercentage: stats.Summary.MaintenanceAssetsPercentage,
			DisposedAssetsPercentage:    stats.Summary.DisposedAssetsPercentage,
			LostAssetsPercentage:        stats.Summary.LostAssetsPercentage,
			GoodConditionPercentage:     stats.Summary.GoodConditionPercentage,
			FairConditionPercentage:     stats.Summary.FairConditionPercentage,
			PoorConditionPercentage:     stats.Summary.PoorConditionPercentage,
			DamagedConditionPercentage:  stats.Summary.DamagedConditionPercentage,
			AssignedAssetsPercentage:    stats.Summary.AssignedAssetsPercentage,
			UnassignedAssetsPercentage:  stats.Summary.UnassignedAssetsPercentage,
			AssetsWithPurchasePrice:     stats.Summary.AssetsWithPurchasePrice,
			PurchasePricePercentage:     stats.Summary.PurchasePricePercentage,
			AssetsWithDataMatrix:        stats.Summary.AssetsWithDataMatrix,
			DataMatrixPercentage:        stats.Summary.DataMatrixPercentage,
			AssetsWithWarranty:          stats.Summary.AssetsWithWarranty,
			WarrantyPercentage:          stats.Summary.WarrantyPercentage,
			TotalCategories:             stats.Summary.TotalCategories,
			TotalLocations:              stats.Summary.TotalLocations,
			AverageAssetsPerDay:         stats.Summary.AverageAssetsPerDay,
			LatestCreationDate:          stats.Summary.LatestCreationDate,
			EarliestCreationDate:        stats.Summary.EarliestCreationDate,
			MostExpensiveAssetValue:     formatPriceToString(stats.Summary.MostExpensiveAssetValue),
			LeastExpensiveAssetValue:    formatPriceToString(stats.Summary.LeastExpensiveAssetValue),
		},
	}

	// Convert ByCategory slice
	response.ByCategory = make([]domain.CategoryStatisticsResponse, len(stats.ByCategory))
	for i, category := range stats.ByCategory {
		response.ByCategory[i] = CategoryStatisticsToResponse(&category)
	}

	// Convert ByLocation slice
	response.ByLocation = make([]domain.LocationStatisticsResponse, len(stats.ByLocation))
	for i, location := range stats.ByLocation {
		response.ByLocation[i] = LocationStatisticsToResponse(&location)
	}

	// Convert CreationTrends slice
	response.CreationTrends = make([]domain.AssetCreationTrendResponse, len(stats.CreationTrends))
	for i, trend := range stats.CreationTrends {
		response.CreationTrends[i] = domain.AssetCreationTrendResponse{
			Date:  trend.Date,
			Count: trend.Count,
		}
	}

	return response
}

// *==================== Update Map conversions ====================
func ToModelAssetUpdateMap(payload *domain.UpdateAssetPayload) map[string]any {
	updates := make(map[string]any)

	if payload.AssetTag != nil {
		updates["asset_tag"] = *payload.AssetTag
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

func MapAssetSortFieldToColumn(field domain.AssetSortField) string {
	columnMap := map[domain.AssetSortField]string{
		domain.AssetSortByAssetTag:      "asset_tag",
		domain.AssetSortByAssetName:     "asset_name",
		domain.AssetSortByBrand:         "brand",
		domain.AssetSortByModel:         "model",
		domain.AssetSortBySerialNumber:  "serial_number",
		domain.AssetSortByPurchaseDate:  "purchase_date",
		domain.AssetSortByPurchasePrice: "purchase_price",
		domain.AssetSortByVendorName:    "vendor_name",
		domain.AssetSortByWarrantyEnd:   "warranty_end",
		domain.AssetSortByStatus:        "status",
		domain.AssetSortByCondition:     "condition_status",
		domain.AssetSortByCreatedAt:     "created_at",
		domain.AssetSortByUpdatedAt:     "updated_at",
	}

	if column, exists := columnMap[field]; exists {
		return column
	}
	return "created_at"
}
