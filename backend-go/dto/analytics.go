package dto

type FarmerStatsResponse struct {
	TotalRevenue    float64 `json:"total_revenue"`
	RevenueGrowth   float64 `json:"revenue_growth"` // Percentage
	TotalOrders     int     `json:"total_orders"`
	OrdersGrowth    float64 `json:"orders_growth"` // Percentage
	TotalCustomers  int     `json:"total_customers"`
	CustomersGrowth float64 `json:"customers_growth"` // Percentage
	TotalProducts   int     `json:"total_products"`
	ProductsGrowth  float64 `json:"products_growth"` // Percentage
}

type FarmerRecentOrder struct {
	ID          uint    `json:"id"`
	ProductName string  `json:"product_name"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"`
	ImageURL    string  `json:"image_url"`
}

type FarmerDashboardResponse struct {
	Stats        FarmerStatsResponse `json:"stats"`
	RecentOrders []FarmerRecentOrder `json:"recent_orders"`
}
