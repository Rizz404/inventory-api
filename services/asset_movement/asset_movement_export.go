package asset_movement

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

// ExportAssetMovementList exports asset movement list to PDF or Excel format
func (s *Service) ExportAssetMovementList(ctx context.Context, payload domain.ExportAssetMovementListPayload, params domain.AssetMovementParams, langCode string) ([]byte, string, error) {
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

	// Get movements without pagination
	movements, err := s.Repo.GetAssetMovementsForExport(ctx, params, langCode)
	if err != nil {
		return nil, "", err
	}

	// Convert to responses (includes translations)
	movementResponses := mapper.AssetMovementsToResponses(movements, langCode)

	switch payload.Format {
	case domain.ExportFormatPDF:
		data, err := s.exportAssetMovementListToPDF(movementResponses, langCode)
		if err != nil {
			return nil, "", domain.ErrInternal(err)
		}
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		filename := fmt.Sprintf("asset_movements_%s.pdf", timestamp)
		return data, filename, nil

	case domain.ExportFormatExcel:
		data, err := s.exportAssetMovementListToExcel(movementResponses)
		if err != nil {
			return nil, "", domain.ErrInternal(err)
		}
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		filename := fmt.Sprintf("asset_movements_%s.xlsx", timestamp)
		return data, filename, nil

	default:
		return nil, "", domain.ErrBadRequest("Invalid export format")
	}
}

// exportAssetMovementListToPDF generates PDF file for asset movement list using gopdf
func (s *Service) exportAssetMovementListToPDF(movements []domain.AssetMovementResponse, langCode string) ([]byte, error) {
	// Get absolute path for fonts and logo
	workDir, _ := os.Getwd()
	fontRegularPath := filepath.Join(workDir, "assets", "fonts", "NotoSansJP-Regular.ttf")
	fontBoldPath := filepath.Join(workDir, "assets", "fonts", "NotoSansJP-Bold.ttf")
	logoPath := filepath.Join(workDir, "assets", "images", "fts-logo.png")

	// Initialize gopdf
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		PageSize: *gopdf.PageSizeA4Landscape,
		Unit:     gopdf.Unit_PT,
	})
	pdf.AddPage()

	// Load fonts
	if err := pdf.AddTTFFont("noto-regular", fontRegularPath); err != nil {
		return nil, fmt.Errorf("failed to load regular font: %w", err)
	}
	if err := pdf.AddTTFFont("noto-bold", fontBoldPath); err != nil {
		return nil, fmt.Errorf("failed to load bold font: %w", err)
	}

	pdf.SetFont("noto-regular", "", 10)

	// Get localized text
	reportTitle := utils.GetLocalizedMessage(utils.PDFAssetMovementReportKey, langCode)
	generatedOnText := utils.GetLocalizedMessage(utils.PDFAssetGeneratedOnKey, langCode)
	totalMovementsText := utils.GetLocalizedMessage(utils.PDFAssetMovementTotalKey, langCode)
	assetTagText := utils.GetLocalizedMessage(utils.PDFAssetAssetTagKey, langCode)
	assetNameText := utils.GetLocalizedMessage(utils.PDFAssetAssetNameKey, langCode)
	fromLocationText := utils.GetLocalizedMessage(utils.PDFAssetMovementFromLocationKey, langCode)
	toLocationText := utils.GetLocalizedMessage(utils.PDFAssetMovementToLocationKey, langCode)
	fromUserText := utils.GetLocalizedMessage(utils.PDFAssetMovementFromUserKey, langCode)
	toUserText := utils.GetLocalizedMessage(utils.PDFAssetMovementToUserKey, langCode)
	movedByText := utils.GetLocalizedMessage(utils.PDFAssetMovementMovedByKey, langCode)
	movementDateText := utils.GetLocalizedMessage(utils.PDFAssetMovementDateKey, langCode)
	notesText := utils.GetLocalizedMessage(utils.PDFAssetMovementNotesKey, langCode)

	// Page setup
	marginLeft := 30.0
	marginTop := 50.0
	pageWidth := 842.0
	pageHeight := 595.0
	contentWidth := pageWidth - (marginLeft * 2)

	// Add company logo if exists
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
	colWidths := []float64{70, 110, 90, 90, 90, 90, 90, 90, 100} // Total: 820
	headers := []string{assetTagText, assetNameText, fromLocationText, toLocationText, fromUserText, toUserText, movedByText, movementDateText, notesText}

	// Draw table header
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

	// Reset for data rows
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("noto-regular", "", 8)
	y += 25

	// Helper function to wrap text
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

	for i, movement := range movements {
		// Calculate row height
		maxLines := 1

		nameLines := wrapText(movement.Asset.AssetName, colWidths[1])
		if len(nameLines) > maxLines {
			maxLines = len(nameLines)
		}

		if movement.Notes != nil && *movement.Notes != "" {
			notesLines := wrapText(*movement.Notes, colWidths[8])
			if len(notesLines) > maxLines {
				maxLines = len(notesLines)
			}
		}

		rowHeight := float64(maxLines) * 12.0
		if rowHeight < 18 {
			rowHeight = 18
		}

		// Check if need new page
		if y+rowHeight > pageHeight-40 {
			pdf.AddPage()
			y = marginTop

			// Redraw header
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

		// Zebra striping
		if i%2 == 1 {
			pdf.SetFillColor(242, 242, 242)
			pdf.RectFromUpperLeftWithStyle(marginLeft, y, contentWidth, rowHeight, "F")
		}

		x = marginLeft
		cellY := y + 5

		// Asset Tag
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, movement.Asset.AssetTag)
		x += colWidths[0]

		// Asset Name (multi-line)
		for lineIdx, line := range nameLines {
			pdf.SetX(x + 3)
			pdf.SetY(cellY + float64(lineIdx)*10)
			pdf.Cell(nil, line)
		}
		x += colWidths[1]

		// From Location
		fromLoc := "-"
		if movement.FromLocation != nil {
			fromLoc = movement.FromLocation.LocationName
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, fromLoc)
		x += colWidths[2]

		// To Location
		toLoc := "-"
		if movement.ToLocation != nil {
			toLoc = movement.ToLocation.LocationName
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, toLoc)
		x += colWidths[3]

		// From User
		fromUser := "-"
		if movement.FromUser != nil {
			fromUser = movement.FromUser.FullName
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, fromUser)
		x += colWidths[4]

		// To User
		toUser := "-"
		if movement.ToUser != nil {
			toUser = movement.ToUser.FullName
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, toUser)
		x += colWidths[5]

		// Moved By
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, movement.MovedBy.FullName)
		x += colWidths[6]

		// Movement Date
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, movement.MovementDate.Format("2006-01-02"))
		x += colWidths[7]

		// Notes (multi-line)
		if movement.Notes != nil && *movement.Notes != "" {
			notesLines := wrapText(*movement.Notes, colWidths[8])
			for lineIdx, line := range notesLines {
				pdf.SetX(x + 3)
				pdf.SetY(cellY + float64(lineIdx)*10)
				pdf.Cell(nil, line)
			}
		} else {
			pdf.SetX(x + 3)
			pdf.SetY(cellY)
			pdf.Cell(nil, "-")
		}

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
	pdf.Cell(nil, fmt.Sprintf("%s: %d", totalMovementsText, len(movements)))

	// Get PDF bytes
	var buf bytes.Buffer
	if err := pdf.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// exportAssetMovementListToExcel generates Excel file for asset movement list
