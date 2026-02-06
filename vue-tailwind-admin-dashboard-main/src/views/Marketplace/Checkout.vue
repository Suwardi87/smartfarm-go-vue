<template>
  <div class="min-h-screen bg-gray-50 py-8">
    <div class="max-w-6xl mx-auto px-4">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900">Checkout</h1>
        <p class="text-gray-600 mt-2">Review your order and choose delivery address</p>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Main Content -->
        <div class="lg:col-span-2">
          <!-- Order Summary -->
          <div class="bg-white rounded-lg shadow-sm p-6 mb-6">
            <h2 class="text-xl font-semibold text-gray-900 mb-4">Order Summary</h2>
            <div class="space-y-4">
              <div
                v-for="item in cart.state.items"
                :key="item.product.id"
                class="flex items-center justify-between py-4 border-b last:border-b-0"
              >
                <div class="flex items-center space-x-4">
                  <img
                    :src="getImageUrl(item.product.image_url)"
                    :alt="item.product.name"
                    class="w-16 h-16 object-cover rounded"
                  />
                  <div>
                    <h3 class="font-semibold text-gray-900">{{ item.product.name }}</h3>
                    <p class="text-sm text-gray-600">Qty: {{ item.quantity }}</p>
                  </div>
                </div>
                <div class="text-right">
                  <p class="font-semibold text-gray-900">
                    Rp{{ (item.product.price * item.quantity).toLocaleString('id-ID') }}
                  </p>
                  <p class="text-sm text-gray-600">Rp{{ item.product.price.toLocaleString('id-ID') }}/unit</p>
                </div>
              </div>
            </div>
          </div>

          <!-- Select Delivery Address -->
          <div class="bg-white rounded-lg shadow-sm p-6 mb-6">
            <div class="flex items-center justify-between mb-4">
              <h2 class="text-xl font-semibold text-gray-900">Delivery Address</h2>
              <button
                @click="goToAddresses"
                class="text-brand-600 hover:text-brand-700 text-sm font-medium"
              >
                Manage Addresses
              </button>
            </div>

            <div v-if="addresses.length === 0" class="text-center py-8">
              <p class="text-gray-600 mb-4">No addresses yet</p>
              <button
                @click="goToAddresses"
                class="inline-flex items-center px-4 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 transition"
              >
                <PlusIcon class="w-5 h-5 mr-2" />
                Add Address
              </button>
            </div>

            <div v-else class="space-y-3">
              <div
                v-for="address in addresses"
                :key="address.id"
                @click="selectedAddressId = address.id"
                class="relative cursor-pointer"
              >
                <div
                  class="p-4 border rounded-lg transition"
                  :class="
                    selectedAddressId === address.id
                      ? 'border-brand-600 bg-brand-50'
                      : 'border-gray-200 hover:border-gray-300'
                  "
                >
                  <div class="flex items-start justify-between">
                    <div class="flex-1">
                      <div class="flex items-center space-x-2">
                        <h3 class="font-semibold text-gray-900">{{ address.label }}</h3>
                        <span
                          v-if="address.is_default"
                          class="inline-block px-2 py-1 bg-green-100 text-green-800 text-xs font-medium rounded"
                        >
                          Default
                        </span>
                      </div>
                      <p class="text-sm text-gray-600 mt-1">{{ address.recipient_name }}</p>
                      <p class="text-sm text-gray-600">{{ address.phone_number }}</p>
                      <p class="text-sm text-gray-600 mt-2">
                        {{ address.street }}, {{ address.city }}, {{ address.province }}
                        {{ address.postal_code }}
                      </p>
                    </div>
                    <div
                      class="ml-4"
                      :class="
                        selectedAddressId === address.id
                          ? 'w-6 h-6 bg-brand-600 rounded-full'
                          : 'w-6 h-6 border-2 border-gray-300 rounded-full'
                      "
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Payment Method -->
          <div class="bg-white rounded-lg shadow-sm p-6">
            <h2 class="text-xl font-semibold text-gray-900 mb-4">Payment Method</h2>
            <div class="space-y-3">
              <div class="p-4 border-2 border-brand-600 rounded-lg bg-brand-50">
                <label class="flex items-center cursor-pointer">
                  <input type="radio" checked disabled class="w-4 h-4" />
                  <span class="ml-3">
                    <span class="font-semibold text-gray-900">Card / Bank Transfer / E-Wallet</span>
                    <p class="text-sm text-gray-600">Pay safely with Midtrans</p>
                  </span>
                </label>
              </div>
            </div>
          </div>
        </div>

        <!-- Order Summary Sidebar -->
        <div class="lg:col-span-1">
          <div class="bg-white rounded-lg shadow-sm p-6 sticky top-8">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Order Total</h3>

            <div class="space-y-3 mb-4 pb-4 border-b">
              <div class="flex justify-between text-gray-600">
                <span>Subtotal</span>
                <span>Rp{{ subtotal.toLocaleString('id-ID') }}</span>
              </div>
              <div class="flex justify-between text-gray-600">
                <span>Shipping</span>
                <span class="text-green-600 font-medium">Free</span>
              </div>
              <div class="flex justify-between text-gray-600">
                <span>Tax</span>
                <span>Rp{{ tax.toLocaleString('id-ID') }}</span>
              </div>
            </div>

            <div class="flex justify-between text-lg font-bold text-gray-900 mb-6">
              <span>Total</span>
              <span>Rp{{ total.toLocaleString('id-ID') }}</span>
            </div>

            <button
              @click="handleCheckout"
              :disabled="isLoading || !selectedAddressId || cart.state.items.length === 0"
              class="w-full py-3 bg-brand-600 text-white rounded-lg hover:bg-brand-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition font-semibold"
            >
              <span v-if="!isLoading">Pay Now - Rp{{ total.toLocaleString('id-ID') }}</span>
              <span v-else class="flex items-center justify-center">
                <LoadingSpinner :is-loading="true" message="Processing..." />
              </span>
            </button>

            <button
              @click="goToCart"
              class="w-full mt-3 py-3 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition font-semibold"
            >
              Back to Cart
            </button>

            <p class="text-xs text-gray-600 mt-4 text-center">
              Your payment is secure and encrypted by Midtrans
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCart } from '@/stores/cart'
import { useToast } from '@/composables/useToast'
import addressService from '@/services/addressService'
import paymentService from '@/services/paymentService'
import orderService from '@/services/orderService'
import PlusIcon from '@/icons/PlusIcon.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import { getImageUrl } from '@/utils/image'

