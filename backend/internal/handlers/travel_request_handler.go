package handlers

import (
	"fmt"
	"net/http"
	"perjalanan-dinas/backend/internal/models"
	"perjalanan-dinas/backend/internal/repository"
	"perjalanan-dinas/backend/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TravelRequestHandler struct {
	repo *repository.Repository
}

func NewTravelRequestHandler(repo *repository.Repository) *TravelRequestHandler {
	return &TravelRequestHandler{repo: repo}
}

type CreateTravelRequestRequest struct {
	EmployeeIDs     []uint `json:"employee_ids" binding:"required,min=1"`
	Purpose         string `json:"purpose" binding:"required"`
	DeparturePlace  string `json:"departure_place"`
	Destination     string `json:"destination" binding:"required"`
	DestinationType string `json:"destination_type" binding:"required,oneof=in_province outside_province abroad"`
	DepartureDate   string `json:"departure_date" binding:"required"` // Format: 2006-01-02
	ReturnDate      string `json:"return_date" binding:"required"`    // Format: 2006-01-02
	Transportation  string `json:"transportation" binding:"required"` // angkutan umum, pesawat, kereta api
}

func (h *TravelRequestHandler) CreateTravelRequest(c *gin.Context) {
	var req CreateTravelRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default departure place if not provided
	if req.DeparturePlace == "" {
		req.DeparturePlace = "Surabaya"
	}

	// Parse dates
	departureDate, err := time.Parse("2006-01-02", req.DepartureDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid departure date format. Use YYYY-MM-DD"})
		return
	}

	returnDate, err := time.Parse("2006-01-02", req.ReturnDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid return date format. Use YYYY-MM-DD"})
		return
	}

	// Validate dates
	if returnDate.Before(departureDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Return date cannot be before departure date"})
		return
	}

	// Calculate duration
	durationDays := repository.CalculateDurationDays(departureDate, returnDate)

	// Validate all employees exist and get first employee for position code
	var firstEmployee *models.Employee
	for i, empID := range req.EmployeeIDs {
		employee, err := h.repo.GetEmployeeByID(empID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Employee with ID %d not found", empID)})
			return
		}
		if i == 0 {
			firstEmployee = employee
		}
	}

	// Get position from first employee
	position := firstEmployee.Position

	// Calculate allowance rate based on destination type
	var allowanceRate int
	switch req.DestinationType {
	case "in_province":
		allowanceRate = position.AllowanceInProvince
	case "outside_province":
		allowanceRate = position.AllowanceOutsideProvince
	case "abroad":
		allowanceRate = position.AllowanceAbroad
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid destination_type"})
		return
	}

	// Calculate total allowance: duration * rate per day * number of employees
	totalAllowance := durationDays * allowanceRate * len(req.EmployeeIDs)

	// Start transaction
	tx := h.repo.GetDB().Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Defer rollback in case of error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get next sequence within transaction
	var config models.NumberingConfig
	if err := tx.First(&config).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get numbering config"})
		return
	}

	// Increment sequence untuk Nota Permintaan (misal: 0001)
	config.LastRequestSequence++
	requestNumber := utils.GenerateRequestNumber(config.LastRequestSequence, position.Code)

	// Berita Acara menggunakan nomor yang SAMA dengan Nota Permintaan
	reportNumber := requestNumber

	// Save sequence
	if err := tx.Save(&config).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sequence"})
		return
	}

	// Create travel request
	// Note: ReportNumber is left empty and will be filled when travel report is created
	travelRequest := &models.TravelRequest{
		Purpose:         req.Purpose,
		DeparturePlace:  req.DeparturePlace,
		Destination:     req.Destination,
		DestinationType: req.DestinationType,
		DepartureDate:   departureDate,
		ReturnDate:      returnDate,
		DurationDays:    durationDays,
		Transportation:  req.Transportation,
		TotalAllowance:  totalAllowance,
		RequestNumber:   requestNumber,
		Status:          "pending",
	}

	if err := tx.Create(travelRequest).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create travel request"})
		return
	}

	// Create travel request employee relations
	for _, empID := range req.EmployeeIDs {
		empRel := &models.TravelRequestEmployee{
			TravelRequestID: travelRequest.ID,
			EmployeeID:      empID,
		}
		if err := tx.Create(empRel).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link employees to travel request"})
			return
		}
	}

	// Get representative config
	repConfig, err := h.repo.GetActiveRepresentativeConfig()
	if err != nil {
		// Fallback to default if config not found
		repConfig = &models.RepresentativeConfig{
			Name:     "M. MACHFUD HIDAYAT",
			Position: "Vice President",
		}
	}

	// Auto-create travel report with different sequence number
	travelReport := &models.TravelReport{
		TravelRequestID:        travelRequest.ID,
		ReportNumber:           reportNumber,
		RepresentativeName:     repConfig.Name,
		RepresentativePosition: repConfig.Position,
	}

	if err := tx.Create(travelReport).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create travel report"})
		return
	}

	// Update travel request with report number
	travelRequest.ReportNumber = reportNumber
	if err := tx.Save(travelRequest).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update travel request"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Reload with employee data
	travelRequest, _ = h.repo.GetTravelRequestByID(travelRequest.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message":        "Travel request created successfully",
		"travel_request": travelRequest,
	})
}

func (h *TravelRequestHandler) GetAllTravelRequests(c *gin.Context) {
	requests, err := h.repo.GetAllTravelRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch travel requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"travel_requests": requests})
}

func (h *TravelRequestHandler) GetTravelRequestByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid travel request ID"})
		return
	}

	request, err := h.repo.GetTravelRequestByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Travel request not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"travel_request": request})
}

func (h *TravelRequestHandler) DeleteTravelRequest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid travel request ID"})
		return
	}

	if err := h.repo.DeleteTravelRequest(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete travel request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Travel request deleted successfully"})
}
