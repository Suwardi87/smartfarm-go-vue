package repositories

import (
	"smartfarm-api/models"
	"time"

	"gorm.io/gorm"
)

type AnalyticsRepository interface {
	LogView(view *models.ProductView) error
	GetTrendingProducts(limit int) ([]models.Product, error)
}

type analyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) AnalyticsRepository {
	return &analyticsRepository{db}
}

func (r *analyticsRepository) LogView(view *models.ProductView) error {
	return r.db.Create(view).Error
}

func (r *analyticsRepository) GetTrendingProducts(limit int) ([]models.Product, error) {
	var products []models.Product
	
	// Query to count views in last 7 days and join with products
	// Using raw SQL or Gorm chaining
	// Simple approach: Group by product_id
	
	subQuery := r.db.Table("product_views").
		Select("product_id, count(*) as views").
		Where("viewed_at > ?", time.Now().AddDate(0, 0, -7)).
		Group("product_id").
		Order("views desc").
		Limit(limit)

	err := r.db.Table("products").
		Select("products.*, pv.views").
		Joins("JOIN (?) as pv ON pv.product_id = products.id", subQuery).
		Order("pv.views desc").
		Find(&products).Error
		
	return products, err
}
