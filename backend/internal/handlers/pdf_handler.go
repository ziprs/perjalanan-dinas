package handlers

import (
	"net/http"
	"perjalanan-dinas/backend/internal/repository"
	"perjalanan-dinas/backend/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PDFHandler struct {
	repo         *repository.Repository
	pdfGenerator *services.PDFGenerator
}

func NewPDFHandler(repo *repository.Repository) *PDFHandler {
	return &PDFHandler{
		repo:         repo,
		pdfGenerator: services.NewPDFGenerator(),
	}
}

// Download Nota Permintaan only
func (h *PDFHandler) DownloadNotaPermintaan(c *gin.Context) {
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

	pdfBytes, err := h.pdfGenerator.GenerateNotaPermintaan(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=nota_permintaan_"+request.RequestNumber+".pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// Download Berita Acara only
func (h *PDFHandler) DownloadBeritaAcara(c *gin.Context) {
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

	report, err := h.repo.GetTravelReportByRequestID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Travel report not found for this request"})
		return
	}

	pdfBytes, err := h.pdfGenerator.GenerateBeritaAcara(request, report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=berita_acara_"+report.ReportNumber+".pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// Download Combined PDF (both documents)
func (h *PDFHandler) DownloadCombinedPDF(c *gin.Context) {
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

	report, err := h.repo.GetTravelReportByRequestID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Travel report not found for this request"})
		return
	}

	pdfBytes, err := h.pdfGenerator.GenerateCombinedPDF(request, report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=perjalanan_dinas_"+request.RequestNumber+".pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
