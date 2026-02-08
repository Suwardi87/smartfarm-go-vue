package repositories

import (
	"log"
	"smartfarm-api/models"
	"time"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	Update(product *models.Product) error
	Delete(id uint) error
	FindAll(query string, limit int, offset int) ([]models.Product, error)
	CountAll(query string) (int64, error)
	FindByID(id uint) (models.Product, error)
	FindByFarmerID(farmerID uint, limit int, offset int) ([]models.Product, error)
	CountByFarmerID(farmerID uint) (int64, error)
	WithTx(tx *gorm.DB) ProductRepository
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *productRepository) FindAll(query string, limit int, offset int) ([]models.Product, error) {
	start := time.Now()
	var products []models.Product
	db := r.db.Preload("Farmer")
	if query != "" {
		q := "%" + query + "%"
		db = db.Where("name LIKE ? OR description LIKE ? OR category LIKE ?", q, q, q)
	}
	err := db.Limit(limit).Offset(offset).Find(&products).Error
	log.Printf("[DEBUG] repo.FindAll DB execution took %v", time.Since(start))
	return products, err
}

func (r *productRepository) CountAll(query string) (int64, error) {
	start := time.Now()
	var count int64
	db := r.db.Model(&models.Product{})
	if query != "" {
		q := "%" + query + "%"
		db = db.Where("name LIKE ? OR description LIKE ? OR category LIKE ?", q, q, q)
	}
	err := db.Count(&count).Error
	log.Printf("[DEBUG] repo.CountAll DB execution took %v", time.Since(start))
	return count, err
}

func (r *productRepository) FindByID(id uint) (models.Product, error) {
	var product models.Product
	err := r.db.Preload("Farmer").First(&product, id).Error
	return product, err
}

func (r *productRepository) FindByFarmerID(farmerID uint, limit int, offset int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Farmer").Where("farmer_id = ?", farmerID).Limit(limit).Offset(offset).Find(&products).Error
	return products, err
}

func (r *productRepository) CountByFarmerID(farmerID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Product{}).Where("farmer_id = ?", farmerID).Count(&count).Error
	return count, err
}

func (r *productRepository) WithTx(tx *gorm.DB) ProductRepository {
	return &productRepository{db: tx}
}
