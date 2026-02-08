import http from "@/lib/http"
import type { Product } from "@/dto/product/Product"
import type { AxiosResponse } from 'axios'

export interface ApiResponse<T> {
  data: T
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
  total_pages: number
}

export function getProducts(page: number = 1, limit: number = 12, q: string = ''): Promise<AxiosResponse<PaginatedResponse<Product>>> {
  return http.get('/products', {
    params: { page, limit, q }
  })
}

export function getProduct(id: number): Promise<AxiosResponse<ApiResponse<Product>>> {
  return http.get(`/products/${id}`)
}

export function createProduct(formData: FormData) {
  return http.post("/products", formData, {
    headers: {
      "Content-Type": "multipart/form-data"
    }
  })
}

export function getFarmerProducts(page: number = 1, limit: number = 10): Promise<AxiosResponse<PaginatedResponse<Product>>> {
  return http.get('/farmer/products', {
    params: { page, limit }
  })
}

export function updateProduct(id: number, formData: FormData) {
  return http.put(`/products/${id}`, formData, {
    headers: {
      "Content-Type": "multipart/form-data"
    }
  })
}

export function deleteProduct(id: number) {
  return http.delete(`/products/${id}`)
}
