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
	GetFarmerDashboardData(farmerID uint) (dto.FarmerDashboardResponse, error)
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

func (s *analyticsService) GetFarmerDashboardData(farmerID uint) (dto.FarmerDashboardResponse, error) {
	rev, orders, customers, products, err := s.repo.GetFarmerStats(farmerID)
	if err != nil {
		return dto.FarmerDashboardResponse{}, err
	}

	recentOrders, err := s.repo.GetFarmerRecentOrders(farmerID, 5)
	if err != nil {
		return dto.FarmerDashboardResponse{}, err
	}

	var recentOrdersResponse []dto.FarmerRecentOrder
	for _, o := range recentOrders {
		// Find the product that belongs to this farmer in the order
		for _, item := range o.OrderItems {
			if item.Product.FarmerID == farmerID {
				recentOrdersResponse = append(recentOrdersResponse, dto.FarmerRecentOrder{
					ID:          o.ID,
					ProductName: item.Product.Name,
					Category:    item.Product.Category,
					Price:       item.Price,
					Status:      o.Status,
					ImageURL:    item.Product.ImageURL,
				})
				break // Only show one product per order for simplicity in "Recent Orders"
			}
		}
	}

	return dto.FarmerDashboardResponse{
		Stats: dto.FarmerStatsResponse{
			TotalRevenue:    rev,
			RevenueGrowth:   12.5, // Dummy for UI
			TotalOrders:     orders,
			OrdersGrowth:    -5.2, // Dummy for UI
			TotalCustomers:  customers,
			CustomersGrowth: 8.1, // Dummy for UI
			TotalProducts:   products,
			ProductsGrowth:  2.0, // Dummy for UI
		},
		RecentOrders: recentOrdersResponse,
	}, nil
}
