package maintenance_record

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

	"github.com/signintech/gopdf"
	"github.com/xuri/excelize/v2"
)

// ExportMaintenanceRecordList exports maintenance record list to PDF or Excel format
func (s *Service) ExportMaintenanceRecordList(ctx context.Context, payload domain.ExportMaintenanceRecordListPayload, params domain.MaintenanceRecordParams, langCode string) ([]byte, string, error) {
	// Override params with payload if provided
	if payload.SearchQuery != nil {
		params.SearchQuery = payload.SearchQuery
	}
	if payload.Filters != nil {
		params.Filters = payload.Filters
	}
	if payload.Sort != nil {
		params.Sort = payload.Sort
	}

	// Get records without pagination
	records, err := s.Repo.GetMaintenanceRecordsForExport(ctx, params, langCode)
	if err != nil {
		return nil, "", err
	}

	// Convert to responses
	recordResponses := mapper.MaintenanceRecordsToListResponses(records, langCode)

	switch payload.Format {
	case domain.ExportFormatPDF:
		data, err := s.exportMaintenanceRecordListToPDF(recordResponses, langCode)
		if err != nil {
			return nil, "", domain.ErrInternal(err)
		}
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		filename := fmt.Sprintf("maintenance_records_%s.pdf", timestamp)
		return data, filename, nil

	case domain.ExportFormatExcel:
		data, err := s.exportMaintenanceRecordListToExcel(recordResponses)
		if err != nil {
			return nil, "", domain.ErrInternal(err)
		}
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		filename := fmt.Sprintf("maintenance_records_%s.xlsx", timestamp)
		return data, filename, nil

	default:
		return nil, "", domain.ErrBadRequest("Invalid export format")
	}
}

