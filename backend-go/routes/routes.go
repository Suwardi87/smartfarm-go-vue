package routes

import (
	"log"
	"smartfarm-api/controllers"
	"smartfarm-api/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// Define your routes here
	// Public Routes
	r.POST("/signup", controllers.Register)
	r.POST("/signin", controllers.Login)
	r.POST("/logout", controllers.Logout)
	r.GET("/products", controllers.GetAllProducts)
	r.GET("/products/:id", controllers.GetProductByID)
	r.POST("/payments/webhook", controllers.PaymentWebhook)

	// Static for images
	r.Static("/uploads", "./uploads")

	// Protected Routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/me", controllers.Me)
		protected.PUT("/me", controllers.UpdateProfile)

		// Address Routes
		protected.POST("/addresses", controllers.CreateAddress)
		protected.GET("/addresses", controllers.GetMyAddresses)
		protected.PUT("/addresses/:id", controllers.UpdateAddress)
		protected.DELETE("/addresses/:id", controllers.DeleteAddress)
		protected.POST("/addresses/:id/default", controllers.SetDefaultAddress)

		// Product Routes
		protected.POST("/products", controllers.CreateProduct)
		protected.GET("/farmer/products", controllers.GetFarmerProducts)
		protected.PUT("/products/:id", controllers.UpdateProduct)
		protected.DELETE("/products/:id", controllers.DeleteProduct)

		// Order Routes
		protected.POST("/orders", controllers.CreateOrder)
		protected.GET("/orders", controllers.GetMyOrders)

		// Subscription Routes
		protected.POST("/subscriptions", controllers.CreateSubscription)
		protected.GET("/subscriptions", controllers.GetMySubscriptions)

		// Analytics Routes
		protected.GET("/analytics/trending", controllers.GetTrendingProducts)
		protected.GET("/analytics/farmer", controllers.GetFarmerDashboard)

		// Payment Routes
		protected.POST("/payments", controllers.CreatePayment)
		protected.POST("/payments/mock-success", controllers.MockPaymentSuccess)
		protected.GET("/payments/orders/:order_id", controllers.GetPaymentStatus)

	}

	// Log all routes
	for _, route := range r.Routes() {
		log.Printf("[Route] %s %s", route.Method, route.Path)
	}

	return r
}
