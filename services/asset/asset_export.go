package asset

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"github.com/Rizz404/inventory-api/internal/utils"

	// ! Old Maroto chart imports - commented for reference
	// "github.com/go-echarts/go-echarts/v2/charts"
	// "github.com/go-echarts/go-echarts/v2/opts"
	"github.com/signintech/gopdf"
	"github.com/xuri/excelize/v2"
	// Maroto v2 - commented out, kept for reference if needed later
	// "github.com/johnfercher/maroto/v2"
	// "github.com/johnfercher/maroto/v2/pkg/components/col"
	// "github.com/johnfercher/maroto/v2/pkg/components/image"
	// "github.com/johnfercher/maroto/v2/pkg/components/row"
	// "github.com/johnfercher/maroto/v2/pkg/components/text"
	// "github.com/johnfercher/maroto/v2/pkg/config"
	// "github.com/johnfercher/maroto/v2/pkg/consts/align"
	// "github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	// "github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	// "github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	// "github.com/johnfercher/maroto/v2/pkg/core/entity"
	// "github.com/johnfercher/maroto/v2/pkg/props"
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
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		filename := fmt.Sprintf("asset_list_%s.pdf", timestamp)
		return data, filename, nil

	case domain.ExportFormatExcel:
		data, err := s.exportAssetListToExcel(assetResponses)
		if err != nil {
			return nil, "", domain.ErrInternal(err)
		}
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		filename := fmt.Sprintf("asset_list_%s.xlsx", timestamp)
		return data, filename, nil

	default:
		return nil, "", domain.ErrBadRequest("Invalid export format")
	}
}

