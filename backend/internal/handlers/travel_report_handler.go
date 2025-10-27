package handlers

import (
	"net/http"
	"perjalanan-dinas/backend/internal/models"
	"perjalanan-dinas/backend/internal/repository"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TravelReportHandler struct {
	repo *repository.Repository
}

func NewTravelReportHandler(repo *repository.Repository) *TravelReportHandler {
	return &TravelReportHandler{repo: repo}
}

type CreateTravelReportRequest struct {
	TravelRequestID        uint              `json:"travel_request_id" binding:"required"`
	RepresentativeName     string            `json:"representative_name" binding:"required"`
	RepresentativePosition string            `json:"representative_position" binding:"required"`
	VisitProofs            []VisitProofInput `json:"visit_proofs" binding:"required"`
}

type VisitProofInput struct {
	Date           string `json:"date" binding:"required"`           // Format: 2006-01-02
	DepartFrom     string `json:"depart_from" binding:"required"`
	StayOrStopAt   string `json:"stay_or_stop_at"`
	ArriveAt       string `json:"arrive_at" binding:"required"`
	SignatureProof string `json:"signature_proof"`
}

func (h *TravelReportHandler) CreateTravelReport(c *gin.Context) {
	var req CreateTravelReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get travel request
	travelRequest, err := h.repo.GetTravelRequestByID(req.TravelRequestID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Travel request not found"})
		return
	}

	// Check if report already exists
	existingReport, err := h.repo.GetTravelReportByRequestID(req.TravelRequestID)
	if err == nil && existingReport.ID > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Travel report already exists for this request"})
		return
	}

	// Create travel report
	travelReport := &models.TravelReport{
		TravelRequestID:        req.TravelRequestID,
		ReportNumber:           travelRequest.ReportNumber,
		RepresentativeName:     req.RepresentativeName,
		RepresentativePosition: req.RepresentativePosition,
	}

	if err := h.repo.CreateTravelReport(travelReport); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create travel report"})
		return
	}

	// Create visit proofs
	for _, proofInput := range req.VisitProofs {
		date, err := time.Parse("2006-01-02", proofInput.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format in visit proof. Use YYYY-MM-DD"})
			return
		}

		proof := &models.VisitProof{
			TravelReportID: travelReport.ID,
			Date:           date,
			DepartFrom:     proofInput.DepartFrom,
			StayOrStopAt:   proofInput.StayOrStopAt,
			ArriveAt:       proofInput.ArriveAt,
			SignatureProof: proofInput.SignatureProof,
		}

		if err := h.repo.CreateVisitProof(proof); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create visit proof"})
			return
		}
	}

	// Update travel request status
	travelRequest.Status = "completed"
	h.repo.UpdateTravelRequest(travelRequest)

	// Reload with visit proofs
	travelReport, _ = h.repo.GetTravelReportByRequestID(req.TravelRequestID)

	c.JSON(http.StatusCreated, gin.H{
		"message":       "Travel report created successfully",
		"travel_report": travelReport,
	})
}

func (h *TravelReportHandler) GetTravelReportByRequestID(c *gin.Context) {
	idStr := c.Param("request_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid travel request ID"})
		return
	}

	report, err := h.repo.GetTravelReportByRequestID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Travel report not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"travel_report": report})
}
