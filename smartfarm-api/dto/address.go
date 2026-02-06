package dto

type CreateAddressRequest struct {
	Label         string `json:"label" binding:"required"`
	RecipientName string `json:"recipient_name" binding:"required"`
	PhoneNumber   string `json:"phone_number" binding:"required"`
	Street        string `json:"street" binding:"required"`
	City          string `json:"city" binding:"required"`
	Province      string `json:"province" binding:"required"`
	PostalCode    string `json:"postal_code" binding:"required"`
	IsDefault     bool   `json:"is_default"`
}

type UpdateAddressRequest struct {
	Label         string `json:"label" binding:"required"`
	RecipientName string `json:"recipient_name" binding:"required"`
	PhoneNumber   string `json:"phone_number" binding:"required"`
	Street        string `json:"street" binding:"required"`
	City          string `json:"city" binding:"required"`
	Province      string `json:"province" binding:"required"`
	PostalCode    string `json:"postal_code" binding:"required"`
	IsDefault     bool   `json:"is_default"`
}
