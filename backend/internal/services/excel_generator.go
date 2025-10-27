package services

import (
	"fmt"

	"perjalanan-dinas/backend/internal/models"

	"github.com/xuri/excelize/v2"
)

type ExcelGenerator struct{}

type EmployeeAllowanceRow struct {
	NIP                  string
	Name                 string
	Position             string
	TotalTrips           int
	DaysInProvince       int
	DaysOutsideProvince  int
	DaysAbroad           int
	TotalAllowance       float64
}

func NewExcelGenerator() *ExcelGenerator {
	return &ExcelGenerator{}
}

func (eg *ExcelGenerator) GenerateMonthlyAllowanceReport(requests []models.TravelRequest, year int, month int) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Rekap Iuran"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to create sheet: %w", err)
	}
	f.SetActiveSheet(index)

	// Delete default sheet
	f.DeleteSheet("Sheet1")

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 15)  // NIP
	f.SetColWidth(sheetName, "B", "B", 30)  // NAMA
	f.SetColWidth(sheetName, "C", "C", 35)  // JABATAN
	f.SetColWidth(sheetName, "D", "D", 12)  // JUMLAH TRIP
	f.SetColWidth(sheetName, "E", "G", 18)  // Days columns
	f.SetColWidth(sheetName, "H", "H", 20)  // TOTAL IURAN

	// Create header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Bold:  true,
			Size:  11,
			Color: "#FFFFFF",
		},
	})

	// Create currency style
	currencyStyle, _ := f.NewStyle(&excelize.Style{
		NumFmt: 177, // Indonesian Rupiah format
	})

	// Create center alignment style
	centerStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	// Aggregate data by employee
	employeeMap := make(map[uint]*EmployeeAllowanceRow)
	hasInProvince := false
	hasOutsideProvince := false
	hasAbroad := false

	for _, request := range requests {
		// Filter by month and year
		departureDate := request.DepartureDate
		if departureDate.Year() != year || int(departureDate.Month()) != month {
			continue
		}

		for _, empRel := range request.TravelRequestEmployees {
			emp := empRel.Employee
			if _, exists := employeeMap[emp.ID]; !exists {
				employeeMap[emp.ID] = &EmployeeAllowanceRow{
					NIP:      emp.NIP,
					Name:     emp.Name,
					Position: emp.Position.Title,
				}
			}

			row := employeeMap[emp.ID]
			row.TotalTrips++
			row.TotalAllowance += float64(request.TotalAllowance) / float64(len(request.TravelRequestEmployees))

			// Count days by destination type
			switch request.DestinationType {
			case "in_province":
				row.DaysInProvince += request.DurationDays
				hasInProvince = true
			case "outside_province":
				row.DaysOutsideProvince += request.DurationDays
				hasOutsideProvince = true
			case "abroad":
				row.DaysAbroad += request.DurationDays
				hasAbroad = true
			}
		}
	}

	// Build dynamic headers based on data
	headers := []string{"NIP", "NAMA", "JABATAN", "JUMLAH TRIP"}
	col := 'E'
	daysCols := make(map[string]rune)

	if hasInProvince {
		headers = append(headers, "HARI DALAM PROVINSI")
		daysCols["in_province"] = col
		col++
	}
	if hasOutsideProvince {
		headers = append(headers, "HARI LUAR PROVINSI")
		daysCols["outside_province"] = col
		col++
	}
	if hasAbroad {
		headers = append(headers, "HARI LUAR NEGERI")
		daysCols["abroad"] = col
		col++
	}

	headers = append(headers, "TOTAL IURAN")
	totalCol := col

	// Write headers
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Write data
	row := 2
	for _, empRow := range employeeMap {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), empRow.NIP)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), empRow.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), empRow.Position)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), empRow.TotalTrips)
		f.SetCellStyle(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), centerStyle)

		// Write days based on what's available
		if hasInProvince {
			cell := fmt.Sprintf("%c%d", daysCols["in_province"], row)
			f.SetCellValue(sheetName, cell, empRow.DaysInProvince)
			f.SetCellStyle(sheetName, cell, cell, centerStyle)
		}
		if hasOutsideProvince {
			cell := fmt.Sprintf("%c%d", daysCols["outside_province"], row)
			f.SetCellValue(sheetName, cell, empRow.DaysOutsideProvince)
			f.SetCellStyle(sheetName, cell, cell, centerStyle)
		}
		if hasAbroad {
			cell := fmt.Sprintf("%c%d", daysCols["abroad"], row)
			f.SetCellValue(sheetName, cell, empRow.DaysAbroad)
			f.SetCellStyle(sheetName, cell, cell, centerStyle)
		}

		// Write total allowance
		totalCell := fmt.Sprintf("%c%d", totalCol, row)
		f.SetCellValue(sheetName, totalCell, empRow.TotalAllowance)
		f.SetCellStyle(sheetName, totalCell, totalCell, currencyStyle)

		row++
	}

	// Generate file
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write excel file: %w", err)
	}

	return buffer.Bytes(), nil
}

func getMonthName(month int) string {
	months := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	if month < 1 || month > 12 {
		return ""
	}
	return months[month-1]
}
