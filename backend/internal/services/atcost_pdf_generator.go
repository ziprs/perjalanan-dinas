package services

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"perjalanan-dinas/backend/internal/models"
	"perjalanan-dinas/backend/internal/repository"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type AtCostPDFGenerator struct {
	repo *repository.Repository
}

func NewAtCostPDFGenerator(repo *repository.Repository) *AtCostPDFGenerator {
	return &AtCostPDFGenerator{repo: repo}
}

// GenerateNotaAtCost generates the Nota At-Cost PDF
func (g *AtCostPDFGenerator) GenerateNotaAtCost(claimID uint) (string, error) {
	// Get claim with all relationships
	claim, err := g.repo.GetAtCostClaimByID(claimID)
	if err != nil {
		return "", fmt.Errorf("failed to get claim: %w", err)
	}

	// Initialize PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()
	pdf.SetAutoPageBreak(true, 15)

	// Draw header with logo
	g.drawHeader(pdf, claim)

	// Draw NOTA title and intro
	g.drawNotaIntro(pdf, claim)

	// Draw table 1: Bukti Perjalanan Dinas
	pdf.Ln(5)
	g.drawTravelProofSection(pdf, claim)

	// Draw table 2: Bukti Pembelian Akomodasi
	pdf.Ln(8)
	g.drawAccommodationSection(pdf, claim)

	// Draw closing and signatures
	pdf.Ln(8)
	g.drawClosingAndSignatures(pdf, claim)

	// Save PDF
	outputDir := "./uploads/pdfs/atcost"
	filename := fmt.Sprintf("nota_atcost_%d_%s.pdf", claim.ID, time.Now().Format("20060102_150405"))
	outputPath := filepath.Join(outputDir, filename)

	if err := pdf.OutputFileAndClose(outputPath); err != nil {
		return "", fmt.Errorf("failed to save PDF: %w", err)
	}

	return outputPath, nil
}

func (g *AtCostPDFGenerator) drawHeader(pdf *gofpdf.Fpdf, claim *models.AtCostClaim) {
	// Save current Y position
	startY := pdf.GetY()

	// Add Bank Jatim logo on the top right corner
	logoPath := "./assets/images/bank jatim.png"
	pdf.ImageOptions(logoPath, 165, startY, 25, 0, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	// Add some space below logo before text
	pdf.Ln(8)

	// Draw Nomor on the left and Date on the right - SAME LINE
	pdf.SetFont("Arial", "", 10)

	// Nomor on the left
	pdf.Cell(30, 5, "Nomor")
	pdf.Cell(5, 5, ":")
	pdf.Cell(75, 5, claim.ClaimNumber)

	// Date on the right side - same line
	currentDate := time.Now().Format("2 January 2006")
	pdf.CellFormat(0, 5, "Surabaya, "+currentDate, "", 1, "R", false, 0, "")
	pdf.Ln(8)
}

func (g *AtCostPDFGenerator) drawNotaIntro(pdf *gofpdf.Fpdf, claim *models.AtCostClaim) {
	// NOTA Title - centered and bold
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 8, "N O T A", "", 1, "C", false, 0, "")
	pdf.Ln(8)

	// Formal letter format
	pdf.SetFont("Arial", "", 10)

	// Kepada
	pdf.Cell(30, 5, "Kepada Yth.")
	pdf.Cell(5, 5, ":")
	pdf.Cell(0, 5, "VP Divisi Human Capital")
	pdf.Ln(5)

	// Dari
	pdf.Cell(30, 5, "Dari")
	pdf.Cell(5, 5, ":")
	pdf.Cell(0, 5, "VP Divisi Digital Banking")
	pdf.Ln(5)

	// Perihal
	pdf.Cell(30, 5, "Perihal")
	pdf.Cell(5, 5, ":")
	pdf.MultiCell(0, 5, "Penyerahan Bukti Surat Perjalanan Dinas dan Permohonan Penggantian Biaya Transportasi", "", "L", false)
	pdf.Ln(2)

	// Divider line
	pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
	pdf.Ln(5)

	// Opening paragraph with JUSTIFIED alignment
	tr := claim.TravelRequest
	departureDate := tr.DepartureDate.Format("2 January 2006")
	returnDate := tr.ReturnDate.Format("2 January 2006")

	openingText := fmt.Sprintf("Menindaklanjuti Nota ke Divisi Human Capital Nomor: %s tanggal %s terkait Surat Permohonan Perjalanan Dinas ke %s pada tanggal %s s/d %s, maka bersama ini disampaikan hal-hal sebagai berikut :",
		tr.RequestNumber, departureDate, tr.Destination, departureDate, returnDate)

	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 5, openingText, "", "J", false) // Changed to "J" for justify
	pdf.Ln(3)
}

