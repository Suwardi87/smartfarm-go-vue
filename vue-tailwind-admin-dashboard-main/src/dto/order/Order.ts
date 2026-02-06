export interface OrderItem {
  id: number
  order_id: number
  product_id: number
  product_name: string
  quantity: number
  price: number
  product?: {
    id: number
    name: string
    image_url: string
    price: number
  }
}

export interface Order {
  id: number
  user_id: number
  status: string // "pending", "paid", "shipped", "delivered"
  total_price: number
  created_at: string
  updated_at: string
  items?: OrderItem[]
}

export interface Subscription {
  id: number
  user_id: number
  product_id: number
  frequency: string // "weekly", "monthly"
  status: string // "active", "paused", "cancelled"
  created_at: string
  updated_at: string
  product?: {
    id: number
    name: string
    image_url: string
    price: number
  }
}
