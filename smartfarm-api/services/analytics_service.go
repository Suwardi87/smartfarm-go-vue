package services

import (
	"smartfarm-api/dto"
	"smartfarm-api/models"
	"smartfarm-api/repositories"
	"time"
)

type AnalyticsService interface {
	LogView(productID uint, userID uint) error
	GetTrendingProducts() ([]dto.ProductResponse, error)
}

type analyticsService struct {
	repo repositories.AnalyticsRepository
}

func NewAnalyticsService(repo repositories.AnalyticsRepository) AnalyticsService {
	return &analyticsService{repo}
}

func (s *analyticsService) LogView(productID uint, userID uint) error {
	view := models.ProductView{
		ProductID: productID,
		UserID:    userID,
		ViewedAt:  time.Now(),
	}
	return s.repo.LogView(&view)
}

func (s *analyticsService) GetTrendingProducts() ([]dto.ProductResponse, error) {
	products, err := s.repo.GetTrendingProducts(5) // Top 5
	if err != nil {
		return nil, err
	}

	var responses []dto.ProductResponse
	// We need to verify if product service map function is reusable or duplicate it
	// For simplicity, I'll duplicate the mapping logic or make it shared.
	// I'll duplicate for now to avoid complexity of shared package dependency cycle (though utils is fine).

	for _, p := range products {
		var harvestDateStr string
		if p.HarvestDate != nil {
			harvestDateStr = p.HarvestDate.Format("2006-01-02")
		}

		responses = append(responses, dto.ProductResponse{
			ID:                 p.ID,
			Name:               p.Name,
			Description:        p.Description,
			Price:              p.Price,
			Stock:              p.Stock,
			ImageURL:           p.ImageURL,
			Category:           p.Category,
			FarmerID:           p.FarmerID,
			IsPreOrder:         p.IsPreOrder,
			HarvestDate:        harvestDateStr,
			IsSubscription:     p.IsSubscription,
			SubscriptionPeriod: p.SubscriptionPeriod,
			Views:              p.Views,
		})
	}

	return responses, nil
}
