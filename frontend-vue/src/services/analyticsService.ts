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



export function getRecommendations(): Promise<AxiosResponse<ApiResponse<CommodityTrend[]>>> {
  return http.get('/analytics/trending')
}
