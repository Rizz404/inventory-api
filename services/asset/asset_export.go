package asset

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/Rizz404/inventory-api/internal/utils"
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
		data, err := s.exportAssetListToPDF(assetResponses, payload.IncludeDataMatrixImage, langCode)
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
func (s *Service) exportAssetListToPDF(assets []domain.AssetResponse, includeDataMatrix bool, langCode string) ([]byte, error) {
	// Get absolute path for assets
	workDir, _ := os.Getwd()
	logoPath := filepath.Join(workDir, "assets", "images", "company-logo.png")

	// Configure PDF - use default fonts for now
	// TODO: Custom Japanese fonts need proper Maroto v2 configuration
	cfgBuilder := config.NewBuilder().
		WithPageSize(pagesize.A4).
		WithOrientation(orientation.Horizontal).
		WithLeftMargin(10).
		WithTopMargin(15).
		WithRightMargin(10).
		WithBottomMargin(15)

	cfg := cfgBuilder.Build()
	mrt := maroto.New(cfg)

	// Get localized text
	reportTitle := utils.GetLocalizedMessage(utils.PDFAssetListReportKey, langCode)
	generatedOnText := utils.GetLocalizedMessage(utils.PDFGeneratedOnKey, langCode)
	totalAssetsText := utils.GetLocalizedMessage(utils.PDFTotalAssetsKey, langCode)
	assetTagText := utils.GetLocalizedMessage(utils.PDFAssetTagKey, langCode)
	assetNameText := utils.GetLocalizedMessage(utils.PDFAssetNameKey, langCode)
	categoryText := utils.GetLocalizedMessage(utils.PDFCategoryKey, langCode)
	brandText := utils.GetLocalizedMessage(utils.PDFBrandKey, langCode)
	modelText := utils.GetLocalizedMessage(utils.PDFModelKey, langCode)
	statusText := utils.GetLocalizedMessage(utils.PDFStatusKey, langCode)
	conditionText := utils.GetLocalizedMessage(utils.PDFConditionKey, langCode)
	locationText := utils.GetLocalizedMessage(utils.PDFLocationKey, langCode)

	// Color definitions
	headerBgColor := &props.Color{Red: 68, Green: 114, Blue: 196}    // Professional blue
	headerTextColor := &props.Color{Red: 255, Green: 255, Blue: 255} // White
	zebraColor := &props.Color{Red: 242, Green: 242, Blue: 242}      // Light gray

	// Add header with logo and title
	headerRow := row.New(12)

	// Add logo if exists
	if _, err := os.Stat(logoPath); err == nil {
		headerRow.Add(
			col.New(3).Add(
				image.NewFromFile(logoPath, props.Rect{
					Left:    0,
					Top:     0,
					Percent: 80,
					Center:  false,
				}),
			),
		)
		headerRow.Add(
			col.New(9).Add(
				text.New(reportTitle, props.Text{
					Size:  14,
					Style: fontstyle.Bold,
					Align: align.Right,
					Top:   4,
				}),
			),
		)
	} else {
		// No logo, just centered title
		headerRow.Add(
			col.New(12).Add(
				text.New(reportTitle, props.Text{
					Size:  14,
					Style: fontstyle.Bold,
					Align: align.Center,
					Top:   4,
				}),
			),
		)
	}

	mrt.AddRows(headerRow)

	// Add subtitle with generation date
	mrt.AddRows(
		row.New(8).Add(
			col.New(12).Add(
				text.New(fmt.Sprintf("%s: %s", generatedOnText, time.Now().Format("2006-01-02 15:04:05")), props.Text{
					Size:  9,
					Align: align.Center,
				}),
			),
		),
	)

	// Add spacing
	mrt.AddRows(row.New(3))

	// Add table header with professional styling
	tableHeaderRow := row.New(10)
	tableHeaderRow.Add(
		col.New(2).Add(text.New(assetTagText, props.Text{
			Style: fontstyle.Bold,
			Size:  9,
			Align: align.Center,
			Color: headerTextColor,
		})).WithStyle(&props.Cell{BackgroundColor: headerBgColor}),
	)
	tableHeaderRow.Add(
		col.New(3).Add(text.New(assetNameText, props.Text{
			Style: fontstyle.Bold,
			Size:  9,
			Align: align.Center,
			Color: headerTextColor,
		})).WithStyle(&props.Cell{BackgroundColor: headerBgColor}),
	)
	tableHeaderRow.Add(
		col.New(2).Add(text.New(categoryText, props.Text{
			Style: fontstyle.Bold,
			Size:  9,
			Align: align.Center,
			Color: headerTextColor,
		})).WithStyle(&props.Cell{BackgroundColor: headerBgColor}),
	)
	tableHeaderRow.Add(
		col.New(1).Add(text.New(brandText, props.Text{
			Style: fontstyle.Bold,
			Size:  9,
			Align: align.Center,
			Color: headerTextColor,
		})).WithStyle(&props.Cell{BackgroundColor: headerBgColor}),
	)
	tableHeaderRow.Add(
		col.New(1).Add(text.New(modelText, props.Text{
			Style: fontstyle.Bold,
			Size:  9,
			Align: align.Center,
			Color: headerTextColor,
		})).WithStyle(&props.Cell{BackgroundColor: headerBgColor}),
	)
	tableHeaderRow.Add(
		col.New(1).Add(text.New(statusText, props.Text{
			Style: fontstyle.Bold,
			Size:  9,
			Align: align.Center,
			Color: headerTextColor,
		})).WithStyle(&props.Cell{BackgroundColor: headerBgColor}),
	)
	tableHeaderRow.Add(
		col.New(1).Add(text.New(conditionText, props.Text{
			Style: fontstyle.Bold,
			Size:  9,
			Align: align.Center,
			Color: headerTextColor,
		})).WithStyle(&props.Cell{BackgroundColor: headerBgColor}),
	)
	tableHeaderRow.Add(
		col.New(1).Add(text.New(locationText, props.Text{
			Style: fontstyle.Bold,
			Size:  9,
			Align: align.Center,
			Color: headerTextColor,
		})).WithStyle(&props.Cell{BackgroundColor: headerBgColor}),
	)
	mrt.AddRows(tableHeaderRow)

	// Add asset data rows with zebra striping and conditional formatting
	for i, asset := range assets {
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

		// Zebra striping: alternate row background color
		var rowBgColor *props.Color
		if i%2 == 1 {
			rowBgColor = zebraColor
		}

		// Conditional formatting for status and condition
		statusColor := getStatusColor(asset.Status)
		conditionColor := getConditionColor(asset.Condition)

		dataRow.Add(
			col.New(2).Add(text.New(asset.AssetTag, props.Text{
				Size: 8,
			})).WithStyle(&props.Cell{BackgroundColor: rowBgColor}),
		)
		dataRow.Add(
			col.New(3).Add(text.New(asset.AssetName, props.Text{
				Size: 8,
			})).WithStyle(&props.Cell{BackgroundColor: rowBgColor}),
		)
		dataRow.Add(
			col.New(2).Add(text.New(categoryName, props.Text{
				Size: 8,
			})).WithStyle(&props.Cell{BackgroundColor: rowBgColor}),
		)
		dataRow.Add(
			col.New(1).Add(text.New(brand, props.Text{
				Size: 8,
			})).WithStyle(&props.Cell{BackgroundColor: rowBgColor}),
		)
		dataRow.Add(
			col.New(1).Add(text.New(model, props.Text{
				Size: 8,
			})).WithStyle(&props.Cell{BackgroundColor: rowBgColor}),
		)
		dataRow.Add(
			col.New(1).Add(text.New(string(asset.Status), props.Text{
				Size:  8,
				Color: statusColor,
				Style: fontstyle.Bold,
			})).WithStyle(&props.Cell{BackgroundColor: rowBgColor}),
		)
		dataRow.Add(
			col.New(1).Add(text.New(string(asset.Condition), props.Text{
				Size:  8,
				Color: conditionColor,
				Style: fontstyle.Bold,
			})).WithStyle(&props.Cell{BackgroundColor: rowBgColor}),
		)
		dataRow.Add(
			col.New(1).Add(text.New(locationName, props.Text{
				Size: 8,
			})).WithStyle(&props.Cell{BackgroundColor: rowBgColor}),
		)
		mrt.AddRows(dataRow)
	}

	// Add spacing
	mrt.AddRows(row.New(3))

	// Add footer with total count
	mrt.AddRows(
		row.New(10).Add(
			col.New(12).Add(
				text.New(fmt.Sprintf("%s: %d", totalAssetsText, len(assets)), props.Text{
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

// getStatusColor returns color based on asset status
func getStatusColor(status domain.AssetStatus) *props.Color {
	switch status {
	case domain.StatusActive:
		return &props.Color{Red: 34, Green: 139, Blue: 34} // Forest green
	case domain.StatusMaintenance:
		return &props.Color{Red: 255, Green: 140, Blue: 0} // Dark orange
	case domain.StatusDisposed:
		return &props.Color{Red: 128, Green: 128, Blue: 128} // Gray
	case domain.StatusLost:
		return &props.Color{Red: 220, Green: 20, Blue: 60} // Crimson red
	default:
		return &props.Color{Red: 0, Green: 0, Blue: 0} // Black
	}
}

// getConditionColor returns color based on asset condition
func getConditionColor(condition domain.AssetCondition) *props.Color {
	switch condition {
	case domain.ConditionGood:
		return &props.Color{Red: 34, Green: 139, Blue: 34} // Forest green
	case domain.ConditionFair:
		return &props.Color{Red: 255, Green: 215, Blue: 0} // Gold
	case domain.ConditionPoor:
		return &props.Color{Red: 255, Green: 140, Blue: 0} // Dark orange
	case domain.ConditionDamaged:
		return &props.Color{Red: 220, Green: 20, Blue: 60} // Crimson red
	default:
		return &props.Color{Red: 0, Green: 0, Blue: 0} // Black
	}
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
