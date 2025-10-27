package handlers

import (
	"net/http"
	"perjalanan-dinas/backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type PositionHandler struct {
	repo *repository.Repository
}

func NewPositionHandler(repo *repository.Repository) *PositionHandler {
	return &PositionHandler{repo: repo}
}

func (h *PositionHandler) GetAllPositions(c *gin.Context) {
	positions, err := h.repo.GetAllPositions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch positions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"positions": positions})
}
