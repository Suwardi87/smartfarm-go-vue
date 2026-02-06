<template>
  <MarketplaceLayout>
    <!-- CONTENT -->
    <div
      v-if="product"
      class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10"
    >
      <div class="grid grid-cols-1 md:grid-cols-2 gap-10">

        <!-- IMAGE -->
        <div class="space-y-4">
          <div class="aspect-square bg-gray-100 rounded-2xl overflow-hidden relative">
            <img
              v-if="product.image_url"
              :src="getImageUrl(product.image_url)"
              :alt="product.name"
              class="w-full h-full object-cover"
            />
            <div
              v-else
              class="absolute inset-0 flex items-center justify-center text-gray-400 text-6xl"
            >
              🥦
            </div>

            <!-- BADGES -->
            <div class="absolute top-4 left-4 flex flex-col gap-2">
              <span
                v-if="product.is_pre_order"
                class="px-3 py-1 bg-blue-500 text-white text-sm font-bold rounded-lg shadow-sm"
              >
                PRE-ORDER
              </span>
              <span
                v-if="product.is_subscription"
                class="px-3 py-1 bg-purple-500 text-white text-sm font-bold rounded-lg shadow-sm"
              >
                LANGGANAN
              </span>
            </div>
          </div>
        </div>

        <!-- INFO -->
        <div>
          <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">
            {{ product.name }}
          </h1>

          <div
            v-if="product.farmer_name"
            class="flex items-center gap-2 mb-6"
          >
            <div class="w-8 h-8 rounded-full bg-brand-100 flex items-center justify-center text-brand-600">
              👨‍🌾
            </div>
            <span class="text-gray-600 dark:text-gray-300">
              Petani:
              <span class="font-semibold text-gray-900 dark:text-white">
                {{ product.farmer_name }}
              </span>
            </span>
          </div>

          <!-- PRICE -->
          <div class="mb-8">
            <span class="text-4xl font-bold text-brand-600">
              {{ formatRupiah(product.price) }}
            </span>
            <span
              v-if="product.is_subscription"
              class="text-gray-500 text-lg"
            >
              / {{ product.subscription_period }}
            </span>
            <span
              v-else
              class="text-gray-500 text-lg"
            >
              / unit
            </span>
          </div>

          <!-- DESCRIPTION -->
          <div class="mb-8 text-gray-600 dark:text-gray-300">
            <h3 class="font-bold text-lg text-gray-900 dark:text-white mb-2">
              Deskripsi
            </h3>
            <p>{{ product.description || '-' }}</p>
          </div>

          <!-- HARVEST DATE -->
          <div
            v-if="product.is_pre_order && product.harvest_date"
            class="bg-blue-50 dark:bg-blue-900/20 p-4 rounded-xl mb-8 border border-blue-100 dark:border-blue-800"
          >
            <h4 class="font-bold text-blue-800 dark:text-blue-300 mb-1">
              📅 Estimasi Panen
            </h4>
            <p class="text-blue-600 dark:text-blue-200">
              Siap dikirim:
              <span class="font-semibold">
                {{ formatDate(product.harvest_date) }}
              </span>
            </p>
          </div>

          <!-- ACTION -->
          <div class="flex gap-4 mb-4">
            <div class="flex items-center border border-gray-300 dark:border-gray-600 rounded-lg">
              <button
                @click="decreaseQty"
                class="px-4 py-3 hover:bg-gray-100 dark:hover:bg-gray-700"
              >
                -
              </button>
              <input
                type="number"
                v-model="quantity"
                min="1"
                class="w-16 text-center bg-transparent border-none focus:ring-0"
              />
              <button
                @click="increaseQty"
                class="px-4 py-3 hover:bg-gray-100 dark:hover:bg-gray-700"
              >
                +
              </button>
            </div>

            <button
              @click="addToCart"
              class="flex-1 bg-brand-600 text-white font-bold py-3 px-8 rounded-lg hover:bg-brand-700 transition"
            >
              Tambah ke Keranjang
            </button>
          </div>

          <p class="text-sm text-gray-500">
            Stok tersedia: {{ product.stock }} unit
          </p>
        </div>
      </div>
    </div>

    <!-- LOADING -->
    <div
      v-else
      class="flex justify-center items-center min-h-[50vh]"
    >
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-brand-600"></div>
    </div>
  </MarketplaceLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import MarketplaceLayout from '@/components/layout/MarketplaceLayout.vue'
import { getProduct } from '@/services/productService'
import { getImageUrl } from '@/utils/image'
import type { Product } from '@/dto/product/Product'
import { useCart } from '@/stores/cart'
import { useToast } from '@/composables/useToast'
import { normalizeProduct } from '@/utils/transformers'

/* state */
const route = useRoute()
const cart = useCart()
const { showSuccess } = useToast()
const product = ref<Product | null>(null)
const quantity = ref(1)

/* fetch by id */
onMounted(async () => {
  const id = Number(route.params.id)
  if (!id) return

  try {
    const response = await getProduct(id)
    const rawData = (response.data as any).data
    product.value = normalizeProduct(rawData)
  } catch (error) {
    console.error('Failed to fetch product', error)
  }
})

/* cart */
const addToCart = () => {
  if (product.value) {
    cart.addItem(product.value, quantity.value)
    showSuccess(`${product.value.name} ditambahkan ke keranjang`)
  }
}

/* qty */
const increaseQty = () => quantity.value++
const decreaseQty = () => {
  if (quantity.value > 1) quantity.value--
}

/* helpers */
const formatRupiah = (value: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(value)
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
}

</script>
