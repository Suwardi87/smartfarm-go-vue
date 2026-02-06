package dto

type CreatePaymentRequest struct {
	OrderID   uint    `json:"order_id" binding:"required"`
	AddressID uint    `json:"address_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required"`
}

type PaymentWebhookRequest struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionID     string `json:"transaction_id"`
	StatusCode        string `json:"status_code"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
}