func (g *AtCostPDFGenerator) drawTravelProofSection(pdf *gofpdf.Fpdf, claim *models.AtCostClaim) {
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 5, "1. Penyampaian Bukti Perjalanan Dinas (terlampir) untuk pegawai sebagai berikut :")
	pdf.Ln(6)

	// Table header
	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(240, 240, 240)

	// Column widths: No(8), Nama(45), Jabatan(42), Tanggal(38), Tujuan(47) = 180mm
	// First row - with vertical centering
	pdf.CellFormat(8, 12, "No.", "1", 0, "CM", true, 0, "")
	pdf.CellFormat(45, 12, "Nama Pegawai (NIP)", "1", 0, "CM", true, 0, "")
	pdf.CellFormat(42, 12, "Jabatan", "1", 0, "CM", true, 0, "")

	// Date column header - merged and centered
	currentX := pdf.GetX()
	currentY := pdf.GetY()
	pdf.CellFormat(38, 6, "Tanggal Berangkat -", "LTR", 0, "C", true, 0, "")
	pdf.CellFormat(47, 12, "Tujuan", "1", 1, "CM", true, 0, "")

	// Second header row for date continuation
	pdf.SetXY(currentX, currentY+6)
	pdf.CellFormat(38, 6, "Tanggal Pulang", "LRB", 0, "C", true, 0, "")
	pdf.Ln(6)

	// Table rows
	pdf.SetFont("Arial", "", 8)
	tr := claim.TravelRequest

	for i, item := range claim.ClaimItems {
		startX := pdf.GetX()
		startY := pdf.GetY()

		// No.
		pdf.SetXY(startX, startY)
		pdf.CellFormat(8, 14, fmt.Sprintf("%d.", i+1), "1", 0, "C", false, 0, "")

		// Name with NIP - centered vertically, showing actual NIP
		pdf.SetXY(startX+8, startY+3.5)
		nameText := fmt.Sprintf("%s\n(%s)", item.Employee.Name, item.Employee.NIP)
		pdf.MultiCell(45, 4, nameText, "", "L", false)

		// Position - moved HIGHER UP to be truly centered (same as name position)
		pdf.SetXY(startX+53, startY+3.5)
		pdf.MultiCell(42, 4, item.Employee.Position.Title, "", "L", false)

		// Dates
		pdf.SetXY(startX+95, startY+1)
		dateRange := fmt.Sprintf("%s s/d\n%s\n(%d hari)",
			tr.DepartureDate.Format("2 Jan 2006"),
			tr.ReturnDate.Format("2 Jan 2006"),
			tr.DurationDays)
		pdf.MultiCell(38, 4, dateRange, "", "C", false)

		// Destination/Purpose - moved HIGHER UP to match position
		pdf.SetXY(startX+133, startY+3.5)
		pdf.MultiCell(47, 4, tr.Purpose, "", "L", false)

		// Draw borders manually for alignment
		currentY := pdf.GetY()
		if currentY-startY < 14 {
			currentY = startY + 14
		}

		// Draw cell borders
		pdf.Rect(startX+8, startY, 45, currentY-startY, "D")
		pdf.Rect(startX+53, startY, 42, currentY-startY, "D")
		pdf.Rect(startX+95, startY, 38, currentY-startY, "D")
		pdf.Rect(startX+133, startY, 47, currentY-startY, "D")

		// Move to next row
		pdf.SetXY(startX, currentY)
	}
}