// exportAssetListToPDF generates PDF file for asset list using gopdf (supports Unicode/CJK)
func (s *Service) exportAssetListToPDF(assets []domain.AssetResponse, includeDataMatrix bool, langCode string) ([]byte, error) {
	// Get absolute path for fonts and logo
	workDir, _ := os.Getwd()
	fontRegularPath := filepath.Join(workDir, "assets", "fonts", "NotoSansJP-Regular.ttf")
	fontBoldPath := filepath.Join(workDir, "assets", "fonts", "NotoSansJP-Bold.ttf")
	logoPath := filepath.Join(workDir, "assets", "images", "company-logo.png")

	// Initialize gopdf
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		PageSize: *gopdf.PageSizeA4Landscape, // A4 Landscape: 842 x 595
		Unit:     gopdf.Unit_PT,
	})
	pdf.AddPage()

	// Load fonts - Always use Noto Sans for Unicode support (works for all languages)
	if err := pdf.AddTTFFont("noto-regular", fontRegularPath); err != nil {
		return nil, fmt.Errorf("failed to load regular font: %w", err)
	}
	if err := pdf.AddTTFFont("noto-bold", fontBoldPath); err != nil {
		return nil, fmt.Errorf("failed to load bold font: %w", err)
	}

	pdf.SetFont("noto-regular", "", 10)

	// Get localized text
	reportTitle := utils.GetLocalizedMessage(utils.PDFAssetListReportKey, langCode)
	generatedOnText := utils.GetLocalizedMessage(utils.PDFAssetGeneratedOnKey, langCode)
	totalAssetsText := utils.GetLocalizedMessage(utils.PDFAssetTotalAssetsKey, langCode)
	assetTagText := utils.GetLocalizedMessage(utils.PDFAssetAssetTagKey, langCode)
	assetNameText := utils.GetLocalizedMessage(utils.PDFAssetAssetNameKey, langCode)
	categoryText := utils.GetLocalizedMessage(utils.PDFAssetCategoryKey, langCode)
	brandText := utils.GetLocalizedMessage(utils.PDFAssetBrandKey, langCode)
	modelText := utils.GetLocalizedMessage(utils.PDFAssetModelKey, langCode)
	statusText := utils.GetLocalizedMessage(utils.PDFAssetStatusKey, langCode)
	conditionText := utils.GetLocalizedMessage(utils.PDFAssetConditionKey, langCode)
	locationText := utils.GetLocalizedMessage(utils.PDFAssetLocationKey, langCode)

	// Page setup (A4 Landscape: 842 x 595 points)
	marginLeft := 30.0
	marginTop := 50.0
	pageWidth := 842.0
	pageHeight := 595.0
	contentWidth := pageWidth - (marginLeft * 2)

	// Add company logo if exists
	currentY := marginTop
	if _, err := os.Stat(logoPath); err == nil {
		// Logo dimensions: 60x60 with proper Rect
		rect := &gopdf.Rect{W: 60, H: 60}
		pdf.Image(logoPath, marginLeft, currentY-10, rect)

		// Title next to logo
		pdf.SetFont("noto-bold", "", 16)
		pdf.SetX(marginLeft + 70)
		pdf.SetY(currentY + 15)
		pdf.Cell(nil, reportTitle)

		currentY += 50
	} else {
		// No logo, centered title
		pdf.SetFont("noto-bold", "", 16)
		titleWidth, _ := pdf.MeasureTextWidth(reportTitle)
		pdf.SetX((pageWidth - titleWidth) / 2)
		pdf.SetY(currentY)
		pdf.Cell(nil, reportTitle)

		currentY += 30
	}

	// Subtitle with date
	pdf.SetFont("noto-regular", "", 10)
	dateText := fmt.Sprintf("%s: %s", generatedOnText, time.Now().Format("2006-01-02 15:04:05"))
	dateWidth, _ := pdf.MeasureTextWidth(dateText)
	pdf.SetX((pageWidth - dateWidth) / 2)
	pdf.SetY(currentY)
	pdf.Cell(nil, dateText)

	currentY += 25

	// Table setup
	startY := currentY

	// Column widths (optimized for A4 landscape: 782 total usable width)
	colWidths := []float64{70, 140, 100, 80, 80, 80, 80, 100} // Total: 730
	headers := []string{assetTagText, assetNameText, categoryText, brandText, modelText, statusText, conditionText, locationText}

	// Draw table header with proper background
	pdf.SetFillColor(68, 114, 196) // Blue background
	pdf.RectFromUpperLeftWithStyle(marginLeft, startY, contentWidth, 25, "F")

	pdf.SetTextColor(255, 255, 255) // White text
	pdf.SetFont("noto-bold", "", 10)

	x := marginLeft
	y := startY
	for i, header := range headers {
		pdf.SetX(x + 3)
		pdf.SetY(y + 8)
		pdf.Cell(nil, header)
		x += colWidths[i]
	}

	// Reset for data rows
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("noto-regular", "", 9)
	y += 25

	// Helper function to wrap text if needed
	wrapText := func(text string, maxWidth float64) []string {
		words := []string{}
		currentLine := ""

		for _, char := range text {
			testLine := currentLine + string(char)
			width, _ := pdf.MeasureTextWidth(testLine)

			if width > maxWidth-10 { // -10 for padding
				if currentLine != "" {
					words = append(words, currentLine)
				}
				currentLine = string(char)
			} else {
				currentLine = testLine
			}
		}
		if currentLine != "" {
			words = append(words, currentLine)
		}

		if len(words) == 0 {
			return []string{text}
		}
		return words
	}

	for i, asset := range assets {
		// Calculate row height based on content (check for multi-line text)
		maxLines := 1

		// Check Asset Name
		pdf.SetFont("noto-regular", "", 9)
		nameLines := wrapText(asset.AssetName, colWidths[1])
		if len(nameLines) > maxLines {
			maxLines = len(nameLines)
		}

		// Check Category
		if asset.Category != nil {
			catLines := wrapText(asset.Category.CategoryName, colWidths[2])
			if len(catLines) > maxLines {
				maxLines = len(catLines)
			}
		}

		rowHeight := float64(maxLines) * 14.0
		if rowHeight < 20 {
			rowHeight = 20
		}

		// Check if need new page
		if y+rowHeight > pageHeight-40 {
			pdf.AddPage()

			// Redraw header on new page
			y = marginTop
			pdf.SetFillColor(68, 114, 196)
			pdf.RectFromUpperLeftWithStyle(marginLeft, y, contentWidth, 25, "F")
			pdf.SetTextColor(255, 255, 255)
			pdf.SetFont("noto-bold", "", 10)

			x = marginLeft
			for j, header := range headers {
				pdf.SetX(x + 3)
				pdf.SetY(y + 8)
				pdf.Cell(nil, header)
				x += colWidths[j]
			}

			y += 25
			pdf.SetTextColor(0, 0, 0)
			pdf.SetFont("noto-regular", "", 9)
		}

		// Zebra striping
		if i%2 == 1 {
			pdf.SetFillColor(242, 242, 242)
			pdf.RectFromUpperLeftWithStyle(marginLeft, y, contentWidth, rowHeight, "F")
		}

		x = marginLeft
		cellY := y + 6

		// Asset Tag
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, asset.AssetTag)
		x += colWidths[0]

		// Asset Name (multi-line support)
		for lineIdx, line := range nameLines {
			pdf.SetX(x + 3)
			pdf.SetY(cellY + float64(lineIdx)*12)
			pdf.Cell(nil, line)
		}
		x += colWidths[1]

		// Category (multi-line support)
		categoryName := ""
		if asset.Category != nil {
			categoryName = asset.Category.CategoryName
			catLines := wrapText(categoryName, colWidths[2])
			for lineIdx, line := range catLines {
				pdf.SetX(x + 3)
				pdf.SetY(cellY + float64(lineIdx)*12)
				pdf.Cell(nil, line)
			}
		}
		x += colWidths[2]

		// Brand
		brand := ""
		if asset.Brand != nil {
			brand = *asset.Brand
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, brand)
		x += colWidths[3]

		// Model
		model := ""
		if asset.Model != nil {
			model = *asset.Model
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, model)
		x += colWidths[4]

		// Status (with color)
		switch asset.Status {
		case domain.StatusActive:
			pdf.SetTextColor(34, 139, 34) // Forest green
		case domain.StatusMaintenance:
			pdf.SetTextColor(255, 140, 0) // Dark orange
		case domain.StatusDisposed:
			pdf.SetTextColor(128, 128, 128) // Gray
		case domain.StatusLost:
			pdf.SetTextColor(220, 20, 60) // Crimson red
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, string(asset.Status))
		pdf.SetTextColor(0, 0, 0) // Reset color
		x += colWidths[5]

		// Condition (with color)
		switch asset.Condition {
		case domain.ConditionGood:
			pdf.SetTextColor(34, 139, 34) // Forest green
		case domain.ConditionFair:
			pdf.SetTextColor(255, 215, 0) // Gold
		case domain.ConditionPoor:
			pdf.SetTextColor(255, 140, 0) // Dark orange
		case domain.ConditionDamaged:
			pdf.SetTextColor(220, 20, 60) // Crimson red
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, string(asset.Condition))
		pdf.SetTextColor(0, 0, 0) // Reset color
		x += colWidths[6]

		// Location
		locationName := ""
		if asset.Location != nil {
			locationName = asset.Location.LocationName
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, locationName)

		y += rowHeight
	}

	// Footer - Total count
	y += 15
	if y > pageHeight-40 {
		pdf.AddPage()
		y = marginTop
	}
	pdf.SetFont("noto-bold", "", 11)
	pdf.SetX(marginLeft)
	pdf.SetY(y)
	pdf.Cell(nil, fmt.Sprintf("%s: %d", totalAssetsText, len(assets)))

	// Get PDF bytes
	var buf bytes.Buffer
	if err := pdf.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ! Old Maroto color helpers - commented for reference
