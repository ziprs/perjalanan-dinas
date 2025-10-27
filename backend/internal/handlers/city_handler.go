package handlers

import (
	"net/http"
	"perjalanan-dinas/backend/internal/database"

	"github.com/gin-gonic/gin"
)

type CityHandler struct{}

func NewCityHandler() *CityHandler {
	return &CityHandler{}
}

type CityResponse struct {
	Name            string `json:"name"`
	DestinationType string `json:"destination_type"`
}

func (h *CityHandler) GetAllCities(c *gin.Context) {
	cities := make([]CityResponse, len(database.AllCities))

	for i, city := range database.AllCities {
		cities[i] = CityResponse{
			Name:            city.Name,
			DestinationType: city.DestinationType,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"cities": cities,
	})
}
