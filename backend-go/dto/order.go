package dto

type CreateOrderRequest struct {
	Items     []OrderItemRequest `json:"items" binding:"required"`
	AddressID uint               `json:"address_id"`
}

type OrderItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

type OrderResponse struct {
	ID           uint                `json:"id"`
	UserID       uint                `json:"user_id"`
	TotalPrice   float64             `json:"total_price"`
	Status       string              `json:"status"`
	Type         string              `json:"type"`
	PaymentProof string              `json:"payment_proof"`
	Items        []OrderItemResponse `json:"items"`
	CreatedAt    string              `json:"created_at"`
}

type OrderItemResponse struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	SubTotal    float64 `json:"sub_total"`
}

// Subscription DTOs
type CreateSubscriptionRequest struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Frequency string `json:"frequency" binding:"required,oneof=weekly monthly"`
	Duration  int    `json:"duration" binding:"required,min=1"` // Number of periods (e.g., 4 weeks)
}

type SubscriptionResponse struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	ProductName string `json:"product_name"`
	Frequency   string `json:"frequency"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Status      string `json:"status"`
}