// exportMaintenanceRecordListToPDF generates PDF file for maintenance record list using gopdf
func (s *Service) exportMaintenanceRecordListToPDF(records []domain.MaintenanceRecordListResponse, langCode string) ([]byte, error) {
	workDir, _ := os.Getwd()
	fontRegularPath := filepath.Join(workDir, "assets", "fonts", "NotoSansJP-Regular.ttf")
	fontBoldPath := filepath.Join(workDir, "assets", "fonts", "NotoSansJP-Bold.ttf")
	logoPath := filepath.Join(workDir, "assets", "images", "company-logo.png")

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		PageSize: *gopdf.PageSizeA4Landscape,
		Unit:     gopdf.Unit_PT,
	})
	pdf.AddPage()

	if err := pdf.AddTTFFont("noto-regular", fontRegularPath); err != nil {
		return nil, fmt.Errorf("failed to load regular font: %w", err)
	}
	if err := pdf.AddTTFFont("noto-bold", fontBoldPath); err != nil {
		return nil, fmt.Errorf("failed to load bold font: %w", err)
	}

	pdf.SetFont("noto-regular", "", 10)

	// Localized text
	reportTitle := utils.GetLocalizedMessage(utils.PDFMaintenanceRecordReportKey, langCode)
	generatedOnText := utils.GetLocalizedMessage(utils.PDFAssetGeneratedOnKey, langCode)
	totalRecordsText := utils.GetLocalizedMessage(utils.PDFMaintenanceRecordTotalKey, langCode)
	assetTagText := utils.GetLocalizedMessage(utils.PDFAssetAssetTagKey, langCode)
	assetNameText := utils.GetLocalizedMessage(utils.PDFAssetAssetNameKey, langCode)
	titleText := utils.GetLocalizedMessage(utils.PDFMaintenanceRecordTitleKey, langCode)
	maintenanceDateText := utils.GetLocalizedMessage(utils.PDFMaintenanceRecordDateKey, langCode)
	completionDateText := utils.GetLocalizedMessage(utils.PDFMaintenanceRecordCompletionKey, langCode)
	performerText := utils.GetLocalizedMessage(utils.PDFMaintenanceRecordPerformerKey, langCode)
	resultText := utils.GetLocalizedMessage(utils.PDFMaintenanceRecordResultKey, langCode)
	costText := utils.GetLocalizedMessage(utils.PDFMaintenanceRecordCostKey, langCode)

	marginLeft := 30.0
	marginTop := 50.0
	pageWidth := 842.0
	pageHeight := 595.0
	contentWidth := pageWidth - (marginLeft * 2)

	currentY := marginTop
	if _, err := os.Stat(logoPath); err == nil {
		rect := &gopdf.Rect{W: 60, H: 60}
		pdf.Image(logoPath, marginLeft, currentY-10, rect)

		pdf.SetFont("noto-bold", "", 16)
		pdf.SetX(marginLeft + 70)
		pdf.SetY(currentY + 15)
		pdf.Cell(nil, reportTitle)

		currentY += 50
	} else {
		pdf.SetFont("noto-bold", "", 16)
		titleWidth, _ := pdf.MeasureTextWidth(reportTitle)
		pdf.SetX((pageWidth - titleWidth) / 2)
		pdf.SetY(currentY)
		pdf.Cell(nil, reportTitle)

		currentY += 30
	}

	pdf.SetFont("noto-regular", "", 10)
	dateText := fmt.Sprintf("%s: %s", generatedOnText, time.Now().Format("2006-01-02 15:04:05"))
	dateWidth, _ := pdf.MeasureTextWidth(dateText)
	pdf.SetX((pageWidth - dateWidth) / 2)
	pdf.SetY(currentY)
	pdf.Cell(nil, dateText)

	currentY += 25

	startY := currentY
	colWidths := []float64{70, 110, 140, 80, 80, 100, 70, 70}
	headers := []string{assetTagText, assetNameText, titleText, maintenanceDateText, completionDateText, performerText, resultText, costText}

	pdf.SetFillColor(68, 114, 196)
	pdf.RectFromUpperLeftWithStyle(marginLeft, startY, contentWidth, 25, "F")

	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("noto-bold", "", 9)

	x := marginLeft
	y := startY
	for i, header := range headers {
		pdf.SetX(x + 3)
		pdf.SetY(y + 8)
		pdf.Cell(nil, header)
		x += colWidths[i]
	}

	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("noto-regular", "", 8)
	y += 25

	wrapText := func(text string, maxWidth float64) []string {
		words := []string{}
		currentLine := ""

		for _, char := range text {
			testLine := currentLine + string(char)
			width, _ := pdf.MeasureTextWidth(testLine)

			if width > maxWidth-10 {
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

	for i, record := range records {
		maxLines := 1

		nameLines := wrapText(record.Asset.AssetName, colWidths[1])
		if len(nameLines) > maxLines {
			maxLines = len(nameLines)
		}

		titleLines := wrapText(record.Title, colWidths[2])
		if len(titleLines) > maxLines {
			maxLines = len(titleLines)
		}

		rowHeight := float64(maxLines) * 12.0
		if rowHeight < 18 {
			rowHeight = 18
		}

		if y+rowHeight > pageHeight-40 {
			pdf.AddPage()
			y = marginTop

			pdf.SetFillColor(68, 114, 196)
			pdf.RectFromUpperLeftWithStyle(marginLeft, y, contentWidth, 25, "F")
			pdf.SetTextColor(255, 255, 255)
			pdf.SetFont("noto-bold", "", 9)

			x = marginLeft
			for j, header := range headers {
				pdf.SetX(x + 3)
				pdf.SetY(y + 8)
				pdf.Cell(nil, header)
				x += colWidths[j]
			}

			y += 25
			pdf.SetTextColor(0, 0, 0)
			pdf.SetFont("noto-regular", "", 8)
		}

		if i%2 == 1 {
			pdf.SetFillColor(242, 242, 242)
			pdf.RectFromUpperLeftWithStyle(marginLeft, y, contentWidth, rowHeight, "F")
		}

		x = marginLeft
		cellY := y + 5

		// Asset Tag
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, record.Asset.AssetTag)
		x += colWidths[0]

		// Asset Name (multi-line)
		for lineIdx, line := range nameLines {
			pdf.SetX(x + 3)
			pdf.SetY(cellY + float64(lineIdx)*10)
			pdf.Cell(nil, line)
		}
		x += colWidths[1]

		// Title (multi-line)
		for lineIdx, line := range titleLines {
			pdf.SetX(x + 3)
			pdf.SetY(cellY + float64(lineIdx)*10)
			pdf.Cell(nil, line)
		}
		x += colWidths[2]

		// Maintenance Date
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, record.MaintenanceDate.Format("2006-01-02"))
		x += colWidths[3]

		// Completion Date
		completionText := "-"
		if record.CompletionDate != nil {
			completionText = record.CompletionDate.Format("2006-01-02")
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, completionText)
		x += colWidths[4]

		// Performer
		performer := "-"
		if record.PerformedByUser != nil {
			performer = record.PerformedByUser.FullName
		} else if record.PerformedByVendor != nil {
			performer = *record.PerformedByVendor
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, performer)
		x += colWidths[5]

		// Result (with color)
		switch record.Result {
		case domain.ResultSuccess:
			pdf.SetTextColor(34, 139, 34)
		case domain.ResultPartial:
			pdf.SetTextColor(255, 140, 0)
		case domain.ResultFailed:
			pdf.SetTextColor(220, 20, 60)
		case domain.ResultRescheduled:
			pdf.SetTextColor(128, 128, 128)
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, string(record.Result))
		pdf.SetTextColor(0, 0, 0)
		x += colWidths[6]

		// Cost
		costText := "-"
		if record.ActualCost != nil && record.ActualCost.Valid {
			value, _ := record.ActualCost.Float64()
			costText = fmt.Sprintf("$%.2f", value)
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, costText)

		y += rowHeight
	}

	// Footer
	y += 15
	if y > pageHeight-40 {
		pdf.AddPage()
		y = marginTop
	}
	pdf.SetFont("noto-bold", "", 11)
	pdf.SetX(marginLeft)
	pdf.SetY(y)
	pdf.Cell(nil, fmt.Sprintf("%s: %d", totalRecordsText, len(records)))

	var buf bytes.Buffer
	if err := pdf.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// exportMaintenanceRecordListToExcel generates Excel file for maintenance record list
func (s *Service) exportMaintenanceRecordListToExcel(records []domain.MaintenanceRecordListResponse) ([]byte, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error closing Excel file:", err)
		}
	}()

	sheetName := "Maintenance Records"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	f.SetActiveSheet(index)

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

	headers := []string{
		"Asset Tag", "Asset Name", "Title", "Maintenance Date",
		"Completion Date", "Duration (min)", "Performed By User", "Performed By Vendor",
		"Result", "Actual Cost", "Notes",
	}

	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	for row, record := range records {
		rowNum := row + 2

		completionDate := ""
		if record.CompletionDate != nil {
			completionDate = record.CompletionDate.Format("2006-01-02")
		}

		duration := ""
		if record.DurationMinutes != nil {
			duration = fmt.Sprintf("%d", *record.DurationMinutes)
		}

		performedByUser := ""
		if record.PerformedByUser != nil {
			performedByUser = record.PerformedByUser.FullName
		}

		performedByVendor := ""
		if record.PerformedByVendor != nil {
			performedByVendor = *record.PerformedByVendor
		}

		actualCost := ""
		if record.ActualCost != nil && record.ActualCost.Valid {
			value, _ := record.ActualCost.Float64()
			actualCost = fmt.Sprintf("%.2f", value)
		}

		notes := ""
		if record.Notes != nil {
			notes = *record.Notes
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), record.Asset.AssetTag)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), record.Asset.AssetName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), record.Title)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), record.MaintenanceDate.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), completionDate)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), duration)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), performedByUser)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), performedByVendor)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", rowNum), string(record.Result))
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", rowNum), actualCost)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", rowNum), notes)
	}

	for col := 1; col <= len(headers); col++ {
		colName, _ := excelize.ColumnNumberToName(col)
		f.SetColWidth(sheetName, colName, colName, 16)
	}

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
