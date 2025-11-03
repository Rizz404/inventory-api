package user

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

// ExportUserList exports user list to PDF or Excel format
func (s *Service) ExportUserList(ctx context.Context, payload domain.ExportUserListPayload, params domain.UserParams, langCode string) ([]byte, string, error) {
	if payload.SearchQuery != nil {
		params.SearchQuery = payload.SearchQuery
	}
	if payload.Filters != nil {
		params.Filters = payload.Filters
	}
	if payload.Sort != nil {
		params.Sort = payload.Sort
	}

	users, err := s.Repo.GetUsersForExport(ctx, params)
	if err != nil {
		return nil, "", err
	}

	// Convert to responses
	userResponses := mapper.UsersToListResponses(users)

	switch payload.Format {
	case domain.ExportFormatPDF:
		data, err := s.exportUserListToPDF(userResponses, langCode)
		if err != nil {
			return nil, "", domain.ErrInternal(err)
		}
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		filename := fmt.Sprintf("users_%s.pdf", timestamp)
		return data, filename, nil

	case domain.ExportFormatExcel:
		data, err := s.exportUserListToExcel(userResponses)
		if err != nil {
			return nil, "", domain.ErrInternal(err)
		}
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		filename := fmt.Sprintf("users_%s.xlsx", timestamp)
		return data, filename, nil

	default:
		return nil, "", domain.ErrBadRequest("Invalid export format")
	}
}

// exportUserListToPDF generates PDF file for user list using gopdf
func (s *Service) exportUserListToPDF(users []domain.UserListResponse, langCode string) ([]byte, error) {
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
	reportTitle := utils.GetLocalizedMessage(utils.PDFUserReportKey, langCode)
	generatedOnText := utils.GetLocalizedMessage(utils.PDFAssetGeneratedOnKey, langCode)
	totalUsersText := utils.GetLocalizedMessage(utils.PDFUserTotalKey, langCode)
	nameText := utils.GetLocalizedMessage(utils.PDFUserNameKey, langCode)
	emailText := utils.GetLocalizedMessage(utils.PDFUserEmailKey, langCode)
	fullNameText := utils.GetLocalizedMessage(utils.PDFUserFullNameKey, langCode)
	roleText := utils.GetLocalizedMessage(utils.PDFUserRoleKey, langCode)
	employeeIDText := utils.GetLocalizedMessage(utils.PDFUserEmployeeIDKey, langCode)
	isActiveText := utils.GetLocalizedMessage(utils.PDFUserIsActiveKey, langCode)

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
	colWidths := []float64{100, 180, 140, 80, 90, 70}
	headers := []string{nameText, emailText, fullNameText, roleText, employeeIDText, isActiveText}

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

	for i, user := range users {
		maxLines := 1

		emailLines := wrapText(user.Email, colWidths[1])
		if len(emailLines) > maxLines {
			maxLines = len(emailLines)
		}

		fullNameLines := wrapText(user.FullName, colWidths[2])
		if len(fullNameLines) > maxLines {
			maxLines = len(fullNameLines)
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

		// Name
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, user.Name)
		x += colWidths[0]

		// Email (multi-line)
		for lineIdx, line := range emailLines {
			pdf.SetX(x + 3)
			pdf.SetY(cellY + float64(lineIdx)*10)
			pdf.Cell(nil, line)
		}
		x += colWidths[1]

		// Full Name (multi-line)
		for lineIdx, line := range fullNameLines {
			pdf.SetX(x + 3)
			pdf.SetY(cellY + float64(lineIdx)*10)
			pdf.Cell(nil, line)
		}
		x += colWidths[2]

		// Role (with color)
		switch user.Role {
		case domain.RoleAdmin:
			pdf.SetTextColor(220, 20, 60)
		case domain.RoleStaff:
			pdf.SetTextColor(255, 140, 0)
		case domain.RoleEmployee:
			pdf.SetTextColor(34, 139, 34)
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, string(user.Role))
		pdf.SetTextColor(0, 0, 0)
		x += colWidths[3]

		// Employee ID
		employeeID := "-"
		if user.EmployeeID != nil {
			employeeID = *user.EmployeeID
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, employeeID)
		x += colWidths[4]

		// Is Active (with color)
		activeText := "No"
		if user.IsActive {
			activeText = "Yes"
			pdf.SetTextColor(34, 139, 34)
		} else {
			pdf.SetTextColor(220, 20, 60)
		}
		pdf.SetX(x + 3)
		pdf.SetY(cellY)
		pdf.Cell(nil, activeText)
		pdf.SetTextColor(0, 0, 0)

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
	pdf.Cell(nil, fmt.Sprintf("%s: %d", totalUsersText, len(users)))

	var buf bytes.Buffer
	if err := pdf.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// exportUserListToExcel generates Excel file for user list
func (s *Service) exportUserListToExcel(users []domain.UserListResponse) ([]byte, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error closing Excel file:", err)
		}
	}()

	sheetName := "Users"
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
		"Name", "Email", "Full Name", "Role",
		"Employee ID", "Preferred Language", "Is Active", "Created At",
	}

	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	for row, user := range users {
		rowNum := row + 2

		employeeID := ""
		if user.EmployeeID != nil {
			employeeID = *user.EmployeeID
		}

		isActive := "No"
		if user.IsActive {
			isActive = "Yes"
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), user.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), user.Email)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), user.FullName)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), string(user.Role))
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), employeeID)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), user.PreferredLang)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), isActive)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), user.CreatedAt.Format("2006-01-02 15:04:05"))
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
