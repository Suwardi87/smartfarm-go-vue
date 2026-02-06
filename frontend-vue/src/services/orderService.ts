import http from "@/lib/http"

export interface OrderItemRequest {
    product_id: number
    quantity: number
}

export interface CreateOrderRequest {
    items: OrderItemRequest[]
    address_id?: number
}

export function createOrder(payload: CreateOrderRequest) {
    return http.post("/orders", payload)
}

export function getMyOrders() {
    return http.get("/orders")
}

export interface CreateSubscriptionRequest {
    product_id: number
    frequency: string // "weekly", "monthly"
    duration: number
}

export function createSubscription(payload: CreateSubscriptionRequest) {
    return http.post("/subscriptions", payload)
}

export function getMySubscriptions() {
    return http.get("/subscriptions")
}
// Class-based service for Checkout page compatibility
class OrderService {
    async createOrder(payload: CreateOrderRequest) {
        const { data } = await http.post("/orders", payload)
        return data.data
    }

    async getMyOrders() {
        const { data } = await http.get("/orders")
        return data.data
    }

    async getMySubscriptions() {
        const { data } = await http.get("/subscriptions")
        return data.data
    }
}

export default new OrderService()
