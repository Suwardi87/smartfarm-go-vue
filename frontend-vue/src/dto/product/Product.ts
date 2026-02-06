export interface User {
    ID: number
    Name: string
    Email: string
}

export interface Product {
  id: number
  name: string
  description: string
  price: number
  stock: number
  image_url: string
  category: string

  farmer_id: number
  farmer_name: string

  is_pre_order: boolean
  harvest_date?: string

  is_subscription: boolean
  subscription_period?: string
}