const router = useRouter()
const cart = useCart()
const { showSuccess, showError } = useToast()

const addresses = ref<any[]>([])
const selectedAddressId = ref<number | null>(null)
const isLoading = ref(false)

const subtotal = computed(() => {
  return cart.state.items.reduce((sum: number, item: any) => sum + item.product.price * item.quantity, 0)
})

const tax = computed(() => Math.round(subtotal.value * 0.1))

const total = computed(() => subtotal.value + tax.value)

onMounted(async () => {
  // Load addresses
  try {
    const data = await addressService.getMyAddresses()
    addresses.value = data
    // Auto-select default address
    const defaultAddress = addresses.value.find((a) => a.is_default)
    if (defaultAddress) {
      selectedAddressId.value = defaultAddress.id
    }
  } catch (error) {
    showError('Failed to load addresses')
  }

  // Load Snap script
  try {
    await paymentService.loadSnapScript()
  } catch (error) {
    showError('Failed to load payment gateway')
  }
})

const handleCheckout = async () => {
  if (!selectedAddressId.value) {
    showError('Please select a delivery address')
    return
  }

  isLoading.value = true

  try {
    // Create order first
    const orderReq = {
      items: cart.state.items.map((item: any) => ({
        product_id: item.product.id,
        quantity: item.quantity,
        price: item.product.price,
      })),
      address_id: selectedAddressId.value,
    }

    const orderResult = await orderService.createOrder(orderReq)
    const orderId = orderResult.id

    // Create payment
    const paymentResult = await paymentService.createPayment(
      orderId,
      selectedAddressId.value,
      total.value
    )

    showSuccess('Order created successfully!')

    // Show Midtrans Snap payment modal
    await paymentService.showPayment(paymentResult.snap_token, paymentResult.payment_id, {
      onSuccess: async () => {
        showSuccess('Payment successful! Your order is being processed.')
        cart.clearCart()
        await router.push(`/orders`)
      },
      onPending: () => {
        showSuccess('Payment pending. Please check your email for updates.')
      },
      onError: () => {
        showError('Payment failed. Please try again.')
      },
    })
  } catch (error: any) {
    showError(error.response?.data?.error || 'Checkout failed')
  } finally {
    isLoading.value = false
  }
}

const goToCart = () => {
  router.push('/cart')
}

const goToAddresses = () => {
  router.push('/addresses')
}

</script>
