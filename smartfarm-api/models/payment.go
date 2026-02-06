package models

import "time"

type Payment struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	OrderID       uint      `gorm:"index" json:"order_id"`
	Order         *Order    `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	UserID        uint      `gorm:"index" json:"user_id"`
	Amount        float64   `json:"amount"`
	Status        string    `gorm:"type:enum('pending','success','failed','expired')" json:"status"`
	PaymentMethod string    `gorm:"type:varchar(50)" json:"payment_method"` // credit_card, bank_transfer, ewallet, etc
	TransactionID string    `gorm:"unique;type:varchar(255)" json:"transaction_id"`
	SnapToken     string    `gorm:"type:text" json:"snap_token"` // Midtrans snap token
	SnapURL       string    `gorm:"type:text" json:"snap_url"`   // Midtrans snap URL
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Payment) TableName() string {
	return "payments"
}
