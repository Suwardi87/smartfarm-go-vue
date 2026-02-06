package controllers

import (
	"log"
	"net/http"
	"smartfarm-api/dto"
	"smartfarm-api/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreatePayment(c *gin.Context) {
	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(uint)

	payment, token, err := services.CreatePayment(userID, req)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to pay for this order"})
		} else {
			// Pass the specific error message (e.g. from Midtrans or Address check)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"payment_id": payment.ID,
			"snap_token": token,
			"amount":     payment.Amount,
		},
	})
}

func PaymentWebhook(c *gin.Context) {
	var webhookData map[string]interface{}
	if err := c.BindJSON(&webhookData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.ProcessPaymentWebhook(webhookData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}

func GetPaymentStatus(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("order_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	userID := c.MustGet("userID").(uint)

	payment, err := services.GetPaymentByOrderID(userID, uint(orderID))
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"payment_id": payment.ID,
			"status":     payment.Status,
			"amount":     payment.Amount,
		},
	})
}

func MockPaymentSuccess(c *gin.Context) {
	log.Println("[MockPaymentSuccess] Received request")
	var req struct {
		PaymentID uint `json:"payment_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[MockPaymentSuccess] Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[MockPaymentSuccess] Confirming payment for ID: %d", req.PaymentID)
	if err := services.ConfirmMockPayment(req.PaymentID); err != nil {
		log.Printf("[MockPaymentSuccess] Error confirmed mock payment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[MockPaymentSuccess] Success for PaymentID: %d", req.PaymentID)
	c.JSON(http.StatusOK, gin.H{"message": "Mock payment confirmed successfully"})
}
