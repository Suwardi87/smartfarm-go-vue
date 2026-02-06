export interface Address {
  id: number
  user_id: number
  label: string
  recipient_name: string
  phone_number: string
  street: string
  city: string
  province: string
  postal_code: string
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface CreateAddressRequest {
  label: string
  recipient_name: string
  phone_number: string
  street: string
  city: string
  province: string
  postal_code: string
  is_default?: boolean
}

export interface UpdateAddressRequest extends CreateAddressRequest {}
