import http from "@/lib/http"
import type { Product } from "@/dto/product/Product"
import type { AxiosResponse } from 'axios'

interface ApiResponse<T> {
  data: T
}

export function getProducts(): Promise<AxiosResponse<ApiResponse<Product[]>>> {
  return http.get('/products')
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

export function getFarmerProducts(): Promise<AxiosResponse<ApiResponse<Product[]>>> {
  return http.get('/farmer/products')
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
