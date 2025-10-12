package asset

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/xuri/excelize/v2"
)

// ExportAssetList exports asset list to PDF or Excel format
func (s *Service) ExportAssetList(ctx context.Context, payload *domain.ExportAssetListPayload, langCode string) ([]byte, string, error) {
	// Build params from payload
	params := domain.AssetParams{
		SearchQuery: payload.SearchQuery,
		Filters:     payload.Filters,
		Sort:        payload.Sort,
	}

	// Get assets without pagination
	assets, err := s.Repo.GetAssetsForExport(ctx, params, langCode)
	if err != nil {
		return nil, "", err
	}

	// Convert to responses
	assetResponses := mapper.AssetsToResponses(assets, langCode)

	switch payload.Format {
	case domain.ExportFormatPDF:
		data, err := s.exportAssetListToPDF(assetResponses, payload.IncludeDataMatrixImage)
		if err != nil {
			return nil, "", domain.ErrInternal(err)
		}
		return data, "asset_list.pdf", nil

	case domain.ExportFormatExcel:
		data, err := s.exportAssetListToExcel(assetResponses)
		if err != nil {
			return nil, "", domain.ErrInternal(err)
		}
		return data, "asset_list.xlsx", nil

	default:
		return nil, "", domain.ErrBadRequest("Invalid export format")
	}
}

// exportAssetListToPDF generates PDF file for asset list
func (s *Service) exportAssetListToPDF(assets []domain.AssetResponse, includeDataMatrix bool) ([]byte, error) {
	cfg := config.NewBuilder().
		WithPageSize(pagesize.A4).
		WithOrientation(orientation.Horizontal).
		WithLeftMargin(10).
		WithTopMargin(10).
		WithRightMargin(10).
		Build()

	mrt := maroto.New(cfg)

	// Add title
	mrt.AddRows(
		row.New(15).Add(
			col.New(12).Add(
				text.New("Asset List Report", props.Text{
					Size:  16,
					Style: fontstyle.Bold,
					Align: align.Center,
				}),
			),
		),
		row.New(8).Add(
			col.New(12).Add(
				text.New(fmt.Sprintf("Generated on: %s", time.Now().Format("2006-01-02 15:04:05")), props.Text{
					Size:  10,
					Align: align.Center,
				}),
			),
		),
	)

	// Add table header
	headerRow := row.New(10)
	headerRow.Add(
		col.New(2).Add(text.New("Asset Tag", props.Text{Style: fontstyle.Bold, Size: 9})),
		col.New(3).Add(text.New("Asset Name", props.Text{Style: fontstyle.Bold, Size: 9})),
		col.New(2).Add(text.New("Category", props.Text{Style: fontstyle.Bold, Size: 9})),
		col.New(1).Add(text.New("Brand", props.Text{Style: fontstyle.Bold, Size: 9})),
		col.New(1).Add(text.New("Model", props.Text{Style: fontstyle.Bold, Size: 9})),
		col.New(1).Add(text.New("Status", props.Text{Style: fontstyle.Bold, Size: 9})),
		col.New(1).Add(text.New("Condition", props.Text{Style: fontstyle.Bold, Size: 9})),
		col.New(1).Add(text.New("Location", props.Text{Style: fontstyle.Bold, Size: 9})),
	)
	mrt.AddRows(headerRow)

	// Add asset data rows
	for _, asset := range assets {
		dataRow := row.New(8)

		categoryName := ""
		if asset.Category != nil {
			categoryName = asset.Category.CategoryName
		}

		locationName := ""
		if asset.Location != nil {
			locationName = asset.Location.LocationName
		}

		brand := ""
		if asset.Brand != nil {
			brand = *asset.Brand
		}

		model := ""
		if asset.Model != nil {
			model = *asset.Model
		}

		dataRow.Add(
			col.New(2).Add(text.New(asset.AssetTag, props.Text{Size: 8})),
			col.New(3).Add(text.New(asset.AssetName, props.Text{Size: 8})),
			col.New(2).Add(text.New(categoryName, props.Text{Size: 8})),
			col.New(1).Add(text.New(brand, props.Text{Size: 8})),
			col.New(1).Add(text.New(model, props.Text{Size: 8})),
			col.New(1).Add(text.New(string(asset.Status), props.Text{Size: 8})),
			col.New(1).Add(text.New(string(asset.Condition), props.Text{Size: 8})),
			col.New(1).Add(text.New(locationName, props.Text{Size: 8})),
		)
		mrt.AddRows(dataRow)
	}

	// Add footer with total count
	mrt.AddRows(
		row.New(10).Add(
			col.New(12).Add(
				text.New(fmt.Sprintf("Total Assets: %d", len(assets)), props.Text{
					Size:  10,
					Style: fontstyle.Bold,
					Align: align.Right,
				}),
			),
		),
	)

	// Generate PDF
	document, err := mrt.Generate()
	if err != nil {
		return nil, err
	}

	return document.GetBytes(), nil
}

