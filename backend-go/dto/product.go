package dto

import "mime/multipart"

type CreateProductRequest struct {
	Name               string                `form:"name" binding:"required"`
	Description        string                `form:"description"`
	Price              float64               `form:"price" binding:"required"`
	Stock              int                   `form:"stock" binding:"required"`
	Image              *multipart.FileHeader `form:"image"`
	Category           string                `form:"category"`
	IsPreOrder         bool                  `form:"is_pre_order"`
	HarvestDate        string                `form:"harvest_date"` // YYYY-MM-DD
	IsSubscription     bool                  `form:"is_subscription"`
	SubscriptionPeriod string                `form:"subscription_period"`
}

type UpdateProductRequest struct {
	Name               string                `form:"name"`
	Description        string                `form:"description"`
	Price              float64               `form:"price"`
	Stock              int                   `form:"stock"`
	Image              *multipart.FileHeader `form:"image"`
	Category           string                `form:"category"`
	IsPreOrder         bool                  `form:"is_pre_order"`
	HarvestDate        string                `form:"harvest_date"`
	IsSubscription     bool                  `form:"is_subscription"`
	SubscriptionPeriod string                `form:"subscription_period"`
}

type ProductResponse struct {
	ID                 uint    `json:"id"`
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	Price              float64 `json:"price"`
	Stock              int     `json:"stock"`
	ImageURL           string  `json:"image_url"`
	Category           string  `json:"category"`
	FarmerID           uint    `json:"farmer_id"`
	FarmerName         string  `json:"farmer_name"`
	IsPreOrder         bool    `json:"is_pre_order"`
	HarvestDate        string  `json:"harvest_date,omitempty"`
	IsSubscription     bool    `json:"is_subscription"`
	SubscriptionPeriod string  `json:"subscription_period,omitempty"`
	Views              int     `json:"views,omitempty"`
}

type PaginatedProductResponse struct {
	Data       []ProductResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}