func (g *AtCostPDFGenerator) drawAccommodationSection(pdf *gofpdf.Fpdf, claim *models.AtCostClaim) {
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 5, "2. Penyampaian Bukti Pembelian Akomodasi (terlampir) sebagai berikut :")
	pdf.Ln(6)

	// Table header with vertical centering
	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(240, 240, 240)

	// Column widths: No(8), Nama(45), Jabatan(42), Transport(28), Penginapan(28), Total(29) = 180mm
	pdf.CellFormat(8, 8, "No.", "1", 0, "CM", true, 0, "")
	pdf.CellFormat(45, 8, "Nama Pegawai (NIP)", "1", 0, "CM", true, 0, "")
	pdf.CellFormat(42, 8, "Jabatan", "1", 0, "CM", true, 0, "")
	pdf.CellFormat(28, 8, "Transportasi", "1", 0, "CM", true, 0, "")
	pdf.CellFormat(28, 8, "Penginapan", "1", 0, "CM", true, 0, "")
	pdf.CellFormat(29, 8, "Total", "1", 1, "CM", true, 0, "")

	// Table rows
	pdf.SetFont("Arial", "", 8)

	for i, item := range claim.ClaimItems {
		startX := pdf.GetX()
		startY := pdf.GetY()

		// No.
		pdf.SetXY(startX, startY)
		pdf.CellFormat(8, 14, fmt.Sprintf("%d.", i+1), "1", 0, "C", false, 0, "")

		// Name with NIP - centered vertically, showing actual NIP
		pdf.SetXY(startX+8, startY+3.5)
		nameText := fmt.Sprintf("%s\n(%s)", item.Employee.Name, item.Employee.NIP)
		pdf.MultiCell(45, 4, nameText, "", "L", false)

		// Position - moved HIGHER UP to be truly centered (same as name position)
		pdf.SetXY(startX+53, startY+3.5)
		pdf.MultiCell(42, 4, item.Employee.Position.Title, "", "L", false)

		// Calculate row height
		currentY := pdf.GetY()
		if currentY-startY < 14 {
			currentY = startY + 14
		}
		rowHeight := currentY - startY

		// Draw borders for text columns
		pdf.Rect(startX+8, startY, 45, rowHeight, "D")
		pdf.Rect(startX+53, startY, 42, rowHeight, "D")

		// Transport cost
		pdf.SetXY(startX+95, startY)
		pdf.CellFormat(28, rowHeight, formatCurrency(item.TransportCost), "1", 0, "R", false, 0, "")

		// Accommodation cost
		accommodationStr := formatCurrency(item.AccommodationCost)
		if item.AccommodationCost == 0 {
			accommodationStr = "-"
		}
		pdf.SetXY(startX+123, startY)
		pdf.CellFormat(28, rowHeight, accommodationStr, "1", 0, "C", false, 0, "")

		// Total
		pdf.SetXY(startX+151, startY)
		pdf.CellFormat(29, rowHeight, formatCurrency(item.TotalCost), "1", 0, "R", false, 0, "")

		// Move to next row
		pdf.SetXY(startX, currentY)
	}

	// Total row
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(220, 220, 220)
	pdf.CellFormat(151, 8, "TOTAL", "1", 0, "C", true, 0, "")
	pdf.CellFormat(29, 8, formatCurrency(claim.TotalAmount), "1", 1, "R", true, 0, "")
}

func (g *AtCostPDFGenerator) drawClosingAndSignatures(pdf *gofpdf.Fpdf, claim *models.AtCostClaim) {
	// Closing paragraph with justify
	pdf.SetFont("Arial", "", 10)
	closingText := "Demikian dapat diproses sesuai dengan ketentuan yang berlaku. Atas perhatian dan kerjasamanya disampaikan terimakasih."
	pdf.MultiCell(0, 5, closingText, "", "J", false)
	pdf.Ln(10)

	// Calculate center position for signature
	// Page width = 210mm, margins = 15mm left + 15mm right = 180mm usable
	// We want signature block centered

	// Signature section - centered
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, "DIVISI DIGITAL BANKING", "", 1, "C", false, 0, "")
	pdf.Ln(15)

	// Representative name - centered, bold, and underlined
	pdf.SetFont("Arial", "BU", 10)
	repName := strings.ToUpper(claim.RepresentativeName)
	pdf.CellFormat(0, 5, repName, "", 1, "C", false, 0, "")
	pdf.Ln(2)

	// Representative position - centered (closer to name)
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, claim.RepresentativePosition, "", 1, "C", false, 0, "")
}

// GenerateCombinedPDF generates combined PDF with nota + all receipts
func (g *AtCostPDFGenerator) GenerateCombinedPDF(claimID uint) (string, error) {
	// Generate nota PDF
	notaPath, err := g.GenerateNotaAtCost(claimID)
	if err != nil {
		return "", err
	}

	// Get claim to access receipts
	claim, err := g.repo.GetAtCostClaimByID(claimID)
	if err != nil {
		return "", err
	}

	// Collect all receipt file paths
	var receiptPaths []string
	for _, item := range claim.ClaimItems {
		for _, receipt := range item.Receipts {
			if receipt.FilePath != "" {
				receiptPaths = append(receiptPaths, receipt.FilePath)
			}
		}
	}

	// If no receipts, just return the nota
	if len(receiptPaths) == 0 {
		return notaPath, nil
	}

	// Combine PDFs using pdfunite
	outputDir := "./uploads/pdfs/atcost"
	combinedFilename := fmt.Sprintf("combined_atcost_%d_%s.pdf", claim.ID, time.Now().Format("20060102_150405"))
	combinedPath := filepath.Join(outputDir, combinedFilename)

	// Build pdfunite command: pdfunite input1.pdf input2.pdf ... output.pdf
	args := []string{notaPath}
	args = append(args, receiptPaths...)
	args = append(args, combinedPath)

	cmd := exec.Command("pdfunite", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to combine PDFs: %s - %w", string(output), err)
	}

	return combinedPath, nil
}

// Helper function to format currency
func formatCurrency(amount int) string {
	if amount == 0 {
		return "Rp 0,-"
	}

	// Convert to string
	s := fmt.Sprintf("%d", amount)

	// Add thousands separator
	var result strings.Builder
	n := len(s)
	for i, digit := range s {
		if i > 0 && (n-i)%3 == 0 {
			result.WriteRune('.')
		}
		result.WriteRune(digit)
	}

	return "Rp " + result.String() + ",-"
}