func (s *Service) exportAssetMovementListToExcel(movements []domain.AssetMovementResponse) ([]byte, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error closing Excel file:", err)
		}
	}()

	sheetName := "Asset Movements"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

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
		"Asset Tag", "Asset Name", "From Location", "To Location",
		"From User", "To User", "Moved By", "Movement Date", "Notes",
	}

	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Add data
	for row, movement := range movements {
		rowNum := row + 2

		fromLoc := ""
		if movement.FromLocation != nil {
			fromLoc = movement.FromLocation.LocationName
		}

		toLoc := ""
		if movement.ToLocation != nil {
			toLoc = movement.ToLocation.LocationName
		}

		fromUser := ""
		if movement.FromUser != nil {
			fromUser = movement.FromUser.FullName
		}

		toUser := ""
		if movement.ToUser != nil {
			toUser = movement.ToUser.FullName
		}

		notes := ""
		if movement.Notes != nil {
			notes = *movement.Notes
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), movement.Asset.AssetTag)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), movement.Asset.AssetName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), fromLoc)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), toLoc)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), fromUser)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), toUser)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), movement.MovedBy.FullName)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), movement.MovementDate.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", rowNum), notes)
	}

	// Auto-fit columns
	for col := 1; col <= len(headers); col++ {
		colName, _ := excelize.ColumnNumberToName(col)
		f.SetColWidth(sheetName, colName, colName, 18)
	}

	// Save to buffer
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
