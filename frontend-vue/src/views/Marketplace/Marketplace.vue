<template>
  <MarketplaceLayout>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
        <div>
          <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Katalog Produk</h1>
          <p class="text-gray-600 dark:text-gray-400">Temukan panen terbaik dari grup tani lokal kami.</p>
        </div>

        <!-- Category Filters -->
        <div class="flex flex-wrap gap-2">
          <button
            v-for="cat in categories"
            :key="cat.id"
            @click="activeCategory = cat.id"
            :class="[
              'px-4 py-2 rounded-full text-sm font-medium transition',
              activeCategory === cat.id
                ? 'bg-brand-600 text-white'
                : 'bg-white text-gray-700 border border-gray-200 hover:bg-gray-50 dark:bg-gray-800 dark:text-gray-300 dark:border-gray-700'
            ]"
          >
            {{ cat.name }}
          </button>
        </div>
      </div>

      <!-- Loading State -->
      <LoadingSpinner :isLoading="isLoading" message="Memuat produk katalog..." />

      <!-- Product Grid -->
      <div v-if="!isLoading">
        <div v-if="filteredProducts.length > 0" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          <ProductCard
            v-for="product in filteredProducts"
            :key="product.id"
            :product="product"
            :label="getProductLabel(product)"
            :labelColor="getProductLabelColor(product)"
            @add-to-cart="handleAddToCart"
          />
        </div>

        <!-- Pagination / Load More -->
        <div v-if="currentPage < totalPages" class="mt-12 text-center">
          <button
            @click="loadMore"
            :disabled="isLoadMore"
            class="inline-flex items-center px-8 py-3 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl font-semibold text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 transition shadow-sm disabled:opacity-50"
          >
            <span v-if="isLoadMore" class="mr-2">
              <svg class="animate-spin h-5 w-5 text-gray-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
            </span>
            {{ isLoadMore ? 'Memuat...' : 'Muat Lebih Banyak' }}
          </button>
          <p class="mt-4 text-sm text-gray-500">Menampilkan {{ products.length }} dari {{ totalProducts }} produk</p>
        </div>
        
        <!-- Empty State -->
        <div v-else class="text-center py-20 bg-gray-50 dark:bg-gray-900/50 rounded-3xl">
          <div class="text-6xl mb-4">🥦</div>
          <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Tidak ada produk ditemukan</h3>
          <p class="text-gray-500">Coba pilih kategori lain atau cek lagi nanti.</p>
        </div>
      </div>
    </div>
  </MarketplaceLayout>
</template>

<script setup lang="ts">
import MarketplaceLayout from '@/components/layout/MarketplaceLayout.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import ProductCard from '@/components/marketplace/ProductCard.vue'
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getProducts } from '@/services/productService'
import type { Product } from '@/dto/product/Product'
import { useCart } from '@/stores/cart'
import { useToast } from '@/composables/useToast'
import { normalizeProduct } from '@/utils/transformers'

const route = useRoute()
const products = ref<Product[]>([])
const totalProducts = ref(0)
const isLoading = ref(true)
const isLoadMore = ref(false)
const currentPage = ref(1)
const totalPages = ref(1)
const limit = 12

const { addItem } = useCart()
const { showSuccess } = useToast()

const activeCategory = ref('all')

const categories = [
  { id: 'all', name: 'Semua Produk' },
  { id: 'fresh', name: 'Ready (Panen)' },
  { id: 'preorder', name: 'Pre-Order' },
  { id: 'subscription', name: 'Langganan' }
]

onMounted(async () => {
  // Sync category and search from query params
  if (route.query.type) {
    activeCategory.value = route.query.type as string
  }

  fetchProducts()
})

const fetchProducts = async (page = 1, append = false) => {
  try {
    if (append) {
      isLoadMore.value = true
    } else {
      isLoading.value = true
      products.value = []
    }

    const q = route.query.q as string || ''
    const response = await getProducts(page, limit, q)
    const paginatedData = response.data
    
    const rawData = paginatedData.data as any[] ?? []
    const normalizedData = rawData.map(normalizeProduct)
    
    if (append) {
      products.value = [...products.value, ...normalizedData]
    } else {
      products.value = normalizedData
    }
    
    currentPage.value = paginatedData.page
    totalPages.value = paginatedData.total_pages
    totalProducts.value = paginatedData.total
  } catch (error) {
    console.error("Failed to fetch products", error)
  } finally {
    isLoading.value = false
    isLoadMore.value = false
  }
}

const loadMore = () => {
  if (currentPage.value < totalPages.value) {
    fetchProducts(currentPage.value + 1, true)
  }
}

// Watch for search query changes
watch(() => route.query.q, () => {
  fetchProducts()
})

// Update activeCategory when query changes
watch(() => route.query.type, (newType) => {
  if (newType) activeCategory.value = newType as string
})

const filteredProducts = computed(() => {
  if (activeCategory.value === 'all') return products.value
  if (activeCategory.value === 'fresh') return products.value.filter(p => !p.is_pre_order && !p.is_subscription)
  if (activeCategory.value === 'preorder') return products.value.filter(p => p.is_pre_order)
  if (activeCategory.value === 'subscription') return products.value.filter(p => p.is_subscription)
  return products.value
})

const getProductLabel = (product: Product) => {
  if (product.is_subscription) return 'SUBSCRIPTION'
  if (product.is_pre_order) return 'PRE-ORDER'
  return 'FRESH'
}

const getProductLabelColor = (product: Product) => {
  if (product.is_subscription) return 'bg-purple-500'
  if (product.is_pre_order) return 'bg-blue-500'
  return 'bg-brand-500'
}

const handleAddToCart = (product: Product) => {
  addItem(product)
  showSuccess(`${product.name} ditambahkan ke keranjang`)
}
</script>
