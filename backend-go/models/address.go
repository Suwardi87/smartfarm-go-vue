package models

import "time"

type Address struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index" json:"user_id"`
	Label         string    `gorm:"type:varchar(50)" json:"label"` // "Rumah", "Kantor", dll
	RecipientName string    `gorm:"type:varchar(100)" json:"recipient_name"`
	PhoneNumber   string    `gorm:"type:varchar(20)" json:"phone_number"`
	Street        string    `gorm:"type:text" json:"street"`
	City          string    `gorm:"type:varchar(100)" json:"city"`
	Province      string    `gorm:"type:varchar(100)" json:"province"`
	PostalCode    string    `gorm:"type:varchar(20)" json:"postal_code"`
	IsDefault     bool      `gorm:"default:false" json:"is_default"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Address) TableName() string {
	return "addresses"
}