// exportAssetListToExcel generates Excel file for asset list
func (s *Service) exportAssetListToExcel(assets []domain.AssetResponse) ([]byte, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error closing Excel file:", err)
		}
	}()

	sheetName := "Assets"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	// Set active sheet
	f.SetActiveSheet(index)

	// Create header style
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}

	// Set headers
	headers := []string{
		"Asset Tag", "Asset Name", "Category", "Brand", "Model",
		"Serial Number", "Purchase Date", "Purchase Price", "Vendor",
		"Warranty End", "Status", "Condition", "Location", "Assigned To",
	}

	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Add data
	for row, asset := range assets {
		rowNum := row + 2 // Start from row 2 (after header)

		categoryName := ""
		if asset.Category != nil {
			categoryName = asset.Category.CategoryName
		}

		locationName := ""
		if asset.Location != nil {
			locationName = asset.Location.LocationName
		}

		assignedToName := ""
		if asset.AssignedTo != nil {
			assignedToName = asset.AssignedTo.FullName
		}

		brand := ""
		if asset.Brand != nil {
			brand = *asset.Brand
		}

		model := ""
		if asset.Model != nil {
			model = *asset.Model
		}

		serialNumber := ""
		if asset.SerialNumber != nil {
			serialNumber = *asset.SerialNumber
		}

		purchaseDate := ""
		if asset.PurchaseDate != nil {
			purchaseDate = asset.PurchaseDate.Format("2006-01-02")
		}

		purchasePrice := ""
		if asset.PurchasePrice != nil && asset.PurchasePrice.Valid {
			value, _ := asset.PurchasePrice.Float64()
			purchasePrice = fmt.Sprintf("%.2f", value)
		}

		vendor := ""
		if asset.VendorName != nil {
			vendor = *asset.VendorName
		}

		warrantyEnd := ""
		if asset.WarrantyEnd != nil {
			warrantyEnd = asset.WarrantyEnd.Format("2006-01-02")
		}

		// Set cell values
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), asset.AssetTag)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), asset.AssetName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), categoryName)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), brand)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), model)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), serialNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), purchaseDate)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), purchasePrice)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", rowNum), vendor)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", rowNum), warrantyEnd)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", rowNum), string(asset.Status))
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", rowNum), string(asset.Condition))
		f.SetCellValue(sheetName, fmt.Sprintf("M%d", rowNum), locationName)
		f.SetCellValue(sheetName, fmt.Sprintf("N%d", rowNum), assignedToName)
	}

	// Auto-fit columns
	for col := 1; col <= len(headers); col++ {
		colName, _ := excelize.ColumnNumberToName(col)
		f.SetColWidth(sheetName, colName, colName, 15)
	}

	// Save to buffer
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// ExportAssetStatistics exports asset statistics to PDF with charts
func (s *Service) ExportAssetStatistics(ctx context.Context, langCode string) ([]byte, string, error) {
	// Get statistics
	stats, err := s.Repo.GetAssetStatistics(ctx)
	if err != nil {
		return nil, "", err
	}

	// Generate PDF with charts
	data, err := s.exportAssetStatisticsToPDF(stats)
	if err != nil {
		return nil, "", domain.ErrInternal(err)
	}

	return data, "asset_statistics.pdf", nil
}

