package services

import (
	"bytes"
	"fmt"
	"perjalanan-dinas/backend/internal/models"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type PDFGenerator struct{}

func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{}
}

// GenerateNotaPermintaan generates the "Nota Permintaan Surat Tugas Perjalanan Dinas" PDF
func (pg *PDFGenerator) GenerateNotaPermintaan(request *models.TravelRequest) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Add logos at the top
	// BPD logo on the left
	pdf.Image("assets/images/bpd.png", 15, 10, 20, 0, false, "", 0, "")
	// Bank Jatim logo on the right
	pdf.Image("assets/images/bank jatim.png", 155, 10, 40, 0, false, "", 0, "")

	// Title
	pdf.SetY(35)
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 7, "NOTA PERMINTAAN", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 7, "SURAT TUGAS PERJALANAN DINAS", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Document details
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(30, 6, "Nomor", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.RequestNumber, "", 1, "L", false, 0, "")

	pdf.CellFormat(30, 6, "Kepada", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, "Divisi Human Capital", "", 1, "L", false, 0, "")

	pdf.CellFormat(30, 6, "Dari", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, "Divisi Digital Banking", "", 1, "L", false, 0, "")

	pdf.CellFormat(30, 6, "Perihal", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, "Permohonan Surat Perjalanan Dinas", "", 1, "L", false, 0, "")
	pdf.Ln(5)

	// Body text
	pdf.MultiCell(0, 6, "Dengan ini kami mohon bantuan Saudara untuk membuatkan Surat Tugas Perjalanan Dinas sehubungan dengan penugasan kepada:", "", "L", false)
	pdf.Ln(2)

	// Employee table
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(10, 8, "NO", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "NIP", "1", 0, "C", true, 0, "")
	pdf.CellFormat(60, 8, "NAMA", "1", 0, "C", true, 0, "")
	pdf.CellFormat(80, 8, "JABATAN", "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 9)
	empColWidths := []float64{10, 30, 60, 80}
	for i, empRel := range request.TravelRequestEmployees {
		drawEmployeeRow(pdf, i+1, empRel.Employee.NIP, empRel.Employee.Name, empRel.Employee.Position.Title, empColWidths, 9)
	}
	pdf.Ln(5)

	// Travel details
	pdf.SetFont("Arial", "", 11)

	// Maksud perjalanan dinas
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(5, 6, "")
	pdf.Cell(5, 6, "-")
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(60, 6, "Maksud perjalanan dinas", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.MultiCell(0, 6, request.Purpose, "", "L", false)

	// Tempat berangkat dan tujuan
	pdf.Cell(5, 6, "")
	pdf.Cell(5, 6, "-")
	pdf.CellFormat(60, 6, "Tempat berangkat", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.DeparturePlace, "", 1, "L", false, 0, "")

	pdf.Cell(10, 6, "")
	pdf.CellFormat(60, 6, "Tempat tujuan", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.Destination, "", 1, "L", false, 0, "")

	// Lama perjalanan dinas
	pdf.Cell(5, 6, "")
	pdf.Cell(5, 6, "-")
	pdf.CellFormat(60, 6, "Lama perjalanan dinas", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf("%d (%s) hari", request.DurationDays, numberToWords(request.DurationDays)), "", 1, "L", false, 0, "")

	pdf.Cell(10, 6, "")
	pdf.CellFormat(60, 6, "Tanggal berangkat", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.DepartureDate.Format("02 January 2006"), "", 1, "L", false, 0, "")

	pdf.Cell(10, 6, "")
	pdf.CellFormat(60, 6, "Tanggal kembali", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.ReturnDate.Format("02 January 2006"), "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// Angkutan yang digunakan
	pdf.Cell(5, 6, "")
	pdf.Cell(5, 6, "-")
	pdf.CellFormat(60, 6, "Angkutan yang digunakan", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.Transportation, "", 1, "L", false, 0, "")
	pdf.Ln(8)

	// Closing
	pdf.MultiCell(0, 6, "Atas perhatian dan kerjasamanya disampaikan terima kasih.", "", "L", false)
	pdf.Ln(10)

	// Signature section - right aligned with consistent positioning
	x2 := 110.0
	pdf.SetX(x2)
	pdf.CellFormat(90, 6, fmt.Sprintf("Surabaya, %s", time.Now().Format("02 January 2006")), "", 1, "C", false, 0, "")
	pdf.SetX(x2)
	pdf.CellFormat(90, 6, "DIVISI DIGITAL BANKING", "", 1, "C", false, 0, "")
	pdf.Ln(20)

	// Use representative from TravelReport, fallback to defaults if not available
	repName := "M. MACHFUD HIDAYAT"
	repPosition := "Vice President"
	if request.TravelReport != nil {
		repName = request.TravelReport.RepresentativeName
		repPosition = request.TravelReport.RepresentativePosition
	}

	pdf.SetFont("Arial", "BU", 11)
	pdf.SetX(x2)
	pdf.CellFormat(90, 6, repName, "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.SetX(x2)
	pdf.CellFormat(90, 6, repPosition, "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Tindasan
	pdf.SetFont("Arial", "", 8)
	pdf.CellFormat(0, 6, "Tindasan :", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 6, "- Arsip", "", 1, "L", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GenerateBeritaAcara generates the "Berita Acara Perjalanan Dinas" PDF
func (pg *PDFGenerator) GenerateBeritaAcara(request *models.TravelRequest, report *models.TravelReport) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Add logos at the top
	// BPD logo on the left
	pdf.Image("assets/images/bpd.png", 15, 10, 30, 0, false, "", 0, "")
	// Bank Jatim logo on the right
	pdf.Image("assets/images/bank jatim.png", 165, 10, 30, 0, false, "", 0, "")

	// Title
	pdf.SetY(35)
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 7, "BERITA ACARA PERJALANAN DINAS", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 6, fmt.Sprintf("No: %s", report.ReportNumber), "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Opening text
	pdf.SetFont("Arial", "", 11)
	pdf.MultiCell(0, 6, "PT. Bank Pembangunan Daerah Jawa Timur, dengan ini menugaskan kepada:", "", "L", false)
	pdf.Ln(2)

	// Employee table
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(12, 8, "No.", "1", 0, "C", true, 0, "")
	pdf.CellFormat(28, 8, "NIP", "1", 0, "C", true, 0, "")
	pdf.CellFormat(50, 8, "Nama", "1", 0, "C", true, 0, "")
	pdf.CellFormat(90, 8, "Jabatan", "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 9)
	empColWidths := []float64{12, 28, 50, 90}
	for i, empRel := range request.TravelRequestEmployees {
		drawEmployeeRow(pdf, i+1, empRel.Employee.NIP, empRel.Employee.Name, empRel.Employee.Position.Title, empColWidths, 9)
	}
	pdf.Ln(3)

	// Travel details with numbering
	pdf.SetFont("Arial", "", 11)

	// 1. Maksud perjalanan dinas
	pdf.CellFormat(10, 6, "1.", "", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, "Maksud perjalanan dinas", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.MultiCell(0, 6, request.Purpose, "", "L", false)

	// 2. Tempat berangkat dan tujuan
	pdf.CellFormat(10, 6, "2.", "", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, "a. Tempat berangkat", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.DeparturePlace, "", 1, "L", false, 0, "")

	pdf.Cell(10, 6, "")
	pdf.CellFormat(60, 6, "b. Tempat tujuan", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.Destination, "", 1, "L", false, 0, "")

	// 3. Lama perjalanan dinas
	pdf.CellFormat(10, 6, "3.", "", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, "Lama perjalanan dinas", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf("%d (%s) hari", request.DurationDays, numberToWords(request.DurationDays)), "", 1, "L", false, 0, "")

	pdf.Cell(10, 6, "")
	pdf.CellFormat(60, 6, "Tanggal berangkat", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.DepartureDate.Format("02 January 2006"), "", 1, "L", false, 0, "")

	pdf.Cell(10, 6, "")
	pdf.CellFormat(60, 6, "Tanggal kembali", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.ReturnDate.Format("02 January 2006"), "", 1, "L", false, 0, "")

	// 4. Angkutan yang digunakan
	pdf.CellFormat(10, 6, "4.", "", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, "Angkutan yang digunakan", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.Transportation, "", 1, "L", false, 0, "")
	pdf.Ln(5)

	// Closing
	pdf.MultiCell(0, 6, "Atas perhatian dan kerjasamanya disampaikan terimakasih", "", "L", false)
	pdf.Ln(10)

	// Signature section - two columns
	pdf.SetFont("Arial", "", 10)

	// Left column - Only first employee with signature space
	x1 := 15.0
	y1 := pdf.GetY()
	pdf.SetXY(x1, y1)
	pdf.CellFormat(90, 6, "Pegawai Yang Bersangkutan", "", 1, "C", false, 0, "")
	pdf.SetX(x1)
	pdf.Ln(22) // Space for signature
	pdf.SetX(x1)

	// Only show first employee
	if len(request.TravelRequestEmployees) > 0 {
		empRel := request.TravelRequestEmployees[0]
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(90, 5, fmt.Sprintf("(%s)", empRel.Employee.Name), "", 1, "C", false, 0, "")
		pdf.SetX(x1)
		pdf.SetFont("Arial", "", 9)
		pdf.CellFormat(90, 4, empRel.Employee.Position.Title, "", 1, "C", false, 0, "")
	}

	// Right column - Representative
	x2 := 110.0
	pdf.SetXY(x2, y1)
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(90, 6, fmt.Sprintf("Surabaya, %s", time.Now().Format("02 January 2006")), "", 1, "C", false, 0, "")
	pdf.SetX(x2)
	pdf.CellFormat(90, 6, "DIVISI DIGITAL BANKING", "", 1, "C", false, 0, "")
	pdf.SetX(x2)
	pdf.Ln(15)
	pdf.SetX(x2)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(90, 6, report.RepresentativeName, "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.SetX(x2)
	pdf.CellFormat(90, 6, report.RepresentativePosition, "", 1, "C", false, 0, "")

	pdf.Ln(10)

	// Visit proof table on same page
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(240, 240, 240)

	colWidths := []float64{25, 30, 35, 30, 60}
	headers := []string{"Tanggal", "Berangkat\ndari", "Bermalam/\nsinggah di", "Datang di", "Tanda-tangan\nKP/Cabang/Lembaga"}

	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Visit proof rows
	pdf.SetFont("Arial", "", 8)
	if len(report.VisitProofs) > 0 {
		for _, proof := range report.VisitProofs {
			pdf.CellFormat(colWidths[0], 8, proof.Date.Format("02/01/2006"), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidths[1], 8, proof.DepartFrom, "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidths[2], 8, proof.StayOrStopAt, "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidths[3], 8, proof.ArriveAt, "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidths[4], 8, "", "1", 1, "C", false, 0, "")
		}
	} else {
		// Empty rows for manual filling
		for i := 0; i < 5; i++ {
			for _, width := range colWidths {
				pdf.CellFormat(width, 8, "", "1", 0, "C", false, 0, "")
			}
			pdf.Ln(-1)
		}
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// addNotaPermintaanPage adds Nota Permintaan page to existing PDF
func (pg *PDFGenerator) addNotaPermintaanPage(pdf *gofpdf.Fpdf, request *models.TravelRequest) error {
	pdf.AddPage()

	// Add logos at the top
	// BPD logo on the left
	pdf.Image("assets/images/bpd.png", 15, 10, 20, 0, false, "", 0, "")
	// Bank Jatim logo on the right
	pdf.Image("assets/images/bank jatim.png", 155, 10, 40, 0, false, "", 0, "")

	// Title
	pdf.SetY(35)
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 7, "NOTA PERMINTAAN", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 7, "SURAT TUGAS PERJALANAN DINAS", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Document details
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(30, 6, "Nomor", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.RequestNumber, "", 1, "L", false, 0, "")

	pdf.CellFormat(30, 6, "Kepada", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, "Divisi Human Capital", "", 1, "L", false, 0, "")

	pdf.CellFormat(30, 6, "Dari", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, "Divisi Digital Banking", "", 1, "L", false, 0, "")

	pdf.CellFormat(30, 6, "Perihal", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, "Permohonan Surat Perjalanan Dinas", "", 1, "L", false, 0, "")
	pdf.Ln(5)

	// Body text
	pdf.MultiCell(0, 6, "Dengan ini kami mohon bantuan Saudara untuk membuatkan Surat Tugas Perjalanan Dinas sehubungan dengan penugasan kepada:", "", "L", false)
	pdf.Ln(2)

	// Employee table
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(10, 8, "NO", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "NIP", "1", 0, "C", true, 0, "")
	pdf.CellFormat(60, 8, "NAMA", "1", 0, "C", true, 0, "")
	pdf.CellFormat(80, 8, "JABATAN", "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 9)
	empColWidths := []float64{10, 30, 60, 80}
	for i, empRel := range request.TravelRequestEmployees {
		drawEmployeeRow(pdf, i+1, empRel.Employee.NIP, empRel.Employee.Name, empRel.Employee.Position.Title, empColWidths, 9)
	}
	pdf.Ln(5)

	// Travel details
	pdf.SetFont("Arial", "", 11)

	// Maksud perjalanan dinas
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(5, 6, "")
	pdf.Cell(5, 6, "-")
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(60, 6, "Maksud perjalanan dinas", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.MultiCell(0, 6, request.Purpose, "", "L", false)

	// Tempat berangkat dan tujuan
	pdf.Cell(5, 6, "")
	pdf.Cell(5, 6, "-")
	pdf.CellFormat(60, 6, "Tempat berangkat", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.DeparturePlace, "", 1, "L", false, 0, "")

	pdf.Cell(10, 6, "")
	pdf.CellFormat(60, 6, "Tempat tujuan", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.Destination, "", 1, "L", false, 0, "")

	// Lama perjalanan dinas
	pdf.Cell(5, 6, "")
	pdf.Cell(5, 6, "-")
	pdf.CellFormat(60, 6, "Lama perjalanan dinas", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf("%d (%s) hari", request.DurationDays, numberToWords(request.DurationDays)), "", 1, "L", false, 0, "")

	pdf.Cell(10, 6, "")
	pdf.CellFormat(60, 6, "Tanggal berangkat", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.DepartureDate.Format("02 January 2006"), "", 1, "L", false, 0, "")

	pdf.Cell(10, 6, "")
	pdf.CellFormat(60, 6, "Tanggal kembali", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.ReturnDate.Format("02 January 2006"), "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// Angkutan yang digunakan
	pdf.Cell(5, 6, "")
	pdf.Cell(5, 6, "-")
	pdf.CellFormat(60, 6, "Angkutan yang digunakan", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.Transportation, "", 1, "L", false, 0, "")
	pdf.Ln(8)

	// Closing
	pdf.MultiCell(0, 6, "Atas perhatian dan kerjasamanya disampaikan terima kasih.", "", "L", false)
	pdf.Ln(10)

	// Signature section - right aligned with consistent positioning
	x2 := 110.0
	pdf.SetX(x2)
	pdf.CellFormat(90, 6, fmt.Sprintf("Surabaya, %s", time.Now().Format("02 January 2006")), "", 1, "C", false, 0, "")
	pdf.SetX(x2)
	pdf.CellFormat(90, 6, "DIVISI DIGITAL BANKING", "", 1, "C", false, 0, "")
	pdf.Ln(20)

	// Use representative from TravelReport, fallback to defaults if not available
	repName := "M. MACHFUD HIDAYAT"
	repPosition := "Vice President"
	if request.TravelReport != nil {
		repName = request.TravelReport.RepresentativeName
		repPosition = request.TravelReport.RepresentativePosition
	}

	pdf.SetFont("Arial", "BU", 11)
	pdf.SetX(x2)
	pdf.CellFormat(90, 6, repName, "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.SetX(x2)
	pdf.CellFormat(90, 6, repPosition, "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Tindasan
	pdf.SetFont("Arial", "", 8)
	pdf.CellFormat(0, 6, "Tindasan :", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 6, "- Arsip", "", 1, "L", false, 0, "")

	return nil
}

// addBeritaAcaraPage adds Berita Acara page to existing PDF
func (pg *PDFGenerator) addBeritaAcaraPage(pdf *gofpdf.Fpdf, request *models.TravelRequest, report *models.TravelReport) error {
	pdf.AddPage()

	// Add logos at the top
	// BPD logo on the left
	pdf.Image("assets/images/bpd.png", 15, 10, 30, 0, false, "", 0, "")
	// Bank Jatim logo on the right
	pdf.Image("assets/images/bank jatim.png", 165, 10, 30, 0, false, "", 0, "")

	// Title
	pdf.SetY(35)
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 7, "BERITA ACARA PERJALANAN DINAS", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 6, fmt.Sprintf("No: %s", report.ReportNumber), "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Opening text
	pdf.SetFont("Arial", "", 11)
	pdf.MultiCell(0, 6, "PT. Bank Pembangunan Daerah Jawa Timur, dengan ini menugaskan kepada:", "", "L", false)
	pdf.Ln(2)

	// Employee table
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(12, 8, "No.", "1", 0, "C", true, 0, "")
	pdf.CellFormat(28, 8, "NIP", "1", 0, "C", true, 0, "")
	pdf.CellFormat(50, 8, "Nama", "1", 0, "C", true, 0, "")
	pdf.CellFormat(90, 8, "Jabatan", "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 9)
	empColWidths := []float64{12, 28, 50, 90}
	for i, empRel := range request.TravelRequestEmployees {
		drawEmployeeRow(pdf, i+1, empRel.Employee.NIP, empRel.Employee.Name, empRel.Employee.Position.Title, empColWidths, 9)
	}
	pdf.Ln(3)

	// Travel details with numbering
	pdf.SetFont("Arial", "", 11)

	// 1. Maksud perjalanan dinas
	pdf.CellFormat(10, 6, "1.", "", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, "Maksud perjalanan dinas", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.MultiCell(0, 6, request.Purpose, "", "L", false)

	// 2. Tempat berangkat dan tujuan
	pdf.CellFormat(10, 6, "2.", "", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, "a. Tempat berangkat", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.DeparturePlace, "", 1, "L", false, 0, "")

	pdf.Cell(10, 6, "")
	pdf.CellFormat(60, 6, "b. Tempat tujuan", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.Destination, "", 1, "L", false, 0, "")

	// 3. Lama perjalanan dinas
	pdf.CellFormat(10, 6, "3.", "", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, "Lama perjalanan dinas", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf("%d (%s) hari", request.DurationDays, numberToWords(request.DurationDays)), "", 1, "L", false, 0, "")

	pdf.Cell(10, 6, "")
	pdf.CellFormat(60, 6, "Tanggal berangkat", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.DepartureDate.Format("02 January 2006"), "", 1, "L", false, 0, "")

	pdf.Cell(10, 6, "")
	pdf.CellFormat(60, 6, "Tanggal kembali", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.ReturnDate.Format("02 January 2006"), "", 1, "L", false, 0, "")

	// 4. Angkutan yang digunakan
	pdf.CellFormat(10, 6, "4.", "", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, "Angkutan yang digunakan", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, request.Transportation, "", 1, "L", false, 0, "")
	pdf.Ln(5)

	// Closing
	pdf.MultiCell(0, 6, "Atas perhatian dan kerjasamanya disampaikan terimakasih", "", "L", false)
	pdf.Ln(10)

	// Signature section - two columns
	pdf.SetFont("Arial", "", 10)

	// Left column - Only first employee with signature space
	x1 := 15.0
	y1 := pdf.GetY()
	pdf.SetXY(x1, y1)
	pdf.CellFormat(90, 6, "Pegawai Yang Bersangkutan", "", 1, "C", false, 0, "")
	pdf.SetX(x1)
	pdf.Ln(22) // Space for signature
	pdf.SetX(x1)

	// Only show first employee
	if len(request.TravelRequestEmployees) > 0 {
		empRel := request.TravelRequestEmployees[0]
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(90, 5, fmt.Sprintf("(%s)", empRel.Employee.Name), "", 1, "C", false, 0, "")
		pdf.SetX(x1)
		pdf.SetFont("Arial", "", 9)
		pdf.CellFormat(90, 4, empRel.Employee.Position.Title, "", 1, "C", false, 0, "")
	}

	// Right column - Representative
	x2 := 110.0
	pdf.SetXY(x2, y1)
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(90, 6, fmt.Sprintf("Surabaya, %s", time.Now().Format("02 January 2006")), "", 1, "C", false, 0, "")
	pdf.SetX(x2)
	pdf.CellFormat(90, 6, "PT. BANK PEMBANGUNAN DAERAH JAWA TIMUR", "", 1, "C", false, 0, "")
	pdf.SetX(x2)
	pdf.Ln(15)
	pdf.SetX(x2)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(90, 6, "(________________________)", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	

	pdf.Ln(10)

	// Visit proof table on same page
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(240, 240, 240)

	colWidths := []float64{25, 30, 35, 30, 60}
	headers := []string{"Tanggal", "Berangkat\ndari", "Bermalam/\nsinggah di", "Datang di", "Tanda-tangan\nKP/Cabang/Lembaga"}

	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Visit proof rows
	pdf.SetFont("Arial", "", 8)
	if len(report.VisitProofs) > 0 {
		for _, proof := range report.VisitProofs {
			pdf.CellFormat(colWidths[0], 8, proof.Date.Format("02/01/2006"), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidths[1], 8, proof.DepartFrom, "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidths[2], 8, proof.StayOrStopAt, "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidths[3], 8, proof.ArriveAt, "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidths[4], 8, "", "1", 1, "C", false, 0, "")
		}
	} else {
		// Empty rows for manual filling
		for i := 0; i < 5; i++ {
			for _, width := range colWidths {
				pdf.CellFormat(width, 8, "", "1", 0, "C", false, 0, "")
			}
			pdf.Ln(-1)
		}
	}

	return nil
}

// GenerateCombinedPDF generates both documents in a single PDF
func (pg *PDFGenerator) GenerateCombinedPDF(request *models.TravelRequest, report *models.TravelReport) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 15)

	// Page 1: Nota Permintaan
	err := pg.addNotaPermintaanPage(pdf, request)
	if err != nil {
		return nil, err
	}

	// Page 2: Berita Acara
	err = pg.addBeritaAcaraPage(pdf, request, report)
	if err != nil {
		return nil, err
	}

	// Output to buffer
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Helper function to draw employee table row with text wrapping for position
func drawEmployeeRow(pdf *gofpdf.Fpdf, no int, nip, name, position string, colWidths []float64, fontSize float64) {
	currentY := pdf.GetY()
	x := 15.0 // Left margin

	// First pass: calculate text height for position using SplitLines
	pdf.SetFont("Arial", "", fontSize)
	lines := pdf.SplitLines([]byte(position), colWidths[3]-2)
	numLines := len(lines)
	lineHeight := 3.0
	textHeight := float64(numLines) * lineHeight

	// Determine row height (minimum 8mm)
	rowHeight := 8.0
	if textHeight+2 > rowHeight {
		rowHeight = textHeight + 2
	}

	// Draw NO cell with border
	pdf.SetXY(x, currentY)
	pdf.CellFormat(colWidths[0], rowHeight, fmt.Sprintf("%d", no), "1", 0, "C", false, 0, "")

	// Draw NIP cell with border
	pdf.CellFormat(colWidths[1], rowHeight, nip, "1", 0, "C", false, 0, "")

	// Draw NAME cell with border
	pdf.CellFormat(colWidths[2], rowHeight, name, "1", 0, "L", false, 0, "")

	// Position cell - draw border first
	posX := pdf.GetX()
	pdf.Rect(posX, currentY, colWidths[3], rowHeight, "D")

	// Draw position text with vertical centering
	yOffset := (rowHeight - textHeight) / 2
	if yOffset < 0 {
		yOffset = 0.5
	}
	pdf.SetXY(posX+1, currentY+yOffset)
	pdf.MultiCell(colWidths[3]-2, lineHeight, position, "", "L", false)

	// Move to next row
	pdf.SetXY(x, currentY+rowHeight)
}

// Helper function to convert numbers to Indonesian words
func numberToWords(num int) string {
	words := map[int]string{
		1: "satu", 2: "dua", 3: "tiga", 4: "empat", 5: "lima",
		6: "enam", 7: "tujuh", 8: "delapan", 9: "sembilan", 10: "sepuluh",
		11: "sebelas", 12: "dua belas", 13: "tiga belas", 14: "empat belas", 15: "lima belas",
		16: "enam belas", 17: "tujuh belas", 18: "delapan belas", 19: "sembilan belas", 20: "dua puluh",
	}

	if word, ok := words[num]; ok {
		return word
	}

	if num < 100 {
		tens := (num / 10) * 10
		ones := num % 10
		if ones == 0 {
			if num == 10 {
				return "sepuluh"
			}
			return words[num/10] + " puluh"
		}
		return words[tens/10] + " puluh " + words[ones]
	}

	return fmt.Sprintf("%d", num)
}
