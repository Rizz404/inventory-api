package scan_log

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

// ExportScanLogList exports scan log list to PDF or Excel format
func (s *Service) ExportScanLogList(ctx context.Context, payload domain.ExportScanLogListPayload, params domain.ScanLogParams, langCode string) ([]byte, string, error) {
	if payload.SearchQuery != nil {
		params.SearchQuery = payload.SearchQuery
	}
	if payload.Filters != nil {
		params.Filters = payload.Filters
	}
	if payload.Sort != nil {
		params.Sort = payload.Sort
	}

	logs, err := s.Repo.GetScanLogsForExport(ctx, params)
	if err != nil {
		return nil, "", err
	}

	// Convert to responses
	logResponses := mapper.ScanLogsToListResponses(logs)

	switch payload.Format {
	case domain.ExportFormatPDF:
		data, err := s.exportScanLogListToPDF(logResponses, langCode)
		if err != nil {
			return nil, "", domain.ErrInternal(err)
		}
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		filename := fmt.Sprintf("scan_logs_%s.pdf", timestamp)
		return data, filename, nil

	case domain.ExportFormatExcel:
		data, err := s.exportScanLogListToExcel(logResponses)
		if err != nil {
			return nil, "", domain.ErrInternal(err)
		}
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		filename := fmt.Sprintf("scan_logs_%s.xlsx", timestamp)
		return data, filename, nil

	default:
		return nil, "", domain.ErrBadRequest("Invalid export format")
	}
}

// exportScanLogListToPDF generates PDF file for scan log list using gopdf
func (s *Service) exportScanLogListToPDF(logs []domain.ScanLogListResponse, langCode string) ([]byte, error) {
	workDir, _ := os.Getwd()
	fontRegularPath := filepath.Join(workDir, "assets", "fonts", "NotoSansJP-Regular.ttf")
	fontBoldPath := filepath.Join(workDir, "assets", "fonts", "NotoSansJP-Bold.ttf")
	logoPath := filepath.Join(workDir, "assets", "images", "fts-logo.png")

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
	reportTitle := utils.GetLocalizedMessage(utils.PDFScanLogReportKey, langCode)
	generatedOnText := utils.GetLocalizedMessage(utils.PDFAssetGeneratedOnKey, langCode)
	totalLogsText := utils.GetLocalizedMessage(utils.PDFScanLogTotalKey, langCode)
	scannedValueText := utils.GetLocalizedMessage(utils.PDFScanLogScannedValueKey, langCode)
	scanMethodText := utils.GetLocalizedMessage(utils.PDFScanLogMethodKey, langCode)
	scanTimestampText := utils.GetLocalizedMessage(utils.PDFScanLogTimestampKey, langCode)
	scanResultText := utils.GetLocalizedMessage(utils.PDFScanLogResultKey, langCode)
	scannedByText := utils.GetLocalizedMessage(utils.PDFScanLogScannedByKey, langCode)
	coordinatesText := utils.GetLocalizedMessage(utils.PDFScanLogCoordinatesKey, langCode)

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
	colWidths := []float64{120, 100, 120, 90, 110, 120}
	headers := []string{scannedValueText, scanMethodText, scanTimestampText, scanResultText, scannedByText, coordinatesText}

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

	for i, log := range logs {
		maxLines := 1

		valueLines := wrapText(log.ScannedValue, colWidths[0])
		if len(valueLines) > maxLines {
			maxLines = len(valueLines)
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

		// Scanned Value (multi-line)
		for lineIdx, line := range valueLines {
			pdf.SetX(x + 3)
			pdf.SetY(cellY + float64(lineIdx)*10)
			pdf.Cell(nil, line)
		}
		x += colWidths[0]

		// Scan Method
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, string(log.ScanMethod))
		x += colWidths[1]

		// Scan Timestamp
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, log.ScanTimestamp.Format("2006-01-02 15:04"))
		x += colWidths[2]

		// Scan Result (with color)
		switch log.ScanResult {
		case domain.ScanResultSuccess:
			pdf.SetTextColor(34, 139, 34)
		case domain.ScanResultInvalidID:
			pdf.SetTextColor(255, 140, 0)
		case domain.ScanResultAssetNotFound:
			pdf.SetTextColor(220, 20, 60)
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, string(log.ScanResult))
		pdf.SetTextColor(0, 0, 0)
		x += colWidths[3]

		// Scanned By ID (just show ID for brevity)
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		scannedByDisplay := log.ScannedByID
		if len(scannedByDisplay) > 12 {
			scannedByDisplay = scannedByDisplay[:12] + "..."
		}
		pdf.Cell(nil, scannedByDisplay)
		x += colWidths[4]

		// Coordinates
		coordText := "-"
		if log.ScanLocationLat != nil && log.ScanLocationLng != nil {
			coordText = fmt.Sprintf("%.4f, %.4f", *log.ScanLocationLat, *log.ScanLocationLng)
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, coordText)

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
	pdf.Cell(nil, fmt.Sprintf("%s: %d", totalLogsText, len(logs)))

	var buf bytes.Buffer
	if err := pdf.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// exportScanLogListToExcel generates Excel file for scan log list
func (s *Service) exportScanLogListToExcel(logs []domain.ScanLogListResponse) ([]byte, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error closing Excel file:", err)
		}
	}()

	sheetName := "Scan Logs"
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
		"Asset ID", "Scanned Value", "Scan Method", "Scanned By",
		"Scan Timestamp", "Scan Result", "Latitude", "Longitude",
	}

	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	for row, log := range logs {
		rowNum := row + 2

		assetID := ""
		if log.AssetID != nil {
			assetID = *log.AssetID
		}

		lat := ""
		if log.ScanLocationLat != nil {
			lat = fmt.Sprintf("%.6f", *log.ScanLocationLat)
		}

		lng := ""
		if log.ScanLocationLng != nil {
			lng = fmt.Sprintf("%.6f", *log.ScanLocationLng)
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), assetID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), log.ScannedValue)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), string(log.ScanMethod))
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), log.ScannedByID)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), log.ScanTimestamp.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), string(log.ScanResult))
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), lat)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), lng)
	}

	for col := 1; col <= len(headers); col++ {
		colName, _ := excelize.ColumnNumberToName(col)
		f.SetColWidth(sheetName, colName, colName, 18)
	}

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