// exportAssetStatisticsToPDF generates PDF file with charts for statistics
func (s *Service) exportAssetStatisticsToPDF(stats domain.AssetStatistics) ([]byte, error) {
	cfg := config.NewBuilder().
		WithPageSize(pagesize.A4).
		WithOrientation(orientation.Vertical).
		WithLeftMargin(15).
		WithTopMargin(15).
		WithRightMargin(15).
		Build()

	mrt := maroto.New(cfg)

	// Add title
	mrt.AddRows(
		row.New(20).Add(
			col.New(12).Add(
				text.New("Asset Statistics Report", props.Text{
					Size:  18,
					Style: fontstyle.Bold,
					Align: align.Center,
				}),
			),
		),
		row.New(8).Add(
			col.New(12).Add(
				text.New(fmt.Sprintf("Generated on: %s", time.Now().Format("2006-01-02 15:04:05")), props.Text{
					Size:  10,
					Align: align.Center,
				}),
			),
		),
	)

	// Generate charts and add them to PDF
	// 1. Status Distribution Chart
	statusChartPath, err := s.generateStatusChart(stats.ByStatus)
	if err == nil && statusChartPath != "" {
		defer os.Remove(statusChartPath) // Clean up temp file

		mrt.AddRows(
			row.New(10).Add(
				col.New(12).Add(
					text.New("Status Distribution", props.Text{
						Size:  14,
						Style: fontstyle.Bold,
					}),
				),
			),
			row.New(70).Add(
				col.New(12).Add(
					image.NewFromFile(statusChartPath, props.Rect{
						Center:  true,
						Percent: 80,
					}),
				),
			),
		)
	}

	// 2. Condition Distribution Chart
	conditionChartPath, err := s.generateConditionChart(stats.ByCondition)
	if err == nil && conditionChartPath != "" {
		defer os.Remove(conditionChartPath)

		mrt.AddRows(
			row.New(10).Add(
				col.New(12).Add(
					text.New("Condition Distribution", props.Text{
						Size:  14,
						Style: fontstyle.Bold,
					}),
				),
			),
			row.New(70).Add(
				col.New(12).Add(
					image.NewFromFile(conditionChartPath, props.Rect{
						Center:  true,
						Percent: 80,
					}),
				),
			),
		)
	}

	// 3. Summary Statistics Table
	mrt.AddRows(
		row.New(10).Add(
			col.New(12).Add(
				text.New("Summary Statistics", props.Text{
					Size:  14,
					Style: fontstyle.Bold,
				}),
			),
		),
		row.New(8).Add(
			col.New(6).Add(text.New("Total Assets:", props.Text{Style: fontstyle.Bold, Size: 10})),
			col.New(6).Add(text.New(fmt.Sprintf("%d", stats.Summary.TotalAssets), props.Text{Size: 10})),
		),
		row.New(8).Add(
			col.New(6).Add(text.New("Total Categories:", props.Text{Style: fontstyle.Bold, Size: 10})),
			col.New(6).Add(text.New(fmt.Sprintf("%d", stats.Summary.TotalCategories), props.Text{Size: 10})),
		),
		row.New(8).Add(
			col.New(6).Add(text.New("Total Locations:", props.Text{Style: fontstyle.Bold, Size: 10})),
			col.New(6).Add(text.New(fmt.Sprintf("%d", stats.Summary.TotalLocations), props.Text{Size: 10})),
		),
		row.New(8).Add(
			col.New(6).Add(text.New("Active Assets:", props.Text{Style: fontstyle.Bold, Size: 10})),
			col.New(6).Add(text.New(fmt.Sprintf("%d (%.2f%%)", stats.ByStatus.Active, stats.Summary.ActiveAssetsPercentage), props.Text{Size: 10})),
		),
		row.New(8).Add(
			col.New(6).Add(text.New("Assigned Assets:", props.Text{Style: fontstyle.Bold, Size: 10})),
			col.New(6).Add(text.New(fmt.Sprintf("%d (%.2f%%)", stats.ByAssignment.Assigned, stats.Summary.AssignedAssetsPercentage), props.Text{Size: 10})),
		),
	)

	// Add value statistics if available
	if stats.ValueStatistics.TotalValue != nil {
		mrt.AddRows(
			row.New(8).Add(
				col.New(6).Add(text.New("Total Value:", props.Text{Style: fontstyle.Bold, Size: 10})),
				col.New(6).Add(text.New(fmt.Sprintf("$%.2f", *stats.ValueStatistics.TotalValue), props.Text{Size: 10})),
			),
		)
	}

	if stats.ValueStatistics.AverageValue != nil {
		mrt.AddRows(
			row.New(8).Add(
				col.New(6).Add(text.New("Average Value:", props.Text{Style: fontstyle.Bold, Size: 10})),
				col.New(6).Add(text.New(fmt.Sprintf("$%.2f", *stats.ValueStatistics.AverageValue), props.Text{Size: 10})),
			),
		)
	}

	// Generate PDF
	document, err := mrt.Generate()
	if err != nil {
		return nil, err
	}

	return document.GetBytes(), nil
}

// generateStatusChart creates a pie chart for status distribution
func (s *Service) generateStatusChart(statusStats domain.AssetStatusStatistics) (string, error) {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Asset Status Distribution",
		}),
	)

	items := []opts.PieData{
		{Name: "Active", Value: statusStats.Active},
		{Name: "Maintenance", Value: statusStats.Maintenance},
		{Name: "Disposed", Value: statusStats.Disposed},
		{Name: "Lost", Value: statusStats.Lost},
	}

	pie.AddSeries("Status", items).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:      opts.Bool(true),
				Formatter: "{b}: {c} ({d}%)",
			}),
		)

	// Create temp file
	tmpFile, err := os.CreateTemp("", "status_chart_*.html")
	if err != nil {
		return "", err
	}
	tmpFile.Close()

	// Render chart to file
	f, err := os.Create(tmpFile.Name())
	if err != nil {
		return "", err
	}
	defer f.Close()

	err = pie.Render(f)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

// generateConditionChart creates a pie chart for condition distribution
func (s *Service) generateConditionChart(conditionStats domain.AssetConditionStatistics) (string, error) {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Asset Condition Distribution",
		}),
	)

	items := []opts.PieData{
		{Name: "Good", Value: conditionStats.Good},
		{Name: "Fair", Value: conditionStats.Fair},
		{Name: "Poor", Value: conditionStats.Poor},
		{Name: "Damaged", Value: conditionStats.Damaged},
	}

	pie.AddSeries("Condition", items).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:      opts.Bool(true),
				Formatter: "{b}: {c} ({d}%)",
			}),
		)

	// Create temp file
	tmpFile, err := os.CreateTemp("", "condition_chart_*.html")
	if err != nil {
		return "", err
	}
	tmpFile.Close()

	// Render chart to file
	f, err := os.Create(tmpFile.Name())
	if err != nil {
		return "", err
	}
	defer f.Close()

	err = pie.Render(f)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}
