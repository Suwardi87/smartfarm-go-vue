package controllers

import (
	"log"
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

func GetFarmerDashboard(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	role := c.MustGet("role").(string)

	// DEBUG LOG
	log.Printf("[Dashboard] Access attempt - UserID: %d, Role: %s", userID, role)

	if role != "petani" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied: farmers only", "your_role": role})
		return
	}

	data, err := analyticsService.GetFarmerDashboardData(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// Global accessor for Logging views from other controllers
func LogProductView(productID uint, userID uint) {
	if analyticsService != nil {
		analyticsService.LogView(productID, userID)
	}
}
