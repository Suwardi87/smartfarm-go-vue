package repositories

import (
	"smartfarm-api/config"
	"smartfarm-api/models"
)

func CreatePayment(payment *models.Payment) error {
	return config.DB.Create(payment).Error
}

func GetPaymentByID(id uint) (*models.Payment, error) {
	var payment models.Payment
	err := config.DB.First(&payment, id).Error
	return &payment, err
}

func GetPaymentByTransactionID(transactionID string) (*models.Payment, error) {
	var payment models.Payment
	err := config.DB.Where("transaction_id = ?", transactionID).First(&payment).Error
	return &payment, err
}

func GetPaymentByOrderID(orderID uint) (*models.Payment, error) {
	var payment models.Payment
	err := config.DB.Where("order_id = ?", orderID).First(&payment).Error
	return &payment, err
}

func UpdatePayment(payment *models.Payment) error {
	return config.DB.Save(payment).Error
}
