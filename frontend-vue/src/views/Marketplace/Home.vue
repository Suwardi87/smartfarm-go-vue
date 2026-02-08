<template>
  <MarketplaceLayout>
    <!-- Loading Spinner -->
    <LoadingSpinner :isLoading="isLoading" message="Memuat produk..." />

    <!-- Hero Section -->
    <section v-if="!isLoading" class="mb-12 text-center py-12 bg-green-50 dark:bg-green-900/20 rounded-3xl">
      <h1 class="text-4xl md:text-5xl font-bold mb-4 text-gray-900 dark:text-white">
        Panen Segar Langsung dari Petani
      </h1>
      <p class="text-lg text-gray-600 dark:text-gray-300 max-w-2xl mx-auto mb-8">
        Dukung petani lokal dan nikmati sayuran segar dengan harga stabil melalui sistem Pre-Order dan Langganan.
      </p>
      <div class="flex gap-4 justify-center">
        <button class="px-6 py-3 bg-brand-600 text-white rounded-xl font-semibold hover:bg-brand-700 transition shadow-lg shadow-brand-500/30">
          Belanja Sekarang
        </button>
        <button class="px-6 py-3 bg-white text-gray-700 border border-gray-200 rounded-xl font-semibold hover:bg-gray-50 dark:bg-gray-800 dark:text-white dark:border-gray-700">
          Cara Kerja
        </button>
      </div>
    </section>

    <!-- Features Grid -->
    <section class="mb-16 grid grid-cols-1 md:grid-cols-3 gap-8">
      <div class="p-6 bg-white dark:bg-gray-800 rounded-2xl border border-gray-100 dark:border-gray-700 shadow-sm hover:shadow-md transition">
        <div class="text-4xl mb-4">📅</div>
        <h3 class="text-xl font-bold mb-2">Pre-Order Panen</h3>
        <p class="text-gray-500 dark:text-gray-400">Pesan sebelum panen untuk harga lebih murah dan kepastian stok.</p>
      </div>
      <div class="p-6 bg-white dark:bg-gray-800 rounded-2xl border border-gray-100 dark:border-gray-700 shadow-sm hover:shadow-md transition">
        <div class="text-4xl mb-4">📦</div>
        <h3 class="text-xl font-bold mb-2">Langganan Mingguan</h3>
        <p class="text-gray-500 dark:text-gray-400">Paket sayur segar dikirim rutin ke rumah Anda setiap minggu.</p>
      </div>
      <div class="p-6 bg-white dark:bg-gray-800 rounded-2xl border border-gray-100 dark:border-gray-700 shadow-sm hover:shadow-md transition">
        <div class="text-4xl mb-4">🚜</div>
        <h3 class="text-xl font-bold mb-2">Dukung Petani</h3>
        <p class="text-gray-500 dark:text-gray-400">Membantu petani merencanakan tanam sesuai kebutuhan Anda.</p>
      </div>
    </section>

    <!-- Fresh Products -->
    <section class="mb-12" v-if="freshProducts.length > 0">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-bold text-gray-800 dark:text-white">Panen Minggu Ini (Ready)</h2>
        <router-link to="/marketplace?type=fresh" class="text-brand-600 font-medium hover:underline">Lihat Semua →</router-link>
      </div>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <ProductCard v-for="product in freshProducts" :key="product.id" :product="product" label="FRESH" @add-to-cart="handleAddToCart" />
      </div>
    </section>

     <!-- Pre-Order Products -->
    <section class="mb-12" v-if="preOrderProducts.length > 0">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-bold text-gray-800 dark:text-white">Open Pre-Order</h2>
        <router-link to="/marketplace?type=preorder" class="text-brand-600 font-medium hover:underline">Lihat Semua →</router-link>
      </div>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <ProductCard v-for="product in preOrderProducts" :key="product.id" :product="product" label="PRE-ORDER" labelColor="bg-blue-500" @add-to-cart="handleAddToCart" />
      </div>
    </section>

     <!-- Subscription Products -->
    <section class="mb-12" v-if="subscriptionProducts.length > 0">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-bold text-gray-800 dark:text-white">Paket Langganan</h2>
        <router-link to="/marketplace?type=subscription" class="text-brand-600 font-medium hover:underline">Lihat Semua →</router-link>
      </div>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <ProductCard v-for="product in subscriptionProducts" :key="product.id" :product="product" label="SUBSCRIPTION" labelColor="bg-purple-500" @add-to-cart="handleAddToCart" />
      </div>
    </section>

  </MarketplaceLayout>
</template>

<script setup lang="ts">
import MarketplaceLayout from '@/components/layout/MarketplaceLayout.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import ProductCard from '@/components/marketplace/ProductCard.vue'
import { ref, onMounted, computed } from 'vue'
import { getProducts } from '@/services/productService'
import type { Product } from '@/dto/product/Product'
import { useCart } from '@/stores/cart'
import { useToast } from '@/composables/useToast'
import { normalizeProduct } from '@/utils/transformers'

const products = ref<Product[]>([])
const isLoading = ref(true)
const { addItem } = useCart()
const { showSuccess } = useToast()

onMounted(async () => {
  try {
    isLoading.value = true
    const response = await getProducts(1, 12) // Limit 8 for Home
    const paginatedData = response.data
    const rawData = paginatedData.data as any[] ?? []
    products.value = rawData.map(normalizeProduct)
  } catch (error) {
    console.error("Failed to fetch products", error)
  } finally {
    isLoading.value = false
  }
})

const freshProducts = computed(() =>
  products.value.filter(p => !p.is_pre_order && !p.is_subscription)
)

const preOrderProducts = computed(() =>
  products.value.filter(p => p.is_pre_order)
)

const subscriptionProducts = computed(() =>
  products.value.filter(p => p.is_subscription)
)

const filteredProducts = computed(() =>
  products.value.filter(p => p.stock > 0)
)

const handleAddToCart = (product: Product) => {
  addItem(product)
  showSuccess(`${product.name} ditambahkan ke keranjang`)
}

</script>
