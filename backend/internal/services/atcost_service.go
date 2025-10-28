package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"perjalanan-dinas/backend/internal/models"
	"perjalanan-dinas/backend/internal/repository"
	"perjalanan-dinas/backend/internal/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

type AtCostService struct {
	repo          *repository.Repository
	parser        *ReceiptParser
	pdfExtractor  *PDFExtractor
	uploadDir     string
}

func NewAtCostService(repo *repository.Repository) *AtCostService {
	uploadDir := "./uploads/receipts"
	// Create upload directory if not exists
	os.MkdirAll(uploadDir, 0755)

	return &AtCostService{
		repo:         repo,
		parser:       NewReceiptParser(),
		pdfExtractor: NewPDFExtractor(),
		uploadDir:    uploadDir,
	}
}

// CreateAtCostClaimRequest represents the request to create a claim
type CreateAtCostClaimRequest struct {
	TravelRequestID uint                       `json:"travel_request_id"`
	ClaimItems      []CreateAtCostClaimItemReq `json:"claim_items"`
}

type CreateAtCostClaimItemReq struct {
	EmployeeID        uint                      `json:"employee_id"`
	TransportCost     int                       `json:"transport_cost"`
	AccommodationCost int                       `json:"accommodation_cost"`
	Receipts          []CreateAtCostReceiptReq  `json:"receipts"`
}

type CreateAtCostReceiptReq struct {
	Type            string `json:"type"` // flight, hotel, train
	ReceiptNumber   string `json:"receipt_number"`
	ReceiptDate     string `json:"receipt_date"` // YYYY-MM-DD
	Vendor          string `json:"vendor"`
	Description     string `json:"description"`
	Amount          int    `json:"amount"`
	PassengerName   string `json:"passenger_name"`
	RouteOrLocation string `json:"route_or_location"`
	FilePath        string `json:"file_path"` // Set after upload
	FileName        string `json:"file_name"`
	ParsedData      string `json:"parsed_data,omitempty"`
}

// ProcessReceiptUpload handles file upload and parsing
func (s *AtCostService) ProcessReceiptUpload(file *multipart.FileHeader) (*ReceiptData, string, error) {
	// Validate file
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".pdf") {
		return nil, "", fmt.Errorf("only PDF files are allowed")
	}

	if file.Size > 10*1024*1024 { // 10MB limit
		return nil, "", fmt.Errorf("file size exceeds 10MB limit")
	}

	// Generate unique filename
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s", timestamp, file.Filename)
	filepath := filepath.Join(s.uploadDir, filename)

	// Save file
	src, err := file.Open()
	if err != nil {
		return nil, "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filepath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, "", fmt.Errorf("failed to save file: %w", err)
	}

	// Extract text from PDF
	text, err := s.pdfExtractor.ExtractTextClean(filepath)
	if err != nil {
		return nil, filepath, fmt.Errorf("failed to extract text from PDF: %w", err)
	}

	// Parse receipt
	receiptData, err := s.parser.ParseReceiptText(text)
	if err != nil {
		return nil, filepath, fmt.Errorf("failed to parse receipt: %w", err)
	}

	return receiptData, filepath, nil
}

