package repositories

import (
	"smartfarm-api/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	Update(product *models.Product) error
	Delete(id uint) error
	FindAll() ([]models.Product, error)
	FindByID(id uint) (models.Product, error)
	FindByFarmerID(farmerID uint) ([]models.Product, error)
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

func (r *productRepository) FindAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Farmer").Find(&products).Error
	return products, err
}

func (r *productRepository) FindByID(id uint) (models.Product, error) {
	var product models.Product
	err := r.db.Preload("Farmer").First(&product, id).Error
	return product, err
}

func (r *productRepository) FindByFarmerID(farmerID uint) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Farmer").Where("farmer_id = ?", farmerID).Find(&products).Error
	return products, err
}

func (r *productRepository) WithTx(tx *gorm.DB) ProductRepository {
	return &productRepository{db: tx}
}