// func getStatusColor(status domain.AssetStatus) *props.Color {
// 	switch status {
// 	case domain.AssetStatusAvailable:
// 		return &props.Color{Red: 34, Green: 139, Blue: 34} // Forest green
// 	case domain.AssetStatusInUse:
// 		return &props.Color{Red: 255, Green: 140, Blue: 0} // Dark orange
// 	case domain.AssetStatusUnderMaintenance:
// 		return &props.Color{Red: 128, Green: 128, Blue: 128} // Gray
// 	case domain.AssetStatusRetired:
// 		return &props.Color{Red: 220, Green: 20, Blue: 60} // Crimson red
// 	default:
// 		return &props.Color{Red: 0, Green: 0, Blue: 0} // Black
// 	}
// }
//
// func getConditionColor(condition domain.AssetCondition) *props.Color {
// 	switch condition {
// 	case domain.AssetConditionExcellent:
// 		return &props.Color{Red: 34, Green: 139, Blue: 34} // Forest green
// 	case domain.AssetConditionGood:
// 		return &props.Color{Red: 255, Green: 215, Blue: 0} // Gold
// 	case domain.AssetConditionFair:
// 		return &props.Color{Red: 255, Green: 140, Blue: 0} // Dark orange
// 	case domain.AssetConditionPoor:
// 		return &props.Color{Red: 220, Green: 20, Blue: 60} // Crimson red
// 	default:
// 		return &props.Color{Red: 0, Green: 0, Blue: 0} // Black
// 	}
// }

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

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("asset_statistics_%s.pdf", timestamp)
	return data, filename, nil
}

// ! Old Maroto implementation - Statistics PDF export - commented for reference
// TODO: Implement statistics PDF using gopdf
// exportAssetStatisticsToPDF generates PDF file with charts for statistics
func (s *Service) exportAssetStatisticsToPDF(stats domain.AssetStatistics) ([]byte, error) {
	return nil, fmt.Errorf("statistics PDF export not yet implemented with gopdf")
}

// ! OLD MAROTO CODE BELOW - commented for future reference
/*
func (s *Service) exportAssetStatisticsToPDF_OLD(stats domain.AssetStatistics) ([]byte, error) {
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
*/
