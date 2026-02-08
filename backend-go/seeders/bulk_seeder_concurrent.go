package seeders

import (
	"fmt"
	"log"
	"math/rand"
	"smartfarm-api/models"
	"sync"
	"time"

	"gorm.io/gorm"
)

// SeedBulkConcurrent seeds 100,000 products using goroutines for faster performance
func SeedBulkConcurrent(db *gorm.DB) {
	log.Println("🚀 Memulai CONCURRENT Seeding 100.000 data...")
	start := time.Now()

	// Find or create farmer
	var farmer models.User
	if err := db.Where("role = ?", "petani").First(&farmer).Error; err != nil {
		farmer = models.User{
			Name:     "Pak Budi (Petani)",
			Email:    "budi@petani.com",
			Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
			Role:     "petani",
		}
		db.Create(&farmer)
	}

	categories := []string{"Vegetables", "Fruits", "Packages", "Herbs", "Hydroponics"}
	productNames := []string{"Tomat", "Bayam", "Kangkung", "Wortel", "Sawi", "Melon", "Semangka", "Cabai", "Bawang", "Selada"}
	adjectives := []string{"Segar", "Organik", "Super", "Pilihan", "Kebun", "Hidroponik", "Premium", "Manis", "Renyah", "Lokal"}

	totalProducts := 100000
	batchSize := 5000
	numBatches := totalProducts / batchSize

	var wg sync.WaitGroup
	var mu sync.Mutex
	successCount := 0

	// Process batches concurrently
	for i := 0; i < numBatches; i++ {
		wg.Add(1)

		go func(batchIndex int) {
			defer wg.Done()

			offset := batchIndex * batchSize
			var products []models.Product

			// Generate products for this batch
			for j := 0; j < batchSize; j++ {
				name := fmt.Sprintf("%s %s #%d",
					productNames[rand.Intn(len(productNames))],
					adjectives[rand.Intn(len(adjectives))],
					offset+j)

				products = append(products, models.Product{
					Name:           name,
					Description:    fmt.Sprintf("Deskripsi produk berkualitas %s yang dihasilkan langsung dari kebun kami.", name),
					Price:          float64((rand.Intn(20) + 1) * 5000),
					Stock:          rand.Intn(100) + 10,
					Category:       categories[rand.Intn(len(categories))],
					FarmerID:       farmer.ID,
					ImageURL:       "https://images.unsplash.com/photo-1615485499978-f6952875f566?auto=format&fit=crop&q=80&w=400",
					IsPreOrder:     rand.Float32() < 0.1,
					IsSubscription: rand.Float32() < 0.05,
					CreatedAt:      time.Now().Add(-time.Duration(rand.Intn(60)) * time.Hour * 24),
				})
			}

			// Insert batch to database
			if err := db.Create(&products).Error; err != nil {
				log.Printf("❌ Batch %d gagal: %v", batchIndex+1, err)
				return
			}

			// Update success counter (thread-safe)
			mu.Lock()
			successCount++
			log.Printf("✅ Batch %d/%d selesai (%d produk)", successCount, numBatches, len(products))
			mu.Unlock()
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("✅ SELESAI! 100.000 data berhasil dibuat dalam waktu %v", elapsed)
	log.Printf("📊 Performance: %.0f produk/detik", float64(totalProducts)/elapsed.Seconds())
}
