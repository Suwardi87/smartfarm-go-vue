package repositories

import (
	"smartfarm-api/models"
	"time"

	"gorm.io/gorm"
)

type AnalyticsRepository interface {
	LogView(view *models.ProductView) error
	GetTrendingProducts(limit int) ([]models.Product, error)
	GetFarmerStats(farmerID uint) (float64, int, int, int, error)
	GetFarmerRecentOrders(farmerID uint, limit int) ([]models.Order, error)
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

func (r *analyticsRepository) GetFarmerStats(farmerID uint) (float64, int, int, int, error) {
	var totalRevenue float64
	var totalOrders int64
	var totalCustomers int64
	var totalProducts int64

	// Revenue and Orders
	r.db.Table("order_items").
		Joins("join products on products.id = order_items.product_id").
		Joins("join orders on orders.id = order_items.order_id").
		Where("products.farmer_id = ? AND orders.status = ?", farmerID, "paid").
		Select("SUM(order_items.price * order_items.quantity)").
		Scan(&totalRevenue)

	r.db.Table("order_items").
		Joins("join products on products.id = order_items.product_id").
		Joins("join orders on orders.id = order_items.order_id").
		Where("products.farmer_id = ? AND orders.status = ?", farmerID, "paid").
		Distinct("order_items.order_id").
		Count(&totalOrders)

	// Unique Customers
	r.db.Table("order_items").
		Joins("join products on products.id = order_items.product_id").
		Joins("join orders on orders.id = order_items.order_id").
		Where("products.farmer_id = ? AND orders.status = ?", farmerID, "paid").
		Distinct("orders.user_id").
		Count(&totalCustomers)

	// Total Products
	r.db.Model(&models.Product{}).Where("farmer_id = ?", farmerID).Count(&totalProducts)

	return totalRevenue, int(totalOrders), int(totalCustomers), int(totalProducts), nil
}

func (r *analyticsRepository) GetFarmerRecentOrders(farmerID uint, limit int) ([]models.Order, error) {
	var orders []models.Order
	// This is a bit complex because one order might have products from multiple farmers.
	// We only show orders that contain at least one product from this farmer.
	err := r.db.Preload("OrderItems.Product").
		Joins("join order_items on order_items.order_id = orders.id").
		Joins("join products on products.id = order_items.product_id").
		Where("products.farmer_id = ?", farmerID).
		Distinct("orders.id", "orders.*"). // Avoid duplicate orders
		Order("orders.created_at desc").
		Limit(limit).
		Find(&orders).Error

	return orders, err
}
