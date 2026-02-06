<template>
  <MarketplaceLayout>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-8">Keranjang Belanja</h1>

      <div v-if="items.length > 0" class="flex flex-col lg:flex-row gap-8">
        <!-- Cart Items -->
        <div class="flex-1 space-y-4">
          <div v-for="(item, index) in items" :key="item.product.id" class="bg-white dark:bg-gray-800 p-4 rounded-xl border border-gray-200 dark:border-gray-700 flex gap-4 items-center">

            <!-- Image -->
            <div class="w-20 h-20 bg-gray-100 rounded-lg overflow-hidden flex-shrink-0 relative">
               <img v-if="item.product.image_url" :src="getImageUrl(item.product.image_url)" class="w-full h-full object-cover" />
               <div v-else class="absolute inset-0 flex items-center justify-center text-2xl">🥦</div>
            </div>

            <!-- Info -->
            <div class="flex-1 min-w-0">
              <h3 class="font-bold text-gray-900 dark:text-white truncate">{{ item.product.name }}</h3>
              <p class="text-sm text-gray-500">{{ formatRupiah(item.product.price) }} / unit</p>
              <div v-if="item.product.is_pre_order" class="inline-block mt-1 px-2 py-0.5 bg-blue-100 text-blue-700 text-xs rounded font-bold">PRE-ORDER</div>
            </div>

            <!-- Qty -->
            <div class="flex items-center border border-gray-300 dark:border-gray-600 rounded-lg h-9">
                <button @click="updateQuantity(item.product.id, item.quantity - 1)" class="px-2 hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-300 transition">-</button>
                <input type="number" :value="item.quantity" @input="(e: any) => updateQuantity(item.product.id, parseInt(e.target.value))" class="w-10 text-center bg-transparent border-none focus:ring-0 text-sm p-0 appearance-none" min="1" />
                <button @click="updateQuantity(item.product.id, item.quantity + 1)" class="px-2 hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-300 transition">+</button>
             </div>

             <!-- Subtotal & Remove -->
             <div class="text-right min-w-[100px]">
                <div class="font-bold text-brand-600">{{ formatRupiah(item.product.price * item.quantity) }}</div>
                <button @click="removeItem(item.product.id)" class="text-xs text-red-500 hover:text-red-700 mt-1 underline">Hapus</button>
             </div>
          </div>
        </div>

        <!-- Summary -->
        <div class="w-full lg:w-80">
          <div class="bg-white dark:bg-gray-800 p-6 rounded-xl border border-gray-100 dark:border-gray-700 shadow-sm sticky top-24">
            <h3 class="font-bold text-lg mb-4 text-gray-900 dark:text-white">Ringkasan Pesanan</h3>

            <div class="space-y-2 mb-4 text-sm">
              <div class="flex justify-between text-gray-600 dark:text-gray-400">
                <span>Total Item</span>
                <span>{{ totalItems }} items</span>
              </div>
              <div class="flex justify-between font-bold text-gray-900 dark:text-white text-lg border-t border-gray-200 dark:border-gray-700 pt-2 mt-2">
                <span>Total Harga</span>
                <span>{{ formatRupiah(totalPrice) }}</span>
              </div>
            </div>

            <button @click="goToCheckout" :disabled="isProcessing" class="w-full bg-brand-600 text-white font-bold py-3 rounded-lg hover:bg-brand-700 transition shadow-lg shadow-brand-500/30 disabled:opacity-50 disabled:cursor-not-allowed">
              {{ isProcessing ? 'Processing...' : 'Proceed to Checkout' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-20">
        <div class="text-6xl mb-4">🛒</div>
        <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Keranjang anda kosong</h2>
        <p class="text-gray-500 mb-8">Yuk mulai belanja sayuran segar untuk kebutuhan Anda!</p>
        <router-link to="/" class="inline-block px-6 py-3 bg-brand-600 text-white rounded-xl font-semibold hover:bg-brand-700 transition">
          Mulai Belanja
        </router-link>
      </div>

    </div>
  </MarketplaceLayout>
</template>

<script setup lang="ts">
import MarketplaceLayout from '@/components/layout/MarketplaceLayout.vue'
import { useCart } from '@/stores/cart'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { createOrder, createSubscription } from '@/services/orderService'
import { useToast } from '@/composables/useToast'
import { getImageUrl } from '@/utils/image'

const cart = useCart()
const router = useRouter()
const { showSuccess, showError } = useToast()
const items = computed(() => cart.state.items)
const { updateQuantity, removeItem, totalItems, totalPrice, clearCart } = cart
const isProcessing = ref(false)

const formatRupiah = (value: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value)
}



const handleCheckout = async () => {
    router.push('/checkout')
}

const goToCheckout = () => {
    router.push('/checkout')
}
</script>
