package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string  `gorm:"type:varchar(255);not null;index" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Price       float64 `gorm:"type:decimal(10,2);not null;index" json:"price"`
	Stock       int     `gorm:"not null" json:"stock"`
	ImageURL    string  `gorm:"type:varchar(255)" json:"image_url"`
	Category    string  `gorm:"type:varchar(100);index" json:"category"`

	// Relation to Farmer (User)
	FarmerID uint `json:"farmer_id"`
	Farmer   User `gorm:"foreignKey:FarmerID" json:"farmer,omitempty"`

	// Features
	IsPreOrder  bool       `json:"is_pre_order"`
	HarvestDate *time.Time `json:"harvest_date"` // Nullable, only for pre-order

	IsSubscription     bool   `json:"is_subscription"`
	SubscriptionPeriod string `json:"subscription_period"` // "weekly", "monthly"

	Views int `gorm:"-" json:"views"` // Transient field for analytics
}
