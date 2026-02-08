import http from "@/lib/http"
import type { AxiosResponse } from "axios"
interface ApiResponse<T> {
  data: T
}

export interface CommodityTrend {
  id: number
  name: string
  price: number
  views?: number
}

export interface FarmerStats {
  total_revenue: number
  revenue_growth: number
  total_orders: number
  orders_growth: number
  total_customers: number
  customers_growth: number
  total_products: number
  products_growth: number
}

export interface FarmerRecentOrder {
  id: number
  product_name: string
  category: string
  price: number
  status: string
  image_url: string
}

export interface FarmerDashboardData {
  stats: FarmerStats
  recent_orders: FarmerRecentOrder[]
}

export function getRecommendations(): Promise<AxiosResponse<ApiResponse<CommodityTrend[]>>> {
  return http.get('/analytics/trending')
}

export function getFarmerDashboardData(): Promise<AxiosResponse<ApiResponse<FarmerDashboardData>>> {
  return http.get('/analytics/farmer')
}
