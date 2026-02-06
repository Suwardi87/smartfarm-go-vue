export interface CreatePaymentRequest {
  order_id: number
  address_id: number
  amount: number
}

export interface PaymentResponse {
  payment_id: number
  snap_token: string
  amount: number
}

export interface PaymentStatus {
  payment_id: number
  status: string
  amount: number
}
