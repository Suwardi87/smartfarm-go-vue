package controllers

import (
	"net/http"
	"smartfarm-api/config"
	"smartfarm-api/dto"
	"smartfarm-api/repositories"
	"smartfarm-api/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

var productService services.ProductService

func init() {
	// Note: In a real app, use Dependency Injection or a proper setup function.
	// This init is a temporary shortcut or depends on config.DB being ready.
	// But config.DB is init in main. Better to setup in SetupRoutes or main.
}

// Helper to init service manually if not using DI framework
func InitProductController() {
	db := config.DB
	repo := repositories.NewProductRepository(db)
	productService = services.NewProductService(repo)
}

func CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Assuming Auth middleware sets "userID"
	// For testing without auth, we might hardcode or check header
	// userID := c.GetUint("userID")
	// For now, let's hardcode active user or get from context if available

	// Simulate user ID 1 (Farmer) for now if middleware not ready
	userID := uint(1)

	res, err := productService.CreateProduct(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": res})
}

func GetAllProducts(c *gin.Context) {
	products, err := productService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": products})
}

func GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	product, err := productService.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Log View (Async to not block response)
	// Mock UserID 1 or 0 (Guest)
	userID := uint(1)
	go func() {
		if analyticsService != nil {
			analyticsService.LogView(product.ID, userID)
		}
	}()

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func GetFarmerProducts(c *gin.Context) {
	// Ambil userID dari context (diatur oleh AuthMiddleware)
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := val.(uint)
	if !ok {
		// handle float64 if coming from JWT claims directly in some middleware setups
		if fID, ok := val.(float64); ok {
			userID = uint(fID)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
			return
		}
	}

	products, err := productService.FindProductsByFarmerID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.CreateProductRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	val, _ := c.Get("userID")
	userID := val.(uint) // Assuming AuthMiddleware is used

	res, err := productService.UpdateProduct(uint(id), req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	val, _ := c.Get("userID")
	userID := val.(uint)

	err = productService.DeleteProduct(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
