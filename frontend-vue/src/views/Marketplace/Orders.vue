<template>
  <MarketplaceLayout>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">Riwayat Pesanan</h1>
        <p class="text-gray-600 dark:text-gray-400">Lihat semua pesanan dan status pengiriman Anda</p>
      </div>

      <!-- Tabs -->
      <div class="flex gap-4 mb-8 border-b border-gray-200 dark:border-gray-700">
        <button
          @click="activeTab = 'orders'"
          :class="[
            'px-4 py-3 font-medium border-b-2 transition',
            activeTab === 'orders'
              ? 'border-brand-600 text-brand-600'
              : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-300'
          ]"
        >
          Pesanan Reguler ({{ orders.length }})
        </button>
        <button
          @click="activeTab = 'subscriptions'"
          :class="[
            'px-4 py-3 font-medium border-b-2 transition',
            activeTab === 'subscriptions'
              ? 'border-brand-600 text-brand-600'
              : 'border-transparent text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-300'
          ]"
        >
          Langganan ({{ subscriptions.length }})
        </button>
      </div>

      <!-- Loading -->
      <LoadingSpinner :isLoading="isLoading" message="Memuat pesanan..." />

      <!-- Orders Tab -->
      <div v-if="!isLoading && activeTab === 'orders'" class="space-y-6">
        <div v-if="orders.length === 0" class="text-center py-12">
          <div class="text-6xl mb-4">📦</div>
          <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Belum ada pesanan</h3>
          <p class="text-gray-500 dark:text-gray-400 mb-6">Mulai belanja sekarang dan cek status pesanan di sini</p>
          <router-link to="/" class="inline-block px-6 py-3 bg-brand-600 text-white rounded-lg hover:bg-brand-700 transition">
            Belanja Sekarang
          </router-link>
        </div>

        <div v-else v-for="order in orders" :key="order.id" class="bg-white dark:bg-gray-800 p-6 rounded-xl border border-gray-200 dark:border-gray-700">
          <!-- Header -->
          <div class="flex items-center justify-between mb-6 pb-6 border-b border-gray-200 dark:border-gray-700">
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400">Nomor Pesanan</p>
              <h3 class="text-lg font-bold text-gray-900 dark:text-white">#{{ order.id }}</h3>
            </div>
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400 mb-1">Status</p>
              <span :class="getStatusBadgeClass(order.status)" class="px-3 py-1 rounded-full text-sm font-semibold">
                {{ getStatusLabel(order.status) }}
              </span>
            </div>
            <div class="text-right">
              <p class="text-sm text-gray-500 dark:text-gray-400">Total</p>
              <p class="text-xl font-bold text-brand-600">{{ formatRupiah(order.total_price) }}</p>
            </div>
          </div>

          <!-- Items -->
          <div class="mb-6">
            <p class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-4">Produk:</p>
            <div v-for="item in order.items" :key="item.id" class="flex gap-4 mb-4">
              <div class="w-16 h-16 bg-gray-100 dark:bg-gray-700 rounded-lg overflow-hidden flex-shrink-0 relative">
                <img v-if="item.product?.image_url" :src="getImageUrl(item.product.image_url)" :alt="item.product_name" class="w-full h-full object-cover" />
                <div v-else class="absolute inset-0 flex items-center justify-center text-xl bg-gray-100 dark:bg-gray-700">🛒</div>
              </div>
              <div class="flex-1">
                <p class="font-semibold text-gray-900 dark:text-white">{{ item.product_name || item.product?.name || 'Produk' }}</p>
                <p class="text-sm text-gray-500">{{ item.quantity }} x {{ formatRupiah(item.price) }}</p>
              </div>
              <p class="font-bold text-gray-900 dark:text-white">{{ formatRupiah(item.price * item.quantity) }}</p>
            </div>
          </div>

          <!-- Date -->
          <div class="text-sm text-gray-500 dark:text-gray-400">
            Dipesan pada: {{ formatDate(order.created_at) }}
          </div>
        </div>
      </div>

      <!-- Subscriptions Tab -->
      <div v-if="!isLoading && activeTab === 'subscriptions'" class="space-y-6">
        <div v-if="subscriptions.length === 0" class="text-center py-12">
          <div class="text-6xl mb-4">🔄</div>
          <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Belum ada langganan</h3>
          <p class="text-gray-500 dark:text-gray-400 mb-6">Mulai langganan produk favorit Anda sekarang</p>
          <router-link to="/" class="inline-block px-6 py-3 bg-brand-600 text-white rounded-lg hover:bg-brand-700 transition">
            Lihat Produk
          </router-link>
        </div>

        <div v-else v-for="sub in subscriptions" :key="sub.id" class="bg-white dark:bg-gray-800 p-6 rounded-xl border border-gray-200 dark:border-gray-700">
          <!-- Header -->
          <div class="flex items-center justify-between mb-6 pb-6 border-b border-gray-200 dark:border-gray-700">
            <div v-if="sub.product" class="flex gap-4 flex-1">
              <div class="w-16 h-16 bg-gray-100 dark:bg-gray-700 rounded-lg overflow-hidden flex-shrink-0">
                <img :src="getImageUrl(sub.product.image_url)" :alt="sub.product.name" class="w-full h-full object-cover" />
              </div>
              <div>
                <p class="font-bold text-gray-900 dark:text-white">{{ sub.product.name }}</p>
                <p class="text-sm text-gray-500">{{ formatRupiah(sub.product.price) }} / {{ sub.frequency }}</p>
              </div>
            </div>
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400 mb-1">Status</p>
              <span :class="getStatusBadgeClass(sub.status)" class="px-3 py-1 rounded-full text-sm font-semibold">
                {{ getStatusLabel(sub.status) }}
              </span>
            </div>
          </div>

          <!-- Info -->
          <div class="text-sm text-gray-500 dark:text-gray-400">
            Dimulai: {{ formatDate(sub.created_at) }}
          </div>
        </div>
      </div>
    </div>
  </MarketplaceLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import MarketplaceLayout from '@/components/layout/MarketplaceLayout.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import { getMyOrders, getMySubscriptions } from '@/services/orderService'
import type { Order, Subscription } from '@/dto/order/Order'
import { useToast } from '@/composables/useToast'
import { getImageUrl } from '@/utils/image'

const activeTab = ref<'orders' | 'subscriptions'>('orders')
const orders = ref<Order[]>([])
const subscriptions = ref<Subscription[]>([])
const isLoading = ref(false)
const { showError } = useToast()


onMounted(async () => {
  isLoading.value = true
  try {
    const [ordersRes, subsRes] = await Promise.all([
      getMyOrders(),
      getMySubscriptions()
    ])
    orders.value = ordersRes.data.data || []
    subscriptions.value = subsRes.data.data || []
  } catch (error: any) {
    console.error('Failed to fetch orders', error)
    showError('Gagal memuat pesanan. Silakan coba lagi.')
  } finally {
    isLoading.value = false
  }
})

const formatRupiah = (value: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value)
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}



const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    pending: 'Menunggu Pembayaran',
    paid: 'Sudah Dibayar',
    shipped: 'Sedang Dikirim',
    delivered: 'Terkirim',
    active: 'Aktif',
    paused: 'Dijeda',
    cancelled: 'Dibatalkan'
  }
  return labels[status] || status
}

const getStatusBadgeClass = (status: string) => {
  const classes: Record<string, string> = {
    pending: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-200',
    paid: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-200',
    shipped: 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-200',
    delivered: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-200',
    active: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-200',
    paused: 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200',
    cancelled: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-200'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}
</script>
