package seeders

import (
	"fmt"
	"log"
	"math/rand"
	"smartfarm-api/models"
	"time"

	"gorm.io/gorm"
)

func SeedBulk(db *gorm.DB) {
	start := time.Now()
	log.Println("🚀 Memulai proses Seeding 100.000 data...")

	// 1. Ambil Farmer ID
	var farmer models.User
	db.Where("role = ?", "petani").First(&farmer)
	if farmer.ID == 0 {
		log.Println("❌ Gagal: Tidak ada petani terdaftar untuk seeding produk.")
		return
	}

	categories := []string{"Vegetables", "Fruits", "Packages", "Herbs", "Hydroponics"}
	productNames := []string{"Tomat", "Bayam", "Kangkung", "Wortel", "Sawi", "Melon", "Semangka", "Cabai", "Bawang", "Selada"}
	adjectives := []string{"Segar", "Organik", "Super", "Pilihan", "Kebun", "Hidroponik", "Premium", "Manis", "Renyah", "Lokal"}

	totalProducts := 100000
	batchSize := 5000

	for i := 0; i < totalProducts; i += batchSize {
		var products []models.Product
		currentBatchSize := batchSize
		if i+batchSize > totalProducts {
			currentBatchSize = totalProducts - i
		}

		for j := 0; j < currentBatchSize; j++ {
			name := fmt.Sprintf("%s %s #%d", productNames[rand.Intn(len(productNames))], adjectives[rand.Intn(len(adjectives))], i+j)
			products = append(products, models.Product{
				Name:           name,
				Description:    fmt.Sprintf("Deskripsi produk berkualitas %s yang dihasilkan langsung dari kebun kami.", name),
				Price:          float64((rand.Intn(20) + 1) * 5000),
				Stock:          rand.Intn(100) + 10,
				Category:       categories[rand.Intn(len(categories))],
				FarmerID:       farmer.ID,
				ImageURL:       "https://images.unsplash.com/photo-1615485499978-f6952875f566?auto=format&fit=crop&q=80&w=400",
				IsPreOrder:     rand.Float32() < 0.1,  // 10% pre-order
				IsSubscription: rand.Float32() < 0.05, // 5% sub
				CreatedAt:      time.Now().Add(-time.Duration(rand.Intn(60)) * time.Hour * 24),
			})
		}

		err := db.CreateInBatches(products, 1000).Error
		if err != nil {
			log.Printf("❌ Gagal meyimpan batch %d : %v", i/batchSize, err)
		} else {
			log.Printf("📦 Berhasil menyimpan %d dari %d data...", i+currentBatchSize, totalProducts)
		}
	}

	elapsed := time.Since(start)
	log.Printf("✅ SELESAI! 100.000 data berhasil dibuat dalam waktu %s", elapsed)
}
