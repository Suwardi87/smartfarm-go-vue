import http from "@/lib/http"
import type { Address, CreateAddressRequest, UpdateAddressRequest } from "@/dto/address/Address"
import { normalizeAddress } from "@/utils/transformers"

export function getMyAddresses() {
  return http.get<{ data: Address[] }>("/addresses")
}

export function createAddress(payload: CreateAddressRequest) {
  return http.post<{ data: Address }>("/addresses", payload)
}

export function updateAddress(addressId: number, payload: UpdateAddressRequest) {
  return http.put<{ data: Address }>(`/addresses/${addressId}`, payload)
}

export function deleteAddress(addressId: number) {
  return http.delete(`/addresses/${addressId}`)
}

export function setDefaultAddress(addressId: number) {
  return http.post(`/addresses/${addressId}/default`)
}
// Class-based service for Checkout page compatibility
class AddressService {
  async getMyAddresses() {
    const { data } = await http.get<{ data: Address[] }>("/addresses")
    return (data.data || []).map(normalizeAddress) as Address[]
  }

  async createAddress(payload: CreateAddressRequest) {
    const { data } = await http.post<{ data: Address }>("/addresses", payload)
    return normalizeAddress(data.data) as Address
  }

  async updateAddress(addressId: number, payload: UpdateAddressRequest) {
    const { data } = await http.put<{ data: Address }>(`/addresses/${addressId}`, payload)
    return normalizeAddress(data.data) as Address
  }

  async deleteAddress(addressId: number) {
    return http.delete(`/addresses/${addressId}`)
  }

  async setDefaultAddress(addressId: number) {
    return http.post(`/addresses/${addressId}/default`)
  }
}

export default new AddressService()
