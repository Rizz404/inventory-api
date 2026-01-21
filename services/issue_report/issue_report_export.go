package issue_report

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

// ExportIssueReportList exports issue report list to PDF or Excel format
func (s *Service) ExportIssueReportList(ctx context.Context, payload domain.ExportIssueReportListPayload, params domain.IssueReportParams, langCode string) ([]byte, string, error) {
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

	// Get reports without pagination
	reports, err := s.Repo.GetIssueReportsForExport(ctx, params, langCode)
	if err != nil {
		return nil, "", err
	}

	// Convert to responses (includes translations)
	reportResponses := mapper.IssueReportsToResponses(reports, langCode)

	switch payload.Format {
	case domain.ExportFormatPDF:
		data, err := s.exportIssueReportListToPDF(reportResponses, langCode)
		if err != nil {
			return nil, "", domain.ErrInternal(err)
		}
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		filename := fmt.Sprintf("issue_reports_%s.pdf", timestamp)
		return data, filename, nil

	case domain.ExportFormatExcel:
		data, err := s.exportIssueReportListToExcel(reportResponses)
		if err != nil {
			return nil, "", domain.ErrInternal(err)
		}
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		filename := fmt.Sprintf("issue_reports_%s.xlsx", timestamp)
		return data, filename, nil

	default:
		return nil, "", domain.ErrBadRequest("Invalid export format")
	}
}

// exportIssueReportListToPDF generates PDF file for issue report list using gopdf
func (s *Service) exportIssueReportListToPDF(reports []domain.IssueReportResponse, langCode string) ([]byte, error) {
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
	reportTitle := utils.GetLocalizedMessage(utils.PDFIssueReportReportKey, langCode)
	generatedOnText := utils.GetLocalizedMessage(utils.PDFAssetGeneratedOnKey, langCode)
	totalReportsText := utils.GetLocalizedMessage(utils.PDFIssueReportTotalKey, langCode)
	assetTagText := utils.GetLocalizedMessage(utils.PDFAssetAssetTagKey, langCode)
	assetNameText := utils.GetLocalizedMessage(utils.PDFAssetAssetNameKey, langCode)
	titleText := utils.GetLocalizedMessage(utils.PDFIssueReportTitleKey, langCode)
	issueTypeText := utils.GetLocalizedMessage(utils.PDFIssueReportTypeKey, langCode)
	priorityText := utils.GetLocalizedMessage(utils.PDFIssueReportPriorityKey, langCode)
	statusText := utils.GetLocalizedMessage(utils.PDFAssetStatusKey, langCode)
	reportedByText := utils.GetLocalizedMessage(utils.PDFIssueReportReportedByKey, langCode)
	reportedDateText := utils.GetLocalizedMessage(utils.PDFIssueReportReportedDateKey, langCode)

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
	colWidths := []float64{70, 100, 120, 80, 70, 70, 100, 80}
	headers := []string{assetTagText, assetNameText, titleText, issueTypeText, priorityText, statusText, reportedByText, reportedDateText}

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

	for i, report := range reports {
		maxLines := 1

		nameLines := wrapText(report.Asset.AssetName, colWidths[1])
		if len(nameLines) > maxLines {
			maxLines = len(nameLines)
		}

		titleLines := wrapText(report.Title, colWidths[2])
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
		pdf.Cell(nil, report.Asset.AssetTag)
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

		// Issue Type
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, report.IssueType)
		x += colWidths[3]

		// Priority (with color)
		switch report.Priority {
		case domain.PriorityLow:
			pdf.SetTextColor(34, 139, 34)
		case domain.PriorityMedium:
			pdf.SetTextColor(255, 215, 0)
		case domain.PriorityHigh:
			pdf.SetTextColor(255, 140, 0)
		case domain.PriorityCritical:
			pdf.SetTextColor(220, 20, 60)
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, string(report.Priority))
		pdf.SetTextColor(0, 0, 0)
		x += colWidths[4]

		// Status (with color)
		switch report.Status {
		case domain.IssueStatusOpen:
			pdf.SetTextColor(255, 140, 0)
		case domain.IssueStatusInProgress:
			pdf.SetTextColor(255, 215, 0)
		case domain.IssueStatusResolved:
			pdf.SetTextColor(34, 139, 34)
		case domain.IssueStatusClosed:
			pdf.SetTextColor(128, 128, 128)
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, string(report.Status))
		pdf.SetTextColor(0, 0, 0)
		x += colWidths[5]

		// Reported By
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, report.ReportedBy.FullName)
		x += colWidths[6]

		// Reported Date
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, report.ReportedDate.Format("2006-01-02"))

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
	pdf.Cell(nil, fmt.Sprintf("%s: %d", totalReportsText, len(reports)))

	var buf bytes.Buffer
	if err := pdf.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// exportIssueReportListToExcel generates Excel file for issue report list
func (s *Service) exportIssueReportListToExcel(reports []domain.IssueReportResponse) ([]byte, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error closing Excel file:", err)
		}
	}()

	sheetName := "Issue Reports"
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
		"Asset Tag", "Asset Name", "Title", "Issue Type",
		"Priority", "Status", "Reported By", "Reported Date",
		"Resolved By", "Resolved Date", "Description", "Resolution Notes",
	}

	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	for row, report := range reports {
		rowNum := row + 2

		resolvedBy := ""
		if report.ResolvedBy != nil {
			resolvedBy = report.ResolvedBy.FullName
		}

		resolvedDate := ""
		if report.ResolvedDate != nil {
			resolvedDate = report.ResolvedDate.Format("2006-01-02")
		}

		description := ""
		if report.Description != nil {
			description = *report.Description
		}

		resolutionNotes := ""
		if report.ResolutionNotes != nil {
			resolutionNotes = *report.ResolutionNotes
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), report.Asset.AssetTag)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), report.Asset.AssetName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), report.Title)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), report.IssueType)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), string(report.Priority))
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), string(report.Status))
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), report.ReportedBy.FullName)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), report.ReportedDate.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", rowNum), resolvedBy)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", rowNum), resolvedDate)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", rowNum), description)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", rowNum), resolutionNotes)
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
