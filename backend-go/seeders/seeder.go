package seeders

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"smartfarm-api/models"
)

func Seed(db *gorm.DB) {
	seedUsers(db)
	seedProducts(db)
	log.Println("✅ Database seeded successfully!")
}

func seedUsers(db *gorm.DB) {
	password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	users := []models.User{
		{Name: "Admin SmartFarm", Email: "admin@smartfarm.com", Password: string(password), Role: "admin"},
		{Name: "Pak Budi (Petani)", Email: "petani@smartfarm.com", Password: string(password), Role: "petani"},
		{Name: "Siti (Pembeli)", Email: "pembeli@smartfarm.com", Password: string(password), Role: "pembeli"},
	}

	for _, user := range users {
		if err := db.FirstOrCreate(&user, models.User{Email: user.Email}).Error; err != nil {
			log.Printf("Failed to seed user %s: %v", user.Email, err)
		}
	}
}

func seedProducts(db *gorm.DB) {
	// Need Farmer ID
	var farmer models.User
	db.Where("email = ?", "petani@smartfarm.com").First(&farmer)

	if farmer.ID == 0 {
		log.Println("Skipping product seeding: Farmer not found")
		return
	}

	futureDate := time.Now().AddDate(0, 1, 0) // 1 month from now

	products := []models.Product{
		// FRESH
		{
			Name: "Bayam Organik Segar", Description: "Bayam hijau segar langsung dari kebun, bebas pestisida.",
			Price: 5000, Stock: 50, Category: "Vegetables", FarmerID: farmer.ID,
			ImageURL: "https://images.unsplash.com/photo-1576045057995-568f588f82fb?auto=format&fit=crop&q=80&w=400",
		},
		{
			Name: "Wortel Brastagi", Description: "Wortel manis dan renyah, cocok untuk jus atau sop.",
			Price: 12000, Stock: 30, Category: "Vegetables", FarmerID: farmer.ID,
			ImageURL: "https://images.unsplash.com/photo-1598170845058-32b9d6a5da37?auto=format&fit=crop&q=80&w=400",
		},
		// PRE-ORDER
		{
			Name: "Melon Golden (Pre-Order)", Description: "Melon manis varietas Golden. Booking sekarang untuk harga lebih murah!",
			Price: 35000, Stock: 20, Category: "Fruits", FarmerID: farmer.ID,
			IsPreOrder: true, HarvestDate: &futureDate,
			ImageURL: "https://images.unsplash.com/photo-1571575173700-afb9492e6a50?auto=format&fit=crop&q=80&w=400",
		},
		// SUBSCRIPTION
		{
			Name: "Paket Sayur Mingguan (Keluarga)", Description: "Paket lengkap (Bayam, Kangkung, Wortel, Tomat) dikirim setiap minggu.",
			Price: 75000, Stock: 100, Category: "Packages", FarmerID: farmer.ID,
			IsSubscription: true, SubscriptionPeriod: "weekly",
			ImageURL: "https://images.unsplash.com/photo-1615485499978-f6952875f566?auto=format&fit=crop&q=80&w=400",
		},
	}

	for _, p := range products {
		if err := db.FirstOrCreate(&p, models.Product{Name: p.Name}).Error; err != nil {
			log.Printf("Failed to seed product %s: %v", p.Name, err)
		}
	}
}
