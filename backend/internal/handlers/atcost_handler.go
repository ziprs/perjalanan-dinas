package handlers

import (
	"encoding/json"
	"net/http"
	"perjalanan-dinas/backend/internal/repository"
	"perjalanan-dinas/backend/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AtCostHandler struct {
	service      *services.AtCostService
	pdfGenerator *services.AtCostPDFGenerator
}

func NewAtCostHandler(repo *repository.Repository) *AtCostHandler {
	return &AtCostHandler{
		service:      services.NewAtCostService(repo),
		pdfGenerator: services.NewAtCostPDFGenerator(repo),
	}
}

// UploadReceipt handles PDF receipt upload and parsing
func (h *AtCostHandler) UploadReceipt(c *gin.Context) {
	file, err := c.FormFile("receipt")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Process upload and parse
	receiptData, filepath, err := h.service.ProcessReceiptUpload(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert parsed data to JSON string
	parser := services.NewReceiptParser()
	parsedJSON, _ := parser.ToJSON(receiptData)

	c.JSON(http.StatusOK, gin.H{
		"message":     "Receipt uploaded and parsed successfully",
		"file_path":   filepath,
		"file_name":   file.Filename,
		"parsed_data": receiptData,
		"parsed_json": parsedJSON,
	})
}

// CreateAtCostClaim creates a new At-Cost claim
func (h *AtCostHandler) CreateAtCostClaim(c *gin.Context) {
	var req services.CreateAtCostClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claim, err := h.service.CreateAtCostClaim(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "At-Cost claim created successfully",
		"claim":   claim,
	})
}

// GetAtCostClaim retrieves a single At-Cost claim by ID
func (h *AtCostHandler) GetAtCostClaim(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid claim ID"})
		return
	}

	claim, err := h.service.GetAtCostClaimByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Claim not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"claim": claim})
}

// GetAtCostClaimByTravelRequest retrieves claim by travel request ID
func (h *AtCostHandler) GetAtCostClaimByTravelRequest(c *gin.Context) {
	idStr := c.Param("travel_request_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid travel request ID"})
		return
	}

	claim, err := h.service.GetAtCostClaimByTravelRequestID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Claim not found for this travel request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"claim": claim})
}

// GetAllAtCostClaims retrieves all At-Cost claims
func (h *AtCostHandler) GetAllAtCostClaims(c *gin.Context) {
	claims, err := h.service.GetAllAtCostClaims()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve claims"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"claims": claims})
}

// UpdateClaimStatus updates the status of a claim
func (h *AtCostHandler) UpdateClaimStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid claim ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=pending approved rejected"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateClaimStatus(uint(id), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

// DeleteAtCostClaim deletes a claim
func (h *AtCostHandler) DeleteAtCostClaim(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid claim ID"})
		return
	}

	if err := h.service.DeleteAtCostClaim(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete claim"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Claim deleted successfully"})
}

// DownloadReceipt downloads a receipt file
func (h *AtCostHandler) DownloadReceipt(c *gin.Context) {
	idStr := c.Param("receipt_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receipt ID"})
		return
	}

	filepath, err := h.service.GetReceiptFile(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
		return
	}

	c.File(filepath)
}

// ParseReceiptManual allows manual text input for parsing (for testing)
func (h *AtCostHandler) ParseReceiptManual(c *gin.Context) {
	var req struct {
		Text string `json:"text" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parser := services.NewReceiptParser()
	data, err := parser.ParseReceiptText(req.Text)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parsedJSON, _ := json.Marshal(data)

	c.JSON(http.StatusOK, gin.H{
		"message":     "Receipt parsed successfully",
		"parsed_data": data,
		"parsed_json": string(parsedJSON),
	})
}

// DownloadNotaAtCost generates and downloads the Nota At-Cost PDF
func (h *AtCostHandler) DownloadNotaAtCost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid claim ID"})
		return
	}

	filepath, err := h.pdfGenerator.GenerateNotaAtCost(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.File(filepath)
}

// DownloadCombinedAtCost generates and downloads combined PDF (nota + receipts)
func (h *AtCostHandler) DownloadCombinedAtCost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid claim ID"})
		return
	}

	filepath, err := h.pdfGenerator.GenerateCombinedPDF(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.File(filepath)
}
