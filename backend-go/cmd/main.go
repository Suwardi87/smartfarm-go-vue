package main

import (
	"log"

	"os"
	"smartfarm-api/config"
	"smartfarm-api/controllers"
	"smartfarm-api/routes"
	"smartfarm-api/seeders"
	"smartfarm-api/services"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env not found, using system env")
	}

	config.ConnectDatabase()

	// Check for seed command
	if len(os.Args) > 1 {
		if os.Args[1] == "seed" {
			seeders.Seed(config.DB)
			return
		}
		if os.Args[1] == "seed-bulk" {
			seeders.CleanOldProducts(config.DB) // 🧹 Clean old products first
			seeders.SeedBulk(config.DB)
			return
		}
		if os.Args[1] == "seed-bulk-concurrent" {
			seeders.CleanOldProducts(config.DB)   // 🧹 Clean old products first
			seeders.SeedBulkConcurrent(config.DB) // 🚀 Concurrent seeding!
			return
		}
		if os.Args[1] == "cleanup" {
			seeders.CleanOldProducts(config.DB)
			return
		}
	}

	// Init services/controllers
	controllers.InitProductController()
	controllers.InitOrderController()
	controllers.InitAnalyticsController()
	services.InitPaymentService()

	r := routes.SetupRoutes()
	r.Run(":8080")
}
