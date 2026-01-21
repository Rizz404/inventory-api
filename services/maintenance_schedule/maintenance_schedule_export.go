package maintenance_schedule

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

// ExportMaintenanceScheduleList exports maintenance schedule list to PDF or Excel format
func (s *Service) ExportMaintenanceScheduleList(ctx context.Context, payload domain.ExportMaintenanceScheduleListPayload, params domain.MaintenanceScheduleParams, langCode string) ([]byte, string, error) {
	if payload.SearchQuery != nil {
		params.SearchQuery = payload.SearchQuery
	}
	if payload.Filters != nil {
		params.Filters = payload.Filters
	}
	if payload.Sort != nil {
		params.Sort = payload.Sort
	}

	schedules, err := s.Repo.GetMaintenanceSchedulesForExport(ctx, params, langCode)
	if err != nil {
		return nil, "", err
	}

	// Convert to responses (includes translations)
	scheduleResponses := mapper.MaintenanceSchedulesToResponses(schedules, langCode)

	switch payload.Format {
	case domain.ExportFormatPDF:
		data, err := s.exportMaintenanceScheduleListToPDF(scheduleResponses, langCode)
		if err != nil {
			return nil, "", domain.ErrInternal(err)
		}
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		filename := fmt.Sprintf("maintenance_schedules_%s.pdf", timestamp)
		return data, filename, nil

	case domain.ExportFormatExcel:
		data, err := s.exportMaintenanceScheduleListToExcel(scheduleResponses)
		if err != nil {
			return nil, "", domain.ErrInternal(err)
		}
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		filename := fmt.Sprintf("maintenance_schedules_%s.xlsx", timestamp)
		return data, filename, nil

	default:
		return nil, "", domain.ErrBadRequest("Invalid export format")
	}
}

// exportMaintenanceScheduleListToPDF generates PDF file for maintenance schedule list using gopdf
func (s *Service) exportMaintenanceScheduleListToPDF(schedules []domain.MaintenanceScheduleResponse, langCode string) ([]byte, error) {
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
	reportTitle := utils.GetLocalizedMessage(utils.PDFMaintenanceScheduleReportKey, langCode)
	generatedOnText := utils.GetLocalizedMessage(utils.PDFAssetGeneratedOnKey, langCode)
	totalSchedulesText := utils.GetLocalizedMessage(utils.PDFMaintenanceScheduleTotalKey, langCode)
	assetTagText := utils.GetLocalizedMessage(utils.PDFAssetAssetTagKey, langCode)
	assetNameText := utils.GetLocalizedMessage(utils.PDFAssetAssetNameKey, langCode)
	titleText := utils.GetLocalizedMessage(utils.PDFMaintenanceScheduleTitleKey, langCode)
	maintenanceTypeText := utils.GetLocalizedMessage(utils.PDFMaintenanceScheduleTypeKey, langCode)
	nextScheduledDateText := utils.GetLocalizedMessage(utils.PDFMaintenanceScheduleNextDateKey, langCode)
	isRecurringText := utils.GetLocalizedMessage(utils.PDFMaintenanceScheduleRecurringKey, langCode)
	stateText := utils.GetLocalizedMessage(utils.PDFMaintenanceScheduleStateKey, langCode)
	estimatedCostText := utils.GetLocalizedMessage(utils.PDFMaintenanceScheduleCostKey, langCode)

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
	colWidths := []float64{70, 100, 120, 90, 80, 70, 70, 80}
	headers := []string{assetTagText, assetNameText, titleText, maintenanceTypeText, nextScheduledDateText, isRecurringText, stateText, estimatedCostText}

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

	for i, schedule := range schedules {
		maxLines := 1

		nameLines := wrapText(schedule.Asset.AssetName, colWidths[1])
		if len(nameLines) > maxLines {
			maxLines = len(nameLines)
		}

		titleLines := wrapText(schedule.Title, colWidths[2])
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
		pdf.Cell(nil, schedule.Asset.AssetTag)
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

		// Maintenance Type
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, string(schedule.MaintenanceType))
		x += colWidths[3]

		// Next Scheduled Date
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, schedule.NextScheduledDate.Format("2006-01-02"))
		x += colWidths[4]

		// Is Recurring
		recurringText := "No"
		if schedule.IsRecurring {
			recurringText = "Yes"
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, recurringText)
		x += colWidths[5]

		// State (with color)
		switch schedule.State {
		case domain.StateActive:
			pdf.SetTextColor(34, 139, 34)
		case domain.StatePaused:
			pdf.SetTextColor(255, 140, 0)
		case domain.StateStopped:
			pdf.SetTextColor(220, 20, 60)
		case domain.StateCompleted:
			pdf.SetTextColor(128, 128, 128)
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, string(schedule.State))
		pdf.SetTextColor(0, 0, 0)
		x += colWidths[6]

		// Estimated Cost
		costText := "-"
		if schedule.EstimatedCost != nil && schedule.EstimatedCost.Valid {
			value, _ := schedule.EstimatedCost.Float64()
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
	pdf.Cell(nil, fmt.Sprintf("%s: %d", totalSchedulesText, len(schedules)))

	var buf bytes.Buffer
	if err := pdf.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// exportMaintenanceScheduleListToExcel generates Excel file for maintenance schedule list
func (s *Service) exportMaintenanceScheduleListToExcel(schedules []domain.MaintenanceScheduleResponse) ([]byte, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error closing Excel file:", err)
		}
	}()

	sheetName := "Maintenance Schedules"
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
		"Asset Tag", "Asset Name", "Title", "Maintenance Type",
		"Next Scheduled Date", "Last Executed Date", "Is Recurring", "Interval",
		"State", "Auto Complete", "Estimated Cost", "Created By", "Description",
	}

	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	for row, schedule := range schedules {
		rowNum := row + 2

		lastExecutedDate := ""
		if schedule.LastExecutedDate != nil {
			lastExecutedDate = schedule.LastExecutedDate.Format("2006-01-02")
		}

		isRecurring := "No"
		if schedule.IsRecurring {
			isRecurring = "Yes"
		}

		interval := "-"
		if schedule.IntervalValue != nil && schedule.IntervalUnit != nil {
			interval = fmt.Sprintf("%d %s", *schedule.IntervalValue, *schedule.IntervalUnit)
		}

		autoComplete := "No"
		if schedule.AutoComplete {
			autoComplete = "Yes"
		}

		estimatedCost := ""
		if schedule.EstimatedCost != nil && schedule.EstimatedCost.Valid {
			value, _ := schedule.EstimatedCost.Float64()
			estimatedCost = fmt.Sprintf("%.2f", value)
		}

		description := ""
		if schedule.Description != nil {
			description = *schedule.Description
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), schedule.Asset.AssetTag)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), schedule.Asset.AssetName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), schedule.Title)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), string(schedule.MaintenanceType))
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), schedule.NextScheduledDate.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), lastExecutedDate)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), isRecurring)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), interval)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", rowNum), string(schedule.State))
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", rowNum), autoComplete)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", rowNum), estimatedCost)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", rowNum), schedule.CreatedBy.FullName)
		f.SetCellValue(sheetName, fmt.Sprintf("M%d", rowNum), description)
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
