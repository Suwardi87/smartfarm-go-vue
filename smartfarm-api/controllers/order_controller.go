package controllers

import (
	"log"
	"net/http"
	"smartfarm-api/config"
	"smartfarm-api/dto"
	"smartfarm-api/repositories"
	"smartfarm-api/services"
	"strings"

	"github.com/gin-gonic/gin"
)

var orderService services.OrderService

func InitOrderController() {
	db := config.DB
	orderRepo := repositories.NewOrderRepository(db)
	productRepo := repositories.NewProductRepository(db) // Need product repo too
	orderService = services.NewOrderService(orderRepo, productRepo)
}

func CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(uint)

	res, err := orderService.CreateOrder(req, userID)
	if err != nil {
		// Use 400 for business logic errors
		status := http.StatusInternalServerError
		errMsg := err.Error()
		if errMsg == "cart is empty" ||
			errMsg == "product not found" ||
			errMsg == "unauthorized" ||
			strings.HasPrefix(errMsg, "insufficient stock") {
			status = http.StatusBadRequest
		} else {
			log.Printf("[OrderController] CreateOrder Error: %v", err)
		}
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": res})
}

func GetMyOrders(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	orders, err := orderService.GetMyOrders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": orders})
}

// Subscriptions
func CreateSubscription(c *gin.Context) {
	var req dto.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.MustGet("userID").(uint)
	res, err := orderService.CreateSubscription(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": res})
}

func GetMySubscriptions(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	subs, err := orderService.GetMySubscriptions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": subs})
}