// CreateAtCostClaim creates a new At-Cost claim with all items and receipts
func (s *AtCostService) CreateAtCostClaim(req *CreateAtCostClaimRequest) (*models.AtCostClaim, error) {
	// Get travel request
	travelRequest, err := s.repo.GetTravelRequestByID(req.TravelRequestID)
	if err != nil {
		return nil, fmt.Errorf("travel request not found: %w", err)
	}

	// Get representative config
	repConfig, err := s.repo.GetActiveRepresentativeConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get representative config: %w", err)
	}

	// Generate NEW claim number with different sequence from travel request
	// Extract position code from travel request number
	positionCode := utils.ExtractPositionCodeFromRequestNumber(travelRequest.RequestNumber)
	if positionCode == "" {
		return nil, fmt.Errorf("failed to extract position code from request number")
	}

	// Get next sequence for at-cost (different from travel request sequence)
	seq, err := s.repo.GetNextRequestSequence()
	if err != nil {
		return nil, fmt.Errorf("failed to get next sequence: %w", err)
	}

	// Generate new claim number with format: 064/{seq}/DIB/{code}/Nota
	claimNumber := utils.GenerateRequestNumber(seq, positionCode)

	// Calculate total amount
	var totalAmount int
	for _, item := range req.ClaimItems {
		totalAmount += item.TransportCost + item.AccommodationCost
	}

	// Begin transaction
	var claim *models.AtCostClaim
	err = s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		// Create claim
		claim = &models.AtCostClaim{
			TravelRequestID:        req.TravelRequestID,
			ClaimNumber:            claimNumber,
			RepresentativeName:     repConfig.Name,
			RepresentativePosition: repConfig.Position,
			Status:                 "pending",
			TotalAmount:            totalAmount,
		}

		if err := tx.Create(claim).Error; err != nil {
			return fmt.Errorf("failed to create claim: %w", err)
		}

		// Create claim items
		for _, itemReq := range req.ClaimItems {
			totalCost := itemReq.TransportCost + itemReq.AccommodationCost

			claimItem := &models.AtCostClaimItem{
				AtCostClaimID:     claim.ID,
				EmployeeID:        itemReq.EmployeeID,
				TransportCost:     itemReq.TransportCost,
				AccommodationCost: itemReq.AccommodationCost,
				TotalCost:         totalCost,
			}

			if err := tx.Create(claimItem).Error; err != nil {
				return fmt.Errorf("failed to create claim item: %w", err)
			}

			// Create receipts
			for _, receiptReq := range itemReq.Receipts {
				receiptDate, _ := time.Parse("2006-01-02", receiptReq.ReceiptDate)

				receipt := &models.AtCostReceipt{
					ClaimItemID:     claimItem.ID,
					ReceiptNumber:   receiptReq.ReceiptNumber,
					ReceiptDate:     receiptDate,
					Vendor:          receiptReq.Vendor,
					Type:            receiptReq.Type,
					Description:     receiptReq.Description,
					Amount:          receiptReq.Amount,
					FilePath:        receiptReq.FilePath,
					FileName:        receiptReq.FileName,
					PassengerName:   receiptReq.PassengerName,
					RouteOrLocation: receiptReq.RouteOrLocation,
					ParsedData:      receiptReq.ParsedData,
				}

				if err := tx.Create(receipt).Error; err != nil {
					return fmt.Errorf("failed to create receipt: %w", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload with all relationships
	return s.repo.GetAtCostClaimByID(claim.ID)
}

// GetAtCostClaimByID retrieves a claim with all relationships
func (s *AtCostService) GetAtCostClaimByID(id uint) (*models.AtCostClaim, error) {
	return s.repo.GetAtCostClaimByID(id)
}

// GetAtCostClaimByTravelRequestID retrieves claim by travel request ID
func (s *AtCostService) GetAtCostClaimByTravelRequestID(travelRequestID uint) (*models.AtCostClaim, error) {
	return s.repo.GetAtCostClaimByTravelRequestID(travelRequestID)
}

// GetAllAtCostClaims retrieves all claims
func (s *AtCostService) GetAllAtCostClaims() ([]models.AtCostClaim, error) {
	return s.repo.GetAllAtCostClaims()
}

// UpdateClaimStatus updates the status of a claim
func (s *AtCostService) UpdateClaimStatus(id uint, status string) error {
	claim, err := s.repo.GetAtCostClaimByID(id)
	if err != nil {
		return err
	}

	claim.Status = status
	return s.repo.UpdateAtCostClaim(claim)
}

// DeleteAtCostClaim deletes a claim and all associated data
func (s *AtCostService) DeleteAtCostClaim(id uint) error {
	// Get claim to access file paths
	claim, err := s.repo.GetAtCostClaimByID(id)
	if err != nil {
		return err
	}

	// Delete physical receipt files
	for _, item := range claim.ClaimItems {
		for _, receipt := range item.Receipts {
			if receipt.FilePath != "" {
				os.Remove(receipt.FilePath) // Ignore errors
			}
		}
	}

	// Delete from database (cascades to items and receipts)
	return s.repo.DeleteAtCostClaim(id)
}

// GetReceiptFile returns the file path for download
func (s *AtCostService) GetReceiptFile(receiptID uint) (string, error) {
	var receipt models.AtCostReceipt
	if err := s.repo.GetDB().First(&receipt, receiptID).Error; err != nil {
		return "", err
	}
	return receipt.FilePath, nil
}
