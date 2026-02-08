package seeders

import (
	"log"
	"smartfarm-api/models"

	"gorm.io/gorm"
)

// CleanOldProducts menghapus produk dengan image URL lokal (bukan Unsplash)
func CleanOldProducts(db *gorm.DB) {
	log.Println("🧹 Membersihkan produk dengan gambar lokal...")

	// Hapus produk yang image_url-nya bukan dari Unsplash
	result := db.Where("image_url NOT LIKE ?", "https://images.unsplash.com%").Delete(&models.Product{})

	if result.Error != nil {
		log.Printf("❌ Gagal menghapus produk lama: %v", result.Error)
		return
	}

	log.Printf("✅ Berhasil menghapus %d produk dengan gambar lokal", result.RowsAffected)
}
