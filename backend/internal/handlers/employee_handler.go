package handlers

import (
	"net/http"
	"perjalanan-dinas/backend/internal/models"
	"perjalanan-dinas/backend/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	repo *repository.Repository
}

func NewEmployeeHandler(repo *repository.Repository) *EmployeeHandler {
	return &EmployeeHandler{repo: repo}
}

type CreateEmployeeRequest struct {
	NIP        string `json:"nip" binding:"required"`
	Name       string `json:"name" binding:"required"`
	PositionID uint   `json:"position_id" binding:"required"`
}

func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate position exists
	_, err := h.repo.GetPositionByID(req.PositionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid position_id"})
		return
	}

	employee := &models.Employee{
		NIP:        req.NIP,
		Name:       req.Name,
		PositionID: req.PositionID,
	}

	if err := h.repo.CreateEmployee(employee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create employee"})
		return
	}

	// Reload with position data
	employee, _ = h.repo.GetEmployeeByID(employee.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Employee created successfully",
		"employee": employee,
	})
}

func (h *EmployeeHandler) GetAllEmployees(c *gin.Context) {
	employees, err := h.repo.GetAllEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employees"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"employees": employees})
}

func (h *EmployeeHandler) GetEmployeeByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	employee, err := h.repo.GetEmployeeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"employee": employee})
}

func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	employee, err := h.repo.GetEmployeeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	var req CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate position exists
	_, err = h.repo.GetPositionByID(req.PositionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid position_id"})
		return
	}

	employee.NIP = req.NIP
	employee.Name = req.Name
	employee.PositionID = req.PositionID

	if err := h.repo.UpdateEmployee(employee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee"})
		return
	}

	// Reload with position data
	employee, _ = h.repo.GetEmployeeByID(employee.ID)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Employee updated successfully",
		"employee": employee,
	})
}

func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	if err := h.repo.DeleteEmployee(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete employee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}
