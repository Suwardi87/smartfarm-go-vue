package controllers

import (
	"net/http"
	"smartfarm-api/config"
	"smartfarm-api/repositories"
	"smartfarm-api/services"

	"github.com/gin-gonic/gin"
)

var analyticsService services.AnalyticsService

func InitAnalyticsController() {
	analyticsRepo := repositories.NewAnalyticsRepository(config.DB)
	analyticsService = services.NewAnalyticsService(analyticsRepo)
}

func GetTrendingProducts(c *gin.Context) {
	products, err := analyticsService.GetTrendingProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": products})
}

// Global accessor for Logging views from other controllers
func LogProductView(productID uint, userID uint) {
	if analyticsService != nil {
		analyticsService.LogView(productID, userID)
	}
}
