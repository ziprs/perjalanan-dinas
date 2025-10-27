package handlers

import (
	"net/http"
	"perjalanan-dinas/backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type RepresentativeHandler struct {
	repo *repository.Repository
}

func NewRepresentativeHandler(repo *repository.Repository) *RepresentativeHandler {
	return &RepresentativeHandler{repo: repo}
}

type UpdateRepresentativeRequest struct {
	Name     string `json:"name" binding:"required"`
	Position string `json:"position" binding:"required"`
}

func (h *RepresentativeHandler) GetRepresentativeConfig(c *gin.Context) {
	config, err := h.repo.GetRepresentativeConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get representative config"})
		return
	}

	c.JSON(http.StatusOK, config)
}

func (h *RepresentativeHandler) UpdateRepresentativeConfig(c *gin.Context) {
	var req UpdateRepresentativeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing config
	config, err := h.repo.GetRepresentativeConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get representative config"})
		return
	}

	// Update fields
	config.Name = req.Name
	config.Position = req.Position

	if err := h.repo.UpdateRepresentativeConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update representative config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Representative config updated successfully",
		"data":    config,
	})
}
