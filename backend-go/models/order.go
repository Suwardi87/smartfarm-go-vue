package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	TotalPrice   float64  `gorm:"type:decimal(15,2)" json:"total_price"`
	Status       string   `gorm:"type:enum('pending','paid','shipped','completed','cancelled');default:'pending'" json:"status"`
	Type         string   `gorm:"type:enum('regular','preorder');default:'regular'" json:"type"`
	PaymentProof string   `gorm:"type:varchar(255)" json:"payment_proof"`
	AddressID    *uint    `json:"address_id"`
	Address      *Address `gorm:"foreignKey:AddressID" json:"address,omitempty"`
	PaymentID    *uint    `json:"payment_id"`
	Payment      *Payment `gorm:"foreignKey:PaymentID" json:"payment,omitempty"`

	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
}

type OrderItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  int     `json:"quantity"`
	Price     float64 `gorm:"type:decimal(10,2)" json:"price"`
}
