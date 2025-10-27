package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"perjalanan-dinas/backend/internal/repository"
	"perjalanan-dinas/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type ExcelHandler struct {
	repo      *repository.Repository
	excelGen  *services.ExcelGenerator
}

func NewExcelHandler(repo *repository.Repository) *ExcelHandler {
	return &ExcelHandler{
		repo:     repo,
		excelGen: services.NewExcelGenerator(),
	}
}

func (h *ExcelHandler) ExportMonthlyAllowance(c *gin.Context) {
	// Get year and month from query parameters
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	// Default to current month if not provided
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	if yearStr != "" {
		y, err := strconv.Atoi(yearStr)
		if err != nil || y < 2000 || y > 2100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year parameter"})
			return
		}
		year = y
	}

	if monthStr != "" {
		m, err := strconv.Atoi(monthStr)
		if err != nil || m < 1 || m > 12 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month parameter"})
			return
		}
		month = m
	}

	// Get all travel requests
	requests, err := h.repo.GetAllTravelRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch travel requests"})
		return
	}

	// Generate Excel file
	excelData, err := h.excelGen.GenerateMonthlyAllowanceReport(requests, year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate Excel file"})
		return
	}

	// Set headers for file download
	monthNames := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	filename := fmt.Sprintf("Rekap_Iuran_%s_%d.xlsx", monthNames[month-1], year)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelData)
}
