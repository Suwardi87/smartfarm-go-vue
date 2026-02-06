package repositories

import (
	"smartfarm-api/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order) error
	FindByID(id uint) (models.Order, error)
	FindByUserID(userID uint) ([]models.Order, error)
	FindAll() ([]models.Order, error)
	UpdateStatus(id uint, status string) error
	Update(order *models.Order) error
	UpdatePaymentInfo(id uint, paymentID uint, addressID uint) error
	CreateSubscription(sub *models.Subscription) error
	FindSubscriptionsByUserID(userID uint) ([]models.Subscription, error)
	WithTx(tx *gorm.DB) OrderRepository
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) FindByID(id uint) (models.Order, error) {
	var order models.Order
	err := r.db.Preload("OrderItems.Product").First(&order, id).Error
	return order, err
}

func (r *orderRepository) FindByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("OrderItems.Product").Where("user_id = ?", userID).Order("created_at desc").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) FindAll() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("User").Preload("OrderItems.Product").Order("created_at desc").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *orderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) UpdatePaymentInfo(id uint, paymentID uint, addressID uint) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).Updates(map[string]interface{}{
		"payment_id": paymentID,
		"address_id": addressID,
	}).Error
}

func (r *orderRepository) CreateSubscription(sub *models.Subscription) error {
	return r.db.Create(sub).Error
}

func (r *orderRepository) FindSubscriptionsByUserID(userID uint) ([]models.Subscription, error) {
	var subs []models.Subscription
	err := r.db.Preload("Product").Where("user_id = ?", userID).Order("created_at desc").Find(&subs).Error
	return subs, err
}

func (r *orderRepository) WithTx(tx *gorm.DB) OrderRepository {
	return &orderRepository{db: tx}
}
